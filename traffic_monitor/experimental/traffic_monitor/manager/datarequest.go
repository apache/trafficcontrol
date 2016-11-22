package manager

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
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/cache"
	ds "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	dsdata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservicedata"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/srvhttp"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// JSONEvents represents the structure we wish to serialize to JSON, for Events.
type JSONEvents struct {
	Events []Event `json:"events"`
}

// CacheState represents the available state of a cache.
type CacheState struct {
	Value bool `json:"value"`
}

// APIPeerStates contains the data to be returned for an API call to get the peer states of a Traffic Monitor. This contains common API data returned by most endpoints, and a map of peers, to caches' states.
type APIPeerStates struct {
	srvhttp.CommonAPIData
	Peers map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState `json:"peers"`
}

// CacheStatus contains summary stat data about the given cache.
// TODO make fields nullable, so error fields can be omitted, letting API callers still get updates for unerrored fields
type CacheStatus struct {
	Type        *string  `json:"type,omitempty"`
	Status      *string  `json:"status,omitempty"`
	LoadAverage *float64 `json:"load_average,omitempty"`
	// QueryTimeMilliseconds is the time it took this app to perform a complete query and process the data, end-to-end, for the latest health query.
	QueryTimeMilliseconds *int64 `json:"query_time_ms,omitempty"`
	// HealthTimeMilliseconds is the time it took to make the HTTP request and get back the full response, for the latest health query.
	HealthTimeMilliseconds *int64 `json:"health_time_ms,omitempty"`
	// StatTimeMilliseconds is the time it took to make the HTTP request and get back the full response, for the latest stat query.
	StatTimeMilliseconds *int64 `json:"stat_time_ms,omitempty"`
	// StatSpanMilliseconds is the length of time between completing the most recent two stat queries. This can be used as a rough gauge of the end-to-end query processing time.
	StatSpanMilliseconds *int64 `json:"stat_span_ms,omitempty"`
	// HealthSpanMilliseconds is the length of time between completing the most recent two health queries. This can be used as a rough gauge of the end-to-end query processing time.
	HealthSpanMilliseconds *int64   `json:"health_span_ms,omitempty"`
	BandwidthKbps          *float64 `json:"bandwidth_kbps,omitempty"`
	ConnectionCount        *int64   `json:"connection_count,omitempty"`
}

// CacheStatFilter fulfills the cache.Filter interface, for filtering stats. See the `NewCacheStatFilter` documentation for details on which query parameters are used to filter.
type CacheStatFilter struct {
	historyCount int
	statsToUse   map[string]struct{}
	wildcard     bool
	cacheType    enum.CacheType
	hosts        map[enum.CacheName]struct{}
	cacheTypes   map[enum.CacheName]enum.CacheType
}

// UseCache returns whether the given cache is in the filter.
func (f *CacheStatFilter) UseCache(name enum.CacheName) bool {
	if _, inHosts := f.hosts[name]; len(f.hosts) != 0 && !inHosts {
		return false
	}
	if f.cacheType != enum.CacheTypeInvalid && f.cacheTypes[name] != f.cacheType {
		return false
	}
	return true
}

// UseStat returns whether the given stat is in the filter.
func (f *CacheStatFilter) UseStat(statName string) bool {
	if len(f.statsToUse) == 0 {
		return true
	}
	if !f.wildcard {
		_, ok := f.statsToUse[statName]
		return ok
	}
	for statToUse := range f.statsToUse {
		if strings.Contains(statName, statToUse) {
			return true
		}
	}
	return false
}

// WithinStatHistoryMax returns whether the given history index is less than the max history of this filter.
func (f *CacheStatFilter) WithinStatHistoryMax(n int) bool {
	if f.historyCount == 0 {
		return true
	}
	if n <= f.historyCount {
		return true
	}
	return false
}

