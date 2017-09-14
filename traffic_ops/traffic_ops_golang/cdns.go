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
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tcstructs"
)

const CDNsPrivLevel = 10

func cdnsHandler(db *sql.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		q := r.URL.Query()
		resp, err := getCDNsResponse(q, db, privLevel)
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

func getCDNsResponse(q url.Values, db *sql.DB, privLevel int) (*tcstructs.CDNsResponse, error) {
	cdns, err := getCDNs(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting cdns response: %v", err)
	}

	resp := tcstructs.CDNsResponse{
		Response: cdns,
	}
	return &resp, nil
}

func getCDNs(v url.Values, db *sql.DB, privLevel int) ([]tcstructs.CDN, error) {
	rows, err := db.Query(selectCDNsQuery())
	if err != nil {
		//TODO: drichardson - send back an alert if the Query Count is larger than 1
		//                    Test for bad Query Parameters
		return nil, err
	}
	defer rows.Close()

	cdns := []tcstructs.CDN{}
	for rows.Next() {
		s := tcstructs.CDN{}
		if err = rows.Scan(&s.DNSSECEnabled, &s.DomainName, &s.ID, &s.LastUpdated, &s.Name); err != nil {
			return nil, fmt.Errorf("getting cdns: %v", err)
		}
		cdns = append(cdns, s)
	}
	return cdns, nil
}

func selectCDNsQuery() string {
	return `
SELECT
 dnssec_enabled,
 domain_name,
 id,
 last_updated,
 name
FROM cdn c
`
}
