package v4

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
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

type topologyTestCase struct {
	testCaseDescription string
	tc.Topology
}

func TestTopologies(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, CDNs, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, ServerCapabilities, ServerServerCapabilitiesForTopologies, Topologies, Tenants, DeliveryServices, TopologyBasedDeliveryServiceRequiredCapabilities}, func() {
		GetTestTopologies(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		rfcTime := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, rfcTime)
		header.Set(rfc.IfUnmodifiedSince, rfcTime)
		UpdateTestTopologies(t)
		UpdateTestTopologiesWithHeaders(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestTopologiesWithHeaders(t, header)
		ValidationTestTopologies(t)
		UpdateValidateTopologyORGServerCacheGroup(t)
		EdgeParentOfEdgeSucceedsWithWarning(t)
		UpdateTopologyName(t)
		GetTopologyWithNonExistentName(t)
		CreateTopologyWithInvalidCacheGroup(t)
		CreateTopologyWithInvalidParentNumber(t)
		CreateTopologyWithoutDescription(t)
		CreateTopologyWithoutName(t)
		CreateTopologyWithoutServers(t)
		CreateTopologyWithDuplicateParents(t)
		CreateTopologyWithNodeAsParentOfItself(t)
		CreateTopologyWithOrgLocAsChildNode(t)
		CreateTopologyWithExistingName(t)
		CreateTopologyWithMidLocTypeWithoutChild(t)
		CRUDTopologyReadOnlyUser(t)
		UpdateTopologyWithCachegroupAssignedToBecomeParentOfItself(t)
		UpdateTopologyWithNoServers(t)
		UpdateTopologyWithInvalidParentNumber(t)
		UpdateTopologyWithMidLocTypeWithoutChild(t)
		UpdateTopologyWithOrgLocAsChildNode(t)
		UpdateTopologyWithSameParentAndSecondaryParent(t)
		DeleteTopologyWithNonExistentName(t)
		DeleteTopologyBeingUsedByDeliveryService(t)
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
		t.Log("Response: ", postResponse)
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

func UpdateTestTopologiesWithHeaders(t *testing.T, header http.Header) {
	originalName := "top-used-by-cdn1-and-cdn2"
	newName := "blah"

	// Retrieve the Topology by name so we can get the id for Update()
	resp, _, err := TOSession.GetTopologyWithHdr(originalName, nil)
	if err != nil {
		t.Errorf("cannot GET Topology by name: '%s', %v", originalName, err)
	}
	if (resp) != nil {
		resp.Name = newName
		_, reqInf, err := TOSession.UpdateTopology(originalName, *resp, header)
		if err == nil {
			t.Errorf("Expected error about Precondition Failed, got none")
		}
		if reqInf.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
		}
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
	updateResponse, _, err := TOSession.UpdateTopology(topology.Name, topology, nil)
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
	_, _, err = TOSession.UpdateTopology(top.Name, *top, nil)
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
	dses, _, err := TOSession.GetDeliveryServicesV4(nil, params)
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
	_, _, err = TOSession.UpdateTopology(top.Name, *top, nil)
	if err == nil {
		t.Errorf("making invalid update to topology (cachegroup contains only servers from cdn1 while the topology is assigned to delivery services in cdn1 and cdn2) - expected: error, actual: nil")
	}
}

func UpdateValidateTopologyORGServerCacheGroup(t *testing.T) {
	params := url.Values{}
	params.Set("xmlId", "ds-top")

	//Get the correct DS
	remoteDS, _, err := TOSession.GetDeliveryServicesV4(nil, params)
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
	_, _, err = TOSession.UpdateTopology(*remoteDS[0].Topology, *resp, nil)
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

func UpdateTopologyName(t *testing.T) {
	currentTopologyName := "top-used-by-cdn1-and-cdn2"

	// Get details on existing topology
	resp, _, err := TOSession.GetTopologyWithHdr(currentTopologyName, nil)
	if err != nil {
		t.Errorf("unable to get topology with name: %v", currentTopologyName)
	}
	newTopologyName := "test-topology"
	resp.Name = newTopologyName

	// Update topology with new name
	updateResponse, _, err := TOSession.UpdateTopology(currentTopologyName, *resp, nil)
	if err != nil {
		t.Errorf("cannot PUT topology: %v - %v", err, updateResponse)
	}
	if updateResponse.Response.Name != newTopologyName {
		t.Errorf("update topology name failed, expected: %v but got:%v", newTopologyName, updateResponse.Response.Name)
	}

	//To check whether the primary key change trickled down to DS table
	resp1, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr("top-ds-in-cdn2", nil)
	if err != nil {
		t.Errorf("failed to get details on DS: %v", err)
	}
	if *resp1[0].Topology != newTopologyName {
		t.Errorf("topology name change failed to trickle to delivery service table, expected: %v but got:%v", newTopologyName, *resp1[0].Topology)
	}

	// Set everything back as it was for further testing.
	resp.Name = currentTopologyName
	r, _, err := TOSession.UpdateTopology(newTopologyName, *resp, nil)
	if err != nil {
		t.Errorf("cannot PUT topology: %v - %v", err, r)
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

func GetTopologyWithNonExistentName(t *testing.T) {
	resp, reqInf, _ := TOSession.GetTopologyWithHdr("non-existent-topology", nil)
	if resp != nil {
		t.Errorf("expected nothing in the response, but got a topology with name %s", resp.Name)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 response code, but got %d", reqInf.StatusCode)
	}
}

func CreateTopologyWithInvalidCacheGroup(t *testing.T) {
	nodes := make([]tc.TopologyNode, 0)
	node := tc.TopologyNode{
		Cachegroup: "non-existent-cachegroup",
		Parents:    nil,
	}
	nodes = append(nodes, node)
	top := tc.Topology{
		Description: "blah",
		Name:        "invalid-cachegroup-topology",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a '400 Bad Request' response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("expected error about the cachegroup name not being valid, but got none")
	}
}

func CreateTopologyWithInvalidParentNumber(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}
	nodes := make([]tc.TopologyNode, 0)

	cachegroupName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}
	node := tc.TopologyNode{
		Cachegroup: cachegroupName,
		Parents:    []int{100},
	}
	nodes = append(nodes, node)
	top := tc.Topology{
		Description: "blah",
		Name:        "invalid-parent-topology",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a '400 Bad Request' response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("expected error about the parent not being valid, but got none")
	}
}

func CreateTopologyWithoutDescription(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}
	nodes := make([]tc.TopologyNode, 0)

	cachegroupName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}
	node := tc.TopologyNode{
		Cachegroup: cachegroupName,
		Parents:    nil,
	}
	nodes = append(nodes, node)
	top := tc.Topology{
		Name:  "topology-without-description",
		Nodes: nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 response code, but got %d", reqInf.StatusCode)
	}
	if err != nil {
		t.Errorf("no error expected about description being empty, but got %v", err)
	}
	_, _, err = TOSession.DeleteTopology(top.Name)
	if err != nil {
		t.Errorf("couldn't delete topology with name %s: %v", top.Name, err)
	}
}

func CreateTopologyWithoutName(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}
	nodes := make([]tc.TopologyNode, 0)

	cachegroupName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}
	node := tc.TopologyNode{
		Cachegroup: cachegroupName,
		Parents:    nil,
	}
	nodes = append(nodes, node)
	top := tc.Topology{
		Description: "description",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology name being empty expected, but got none")
	}
}

