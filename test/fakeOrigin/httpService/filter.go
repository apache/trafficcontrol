package httpService

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
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// BodyInterceptor is a container for a http writer body so we can write headers after the real writes are issued
type BodyInterceptor struct {
	w            http.ResponseWriter
	body         bytes.Buffer
	responseCode int
}

// WriteHeader doesn't actually write a header, just a response code.  Thanks go.
func (i *BodyInterceptor) WriteHeader(rc int) {
	i.responseCode = rc
}

// Write is used for interface compatability with http response writer
func (i *BodyInterceptor) Write(b []byte) (int, error) {
	i.body.Write(b)
	return i.body.Len(), nil
}

// Header is used for interface compatability with http response writer
func (i *BodyInterceptor) Header() http.Header {
	return i.w.Header()
}

// RealWrite is called to perform the actual write operation already done further in the chain
func (i *BodyInterceptor) RealWrite() (int, error) {
	if i.responseCode != 0 {
		i.w.WriteHeader(i.responseCode)
	}
	c := i.body.Len()
	io.Copy(i.w, &i.body)
	return c, nil
}

// Body is used for interface compatability with http response writer
func (i *BodyInterceptor) Body() []byte {
	return i.body.Bytes()
}

// ParseHTTPDate parses the given RFC7231§7.1.1 HTTP-date
func ParseHTTPDate(d string) (time.Time, bool) {
	if t, err := time.Parse(time.RFC1123, d); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.RFC850, d); err == nil {
		return t, true
	}
	if t, err := time.Parse(time.ANSIC, d); err == nil {
		return t, true
	}
	return time.Time{}, false
}

func AddFullDefaultHeader(w http.ResponseWriter, r *http.Request, newKey string, newVals []string) bool {
	if w.Header().Get(newKey) != "" || len(newVals) == 0 {
		return false
	}
	w.Header().Set(newKey, newVals[0])
	for _, val := range newVals[1:] {
		w.Header().Add(newKey, val)
	}
	return true
}

func GenerateETag(source string) string {
	h := fnv.New32a()
	h.Write([]byte(source))
	return fmt.Sprintf("\"%d\"", h.Sum32())
}

func logfo(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		iw := &BodyInterceptor{w: w}
		handler.ServeHTTP(iw, r)
		size := iw.body.Len()
		rc := iw.responseCode
		if rc == 0 {
			rc = http.StatusOK
		}
		iw.RealWrite()
		finishTime := time.Now()
		remoteAddr := r.RemoteAddr
		finishTS := finishTime.Format(time.RFC1123)
		method := r.Method
		rURI := r.URL.EscapedPath()
		proto := r.Proto
		dur := finishTime.Sub(startTime)
		refer := strings.Replace(r.Referer(), `"`, `\"`, -1)
		uas := strings.Replace(r.UserAgent(), `"`, `\"`, -1)
		im := strings.Replace(r.Header.Get("If-Match"), `"`, `\"`, -1)
		inm := strings.Replace(r.Header.Get("If-None-Match"), `"`, `\"`, -1)
		ims := strings.Replace(r.Header.Get("If-Modified-Since"), `"`, `\"`, -1)
		ius := strings.Replace(r.Header.Get("If-Unmodified-Since"), `"`, `\"`, -1)
		ir := strings.Replace(r.Header.Get("If-Range"), `"`, `\"`, -1)
		fmt.Println("hit in log handler")
		fmt.Printf("%s - [%s] \"%s %s %s\" %d %d %d \"%s\" \"%s\" \"%s\" \"%s\" \"%s\" \"%s\" \"%s\"\n", remoteAddr, finishTS, method, rURI, proto, rc, size, dur, refer, uas, im, inm, ims, ius, ir)
	})
}

func strictTransportSecurity(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := &BodyInterceptor{w: w}
		handler.ServeHTTP(iw, r)
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		iw.RealWrite()
	})
}
func originHeaderManipulation(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := &BodyInterceptor{w: w}
		handler.ServeHTTP(iw, r)
		for key, value := range r.Header {
			if strings.HasPrefix(key, "Fakeorigin-") {
				newKey := strings.TrimPrefix(key, "Fakeorigin-")
				// Intentionally stomp on anything else you may otherwise have known
				w.Header().Set(newKey, value[0])
				for _, val := range value[1:len(value)] {
					w.Header().Add(newKey, val)
				}
			}
		}
		iw.RealWrite()
	})
}

func checkIfMatch(req, origineTag, diskpath string) bool {
	if req == "*" {
		if _, err := os.Stat(diskpath); err == nil {
			return true
		}
	}
	etags := strings.Split(req, ",")
	for _, etag := range etags {
		if strings.TrimSpace(etag) == origineTag {
			return true
		}
	}
	return false
}

func checkIfNoneMatch(req, origineTag, diskpath string) bool {
	if req == "*" {
		if _, err := os.Stat(diskpath); err == nil {
			return false
		}
	}
	etags := strings.Split(req, ",")
	for _, etag := range etags {
		if strings.TrimPrefix(strings.TrimSpace(etag), "W/") == origineTag {
			return false
		}
	}
	return true
}

