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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/traffic_monitor/config"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

func TestCalcAvailabilityThresholds(t *testing.T) {
	result := cache.Result{
		ID:            "myCacheName",
		Error:         nil,
		Miscellaneous: map[string]interface{}{},
		Statistics: cache.Statistics{
			Loadavg: cache.Loadavg{
				One:              5.43,
				Five:             4.32,
				Fifteen:          3.21,
				CurrentProcesses: 3,
				TotalProcesses:   1234,
				LatestPID:        32109,
			},
			Interfaces: map[string]cache.Interface{
				"bond0": cache.Interface{
					Speed:    20000,
					BytesIn:  1234567891011121,
					BytesOut: 12345678910111213,
				},
				"eth0": cache.Interface{
					Speed:    30000,
					BytesIn:  1234567891011121,
					BytesOut: 12345678910111213,
				},
			},
			NotAvailable: false,
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
	toData.ServerTypes[tc.CacheName(result.ID)] = tc.CacheTypeEdge
	toData.DeliveryServiceServers["myDS"] = []tc.CacheName{tc.CacheName(result.ID)}
	toData.ServerCachegroups[tc.CacheName(result.ID)] = "myCG"

	localCacheStatusThreadsafe := threadsafe.NewCacheAvailableStatus()
	localStates := peer.NewCRStatesThreadsafe()
	events := NewThreadsafeEvents(200)

	// test that a normal stat poll over the kbps threshold marks down

	pollerName := "stat"
	results := []cache.Result{result}

	// Ensure that if the interfaces haven't been reported yet that CalcAvailability doesn't panic
	original := results[0].Statistics.Interfaces
	results[0].Statistics.Interfaces = make(map[string]cache.Interface)
	CalcAvailability(results, pollerName, statResultHistory, mc, toData, localCacheStatusThreadsafe, localStates, events, config.Both)
	results[0].Statistics.Interfaces = original

	CalcAvailability(results, pollerName, statResultHistory, mc, toData, localCacheStatusThreadsafe, localStates, events, config.Both)

	localCacheStatuses := localCacheStatusThreadsafe.Get()
	if _, ok := localCacheStatuses[tc.CacheName(result.ID)]; !ok {
		t.Fatalf("expected: localCacheStatus[cacheName], actual: missing")
	}
	for interfaceName, localCacheStatus := range localCacheStatuses[tc.CacheName(result.ID)] {
		if interfaceName == tc.CacheInterfacesAggregate {
			continue
		}
		if localCacheStatus.Available.IPv4 {
			t.Fatalf("localCacheStatus.Available.IPv4 (%s) over kbps threshold expected: false, actual: true", interfaceName)
		} else if localCacheStatus.Available.IPv6 {
			t.Fatalf("localCacheStatus.Available.IPv6 (%s) over kbps threshold expected: false, actual: true", interfaceName)
		} else if localCacheStatus.Status != string(tc.CacheStatusReported) {
			t.Fatalf("localCacheStatus.Status (%s) expected %v actual %v", interfaceName, "todo", localCacheStatus.Status)
		} else if localCacheStatus.UnavailableStat != "availableBandwidthInKbps" {
			t.Fatalf("localCacheStatus.UnavailableStat (%s) expected %v actual %v", interfaceName, "availableBandwidthInKbps", localCacheStatus.UnavailableStat)
		} else if localCacheStatus.Poller != pollerName {
			t.Fatalf("localCacheStatus.Poller (%s) expected %v actual %v", interfaceName, pollerName, localCacheStatus.Poller)
		} else if !strings.Contains(localCacheStatus.Why, "availableBandwidthInKbps too low") {
			t.Fatalf("localCacheStatus.Why (%s) expected 'availableBandwidthInKbps too low' actual %v", interfaceName, localCacheStatus.Why)
		}
	}

	// test that the health poll didn't override the stat poll threshold markdown and mark available
	// https://github.com/apache/trafficcontrol/issues/3646

	healthResult := result
	healthResultInf := result.Statistics.Interfaces["bond0"]
	healthResultInf.BytesOut = 12345680160111212
	healthResult.Statistics.Interfaces["bond0"] = healthResultInf

	GetVitals(&healthResult, &result, nil)
	healthPollerName := "health"
	healthResults := []cache.Result{healthResult}
	CalcAvailability(healthResults, healthPollerName, nil, mc, toData, localCacheStatusThreadsafe, localStates, events, config.Both)

	localCacheStatuses = localCacheStatusThreadsafe.Get()
	localCacheStatus := localCacheStatuses[tc.CacheName(result.ID)]["bond0"]
	if localCacheStatus.Available.IPv4 {
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