func CreateTopologyWithoutServers(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}
	nodes := make([]tc.TopologyNode, 0)

	cachegroupName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) == 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}
	node := tc.TopologyNode{
		Cachegroup: cachegroupName,
		Parents:    nil,
	}
	nodes = append(nodes, node)
	top := tc.Topology{
		Name:        "topology_without_servers",
		Description: "description",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology containing no servers expected, but got none")
	}
}

func CreateTopologyWithDuplicateParents(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}

	cachegroupName := ""
	parentName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && parentName != *cg.Name && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}

	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil {
				parentName = *cg.Name
				break
			}
		}
	}

	nodes := []tc.TopologyNode{
		{
			Cachegroup: parentName,
			Parents:    []int{},
		},
		{
			Cachegroup: cachegroupName,
			Parents:    []int{0, 0},
		},
	}
	top := tc.Topology{
		Name:        "topology_with_duplicate_parents",
		Description: "description",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having duplicate parent expected, but got none")
	}
}

func CreateTopologyWithNodeAsParentOfItself(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}

	cachegroupName := ""
	parentName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}

	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && parentName != *cg.Name {
				parentName = *cg.Name
				break
			}
		}
	}

	nodes := []tc.TopologyNode{
		{
			Cachegroup: parentName,
			Parents:    []int{},
		},
		{
			Cachegroup: cachegroupName,
			Parents:    []int{0, 1},
		},
	}
	top := tc.Topology{
		Name:        "topology_with_node_as_parent_of_itself",
		Description: "description",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having node as parent of itself expected, but got none")
	}
}

