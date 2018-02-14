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
	"fmt"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

//we need a type alias to define functions on
type TODivision tc.Division

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TODivision(tc.Division{})

func GetRefType() *TODivision {
	return &refType
}

func (cdn *TODivision) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
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

	query := selectDivisionsQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err := db.NamedQuery(query, queryValues)
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

func getDivisionsResponse(params map[string]string, db *sqlx.DB) (*tc.DivisionsResponse, []error, tc.ApiErrorType) {
	divisions, errs, errType := getDivisions(params, db)
	if len(errs) > 0 {
		return nil, errs, errType
	}

	resp := tc.DivisionsResponse{
		Response: divisions,
	}
	return &resp, nil, tc.NoError
}

func getDivisions(params map[string]string, db *sqlx.DB) ([]tc.Division, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":   dbhelpers.WhereColumnInfo{"id", api.IsInt},
		"name": dbhelpers.WhereColumnInfo{"name", nil},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectDivisionsQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, []error{err}, tc.SystemError
	}
	defer rows.Close()

	o := []tc.Division{}
	for rows.Next() {
		var d tc.Division
		if err = rows.StructScan(&d); err != nil {
			return nil, []error{fmt.Errorf("getting divisions: %v", err)}, tc.SystemError
		}
		o = append(o, d)
	}
	return o, nil, tc.NoError
}

func selectDivisionsQuery() string {

	query := `SELECT
id,
last_updated,
name 

FROM division d`
	return query
}
