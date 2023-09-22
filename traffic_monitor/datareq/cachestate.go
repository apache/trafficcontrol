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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"

	jsoniter "github.com/json-iterator/go"
)

// NotFoundStatus is the status value of an interface that has not been found in
// the polled interface data.
const NotFoundStatus = "unavailable - interface not found"

// OnlineStatus is the status value of all interfaces that are associated with an ONLINE server.
const OnlineStatus = "available - server ONLINE"

// CacheStatus contains summary stat data about the given cache.
type CacheStatus struct {
	Type        *string  `json:"type,omitempty"`
	LoadAverage *float64 `json:"load_average,omitempty"`
	// QueryTimeMilliseconds is the time it took this app to perform a complete
	// query and process the data, end-to-end, for the latest health query.
	QueryTimeMilliseconds *int64 `json:"query_time_ms,omitempty"`
	// HealthTimeMilliseconds is the time it took to make the HTTP request and
	// get back the full response, for the latest health query.
	HealthTimeMilliseconds *int64 `json:"health_time_ms,omitempty"`
	// StatTimeMilliseconds is the time it took to make the HTTP request and get
	// back the full response, for the latest stat query.
	StatTimeMilliseconds *int64 `json:"stat_time_ms,omitempty"`
	// StatSpanMilliseconds is the length of time between completing the most
	// recent two stat queries. This can be used as a rough gauge of the
	// end-to-end query processing time.
	StatSpanMilliseconds *int64 `json:"stat_span_ms,omitempty"`
	// HealthSpanMilliseconds is the length of time between completing the most
	// recent two health queries. This can be used as a rough gauge of the
	// end-to-end query processing time.
	HealthSpanMilliseconds *int64 `json:"health_span_ms,omitempty"`

	Status                *string  `json:"status,omitempty"`
	StatusPoller          *string  `json:"status_poller,omitempty"`
	BandwidthKbps         *float64 `json:"bandwidth_kbps,omitempty"`
	BandwidthCapacityKbps *float64 `json:"bandwidth_capacity_kbps,omitempty"`
	ConnectionCount       *int64   `json:"connection_count,omitempty"`
	IPv4Available         *bool    `json:"ipv4_available,omitempty"`
	IPv6Available         *bool    `json:"ipv6_available,omitempty"`
	CombinedAvailable     *bool    `json:"combined_available,omitempty"`

	Interfaces *map[string]CacheInterfaceStatus `json:"interfaces,omitempty"`
}

// CacheInterfaceStatus represents the status of a single network interface of a
// cache server.
type CacheInterfaceStatus struct {
	Status        string  `json:"status"`
	StatusPoller  string  `json:"status_poller"`
	BandwidthKbps float64 `json:"bandwidth_kbps"`
	Available     bool    `json:"available"`
}

func srvAPICacheStates(
	toData todata.TODataThreadsafe,
	statInfoHistory threadsafe.ResultInfoHistory,
	statResultHistory threadsafe.ResultStatHistory,
	healthHistory threadsafe.ResultHistory,
	lastHealthDurations threadsafe.DurationMap,
	localCacheStatus threadsafe.CacheAvailableStatus,
	statMaxKbpses threadsafe.CacheKbpses,
	monitorConfig threadsafe.TrafficMonitorConfigMap,
) ([]byte, error) {
	json := jsoniter.ConfigFastest
	return json.Marshal(createCacheStatuses(toData.Get().ServerTypes, statInfoHistory.Get(), statResultHistory, healthHistory.Get(), lastHealthDurations.Get(), localCacheStatus, statMaxKbpses, monitorConfig.Get().TrafficServer))
}

// interfaceStatus returns the status of the given interface, both qualitatively
// as a human-readable string and as a boolean indicating its availability.
func interfaceStatus(inf tc.ServerInterfaceInfo, result cache.Result) (string, bool) {
	var cacheError = ""
	if result.Error != nil {
		cacheError = "; " + result.Error.Error()
	}
	vitalsMap := result.InterfaceVitals
	vitals, ok := vitalsMap[inf.Name]
	if !ok {
		return "not found in health polling data" + cacheError, false
	}
	if inf.MaxBandwidth != nil && *inf.MaxBandwidth < uint64(vitals.KbpsOut) {
		return "maximum bandwidth exceeded" + cacheError, false
	}
	return "available", true
}

