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

	"github.com/apache/trafficcontrol/v6/lib/go-log"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v6/traffic_monitor/dsdata"
)

// UnpolledCaches is a structure containing a map of caches which have yet to be polled, which is threadsafe for multiple readers and one writer.
// This could be made lock-free, if the performance was necessary
type UnpolledCaches struct {
	unpolledCaches *map[tc.CacheName]struct{}
	seenCaches     *map[tc.CacheName]time.Time
	allCaches      *map[tc.CacheName]struct{}
	initialized    *bool
	m              *sync.RWMutex
}

// NewUnpolledCaches returns a new UnpolledCaches object.
func NewUnpolledCaches() UnpolledCaches {
	b := false
	return UnpolledCaches{
		m:              &sync.RWMutex{},
		unpolledCaches: &map[tc.CacheName]struct{}{},
		allCaches:      &map[tc.CacheName]struct{}{},
		seenCaches:     &map[tc.CacheName]time.Time{},
		initialized:    &b,
	}
}

// UnpolledCaches returns a map of caches not yet polled. Callers MUST NOT modify. If mutation is necessary, copy the map
func (t *UnpolledCaches) UnpolledCaches() map[tc.CacheName]struct{} {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.unpolledCaches
}

// setUnpolledCaches sets the internal unpolled caches map. This is only safe for one thread of execution. This MUST NOT be called from multiple threads.
func (t *UnpolledCaches) setUnpolledCaches(v map[tc.CacheName]struct{}) {
	t.m.Lock()
	*t.initialized = true
	*t.unpolledCaches = v
	t.m.Unlock()
}

// setUnpolledCaches sets the internal unpolled caches map. This is only safe for one thread of execution. This MUST NOT be called from multiple threads.
func (t *UnpolledCaches) setSeenCaches(v map[tc.CacheName]time.Time) {
	t.m.Lock()
	*t.seenCaches = v
	t.m.Unlock()
}

// SetNewCaches takes a list of new caches, which may overlap with the existing caches, diffs them, removes any `unpolledCaches` which aren't in the new list, and sets the list of `polledCaches` (which is only used by this func) to the `newCaches`. This is threadsafe with one writer, along with `setUnpolledCaches`.
func (t *UnpolledCaches) SetNewCaches(newCaches map[tc.CacheName]struct{}) {
	unpolledCaches := copyCaches(t.UnpolledCaches())
	allCaches := copyCaches(*t.allCaches) // not necessary to lock `allCaches`, as the single-writer is the only thing that accesses it.
	seenCaches := copyCachesTime(*t.seenCaches)
	for cache := range unpolledCaches {
		if _, ok := newCaches[cache]; !ok {
			delete(unpolledCaches, cache)
			delete(seenCaches, cache)
		}
	}
	for cache := range allCaches {
		if _, ok := newCaches[cache]; !ok {
			delete(allCaches, cache)
		}
	}
	for cache := range newCaches {
		if _, ok := allCaches[cache]; !ok {
			unpolledCaches[cache] = struct{}{}
			allCaches[cache] = struct{}{}
		}
	}
	*t.allCaches = allCaches
	t.setUnpolledCaches(unpolledCaches)
	t.setSeenCaches(seenCaches)
}

// Any returns whether there are any caches marked as not polled. Also returns true if SetNewCaches() has never been called (assuming there exist caches, if this hasn't been initialized, we couldn't have polled any of them).
func (t *UnpolledCaches) Any() bool {
	t.m.Lock()
	defer t.m.Unlock()
	return !(*t.initialized) || len(*t.unpolledCaches) > 0
}

// copyCaches performs a deep copy of the given map.
func copyCaches(a map[tc.CacheName]struct{}) map[tc.CacheName]struct{} {
	b := map[tc.CacheName]struct{}{}
	for k := range a {
		b[k] = struct{}{}
	}
	return b
}

func copyCachesTime(a map[tc.CacheName]time.Time) map[tc.CacheName]time.Time {
	b := map[tc.CacheName]time.Time{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

const PolledBytesPerSecTimeout = time.Second * 10

// SetPolled sets cache which have been polled. This is used to determine when the app has fully started up, and we can start serving. Serving Traffic Router with caches as 'down' which simply haven't been polled yet would be bad. Therefore, a cache is set as 'polled' if it has received different bandwidths from two different ATS ticks, OR if the cache is marked as down (and thus we won't get a bandwidth).
// This is threadsafe for one writer, along with `Set`.
// This is fast if there are no unpolled caches. Moreover, its speed is a function of the number of unpolled caches, not the number of caches total.
func (t *UnpolledCaches) SetPolled(results []cache.Result, lastStats dsdata.LastStats) {
	unpolledCaches := copyCaches(t.UnpolledCaches())
	seenCaches := copyCachesTime(*t.seenCaches)
	numUnpolledCaches := len(unpolledCaches)
	if numUnpolledCaches == 0 {
		return
	}
	for cache := range unpolledCaches {
	innerLoop:
		for _, result := range results {
			if result.ID != string(cache) {
				continue
			}

			// TODO fix "whether a cache has ever been polled" to be generic somehow. The result.System.NotAvailable check is duplicated in health.EvalCache, and is fragile. What if another "successfully polled but unavailable" flag were added?
			if !result.Available || result.Error != nil || result.Statistics.NotAvailable {
				log.Debugf("polled %v\n", cache)
				delete(unpolledCaches, cache)
				delete(seenCaches, cache)
				break innerLoop
			}
		}
		lastStat, ok := lastStats.Caches[cache]
		if !ok {
			continue
		}

		if lastStat.Bytes.PerSec != 0 {
			log.Debugf("polled %v\n", cache)
			delete(unpolledCaches, cache)
			delete(seenCaches, cache)
		} else {
			if _, ok := seenCaches[cache]; !ok {
				seenCaches[cache] = lastStat.Bytes.Time
			}
		}

		if seenTime, ok := seenCaches[cache]; ok && time.Since(seenTime) > PolledBytesPerSecTimeout {
			log.Debugf("polled %v (byte change timed out)\n", cache)
			delete(unpolledCaches, cache)
			delete(seenCaches, cache)
		}
	}

	if len(unpolledCaches) == numUnpolledCaches {
		return
	}
	t.setUnpolledCaches(unpolledCaches)
	t.setSeenCaches(seenCaches)
	if len(unpolledCaches) != 0 {
		log.Infof("remaining unpolled %v\n", unpolledCaches)
	} else {
		log.Infof("all caches polled, ready to serve!\n")
	}
}
