package atscfg

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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestMakeIPAllowDotConfig(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"
	params := map[string][]string{
		"purge_allow_ip":       []string{"192.168.2.99"},
		ParamCoalesceMaskLenV4: []string{"24"},
		ParamCoalesceNumberV4:  []string{"3"},
		ParamCoalesceMaskLenV6: []string{"48"},
		ParamCoalesceNumberV6:  []string{"4"},
	}

	svs := []tc.ServerNullable{
		*makeIPAllowChild("child0", "192.168.2.1", "2001:DB8:1::1/64"),
		*makeIPAllowChild("child1", "192.168.2.100/30", "2001:DB8:2::1/64"),
		*makeIPAllowChild("child2", "192.168.2.150", ""),
		*makeIPAllowChild("child3", "", "2001:DB8:2::2/64"),
		*makeIPAllowChild("child4", "", "192.168.2.155/32"),
		*makeIPAllowChild("child5", "", "2001:DB8:3::1"),
		*makeIPAllowChild("child6", "", "2001:DB8:2::3"),
		*makeIPAllowChild("child7", "", "2001:DB8:2::4"),
		*makeIPAllowChild("child8", "", "2001:DB8:2::5/64"),
	}

	expecteds := []string{
		"127.0.0.1",
		"::1",
		"0.0.0.0-255.255.255.255",
		"::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"172.16.0.0-172.31.255.255",
		"10.0.0.0-10.255.255.255",
		"2001:db8:3::1",
		"192.168.2.0-192.168.2.255",
		"192.168.2.99",
		"2001:db8:1::-2001:db8:1:0:ffff:ffff:ffff:ffff",
		"2001:db8:2::-2001:db8:2:ffff:ffff:ffff:ffff:ffff",
	}

	cgs := []tc.CacheGroupNullable{
		tc.CacheGroupNullable{
			Name: util.StrPtr("cg0"),
		},
	}

	sv := &tc.ServerNullable{}
	sv.HostName = util.StrPtr("server0")
	sv.Type = string(tc.CacheTypeMid)
	sv.Cachegroup = cgs[0].Name
	svs = append(svs, *sv)
	txt := MakeIPAllowDotConfig(toToolName, toURL, params, sv, svs, cgs)

	lines := strings.Split(txt, "\n")

	if len(lines) == 0 {
		t.Fatalf("expected: lines actual: no lines\n")
	}

	commentLine := lines[0]
	commentLine = strings.TrimSpace(commentLine)
	if !strings.HasPrefix(commentLine, "#") {
		t.Errorf("expected: comment line starting with '#', actual: '%v'\n", commentLine)
	}
	if !strings.Contains(commentLine, toToolName) {
		t.Errorf("expected: comment line containing toolName '%v', actual: '%v'\n", toToolName, commentLine)
	}
	if !strings.Contains(commentLine, toURL) {
		t.Errorf("expected: comment line containing toURL '%v', actual: '%v'\n", toURL, commentLine)
	}

	lines = lines[1:] // remove comment line

	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected %+v actual '%v'\n", expected, txt)
		}
	}
}

func TestMakeIPAllowDotConfigEdge(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"
	params := map[string][]string{
		ParamCoalesceMaskLenV4: []string{"24"},
		ParamCoalesceNumberV4:  []string{"3"},
		ParamCoalesceMaskLenV6: []string{"48"},
		ParamCoalesceNumberV6:  []string{"4"},
	}

	svs := []tc.ServerNullable{
		*makeIPAllowChild("child0", "192.168.2.1", "2001:DB8:1::1/64"),
		*makeIPAllowChild("child1", "192.168.2.100/30", "2001:DB8:2::1/64"),
		*makeIPAllowChild("child2", "192.168.2.150", ""),
		*makeIPAllowChild("child3", "", "2001:DB8:2::2/64"),
		*makeIPAllowChild("child4", "", "192.168.2.155/32"),
		*makeIPAllowChild("child5", "", "2001:DB8:3::1"),
		*makeIPAllowChild("child6", "", "2001:DB8:2::3"),
		*makeIPAllowChild("child7", "", "2001:DB8:2::4"),
		*makeIPAllowChild("child8", "", "2001:DB8:2::5/64"),
	}

	expecteds := []string{
		"127.0.0.1",
		"::1",
		"0.0.0.0-255.255.255.255",
		"::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
	}

	notExpecteds := []string{
		"2001:db8",
		"192.168.2",
	}

	cgs := []tc.CacheGroupNullable{
		tc.CacheGroupNullable{
			Name: util.StrPtr("cg0"),
		},
	}

	sv := &tc.ServerNullable{}
	sv.HostName = util.StrPtr("server0")
	sv.Type = string(tc.CacheTypeEdge)
	sv.Cachegroup = cgs[0].Name
	svs = append(svs, *sv)
	txt := MakeIPAllowDotConfig(toToolName, toURL, params, sv, svs, cgs)

	lines := strings.Split(txt, "\n")

	if len(lines) == 0 {
		t.Fatalf("expected: lines actual: no lines\n")
	}

	commentLine := lines[0]
	commentLine = strings.TrimSpace(commentLine)
	if !strings.HasPrefix(commentLine, "#") {
		t.Errorf("expected: comment line starting with '#', actual: '%v'\n", commentLine)
	}
	if !strings.Contains(commentLine, toToolName) {
		t.Errorf("expected: comment line containing toolName '%v', actual: '%v'\n", toToolName, commentLine)
	}
	if !strings.Contains(commentLine, toURL) {
		t.Errorf("expected: comment line containing toURL '%v', actual: '%v'\n", toURL, commentLine)
	}

	lines = lines[1:] // remove comment line

	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected %+v actual '%v'\n", expected, txt)
		}
	}

	for _, notExpected := range notExpecteds {
		if strings.Contains(txt, notExpected) {
			t.Errorf("expected NOT %+v actual '%v'\n", notExpected, txt)
		}
	}
}

