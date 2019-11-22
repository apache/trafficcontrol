package iso

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
	"net"
	"strings"
	"testing"
)

func TestISORequest(t *testing.T) {
	cases := []struct {
		input            string
		expected         isoRequest
		expectedValidate bool
	}{
		{
			`{
				"dhcp": "no",
				"stream": "yes",
				"osversionDir": "centos72",
				"hostName": "db",
				"domainName": "infra.ciab.test",
				"ipAddress": "172.20.0.4",
				"interfaceMtu": 1500,
				"interfaceName": "eth0",
				"ip6Address": null,
				"ip6Gateway": "",
				"ipGateway": "172.20.0.1",
				"ipNetmask": "255.255.0.0",
				"mgmtIpAddress": "",
				"mgmtIpGateway": "",
				"mgmtIpNetmask": "",
				"disk": "sda",
				"rootPass": "12345678"
			}`,
			isoRequest{
				DHCP:          boolStr{true, false},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{172, 20, 0, 4},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{172, 20, 0, 1},
				IPNetmask:     net.IP{255, 255, 0, 0},
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
			true,
		},

		{
			`{
				"dhcp": null,
				"stream": "",
				"osversionDir": "",
				"hostName": "",
				"domainName": "",
				"ipAddress": "",
				"interfaceMtu": 0,
				"interfaceName": "",
				"ip6Address": null,
				"ip6Gateway": "",
				"ipGateway": "",
				"ipNetmask": "",
				"mgmtIpAddress": "",
				"mgmtIpGateway": "",
				"mgmtIpNetmask": "",
				"disk": "",
				"rootPass": ""
			}`,
			isoRequest{
				DHCP:          boolStr{false, false},
				Stream:        boolStr{false, false},
				OSVersionDir:  "",
				HostName:      "",
				DomainName:    "",
				IPAddr:        net.IP{},
				InterfaceMTU:  0,
				InterfaceName: "",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "",
				RootPass:      "",
			},
			false,
		},
	}

	for _, tc := range cases {
		var got isoRequest
		if err := json.NewDecoder(strings.NewReader(tc.input)).Decode(&got); err != nil {
			t.Fatalf("unexpected error decoding input: %s", err)
		}
		if !got.equal(tc.expected) {
			t.Fatalf("got isoRequest not equal to expected\ngot:\n%+v\nexpected:\n%+v", got, tc.expected)
		}

		gotErrs := got.validate()
		if tc.expectedValidate != (len(gotErrs) == 0) {
			t.Fatalf("isoRequest.validate() = %v; expected errors %v", gotErrs, tc.expectedValidate)
		}

		t.Logf("newISORequest() = %+v\nisoRequest.validate() = %+v", got, gotErrs)
	}
}
