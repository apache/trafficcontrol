package topology

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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func checkUniqueCacheGroupNames(nodes []tc.TopologyNode) error {
	cacheGroupNames := map[string]bool{}
	for _, node := range nodes {
		if _, exists := cacheGroupNames[node.Cachegroup]; exists {
			return fmt.Errorf("cachegroup %v cannot be used more than once in the topology", node.Cachegroup)
		}
		cacheGroupNames[node.Cachegroup] = true
	}
	return nil
}

func checkForDuplicateParents(nodes []tc.TopologyNode, index int) error {
	parents := nodes[index].Parents
	if len(parents) != 2 || parents[0] != parents[1] {
		return nil
	}
	return fmt.Errorf("cachegroup %v cannot be both a primary and secondary parent of cachegroup %v", nodes[parents[0]].Cachegroup, nodes[index].Cachegroup)
}

func checkForSelfParents(nodes []tc.TopologyNode, index int) error {
	for _, parentIndex := range nodes[index].Parents {
		if index == parentIndex {
			return fmt.Errorf("cachegroup %v cannot be a parent of itself", index)
		}
	}
	return nil
}

func checkForEdgeParents(nodes []tc.TopologyNode, cachegroups []tc.CacheGroupNullable, nodeIndex int) error {
	node := nodes[nodeIndex]
	errs := make([]error, len(node.Parents))
	for parentIndex, cachegroupIndex := range node.Parents {
		if cachegroupIndex < 0 || cachegroupIndex >= len(cachegroups) {
			errs = append(errs, fmt.Errorf("parent %d of cachegroup %s refers to a cachegroup at index %d, but no such cachegroup exists", parentIndex, node.Cachegroup, cachegroupIndex))
			break
		}
		cacheGroupType := cachegroups[cachegroupIndex].Type
		if *cacheGroupType == tc.CacheGroupEdgeTypeName {
			errs = append(errs, fmt.Errorf("cachegroup %v's type is %v; it cannot be a parent of %v", nodes[cachegroupIndex].Cachegroup, tc.CacheGroupEdgeTypeName, node.Cachegroup))
		}
	}
	return util.JoinErrs(errs)
}

func checkForLeafMids(nodes []tc.TopologyNode, cacheGroups []tc.CacheGroupNullable) []tc.TopologyNode {
	isLeafMid := make([]bool, len(nodes))
	for index := range isLeafMid {
		isLeafMid[index] = true
	}
	for index, node := range nodes {
		if *cacheGroups[index].Type == tc.CacheGroupEdgeTypeName {
			isLeafMid[index] = false
		}
		for _, parentIndex := range node.Parents {
			if !isLeafMid[parentIndex] {
				continue
			}
			isLeafMid[parentIndex] = false
		}
	}

	var leafMids []tc.TopologyNode
	for index, node := range nodes {
		if isLeafMid[index] {
			leafMids = append(leafMids, node)
		}
	}
	return leafMids
}

func checkForCycles(nodes []tc.TopologyNode) error {
	components := tarjan(nodes)
	var errs []error
	for _, component := range components {
		if len(component) > 1 {
			errString := "cycle detected between cachegroups "
			var node tc.TopologyNode
			for _, node = range component {
				errString += node.Cachegroup + ", "
			}
			length := len(errString)
			cachegroupNameLength := len(node.Cachegroup)
			errString = errString[0:length-2-cachegroupNameLength-2] + " and " + errString[length-2-cachegroupNameLength:length-2]
			errs = append(errs, fmt.Errorf(errString))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	errs = append([]error{fmt.Errorf("topology cannot have cycles")}, errs...)
	return util.JoinErrs(errs)
}
