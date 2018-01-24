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

// ASNsPrivLevel ...
const ASNsPrivLevel = 10

// ASNsHandler ...
func ASNsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		resp, errs, errType := getASNsResponse(params, db)
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

func getASNsResponse(parameters map[string]string, db *sqlx.DB) (*tc.ASNsResponse, []error, tc.ApiErrorType) {
	asns, errs, errType := getASNs(parameters, db)
	if len(errs) > 0 {
		return nil, errs, errType
	}

	resp := tc.ASNsResponse{
		Response: asns,
	}
	return &resp, nil, tc.NoError
}

func getASNs(parameters map[string]string, db *sqlx.DB) ([]tc.ASN, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"asn":        dbhelpers.WhereColumnInfo{"a.asn", api.IsInt},
		"id":         dbhelpers.WhereColumnInfo{"a.id", api.IsInt},
		"cachegroup": dbhelpers.WhereColumnInfo{"cg.id", api.IsInt},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}
	query := selectASNsQuery() + where + orderBy
	log.Debugln("Query is ", query)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, []error{err}, tc.SystemError
	}
	defer rows.Close()

	ASNs := []tc.ASN{}
	for rows.Next() {
		var s tc.ASN
		if err = rows.StructScan(&s); err != nil {
			return nil, []error{fmt.Errorf("getting ASNs: %v", err)}, tc.SystemError
		}
		ASNs = append(ASNs, s)
	}
	return ASNs, nil, tc.NoError
}

func selectASNsQuery() string {

	query := `SELECT
a.asn,
cg.name as cachegroup,
a.cachegroup as cachegroup_id,
a.id,
a.last_updated

FROM asn a
JOIN cachegroup cg ON cg.id = a.cachegroup`
	return query
}
