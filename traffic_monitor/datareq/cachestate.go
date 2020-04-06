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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/ds"
	"github.com/apache/trafficcontrol/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"

	"github.com/json-iterator/go"
)

// CacheStatus contains summary stat data about the given cache.
// TODO make fields nullable, so error fields can be omitted, letting API callers still get updates for unerrored fields
type CacheStatus struct {
	Type         *string  `json:"type,omitempty"`
	Status       *string  `json:"status,omitempty"`
	StatusPoller *string  `json:"status_poller,omitempty"`
	LoadAverage  *float64 `json:"load_average,omitempty"`
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
	BandwidthCapacityKbps  *float64 `json:"bandwidth_capacity_kbps,omitempty"`
	ConnectionCount        *int64   `json:"connection_count,omitempty"`
	IPv4Available          *bool    `json:"ipv4_available,omitempty"`
	IPv6Available          *bool    `json:"ipv6_available,omitempty"`
	CombinedAvailable      *bool    `json:"combined_available,omitempty"`
}

func srvAPICacheStates(
	toData todata.TODataThreadsafe,
	statInfoHistory threadsafe.ResultInfoHistory,
	statResultHistory threadsafe.ResultStatHistory,
	healthHistory threadsafe.ResultHistory,
	lastHealthDurations threadsafe.DurationMap,
	localStates peer.CRStatesThreadsafe,
	lastStats threadsafe.LastStats,
	localCacheStatus threadsafe.CacheAvailableStatus,
	statMaxKbpses threadsafe.CacheKbpses,
	monitorConfig threadsafe.TrafficMonitorConfigMap,
) ([]byte, error) {
	json := jsoniter.ConfigFastest
	return json.Marshal(createCacheStatuses(toData.Get().ServerTypes, statInfoHistory.Get(), statResultHistory, healthHistory.Get(), lastHealthDurations.Get(), localStates.Get().Caches, lastStats.Get(), localCacheStatus, statMaxKbpses, monitorConfig.Get().TrafficServer))
}

