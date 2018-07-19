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

// tenant.go defines the TOTenant object and methods/functions required for the api/.../tenants endpoints

import (
	"errors"
	"fmt"
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

// TOTenant provides a local type against which to define methods
type TOTenant struct {
	ReqInfo *api.APIInfo `json:"-"`
	tc.TenantNullable
}

func GetTypeSingleton() api.CRUDFactory {
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

// Validate fulfills the api.Validator interface
func (ten TOTenant) Validate() error {
	errs := validation.Errors{
		"name":       validation.Validate(ten.Name, validation.Required),
		"active":     validation.Validate(ten.Active), // only validate it's boolean
		"parentId":   validation.Validate(ten.ParentID, validation.Required, validation.Min(1)),
		"parentName": nil,
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
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

	tenantID := ten.ReqInfo.User.TenantID
	if tenantID == auth.TenantIDInvalid {
		// NOTE: work around issue where user has no tenancy assigned.  If tenancy turned off, there
		// should be no tenancy restrictions.  This should be removed once tenant_id NOT NULL constraints
		// are in place
		enabled, err := IsTenancyEnabledTx(ten.ReqInfo.Tx.Tx)
		if err != nil {
			log.Infof("error checking tenancy: %v", err)
			return nil, nil, tc.SystemError
		}

		if enabled {
			// tenancy enabled, but user doesn't belong to one -- return empty list
			return nil, nil, tc.NoError
		}
		// give it the root tenant -- since tenancy turned off,  does not matter what it is
		tenantID = 1
	}

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"active":      dbhelpers.WhereColumnInfo{Column: "q.active", Checker: nil},
		"id":          dbhelpers.WhereColumnInfo{Column: "q.id", Checker: api.IsInt},
		"name":        dbhelpers.WhereColumnInfo{Column: "q.name", Checker: nil},
		"parent_id":   dbhelpers.WhereColumnInfo{Column: "q.parent_id", Checker: api.IsInt},
		"parent_name": dbhelpers.WhereColumnInfo{Column: "p.name", Checker: nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery(tenantID) + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := ten.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying tenants: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	tenants := []interface{}{}
	tenantNames := make(map[int]*string)
	for rows.Next() {
		var t TOTenant
		if err = rows.StructScan(&t); err != nil {
			log.Errorf("error parsing Tenant rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		if t.ID == nil || t.Name == nil {
			log.Errorf("tenant with no id and/or name: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}

		tenantNames[*t.ID] = t.Name
		tenants = append(tenants, t)
	}
	// fill in parent names
	for _, i := range tenants {
		t := i.(TOTenant)
		if t.ParentID == nil || tenantNames[*t.ParentID] == nil {
			// root tenant has no parent
			continue
		}
		p := *tenantNames[*t.ParentID]
		t.ParentName = &p
	}

	return tenants, []error{}, tc.NoError
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
		ok, err = IsResourceAuthorizedToUserTx(*ten.ID, user, ten.ReqInfo.Tx.Tx)
		if !ok || err != nil {
			return ok, err
		}

		if ten.ParentID == nil || *ten.ParentID == 0 {
			// not changing parent
			return ok, err
		}

		// get current parentID to check if it's being changed
		var parentID int
		tx := ten.ReqInfo.Tx.Tx
		err = tx.QueryRow(`SELECT parent_id FROM tenant WHERE id = ` + strconv.Itoa(*ten.ID)).Scan(&parentID)
		if err != nil {
			return false, err
		}
		if parentID == *ten.ParentID {
			// parent not being changed
			return ok, err
		}
	}

	// creating new tenant -- must specify a parent
	if ten.ParentID == nil || *ten.ParentID == 0 {
		return false, err
	}

	// check if authorized on new parent tenant
	return IsResourceAuthorizedToUserTx(*ten.ParentID, user, ten.ReqInfo.Tx.Tx)
}

//Update implements the Updater interface.
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
			log.Error.Printf("could not scan lastUpdated from update: %s\n", err)
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
