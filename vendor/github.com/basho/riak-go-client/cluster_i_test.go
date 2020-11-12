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

// +build integration

package riak

import (
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestExecuteCommandOnCluster(t *testing.T) {
	nodeOpts := &NodeOptions{
		RemoteAddress: getRiakAddress(),
	}
	var node *Node
	var err error
	if node, err = NewNode(nodeOpts); err != nil {
		t.Fatal(err)
	}
	if node == nil {
		t.Fatal()
	}
	nodes := []*Node{node}
	opts := &ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: 3,
	}

	if expected, actual := 1, len(opts.Nodes); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := node, opts.Nodes[0]; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	cluster, err := NewCluster(opts)
	if err != nil {
		t.Error(err.Error())
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	if expected, actual := node, cluster.nodes[0]; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	if err := cluster.Start(); err != nil {
		t.Error(err.Error())
	}

	command := &PingCommand{}
	if err := cluster.Execute(command); err != nil {
		t.Error(err.Error())
	}

	if expected, actual := true, command.Success(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestExecuteConcurrentCommandsOnCluster(t *testing.T) {
	// NB: we're executing more commands than there are connections
	// so this tests command retries
	maxCount := uint16(8)
	nodeOpts := &NodeOptions{
		MinConnections: 1,
		MaxConnections: 4,
		RemoteAddress:  getRiakAddress(),
	}

	var node *Node
	var err error
	if node, err = NewNode(nodeOpts); err != nil {
		t.Error(err.Error())
	}
	if node == nil {
		t.FailNow()
	}

	nodes := []*Node{node}
	opts := &ClusterOptions{
		Nodes: nodes,
	}

	cluster, err := NewCluster(opts)
	if err != nil {
		t.Error(err.Error())
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	if err := cluster.Start(); err != nil {
		t.Error(err.Error())
	}

	pingChan := make(chan *PingCommand)
	for i := uint16(0); i < maxCount; i++ {
		go func() {
			command := &PingCommand{}
			if err := cluster.Execute(command); err != nil {
				t.Error(err.Error())
			}
			pingChan <- command
		}()
	}

	j := uint16(0)
	for i := uint16(0); i < maxCount; i++ {
		pingCommand := <-pingChan
		if expected, actual := true, pingCommand.Success(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		j++
	}
	if expected, actual := maxCount, j; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestExecuteConcurrentCommandsOnClusterWithMaxConnections(t *testing.T) {
	nodeOpts := &NodeOptions{
		RemoteAddress: getRiakAddress(),
	}

	var node *Node
	var err error
	if node, err = NewNode(nodeOpts); err != nil {
		t.Error(err.Error())
	}
	if node == nil {
		t.FailNow()
	}

	nodes := []*Node{node}
	opts := &ClusterOptions{
		Nodes: nodes,
	}

	cluster, err := NewCluster(opts)
	if err != nil {
		t.Error(err.Error())
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	if err := cluster.Start(); err != nil {
		t.Error(err.Error())
	}

	count := nodeOpts.MaxConnections
	pingChan := make(chan *PingCommand)
	for i := uint16(0); i < count; i++ {
		go func() {
			command := &PingCommand{}
			if err := cluster.Execute(command); err != nil {
				t.Error(err.Error())
			}
			pingChan <- command
		}()
	}

	j := uint16(0)
	for i := uint16(0); i < count; i++ {
		pingCommand := <-pingChan
		if expected, actual := true, pingCommand.Success(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		j++
	}
	if expected, actual := count, j; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestExecuteCommandThreeTimesOnDifferentNodes(t *testing.T) {
	nodeCount := 3
	listenerChan := make(chan bool, nodeCount)
	listeners := make([]*testListener, nodeCount)
	defer func() {
		for _, s := range listeners {
			s.stop()
		}
	}()

	nodes := make([]*Node, nodeCount)
	for i := 0; i < nodeCount; i++ {

		var onConn = func(c net.Conn) bool {
			handleClientMessageWithRiakError(t, c, 1, listenerChan)
			return true
		}
		o := &testListenerOpts{
			test:   t,
			onConn: onConn,
		}
		tl := newTestListener(o)
		tl.start()

		listeners[i] = tl

		nodeOptions := &NodeOptions{
			RemoteAddress:  tl.addr.String(),
			MinConnections: 0,
			MaxConnections: 1,
		}
		if node, err := NewNode(nodeOptions); err == nil {
			nodes[i] = node
		} else {
			t.Fatal(err)
		}
	}

	clusterOptions := &ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: 3,
	}
	cluster, err := NewCluster(clusterOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := cluster.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err := NewFetchValueCommandBuilder().
		WithBucket("b").
		WithKey("k").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	cluster.Execute(cmd)

	j := 0
	for j = 0; j < nodeCount; {
		select {
		case <-listenerChan:
			j++
		case <-time.After(5 * time.Second):
			t.Fatal("test timed out")
		}
	}
	if expected, actual := nodeCount, j; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestAsyncExecuteCommandOnCluster(t *testing.T) {
	nodeOpts := &NodeOptions{
		RemoteAddress: getRiakAddress(),
	}

	var node *Node
	var err error
	if node, err = NewNode(nodeOpts); err != nil {
		t.Fatal(err.Error())
	}
	if node == nil {
		t.FailNow()
	}

	nodes := []*Node{node}
	opts := &ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: 3,
	}

	cluster, err := NewCluster(opts)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	if err := cluster.Start(); err != nil {
		t.Fatal(err.Error())
	}

	command := &PingCommand{}
	args := &Async{
		Command: command,
		Done:    make(chan Command),
	}
	if err := cluster.ExecuteAsync(args); err != nil {
		t.Fatal(err.Error())
	}

	done := <-args.Done
	pingDone := done.(*PingCommand)

	if expected, actual := true, command == pingDone; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := true, command.Success(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	command = &PingCommand{}
	wg := &sync.WaitGroup{}
	args = &Async{
		Command: command,
		Wait:    wg,
	}
	if err := cluster.ExecuteAsync(args); err != nil {
		t.Fatal(err.Error())
	}

	wg.Wait()
	if expected, actual := true, command.Success(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestEnqueueCommandsAndRetryFromQueue(t *testing.T) {
	pingCommandCount := uint16(8)
	stateChan := make(chan state)
	doneChan := make(chan bool)
	var node *Node

	go func() {
		pingCommands := make([]*PingCommand, pingCommandCount)

		var err error
		nodeOpts := &NodeOptions{
			RemoteAddress:  "127.0.0.1:13339",
			MinConnections: 0,
		}
		node, err = NewNode(nodeOpts)
		if err != nil {
			t.Fatal(err)
		}
		if node == nil {
			t.FailNow()
		}
		origNodeSetStateFunc := node.setStateFunc
		node.setStateFunc = func(sd *stateData, st state) {
			origNodeSetStateFunc(&node.stateData, st)
			logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "sending state '%v' down stateChan", st)
			stateChan <- st
			logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "sent state '%v' down stateChan", st)
		}
		nodes := []*Node{node}
		clusterOpts := &ClusterOptions{
			Nodes:             nodes,
			ExecutionAttempts: 3,
			QueueMaxDepth:     pingCommandCount,
		}
		cluster, err := NewCluster(clusterOpts)
		if err != nil {
			t.Fatal(err.Error())
		}
		if err := cluster.Start(); err != nil {
			t.Fatal(err.Error())
		}
		wg := &sync.WaitGroup{}
		for i := uint16(0); i < pingCommandCount; i++ {
			ping := &PingCommand{}
			pingCommands[i] = ping
			args := &Async{
				Command: ping,
				Wait:    wg,
			}
			if err := cluster.ExecuteAsync(args); err != nil {
				t.Error(err)
			}
		}
		wg.Wait()
		node.setStateFunc = origNodeSetStateFunc
		close(stateChan)
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
		for _, pc := range pingCommands {
			if expected, actual := true, pc.Success(); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}

		doneChan <- true
	}()

	go func() {
		o := &testListenerOpts{
			test: t,
			host: "127.0.0.1",
			port: 13339,
		}
		tl := newTestListener(o)
		defer tl.stop()

		listenerStarted := false
		for {
			logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "waiting for state on stateChan")
			if nodeState, ok := <-stateChan; ok {
				logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "got nodeState: '%v'", nodeState)
				if !listenerStarted && node.isCurrentState(nodeHealthChecking) {
					logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "starting listener")
					listenerStarted = true
					tl.start()
					logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "listener is started")
				}
			} else {
				logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "stateChan CLOSED")
				break
			}
		}
	}()

	logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "WAITING")
	done := <-doneChan
	logDebug("[TestEnqueueCommandsAndRetryFromQueue]", "DONE: %v", done)
}

func TestRecoverFromReadTimeout(t *testing.T) {
	var connects int32 = 0
	var nc int32 = 2

	var onConn = func(c net.Conn) bool {
		var j int32 = atomic.AddInt32(&connects, 1)
		if j%nc == 0 {
			time.Sleep(time.Second * 1)
		}
		if readWriteResp(t, c, false) {
			return false // connection is not done
		} else {
			return true // close connection
		}
	}

	nodes := make([]*Node, 2)

	for i := 0; i < int(nc); i++ {
		o := &testListenerOpts{
			test:   t,
			onConn: onConn,
		}
		tl := newTestListener(o)
		defer tl.stop()
		tl.start()

		nodeOpts := &NodeOptions{
			MinConnections:      1,
			MaxConnections:      2,
			TempNetErrorRetries: 8,
			RequestTimeout:      time.Millisecond * 100,
			RemoteAddress:       tl.addr.String(),
		}
		var node *Node
		var err error
		if node, err = NewNode(nodeOpts); err != nil {
			t.Fatal(err)
		}
		if node == nil {
			t.Fatal()
		}
		nodes[i] = node
	}

	opts := &ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: 3,
	}

	cluster, err := NewCluster(opts)
	if err != nil {
		t.Error(err.Error())
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	if err := cluster.Start(); err != nil {
		t.Error(err.Error())
	}

	for i := 0; i < 12; i++ {
		var cmd Command
		if i%2 == 0 {
			cmd = &PingCommand{}
		} else {
			cmd = &GetServerInfoCommand{}
		}
		if err := cluster.Execute(cmd); err != nil {
			t.Error(err.Error())
		}
	}
}
