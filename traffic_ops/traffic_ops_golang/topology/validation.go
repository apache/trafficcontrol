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
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
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

// checkForEdgeParents adds a warning + returns a nil error if
// an edge parents an edge, and returns an error if an edge parents a non-edge cachegroup.
func checkForEdgeParents(topology tc.TopologyV5, cacheGroups []tc.CacheGroupNullable, nodeIndex int) (tc.Alerts, error) {
	var alerts tc.Alerts
	node := topology.Nodes[nodeIndex]
	errs := make([]error, len(node.Parents))
	for parentIndex, parentCacheGroupIndex := range node.Parents {
		if parentCacheGroupIndex < 0 || parentCacheGroupIndex >= len(topology.Nodes) {
			errs = append(errs, fmt.Errorf("parent %d of cachegroup %s refers to a cachegroup at index %d, but no such cachegroup exists", parentIndex, node.Cachegroup, parentCacheGroupIndex))
			break
		}
		parentCacheGroupType := *cacheGroups[parentCacheGroupIndex].Type
		if parentCacheGroupType != tc.CacheGroupEdgeTypeName {
			continue
		}
		switch cacheGroupType := *cacheGroups[nodeIndex].Type; cacheGroupType {
		case tc.CacheGroupEdgeTypeName:
			parentTerm := "parent"
			if parentIndex == 1 {
				parentTerm = "secondary " + parentTerm
			}
			alerts.AddNewAlert(tc.WarnLevel, fmt.Sprintf(
				"%s-typed cachegroup %s is a %s of %s, unexpected behavior may result",
				parentCacheGroupType,
				topology.Nodes[parentCacheGroupIndex].Cachegroup,
				parentTerm,
				node.Cachegroup))
		default:
			errs = append(errs, fmt.Errorf(
				"cachegroup %s's type is %s; it cannot parent a %s-typed cachegroup %s",
				topology.Nodes[parentCacheGroupIndex].Cachegroup,
				parentCacheGroupType,
				cacheGroupType,
				node.Cachegroup))
		}
	}
	return alerts, util.JoinErrs(errs)
}

// checkForEdgeParents returns an error if an index given in the parents array, adds a warning + returns a nil error if
// an edge parents an edge, and returns an error if an edge parents a non-edge cachegroup.
func (topology *TOTopology) checkForEdgeParents(cacheGroups []tc.CacheGroupNullable, nodeIndex int) error {
	node := topology.Nodes[nodeIndex]
	errs := make([]error, len(node.Parents))
	for parentIndex, parentCacheGroupIndex := range node.Parents {
		if parentCacheGroupIndex < 0 || parentCacheGroupIndex >= len(topology.Nodes) {
			errs = append(errs, fmt.Errorf("parent %d of cachegroup %s refers to a cachegroup at index %d, but no such cachegroup exists", parentIndex, node.Cachegroup, parentCacheGroupIndex))
			break
		}
		parentCacheGroupType := *cacheGroups[parentCacheGroupIndex].Type
		if parentCacheGroupType != tc.CacheGroupEdgeTypeName {
			continue
		}
		switch cacheGroupType := *cacheGroups[nodeIndex].Type; cacheGroupType {
		case tc.CacheGroupEdgeTypeName:
			parentTerm := "parent"
			if parentIndex == 1 {
				parentTerm = "secondary " + parentTerm
			}
			topology.Alerts.AddNewAlert(tc.WarnLevel, fmt.Sprintf(
				"%s-typed cachegroup %s is a %s of %s, unexpected behavior may result",
				parentCacheGroupType,
				topology.Nodes[parentCacheGroupIndex].Cachegroup,
				parentTerm,
				node.Cachegroup))
		default:
			errs = append(errs, fmt.Errorf(
				"cachegroup %s's type is %s; it cannot parent a %s-typed cachegroup %s",
				topology.Nodes[parentCacheGroupIndex].Cachegroup,
				parentCacheGroupType,
				cacheGroupType,
				node.Cachegroup))
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

func checkForCycles(nodes []tc.TopologyNode) ([]string, error) {
	components := tarjan(nodes)
	var (
		errs        []error
		cacheGroups []string
	)
	for _, component := range components {
		if len(component) > 1 {
			errString := "cycle detected between cachegroups "
			var node tc.TopologyNode
			for _, node = range component {
				cacheGroups = append(cacheGroups, node.Cachegroup)
				errString += node.Cachegroup + ", "
			}
			length := len(errString)
			cachegroupNameLength := len(node.Cachegroup)
			errString = errString[0:length-2-cachegroupNameLength-2] + " and " + errString[length-2-cachegroupNameLength:length-2]
			errs = append(errs, fmt.Errorf(errString))
		}
	}
	if len(errs) == 0 {
		return nil, nil
	}
	errs = append([]error{fmt.Errorf("topology cannot have cycles")}, errs...)
	return cacheGroups, util.JoinErrs(errs)
}

func checkForCyclesAcrossTopologies(info *api.Info, topologyNodes []tc.TopologyNode, name string) error {
	var (
		nodes                  []tc.TopologyNode
		topologiesByCacheGroup map[string][]string
		cacheGroups            []string
		err                    error
	)
	if nodes, topologiesByCacheGroup, err = nodesInOtherTopologies(info, topologyNodes); err != nil {
		return err
	}
	if cacheGroups, err = checkForCycles(nodes); err == nil {
		return nil
	}
	if cacheGroups == nil {
		return fmt.Errorf("unable to check topology %s for cycles across all topologies", name)
	}
	var involvedTopologies []string
	includedTopology := map[string]bool{}
	for _, cacheGroup := range cacheGroups {
		for _, topology := range topologiesByCacheGroup[cacheGroup] {
			if _, alreadyIncluded := includedTopology[topology]; alreadyIncluded {
				continue
			}

			involvedTopologies = append(involvedTopologies, topology)
			includedTopology[topology] = true
		}
	}
	return fmt.Errorf("cycles exist between topology %s and topologies [%s]: %v", name, strings.Join(involvedTopologies, ", "), err)
}
