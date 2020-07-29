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

	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/config"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"

	"github.com/apache/trafficcontrol/lib/go-tc"
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
		InterfaceVitals: map[string]cache.Vitals{},
		PollID:          42,
		PollFinished:    make(chan uint64, 1),
		PrecomputedData: cache.PrecomputedData{},
		Available:       true,
		UsingIPv4:       false,
	}
	GetVitals(&result, nil, nil)

	totalBytesOut := result.Statistics.Interfaces["bond0"].BytesOut + result.Statistics.Interfaces["eth0"].BytesOut
	if totalBytesOut != result.Vitals.BytesOut {
		t.Errorf("Incorrect calculated BytesOut; expected: %d, got: %d", totalBytesOut, result.Vitals.BytesOut)
	}
	prevIV := map[string]cache.Vitals{}
	prevIV["eth0"] = cache.Vitals{BytesOut: result.Vitals.BytesOut - 1250000000}  // 10 gigabits
	prevIV["bond0"] = cache.Vitals{BytesOut: result.Vitals.BytesOut - 1250000000} // 10 gigabits

	prevResult := cache.Result{
		Time:            result.Time.Add(time.Second * -1),
		Vitals:          cache.Vitals{BytesOut: result.Vitals.BytesOut - 1250000000}, // 10 gigabits
		InterfaceVitals: prevIV,
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
	localCacheStatus, ok := localCacheStatuses[result.ID]
	if !ok {
		t.Fatalf("expected: localCacheStatus[cacheName], actual: missing")
	}

	if !strings.Contains(localCacheStatus.Why, "availableBandwidthInKbps too low") {
		t.Errorf("localCacheStatus.Why expected 'availableBandwidthInKbps too low' actual %s", localCacheStatus.Why)
	} else if !strings.HasPrefix(localCacheStatus.UnavailableStat, "availableBandwidthInKbps") { // only check prefix because we don't care about the specific interfaces right now
		t.Errorf("localCacheStatus.UnavailableStat expected it to start with: 'availableBandwidthInKbps', actual: '%s'", localCacheStatus.UnavailableStat)
	}
	if localCacheStatus.Available.IPv4 {
		t.Errorf("localCacheStatus.Available.IPv4 over kbps threshold expected: false, actual: true")
	}
	if localCacheStatus.Available.IPv6 {
		t.Error("localCacheStatus.Available.IPv6 over kbps threshold expected: false, actual: true")
	}
	if localCacheStatus.Status != string(tc.CacheStatusReported) {
		t.Errorf("localCacheStatus.Status expected: 'todo', actual: '%s'", localCacheStatus.Status)
	}
	if localCacheStatus.Poller != pollerName {
		t.Errorf("localCacheStatus.Poller expected '%s' actual '%s'", pollerName, localCacheStatus.Poller)
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
	if _, ok := localCacheStatuses[result.ID]; !ok {
		t.Fatalf("expected: localCacheStatus[cacheName], actual: missing")
	}

	localCacheStatus = localCacheStatuses[result.ID]
	if !strings.Contains(localCacheStatus.Why, "availableBandwidthInKbps too low") {
		t.Errorf("localCacheStatus.Why expected: 'availableBandwidthInKbps too low' actual: '%s'", localCacheStatus.Why)
	} else if !strings.HasPrefix(localCacheStatus.UnavailableStat, "availableBandwidthInKbps") { // only check prefix because we don't care about the specific interfaces right now
		t.Errorf("localCacheStatus.UnavailableStat expected it to start with: 'availableBandwidthInKbps', actual: '%s'", localCacheStatus.UnavailableStat)
	}

	if localCacheStatus.Available.IPv4 {
		t.Fatal("localCacheStatus.Available.IPv4 over kbps threshold expected: false, actual: true")
	} else if localCacheStatus.Available.IPv6 {
		t.Fatal("localCacheStatus.Available.IPv6 over kbps threshold expected: false, actual: true")
	} else if localCacheStatus.Status != string(tc.CacheStatusReported) {
		t.Fatalf("localCacheStatus.Status expected: 'todo' actual: '%s'", localCacheStatus.Status)
	} else if localCacheStatus.Poller != healthPollerName {
		t.Fatalf("localCacheStatus.Poller expected: '%s' actual '%s'", healthPollerName, localCacheStatus.Poller)
	}
}

func TestEvalInterface(t *testing.T) {
	result := cache.ResultInfo{
		Available:       true,
		Error:           nil,
		ID:              "test",
		PollID:          1,
		RequestTime:     time.Second,
		Statistics:      cache.Statistics{},
		Time:            time.Now(),
		UsingIPv4:       true,
		Vitals:          cache.Vitals{},
		InterfaceVitals: map[string]cache.Vitals{},
	}

	var infMaxKbps uint64 = 200
	mc := tc.TrafficMonitorConfigMap{
		Profile: map[string]tc.TMProfile{
			"testProfile": {},
		},
		TrafficServer: map[string]tc.TrafficServer{
			"test": {
				Profile: "testProfile",
				Interfaces: []tc.ServerInterfaceInfo{
					{
						Monitor:      true,
						MaxBandwidth: &infMaxKbps,
						Name:         "testInterface",
						IPAddresses: []tc.ServerIPAddress{
							{
								Address:        "::1",
								ServiceAddress: true,
							},
						},
					},
				},
			},
		},
	}

	available, why := EvalInterface(result.InterfaceVitals, mc.TrafficServer["test"].Interfaces[0])
	if available {
		t.Error("Expected unpolled interface to be unavailable, but it wasn't")
	}

	if why != "not found in polled data" {
		t.Errorf("Incorrect reason for unpolled interface's availability; expected: 'not found in polled data', got: '%s'", why)
	}

	result.InterfaceVitals["testInterface"] = cache.Vitals{
		KbpsOut: 199,
	}
	available, why = EvalInterface(result.InterfaceVitals, mc.TrafficServer["test"].Interfaces[0])
	if !available {
		t.Error("Expected polled interface within threshold to be available, but it wasn't")
	}

	if why != "" {
		t.Errorf("Expected no reason why polled interface within threshold is unavailable, got: '%s'", why)
	}

	result.InterfaceVitals["testInterface"] = cache.Vitals{
		KbpsOut: 201,
	}
	available, why = EvalInterface(result.InterfaceVitals, mc.TrafficServer["test"].Interfaces[0])
	if available {
		t.Error("Expected interface exceeding threshold to be unavailable, but it wasn't")
	}

	if why != "maximum bandwidth exceeded" {
		t.Errorf("Incorrect reason for interface exceeding threshold to be unavailable; expected: 'maximum bandwidth exceeded', got: '%s'", why)
	}
}
