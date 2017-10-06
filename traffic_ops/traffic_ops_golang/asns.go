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

const ASNSPrivLevel = 10

func ASNsHandler(db *sqlx.DB) AuthRegexHandlerFunc {
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
		resp, err := getASNsResponse(q, db, privLevel)
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

func getASNsResponse(q url.Values, db *sqlx.DB, privLevel int) (*tc.ASNsResponse, error) {
	asns, err := getASNs(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting asns response: %v", err)
	}

	resp := tc.ASNsResponse{
		Response: asns,
	}
	return &resp, nil
}

func getASNs(v url.Values, db *sqlx.DB, privLevel int) ([]tc.ASN, error) {

	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]string{
		"asn":        "asn",
		"id":         "id",
		"cachegroup": "cachegroup",
	}

	query, queryValues := BuildQuery(v, selectASNsQuery(), queryParamsToQueryCols)

	rows, err = db.NamedQuery(query, queryValues)

	if err != nil {
		return nil, err
	}
	ASNs := []tc.ASN{}

	defer rows.Close()
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
asn,
cachegroup,
id,
last_updated

FROM asn`
	return query
}