// createCacheStatuses builds a map of cache server hostnames to their
// respective status by examining the calculated availability and statistics of
// each cache server and its network interfaces.
func createCacheStatuses(
	cacheTypes map[tc.CacheName]tc.CacheType,
	statInfoHistory cache.ResultInfoHistory,
	statResultHistory threadsafe.ResultStatHistory,
	healthHistory map[tc.CacheName][]cache.Result,
	lastHealthDurations map[tc.CacheName]time.Duration,
	localCacheStatusThreadsafe threadsafe.CacheAvailableStatus,
	statMaxKbpses threadsafe.CacheKbpses,
	servers map[string]tc.TrafficServer,
) map[string]CacheStatus {
	conns := createCacheConnections(statResultHistory)
	statii := make(map[string]CacheStatus, len(servers))
	localCacheStatus := localCacheStatusThreadsafe.Get().Copy() // TODO test whether copy is necessary
	maxKbpses := statMaxKbpses.Get()

	for cacheName, serverInfo := range servers {
		interfaceStatuses := make(map[string]CacheInterfaceStatus, len(serverInfo.Interfaces))

		var totalMaxKbps float64 = 0
		maxKbps, maxKbpsOk := maxKbpses[cacheName]
		if !maxKbpsOk {
			log.Infof("Cache server '%s' not in max kbps cache", cacheName)
		} else {
			totalMaxKbps = float64(maxKbps)
		}

		health, healthOk := healthHistory[tc.CacheName(cacheName)]
		if !healthOk {
			log.Infof("Cache server '%s' not in max kbps cache", cacheName)
		} else if len(health) < 1 {
			log.Infof("No health data history for cache server '%s'", cacheName)
			healthOk = false
		}

		var totalKbps float64 = 0

		cacheStatus, statusOk := localCacheStatus[cacheName]
		poller := "unknown"
		if !statusOk {
			log.Warnf("No cache status found for cache '%s'", cacheName)
		} else {
			poller = cacheStatus.Poller
		}
		for _, inf := range serverInfo.Interfaces {
			interfaceName := inf.Name

			infStatus := CacheInterfaceStatus{
				Available:     false,
				BandwidthKbps: 0,
				Status:        NotFoundStatus,
				StatusPoller:  poller,
			}

			if healthOk {
				infStatus.Status, infStatus.Available = interfaceStatus(inf, health[0])
				if infVit, ok := health[0].InterfaceVitals[inf.Name]; ok {
					infStatus.BandwidthKbps = float64(infVit.KbpsOut)
					totalKbps += infStatus.BandwidthKbps
				} else {
					log.Infof("Cache server '%s' interface '%s' not in last health measurement.", cacheName, inf.Name)
				}
			}

			if serverInfo.ServerStatus == tc.CacheStatusOnline.String() {
				infStatus.Status = OnlineStatus
				infStatus.Available = true
			}

			interfaceStatuses[interfaceName] = infStatus
		}

		var connections int64 = 0
		connectionsVal, ok := conns[cacheName]
		if !ok {
			log.Infof("Cache server '%s' not in connections.", cacheName)
		} else {
			connections = connectionsVal
		}

		cacheTypeStr := ""
		if cacheType, ok := cacheTypes[tc.CacheName(cacheName)]; !ok {
			log.Infof("Error getting cache type for %v: not in types\n", cacheName)
		} else {
			cacheTypeStr = string(cacheType)
		}

		loadAverage := 0.0
		if infoHistory, ok := statInfoHistory[tc.CacheName(cacheName)]; !ok {
			log.Infof("createCacheStatuses stat info history missing cache %s\n", cacheName)
		} else if len(infoHistory) < 1 {
			log.Infof("createCacheStatuses stat info history empty for cache %s\n", cacheName)
		} else {
			loadAverage = infoHistory[0].Vitals.LoadAvg
		}

		healthQueryTime, err := latestQueryTimeMS(cacheName, lastHealthDurations)
		if err != nil {
			log.Infof("Error getting cache %v health query time: %v\n", cacheName, err)
		}

		statTime, err := latestResultInfoTimeMS(tc.CacheName(cacheName), statInfoHistory)
		if err != nil {
			log.Infof("Error getting cache %v stat result time: %v\n", cacheName, err)
		}

		healthTime, err := latestResultTimeMS(tc.CacheName(cacheName), healthHistory)
		if err != nil {
			log.Infof("Error getting cache %v health result time: %v\n", cacheName, err)
		}

		statSpan, err := infoResultSpanMS(tc.CacheName(cacheName), statInfoHistory)
		if err != nil {
			log.Infof("Error getting cache %v stat span: %v\n", cacheName, err)
		}

		healthSpan, err := resultSpanMS(tc.CacheName(cacheName), healthHistory)
		if err != nil {
			log.Infof("Error getting cache %v health span: %v\n", cacheName, err)
		}

		if serverInfo.ServerStatus == tc.CacheStatusOnline.String() {
			cacheStatus.Why = "ONLINE - available"
			cacheStatus.Available.IPv4 = serverInfo.IPv4() != ""
			cacheStatus.Available.IPv6 = serverInfo.IPv6() != ""
			cacheStatus.ProcessedAvailable = cacheStatus.Available.IPv4 || cacheStatus.Available.IPv6
		}

		statii[cacheName] = CacheStatus{
			Type:                   &cacheTypeStr,
			LoadAverage:            &loadAverage,
			QueryTimeMilliseconds:  &healthQueryTime,
			StatTimeMilliseconds:   &statTime,
			HealthTimeMilliseconds: &healthTime,
			StatSpanMilliseconds:   &statSpan,
			HealthSpanMilliseconds: &healthSpan,
			BandwidthKbps:          &totalKbps,
			BandwidthCapacityKbps:  &totalMaxKbps,
			ConnectionCount:        &connections,
			Status:                 &cacheStatus.Why,
			StatusPoller:           &poller,
			IPv4Available:          &cacheStatus.Available.IPv4,
			IPv6Available:          &cacheStatus.Available.IPv6,
			CombinedAvailable:      &cacheStatus.ProcessedAvailable,
			Interfaces:             &interfaceStatuses,
		}
	}
	return statii
}