func checkIfRange(req, origineTag string, lastUpdated time.Time) bool {
	if req == "" {
		return true
	}
	reqTime, timeOk := ParseHTTPDate(req)
	if !timeOk && strings.TrimPrefix(strings.TrimSpace(req), "W/") == origineTag {
		return true
	}
	if timeOk && reqTime == lastUpdated {
		return true
	}
	return false
}

func checkIsFullRange(rr string, size int) bool {
	if rr == "" {
		return false
	}
	nospaces := strings.Replace(rr, " ", "", -1)
	if strings.HasSuffix(nospaces, "0-") {
		return true
	}
	i := strconv.Itoa(size - 1)
	if strings.HasSuffix(nospaces, "0-"+i) {
		return true
	}

	return false
}

// https://tools.ietf.org/html/rfc7232 - If-Match, If-None-Match, If-Modified-Since, If-Unmodified-Since, If-Range
// https://tools.ietf.org/html/rfc7233 - Range Requests, If-Range
func cacheOptimization(handler http.Handler, startTime time.Time, ep httpEndpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		diskpath := ep.OutputDir + "/" + strings.TrimPrefix(path.Base(r.RequestURI), "/"+ep.ID)
		eTag := GenerateETag(ep.OutputDir + ep.DiskID + ep.LastTranscodeTime.Format(time.RFC1123) + path.Base(r.RequestURI))
		w.Header().Set("ETag", eTag)
		w.Header().Set("Last-Modified", ep.LastTranscodeTime.Format(time.RFC1123))
		if r.Header.Get("If-Match") != "" && !checkIfMatch(r.Header.Get("If-Match"), eTag, diskpath) {
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Write([]byte(strconv.Itoa(http.StatusPreconditionFailed) + " " + http.StatusText(http.StatusPreconditionFailed)))
		}
		if r.Header.Get("If-None-Match") != "" && !checkIfNoneMatch(r.Header.Get("If-None-Match"), eTag, diskpath) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				w.WriteHeader(http.StatusNotModified)
				w.Write([]byte(strconv.Itoa(http.StatusNotModified) + " " + http.StatusText(http.StatusNotModified)))
			} else {
				w.WriteHeader(http.StatusPreconditionFailed)
				w.Write([]byte(strconv.Itoa(http.StatusPreconditionFailed) + " " + http.StatusText(http.StatusPreconditionFailed)))
			}
		}
		imsTime, timeOk := ParseHTTPDate(r.Header.Get("If-Modified-Since"))
		if r.Header.Get("If-None-Match") == "" && (r.Method == http.MethodGet || r.Method == http.MethodHead) && timeOk {
			if imsTime.Before(ep.LastTranscodeTime) || imsTime == ep.LastTranscodeTime {
				w.WriteHeader(http.StatusNotModified)
				w.Write([]byte(strconv.Itoa(http.StatusNotModified) + " " + http.StatusText(http.StatusNotModified)))
			}
		}
		iusTime, timeOk := ParseHTTPDate(r.Header.Get("If-Unmodified-Since"))
		if r.Header.Get("If-Match") == "" && (r.Method == http.MethodGet || r.Method == http.MethodHead) && timeOk {
			if ep.LastTranscodeTime.After(iusTime) {
				w.WriteHeader(http.StatusPreconditionFailed)
				w.Write([]byte(strconv.Itoa(http.StatusPreconditionFailed) + " " + http.StatusText(http.StatusPreconditionFailed)))
			}
		}
		iw := &BodyInterceptor{w: w}
		handler.ServeHTTP(iw, r)
		w.Header().Set("Accept-Ranges", "bytes")
		rrange := r.Header.Get("Range")
		irrange := r.Header.Get("If-Range")
		// TODO: ensure this doesn't trigger on anything but 200
		if rrange != "" && checkIfRange(irrange, eTag, ep.LastTranscodeTime) && !checkIsFullRange(rrange, len(iw.Body())) {
			// Generate a 206 Paritial Content Range Request
			if ranges, err := parseRange(rrange, uint64(len(iw.Body()))); err != nil {
				iw.body.Reset()
				iw.responseCode = http.StatusRequestedRangeNotSatisfiable
			} else {
				b, headers, err := clipToRange(ranges, iw.body.Bytes(), w.Header().Get("Content-Type"))
				if err != nil {
					iw.body.Reset()
					iw.responseCode = http.StatusRequestedRangeNotSatisfiable
				} else {
					AddFullDefaultHeader(w, r, "Cache-Control", []string{})
					AddFullDefaultHeader(w, r, "Expires", []string{(time.Now().Add(time.Minute * time.Duration(10))).Format(time.RFC1123)})
					AddFullDefaultHeader(w, r, "Content-Location", []string{})
					AddFullDefaultHeader(w, r, "Vary", []string{})
					iw.responseCode = http.StatusPartialContent
					// Reset the body in the body interceptor chain since we're explicitly clipping it to the requested ranges, otherwise it's just appending
					iw.body.Reset()
					iw.body.Write(b)
					for key, val := range headers {
						w.Header().Set(key, val)
					}
				}
			}
		}
		iw.RealWrite()
	})
}
