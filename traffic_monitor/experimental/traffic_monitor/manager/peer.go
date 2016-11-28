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
	"sort"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
)

// StartPeerManager listens for peer results, and when it gets one, it adds it to the peerStates list, and optimistically combines the good results into combinedStates
func StartPeerManager(
	peerChan <-chan peer.Result,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
) peer.CRStatesThreadsafe {
	combinedStates := peer.NewCRStatesThreadsafe()
	go func() {
		for crStatesResult := range peerChan {
			peerStates.Set(crStatesResult.ID, crStatesResult.PeerStats)
			combineCrStates(peerStates.Get(), localStates.Get(), combinedStates)
			crStatesResult.PollFinished <- crStatesResult.PollID
		}
	}()
	return combinedStates
}

// TODO JvD: add deliveryservice stuff
func combineCrStates(peerStates map[enum.TrafficMonitorName]peer.Crstates, localStates peer.Crstates, combinedStates peer.CRStatesThreadsafe) {
	for cacheName, localCacheState := range localStates.Caches { // localStates gets pruned when servers are disabled, it's the source of truth
		downVotes := 0 // TODO JvD: change to use parameter when deciding to be optimistic or pessimistic.
		available := false
		if localCacheState.IsAvailable {
			// log.Infof(cacheName, " is available locally - setting to IsAvailable: true")
			available = true // we don't care about the peers, we got a "good one", and we're optimistic
		} else {
			downVotes++ // localStates says it's not happy
			for _, peerCrStates := range peerStates {
				if peerCrStates.Caches[cacheName].IsAvailable {
					// log.Infoln(cacheName, "- locally we think it's down, but", peerName, "says IsAvailable: ", peerCrStates.Caches[cacheName].IsAvailable, "trusting the peer.")
					available = true // we don't care about the peers, we got a "good one", and we're optimistic
					break            // one peer that thinks we're good is all we need.
				} else {
					// log.Infoln(cacheName, "- locally we think it's down, and", peerName, "says IsAvailable: ", peerCrStates.Caches[cacheName].IsAvailable, "down voting")
					downVotes++ // peerStates for this peer doesn't like it
				}
			}
		}
		if downVotes > len(peerStates) {
			// log.Infoln(cacheName, "-", downVotes, "down votes, setting to IsAvailable: false")
			available = false
		}
		combinedStates.SetCache(cacheName, peer.IsAvailable{IsAvailable: available})
	}

	for deliveryServiceName, localDeliveryService := range localStates.Deliveryservice {
		deliveryService := peer.Deliveryservice{IsAvailable: false, DisabledLocations: []enum.CacheName{}} // important to initialize DisabledLocations, so JSON is `[]` not `null`
		if localDeliveryService.IsAvailable {
			deliveryService.IsAvailable = true
		}
		deliveryService.DisabledLocations = localDeliveryService.DisabledLocations

		for peerName, iPeerStates := range peerStates {
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
