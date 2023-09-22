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
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
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

func selectMaxLastUpdatedQuery(where string) string {
	tableName := "tenant"
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` q ` + where +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
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

	// Check if tenant is tenable
	authorized, err := isTenantAuthorizedV5(&tenant, inf.User, tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("checking tenant authorized: "+err.Error()))
		return
	}
	if !authorized {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
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
	changeLogMsg := fmt.Sprintf("TENANT: %s, ID: %d, ACTION: Created tenant", *tenant.Name, *tenant.ID)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
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

	var maxTime *time.Time
	var usrErr error
	var syErr error

	var tenants []tc.TenantV5

	api.DefaultSort(inf, "name")
	tenants, usrErr, syErr, code, maxTime = getTenants(inf.Tx, inf.Params, useIMS, r.Header, tenantID)
	if usrErr != nil {
		api.HandleErr(w, r, tx, code, fmt.Errorf("read tenant: get tenant: "+usrErr.Error()), nil)
	}
	if syErr != nil {
		api.HandleErr(w, r, tx, code, nil, fmt.Errorf("read tenant: get tenant: "+syErr.Error()))
	}
	if maxTime != nil && api.SetLastModifiedHeader(r, useIMS) {
		api.AddLastModifiedHdr(w, *maxTime)
		w.WriteHeader(http.StatusNotModified)
		return
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

func getTenants(tx *sqlx.Tx, params map[string]string, useIMS bool, header http.Header, id int) ([]tc.TenantV5, error, error, int, *time.Time) {
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
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(tx, header, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			code = http.StatusNotModified
			return tenants, nil, nil, code, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	// Case where we need to run the second query
	query := selectQuery(id) + where + orderBy + pagination
	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, err, http.StatusInternalServerError, nil
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
			return nil, nil, err, http.StatusInternalServerError, nil
		}
		tenants = append(tenants, t)
	}

	return tenants, nil, nil, http.StatusOK, nil
}

// UpdateTenant [Version : V5] function Process the *http.Request and updates the tenant.
func UpdateTenant(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	defer r.Body.Close()
	var tenant tc.TenantV5

	tenant, readValErr := readAndValidateJsonStruct(r)
	if readValErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, readValErr, nil)
		return
	}

	if id, ok := inf.Params["id"]; !ok {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("missing key: id"), nil)
		return
	} else {
		idNum, err := strconv.Atoi(id)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("couldn't convert ID into a numeric value: "+err.Error()), nil)
			return
		}
		tenant.ID = &idNum

		existingLastUpdated, found, err := api.GetLastUpdated(inf.Tx, idNum, "tenant")
		if err == nil && found == false {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no tenant found with this id"), nil)
			return
		}
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusNotFound, nil, err)
			return
		}
		if !api.IsUnmodified(r.Header, *existingLastUpdated) {
			api.HandleErr(w, r, tx, http.StatusPreconditionFailed, api.ResourceModifiedError, nil)
			return
		}

		// Check if tenant is tenable
		authorized, err := isTenantAuthorizedV5(&tenant, inf.User, tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("checking tenant authorized: "+err.Error()))
			return
		}
		if !authorized {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
			return
		}

		//Check if tenant is updatable
		userErr, sysErr, statusCode := isUpdatableV5(&tenant, inf.Tx)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}

		rows, err := inf.Tx.NamedQuery(updateQuery(), tenant)
		if err != nil {
			userErr, sysErr, errCode = api.ParseDBError(err)
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no tenant found with this id"), nil)
			return
		}
		lastUpdated := time.Time{}
		if err := rows.Scan(&lastUpdated); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning lastUpdated from tenant insert: "+err.Error()))
			return
		}
		tenant.LastUpdated = &lastUpdated
		if rows.Next() {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("tenant update affected too many rows: >1"))
			return
		}

		alerts := tc.CreateAlerts(tc.SuccessLevel, "tenant was updated.")
		api.WriteAlertsObj(w, r, http.StatusOK, alerts, tenant)
		changeLogMsg := fmt.Sprintf("TENANT: %s, ID: %d, ACTION: Updated tenant", *tenant.Name, *tenant.ID)
		api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
		return
	}
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

	useIMS := false
	var tenantV5 []tc.TenantV5

	tenantV5, usrErr, syErr, code, maxTime := getTenants(inf.Tx, inf.Params, useIMS, r.Header, id)
	if userErr != nil {
		api.HandleErr(w, r, tx, code, fmt.Errorf("delete tenant: get tenant: "+usrErr.Error()), nil)
	}
	if sysErr != nil {
		api.HandleErr(w, r, tx, code, nil, fmt.Errorf("delete tenant: get tenant: "+syErr.Error()))
	}
	if maxTime != nil && api.SetLastModifiedHeader(r, useIMS) {
		api.AddLastModifiedHdr(w, *maxTime)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Check if tenant is tenable
	authorized, err := isTenantAuthorizedV5(&tenantV5[0], inf.User, tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("checking tenant authorized: "+err.Error()))
		return
	}
	if !authorized {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
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
	changeLogMsg := fmt.Sprintf("TENANT: ID: %d, ACTION: Deleted tenant", id)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
	return
}

