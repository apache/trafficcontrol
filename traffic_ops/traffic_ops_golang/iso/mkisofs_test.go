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
	"bytes"
	"net"
	"strings"
	"testing"
)

func TestWriteNetworkCfg(t *testing.T) {
	cases := []struct {
		name        string
		input       isoRequest
		nameservers []string
		expected    string
	}{
		{
			"empty",
			isoRequest{},
			nil,

			`
IPADDR=""
NETMASK=""
GATEWAY=""
DEVICE=""
MTU="0"
NAMESERVER=""
HOSTNAME=""
NETWORKING_IPV6="yes"
IPV6ADDR=""
IPV6_DEFAULTGW=""
DHCP="no"
`,
		},

		{
			"no domain",
			isoRequest{
				IPAddr:        net.IP{192, 168, 1, 2},
				IPNetmask:     net.IP{255, 255, 255, 0},
				IPGateway:     net.IP{192, 168, 1, 255},
				InterfaceName: "eth0",
				InterfaceMTU:  1500,
				HostName:      "test.server",
				DomainName:    "",
				IP6Address:    net.ParseIP("beef::1"),
				IP6Gateway:    net.ParseIP("::1"),
				DHCP:          boolStr{true, true},
			},
			[]string{"8.8.8.8", "1.1.1.1"},

			`
IPADDR="192.168.1.2"
NETMASK="255.255.255.0"
GATEWAY="192.168.1.255"
DEVICE="eth0"
MTU="1500"
NAMESERVER="8.8.8.8,1.1.1.1"
HOSTNAME="test.server"
NETWORKING_IPV6="yes"
IPV6ADDR="beef::1"
IPV6_DEFAULTGW="::1"
DHCP="yes"
`,
		},

		{
			"non-bonded",
			isoRequest{
				IPAddr:        net.IP{192, 168, 1, 2},
				IPNetmask:     net.IP{255, 255, 255, 0},
				IPGateway:     net.IP{192, 168, 1, 255},
				InterfaceName: "eth0",
				InterfaceMTU:  1500,
				HostName:      "test.server",
				DomainName:    "example.com",
				IP6Address:    net.ParseIP("beef::1"),
				IP6Gateway:    net.ParseIP("::1"),
				DHCP:          boolStr{true, true},
			},
			[]string{"8.8.8.8", "1.1.1.1"},

			`
IPADDR="192.168.1.2"
NETMASK="255.255.255.0"
GATEWAY="192.168.1.255"
DEVICE="eth0"
MTU="1500"
NAMESERVER="8.8.8.8,1.1.1.1"
HOSTNAME="test.server.example.com"
NETWORKING_IPV6="yes"
IPV6ADDR="beef::1"
IPV6_DEFAULTGW="::1"
DHCP="yes"
`,
		},

		{
			"bonded",
			isoRequest{
				IPAddr:        net.IP{192, 168, 1, 2},
				IPNetmask:     net.IP{255, 255, 255, 0},
				IPGateway:     net.IP{192, 168, 1, 255},
				InterfaceName: "bond01",
				InterfaceMTU:  1500,
				HostName:      "test.server",
				DomainName:    "",
				IP6Address:    net.ParseIP("beef::1"),
				IP6Gateway:    net.ParseIP("::1"),
				DHCP:          boolStr{true, true},
			},
			[]string{"8.8.8.8", "1.1.1.1"},

			`
IPADDR="192.168.1.2"
NETMASK="255.255.255.0"
GATEWAY="192.168.1.255"
BOND_DEVICE="bond01"
MTU="1500"
NAMESERVER="8.8.8.8,1.1.1.1"
HOSTNAME="test.server"
NETWORKING_IPV6="yes"
IPV6ADDR="beef::1"
IPV6_DEFAULTGW="::1"
BONDING_OPTS="miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4"
DHCP="yes"
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writeNetworkCfg(&w, tc.input, tc.nameservers); err != nil {
				t.Fatalf("writeNetworkCfg() err = %v", err)
			}
			got := w.String()
			expected := strings.TrimSpace(tc.expected)

			if got != expected {
				t.Fatalf("writeNetworkCfg() got != expected\n got:\n%s\n expected:\n%s", got, expected)
			}
			t.Logf("writeNetworkCfg():\n%s", got)
		})
	}
}

func TestWriteMgmtNetworkCfg(t *testing.T) {
	cases := []struct {
		name     string
		input    isoRequest
		expected string
	}{
		{
			"empty",
			isoRequest{},

			`
IPADDR=""
NETMASK=""
GATEWAY=""
DEVICE=""
`,
		},

		{
			"IPv4",
			isoRequest{
				MgmtIPAddress: net.IP{192, 168, 2, 3},
				MgmtIPNetmask: net.IP{255, 255, 255, 255},
				MgmtIPGateway: net.IP{192, 168, 1, 255},
				MgmtInterface: "eth0",
			},

			`
IPADDR="192.168.2.3"
NETMASK="255.255.255.255"
GATEWAY="192.168.1.255"
DEVICE="eth0"
`,
		},

		{
			"IPv6",
			isoRequest{
				MgmtIPAddress: net.ParseIP("beef::1"),
				MgmtIPNetmask: net.IP{255, 255, 255, 255},
				MgmtIPGateway: net.IP{192, 168, 1, 255},
				MgmtInterface: "eth0",
			},

			`
IPV6ADDR="beef::1"
NETMASK="255.255.255.255"
GATEWAY="192.168.1.255"
DEVICE="eth0"
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writeMgmtNetworkCfg(&w, tc.input); err != nil {
				t.Fatalf("writeMgmtNetworkCfg() err = %v", err)
			}
			got := w.String()
			expected := strings.TrimSpace(tc.expected)

			if got != expected {
				t.Fatalf("writeMgmtNetworkCfg() got != expected\n got:\n%s\n expected:\n%s", got, expected)
			}
			t.Logf("writeMgmtNetworkCfg():\n%s", got)
		})
	}
}

