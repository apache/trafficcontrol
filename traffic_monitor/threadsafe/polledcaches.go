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
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
)

// UnpolledCaches is a structure containing a map of caches names (which have yet to be polled) to
// booleans (which, if true, means the cache is directly polled), which is threadsafe for multiple
// readers and one writer. This could be made lock-free, if the performance was necessary
type UnpolledCaches struct {
	unpolledCaches map[tc.CacheName]bool
	seenCaches     map[tc.CacheName]time.Time
	allCaches      map[tc.CacheName]bool
	initialized    *bool
	m              *sync.RWMutex
}

// NewUnpolledCaches returns a new UnpolledCaches object.
func NewUnpolledCaches() UnpolledCaches {
	b := false
	return UnpolledCaches{
		m:              &sync.RWMutex{},
		unpolledCaches: map[tc.CacheName]bool{},
		allCaches:      map[tc.CacheName]bool{},
		seenCaches:     map[tc.CacheName]time.Time{},
		initialized:    &b,
	}
}

// UnpolledCaches returns a map of caches not yet polled. Callers MUST NOT modify. If mutation is necessary, copy the map
func (t *UnpolledCaches) UnpolledCaches() map[tc.CacheName]bool {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.unpolledCaches
}

// SetNewCaches takes a list of new caches, which may overlap with the existing caches, diffs them, removes any `unpolledCaches` which aren't in the new list, and sets the list of `polledCaches` (which is only used by this func) to the `newCaches`.
func (t *UnpolledCaches) SetNewCaches(newCaches map[tc.CacheName]bool) {
	t.m.Lock()
	defer t.m.Unlock()
	for cache := range t.unpolledCaches {
		if _, ok := newCaches[cache]; !ok {
			delete(t.unpolledCaches, cache)
			delete(t.seenCaches, cache)
		}
	}
	for cache := range t.allCaches {
		if _, ok := newCaches[cache]; !ok {
			delete(t.allCaches, cache)
		}
	}
	for cache, v := range newCaches {
		if _, ok := t.allCaches[cache]; !ok {
			t.unpolledCaches[cache] = v
			t.allCaches[cache] = v
		}
	}
	*t.initialized = true
}

// Any returns whether there are any caches marked as not polled. Also returns true if SetNewCaches() has never been called (assuming there exist caches, if this hasn't been initialized, we couldn't have polled any of them).
func (t *UnpolledCaches) Any() bool {
	t.m.RLock()
	defer t.m.RUnlock()
	return !(*t.initialized) || len(t.unpolledCaches) > 0
}

// AnyDirectlyPolled returns whether there are any directly-polled caches marked as not polled.
// Also returns true if SetNewCaches() has never been called (assuming there exist caches, if this
// hasn't been initialized, we couldn't have polled any of them).
func (t *UnpolledCaches) AnyDirectlyPolled() bool {
	t.m.RLock()
	defer t.m.RUnlock()
	if !*t.initialized {
		return true
	}
	for _, directlyPolled := range t.unpolledCaches {
		if directlyPolled {
			return true
		}
	}
	return false
}

const PolledBytesPerSecTimeout = time.Second * 10

// SetPolled sets cache which have been polled. This is used to determine when the app has fully started up, and we can start serving. Serving Traffic Router with caches as 'down' which simply haven't been polled yet would be bad. Therefore, a cache is set as 'polled' if it has received different bandwidths from two different ATS ticks, OR if the cache is marked as down (and thus we won't get a bandwidth).
// This is threadsafe for one writer, along with `Set`.
// This is fast if there are no unpolled caches. Moreover, its speed is a function of the number of unpolled caches, not the number of caches total.
func (t *UnpolledCaches) SetPolled(results []cache.Result, lastStats dsdata.LastStats) {
	t.m.Lock()
	defer t.m.Unlock()
	numUnpolledCaches := len(t.unpolledCaches)
	if numUnpolledCaches == 0 {
		return
	}
	for cache := range t.unpolledCaches {
	innerLoop:
		for _, result := range results {
			if result.ID != string(cache) {
				continue
			}

			// TODO fix "whether a cache has ever been polled" to be generic somehow. The result.System.NotAvailable check is duplicated in health.EvalCache, and is fragile. What if another "successfully polled but unavailable" flag were added?
			if !result.Available || result.Error != nil || result.Statistics.NotAvailable {
				log.Debugf("polled %v\n", cache)
				delete(t.unpolledCaches, cache)
				delete(t.seenCaches, cache)
				break innerLoop
			}
		}
		lastStat, ok := lastStats.Caches[cache]
		if !ok {
			continue
		}

		if lastStat.Bytes.PerSec != 0 {
			log.Debugf("polled %v\n", cache)
			delete(t.unpolledCaches, cache)
			delete(t.seenCaches, cache)
		} else {
			if _, ok := t.seenCaches[cache]; !ok {
				t.seenCaches[cache] = lastStat.Bytes.Time
			}
		}

		if seenTime, ok := t.seenCaches[cache]; ok && time.Since(seenTime) > PolledBytesPerSecTimeout {
			log.Debugf("polled %v (byte change timed out)\n", cache)
			delete(t.unpolledCaches, cache)
			delete(t.seenCaches, cache)
		}
	}

	if len(t.unpolledCaches) == numUnpolledCaches {
		return
	}
	if len(t.unpolledCaches) != 0 {
		log.Infof("remaining unpolled %v\n", t.unpolledCaches)
	} else {
		log.Infof("all caches polled, ready to serve!\n")
	}
}

// SetHealthPolled sets caches which have been *health* polled (as opposed to *stat* polled).
// This is used, when stat polling is disabled, to determine when the app has fully started up,
// and we can start serving. Serving Traffic Router with caches as 'down' which simply haven't
// been polled yet would be bad. Therefore, a cache is set as 'polled' if it has given TM two
// results, OR if the cache is marked as down (and thus we don't need a 2nd result).
func (t *UnpolledCaches) SetHealthPolled(results []cache.Result) {
	t.m.Lock()
	defer t.m.Unlock()
	numUnpolledCaches := len(t.unpolledCaches)
	if numUnpolledCaches == 0 {
		return
	}
	for cache := range t.unpolledCaches {
	innerLoop:
		for _, result := range results {
			if result.ID != string(cache) {
				continue
			}

			if !result.Available || result.Error != nil || result.Statistics.NotAvailable {
				log.Debugf("polled %v\n", cache)
				delete(t.unpolledCaches, cache)
				delete(t.seenCaches, cache)
				break innerLoop
			} else {
				// if this cache has already been seen once before, consider it polled
				if _, ok := t.seenCaches[cache]; ok {
					delete(t.unpolledCaches, cache)
					delete(t.seenCaches, cache)
				} else {
					t.seenCaches[cache] = time.Time{}
				}
			}
		}
	}
}

// SetRemotePolled sets caches which have been *remote health* polled (as opposed to *locally health* polled).
func (t *UnpolledCaches) SetRemotePolled(results map[tc.CacheName]tc.IsAvailable) {
	t.m.Lock()
	defer t.m.Unlock()
	numUnpolledCaches := len(t.unpolledCaches)
	if numUnpolledCaches == 0 {
		return
	}
	for cache := range t.unpolledCaches {
	innerLoop:
		for cacheName := range results {
			if cacheName != cache {
				continue
			}
			delete(t.unpolledCaches, cache)
			delete(t.seenCaches, cache)
			break innerLoop
		}
	}
}
