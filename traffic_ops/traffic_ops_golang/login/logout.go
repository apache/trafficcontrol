package login

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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tocookie"
)

func LogoutHandler(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
		tx := inf.Tx.Tx
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
		defer inf.Close()

		cookie := tocookie.GetCookie(inf.User.UserName, 0, secret)
		http.SetCookie(w, cookie)
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(0),
			MaxAge:   0,
			HttpOnly: true, // prevents the cookie being accessed by Javascript. DO NOT remove, security vulnerability
		})
		resp := struct {
			tc.Alerts
		}{tc.CreateAlerts(tc.SuccessLevel, "You are logged out.")}

		respBts, err := json.Marshal(resp)
		if err != nil {
			errCode = http.StatusInternalServerError
			sysErr = fmt.Errorf("Marshaling response: %v", err)
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}

		w.Header().Set(rfc.ContentType, rfc.ApplicationJSON)
		api.WriteAndLogErr(w, r, append(respBts, '\n'))
	}
}
