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
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/tocookie"
	"net/http"
	"time"
)

const ServerName = "traffic_ops_golang" + "/" + Version

func wrapHeaders(h RegexHandlerFunc) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Server-Name", ServerName)
		iw := &BodyInterceptor{w: w}
		w = iw
		h(w, r, p)
		sha := sha512.Sum512(iw.body)
		w.Header().Set("Whole-Content-SHA512", base64.StdEncoding.EncodeToString(sha[:]))
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
			log.EventfRaw(`%s - %s [%s] "%v %v HTTP/1.1" %v %v %v "%v"`, r.RemoteAddr, username, time.Now().Format(AccessLogTimeFormat), r.Method, r.URL.Path, iw.code, iw.byteCount, int(time.Now().Sub(start)/time.Millisecond), r.UserAgent())
		}()

		handleUnauthorized := func(reason string) {
			status := http.StatusUnauthorized
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
			log.Infof("%v %v %v %v returned unauthorized: %v\n", r.RemoteAddr, r.Method, r.URL.Path, username, reason)
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
		http.SetCookie(w, &http.Cookie{Name: tocookie.Name, Value: newCookieVal, Path: "/", HttpOnly: true})

		h(w, r, p)
	}
}

const AccessLogTimeFormat = "02/Jan/2006:15:04:05 -0700"

func wrapAccessLog(secret string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iw := &Interceptor{w: w}
		user := "-"
		cookie, err := r.Cookie(tocookie.Name)
		if err == nil && cookie != nil {
			cookie, err := tocookie.Parse(secret, cookie.Value)
			if err == nil {
				user = cookie.AuthData
			}
		}
		start := time.Now()
		defer func() {
			log.EventfRaw(`%s - %s [%s] "%v %v HTTP/1.1" %v %v %v "%v"`, r.RemoteAddr, user, time.Now().Format(AccessLogTimeFormat), r.Method, r.URL.Path, iw.code, iw.byteCount, int(time.Now().Sub(start)/time.Millisecond), r.UserAgent())
		}()
		h.ServeHTTP(iw, r)
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
	if i.code == 0 {
		i.code = 200
	}
	return wi, werr
}

func (i *Interceptor) Header() http.Header {
	return i.w.Header()
}

type BodyInterceptor struct {
	w    http.ResponseWriter
	body []byte
}

func (i *BodyInterceptor) WriteHeader(rc int) {
	i.w.WriteHeader(rc)
}

func (i *BodyInterceptor) Write(b []byte) (int, error) {
	i.body = append(i.body, b...)
	wi, werr := i.w.Write(b)
	return wi, werr
}

func (i *BodyInterceptor) Header() http.Header {
	return i.w.Header()
}
