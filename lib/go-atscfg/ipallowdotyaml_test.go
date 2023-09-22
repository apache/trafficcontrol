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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestMakeIPAllowDotYAML(t *testing.T) {
	hdr := "myHeaderComment"

	params := makeParamsFromMapArr("serverProfile", IPAllowConfigFileName, map[string][]string{
		"purge_allow_ip":       []string{"192.168.2.99"},
		ParamCoalesceMaskLenV4: []string{"24"},
		ParamCoalesceNumberV4:  []string{"3"},
		ParamCoalesceMaskLenV6: []string{"48"},
		ParamCoalesceNumberV6:  []string{"4"},
	})

	svs := []Server{
		*makeIPAllowChild("child0", "192.168.2.1", "2001:DB8:1::1/64", tc.MonitorTypeName),
		*makeIPAllowChild("child1", "192.168.2.100/30", "2001:DB8:2::1/64", tc.MonitorTypeName),
		*makeIPAllowChild("child2", "192.168.2.150", "", tc.MonitorTypeName),
		*makeIPAllowChild("child3", "", "2001:DB8:2::2/64", tc.MonitorTypeName),
		*makeIPAllowChild("child4", "", "192.168.2.155/32", tc.MonitorTypeName),
		*makeIPAllowChild("child5", "", "2001:DB8:3::1", tc.MonitorTypeName),
		*makeIPAllowChild("child6", "", "2001:DB8:2::3", tc.MonitorTypeName),
		*makeIPAllowChild("child7", "", "2001:DB8:2::4", tc.MonitorTypeName),
		*makeIPAllowChild("child8", "", "2001:DB8:2::5/64", tc.MonitorTypeName),
	}

	expecteds := []string{
		"127.0.0.1",
		"::1",
		"0.0.0.0/0",
		"::/0",
		"172.16.0.0/12",
		"10.0.0.0/8",
		"2001:db8:3::1",
		"192.168.2.0/24",
		"192.168.2.99",
		"2001:db8:1::/64",
		"2001:db8:2::/48",
	}

	cgs := []tc.CacheGroupNullableV5{
		tc.CacheGroupNullableV5{
			Name: util.Ptr("cg0"),
		},
	}

	sv := &Server{}
	sv.HostName = "server0"
	sv.Type = string(tc.CacheTypeMid)
	sv.CacheGroup = *cgs[0].Name
	svs = append(svs, *sv)

	topologies := []tc.TopologyV5{}

	cfg, err := MakeIPAllowDotYAML(params, sv, svs, cgs, topologies, &IPAllowDotYAMLOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	lines := strings.Split(txt, "\n")

	if len(lines) == 0 {
		t.Fatalf("expected: lines actual: no lines\n")
	}

	commentLine := lines[0]
	commentLine = strings.TrimSpace(commentLine)
	if !strings.HasPrefix(commentLine, "#") {
		t.Errorf("expected: comment line starting with '#', actual: '%v'\n", commentLine)
	}
	if !strings.Contains(commentLine, hdr) {
		t.Errorf("expected: comment line containing header comment '%v', actual: '%v'\n", hdr, commentLine)
	}

	lines = lines[1:] // remove comment line

	groups := strings.Split(txt, `apply: in`)

	/* Test that PUSH and PURGE are denied ere the allowance of anything else. */
	{
		for _, group := range groups {
			if strings.Contains(group, "ALL") &&
				strings.Contains(group, "ip_allow") &&
				!(strings.Contains(group, `ip_addrs: ::1`) ||
					strings.Contains(group, `ip_addrs: 127.0`) ||
					strings.Contains(group, `ip_addrs: 192.168.2.99`)) {
				t.Fatalf("Expected the only rules allowing ALL (i.e. PUSH and PURGE) to be localhost and purge_allow_ip, actual: rule '%v'", group)
			}
		}
	}

	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected %+v actual '%v'\n", expected, txt)
		}
	}
}

