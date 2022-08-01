package v3

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
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

type topologyTestCase struct {
	testCaseDescription string
	tc.Topology
}

func TestTopologies(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, CDNs, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, ServerCapabilities, ServerServerCapabilities, Topologies, Tenants, DeliveryServices, DeliveryServicesRequiredCapabilities}, func() {
		GetTestTopologies(t)
		UpdateTestTopologies(t)
		ValidationTestTopologies(t)
		UpdateValidateTopologyORGServerCacheGroup(t)
		EdgeParentOfEdgeSucceedsWithWarning(t)
	})
}

func CreateTestTopologies(t *testing.T) {
	var (
		postResponse *tc.TopologyResponse
		err          error
	)
	for _, topology := range testData.Topologies {
		if postResponse, _, err = TOSession.CreateTopology(topology); err != nil {
			t.Fatalf("could not CREATE topology: %v", err)
		}
		postResponse.Response.LastUpdated = nil
		if !reflect.DeepEqual(topology, postResponse.Response) {
			t.Fatalf("Topology in response should be the same as the one POSTed. expected: %v\nactual: %v", topology, postResponse.Response)
		}
	}
}

func GetTestTopologies(t *testing.T) {
	if len(testData.Topologies) < 1 {
		t.Fatalf("test data has no topologies, can't test")
	}
	topos, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("expected GET error to be nil, actual: %v", err)
	}
	if len(topos) != len(testData.Topologies) {
		t.Errorf("expected topologies GET to return %v topologies, actual %v", len(testData.Topologies), len(topos))
	}
}

func EdgeParentOfEdgeSucceedsWithWarning(t *testing.T) {
	testCase := topologyTestCase{testCaseDescription: "an edge parenting a mid", Topology: tc.Topology{
		Name:        "edge-parent-of-edge",
		Description: "An edge is a parent, which is technically valid, but we will warn the user in case it was a mistake",
		Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1}},
			{Cachegroup: "cachegroup2", Parents: []int{}},
		}}}
	response, _, err := TOSession.CreateTopology(testCase.Topology)
	if err != nil {
		t.Fatalf("expected POST with %v to succeed, actual: nil", testCase.testCaseDescription)
	}
	containsWarning := false
	for _, alert := range response.Alerts.Alerts {
		if alert.Level == "warning" {
			containsWarning = true
		}
	}
	if !containsWarning {
		t.Fatalf("expected a warning-level alert message in the response, actual: %v", response.Alerts)
	}
	delResp, _, err := TOSession.DeleteTopology(testCase.Topology.Name)
	if err != nil {
		t.Fatalf("cannot DELETE topology: %v - %v", err, delResp)
	}
}

