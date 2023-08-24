package apitenant

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

// tenant.go defines the TOTenant object and methods/functions required for the api/.../tenants endpoints

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	time "time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

const rootName = `root`

// TOTenant provides a local type against which to define methods
type TOTenant struct {
	api.APIInfoImpl `json:"-"`
	tc.TenantNullable
}

func (ten *TOTenant) GetLastUpdated() (*time.Time, bool, error) {
	return api.GetLastUpdated(ten.APIInfo().Tx, *ten.ID, "tenant")
}

func (ten *TOTenant) SetLastUpdated(t tc.TimeNoMod) { ten.LastUpdated = &t }
func (ten *TOTenant) InsertQuery() string           { return insertQuery() }
func (ten *TOTenant) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` q ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}
func (ten *TOTenant) NewReadObj() interface{} { return &tc.TenantNullable{} }
func (ten *TOTenant) SelectQuery() string {
	return selectQuery(ten.APIInfo().User.TenantID)
}
func (ten *TOTenant) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"active":      dbhelpers.WhereColumnInfo{Column: "q.active", Checker: nil},
		"id":          dbhelpers.WhereColumnInfo{Column: "q.id", Checker: api.IsInt},
		"name":        dbhelpers.WhereColumnInfo{Column: "q.name", Checker: nil},
		"parent_id":   dbhelpers.WhereColumnInfo{Column: "q.parent_id", Checker: api.IsInt},
		"parent_name": dbhelpers.WhereColumnInfo{Column: "p.name", Checker: nil},
	}
}
func (ten *TOTenant) UpdateQuery() string { return updateQuery() }

// GetID wraps the ID member with null checking
// Part of the Identifier interface
func (ten TOTenant) GetID() (int, bool) {
	if ten.ID == nil {
		return 0, false
	}
	return *ten.ID, true
}

// GetKeyFieldsInfo identifies types of the key fields
func (ten TOTenant) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{Field: "id", Func: api.GetIntKey}}
}

// GetKeys returns values of keys
func (ten TOTenant) GetKeys() (map[string]interface{}, bool) {
	var id int
	if ten.ID != nil {
		id = *ten.ID
	}
	return map[string]interface{}{"id": id}, true
}

// GetAuditName returns a unique identifier
// Part of the Identifier interface
func (ten TOTenant) GetAuditName() string {
	if ten.Name != nil {
		return *ten.Name
	}
	id, _ := ten.GetID()
	return strconv.Itoa(id)
}

// GetType returns the name of the type for messages
// Part of the Identifier interface
func (ten TOTenant) GetType() string {
	return "tenant"
}

// SetKeys allows CreateHandler to assign id once object is created.
// Part of the Identifier interface
func (ten *TOTenant) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	ten.ID = &i
}

// Validate fulfills the api.Validator interface.
func (ten TOTenant) Validate() (error, error) {
	errs := validation.Errors{
		"name":       validation.Validate(ten.Name, validation.Required),
		"active":     validation.Validate(ten.Active), // only validate it's boolean
		"parentId":   validation.Validate(ten.ParentID, validation.Required, validation.Min(1)),
		"parentName": nil,
	}
	return util.JoinErrs(tovalidate.ToErrors(errs)), nil
}

func (ten *TOTenant) Create() (error, error, int) { return api.GenericCreate(ten) }

func (ten *TOTenant) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	if ten.APIInfo().User.TenantID == auth.TenantIDInvalid {
		return nil, nil, nil, http.StatusOK, nil
	}
	api.DefaultSort(ten.APIInfo(), "name")
	tenants, userErr, sysErr, errCode, maxTime := api.GenericRead(h, ten, useIMS)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, nil
	}

	tenantNames := map[int]*string{}
	for _, it := range tenants {
		t := it.(*tc.TenantNullable)
		tenantNames[*t.ID] = t.Name
	}
	for _, it := range tenants {
		t := it.(*tc.TenantNullable)
		if t.ParentID == nil || tenantNames[*t.ParentID] == nil {
			// root tenant has no parent
			continue
		}
		p := *tenantNames[*t.ParentID]
		t.ParentName = &p // copy
	}
	return tenants, nil, nil, errCode, maxTime
}

// IsTenantAuthorized implements the Tenantable interface for TOTenant
// returns true if the user has access on this tenant and on the ParentID if changed.
func (ten *TOTenant) IsTenantAuthorized(user *auth.CurrentUser) (bool, error) {
	var ok = false
	var err error

	if ten == nil {
		// should never happen
		return ok, err
	}

	if ten.ID != nil && *ten.ID != 0 {
		// modifying an existing tenant
		ok, err = tenant.IsResourceAuthorizedToUserTx(*ten.ID, user, ten.APIInfo().Tx.Tx)
		if !ok || err != nil {
			return ok, err
		}

		if ten.ParentID == nil || *ten.ParentID == 0 {
			// not changing parent
			return ok, err
		}

		// get current parentID to check if it's being changed
		var parentID int
		tx := ten.APIInfo().Tx.Tx
		// If it's the root tenant, don't check for parent
		if ten.Name != nil && *ten.Name != rootName {
			err = tx.QueryRow(`SELECT parent_id FROM tenant WHERE id = ` + strconv.Itoa(*ten.ID)).Scan(&parentID)
			if err != nil {
				return false, err
			}
			if parentID == *ten.ParentID {
				// parent not being changed
				return ok, err
			}
		}
	}

	// creating new tenant -- must specify a parent
	if ten.ParentID == nil || *ten.ParentID == 0 {
		return false, err
	}

	return tenant.IsResourceAuthorizedToUserTx(*ten.ParentID, user, ten.APIInfo().Tx.Tx)
}

// Update wraps tenant validation and the generic API Update call into a single call.
func (ten *TOTenant) Update(h http.Header) (error, error, int) {

	userErr, sysErr, statusCode := ten.isUpdatable()
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
	}

	return api.GenericUpdate(h, ten)
}

// isUpdatable peforms validation on the fields for the Tenant, such as ensuring
// the tenant cannot be modified if it is root, or that it cannot convert its own child
// to its own parent. This is different than the basic validation rules performed in
// Validate() as it pertains to specific business logic, not generic API rules.
func (ten *TOTenant) isUpdatable() (error, error, int) {
	if ten.Name != nil && *ten.Name == rootName {
		return errors.New("trying to change the root tenant, which is immutable"), nil, http.StatusBadRequest
	}

	// Perform SelectQuery
	vals := []tc.TenantNullable{}
	query := selectQuery(*ten.ID)
	rows, err := ten.APIInfo().Tx.Queryx(query)
	if err != nil {
		return nil, errors.New("querying " + ten.GetType() + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var v tc.TenantNullable
		if err = rows.StructScan(&v); err != nil {
			return nil, errors.New("scanning " + ten.GetType() + ": " + err.Error()), http.StatusInternalServerError
		}
		vals = append(vals, v)
	}

	// Ensure the new desired ParentID does not exist in the susequent list of Children
	for _, val := range vals {
		if *ten.ParentID == *val.ID {
			return errors.New("trying to set existing child as new parent"), nil, http.StatusBadRequest
		}
	}
	return nil, nil, http.StatusOK
}

func (ten *TOTenant) Delete() (error, error, int) {
	result, err := ten.APIInfo().Tx.NamedExec(deleteQuery(), ten)
	if err != nil {
		return parseDeleteErr(err, *ten.ID, ten.APIInfo().Tx.Tx) // this is why we can't use api.GenericDelete
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return nil, errors.New("deleting " + ten.GetType() + ": getting rows affected: " + err.Error()), http.StatusInternalServerError
	} else if rowsAffected < 1 {
		return errors.New("no " + ten.GetType() + " with that id found"), nil, http.StatusNotFound
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf(ten.GetType()+" delete affected too many rows: %d", rowsAffected), http.StatusInternalServerError
	}
	return nil, nil, http.StatusOK
}

// parseDeleteErr takes the tenant delete error, and returns the appropriate user error, system error, and http status code.
func parseDeleteErr(err error, id int, tx *sql.Tx) (error, error, int) {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return nil, errors.New("deleting tenant: " + err.Error()), http.StatusInternalServerError
	}
	// TODO fix this to check for other Postgres errors besides key violations
	existing := ""
	switch pqErr.Table {
	case "tenant":
		existing = "child tenants"
	case "tm_user":
		existing = "users"
	case "deliveryservice":
		existing = "deliveryservices"
	case "origin":
		existing = "origins"
	default:
		existing = pqErr.Table
	}
	return errors.New("Tenant '" + strconv.Itoa(id) + "' has " + existing + ". Please update these " + existing + " and retry."), nil, http.StatusBadRequest
}

// selectQuery returns a query on the tenant table that limits to tenants within the realm of the tenantID.  It's intended
// to be extensible by adding WHERE/ORDERBY/etc clauses at the end to further refine the query.
func selectQuery(tenantID int) string {
	query := `
