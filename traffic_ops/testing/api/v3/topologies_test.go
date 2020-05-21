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
	"reflect"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

type topologyTestCase struct {
	reasonToFail string
	tc.Topology
}

func TestTopologies(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, Topologies}, func() {
		UpdateTestTopologies(t)
		ValidationTestTopologies(t)
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

func ValidationTestTopologies(t *testing.T) {
	invalidTopologyTestCases := []topologyTestCase{
		{reasonToFail: "no nodes", Topology: tc.Topology{Name: "empty-top", Description: "Invalid because there are no nodes", Nodes: []tc.TopologyNode{}}},
		{reasonToFail: "a node listing itself as a parent", Topology: tc.Topology{Name: "self-parent", Description: "Invalid because a node lists itself as a parent", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1}},
			{Cachegroup: "parentCachegroup", Parents: []int{1}},
		}}},
		{reasonToFail: "duplicate parents", Topology: tc.Topology{}},
		{reasonToFail: "too many parents", Topology: tc.Topology{Name: "duplicate-parents", Description: "Invalid because a node lists the same parent twice", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1, 1}},
			{Cachegroup: "parentCachegroup", Parents: []int{}},
		}}},
		{reasonToFail: "too many parents", Topology: tc.Topology{Name: "too-many-parents", Description: "Invalid because a node has more than 2 parents", Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{}},
			{Cachegroup: "parentCachegroup2", Parents: []int{}},
			{Cachegroup: "cachegroup1", Parents: []int{0, 1, 2}},
		}}},
		{reasonToFail: "a parent edge", Topology: tc.Topology{Name: "parent-edge", Description: "Invalid because an edge is a parent", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1}},
			{Cachegroup: "cachegroup2", Parents: []int{}},
		}}},
		{reasonToFail: "a leaf mid", Topology: tc.Topology{Name: "leaf-mid", Description: "Invalid because a mid is a leaf node", Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{1}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{}},
		}}},
		{reasonToFail: "cyclical nodes", Topology: tc.Topology{Name: "cyclical-nodes", Description: "Invalid because it contains cycles", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1, 2}},
			{Cachegroup: "parentCachegroup", Parents: []int{2}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{1}},
		}}},
		{reasonToFail: "a nonexistent cache group", Topology: tc.Topology{Name: "nonexistent-cg", Description: "Invalid because it references a cache group that does not exist", Nodes: []tc.TopologyNode{
			{Cachegroup: "legitcachegroup", Parents: []int{0}},
		}}},
		{reasonToFail: "an out-of-bounds parent index", Topology: tc.Topology{Name: "oob-parent", Description: "Invalid because it contains a parent", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{7}},
		}}},
	}
	var statusCode int
	for _, testCase := range invalidTopologyTestCases {
		_, reqInf, err := TOSession.CreateTopology(testCase.Topology)
		if err == nil {
			t.Fatalf("expected POST with %v to return an error, actual: nil", testCase.reasonToFail)
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
	topologiesCount := len(testData.Topologies)
	for index := range testData.Topologies {
		topology := testData.Topologies[(index+1)%topologiesCount]
		topology.Name = testData.Topologies[index].Name // We cannot update a topology's name
		if err := updateSingleTopology(topology); err != nil {
			t.Fatalf(err.Error())
		}
	}
	// Revert test topologies
	for _, topology := range testData.Topologies {
		if err := updateSingleTopology(topology); err != nil {
			t.Fatalf(err.Error())
		}
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
