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
	"github.com/apache/trafficcontrol/traffic_monitor/towrap"
	"sync"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

func TestUpdateStatusAnycast(t *testing.T) {

	localStates := peer.NewCRStatesThreadsafe()
	sameIpServers := map[tc.CacheName]map[tc.CacheName]bool{}

	toData := todata.NewThreadsafe()
	towrap.TrafficOpsSessionThreadsafe{}
	toData.Update()

	TODataThreadsafe{m: &sync.RWMutex{}, toData: New()}

	localStates.AddCache("available", tc.IsAvailable{IsAvailable: true, Ipv4Available: true, Ipv6Available: true})
	localStates.AddCache("tooHigh", tc.IsAvailable{IsAvailable: false, Ipv4Available: true, Ipv6Available: true, Status: "too high"})
	toData.
	crStates := updateStatusAnycast(localStates, toData)

	/*
			localStatesC := localStates.Get()
		toDataC := toData.Get()

		for cache, _ := range localStatesC.Caches {
			if _, ok := toDataC.SameIpServers[cache]; ok {
				// all servers with same ip must be available if they are in reported state
				allAvailableV4 := true
				allAvailableV6 := true
				allIsAvailable := true
				for partner, _ := range toDataC.SameIpServers[cache] {
					if partnerState, ok := localStatesC.Caches[partner]; ok {
						// a partner host is reported but is marked down for too high traffic or load
						// this host also needs to be marked down to divert all traffic for their
						// common anycast ip
						if tc.CacheStatusFromString(partnerState.Status) == tc.CacheStatusReported &&
							strings.Contains(partnerState.Status, "too high") {

							if !partnerState.Ipv4Available {
								allAvailableV4 = false
							}
							if !partnerState.Ipv6Available {
								allAvailableV6 = false
							}
							if !partnerState.IsAvailable {
								allIsAvailable = false
							}
							if !allAvailableV4 && !allAvailableV6 && !allIsAvailable {
								break
							}
						}
					}
				}
				newIsAvailable := tc.IsAvailable{}
				newIsAvailable.DirectlyPolled = localStatesC.Caches[cache].DirectlyPolled
				newIsAvailable.Status = localStatesC.Caches[cache].Status
				newIsAvailable.LastPoll = localStatesC.Caches[cache].LastPoll
				newIsAvailable.LastPollV6 = localStatesC.Caches[cache].LastPollV6
				newIsAvailable.IsAvailable = localStatesC.Caches[cache].IsAvailable && allIsAvailable
				newIsAvailable.Ipv4Available = localStatesC.Caches[cache].Ipv4Available && allAvailableV4
				newIsAvailable.Ipv6Available = localStatesC.Caches[cache].Ipv6Available && allAvailableV6

				localStatesC.Caches[cache] = newIsAvailable
			}
		}
		return localStatesC

	*/
}
