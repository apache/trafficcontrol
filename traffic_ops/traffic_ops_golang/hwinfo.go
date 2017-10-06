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

const HWInfoPrivLevel = 10

func hwInfoHandler(db *sqlx.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		q := r.URL.Query()
		resp, err := getHWInfoResponse(q, db, privLevel)
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

func getHWInfoResponse(q url.Values, db *sqlx.DB, privLevel int) (*tc.HWInfoResponse, error) {
	hwInfo, err := getHWInfo(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting hwInfo response: %v", err)
	}

	resp := tc.HWInfoResponse{
		Response: hwInfo,
	}
	return &resp, nil
}

func getHWInfo(v url.Values, db *sqlx.DB, privLevel int) ([]tc.HWInfo, error) {

	var rows *sqlx.Rows
	var err error

	rows, err = db.Queryx(selectHWInfoQuery())

	if err != nil {
		//TODO: drichardson - send back an alert if the Query Count is larger than 1
		//                    Test for bad Query Parameters
		return nil, err
	}
	hwInfo := []tc.HWInfo{}

	defer rows.Close()
	for rows.Next() {
		var s tc.HWInfo
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting hwInfo: %v", err)
		}
		hwInfo = append(hwInfo, s)
	}
	return hwInfo, nil
}

func selectHWInfoQuery() string {

	query := `SELECT
    id,
    serverid,
    description,
    val,
    last_updated

FROM hwInfo c`
	return query
}
