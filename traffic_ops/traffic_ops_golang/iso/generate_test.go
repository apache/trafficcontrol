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
	"net"
	"testing"
)

func TestISORequest_validate(t *testing.T) {
	cases := []struct {
		name             string
		expectedValidate bool
		input            isoRequest
	}{
		{
			"valid with dhcp false",
			true,
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
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
		{
			"valid with dhcp true",
			true,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
		{
			"invalid with dhcp false",
			false,
			isoRequest{
				DHCP:          boolStr{true, false},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{},
				MgmtIPGateway: net.IP{},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},

		{
			"valid with mgmt addr",
			true,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "not empty",
				MgmtIPAddress: net.IP{192, 168, 0, 1},
				MgmtIPGateway: net.IP{192, 168, 0, 2},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},
		{
			"invalid with mgmt addr",
			false,
			isoRequest{
				DHCP:          boolStr{true, true},
				Stream:        boolStr{true, true},
				OSVersionDir:  "centos72",
				HostName:      "db",
				DomainName:    "infra.ciab.test",
				IPAddr:        net.IP{},
				InterfaceMTU:  1500,
				InterfaceName: "eth0",
				IP6Address:    net.IP{},
				IP6Gateway:    net.IP{},
				IPGateway:     net.IP{},
				IPNetmask:     net.IP{},
				MgmtInterface: "",
				MgmtIPAddress: net.IP{192, 168, 0, 1},
				MgmtIPGateway: net.IP{192, 168, 0, 2},
				MgmtIPNetmask: net.IP{},
				Disk:          "sda",
				RootPass:      "12345678",
			},
		},

		{
			"invalid with zero values",
			false,
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
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotErrs := tc.input.validate()
			if tc.expectedValidate != (len(gotErrs) == 0) {
				t.Fatalf("isoRequest.validate() = %v; expected errors %v", gotErrs, tc.expectedValidate)
			}
			t.Logf("isoRequest.validate() = %+v", gotErrs)
		})
	}
}

func TestBoolStr_UnmarshalText(t *testing.T) {
	cases := []struct {
		input    string
		expected boolStr
	}{
		{
			`no`,
			boolStr{isSet: true, val: false},
		},
		{
			`No`,
			boolStr{isSet: true, val: false},
		},
		{
			`YES`,
			boolStr{isSet: true, val: true},
		},
		{
			`other`,
			boolStr{isSet: false, val: false},
		},
		{
			``,
			boolStr{isSet: false, val: false},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			var got boolStr
			if err := got.UnmarshalText([]byte(tc.input)); err != nil {
				t.Fatal(err)
			}

			if got != tc.expected {
				t.Fatalf("got %+v; expected %+v", got, tc.expected)
			}
			t.Logf("got %+v", got)
		})
	}
}
