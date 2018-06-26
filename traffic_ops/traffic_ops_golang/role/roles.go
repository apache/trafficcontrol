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
	"errors"
	"fmt"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TORole struct{
	ReqInfo *api.APIInfo `json:"-"`
	v13.Role
}

func GetTypeSingleton() func(reqInfo *api.APIInfo)api.CRUDer {
	return func(reqInfo *api.APIInfo)api.CRUDer {
		toReturn := TORole{reqInfo, v13.Role{}}
		return &toReturn
	}
}

func (role TORole) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
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
func (role TORole) Validate() []error {
	errs := validation.Errors{
		"name":        validation.Validate(role.Name, validation.Required),
		"description": validation.Validate(role.Description, validation.Required),
		"privLevel":   validation.Validate(role.PrivLevel, validation.Required)}

	errsToReturn := tovalidate.ToErrors(errs)
	checkCaps := `SELECT cap FROM UNNEST($1::text[]) AS cap WHERE NOT cap =  ANY(ARRAY(SELECT c.name FROM capability AS c WHERE c.name = ANY($1)))`
	var badCaps []string
	if role.ReqInfo.Tx != nil {
		err := role.ReqInfo.Tx.Select(&badCaps, checkCaps, pq.Array(role.Capabilities))
		if err != nil {
			log.Errorf("got error from selecting bad capabilities: %v", err)
			return []error{tc.DBError}
		}
		if len(badCaps) > 0 {
			errsToReturn = append(errsToReturn, fmt.Errorf("can not add non-existent capabilities: %v", badCaps))
		}
	}
	return errsToReturn
}

//The TORole implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a role with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted role and have
//to be added to the struct
func (role *TORole) Create() (error, tc.ApiErrorType) {
	if *role.PrivLevel > role.ReqInfo.User.PrivLevel {
		return errors.New("can not create a role with a higher priv level than your own"), tc.ForbiddenError
	}
	resultRows, err := role.ReqInfo.Tx.NamedQuery(insertQuery(), role)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a role with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	defer resultRows.Close()

	var id int
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id); err != nil {
			log.Error.Printf("could not scan id from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no role was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from role insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	role.SetKeys(map[string]interface{}{"id": id})
	//after we have role ID we can associate the capabilities:
	err, errType := role.createRoleCapabilityAssociations(role.ReqInfo.Tx)
	if err != nil {
		return err, errType
	}

	return nil, tc.NoError
}

func (role *TORole) createRoleCapabilityAssociations(tx *sqlx.Tx) (error, tc.ApiErrorType) {
	result, err := tx.Exec(associateCapabilities(), role.ID, pq.Array(role.Capabilities))
	if err != nil {
		log.Errorf("received non pq error: %++v from create execution", err)
		return tc.DBError, tc.SystemError
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	}
	expected := len(*role.Capabilities)
	if int(rows) != expected {
		log.Errorf("wrong number of role_capability rows created: %d expected: %d", rows, expected)
	}
	return nil, tc.NoError
}

func (role *TORole) deleteRoleCapabilityAssociations(tx *sqlx.Tx) (error, tc.ApiErrorType) {
	result, err := tx.Exec(deleteAssociatedCapabilities(), role.ID)
	if err != nil {

		log.Errorf("received error: %++v from create execution", err)
		return tc.DBError, tc.SystemError

	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Errorf("could not check result after inserting role_capability relations: %v", err)
	}
	return nil, tc.NoError
}

func (role *TORole) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"name": dbhelpers.WhereColumnInfo{"name", nil},
		"id":   dbhelpers.WhereColumnInfo{"id", api.IsInt},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := role.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying Roles: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	Roles := []interface{}{}
	for rows.Next() {
		var r TORole
		var caps []string
		if err = rows.Scan(&r.ID, &r.Name, &r.Description, &r.PrivLevel, pq.Array(&caps)); err != nil {
			log.Errorf("error parsing Role rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		r.Capabilities = &caps
		Roles = append(Roles, r)
	}

	return Roles, []error{}, tc.NoError
}

//The TORole implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a role with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (role *TORole) Update() (error, tc.ApiErrorType) {
	if *role.PrivLevel > role.ReqInfo.User.PrivLevel {
		return errors.New("can not create a role with a higher priv level than your own"), tc.ForbiddenError
	}

	log.Debugf("about to run exec query: %s with role: %++v\n", updateQuery(), role)
	result, err := role.ReqInfo.Tx.NamedExec(updateQuery(), role)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a role with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return tc.DBError, tc.SystemError
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("received error: %++v from checking result of update", err)
		return tc.DBError, tc.SystemError
	}

	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no role found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	//remove associations
	err, errType := role.deleteRoleCapabilityAssociations(role.ReqInfo.Tx)
	if err != nil {
		return err, errType
	}
	//create new associations
	err, errType = role.createRoleCapabilityAssociations(role.ReqInfo.Tx)
	if err != nil {
		return err, errType
	}

	return nil, tc.NoError
}

//The Role implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (role *TORole) Delete() (error, tc.ApiErrorType) {
	assignedUsers := 0
	err := role.ReqInfo.Tx.Get(&assignedUsers, "SELECT COUNT(id) FROM tm_user WHERE role=$1", role.ID)
	if err != nil {
		log.Errorf("received error: %++v from assigned users check", err)
		return tc.DBError, tc.SystemError
	}
	if assignedUsers != 0 {
		return fmt.Errorf("can not delete a role with %d assigned users", assignedUsers), tc.DataConflictError
	}

	log.Debugf("about to run exec query: %s with role: %++v", deleteQuery(), role)
	result, err := role.ReqInfo.Tx.NamedExec(deleteQuery(), role)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no role with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	//remove associations
	err, errType := role.deleteRoleCapabilityAssociations(role.ReqInfo.Tx)
	if err != nil {
		return err, errType
	}

	return nil, tc.NoError
}

func selectQuery() string {
	query := `SELECT
id,
name,
description,
priv_level,
ARRAY(SELECT rc.cap_name FROM role_capability AS rc WHERE rc.role_id=id) AS capabilities

FROM role`
	return query
}

func updateQuery() string {
	query := `UPDATE
role SET
name=:name,
description=:description
WHERE id=:id`
	return query
}

func deleteAssociatedCapabilities() string {
	query := `DELETE FROM role_capability
WHERE role_id=$1`
	return query
}

func associateCapabilities() string {
	query := `INSERT INTO role_capability (
role_id,
cap_name) WITH
	q1 AS ( SELECT * FROM (VALUES ($1::bigint)) AS role_id ),
	q2 AS (SELECT UNNEST($2::text[]))
	SELECT * FROM q1,q2`
	return query
}

func insertQuery() string {
	query := `INSERT INTO role (
name,
description,
priv_level) VALUES (
:name,
:description,
:priv_level) RETURNING id`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM role
WHERE id=:id`
	return query
}
