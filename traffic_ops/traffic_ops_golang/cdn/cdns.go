package cdn

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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//we need a type alias to define functions on
type TOCDN v13.CDNNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOCDN{}

func GetRefType() *TOCDN {
	return &refType
}

func (cdn TOCDN) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (cdn TOCDN) GetKeys() (map[string]interface{}, bool) {
	if cdn.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *cdn.ID}, true
}

func (cdn TOCDN) GetAuditName() string {
	if cdn.Name != nil {
		return *cdn.Name
	}
	if cdn.ID != nil {
		return strconv.Itoa(*cdn.ID)
	}
	return "0"
}

func (cdn TOCDN) GetType() string {
	return "cdn"
}

func (cdn *TOCDN) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	cdn.ID = &i
}

func isValidCDNchar(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	if r == '.' || r == '-' {
		return true
	}
	return false
}

// IsValidCDNName returns true if the name contains only characters valid for a CDN name
func IsValidCDNName(str string) bool {
	i := strings.IndexFunc(str, func(r rune) bool { return !isValidCDNchar(r) })
	return i == -1
}

// Validate fulfills the api.Validator interface
func (cdn TOCDN) Validate(db *sqlx.DB) []error {
	validName := validation.NewStringRule(IsValidCDNName, "invalid characters found - Use alphanumeric . or - .")
	validDomainName := validation.NewStringRule(govalidator.IsDNSName, "not a valid domain name")
	errs := validation.Errors{
		"name":       validation.Validate(cdn.Name, validation.Required, validName),
		"domainName": validation.Validate(cdn.DomainName, validation.Required, validDomainName),
	}
	return tovalidate.ToErrors(errs)
}

//The TOCDN implementation of the Creator interface
//all implementations of Creator should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cdn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
//The insert sql returns the id and lastUpdated values of the newly inserted cdn and have
//to be added to the struct
func (cdn *TOCDN) Create(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	// make sure that cdn.DomainName is lowercase
	*cdn.DomainName = strings.ToLower(*cdn.DomainName)
	resultRows, err := tx.NamedQuery(insertQuery(), cdn)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a cdn with " + err.Error()), eType
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
		err = errors.New("no cdn was inserted, no id was returned")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	} else if rowsAffected > 1 {
		err = errors.New("too many ids returned from cdn insert")
		log.Errorln(err)
		return tc.DBError, tc.SystemError
	}
	cdn.SetKeys(map[string]interface{}{"id": id})
	cdn.LastUpdated = &lastUpdated
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

