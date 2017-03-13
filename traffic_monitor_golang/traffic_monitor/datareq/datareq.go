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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/config"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/health"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
)

// MakeDispatchMap returns the map of paths to http.HandlerFuncs for dispatching.
func MakeDispatchMap(
	opsConfig threadsafe.OpsConfig,
	toSession towrap.ITrafficOpsSession,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
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
	unpolledCaches threadsafe.UnpolledCaches,
	monitorConfig threadsafe.TrafficMonitorConfigMap,
) map[string]http.HandlerFunc {

	// wrap composes all universal wrapper functions. Right now, it's only the UnpolledCheck, but there may be others later. For example, security headers.
	wrap := func(f http.HandlerFunc) http.HandlerFunc {
		return wrapUnpolledCheck(unpolledCaches, errorCount, f)
	}

	dispatchMap := map[string]http.HandlerFunc{
		"/publish/CrConfig": wrap(WrapAgeErr(errorCount, func() ([]byte, time.Time, error) {
			return srvTRConfig(opsConfig, toSession)
		}, ContentTypeJSON)),
		"/publish/CrStates": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			bytes, err := srvTRState(params, localStates, combinedStates)
			return WrapErrCode(errorCount, path, bytes, err)
		}, ContentTypeJSON)),
		"/publish/CacheStats": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvCacheStats(params, errorCount, path, toData, statResultHistory, statInfoHistory, monitorConfig, combinedStates, statMaxKbpses)
		}, ContentTypeJSON)),
		"/publish/DsStats": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvDSStats(params, errorCount, path, toData, dsStats)
		}, ContentTypeJSON)),
		"/publish/EventLog": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvEventLog(events)
		}, ContentTypeJSON)),
		"/publish/PeerStates": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvPeerStates(params, errorCount, path, toData, peerStates)
		}, ContentTypeJSON)),
		"/publish/Stats": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvStats(staticAppData, healthPollInterval, lastHealthDurations, fetchCount, healthIteration, errorCount, peerStates)
		}, ContentTypeJSON)),
		"/publish/ConfigDoc": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvConfigDoc(opsConfig)
		}, ContentTypeJSON)),
		"/publish/StatSummary": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvStatSummary(params, errorCount, path, toData, statResultHistory)
		}, ContentTypeJSON)),
		"/api/cache-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheCount(localStates)
		}, ContentTypeJSON)),
		"/api/cache-available-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheAvailableCount(localStates)
		}, ContentTypeJSON)),
		"/api/cache-down-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheDownCount(localStates, monitorConfig)
		}, ContentTypeJSON)),
		"/api/version": wrap(WrapBytes(func() []byte {
			return srvAPIVersion(staticAppData)
		}, ContentTypeJSON)),
		"/api/traffic-ops-uri": wrap(WrapBytes(func() []byte {
			return srvAPITrafficOpsURI(opsConfig)
		}, ContentTypeJSON)),
		"/api/cache-statuses": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvAPICacheStates(toData, statInfoHistory, statResultHistory, healthHistory, lastHealthDurations, localStates, lastStats, localCacheStatus, statMaxKbpses)
		}, ContentTypeJSON)),
		"/api/bandwidth-kbps": wrap(WrapBytes(func() []byte {
			return srvAPIBandwidthKbps(toData, lastStats)
		}, ContentTypeJSON)),
		"/api/bandwidth-capacity-kbps": wrap(WrapBytes(func() []byte {
			return srvAPIBandwidthCapacityKbps(statMaxKbpses)
		}, ContentTypeJSON)),
		"/api/monitor-config": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvMonitorConfig(monitorConfig)
		}, ContentTypeJSON)),
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

// WrapErrCode takes the body, err, and log context (errorCount, reqPath). It logs and deals with any error, and returns the appropriate bytes and response code for the `srvhttp`. It notably returns InternalServerError status on any error, for security reasons.
func WrapErrCode(errorCount threadsafe.Uint, reqPath string, body []byte, err error) ([]byte, int) {
	if err == nil {
		return body, http.StatusOK
	}
	HandleErr(errorCount, reqPath, err)
	return nil, http.StatusInternalServerError
}

// WrapBytes takes a function which cannot error and returns only bytes, and wraps it as a http.HandlerFunc. The errContext is logged if the write fails, and should be enough information to trace the problem (function name, endpoint, request parameters, etc).
func WrapBytes(f func() []byte, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		log.Write(w, f(), r.URL.EscapedPath())
	}
}

// WrapErr takes a function which returns bytes and an error, and wraps it as a http.HandlerFunc. If the error is nil, the bytes are written with Status OK. Else, the error is logged, and InternalServerError is returned as the response code. If you need to return a different response code (for example, StatusBadRequest), call wrapRespCode.
func WrapErr(errorCount threadsafe.Uint, f func() ([]byte, error), contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := f()
		_, code := WrapErrCode(errorCount, r.URL.EscapedPath(), bytes, err)
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
		bytes, code := f(r.URL.Query(), r.URL.EscapedPath())
		if len(bytes) > 0 {
			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(code)
			if _, err := w.Write(bytes); err != nil {
				log.Warnf("received error writing data request %v: %v\n", r.URL.EscapedPath(), err)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte("Internal Server Error")); err != nil {
				log.Warnf("received error writing data request %v: %v\n", r.URL.EscapedPath(), err)
			}
		}
	}
}

func WrapAgeErr(errorCount threadsafe.Uint, f func() ([]byte, time.Time, error), contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, contentTime, err := f()
		_, code := WrapErrCode(errorCount, r.URL.EscapedPath(), bytes, err)
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Age", fmt.Sprintf("%.0f", time.Since(contentTime).Seconds()))
		w.WriteHeader(code)
		log.Write(w, bytes, r.URL.EscapedPath())
	}
}

// WrapUnpolledCheck wraps an http.HandlerFunc, returning ServiceUnavailable if all caches have't been polled; else, calling the wrapped func. Once all caches have been polled, we never return a 503 again, even if the CRConfig has been changed and new, unpolled caches exist. This is because, before those new caches existed in the CRConfig, they weren't being routed to, so it doesn't break anything to continue not routing to them until they're polled, while still serving polled caches as available. Whereas, on startup, if we were to return data with some caches unpolled, we would be telling clients that existing, potentially-available caches are unavailable, simply because we hadn't polled them yet.
func wrapUnpolledCheck(unpolledCaches threadsafe.UnpolledCaches, errorCount threadsafe.Uint, f http.HandlerFunc) http.HandlerFunc {
	polledAll := false
	return func(w http.ResponseWriter, r *http.Request) {
		if !polledAll && unpolledCaches.Any() {
			HandleErr(errorCount, r.URL.EscapedPath(), fmt.Errorf("service still starting, some caches unpolled: %v", unpolledCaches.UnpolledCaches()))
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Write(w, []byte("Service Unavailable"), r.URL.EscapedPath())
			return
		}
		polledAll = true
		f(w, r)
	}
}

const ContentTypeJSON = "application/json"

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
