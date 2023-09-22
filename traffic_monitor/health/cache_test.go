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

	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// TestNoMonitoredInterfacesGetVitals assures that GetVitals
// does not fail even if no interfaces are marked to be monitored
func TestNoMonitoredInterfacesGetVitals(t *testing.T) {
	serverID := "no-monitored"
	fakeRequestTime := time.Now()
	zeroValueVitals := cache.Vitals{}

	// Interfaces to monitor are marked true (none)
	tmcm := tc.TrafficMonitorConfigMap{
		TrafficServer: map[string]tc.TrafficServer{
			serverID: {
				Interfaces: []tc.ServerInterfaceInfo{
					{
						Name:    "bond0",
						Monitor: false,
					},
					{
						Name:    "bond1",
						Monitor: false,
					},
					{
						Name:    "lo",
						Monitor: false,
					},
				},
			},
		},
	}

	// multiple interfaces, plus extra
	firstResult := cache.Result{
		ID:            serverID,
		Error:         nil,
		Miscellaneous: map[string]interface{}{},
		Statistics: cache.Statistics{
			Interfaces: map[string]cache.Interface{
				"bond0": {
					Speed:    100000,
					BytesIn:  570791700709,
					BytesOut: 4212211168526,
				},
				"bond1": {
					Speed:    100000,
					BytesIn:  1989352297218,
					BytesOut: 10630690813,
				},
				"lo": {
					Speed:    0,
					BytesIn:  181882394,
					BytesOut: 181882394,
				},
				"em5": {
					Speed:    0,
					BytesIn:  0,
					BytesOut: 0,
				},
			},
		},
		Time:            fakeRequestTime,
		RequestTime:     time.Second,
		Vitals:          cache.Vitals{},
		InterfaceVitals: nil,
		PollID:          42,
		PollFinished:    make(chan uint64, 1),
		PrecomputedData: cache.PrecomputedData{},
		Available:       true,
		UsingIPv4:       false,
	}
	GetVitals(&firstResult, nil, &tmcm)

	// No interfaces were selected to be monitored so none
	// should have been added later
	if len(firstResult.InterfaceVitals) > 0 {
		t.Errorf("InterfaceVitals map should be empty. expected: %v actual: %v:", 0, len(firstResult.InterfaceVitals))
	}

	// No interfaces were selected to be monitored so no vitals
	// should have been calculated
	if firstResult.Vitals != zeroValueVitals {
		t.Errorf("Vitals should have zero values. expected: %v actual: %v:", zeroValueVitals, firstResult.Vitals)
	}

	secondResult := firstResult
	secondResult.Time = fakeRequestTime.Add(5 * time.Second)

	GetVitals(&secondResult, &firstResult, &tmcm)

	// No interfaces were selected to be monitored so none
	// should have been added later
	if len(secondResult.InterfaceVitals) > 0 {
		t.Errorf("InterfaceVitals map should be empty. expected: %v actual: %v:", 0, len(secondResult.InterfaceVitals))
	}

	// No interfaces were selected to be monitored so no vitals
	// should have been calculated
	if secondResult.Vitals != zeroValueVitals {
		t.Errorf("Vitals should have zero values. expected: %v actual: %v:", zeroValueVitals, secondResult.Vitals)
	}

	// The previous results should not have been impacted
	if firstResult.Vitals != zeroValueVitals {
		t.Errorf("Vitals should have zero values. expected: %v actual: %v:", zeroValueVitals, firstResult.Vitals)
	}
}

