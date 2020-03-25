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

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCreateServerHealthPollURL(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{IP: "192.0.2.42", IP6: "1::3:3:7", InterfaceName: "george"}

	expectedV4 := `http://` + srv.IP + `/_astats?application=system&inf.name=` + srv.InterfaceName
	expectedV6 := `http://[` + srv.IP6 + `]/_astats?application=system&inf.name=` + srv.InterfaceName
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
	srv := tc.TrafficServer{IP: "192.0.2.42", IP6: "1::3:3:7", InterfaceName: "george"}

	expectedV4 := `http://` + srv.IP + `:1234/_astats?application=system&inf.name=` + srv.InterfaceName
	expectedV6 := `http://[` + srv.IP6 + `]:1234/_astats?application=system&inf.name=` + srv.InterfaceName
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
	srv := tc.TrafficServer{IP: "192.0.2.42", IP6: "1::3:3:7", Port: 5678, HTTPSPort: 910, InterfaceName: "george"}

	expectedV4 := `http://` + srv.IP + ":" + strconv.Itoa(srv.Port) + `/_astats?application=system&inf.name=` + srv.InterfaceName
	expectedV6 := `http://[` + srv.IP6 + "]:" + strconv.Itoa(srv.Port) + `/_astats?application=system&inf.name=` + srv.InterfaceName
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
	srv := tc.TrafficServer{IP: "192.0.2.42", IP6: "1::3:3:7", Port: 5678, HTTPSPort: 910, InterfaceName: "george"}

	expectedV4 := `https://` + srv.IP + ":" + strconv.Itoa(srv.HTTPSPort) + `/_astats?application=system&inf.name=` + srv.InterfaceName
	expectedV6 := `https://[` + srv.IP6 + "]:" + strconv.Itoa(srv.HTTPSPort) + `/_astats?application=system&inf.name=` + srv.InterfaceName
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
	srv := tc.TrafficServer{IP: "192.0.2.42", IP6: "1::3:3:7", Port: 5678, HTTPSPort: 910, InterfaceName: "george"}

	expectedV4 := `http://` + srv.IP + `:1234/_astats?application=system&inf.name=` + srv.InterfaceName
	expectedV6 := `http://[` + srv.IP6 + `]:1234/_astats?application=system&inf.name=` + srv.InterfaceName
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
	srv := tc.TrafficServer{IP: "192.0.2.42", IP6: "1::3:3:7", InterfaceName: "george"}

	expectedV4 := `http://` + srv.IP + `/_astats?application=&inf.name=` + srv.InterfaceName
	expectedV6 := `http://[` + srv.IP6 + `]/_astats?application=&inf.name=` + srv.InterfaceName

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
