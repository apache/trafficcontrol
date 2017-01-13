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
	"sync"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/cache"
)

// ResultStatHistory provides safe access for multiple goroutines readers and a single writer to a stored HistoryHistory object.
// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
type ResultStatHistory struct {
	history *cache.ResultStatHistory
	m       *sync.RWMutex
}

// NewResultStatHistory returns a new ResultStatHistory safe for multiple readers and a single writer.
func NewResultStatHistory() ResultStatHistory {
	h := cache.ResultStatHistory{}
	return ResultStatHistory{m: &sync.RWMutex{}, history: &h}
}

// Get returns the ResultStatHistory. Callers MUST NOT modify. If mutation is necessary, call ResultStatHistory.Copy()
func (h *ResultStatHistory) Get() cache.ResultStatHistory {
	h.m.RLock()
	defer h.m.RUnlock()
	return *h.history
}

// Set sets the internal ResultStatHistory. This is only safe for one thread of execution. This MUST NOT be called from multiple threads.
func (h *ResultStatHistory) Set(v cache.ResultStatHistory) {
	h.m.Lock()
	*h.history = v
	h.m.Unlock()
}

// ResultStatHistory provides safe access for multiple goroutines readers and a single writer to a stored HistoryHistory object.
// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
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