// NewCacheStatFilter takes the HTTP query parameters and creates a CacheStatFilter which fulfills the `cache.Filter` interface, filtering according to the query parameters passed.
// Query parameters used are `hc`, `stats`, `wildcard`, `type`, and `hosts`.
// If `hc` is 0, all history is returned. If `hc` is empty, 1 history is returned.
// If `stats` is empty, all stats are returned.
// If `wildcard` is empty, `stats` is considered exact.
// If `type` is empty, all cache types are returned.
func NewCacheStatFilter(params url.Values, cacheTypes map[enum.CacheName]enum.CacheType) (cache.Filter, error) {
	validParams := map[string]struct{}{"hc": struct{}{}, "stats": struct{}{}, "wildcard": struct{}{}, "type": struct{}{}, "hosts": struct{}{}}
	if len(params) > len(validParams) {
		return nil, fmt.Errorf("invalid query parameters")
	}
	for param := range params {
		if _, ok := validParams[param]; !ok {
			return nil, fmt.Errorf("invalid query parameter '%v'", param)
		}
	}

	historyCount := 1
	if paramHc, exists := params["hc"]; exists && len(paramHc) > 0 {
		v, err := strconv.Atoi(paramHc[0])
		if err == nil {
			historyCount = v
		}
	}

	statsToUse := map[string]struct{}{}
	if paramStats, exists := params["stats"]; exists && len(paramStats) > 0 {
		commaStats := strings.Split(paramStats[0], ",")
		for _, stat := range commaStats {
			statsToUse[stat] = struct{}{}
		}
	}

	wildcard := false
	if paramWildcard, exists := params["wildcard"]; exists && len(paramWildcard) > 0 {
		wildcard, _ = strconv.ParseBool(paramWildcard[0]) // ignore errors, error => false
	}

	cacheType := enum.CacheTypeInvalid
	if paramType, exists := params["type"]; exists && len(paramType) > 0 {
		cacheType = enum.CacheTypeFromString(paramType[0])
		if cacheType == enum.CacheTypeInvalid {
			return nil, fmt.Errorf("invalid query parameter type '%v' - valid types are: {edge, mid}", paramType[0])
		}
	}

	hosts := map[enum.CacheName]struct{}{}
	if paramHosts, exists := params["hosts"]; exists && len(paramHosts) > 0 {
		commaHosts := strings.Split(paramHosts[0], ",")
		for _, host := range commaHosts {
			hosts[enum.CacheName(host)] = struct{}{}
		}
	}
	// parameters without values are considered hosts, e.g. `?my-cache-0`
	for maybeHost, val := range params {
		if len(val) == 0 || (len(val) == 1 && val[0] == "") {
			hosts[enum.CacheName(maybeHost)] = struct{}{}
		}
	}

	return &CacheStatFilter{
		historyCount: historyCount,
		statsToUse:   statsToUse,
		wildcard:     wildcard,
		cacheType:    cacheType,
		hosts:        hosts,
		cacheTypes:   cacheTypes,
	}, nil
}

// DSStatFilter fulfills the cache.Filter interface, for filtering stats. See the `NewDSStatFilter` documentation for details on which query parameters are used to filter.
type DSStatFilter struct {
	historyCount     int
	statsToUse       map[string]struct{}
	wildcard         bool
	dsType           enum.DSType
	deliveryServices map[enum.DeliveryServiceName]struct{}
	dsTypes          map[enum.DeliveryServiceName]enum.DSType
}

// UseDeliveryService returns whether the given delivery service is in this filter.
func (f *DSStatFilter) UseDeliveryService(name enum.DeliveryServiceName) bool {
	if _, inDSes := f.deliveryServices[name]; len(f.deliveryServices) != 0 && !inDSes {
		return false
	}
	if f.dsType != enum.DSTypeInvalid && f.dsTypes[name] != f.dsType {
		return false
	}
	return true
}

// UseStat returns whether the given stat is in this filter.
func (f *DSStatFilter) UseStat(statName string) bool {
	if len(f.statsToUse) == 0 {
		return true
	}
	if !f.wildcard {
		_, ok := f.statsToUse[statName]
		return ok
	}
	for statToUse := range f.statsToUse {
		if strings.Contains(statName, statToUse) {
			return true
		}
	}
	return false
}

// WithinStatHistoryMax returns whether the given history index is less than the max history of this filter.
func (f *DSStatFilter) WithinStatHistoryMax(n int) bool {
	if f.historyCount == 0 {
		return true
	}
	if n <= f.historyCount {
		return true
	}
	return false
}

// NewDSStatFilter takes the HTTP query parameters and creates a cache.Filter, filtering according to the query parameters passed.
// Query parameters used are `hc`, `stats`, `wildcard`, `type`, and `deliveryservices`.
// If `hc` is 0, all history is returned. If `hc` is empty, 1 history is returned.
// If `stats` is empty, all stats are returned.
// If `wildcard` is empty, `stats` is considered exact.
// If `type` is empty, all types are returned.
func NewDSStatFilter(params url.Values, dsTypes map[enum.DeliveryServiceName]enum.DSType) (dsdata.Filter, error) {
	validParams := map[string]struct{}{"hc": struct{}{}, "stats": struct{}{}, "wildcard": struct{}{}, "type": struct{}{}, "deliveryservices": struct{}{}}
	if len(params) > len(validParams) {
		return nil, fmt.Errorf("invalid query parameters")
	}
	for param := range params {
		if _, ok := validParams[param]; !ok {
			return nil, fmt.Errorf("invalid query parameter '%v'", param)
		}
	}

	historyCount := 1
	if paramHc, exists := params["hc"]; exists && len(paramHc) > 0 {
		v, err := strconv.Atoi(paramHc[0])
		if err == nil {
			historyCount = v
		}
	}

	statsToUse := map[string]struct{}{}
	if paramStats, exists := params["stats"]; exists && len(paramStats) > 0 {
		commaStats := strings.Split(paramStats[0], ",")
		for _, stat := range commaStats {
			statsToUse[stat] = struct{}{}
		}
	}

	wildcard := false
	if paramWildcard, exists := params["wildcard"]; exists && len(paramWildcard) > 0 {
		wildcard, _ = strconv.ParseBool(paramWildcard[0]) // ignore errors, error => false
	}

	dsType := enum.DSTypeInvalid
	if paramType, exists := params["type"]; exists && len(paramType) > 0 {
		dsType = enum.DSTypeFromString(paramType[0])
		if dsType == enum.DSTypeInvalid {
			return nil, fmt.Errorf("invalid query parameter type '%v' - valid types are: {http, dns}", paramType[0])
		}
	}

	deliveryServices := map[enum.DeliveryServiceName]struct{}{}
	// TODO rename 'hosts' to 'names' for consistency
	if paramNames, exists := params["deliveryservices"]; exists && len(paramNames) > 0 {
		commaNames := strings.Split(paramNames[0], ",")
		for _, name := range commaNames {
			deliveryServices[enum.DeliveryServiceName(name)] = struct{}{}
		}
	}
	// parameters without values are considered names, e.g. `?my-cache-0` or `?my-delivery-service`
	for maybeName, val := range params {
		if len(val) == 0 || (len(val) == 1 && val[0] == "") {
			deliveryServices[enum.DeliveryServiceName(maybeName)] = struct{}{}
		}
	}

	return &DSStatFilter{
		historyCount:     historyCount,
		statsToUse:       statsToUse,
		wildcard:         wildcard,
		dsType:           dsType,
		deliveryServices: deliveryServices,
		dsTypes:          dsTypes,
	}, nil
}

