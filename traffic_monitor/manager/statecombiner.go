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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/health"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

// StartStateCombiner starts the State Combiner goroutine, and returns the threadsafe CombinedStates, and a func to signal to combine states.
func StartStateCombiner(events health.ThreadsafeEvents, peerStates peer.CRStatesPeersThreadsafe, localStates peer.CRStatesThreadsafe, toData todata.TODataThreadsafe) (peer.CRStatesThreadsafe, func()) {
	combinedStates := peer.NewCRStatesThreadsafe()

	// the chan buffer just reduces the number of goroutines on our infinite buffer hack in combineState(), no real writer will block, since combineState() writes in a goroutine.
	combineStateChan := make(chan struct{}, 1)
	combineState := func() {
		select {
		case combineStateChan <- struct{}{}:
		default:
		}
	}

	go func() {
		overrideMap := map[tc.CacheName]bool{}
		for range combineStateChan {
			combineCrStates(events, peerStates.GetCRStatesPeersInfo(), localStates.Get(), combinedStates, overrideMap, toData.Get())
		}
	}()

	return combinedStates, combineState
}

func combineCacheState(
	cacheName tc.CacheName,
	localCacheState tc.IsAvailable,
	events health.ThreadsafeEvents,
	peerCrStatesInfo peer.CRStatesPeersInfo,
	combinedStates peer.CRStatesThreadsafe,
	overrideMap map[tc.CacheName]bool,
	toData todata.TOData,
) {

	overrideCondition := ""
	available := localCacheState.Ipv4Available || localCacheState.Ipv6Available
	ipv4Available := localCacheState.Ipv4Available
	ipv6Available := localCacheState.Ipv6Available
	override := overrideMap[cacheName]

	if localCacheState.Ipv4Available && localCacheState.Ipv6Available {
		// we don't care about the peers, we got a "good one", and we're optimistic
		if override {
			overrideCondition = "cleared; healthy locally"
			overrideMap[cacheName] = false
		}
	} else if !peerCrStatesInfo.HasAvailablePeers() {
		if override {
			overrideCondition = "irrelevant; no peers online"
			overrideMap[cacheName] = false
		}
	} else {
		onlineOnPeers := make([]string, 0)
		ipv4OnlineOnPeers := make([]string, 0)
		ipv6OnlineOnPeers := make([]string, 0)

		for peerName, peerCrStates := range peerCrStatesInfo.GetCrStates() {
			if peerCrStatesInfo.GetPeerAvailability(peerName) {
				if peerCrStates.Caches[cacheName].IsAvailable {
					onlineOnPeers = append(onlineOnPeers, peerName.String())
				}
				if peerCrStates.Caches[cacheName].Ipv4Available {
					ipv4OnlineOnPeers = append(ipv4OnlineOnPeers, peerName.String())
				}
				if peerCrStates.Caches[cacheName].Ipv6Available {
					ipv6OnlineOnPeers = append(ipv6OnlineOnPeers, peerName.String())
				}
			}
		}

		if len(onlineOnPeers) > 0 {
			available = true
			ipv4Available = ipv4Available || len(ipv4OnlineOnPeers) > 0 // optimistically accept true from local or peer
			ipv6Available = ipv6Available || len(ipv6OnlineOnPeers) > 0 // optimistically accept true from local or peer

			if !override {
				overrideCondition = fmt.Sprintf("detected; healthy on (at least) %s", strings.Join(onlineOnPeers, ", "))
				overrideMap[cacheName] = true
			}
		} else if override {
			overrideCondition = "irrelevant; not online on any peers"
			overrideMap[cacheName] = false
		}
	}

	if overrideCondition != "" {
		events.Add(
			health.Event{
				Time:          health.Time(time.Now()),
				Description:   fmt.Sprintf("Health protocol override condition %s", overrideCondition),
				Name:          cacheName.String(),
				Hostname:      cacheName.String(),
				Type:          toData.ServerTypes[cacheName].String(),
				Available:     available,
				IPv4Available: ipv4Available,
				IPv6Available: ipv6Available})
	}

	combinedStates.AddCache(cacheName, tc.IsAvailable{IsAvailable: available, Ipv4Available: ipv4Available, Ipv6Available: ipv6Available, DirectlyPolled: localCacheState.DirectlyPolled, Status: localCacheState.Status, LastPoll: localCacheState.LastPoll})
}

