package tenant

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
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// DeliveryServiceTenantInfo provides only deliveryservice info needed here
type DeliveryServiceTenantInfo tc.DeliveryServiceNullable

// IsTenantAuthorized returns true if the user has tenant access on this tenant
func (dsInfo DeliveryServiceTenantInfo) IsTenantAuthorized(user *auth.CurrentUser, tx *sql.Tx) (bool, error) {
	if dsInfo.TenantID == nil {
		return false, errors.New("TenantID is nil")
	}
	return IsResourceAuthorizedToUserTx(*dsInfo.TenantID, user, tx)
}

// returns tenant information for a deliveryservice
func GetDeliveryServiceTenantInfo(xmlID string, tx *sql.Tx) (*DeliveryServiceTenantInfo, error) {
	ds := DeliveryServiceTenantInfo{XMLID: util.StrPtr(xmlID)}
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where xml_id = $1`, &ds.XMLID).Scan(&ds.TenantID); err != nil {
		if err == sql.ErrNoRows {
			return &ds, errors.New("a deliveryservice with xml_id '" + xmlID + "' was not found")
		}
		return nil, errors.New("querying tenant id from delivery service: " + err.Error())
	}
	return &ds, nil
}

// Check checks that the given user has access to the given XMLID. Returns a user error, system error,
// and the HTTP status code to be returned to the user if an error occurred. On success, the user error
// and system error will both be nil, and the error code should be ignored.
func Check(user *auth.CurrentUser, XMLID string, tx *sql.Tx) (error, error, int) {
	dsInfo, err := GetDeliveryServiceTenantInfo(XMLID, tx)
	if err != nil {
		if dsInfo == nil {
			return nil, errors.New("deliveryservice lookup failure: " + err.Error()), http.StatusInternalServerError
		}
		return errors.New("no such deliveryservice: '" + XMLID + "'"), nil, http.StatusBadRequest
	}
	hasAccess, err := dsInfo.IsTenantAuthorized(user, tx)
	if err != nil {
		return nil, errors.New("user tenancy check failure: " + err.Error()), http.StatusInternalServerError
	}
	if !hasAccess {
		return nil, errors.New("Access to this resource is not authorized"), http.StatusForbidden
	}
	return nil, nil, http.StatusOK
}

// Check checks that the given user has access to the given delivery service. Returns a user error, a system error, and an HTTP error code. If both the user and system error are nil, the error code should be ignored.
func CheckID(tx *sql.Tx, user *auth.CurrentUser, dsID int) (error, error, int) {
	ok, err := IsTenancyEnabledTx(tx)
	if err != nil {
		return nil, errors.New("checking tenancy enabled: " + err.Error()), http.StatusInternalServerError
	}
	if !ok {
		return nil, nil, http.StatusOK
	}

	dsTenantID, ok, err := getDSTenantIDByIDTx(tx, dsID)
	if err != nil {
		return nil, errors.New("checking tenant: " + err.Error()), http.StatusInternalServerError
	}
	if !ok {
		return errors.New("delivery service " + strconv.Itoa(dsID) + " not found"), nil, http.StatusNotFound
	}
	if dsTenantID == nil {
		return nil, nil, http.StatusOK
	}

	authorized, err := IsResourceAuthorizedToUserTx(*dsTenantID, user, tx)
	if err != nil {
		return nil, errors.New("checking tenant: " + err.Error()), http.StatusInternalServerError
	}
	if !authorized {
		return errors.New("not authorized on this tenant"), nil, http.StatusForbidden
	}
	return nil, nil, http.StatusOK
}

// returns a Tenant list that the specified user has access too.
// NOTE: This method does not use the use_tenancy parameter and if this method is being used
// to control tenancy the parameter must be checked. The method IsResourceAuthorizedToUser checks the use_tenancy parameter
// and should be used for this purpose in most cases.
func GetUserTenantList(user auth.CurrentUser, db *sqlx.DB) ([]tc.TenantNullable, error) {
	query := `WITH RECURSIVE q AS (SELECT id, name, active, parent_id, last_updated FROM tenant WHERE id = $1
	UNION SELECT t.id, t.name, t.active, t.parent_id, t.last_updated  FROM tenant t JOIN q ON q.id = t.parent_id)
	SELECT id, name, active, parent_id, last_updated FROM q;`

	log.Debugln("\nQuery: ", query)

	var (
		tenantID, parentID int
		name               string
		active             bool
		lastUpdated        tc.TimeNoMod
	)
	rows, err := db.Query(query, user.TenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []tc.TenantNullable{}

	for rows.Next() {
		if err := rows.Scan(&tenantID, &name, &active, &parentID, &lastUpdated); err != nil {
			return nil, err
		}

		tenants = append(tenants, tc.TenantNullable{ID: &tenantID, Name: &name, Active: &active, ParentID: &parentID})
	}

	return tenants, nil
}

// returns a TenantID list that the specified user has access too.
// NOTE: This method does not use the use_tenancy parameter and if this method is being used
// to control tenancy the parameter must be checked. The method IsResourceAuthorizedToUser checks the use_tenancy parameter
// and should be used for this purpose in most cases.
func GetUserTenantIDList(user auth.CurrentUser, db *sqlx.DB) ([]int, error) {
	query := `WITH RECURSIVE q AS (SELECT id, name, active, parent_id FROM tenant WHERE id = $1
	UNION SELECT t.id, t.name, t.active, t.parent_id  FROM tenant t JOIN q ON q.id = t.parent_id)
	SELECT id FROM q;`

	log.Debugln("\nQuery: ", query)

	var tenantID int

	rows, err := db.Query(query, user.TenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []int{}

	for rows.Next() {
		if err := rows.Scan(&tenantID); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenantID)
	}

	return tenants, nil
}

func GetUserTenantIDListTx(user *auth.CurrentUser, tx *sqlx.Tx) ([]int, error) {
	query := `WITH RECURSIVE q AS (SELECT id, name, active, parent_id FROM tenant WHERE id = $1
	UNION SELECT t.id, t.name, t.active, t.parent_id  FROM tenant t JOIN q ON q.id = t.parent_id)
	SELECT id FROM q;`

	log.Debugln("\nQuery: ", query)

	var tenantID int

	rows, err := tx.Query(query, user.TenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []int{}

	for rows.Next() {
		if err := rows.Scan(&tenantID); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenantID)
	}

	return tenants, nil
}

// IsTenancyEnabled returns true if tenancy is enabled or false otherwise
func IsTenancyEnabled(db *sqlx.DB) bool {
	query := `SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY`
	var useTenancy bool
	err := db.QueryRow(query).Scan(&useTenancy)
	if err != nil {
		log.Errorf("Error checking if tenancy is enabled: %v", err)
		return false
	}
	return useTenancy
}

func IsTenancyEnabledTx(tx *sql.Tx) (bool, error) {
	query := `SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY`
	useTenancy := false
	if err := tx.QueryRow(query).Scan(&useTenancy); err != nil {
		return false, errors.New("checking if tenancy is enabled: " + err.Error())
	}
	return useTenancy, nil
}

// returns a boolean value describing if the user has access to the provided resource tenant id and an error
// if use_tenancy is set to false (0 in the db) this method will return true allowing access.
func IsResourceAuthorizedToUser(resourceTenantID int, user *auth.CurrentUser, db *sqlx.DB) (bool, error) {
	// $1 is the user tenant ID and $2 is the resource tenant ID
	query := `WITH RECURSIVE q AS (SELECT id, active FROM tenant WHERE id = $1
	UNION SELECT t.id, t.active FROM TENANT t JOIN q ON q.id = t.parent_id),
	tenancy AS (SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY)
	SELECT id, active, tenancy.value AS use_tenancy FROM tenancy, q WHERE id = $2 UNION ALL SELECT -1, false, tenancy.value AS use_tenancy FROM tenancy FETCH FIRST 1 ROW ONLY;`

	var tenantID int
	var active bool
	var useTenancy bool

	log.Debugln("\nQuery: ", query)
	err := db.QueryRow(query, user.TenantID, resourceTenantID).Scan(&tenantID, &active, &useTenancy)

	switch {
	case err != nil:
		log.Errorf("Error checking user tenant %v access on resourceTenant  %v: %v", user.TenantID, resourceTenantID, err.Error())
		return false, err
	default:
		if !useTenancy {
			return true, nil
		}
		if active && tenantID == resourceTenantID {
			return true, nil
		}
		return false, nil
	}
}

// returns a boolean value describing if the user has access to the provided resource tenant id and an error
// if use_tenancy is set to false (0 in the db) this method will return true allowing access.
func IsResourceAuthorizedToUserTx(resourceTenantID int, user *auth.CurrentUser, tx *sql.Tx) (bool, error) {
	// $1 is the user tenant ID and $2 is the resource tenant ID
	query := `WITH RECURSIVE q AS (SELECT id, active FROM tenant WHERE id = $1
	UNION SELECT t.id, t.active FROM TENANT t JOIN q ON q.id = t.parent_id),
	tenancy AS (SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY)
	SELECT id, active, tenancy.value AS use_tenancy FROM tenancy, q WHERE id = $2 UNION ALL SELECT -1, false, tenancy.value AS use_tenancy FROM tenancy FETCH FIRST 1 ROW ONLY;`

	var tenantID int
	var active bool
	var useTenancy bool

	log.Debugln("\nQuery: ", query)
	err := tx.QueryRow(query, user.TenantID, resourceTenantID).Scan(&tenantID, &active, &useTenancy)

	switch {
	case err != nil:
		log.Errorf("Error checking user tenant %v access on resourceTenant  %v: %v", user.TenantID, resourceTenantID, err.Error())
		return false, err
	default:
		if !useTenancy {
			return true, nil
		}
		if active && tenantID == resourceTenantID {
			return true, nil
		} else {
			fmt.Printf("default")
			return false, nil
		}
	}
}

// TOTenant provides a local type against which to define methods
type TOTenant struct {
	ReqInfo *api.APIInfo `json:"-"`
	tc.TenantNullable
}

func GetTypeSingleton() func(reqInfo *api.APIInfo) api.CRUDer {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOTenant{reqInfo, tc.TenantNullable{}}
		return &toReturn
	}
}

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
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
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

// Validate fulfills the api.Validator interface
func (ten TOTenant) Validate() []error {
	errs := validation.Errors{
		"name":     validation.Validate(ten.Name, validation.Required),
		"active":   validation.Validate(ten.Active), // only validate it's boolean
		"parentId": validation.Validate(ten.ParentID, validation.Required, validation.Min(1)),
	}
	return tovalidate.ToErrors(errs)
}

// Create implements the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a tenant with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted tenant and have
//to be added to the struct
func (ten *TOTenant) Create() (error, tc.ApiErrorType) {
	resultRows, err := ten.ReqInfo.Tx.NamedQuery(insertQuery(), ten)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a tenant with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received non pq error: %++v from create execution", err)
		return tc.DBError, tc.SystemError
	}
	defer resultRows.Close()

	var id int
	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no tenant was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from tenant insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	ten.SetKeys(map[string]interface{}{"id": id})
	ten.LastUpdated = &lastUpdated

	return nil, tc.NoError
}

// Read implements the tc.Reader interface
func (ten *TOTenant) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"active":      dbhelpers.WhereColumnInfo{Column: "t.active", Checker: nil},
		"id":          dbhelpers.WhereColumnInfo{Column: "t.id", Checker: api.IsInt},
		"name":        dbhelpers.WhereColumnInfo{Column: "t.name", Checker: nil},
		"parent_id":   dbhelpers.WhereColumnInfo{Column: "t.parentID", Checker: api.IsInt},
		"parent_name": dbhelpers.WhereColumnInfo{Column: "p.name", Checker: api.IsInt},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := ten.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying tenants: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	tenants := []interface{}{}
	for rows.Next() {
		var s TOTenant
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing Tenant rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		tenants = append(tenants, s)
	}

	return tenants, []error{}, tc.NoError
}

//The TOTenant implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a tenant with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (ten *TOTenant) Update() (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with tenant: %++v", updateQuery(), ten)
	resultRows, err := ten.ReqInfo.Tx.NamedQuery(updateQuery(), ten)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a tenant with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received error: %++v from update execution", err)
		return tc.DBError, tc.SystemError
	}
	defer resultRows.Close()

	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			log.Error.Printf("could not scan lastUpdated from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	log.Debugf("lastUpdated: %++v", lastUpdated)
	ten.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no tenant found with this id"), tc.DataMissingError
		}
		return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
	}
	return nil, tc.NoError
}

//Delete implements the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (ten *TOTenant) Delete() (error, tc.ApiErrorType) {
	if ten.ID == nil {
		// should never happen...
		return errors.New("invalid tenant: id is nil"), tc.SystemError
	}

	log.Debugf("about to run exec query: %s with tenant: %++v", deleteQuery(), ten)
	result, err := ten.ReqInfo.Tx.NamedExec(deleteQuery(), ten)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err = fmt.Errorf("pqErr is %++v\n", pqErr)
			var existing string
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

			// another query to get tenant name for the error message
			name := strconv.Itoa(*ten.ID)
			if err := ten.ReqInfo.Tx.QueryRow(`SELECT name FROM tenant WHERE id = $1`, *ten.ID).Scan(&name); err != nil {
				// use ID as a backup for name the error -- this should never happen
				log.Debugf("error getting tenant name: %++v", err)
			}

			err = errors.New("Tenant '" + name + "' has " + existing + ". Please update these " + existing + " and retry.")
			return err, tc.DataConflictError
		}
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no tenant with that id found"), tc.DataMissingError
		}
		return fmt.Errorf("this delete affected too many rows: %d", rowsAffected), tc.SystemError
	}

	return nil, tc.NoError
}

func selectQuery() string {
	query := `SELECT
t.active AS active,
t.name AS name,
t.id AS id,
t.last_updated AS last_updated,
t.parent_id AS parent_id,
p.name AS parent_name

FROM tenant AS t
LEFT OUTER JOIN tenant AS p
ON t.parent_id = p.id`
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

// getDSTenantIDByIDTx returns the tenant ID, whether the delivery service exists, and any error.
// Note the id may be nil, even if true is returned, if the delivery service exists but its tenant_id field is null.
// TODO move somewhere generic
func getDSTenantIDByIDTx(tx *sql.Tx, id int) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where id = $1`, id).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service ID '%v': %v", id, err)
	}
	return tenantID, true, nil
}
