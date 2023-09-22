package role

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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type roleError string

func (e roleError) Error() string {
	return string(e)
}

const cannotModifyAdminError roleError = "the '" + tc.AdminRoleName + "' Role cannot be deleted or modified"

const isAdminQuery = `SELECT name='` + tc.AdminRoleName + `' FROM public.role WHERE id=$1`

type TORole struct {
	api.APIInfoImpl `json:"-"`
	tc.Role
	LastUpdated    *tc.TimeNoMod   `json:"-"`
	PQCapabilities *pq.StringArray `json:"-" db:"capabilities"`
}

func updateRoleQuery() string {
	return `UPDATE
role SET
name=$1,
description=$2
WHERE name=$3 RETURNING last_updated`
}

func (v *TORole) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(v.APIInfo().Tx, *v.ID, "role")
}

func (v *TORole) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` r ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}

func (v *TORole) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TORole) InsertQuery() string           { return insertQuery() }
func (v *TORole) NewReadObj() interface{}       { return &TORole{} }
func (v *TORole) SelectQuery() string           { return selectQuery() }
func (v *TORole) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":      dbhelpers.WhereColumnInfo{Column: "name"},
		"id":        dbhelpers.WhereColumnInfo{Column: "id", Checker: api.IsInt},
		"privLevel": dbhelpers.WhereColumnInfo{Column: "priv_level", Checker: api.IsInt}}
}
func (v *TORole) UpdateQuery() string { return updateQuery() }
func (v *TORole) DeleteQuery() string { return deleteQuery() }

func (role TORole) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// Implementation of the Identifier, Validator interface functions
func (role TORole) GetKeys() (map[string]interface{}, bool) {
	if role.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *role.ID}, true
}

func (role TORole) GetAuditName() string {
	if role.Name != nil {
		return *role.Name
	}
	if role.ID != nil {
		return strconv.Itoa(*role.ID)
	}
	return "0"
}

func (role TORole) GetType() string {
	return "role"
}

func (role *TORole) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	role.ID = &i
}

// Validate fulfills the api.Validator interface
func (role TORole) Validate() (error, error) {
	errs := validation.Errors{
		"name":        validation.Validate(role.Name, validation.Required),
		"description": validation.Validate(role.Description, validation.Required),
		"privLevel":   validation.Validate(role.PrivLevel, validation.NotNil)}

	errsToReturn := tovalidate.ToErrors(errs)
	checkCaps := `SELECT cap FROM UNNEST($1::text[]) AS cap WHERE NOT cap =  ANY(ARRAY(SELECT c.name FROM capability AS c WHERE c.name = ANY($1)))`
	var badCaps []string
	if role.ReqInfo.Tx != nil {
		err := role.ReqInfo.Tx.Select(&badCaps, checkCaps, pq.Array(role.Capabilities))
		if err != nil {
			return nil, fmt.Errorf("got error from selecting bad capabilities: %w", err)
		}
		if len(badCaps) > 0 {
			errsToReturn = append(errsToReturn, fmt.Errorf("can not add non-existent capabilities: %v", badCaps))
		}
	}
	return util.JoinErrs(errsToReturn), nil
}

func (role *TORole) Create() (error, error, int) {
	if *role.PrivLevel > role.ReqInfo.User.PrivLevel {
		return errors.New("can not create a role with a higher priv level than your own"), nil, http.StatusBadRequest
	}
	if role.Capabilities != nil && *role.Capabilities != nil {
		caps := *role.Capabilities
		missing := role.ReqInfo.User.MissingPermissions(caps...)
		if len(missing) != 0 {
			return fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil, http.StatusForbidden
		}
	}
	userErr, sysErr, errCode := api.GenericCreate(role)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	//after we have role ID we can associate the capabilities:
	if role.Capabilities != nil && len(*role.Capabilities) > 0 {
		userErr, sysErr, errCode = role.createRoleCapabilityAssociations(role.ReqInfo.Tx)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
	}
	return nil, nil, http.StatusOK
}

func (role *TORole) createRoleCapabilityAssociations(tx *sqlx.Tx) (error, error, int) {
	result, err := tx.Exec(associateCapabilities(), role.ID, pq.Array(role.Capabilities))
	if err != nil {
		return nil, errors.New("creating role capabilities: " + err.Error()), http.StatusInternalServerError
	}

	if rows, err := result.RowsAffected(); err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	} else if expected := len(*role.Capabilities); int(rows) != expected {
		log.Errorf("wrong number of role_capability rows created: %d expected: %d", rows, expected)
	}
	return nil, nil, http.StatusOK
}

func (role *TORole) deleteRoleCapabilityAssociations(tx *sqlx.Tx) (error, error, int) {
	result, err := tx.Exec(deleteAssociatedCapabilities(), role.Name)
	if err != nil {
		return nil, errors.New("deleting role capabilities: " + err.Error()), http.StatusInternalServerError
	}

	if _, err = result.RowsAffected(); err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	}
	// TODO verify expected row count shouldn't be checked?
	return nil, nil, http.StatusOK
}

func (role *TORole) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	api.DefaultSort(role.APIInfo(), "name")
	vals, userErr, sysErr, errCode, maxTime := api.GenericRead(h, role, useIMS)
	if errCode == http.StatusNotModified {
		return []interface{}{}, nil, nil, errCode, maxTime
	}
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, maxTime
	}

	returnable := []interface{}{}
	for _, val := range vals {
		rl := val.(*TORole)
		caps := ([]string)(*rl.PQCapabilities)
		rl.Capabilities = &caps
		returnable = append(returnable, rl)
	}
	return returnable, nil, nil, http.StatusOK, maxTime
}

func (role *TORole) Update(h http.Header) (error, error, int) {

	if ok, err := dbhelpers.RoleExists(role.ReqInfo.Tx.Tx, *role.ID); err != nil {
		return nil, fmt.Errorf("verifying Role exists: %w", err), http.StatusInternalServerError
	} else if !ok {
		return errors.New("role not found"), nil, http.StatusNotFound
	}

	var isAdmin bool
	if err := role.ReqInfo.Tx.Get(&isAdmin, isAdminQuery, role.ID); err != nil {
		return nil, fmt.Errorf("checking if Role to be modified is '%s': %w", tc.AdminRoleName, err), http.StatusInternalServerError
	}
	if isAdmin {
		return cannotModifyAdminError, nil, http.StatusBadRequest
	}

	if *role.PrivLevel > role.ReqInfo.User.PrivLevel {
		return errors.New("can not create a role with a higher priv level than your own"), nil, http.StatusForbidden
	}
	if role.Capabilities != nil && *role.Capabilities != nil {
		caps := *role.Capabilities
		missing := role.ReqInfo.User.MissingPermissions(caps...)
		if len(missing) != 0 {
			return fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil, http.StatusForbidden
		}
	}
	userErr, sysErr, errCode := api.GenericUpdate(h, role)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	// TODO cascade delete, to automatically do this in SQL?
	if role.Capabilities != nil && *role.Capabilities != nil {
		userErr, sysErr, errCode = role.deleteRoleCapabilityAssociations(role.ReqInfo.Tx)
		if userErr != nil || sysErr != nil {
			return userErr, sysErr, errCode
		}
		return role.createRoleCapabilityAssociations(role.ReqInfo.Tx)
	}
	return nil, nil, http.StatusOK
}

func (role *TORole) Delete() (error, error, int) {

	if ok, err := dbhelpers.RoleExists(role.ReqInfo.Tx.Tx, *role.ID); err != nil {
		return nil, fmt.Errorf("verifying Role exists: %w", err), http.StatusInternalServerError
	} else if !ok {
		return errors.New("role not found"), nil, http.StatusNotFound
	}

	assignedUsers := 0
	if err := role.ReqInfo.Tx.Get(&assignedUsers, "SELECT COUNT(id) FROM public.tm_user WHERE role=$1", role.ID); err != nil {
		return nil, errors.New("role delete counting assigned users: " + err.Error()), http.StatusInternalServerError
	} else if assignedUsers != 0 {
		return fmt.Errorf("can not delete a role with %d assigned users", assignedUsers), nil, http.StatusBadRequest
	}

	var isAdmin bool
	if err := role.ReqInfo.Tx.Get(&isAdmin, isAdminQuery, role.ID); err != nil {
		return nil, fmt.Errorf("checking if Role to be deleted is '%s': %w", tc.AdminRoleName, err), http.StatusInternalServerError
	}
	if isAdmin {
		return cannotModifyAdminError, nil, http.StatusBadRequest
	}

	userErr, sysErr, errCode := api.GenericDelete(role)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	return role.deleteRoleCapabilityAssociations(role.ReqInfo.Tx)
}

func selectQuery() string {
	return `SELECT
