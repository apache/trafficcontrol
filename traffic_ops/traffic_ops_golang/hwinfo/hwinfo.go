package hwinfo

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

func HWInfoHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		resp, errs, errType := getHWInfoResponse(params, db)
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

func getHWInfoResponse(params map[string]string, db *sqlx.DB) (*tc.HWInfoResponse, []error, tc.ApiErrorType) {
	hwInfo, errs, errType := getHWInfo(params, db)
	if len(errs) > 0 {
		return nil, errs, errType
	}

	resp := tc.HWInfoResponse{
		Response: hwInfo,
	}
	return &resp, nil, tc.NoError
}

func getHWInfo(params map[string]string, db *sqlx.DB) ([]tc.HWInfo, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":             dbhelpers.WhereColumnInfo{"h.id", api.IsInt},
		"serverHostName": dbhelpers.WhereColumnInfo{"s.host_name", nil},
		"serverId":       dbhelpers.WhereColumnInfo{"s.id", api.IsInt}, // TODO: this can be either s.id or h.serverid not sure what makes the most sense
		"description":    dbhelpers.WhereColumnInfo{"h.description", nil},
		"val":            dbhelpers.WhereColumnInfo{"h.val", nil},
		"lastUpdated":    dbhelpers.WhereColumnInfo{"h.last_updated", nil}, //TODO: this doesn't appear to work needs debugging
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}

	query := selectHWInfoQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, []error{err}, tc.SystemError
	}
	defer rows.Close()

	hwInfo := []tc.HWInfo{}
	for rows.Next() {
		var s tc.HWInfo
		if err = rows.StructScan(&s); err != nil {
			return nil, []error{fmt.Errorf("getting hwInfo: %v", err)}, tc.SystemError
		}
		hwInfo = append(hwInfo, s)
	}
	return hwInfo, nil, tc.NoError
}

func selectHWInfoQuery() string {

	query := `SELECT
	s.host_name as serverhostname,
    h.id,
    h.serverid,
    h.description,
    h.val,
    h.last_updated

FROM hwInfo h

JOIN server s ON s.id = h.serverid`
	return query
}
