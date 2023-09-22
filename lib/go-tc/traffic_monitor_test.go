package tc

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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func makeFakeStats(text string) map[string][]ResultStatVal {
	return map[string][]ResultStatVal{
		text + "stat1": {
			{
				Span: 50,
				Time: time.Now(),
				Val:  50,
			},
		},
		text + "stat2": {
			{
				Span: 50,
				Time: time.Now(),
				Val:  50,
			},
			{
				Span: 50,
				Time: time.Now(),
				Val:  50,
			},
		},
		text + "stat3": {
			{
				Span: 50,
				Time: time.Now(),
				Val:  50,
			},
			{
				Span: 50,
				Time: time.Now(),
				Val:  50,
			},
			{
				Span: 50,
				Time: time.Now(),
				Val:  50,
			},
		},
	}
}

func makeFakeInterfaces() map[string]map[string][]ResultStatVal {
	return map[string]map[string][]ResultStatVal{
		"interface1": makeFakeStats("interf"),
		"interface2": makeFakeStats("interf"),
		"interface3": makeFakeStats("interf"),
	}
}
func TestLegacyStatsConversion(t *testing.T) {
	stats := Stats{
		CommonAPIData: CommonAPIData{},
		Caches:        make(map[string]ServerStats),
	}
	config := TrafficMonitorConfigMap{
		TrafficServer: make(map[string]TrafficServer),
	}
	for _, cacheName := range []string{"cache1", "cache2", "cache3"} {
		stats.Caches[cacheName] = ServerStats{
			Interfaces: makeFakeInterfaces(),
			Stats:      makeFakeStats(""),
		}
		interfaces := []ServerInterfaceInfo{}
		for name, _ := range stats.Caches[cacheName].Interfaces {
			interfaces = append(interfaces, ServerInterfaceInfo{
				IPAddresses: []ServerIPAddress{
					{
						Address:        "192.168.0.8",
						Gateway:        util.StrPtr("192.168.0.1"),
						ServiceAddress: true,
					},
				},
				MaxBandwidth: util.Uint64Ptr(1500),
				Monitor:      false,
				MTU:          util.UInt64Ptr(1500),
				Name:         name,
			})
		}
		interfaces[0].Monitor = true
		config.TrafficServer[cacheName] = TrafficServer{
			Interfaces: interfaces,
		}
	}

	issues, legacyStats := stats.ToLegacy(config)

	if len(issues) != 0 {
		t.Error("expect no issues")
	}

	if legacyStats.CommonAPIData != stats.CommonAPIData {
		t.Error("expected CommonAPIData to be the same")
	}

	if len(legacyStats.Caches) != len(stats.Caches) {
		t.Errorf("expected %v caches, got %v", len(stats.Caches), len(legacyStats.Caches))
	}

	for cacheName, legacyCache := range legacyStats.Caches {
		cache, ok := stats.Caches[string(cacheName)]
		if !ok {
			t.Errorf("new interface %v found in upgraded stats, but not in legacy stats", cacheName)
		}
		interf := cache.Interfaces["interface1"]
		if len(interf)+len(cache.Stats) != len(legacyCache) {
			t.Errorf("expected %v stats, got %v", len(interf)+len(cache.Stats), len(legacyCache))
		}
	}
}

func TestLegacyStatsNilConversion(t *testing.T) {
	stats := Stats{
		CommonAPIData: CommonAPIData{},
		Caches:        nil,
	}
	config := TrafficMonitorConfigMap{
		TrafficServer: nil,
	}
	issues, legacyStats := stats.ToLegacy(config)

	if legacyStats.CommonAPIData != stats.CommonAPIData {
		t.Error("expected CommonAPIData to be the same")
	}

	if len(issues) != 0 {
		t.Error("expect no issues")
	}
}

func TestLegacyMonitorConfigValid(t *testing.T) {
	mc := (*LegacyTrafficMonitorConfigMap)(nil)
	if LegacyMonitorConfigValid(mc) == nil {
		t.Errorf("MonitorCopnfigValid(nil) expected: error, actual: nil")
	}
	mc = &LegacyTrafficMonitorConfigMap{}
	if LegacyMonitorConfigValid(mc) == nil {
		t.Errorf("MonitorConfigValid({}) expected: error, actual: nil")
	}

	validMC := &LegacyTrafficMonitorConfigMap{
		TrafficServer:   map[string]LegacyTrafficServer{"a": {}},
		CacheGroup:      map[string]TMCacheGroup{"a": {}},
		TrafficMonitor:  map[string]TrafficMonitor{"a": {}},
		DeliveryService: map[string]TMDeliveryService{"a": {}},
		Profile:         map[string]TMProfile{"a": {}},
		Config: map[string]interface{}{
			"peers.polling.interval":  42.0,
			"health.polling.interval": 24.0,
		},
	}
	if err := LegacyMonitorConfigValid(validMC); err != nil {
		t.Errorf("MonitorConfigValid(%++v) expected: nil, actual: %+v", validMC, err)
	}
}