func createCacheStatuses(
	cacheTypes map[string]tc.CacheType,
	statInfoHistory cache.ResultInfoHistory,
	statResultHistory threadsafe.ResultStatHistory,
	healthHistory map[string][]cache.Result,
	lastHealthDurations map[string]time.Duration,
	cacheStates map[string]tc.IsAvailable,
	lastStats dsdata.LastStats,
	localCacheStatusThreadsafe threadsafe.CacheAvailableStatus,
	statMaxKbpses threadsafe.CacheKbpses,
	servers map[string]tc.TrafficServer,
) map[string]CacheStatus {
	conns := createCacheConnections(statResultHistory)
	statii := map[string]CacheStatus{}
	localCacheStatus := localCacheStatusThreadsafe.Get().Copy() // TODO test whether copy is necessary
	maxKbpses := statMaxKbpses.Get()

	for cacheName, serverInfo := range servers {
		status, statusPoller, ipv4, ipv6, combinedStatus := cacheStatusAndPoller(cacheName, serverInfo, localCacheStatus)

		cacheTypeStr := ""
		if cacheType, ok := cacheTypes[cacheName]; !ok {
			log.Infof("Error getting cache type for %v: not in types\n", cacheName)
		} else {
			cacheTypeStr = string(cacheType)
		}

		loadAverage := 0.0
		if infoHistory, ok := statInfoHistory[cacheName]; !ok {
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

		statTime, err := latestResultInfoTimeMS(cacheName, statInfoHistory)
		if err != nil {
			log.Infof("Error getting cache %v stat result time: %v\n", cacheName, err)
		}

		healthTime, err := latestResultTimeMS(cacheName, healthHistory)
		if err != nil {
			log.Infof("Error getting cache %v health result time: %v\n", cacheName, err)
		}

		statSpan, err := infoResultSpanMS(cacheName, statInfoHistory)
		if err != nil {
			log.Infof("Error getting cache %v stat span: %v\n", cacheName, err)
		}

		healthSpan, err := resultSpanMS(cacheName, healthHistory)
		if err != nil {
			log.Infof("Error getting cache %v health span: %v\n", cacheName, err)
		}

		var kbps *float64
		if lastStat, ok := lastStats.Caches[cacheName]; !ok {
			log.Infof("cache not in last kbps cache %s\n", cacheName)
		} else {
			kbpsVal := lastStat.Bytes.PerSec / float64(ds.BytesPerKilobit)
			kbps = &kbpsVal
		}

		var maxKbps *float64
		if v, ok := maxKbpses[cacheName]; !ok {
			log.Infof("cache not in max kbps cache %s\n", cacheName)
		} else {
			fv := float64(v)
			maxKbps = &fv
		}

		var connections *int64
		connectionsVal, ok := conns[cacheName]
		if !ok {
			log.Infof("cache not in connections %s\n", cacheName)
		} else {
			connections = &connectionsVal
		}

		statii[cacheName] = CacheStatus{
			Type:                   &cacheTypeStr,
			LoadAverage:            &loadAverage,
			QueryTimeMilliseconds:  &healthQueryTime,
			StatTimeMilliseconds:   &statTime,
			HealthTimeMilliseconds: &healthTime,
			StatSpanMilliseconds:   &statSpan,
			HealthSpanMilliseconds: &healthSpan,
			BandwidthKbps:          kbps,
			BandwidthCapacityKbps:  maxKbps,
			ConnectionCount:        connections,
			Status:                 &status,
			StatusPoller:           &statusPoller,
			IPv4Available:          &ipv4,
			IPv6Available:          &ipv6,
			CombinedAvailable:      &combinedStatus,
		}
	}
	return statii
}

//cacheStatusAndPoller returns the the reason why a cache is unavailable (or that is available), the poller, and 3 booleans in order:
// IPv4 availability, IPv6 availability and Processed availability which is what the monitor reports based on the PollingProtocol chosen (ipv4only,ipv6only or both)
func cacheStatusAndPoller(server string, serverInfo tc.TrafficServer, localCacheStatus cache.AvailableStatuses) (string, string, bool, bool, bool) {
	switch status := tc.CacheStatusFromString(serverInfo.ServerStatus); status {
	case tc.CacheStatusAdminDown:
		fallthrough
	case tc.CacheStatusOnline:
		fallthrough
	case tc.CacheStatusOffline:
		return status.String(), "", false, false, false
	}

	statusVal, ok := localCacheStatus[server]
	if !ok {
		log.Infof("cache not in statuses %s\n", server)
		return "ERROR - not in statuses", "", false, false, false
	}

	if statusVal.Why != "" {
		return statusVal.Why, statusVal.Poller, statusVal.Available.IPv4, statusVal.Available.IPv6, statusVal.ProcessedAvailable
	}
	if statusVal.ProcessedAvailable {
		return fmt.Sprintf("%s - available", statusVal.Status), statusVal.Poller, statusVal.Available.IPv4, statusVal.Available.IPv6, statusVal.ProcessedAvailable
	}
	return fmt.Sprintf("%s - unavailable", statusVal.Status), statusVal.Poller, statusVal.Available.IPv4, statusVal.Available.IPv6, statusVal.ProcessedAvailable
}

func createCacheConnections(statResultHistory threadsafe.ResultStatHistory) map[string]int64 {
	conns := map[string]int64{}
	statResultHistory.Range(func(server string, history threadsafe.ResultStatValHistory) bool {
		// for server, history := range statResultHistory {
		vals := history.Load("proxy.process.http.current_client_connections")
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
func infoResultSpanMS(cacheName string, history cache.ResultInfoHistory) (int64, error) {
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
func resultSpanMS(cacheName string, history map[string][]cache.Result) (int64, error) {
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

func latestQueryTimeMS(cacheName string, lastDurations map[string]time.Duration) (int64, error) {
	queryTime, ok := lastDurations[cacheName]
	if !ok {
		return 0, fmt.Errorf("cache %v not in last durations\n", cacheName)
	}
	return int64(queryTime / time.Millisecond), nil
}

// latestResultTimeMS returns the length of time in milliseconds that it took to request the most recent non-errored result.
func latestResultTimeMS(cacheName string, history map[string][]cache.Result) (int64, error) {

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
func latestResultInfoTimeMS(cacheName string, history cache.ResultInfoHistory) (int64, error) {
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
