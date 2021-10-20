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
	"testing"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
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
