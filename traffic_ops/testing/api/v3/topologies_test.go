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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"testing"
)

func TestTopologies(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, Topologies}, func() {
		ValidationTestTopologies(t)
	})
}

func CreateTestTopologies(t *testing.T) {
	for _, top := range testData.Topologies {
		resp, _, err := TOSession.CreateTopology(top)
		if err != nil {
			t.Errorf("could not CREATE topology: %v", err)
		}
		t.Log("Response: ", resp)
	}
}

func ValidationTestTopologies(t *testing.T) {
	invalidTopologies := []tc.Topology{
		{Name: "empty-top", Description: "Invalid because there are no nodes", Nodes: []tc.TopologyNode{}},
		{Name: "self-parent", Description: "Invalid because a node lists itself as a parent", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1}},
			{Cachegroup: "parentCachegroup", Parents: []int{1}},
		}},
		{Name: "duplicate-parents", Description: "Invalid because a node lists the same parent twice", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1, 1}},
			{Cachegroup: "parentCachegroup", Parents: []int{}},
		}},
		{Name: "too-many-parents", Description: "Invalid because a node has more than 2 parents", Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{}},
			{Cachegroup: "parentCachegroup2", Parents: []int{}},
			{Cachegroup: "cachegroup1", Parents: []int{0, 1, 2}},
		}},
		{Name: "parent-edge", Description: "Invalid because an edge is a parent", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1}},
			{Cachegroup: "cachegroup2", Parents: []int{}},
		}},
		{Name: "leaf-mid", Description: "Invalid because a mid is a leaf node", Nodes: []tc.TopologyNode{
			{Cachegroup: "parentCachegroup", Parents: []int{1}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{}},
		}},
		{Name: "cyclical-nodes", Description: "Invalid because it contains cycles", Nodes: []tc.TopologyNode{
			{Cachegroup: "cachegroup1", Parents: []int{1, 2}},
			{Cachegroup: "parentCachegroup", Parents: []int{2}},
			{Cachegroup: "secondaryCachegroup", Parents: []int{1}},
		}},
	}
	expectations := []string{
		"no nodes",
		"a node listing itself as a parent",
		"duplicate parents",
		"too many parents",
		"a parent edge",
		"a leaf mid",
		"cyclical nodes",
	}
	for index, invalidTopology := range invalidTopologies {
		if _, _, err := TOSession.CreateTopology(invalidTopology); err == nil {
			t.Errorf("expected POST with %v to return an error, actual: nil", expectations[index])
		}
	}
}

func DeleteTestTopologies(t *testing.T) {
	for _, top := range testData.Topologies {
		delResp, _, err := TOSession.DeleteTopology(top.Name)
		if err != nil {
			t.Errorf("cannot DELETE topology: %v - %v", err, delResp)
		}

		topology, _, err := TOSession.GetTopology(top.Name)
		if err == nil {
			t.Errorf("expected error trying to GET deleted topology: %s, actual: nil", top.Name)
		}
		if topology != nil {
			t.Errorf("expected nil trying to GET deleted topology: %s, actual: non-nil", top.Name)
		}
	}
}