id,
name,
description,
priv_level,
ARRAY(SELECT rc.cap_name FROM role_capability AS rc WHERE rc.role_id=id) AS capabilities
FROM role`
}

func updateQuery() string {
	return `UPDATE
role SET
name=:name,
description=:description
WHERE id=:id RETURNING last_updated`
}

func deleteAssociatedCapabilities() string {
	return `DELETE FROM role_capability
WHERE role_id=(SELECT id from role r WHERE r.name=$1)`
}

func associateCapabilities() string {
	return `INSERT INTO role_capability (
role_id,
cap_name) WITH
	q1 AS ( SELECT * FROM (VALUES ($1::bigint)) AS role_id ),
	q2 AS (SELECT UNNEST($2::text[]))
	SELECT * FROM q1,q2`
}

func insertQuery() string {
	return `INSERT INTO role (
name,
description,
priv_level
) VALUES (
:name,
:description,
:priv_level
)
RETURNING id, last_updated`
}

func deleteQuery() string {
	return `DELETE FROM role WHERE id = :id`
}

// Update will modify the role identified by the role name.
func Update(w http.ResponseWriter, r *http.Request) {
	var roleV4 tc.RoleV4

	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	if err := json.NewDecoder(r.Body).Decode(&roleV4); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if err := roleV4.Validate(); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	currentRoleName := inf.Params["name"]
	if currentRoleName == tc.AdminRoleName {
		api.HandleErr(w, r, tx, http.StatusBadRequest, cannotModifyAdminError, nil)
		return
	}

	missing := inf.User.MissingPermissions(roleV4.Permissions...)
	if len(missing) != 0 {
		api.HandleErr(w, r, tx, http.StatusForbidden, fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil)
		return
	}

	roleID, ok, err := dbhelpers.GetRoleIDFromName(tx, currentRoleName)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	} else if !ok {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no such role"), nil)
		return
	}

	existingLastUpdated, found, err := api.GetLastUpdatedByName(inf.Tx, currentRoleName, "role")
	if err == nil && found == false {
		api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no such role"), nil)
		return
	}
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	}

	if !api.IsUnmodified(r.Header, *existingLastUpdated) {
		api.HandleErr(w, r, tx, http.StatusPreconditionFailed, api.ResourceModifiedError, nil)
		return
	}
	err = tx.QueryRow(updateRoleQuery(), roleV4.Name, roleV4.Description, currentRoleName).Scan(&roleV4.LastUpdated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no such role"), nil)
			return
		}
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, fmt.Errorf("updating role and scanning lastUpdated : %w", sysErr))
		return
	}

	userErr, sysErr, errCode = deleteRoleCapabilityAssociations(inf.Tx, roleV4.Name)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	userErr, sysErr, errCode = createRoleCapabilityAssociations(inf.Tx, roleID, &roleV4.Permissions)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "role was updated.")
	// to return empty array instead of null
	if roleV4.Permissions == nil {
		roleV4.Permissions = []string{}
	}
	var roleResponse interface{}
	roleResponse = tc.RoleV4{
		Name:        roleV4.Name,
		Permissions: roleV4.Permissions,
		Description: roleV4.Description,
		LastUpdated: roleV4.LastUpdated,
	}
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, roleResponse)
	changeLogMsg := fmt.Sprintf("ROLE: %s, ID: %d, ACTION: Updated Role", roleV4.Name, roleID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func deleteRoleQuery() string {
	return `DELETE FROM role WHERE name = $1`
}

func readQuery() string {
	return `SELECT