func ExampleHealthThreshold_String() {
	ht := HealthThreshold{Comparator: ">=", Val: 500}
	fmt.Println(ht)
	// Output: >=500.000000
}

func ExampleTMParameters_UnmarshalJSON() {
	const data = `{
		"health.connection.timeout": 5,
		"health.polling.url": "https://example.com/",
		"health.polling.format": "stats_over_http",
		"history.count": 1,
		"health.threshold.bandwidth": ">50",
		"health.threshold.foo": "<=500"
	}`

	var params TMParameters
	if err := json.Unmarshal([]byte(data), &params); err != nil {
		fmt.Printf("Failed to unmarshal: %v\n", err)
		return
	}
	fmt.Printf("timeout: %d\n", params.HealthConnectionTimeout)
	fmt.Printf("url: %s\n", params.HealthPollingURL)
	fmt.Printf("format: %s\n", params.HealthPollingFormat)
	fmt.Printf("history: %d\n", params.HistoryCount)
	fmt.Printf("# of Thresholds: %d - foo: %s, bandwidth: %s\n", len(params.Thresholds), params.Thresholds["foo"], params.Thresholds["bandwidth"])

	// Output: timeout: 5
	// url: https://example.com/
	// format: stats_over_http
	// history: 1
	// # of Thresholds: 2 - foo: <=500.000000, bandwidth: >50.000000
}

func ExampleTrafficMonitorConfigMap_Valid() {
	mc := &TrafficMonitorConfigMap{
		CacheGroup: map[string]TMCacheGroup{"a": {}},
		Config: map[string]interface{}{
			"peers.polling.interval":  0.0,
			"health.polling.interval": 0.0,
		},
		DeliveryService: map[string]TMDeliveryService{"a": {}},
		Profile:         map[string]TMProfile{"a": {}},
		TrafficMonitor:  map[string]TrafficMonitor{"a": {}},
		TrafficServer:   map[string]TrafficServer{"a": {}},
	}

	fmt.Printf("Validity error: %v", mc.Valid())

	// Output: Validity error: <nil>
}

