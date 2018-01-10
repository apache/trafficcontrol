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

		resp, err := getSystemInfoResponse(db, privLevel)
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
func getSystemInfoResponse(db *sqlx.DB, privLevel int) (*tc.SystemInfoResponse, error) {
	info, err := getSystemInfo(db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting SystemInfo: %v", err)
	}

	resp := tc.SystemInfoResponse{}
	resp.Response.Parameters = info
	return &resp, nil
}

func getSystemInfo(db *sqlx.DB, privLevel int) (map[string]string, error) {
	// system info returns all global parameters
	query := `SELECT
p.name,
p.secure,
p.value
FROM parameter p
WHERE p.config_file='global'`

	rows, err := db.Queryx(query)

	if err != nil {
		return nil, fmt.Errorf("querying: %v", err)
	}
	defer rows.Close()

	info := make(map[string]string)
	for rows.Next() {
		p := tc.Parameter{}
		if err = rows.StructScan(&p); err != nil {
			return nil, fmt.Errorf("getting system_info: %v", err)
		}
		if p.Secure && privLevel < auth.PrivLevelAdmin {
			// Secure params only visible to admin
			continue
		}
		info[p.Name] = p.Value
	}
	if err != nil {
		return nil, err
	}

	return info, nil
}
