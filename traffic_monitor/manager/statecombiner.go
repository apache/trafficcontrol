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
	"github.com/apache/trafficcontrol/lib/go-tc/tce"
	"sort"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/health"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

// StartStateCombiner starts the State Combiner goroutine, and returns the threadsafe CombinedStates, and a func to signal to combine states.
func StartStateCombiner(events health.ThreadsafeEvents, peerStates peer.CRStatesPeersThreadsafe, localStates peer.CRStatesThreadsafe, toData todata.TODataThreadsafe) (peer.CRStatesThreadsafe, func()) {
	combinedStates := peer.NewCRStatesThreadsafe()

	// the chan buffer just reduces the number of goroutines on our infinite buffer hack in combineState(), no real writer will block, since combineState() writes in a goroutine.
	combineStateChan := make(chan struct{}, 5)
	combineState := func() {
		go func() { combineStateChan <- struct{}{} }()
	}

	drain := func(c <-chan struct{}) {
	outer:
		for {
			select {
			case <-c:
				continue
			default:
				break outer
			}
		}
	}

	go func() {
		overrideMap := map[tce.CacheName]bool{}
		for range combineStateChan {
			drain(combineStateChan)
			combineCrStates(events, true, peerStates, localStates.Get(), combinedStates, overrideMap, toData.Get())
		}
	}()

	return combinedStates, combineState
}

func combineCacheState(cacheName tce.CacheName, localCacheState tc.IsAvailable, events health.ThreadsafeEvents, peerOptimistic bool, peerStates peer.CRStatesPeersThreadsafe, localStates tc.CRStates, combinedStates peer.CRStatesThreadsafe, overrideMap map[tce.CacheName]bool, toData todata.TOData) {
	overrideCondition := ""
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
		events.Add(health.Event{Time: health.Time(time.Now()), Description: fmt.Sprintf("Health protocol override condition %s", overrideCondition), Name: cacheName.String(), Hostname: cacheName.String(), Type: toData.ServerTypes[cacheName].String(), Available: available})
	}

	combinedStates.AddCache(cacheName, tc.IsAvailable{IsAvailable: available})
}

func combineDSState(
	deliveryServiceName tce.DeliveryServiceName,
	localDeliveryService tc.CRStatesDeliveryService,
	events health.ThreadsafeEvents,
	peerOptimistic bool,
	peerStates peer.CRStatesPeersThreadsafe,
	localStates tc.CRStates,
	combinedStates peer.CRStatesThreadsafe,
	overrideMap map[tce.CacheName]bool,
	toData todata.TOData,
) {
	deliveryService := tc.CRStatesDeliveryService{IsAvailable: false, DisabledLocations: []tce.CacheGroupName{}} // important to initialize DisabledLocations, so JSON is `[]` not `null`
	if localDeliveryService.IsAvailable {
		deliveryService.IsAvailable = true
	}
	deliveryService.DisabledLocations = localDeliveryService.DisabledLocations

	for peerName, iPeerStates := range peerStates.GetCrstates() {
		peerDeliveryService, ok := iPeerStates.DeliveryService[deliveryServiceName]
		if !ok {
			log.Infof("local delivery service %s not found in peer %s\n", deliveryServiceName, peerName)
			continue
		}
		if peerDeliveryService.IsAvailable {
			deliveryService.IsAvailable = true
		}
		deliveryService.DisabledLocations = intersection(deliveryService.DisabledLocations, peerDeliveryService.DisabledLocations)
	}
	combinedStates.SetDeliveryService(deliveryServiceName, deliveryService)
}

// pruneCombinedDSState deletes delivery services in combined states which have been removed from localStates and peerStates
func pruneCombinedDSState(combinedStates peer.CRStatesThreadsafe, localStates tc.CRStates, peerStates peer.CRStatesPeersThreadsafe) {
	combinedCRStates := combinedStates.Get()

	// remove any DS in combinedStates NOT in local states or peer states
	for deliveryServiceName := range combinedCRStates.DeliveryService {
		inPeer := false
		inLocal := false
		for _, iPeerStates := range peerStates.GetCrstates() {
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
	for cacheName, _ := range combinedCaches {
		if _, ok := localStates.Caches[cacheName]; !ok {
			combinedStates.DeleteCache(cacheName)
		}
	}
}

func combineCrStates(events health.ThreadsafeEvents, peerOptimistic bool, peerStates peer.CRStatesPeersThreadsafe, localStates tc.CRStates, combinedStates peer.CRStatesThreadsafe, overrideMap map[tce.CacheName]bool, toData todata.TOData) {
	for cacheName, localCacheState := range localStates.Caches { // localStates gets pruned when servers are disabled, it's the source of truth
		combineCacheState(cacheName, localCacheState, events, peerOptimistic, peerStates, localStates, combinedStates, overrideMap, toData)
	}

	for deliveryServiceName, localDeliveryService := range localStates.DeliveryService {
		combineDSState(deliveryServiceName, localDeliveryService, events, peerOptimistic, peerStates, localStates, combinedStates, overrideMap, toData)
	}

	pruneCombinedDSState(combinedStates, localStates, peerStates)
	pruneCombinedCaches(combinedStates, localStates)
}

// CacheNameSlice is a slice of cache names, which fulfills the `sort.Interface` interface.
type CacheGroupNameSlice []tce.CacheGroupName

func (p CacheGroupNameSlice) Len() int           { return len(p) }
func (p CacheGroupNameSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p CacheGroupNameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// intersection returns strings in both a and b.
// Note this modifies a and b. Specifically, it sorts them. If that isn't acceptable, pass copies of your real data.
func intersection(a []tce.CacheGroupName, b []tce.CacheGroupName) []tce.CacheGroupName {
	sort.Sort(CacheGroupNameSlice(a))
	sort.Sort(CacheGroupNameSlice(b))
	c := []tce.CacheGroupName{} // important to initialize, so JSON is `[]` not `null`
	for _, s := range a {
		i := sort.Search(len(b), func(i int) bool { return b[i] >= s })
		if i < len(b) && b[i] == s {
			c = append(c, s)
		}
	}
	return c
}
