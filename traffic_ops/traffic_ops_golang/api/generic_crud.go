package api

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

type GenericCreator interface {
	GetType() string
	APIInfo() *APIInfo
	SetKeys(map[string]interface{})
	SetLastUpdated(tc.TimeNoMod)
	InsertQuery() string
}

type GenericReader interface {
	GetType() string
	APIInfo() *APIInfo
	ParamColumns() map[string]dbhelpers.WhereColumnInfo
	NewReadObj() interface{}
	SelectQuery() string
}

type GenericUpdater interface {
	GetType() string
	APIInfo() *APIInfo
	SetLastUpdated(tc.TimeNoMod)
	UpdateQuery() string
}

type GenericDeleter interface {
	GetType() string
	APIInfo() *APIInfo
	DeleteQuery() string
}

// GenericCreate does a Create (POST) for the given GenericCreator object and type. This exists as a generic function, for the common use case of a single "id" key and a lastUpdated field.
func GenericCreate(val GenericCreator) (error, error, int) {
	resultRows, err := val.APIInfo().Tx.NamedQuery(val.InsertQuery(), val)
	if err != nil {
		return ParseDBError(err)
	}
	defer resultRows.Close()

	id := 0
	lastUpdated := tc.TimeNoMod{}
	if err = ScanIntoItems(resultRows.Rows, []interface{}{&id, &lastUpdated}); err != nil {
		return nil, fmt.Errorf("generic create for %s: %v", val.GetType(), err), http.StatusInternalServerError
	}
	val.SetKeys(map[string]interface{}{"id": id})
	val.SetLastUpdated(lastUpdated)
	return nil, nil, http.StatusOK
}

func GenericRead(val GenericReader) ([]interface{}, error, error, int) {
	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(val.APIInfo().Params, val.ParamColumns())
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}

	query := val.SelectQuery() + where + orderBy
	rows, err := val.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("querying " + val.GetType() + ": " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	vals := []interface{}{}
	for rows.Next() {
		v := val.NewReadObj()
		if err = rows.StructScan(v); err != nil {
			return nil, nil, errors.New("scanning " + val.GetType() + ": " + err.Error()), http.StatusInternalServerError
		}
		vals = append(vals, v)
	}
	return vals, nil, nil, http.StatusOK
}

// GenericUpdate handles the common update case, where the update returns the new last_modified time.
func GenericUpdate(val GenericUpdater) (error, error, int) {
	rows, err := val.APIInfo().Tx.NamedQuery(val.UpdateQuery(), val)
	if err != nil {
		return ParseDBError(err)
	}
	defer rows.Close()

	lastUpdated := tc.TimeNoMod{}
	if err = ScanIntoItem(rows.Rows, &lastUpdated); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no %s found with this id", val.GetType()), http.StatusNotFound
		}
		return nil, fmt.Errorf("generic update for %s: %v", val.GetType(), err), http.StatusInternalServerError
	}
	val.SetLastUpdated(lastUpdated)

	return nil, nil, http.StatusOK
}

// GenericDelete does a Delete (DELETE) for the given GenericDeleter object and type. This exists as a generic function, for the common use case of a simple delete with query parameters defined in the sqlx struct tags.
func GenericDelete(val GenericDeleter) (error, error, int) {
	result, err := val.APIInfo().Tx.NamedExec(val.DeleteQuery(), val)
	if err != nil {
		return nil, errors.New("deleting " + val.GetType() + ": " + err.Error()), http.StatusInternalServerError
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return nil, errors.New("deleting " + val.GetType() + ": getting rows affected: " + err.Error()), http.StatusInternalServerError
	} else if rowsAffected < 1 {
		return errors.New("no " + val.GetType() + " with that id found"), nil, http.StatusNotFound
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf(val.GetType()+" delete affected too many rows: %d", rowsAffected), http.StatusInternalServerError
	}
	return nil, nil, http.StatusOK
}
