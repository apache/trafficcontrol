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
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/health"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

func srvTRState(
	params url.Values,
	localStates peer.CRStatesThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	distributedPollingEnabled bool,
	toData todata.TODataThreadsafe,
) ([]byte, int, error) {
	_, raw := params["raw"]     // peer polling case
	_, local := params["local"] // distributed peer polling case
	if raw {
		data, err := srvTRStateSelf(localStates, distributedPollingEnabled, toData)
		return data, http.StatusOK, err
	}

	// This covers the case where we have lost connectivity to all peers, but multiple peers exist. In this case, it is
	// more likely that the local machine has lost all connectivity than both peers losing connectivity or crashing. If
	// the peers really did crash, the health protocol is essentially broken, and serving a 503 will cause Traffic Router
	// to use the last good state fetched from a Traffic Monitor within the CDN. If the peers are simply unreachable from
	// this Traffic Monitor, serving 503s until connectivity is restored will cause Traffic Router to ignore this instance
	// until the health protocol can be relied upon once again.
	if peerStates.OptimisticQuorumEnabled() {
		optimisticQuorum, peersAvailable, peerCount, minimum := peerStates.HasOptimisticQuorum()
		log.Debugf("optimisticQuorum=%v, peerCount=%v, peersAvailable=%v, minimum=%v", optimisticQuorum, peerCount, peersAvailable, minimum)

		if !optimisticQuorum {
			return nil, http.StatusServiceUnavailable, fmt.Errorf("number of peers available (%d/%d) is less than the minimum number of %d required for optimistic peer quorum", peersAvailable, peerCount, minimum)
		}
	}

	data, err := srvTRStateDerived(combinedStates, local && distributedPollingEnabled, toData)

	return data, http.StatusOK, err
}

func srvTRStateDerived(combinedStates peer.CRStatesThreadsafe, directlyPolledOnly bool, toData todata.TODataThreadsafe) ([]byte, error) {
	if !directlyPolledOnly {
		combinedStatesC := updateStatusSameIpServers(combinedStates, toData)
		return tc.CRStatesMarshall(combinedStatesC)
	}
	unfiltered := updateStatusSameIpServers(combinedStates, toData)
	return tc.CRStatesMarshall(filterDirectlyPolledCaches(unfiltered))
}

func filterDirectlyPolledCaches(crstates tc.CRStates) tc.CRStates {
	filtered := tc.CRStates{
		Caches:          make(map[tc.CacheName]tc.IsAvailable),
		DeliveryService: crstates.DeliveryService,
	}
	for cacheName, availability := range crstates.Caches {
		if availability.DirectlyPolled {
			filtered.Caches[cacheName] = availability
		}
	}
	return filtered
}

func srvTRStateSelf(localStates peer.CRStatesThreadsafe, directlyPolledOnly bool, toData todata.TODataThreadsafe) ([]byte, error) {
	if !directlyPolledOnly {
		localStatesC := updateStatusSameIpServers(localStates, toData)
		return tc.CRStatesMarshall(localStatesC)
	}
	unfiltered := updateStatusSameIpServers(localStates, toData)
	return tc.CRStatesMarshall(filterDirectlyPolledCaches(unfiltered))
}

func updateStatusSameIpServers(localStates peer.CRStatesThreadsafe, toData todata.TODataThreadsafe) tc.CRStates {
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
					// a partner host is reported but is marked down for exceeding a threshold
					// this host also needs to be marked down to divert all traffic for their
					// common ip
					if strings.Contains(partnerState.Status, string(tc.CacheStatusReported)) &&
						strings.Contains(partnerState.Status, health.TooHigh.String()) {
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
}
