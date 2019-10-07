package routing

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
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/apache/trafficcontrol/lib/go-log"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/about"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tocookie"
)

// ServerName - the server identifier
var ServerName = "traffic_ops_golang" + "/" + about.About.Version

// AuthBase ...
type AuthBase struct {
	secret   string
	override Middleware
}

// GetWrapper ...
func (a AuthBase) GetWrapper(privLevelRequired int) Middleware {
	if a.override != nil {
		return a.override
	}
	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, userErr, sysErr, errCode := api.GetUserFromReq(w, r, a.secret)
			if userErr != nil || sysErr != nil {
				api.HandleErr(w, r, nil, errCode, userErr, sysErr)
				return
			}
			if user.PrivLevel < privLevelRequired {
				api.HandleErr(w, r, nil, http.StatusForbidden, errors.New("Forbidden."), nil)
				return
			}
			api.AddUserToReq(r, user)
			handlerFunc(w, r)
		}
	}
}

func timeOutWrapper(timeout time.Duration) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			http.TimeoutHandler(h, timeout, "server timed out").ServeHTTP(w, r)
		}
	}
}

func wrapHeaders(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Server-Name", ServerName)
		iw := &BodyInterceptor{w: w}
		h(iw, r)

		sha := sha512.Sum512(iw.Body())
		w.Header().Set("Whole-Content-SHA512", base64.StdEncoding.EncodeToString(sha[:]))

		gzipResponse(w, r, iw.Body())

	}
}

func wrapPanicRecover(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				api.HandleErr(w, r, nil, http.StatusInternalServerError, nil, fmt.Errorf("panic: (err: %v) stacktrace:\n%s\n", err, stacktrace()))
				return
			}
		}()
		h(w, r)
	}
}

func stacktrace() []byte {
	initialBufSize := 1024
	buf := make([]byte, initialBufSize)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, len(buf)*2)
	}
}

// AccessLogTimeFormat ...
const AccessLogTimeFormat = "02/Jan/2006:15:04:05 -0700"

func getWrapAccessLog(secret string) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return wrapAccessLog(secret, h)
	}
}

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
			log.EventfRaw(`%s - %s [%s] "%v %v?%v %s" %v %v %v "%v"`, r.RemoteAddr, user, time.Now().Format(AccessLogTimeFormat), r.Method, r.URL.Path, r.URL.RawQuery, r.Proto, iw.code, iw.byteCount, int(time.Now().Sub(start)/time.Millisecond), r.UserAgent())
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

	ctx := r.Context()
	val := ctx.Value(tc.StatusKey)
	status, ok := val.(int)
	if ok { //if not we assume it is a 200
		w.WriteHeader(status)
	}

	w.Write(bytes)
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
		return nil, fmt.Errorf("gzipping bytes: %v", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("closing gzip writer: %v", err)
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

// Interceptor ...
type Interceptor struct {
	w         http.ResponseWriter
	code      int
	byteCount int
}

// WriteHeader ...
func (i *Interceptor) WriteHeader(rc int) {
	i.w.WriteHeader(rc)
	i.code = rc
}

// Write ...
func (i *Interceptor) Write(b []byte) (int, error) {
	wi, werr := i.w.Write(b)
	i.byteCount += wi
	if i.code == 0 {
		i.code = 200
	}
	return wi, werr
}

// Header ...
func (i *Interceptor) Header() http.Header {
	return i.w.Header()
}

// BodyInterceptor fulfills the Writer interface, but records the body and doesn't actually write. This allows performing operations on the entire body written by a handler, for example, compressing or hashing. To actually write, call `RealWrite()`. Note this means `len(b)` and `nil` are always returned by `Write()`, any real write errors will be returned by `RealWrite()`.
type BodyInterceptor struct {
	w    http.ResponseWriter
	body []byte
}

// WriteHeader ...
func (i *BodyInterceptor) WriteHeader(rc int) {
	i.w.WriteHeader(rc)
}

// Write ...
func (i *BodyInterceptor) Write(b []byte) (int, error) {
	i.body = append(i.body, b...)
	return len(b), nil
}

// Header ...
func (i *BodyInterceptor) Header() http.Header {
	return i.w.Header()
}

// RealWrite ...
func (i *BodyInterceptor) RealWrite(b []byte) (int, error) {
	wi, werr := i.w.Write(i.body)
	return wi, werr
}

// Body ...
func (i *BodyInterceptor) Body() []byte {
	return i.body
}

func NotImplementedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"alerts":[{"level":"error","text":"The requested api version is not implemented by this server. If you are using a newer client with an older server, you will need to use an older client version or upgrade your server."}]}`))
	})
}
