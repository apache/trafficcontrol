package util

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package rfc contains functions implementing RFC 7234, 2616, and other RFCs.
// When changing functions, be sure they still conform to the corresponding RFC.
// When adding symbols, document the RFC and section they correspond to.

import (
	"fmt"
	"net"
	"testing"
)

func TestCoalesceIPs(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.2"),
		net.ParseIP("192.168.1.3"),
		net.ParseIP("192.168.2.1"),
		net.ParseIP("192.168.2.2"),
		net.ParseIP("192.168.2.3"),
		net.ParseIP("192.168.2.4"),
	}

	nets := CoalesceIPs(ips, 2, 24)

	for _, ipnet := range nets {
		if ipnet.String() != "192.168.1.0/24" && ipnet.String() != "192.168.2.0/24" {
			t.Errorf("expected '192.168.1.0/24' and '192.168.2.0/24', actual: %+v\n", ipnet)
		}
	}

	nets = CoalesceIPs(nil, 0, 0)
	if nets != nil {
		t.Errorf("expected nil output when passing in no IPs, got: %+v", nets)
	}
	nets = CoalesceIPs([]net.IP{}, 0, 0)
	if nets != nil {
		t.Errorf("expected nil output when passing in no IPs, got: %+v", nets)
	}
}

func TestCoalesceIPsSmallerThanNum(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.2"),
		net.ParseIP("192.168.1.3"),
		net.ParseIP("192.168.2.1"),
		net.ParseIP("192.168.2.2"),
		net.ParseIP("192.168.2.3"),
		net.ParseIP("192.168.2.4"),
	}

	nets := CoalesceIPs(ips, 4, 24)

	expecteds := map[string]struct{}{
		"192.168.1.1/32": {},
		"192.168.1.2/32": {},
		"192.168.1.3/32": {},
		"192.168.2.0/24": {},
	}

	for _, ipnet := range nets {
		if _, ok := expecteds[ipnet.String()]; !ok {
			t.Errorf("expected: %+v actual: %+v\n", expecteds, ipnet)
		}
		delete(expecteds, ipnet.String())
	}
}

func TestCoalesceIPsV6(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::2"),
		net.ParseIP("2001:db8::3"),
		net.ParseIP("2001:db8::4:1"),
		net.ParseIP("2001:db8::4:2"),
		net.ParseIP("2001:db8::4:3"),
	}

	nets := CoalesceIPs(ips, 3, 112)

	expecteds := map[string]struct{}{
		"2001:db8::/112":    {},
		"2001:db8::4:0/112": {},
	}

	for _, ipnet := range nets {
		if _, ok := expecteds[ipnet.String()]; !ok {
			t.Errorf("expected: %+v actual: %+v\n", expecteds, ipnet)
		}
		delete(expecteds, ipnet.String())
	}
}

func TestCoalesceIPsV6SmallerThanNum(t *testing.T) {
	ips := []net.IP{
		net.ParseIP("2001:db8::1"),
		net.ParseIP("2001:db8::2"),
		net.ParseIP("2001:db8::3"),
		net.ParseIP("2001:db8::4:1"),
		net.ParseIP("2001:db8::4:2"),
		net.ParseIP("2001:db8::4:3"),
		net.ParseIP("2001:db8::4:4"),
	}

	nets := CoalesceIPs(ips, 4, 112)

	expecteds := map[string]struct{}{
		"2001:db8::1/128":   {},
		"2001:db8::2/128":   {},
		"2001:db8::3/128":   {},
		"2001:db8::4:0/112": {},
	}

	for _, ipnet := range nets {
		if _, ok := expecteds[ipnet.String()]; !ok {
			t.Errorf("expected: %+v actual: %+v\n", expecteds, ipnet)
		}
		delete(expecteds, ipnet.String())
	}
}

func TestRangeStr(t *testing.T) {
	inputExpecteds := map[string]string{
		"192.168.1.0/24":     "192.168.1.0-192.168.1.255",
		"192.168.1.0/16":     "192.168.0.0-192.168.255.255",
		"192.168.1.42/32":    "192.168.1.42",
		"2001:db8::4:42/128": "2001:db8::4:42",
		"2001:db8::4:0/112":  "2001:db8::4:0-2001:db8::4:ffff",
	}
	for input, expected := range inputExpecteds {
		_, ipn, err := net.ParseCIDR(input)
		if err != nil {
			t.Fatal(err.Error())
		}

		//	t.Errorf("ipn: " + ipn.String())

		actual := RangeStr(ipn)
		if expected != actual {
			t.Errorf("expected: '" + expected + "' actual '" + actual + "'")
		}
	}
}

