package user

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jmoiron/sqlx"
)

type TOUser struct {
	api.APIInfoImpl `json:"-"`
	tc.User
}

func (user TOUser) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

func (user TOUser) GetKeys() (map[string]interface{}, bool) {
	if user.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *user.ID}, true
}

func (user TOUser) GetAuditName() string {
	if user.Username != nil {
		return *user.Username
	}
	if user.ID != nil {
		return strconv.Itoa(*user.ID)
	}
	return "unknown"
}

func (user TOUser) GetType() string {
	return "user"
}

func (user *TOUser) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) // non-panicking type assertion
	user.ID = &i
}

func (user *TOUser) SetLastUpdated(t tc.TimeNoMod) {
	user.LastUpdated = &t
}

func (user *TOUser) NewReadObj() interface{} {
	return &tc.User{}
}

func (user *TOUser) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":       dbhelpers.WhereColumnInfo{Column: "u.id", Checker: api.IsInt},
		"role":     dbhelpers.WhereColumnInfo{Column: "r.name"},
		"tenant":   dbhelpers.WhereColumnInfo{Column: "t.name"},
		"username": dbhelpers.WhereColumnInfo{Column: "u.username"},
	}
}

func (user *TOUser) Validate() (error, error) {

	validateErrs := validation.Errors{
		"email":    validation.Validate(user.Email, validation.Required, is.Email),
		"fullName": validation.Validate(user.FullName, validation.Required),
		"role":     validation.Validate(user.Role, validation.Required),
		"username": validation.Validate(user.Username, validation.Required),
		"tenantID": validation.Validate(user.TenantID, validation.Required),
	}

	// Password is not required for update
	if user.LocalPassword != nil {
		_, err := auth.IsGoodLoginPair(*user.Username, *user.LocalPassword)
		if err != nil {
			return err, nil
		}
	}

	return util.JoinErrs(tovalidate.ToErrors(validateErrs)), nil
}