// TestDualHomingMonitoredInterfacesGetVitals ensures cache servers
// with multiple interfaces correctly calculate bandwidth based on
// whether the interfaces are marked as "Monitor this interface"
func TestDualHomingMonitoredInterfacesGetVitals(t *testing.T) {

	serverID := "dual-homed"
	fakeRequestTime := time.Now()

	// Interfaces to monitor are marked true
	tmcm := tc.TrafficMonitorConfigMap{
		TrafficServer: map[string]tc.TrafficServer{
			serverID: {
				Interfaces: []tc.ServerInterfaceInfo{
					{
						Name:    "bond0",
						Monitor: true,
					},
					{
						Name:    "bond1",
						Monitor: true,
					},
					{
						Name:    "lo",
						Monitor: false,
					},
				},
			},
		},
	}

	// multiple interfaces, plus extras
	firstResult := cache.Result{
		ID:            serverID,
		Error:         nil,
		Miscellaneous: map[string]interface{}{},
		Statistics: cache.Statistics{
			Interfaces: map[string]cache.Interface{
				"bond0": {
					Speed:    100000,
					BytesIn:  570791700709,
					BytesOut: 4212211168526,
				},
				"bond1": {
					Speed:    100000,
					BytesIn:  1989352297218,
					BytesOut: 10630690813,
				},
				"p1p1": {
					Speed:    100000,
					BytesIn:  570793589545,
					BytesOut: 4212220919951,
				},
				"p3p1": {
					Speed:    100000,
					BytesIn:  1989354450479,
					BytesOut: 10630690813,
				},
				"lo": {
					Speed:    0,
					BytesIn:  181882394,
					BytesOut: 181882394,
				},
				"em5": {
					Speed:    0,
					BytesIn:  0,
					BytesOut: 0,
				},
				"em6": {
					Speed:    0,
					BytesIn:  0,
					BytesOut: 0,
				},
			},
		},
		Time:            fakeRequestTime,
		RequestTime:     time.Second,
		Vitals:          cache.Vitals{},
		InterfaceVitals: nil,
		PollID:          42,
		PollFinished:    make(chan uint64, 1),
		PrecomputedData: cache.PrecomputedData{},
		Available:       true,
		UsingIPv4:       false,
	}
	GetVitals(&firstResult, nil, &tmcm)

	// Two interfaces were selected to be monitored so they
	// should have been added later
	if len(firstResult.InterfaceVitals) != 2 {
		t.Errorf("InterfaceVitals map should not be empty. expected: %v actual: %v:", 2, len(firstResult.InterfaceVitals))
	}

	expectedFirstVitals := cache.Vitals{
		LoadAvg:    0,
		BytesIn:    2560143997927,
		BytesOut:   4222841859339,
		KbpsOut:    0,
		MaxKbpsOut: 200000000,
	}
	// Only two interfaces were selected to be monitored so vitals
	// should have been calculated based on those two (bond0 and bond1)
	if firstResult.Vitals != expectedFirstVitals {
		t.Errorf("Vitals do not match expected output. expected: %v actual: %v:", expectedFirstVitals, firstResult.Vitals)
	}

	secondResult := firstResult
	secondResult.Statistics.Interfaces = map[string]cache.Interface{
		"bond0": {
			Speed:    100000,
			BytesIn:  572608907987,
			BytesOut: 4227149141326,
		},
		"bond1": {
			Speed:    100000,
			BytesIn:  1996376171468,
			BytesOut: 10630696953,
		},
		"p1p1": {
			Speed:    100000,
			BytesIn:  572609282353,
			BytesOut: 4227157881921,
		},
		"p3p1": {
			Speed:    100000,
			BytesIn:  1996378204692,
			BytesOut: 10630696953,
		},
		"lo": {
			Speed:    0,
			BytesIn:  181882394,
			BytesOut: 181882394,
		},
		"em5": {
			Speed:    0,
			BytesIn:  0,
			BytesOut: 0,
		},
		"em6": {
			Speed:    0,
			BytesIn:  0,
			BytesOut: 0,
		},
	}
	secondResult.Time = fakeRequestTime.Add(5 * time.Second)
	secondResult.Vitals = cache.Vitals{}

	GetVitals(&secondResult, &firstResult, &tmcm)

	// Two interfaces were selected to be monitored so they
	// should have been added later
	if len(secondResult.InterfaceVitals) != 2 {
		t.Errorf("InterfaceVitals map should not be empty. expected: %v actual: %v:", 2, len(secondResult.InterfaceVitals))
	}

	expectedSecondVitals := cache.Vitals{
		LoadAvg:    0,
		BytesIn:    2568985079455,
		BytesOut:   4237779838279,
		KbpsOut:    23900766,
		MaxKbpsOut: 200000000,
	}

	// Only two interfaces were selected to be monitored so vitals
	// should have been calculated based on those two (bond0 and bond1)
	if secondResult.Vitals != expectedSecondVitals {
		t.Errorf("Vitals do not match expected output. expected: %v actual: %v:", expectedSecondVitals, secondResult.Vitals)
	}

	// Previous result values should have been altered
	if firstResult.Vitals != expectedFirstVitals {
		t.Errorf("Vitals do not match expected output. expected: %v actual: %v:", expectedFirstVitals, firstResult.Vitals)
	}
}

