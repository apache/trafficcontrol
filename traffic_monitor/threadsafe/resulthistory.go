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

// ResultHistory provides safe access for multiple goroutines readers and a single writer to a stored ResultHistory object.
// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
type ResultHistory struct {
	resultHistory *cache.ResultHistory
	m             *sync.RWMutex
}

// NewResultHistory returns a new ResultHistory safe for multiple readers and a single writer.
func NewResultHistory() ResultHistory {
	h := cache.ResultHistory{}
	return ResultHistory{m: &sync.RWMutex{}, resultHistory: &h}
}

// Get returns the ResultHistory. Callers MUST NOT modify. If mutation is necessary, call ResultHistory.Copy()
func (h *ResultHistory) Get() cache.ResultHistory {
	h.m.RLock()
	defer h.m.RUnlock()
	return *h.resultHistory
}

// Set sets the internal ResultHistory. This is only safe for one thread of execution. This MUST NOT be called from multiple threads.
func (h *ResultHistory) Set(v cache.ResultHistory) {
	h.m.Lock()
	*h.resultHistory = v
	h.m.Unlock()
}