func TestFirstIP(t *testing.T) {
	inputExpecteds := map[string]string{
		"192.168.1.0/24":     "192.168.1.0",
		"192.168.1.0/16":     "192.168.0.0",
		"192.168.1.42/32":    "192.168.1.42",
		"2001:db8::4:42/128": "2001:db8::4:42",
		"2001:db8::4:0/112":  "2001:db8::4:0",
	}
	for input, expected := range inputExpecteds {
		_, ipn, err := net.ParseCIDR(input)
		if err != nil {
			t.Fatal(err.Error())
		}

		//	t.Errorf("ipn: " + ipn.String())

		actual := FirstIP(ipn).String()
		if expected != actual {
			t.Errorf("expected: '" + expected + "' actual '" + actual + "'")
		}
	}
}

func TestLastIP(t *testing.T) {
	inputExpecteds := map[string]string{
		"192.168.1.0/24":     "192.168.1.255",
		"192.168.1.0/16":     "192.168.255.255",
		"192.168.1.42/32":    "192.168.1.42",
		"2001:db8::4:42/128": "2001:db8::4:42",
		"2001:db8::4:0/112":  "2001:db8::4:ffff",
	}
	for input, expected := range inputExpecteds {
		_, ipn, err := net.ParseCIDR(input)
		if err != nil {
			t.Fatal(err.Error())
		}

		//	t.Errorf("ipn: " + ipn.String())

		actual := LastIP(ipn).String()
		if expected != actual {
			t.Errorf("expected: '" + expected + "' actual '" + actual + "'")
		}
	}
}

func TestIP4ToNum(t *testing.T) {
	var tests = []struct {
		ip     string
		number uint32
	}{
		{"127.0.0.1", uint32(2130706433)},
		{"127.0.0.4", uint32(2130706436)},
		{"127.255.255.255", uint32(2147483647)},
	}
	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			n, err := IP4ToNum(tt.ip)
			if err != nil {
				t.Errorf("unexpected error: %v", err)

			}
			if n != tt.number {
				t.Errorf("got %v, want %v", n, tt.number)
			}
		})
	}

	_, err := IP4ToNum("invalid")
	if err == nil {
		t.Error("Expected an error passing invalid input, but didn't get one")
	}
	_, err = IP4ToNum("127.0.0.invalid")
	if err == nil {
		t.Error("Expected an error passing invalid input, but didn't get one")
	}
}

func TestIP4InRange(t *testing.T) {
	var tests = []struct {
		ip      string
		ipRange string
		inRange bool
	}{
		{"111.0.0.1", "127.0.0.0-127.255.255.255", false},
		{"128.0.0.1", "127.0.0.0-127.255.255.255", false},
		{"127.0.0.1", "127.0.0.0-127.255.255.255", true},
		{"127.0.0.1", "127.0.0.1", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v in range %v", tt.ip, tt.inRange), func(t *testing.T) {
			exists, err := IP4InRange(tt.ip, tt.ipRange)
			if err != nil {
				t.Errorf("unexpected error: %v", err)

			}
			if exists != tt.inRange {
				t.Errorf("got %v, want %v", exists, tt.inRange)
			}
		})
	}

	_, err := IP4InRange("127.0.0.1", "127.0.0.0-127.0.0.2-127.0.1.2")
	if err == nil {
		t.Error("Expected an error passing an invalid range, but didn't get one")
	}
	_, err = IP4InRange("invalid", "127.0.0.0-127.0.0.2")
	if err == nil {
		t.Error("Expected an error passing an invalid ip, but didn't get one")
	}
	_, err = IP4InRange("127.0.0.1", "invalid-127.0.0.2")
	if err == nil {
		t.Error("Expected an error passing an invalid range start, but didn't get one")
	}
	_, err = IP4InRange("127.0.0.1", "127.0.0.0-invalid")
	if err == nil {
		t.Error("Expected an error passing an invalid range end, but didn't get one")
	}

}
