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

	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
)

// CacheKbpses wraps a map of cache kbps measurements to be safe for multiple
// reader goroutines and one writer.
type CacheKbpses struct {
	v *cache.Kbpses
	m *sync.RWMutex
}

// NewCacheKbpses creates and returns a new, empty CacheKbpses,
// initializing internal pointer values.
func NewCacheKbpses() CacheKbpses {
	v := cache.Kbpses{}
	return CacheKbpses{m: &sync.RWMutex{}, v: &v}
}

// Get returns the internal map of cache kbps measurements. The returned map
// MUST NOT be modified. If modification is necessary, copy.
func (o *CacheKbpses) Get() cache.Kbpses {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.v
}

// Set sets the internal map of cache kbps measurements. This MUST NOT be called
// by multiple goroutines.
func (o *CacheKbpses) Set(v cache.Kbpses) {
	o.m.Lock()
	*o.v = v
	o.m.Unlock()
}
