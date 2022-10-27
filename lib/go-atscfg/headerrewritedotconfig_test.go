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
	"reflect"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestMakeHeaderRewriteDotConfig(t *testing.T) {
	xmlID := "xml-id"
	fileName := "hdr_rw_" + xmlID + ".config"
	cdnName := "mycdn"
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = &cdnName

	server.HostName = util.StrPtr("my-edge")
	server.ID = util.IntPtr(990)
	server.Status = util.StrPtr(string(tc.CacheStatusReported))
	server.CDNName = &cdnName

	ds := makeGenericDS()
	ds.EdgeHeaderRewrite = util.StrPtr("edgerewrite")
	ds.ID = util.IntPtr(240)
	ds.XMLID = &xmlID
	ds.MaxOriginConnections = util.IntPtr(42)
	ds.MidHeaderRewrite = util.StrPtr("midrewrite")
	ds.CDNName = &cdnName
	dsType := tc.DSTypeHTTPLive
	ds.Type = &dsType
	ds.ServiceCategory = util.StrPtr("servicecategory")

	sv1 := makeGenericServer()
	sv1.HostName = util.StrPtr("my-edge-1")
	sv1.CDNName = &cdnName
	sv1.ID = util.IntPtr(991)
	sv1Status := string(tc.CacheStatusOnline)
	sv1.Status = &sv1Status

	sv2 := makeGenericServer()
	sv2.HostName = util.StrPtr("my-edge-2")
	sv2.CDNName = &cdnName
	sv2.ID = util.IntPtr(992)
	sv2Status := string(tc.CacheStatusOffline)
	sv2.Status = &sv2Status

	servers := []Server{*server, *sv1, *sv2}
	dses := []DeliveryService{*ds}

	dss := makeDSS(servers, dses)

	topologies := []tc.Topology{}
	serverParams := makeHdrRwServerParams()
	cgs := []tc.CacheGroupNullable{}
	serverCaps := map[int]map[ServerCapability]struct{}{}
	dsRequiredCaps := map[int]map[ServerCapability]struct{}{}

	cfg, err := MakeHeaderRewriteDotConfig(fileName, dses, dss, server, servers, cgs, serverParams, serverCaps, dsRequiredCaps, topologies, &HeaderRewriteDotConfigOpts{HdrComment: hdr})

	if err != nil {
		t.Errorf("error expected nil, actual '%v'\n", err)
	}

	txt := cfg.Text

	if !strings.Contains(txt, "edgerewrite") {
		t.Errorf("expected 'edgerewrite' actual '%v'\n", txt)
	}

	if strings.Contains(txt, "midrewrite") {
		t.Errorf("expected no 'midrewrite' actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected origin_max_connections on edge header rewrite that doesn't use the mids, actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "21") { // 21, because max is 42, and there are 2 not-offline mids, so 42/2=21
		t.Errorf("expected origin_max_connections of 21, actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "xml-id|servicecategory") {
		t.Errorf("expected 'xml-id|servicecategory' actual '%v'\n", txt)
	}
}

func TestMakeHeaderRewriteDotConfigNoMaxOriginConnections(t *testing.T) {
	xmlID := "xml-id"
	fileName := "hdr_rw_" + xmlID + ".config"
	cdnName := "mycdn"
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = &cdnName

	server.HostName = util.StrPtr("my-edge")
	server.ID = util.IntPtr(990)
	serverStatus := string(tc.CacheStatusReported)
	server.Status = &serverStatus
	server.CDNName = &cdnName

	ds := makeGenericDS()
	ds.EdgeHeaderRewrite = util.StrPtr("edgerewrite")
	ds.ID = util.IntPtr(240)
	ds.XMLID = &xmlID
	ds.MaxOriginConnections = util.IntPtr(42)
	ds.MidHeaderRewrite = util.StrPtr("midrewrite")
	ds.CDNName = &cdnName
	dsType := tc.DSTypeHTTP
	ds.Type = &dsType
	ds.ServiceCategory = util.StrPtr("servicecategory")

	sv1 := makeGenericServer()
	sv1.HostName = util.StrPtr("my-edge-1")
	sv1.CDNName = &cdnName
	sv1.ID = util.IntPtr(991)
	sv1Status := string(tc.CacheStatusOnline)
	sv1.Status = &sv1Status

	sv2 := makeGenericServer()
	sv2.HostName = util.StrPtr("my-edge-2")
	sv2.CDNName = &cdnName
	sv2.ID = util.IntPtr(992)
	sv2Status := string(tc.CacheStatusOffline)
	sv2.Status = &sv2Status

	servers := []Server{*server, *sv1, *sv2}
	dses := []DeliveryService{*ds}

	dss := makeDSS(servers, dses)

	topologies := []tc.Topology{}
	serverParams := makeHdrRwServerParams()
	cgs := []tc.CacheGroupNullable{}
	serverCaps := map[int]map[ServerCapability]struct{}{}
	dsRequiredCaps := map[int]map[ServerCapability]struct{}{}

	cfg, err := MakeHeaderRewriteDotConfig(fileName, dses, dss, server, servers, cgs, serverParams, serverCaps, dsRequiredCaps, topologies, &HeaderRewriteDotConfigOpts{HdrComment: hdr})

	if err != nil {
		t.Errorf("error expected nil, actual '%v'\n", err)
	}

	txt := cfg.Text

	if strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected no origin_max_connections on DS that uses the mid, actual '%v'\n", txt)
	}
}

func TestGetCachegroupsInSameTopologyTier(t *testing.T) {
	allCachegroups := []tc.CacheGroupNullable{
		{
			Name: util.StrPtr("edge1"),
			Type: util.StrPtr(tc.CacheGroupEdgeTypeName),
		},
		{
			Name: util.StrPtr("edge2"),
			Type: util.StrPtr(tc.CacheGroupEdgeTypeName),
		},
		{
			Name: util.StrPtr("deep1"),
			Type: util.StrPtr(tc.CacheGroupEdgeTypeName),
		},
		{
			Name: util.StrPtr("mid1"),
			Type: util.StrPtr(tc.CacheGroupMidTypeName),
		},
		{
			Name: util.StrPtr("mid2"),
			Type: util.StrPtr(tc.CacheGroupMidTypeName),
		},
		{
			Name: util.StrPtr("org1"),
			Type: util.StrPtr(tc.CacheGroupOriginTypeName),
		},
		{
			Name: util.StrPtr("org2"),
			Type: util.StrPtr(tc.CacheGroupOriginTypeName),
		},
	}
	type testCase struct {
		cachegroup  string
		cachegroups []tc.CacheGroupNullable
		topology    tc.Topology
		expected    map[string]bool
	}
	testCases := []testCase{
		{
			cachegroup:  "edge1",
			cachegroups: allCachegroups,
			topology: tc.Topology{
				Nodes: []tc.TopologyNode{
					{
						// 0
						Cachegroup: "edge1",
						Parents:    []int{3},
					},
					{
						// 1
						Cachegroup: "deep1",
						Parents:    []int{0},
					},
					{
						// 2
						Cachegroup: "edge2",
						Parents:    []int{4},
					},
					{
						// 3
						Cachegroup: "mid1",
						Parents:    []int{},
					},
					{
						// 4
						Cachegroup: "mid2",
						Parents:    []int{},
					},
				},
			},
			expected: map[string]bool{"edge1": true, "edge2": true},
		},
		{
			cachegroup:  "deep1",
			cachegroups: allCachegroups,
			topology: tc.Topology{
				Nodes: []tc.TopologyNode{
					{
						// 0
						Cachegroup: "edge1",
						Parents:    []int{3},
					},
					{
						// 1
						Cachegroup: "deep1",
						Parents:    []int{0},
					},
					{
						// 2
						Cachegroup: "edge2",
						Parents:    []int{4},
					},
					{
						// 3
						Cachegroup: "mid1",
						Parents:    []int{},
					},
					{
						// 4
						Cachegroup: "mid2",
						Parents:    []int{},
					},
				},
			},
			expected: map[string]bool{"deep1": true},
		},
		{
			cachegroup:  "mid1",
			cachegroups: allCachegroups,
			topology: tc.Topology{
				Nodes: []tc.TopologyNode{
					{
						// 0
						Cachegroup: "edge1",
						Parents:    []int{3},
					},
					{
						// 1
						Cachegroup: "deep1",
						Parents:    []int{0},
					},
					{
						// 2
						Cachegroup: "edge2",
						Parents:    []int{4},
					},
					{
						// 3
						Cachegroup: "mid1",
						Parents:    []int{},
					},
					{
						// 4
						Cachegroup: "mid2",
						Parents:    []int{},
					},
				},
			},
			expected: map[string]bool{"mid1": true, "mid2": true},
		},
		{
			cachegroup:  "edge2",
			cachegroups: allCachegroups,
			topology: tc.Topology{
				Nodes: []tc.TopologyNode{
					{
						// 0
						Cachegroup: "edge1",
						Parents:    []int{3},
					},
					{
						// 1
						Cachegroup: "deep1",
						Parents:    []int{0},
					},
					{
						// 2
						Cachegroup: "edge2",
						Parents:    []int{4},
					},
					{
						// 3
						Cachegroup: "mid1",
						Parents:    []int{5},
					},
					{
						// 4
						Cachegroup: "mid2",
						Parents:    []int{6},
					},
					{
						// 5
						Cachegroup: "org1",
						Parents:    []int{},
					},
					{
						// 5
						Cachegroup: "org2",
						Parents:    []int{},
					},
				},
			},
			expected: map[string]bool{"edge1": true, "edge2": true},
		},
	}
	for _, tc := range testCases {
		actual := getCachegroupsInSameTopologyTier(tc.cachegroup, tc.cachegroups, tc.topology)
		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("getting cachegroups in same topology tier -- expected: %v, actual: %v", tc.expected, actual)
		}
	}
}

