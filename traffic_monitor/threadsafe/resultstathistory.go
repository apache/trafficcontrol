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
	"sync"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/srvhttp"

	"github.com/json-iterator/go"
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

// Set sets the internal ResultInfoHistory. This is only safe for one thread of execution. This MUST NOT be called from multiple threads.
func (h *ResultInfoHistory) Set(v cache.ResultInfoHistory) {
	h.m.Lock()
	*h.history = v
	h.m.Unlock()
}

type ResultStatHistory struct{ *sync.Map } // map[tc.CacheName]ResultStatValHistory

func NewResultStatHistory() ResultStatHistory {
	return ResultStatHistory{&sync.Map{}}
}

func (h ResultStatHistory) LoadOrStore(cache tc.CacheName) ResultStatValHistory {
	// TODO change to use sync.Pool?
	v, _ := h.Map.LoadOrStore(cache, NewResultStatValHistory())
	return v.(ResultStatValHistory)
}

// Range behaves like sync.Map.Range. It calls f for every value in the map; if f returns false, the iteration is stopped.
func (h ResultStatHistory) Range(f func(cache tc.CacheName, val ResultStatValHistory) bool) {
	h.Map.Range(func(k, v interface{}) bool {
		return f(k.(tc.CacheName), v.(ResultStatValHistory))
	})
}

// ResultStatValHistory is threadsafe for one writer. Specifically, because a CompareAndSwap is not provided, it's not possible to Load and Store without a race condition.
// If multiple writers were necessary, it wouldn't be difficult to add a CompareAndSwap, internally storing an atomically-accessed pointer to the slice.
type ResultStatValHistory struct{ *sync.Map } //  map[string][]ResultStatVal

func NewResultStatValHistory() ResultStatValHistory { return ResultStatValHistory{&sync.Map{}} }

// Load returns the []ResultStatVal for the given stat. If the given stat does not exist, nil is returned.
func (h ResultStatValHistory) Load(stat string) []cache.ResultStatVal {
	v, ok := h.Map.Load(stat)
	if !ok {
		return nil
	}
	return v.([]cache.ResultStatVal)
}

// Range behaves like sync.Map.Range. It calls f for every value in the map; if f returns false, the iteration is stopped.
func (h ResultStatValHistory) Range(f func(stat string, val []cache.ResultStatVal) bool) {
	h.Map.Range(func(k, v interface{}) bool {
		return f(k.(string), v.([]cache.ResultStatVal))
	})
}

// Store stores the given []ResultStatVal in the ResultStatValHistory for the given stat. Store is threadsafe for only one writer.
// Specifically, if there are multiple writers, there's a race, that one writer could Load(), another writer could Store() underneath it, and the first writer would then Store() having lost values.
// To safely use ResultStatValHistory with multiple writers, a CompareAndSwap function would have to be added.
func (h ResultStatValHistory) Store(stat string, vals []cache.ResultStatVal) {
	h.Map.Store(stat, vals)
}

func (a ResultStatHistory) Add(r cache.Result, limit uint64) error {
	errStrs := ""
	resultHistory := a.LoadOrStore(tc.CacheName(r.ID))
	if limit == 0 {
		log.Warnln("ResultStatHistory.Add got limit 0 - setting to 1")
		limit = 1
	}

	for statName, statVal := range r.Miscellaneous {
		statHistory := resultHistory.Load(statName)
		if len(statHistory) == 0 {
			statHistory = make([]cache.ResultStatVal, 0, limit) // initialize to the limit, to avoid multiple allocations. TODO put in .Load(statName, defaultSize)?
		}

		// TODO check len(statHistory) == 0 before indexing, potential panic?

		ok, err := newStatEqual(statHistory, statVal)

		// If the new stat value is the same as the last, update the time and increment the span. Span is the number of polls the latest value has been the same, and hence the length of time it's been the same is span*pollInterval.
		if err != nil {
			errStrs += "cannot add stat " + statName + ": " + err.Error() + "; "
		} else if ok {
			statHistory[0].Time = r.Time
			statHistory[0].Span++
		} else {
			resultVal := cache.ResultStatVal{
				Val:  statVal,
				Time: r.Time,
				Span: 1,
			}

			if len(statHistory) > int(limit) {
				statHistory = statHistory[:int(limit)]
			} else if len(statHistory) < int(limit) {
				statHistory = append(statHistory, cache.ResultStatVal{})
			}
			// shift all values to the right, in order to put the new val at the beginning. Faster than allocating memory again
			for i := len(statHistory) - 1; i >= 1; i-- {
				statHistory[i] = statHistory[i-1]
			}
			statHistory[0] = resultVal // new result at the beginning
		}
		resultHistory.Store(statName, statHistory)
	}

	if errStrs != "" {
		return errors.New("some stats could not be added: " + errStrs[:len(errStrs)-2])
	}
	return nil
}

