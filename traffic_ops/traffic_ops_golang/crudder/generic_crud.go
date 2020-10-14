package crudder

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
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	ims "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"
)

type GenericCreator interface {
	GetType() string
	APIInfo() *api.APIInfo
	SetKeys(map[string]interface{})
	SetLastUpdated(tc.TimeNoMod)
	InsertQuery() string
}

type GenericReader interface {
	GetType() string
	APIInfo() *api.APIInfo
	ParamColumns() map[string]dbhelpers.WhereColumnInfo
	NewReadObj() interface{}
	SelectQuery() string
	SelectMaxLastUpdatedQuery(where string, orderBy string, pagination string, tableName string) string
}

type GenericUpdater interface {
	GetType() string
	APIInfo() *api.APIInfo
	SetLastUpdated(tc.TimeNoMod)
	UpdateQuery() string
	GetLastUpdated() (*time.Time, bool, error)
}

type GenericDeleter interface {
	GetType() string
	APIInfo() *api.APIInfo
	DeleteQuery() string
}

// GenericOptionsDeleter can use any key listed in DeleteKeyOptions() to delete a resource.
type GenericOptionsDeleter interface {
	GetType() string
	APIInfo() *api.APIInfo
	DeleteKeyOptions() map[string]dbhelpers.WhereColumnInfo
	DeleteQueryBase() string
}

// GenericCreate does a Create (POST) for the given GenericCreator object and type. This exists as a generic function, for the common use case of a single "id" key and a lastUpdated field.
func GenericCreate(val GenericCreator) api.Errors {
	resultRows, err := val.APIInfo().Tx.NamedQuery(val.InsertQuery(), val)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	var id interface{}
	lastUpdated := tc.TimeNoMod{}
	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&id, &lastUpdated); err != nil {
			return api.Errors{
				Code:        http.StatusInternalServerError,
				SystemError: fmt.Errorf("%v create scanning: %v", val.GetType(), err),
			}
		}
	}

	switch id.(type) {
	case int64:
		// ahhh generics in a language without generics
		// sql.Driver values return int64s from integral database data
		// but our objects all use ambiguous-width ints for their IDs
		// so naturally this suffers from overflow issues, but c'est la vie
		id = int(id.(int64))
	default:
	}
	if rowsAffected == 0 {
		return api.Errors{
			SystemError: fmt.Errorf("%s create: no %s was inserted, no id was returned", val.GetType(), val.GetType()),
			Code:        http.StatusInternalServerError,
		}
	}
	if rowsAffected > 1 {
		return api.Errors{
			SystemError: fmt.Errorf("too many ids returned from %s insert", val.GetType()),
			Code:        http.StatusInternalServerError,
		}
	}
	val.SetKeys(map[string]interface{}{"id": id})
	val.SetLastUpdated(lastUpdated)
	return api.NewErrors()
}

// GenericCreateNameBasedID does a Create (POST) for the given GenericCreator object and type. This exists as a generic function, for the use case of a single "name" key (not a numerical "id" key) and a lastUpdated field.
func GenericCreateNameBasedID(val GenericCreator) api.Errors {
	resultRows, err := val.APIInfo().Tx.NamedQuery(val.InsertQuery(), val)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	lastUpdated := tc.TimeNoMod{}
	rowsAffected := 0
	errs := api.NewErrors()
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.Scan(&lastUpdated); err != nil {
			errs.SystemError = fmt.Errorf("%s create scanning: %v", val.GetType(), err)
			errs.Code = http.StatusInternalServerError
			return errs
		}
	}
	if rowsAffected == 0 {
		errs.Code = http.StatusInternalServerError
		errs.SystemError = fmt.Errorf("%s create: no %s was inserted, no row was returned", val.GetType(), val.GetType())
		return errs
	} else if rowsAffected > 1 {
		errs.Code = http.StatusInternalServerError
		errs.SystemError = fmt.Errorf("too many rows returned from %s insert", val.GetType())
		return errs
	}
	val.SetLastUpdated(lastUpdated)
	return errs
}

// TryIfModifiedSinceQuery checks to see the max time that an entity was changed, and then returns a boolean (which tells us whether or not to run the main query for the endpoint)
// along with the max time
// If the returned boolean is false, there is no need to run the main query for the GET API endpoint, and we return a 304 status
func TryIfModifiedSinceQuery(val GenericReader, h http.Header, where string, orderBy string, pagination string, queryValues map[string]interface{}) (bool, time.Time) {
	var max time.Time
	var imsDate time.Time
	var ok bool
	imsDateHeader := []string{}
	runSecond := true
	dontRunSecond := false
	if h == nil {
		return runSecond, max
	}
	imsDateHeader = h[rfc.IfModifiedSince]
	if len(imsDateHeader) == 0 {
		return runSecond, max
	}
	if imsDate, ok = rfc.ParseHTTPDate(imsDateHeader[0]); !ok {
		log.Warnf("IMS request header date '%s' not parsable", imsDateHeader[0])
		return runSecond, max
	}
	// ToDo: Remove orderBy, pagination from all the implementations, and eventually remove it from the function definition
	query := val.SelectMaxLastUpdatedQuery(where, "", "", val.GetType())
	rows, err := val.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		log.Warnf("Couldn't get the max last updated time: %v", err)
		return runSecond, max
	}
	if err == sql.ErrNoRows {
		return dontRunSecond, max
	}
	defer rows.Close()
	// This should only ever contain one row
	if rows.Next() {
		v := &ims.LatestTimestamp{}
		if err = rows.StructScan(v); err != nil || v == nil {
			log.Warnf("Failed to parse the max time stamp into a struct %v", err)
			return runSecond, max
		}
		if v.LatestTime != nil {
			max = v.LatestTime.Time
			// The request IMS time is later than the max of (lastUpdated, deleted_time)
			if imsDate.After(v.LatestTime.Time) {
				return dontRunSecond, max
			}
		}
	}
	return runSecond, max
}