func TestGetAssignedTierPeers(t *testing.T) {
	allCachegroups := []tc.CacheGroupNullable{
		{
			Name:       util.StrPtr("edge1"),
			ParentName: util.Ptr("mid1"),
			Type:       util.StrPtr(tc.CacheGroupEdgeTypeName),
		},
		{
			Name:       util.StrPtr("edge2"),
			ParentName: util.Ptr("mid2"),
			Type:       util.StrPtr(tc.CacheGroupEdgeTypeName),
		},
		{
			Name:       util.StrPtr("mid1"),
			ParentName: util.Ptr("org1"),
			Type:       util.StrPtr(tc.CacheGroupMidTypeName),
		},
		{
			Name:       util.StrPtr("mid2"),
			ParentName: util.Ptr("org1"),
			Type:       util.StrPtr(tc.CacheGroupMidTypeName),
		},
		{
			Name: util.StrPtr("org1"),
			Type: util.StrPtr(tc.CacheGroupOriginTypeName),
		},
	}

	edges := []Server{
		{
			Cachegroup: util.Ptr("edge1"),
			CDNName:    util.Ptr("mycdn"),
			HostName:   util.Ptr("edgeCache1"),
			ID:         util.Ptr(1),
			Status:     util.Ptr(string(tc.CacheStatusReported)),
		},
		{
			Cachegroup: util.Ptr("edge2"),
			CDNName:    util.Ptr("mycdn"),
			HostName:   util.Ptr("edgeCache2"),
			ID:         util.Ptr(2),
			Status:     util.Ptr(string(tc.CacheStatusReported)),
		},
	}
	mids := []Server{
		{
			Cachegroup: util.Ptr("mid1"),
			CDNName:    util.Ptr("mycdn"),
			HostName:   util.Ptr("midCache1"),
			ID:         util.Ptr(3),
			Status:     util.Ptr(string(tc.CacheStatusReported)),
			Type:       tc.MidTypePrefix,
		},
		{
			Cachegroup: util.Ptr("mid2"),
			CDNName:    util.Ptr("mycdn"),
			HostName:   util.Ptr("midCache2"),
			ID:         util.Ptr(4),
			Status:     util.Ptr(string(tc.CacheStatusReported)),
			Type:       tc.MidTypePrefix,
		},
	}
	allServers := append(edges, mids...)

	topology := tc.Topology{
		Name: "mytopology",
		Nodes: []tc.TopologyNode{
			{
				Cachegroup: "edge1",
				Parents:    []int{2},
			},
			{
				Cachegroup: "edge2",
				Parents:    []int{2},
			},
			{
				Cachegroup: "org1",
			},
		},
	}

	allDeliveryServices := []DeliveryService{{}, {}, {}, {}}
	allDeliveryServices[0].ID = util.Ptr(1)
	allDeliveryServices[0].CDNName = util.Ptr("mycdn")
	allDeliveryServices[1].ID = util.Ptr(2)
	allDeliveryServices[1].Regional = true
	allDeliveryServices[1].CDNName = util.Ptr("mycdn")
	allDeliveryServices[2].Topology = util.Ptr(topology.Name)
	allDeliveryServices[3].Topology = util.Ptr(topology.Name)
	allDeliveryServices[3].Regional = true

	type testCase struct {
		server                 *Server
		ds                     *DeliveryService
		topology               tc.Topology
		deliveryServiceServers []DeliveryServiceServer
		dsRequiredCapabilities map[ServerCapability]struct{}
		servers                []Server
		cacheGroups            []tc.CacheGroupNullable
		serverCapabilities     map[int]map[ServerCapability]struct{}

		expectedHostnames []string
	}
	testCases := []testCase{
		// topology
		{
			ds:          &allDeliveryServices[2],
			server:      &allServers[0],
			topology:    topology,
			cacheGroups: allCachegroups,
			servers:     allServers,

			expectedHostnames: []string{"edgeCache1", "edgeCache2"},
		},
		{
			ds:          &allDeliveryServices[3],
			server:      &allServers[0],
			topology:    topology,
			cacheGroups: allCachegroups,
			servers:     allServers,

			expectedHostnames: []string{"edgeCache1"},
		},
		// mid
		{
			server:  &allServers[2],
			ds:      &allDeliveryServices[0],
			servers: allServers,
			deliveryServiceServers: []DeliveryServiceServer{
				{
					Server:          1,
					DeliveryService: 1,
				},
				{
					Server:          2,
					DeliveryService: 1,
				},
				{
					Server:          3,
					DeliveryService: 1,
				},
				{
					Server:          4,
					DeliveryService: 1,
				},
			},
			cacheGroups: allCachegroups,

			expectedHostnames: []string{"midCache1", "midCache2"},
		},
		{
			server:  &allServers[2],
			ds:      &allDeliveryServices[1],
			servers: allServers,
			deliveryServiceServers: []DeliveryServiceServer{
				{
					Server:          1,
					DeliveryService: 2,
				},
				{
					Server:          2,
					DeliveryService: 2,
				},
				{
					Server:          3,
					DeliveryService: 2,
				},
				{
					Server:          4,
					DeliveryService: 2,
				},
			},
			cacheGroups: allCachegroups,

			expectedHostnames: []string{"midCache1"},
		},
		// edge
		{
			server:  &edges[0],
			ds:      &allDeliveryServices[0],
			servers: edges,
			deliveryServiceServers: []DeliveryServiceServer{
				{
					Server:          1,
					DeliveryService: 1,
				},
				{
					Server:          2,
					DeliveryService: 1,
				},
				{
					Server:          3,
					DeliveryService: 1,
				},
				{
					Server:          4,
					DeliveryService: 1,
				},
			},
			cacheGroups: allCachegroups,

			expectedHostnames: []string{"edgeCache1", "edgeCache2"},
		},
		{
			server:  &edges[0],
			ds:      &allDeliveryServices[1],
			servers: edges,
			deliveryServiceServers: []DeliveryServiceServer{
				{
					Server:          1,
					DeliveryService: 2,
				},
				{
					Server:          2,
					DeliveryService: 2,
				},
				{
					Server:          3,
					DeliveryService: 2,
				},
				{
					Server:          4,
					DeliveryService: 2,
				},
			},
			cacheGroups: allCachegroups,

			expectedHostnames: []string{"edgeCache1"},
		},
	}

	for _, tc := range testCases {
		actualServers, _ := getAssignedTierPeers(tc.server, tc.ds, tc.topology, tc.servers, tc.deliveryServiceServers, tc.cacheGroups, tc.serverCapabilities, tc.dsRequiredCapabilities)
		actualHostnames := []string{}
		for _, as := range actualServers {
			actualHostnames = append(actualHostnames, *as.HostName)
		}
		if !reflect.DeepEqual(tc.expectedHostnames, actualHostnames) {
			t.Errorf("getting servers in same cachegroup tier -- expected: %v, actual: %v", tc.expectedHostnames, actualHostnames)
		}
	}
}