func TestWriteDiskCfg(t *testing.T) {
	cases := []struct {
		name     string
		input    isoRequest
		expected string
	}{
		{
			"empty",
			isoRequest{},
			`
boot_drives=""
`,
		},

		{
			"non-empty",
			isoRequest{
				Disk: "sda1",
			},
			`
boot_drives="sda1"
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writeDiskCfg(&w, tc.input); err != nil {
				t.Fatalf("writeDiskCfg() err = %v", err)
			}
			got := w.String()
			expected := strings.TrimSpace(tc.expected)

			if got != expected {
				t.Fatalf("writeDiskCfg() got != expected\n got:\n%s\n expected:\n%s", got, expected)
			}
			t.Logf("writeDiskCfg():\n%s", got)
		})
	}
}

func TestWritePasswordCfg(t *testing.T) {
	cases := []struct {
		name     string
		input    isoRequest
		salt     string
		expected string
	}{
		{
			"empty",
			isoRequest{},
			"salt",
			"rootpw --iscrypted $1$salt$UsdFqFVB.FsuinRDK5eE..\n",
		},

		{
			"non-empty",
			isoRequest{
				RootPass: "Traffic Ops",
			},
			"salt",
			"rootpw --iscrypted $1$salt$17HeaymOIi.65dl76MkK01\n",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writePasswordCfg(&w, tc.input, tc.salt); err != nil {
				t.Fatalf("writePasswordCfg() err = %v", err)
			}
			got := w.String()

			if got != tc.expected {
				t.Fatalf("writePasswordCfg() got != expected\n got:\n%s\n expected:\n%s", got, tc.expected)
			}
			t.Logf("writePasswordCfg():\n%q", got)
		})
	}
}

func TestWritePasswordCfg_rndSalt(t *testing.T) {
	cases := []struct {
		name  string
		input isoRequest
	}{
		{
			"empty",
			isoRequest{},
		},

		{
			"non-empty",
			isoRequest{
				RootPass: "Traffic Ops",
			},
		},
		{
			"long",
			isoRequest{
				RootPass: "this is a long password made longer even now",
			},
		},
	}

	const (
		expectedPrefix = "rootpw --iscrypted $1$"
		expectedPWLen  = 32
	)

	// Ensure use of random salt generates correct looking passwords.

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var w bytes.Buffer
			if err := writePasswordCfg(&w, tc.input, ""); err != nil {
				t.Fatalf("writePasswordCfg() err = %v", err)
			}
			got := w.String()

			if !strings.HasPrefix(got, expectedPrefix) {
				t.Fatalf("writePasswordCfg() got: %q\nexpected prefix of: %q", got, expectedPrefix)
			}
			if pwLen := len(got) - len(expectedPrefix); pwLen != expectedPWLen {
				t.Fatalf("writePasswordCfg() got: %q with password length %d\nexpected password length of at least: %d", got, pwLen, expectedPWLen)
			}

			t.Logf("writePasswordCfg():\n%q", got)
		})
	}
}
