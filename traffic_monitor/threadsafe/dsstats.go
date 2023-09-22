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

// DSStats wraps a deliveryservice.Stats object to be safe for multiple reader goroutines and a single writer.
type DSStats struct {
	dsStats *dsdata.Stats
	m       *sync.RWMutex
}

// DSStatsReader permits reading of a dsdata.Stats object, but not writing. This is designed so a Stats object can safely be passed to multiple goroutines, without worry one may unsafely write.
type DSStatsReader interface {
	Get() dsdata.StatsReadonly
}

// NewDSStats returns a deliveryservice.Stats object wrapped to be safe for multiple readers and a single writer.
func NewDSStats() DSStats {
	return DSStats{m: &sync.RWMutex{}, dsStats: dsdata.NewStats(0)}
}

// Get returns a Stats object safe for reading by multiple goroutines
func (o *DSStats) Get() dsdata.StatsReadonly {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.dsStats
}

// Set sets the internal Stats object. This MUST NOT be called by multiple goroutines.
func (o *DSStats) Set(newDsStats dsdata.Stats) {
	o.m.Lock()
	*o.dsStats = newDsStats
	o.m.Unlock()
}