func TestMakeHeaderRewriteMidDotConfig(t *testing.T) {
	cdnName := "mycdn"
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = util.StrPtr(cdnName)
	server.Cachegroup = util.StrPtr("edgeCG")
	server.HostName = util.StrPtr("myserver")
	server.Status = util.StrPtr(string(tc.CacheStatusReported))
	server.Type = string(tc.CacheTypeMid)

	ds := makeGenericDS()
	ds.EdgeHeaderRewrite = util.StrPtr("edgerewrite")
	ds.ID = util.IntPtr(24)
	ds.XMLID = util.StrPtr("ds0")
	ds.MaxOriginConnections = util.IntPtr(42)
	ds.MidHeaderRewrite = util.StrPtr("midrewrite")
	ds.CDNName = &cdnName
	dsType := tc.DSTypeHTTP
	ds.Type = &dsType
	ds.ServiceCategory = util.StrPtr("servicecategory")

	mid0 := makeGenericServer()
	mid0.CDNName = &cdnName
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0Status := string(tc.CacheStatusReported)
	mid0.Status = &mid0Status

	mid1 := makeGenericServer()
	mid1.CDNName = &cdnName
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.HostName = util.StrPtr("mymid1")
	mid1Status := string(tc.CacheStatusOnline)
	mid1.Status = &mid1Status

	mid2 := makeGenericServer()
	mid2.CDNName = &cdnName
	mid2.Cachegroup = util.StrPtr("midCG")
	mid2.HostName = util.StrPtr("mymid2")
	mid2Status := string(tc.CacheStatusOffline)
	mid2.Status = &mid2Status

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	cgs := []tc.CacheGroupNullable{*eCG, *mCG}
	servers := []Server{*server, *mid0, *mid1, *mid2}
	dses := []DeliveryService{*ds}
	dss := makeDSS(servers, dses)

	fileName := "hdr_rw_mid_" + *ds.XMLID + ".config"

	topologies := []tc.Topology{}
	serverParams := makeHdrRwServerParams()
	serverCaps := map[int]map[ServerCapability]struct{}{}
	dsRequiredCaps := map[int]map[ServerCapability]struct{}{}

	cfg, err := MakeHeaderRewriteDotConfig(fileName, dses, dss, server, servers, cgs, serverParams, serverCaps, dsRequiredCaps, topologies, &HeaderRewriteDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Error(err)
	}

	txt := cfg.Text

	if !strings.Contains(txt, "midrewrite") {
		t.Errorf("expected no 'midrewrite' actual '%v'\n", txt)
	}

	if strings.Contains(txt, "edgerewrite") {
		t.Errorf("expected 'edgerewrite' actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected origin_max_connections on edge header rewrite that uses the mids, actual '%v'\n", txt)
	}

	if !strings.Contains(txt, "21") { // 21, because max is 42, and there are 2 not-offline mids, so 42/2=21
		t.Errorf("expected origin_max_connections of 21, actual '%v'\n", txt)
	}
}

