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

package datareq

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/health"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/towrap"
)

// MakeDispatchMap returns the map of paths to http.HandlerFuncs for dispatching.
func MakeDispatchMap(
	opsConfig threadsafe.OpsConfig,
	toSession towrap.TrafficOpsSessionThreadsafe,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	distributedPeerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	statInfoHistory threadsafe.ResultInfoHistory,
	statResultHistory threadsafe.ResultStatHistory,
	statMaxKbpses threadsafe.CacheKbpses,
	healthHistory threadsafe.ResultHistory,
	dsStats threadsafe.DSStatsReader,
	events health.ThreadsafeEvents,
	staticAppData config.StaticAppData,
	healthPollInterval time.Duration,
	lastHealthDurations threadsafe.DurationMap,
	fetchCount threadsafe.Uint,
	healthIteration threadsafe.Uint,
	errorCount threadsafe.Uint,
	toData todata.TODataThreadsafe,
	localCacheStatus threadsafe.CacheAvailableStatus,
	lastStats threadsafe.LastStats,
	statUnpolledCaches threadsafe.UnpolledCaches,
	healthUnpolledCaches threadsafe.UnpolledCaches,
	monitorConfig threadsafe.TrafficMonitorConfigMap,
	statPollingEnabled bool,
	distributedPollingEnabled bool,
) map[string]http.HandlerFunc {

	// wrap composes all universal wrapper functions. Right now, it's only the UnpolledCheck, but there may be others later. For example, security headers.
	wrap := func(f http.HandlerFunc) http.HandlerFunc {
		if statPollingEnabled {
			return wrapUnpolledCheck(statUnpolledCaches, errorCount, f)
		} else {
			return wrapUnpolledCheck(healthUnpolledCaches, errorCount, f)
		}
	}

	dispatchMap := map[string]http.HandlerFunc{
		"/publish/CrConfig": wrap(WrapAgeErr(errorCount, func() ([]byte, time.Time, error) {
			return srvTRConfig(opsConfig, toSession)
		}, rfc.ApplicationJSON)),
		"/publish/CrStates": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			bytes, statusCode, err := srvTRState(params, localStates, combinedStates, peerStates, distributedPollingEnabled)
			return WrapErrStatusCode(errorCount, path, bytes, statusCode, err)
		}, rfc.ApplicationJSON)),
		"/publish/CacheStatsNew": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvCacheStats(params, errorCount, path, toData, statResultHistory, statInfoHistory, monitorConfig, combinedStates, statMaxKbpses)
		}, rfc.ApplicationJSON)),
		"/publish/CacheStats": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvLegacyCacheStats(params, errorCount, path, toData, statResultHistory, statInfoHistory, monitorConfig, combinedStates, statMaxKbpses)
		}, rfc.ApplicationJSON)),
		"/publish/DsStats": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvDSStats(params, errorCount, path, toData, dsStats)
		}, rfc.ApplicationJSON)),
		"/publish/EventLog": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvEventLog(events)
		}, rfc.ApplicationJSON)),
		"/publish/PeerStates": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvPeerStates(params, errorCount, path, toData, peerStates)
		}, rfc.ApplicationJSON)),
		"/publish/DistributedPeerStates": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvPeerStates(params, errorCount, path, toData, distributedPeerStates)
		}, rfc.ApplicationJSON)),
		"/publish/Stats": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvStats(staticAppData, healthPollInterval, lastHealthDurations, fetchCount, healthIteration, errorCount, peerStates)
		}, rfc.ApplicationJSON)),
		"/publish/ConfigDoc": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvConfigDoc(opsConfig)
		}, rfc.ApplicationJSON)),
		"/publish/StatSummary": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvStatSummary(params, errorCount, path, toData, statResultHistory)
		}, rfc.ApplicationJSON)),
		"/api/cache-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheCount(localStates)
		}, rfc.ApplicationJSON)),
		"/api/cache-available-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheAvailableCount(localStates)
		}, rfc.ApplicationJSON)),
		"/api/cache-down-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheDownCount(localStates, monitorConfig)
		}, rfc.ApplicationJSON)),
		"/api/version": wrap(WrapBytes(func() []byte {
			return srvAPIVersion(staticAppData)
		}, rfc.ContentTypeTextPlain)),
		"/api/traffic-ops-uri": wrap(WrapBytes(func() []byte {
			return srvAPITrafficOpsURI(opsConfig)
		}, rfc.ContentTypeURIList)),
		"/api/cache-statuses": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvAPICacheStates(toData, statInfoHistory, statResultHistory, healthHistory, lastHealthDurations, localCacheStatus, statMaxKbpses, monitorConfig)
		}, rfc.ApplicationJSON)),
		"/api/bandwidth-kbps": wrap(WrapBytes(func() []byte {
			return srvAPIBandwidthKbps(toData, lastStats)
		}, rfc.ApplicationJSON)),
		"/api/bandwidth-capacity-kbps": wrap(WrapBytes(func() []byte {
			return srvAPIBandwidthCapacityKbps(statMaxKbpses)
		}, rfc.ApplicationJSON)),
		"/api/monitor-config": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvMonitorConfig(monitorConfig)
		}, rfc.ApplicationJSON)),
		"/api/crconfig-history": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvAPICRConfigHist(toSession)
		}, rfc.ApplicationJSON)),
	}
	return addTrailingSlashEndpoints(dispatchMap)
}

