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
	"math/rand"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/health"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

func TestCombineCacheState(t *testing.T) {
	cacheName := tc.CacheName("testCache")
	localCacheStates := []tc.IsAvailable{
		tc.IsAvailable{
			IsAvailable:   true,
			Ipv4Available: true,
			Ipv6Available: true,
		},
		tc.IsAvailable{
			IsAvailable:   true,
			Ipv4Available: false,
			Ipv6Available: true,
		},
		tc.IsAvailable{
			IsAvailable:   true,
			Ipv4Available: true,
			Ipv6Available: false,
		},
		tc.IsAvailable{
			IsAvailable:   false,
			Ipv4Available: false,
			Ipv6Available: false,
		},
	}
	events := health.NewThreadsafeEvents(1)
	peerStates := peer.NewCRStatesPeersThreadsafe(1)
	peerStates.SetTimeout(time.Duration(rand.Int63()))
	peerResult := peer.Result{
		ID:        tc.TrafficMonitorName("TestTM-01"),
		Available: true,
		PeerStates: tc.CRStates{
			Caches: map[tc.CacheName]tc.IsAvailable{
				tc.CacheName(cacheName): tc.IsAvailable{
					IsAvailable:   true,
					Ipv4Available: true,
					Ipv6Available: true,
				},
			},
		},
		Time: time.Now(),
	}
	peerStates.Set(peerResult)
	peerSet := map[tc.TrafficMonitorName]struct{}{
		tc.TrafficMonitorName("TestTM-01"): struct{}{},
	}
	peerStates.SetPeers(peerSet)
	peerStates.SetTimeout(time.Duration(rand.Int()))

	combinedStates := peer.NewCRStatesThreadsafe()
	overrideMap := map[tc.CacheName]bool{}
	overrideMap[cacheName] = false
	toData := todata.TOData{}
	toData.ServerTypes = map[tc.CacheName]tc.CacheType{
		cacheName: tc.CacheTypeEdge,
	}

	for _, localCacheState := range localCacheStates {
		combineCacheState(cacheName, localCacheState, events, peerStates.GetCRStatesPeersInfo(), combinedStates, overrideMap, toData)

		if !combinedStates.Get().Caches[cacheName].IsAvailable {
			t.Fatalf("cache is unavailable and should be available")
		}
		if !combinedStates.Get().Caches[cacheName].Ipv4Available {
			t.Fatalf("cache IPv4 is unavailable and should be available")
		}
		if !combinedStates.Get().Caches[cacheName].Ipv6Available {
			t.Fatalf("cache IPv6 is unavailable and should be available")
		}
	}
}

func TestCombineCacheStateCacheDown(t *testing.T) {
	cacheName := tc.CacheName("testCache")
	localCacheState := tc.IsAvailable{
		IsAvailable:   false,
		Ipv4Available: false,
		Ipv6Available: false,
	}

	events := health.NewThreadsafeEvents(1)
	peerStates := peer.NewCRStatesPeersThreadsafe(1)
	peerStates.SetTimeout(time.Duration(rand.Int63()))
	peerResult := peer.Result{
		ID:        tc.TrafficMonitorName("TestTM-01"),
		Available: true,
		PeerStates: tc.CRStates{
			Caches: map[tc.CacheName]tc.IsAvailable{
				tc.CacheName(cacheName): tc.IsAvailable{
					IsAvailable:   true,
					Ipv4Available: false,
					Ipv6Available: true,
				},
			},
		},
		Time: time.Now(),
	}
	peerStates.Set(peerResult)
	peerSet := map[tc.TrafficMonitorName]struct{}{
		tc.TrafficMonitorName("TestTM-01"): struct{}{},
	}
	peerStates.SetPeers(peerSet)
	peerStates.SetTimeout(time.Duration(rand.Int()))

	combinedStates := peer.NewCRStatesThreadsafe()
	overrideMap := map[tc.CacheName]bool{}
	overrideMap[cacheName] = false
	toData := todata.TOData{}
	toData.ServerTypes = map[tc.CacheName]tc.CacheType{
		cacheName: tc.CacheTypeEdge,
	}

	combineCacheState(cacheName, localCacheState, events, peerStates.GetCRStatesPeersInfo(), combinedStates, overrideMap, toData)

	if !combinedStates.Get().Caches[cacheName].IsAvailable {
		t.Fatalf("cache is unavailable and should be available")
	}
	if combinedStates.Get().Caches[cacheName].Ipv4Available {
		t.Fatalf("cache IPv4 is available and should be unavailable")
	}
	if !combinedStates.Get().Caches[cacheName].Ipv6Available {
		t.Fatalf("cache IPv6 is unavailable and should be available")
	}
}
