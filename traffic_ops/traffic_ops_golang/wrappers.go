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
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/tocookie"
)

const ServerName = "traffic_ops_golang" + "/" + Version

type AuthRegexHandlerFunc func(w http.ResponseWriter, r *http.Request, params PathParams, user string, privLevel int)

func wrapHeaders(h RegexHandlerFunc) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Server-Name", ServerName)
		iw := &BodyInterceptor{w: w}
		h(iw, r, p)

		sha := sha512.Sum512(iw.Body())
		w.Header().Set("Whole-Content-SHA512", base64.StdEncoding.EncodeToString(sha[:]))

		gzipResponse(w, r, iw.Body())

	}
}

func handlerToAuthHandler(h RegexHandlerFunc) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, user string, privLevel int) { h(w, r, p) }
}

func wrapAuth(h RegexHandlerFunc, noAuth bool, secret string, privLevelStmt *sql.Stmt, privLevelRequired int) RegexHandlerFunc {
	return wrapAuthWithData(handlerToAuthHandler(h), noAuth, secret, privLevelStmt, privLevelRequired)
}

func wrapAuthWithData(h AuthRegexHandlerFunc, noAuth bool, secret string, privLevelStmt *sql.Stmt, privLevelRequired int) RegexHandlerFunc {
	if noAuth {
		return func(w http.ResponseWriter, r *http.Request, p PathParams) {
			h(w, r, p, "", PrivLevelInvalid)
		}
	}
	return func(w http.ResponseWriter, r *http.Request, p PathParams) {
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
		privLevel := PrivLevel(privLevelStmt, username)
		if privLevel < privLevelRequired {
			handleUnauthorized("insufficient privileges")
			return
		}

		newCookieVal := tocookie.Refresh(oldCookie, secret)
		http.SetCookie(w, &http.Cookie{Name: tocookie.Name, Value: newCookieVal, Path: "/", HttpOnly: true})

		h(w, r, p, username, privLevel)
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

// gzipResponse takes a function which cannot error and returns only bytes, and wraps it as a http.HandlerFunc. The errContext is logged if the write fails, and should be enough information to trace the problem (function name, endpoint, request parameters, etc).
func gzipResponse(w http.ResponseWriter, r *http.Request, bytes []byte) {

	bytes, err := gzipIfAccepts(r, w, bytes)
	if err != nil {
		log.Errorf("gzipping request '%v': %v\n", r.URL.EscapedPath(), err)
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		if _, err := w.Write([]byte(http.StatusText(code))); err != nil {
			log.Warnf("received error writing data request %v: %v\n", r.URL.EscapedPath(), err)
		}
		return
	}

	w.Write(bytes)
}

// wrapBytes takes a function which cannot error and returns only bytes, and wraps it as a http.HandlerFunc. The errContext is logged if the write fails, and should be enough information to trace the problem (function name, endpoint, request parameters, etc).
//TODO: drichardson - refactor these to a generic area
func wrapBytes(f func() []byte, contentType string) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams) {
		bytes := f()
		bytes, err := gzipIfAccepts(r, w, bytes)
		if err != nil {
			log.Errorf("gzipping request '%v': %v\n", r.URL.EscapedPath(), err)
			code := http.StatusInternalServerError
			w.WriteHeader(code)
			if _, err := w.Write([]byte(http.StatusText(code))); err != nil {
				log.Warnf("received error writing data request %v: %v\n", r.URL.EscapedPath(), err)
			}
			return
		}

		w.Header().Set(tc.ContentType, contentType)
		log.Write(w, bytes, r.URL.EscapedPath())
	}
}

// gzipIfAccepts gzips the given bytes, writes a `Content-Encoding: gzip` header to the given writer, and returns the gzipped bytes, if the Request supports GZip (has an Accept-Encoding header). Else, returns the bytes unmodified. Note the given bytes are NOT written to the given writer. It is assumed the bytes may need to pass thru other middleware before being written.
//TODO: drichardson - refactor these to a generic area
func gzipIfAccepts(r *http.Request, w http.ResponseWriter, b []byte) ([]byte, error) {
	// TODO this could be made more efficient by wrapping ResponseWriter with the GzipWriter, and letting callers writer directly to it - but then we'd have to deal with Closing the gzip.Writer.
	if len(b) == 0 || !acceptsGzip(r) {
		return b, nil
	}
	w.Header().Set(tc.ContentEncoding, tc.Gzip)

	buf := bytes.Buffer{}
	zw := gzip.NewWriter(&buf)

	if _, err := zw.Write(b); err != nil {
		return nil, fmt.Errorf("gzipping bytes: %v")
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("closing gzip writer: %v")
	}

	return buf.Bytes(), nil
}

func acceptsGzip(r *http.Request) bool {
	encodingHeaders := r.Header["Accept-Encoding"] // headers are case-insensitive, but Go promises to Canonical-Case requests
	for _, encodingHeader := range encodingHeaders {
		encodingHeader = stripAllWhitespace(encodingHeader)
		encodings := strings.Split(encodingHeader, ",")
		for _, encoding := range encodings {
			if strings.ToLower(encoding) == tc.Gzip { // encoding is case-insensitive, per the RFC
				return true
			}
		}
	}
	return false
}

func stripAllWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
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

// BodyInterceptor fulfills the Writer interface, but records the body and doesn't actually write. This allows performing operations on the entire body written by a handler, for example, compressing or hashing. To actually write, call `RealWrite()`. Note this means `len(b)` and `nil` are always returned by `Write()`, any real write errors will be returned by `RealWrite()`.
type BodyInterceptor struct {
	w    http.ResponseWriter
	body []byte
}

func (i *BodyInterceptor) WriteHeader(rc int) {
	i.w.WriteHeader(rc)
}
func (i *BodyInterceptor) Write(b []byte) (int, error) {
	i.body = append(i.body, b...)
	return len(b), nil
}
func (i *BodyInterceptor) Header() http.Header {
	return i.w.Header()
}
func (i *BodyInterceptor) RealWrite(b []byte) (int, error) {
	wi, werr := i.w.Write(i.body)
	return wi, werr
}
func (i *BodyInterceptor) Body() []byte {
	return i.body
}