// This is the "spirit" of how TM1.0 works; hack to extract a path argument to filter data (/publish/SomeEndpoint/:argument).
func getPathArgument(path string) string {
	pathParts := strings.Split(path, "/")
	if len(pathParts) >= 4 {
		return pathParts[3]
	}

	return ""
}

// HandleErr takes an error, and the request type it came from, and logs. It is ok to call with a nil error, in which case this is a no-op.
func HandleErr(errorCount threadsafe.Uint, reqPath string, err error) {
	if err == nil {
		return
	}
	errorCount.Inc()
	log.Errorf("Request Error: %v\n", fmt.Errorf(reqPath+": %v", err))
}

// WrapErrStatusCode takes the body, err, status code, and log context (errorCount, reqPath). It logs and deals with any error, and returns the appropriate bytes and response code for the `srvhttp`. It notably returns InternalServerError status on any error, for security reasons.
func WrapErrStatusCode(errorCount threadsafe.Uint, reqPath string, body []byte, statusCode int, err error) ([]byte, int) {
	if err == nil {
		return body, http.StatusOK
	}

	HandleErr(errorCount, reqPath, err)

	code := http.StatusInternalServerError

	if statusCode > 0 {
		code = statusCode
	}

	return nil, code
}

// WrapErrCode calls the WrapErrStatusCode function with a hardcoded 500. This is a convenience function for callers that do not want to provide a status code.
func WrapErrCode(errorCount threadsafe.Uint, reqPath string, body []byte, err error) ([]byte, int) {
	return WrapErrStatusCode(errorCount, reqPath, body, http.StatusInternalServerError, err)
}

// WrapBytes takes a function which cannot error and returns only bytes, and wraps it as a http.HandlerFunc. The errContext is logged if the write fails, and should be enough information to trace the problem (function name, endpoint, request parameters, etc).
func WrapBytes(f func() []byte, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		w.Header().Set("Content-Type", contentType)
		log.Write(w, bytes, r.URL.EscapedPath())
	}
}

// WrapErr takes a function which returns bytes and an error, and wraps it as a http.HandlerFunc. If the error is nil, the bytes are written with Status OK. Else, the error is logged, and InternalServerError is returned as the response code. If you need to return a different response code (for example, StatusBadRequest), call wrapRespCode.
func WrapErr(errorCount threadsafe.Uint, f func() ([]byte, error), contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := f()
		code := http.StatusOK
		if err != nil {
			bytes, code = WrapErrCode(errorCount, r.URL.EscapedPath(), bytes, err)
		} else {
			bytes, err = gzipIfAccepts(r, w, bytes)
			bytes, code = WrapErrCode(errorCount, r.URL.EscapedPath(), bytes, err)
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(code)
		log.Write(w, bytes, r.URL.EscapedPath())
	}
}

// SrvFunc is a function which takes URL parameters, and returns the requested data, and a response code. Note it does not take the full http.Request, and does not have the path. SrvFunc functions should be called via dispatch, and any additional data needed should be closed via a lambda.
// TODO split params and path into 2 separate wrappers?
// TODO change to simply take the http.Request?
type SrvFunc func(params url.Values, path string) ([]byte, int)

// WrapParams takes a SrvFunc and wraps it as an http.HandlerFunc. Note if the SrvFunc returns 0 bytes, an InternalServerError is returned, and the response code is ignored, for security reasons. If an error response code is necessary, return bytes to that effect, for example, "Bad Request". DO NOT return informational messages regarding internal server errors; these should be logged, and only a 500 code returned to the client, for security reasons.
func WrapParams(f SrvFunc, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusInternalServerError
		bytes, statusCode := f(r.URL.Query(), r.URL.EscapedPath())

		if statusCode > 0 {
			code = statusCode
		}

		if len(bytes) == 0 {
			w.WriteHeader(code)
			if _, err := w.Write([]byte(http.StatusText(code))); err != nil {
				log.Warnf("received error writing data request %v: %v\n", r.URL.EscapedPath(), err)
			}
			return
		}

		bytes, err := gzipIfAccepts(r, w, bytes)

		if err != nil {
			log.Errorf("gzipping '%v': %v\n", r.URL.EscapedPath(), err)
			w.WriteHeader(code)
			if _, err := w.Write([]byte(http.StatusText(code))); err != nil {
				log.Warnf("received error writing data request %v: %v\n", r.URL.EscapedPath(), err)
			}
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(code)
		if _, err := w.Write(bytes); err != nil {
			log.Warnf("received error writing data request %v: %v\n", r.URL.EscapedPath(), err)
		}
	}
}

