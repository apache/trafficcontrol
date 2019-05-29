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
	srv := tc.TrafficServer{IP: "192.0.2.42", InterfaceName: "george"}

	expected := `http://` + srv.IP + `/_astats?application=system&inf.name=` + srv.InterfaceName
	actual := createServerHealthPollURL(tmpl, srv)

	if expected != actual {
		t.Errorf("expected createServerHealthPollURL '" + expected + "' actual: '" + actual + "'")
	}
}

func TestCreateServerHealthPollURLTemplatePort(t *testing.T) {
	tmpl := `http://${hostname}:1234/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{IP: "192.0.2.42", InterfaceName: "george"}

	expected := `http://` + srv.IP + `:1234/_astats?application=system&inf.name=` + srv.InterfaceName
	actual := createServerHealthPollURL(tmpl, srv)

	if expected != actual {
		t.Errorf("expected createServerHealthPollURL '" + expected + "' actual: '" + actual + "'")
	}
}

func TestCreateServerHealthPollURLServerPort(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{IP: "192.0.2.42", Port: 5678, HTTPSPort: 910, InterfaceName: "george"}

	expected := `http://` + srv.IP + ":" + strconv.Itoa(srv.Port) + `/_astats?application=system&inf.name=` + srv.InterfaceName
	actual := createServerHealthPollURL(tmpl, srv)

	if expected != actual {
		t.Errorf("expected createServerHealthPollURL '" + expected + "' actual: '" + actual + "'")
	}
}

func TestCreateServerHealthPollURLServerPortHTTPS(t *testing.T) {
	tmpl := `hTTps://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{IP: "192.0.2.42", Port: 5678, HTTPSPort: 910, InterfaceName: "george"}

	expected := `https://` + srv.IP + ":" + strconv.Itoa(srv.HTTPSPort) + `/_astats?application=system&inf.name=` + srv.InterfaceName
	actual := createServerHealthPollURL(tmpl, srv)

	if expected != actual {
		t.Errorf("expected createServerHealthPollURL '" + expected + "' actual: '" + actual + "'")
	}
}

func TestCreateServerHealthPollURLTemplateAndServerPort(t *testing.T) {

	// if both template and server ports exist, template takes precedence

	tmpl := `http://${hostname}:1234/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{IP: "192.0.2.42", Port: 5678, HTTPSPort: 910, InterfaceName: "george"}

	expected := `http://` + srv.IP + `:1234/_astats?application=system&inf.name=` + srv.InterfaceName
	actual := createServerHealthPollURL(tmpl, srv)

	if expected != actual {
		t.Errorf("expected createServerHealthPollURL '" + expected + "' actual: '" + actual + "'")
	}
}

func TestCreateServerStatPollURL(t *testing.T) {
	tmpl := `http://${hostname}/_astats?application=&inf.name=${interface_name}`
	srv := tc.TrafficServer{IP: "192.0.2.42", InterfaceName: "george"}

	expected := `http://` + srv.IP + `/_astats?application=&inf.name=` + srv.InterfaceName
	actual := createServerStatPollURL(createServerHealthPollURL(tmpl, srv))

	if expected != actual {
		t.Errorf("expected createServerStatPollURL '" + expected + "' actual: '" + actual + "'")
	}
}
