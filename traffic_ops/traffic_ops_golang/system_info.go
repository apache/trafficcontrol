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

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
)

const SystemInfoPrivLevel = 10

func systemInfoHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		privLevel := user.PrivLevel

		q := r.URL.Query()
		resp, err := getSystemInfoResponse(q, db, privLevel)
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
func getSystemInfoResponse(q url.Values, db *sqlx.DB, privLevel int) (*tc.SystemInfoResponse, error) {
	info, err := getSystemInfo(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting SystemInfo: %v", err)
	}

	resp := tc.SystemInfoResponse{}
	resp.Response.Parameters = info
	return &resp, nil
}

func getSystemInfo(_ url.Values, db *sqlx.DB, privLevel int) (map[string]string, error) {
	// system info returns all global parameters
	// no parameters on the url, but use that mechanism to get the right params from the db

	v := url.Values{}
	v.Set("config_file", "global")
	params, err := getParameters(v, db, SystemInfoPrivLevel)
	if err != nil {
		return nil, err
	}

	info := make(map[string]string, len(params))
	for _, p := range params {
		info[p.Name] = p.Value
	}
	return info, nil
}