WITH RECURSIVE q AS (
SELECT id, name, active, parent_id, last_updated FROM tenant WHERE id = ` + strconv.Itoa(tenantID) + `
UNION SELECT t.id, t.name, t.active, t.parent_id, t.last_updated FROM tenant t JOIN q ON q.id = t.parent_id)
SELECT q.id AS id, q.name AS name, q.active AS active, q.parent_id AS parent_id, q.last_updated AS last_updated,
p.name AS parent_name FROM q LEFT JOIN tenant p ON q.parent_id = p.id
`

	return query
}

func updateQuery() string {
	query := `UPDATE
tenant SET
active=:active,
name=:name,
parent_id=:parent_id

WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO tenant (
name,
active,
parent_id
) VALUES (
:name,
:active,
:parent_id
) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM tenant
WHERE id=:id`
	return query
}

// Downgrade will convert an instance of CacheGroupNullableV5 to CacheGroupNullable.
// Note that this function does a shallow copy of the requested and original Cache Group structures.
func Downgrade(tV5 tc.TenantV5) TOTenant {
	var t TOTenant
	t.ID = util.CopyIfNotNil(tV5.ID)
	t.Name = util.CopyIfNotNil(tV5.Name)
	t.Active = util.CopyIfNotNil(tV5.Active)
	if tV5.LastUpdated != nil {
		t.LastUpdated = tc.TimeNoModFromTime(*tV5.LastUpdated)
	}
	t.ParentID = util.CopyIfNotNil(tV5.ParentID)
	t.ParentName = util.CopyIfNotNil(tV5.ParentName)
	return t
}

// Upgrade will convert an instance of CacheGroupNullable to CacheGroupNullableV5.
// Note that this function does a shallow copy of the requested and original Cache Group structures.
func (t TOTenant) Upgrade() (tc.TenantV5, error) {
	var tV5 tc.TenantV5
	tV5.ID = util.CopyIfNotNil(t.ID)
	tV5.Name = util.CopyIfNotNil(t.Name)
	tV5.Active = util.CopyIfNotNil(t.Active)
	if t.LastUpdated != nil {
		tV5.LastUpdated = &t.LastUpdated.Time
		t, err := util.ConvertTimeFormat(*tV5.LastUpdated, time.RFC3339)
		if err != nil {
			return tV5, err
		}
		tV5.LastUpdated = t
	}
	tV5.ParentID = util.CopyIfNotNil(t.ParentID)
	tV5.ParentName = util.CopyIfNotNil(t.ParentName)
	return tV5, nil
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(t) from (
		SELECT max(tenant.last_updated) as t from tenant + q ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='tenant') as res`
}