func ValidationTestTopologies(t *testing.T) {
	invalidTopologyTestCases := []topologyTestCase{
		{testCaseDescription: "no nodes", Topology: tc.Topology{Name: "empty-top", Description: "Invalid because there are no nodes", Nodes: []tc.TopologyNode{}}},
		{testCaseDescription: "a node listing itself as a parent", Topology: tc.Topology{Name: "self-parent", Description: "Invalid because a node lists itself as a parent", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1}},
			{Cachegroup: "parentCachegroup", Parents: []int{1}},
		}}},
		{testCaseDescription: "duplicate parents", Topology: tc.Topology{}},
		{testCaseDescription: "too many parents", Topology: tc.Topology{Name: "duplicate-parents", Description: "Invalid because a node lists the same parent twice", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1, 1}},
			{Cachegroup: "parentCachegroup", Parents: []int{}},
		}}},
		{testCaseDescription: "too many parents", Topology: tc.Topology{Name: "too-many-parents", Description: "Invalid because a node has more than 2 parents", Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{}},
			{Cachegroup: "parentCachegroup2", Parents: []int{}},
			{Cachegroup: "cachegroup1", Parents: []int{0, 1, 2}},
		}}},
		{testCaseDescription: "an edge parenting a mid", Topology: tc.Topology{Name: "edge-parent-of-mid", Description: "Invalid because an edge is a parent of a mid", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1}},
			{Cachegroup: "parentCachegroup", Parents: []int{2}},
			{Cachegroup: "cachegroup2", Parents: []int{}},
		}}},
		{testCaseDescription: "a leaf mid", Topology: tc.Topology{Name: "leaf-mid", Description: "Invalid because a mid is a leaf node", Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{1}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{}},
		}}},
		{testCaseDescription: "cyclical nodes", Topology: tc.Topology{Name: "cyclical-nodes", Description: "Invalid because it contains cycles", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1, 2}},
			{Cachegroup: "parentCachegroup", Parents: []int{2}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{1}},
		}}},
		{testCaseDescription: "a cycle across topologies", Topology: tc.Topology{Name: "cycle-with-4-tier-topology", Description: `Invalid because it contains a cycle when combined with the "4-tiers" topology`, Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{1}},
			{Cachegroup: "parentCachegroup2", Parents: []int{}},
			{Cachegroup: "cachegroup1", Parents: []int{0}},
		}}},
		{testCaseDescription: "a cycle across cache groups", Topology: tc.Topology{Name: "cycle-with-non-topology-cachegroups", Description: "Invalid because it contains a cycle when combined with a topology constructed from cache group parentage", Nodes: []tc.TopologyNode{
			{Cachegroup: "edge-parent1", Parents: []int{1}},
			{Cachegroup: "has-edge-parent1", Parents: []int{}},
		}}},
		{testCaseDescription: "a nonexistent cache group", Topology: tc.Topology{Name: "nonexistent-cg", Description: "Invalid because it references a cache group that does not exist", Nodes: []tc.TopologyNode{
			{Cachegroup: "legitcachegroup", Parents: []int{}},
		}}},
		{testCaseDescription: "an out-of-bounds parent index", Topology: tc.Topology{Name: "oob-parent", Description: "Invalid because it contains an out-of-bounds parent", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{7}},
		}}},
		{testCaseDescription: "a cachegroup containing no servers", Topology: tc.Topology{Name: "empty-cg", Description: `Invalid because it contains a cachegroup, fallback3, that contains no servers`, Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{}},
			{Cachegroup: "parentCachegroup2", Parents: []int{}},
			{Cachegroup: "fallback3", Parents: []int{0, 1}},
		}}},
	}
	var statusCode int
	for _, testCase := range invalidTopologyTestCases {
		_, reqInf, err := TOSession.CreateTopology(testCase.Topology)
		if err == nil {
			t.Fatalf("expected POST with %v to return an error, actual: nil", testCase.testCaseDescription)
		}
		statusCode = reqInf.StatusCode
		if statusCode < 400 || statusCode >= 500 {
			t.Fatalf("Expected a 400-level status code for topology %s but got %d", testCase.Topology.Name, statusCode)
		}
	}
}

func updateSingleTopology(topology tc.Topology) error {
	updateResponse, _, err := TOSession.UpdateTopology(topology.Name, topology)
	if err != nil {
		return fmt.Errorf("cannot PUT topology: %v - %v", err, updateResponse)
	}
	updateResponse.Response.LastUpdated = nil
	if !reflect.DeepEqual(topology, updateResponse.Response) {
		return fmt.Errorf("Topologies should be equal after updating. expected: %v\nactual: %v", topology, updateResponse.Response)
	}
	return nil
}

func UpdateTestTopologies(t *testing.T) {
	for _, topology := range testData.Topologies {
		if err := updateSingleTopology(topology); err != nil {
			t.Fatalf(err.Error())
		}
	}

	// attempt to add cachegroup that doesn't meet DS required capabilities
	top, _, err := TOSession.GetTopologyWithHdr("top-for-ds-req", nil)
	if err != nil {
		t.Fatalf("cannot GET topology: %v", err)
	}
	top.Nodes = append(top.Nodes, tc.TopologyNode{Cachegroup: "cachegroup1", Parents: []int{0}})
	_, _, err = TOSession.UpdateTopology(top.Name, *top)
	if err == nil {
		t.Errorf("making invalid update to topology - expected: error, actual: nil")
	}

	// attempt to add a cachegroup that only has caches in one CDN while the topology is assigned to DSes from multiple CDNs
	top, _, err = TOSession.GetTopologyWithHdr("top-used-by-cdn1-and-cdn2", nil)
	if err != nil {
		t.Fatalf("cannot GET topology: %v", err)
	}
	params := url.Values{}
	params.Add("topology", "top-used-by-cdn1-and-cdn2")
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET delivery services: %v", err)
	}
	if len(dses) < 2 {
		t.Fatalf("expected at least 2 delivery services assigned to topology top-used-by-cdn1-and-cdn2, actual: %d", len(dses))
	}
	foundCDN1 := false
	foundCDN2 := false
	for _, ds := range dses {
		if *ds.CDNName == "cdn1" {
			foundCDN1 = true
		} else if *ds.CDNName == "cdn2" {
			foundCDN2 = true
		}
	}
	if !foundCDN1 || !foundCDN2 {
		t.Fatalf("expected delivery services assigned to topology top-used-by-cdn1-and-cdn2 to be assigned to cdn1 and cdn2")
	}
	cgs, _, err := TOSession.GetCacheGroupNullableByNameWithHdr("cdn1-only", nil)
	if err != nil {
		t.Fatalf("unable to GET cachegroup by name: %v", err)
	}
	if len(cgs) != 1 {
		t.Fatalf("expected: to get 1 cachegroup named 'cdn1-only', actual: got %d", len(cgs))
	}
	params = url.Values{}
	params.Add("cachegroup", strconv.Itoa(*cgs[0].ID))
	servers, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("unable to GET servers by cachegroup: %v", err)
	}
	for _, s := range servers.Response {
		if *s.Cachegroup != "cdn1-only" {
			t.Fatalf("GET servers by cachegroup 'cdn1-only' - expected: only servers in cachegroup 'cdn1-only', actual: got server in %s", *s.Cachegroup)
		}
		if *s.CDNName != "cdn1" {
			t.Fatalf("expected: servers in cachegroup 'cdn1-only' to only be in cdn1, actual: servers in cdn %s", *s.CDNName)
		}
	}
	top.Nodes = append(top.Nodes, tc.TopologyNode{
		Cachegroup: "cdn1-only",
		Parents:    []int{0},
	})
	_, _, err = TOSession.UpdateTopology(top.Name, *top)
	if err == nil {
		t.Errorf("making invalid update to topology (cachegroup contains only servers from cdn1 while the topology is assigned to delivery services in cdn1 and cdn2) - expected: error, actual: nil")
	}
}

