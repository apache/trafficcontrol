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

const PhysLocationsPrivLevel = 10

func physLocationsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		ctx := r.Context()
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		// Load the PathParams into the query parameters for pass through
		q := r.URL.Query()
		for k, v := range pathParams {
			q.Set(k, v)
		}
		resp, err := getPhysLocationsResponse(q, db)
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

func getPhysLocationsResponse(q url.Values, db *sqlx.DB) (*tc.PhysLocationsResponse, error) {
	physLocations, err := getPhysLocations(q, db)
	if err != nil {
		return nil, fmt.Errorf("getting physLocations response: %v", err)
	}

	resp := tc.PhysLocationsResponse{
		Response: physLocations,
	}
	return &resp, nil
}

func getPhysLocations(v url.Values, db *sqlx.DB) ([]tc.PhysLocation, error) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]string{
		"id":       "pl.id",
	}

	query, queryValues := BuildQuery(v, selectPhysLocationsQuery(), queryParamsToQueryCols)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	physLocations := []tc.PhysLocation{}
	for rows.Next() {
		var s tc.PhysLocation
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting physLocations: %v", err)
		}
		physLocations = append(physLocations, s)
	}
	return physLocations, nil
}

func selectPhysLocationsQuery() string {

query := `SELECT
pl.id,
pl.name,
pl.short_name,
pl.address,
pl.city,
pl.state,
pl.zip,
pl.region as regionid,
COALESCE(pl.poc, '') as poc,
COALESCE(pl.phone, '') as phone,
COALESCE(pl.email, '') as email,
COALESCE(pl.comments, '') as comments,
r.name as regionname,
pl.last_updated
FROM phys_location pl
JOIN region r ON pl.region = r.id` 
return query
}
