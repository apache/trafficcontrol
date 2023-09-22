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
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestCreateServerHealthPollURL(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{
		Interfaces: []tc.ServerInterfaceInfo{
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "192.0.2.42",
						Gateway:        nil,
						ServiceAddress: true,
					},
					{
						Address:        "1::3:3:7",
						Gateway:        nil,
						ServiceAddress: true,
					},
				},
				MaxBandwidth: nil,
				MTU:          nil,
				Monitor:      true,
				Name:         "george",
			},
		},
	}

	expectedV4 := "http://192.0.2.42/_astats?application=system&inf.name=george"
	expectedV6 := "http://[1::3:3:7]/_astats?application=system&inf.name=george"
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("incorrect IPv4 polling URL; expected: '%s', actual: '%s'", expectedV4, actualV4)
	}

	if expectedV6 != actualV6 {
		t.Errorf("incorrect IPv6 polling URL; expected: '%s', actual: '%s'", expectedV6, actualV6)
	}
}

func TestCreateServerHealthPollURLTemplatePort(t *testing.T) {
	tmpl := "http://${hostname}:1234/_astats?application=&inf.name=${interface_name}"
	srv := tc.TrafficServer{
		Interfaces: []tc.ServerInterfaceInfo{
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "192.0.2.42",
						Gateway:        nil,
						ServiceAddress: true,
					},
					{
						Address:        "1::3:3:7",
						Gateway:        nil,
						ServiceAddress: true,
					},
				},
				MaxBandwidth: nil,
				MTU:          nil,
				Monitor:      true,
				Name:         "george",
			},
		},
		Port: 1234,
	}

	expectedV4 := "http://192.0.2.42:1234/_astats?application=system&inf.name=george"
	expectedV6 := "http://[1::3:3:7]:1234/_astats?application=system&inf.name=george"
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("incorrect IPv4 polling URL; expected: '%s', actual: '%s'", expectedV4, actualV4)
	}

	if expectedV6 != actualV6 {
		t.Errorf("incorrect IPv6 polling URL; expected: '%s', actual: '%s'", expectedV6, actualV6)
	}
}

func TestCreateServerHealthPollURLServerPort(t *testing.T) {
	tmpl := "http://${hostname}/_astats?application=&inf.name=${interface_name}"
	srv := tc.TrafficServer{
		HTTPSPort: 910,
		Interfaces: []tc.ServerInterfaceInfo{
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "192.0.2.42",
						Gateway:        nil,
						ServiceAddress: true,
					},
					{
						Address:        "1::3:3:7",
						Gateway:        nil,
						ServiceAddress: true,
					},
				},
				MaxBandwidth: nil,
				MTU:          nil,
				Monitor:      true,
				Name:         "george",
			},
		},
		Port: 5678,
	}

	expectedV4 := "http://192.0.2.42:5678/_astats?application=system&inf.name=george"
	expectedV6 := "http://[1::3:3:7]:5678/_astats?application=system&inf.name=george"
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("incorrect IPv4 polling URL; expected: '%s', actual: '%s'", expectedV4, actualV4)
	}

	if expectedV6 != actualV6 {
		t.Errorf("incorrect IPv6 polling URL; expected: '%s', actual: '%s'", expectedV6, actualV6)
	}

}

func TestCreateServerHealthPollURLServerPortHTTPS(t *testing.T) {
	tmpl := "hTTps://${hostname}/_astats?application=&inf.name=${interface_name}"
	srv := tc.TrafficServer{
		HTTPSPort: 910,
		Interfaces: []tc.ServerInterfaceInfo{
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "192.0.2.42",
						Gateway:        nil,
						ServiceAddress: true,
					},
					{
						Address:        "1::3:3:7",
						Gateway:        nil,
						ServiceAddress: true,
					},
				},
				MaxBandwidth: nil,
				MTU:          nil,
				Monitor:      true,
				Name:         "george",
			},
		},
		Port: 5678,
	}

	expectedV4 := "https://192.0.2.42:910/_astats?application=system&inf.name=george"
	expectedV6 := "https://[1::3:3:7]:910/_astats?application=system&inf.name=george"
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("incorrect IPv4 polling URL; expected: '%s', actual: '%s'", expectedV4, actualV4)
	}

	if expectedV6 != actualV6 {
		t.Errorf("incorrect IPv6 polling URL; expected: '%s', actual: '%s'", expectedV6, actualV6)
	}
}

func TestCreateServerHealthPollURLTemplateAndServerPort(t *testing.T) {

	// if both template and server ports exist, template takes precedence

	tmpl := "http://${hostname}:1234/_astats?application=&inf.name=${interface_name}"
	srv := tc.TrafficServer{
		HTTPSPort: 910,
		Interfaces: []tc.ServerInterfaceInfo{
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "192.0.2.42",
						Gateway:        nil,
						ServiceAddress: true,
					},
					{
						Address:        "1::3:3:7",
						Gateway:        nil,
						ServiceAddress: true,
					},
				},
				MaxBandwidth: nil,
				MTU:          nil,
				Monitor:      true,
				Name:         "george",
			},
		},
		Port: 5678,
	}

	expectedV4 := "http://192.0.2.42:1234/_astats?application=system&inf.name=george"
	expectedV6 := "http://[1::3:3:7]:1234/_astats?application=system&inf.name=george"
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("incorrect IPv4 polling URL; expected: '%s', actual: '%s'", expectedV4, actualV4)
	}

	if expectedV6 != actualV6 {
		t.Errorf("incorrect IPv6 polling URL; expected: '%s', actual: '%s'", expectedV6, actualV6)
	}
}