func UpdateValidateTopologyORGServerCacheGroup(t *testing.T) {
	params := url.Values{}
	params.Set("xmlId", "ds-top")

	//Get the correct DS
	remoteDS, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Errorf("cannot GET Delivery Services: %v", err)
	}

	//Assign ORG server to DS
	assignServer := []string{"denver-mso-org-01"}
	_, _, err = TOSession.AssignServersToDeliveryService(assignServer, *remoteDS[0].XMLID)
	if err != nil {
		t.Errorf("cannot assign server to Delivery Services: %v", err)
	}

	//Get Topology node to update and remove ORG server nodes
	origTopo := *remoteDS[0].Topology
	resp, _, err := TOSession.GetTopologyWithHdr(origTopo, nil)
	if err != nil {
		t.Fatalf("couldn't find any topologies: %v", err)
	}

	// remove org server cachegroup
	var p []int
	newNodes := []tc.TopologyNode{{Id: 0, Cachegroup: "topology-edge-cg-01", Parents: p, LastUpdated: nil}}
	if *remoteDS[0].Topology == resp.Name {
		resp.Nodes = newNodes
	}
	_, _, err = TOSession.UpdateTopology(*remoteDS[0].Topology, *resp)
	if err == nil {
		t.Fatalf("shouldnot UPDATE topology:%v to %v, but update was a success", *remoteDS[0].Topology, newNodes[0].Cachegroup)
	} else if !strings.Contains(err.Error(), "ORG servers are assigned to delivery services that use this topology, and their cachegroups cannot be removed:") {
		t.Errorf("expected error messsage containing: \"ORG servers are assigned to delivery services that use this topology, and their cachegroups cannot be removed\", got:%s", err.Error())

	}

	//Remove org server assignment and reset DS back to as it was for further testing
	params.Set("hostName", "denver-mso-org-01")
	serverResp, _, err := TOSession.GetServersWithHdr(&params, nil)
	if len(serverResp.Response) == 0 {
		t.Fatal("no servers in response, quitting")
	}
	if serverResp.Response[0].ID == nil {
		t.Fatal("ID of the response server is nil, quitting")
	}
	_, _, err = TOSession.DeleteDeliveryServiceServer(*remoteDS[0].ID, *serverResp.Response[0].ID)
	if err != nil {
		t.Errorf("cannot delete assigned server from Delivery Services: %v", err)
	}
}

func DeleteTestTopologies(t *testing.T) {
	for _, top := range testData.Topologies {
		delResp, _, err := TOSession.DeleteTopology(top.Name)
		if err != nil {
			t.Fatalf("cannot DELETE topology: %v - %v", err, delResp)
		}
		deleteLog, _, err := TOSession.GetLogsByLimit(1)
		if err != nil {
			t.Fatalf("unable to get latest audit log entry")
		}
		if len(deleteLog) != 1 {
			t.Fatalf("log entry length - expected: 1, actual: %d", len(deleteLog))
		}
		if !strings.Contains(*deleteLog[0].Message, top.Name) {
			t.Errorf("topology deletion audit log entry - expected: message containing topology name '%s', actual: %s", top.Name, *deleteLog[0].Message)
		}

		topology, _, err := TOSession.GetTopology(top.Name)
		if err == nil {
			t.Fatalf("expected error trying to GET deleted topology: %s, actual: nil", top.Name)
		}
		if topology != nil {
			t.Fatalf("expected nil trying to GET deleted topology: %s, actual: non-nil", top.Name)
		}
	}
}
