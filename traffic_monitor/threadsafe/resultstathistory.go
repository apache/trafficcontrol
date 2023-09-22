// TODO rename
package threadsafe

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
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/srvhttp"

	jsoniter "github.com/json-iterator/go"
)

// InterfaceStatNames is a "set" of the names of all the statistics that may be
// found on a network interface of a cache server.
const (
	InterfaceStatNameBytesIn  = "inBytes"
	InterfaceStatNameBytesOut = "outBytes"
	InterfaceStatNameSpeed    = "speed"
)

// ResultStatHistory provides safe access for multiple goroutines readers and a single writer to a stored HistoryHistory object.
// This could be made lock-free, if the performance was necessary
type ResultInfoHistory struct {
	history *cache.ResultInfoHistory
	m       *sync.RWMutex
}

// NewResultInfoHistory returns a new ResultInfoHistory safe for multiple readers and a single writer.
func NewResultInfoHistory() ResultInfoHistory {
	h := cache.ResultInfoHistory{}
	return ResultInfoHistory{m: &sync.RWMutex{}, history: &h}
}

// Get returns the ResultInfoHistory. Callers MUST NOT modify. If mutation is necessary, call ResultInfoHistory.Copy()
func (h *ResultInfoHistory) Get() cache.ResultInfoHistory {
	h.m.RLock()
	defer h.m.RUnlock()
	return *h.history
}

// Set sets the internal ResultInfoHistory. This is only safe for one thread of
// execution. This MUST NOT be called from multiple threads.
func (h *ResultInfoHistory) Set(v cache.ResultInfoHistory) {
	h.m.Lock()
	*h.history = v
	h.m.Unlock()
}

// ResultStatHistory is a thread-safe mapping of cache server hostnames to
// CacheStatHistory objects containing statistics for those cache servers.
type ResultStatHistory struct{ *sync.Map } // map[string]CacheStatHistory

// NewResultStatHistory constructs a new, empty ResultStatHistory.
func NewResultStatHistory() ResultStatHistory {
	return ResultStatHistory{&sync.Map{}}
}

// LoadOrStore returns the stored CacheStatHistory for the given cache server
// hostname if it has already been stored. If it has not already been stored, a
// new, empty CacheStatHistory object is created, stored under the given
// hostname, and returned.
func (h ResultStatHistory) LoadOrStore(hostname string) CacheStatHistory {
	// TODO change to use sync.Pool?
	v, _ := h.Map.LoadOrStore(hostname, NewCacheStatHistory())
	rv, ok := v.(CacheStatHistory)
	if !ok {
		log.Errorf("Failed to load or store stat history for '%s': invalid stored type.", hostname)
		return NewCacheStatHistory()
	}

	return rv
}

// Range behaves like sync.Map.Range. It calls f for every value in the map; if
// f returns false, the iteration is stopped.
func (h ResultStatHistory) Range(f func(cacheName string, val CacheStatHistory) bool) {
	h.Map.Range(func(k, v interface{}) bool {
		i, ok := v.(CacheStatHistory)
		if !ok {
			log.Warnf("Non-CacheStatHistory object (%T) found in ResultStatHistory during Range.", v)
			return true
		}
		cacheName, ok := k.(string)
		if !ok {
			log.Warnf("Non-string object (%T) found as key in ResultStatHistory during Range.", k)
			return true
		}
		return f(cacheName, i)
	})
}

// interfaceStat is just a convenience structure used only for passing data
// about a single statistic for a network interface into
// compareAndAppendStatForInterface.
type interfaceStat struct {
	InterfaceName string
	Stat          interface{}
	StatName      string
	Time          time.Time
}

// compareAndAppendStatForInterface is a little helper function used to compare
// a single stat for a single network interface to its historical values and do
// the appropriate appending and management of the history to ensure it never
// exceeds `limit`.
func compareAndAppendStatForInterface(history []tc.ResultStatVal, errs strings.Builder, limit uint64, stat interfaceStat) []tc.ResultStatVal {
	if history == nil {
		history = make([]tc.ResultStatVal, 0, limit)
	}

	ok, err := newStatEqual(history, stat.Stat)
	if err != nil {
		errs.WriteString(stat.InterfaceName)
		errs.Write([]byte(": cannot add stat "))
		errs.WriteString(stat.StatName)
		errs.Write([]byte(": "))
		errs.WriteString(err.Error())
		errs.Write([]byte("; "))
	} else if ok {
		history[0].Time = stat.Time
		history[0].Span++
	} else {
		if uint64(len(history)) > limit {
			history = history[:limit]
		} else if uint64(len(history)) < limit {
			history = append(history, tc.ResultStatVal{})
		}

		for i := len(history) - 1; i >= 1; i-- {
			history[i] = history[i-1]
		}
		history[0] = tc.ResultStatVal{
			Val:  stat.Stat,
			Time: stat.Time,
			Span: 1,
		}
	}
	return history
}