// PeerStateFilter fulfills the cache.Filter interface, for filtering stats. See the `NewPeerStateFilter` documentation for details on which query parameters are used to filter.
type PeerStateFilter struct {
	historyCount int
	cachesToUse  map[enum.CacheName]struct{}
	peersToUse   map[enum.TrafficMonitorName]struct{}
	wildcard     bool
	cacheType    enum.CacheType
	cacheTypes   map[enum.CacheName]enum.CacheType
}

// UsePeer returns whether the given Traffic Monitor peer is in this filter.
func (f *PeerStateFilter) UsePeer(name enum.TrafficMonitorName) bool {
	if _, inPeers := f.peersToUse[name]; len(f.peersToUse) != 0 && !inPeers {
		return false
	}
	return true
}

// UseCache returns whether the given cache is in this filter.
func (f *PeerStateFilter) UseCache(name enum.CacheName) bool {
	if f.cacheType != enum.CacheTypeInvalid && f.cacheTypes[name] != f.cacheType {
		return false
	}

	if len(f.cachesToUse) == 0 {
		return true
	}

	if !f.wildcard {
		_, ok := f.cachesToUse[name]
		return ok
	}
	for cacheToUse := range f.cachesToUse {
		if strings.Contains(string(name), string(cacheToUse)) {
			return true
		}
	}
	return false
}

// WithinStatHistoryMax returns whether the given history index is less than the max history of this filter.
func (f *PeerStateFilter) WithinStatHistoryMax(n int) bool {
	if f.historyCount == 0 {
		return true
	}
	if n <= f.historyCount {
		return true
	}
	return false
}

// NewPeerStateFilter takes the HTTP query parameters and creates a cache.Filter, filtering according to the query parameters passed.
// Query parameters used are `hc`, `stats`, `wildcard`, `typep`, and `hosts`. The `stats` param filters caches. The `hosts` param filters peer Traffic Monitors. The `type` param filters cache types (edge, mid).
// If `hc` is 0, all history is returned. If `hc` is empty, 1 history is returned.
// If `stats` is empty, all stats are returned.
// If `wildcard` is empty, `stats` is considered exact.
// If `type` is empty, all cache types are returned.
func NewPeerStateFilter(params url.Values, cacheTypes map[enum.CacheName]enum.CacheType) (*PeerStateFilter, error) {
	// TODO change legacy `stats` and `hosts` to `caches` and `monitors` (or `peers`).
	validParams := map[string]struct{}{"hc": struct{}{}, "stats": struct{}{}, "wildcard": struct{}{}, "type": struct{}{}, "peers": struct{}{}}
	if len(params) > len(validParams) {
		return nil, fmt.Errorf("invalid query parameters")
	}
	for param := range params {
		if _, ok := validParams[param]; !ok {
			return nil, fmt.Errorf("invalid query parameter '%v'", param)
		}
	}

	historyCount := 1
	if paramHc, exists := params["hc"]; exists && len(paramHc) > 0 {
		v, err := strconv.Atoi(paramHc[0])
		if err == nil {
			historyCount = v
		}
	}

	cachesToUse := map[enum.CacheName]struct{}{}
	// TODO rename 'stats' to 'caches'
	if paramStats, exists := params["stats"]; exists && len(paramStats) > 0 {
		commaStats := strings.Split(paramStats[0], ",")
		for _, stat := range commaStats {
			cachesToUse[enum.CacheName(stat)] = struct{}{}
		}
	}

	wildcard := false
	if paramWildcard, exists := params["wildcard"]; exists && len(paramWildcard) > 0 {
		wildcard, _ = strconv.ParseBool(paramWildcard[0]) // ignore errors, error => false
	}

	cacheType := enum.CacheTypeInvalid
	if paramType, exists := params["type"]; exists && len(paramType) > 0 {
		cacheType = enum.CacheTypeFromString(paramType[0])
		if cacheType == enum.CacheTypeInvalid {
			return nil, fmt.Errorf("invalid query parameter type '%v' - valid types are: {edge, mid}", paramType[0])
		}
	}

	peersToUse := map[enum.TrafficMonitorName]struct{}{}
	if paramNames, exists := params["peers"]; exists && len(paramNames) > 0 {
		commaNames := strings.Split(paramNames[0], ",")
		for _, name := range commaNames {
			peersToUse[enum.TrafficMonitorName(name)] = struct{}{}
		}
	}
	// parameters without values are considered names, e.g. `?my-cache-0` or `?my-delivery-service`
	for maybeName, val := range params {
		if len(val) == 0 || (len(val) == 1 && val[0] == "") {
			peersToUse[enum.TrafficMonitorName(maybeName)] = struct{}{}
		}
	}

	return &PeerStateFilter{
		historyCount: historyCount,
		cachesToUse:  cachesToUse,
		wildcard:     wildcard,
		cacheType:    cacheType,
		peersToUse:   peersToUse,
		cacheTypes:   cacheTypes,
	}, nil
}

