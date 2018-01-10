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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

const CDNsPrivLevel = 10

func GetHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		q := r.URL.Query()
		for k, v := range pathParams {
			if k == `id` {
				if _, err := strconv.Atoi(v); err != nil {
					log.Errorf("Expected {id} to be an integer: %s", v)
					handleErrs(http.StatusNotFound, errors.New("Resource not found.")) //matches perl response
					return
				}
			}
			q.Set(k, v)
		}

		resp, err := getCDNsResponse(q, db)

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

func getCDNsResponse(q url.Values, db *sqlx.DB) (*tc.CDNsResponse, error) {
	CDNs, err := getCDNs(q, db)
	if err != nil {
		return nil, fmt.Errorf("getting CDNs response: %v", err)
	}

	resp := tc.CDNsResponse{
		Response: CDNs,
	}
	return &resp, nil
}

func getCDNs(v url.Values, db *sqlx.DB) ([]tc.CDN, error) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]string{
		"domainName":    "domain_name",
		"dnssecEnabled": "dnssec_enabled",
		"id":            "id",
		"name":          "name",
	}

	query, queryValues := dbhelpers.BuildQuery(v, selectCDNsQuery(), queryParamsToQueryCols)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	CDNs := []tc.CDN{}
	for rows.Next() {
		var s tc.CDN
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting CDNs: %v", err)
		}
		CDNs = append(CDNs, s)
	}
	return CDNs, nil
}

func selectCDNsQuery() string {
	query := `SELECT
dnssec_enabled,
domain_name,
id,
last_updated,
name

FROM cdn c`
	return query
}
