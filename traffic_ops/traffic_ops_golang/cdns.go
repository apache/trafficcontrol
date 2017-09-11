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

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tcstructs"
	"github.com/jmoiron/sqlx"
)

const CdnsPrivLevel = 10

func cdnsHandler(db *sqlx.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		q := r.URL.Query()
		resp, err := getCdnsResponse(q, db, privLevel)
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

func getCdnsResponse(q url.Values, db *sqlx.DB, privLevel int) (*tcstructs.CdnsResponse, error) {
	cdns, err := getCdns(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting cdns response: %v", err)
	}

	resp := tcstructs.CdnsResponse{
		Response: cdns,
	}
	return &resp, nil
}

func getCdns(v url.Values, db *sqlx.DB, privLevel int) ([]tcstructs.Cdn, error) {

	var rows *sqlx.Rows
	var err error

	rows, err = db.Queryx(selectCdnsQuery())

	if err != nil {
		//TODO: drichardson - send back an alert if the Query Count is larger than 1
		//                    Test for bad Query Parameters
		return nil, err
	}
	cdns := []tcstructs.Cdn{}

	defer rows.Close()
	for rows.Next() {
		var s tcstructs.Cdn
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting cdns: %v", err)
		}
		cdns = append(cdns, s)
	}
	return cdns, nil
}

func selectCdnsQuery() string {

	query := `SELECT
dnssec_enabled,
domain_name,
id,
last_updated,
name 

FROM cdn c`
	return query
}