func (user *TOUser) postValidate() error {
	validateErrs := validation.Errors{
		"localPasswd": validation.Validate(user.LocalPassword, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}

func postValidateV40(user tc.UserV4) error {
	validateErrs := validation.Errors{
		"localPasswd": validation.Validate(user.LocalPassword, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}

// Note: Not using GenericCreate because Scan also needs to scan tenant and rolename
func (user *TOUser) Create() (error, error, int) {

	// PUT and POST validation differs slightly
	err := user.postValidate()
	if err != nil {
		return err, nil, http.StatusBadRequest
	}

	// make sure the user cannot create someone with a higher priv_level than themselves
	if usrErr, sysErr, code := user.privCheck(); code != http.StatusOK {
		return usrErr, sysErr, code
	}
	var caps []string
	if user.Role != nil {
		caps, err = dbhelpers.GetCapabilitiesFromRoleID(user.ReqInfo.Tx.Tx, *user.Role)
	} else if user.RoleName != nil {
		caps, err = dbhelpers.GetCapabilitiesFromRoleName(user.ReqInfo.Tx.Tx, *user.RoleName)
	}
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	missing := user.ReqInfo.User.MissingPermissions(caps...)
	if len(missing) != 0 {
		return fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil, http.StatusForbidden
	}
	// Convert password to SCRYPT
	*user.LocalPassword, err = auth.DerivePassword(*user.LocalPassword)
	if err != nil {
		return err, nil, http.StatusBadRequest
	}

	resultRows, err := user.ReqInfo.Tx.NamedQuery(user.InsertQuery(), user)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	var id int
	var lastUpdated tc.TimeNoMod
	var tenant string
	var rolename string

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err = resultRows.Scan(&id, &lastUpdated, &tenant, &rolename); err != nil {
			return nil, fmt.Errorf("could not scan after insert: %s\n)", err), http.StatusInternalServerError
		}
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no user was inserted, nothing was returned"), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf("too many rows affected from user insert"), http.StatusInternalServerError
	}

	user.ID = &id
	user.LastUpdated = &lastUpdated
	user.Tenant = &tenant
	user.RoleName = &rolename
	user.LocalPassword = nil

	return nil, nil, http.StatusOK
}

// This is not using GenericRead because of this tenancy check. Maybe we can add tenancy functionality to the generic case?
func (this *TOUser) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	var query string

	inf := this.APIInfo()
	api.DefaultSort(inf, "username")
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, this.ParamColumns())
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	tenantIDs, err := tenant.GetUserTenantIDListTx(inf.Tx.Tx, inf.User.TenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting tenant list for user: %w", err), http.StatusInternalServerError, nil
	}
	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "u.tenant_id", tenantIDs)

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(this.APIInfo().Tx, h, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return []interface{}{}, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	groupBy := "\n" + `GROUP BY u.id, r.name, t.name`
	orderBy = groupBy + orderBy

	version := inf.Version
	if version == nil {
		return nil, nil, fmt.Errorf("TOUsers.Read called with invalid API version"), http.StatusInternalServerError, nil
	}
	if version.Major >= 4 {
		query = this.SelectQuery40() + where + orderBy + pagination
	} else {
		query = this.SelectQuery() + where + orderBy + pagination
	}

	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, fmt.Errorf("querying users : %w", err), http.StatusInternalServerError, nil
	}
	defer rows.Close()

	type UserGet struct {
		RoleName *string `json:"rolename" db:"rolename"`
		tc.User
	}
	type UserGet40 struct {
		UserGet
		ChangeLogCount    *int       `json:"changeLogCount" db:"change_log_count"`
		LastAuthenticated *time.Time `json:"lastAuthenticated" db:"last_authenticated"`
	}

	user := &UserGet{}
	user40 := &UserGet40{}
	users := []interface{}{}
	for rows.Next() {
		if version.Major >= 4 {
			if err = rows.StructScan(user40); err != nil {
				return nil, nil, fmt.Errorf("parsing user rows: %w", err), http.StatusInternalServerError, nil
			}
			users = append(users, *user40)
		} else {
			if err = rows.StructScan(user); err != nil {
				return nil, nil, fmt.Errorf("parsing user rows: %w", err), http.StatusInternalServerError, nil
			}
			users = append(users, *user)
		}
	}

	return users, nil, nil, http.StatusOK, &maxTime
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(u.last_updated) as t FROM tm_user u
		LEFT JOIN tenant t ON u.tenant_id = t.id
		LEFT JOIN role r ON u.role = r.id ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='tm_user') as res`
}

func (user *TOUser) privCheck() (error, error, int) {
	var requestedPrivLevel int
	var err error
	if user.Role == nil {
		requestedPrivLevel, _, err = dbhelpers.GetPrivLevelFromRole(user.ReqInfo.Tx.Tx, *user.RoleName)
	} else {
		requestedPrivLevel, _, err = dbhelpers.GetPrivLevelFromRoleID(user.ReqInfo.Tx.Tx, *user.Role)
	}
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if user.ReqInfo.User.PrivLevel < requestedPrivLevel {
		return fmt.Errorf("user cannot update a user with a role more privileged than themselves"), nil, http.StatusForbidden
	}

	return nil, nil, http.StatusOK
}

func (user *TOUser) Update(h http.Header) (error, error, int) {

	// make sure current user cannot update their own role to a new value
	if user.ReqInfo.User.ID == *user.ID && user.ReqInfo.User.Role != *user.Role {
		return fmt.Errorf("users cannot update their own role"), nil, http.StatusBadRequest
	}

	// make sure the user cannot update someone with a higher priv_level than themselves
	if usrErr, sysErr, code := user.privCheck(); code != http.StatusOK {
		return usrErr, sysErr, code
	}

	var caps []string
	var err error
	if user.Role != nil {
		caps, err = dbhelpers.GetCapabilitiesFromRoleID(user.ReqInfo.Tx.Tx, *user.Role)
	} else if user.RoleName != nil {
		caps, err = dbhelpers.GetCapabilitiesFromRoleName(user.ReqInfo.Tx.Tx, *user.RoleName)
	}
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	missing := user.ReqInfo.User.MissingPermissions(caps...)
	if len(missing) != 0 {
		return fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil, http.StatusForbidden
	}

	if user.LocalPassword != nil {
		var err error
		*user.LocalPassword, err = auth.DerivePassword(*user.LocalPassword)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
	}
	userErr, sysErr, errCode := api.CheckIfUnModified(h, user.ReqInfo.Tx, *user.ID, "tm_user")
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	resultRows, err := user.ReqInfo.Tx.NamedQuery(user.UpdateQuery(), user)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	var lastUpdated tc.TimeNoMod
	var tenant string
	var rolename string

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated, &tenant, &rolename); err != nil {
			return nil, fmt.Errorf("could not scan lastUpdated from insert: %s\n", err), http.StatusInternalServerError
		}
	}

	user.LastUpdated = &lastUpdated
	user.Tenant = &tenant
	user.RoleName = &rolename
	user.LocalPassword = nil

	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return fmt.Errorf("no user found with this id"), nil, http.StatusNotFound
		}
		return nil, fmt.Errorf("this update affected too many rows: %d", rowsAffected), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

func (u *TOUser) IsTenantAuthorized(user *auth.CurrentUser) (bool, error) {

	// Delete: only id is given
	// Create: only tenant id
	// Update: id and tenant id
	//	id is associated with old tenant id
	//	we need to also check new tenant id

	tx := u.ReqInfo.Tx.Tx

	if u.ID != nil { // old tenant id (only on update or delete)

		var tenantID int
		if err := tx.QueryRow(`SELECT tenant_id from tm_user WHERE id = $1`, *u.ID).Scan(&tenantID); err != nil {
			if err != sql.ErrNoRows {
				return false, err
			}

			// At this point, tenancy isn't technically 'true', but I can't return a resource not found error here.
			// Letting it continue will let it run into a 404 when it tries to update.
			return true, nil
		}

		//log.Debugf("%d with tenancy %d trying to access %d with tenancy %d", user.ID, user.TenantID, *u.ID, tenantID)
		authorized, err := tenant.IsResourceAuthorizedToUserTx(tenantID, user, tx)
		if err != nil {
			return false, err
		}
		if !authorized {
			return false, nil

		}
	}

	if u.TenantID != nil { // new tenant id (only on create or udpate)

		//log.Debugf("%d with tenancy %d trying to access %d", user.ID, user.TenantID, *u.TenantID)
		authorized, err := tenant.IsResourceAuthorizedToUserTx(*u.TenantID, user, tx)
		if err != nil {
			return false, err
		}
		if !authorized {
			return false, nil
		}
	}

	return true, nil
}

func (user *TOUser) SelectQuery() string {
	return `
	SELECT
	u.id,
	u.username as username,
	u.public_ssh_key,
	u.role,
	r.name as rolename,
	u.company,
	u.email,
	u.full_name,
	u.new_user,
	u.address_line1,
	u.address_line2,
	u.city,
	u.state_or_province,
	u.phone_number,
	u.postal_code,
	u.country,
	u.registration_sent,
	u.tenant_id,
	t.name as tenant,
	u.last_updated
	FROM tm_user u
	LEFT JOIN tenant t ON u.tenant_id = t.id
	LEFT JOIN role r ON u.role = r.id`
}

func (user *TOUser) SelectQuery40() string {
	return `
	SELECT
	u.id,
	u.username as username,
	u.public_ssh_key,
	u.role,
	r.name as rolename,
	u.company,
	u.email,
	u.full_name,
	u.new_user,
	u.address_line1,
	u.address_line2,
	u.city,
	u.state_or_province,
	u.phone_number,
	u.postal_code,
	u.country,
	u.registration_sent,
	u.tenant_id,
	t.name as tenant,
	u.last_updated,
	u.last_authenticated,
	(SELECT count(l.tm_user) FROM log as l WHERE l.tm_user = u.id) as change_log_count
	FROM tm_user u
	LEFT JOIN tenant t ON u.tenant_id = t.id
	LEFT JOIN role r ON u.role = r.id`
}

func (user *TOUser) UpdateQuery() string {
	return `
	UPDATE tm_user u SET
	username=:username,
	public_ssh_key=:public_ssh_key,
	role=:role,
	company=:company,
	email=:email,
	full_name=:full_name,
	new_user=COALESCE(:new_user, FALSE),
	address_line1=:address_line1,
	address_line2=:address_line2,
	city=:city,
	state_or_province=:state_or_province,
	phone_number=:phone_number,
	postal_code=:postal_code,
	country=:country,
	tenant_id=:tenant_id,
	local_passwd=COALESCE(:local_passwd, local_passwd)
	WHERE id=:id
	RETURNING last_updated,
	 (SELECT t.name FROM tenant t WHERE id = u.tenant_id),
	 (SELECT r.name FROM role r WHERE id = u.role)`
}

func UpdateQueryV40() string {
	return `
	UPDATE tm_user u SET
	username=:username,
	public_ssh_key=:public_ssh_key,
	role=(SELECT id FROM role WHERE role.name = :role),
	company=:company,
	email=:email,
	full_name=:full_name,
	new_user=COALESCE(:new_user, FALSE),
	address_line1=:address_line1,
	address_line2=:address_line2,
	city=:city,
	state_or_province=:state_or_province,
	phone_number=:phone_number,
	postal_code=:postal_code,
	country=:country,
	tenant_id=:tenant_id,
	local_passwd=COALESCE(:local_passwd, local_passwd),
	ucdn=:ucdn
	WHERE id=:id
	RETURNING last_updated,
	 (SELECT t.name FROM tenant t WHERE id = u.tenant_id),
	 (SELECT r.name FROM role r WHERE id = u.role)`
}

func InsertQueryV40() string {
	return `
	INSERT INTO tm_user (
	username,
	public_ssh_key,
	role,
	company,
	email,
	full_name,
	new_user,
	address_line1,
	address_line2,
	city,
	state_or_province,
	phone_number,
	postal_code,
	country,
	tenant_id,
	local_passwd,
	ucdn
	) VALUES (
	:username,
	:public_ssh_key,
	(SELECT id FROM role WHERE name = :role),
	:company,
	:email,
	:full_name,
	COALESCE(:new_user, FALSE),
	:address_line1,
	:address_line2,
	:city,
	:state_or_province,
	:phone_number,
	:postal_code,
	:country,
	:tenant_id,
	:local_passwd,
	:ucdn
	) RETURNING id, last_updated,
	(SELECT t.name FROM tenant t WHERE id = tm_user.tenant_id),
	(SELECT r.name FROM role r WHERE id = tm_user.role)`
}

func (user *TOUser) DeleteQuery() string {
	return `DELETE FROM tm_user WHERE id = :id`
}

const readBaseQuery = `
SELECT
	u.id,
	u.username AS username,
	u.public_ssh_key,
	u.company,
	u.email,
	u.full_name,
	u.new_user,
	u.address_line1,
	u.address_line2,
	u.city,
	u.state_or_province,
	u.phone_number,
	u.postal_code,
	u.country,
	u.registration_sent,
	u.tenant_id,
	t.name AS tenant,
	u.last_updated,
	u.ucdn,`

const readQuery = readBaseQuery + `
u.last_authenticated,
(SELECT count(l.tm_user) FROM log as l WHERE l.tm_user = u.id) as change_log_count,
r.name as role
FROM tm_user u
LEFT JOIN tenant t ON u.tenant_id = t.id
LEFT JOIN role r ON u.role = r.id
LEFT JOIN role_capability rc ON rc.role_id = r.id
`

const legacyReadQuery = readBaseQuery + `
	r.name AS rolename,
	u.role
