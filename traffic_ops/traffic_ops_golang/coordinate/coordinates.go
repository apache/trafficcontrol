package coordinate

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
	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOCoordinate struct {
	ReqInfo *api.APIInfo `json:"-"`
	v13.CoordinateNullable
}

func GetTypeSingleton() api.CRUDFactory {
	return func(reqInfo *api.APIInfo) api.CRUDer {
		toReturn := TOCoordinate{reqInfo, v13.CoordinateNullable{}}
		return &toReturn
	}
}

func (coordinate TOCoordinate) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (coordinate TOCoordinate) GetKeys() (map[string]interface{}, bool) {
	if coordinate.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *coordinate.ID}, true
}

func (coordinate TOCoordinate) GetAuditName() string {
	if coordinate.Name != nil {
		return *coordinate.Name
	}
	if coordinate.ID != nil {
		return strconv.Itoa(*coordinate.ID)
	}
	return "0"
}

func (coordinate TOCoordinate) GetType() string {
	return "coordinate"
}

func (coordinate *TOCoordinate) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	coordinate.ID = &i
}

func isValidCoordinateChar(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	if r == '.' || r == '-' || r == '_' {
		return true
	}
	return false
}

// IsValidCoordinateName returns true if the name contains only characters valid for a Coordinate name
func IsValidCoordinateName(str string) bool {
	i := strings.IndexFunc(str, func(r rune) bool { return !isValidCoordinateChar(r) })
	return i == -1
}

// Validate fulfills the api.Validator interface
func (coordinate TOCoordinate) Validate() []error {
	validName := validation.NewStringRule(IsValidCoordinateName, "invalid characters found - Use alphanumeric . or - or _ .")
	latitudeErr := "Must be a floating point number within the range +-90"
	longitudeErr := "Must be a floating point number within the range +-180"
	errs := validation.Errors{
		"name":      validation.Validate(coordinate.Name, validation.Required, validName),
		"latitude":  validation.Validate(coordinate.Latitude, validation.Min(-90.0).Error(latitudeErr), validation.Max(90.0).Error(latitudeErr)),
		"longitude": validation.Validate(coordinate.Longitude, validation.Min(-180.0).Error(longitudeErr), validation.Max(180.0).Error(longitudeErr)),
	}
	return tovalidate.ToErrors(errs)
}

//The TOCoordinate implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a coordinate with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted coordinate and have
//to be added to the struct
func (coordinate *TOCoordinate) Create() (error, tc.ApiErrorType) {
	resultRows, err := coordinate.ReqInfo.Tx.NamedQuery(insertQuery(), coordinate)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a coordinate with " + err.Error()), eType
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
		err = errors.New("no coordinate was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from coordinate insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	coordinate.SetKeys(map[string]interface{}{"id": id})
	coordinate.LastUpdated = &lastUpdated
	return nil, tc.NoError
}

func (coordinate *TOCoordinate) Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

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

	rows, err := coordinate.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying Coordinate: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	Coordinates := []interface{}{}
	for rows.Next() {
		var s TOCoordinate
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing Coordinate rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		Coordinates = append(Coordinates, s)
	}

	return Coordinates, []error{}, tc.NoError
}

//The TOCoordinate implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a coordinate with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (coordinate *TOCoordinate) Update() (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with coordinate: %++v", updateQuery(), coordinate)
	resultRows, err := coordinate.ReqInfo.Tx.NamedQuery(updateQuery(), coordinate)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a coordinate with " + err.Error()), eType
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
	coordinate.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no coordinate found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

//The Coordinate implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (coordinate *TOCoordinate) Delete() (error, tc.ApiErrorType) {
	log.Debugf("about to run exec query: %s with coordinate: %++v", deleteQuery(), coordinate)
	result, err := coordinate.ReqInfo.Tx.NamedExec(deleteQuery(), coordinate)
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
			return errors.New("no coordinate with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this delete affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	return nil, tc.NoError
}

func selectQuery() string {
	query := `SELECT
id,
latitude,
longitude,
last_updated,
name

FROM coordinate c`
	return query
}

func updateQuery() string {
	query := `UPDATE
coordinate SET
latitude=:latitude,
longitude=:longitude,
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO coordinate (
latitude,
longitude,
name) VALUES (
:latitude,
:longitude,
:name) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM coordinate
WHERE id=:id`
	return query
}
