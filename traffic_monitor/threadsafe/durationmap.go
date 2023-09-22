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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// DurationMap wraps a map[tc.CacheName]time.Duration in an object safe for a single writer and multiple readers.
type DurationMap struct {
	durationMap *map[tc.CacheName]time.Duration
	m           *sync.RWMutex
}

// CopyDurationMap copies this duration map.
func CopyDurationMap(a map[tc.CacheName]time.Duration) map[tc.CacheName]time.Duration {
	b := make(map[tc.CacheName]time.Duration, len(a))
	for k, v := range a {
		b[k] = v
	}
	return b
}

// NewDurationMap returns a new DurationMap safe for multiple readers and a single writer goroutine.
func NewDurationMap() DurationMap {
	m := map[tc.CacheName]time.Duration{}
	return DurationMap{m: &sync.RWMutex{}, durationMap: &m}
}

// Get returns the duration map. Callers MUST NOT mutate. If mutation is necessary, call DurationMap.Copy().
func (o *DurationMap) Get() map[tc.CacheName]time.Duration {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.durationMap
}

// Set sets the internal duration map. This MUST NOT be called by multiple goroutines.
func (o *DurationMap) Set(d map[tc.CacheName]time.Duration) {
	o.m.Lock()
	*o.durationMap = d
	o.m.Unlock()
}