// HandleErr takes an error, and the request type it came from, and logs. It is ok to call with a nil error, in which case this is a no-op.
func HandleErr(errorCount UintThreadsafe, reqPath string, err error) {
	if err == nil {
		return
	}
	errorCount.Inc()
	log.Errorf("Request Error: %v\n", fmt.Errorf(reqPath+": %v", err))
}

// WrapErrCode takes the body, err, and log context (errorCount, reqPath). It logs and deals with any error, and returns the appropriate bytes and response code for the `srvhttp`. It notably returns InternalServerError status on any error, for security reasons.
func WrapErrCode(errorCount UintThreadsafe, reqPath string, body []byte, err error) ([]byte, int) {
	if err == nil {
		return body, http.StatusOK
	}
	HandleErr(errorCount, reqPath, err)
	return nil, http.StatusInternalServerError
}

// WrapBytes takes a function which cannot error and returns only bytes, and wraps it as a http.HandlerFunc. The errContext is logged if the write fails, and should be enough information to trace the problem (function name, endpoint, request parameters, etc).
func WrapBytes(f func() []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Write(w, f(), r.URL.EscapedPath())
	}
}

// WrapErr takes a function which returns bytes and an error, and wraps it as a http.HandlerFunc. If the error is nil, the bytes are written with Status OK. Else, the error is logged, and InternalServerError is returned as the response code. If you need to return a different response code (for example, StatusBadRequest), call wrapRespCode.
func WrapErr(errorCount UintThreadsafe, f func() ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := f()
		_, code := WrapErrCode(errorCount, r.URL.EscapedPath(), bytes, err)
		w.WriteHeader(code)
		log.Write(w, bytes, r.URL.EscapedPath())
	}
}

// SrvFunc is a function which takes URL parameters, and returns the requested data, and a response code. Note it does not take the full http.Request, and does not have the path. SrvFunc functions should be called via dispatch, and any additional data needed should be closed via a lambda.
// TODO split params and path into 2 separate wrappers?
// TODO change to simply take the http.Request?
type SrvFunc func(params url.Values, path string) ([]byte, int)

// WrapParams takes a SrvFunc and wraps it as an http.HandlerFunc. Note if the SrvFunc returns 0 bytes, an InternalServerError is returned, and the response code is ignored, for security reasons. If an error response code is necessary, return bytes to that effect, for example, "Bad Request". DO NOT return informational messages regarding internal server errors; these should be logged, and only a 500 code returned to the client, for security reasons.
func WrapParams(f SrvFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, code := f(r.URL.Query(), r.URL.EscapedPath())
		if len(bytes) > 0 {
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

func srvTRConfig(opsConfig OpsConfigThreadsafe, toSession towrap.ITrafficOpsSession) ([]byte, error) {
	cdnName := opsConfig.Get().CdnName
	if toSession == nil {
		return nil, fmt.Errorf("Unable to connect to Traffic Ops")
	}
	if cdnName == "" {
		return nil, fmt.Errorf("No CDN Configured")
	}
	return toSession.CRConfigRaw(cdnName)
}

func makeWrapAll(errorCount UintThreadsafe, unpolledCaches UnpolledCachesThreadsafe) func(http.HandlerFunc) http.HandlerFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return wrapUnpolledCheck(unpolledCaches, errorCount, f)
	}
}

func makeCrConfigHandler(wrapper func(http.HandlerFunc) http.HandlerFunc, errorCount UintThreadsafe, opsConfig OpsConfigThreadsafe, toSession towrap.ITrafficOpsSession) http.HandlerFunc {
	return wrapper(WrapErr(errorCount, func() ([]byte, error) {
		return srvTRConfig(opsConfig, toSession)
	}))
}

func srvTRState(params url.Values, localStates peer.CRStatesThreadsafe, combinedStates peer.CRStatesThreadsafe) ([]byte, error) {
	if _, raw := params["raw"]; raw {
		return srvTRStateSelf(localStates)
	}
	return srvTRStateDerived(combinedStates)
}

func srvTRStateDerived(combinedStates peer.CRStatesThreadsafe) ([]byte, error) {
	return peer.CrstatesMarshall(combinedStates.Get())
}

func srvTRStateSelf(localStates peer.CRStatesThreadsafe) ([]byte, error) {
	return peer.CrstatesMarshall(localStates.Get())
}

// TODO remove error params, handle by returning an error? How, since we need to return a non-standard code?
func srvCacheStats(params url.Values, errorCount UintThreadsafe, errContext string, toData todata.TODataThreadsafe, statHistory ResultHistoryThreadsafe) ([]byte, int) {
	filter, err := NewCacheStatFilter(params, toData.Get().ServerTypes)
	if err != nil {
		HandleErr(errorCount, errContext, err)
		return []byte(err.Error()), http.StatusBadRequest
	}
	bytes, err := cache.StatsMarshall(statHistory.Get(), filter, params)
	return WrapErrCode(errorCount, errContext, bytes, err)
}

