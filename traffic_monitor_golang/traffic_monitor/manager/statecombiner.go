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

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/health"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
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
		overrideMap := map[enum.CacheName]bool{}
		for range combineStateChan {
			drain(combineStateChan)
			combineCrStates(events, true, peerStates, localStates.Get(), combinedStates, overrideMap, toData.Get())
		}
	}()

	return combinedStates, combineState
}

func combineCacheState(cacheName enum.CacheName, localCacheState peer.IsAvailable, events health.ThreadsafeEvents, peerOptimistic bool, peerStates peer.CRStatesPeersThreadsafe, localStates peer.Crstates, combinedStates peer.CRStatesThreadsafe, overrideMap map[enum.CacheName]bool, toData todata.TOData) {
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

	combinedStates.SetCache(cacheName, peer.IsAvailable{IsAvailable: available})
}

func combineDSState(deliveryServiceName enum.DeliveryServiceName, localDeliveryService peer.Deliveryservice, events health.ThreadsafeEvents, peerOptimistic bool, peerStates peer.CRStatesPeersThreadsafe, localStates peer.Crstates, combinedStates peer.CRStatesThreadsafe, overrideMap map[enum.CacheName]bool, toData todata.TOData) {
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

func combineCrStates(events health.ThreadsafeEvents, peerOptimistic bool, peerStates peer.CRStatesPeersThreadsafe, localStates peer.Crstates, combinedStates peer.CRStatesThreadsafe, overrideMap map[enum.CacheName]bool, toData todata.TOData) {
	for cacheName, localCacheState := range localStates.Caches { // localStates gets pruned when servers are disabled, it's the source of truth
		combineCacheState(cacheName, localCacheState, events, peerOptimistic, peerStates, localStates, combinedStates, overrideMap, toData)
	}
	for deliveryServiceName, localDeliveryService := range localStates.Deliveryservice {
		combineDSState(deliveryServiceName, localDeliveryService, events, peerOptimistic, peerStates, localStates, combinedStates, overrideMap, toData)
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
