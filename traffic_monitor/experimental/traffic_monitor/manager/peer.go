package manager

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
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/util"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/health"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
)

// StartPeerManager listens for peer results, and when it gets one, it adds it to the peerStates list, and optimistically combines the good results into combinedStates
func StartPeerManager(
	peerChan <-chan peer.Result,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	events threadsafe.Events,
	peerOptimistic bool,
	toData todata.TODataThreadsafe,
) peer.CRStatesThreadsafe {
	combinedStates := peer.NewCRStatesThreadsafe()
	overrideMap := map[enum.CacheName]bool{}

	go func() {
		for peerResult := range peerChan {
			comparePeerState(events, peerResult, peerStates)
			peerStates.Set(peerResult)
			combineCrStates(events, peerOptimistic, peerStates, localStates.Get(), combinedStates, overrideMap, toData)
			peerResult.PollFinished <- peerResult.PollID
		}
	}()
	return combinedStates
}

func comparePeerState(events threadsafe.Events, result peer.Result, peerStates peer.CRStatesPeersThreadsafe) {
	if result.Available != peerStates.GetPeerAvailability(result.ID) {
		events.Add(health.Event{Time: result.Time, Unix: result.Time.Unix(), Description: util.JoinErrorsString(result.Errors), Name: result.ID.String(), Hostname: result.ID.String(), Type: "Peer", Available: result.Available})
	}
}

// TODO JvD: add deliveryservice stuff
func combineCrStates(events threadsafe.Events, peerOptimistic bool, peerStates peer.CRStatesPeersThreadsafe, localStates peer.Crstates, combinedStates peer.CRStatesThreadsafe, overrideMap map[enum.CacheName]bool, toData todata.TODataThreadsafe) {
	toDataCopy := toData.Get()

	for cacheName, localCacheState := range localStates.Caches { // localStates gets pruned when servers are disabled, it's the source of truth
		var overrideCondition string
		available := false
		override := overrideMap[cacheName]

		if localCacheState.IsAvailable {
			available = true // we don't care about the peers, we got a "good one", and we're optimistic

			if override {
				overrideCondition = "cleared; healthy locally"
				overrideMap[cacheName] = false
			}
		} else if peerOptimistic {
			if !peerStates.HasAvailablePeers() {
				if override {
					overrideCondition = "irrelevant; no peers online"
					overrideMap[cacheName] = false
				}
			} else {
				onlineOnPeers := make([]string, 0)

				for peer, peerCrStates := range peerStates.GetCrstates() {
					if peerStates.GetPeerAvailability(peer) {
						if peerCrStates.Caches[cacheName].IsAvailable {
							onlineOnPeers = append(onlineOnPeers, peer.String())
						}
					}
				}

				if len(onlineOnPeers) > 0 {
					available = true

					if !override {
						overrideCondition = fmt.Sprintf("detected; healthy on (at least) %s", strings.Join(onlineOnPeers, ", "))
						overrideMap[cacheName] = true
					}
				} else {
					if override {
						overrideCondition = "irrelevant; not online on any peers"
						overrideMap[cacheName] = false
					}
				}
			}
		}

		if overrideCondition != "" {
			events.Add(health.Event{Time: time.Now(), Unix: time.Now().Unix(), Description: fmt.Sprintf("Health protocol override condition %s", overrideCondition), Name: cacheName.String(), Hostname: cacheName.String(), Type: toDataCopy.ServerTypes[cacheName].String(), Available: available})
		}

		combinedStates.SetCache(cacheName, peer.IsAvailable{IsAvailable: available})
	}

	for deliveryServiceName, localDeliveryService := range localStates.Deliveryservice {
		deliveryService := peer.Deliveryservice{IsAvailable: false, DisabledLocations: []enum.CacheName{}} // important to initialize DisabledLocations, so JSON is `[]` not `null`
		if localDeliveryService.IsAvailable {
			deliveryService.IsAvailable = true
		}
		deliveryService.DisabledLocations = localDeliveryService.DisabledLocations

		for peerName, iPeerStates := range peerStates.GetCrstates() {
			peerDeliveryService, ok := iPeerStates.Deliveryservice[deliveryServiceName]
			if !ok {
				log.Warnf("local delivery service %s not found in peer %s\n", deliveryServiceName, peerName)
				continue
			}
			if peerDeliveryService.IsAvailable {
				deliveryService.IsAvailable = true
			}
			deliveryService.DisabledLocations = intersection(deliveryService.DisabledLocations, peerDeliveryService.DisabledLocations)
		}
		combinedStates.SetDeliveryService(deliveryServiceName, deliveryService)
	}
}

// CacheNameSlice is a slice of cache names, which fulfills the `sort.Interface` interface.
type CacheNameSlice []enum.CacheName

func (p CacheNameSlice) Len() int           { return len(p) }
func (p CacheNameSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p CacheNameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// intersection returns strings in both a and b.
// Note this modifies a and b. Specifically, it sorts them. If that isn't acceptable, pass copies of your real data.
func intersection(a []enum.CacheName, b []enum.CacheName) []enum.CacheName {
	sort.Sort(CacheNameSlice(a))
	sort.Sort(CacheNameSlice(b))
	c := []enum.CacheName{} // important to initialize, so JSON is `[]` not `null`
	for _, s := range a {
		i := sort.Search(len(b), func(i int) bool { return b[i] >= s })
		if i < len(b) && b[i] == s {
			c = append(c, s)
		}
	}
	return c
}