func TestMakeHeaderRewriteMidDotConfigNoMaxConns(t *testing.T) {
	cdnName := "mycdn"
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = util.StrPtr(cdnName)
	server.Cachegroup = util.StrPtr("edgeCG")
	server.HostName = util.StrPtr("myserver")
	server.Status = util.StrPtr(string(tc.CacheStatusReported))
	server.Type = string(tc.CacheTypeMid)

	ds := makeGenericDS()
	ds.EdgeHeaderRewrite = util.StrPtr("edgerewrite")
	ds.ID = util.IntPtr(24)
	ds.XMLID = util.StrPtr("ds0")
	ds.MidHeaderRewrite = util.StrPtr("midrewrite")
	ds.CDNName = &cdnName
	dsType := tc.DSTypeHTTP
	ds.Type = &dsType
	ds.ServiceCategory = util.StrPtr("servicecategory")

	mid0 := makeGenericServer()
	mid0.Cachegroup = util.StrPtr("midCG")
	mid0.HostName = util.StrPtr("mymid0")
	mid0Status := string(tc.CacheStatusReported)
	mid0.Status = &mid0Status

	mid1 := makeGenericServer()
	mid1.Cachegroup = util.StrPtr("midCG")
	mid1.HostName = util.StrPtr("mymid1")
	mid1Status := string(tc.CacheStatusOnline)
	mid1.Status = &mid1Status

	mid2 := makeGenericServer()
	mid2.Cachegroup = util.StrPtr("midCG")
	mid2.HostName = util.StrPtr("mymid2")
	mid2Status := string(tc.CacheStatusOffline)
	mid2.Status = &mid2Status

	eCG := &tc.CacheGroupNullable{}
	eCG.Name = server.Cachegroup
	eCG.ID = server.CachegroupID
	eCG.ParentName = mid0.Cachegroup
	eCG.ParentCachegroupID = mid0.CachegroupID
	eCGType := tc.CacheGroupEdgeTypeName
	eCG.Type = &eCGType

	mCG := &tc.CacheGroupNullable{}
	mCG.Name = mid0.Cachegroup
	mCG.ID = mid0.CachegroupID
	mCGType := tc.CacheGroupMidTypeName
	mCG.Type = &mCGType

	cgs := []tc.CacheGroupNullable{*eCG, *mCG}
	servers := []Server{*server, *mid0, *mid1, *mid2}
	dses := []DeliveryService{*ds}
	dss := makeDSS(servers, dses)

	fileName := "hdr_rw_mid_" + *ds.XMLID + ".config"

	topologies := []tc.Topology{}
	serverParams := makeHdrRwServerParams()
	serverCaps := map[int]map[ServerCapability]struct{}{}
	dsRequiredCaps := map[int]map[ServerCapability]struct{}{}

	cfg, err := MakeHeaderRewriteDotConfig(fileName, dses, dss, mid0, servers, cgs, serverParams, serverCaps, dsRequiredCaps, topologies, &HeaderRewriteDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Error(err)
	}

	txt := cfg.Text

	if strings.Contains(txt, "origin_max_connections") {
		t.Errorf("expected no origin_max_connections on edge-only DS, actual '%v'\n", txt)
	}
}

func makeHdrRwServerParams() []tc.Parameter {
	serverParams := []tc.Parameter{
		tc.Parameter{
			Name:       "trafficserver",
			ConfigFile: "package",
			Value:      "7",
			Profiles:   []byte(`["global"]`),
		},
	}
	return serverParams
}
