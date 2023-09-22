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

// CacheAvailableStatus wraps a map of cache available statuses to be safe for
// multiple reader goroutines and one writer.
type CacheAvailableStatus struct {
	caches *cache.AvailableStatuses
	m      *sync.RWMutex
}

// NewCacheAvailableStatus creates and returns a new CacheAvailableStatus,
// initializing internal pointer values.
func NewCacheAvailableStatus() CacheAvailableStatus {
	c := cache.AvailableStatuses(map[string]cache.AvailableStatus{})
	return CacheAvailableStatus{m: &sync.RWMutex{}, caches: &c}
}

// Get returns the internal map of cache statuses. The returned map MUST NOT be
// modified. If modification is necessary, copy.
func (o *CacheAvailableStatus) Get() cache.AvailableStatuses {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.caches
}

// Set sets the internal map of cache availability. This MUST NOT be called by
// multiple goroutines.
func (o *CacheAvailableStatus) Set(v cache.AvailableStatuses) {
	o.m.Lock()
	*o.caches = v
	o.m.Unlock()
}