// CreateTenant [Version : V5] function Process the *http.Request and creates new tenant
func CreateTenant(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	defer r.Body.Close()

	tenant, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	resultRows, err := inf.Tx.NamedQuery(insertQuery(), tenant)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	var id int
	lastUpdated := time.Time{}
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("tenant create scanning: "+err.Error()))
			return
		}
	}

	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("tenant create: no tenant was inserted, no id was returned"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("too many ids returned from tenant insert"))
		return
	}
	tenant.ID = &id
	tenant.LastUpdated = &lastUpdated

	alerts := tc.CreateAlerts(tc.SuccessLevel, "tenant was created.")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, tenant)
	return
}

// GetTenant [Version : V5] function Process the *http.Request and writes the response. It uses getTenant function.
func GetTenant(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	user, _ := auth.GetCurrentUser(r.Context())
	tenantID := user.TenantID
	if tenantID == auth.TenantIDInvalid {
		return
	}

	code := http.StatusOK
	useIMS := false
	config, e := api.GetConfig(r.Context())
	if e == nil && config != nil {
		useIMS = config.UseIMS
	} else {
		log.Warnf("Couldn't get config %v", e)
	}

	var maxTime time.Time
	var usrErr error
	var syErr error

	var tenants []tc.TenantV5

	api.DefaultSort(inf, "name")
	tenants, usrErr, syErr, code, maxTime = getTenants(inf.Tx, inf.Params, useIMS, r.Header, tenantID)
	if userErr != nil {
		api.HandleErr(w, r, tx, code, fmt.Errorf("read tenant: get tenant: "+usrErr.Error()), nil)
	}
	if sysErr != nil {
		api.HandleErr(w, r, tx, code, nil, fmt.Errorf("read tenant: get tenant: "+syErr.Error()))
	}
	if maxTime != (time.Time{}) && api.SetLastModifiedHeader(r, useIMS) {
		api.AddLastModifiedHdr(w, maxTime)
	}

	tenantNames := map[int]*string{}
	for _, it := range tenants {
		tenantNames[*it.ID] = it.Name
	}
	for _, it := range tenants {
		if it.ParentID == nil || tenantNames[*it.ParentID] == nil {
			// root tenant has no parent
			continue
		}
		p := *tenantNames[*it.ParentID]
		it.ParentName = &p // copy
	}
	api.WriteResp(w, r, tenants)
}

