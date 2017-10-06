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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/jmoiron/sqlx"
)

const RegionsPrivLevel = 10

func regionsHandler(db *sqlx.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		// Load the PathParams into the query parameters for pass through
		q := r.URL.Query()
		for k, v := range p {
			q.Set(k, v)
		}
		resp, err := getRegionsResponse(q, db, privLevel)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getRegionsResponse(q url.Values, db *sqlx.DB, privLevel int) (*tc.RegionsResponse, error) {
	regions, err := getRegions(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting regions response: %v", err)
	}

	resp := tc.RegionsResponse{
		Response: regions,
	}
	return &resp, nil
}

func getRegions(v url.Values, db *sqlx.DB, privLevel int) ([]tc.Region, error) {

	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]string{
		"division": "division",
		"id":       "id",
		"name":     "name",
	}

	query, queryValues := BuildQuery(v, selectRegionsQuery(), queryParamsToQueryCols)

	rows, err = db.NamedQuery(query, queryValues)

	if err != nil {
		return nil, err
	}
	regions := []tc.Region{}

	defer rows.Close()
	for rows.Next() {
		var s tc.Region
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting regions: %v", err)
		}
		regions = append(regions, s)
	}
	return regions, nil
}

func selectRegionsQuery() string {

	query := `SELECT
division,
id,
last_updated,
name 

FROM region`
	return query
}