func srvDSStats(params url.Values, errorCount UintThreadsafe, errContext string, toData todata.TODataThreadsafe, dsStats DSStatsReader) ([]byte, int) {
	filter, err := NewDSStatFilter(params, toData.Get().DeliveryServiceTypes)
	if err != nil {
		HandleErr(errorCount, errContext, err)
		return []byte(err.Error()), http.StatusBadRequest
	}
	bytes, err := json.Marshal(dsStats.Get().JSON(filter, params))
	return WrapErrCode(errorCount, errContext, bytes, err)
}

func srvEventLog(events EventsThreadsafe) ([]byte, error) {
	return json.Marshal(JSONEvents{Events: events.Get()})
}

func srvPeerStates(params url.Values, errorCount UintThreadsafe, errContext string, toData todata.TODataThreadsafe, peerStates peer.CRStatesPeersThreadsafe) ([]byte, int) {
	filter, err := NewPeerStateFilter(params, toData.Get().ServerTypes)
	if err != nil {
		HandleErr(errorCount, errContext, err)
		return []byte(err.Error()), http.StatusBadRequest
	}
	bytes, err := json.Marshal(createAPIPeerStates(peerStates.Get(), filter, params))
	return WrapErrCode(errorCount, errContext, bytes, err)
}

func srvStatSummary() ([]byte, int) {
	return nil, http.StatusNotImplemented
}

func srvStats(staticAppData StaticAppData, healthPollInterval time.Duration, lastHealthDurations DurationMapThreadsafe, fetchCount UintThreadsafe, healthIteration UintThreadsafe, errorCount UintThreadsafe) ([]byte, error) {
	return getStats(staticAppData, healthPollInterval, lastHealthDurations.Get(), fetchCount.Get(), healthIteration.Get(), errorCount.Get())
}

func srvConfigDoc(opsConfig OpsConfigThreadsafe) ([]byte, error) {
	opsConfigCopy := opsConfig.Get()
	// if the password is blank, leave it blank, so callers can see it's missing.
	if opsConfigCopy.Password != "" {
		opsConfigCopy.Password = "*****"
	}
	return json.Marshal(opsConfigCopy)
}

// TODO determine if this should use peerStates
func srvAPICacheCount(localStates peer.CRStatesThreadsafe) []byte {
	return []byte(strconv.Itoa(len(localStates.Get().Caches)))
}

func srvAPICacheAvailableCount(localStates peer.CRStatesThreadsafe) []byte {
	return []byte(strconv.Itoa(cacheAvailableCount(localStates.Get().Caches)))
}

func srvAPICacheDownCount(localStates peer.CRStatesThreadsafe, monitorConfig TrafficMonitorConfigMapThreadsafe) []byte {
	return []byte(strconv.Itoa(cacheDownCount(localStates.Get().Caches, monitorConfig.Get().TrafficServer)))
}

func srvAPIVersion(staticAppData StaticAppData) []byte {
	s := "traffic_monitor-" + staticAppData.Version + "."
	if len(staticAppData.GitRevision) > 6 {
		s += staticAppData.GitRevision[:6]
	} else {
		s += staticAppData.GitRevision
	}
	return []byte(s)
}

func srvAPITrafficOpsURI(opsConfig OpsConfigThreadsafe) []byte {
	return []byte(opsConfig.Get().Url)
}
func srvAPICacheStates(toData todata.TODataThreadsafe, statHistory ResultHistoryThreadsafe, healthHistory ResultHistoryThreadsafe, lastHealthDurations DurationMapThreadsafe, localStates peer.CRStatesThreadsafe, lastStats LastStatsThreadsafe, localCacheStatus CacheAvailableStatusThreadsafe) ([]byte, error) {
	return json.Marshal(createCacheStatuses(toData.Get().ServerTypes, statHistory.Get(), healthHistory.Get(), lastHealthDurations.Get(), localStates.Get().Caches, lastStats.Get(), localCacheStatus))
}

func srvAPIBandwidthKbps(toData todata.TODataThreadsafe, lastStats LastStatsThreadsafe) []byte {
	// serverTypes := toData.Get().ServerTypes
	kbpsStats := lastStats.Get()
	sum := float64(0.0)
	for _, data := range kbpsStats.Caches {
		// if serverTypes[cache] != enum.CacheTypeEdge {
		// 	continue
		// }
		sum += data.Bytes.PerSec / ds.BytesPerKilobit
	}
	return []byte(fmt.Sprintf("%f", sum))
}
func srvAPIBandwidthCapacityKbps(statHistoryThs ResultHistoryThreadsafe) []byte {
	statHistory := statHistoryThs.Get()
	cap := int64(0)
	for _, results := range statHistory {
		if len(results) == 0 {
			continue
		}
		cap += results[0].PrecomputedData.MaxKbps
	}
	return []byte(fmt.Sprintf("%d", cap))
}

