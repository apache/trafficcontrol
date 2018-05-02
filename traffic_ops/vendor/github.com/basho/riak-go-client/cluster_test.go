// Copyright 2015-present Basho Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package riak

import (
	"fmt"
	"net"
	"testing"
)

func TestCreateClusterWithDefaultOptions(t *testing.T) {
	cluster, err := NewCluster(nil)
	if err != nil {
		t.Error(err.Error())
	}
	if cluster.nodes == nil {
		t.Errorf("expected non-nil value")
	}
	if expected, actual := 1, len(cluster.nodes); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	defaultNodeAddr := cluster.nodes[0].addr.String()
	if expected, actual := defaultRemoteAddress, defaultNodeAddr; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := clusterCreated, cluster.getState(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if cluster.nodeManager == nil {
		t.Error("expected cluster to have a node manager")
	}
	if expected, actual := defaultExecutionAttempts, cluster.executionAttempts; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestCreateClusterWithoutDefaultNodeAndStart(t *testing.T) {
	o := &ClusterOptions{
		NoDefaultNode: true,
	}
	cluster, err := NewCluster(o)
	if err != nil {
		t.Error(err.Error())
	}
	if cluster.nodes == nil {
		t.Errorf("expected non-nil value")
	}
	if expected, actual := 0, len(cluster.nodes); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := clusterCreated, cluster.getState(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if cluster.nodeManager == nil {
		t.Error("expected cluster to have a node manager")
	}
	if expected, actual := defaultExecutionAttempts, cluster.executionAttempts; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	err = cluster.Start()
	if err != nil {
		t.Error(err.Error())
	}
	if expected, actual := clusterRunning, cluster.getState(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestAddAndRemoveNodeFromCluster(t *testing.T) {
	var err error
	var c *Cluster
	c, err = NewCluster(nil)
	if err != nil {
		t.Fatal(err)
	}

	node := c.nodes[0]
	// re-adding same node instance won't add it
	if err = c.AddNode(node); err != nil {
		t.Fatal(err)
	}
	if expected, actual := 1, len(c.nodes); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	// add 4 more nodes
	var addrRemoved *net.TCPAddr
	var nodeToRemove *Node
	portToRemove := 10027
	addrRemoved, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", portToRemove))
	for port := 10017; port <= 10047; port += 10 {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		opts := &NodeOptions{
			RemoteAddress: addr,
		}
		var n *Node
		if n, err = NewNode(opts); err != nil {
			t.Fatal(err)
		} else {
			if err = c.AddNode(n); err != nil {
				t.Fatal(err)
			}
			if port == portToRemove {
				nodeToRemove = n
			}
		}
	}
	if expected, actual := 5, len(c.nodes); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := true, c.isCurrentState(clusterCreated); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	for _, n := range c.nodes {
		if expected, actual := true, n.isCurrentState(nodeCreated); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	}
	// remove node with port 10027
	if err = c.RemoveNode(nodeToRemove); err != nil {
		t.Fatal(err)
	}
	if expected, actual := 4, len(c.nodes); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	for _, n := range c.nodes {
		if expected, actual := true, n.isCurrentState(nodeCreated); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if n.addr == addrRemoved {
			t.Errorf("node with addr %v should have been removed", addrRemoved)
		}
	}
}

func TestCreateClusterWithFourNodes(t *testing.T) {
	nodes := make([]*Node, 0, 4)
	for port := 10017; port <= 10047; port += 10 {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		opts := &NodeOptions{
			RemoteAddress: addr,
		}
		if node, err := NewNode(opts); err != nil {
			t.Error(err.Error())
		} else {
			nodes = append(nodes, node)
		}
	}

	opts := &ClusterOptions{
		Nodes: nodes,
	}
	cluster, err := NewCluster(opts)
	if err != nil {
		t.Error(err.Error())
	}
	if expected, actual := 4, len(cluster.nodes); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	for i, node := range cluster.nodes {
		port := 10007 + ((i + 1) * 10)
		expectedAddr := fmt.Sprintf("127.0.0.1:%d", port)
		if expected, actual := expectedAddr, node.addr.String(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	}
	if expected, actual := clusterCreated, cluster.getState(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if cluster.nodeManager == nil {
		t.Error("expected cluster to have a node manager")
	}
}

func ExampleNewCluster() {
	cluster, err := NewCluster(nil)
	if err != nil {
		panic(fmt.Sprintf("Error building cluster object: %s", err.Error()))
	}
	fmt.Println(cluster.nodes[0].addr.String())
	// Output: 127.0.0.1:8087
}