// newStatEqual Returns whether the given stat is equal to the latest stat in history. If len(history)==0, this returns false without error. If the given stat is not a JSON primitive (string, number, bool), this returns an error. We explicitly refuse to compare arrays and objects, for performance.
func newStatEqual(history []cache.ResultStatVal, stat interface{}) (bool, error) {
	if len(history) == 0 {
		return false, nil // if there's no history, it's "not equal", i.e. store this new history
	}
	switch stat.(type) {
	case string:
	case float64:
	case bool:
	default:
		return false, fmt.Errorf("incomparable stat type %T", stat)
	}
	switch history[0].Val.(type) {
	case string:
	case float64:
	case bool:
	default:
		return false, fmt.Errorf("incomparable history stat type %T", stat)
	}
	return stat == history[0].Val, nil
}

// StatsMarshall encodes the stats in JSON, encoding up to historyCount of each stat. If statsToUse is empty, all stats are encoded; otherwise, only the given stats are encoded. If wildcard is true, stats which contain the text in each statsToUse are returned, instead of exact stat names. If cacheType is not CacheTypeInvalid, only stats for the given type are returned. If hosts is not empty, only the given hosts are returned.
func StatsMarshall(statResultHistory ResultStatHistory, statInfo cache.ResultInfoHistory, combinedStates tc.CRStates, monitorConfig tc.TrafficMonitorConfigMap, statMaxKbpses cache.Kbpses, filter cache.Filter, params url.Values) ([]byte, error) {
	stats := cache.Stats{
		CommonAPIData: srvhttp.GetCommonAPIData(params, time.Now()),
		Caches:        map[tc.CacheName]map[string][]cache.ResultStatVal{},
	}

	computedStats := cache.ComputedStats()

	// TODO in 1.0, stats are divided into 'location', 'cache', and 'type'. 'cache' are hidden by default.

	for id, combinedStatesCache := range combinedStates.Caches {
		if !filter.UseCache(id) {
			continue
		}

		cacheStatResultHistory := statResultHistory.LoadOrStore(id)
		cacheStatResultHistory.Range(func(stat string, vals []cache.ResultStatVal) bool {
			stat = "ats." + stat // TM1 prefixes ATS stats with 'ats.'
			if !filter.UseStat(stat) {
				return true
			}
			historyCount := 1
			for _, val := range vals {
				if !filter.WithinStatHistoryMax(historyCount) {
					break
				}
				if _, ok := stats.Caches[id]; !ok {
					stats.Caches[id] = map[string][]cache.ResultStatVal{}
				}
				stats.Caches[id][stat] = append(stats.Caches[id][stat], val)
				historyCount += int(val.Span)
			}
			return true
		})

		serverInfo, ok := monitorConfig.TrafficServer[string(id)]
		if !ok {
			log.Warnf("cache.StatsMarshall server %s missing from monitorConfig\n", id)
		}

		serverProfile, ok := monitorConfig.Profile[serverInfo.Profile]
		if !ok {
			log.Warnf("cache.StatsMarshall server %s missing profile in monitorConfig\n", id)
		}

		for i, resultInfo := range statInfo[id] {
			if !filter.WithinStatHistoryMax(i + 1) {
				break
			}
			if _, ok := stats.Caches[id]; !ok {
				stats.Caches[id] = map[string][]cache.ResultStatVal{}
			}

			t := resultInfo.Time

			for stat, statValF := range computedStats {
				if !filter.UseStat(stat) {
					continue
				}
				stats.Caches[id][stat] = append(stats.Caches[id][stat], cache.ResultStatVal{Val: statValF(resultInfo, serverInfo, serverProfile, combinedStatesCache), Time: t, Span: 1}) // combinedState will default to unavailable
			}
		}
	}

	json := jsoniter.ConfigFastest // TODO make configurable
	return json.Marshal(stats)
}

func pruneStats(history []cache.ResultStatVal, limit uint64) []cache.ResultStatVal {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}