func ExampleLegacyTrafficMonitorConfigMap_Upgrade() {
	lcm := LegacyTrafficMonitorConfigMap{
		CacheGroup: map[string]TMCacheGroup{
			"test": {
				Name: "test",
				Coordinates: MonitoringCoordinates{
					Latitude:  0,
					Longitude: 0,
				},
			},
		},
		Config: map[string]interface{}{
			"foo": "bar",
		},
		DeliveryService: map[string]TMDeliveryService{
			"test": {
				XMLID:              "test",
				TotalTPSThreshold:  -1,
				ServerStatus:       "testStatus",
				TotalKbpsThreshold: -1,
			},
		},
		Profile: map[string]TMProfile{
			"test": {
				Parameters: TMParameters{
					HealthConnectionTimeout: -1,
					HealthPollingURL:        "testURL",
					HealthPollingFormat:     "astats",
					HealthPollingType:       "http",
					HistoryCount:            -1,
					MinFreeKbps:             -1,
					Thresholds: map[string]HealthThreshold{
						"availableBandwidthInKbps": {
							Comparator: "<",
							Val:        -1,
						},
					},
				},
				Name: "test",
				Type: "testType",
			},
		},
		TrafficMonitor: map[string]TrafficMonitor{
			"test": {
				Port:         -1,
				IP6:          "::1",
				IP:           "0.0.0.0",
				HostName:     "test",
				FQDN:         "test.quest",
				Profile:      "test",
				Location:     "test",
				ServerStatus: "testStatus",
			},
		},
		TrafficServer: map[string]LegacyTrafficServer{
			"test": {
				CacheGroup:       "test",
				DeliveryServices: []TSDeliveryService{},
				FQDN:             "test.quest",
				HashID:           "test",
				HostName:         "test",
				HTTPSPort:        -1,
				InterfaceName:    "testInterface",
				IP:               "0.0.0.1",
				IP6:              "::2",
				Port:             -1,
				Profile:          "test",
				ServerStatus:     "testStatus",
				Type:             "testType",
			},
		},
	}

	cm := lcm.Upgrade()
	fmt.Println("# of Cachegroups:", len(cm.CacheGroup))
	fmt.Println("Cachegroup Name:", cm.CacheGroup["test"].Name)
	fmt.Printf("Cachegroup Coordinates: (%v,%v)\n", cm.CacheGroup["test"].Coordinates.Latitude, cm.CacheGroup["test"].Coordinates.Longitude)
	fmt.Println("# of Config parameters:", len(cm.Config))
	fmt.Println(`Config["foo"]:`, cm.Config["foo"])
	fmt.Println("# of DeliveryServices:", len(cm.DeliveryService))
	fmt.Println("DeliveryService XMLID:", cm.DeliveryService["test"].XMLID)
	fmt.Println("DeliveryService TotalTPSThreshold:", cm.DeliveryService["test"].TotalTPSThreshold)
	fmt.Println("DeliveryService ServerStatus:", cm.DeliveryService["test"].ServerStatus)
	fmt.Println("DeliveryService TotalKbpsThreshold:", cm.DeliveryService["test"].TotalKbpsThreshold)
	fmt.Println("# of Profiles:", len(cm.Profile))
	fmt.Println("Profile Name:", cm.Profile["test"].Name)
	fmt.Println("Profile Type:", cm.Profile["test"].Type)
	fmt.Println("Profile HealthConnectionTimeout:", cm.Profile["test"].Parameters.HealthConnectionTimeout)
	fmt.Println("Profile HealthPollingURL:", cm.Profile["test"].Parameters.HealthPollingURL)
	fmt.Println("Profile HealthPollingFormat:", cm.Profile["test"].Parameters.HealthPollingFormat)
	fmt.Println("Profile HealthPollingType:", cm.Profile["test"].Parameters.HealthPollingType)
	fmt.Println("Profile HistoryCount:", cm.Profile["test"].Parameters.HistoryCount)
	fmt.Println("Profile MinFreeKbps:", cm.Profile["test"].Parameters.MinFreeKbps)
	fmt.Println("# of Profile Thresholds:", len(cm.Profile["test"].Parameters.Thresholds))
	fmt.Println("Profile availableBandwidthInKbps Threshold:", cm.Profile["test"].Parameters.Thresholds["availableBandwidthInKbps"])
	fmt.Println("# of TrafficMonitors:", len(cm.TrafficMonitor))
	fmt.Println("TrafficMonitor Port:", cm.TrafficMonitor["test"].Port)
	fmt.Println("TrafficMonitor IP6:", cm.TrafficMonitor["test"].IP6)
	fmt.Println("TrafficMonitor IP:", cm.TrafficMonitor["test"].IP)
	fmt.Println("TrafficMonitor HostName:", cm.TrafficMonitor["test"].HostName)
	fmt.Println("TrafficMonitor FQDN:", cm.TrafficMonitor["test"].FQDN)
	fmt.Println("TrafficMonitor Profile:", cm.TrafficMonitor["test"].Profile)
	fmt.Println("TrafficMonitor Location:", cm.TrafficMonitor["test"].Location)
	fmt.Println("TrafficMonitor ServerStatus:", cm.TrafficMonitor["test"].ServerStatus)
	fmt.Println("# of TrafficServers:", len(cm.TrafficServer))
	fmt.Println("TrafficServer CacheGroup:", cm.TrafficServer["test"].CacheGroup)
	fmt.Println("TrafficServer # of DeliveryServices:", len(cm.TrafficServer["test"].DeliveryServices))
	fmt.Println("TrafficServer FQDN:", cm.TrafficServer["test"].FQDN)
	fmt.Println("TrafficServer HashID:", cm.TrafficServer["test"].HashID)
	fmt.Println("TrafficServer HostName:", cm.TrafficServer["test"].HostName)
	fmt.Println("TrafficServer HTTPSPort:", cm.TrafficServer["test"].HTTPSPort)
	fmt.Println("TrafficServer # of Interfaces:", len(cm.TrafficServer["test"].Interfaces))
	fmt.Println("TrafficServer Interface Name:", cm.TrafficServer["test"].Interfaces[0].Name)
	fmt.Println("TrafficServer # of Interface IP Addresses:", len(cm.TrafficServer["test"].Interfaces[0].IPAddresses))
	fmt.Println("TrafficServer first IP Address:", cm.TrafficServer["test"].Interfaces[0].IPAddresses[0].Address)
	fmt.Println("TrafficServer second IP Address:", cm.TrafficServer["test"].Interfaces[0].IPAddresses[1].Address)
	fmt.Println("TrafficServer Port:", cm.TrafficServer["test"].Port)
	fmt.Println("TrafficServer Profile:", cm.TrafficServer["test"].Profile)
	fmt.Println("TrafficServer ServerStatus:", cm.TrafficServer["test"].ServerStatus)
	fmt.Println("TrafficServer Type:", cm.TrafficServer["test"].Type)

	// Output: # of Cachegroups: 1
	// Cachegroup Name: test
	// Cachegroup Coordinates: (0,0)
	// # of Config parameters: 1
	// Config["foo"]: bar
	// # of DeliveryServices: 1
	// DeliveryService XMLID: test
	// DeliveryService TotalTPSThreshold: -1
	// DeliveryService ServerStatus: testStatus
	// DeliveryService TotalKbpsThreshold: -1
	// # of Profiles: 1
	// Profile Name: test
	// Profile Type: testType
	// Profile HealthConnectionTimeout: -1
	// Profile HealthPollingURL: testURL
	// Profile HealthPollingFormat: astats
	// Profile HealthPollingType: http
	// Profile HistoryCount: -1
	// Profile MinFreeKbps: -1
	// # of Profile Thresholds: 1
	// Profile availableBandwidthInKbps Threshold: <-1.000000
	// # of TrafficMonitors: 1
	// TrafficMonitor Port: -1
	// TrafficMonitor IP6: ::1
	// TrafficMonitor IP: 0.0.0.0
	// TrafficMonitor HostName: test
	// TrafficMonitor FQDN: test.quest
	// TrafficMonitor Profile: test
	// TrafficMonitor Location: test
	// TrafficMonitor ServerStatus: testStatus
	// # of TrafficServers: 1
	// TrafficServer CacheGroup: test
	// TrafficServer # of DeliveryServices: 0
	// TrafficServer FQDN: test.quest
	// TrafficServer HashID: test
	// TrafficServer HostName: test
	// TrafficServer HTTPSPort: -1
	// TrafficServer # of Interfaces: 1
	// TrafficServer Interface Name: testInterface
	// TrafficServer # of Interface IP Addresses: 2
	// TrafficServer first IP Address: 0.0.0.1
	// TrafficServer second IP Address: ::2
	// TrafficServer Port: -1
	// TrafficServer Profile: test
	// TrafficServer ServerStatus: testStatus
	// TrafficServer Type: testType
}