func combineDSState(
	deliveryServiceName tc.DeliveryServiceName,
	localDeliveryService tc.CRStatesDeliveryService,
	peerCrStatesInfo peer.CRStatesPeersInfo,
	combinedStates peer.CRStatesThreadsafe,
) {
	deliveryService := tc.CRStatesDeliveryService{IsAvailable: false, DisabledLocations: []tc.CacheGroupName{}} // important to initialize DisabledLocations, so JSON is `[]` not `null`
	if localDeliveryService.IsAvailable {
		deliveryService.IsAvailable = true
	}

	for peerName, iPeerStates := range peerCrStatesInfo.GetCrStates() {
		peerDeliveryService, ok := iPeerStates.DeliveryService[deliveryServiceName]
		if !ok {
			log.Infof("local delivery service %s not found in peer %s\n", deliveryServiceName, peerName)
			continue
		}
		if peerDeliveryService.IsAvailable {
			deliveryService.IsAvailable = true
		}
	}
	combinedStates.SetDeliveryService(deliveryServiceName, deliveryService)
}

// pruneCombinedDSState deletes delivery services in combined states which have been removed from localStates and peerStates
func pruneCombinedDSState(combinedStates peer.CRStatesThreadsafe, localStates tc.CRStates, peerCrStatesInfo peer.CRStatesPeersInfo) {
	combinedCRStates := combinedStates.Get()

	// remove any DS in combinedStates NOT in local states or peer states
	for deliveryServiceName := range combinedCRStates.DeliveryService {
		inPeer := false
		inLocal := false
		for _, iPeerStates := range peerCrStatesInfo.GetCrStates() {
			if _, ok := iPeerStates.DeliveryService[deliveryServiceName]; ok {
				inPeer = true
				break
			}
		}

		if _, ok := localStates.DeliveryService[deliveryServiceName]; ok {
			inLocal = true
		}

		if !inPeer && !inLocal {
			combinedStates.DeleteDeliveryService(deliveryServiceName)
		}
	}
}

// pruneCombinedCaches deletes caches in combined states which have been removed from localStates.
func pruneCombinedCaches(combinedStates peer.CRStatesThreadsafe, localStates tc.CRStates) {
	combinedCaches := combinedStates.GetCaches()
	for cacheName := range combinedCaches {
		if _, ok := localStates.Caches[cacheName]; !ok {
			combinedStates.DeleteCache(cacheName)
		}
	}
}

func combineCrStates(events health.ThreadsafeEvents, peerCrStatesInfo peer.CRStatesPeersInfo, localStates tc.CRStates, combinedStates peer.CRStatesThreadsafe, overrideMap map[tc.CacheName]bool, toData todata.TOData) {
	for cacheName, localCacheState := range localStates.Caches { // localStates gets pruned when servers are disabled, it's the source of truth
		combineCacheState(cacheName, localCacheState, events, peerCrStatesInfo, combinedStates, overrideMap, toData)
	}

	for deliveryServiceName, localDeliveryService := range localStates.DeliveryService {
		combineDSState(deliveryServiceName, localDeliveryService, peerCrStatesInfo, combinedStates)
	}

	pruneCombinedDSState(combinedStates, localStates, peerCrStatesInfo)
	pruneCombinedCaches(combinedStates, localStates)
}

// CacheGroupNameSlice is a slice of cache names, which fulfills the `sort.Interface` interface.
type CacheGroupNameSlice []tc.CacheGroupName

func (p CacheGroupNameSlice) Len() int           { return len(p) }
func (p CacheGroupNameSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p CacheGroupNameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
