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
	"math"
	"runtime"
	"sort"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"

	"github.com/json-iterator/go"
)

type JSONStats struct {
	Stats Stats `json:"stats"`
}

// Stats contains statistics data about this running app. Designed to be returned via an API endpoint.
type Stats struct {
	MaxMemoryMB                 uint64  `json:"Max Memory (MB),string"`
	GitRevision                 string  `json:"git-revision"`
	ErrorCount                  uint64  `json:"Error Count,string"`
	Uptime                      uint64  `json:"uptime,string"`
	FreeMemoryMB                uint64  `json:"Free Memory (MB),string"`
	TotalMemoryMB               uint64  `json:"Total Memory (MB),string"`
	Version                     string  `json:"version"`
	DeployDir                   string  `json:"deploy-dir"`
	FetchCount                  uint64  `json:"Fetch Count,string"`
	QueryIntervalDelta          int     `json:"Query Interval Delta,string"`
	IterationCount              uint64  `json:"Iteration Count,string"`
	Name                        string  `json:"name"`
	BuildTimestamp              string  `json:"buildTimestamp"`
	QueryIntervalTarget         int     `json:"Query Interval Target,string"`
	QueryIntervalActual         int     `json:"Query Interval Actual,string"`
	SlowestCache                string  `json:"Slowest Cache"`
	LastQueryInterval           int     `json:"Last Query Interval,string"`
	Microthreads                int     `json:"Goroutines"`
	LastGC                      string  `json:"Last Garbage Collection"`
	MemAllocBytes               uint64  `json:"Memory Bytes Allocated"`
	MemTotalBytes               uint64  `json:"Total Bytes Allocated"`
	MemSysBytes                 uint64  `json:"System Bytes Allocated"`
	OldestPolledPeer            string  `json:"Oldest Polled Peer"`
	OldestPolledPeerMs          int64   `json:"Oldest Polled Peer Time (ms)"`
	QueryInterval95thPercentile int64   `json:"Query Interval 95th Percentile (ms)"`
	GCCPUFraction               float64 `json:"gc-cpu-fraction"`
}

func srvStats(staticAppData config.StaticAppData, healthPollInterval time.Duration, lastHealthDurations threadsafe.DurationMap, fetchCount threadsafe.Uint, healthIteration threadsafe.Uint, errorCount threadsafe.Uint, peerStates peer.CRStatesPeersThreadsafe) ([]byte, error) {
	return getStats(staticAppData, healthPollInterval, lastHealthDurations.Get(), fetchCount.Get(), healthIteration.Get(), errorCount.Get(), peerStates)
}

func getStats(staticAppData config.StaticAppData, pollingInterval time.Duration, lastHealthTimes map[tc.CacheName]time.Duration, fetchCount uint64, healthIteration uint64, errorCount uint64, peerStates peer.CRStatesPeersThreadsafe) ([]byte, error) {
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
	s.GCCPUFraction = memStats.GCCPUFraction

	oldestPolledPeer, oldestPolledPeerTime := oldestPeerPollTime(peerStates.GetQueryTimes(), peerStates.GetPeersOnline())
	s.OldestPolledPeer = string(oldestPolledPeer)
	s.OldestPolledPeerMs = time.Now().Sub((oldestPolledPeerTime)).Nanoseconds() / util.MSPerNS

	s.QueryInterval95thPercentile = getCacheTimePercentile(lastHealthTimes, 0.95).Nanoseconds() / util.MSPerNS

	json := jsoniter.ConfigDefault
	return json.Marshal(JSONStats{Stats: s})
}

func getLongestPoll(lastHealthTimes map[tc.CacheName]time.Duration) (tc.CacheName, time.Duration) {
	var longestCache tc.CacheName
	var longestTime time.Duration
	for cache, time := range lastHealthTimes {
		if time > longestTime {
			longestTime = time
			longestCache = cache
		}
	}
	return longestCache, longestTime
}

type Durations []time.Duration

func (s Durations) Len() int {
	return len(s)
}
func (s Durations) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s Durations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// getCacheTimePercentile returns the given percentile of cache result times. The `percentile` should be a decimal percent, for example, for the 95th percentile pass 0.95
func getCacheTimePercentile(lastHealthTimes map[tc.CacheName]time.Duration, percentile float64) time.Duration {
	times := make([]time.Duration, 0, len(lastHealthTimes))
	for _, t := range lastHealthTimes {
		times = append(times, t)
	}
	sort.Sort(Durations(times))

	n := int(float64(len(lastHealthTimes)) * percentile)

	return times[n]
}

func oldestPeerPollTime(peerTimes map[tc.TrafficMonitorName]time.Time, peerOnline map[tc.TrafficMonitorName]bool) (tc.TrafficMonitorName, time.Time) {
	now := time.Now()
	oldestTime := now
	oldestPeer := tc.TrafficMonitorName("")
	for p, t := range peerTimes {
		if !peerOnline[p] {
			continue
		}
		if oldestTime.After(t) {
			oldestTime = t
			oldestPeer = p
		}
	}
	if oldestTime == now {
		oldestTime = time.Time{}
	}
	return oldestPeer, oldestTime
}
