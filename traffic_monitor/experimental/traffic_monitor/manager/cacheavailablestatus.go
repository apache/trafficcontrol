package manager

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
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"sync"
)

// CacheAvailableStatusReported is the status string returned by caches set to "reported" in Traffic Ops.
// TODO put somewhere more generic
const CacheAvailableStatusReported = "REPORTED"

// CacheAvailableStatus is the available status of the given cache. It includes a boolean available/unavailable flag, and a descriptive string.
type CacheAvailableStatus struct {
	Available bool
	Status    string
}

// CacheAvailableStatuses is the available status of each cache.
type CacheAvailableStatuses map[enum.CacheName]CacheAvailableStatus

// CacheAvailableStatusThreadsafe wraps a map of cache available statuses to be safe for multiple reader goroutines and one writer.
type CacheAvailableStatusThreadsafe struct {
	caches *CacheAvailableStatuses
	m      *sync.RWMutex
}

// Copy copies this CacheAvailableStatuses. It does not modify, and thus is safe for multiple reader goroutines.
func (a CacheAvailableStatuses) Copy() CacheAvailableStatuses {
	b := CacheAvailableStatuses(map[enum.CacheName]CacheAvailableStatus{})
	for k, v := range a {
		b[k] = v
	}
	return b
}

// NewCacheAvailableStatusThreadsafe creates and returns a new CacheAvailableStatusThreadsafe, initializing internal pointer values.
func NewCacheAvailableStatusThreadsafe() CacheAvailableStatusThreadsafe {
	c := CacheAvailableStatuses(map[enum.CacheName]CacheAvailableStatus{})
	return CacheAvailableStatusThreadsafe{m: &sync.RWMutex{}, caches: &c}
}

// Get returns the internal map of cache statuses. The returned map MUST NOT be modified. If modification is necessary, copy.
func (o *CacheAvailableStatusThreadsafe) Get() CacheAvailableStatuses {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.caches
}

// Set sets the internal map of cache availability. This MUST NOT be called by multiple goroutines.
func (o *CacheAvailableStatusThreadsafe) Set(v CacheAvailableStatuses) {
	o.m.Lock()
	*o.caches = v
	o.m.Unlock()
}
