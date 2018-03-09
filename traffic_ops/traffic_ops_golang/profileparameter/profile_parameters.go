package profileparameter

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

//we need a type alias to define functions on
type TOProfileParameter tc.ProfileParameterNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOProfileParameter(tc.ProfileParameterNullable{})

func GetRefType() *TOProfileParameter {
	return &refType
}

//Implementation of the Identifier, Validator interface functions
func (pp *TOProfileParameter) GetID() (int, bool) {
	if pp.Profile == nil {
		return 0, false
	}
	return *pp.Profile, true
}

func (pp *TOProfileParameter) GetAuditName() string {
	if pp.Profile != nil {
		return strconv.Itoa(*pp.Profile)
	}
	return "unknown"
}

func (pp *TOProfileParameter) GetType() string {
	return "profileParameter"
}

func (pp *TOProfileParameter) SetID(i int) {
	pp.Profile = &i
}

// Validate fulfills the api.Validator interface
func (pp TOProfileParameter) Validate(db *sqlx.DB) []error {

	errs := validation.Errors{
		"profile":   validation.Validate(pp.Profile, validation.Required),
		"parameter": validation.Validate(pp.Parameter, validation.Required),
	}

	return tovalidate.ToErrors(errs)
}

//The TOProfileParameter implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a profileparameter with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the profile and lastUpdated values of the newly inserted profileparameter and have
//to be added to the struct
func (pp *TOProfileParameter) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	resultRows, err := tx.NamedQuery(insertQuery(), pp)
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

	var profile int
	var lastUpdated tc.TimeNoMod
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&profile, &lastUpdated); err != nil {
			log.Error.Printf("could not scan profile from insert: %s\n", err)
			return tc.DBError, tc.SystemError
		}
	}
	if rowsAffected == 0 {
		err = errors.New("no parameter was inserted, no profile was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	if rowsAffected > 1 {
		err = errors.New("too many ids returned from parameter insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}

	pp.SetID(profile)
	pp.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func insertQuery() string {
	query := `INSERT INTO profile_parameter (
profile,
parameter) VALUES (
:profile,
:parameter) RETURNING profile_profile,last_updated`
	return query
}

func (pp *TOProfileParameter) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	privLevel := user.PrivLevel

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"profile":      dbhelpers.WhereColumnInfo{"pp.profile", api.IsInt},
		"last_updated": dbhelpers.WhereColumnInfo{"pp.last_updated", nil},
		"name":         dbhelpers.WhereColumnInfo{"p.parameter", nil},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying Parameters: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	params := []interface{}{}
	hiddenField := "********"
	for rows.Next() {
		var p tc.ParameterNullable
		if err = rows.StructScan(&p); err != nil {
			log.Errorf("error parsing pp rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		var isSecure bool
		if p.Secure != nil {
			isSecure = *p.Secure
		}

		if isSecure && (privLevel < auth.PrivLevelAdmin) {
			p.Value = &hiddenField
		}
		params = append(params, p)
	}

	return params, []error{}, tc.NoError

}

//The TOProfileParameter implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a parameter with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (pp *TOProfileParameter) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	log.Debugf("about to run exec query: %s with parameter: %++v", updateQuery(), pp)
	resultRows, err := tx.NamedQuery(updateQuery(), pp)
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
	pp.LastUpdated = &lastUpdated
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
func (pp *TOProfileParameter) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	log.Debugf("about to run exec query: %s with parameter: %++v", deleteQuery(), pp)
	result, err := tx.NamedExec(deleteQuery(), pp)
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
p.last_updated,
p.profile,
p.id,
p.parameter
FROM profile_parameter p`
	return query
}

func updateQuery() string {
	query := `UPDATE
profile_parameter SET
profile=:profile,
parameter=:parameter
WHERE id=:id RETURNING last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM profile_parameter
WHERE id=:id`
	return query
}
