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
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

const DivisionsPrivLevel = 10

func divisionsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		q := r.URL.Query()
		resp, err := getDivisionsResponse(q, db)
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

func getDivisionsResponse(q url.Values, db *sqlx.DB) (*tc.DivisionsResponse, error) {
	divisions, err := getDivisions(q, db)
	if err != nil {
		return nil, fmt.Errorf("getting divisions response: %v", err)
	}

	resp := tc.DivisionsResponse{
		Response: divisions,
	}
	return &resp, nil
}

func getDivisions(v url.Values, db *sqlx.DB) ([]tc.Division, error) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]string{
		"id":   "id",
		"name": "name",
	}

	query, queryValues := dbhelpers.BuildQuery(v, selectDivisionsQuery(), queryParamsToSQLCols)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	o := []tc.Division{}
	for rows.Next() {
		var d tc.Division
		if err = rows.StructScan(&d); err != nil {
			return nil, fmt.Errorf("getting divisions: %v", err)
		}
		o = append(o, d)
	}
	return o, nil
}

func selectDivisionsQuery() string {

	query := `SELECT
id,
last_updated,
name 

FROM division d`
	return query
}