func WrapAgeErr(errorCount threadsafe.Uint, f func() ([]byte, time.Time, error), contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, contentTime, err := f()
		code := http.StatusOK
		if err != nil {
			bytes, code = WrapErrCode(errorCount, r.URL.EscapedPath(), bytes, err)
		} else {
			bytes, err = gzipIfAccepts(r, w, bytes)
			bytes, code = WrapErrCode(errorCount, r.URL.EscapedPath(), bytes, err)
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Age", fmt.Sprintf("%.0f", time.Since(contentTime).Seconds()))
		w.WriteHeader(code)

		log.Write(w, bytes, r.URL.EscapedPath())
	}
}

func accessLogTime(t time.Time) float64 {
	return float64(t.UnixMilli()) / 1000.0
}

func accessLogStr(
	timestamp time.Time, // prefix
	remoteAddress string, // chi
	reqMethod string, // cqhm
	reqPath string, // url
	reqRawQuery string,
	statusCode int, // pssc
	respSize int, // b
	reqServeTimeMs int, // ttms
	userAgent string, // uas
) string {
	return fmt.Sprintf("%.3f chi=%s cqhm=%s url=\"%s?%s\" pssc=%d b=%d ttms=%d uas=\"%s\"",
		accessLogTime(timestamp),
		remoteAddress,
		reqMethod,
		reqPath,
		reqRawQuery,
		statusCode,
		respSize,
		reqServeTimeMs,
		userAgent)
}

// WrapUnpolledCheck wraps an http.HandlerFunc, returning ServiceUnavailable if all caches have't been polled; else, calling the wrapped func. Once all caches have been polled, we never return a 503 again, even if the CRConfig has been changed and new, unpolled caches exist. This is because, before those new caches existed in the CRConfig, they weren't being routed to, so it doesn't break anything to continue not routing to them until they're polled, while still serving polled caches as available. Whereas, on startup, if we were to return data with some caches unpolled, we would be telling clients that existing, potentially-available caches are unavailable, simply because we hadn't polled them yet.
func wrapUnpolledCheck(unpolledCaches threadsafe.UnpolledCaches, errorCount threadsafe.Uint, f http.HandlerFunc) http.HandlerFunc {
	polledAll := false
	polledLocal := false
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		iw := &util.Interceptor{W: w}
		defer func() {
			log.Accessln(accessLogStr(time.Now(), r.RemoteAddr, r.Method, r.URL.Path, r.URL.RawQuery, iw.Code, iw.ByteCount, int(time.Now().Sub(start)/time.Millisecond), r.UserAgent()))
		}()
		if !polledAll || !polledLocal {
			polledAll = !unpolledCaches.Any()
			polledLocal = !unpolledCaches.AnyDirectlyPolled()
			rawOrLocal := r.URL.Query().Has("raw") || r.URL.Query().Has("local")
			if (!rawOrLocal && !polledAll) || (rawOrLocal && !polledLocal) {
				HandleErr(errorCount, r.URL.EscapedPath(), fmt.Errorf("service still starting, some caches unpolled: %v", unpolledCaches.UnpolledCaches()))
				iw.WriteHeader(http.StatusServiceUnavailable)
				log.Write(iw, []byte("Service Unavailable"), r.URL.EscapedPath())
				return
			}
		}
		iw.Header().Set(rfc.PermissionsPolicy, "interest-cohort=()")
		f(iw, r)
	}
}

func stripAllWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// addTrailingEndpoints adds endpoints with trailing slashes to the given dispatch map. Without this, Go will match `route` and `route/` differently.
func addTrailingSlashEndpoints(dispatchMap map[string]http.HandlerFunc) map[string]http.HandlerFunc {
	for route, handler := range dispatchMap {
		if strings.HasSuffix(route, "/") {
			continue
		}
		dispatchMap[route+"/"] = handler
	}
	return dispatchMap
}

func acceptsGzip(r *http.Request) bool {
	encodingHeaders := r.Header["Accept-Encoding"] // headers are case-insensitive, but Go promises to Canonical-Case requests
	for _, encodingHeader := range encodingHeaders {
		encodingHeader = stripAllWhitespace(encodingHeader)
		encodings := strings.Split(encodingHeader, ",")
		for _, encoding := range encodings {
			if strings.ToLower(encoding) == "gzip" { // encoding is case-insensitive, per the RFC
				return true
			}
		}
	}
	return false
}

// gzipIfAccepts gzips the given bytes, writes a `Content-Encoding: gzip` header to the given writer, and returns the gzipped bytes, if the Request supports GZip (has an Accept-Encoding header). Else, returns the bytes unmodified. Note the given bytes are NOT written to the given writer. It is assumed the bytes may need to pass thru other middleware before being written.
func gzipIfAccepts(r *http.Request, w http.ResponseWriter, b []byte) ([]byte, error) {
	// TODO this could be made more efficient by wrapping ResponseWriter with the GzipWriter, and letting callers writer directly to it - but then we'd have to deal with Closing the gzip.Writer.
	if len(b) == 0 || !acceptsGzip(r) {
		return b, nil
	}
	w.Header().Set("Content-Encoding", "gzip")

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
