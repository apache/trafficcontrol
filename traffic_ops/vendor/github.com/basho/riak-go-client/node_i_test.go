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
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestCreateNodeWithOptionsAndStart(t *testing.T) {
	o := &testListenerOpts{
		test: t,
	}
	tl := newTestListener(o)
	tl.start()
	defer tl.stop()

	count := uint16(16)
	opts := &NodeOptions{
		RemoteAddress:       tl.addr.String(),
		MinConnections:      count,
		MaxConnections:      count,
		IdleTimeout:         tenSeconds,
		ConnectTimeout:      tenSeconds,
		RequestTimeout:      tenSeconds,
		HealthCheckInterval: time.Millisecond * 500,
		HealthCheckBuilder:  &PingCommandBuilder{},
		TempNetErrorRetries: 128,
	}
	node, err := NewNode(opts)
	if err != nil {
		t.Error(err.Error())
	}
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	if node.addr.Port != int(tl.port) {
		t.Errorf("expected port %d, got: %d", tl.port, node.addr.Port)
	}
	if node.addr.Zone != "" {
		t.Errorf("expected empty zone, got: %s", string(node.addr.Zone))
	}
	if expected, actual := opts.MinConnections, node.cm.minConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := opts.MaxConnections, node.cm.maxConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := opts.IdleTimeout, node.cm.idleTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if err := node.start(); err != nil {
		t.Error(err)
	}
	var f = func(v interface{}) (bool, bool) {
		conn := v.(*connection)
		if conn == nil {
			t.Error("got unexpected nil value")
			return true, false
		}
		if expected, actual := int(tl.port), conn.addr.Port; expected != actual {
			t.Errorf("expected %d, got: %d", expected, actual)
		}
		if conn.addr.Zone != "" {
			t.Errorf("expected empty zone, got: %s", string(conn.addr.Zone))
		}
		if expected, actual := conn.connectTimeout, opts.ConnectTimeout; expected != actual {
			t.Errorf("expected %v, got: %v", expected, actual)
		}
		if expected, actual := conn.requestTimeout, opts.RequestTimeout; expected != actual {
			t.Errorf("expected %v, got: %v", expected, actual)
		}
		if got, want := conn.tempNetErrorRetries, opts.TempNetErrorRetries; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		return false, true
	}
	if err := node.cm.q.iterate(f); err != nil {
		t.Error(err)
	}
	if err := node.stop(); err != nil {
		t.Error(err)
	}
}

func TestRecoverViaDefaultPingHealthCheck(t *testing.T) {
	connects := uint32(0)
	var onConn = func(c net.Conn) bool {
		if atomic.AddUint32(&connects, 1) == 1 {
			logDebug("[TestRecoverViaDefaultPingHealthCheck]", "onConn, closing, connects: %v", connects)
			c.Close()
		} else {
			logDebug("[TestRecoverViaDefaultPingHealthCheck]", "onConn, readWriteResp, connects: %v", connects)
			readWriteResp(t, c, true)
		}
		return true
	}
	o := &testListenerOpts{
		test:   t,
		onConn: onConn,
	}
	tl := newTestListener(o)
	tl.start()
	defer tl.stop()

	doneChan := make(chan struct{})
	stateChan := make(chan state)

	go func() {
		opts := &NodeOptions{
			RemoteAddress:  tl.addr.String(),
			MinConnections: 0,
		}
		node, err := NewNode(opts)
		if err != nil {
			t.Error(err)
		}

		origSetStateFunc := node.setStateFunc
		node.setStateFunc = func(sd *stateData, st state) {
			origSetStateFunc(&node.stateData, st)
			logDebug("[TestRecoverViaDefaultPingHealthCheck]", "sending state '%v' down stateChan", st)
			stateChan <- st
		}

		node.start()

		pingFunc := func() {
			ping := &PingCommand{}
			executed, err := node.execute(ping)
			if executed == false {
				t.Error("expected ping to be executed")
			}
			if err != nil {
				t.Logf("ping err: %v", err.Error())
			}
		}

		go func() {
			pingFunc()
			for {
				select {
				case <-doneChan:
					logDebug("[TestRecoverViaDefaultPingHealthCheck]", "stopping pings")
					return
				case <-time.After(time.Millisecond * 500):
					pingFunc()
				}
			}
		}()

		for {
			select {
			case <-doneChan:
				logDebug("[TestRecoverViaDefaultPingHealthCheck]", "stopping node")
				node.stop()
				return
			case <-time.After(time.Second * 1):
				logDebug("[TestRecoverViaDefaultPingHealthCheck]", "still waiting to stop node...")
			}
		}
	}()

	checkStatesFunc := func(states []state) {
		idx := 0
		for {
			select {
			case nodeState := <-stateChan:
				if expected, actual := states[idx], nodeState; expected != actual {
					t.Errorf("expected %v, got %v", expected, actual)
				} else {
					logDebug("[TestRecoverViaDefaultPingHealthCheck]", "saw state %d", nodeState)
				}
				idx++
				if idx >= len(states) {
					return
				}
			case <-time.After(time.Second * 5):
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, true)
				t.Fatalf("[TestRecoverViaDefaultPingHealthCheck] timeout waiting for stateChan!\n%s", string(buf[0:stackSize]))
				return
			}
		}
	}

	expectedStates := []state{
		nodeRunning, nodeHealthChecking, nodeRunning,
	}

	checkStatesFunc(expectedStates)

	close(doneChan)

	expectedStates = []state{
		nodeShuttingDown, nodeShutdown,
	}

	checkStatesFunc(expectedStates)

	close(stateChan)
}

func TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck(t *testing.T) {
	o := &testListenerOpts{
		test: t,
		host: "127.0.0.1",
		port: 13338,
	}
	tl := newTestListener(o)
	defer tl.stop()

	stateChan := make(chan state)
	recoveredChan := make(chan struct{})

	var err error
	var node *Node

	go func() {
		listenerStarted := false
		nodeIsRunningCount := 0
		for {
			logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "waiting on stateChan...")
			if nodeState, ok := <-stateChan; ok {
				logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "received nodeState: '%v'", nodeState)
				if nodeState == nodeRunning {
					nodeIsRunningCount++
				}
				logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "nodeIsRunningCount: '%v'", nodeIsRunningCount)
				if nodeIsRunningCount == 2 {
					// This is the second time node has entered nodeRunning state, so it must have recovered via the healthcheck
					logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "SUCCESS node recovered via healthcheck")
					close(recoveredChan)
					break
				}
				if !listenerStarted && nodeState == nodeHealthChecking {
					logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "STARTING LISTENER")
					tl.start()
					listenerStarted = true
				}
			} else {
				t.Error("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck] stateChan closed before recovering via healthcheck")
				break
			}
		}
	}()

	opts := &NodeOptions{
		ConnectTimeout: 500 * time.Millisecond,
		RemoteAddress:  "127.0.0.1:13338", // NB: can't use tl.addr since it isn't set until tl.start()
	}
	node, err = NewNode(opts)
	if err != nil {
		t.Fatal(err)
	}
	origSetStateFunc := node.setStateFunc

	go func() {
		node.setStateFunc = func(sd *stateData, st state) {
			origSetStateFunc(&node.stateData, st)
			logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "SENDING state '%v' down stateChan", st)
			stateChan <- st
			logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "SENT state '%v' down stateChan", st)
		}
		node.start()

		pingFunc := func() {
			ping := &PingCommand{}
			if _, perr := node.execute(ping); perr != nil {
				t.Logf("ping err: %v", perr)
			}
		}

		pingFunc()
		for {
			select {
			case <-recoveredChan:
				logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "stopping pings")
				return
			case <-time.After(time.Millisecond * 500):
				pingFunc()
			}
		}
	}()

	select {
	case <-recoveredChan:
		logDebug("[TestRecoverAfterConnectionComesUpViaDefaultPingHealthCheck]", "recovered")
		node.setStateFunc = origSetStateFunc
		node.stop()
		close(stateChan)
	case <-time.After(10 * time.Second):
		t.Error("test timed out")
	}
}