// IsTenantAuthorized implements the Tenantable interface for TOTenant
// returns true if the user has access on this tenant and on the ParentID if changed.
func isTenantAuthorizedV5(tenantV5 *tc.TenantV5, user *auth.CurrentUser, tx *sql.Tx) (bool, error) {
	var ok = false
	var err error

	if tenantV5 == nil {
		// should never happen
		return ok, err
	}

	if tenantV5.ID != nil && *tenantV5.ID != 0 {
		// modifying an existing tenant
		ok, err = tenant.IsResourceAuthorizedToUserTx(*tenantV5.ID, user, tx)
		if !ok || err != nil {
			return ok, err
		}

		if tenantV5.ParentID == nil || *tenantV5.ParentID == 0 {
			// not changing parent
			return ok, err
		}

		// get current parentID to check if it's being changed
		var parentID int
		// If it's the root tenant, don't check for parent
		if tenantV5.Name != nil && *tenantV5.Name != rootName {
			err = tx.QueryRow(`SELECT parent_id FROM tenant WHERE id = ` + strconv.Itoa(*tenantV5.ID)).Scan(&parentID)
			if err != nil {
				return false, err
			}
			if parentID == *tenantV5.ParentID {
				// parent not being changed
				return ok, err
			}
		}
	}

	// creating new tenant -- must specify a parent
	if tenantV5.ParentID == nil || *tenantV5.ParentID == 0 {
		return false, err
	}

	return tenant.IsResourceAuthorizedToUserTx(*tenantV5.ParentID, user, tx)
}

func isUpdatableV5(tenantV5 *tc.TenantV5, tx *sqlx.Tx) (error, error, int) {
	if tenantV5.Name != nil && *tenantV5.Name == rootName {
		return errors.New("trying to change the root tenant, which is immutable"), nil, http.StatusBadRequest
	}

	// Perform SelectQuery
	vals := []tc.TenantNullable{}
	query := selectQuery(*tenantV5.ID)
	rows, err := tx.Queryx(query)
	if err != nil {
		return nil, errors.New("querying tenant: " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var v tc.TenantNullable
		if err = rows.StructScan(&v); err != nil {
			return nil, errors.New("scanning tenant: " + err.Error()), http.StatusInternalServerError
		}
		vals = append(vals, v)
	}

	// Ensure the new desired ParentID does not exist in the susequent list of Children
	for _, val := range vals {
		if *tenantV5.ParentID == *val.ID {
			return errors.New("trying to set existing child as new parent"), nil, http.StatusBadRequest
		}
	}
	return nil, nil, http.StatusOK
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