func getTenants(tx *sqlx.Tx, params map[string]string, useIMS bool, header http.Header, id int) ([]tc.TenantV5, error, error, int, time.Time) {
	tenants := make([]tc.TenantV5, 0)
	code := http.StatusOK

	var maxTime time.Time
	var runSecond bool

	// Query Parameters to Database Query column mappings
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"active":      {Column: "q.active", Checker: nil},
		"name":        {Column: "q.name", Checker: nil},
		"parent_name": {Column: "p.name", Checker: nil},
		"id":          {Column: "q.id", Checker: api.IsInt},
		"parent_id":   {Column: "q.parent_id", Checker: api.IsInt},
	}
	if _, ok := params["orderby"]; !ok {
		params["orderby"] = "name"
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, time.Time{}
	}

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, header, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			code = http.StatusNotModified
			return tenants, nil, nil, code, maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	// Case where we need to run the second query
	query := selectQuery(id) + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, err, http.StatusInternalServerError, time.Time{}
	}
	defer rows.Close()

	for rows.Next() {
		var t tc.TenantV5

		if err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Active,
			&t.ParentID,
			&t.LastUpdated,
			&t.ParentName,
		); err != nil {
			return nil, nil, err, http.StatusInternalServerError, time.Time{}
		}
		tenants = append(tenants, t)
	}

	return tenants, nil, nil, http.StatusOK, maxTime
}