// cacheStatusAndPoller returns the reason why a cache is unavailable (or
// that is available), the poller, and 3 booleans in order: IPv4 availability,
// IPv6 availability and Processed availability which is what the monitor
// reports based on the PollingProtocol chosen (ipv4only,ipv6only or both).
func cacheStatusAndPoller(server string, serverInfo tc.TrafficServer, localCacheStatus cache.AvailableStatuses) (string, string, bool, bool, bool) {
	switch status := tc.CacheStatusFromString(serverInfo.ServerStatus); status {
	case tc.CacheStatusAdminDown:
		fallthrough
	case tc.CacheStatusOnline:
		fallthrough
	case tc.CacheStatusOffline:
		return status.String(), "", false, false, false
	}

	status, ok := localCacheStatus[server]
	if !ok {
		log.Infof("Cache server '%s' not in statuses.", server)
		return "ERROR - not in statuses", "", false, false, false
	}

	var statusStr string
	if status.Why == "" {
		if status.ProcessedAvailable {
			statusStr = status.Status + " - available"
		} else {
			statusStr = status.Status + " - unavailable"
		}
	} else {
		statusStr = status.Why
	}
	return statusStr, status.Poller, status.Available.IPv4, status.Available.IPv6, status.ProcessedAvailable
}

func createCacheConnections(statResultHistory threadsafe.ResultStatHistory) map[string]int64 {
	conns := map[string]int64{}
	statResultHistory.Range(func(server string, history threadsafe.CacheStatHistory) bool {
		// We only want to create connections for each cache
		if _, ok := conns[server]; ok {
			return true
		}
		vals := history.Stats.Load("proxy.process.http.current_client_connections")
		if len(vals) == 0 {
			return true
		}

		v, ok := vals[0].Val.(float64)
		if !ok {
			return true // TODO log warning? error?
		}
		conns[server] = int64(v)
		return true
	})
	return conns
}

// infoResultSpanMS returns the length of time between the most recent two results. That is, how long could the cache have been down before we would have noticed it? Note this returns the time between the most recent two results, irrespective if they errored.
// Note this is unrelated to the Stat Span field.
func infoResultSpanMS(cacheName tc.CacheName, history cache.ResultInfoHistory) (int64, error) {
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

// resultSpanMS returns the length of time between the most recent two results. That is, how long could the cache have been down before we would have noticed it? Note this returns the time between the most recent two results, irrespective if they errored.
// Note this is unrelated to the Stat Span field.
func resultSpanMS(cacheName tc.CacheName, history map[tc.CacheName][]cache.Result) (int64, error) {
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

func latestQueryTimeMS(cacheName string, lastDurations map[tc.CacheName]time.Duration) (int64, error) {
	queryTime, ok := lastDurations[tc.CacheName(cacheName)]
	if !ok {
		return 0, fmt.Errorf("cache %v not in last durations", cacheName)
	}
	return int64(queryTime / time.Millisecond), nil
}

// latestResultTimeMS returns the length of time in milliseconds that it took to request the most recent non-errored result.
func latestResultTimeMS(cacheName tc.CacheName, history map[tc.CacheName][]cache.Result) (int64, error) {

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

// latestResultInfoTimeMS returns the length of time in milliseconds that it took to request the most recent non-errored result info.
func latestResultInfoTimeMS(cacheName tc.CacheName, history cache.ResultInfoHistory) (int64, error) {
	results, ok := history[cacheName]
	if !ok {
		return 0, fmt.Errorf("cache %v has no history", cacheName)
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("cache %v history empty", cacheName)
	}
	result := cache.ResultInfo{}
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
