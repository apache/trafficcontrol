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
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/experimental/tocookie"
	"log" // TODO change to traffic_monitor_golang/common/log
	"net/http"
	"time"
)

func wrapAuth(h RegexHandlerFunc, noAuth bool, secret string, privLevelStmt *sql.Stmt, privLevelRequired int) RegexHandlerFunc {
	if noAuth {
		return h
	}
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		handleUnauthorized := func(reason string) {
			log.Printf("%v %v %v sent 401 - %v\n", time.Now(), r.RemoteAddr, r.URL.Path, reason)
			status := http.StatusUnauthorized
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		cookie, err := r.Cookie(tocookie.Name)
		if err != nil {
			handleUnauthorized("error getting cookie: " + err.Error())
			return
		}

		if cookie == nil {
			handleUnauthorized("no auth cookie")
			return
		}

		oldCookie, err := tocookie.Parse(secret, cookie.Value)
		if err != nil {
			handleUnauthorized("cookie error: " + err.Error())
			return
		}

		username := oldCookie.AuthData
		if !hasPrivLevel(privLevelStmt, username, privLevelRequired) {
			handleUnauthorized("insufficient privileges")
			return
		}

		newCookieVal := tocookie.Refresh(oldCookie, secret)
		http.SetCookie(w, &http.Cookie{Name: tocookie.Name, Value: newCookieVal})
		h(w, r, p)
	}
}

func wrapLogTime(h RegexHandlerFunc) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		start := time.Now()
		defer func() {
			now := time.Now()
			log.Printf("%v %v served %v in %v\n", now, r.RemoteAddr, r.URL.Path, now.Sub(start))
		}()
		h(w, r, p)
	}
}
