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

const ParametersPrivLevel = PrivLevelReadOnly

func parametersHandler(db *sqlx.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		q := r.URL.Query()
		resp, err := getParametersResponse(q, db, privLevel)
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

func getParametersResponse(q url.Values, db *sqlx.DB, privLevel int) (*tc.ParametersResponse, error) {
	parameters, err := getParameters(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting parameters response: %v", err)
	}

	resp := tc.ParametersResponse{
		Response: parameters,
	}
	return &resp, nil
}

func getParameters(v url.Values, db *sqlx.DB, privLevel int) ([]tc.Parameter, error) {

	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]string{
		"config_file":  "p.config_file",
		"id":           "p.id",
		"last_updated": "p.last_updated",
		"name":         "p.name",
		"secure":       "p.secure",
	}

	query, queryValues := BuildQuery(v, selectParametersQuery(), queryParamsToSQLCols)

	log.Debugln("Query is ", query)
	rows, err = db.NamedQuery(query, queryValues)

	if err != nil {
		return nil, fmt.Errorf("querying: %v", err)
	}
	parameters := []tc.Parameter{}

	defer rows.Close()

	for rows.Next() {
		var s tc.Parameter
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting parameters: %v", err)
		}
		if s.Secure && privLevel < PrivLevelAdmin {
			// Secure params only visible to admin
			continue
		}
		parameters = append(parameters, s)
	}
	return parameters, nil
}

func selectParametersQuery() string {

	query := `SELECT
config_file,
id,
last_updated,
name,
value,
secure

FROM parameter p`
	return query
}