FROM tm_user u
LEFT JOIN tenant t ON u.tenant_id = t.id
LEFT JOIN role r ON u.role = r.id
`

// this is necessary because tc.User doesn't read its RoleName field in sql
// driver scans.
type userGet struct {
	RoleName *string `json:"rolename" db:"rolename"`
	tc.User
}

type userGet40 struct {
	userGet
	ChangeLogCount    *int       `json:"changeLogCount" db:"change_log_count"`
	LastAuthenticated *time.Time `json:"lastAuthenticated" db:"last_authenticated"`
}

func read(rows *sqlx.Rows) ([]tc.UserV4, error) {
	if rows == nil {
		return nil, errors.New("cannot read from nil rows")
	}

	users := []tc.UserV4{}
	for rows.Next() {
		var user tc.UserV4
		if err := rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("scanning UserV4 row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func getMaxLastUpdated(where string, queryValues map[string]interface{}, tx *sqlx.Tx) (time.Time, error) {
	query := selectMaxLastUpdatedQuery(where)
	var t time.Time
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return t, fmt.Errorf("query for max user last updated time: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&t); err != nil {
			return t, fmt.Errorf("scanning user max last updated time: %w", err)
		}
	}
	return t, nil
}

// Get is the handler for GET requests made to /users.
func Get(w http.ResponseWriter, r *http.Request) {
	var query string
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	api.DefaultSort(inf, "username")
	params := map[string]dbhelpers.WhereColumnInfo{
		"id":       {Column: "u.id", Checker: api.IsInt},
		"role":     {Column: "r.name"},
		"tenant":   {Column: "t.name"},
		"username": {Column: "u.username"},
	}
	params["company"] = dbhelpers.WhereColumnInfo{Column: "u.company"}
	params["email"] = dbhelpers.WhereColumnInfo{Column: "u.email"}
	params["fullName"] = dbhelpers.WhereColumnInfo{Column: "u.full_name"}
	params["newUser"] = dbhelpers.WhereColumnInfo{Column: "u.new_user"}
	params["city"] = dbhelpers.WhereColumnInfo{Column: "u.city"}
	params["stateOrProvince"] = dbhelpers.WhereColumnInfo{Column: "u.state_or_province"}
	params["country"] = dbhelpers.WhereColumnInfo{Column: "u.country"}
	params["postalCode"] = dbhelpers.WhereColumnInfo{Column: "u.postal_code"}
	params["capability"] = dbhelpers.WhereColumnInfo{Column: "rc.cap_name"}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, params)
	if len(errs) != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}

	tenantIDs, err := tenant.GetUserTenantIDListTx(inf.Tx.Tx, inf.User.TenantID)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("getting tenant list for user: %w", err))
		return
	}
	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "u.tenant_id", tenantIDs)

	if inf.Config.UseIMS {
		runSecond, maxTime := ims.TryIfModifiedSinceQuery(inf.Tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			w.Header().Add(rfc.LastModified, maxTime.Format(rfc.LastModifiedFormat))
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	groupBy := "\n" + `GROUP BY u.id, r.name, t.name`
	orderBy = groupBy + orderBy

	query = readQuery + where + orderBy + pagination

	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("querying Users: %w", err))
		return
	}
	defer log.Close(rows, "reading in Users from the database")

	var response interface{}
	response, err = read(rows)

	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}

	if inf.UseIMS() {
		maxTime, err := getMaxLastUpdated(where, queryValues, inf.Tx)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
		w.Header().Add(rfc.LastModified, maxTime.Format(rfc.LastModifiedFormat))
	}
	api.WriteResp(w, r, response)
}

func validate(user TOUser) error {
	validateErrs := validation.Errors{
		"email":    validation.Validate(user.Email, validation.Required, is.Email),
		"fullName": validation.Validate(user.FullName, validation.Required),
		"role":     validation.Validate(user.Role, validation.Required),
		"username": validation.Validate(user.Username, validation.Required),
		"tenantID": validation.Validate(user.TenantID, validation.Required),
	}

	// Password is not required for update
	if user.LocalPassword != nil {
		_, err := auth.IsGoodLoginPair(*user.Username, *user.LocalPassword)
		if err != nil {
			return err
		}
	}

	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}

func validateUserV4(user tc.UserV4) error {
	validateErrs := validation.Errors{
		"email":    validation.Validate(user.Email, validation.Required, is.Email),
		"fullName": validation.Validate(user.FullName, validation.Required),
		"role":     validation.Validate(user.Role, validation.Required),
		"username": validation.Validate(user.Username, validation.Required),
		"tenantID": validation.Validate(user.TenantID, validation.Required),
	}

	// Password is not required for update
	if user.LocalPassword != nil {
		_, err := auth.IsGoodLoginPair(user.Username, *user.LocalPassword)
		if err != nil {
			return err
		}
	}

	return util.JoinErrs(tovalidate.ToErrors(validateErrs))
}

func Create(w http.ResponseWriter, r *http.Request) {
	var userV4 tc.UserV4
	var err error
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	if err := json.NewDecoder(r.Body).Decode(&userV4); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	if err := validateUserV4(userV4); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	if err := postValidateV40(userV4); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	toUser := TOUser{
		APIInfoImpl: api.APIInfoImpl{ReqInfo: inf},
	}
	toUser.User = userV4.Downgrade()

	authorized, err := toUser.IsTenantAuthorized(inf.User)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("checking tenant authorized: "+err.Error()))
		return
	}
	if !authorized {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}

	// Convert password to SCRYPT
	*userV4.LocalPassword, err = auth.DerivePassword(*userV4.LocalPassword)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	var resultRows *sqlx.Rows
	_, ok, err := dbhelpers.GetRoleIDFromName(tx, userV4.Role)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error fetching ID from role name: %w", err))
		return
	} else if !ok {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("role not found"), nil)
		return
	}

	var caps []string
	caps, err = dbhelpers.GetCapabilitiesFromRoleName(tx, userV4.Role)

	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	missing := inf.User.MissingPermissions(caps...)
	if len(missing) != 0 {
		api.HandleErr(w, r, tx, http.StatusForbidden, fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil)
		return
	}

	resultRows, err = inf.Tx.NamedQuery(InsertQueryV40(), userV4)
	if err != nil {
		userErr, sysErr, statusCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	var id int
	var lastUpdated time.Time
	var tenant string
	var rolename string
	var changeLogMsg string

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err = resultRows.Scan(&id, &lastUpdated, &tenant, &rolename); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("could not scan after insert: %w)", err))
			return
		}
	}

	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("no userV4 was inserted, nothing was returned"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("too many rows affected from userV4 insert"))
		return
	}

	userV4.ID = &id
	userV4.LastUpdated = lastUpdated
	userV4.Tenant = &tenant
	userV4.Role = rolename
	userV4.LocalPassword = nil

	userResponse := tc.UserResponseV4{
		Response: userV4,
		Alerts:   tc.CreateAlerts(tc.SuccessLevel, "user was created."),
	}
	w.Header().Set(rfc.Location, fmt.Sprintf("/api/%s/users?id=%d", inf.Version, *userV4.ID))
	api.WriteAlertsObj(w, r, http.StatusCreated, userResponse.Alerts, userResponse.Response)
	changeLogMsg = fmt.Sprintf("USER: %s, ID: %d, ACTION: Created User", userV4.Username, *userV4.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
	return
}

func (user *TOUser) InsertQuery() string {
	return `
	INSERT INTO tm_user (
	username,
	public_ssh_key,
	role,
	company,
	email,
	full_name,
	new_user,
	address_line1,
	address_line2,
	city,
	state_or_province,
	phone_number,
	postal_code,
	country,
	tenant_id,
	local_passwd
	) VALUES (
	:username,
	:public_ssh_key,
	:role,
	:company,
	:email,
	:full_name,
	COALESCE(:new_user, FALSE),
	:address_line1,
	:address_line2,
	:city,
	:state_or_province,
	:phone_number,
	:postal_code,
	:country,
	:tenant_id,
	:local_passwd
	) RETURNING id, last_updated,
	(SELECT t.name FROM tenant t WHERE id = tm_user.tenant_id),
	(SELECT r.name FROM role r WHERE id = tm_user.role)`
}

// Update is the handler for PUT requests made to /users.
func Update(w http.ResponseWriter, r *http.Request) {
	var userV4 tc.UserV4
	var roleID int
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	idParam, ok := inf.Params["id"]
	if !ok {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("no ID supplied"), nil)
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("couldn't convert id into an int"), nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&userV4); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	if err := validateUserV4(userV4); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	userV4.ID = &id

	roleID, ok, err = dbhelpers.GetRoleIDFromName(inf.Tx.Tx, userV4.Role)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	} else if !ok {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no such role"), nil)
		return
	}
	// make sure current userV4 cannot update their own role to a new value
	if inf.User.ID == *userV4.ID && inf.User.Role != roleID {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, fmt.Errorf("users cannot update their own role"), nil)
		return
	}

	toUser := TOUser{
		APIInfoImpl: api.APIInfoImpl{ReqInfo: inf},
	}

	toUser.User = userV4.Downgrade()

	authorized, err := toUser.IsTenantAuthorized(inf.User)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("checking tenant authorized: "+err.Error()))
		return
	}
	if !authorized {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}
	// make sure the userV4 cannot create someone with a higher priv_level than themselves
	if userErr, sysErr, code := toUser.privCheck(); code != http.StatusOK {
		api.HandleErr(w, r, tx, code, userErr, sysErr)
		return
	}

	if userV4.LocalPassword != nil {
		// Convert password to SCRYPT
		*userV4.LocalPassword, err = auth.DerivePassword(*userV4.LocalPassword)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}
	}

	var caps []string
	caps, err = dbhelpers.GetCapabilitiesFromRoleName(tx, userV4.Role)

	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}
	missing := inf.User.MissingPermissions(caps...)
	if len(missing) != 0 {
		api.HandleErr(w, r, tx, http.StatusForbidden, fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil)
		return
	}

	userErr, sysErr, errCode = api.CheckIfUnModified(r.Header, inf.Tx, id, "tm_user")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	var resultRows *sqlx.Rows
	resultRows, err = inf.Tx.NamedQuery(UpdateQueryV40(), userV4)

	if err != nil {
		api.ParseDBError(err)
		return
	}
	defer resultRows.Close()

	var lastUpdated time.Time
	var tenant string
	var rolename string
	var changeLogMsg string

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated, &tenant, &rolename); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("could not scan lastUpdated from insert: %s\n", err))
			return
		}
	}

	if rowsAffected != 1 {
		if rowsAffected < 1 {
			api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("no user found with this id"), nil)
			return
		}
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("this update affected too many rows: %d", rowsAffected))
		return
	}

	userV4.LastUpdated = lastUpdated
	userV4.Tenant = &tenant
	userV4.Role = rolename
	userV4.LocalPassword = nil

	userResponse := tc.UserResponseV4{
		Response: userV4,
		Alerts:   tc.CreateAlerts(tc.SuccessLevel, "user was updated."),
	}
	api.WriteAlertsObj(w, r, http.StatusOK, userResponse.Alerts, userResponse.Response)
	changeLogMsg = fmt.Sprintf("USER: %s, ID: %d, ACTION: Updated User", userV4.Username, *userV4.ID)

	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
