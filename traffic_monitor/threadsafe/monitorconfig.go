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

	tc "github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// TrafficMonitorConfigMapThreadsafe encapsulates a LegacyTrafficMonitorConfigMap safe for multiple readers and a single writer.
type TrafficMonitorConfigMap struct {
	monitorConfig *tc.TrafficMonitorConfigMap
	m             *sync.RWMutex
}

// NewTrafficMonitorConfigMap returns an encapsulated LegacyTrafficMonitorConfigMap safe for multiple readers and a single writer.
func NewTrafficMonitorConfigMap() TrafficMonitorConfigMap {
	return TrafficMonitorConfigMap{monitorConfig: &tc.TrafficMonitorConfigMap{}, m: &sync.RWMutex{}}
}

// Get returns the LegacyTrafficMonitorConfigMap. Callers MUST NOT modify, it is not threadsafe for mutation. If mutation is necessary, call CopyTrafficMonitorConfigMap().
func (t *TrafficMonitorConfigMap) Get() tc.TrafficMonitorConfigMap {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.monitorConfig
}

// Set sets the LegacyTrafficMonitorConfigMap. This is only safe for one writer. This MUST NOT be called by multiple threads.
func (t *TrafficMonitorConfigMap) Set(c tc.TrafficMonitorConfigMap) {
	t.m.Lock()
	*t.monitorConfig = c
	t.m.Unlock()
}