func CreateTopologyWithOrgLocAsChildNode(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}

	cachegroupName := ""
	parentName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				parentName = *cg.Name
				break
			}
		}
	}

	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && parentName != *cg.Name && cg.Type != nil && *cg.Type == tc.CacheGroupOriginTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}

	nodes := []tc.TopologyNode{
		{
			Cachegroup: parentName,
			Parents:    []int{},
		},
		{
			Cachegroup: cachegroupName,
			Parents:    []int{0},
		},
	}
	top := tc.Topology{
		Name:        "topology_with_orgloc_as_child_node",
		Description: "description",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having ord_loc node as child expected, but got none")
	}
}

func CreateTopologyWithExistingName(t *testing.T) {
	resp, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("could not GET topologies: %v", err)
	}
	if len(resp) == 0 {
		t.Fatalf("expected 1 or more topologies in response, but got 0")
	}
	_, reqInf, err := TOSession.CreateTopology(resp[0])
	if err == nil {
		t.Errorf("expected error about creating topology with same name, but got none")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
}

func CreateTopologyWithMidLocTypeWithoutChild(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}

	parentName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupMidTypeName {
				parentName = *cg.Name
				break
			}
		}
	}

	nodes := []tc.TopologyNode{
		{
			Cachegroup: parentName,
			Parents:    []int{},
		},
	}
	top := tc.Topology{
		Name:        "topology_with_midloc_and_no_child_nodes",
		Description: "description",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.CreateTopology(top)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having mid_loc node and no children expected, but got none")
	}
}

