package parameter

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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	NameQueryParam       = "name"
	SecureQueryParam     = "secure"
	ConfigFileQueryParam = "configFile"
	IDQueryParam         = "id"
	ValueQueryParam      = "value"
)

var (
	HiddenField = "********"
)

//we need a type alias to define functions on
type TOParameter tc.ParameterNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOParameter(tc.ParameterNullable{})

func GetRefType() *TOParameter {
	return &refType
}

func (parameter TOParameter) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{IDQueryParam, api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (parameter TOParameter) GetKeys() (map[string]interface{}, bool) {
	if parameter.ID == nil {
		return map[string]interface{}{IDQueryParam: 0}, false
	}
	return map[string]interface{}{IDQueryParam: *parameter.ID}, true
}

func (parameter *TOParameter) SetKeys(keys map[string]interface{}) {
	i, _ := keys[IDQueryParam].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	parameter.ID = &i
}

func (parameter *TOParameter) GetAuditName() string {
	if parameter.Name != nil {
		return *parameter.Name
	}
	if parameter.ID != nil {
		return strconv.Itoa(*parameter.ID)
	}
	return "unknown"
}

func (parameter *TOParameter) GetType() string {
	return "parameter"
}

// Validate fulfills the api.Validator interface
func (parameter TOParameter) Validate(db *sqlx.DB) []error {

	// Test
	// - Secure Flag is always set to either 1/0
	// - Admin rights only
	// - Do not allow duplicate parameters by name+config_file+value
	errs := validation.Errors{
		NameQueryParam:       validation.Validate(parameter.Name, validation.Required),
		ConfigFileQueryParam: validation.Validate(parameter.ConfigFile, validation.Required),
		ValueQueryParam:      validation.Validate(parameter.Value, validation.Required),
	}

	return tovalidate.ToErrors(errs)
}

//The TOParameter implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a parameter with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted parameter and have
//to be added to the struct
func (pl *TOParameter) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	resultRows, err := tx.NamedQuery(insertQuery(), pl)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a parameter with " + err.Error()), eType
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
		err = errors.New("no parameter was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	if rowsAffected > 1 {
		err = errors.New("too many ids returned from parameter insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}

	pl.SetKeys(map[string]interface{}{IDQueryParam: id})
	pl.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func insertQuery() string {
	query := `INSERT INTO parameter (
name,
config_file,
value,
secure) VALUES (
:name,
:config_file,
:value,
:secure) RETURNING id,last_updated`
	return query
}

func (parameter *TOParameter) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	privLevel := user.PrivLevel

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		ConfigFileQueryParam: dbhelpers.WhereColumnInfo{"p.config_file", nil},
		IDQueryParam:         dbhelpers.WhereColumnInfo{"p.id", api.IsInt},
		NameQueryParam:       dbhelpers.WhereColumnInfo{"p.name", nil},
		SecureQueryParam:     dbhelpers.WhereColumnInfo{"p.secure", api.IsBool},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + ParametersGroupBy() + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying Parameters: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	params := []interface{}{}
	for rows.Next() {
		var p tc.ParameterNullable
		if err = rows.StructScan(&p); err != nil {
			log.Errorf("error parsing Parameter rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		var isSecure bool
		if p.Secure != nil {
			isSecure = *p.Secure
		}

		if isSecure && (privLevel < auth.PrivLevelAdmin) {
			p.Value = &HiddenField
		}
		params = append(params, p)
	}

	return params, []error{}, tc.NoError

}

//The TOParameter implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a parameter with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (pl *TOParameter) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with parameter: %++v", updateQuery(), pl)
	resultRows, err := tx.NamedQuery(updateQuery(), pl)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a parameter with " + err.Error()), eType
			}
			return err, eType
		}
		log.Errorf("received error: %++v from update execution", err)
		return tc.DBError, tc.SystemError
	}
	defer resultRows.Close()

	// get LastUpdated field -- updated by trigger in the db
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
	pl.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no parameter found with this id"), tc.DataMissingError
		}
		return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

//The Parameter implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (pl *TOParameter) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	rollbackTransaction := true
	tx, err := db.Beginx()
	defer func() {
		if tx == nil || !rollbackTransaction {
			return
		}
		err := tx.Rollback()
		if err != nil {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v", err)
		return tc.DBError, tc.SystemError
	}
	log.Debugf("about to run exec query: %s with parameter: %++v", deleteQuery(), pl)
	result, err := tx.NamedExec(deleteQuery(), pl)
	if err != nil {
		log.Errorf("received error: %++v from delete execution", err)
		return tc.DBError, tc.SystemError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tc.DBError, tc.SystemError
	}
	if rowsAffected < 1 {
		return errors.New("no parameter with that id found"), tc.DataMissingError
	}
	if rowsAffected > 1 {
		return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
	}

	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func selectQuery() string {

	query := `SELECT
p.config_file,
p.id,
p.last_updated,
p.name,
p.value,
p.secure,
COALESCE(array_to_json(array_agg(pr.name) FILTER (WHERE pr.name IS NOT NULL)), '[]') AS profiles
FROM parameter p
LEFT JOIN profile_parameter pp ON p.id = pp.parameter
LEFT JOIN profile pr ON pp.profile = pr.id`
	return query
}

func updateQuery() string {
	query := `UPDATE
parameter SET
config_file=:config_file,
id=:id,
name=:name,
value=:value,
secure=:secure
WHERE id=:id RETURNING last_updated`
	return query
}

// ParametersGroupBy ...
func ParametersGroupBy() string {
	groupBy := ` GROUP BY p.config_file, p.id, p.last_updated, p.name, p.value, p.secure`
	return groupBy
}

func deleteQuery() string {
	query := `DELETE FROM parameter
WHERE id=:id`
	return query
}