// WrapUnpolledCheck wraps an http.HandlerFunc, returning ServiceUnavailable if any caches are unpolled; else, calling the wrapped func.
func wrapUnpolledCheck(unpolledCaches UnpolledCachesThreadsafe, errorCount UintThreadsafe, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if unpolledCaches.Any() {
			HandleErr(errorCount, r.URL.EscapedPath(), fmt.Errorf("service still starting, some caches unpolled"))
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Write(w, []byte("Service Unavailable"), r.URL.EscapedPath())
			return
		}
		f(w, r)
	}
}

// MakeDispatchMap returns the map of paths to http.HandlerFuncs for dispatching.
func MakeDispatchMap(
	opsConfig OpsConfigThreadsafe,
	toSession towrap.ITrafficOpsSession,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	statHistory ResultHistoryThreadsafe,
	healthHistory ResultHistoryThreadsafe,
	dsStats DSStatsReader,
	events EventsThreadsafe,
	staticAppData StaticAppData,
	healthPollInterval time.Duration,
	lastHealthDurations DurationMapThreadsafe,
	fetchCount UintThreadsafe,
	healthIteration UintThreadsafe,
	errorCount UintThreadsafe,
	toData todata.TODataThreadsafe,
	localCacheStatus CacheAvailableStatusThreadsafe,
	lastStats LastStatsThreadsafe,
	unpolledCaches UnpolledCachesThreadsafe,
	monitorConfig TrafficMonitorConfigMapThreadsafe,
) map[string]http.HandlerFunc {

	// wrap composes all universal wrapper functions. Right now, it's only the UnpolledCheck, but there may be others later. For example, security headers.
	wrap := func(f http.HandlerFunc) http.HandlerFunc {
		return wrapUnpolledCheck(unpolledCaches, errorCount, f)
	}

	return map[string]http.HandlerFunc{
		"/publish/CrConfig": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvTRConfig(opsConfig, toSession)
		})),
		"/publish/CrStates": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			bytes, err := srvTRState(params, localStates, combinedStates)
			return WrapErrCode(errorCount, path, bytes, err)
		})),
		"/publish/CacheStats": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvCacheStats(params, errorCount, path, toData, statHistory)
		})),
		"/publish/DsStats": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvDSStats(params, errorCount, path, toData, dsStats)
		})),
		"/publish/EventLog": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvEventLog(events)
		})),
		"/publish/PeerStates": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvPeerStates(params, errorCount, path, toData, peerStates)
		})),
		"/publish/StatSummary": wrap(WrapParams(func(params url.Values, path string) ([]byte, int) {
			return srvStatSummary()
		})),
		"/publish/Stats": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvStats(staticAppData, healthPollInterval, lastHealthDurations, fetchCount, healthIteration, errorCount)
		})),
		"/publish/ConfigDoc": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvConfigDoc(opsConfig)
		})),
		"/api/cache-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheCount(localStates)
		})),
		"/api/cache-available-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheAvailableCount(localStates)
		})),
		"/api/cache-down-count": wrap(WrapBytes(func() []byte {
			return srvAPICacheDownCount(localStates, monitorConfig)
		})),
		"/api/version": wrap(WrapBytes(func() []byte {
			return srvAPIVersion(staticAppData)
		})),
		"/api/traffic-ops-uri": wrap(WrapBytes(func() []byte {
			return srvAPITrafficOpsURI(opsConfig)
		})),
		"/api/cache-statuses": wrap(WrapErr(errorCount, func() ([]byte, error) {
			return srvAPICacheStates(toData, statHistory, healthHistory, lastHealthDurations, localStates, lastStats, localCacheStatus)
		})),
		"/api/bandwidth-kbps": wrap(WrapBytes(func() []byte {
			return srvAPIBandwidthKbps(toData, lastStats)
		})),
		"/api/bandwidth-capacity-kbps": wrap(WrapBytes(func() []byte {
			return srvAPIBandwidthCapacityKbps(statHistory)
		})),
	}
}

// latestResultTimeMS returns the length of time in milliseconds that it took to request the most recent non-errored result.
func latestResultTimeMS(cacheName enum.CacheName, history map[enum.CacheName][]cache.Result) (int64, error) {
	results, ok := history[cacheName]
	if !ok {
		return 0, fmt.Errorf("cache %v has no history", cacheName)
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("cache %v history empty", cacheName)
	}
	result := cache.Result{}
	foundResult := false
	for _, r := range results {
		if r.Error == nil {
			result = r
			foundResult = true
			break
		}
	}
	if !foundResult {
		return 0, fmt.Errorf("cache %v No unerrored result", cacheName)
	}
	return int64(result.RequestTime / time.Millisecond), nil
}

func latestQueryTimeMS(cacheName enum.CacheName, lastDurations map[enum.CacheName]time.Duration) (int64, error) {
	queryTime, ok := lastDurations[cacheName]
	if !ok {
		return 0, fmt.Errorf("cache %v not in last durations\n", cacheName)
	}
	return int64(queryTime / time.Millisecond), nil
}

// resultSpanMS returns the length of time between the most recent two results. That is, how long could the cache have been down before we would have noticed it? Note this returns the time between the most recent two results, irrespective if they errored.
func resultSpanMS(cacheName enum.CacheName, history map[enum.CacheName][]cache.Result) (int64, error) {
	results, ok := history[cacheName]
	if !ok {
		return 0, fmt.Errorf("cache %v has no history", cacheName)
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("cache %v history empty", cacheName)
	}
	if len(results) < 2 {
		return 0, fmt.Errorf("cache %v history only has one result, can't compute span between results", cacheName)
	}

	latestResult := results[0]
	penultimateResult := results[1]
	span := latestResult.Time.Sub(penultimateResult.Time)
	return int64(span / time.Millisecond), nil
}