// Add adds the given result to the stored statistics history, keeping only up
// to `limit` number of records for any given stat (oldest records will be
// removed to make way for new ones if the limit would otherwise be exceeded).
//
// If `limit` is zero, it will be treated as though it were one instead.
func (a ResultStatHistory) Add(r cache.Result, limit uint64) error {
	var errStrs strings.Builder
	cacheHistory := a.LoadOrStore(r.ID)
	if limit == 0 {
		log.Warnln("ResultStatHistory.Add got limit 0 - setting to 1")
		limit = 1
	}

	for statName, statVal := range r.Miscellaneous {
		statHistory := cacheHistory.Stats.Load(statName)
		if statHistory == nil {
			statHistory = make([]tc.ResultStatVal, 0, limit)
		}

		ok, err := newStatEqual(statHistory, statVal)
		// If the new stat value is the same as the last, update the time and
		// increment the span. Span is the number of polls the latest value has
		// been the same, and hence the length of time it's been the same is
		// span*pollInterval.
		if err != nil {
			errStrs.Write([]byte("cannot add stat "))
			errStrs.WriteString(statName)
			errStrs.Write([]byte(": "))
			errStrs.WriteString(err.Error())
			errStrs.Write([]byte("; "))
		} else if ok {
			statHistory[0].Time = r.Time
			statHistory[0].Span++
		} else {
			if uint64(len(statHistory)) > limit {
				statHistory = statHistory[:limit]
			} else if uint64(len(statHistory)) < limit {
				statHistory = append(statHistory, tc.ResultStatVal{})
			}

			for i := len(statHistory) - 1; i >= 1; i-- {
				statHistory[i] = statHistory[i-1]
			}
			statHistory[0] = tc.ResultStatVal{
				Val:  statVal,
				Time: r.Time,
				Span: 1,
			}
		}
		cacheHistory.Stats.Store(statName, statHistory)

	}

	stat := interfaceStat{
		Time: r.Time,
	}
	for interfaceName, inf := range r.Interfaces() {
		statHistory, ok := cacheHistory.Interfaces[interfaceName]
		if !ok {
			statHistory = NewResultStatValHistory()
			cacheHistory.Interfaces[interfaceName] = statHistory
		}

		speedHistory := statHistory.Load(InterfaceStatNameSpeed)

		stat.InterfaceName = interfaceName
		stat.Stat = inf.Speed
		stat.StatName = InterfaceStatNameSpeed

		speedHistory = compareAndAppendStatForInterface(speedHistory, errStrs, limit, stat)
		statHistory.Store(InterfaceStatNameSpeed, speedHistory)

		outHistory := statHistory.Load(InterfaceStatNameBytesOut)

		stat.Stat = inf.BytesOut
		stat.StatName = InterfaceStatNameBytesOut

		outHistory = compareAndAppendStatForInterface(outHistory, errStrs, limit, stat)
		statHistory.Store(InterfaceStatNameBytesOut, outHistory)

		inHistory := statHistory.Load(InterfaceStatNameBytesIn)

		stat.Stat = inf.BytesIn
		stat.StatName = InterfaceStatNameBytesIn

		inHistory = compareAndAppendStatForInterface(inHistory, errStrs, limit, stat)
		statHistory.Store(InterfaceStatNameBytesIn, inHistory)
	}

	if errStrs.Len() > 0 {
		errStr := errStrs.String()
		return errors.New("some stats could not be added: " + errStr[:len(errStr)-2])
	}
	return nil
}

// ResultStatValHistory is thread-safe for one writer. Specifically, because a
// CompareAndSwap is not provided, it's not possible to Load and Store without
// a race condition. If multiple writers were necessary, it wouldn't be
// difficult to add a CompareAndSwap, internally storing an atomically-accessed
// pointer to the slice.
type ResultStatValHistory struct{ *sync.Map } //  map[string][]ResultStatVal