func CRUDTopologyReadOnlyUser(t *testing.T) {
	resp, _, err := TOSession.TenantByNameWithHdr("root", nil)
	if err != nil {
		t.Fatalf("couldn't get the root tenant ID: %v", err)
	}
	if resp == nil {
		t.Fatalf("expected a valid tenant response, but got nothing")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	user := tc.User{
		Username:             util.StrPtr("test_user"),
		RegistrationSent:     tc.TimeNoModFromTime(time.Now()),
		LocalPassword:        util.StrPtr("test_pa$$word"),
		ConfirmLocalPassword: util.StrPtr("test_pa$$word"),
		RoleName:             util.StrPtr("read-only user"),
	}
	user.Email = util.StrPtr("email@domain.com")
	user.TenantID = util.IntPtr(resp.ID)
	user.FullName = util.StrPtr("firstName LastName")

	u, _, err := TOSession.CreateUser(&user)
	if err != nil {
		t.Fatalf("could not create read-only user: %v", err)
	}
	client, _, err := toclient.LoginWithAgent(TOSession.URL, "test_user", "test_pa$$word", true, "to-api-v4-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with test_user: %v", err.Error())
	}
	nodes := []tc.TopologyNode{
		{
			Cachegroup: "parentName",
			Parents:    []int{},
		},
	}
	top := tc.Topology{
		Name:        "topology",
		Description: "description",
		Nodes:       nodes,
	}
	// Create
	_, reqInf, err := client.CreateTopology(top)
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 Forbidden error, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("expected error about Read-Only users not being able to create topologies, but got nothing")
		_, _, err = client.DeleteTopology(top.Name)
		if err != nil {
			t.Errorf("could not delete topology %s: %v", top.Name, err)
		}
	}

	// Read
	tops, _, err := client.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("couldn't get topologies: %v", err)
	}
	if len(tops) == 0 {
		t.Fatal("expected to get one or more topologies in the response, but got none")
	}

	// Update
	updatedTop := tops[0]
	updatedTop.Description = "updated description"
	_, reqInf, err = client.UpdateTopology(updatedTop.Name, updatedTop, nil)
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 Forbidden error, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("expected error about Read-Only users not being able to update topologies, but got nothing")
	}

	// Delete
	_, reqInf, err = client.DeleteTopology(tops[0].Name)
	if reqInf.StatusCode != http.StatusForbidden {
		t.Errorf("expected a 403 Forbidden error, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("expected error about Read-Only users not being able to delete topologies, but got nothing")
	}

	if u != nil && u.Response.Username != nil {
		ForceDeleteTestUsersByUsernames(t, []string{"test_user"})
	}
}

func UpdateTopologyWithCachegroupAssignedToBecomeParentOfItself(t *testing.T) {
	tops, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("couldn't get topologies: %v", err)
	}
	if len(tops) == 0 {
		t.Fatal("expected to get one or more topologies in the response, but got none")
	}
	tp := tops[0]
	parents := make([]int, 0)

	// create a list of indices consisting of all the node indices,
	// so that when we assign this parent list wile updating,
	// TO complains about the parent of a node being the same as itself
	for i, _ := range tp.Nodes {
		parents = append(parents, i)
	}
	nodes := tp.Nodes
	for i, _ := range nodes {
		nodes[i].Parents = parents
	}
	tp.Nodes = nodes

	tops[0] = tp
	_, reqInf, err := TOSession.UpdateTopology(tops[0].Name, tops[0], nil)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having parents the same as children expected, but got none")
	}
}

func UpdateTopologyWithSameParentAndSecondaryParent(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}

	cachegroupName := ""
	parentName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				parentName = *cg.Name
				break
			}
		}
	}

	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && parentName != *cg.Name && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}

	nodes := []tc.TopologyNode{
		{
			Cachegroup: parentName,
			Parents:    []int{},
		},
		{
			Cachegroup: cachegroupName,
			Parents:    []int{0, 0},
		},
	}

	tops, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("couldn't get topologies: %v", err)
	}
	if len(tops) == 0 {
		t.Fatal("expected to get one or more topologies in the response, but got none")
	}

	top := tc.Topology{
		Description: tops[0].Description,
		Name:        tops[0].Name,
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.UpdateTopology(tops[0].Name, top, nil)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology's cachegroup having same primary and secondary parents, but got none")
	}
}

func UpdateTopologyWithOrgLocAsChildNode(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}

	cachegroupName := ""
	parentName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				parentName = *cg.Name
				break
			}
		}
	}

	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && parentName != *cg.Name && cg.Type != nil && *cg.Type == tc.CacheGroupOriginTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}

	nodes := []tc.TopologyNode{
		{
			Cachegroup: parentName,
			Parents:    []int{},
		},
		{
			Cachegroup: cachegroupName,
			Parents:    []int{0},
		},
	}

	tops, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("couldn't get topologies: %v", err)
	}
	if len(tops) == 0 {
		t.Fatal("expected to get one or more topologies in the response, but got none")
	}
	top := tc.Topology{
		Name:        "topology_with_orgloc_as_child_node",
		Description: "description",
		Nodes:       nodes,
	}
	_, reqInf, err := TOSession.UpdateTopology(tops[0].Name, top, nil)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having ord_loc node as child expected, but got none")
	}
}