func createCacheStatuses(
	cacheTypes map[enum.CacheName]enum.CacheType,
	statHistory map[enum.CacheName][]cache.Result,
	healthHistory map[enum.CacheName][]cache.Result,
	lastHealthDurations map[enum.CacheName]time.Duration,
	cacheStates map[enum.CacheName]peer.IsAvailable,
	lastStats ds.LastStats,
	localCacheStatusThreadsafe CacheAvailableStatusThreadsafe,
) map[enum.CacheName]CacheStatus {
	conns := createCacheConnections(statHistory)
	statii := map[enum.CacheName]CacheStatus{}
	localCacheStatus := localCacheStatusThreadsafe.Get()

	for cacheName, cacheType := range cacheTypes {
		cacheStatHistory, ok := statHistory[cacheName]
		if !ok {
			log.Warnf("createCacheStatuses stat history missing cache %s\n", cacheName)
			continue
		}

		if len(cacheStatHistory) < 1 {
			log.Warnf("createCacheStatuses stat history empty for cache %s\n", cacheName)
			continue
		}

		log.Debugf("createCacheStatuses NOT empty for cache %s\n", cacheName)

		var loadAverage *float64
		procLoadAvg := cacheStatHistory[0].Astats.System.ProcLoadavg
		if procLoadAvg != "" {
			firstSpace := strings.IndexRune(procLoadAvg, ' ')
			if firstSpace == -1 {
				log.Warnf("WARNING unexpected proc.loadavg '%s' for cache %s\n", procLoadAvg, cacheName)
			} else {
				loadAverageVal, err := strconv.ParseFloat(procLoadAvg[:firstSpace], 64)
				if err != nil {
					log.Warnf("proc.loadavg doesn't contain a float prefix '%s' for cache %s\n", procLoadAvg, cacheName)
				} else {
					loadAverage = &loadAverageVal
				}
			}
		}

		healthQueryTime, err := latestQueryTimeMS(cacheName, lastHealthDurations)
		if err != nil {
			log.Warnf("Error getting cache %v health query time: %v\n", cacheName, err)
		}

		statTime, err := latestResultTimeMS(cacheName, statHistory)
		if err != nil {
			log.Warnf("Error getting cache %v stat result time: %v\n", cacheName, err)
		}

		healthTime, err := latestResultTimeMS(cacheName, healthHistory)
		if err != nil {
			log.Warnf("Error getting cache %v health result time: %v\n", cacheName, err)
		}

		statSpan, err := resultSpanMS(cacheName, statHistory)
		if err != nil {
			log.Warnf("Error getting cache %v stat span: %v\n", cacheName, err)
		}

		healthSpan, err := resultSpanMS(cacheName, healthHistory)
		if err != nil {
			log.Warnf("Error getting cache %v health span: %v\n", cacheName, err)
		}

		var kbps *float64
		lastStat, ok := lastStats.Caches[cacheName]
		if !ok {
			log.Warnf("cache not in last kbps cache %s\n", cacheName)
		} else {
			kbpsVal := lastStat.Bytes.PerSec / float64(ds.BytesPerKilobit)
			kbps = &kbpsVal
		}

		var connections *int64
		connectionsVal, ok := conns[cacheName]
		if !ok {
			log.Warnf("cache not in connections %s\n", cacheName)
		} else {
			connections = &connectionsVal
		}

		var status *string
		statusVal, ok := localCacheStatus[cacheName]
		if !ok {
			log.Warnf("cache not in statuses %s\n", cacheName)
		} else {
			statusString := statusVal.Status + " - "
			if localCacheStatus[cacheName].Available {
				statusString += "available"
			} else {
				statusString += "unavailable"
			}
			status = &statusString
		}

		cacheTypeStr := string(cacheType)
		statii[cacheName] = CacheStatus{
			Type:                   &cacheTypeStr,
			LoadAverage:            loadAverage,
			QueryTimeMilliseconds:  &healthQueryTime,
			StatTimeMilliseconds:   &statTime,
			HealthTimeMilliseconds: &healthTime,
			StatSpanMilliseconds:   &statSpan,
			HealthSpanMilliseconds: &healthSpan,
			BandwidthKbps:          kbps,
			ConnectionCount:        connections,
			Status:                 status,
		}
	}
	return statii
}

func createCacheConnections(statHistory map[enum.CacheName][]cache.Result) map[enum.CacheName]int64 {
	conns := map[enum.CacheName]int64{}
	for server, history := range statHistory {
		for _, result := range history {
			val, ok := result.Astats.Ats["proxy.process.http.current_client_connections"]
			if !ok {
				continue
			}

			v, ok := val.(float64)
			if !ok {
				continue
			}

			conns[server] = int64(v)
			break
		}
	}
	return conns
}

// cacheOfflineCount returns the total caches not available, including marked unavailable, status offline, and status admin_down
func cacheOfflineCount(caches map[enum.CacheName]peer.IsAvailable) int {
	count := 0
	for _, available := range caches {
		if !available.IsAvailable {
			count++
		}
	}
	return count
}