func TestTrafficMonitorConfigMap_Valid(t *testing.T) {
	var mc *TrafficMonitorConfigMap = nil
	err := mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of nil config map")
	} else {
		t.Logf("Got expected error: checking validity of nil config map: %v", err)
	}
	mc = &TrafficMonitorConfigMap{
		CacheGroup: nil,
		Config: map[string]interface{}{
			"peers.polling.interval":  42.0,
			"health.polling.interval": 24.0,
		},
		DeliveryService: map[string]TMDeliveryService{"a": {}},
		Profile:         map[string]TMProfile{"a": {}},
		TrafficMonitor:  map[string]TrafficMonitor{"a": {}},
		TrafficServer:   map[string]TrafficServer{"a": {}},
	}

	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with nil CacheGroup")
	} else {
		t.Logf("Got expected error: checking validity of config map with nil CacheGroup: %v", err)
	}

	mc.CacheGroup = map[string]TMCacheGroup{}
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with no CacheGroups")
	} else {
		t.Logf("Got expected error: checking validity of config map with no CacheGroups: %v", err)
	}

	mc.CacheGroup["a"] = TMCacheGroup{}
	mc.Config = nil
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with nil Config")
	} else {
		t.Logf("Got expected error: checking validity of config map with nil Config: %v", err)
	}

	mc.Config = map[string]interface{}{}
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with empty Config")
	} else {
		t.Logf("Got expected error: checking validity of config map with empty Config: %v", err)
	}

	mc.Config["peers.polling.interval"] = 42.0
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map without health.polling.interval")
	} else {
		t.Logf("Got expected error: checking validity of config map without health.polling.interval: %v", err)
	}

	delete(mc.Config, "peers.polling.interval")
	mc.Config["health.polling.interval"] = 42.0
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map without peers.polling.interval")
	} else {
		t.Logf("Got expected error: checking validity of config map without peers.polling.interval: %v", err)
	}

	mc.Config["peers.polling.interval"] = 42.0
	mc.DeliveryService = nil
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with nil DeliveryService")
	} else {
		t.Logf("Got expected error: checking validity of config map with nil DeliveryService: %v", err)
	}

	mc.DeliveryService = map[string]TMDeliveryService{}
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with no DeliveryServices")
	} else {
		t.Logf("Got expected error: checking validity of config map with no DeliveryServices: %v", err)
	}

	mc.DeliveryService["a"] = TMDeliveryService{}
	mc.Profile = nil
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with nil Profile")
	} else {
		t.Logf("Got expected error: checking validity of config map with nil Profile: %v", err)
	}

	mc.Profile = map[string]TMProfile{}
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with no Profiles")
	} else {
		t.Logf("Got expected error: checking validity of config map with no Profiles: %v", err)
	}

	mc.Profile["a"] = TMProfile{}
	mc.TrafficMonitor = nil
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with nil TrafficMonitor")
	} else {
		t.Logf("Got expected error: checking validity of config map with nil TrafficMonitor: %v", err)
	}

	mc.TrafficMonitor = map[string]TrafficMonitor{}
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with no TrafficMonitors")
	} else {
		t.Logf("Got expected error: checking validity of config map with no TrafficMonitors: %v", err)
	}

	mc.TrafficMonitor["a"] = TrafficMonitor{}
	mc.TrafficServer = nil
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with nil TrafficServer")
	} else {
		t.Logf("Got expected error: checking validity of config map with nil TrafficServer: %v", err)
	}

	mc.TrafficServer = map[string]TrafficServer{}
	err = mc.Valid()
	if err == nil {
		t.Error("Didn't get expected error checking validity of config map with no TrafficServers")
	} else {
		t.Logf("Got expected error: checking validity of config map with no TrafficServers: %v", err)
	}
}