func TestCreateServerStatPollURL(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{
		Interfaces: []tc.ServerInterfaceInfo{
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "192.0.2.42",
						Gateway:        nil,
						ServiceAddress: true,
					},
					{
						Address:        "1::3:3:7",
						Gateway:        nil,
						ServiceAddress: true,
					},
				},
				MaxBandwidth: nil,
				MTU:          nil,
				Monitor:      true,
				Name:         "george",
			},
		},
	}

	expectedV4 := "http://192.0.2.42/_astats?application=&inf.name=george"
	expectedV6 := "http://[1::3:3:7]/_astats?application=&inf.name=george"

	healthV4, healthV6 := createServerHealthPollURLs(tmpl, srv)
	actualV4 := createServerStatPollURL(healthV4)
	actualV6 := createServerStatPollURL(healthV6)

	if expectedV4 != actualV4 {
		t.Errorf("incorrect IPv4 polling URL; expected: '%s', actual: '%s'", expectedV4, actualV4)
	}

	if expectedV6 != actualV6 {
		t.Errorf("incorrect IPv6 polling URL; expected: '%s', actual: '%s'", expectedV6, actualV6)
	}
}

func TestGetCacheGroupsToPoll(t *testing.T) {
	monitors := map[string]tc.TrafficMonitor{
		"tm2": {
			HostName:     "tm2",
			Location:     "tm-group-2",
			ServerStatus: "ONLINE",
		},
		"tm1": {
			HostName:     "tm1",
			Location:     "tm-group-1",
			ServerStatus: "ONLINE",
		},
		"tm3": {
			HostName:     "tm3",
			Location:     "tm-group-3",
			ServerStatus: "ONLINE",
		},
	}
	caches := map[string]tc.TrafficServer{
		"cache2": {
			CacheGroup:   "cache-group-2",
			ServerStatus: "REPORTED",
		},
		"cache1": {
			CacheGroup:   "cache-group-1",
			ServerStatus: "REPORTED",
		},
		"cache3": {
			CacheGroup:   "cache-group-3",
			ServerStatus: "REPORTED",
		},
	}
	cacheGroups := map[string]tc.TMCacheGroup{
		"tm-group-2": {
			Name: "tm-group-2",
			Coordinates: tc.MonitoringCoordinates{
				Latitude:  38.39,
				Longitude: -99.58,
			},
		},
		"tm-group-1": {
			Name: "tm-group-1",
			Coordinates: tc.MonitoringCoordinates{
				Latitude:  37.32,
				Longitude: -121.34,
			},
		},
		"tm-group-3": {
			Name: "tm-group-3",
			Coordinates: tc.MonitoringCoordinates{
				Latitude:  37.22,
				Longitude: -77.53,
			},
		},
		"cache-group-3": {
			Name: "cache-group-3",
			Coordinates: tc.MonitoringCoordinates{
				Latitude:  41.93,
				Longitude: -74.17,
			},
		},
		"cache-group-1": {
			Name: "cache-group-1",
			Coordinates: tc.MonitoringCoordinates{
				Latitude:  35.04,
				Longitude: -120.12,
			},
		},
		"cache-group-2": {
			Name: "cache-group-2",
			Coordinates: tc.MonitoringCoordinates{
				Latitude:  40.92,
				Longitude: -98.49,
			},
		},
	}
	type testCase struct {
		DistributedPolling bool
		TMName             string
		ExpectedTMGroup    string
		ExpectedToPoll     map[string]tc.TMCacheGroup
		ExpectErr          bool
	}

	for _, tc := range []testCase{
		{
			DistributedPolling: true,
			TMName:             "tm1",
			ExpectedTMGroup:    "tm-group-1",
			ExpectedToPoll: map[string]tc.TMCacheGroup{
				"cache-group-1": cacheGroups["cache-group-1"],
			},
			ExpectErr: false,
		},
		{
			DistributedPolling: true,
			TMName:             "tm2",
			ExpectedTMGroup:    "tm-group-2",
			ExpectedToPoll: map[string]tc.TMCacheGroup{
				"cache-group-2": cacheGroups["cache-group-2"],
			},
			ExpectErr: false,
		},
		{
			DistributedPolling: true,
			TMName:             "tm3",
			ExpectedTMGroup:    "tm-group-3",
			ExpectedToPoll: map[string]tc.TMCacheGroup{
				"cache-group-3": cacheGroups["cache-group-3"],
			},
			ExpectErr: false,
		},
		{
			DistributedPolling: false,
			TMName:             "tm3",
			ExpectedTMGroup:    "tm-group-3",
			ExpectedToPoll: map[string]tc.TMCacheGroup{
				"cache-group-1": cacheGroups["cache-group-1"],
				"cache-group-2": cacheGroups["cache-group-2"],
				"cache-group-3": cacheGroups["cache-group-3"],
			},
			ExpectErr: false,
		},
	} {
		tmGroup, tmStatus, toPoll, err := getCacheGroupsToPoll(tc.DistributedPolling, tc.TMName, monitors, caches, cacheGroups)
		if tc.ExpectErr != (err != nil) {
			t.Errorf("getting cachegroups to poll -- expect error: %t, actual error: %v", tc.ExpectErr, err)
		}
		if tc.ExpectedTMGroup != tmGroup {
			t.Errorf("getting TM group -- expected: %s, actual: %s", tc.ExpectedTMGroup, tmGroup)
		}
		if tmStatus != "ONLINE" {
			t.Errorf("expected TM status: ONLINE, actual: %s", tmStatus)
		}
		if !reflect.DeepEqual(tc.ExpectedToPoll, toPoll) {
			t.Errorf("getting cachegroups to poll -- expected: %+v, actual: %+v", tc.ExpectedToPoll, toPoll)
		}
	}
}
