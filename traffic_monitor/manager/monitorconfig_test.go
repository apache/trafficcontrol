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
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func createTrafficServer(withPorts bool) tc.TrafficServer {
	server := tc.TrafficServer{
		Interfaces: []tc.InterfaceInfo{
			{
				IPAddresses: []tc.IPAddress{
					{
						Address:        "192.0.2.24",
						Gateway:        util.StrPtr("192.0.2.0/24"),
						ServiceAddress: true,
					},
					{
						Address:        "1::3:3:7",
						Gateway:        nil,
						ServiceAddress: true,
					},
				},
				MaxBandwidth: util.UInt64Ptr(2500),
				Monitor:      true,
				MTU:          util.UInt64Ptr(9000),
				Name:         "george",
			},
		},
	}
	if withPorts {
		server.Port = 5678
		server.HTTPSPort = 910
	}
	return server
}

func TestCreateServerHealthPollURL(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := createTrafficServer(false)

	interf := tc.GetVIPInterface(srv)
	ip4, ip6 := interf.GetDefaultAddress()
	name := interf.Name

	expectedV4 := `http://` + ip4 + `/_astats?application=system&inf.name=` + name
	expectedV6 := `http://[` + ip6 + `]/_astats?application=system&inf.name=` + name
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("for IPv4 expected createServerHealthPollURL '" + expectedV4 + "' actual: '" + actualV4 + "'")
	}

	if expectedV6 != actualV6 {
		t.Errorf("for IPv6 expected createServerHealthPollURL '" + expectedV6 + "' actual: '" + actualV6 + "'")
	}
}

func TestCreateServerHealthPollURLTemplatePort(t *testing.T) {
	tmpl := `http://${hostname}:1234/_astats?application=&inf.name=${interface_name}`
	srv := createTrafficServer(false)

	interf := tc.GetVIPInterface(srv)
	ip4, ip6 := interf.GetDefaultAddress()
	name := interf.Name
	expectedV4 := `http://` + ip4 + `:1234/_astats?application=system&inf.name=` + name
	expectedV6 := `http://[` + ip6 + `]:1234/_astats?application=system&inf.name=` + name
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("for IPv4 expected createServerHealthPollURL '" + expectedV4 + "' actual: '" + actualV4 + "'")
	}

	if expectedV6 != actualV6 {
		t.Errorf("for IPv6 expected createServerHealthPollURL '" + expectedV6 + "' actual: '" + actualV6 + "'")
	}
}

func TestCreateServerHealthPollURLServerPort(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := createTrafficServer(true)

	interf := tc.GetVIPInterface(srv)
	ip4, ip6 := interf.GetDefaultAddress()
	name := interf.Name
	expectedV4 := `http://` + ip4 + ":" + strconv.Itoa(srv.Port) + `/_astats?application=system&inf.name=` + name
	expectedV6 := `http://[` + ip6 + "]:" + strconv.Itoa(srv.Port) + `/_astats?application=system&inf.name=` + name
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("for IPv4 expected createServerHealthPollURL '" + expectedV4 + "' actual: '" + actualV4 + "'")
	}
	if expectedV6 != actualV6 {
		t.Errorf("for IPv6 expected createServerHealthPollURL '" + expectedV6 + "' actual: '" + actualV6 + "'")
	}

}

func TestCreateServerHealthPollURLServerPortHTTPS(t *testing.T) {
	tmpl := `hTTps://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := createTrafficServer(true)

	interf := tc.GetVIPInterface(srv)
	ip4, ip6 := interf.GetDefaultAddress()
	name := interf.Name
	expectedV4 := `https://` + ip4 + ":" + strconv.Itoa(srv.HTTPSPort) + `/_astats?application=system&inf.name=` + name
	expectedV6 := `https://[` + ip6 + "]:" + strconv.Itoa(srv.HTTPSPort) + `/_astats?application=system&inf.name=` + name
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("for IPv4 expected createServerHealthPollURL '" + expectedV4 + "' actual: '" + actualV4 + "'")
	}

	if expectedV6 != actualV6 {
		t.Errorf("for IPv6 expected createServerHealthPollURL '" + expectedV6 + "' actual: '" + actualV6 + "'")
	}
}

func TestCreateServerHealthPollURLTemplateAndServerPort(t *testing.T) {
	// if both template and server ports exist, template takes precedence
	tmpl := `http://${hostname}:1234/_astats?application=&inf.name=${interface_name}`
	srv := createTrafficServer(true)

	interf := tc.GetVIPInterface(srv)
	ip4, ip6 := interf.GetDefaultAddress()
	name := interf.Name

	expectedV4 := `http://` + ip4 + `:1234/_astats?application=system&inf.name=` + name
	expectedV6 := `http://[` + ip6 + `]:1234/_astats?application=system&inf.name=` + name
	actualV4, actualV6 := createServerHealthPollURLs(tmpl, srv)

	if expectedV4 != actualV4 {
		t.Errorf("for IPv4 expected createServerHealthPollURL '" + expectedV4 + "' actual: '" + actualV4 + "'")
	}

	if expectedV6 != actualV6 {
		t.Errorf("for IPv6 expected createServerHealthPollURL '" + expectedV6 + "' actual: '" + actualV6 + "'")
	}
}

func TestCreateServerStatPollURL(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := createTrafficServer(false)

	interf := tc.GetVIPInterface(srv)
	ip4, ip6 := interf.GetDefaultAddress()
	name := interf.Name

	expectedV4 := `http://` + ip4 + `/_astats?application=&inf.name=` + name
	expectedV6 := `http://[` + ip6 + `]/_astats?application=&inf.name=` + name

	healthV4, healthV6 := createServerHealthPollURLs(tmpl, srv)
	actualV4 := createServerStatPollURL(healthV4)
	actualV6 := createServerStatPollURL(healthV6)

	if expectedV4 != actualV4 {
		t.Errorf("for IPv4 expected createServerStatPollURL '" + expectedV4 + "' actual: '" + actualV4 + "'")
	}

	if expectedV6 != actualV6 {
		t.Errorf("for IPv6 expected createServerStatPollURL '" + expectedV6 + "' actual: '" + actualV6 + "'")
	}
}
