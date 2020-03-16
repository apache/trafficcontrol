package health

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
	"github.com/apache/trafficcontrol/traffic_monitor/config"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

func TestCalcAvailabilityThresholds(t *testing.T) {
	result := cache.Result{
		ID:    "myCacheName",
		Error: nil,
		Astats: cache.Astats{
			Ats: map[string]interface{}{},
			System: cache.AstatsSystem{
				InfName:           "bond0",
				InfSpeed:          20000,
				ProcNetDev:        "bond0: 1234567891011121 123456789101    0    5    0     0          0  9876543 12345678910111213 1234567891011    0 1234    0     0       0          0",
				ProcLoadavg:       "5.43 4.32 3.21 3/1234 32109",
				ConfigLoadRequest: 9,
				LastReloadRequest: 1559237772,
				ConfigReloads:     1,
				LastReload:        1559237773,
				AstatsLoad:        1559237774,
				NotAvailable:      false,
			},
		},
		Time:            time.Now(),
		RequestTime:     time.Second,
		Vitals:          cache.Vitals{},
		PollID:          42,
		PollFinished:    make(chan uint64, 1),
		PrecomputedData: cache.PrecomputedData{},
		Available:       true,
	}
	GetVitals(&result, nil, nil)

	prevResult := cache.Result{
		Time:   result.Time.Add(time.Second * -1),
		Vitals: cache.Vitals{BytesOut: result.Vitals.BytesOut - 1250000000}, // 10 gigabits
	}
	GetVitals(&result, &prevResult, nil)

	statResultHistory := (*threadsafe.ResultStatHistory)(nil)
	mc := tc.TrafficMonitorConfigMap{
		TrafficServer: map[string]tc.TrafficServer{
			string(result.ID): {
				ServerStatus: string(tc.CacheStatusReported),
				Profile:      "myProfileName",
			},
		},
		Profile: map[string]tc.TMProfile{},
	}
	mc.Profile[mc.TrafficServer[string(result.ID)].Profile] = tc.TMProfile{
		Name: mc.TrafficServer[string(result.ID)].Profile,
		Parameters: tc.TMParameters{
			Thresholds: map[string]tc.HealthThreshold{
				"availableBandwidthInKbps": tc.HealthThreshold{
					Val:        15000000,
					Comparator: ">",
				},
			},
		},
	}

	toData := todata.TOData{
		ServerTypes:            map[tc.CacheName]tc.CacheType{},
		DeliveryServiceServers: map[tc.DeliveryServiceName][]tc.CacheName{},
		ServerCachegroups:      map[tc.CacheName]tc.CacheGroupName{},
	}
	toData.ServerTypes[result.ID] = tc.CacheTypeEdge
	toData.DeliveryServiceServers["myDS"] = []tc.CacheName{result.ID}
	toData.ServerCachegroups[result.ID] = "myCG"

	localCacheStatusThreadsafe := threadsafe.NewCacheAvailableStatus()
	localStates := peer.NewCRStatesThreadsafe()
	events := NewThreadsafeEvents(200)

	// test that a normal stat poll over the kbps threshold marks down

	pollerName := "stat"
	results := []cache.Result{result}
	CalcAvailability(results, pollerName, statResultHistory, mc, toData, localCacheStatusThreadsafe, localStates, events, config.Both)

	localCacheStatuses := localCacheStatusThreadsafe.Get()
	if localCacheStatus, ok := localCacheStatuses[result.ID]; !ok {
		t.Fatalf("expected: localCacheStatus[cacheName], actual: missing")
	} else if localCacheStatus.Available.IPv4 {
		t.Fatalf("localCacheStatus.Available.IPv4 over kbps threshold expected: false, actual: true")
	} else if localCacheStatus.Available.IPv6 {
		t.Fatalf("localCacheStatus.Available.IPv6 over kbps threshold expected: false, actual: true")
	} else if localCacheStatus.Status != string(tc.CacheStatusReported) {
		t.Fatalf("localCacheStatus.Status expected %v actual %v", "todo", localCacheStatus.Status)
	} else if localCacheStatus.UnavailableStat != "availableBandwidthInKbps" {
		t.Fatalf("localCacheStatus.UnavailableStat expected %v actual %v", "availableBandwidthInKbps", localCacheStatus.UnavailableStat)
	} else if localCacheStatus.Poller != pollerName {
		t.Fatalf("localCacheStatus.Poller expected %v actual %v", pollerName, localCacheStatus.Poller)
	} else if !strings.Contains(localCacheStatus.Why, "availableBandwidthInKbps too low") {
		t.Fatalf("localCacheStatus.Why expected 'availableBandwidthInKbps too low' actual %v", localCacheStatus.Why)
	}

	// test that the health poll didn't override the stat poll threshold markdown and mark available
	// https://github.com/apache/trafficcontrol/issues/3646

	healthResult := result
	healthResult.Astats.System.ProcNetDev = "bond0: 1234567891011121 123456789101    0    5    0     0          0  9876543 12345680160111212 1234567891011    0 1234    0     0       0          0" // 10Gb more than result
	GetVitals(&healthResult, &result, nil)
	healthPollerName := "health"
	healthResults := []cache.Result{healthResult}
	CalcAvailability(healthResults, healthPollerName, nil, mc, toData, localCacheStatusThreadsafe, localStates, events, config.Both)

	localCacheStatuses = localCacheStatusThreadsafe.Get()
	if localCacheStatus, ok := localCacheStatuses[result.ID]; !ok {
		t.Fatalf("expected: localCacheStatus[cacheName], actual: missing")
	} else if localCacheStatus.Available.IPv4 {
		t.Fatalf("localCacheStatus.Available.IPv4 over kbps threshold expected: false, actual: true")
	} else if localCacheStatus.Available.IPv6 {
		t.Fatalf("localCacheStatus.Available.IPv6 over kbps threshold expected: false, actual: true")
	} else if localCacheStatus.Status != string(tc.CacheStatusReported) {
		t.Fatalf("localCacheStatus.Status expected %v actual %v", "tc.CacheStatusReported", localCacheStatus.Status)
	} else if localCacheStatus.UnavailableStat != "availableBandwidthInKbps" {
		t.Fatalf("localCacheStatus.UnavailableStat expected %v actual %v", "availableBandwidthInKbps", localCacheStatus.UnavailableStat)
	} else if localCacheStatus.Poller != healthPollerName {
		t.Fatalf("localCacheStatus.Poller expected %v actual %v", healthPollerName, localCacheStatus.Poller)
	} else if !strings.Contains(localCacheStatus.Why, "availableBandwidthInKbps too low") {
		t.Fatalf("localCacheStatus.Why expected 'availableBandwidthInKbps too low' actual %v", localCacheStatus.Why)
	}
}