func NewResultStatValHistory() ResultStatValHistory { return ResultStatValHistory{&sync.Map{}} }

// Load returns the []ResultStatVal for the given stat. If the given stat does
// not exist, nil is returned.
func (h ResultStatValHistory) Load(stat string) []tc.ResultStatVal {
	i, ok := h.Map.Load(stat)
	if !ok {
		return nil
	}
	return i.([]tc.ResultStatVal)
}

// Range behaves like sync.Map.Range. It calls f for every value in the map; if
// f returns false, the iteration is stopped.
func (h ResultStatValHistory) Range(f func(stat string, val []tc.ResultStatVal) bool) {
	h.Map.Range(func(k, v interface{}) bool {
		return f(k.(string), v.([]tc.ResultStatVal))
	})
}

// Store stores the given []ResultStatVal in the ResultStatValHistory for the
// given stat. Store is thread-safe for only one writer. Specifically, if there
// are multiple writers, there's a race, that one writer could Load(), another
// writer could Store() underneath it, and the first writer would then Store()
// having lost values. To safely use ResultStatValHistory with multiple writers,
// a CompareAndSwap method would have to be added.
func (h ResultStatValHistory) Store(stat string, vals []tc.ResultStatVal) {
	h.Map.Store(stat, vals)
}

// CacheStatHistory is the type of a single record in a ResultStatHistory map.
// It contains interface statistics as well as historical statistics for each
// of a cache server's polled interfaces.
type CacheStatHistory struct {
	// Interfaces is a map of the names of network interfaces that have been
	// polled for this cache server to historical collections of their polled
	// statistics.
	Interfaces map[string]ResultStatValHistory
	// Stats is a historical collection of all of the cache server's generic
	// (non-interface-dependent) statistics.
	Stats ResultStatValHistory
}

// NewCacheStatHistory constructs a new empty CacheStatHistory.
func NewCacheStatHistory() CacheStatHistory {
	return CacheStatHistory{
		Interfaces: make(map[string]ResultStatValHistory),
		Stats:      NewResultStatValHistory(),
	}
}

// newStatEqual returns whether the given stat is equal to the latest stat in
// history. If len(history)==0, this returns false without error. If the given
// stat is not a JSON primitive (string, number, bool), this returns an error.
// We explicitly refuse to compare arrays and objects, for performance.
func newStatEqual(history []tc.ResultStatVal, stat interface{}) (bool, error) {
	if len(history) == 0 {
		return false, nil // if there's no history, it's "not equal", i.e. store this new history
	}
	switch stat.(type) {
	case bool:
	case float64:
	case int64:
	case string:
	case uint64:
	default:
		return false, fmt.Errorf("incomparable stat type %T", stat)
	}
	switch history[0].Val.(type) {
	case bool:
	case float64:
	case int64:
	case string:
	case uint64:
	default:
		return false, fmt.Errorf("incomparable history stat type %T", stat)
	}
	return stat == history[0].Val, nil
}