func TestTrafficMonitorTransformToMap(t *testing.T) {
	tmc := TrafficMonitorConfig{
		TrafficServers: []TrafficServer{
			{
				HostName: "testHostname",
				Interfaces: []ServerInterfaceInfo{
					{
						Name: "testInterface",
						IPAddresses: []ServerIPAddress{
							{
								Address:        "::1",
								ServiceAddress: true,
							},
						},
					},
				},
			},
		},
		CacheGroups: []TMCacheGroup{
			TMCacheGroup{},
		},
		Config: map[string]interface{}{
			"peers.polling.interval":  5.0,
			"health.polling.interval": 5.0,
		},
		TrafficMonitors: []TrafficMonitor{
			TrafficMonitor{},
		},
		DeliveryServices: []TMDeliveryService{{XMLID: "foo"}},
		Profiles: []TMProfile{
			{
				Name: "test",
				Parameters: TMParameters{
					Thresholds: map[string]HealthThreshold{
						"availableBandwidthInKbps": {
							Comparator: ">",
							Val:        42,
						},
					},
				},
			},
		},
	}

	converted, err := TrafficMonitorTransformToMap(&tmc)
	if err != nil {
		t.Fatalf("Unexpected error converting valid TrafficMonitorConfig to map: %v", err)
	}
	if converted == nil {
		t.Fatal("Null map after conversion")
	}

	if len(converted.TrafficServer) != 1 {
		t.Errorf("Incorrect number of traffic servers after conversion; expected: 1, got: %d", len(converted.TrafficServer))
	}

	if len(converted.DeliveryService) != 1 {
		t.Errorf("Incorrect number of deliveryServices after conversion; expected: 1, got: %d", len(converted.DeliveryService))
	}
	if _, ok := converted.DeliveryService["foo"]; !ok {
		t.Error("Expected delivery service 'foo' to exist in map after conversion, but it didn't")
	}

	if _, ok := converted.TrafficServer["testHostname"]; !ok {
		t.Error("Expected server 'testHostname' to exist in map after conversion, but it didn't")
	} else if len(converted.TrafficServer["testHostname"].Interfaces) != 1 {
		t.Errorf("Incorrect number of interfaces on converted traffic server; expected: 1, got: %d", len(converted.TrafficServer["testHostname"].Interfaces))
	} else if len(converted.TrafficServer["testHostname"].Interfaces[0].IPAddresses) != 1 {
		t.Errorf("Incorrect number of IP addresses on converted traffic server's interface; expected: 1, got: %d", len(converted.TrafficServer["testHostname"].Interfaces[0].IPAddresses))
	}
}
