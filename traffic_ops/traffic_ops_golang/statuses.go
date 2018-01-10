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
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
)

const StatusesPrivLevel = 10

func statusesHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		pathParams, err := api.GetPathParams(r.Context())
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		q := r.URL.Query()

		for k, v := range pathParams {
			if k == `id` {
				if _, err := strconv.Atoi(v); err != nil {
					log.Errorf("Expected {id} to be an integer: %s", v)
					handleErrs(http.StatusBadRequest, err)
					return
				}
			}
			q.Set(k, v)
		}

		resp, err := getStatusesResponse(q, db)

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

func getStatusesResponse(q url.Values, db *sqlx.DB) (*tc.StatusesResponse, error) {
	cdns, err := getStatuses(q, db)
	if err != nil {
		return nil, fmt.Errorf("getting cdns response: %v", err)
	}

	resp := tc.StatusesResponse{
		Response: cdns,
	}
	return &resp, nil
}

func getStatuses(v url.Values, db *sqlx.DB) ([]tc.Status, error) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]string{
		"id":          "id",
		"name":        "name",
		"description": "description",
	}

	query, queryValues := dbhelpers.BuildQuery(v, selectStatusesQuery(), queryParamsToSQLCols)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := []tc.Status{}
	for rows.Next() {
		var s tc.Status
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting statuses: %v", err)
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

func selectStatusesQuery() string {

	query := `SELECT
description,
id,
last_updated,
name 

FROM status c`
	return query
}