func (cdn *TOCDN) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"domainName":    dbhelpers.WhereColumnInfo{"domain_name", nil},
		"dnssecEnabled": dbhelpers.WhereColumnInfo{"dnssec_enabled", nil},
		"id":            dbhelpers.WhereColumnInfo{"id", api.IsInt},
		"name":          dbhelpers.WhereColumnInfo{"name", nil},
	}
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
	if err != nil {
		log.Errorf("Error querying CDNs: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	CDNs := []interface{}{}
	for rows.Next() {
		var s TOCDN
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing CDN rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		CDNs = append(CDNs, s)
	}

	return CDNs, []error{}, tc.NoError
}

//The TOCDN implementation of the Updater interface
//all implementations of Updater should use transactions and return the proper errorType
//ParsePQUniqueConstraintError is used to determine if a cdn with conflicting values exists
//if so, it will return an errorType of DataConflict and the type should be appended to the
//generic error message returned
func (cdn *TOCDN) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	log.Debugf("about to run exec query: %s with cdn: %++v", updateQuery(), cdn)
	// make sure that cdn.DomainName is lowercase
	*cdn.DomainName = strings.ToLower(*cdn.DomainName)
	resultRows, err := tx.NamedQuery(updateQuery(), cdn)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			err, eType := dbhelpers.ParsePQUniqueConstraintError(pqErr)
			if eType == tc.DataConflictError {
				return errors.New("a cdn with " + err.Error()), eType
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
	cdn.LastUpdated = &lastUpdated
	if rowsAffected != 1 {
		if rowsAffected < 1 {
			return errors.New("no cdn found with this id"), tc.DataMissingError
		} else {
			return fmt.Errorf("this update affected too many rows: %d", rowsAffected), tc.SystemError
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Errorln("Could not commit transaction: ", err)
		return tc.DBError, tc.SystemError
	}
	rollbackTransaction = false
	return nil, tc.NoError
}

//The CDN implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (cdn *TOCDN) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
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
	log.Debugf("about to run exec query: %s with cdn: %++v", deleteQuery(), cdn)
	result, err := tx.NamedExec(deleteQuery(), cdn)
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
			return errors.New("no cdn with that id found"), tc.DataMissingError
		} else {
			return fmt.Errorf("this create affected too many rows: %d", rowsAffected), tc.SystemError
		}
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
dnssec_enabled,
domain_name,
id,
last_updated,
name

FROM cdn c`
	return query
}

func updateQuery() string {
	query := `UPDATE
cdn SET
dnssec_enabled=:dnssec_enabled,
domain_name=:domain_name,
name=:name
WHERE id=:id RETURNING last_updated`
	return query
}

func insertQuery() string {
	query := `INSERT INTO cdn (
dnssec_enabled,
domain_name,
name) VALUES (
:dnssec_enabled,
:domain_name,
:name) RETURNING id,last_updated`
	return query
}

func deleteQuery() string {
	query := `DELETE FROM cdn
WHERE id=:id`
	return query
}

func DeleteNameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		ctx := r.Context()
		params, err := api.GetCombinedParams(r)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			log.Errorf("unable to retrieve current user from context: %s", err)
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		cdnName, hasCDNName := params["name"]
		if !hasCDNName || cdnName == "" {
			handleErrs(http.StatusBadRequest, err)
			return
		}
		if ok, err := IsTenantAuthorized(*user, cdnName, db); err != nil {
			handleErrs(http.StatusBadRequest, err)
			return
		} else if !ok {
			handleErrs(http.StatusForbidden, errors.New("not authorized on this tenant"))
			return
		}
		txSuccess := WithTx(db, handleErrs, func(tx *sql.Tx) bool {
			res, err := tx.Exec(`DELETE FROM cdn WHERE name=$1`, cdnName)
			if err != nil {
				log.Errorf("received error: %++v from named delete execution", err)
				tc.HandleErrWithType(tc.DBError, tc.SystemError, handleErrs)
				return false
			}
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				tc.HandleErrWithType(tc.DBError, tc.SystemError, handleErrs)
				return false
			}
			if rowsAffected == 0 {
				tc.HandleErrWithType(errors.New("not found"), tc.DataMissingError, handleErrs)
				return false
			}
			if rowsAffected > 1 {
				tc.HandleErrWithType(fmt.Errorf("affected too many rows: %d", rowsAffected), tc.SystemError, handleErrs)
				return false
			}
			return true
		})
		if !txSuccess {
			return // return; WithTx or its dbFunc called handleErrs which wrote the error response
		}
		WriteSuccess(w, "cdn was deleted", api.ApiChange, api.Deleted, cdnName, *user, db, handleErrs)
	}
}

// WriteSuccess writes the success message to the client, and logs to the change log
func WriteSuccess(w http.ResponseWriter, msg, changeLevel, changeAction, iden string, user auth.CurrentUser, db *sql.DB, handleErrs func(status int, errs ...error)) {
	log.Debugf("changelog for delete on object")
	api.CreateChangeLog(changeLevel, changeAction, TOCDN{Name: &iden}, user, db)
	resp := struct{ tc.Alerts }{tc.CreateAlerts(tc.SuccessLevel, msg)}
	respBts, err := json.Marshal(resp)
	if err != nil {
		handleErrs(http.StatusInternalServerError, err)
		return
	}
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write(respBts)
}

// WithTx takes a function dbFunc and calls it with a transaction.
// Returns whether the transaction succeeded or failed. All failures are written to the handleErrs function, and therefore to the user. If false is returned, the error has been written, and the caller shouldn't write anything else to the client.
// The dbFunc returns true if the transaction should be committed. If dbFunc returns false, it is assumed it responded to the user, and WithTx will make no further response, but only log rollback errors.
// The rowsAffected is the number of rows which should be affected on success. If negative, rows affected are not considered. If greater than or equal to zero, if the database returns any other number of rows affected, the transaction is rolled back and false (failure) returned.
func WithTx(db *sql.DB, handleErrs func(status int, errs ...error), dbFunc func(*sql.Tx) bool) bool {
	tx, err := db.Begin()
	if err != nil {
		log.Errorln("beginning transaction: " + err.Error())
		tc.HandleErrorsWithType([]error{tc.DBError}, tc.SystemError, handleErrs)
		return false
	}
	rollback := true
	defer func() {
		if !rollback {
			return
		}
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Errorln(errors.New("rolling back transaction: " + err.Error()))
		}
	}()
	if !dbFunc(tx) {
		return false
	}
	if err = tx.Commit(); err != nil {
		log.Errorln("committing transaction: " + err.Error())
		tc.HandleErrorsWithType([]error{tc.DBError}, tc.SystemError, handleErrs)
		return false
	}
	rollback = false
	return true
}

func IsTenantAuthorized(user auth.CurrentUser, name string, db *sql.DB) (bool, error) {
	return true, nil // TODO implement
}
