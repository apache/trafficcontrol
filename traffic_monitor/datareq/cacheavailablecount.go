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
)

func srvAPICacheAvailableCount(localStates peer.CRStatesThreadsafe) []byte {
	return []byte(strconv.Itoa(cacheAvailableCount(localStates.Get().Caches)))
}

// cacheOfflineCount returns the total caches not available, including marked unavailable, status offline, and status admin_down
func cacheOfflineCount(caches map[tc.CacheName]tc.IsAvailable) int {
	count := 0
	for _, available := range caches {
		if !available.IsAvailable {
			count++
		}
	}
	return count
}

// cacheAvailableCount returns the total caches available, including marked available and status online
func cacheAvailableCount(caches map[tc.CacheName]tc.IsAvailable) int {
	return len(caches) - cacheOfflineCount(caches)
}