func generateStats(
	statResultHistory ResultStatHistory,
	statInfo cache.ResultInfoHistory,
	combinedStates tc.CRStates,
	monitorConfig tc.TrafficMonitorConfigMap,
	statMaxKbpses cache.Kbpses,
	filter cache.Filter,
	params url.Values,
) tc.Stats {
	stats := tc.Stats{
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
		Caches:        map[string]tc.ServerStats{},
	}

	computedStats := cache.ComputedStats()

	// TODO in 1.0, stats are divided into 'location', 'cache', and 'type'. 'cache' are hidden by default.

	for id, combinedStatesCache := range combinedStates.Caches {
		if !filter.UseCache(id) {
			continue
		}

		cacheId := string(id)

		cacheStatResultHistory := statResultHistory.LoadOrStore(cacheId)
		if _, ok := stats.Caches[cacheId]; !ok {
			stats.Caches[cacheId] = tc.ServerStats{
				Interfaces: make(map[string]map[string][]tc.ResultStatVal),
				Stats:      make(map[string][]tc.ResultStatVal),
			}
		}

		cacheStatResultHistory.Stats.Range(func(stat string, vals []tc.ResultStatVal) bool {
			stat = "ats." + stat // legacy reasons
			if !filter.UseStat(stat) {
				return true
			}

			var historyCount uint64 = 1
			for _, val := range vals {
				if !filter.WithinStatHistoryMax(historyCount) {
					break
				}
				if _, ok := stats.Caches[cacheId].Stats[stat]; !ok {
					stats.Caches[cacheId].Stats[stat] = []tc.ResultStatVal{val}
				} else {
					stats.Caches[cacheId].Stats[stat] = append(stats.Caches[cacheId].Stats[stat], val)
				}
				historyCount += val.Span
			}

			return true
		})

		for interfaceName, interfaceHistory := range cacheStatResultHistory.Interfaces {
			interfaceHistory.Range(func(stat string, vals []tc.ResultStatVal) bool {
				if !filter.UseInterfaceStat(stat) {
					return true
				}

				var historyCount uint64 = 1
				for _, val := range vals {
					if !filter.WithinStatHistoryMax(historyCount) {
						break
					}
					if _, ok := stats.Caches[cacheId].Interfaces[interfaceName]; !ok {
						stats.Caches[cacheId].Interfaces[interfaceName] = map[string][]tc.ResultStatVal{}
					}
					stats.Caches[cacheId].Interfaces[interfaceName][stat] = append(stats.Caches[cacheId].Interfaces[interfaceName][stat], val)
					historyCount += val.Span
				}
				return true
			})
		}

		serverInfo, ok := monitorConfig.TrafficServer[string(id)]
		if !ok {
			log.Warnf("cache.StatsMarshall server %s missing from monitorConfig\n", id)
		}

		serverProfile, ok := monitorConfig.Profile[serverInfo.Profile]
		if !ok {
			log.Warnf("cache.StatsMarshall server %s missing profile in monitorConfig\n", id)
		}

		for i, resultInfo := range statInfo[id] {
			if !filter.WithinStatHistoryMax(uint64(i) + 1) {
				break
			}

			t := resultInfo.Time

			for stat, statValF := range computedStats {
				if !filter.UseStat(stat) {
					continue
				}
				rv := tc.ResultStatVal{
					Span: 1,
					Time: t,
					Val:  statValF(resultInfo, serverInfo, serverProfile, combinedStatesCache),
				}
				stats.Caches[cacheId].Stats[stat] = append(stats.Caches[cacheId].Stats[stat], rv)
			}
		}
	}

	return stats
}

// StatsMarshall encodes the stats in JSON, encoding up to historyCount of each
// stat. If statsToUse is empty, all stats are encoded; otherwise, only the
// given stats are encoded. If `wildcard` is true, stats which contain the text
// in each statsToUse are returned, instead of exact stat names. If cacheType is
// not CacheTypeInvalid, only stats for the given type are returned. If hosts is
// not empty, only the given hosts are returned.
func StatsMarshall(
	statResultHistory ResultStatHistory,
	statInfo cache.ResultInfoHistory,
	combinedStates tc.CRStates,
	monitorConfig tc.TrafficMonitorConfigMap,
	statMaxKbpses cache.Kbpses,
	filter cache.Filter,
	params url.Values,
) ([]byte, error) {
	stats := generateStats(statResultHistory, statInfo, combinedStates, monitorConfig, statMaxKbpses, filter, params)

	json := jsoniter.ConfigFastest // TODO make configurable
	return json.Marshal(stats)
}

// LegacyStatsMarshall encodes the stats in JSON, encoding up to historyCount of each
// stat. If statsToUse is empty, all stats are encoded; otherwise, only the
// given stats are encoded. If `wildcard` is true, stats which contain the text
// in each statsToUse are returned, instead of exact stat names. If cacheType is
// not CacheTypeInvalid, only stats for the given type are returned. If hosts is
// not empty, only the given hosts are returned.
func LegacyStatsMarshall(
	statResultHistory ResultStatHistory,
	statInfo cache.ResultInfoHistory,
	combinedStates tc.CRStates,
	monitorConfig tc.TrafficMonitorConfigMap,
	statMaxKbpses cache.Kbpses,
	filter cache.Filter,
	params url.Values,
) ([]byte, error) {

	stats := generateStats(statResultHistory, statInfo, combinedStates, monitorConfig, statMaxKbpses, filter, params)
	skippedCaches, legacyStats := stats.ToLegacy(monitorConfig)
	if len(skippedCaches) > 0 {
		log.Warnln(strings.Join(skippedCaches, "\n"))
	}

	json := jsoniter.ConfigFastest // TODO make configurable
	return json.Marshal(legacyStats)
}

func pruneStats(history []tc.ResultStatVal, limit uint64) []tc.ResultStatVal {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}
