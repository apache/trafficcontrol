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

	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
)

// LastStats wraps a deliveryservice.LastStats object to be safe for multiple readers and one writer.
type LastStats struct {
	stats *dsdata.LastStats
	m     *sync.RWMutex
}

// NewLastStats returns a wrapped a deliveryservice.LastStats object safe for multiple readers and one writer.
func NewLastStats() LastStats {
	return LastStats{m: &sync.RWMutex{}, stats: dsdata.NewLastStats(0, 0)}
}

// Get returns the last KBPS stats object. Callers MUST NOT modify the object. It is not threadsafe for writing. If the object must be modified, callers must call LastStats.Copy() and modify the copy.
func (o LastStats) Get() dsdata.LastStats {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.stats
}

// Set sets the internal LastStats object. This MUST NOT be called by multiple goroutines.
func (o LastStats) Set(s dsdata.LastStats) {
	o.m.Lock()
	*o.stats = s
	o.m.Unlock()
}
