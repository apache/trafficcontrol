package systeminfo

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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
)

func Handler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		privLevel := user.PrivLevel

		cfg, ctxErr := api.GetConfig(ctx)
		if ctxErr != nil {
			log.Errorln("unable to retrieve config from context: ", ctxErr)
			handleErrs(http.StatusInternalServerError, errors.New("no config found in context"))
		}

		resp, err := getSystemInfoResponse(db, privLevel, time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second)
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
func getSystemInfoResponse(db *sqlx.DB, privLevel int, timeout time.Duration) (*tc.SystemInfoResponse, error) {
	info, err := getSystemInfo(db, privLevel, timeout)
	if err != nil {
		return nil, fmt.Errorf("getting SystemInfo: %v", err)
	}

	resp := tc.SystemInfoResponse{}
	resp.Response.ParametersNullable = info
	return &resp, nil
}

func getSystemInfo(db *sqlx.DB, privLevel int, timeout time.Duration) (map[string]string, error) {
	// system info returns all global parameters
	query := `SELECT
p.name,
p.secure,
p.last_updated,
p.value
FROM parameter p
WHERE p.config_file='global'`
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	rows, err := db.QueryxContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("querying: %v", err)
	}
	defer rows.Close()

	info := make(map[string]string)
	for rows.Next() {
		p := tc.ParameterNullable{}
		if err = rows.StructScan(&p); err != nil {
			return nil, fmt.Errorf("getting system_info: %v", err)
		}

		var isSecure bool
		if p.Secure != nil {
			isSecure = *p.Secure
		}

		name := p.Name
		value := p.Value
		if isSecure && privLevel < auth.PrivLevelAdmin {
			// Secure params only visible to admin
			continue
		}

		if name != nil && value != nil {
			info[*name] = *value
		}
	}
	if err != nil {
		return nil, err
	}

	return info, nil
}