name,
description,
last_updated,
ARRAY(SELECT rc.cap_name FROM role_capability AS rc WHERE rc.role_id=id) AS permissions
FROM role`
}

func createQuery() string {
	return `INSERT INTO role (
name,
description,
priv_level
) VALUES (
$1,
$2,
$3
)
RETURNING id, last_updated`
}

func createRoleCapabilityAssociations(tx *sqlx.Tx, roleID int, permissions *[]string) (error, error, int) {
	result, err := tx.Exec(associateCapabilities(), roleID, pq.Array(permissions))
	if err != nil {
		return nil, fmt.Errorf("creating role capabilities: %w", err), http.StatusInternalServerError
	}

	if rows, err := result.RowsAffected(); err != nil {
		return nil, fmt.Errorf("could not check result after inserting role_capability relations: %w", err), http.StatusInternalServerError
	} else if expected := int64(len(*permissions)); rows != expected {
		return nil, fmt.Errorf("wrong number of role_capability rows created: %d expected: %d", rows, expected), http.StatusInternalServerError
	}
	return nil, nil, http.StatusOK
}

func deleteRoleCapabilityAssociations(tx *sqlx.Tx, roleName string) (error, error, int) {
	_, err := tx.Exec(deleteAssociatedCapabilities(), roleName)
	if err != nil {
		return nil, fmt.Errorf("deleting role capabilities: %w", err), http.StatusInternalServerError
	}
	return nil, nil, http.StatusOK
}

// Delete will delete the role identified by the role name.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	roleName := inf.Params["name"]

	if roleName == tc.AdminRoleName {
		api.HandleErr(w, r, tx, http.StatusBadRequest, cannotModifyAdminError, nil)
		return
	}

	assignedUsers := 0
	if err := inf.Tx.Get(&assignedUsers, "SELECT COUNT(id) FROM tm_user WHERE role= (SELECT id FROM role r WHERE r.name=$1)", roleName); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("role delete counting assigned users: %w", err))
		return
	} else if assignedUsers != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, fmt.Errorf("can not delete a role with %d assigned users", assignedUsers), nil)
		return
	}

	rows, err := tx.Query(deleteRoleQuery(), roleName)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("deleting role: %w", err))
		return
	}
	rows.Close()
	alerts := tc.CreateAlerts(tc.SuccessLevel, "role was deleted.")
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	changeLogMsg := fmt.Sprintf("ROLE: %s, ACTION: Deleted Role", roleName)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Create will create a new role based on the struct supplied.
func Create(w http.ResponseWriter, r *http.Request) {
	var roleID int
	var roleName string
	var roleDesc string
	var privLevel int
	var roleCapabilities []string
	var lastUpdated time.Time
	var roleV4 tc.RoleV4

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	if err := json.NewDecoder(r.Body).Decode(&roleV4); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	if err := roleV4.Validate(); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}
	missing := inf.User.MissingPermissions(roleV4.Permissions...)
	if len(missing) != 0 {
		api.HandleErr(w, r, tx, http.StatusForbidden, fmt.Errorf("cannot request more than assigned permissions, current user needs %s permissions", strings.Join(missing, ",")), nil)
		return
	}
	roleName = roleV4.Name
	roleDesc = roleV4.Description
	privLevel = inf.User.PrivLevel
	roleCapabilities = roleV4.Permissions

	rows, err := tx.Query(createQuery(), roleName, roleDesc, privLevel)
	if err != nil {
		usrErr, sysErr, code := api.ParseDBError(err)
		api.HandleErr(w, r, tx, code, usrErr, fmt.Errorf("creating role: %w", sysErr))
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&roleID, &lastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("role create: scanning role ID: %w", err))
			return
		}
	}

	if len(roleCapabilities) > 0 {
		userErr, sysErr, errCode = createRoleCapabilityAssociations(inf.Tx, roleID, &roleCapabilities)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
	}
	alerts := tc.CreateAlerts(tc.SuccessLevel, "role was created.")
	// to return empty array instead of null
	if roleCapabilities == nil {
		roleCapabilities = []string{}
	}
	var roleResponse interface{}
	capabilities := roleCapabilities
	roleResponse = tc.RoleV4{
		Name:        roleName,
		Permissions: capabilities,
		Description: roleDesc,
		LastUpdated: &lastUpdated,
	}
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, roleResponse)
	changeLogMsg := fmt.Sprintf("ROLE: %s, ID: %d, ACTION: Created Role", roleName, roleID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Get will read the roles and return them to the user.
func Get(w http.ResponseWriter, r *http.Request) {
	var maxTime time.Time
	var runSecond bool
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx
	params := make(map[string]dbhelpers.WhereColumnInfo, 1)
	params["name"] = dbhelpers.WhereColumnInfo{Column: "name"}

	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = "name"
	}
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, params)
	if len(errs) != 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}
	if perm, ok := inf.Params["can"]; ok {
		queryValues["can"] = perm
		where = dbhelpers.AppendWhere(where, "permissions @> :can")
	}
	if inf.Config.UseIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(inf.Tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			api.AddLastModifiedHdr(w, maxTime)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	query := readQuery() + where + orderBy + pagination

	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("querying Roles: %w", err))
		return
	}
	defer log.Close(rows, "reading in Roles from the database")

	var roleV4 tc.RoleV4
	rolesV4 := []tc.RoleV4{}

	for rows.Next() {
		if err = rows.Scan(&roleV4.Name, &roleV4.Description, &roleV4.LastUpdated, pq.Array(&roleV4.Permissions)); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("scanning RoleV4 row: %w", err))
			return
		}
		sort.Strings(roleV4.Permissions)
		rolesV4 = append(rolesV4, roleV4)
	}
	api.WriteResp(w, r, rolesV4)
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) FROM (
		SELECT max(r.last_updated) AS t FROM role r ` + where + ` UNION ALL
		SELECT max(l.last_updated) AS t FROM last_deleted l WHERE l.table_name='role' OR l.table_name='role_capability' UNION ALL
		SELECT max(rc.last_updated) AS t FROM role_capability rc INNER JOIN role ON rc.role_id = role.id)
		AS res`
}