func GenericRead(h http.Header, val GenericReader, useIMS bool) ([]interface{}, api.Errors, *time.Time) {
	vals := []interface{}{}
	var maxTime time.Time
	var runSecond bool
	where, orderBy, pagination, queryValues, es := dbhelpers.BuildWhereAndOrderByAndPagination(val.APIInfo().Params, val.ParamColumns())
	if len(es) > 0 {
		return nil, api.Errors{UserError: util.JoinErrs(es), Code: http.StatusBadRequest}, nil
	}
	if useIMS {
		runSecond, maxTime = TryIfModifiedSinceQuery(val, h, where, orderBy, pagination, queryValues)
		if !runSecond {
			log.Debugln("IMS HIT")
			return vals, api.Errors{Code: http.StatusNotModified}, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	// Case where we need to run the second query
	query := val.SelectQuery() + where + orderBy + pagination
	rows, err := val.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, api.NewSystemError(fmt.Errorf("querying %s: %w", val.GetType(), err)), &maxTime
	}
	defer rows.Close()

	for rows.Next() {
		v := val.NewReadObj()
		if err = rows.StructScan(v); err != nil {
			return nil, api.NewSystemError(fmt.Errorf("scanning %s: %w", val.GetType(), err)), &maxTime
		}
		vals = append(vals, v)
	}
	return vals, api.NewErrors(), &maxTime
}

// GenericUpdate handles the common update case, where the update returns the new last_modified time.
func GenericUpdate(h http.Header, val GenericUpdater) api.Errors {
	existingLastUpdated, found, err := val.GetLastUpdated()
	if err == nil && found == false {
		return api.Errors{UserError: errors.New("no " + val.GetType() + " found with this id"), Code: http.StatusNotFound}
	}
	if err != nil {
		return api.Errors{SystemError: err, Code: http.StatusInternalServerError}
	}
	if !api.IsUnmodified(h, *existingLastUpdated) {
		return api.Errors{UserError: api.ResourceModifiedError, Code: http.StatusPreconditionFailed}
	}

	rows, err := val.APIInfo().Tx.NamedQuery(val.UpdateQuery(), val)
	if err != nil {
		errs := api.ParseDBError(err)
		return errs
	}
	defer rows.Close()

	if !rows.Next() {
		return api.Errors{UserError: errors.New("no " + val.GetType() + " found with this id"), Code: http.StatusNotFound}
	}
	lastUpdated := tc.TimeNoMod{}
	if err := rows.Scan(&lastUpdated); err != nil {
		return api.Errors{SystemError: fmt.Errorf("scanning lastUpdated from %s insert: %w", val.GetType(), err), Code: http.StatusInternalServerError}
	}
	val.SetLastUpdated(lastUpdated)
	if rows.Next() {
		return api.Errors{SystemError: errors.New(val.GetType() + " update affected too many rows: >1"), Code: http.StatusInternalServerError}
	}
	return api.NewErrors()
}

// GenericOptionsDelete does a Delete (DELETE) for the given GenericOptionsDeleter object and type. Unlike
// GenericDelete, there is no requirement that a specific key is used as the parameter.
// GenericOptionsDeleter.DeleteKeyOptions() specifies which keys can be used.
func GenericOptionsDelete(val GenericOptionsDeleter) api.Errors {
	where, _, _, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(val.APIInfo().Params, val.DeleteKeyOptions())
	if len(errs) > 0 {
		return api.Errors{UserError: util.JoinErrs(errs), Code: http.StatusBadRequest}
	}

	query := val.DeleteQueryBase() + where
	tx := val.APIInfo().Tx
	result, err := tx.NamedExec(query, queryValues)
	if err != nil {
		return api.ParseDBError(err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return api.Errors{SystemError: fmt.Errorf("deleting %s: getting rows affected: %w", val.GetType(), err), Code: http.StatusInternalServerError}
	} else if rowsAffected < 1 {
		return api.Errors{UserError: errors.New("no " + val.GetType() + " with that key found"), Code: http.StatusNotFound}
	} else if rowsAffected > 1 {
		return api.Errors{SystemError: fmt.Errorf(val.GetType()+" delete affected too many rows: %d", rowsAffected), Code: http.StatusInternalServerError}
	}

	return api.NewErrors()
}

// GenericDelete does a Delete (DELETE) for the given GenericDeleter object and type. This exists as a generic function, for the common use case of a simple delete with query parameters defined in the sqlx struct tags.
func GenericDelete(val GenericDeleter) api.Errors {
	result, err := val.APIInfo().Tx.NamedExec(val.DeleteQuery(), val)
	if err != nil {
		errs := api.ParseDBError(err)
		return errs
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return api.Errors{SystemError: fmt.Errorf("deleting %s: getting rows affected: %w", val.GetType(), err), Code: http.StatusInternalServerError}
	} else if rowsAffected < 1 {
		return api.Errors{UserError: errors.New("no " + val.GetType() + " with that key found"), Code: http.StatusNotFound}
	} else if rowsAffected > 1 {
		return api.Errors{SystemError: fmt.Errorf(val.GetType()+" delete affected too many rows: %d", rowsAffected), Code: http.StatusInternalServerError}
	}
	return api.NewErrors()
}