// UpdateTenant [Version : V5] function updates the name of the tenant passed.
func UpdateTenant(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	t, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	dgT := Downgrade(t)

	// Check that tenant can be updated
	userErr, sysErr, statusCode := dgT.isUpdatable()
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
	}

	existingLastUpdated, found, err := dgT.GetLastUpdated()
	if err == nil && found == false {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("update tenant: find last updated: no "+dgT.GetType()+" found with this id"), nil)
	}
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("update tenant: find last updated: "+err.Error()), nil)
	}
	if !api.IsUnmodified(r.Header, *existingLastUpdated) {
		api.HandleErr(w, r, tx, http.StatusPreconditionFailed, fmt.Errorf(string("update tenant: check if unmodified: "+api.ResourceModifiedError)), nil)
	}

	rows, err := dgT.APIInfo().Tx.NamedQuery(dgT.UpdateQuery(), dgT)
	if err != nil {
		userErr, sysErr, errCode = api.ParseDBError(err)
		if userErr != nil {
			api.HandleErr(w, r, tx, errCode, fmt.Errorf("update tenant: get rows: "+userErr.Error()), nil)
		}
		if sysErr != nil {
			api.HandleErr(w, r, tx, errCode, nil, fmt.Errorf("update tenant: get rows: "+sysErr.Error()))
		}
	}
	defer rows.Close()

	if !rows.Next() {
		api.HandleErr(w, r, tx, http.StatusNotFound, fmt.Errorf("update tenant: get rows: no "+dgT.GetType()+" found with this id"), nil)
	}
	lastUpdated := tc.TimeNoMod{}
	if err := rows.Scan(&lastUpdated); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("update tenant: get rows: scanning lastUpdated from "+dgT.GetType()+" insert: "+err.Error()))
	}
	dgT.SetLastUpdated(lastUpdated)
	if rows.Next() {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("update tenant: get rows: "+dgT.GetType()+" update affected too many rows: >1"))
	}

	t, err = dgT.Upgrade()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, fmt.Errorf("converting cachegroup: converting cache group upgrade: "+err.Error()), nil)
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "tenant was updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, t)
	return
}

// DeleteTenant [Version : V5] function deletes the tenant passed.
func DeleteTenant(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ID := inf.Params["id"]
	id, err := strconv.Atoi(ID)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusUnprocessableEntity, fmt.Errorf("delete cachegroup: converted to type int: "+err.Error()), nil)
		return
	}

	res, err := tx.Exec("DELETE FROM tenant AS t WHERE t.ID=$1", id)
	if err != nil {
		usrErr, syErr, code := parseDeleteErr(err, id, tx)
		api.HandleErr(w, r, tx, code, usrErr, syErr)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("determining rows affected for delete cachegroup: %w", err))
		return
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("no rows deleted for cachegroup"))
		return
	}

	alertMessage := fmt.Sprint("tenant was deleted.")
	alerts := tc.CreateAlerts(tc.SuccessLevel, alertMessage)
	api.WriteAlerts(w, r, http.StatusOK, alerts)
	return
}

// readAndValidateJsonStruct populates select missing fields and validates JSON body
func readAndValidateJsonStruct(r *http.Request) (tc.TenantV5, error) {
	var ten tc.TenantV5
	if err := json.NewDecoder(r.Body).Decode(&ten); err != nil {
		userErr := fmt.Errorf("error decoding POST request body into TenantV5 struct %w", err)
		return ten, userErr
	}

	// validate JSON body
	rule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters")
	errs := tovalidate.ToErrors(validation.Errors{
		"name":       validation.Validate(ten.Name, validation.Required, rule),
		"active":     validation.Validate(ten.Active), // only validate it's boolean
		"parentId":   validation.Validate(ten.ParentID, validation.Required, validation.Min(1)),
		"parentName": nil,
	})
	if len(errs) > 0 {
		userErr := util.JoinErrs(errs)
		return ten, userErr
	}
	return ten, nil
}
