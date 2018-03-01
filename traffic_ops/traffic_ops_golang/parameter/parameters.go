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
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOParameter tc.ParameterNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOParameter(tc.ParameterNullable{})

func GetRefType() *TOParameter {
	return &refType
}

//Implementation of the Identifier, Validator interface functions
func (parameter *TOParameter) GetID() (int, bool) {
	if parameter.ID == nil {
		return 0, false
	}
	return *parameter.ID, true
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

func (parameter *TOParameter) SetID(i int) {
	parameter.ID = &i
}

func (pl *TOParameter) Validate(db *sqlx.DB) []error {
	errs := []error{}
	name := pl.Name
	if name != nil && len(*name) < 1 {
		errs = append(errs, errors.New(`Parameter 'name' is required.`))
	}
	return errs
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

	pl.SetID(id)
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

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"config_file":  dbhelpers.WhereColumnInfo{"p.config_file", nil},
		"id":           dbhelpers.WhereColumnInfo{"p.id", api.IsInt},
		"last_updated": dbhelpers.WhereColumnInfo{"p.last_updated", nil},
		"name":         dbhelpers.WhereColumnInfo{"p.name", nil},
		"secure":       dbhelpers.WhereColumnInfo{"p.secure", api.IsBool},
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
		var s tc.ParameterNullable
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing Parameter rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		params = append(params, s)
	}

	return params, []error{}, tc.NoError

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

// ParametersGroupBy ...
func ParametersGroupBy() string {
	groupBy := ` GROUP BY p.config_file, p.id, p.last_updated, p.name, p.value, p.secure`
	return groupBy
}
