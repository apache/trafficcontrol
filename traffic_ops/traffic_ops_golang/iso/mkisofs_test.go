package iso

import (
	"bytes"
	"net"
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

			`IPADDR=""
NETMASK=""
GATEWAY=""
DEVICE=""
MTU="0"
NAMESERVER=""
HOSTNAME=""
NETWORKING_IPV6="yes"
IPV6ADDR=""
IPV6_DEFAULTGW=""
DHCP="no"`,
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

			`IPADDR="192.168.1.2"
NETMASK="255.255.255.0"
GATEWAY="192.168.1.255"
DEVICE="eth0"
MTU="1500"
NAMESERVER="8.8.8.8,1.1.1.1"
HOSTNAME="test.server.example.com"
NETWORKING_IPV6="yes"
IPV6ADDR="beef::1"
IPV6_DEFAULTGW="::1"
DHCP="yes"`,
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

			`IPADDR="192.168.1.2"
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
DHCP="yes"`,
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

			if got != tc.expected {
				t.Fatalf("writeNetworkCfg() got != expected\n got:\n%s\n expected:\n%s", got, tc.expected)
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

			`IPADDR=""
NETMASK=""
GATEWAY=""
DEVICE=""`,
		},

		{
			"IPv4",
			isoRequest{
				MgmtIPAddress: net.IP{192, 168, 2, 3},
				MgmtIPNetmask: net.IP{255, 255, 255, 255},
				MgmtIPGateway: net.IP{192, 168, 1, 255},
				MgmtInterface: "eth0",
			},

			`IPADDR="192.168.2.3"
NETMASK="255.255.255.255"
GATEWAY="192.168.1.255"
DEVICE="eth0"`,
		},

		{
			"IPv6",
			isoRequest{
				MgmtIPAddress: net.ParseIP("beef::1"),
				MgmtIPNetmask: net.IP{255, 255, 255, 255},
				MgmtIPGateway: net.IP{192, 168, 1, 255},
				MgmtInterface: "eth0",
			},

			`IPV6ADDR="beef::1"
NETMASK="255.255.255.255"
GATEWAY="192.168.1.255"
DEVICE="eth0"`,
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

			if got != tc.expected {
				t.Fatalf("writeMgmtNetworkCfg() got != expected\n got:\n%s\n expected:\n%s", got, tc.expected)
			}
			t.Logf("writeMgmtNetworkCfg():\n%s", got)
		})
	}
}
