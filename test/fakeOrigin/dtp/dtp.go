package dtp

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
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type LogRecorder struct {
	http.ResponseWriter

	Status       int
	HeaderBytes  int64
	ContentBytes int64
}

func (rec *LogRecorder) WriteHeader(code int) {
	rec.Status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *LogRecorder) Write(bytes []byte) (int, error) {
	rec.ContentBytes += int64(len(bytes))
	return rec.ResponseWriter.Write(bytes)
}

// this is mostly for hijack
func isHandlerType(r *http.Request) bool {
	if strings.Contains(r.URL.EscapedPath(), "~h.") {
		return true
	} else if strings.Contains(r.URL.RawQuery, "~h.") {
		return true
	} else {
		for _, part := range r.Header[`X-Dtp`] {
			if strings.Contains(part, "~h.") {
				return true
			}
		}
	}

	return false
}

func Logger(alog *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now().UnixNano()

		// the logger interferes with hijacking
		if isHandlerType(r) {
			next.ServeHTTP(w, r)
			url := strings.Replace(r.URL.String(), "\n", "", -1)
			url = strings.Replace(url, "\r", "", -1)
			ua := strings.Replace(r.UserAgent(), "\n", "", -1)
			ua = strings.Replace(ua, "\r", "", -1)
			thisRange := strings.Replace(r.Header.Get("Range"), "\n", "", -1)
			thisRange = strings.Replace(thisRange, "\r", "", -1)
			alog.Printf("%.3f %s \"%s\" %d b=%d ttms=%d uas=\"%s\" rr=\"%s\"\n",
				float64(timeStart)/float64(1.e9),
				r.Method,
				url,
				42, // status code -- why not?
				0,  // bytes
				0,  // ttms
				ua,
				thisRange,
			)
			return
		}

		tlsstr := "-"
		if r.TLS != nil {
			tlsstr = tls.CipherSuiteName(r.TLS.CipherSuite)
		}
		/*
			if nil != r.TLS {
				tlsstr = fmt.Sprintf("0x%x", r.TLS.CipherSuite)
			}
		*/

		rec := LogRecorder{w, 200, 0, 0}
		next.ServeHTTP(&rec, r)
		timeStop := time.Now().UnixNano()
		url := strings.Replace(r.URL.String(), "\n", "", -1)
		url = strings.Replace(url, "\r", "", -1)
		ua := strings.Replace(r.UserAgent(), "\n", "", -1)
		ua = strings.Replace(ua, "\r", "", -1)
		thisRange := strings.Replace(r.Header.Get("Range"), "\n", "", -1)
		thisRange = strings.Replace(thisRange, "\r", "", -1)
		alog.Printf("%.3f %s \"%s\" %s %d b=%d ttms=%d uas=\"%s\" rr=\"%s\"\n",
			float64(timeStart)/float64(1.e9),
			r.Method,
			url,
			tlsstr,
			rec.Status,
			rec.ContentBytes,
			int64((timeStop-timeStart)/1000000),
			ua,
			thisRange,
		)

		if GlobalConfig.Log.RequestHeaders {
			alog.Print(r.Header)
		}
		if GlobalConfig.Log.ResponseHeaders {
			alog.Print(w.Header())
		}
	})
}

func DebugLog(s string) {
	if GlobalConfig.Debug {
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		fmt.Println(s)
	}
}

func DebugLogf(format string, args ...interface{}) {
	if GlobalConfig.Debug {
		fmt.Printf(format, args...)
	}
}
