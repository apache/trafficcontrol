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

func TestMakeHostingDotConfig(t *testing.T) {
	server := &tc.Server{HostName: "server0"}
	toToolName := "to0"
	toURL := "trafficops.example.net"
	params := map[string]string{
		ParamRAMDrivePrefix: "ParamRAMDrivePrefix-shouldnotappearinconfig",
		ParamDrivePrefix:    "ParamDrivePrefix-shouldnotappearinconfig",
		"somethingelse":     "somethingelse-shouldnotappearinconfig",
	}
	origins := []string{
		"https://origin0.example.net",
		"http://origin1.example.net",
		"http://origin2.example.net/path0",
		"origin3.example.net/",
		"https://origin4.example.net/",
		"http://origin5.example.net/",
	}
	dses := []tc.DeliveryServiceNullableV30{}
	for _, origin := range origins {
		ds := tc.DeliveryServiceNullableV30{}
		ds.OrgServerFQDN = util.StrPtr(origin)
		dses = append(dses, ds)
	}

	txt := MakeHostingDotConfig(server, toToolName, toURL, params, dses, nil)

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

	originFQDNs := getFQDNs(origins)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strArrContainsSubstr(originFQDNs, line) {
			t.Errorf("expected %+v actual '%v'\n", originFQDNs, line)
		}
		originFQDNs = strArrRemoveSubstr(originFQDNs, line)
	}

	if len(originFQDNs) > 0 {
		t.Errorf("expected %+v actual %v\n", originFQDNs, "missing")
	}
}

func TestMakeHostingDotConfigTopologiesIgnoreDSS(t *testing.T) {
	server := &tc.Server{HostName: "server0", Cachegroup: "edgeCG"}
	toToolName := "to0"
	toURL := "trafficops.example.net"
	params := map[string]string{
		ParamRAMDrivePrefix: "ParamRAMDrivePrefix-shouldnotappearinconfig",
		ParamDrivePrefix:    "ParamDrivePrefix-shouldnotappearinconfig",
		"somethingelse":     "somethingelse-shouldnotappearinconfig",
	}

	dsTopology := tc.DeliveryServiceNullableV30{}
	dsTopology.OrgServerFQDN = util.StrPtr("https://origin0.example.net")
	dsTopology.XMLID = util.StrPtr("ds-topology")
	dsTopology.Topology = util.StrPtr("t0")
	dsTopology.Active = util.BoolPtr(true)

	dsTopologyWithoutServer := tc.DeliveryServiceNullableV30{}
	dsTopologyWithoutServer.OrgServerFQDN = util.StrPtr("https://origin1.example.net")
	dsTopologyWithoutServer.XMLID = util.StrPtr("ds-topology-without-server")
	dsTopologyWithoutServer.Topology = util.StrPtr("t1")
	dsTopologyWithoutServer.Active = util.BoolPtr(true)

	dses := []tc.DeliveryServiceNullableV30{dsTopology, dsTopologyWithoutServer}

	topologies := []tc.Topology{
		tc.Topology{
			Name: "t0",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
		tc.Topology{
			Name: "t1",
			Nodes: []tc.TopologyNode{
				tc.TopologyNode{
					Cachegroup: "otherEdgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNode{
					Cachegroup: "midCG",
				},
			},
		},
	}

	txt := MakeHostingDotConfig(server, toToolName, toURL, params, dses, topologies)

	if !strings.Contains(txt, "origin0") {
		t.Errorf("expected origin0 in topology, actual %v\n", txt)
	}
	if strings.Contains(txt, "origin1") {
		t.Errorf("expected no origin1 not in topology, actual %v\n", txt)
	}
}

func strArrContainsSubstr(arr []string, substr string) bool {
	for _, as := range arr {
		if strings.Contains(as, substr) {
			return true
		}
	}
	return false
}

func strArrRemoveSubstr(arr []string, substr string) []string {
	// this is terribly inefficient, but it's just for testing, so it doesn't matter
	newArr := []string{}
	for _, as := range arr {
		if strings.Contains(as, substr) {
			continue
		}
		newArr = append(newArr, as)
	}
	return newArr
}

func getFQDNs(origins []string) []string {
	newOrigins := []string{}
	for _, origin := range origins {
		origin = strings.TrimLeft(origin, "http://")
		origin = strings.TrimLeft(origin, "https://")
		origin = strings.TrimRight(origin, "/")
	}
	return newOrigins
}
