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

	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// CopyTrafficMonitorConfigMap returns a deep copy of the given TrafficMonitorConfigMap
func CopyTrafficMonitorConfigMap(a *to.TrafficMonitorConfigMap) to.TrafficMonitorConfigMap {
	b := to.TrafficMonitorConfigMap{}
	b.TrafficServer = map[string]to.TrafficServer{}
	b.CacheGroup = map[string]to.TMCacheGroup{}
	b.Config = map[string]interface{}{}
	b.TrafficMonitor = map[string]to.TrafficMonitor{}
	b.DeliveryService = map[string]to.TMDeliveryService{}
	b.Profile = map[string]to.TMProfile{}
	for k, v := range a.TrafficServer {
		b.TrafficServer[k] = v
	}
	for k, v := range a.CacheGroup {
		b.CacheGroup[k] = v
	}
	for k, v := range a.Config {
		b.Config[k] = v
	}
	for k, v := range a.TrafficMonitor {
		b.TrafficMonitor[k] = v
	}
	for k, v := range a.DeliveryService {
		b.DeliveryService[k] = v
	}
	for k, v := range a.Profile {
		b.Profile[k] = v
	}
	return b
}

// TrafficMonitorConfigMapThreadsafe encapsulates a TrafficMonitorConfigMap safe for multiple readers and a single writer.
type TrafficMonitorConfigMap struct {
	monitorConfig *to.TrafficMonitorConfigMap
	m             *sync.RWMutex
}

// NewTrafficMonitorConfigMap returns an encapsulated TrafficMonitorConfigMap safe for multiple readers and a single writer.
func NewTrafficMonitorConfigMap() TrafficMonitorConfigMap {
	return TrafficMonitorConfigMap{monitorConfig: &to.TrafficMonitorConfigMap{}, m: &sync.RWMutex{}}
}

// Get returns the TrafficMonitorConfigMap. Callers MUST NOT modify, it is not threadsafe for mutation. If mutation is necessary, call CopyTrafficMonitorConfigMap().
func (t *TrafficMonitorConfigMap) Get() to.TrafficMonitorConfigMap {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.monitorConfig
}

// Set sets the TrafficMonitorConfigMap. This is only safe for one writer. This MUST NOT be called by multiple threads.
func (t *TrafficMonitorConfigMap) Set(c to.TrafficMonitorConfigMap) {
	t.m.Lock()
	*t.monitorConfig = c
	t.m.Unlock()
}
