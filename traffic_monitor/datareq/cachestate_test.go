package datareq

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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
)

func TestCreateCacheStatusesForKbps(t *testing.T) {
	var cacheTypes map[tc.CacheName]tc.CacheType
	var statInfoHistory cache.ResultInfoHistory
	statResultHistory := threadsafe.NewResultStatHistory()
	healthHistory := make(map[tc.CacheName][]cache.Result, 0)
	cacheResult := cache.Result{
		Available:       true,
		Error:           nil,
		ID:              "1",
		Miscellaneous:   nil,
		PollFinished:    nil,
		PollID:          0,
		PrecomputedData: cache.PrecomputedData{},
		RequestTime:     0,
		Statistics:      cache.Statistics{},
		Time:            time.Now(),
		UsingIPv4:       true,
		Vitals:          cache.Vitals{},
		InterfaceVitals: map[string]cache.Vitals{"interface1": cache.Vitals{
			LoadAvg:    23.2,
			BytesOut:   200,
			BytesIn:    140,
			KbpsOut:    300,
			MaxKbpsOut: 1000,
		}},
	}
	healthHistory["edgeserver"] = []cache.Result{cacheResult}

	var lastHealthDurations map[tc.CacheName]time.Duration
	localCacheStatusThreadsafe := threadsafe.NewCacheAvailableStatus()
	statMaxKbpses := threadsafe.NewCacheKbpses()
	servers := make(map[string]tc.TrafficServer, 0)
	interfaces := make([]tc.ServerInterfaceInfo, 0)
	ipAddresses := make([]tc.ServerIPAddress, 0)
	ipAddresses = append(ipAddresses, tc.ServerIPAddress{
		Address:        "123.24.25.26",
		Gateway:        util.StrPtr("255.255.0.0"),
		ServiceAddress: true,
	})
	i := tc.ServerInterfaceInfo{
		IPAddresses:  ipAddresses,
		MaxBandwidth: util.Uint64Ptr(1000),
		Monitor:      false,
		MTU:          util.Uint64Ptr(9000),
		Name:         "interface1",
	}
	interfaces = append(interfaces, i)

	s := tc.TrafficServer{
		CacheGroup:       "cg",
		DeliveryServices: nil,
		FQDN:             "fqdn",
		HashID:           "hashID",
		HostName:         "hostName",
		HTTPSPort:        443,
		Interfaces:       interfaces,
		Port:             8080,
		Profile:          "profile",
		ServerStatus:     "REPORTED",
		Type:             "EDGE",
	}
	servers["edgeserver"] = s
	result := createCacheStatuses(cacheTypes,
		statInfoHistory,
		statResultHistory,
		healthHistory,
		lastHealthDurations,
		localCacheStatusThreadsafe,
		statMaxKbpses,
		servers)

	if len(result) != 1 {
		t.Fatalf("expected only one cache in result, but got %d", len(result))
	}
	if status, ok := result["edgeserver"]; !ok {
		t.Fatalf("result status did not contain the expected key 'edgeserver'")
	} else {
		if status.BandwidthKbps == nil {
			t.Fatalf("expected a valid value in BandwidthKbps, but got nothing")
		}
		if *status.BandwidthKbps != 300 {
			t.Errorf("expected BandwidthKbps to be equal to the sum of the values in the interfaces (300), but got %f", *status.BandwidthKbps)
		}
	}

	// Add another interface to the server and test again
	ipAddresses = append(ipAddresses, tc.ServerIPAddress{
		Address:        "223.24.25.26",
		Gateway:        util.StrPtr("255.255.0.0"),
		ServiceAddress: true,
	})

	i = tc.ServerInterfaceInfo{
		IPAddresses:  ipAddresses,
		MaxBandwidth: util.Uint64Ptr(1000),
		Monitor:      false,
		MTU:          util.Uint64Ptr(9000),
		Name:         "interface2",
	}
	interfaces = append(interfaces, i)
	s.Interfaces = interfaces
	servers["edgeserver"] = s
	cacheResult.InterfaceVitals["interface2"] = cache.Vitals{
		LoadAvg:    19.23,
		BytesOut:   45,
		BytesIn:    40,
		KbpsOut:    500,
		MaxKbpsOut: 1000,
	}
	healthHistory["edgeserver"] = []cache.Result{cacheResult}

	result = createCacheStatuses(cacheTypes,
		statInfoHistory,
		statResultHistory,
		healthHistory,
		lastHealthDurations,
		localCacheStatusThreadsafe,
		statMaxKbpses,
		servers)

	if len(result) != 1 {
		t.Fatalf("expected only one cache in result, but got %d", len(result))
	}
	if status, ok := result["edgeserver"]; !ok {
		t.Fatalf("result status did not contain the expected key 'edgeserver'")
	} else {
		if status.BandwidthKbps == nil {
			t.Fatalf("expected a valid value in BandwidthKbps, but got nothing")
		}
		if *status.BandwidthKbps != 800 {
			t.Errorf("expected BandwidthKbps to be equal to the sum of the values in the interfaces (800), but got %f", *status.BandwidthKbps)
		}
	}
}