func TestMakeIPAllowDotYAMLEdge(t *testing.T) {
	hdr := "myHeaderComment"

	params := makeParamsFromMapArr("serverProfile", IPAllowConfigFileName, map[string][]string{
		ParamCoalesceMaskLenV4: []string{"24"},
		ParamCoalesceNumberV4:  []string{"3"},
		ParamCoalesceMaskLenV6: []string{"48"},
		ParamCoalesceNumberV6:  []string{"4"},
	})

	svs := []Server{
		*makeIPAllowChild("child0", "192.168.2.1", "2001:DB8:1::1/64", tc.MonitorTypeName),
		*makeIPAllowChild("child1", "192.168.2.100/30", "2001:DB8:2::1/64", tc.MonitorTypeName),
		*makeIPAllowChild("child2", "192.168.2.150", "", tc.MonitorTypeName),
		*makeIPAllowChild("child3", "", "2001:DB8:2::2/64", tc.MonitorTypeName),
		*makeIPAllowChild("child4", "", "192.168.2.155/32", tc.MonitorTypeName),
		*makeIPAllowChild("child5", "", "2001:DB8:3::1", tc.MonitorTypeName),
		*makeIPAllowChild("child6", "", "2001:DB8:2::3", tc.MonitorTypeName),
		*makeIPAllowChild("child7", "", "2001:DB8:2::4", tc.MonitorTypeName),
		*makeIPAllowChild("child8", "", "2001:DB8:2::5/64", tc.MonitorTypeName),
	}

	expecteds := []string{
		"127.0.0.1",
		"::1",
		"0.0.0.0/0",
		"::/0",
	}

	notExpecteds := []string{
		"2001:db8",
		"192.168.2",
	}

	cgs := []tc.CacheGroupNullableV5{
		tc.CacheGroupNullableV5{
			Name: util.Ptr("cg0"),
		},
	}

	sv := &Server{}
	sv.HostName = "server0"
	sv.Type = string(tc.CacheTypeEdge)
	sv.CacheGroup = *cgs[0].Name
	svs = append(svs, *sv)

	topologies := []tc.TopologyV5{}

	cfg, err := MakeIPAllowDotYAML(params, sv, svs, cgs, topologies, &IPAllowDotYAMLOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	lines := strings.Split(txt, "\n")

	if len(lines) == 0 {
		t.Fatalf("expected: lines actual: no lines\n")
	}

	commentLine := lines[0]
	commentLine = strings.TrimSpace(commentLine)
	if !strings.HasPrefix(commentLine, "#") {
		t.Errorf("expected: comment line starting with '#', actual: '%v'\n", commentLine)
	}
	if !strings.Contains(commentLine, hdr) {
		t.Errorf("expected: comment line containing header comment '%v', actual: '%v'\n", hdr, commentLine)
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

func TestMakeIPAllowDotYAMLNonDefaultV6Number(t *testing.T) {
	hdr := "myHeaderComment"
	params := makeParamsFromMapArr("serverProfile", IPAllowConfigFileName, map[string][]string{
		"purge_allow_ip":       []string{"192.168.2.99"},
		ParamCoalesceMaskLenV4: []string{"24"},
		ParamCoalesceNumberV4:  []string{"3"},
		ParamCoalesceMaskLenV6: []string{"48"},
		ParamCoalesceNumberV6:  []string{"100"},
	})

	svs := []Server{
		*makeIPAllowChild("child0", "192.168.2.1", "2001:DB8:1::1/64", tc.MonitorTypeName),
		*makeIPAllowChild("child1", "192.168.2.100/30", "2001:DB8:2::1/64", tc.MonitorTypeName),
		*makeIPAllowChild("child2", "192.168.2.150", "", tc.MonitorTypeName),
		*makeIPAllowChild("child3", "", "2001:DB8:2::2/64", tc.MonitorTypeName),
		*makeIPAllowChild("child4", "", "192.168.2.155/32", tc.MonitorTypeName),
		*makeIPAllowChild("child5", "", "2001:DB8:3::1", tc.MonitorTypeName),
		*makeIPAllowChild("child6", "", "2001:DB8:2::3", tc.MonitorTypeName),
		*makeIPAllowChild("child7", "", "2001:DB8:2::4", tc.MonitorTypeName),
		*makeIPAllowChild("child8", "", "2001:DB8:2::5/64", tc.MonitorTypeName),
	}

	expecteds := []string{
		"127.0.0.1",
		"::1",
		"0.0.0.0/0",
		"::/0",
		"172.16.0.0/12",
		"10.0.0.0/8",
		"2001:db8:3::1",
		"192.168.2.0/24",
		"192.168.2.99",
		"2001:db8:2::3",
		"2001:db8:2::4",
	}

	cgs := []tc.CacheGroupNullableV5{
		tc.CacheGroupNullableV5{
			Name: util.Ptr("cg0"),
		},
	}

	sv := &Server{}
	sv.HostName = "server0"
	sv.Type = string(tc.CacheTypeMid)
	sv.CacheGroup = *cgs[0].Name
	svs = append(svs, *sv)

	topologies := []tc.TopologyV5{}

	cfg, err := MakeIPAllowDotYAML(params, sv, svs, cgs, topologies, &IPAllowDotYAMLOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	lines := strings.Split(txt, "\n")

	if len(lines) == 0 {
		t.Fatalf("expected: lines actual: no lines\n")
	}

	commentLine := lines[0]
	commentLine = strings.TrimSpace(commentLine)
	if !strings.HasPrefix(commentLine, "#") {
		t.Errorf("expected: comment line starting with '#', actual: '%v'\n", commentLine)
	}
	if !strings.Contains(commentLine, hdr) {
		t.Errorf("expected: comment line containing header comment '%v', actual: '%v'\n", hdr, commentLine)
	}

	lines = lines[1:] // remove comment line

	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected %+v actual '%v'\n", expected, txt)
		}
	}
}

func TestMakeIPAllowDotYAMLTopologies(t *testing.T) {
	hdr := "myHeaderComment"

	params := makeParamsFromMapArr("serverProfile", IPAllowConfigFileName, map[string][]string{
		"purge_allow_ip":       []string{"192.168.2.99"},
		ParamCoalesceMaskLenV4: []string{"24"},
		ParamCoalesceNumberV4:  []string{"3"},
		ParamCoalesceMaskLenV6: []string{"48"},
		ParamCoalesceNumberV6:  []string{"4"},
	})

	// make children all MID types, because MIDs would never normally be parented to MIDs with pre-topologies
	svs := []Server{
		*makeIPAllowChild("child0", "192.168.2.1", "2001:DB8:1::1/64", tc.MidTypePrefix),
		*makeIPAllowChild("child1", "192.168.2.100/30", "2001:DB8:2::1/64", tc.MidTypePrefix),
		*makeIPAllowChild("child2", "192.168.2.150", "", tc.MidTypePrefix),
		*makeIPAllowChild("child3", "", "2001:DB8:2::2/64", tc.MidTypePrefix),
		*makeIPAllowChild("child4", "", "192.168.2.155/32", tc.MidTypePrefix),
		*makeIPAllowChild("child5", "", "2001:DB8:3::1", tc.MidTypePrefix),
		*makeIPAllowChild("child6", "", "2001:DB8:2::3", tc.MidTypePrefix),
		*makeIPAllowChild("child7", "", "2001:DB8:2::4", tc.MidTypePrefix),
		*makeIPAllowChild("child8", "", "2001:DB8:2::5/64", tc.MidTypePrefix),
	}

	expecteds := []string{
		"127.0.0.1",
		"::1",
		"0.0.0.0/0",
		"::/0",
		"172.16.0.0/12",
		"10.0.0.0/8",
		"2001:db8:3::1",
		"192.168.2.0/24",
		"192.168.2.99",
		"2001:db8:1::/64",
		"2001:db8:2::/48",
	}

	cgs := []tc.CacheGroupNullableV5{
		tc.CacheGroupNullableV5{
			Name: util.Ptr("midcg"),
		},
		tc.CacheGroupNullableV5{
			Name: util.Ptr("midcg2"),
		},
		tc.CacheGroupNullableV5{
			Name: util.Ptr("childcg"),
		},
	}

	topologies := []tc.TopologyV5{
		tc.TopologyV5{
			Name: "t0",
			Nodes: []tc.TopologyNodeV5{
				tc.TopologyNodeV5{
					Cachegroup: "childcg",
					Parents:    []int{1, 2},
				},
				tc.TopologyNodeV5{
					Cachegroup: "midcg",
				},
				tc.TopologyNodeV5{
					Cachegroup: "midcg2",
				},
			},
		},
	}

	sv := &Server{}
	sv.HostName = "server0"
	sv.Type = string(tc.CacheTypeMid)
	sv.CacheGroup = *cgs[1].Name
	svs = append(svs, *sv)

	//	topologies := []tc.Topology{}

	cfg, err := MakeIPAllowDotYAML(params, sv, svs, cgs, topologies, &IPAllowDotYAMLOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	lines := strings.Split(txt, "\n")

	if len(lines) == 0 {
		t.Fatalf("expected: lines actual: no lines\n")
	}

	commentLine := lines[0]
	commentLine = strings.TrimSpace(commentLine)
	if !strings.HasPrefix(commentLine, "#") {
		t.Errorf("expected: comment line starting with '#', actual: '%v'\n", commentLine)
	}
	if !strings.Contains(commentLine, hdr) {
		t.Errorf("expected: comment line containing header comment '%v', actual: '%v'\n", hdr, commentLine)
	}

	lines = lines[1:] // remove comment line

	groups := strings.Split(txt, `apply: in`)

	/* Test that PUSH and PURGE are denied ere the allowance of anything else. */
	{
		for _, group := range groups {
			if strings.Contains(group, "ALL") && strings.Contains(group, "ip_allow") && !(strings.Contains(group, `ip_addrs: ::1`) ||
				strings.Contains(group, `ip_addrs: 127.0`) ||
				strings.Contains(group, `ip_addrs: 192.168.2.99`)) {
				t.Fatalf("Expected the only rules allowing ALL (i.e. PUSH and PURGE) to be localhost and purge_allow_ip, actual: rule '%v'", group)
			}
		}
	}

	for _, expected := range expecteds {
		if !strings.Contains(txt, expected) {
			t.Errorf("expected %+v actual '%v'\n", expected, txt)
		}
	}
}
