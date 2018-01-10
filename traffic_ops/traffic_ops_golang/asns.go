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
	"net/url"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

const ASNsPrivLevel = 10

func ASNsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		// Load the PathParams into the query parameters for pass through
		q := r.URL.Query()
		for k, v := range pathParams {
			q.Set(k, v)
		}
		resp, err := getASNsResponse(q, db)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
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

func getASNsResponse(q url.Values, db *sqlx.DB) (*tc.ASNsResponse, error) {
	asns, err := getASNs(q, db)
	if err != nil {
		return nil, fmt.Errorf("getting asns response: %v", err)
	}

	resp := tc.ASNsResponse{
		Response: asns,
	}
	return &resp, nil
}

func getASNs(v url.Values, db *sqlx.DB) ([]tc.ASN, error) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]string{
		"asn":        "a.asn",
		"id":         "a.id",
		"cachegroup": "cg.id",
	}

	query, queryValues := dbhelpers.BuildQuery(v, selectASNsQuery(), queryParamsToQueryCols)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ASNs := []tc.ASN{}
	for rows.Next() {
		var s tc.ASN
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting ASNs: %v", err)
		}
		ASNs = append(ASNs, s)
	}
	return ASNs, nil
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