func TestCalcAvailabilityThresholds(t *testing.T) {

	resultID := "myCacheName"

	mc := tc.TrafficMonitorConfigMap{
		TrafficServer: map[string]tc.TrafficServer{
			string(resultID): {
				ServerStatus: string(tc.CacheStatusReported),
				Profile:      "myProfileName",
				Interfaces: []tc.ServerInterfaceInfo{
					{
						Name:    "bond0",
						Monitor: true,
					},
					{
						Name:    "eth0",
						Monitor: true,
					},
					{
						Name:    "lo",
						Monitor: false,
					},
				},
			},
		},
		Profile: map[string]tc.TMProfile{},
	}

	result := cache.Result{
		ID:            resultID,
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
				"bond0": {
					Speed:    20000,
					BytesIn:  1234567891011121,
					BytesOut: 12345678910111213,
				},
				"eth0": {
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
	GetVitals(&result, nil, &mc)

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
	GetVitals(&result, &prevResult, &mc)

	mc.Profile[mc.TrafficServer[string(result.ID)].Profile] = tc.TMProfile{
		Name: mc.TrafficServer[string(result.ID)].Profile,
		Parameters: tc.TMParameters{
			Thresholds: map[string]tc.HealthThreshold{
				"availableBandwidthInKbps": {
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
	toData.ServerCachegroups[tc.CacheName(result.ID)] = "myCG"

	localCacheStatusThreadsafe := threadsafe.NewCacheAvailableStatus()
	localStates := peer.NewCRStatesThreadsafe()
	localStates.SetDeliveryService("myDS", tc.CRStatesDeliveryService{})
	events := NewThreadsafeEvents(200)

	// test that a normal stat poll over the kbps threshold marks down

	pollerName := "stat"
	results := []cache.Result{result}

	// Ensure that if the interfaces haven't been reported yet that CalcAvailability doesn't panic
	original := results[0].Statistics.Interfaces
	statResultHistory := (*threadsafe.ResultStatHistory)(nil)
	results[0].Statistics.Interfaces = make(map[string]cache.Interface)
	CalcAvailability(results, pollerName, statResultHistory, mc, toData, localCacheStatusThreadsafe, localStates, events, config.Both)
	results[0].Statistics.Interfaces = original

	CalcAvailability(results, pollerName, statResultHistory, mc, toData, localCacheStatusThreadsafe, localStates, events, config.Both)

	// ensure that the DisabledLocations is an empty, non-nil slice
	for _, ds := range localStates.GetDeliveryServices() {
		if ds.DisabledLocations == nil {
			t.Error("expected: non-nil DisabledLocations, actual: nil")
		}
		if len(ds.DisabledLocations) > 0 {
			t.Errorf("expected: empty DisabledLocations, actual: %d", len(ds.DisabledLocations))
		}
	}

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
