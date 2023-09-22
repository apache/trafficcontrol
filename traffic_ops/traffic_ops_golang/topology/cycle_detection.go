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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"math"
)

type TarjanNode struct {
	tc.TopologyNode
	Index   *int
	LowLink *int
	OnStack *bool
}

type NodeStack []*TarjanNode
type Graph []*TarjanNode
type Component []tc.TopologyNode

type Tarjan struct {
	Graph      *Graph
	Stack      *NodeStack
	Components []Component
	Index      int
}

func (stack *NodeStack) push(node *TarjanNode) {
	*stack = append(append([]*TarjanNode{}, *stack...), node)
}

func (stack *NodeStack) pop() *TarjanNode {
	length := len(*stack)
	node := (*stack)[length-1]
	*stack = (*stack)[:length-1]
	return node
}

func tarjan(nodes []tc.TopologyNode) [][]tc.TopologyNode {
	structs := Tarjan{
		Graph:      &Graph{},
		Stack:      &NodeStack{},
		Components: []Component{},
		Index:      0,
	}
	for _, node := range nodes {
		tarjanNode := TarjanNode{TopologyNode: node, LowLink: new(int)}
		*tarjanNode.LowLink = 500
		*structs.Graph = append(*structs.Graph, &tarjanNode)
	}
	structs.Stack = &NodeStack{}
	structs.Index = 0
	for _, vertex := range *structs.Graph {
		if vertex.Index == nil {
			structs.strongConnect(vertex)
		}
	}
	var components [][]tc.TopologyNode
	for _, component := range structs.Components {
		var componentArray []tc.TopologyNode
		for _, node := range component {
			componentArray = append(componentArray, node)
		}
		components = append(components, componentArray)
	}
	return components
}

func (structs *Tarjan) nextComponent() (Component, int) {
	var component = Component{}
	index := len(structs.Components)
	structs.Components = append(structs.Components, component)
	return component, index
}

func (structs *Tarjan) strongConnect(node *TarjanNode) {
	stack := structs.Stack
	node.Index = new(int)
	*node.Index = structs.Index
	node.LowLink = new(int)
	*node.LowLink = structs.Index
	structs.Index++
	stack.push(node)
	node.OnStack = new(bool)
	*node.OnStack = true

	for _, parentIndex := range node.Parents {
		parent := (*structs.Graph)[parentIndex]
		if parent.Index == nil {
			structs.strongConnect(parent)
			*(*parent).LowLink = int(math.Min(float64(*node.LowLink), float64(*parent.LowLink)))
		} else if *parent.OnStack {
			*node.LowLink = int(math.Min(float64(*node.LowLink), float64(*parent.Index)))
		}
	}

	if *node.LowLink == *node.Index {
		component, index := structs.nextComponent()
		var otherNode *TarjanNode
		for node != otherNode {
			otherNode = stack.pop()
			*otherNode.OnStack = false
			component = append(component, otherNode.TopologyNode)
		}
		structs.Components[index] = component
	}
}
