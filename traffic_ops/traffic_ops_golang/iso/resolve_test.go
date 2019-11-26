package iso

import (
	"strings"
	"testing"
)

func TestParseResolve(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"empty",
			"",
			nil,
		},
		{
			"comments",
			`# nameserver 127.0.0.1
			 # nameserver 192.168.1.1`,
			nil,
		},
		{
			"invalid IPs",
			`# /etc/resolv.conf
			nameserver not.an.ip
			nameserver
			nameserver 1
			`,
			nil,
		},
		{
			"single",
			`# /etc/resolv.conf
			nameserver 127.0.0.1
			`,
			[]string{"127.0.0.1"},
		},
		{
			"multi",
			`# /etc/resolv.conf
			# IPv4
			nameserver   127.0.0.1
			nameserver 192.168.1.10
			# IPv6
			nameserver   beef::1
			nameserver   ::0
			`,
			[]string{"127.0.0.1", "192.168.1.10", "beef::1", "::0"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseResolve(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("parseResolve() err = %v", err)
			}

			if lGot, lExpected := len(got), len(tc.expected); lGot != lExpected {
				t.Fatalf("got %d nameservers; expected %d", lGot, lExpected)
			}

			for i, expectedNS := range tc.expected {
				if gotNS := got[i]; gotNS != expectedNS {
					t.Errorf("got nameserver[%d] = %q; expected %q", i, gotNS, expectedNS)
				}
			}

			t.Logf("parseResolve(): %v", got)
		})
	}
}
