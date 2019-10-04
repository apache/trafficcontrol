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

import "encoding/json"
import "fmt"
import "net/http"
import "time"

import "github.com/apache/trafficcontrol/lib/go-tc"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tocookie"

func LogoutHandler(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
		tx := inf.Tx.Tx
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, tx, errCode, userErr, sysErr)
			return
		}
		defer inf.Close()

		expiry := time.Now()
		cookie := tocookie.New(inf.User.UserName, expiry, secret)
		httpCookie := http.Cookie{
			Name:  tocookie.Name,
			Value: cookie, Path: "/",
			Expires:  expiry,
			HttpOnly: true,
		}

		http.SetCookie(w, &httpCookie)
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

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		w.Write(append(respBts, '\n'))
	}
}
