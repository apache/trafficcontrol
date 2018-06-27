package division

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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TODivision struct {
	ReqInfo *api.APIInfo `json:"-"`
	tc.DivisionNullable
}

func GetTypeSingleton() func(reqInfo *api.APIInfo) api.CRUDer {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TODivision{reqInfo, tc.DivisionNullable{}}
		return &toReturn
	}
}

func (division TODivision) GetAuditName() string {
	if division.Name != nil {
		return *division.Name
	}
	if division.ID != nil {
		return strconv.Itoa(*division.ID)
	}
	return "unknown"
}

func (division TODivision) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (division TODivision) GetKeys() (map[string]interface{}, bool) {
	if division.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *division.ID}, true
}

func (division *TODivision) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	division.ID = &i
}

func (division TODivision) GetType() string {
	return "division"
}

func (division TODivision) Validate() []error {
	errs := validation.Errors{
		"name": validation.Validate(division.Name, validation.NotNil, validation.Required),
	}
	return tovalidate.ToErrors(errs)
}

//The TODivision implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a division with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted division and have
//to be added to the struct
func (division *TODivision) Create() (error, tc.ApiErrorType) {
	resultRows, err := division.ReqInfo.Tx.NamedQuery(insertQuery(), division)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a division with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received non pq error: %++v from create execution", err)
			return tc.DBError, tc.SystemError
		}
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
		err = errors.New("no division was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from division insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	division.SetKeys(map[string]interface{}{"id": id})
	division.LastUpdated = &lastUpdated
	return nil, tc.NoError
}

func (division *TODivision) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	if strings.HasSuffix(parameters["name"], ".json") {
		parameters["name"] = parameters["name"][:len(parameters["name"])-len(".json")]
	}
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"id":   dbhelpers.WhereColumnInfo{"id", api.IsInt},
		"name": dbhelpers.WhereColumnInfo{"name", nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := division.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying Divisions: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	divisions := []interface{}{}
	for rows.Next() {
		var s tc.Division
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing Division rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		divisions = append(divisions, s)
	}

	return divisions, []error{}, tc.NoError
}

//The TODivision implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a division with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (division *TODivision) Update() (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with division: %++v", updateQuery(), division)
	resultRows, err := division.ReqInfo.Tx.NamedQuery(updateQuery(), division)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a division with " + err.Error()), eType
			}
			return err, eType
		} else {
			log.Errorf("received error: %++v from update execution", err)
			return tc.DBError, tc.SystemError
		}
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
	division.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no division found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//The Division implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (division *TODivision) Delete() (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with division: %++v", deleteQuery(), division)
	result, err := division.ReqInfo.Tx.NamedExec(deleteQuery(), division)
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
			return errors.New("no division with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

func insertQuery() string {
	query := `INSERT INTO division (
name) VALUES (:name) RETURNING id,last_updated`
	return query
}

func selectQuery() string {

	query := `SELECT
id,
last_updated,
name 

FROM division d`
	return query
}

func updateQuery() string {
	query := `UPDATE
division SET
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM division
WHERE id=:id`
	return query
}