func TestMakeIPAllowDotConfigNonDefaultV6Number(t *testing.T) {
	toToolName := "to0"
	toURL := "trafficops.example.net"
	params := map[string][]string{
		"purge_allow_ip":       []string{"192.168.2.99"},
		ParamCoalesceMaskLenV4: []string{"24"},
		ParamCoalesceNumberV4:  []string{"3"},
		ParamCoalesceMaskLenV6: []string{"48"},
		ParamCoalesceNumberV6:  []string{"100"},
	}

	svs := []tc.ServerNullable{
		*makeIPAllowChild("child0", "192.168.2.1", "2001:DB8:1::1/64"),
		*makeIPAllowChild("child1", "192.168.2.100/30", "2001:DB8:2::1/64"),
		*makeIPAllowChild("child2", "192.168.2.150", ""),
		*makeIPAllowChild("child3", "", "2001:DB8:2::2/64"),
		*makeIPAllowChild("child4", "", "192.168.2.155/32"),
		*makeIPAllowChild("child5", "", "2001:DB8:3::1"),
		*makeIPAllowChild("child6", "", "2001:DB8:2::3"),
		*makeIPAllowChild("child7", "", "2001:DB8:2::4"),
		*makeIPAllowChild("child8", "", "2001:DB8:2::5/64"),
	}

	expecteds := []string{
		"127.0.0.1",
		"::1",
		"0.0.0.0-255.255.255.255",
		"::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"172.16.0.0-172.31.255.255",
		"10.0.0.0-10.255.255.255",
		"2001:db8:3::1",
		"192.168.2.0-192.168.2.255",
		"192.168.2.99",
		"2001:db8:2::3",
		"2001:db8:2::4",
	}

	cgs := []tc.CacheGroupNullable{
		tc.CacheGroupNullable{
			Name: util.StrPtr("cg0"),
		},
	}

	sv := &tc.ServerNullable{}
	sv.HostName = util.StrPtr("server0")
	sv.Type = string(tc.CacheTypeMid)
	sv.Cachegroup = cgs[0].Name
	svs = append(svs, *sv)
	txt := MakeIPAllowDotConfig(toToolName, toURL, params, sv, svs, cgs)

	lines := strings.Split(txt, "\n")

	if len(lines) == 0 {
		t.Fatalf("expected: lines actual: no lines\n")
	}

	commentLine := lines[0]
	commentLine = strings.TrimSpace(commentLine)
	if !strings.HasPrefix(commentLine, "#") {
		t.Errorf("expected: comment line starting with '#', actual: '%v'\n", commentLine)
	}
	if !strings.Contains(commentLine, toToolName) {
		t.Errorf("expected: comment line containing toolName '%v', actual: '%v'\n", toToolName, commentLine)
	}
	if !strings.Contains(commentLine, toURL) {
		t.Errorf("expected: comment line containing toURL '%v', actual: '%v'\n", toURL, commentLine)
	}

	lines = lines[1:] // remove comment line

	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected %+v actual '%v'\n", expected, txt)
		}
	}
}

func makeIPAllowChild(name string, ip string, ip6 string) *tc.ServerNullable {
	sv := &tc.ServerNullable{}
	sv.Cachegroup = util.StrPtr("childcg")
	sv.HostName = util.StrPtr("child0")
	sv.Type = tc.MonitorTypeName
	setIPInfo(sv, "eth0", ip, ip6)
	return sv
}
