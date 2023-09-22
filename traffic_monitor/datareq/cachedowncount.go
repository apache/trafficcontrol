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

package datareq

import (
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
)

func srvAPICacheDownCount(localStates peer.CRStatesThreadsafe, monitorConfig threadsafe.TrafficMonitorConfigMap) []byte {
	return []byte(strconv.Itoa(cacheDownCount(localStates.Get().Caches, monitorConfig.Get().TrafficServer)))
}

// cacheOfflineCount returns the total reported caches marked down, excluding status offline and admin_down.
func cacheDownCount(caches map[tc.CacheName]tc.IsAvailable, toServers map[string]tc.TrafficServer) int {
	count := 0
	for cache, available := range caches {
		if !available.IsAvailable && tc.CacheStatusFromString(toServers[string(cache)].ServerStatus) == tc.CacheStatusReported {
			count++
		}
	}
	return count
}
