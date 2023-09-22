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

func TestMakeHostingDotConfig(t *testing.T) {
	cdnName := "cdn0"

	server := makeGenericServer()
	server.HostName = "server0"
	server.CDN = cdnName
	server.Profiles = []string{"serverprofile"}
	hdr := "myHeaderComment"

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       ParamRAMDrivePrefix,
			ConfigFile: HostingConfigParamConfigFile,
			Value:      "ParamRAMDrivePrefix-shouldnotappearinconfig",
			Profiles:   []byte(`["` + server.Profiles[0] + `"]`),
		},
		tc.ParameterV5{
			Name:       ParamDrivePrefix,
			ConfigFile: HostingConfigParamConfigFile,
			Value:      "ParamDrivePrefix-shouldnotappearinconfig",
			Profiles:   []byte(`["` + server.Profiles[0] + `"]`),
		},
		tc.ParameterV5{
			Name:       "somethingelse",
			ConfigFile: HostingConfigParamConfigFile,
			Value:      "somethingelse-shouldnotappearinconfig",
			Profiles:   []byte(`["` + server.Profiles[0] + `"]`),
		},
	}

	origins := []string{
		"https://origin0.example.net",
		"http://origin1.example.net",
		"http://origin2.example.net/path0",
		"origin3.example.net/",
		"https://origin4.example.net/",
		"http://origin5.example.net/",
	}
	dses := []DeliveryService{}
	for _, origin := range origins {
		ds := makeGenericDS()
		ds.CDNName = &cdnName
		ds.OrgServerFQDN = util.Ptr(origin)
		dses = append(dses, *ds)
	}

	servers := []Server{*server}
	dss := makeDSS(servers, dses)

	cfg, err := MakeHostingDotConfig(server, servers, serverParams, dses, dss, nil, &HostingDotConfigOpts{HdrComment: hdr})
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
	cdnName := "cdn0"

	server := makeGenericServer()
	server.HostName = "server0"
	server.CacheGroup = "edgeCG"
	server.CDN = cdnName
	server.CDNID = 400
	server.ID = 899
	server.Profiles = []string{"serverprofile"}
	hdr := "myHeaderComment"

	serverParams := []tc.ParameterV5{
		tc.ParameterV5{
			Name:       ParamRAMDrivePrefix,
			ConfigFile: HostingConfigParamConfigFile,
			Value:      "ParamRAMDrivePrefix-shouldnotappearinconfig",
			Profiles:   []byte(`["` + server.Profiles[0] + `"]`),
		},
		tc.ParameterV5{
			Name:       ParamDrivePrefix,
			ConfigFile: HostingConfigParamConfigFile,
			Value:      "ParamDrivePrefix-shouldnotappearinconfig",
			Profiles:   []byte(`["` + server.Profiles[0] + `"]`),
		},
		tc.ParameterV5{
			Name:       "somethingelse",
			ConfigFile: HostingConfigParamConfigFile,
			Value:      "somethingelse-shouldnotappearinconfig",
			Profiles:   []byte(`["` + server.Profiles[0] + `"]`),
		},
	}

	dsTopology := makeGenericDS()
	dsTopology.OrgServerFQDN = util.Ptr("https://origin0.example.net")
	dsTopology.XMLID = "ds-topology"
	dsTopology.CDNID = 400
	dsTopology.ID = util.Ptr(900)
	dsTopology.Topology = util.Ptr("t0")
	dsTopology.Active = tc.DSActiveStateActive
	dsType := "HTTP_LIVE"
	dsTopology.Type = &dsType

	dsTopologyWithoutServer := makeGenericDS()
	dsTopologyWithoutServer.ID = util.Ptr(901)
	dsTopologyWithoutServer.OrgServerFQDN = util.Ptr("https://origin1.example.net")
	dsTopologyWithoutServer.XMLID = "ds-topology-without-server"
	dsTopologyWithoutServer.CDNID = 400
	dsTopologyWithoutServer.Topology = util.Ptr("t1")
	dsTopologyWithoutServer.Active = tc.DSActiveStateActive
	dsType2 := "HTTP"
	dsTopologyWithoutServer.Type = &dsType2

	dses := []DeliveryService{*dsTopology, *dsTopologyWithoutServer}

	topologies := []tc.TopologyV5{
		tc.TopologyV5{
			Name: "t0",
			Nodes: []tc.TopologyNodeV5{
				tc.TopologyNodeV5{
					Cachegroup: "edgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNodeV5{
					Cachegroup: "midCG",
				},
			},
		},
		tc.TopologyV5{
			Name: "t1",
			Nodes: []tc.TopologyNodeV5{
				tc.TopologyNodeV5{
					Cachegroup: "otherEdgeCG",
					Parents:    []int{1},
				},
				tc.TopologyNodeV5{
					Cachegroup: "midCG",
				},
			},
		},
	}

	servers := []Server{*server}
	dss := makeDSS(servers, dses)

	cfg, err := MakeHostingDotConfig(server, servers, serverParams, dses, dss, topologies, &HostingDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

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