func UpdateTopologyWithMidLocTypeWithoutChild(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}

	parentName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) != 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupMidTypeName {
				parentName = *cg.Name
				break
			}
		}
	}

	nodes := []tc.TopologyNode{
		{
			Cachegroup: parentName,
			Parents:    []int{},
		},
	}
	tops, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("couldn't get topologies: %v", err)
	}
	if len(tops) == 0 {
		t.Fatal("expected to get one or more topologies in the response, but got none")
	}

	top := tc.Topology{
		Name:        tops[0].Name,
		Description: tops[0].Description,
		Nodes:       nodes,
	}

	_, reqInf, err := TOSession.UpdateTopology(tops[0].Name, top, nil)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having mid_loc node and no children expected, but got none")
	}
}

func UpdateTopologyWithInvalidParentNumber(t *testing.T) {
	tops, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("couldn't get topologies: %v", err)
	}
	if len(tops) == 0 {
		t.Fatal("expected to get one or more topologies in the response, but got none")
	}
	tp := tops[0]
	parents := make([]int, 0)
	parents = append(parents, len(tp.Nodes)+1)
	for i, _ := range tp.Nodes {
		tp.Nodes[i].Parents = parents
	}
	_, reqInf, err := TOSession.UpdateTopology(tp.Name, tp, nil)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a '400 Bad Request' response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("expected error about the parent not being valid, but got none")
	}
}

func UpdateTopologyWithNoServers(t *testing.T) {
	cacheGroups, _, err := TOSession.GetCacheGroupsNullableWithHdr(nil)
	if err != nil {
		t.Fatalf("error while getting cachegroups: %v", err)
	}
	if len(cacheGroups) == 0 {
		t.Fatalf("no cachegroups in response")
	}
	nodes := make([]tc.TopologyNode, 0)

	cachegroupName := ""
	params := url.Values{}
	for _, cg := range cacheGroups {
		params["cachegroup"] = []string{strconv.Itoa(*cg.ID)}
		resp, _, _ := TOSession.GetServersWithHdr(&params, nil)
		if len(resp.Response) == 0 {
			if cg.Name != nil && cg.Type != nil && *cg.Type == tc.CacheGroupEdgeTypeName {
				cachegroupName = *cg.Name
				break
			}
		}
	}
	node := tc.TopologyNode{
		Cachegroup: cachegroupName,
		Parents:    nil,
	}
	nodes = append(nodes, node)

	tops, _, err := TOSession.GetTopologiesWithHdr(nil)
	if err != nil {
		t.Fatalf("error getting topologies: %v", err)
	}
	if len(tops) == 0 {
		t.Fatalf("expected 1 or more topologies in response, but got none")
	}
	tops[0].Nodes = nodes
	_, reqInf, err := TOSession.UpdateTopology(tops[0].Name, tops[0], nil)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology having no servers expected, but got none")
	}
}

func DeleteTopologyBeingUsedByDeliveryService(t *testing.T) {
	ds, _, err := TOSession.GetDeliveryServicesV4(nil, nil)
	if err != nil {
		t.Fatalf("couldn't get deliveryservices: %v", err)
	}
	if len(ds) == 0 {
		t.Fatalf("expected one or more ds's in the response, got none")
	}
	topologyName := ""
	for _, d := range ds {
		if d.Topology != nil && *d.Topology != "" {
			topologyName = *d.Topology
			break
		}
	}
	_, reqInf, err := TOSession.DeleteTopology(topologyName)
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology being used by a ds expected, but got none")
	}
}

func DeleteTopologyWithNonExistentName(t *testing.T) {
	_, reqInf, err := TOSession.DeleteTopology("non existent name")
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("expected a 400 response code, but got %d", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("error about topology not being present expected, but got none")
	}
}
