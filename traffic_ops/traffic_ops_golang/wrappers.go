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
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/tocookie"
	"net/http"
	"time"
)

const ServerName = "traffic_ops_golang" + "/" + Version

func wrapHeaders(h RegexHandlerFunc) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Server-Name", ServerName)
		h(w, r, p)
	}
}

func wrapAuth(h RegexHandlerFunc, noAuth bool, secret string, privLevelStmt *sql.Stmt, privLevelRequired int) RegexHandlerFunc {
	if noAuth {
		return h
	}
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		// TODO remove, and make username available to wrapLogTime
		start := time.Now()
		iw := &Interceptor{w: w}
		w = iw
		username := "-"
		defer func() {
			log.Infof(`%s - %s [%s] "%v %v HTTP/1.1" %v 0 0 "%v"\n`, r.RemoteAddr, username, time.Now().Format(AccessLogTimeFormat), r.Method, r.URL.Path, iw.code, time.Now().Sub(start)/time.Millisecond, iw.byteCount, r.UserAgent())
		}()

		handleUnauthorized := func(reason string) {
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

		username = oldCookie.AuthData
		if !hasPrivLevel(privLevelStmt, username, privLevelRequired) {
			handleUnauthorized("insufficient privileges")
			return
		}

		newCookieVal := tocookie.Refresh(oldCookie, secret)
		http.SetCookie(w, &http.Cookie{Name: tocookie.Name, Value: newCookieVal})

		h(w, r, p)
	}
}

const AccessLogTimeFormat = "02/Jan/2006:15:04:05 -0700"

func wrapLogTime(h RegexHandlerFunc) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		start := time.Now()
		iw := &Interceptor{w: w}
		defer func() {
			user := "-" // TODO fix
			log.Infof(`%s - %s [%s] "%v %v HTTP/1.1" %v 0 0 "%v"\n`, r.RemoteAddr, user, time.Now().Format(AccessLogTimeFormat), r.Method, r.URL.Path, iw.code, time.Now().Sub(start)/time.Millisecond, iw.byteCount, r.UserAgent())
		}()
		h(iw, r, p)
	}
}

type Interceptor struct {
	w         http.ResponseWriter
	code      int
	byteCount int
}

func (i *Interceptor) WriteHeader(rc int) {
	i.w.WriteHeader(rc)
	i.code = rc
}

func (i *Interceptor) Write(b []byte) (int, error) {
	wi, werr := i.w.Write(b)
	i.byteCount += wi
	return wi, werr
}

func (i *Interceptor) Header() http.Header {
	return i.w.Header()
}