// cacheAvailableCount returns the total caches available, including marked available and status online
func cacheAvailableCount(caches map[enum.CacheName]peer.IsAvailable) int {
	return len(caches) - cacheOfflineCount(caches)
}

// cacheOfflineCount returns the total reported caches marked down, excluding status offline and admin_down.
func cacheDownCount(caches map[enum.CacheName]peer.IsAvailable, toServers map[string]to.TrafficServer) int {
	count := 0
	for cache, available := range caches {
		if !available.IsAvailable && enum.CacheStatusFromString(toServers[string(cache)].Status) == enum.CacheStatusReported {
			count++
		}
	}
	return count
}

func createAPIPeerStates(peerStates map[enum.TrafficMonitorName]peer.Crstates, filter *PeerStateFilter, params url.Values) APIPeerStates {
	apiPeerStates := APIPeerStates{
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
		Peers:         map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState{},
	}

	for peer, state := range peerStates {
		if !filter.UsePeer(peer) {
			continue
		}
		if _, ok := apiPeerStates.Peers[peer]; !ok {
			apiPeerStates.Peers[peer] = map[enum.CacheName][]CacheState{}
		}
		peerState := apiPeerStates.Peers[peer]
		for cache, available := range state.Caches {
			if !filter.UseCache(cache) {
				continue
			}
			peerState[cache] = []CacheState{CacheState{Value: available.IsAvailable}}
		}
		apiPeerStates.Peers[peer] = peerState
	}
	return apiPeerStates
}

// Stats contains statistics data about this running app. Designed to be returned via an API endpoint.
type Stats struct {
	MaxMemoryMB         uint64 `json:"Max Memory (MB)"`
	GitRevision         string `json:"git-revision"`
	ErrorCount          uint64 `json:"Error Count"`
	Uptime              uint64 `json:"uptime"`
	FreeMemoryMB        uint64 `json:"Free Memory (MB)"`
	TotalMemoryMB       uint64 `json:"Total Memory (MB)"`
	Version             string `json:"version"`
	DeployDir           string `json:"deploy-dir"`
	FetchCount          uint64 `json:"Fetch Count"`
	QueryIntervalDelta  int    `json:"Query Interval Delta"`
	IterationCount      uint64 `json:"Iteration Count"`
	Name                string `json:"name"`
	BuildTimestamp      string `json:"buildTimestamp"`
	QueryIntervalTarget int    `json:"Query Interval Target"`
	QueryIntervalActual int    `json:"Query Interval Actual"`
	SlowestCache        string `json:"Slowest Cache"`
	LastQueryInterval   int    `json:"Last Query Interval"`
	Microthreads        int    `json:"Goroutines"`
	LastGC              string `json:"Last Garbage Collection"`
	MemAllocBytes       uint64 `json:"Memory Bytes Allocated"`
	MemTotalBytes       uint64 `json:"Total Bytes Allocated"`
	MemSysBytes         uint64 `json:"System Bytes Allocated"`
}

func getLongestPoll(lastHealthTimes map[enum.CacheName]time.Duration) (enum.CacheName, time.Duration) {
	var longestCache enum.CacheName
	var longestTime time.Duration
	for cache, time := range lastHealthTimes {
		if time > longestTime {
			longestTime = time
			longestCache = cache
		}
	}
	return longestCache, longestTime
}

func getStats(staticAppData StaticAppData, pollingInterval time.Duration, lastHealthTimes map[enum.CacheName]time.Duration, fetchCount uint64, healthIteration uint64, errorCount uint64) ([]byte, error) {
	longestPollCache, longestPollTime := getLongestPoll(lastHealthTimes)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var s Stats
	s.MaxMemoryMB = memStats.TotalAlloc / (1024 * 1024)
	s.GitRevision = staticAppData.GitRevision
	s.ErrorCount = errorCount
	s.Uptime = uint64(time.Since(staticAppData.StartTime) / time.Second)
	s.FreeMemoryMB = staticAppData.FreeMemoryMB
	s.TotalMemoryMB = memStats.Alloc / (1024 * 1024) // TODO rename to "used memory" if/when nothing is using the JSON entry
	s.Version = staticAppData.Version
	s.DeployDir = staticAppData.WorkingDir
	s.FetchCount = fetchCount
	s.SlowestCache = string(longestPollCache)
	s.IterationCount = healthIteration
	s.Name = staticAppData.Name
	s.BuildTimestamp = staticAppData.BuildTimestamp
	s.QueryIntervalTarget = int(pollingInterval / time.Millisecond)
	s.QueryIntervalActual = int(longestPollTime / time.Millisecond)
	s.QueryIntervalDelta = s.QueryIntervalActual - s.QueryIntervalTarget
	s.LastQueryInterval = int(math.Max(float64(s.QueryIntervalActual), float64(s.QueryIntervalTarget)))
	s.Microthreads = runtime.NumGoroutine()
	s.LastGC = time.Unix(0, int64(memStats.LastGC)).String()
	s.MemAllocBytes = memStats.Alloc
	s.MemTotalBytes = memStats.TotalAlloc
	s.MemSysBytes = memStats.Sys

	return json.Marshal(s)
}
