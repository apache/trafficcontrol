package main

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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

func regionsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		resp, errs, errType := getRegionsResponse(params, db)
		if len(errs) > 0 {
			tc.HandleErrorsWithType(errs, errType, handleErrs)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getRegionsResponse(params map[string]string, db *sqlx.DB) (*tc.RegionsResponse, []error, tc.ApiErrorType) {
	regions, errs, errType := getRegions(params, db)
	if len(errs) > 0 {
		return nil, errs, errType
	}

	resp := tc.RegionsResponse{
		Response: regions,
	}
	return &resp, nil, tc.NoError
}

func getRegions(params map[string]string, db *sqlx.DB) ([]tc.Region, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"division": dbhelpers.WhereColumnInfo{"d.id", api.IsInt},
		"id":       dbhelpers.WhereColumnInfo{"r.id", api.IsInt},
		"name":     dbhelpers.WhereColumnInfo{"r.name", nil},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectRegionsQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, []error{err}, tc.SystemError
	}
	defer rows.Close()

	regions := []tc.Region{}
	for rows.Next() {
		var s tc.Region
		if err = rows.StructScan(&s); err != nil {
			return nil, []error{fmt.Errorf("getting regions: %v", err)}, tc.SystemError
		}
		regions = append(regions, s)
	}
	return regions, nil, tc.NoError
}

func selectRegionsQuery() string {

	query := `SELECT
r.division,
d.name as divisionname,
r.id,
r.last_updated,
r.name

FROM region r
JOIN division d ON r.division = d.id`
	return query
}
