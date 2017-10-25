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
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestSuccessfulConnection(t *testing.T) {
	connChan := make(chan bool)

	var onConn = func(c net.Conn) bool {
		defer c.Close()
		connChan <- true
		return true
	}
	o := &testListenerOpts{
		test:   t,
		onConn: onConn,
	}
	tl := newTestListener(o)
	defer tl.stop()
	tl.start()

	opts := &connectionOptions{
		remoteAddress: tl.addr.(*net.TCPAddr),
	}

	conn, err := newConnection(opts)
	if err != nil {
		t.Error(err)
	}

	if err := conn.connect(); err != nil {
		t.Error(err)
	}

	sawConnection := <-connChan

	if err := conn.close(); err != nil {
		t.Error(err)
	}

	if !sawConnection {
		t.Error("did not connect")
	}
}

func TestConnectionClosed(t *testing.T) {
	var onConn = func(c net.Conn) bool {
		if err := c.Close(); err != nil {
			t.Error(err)
		}
		return true
	}
	o := &testListenerOpts{
		test:   t,
		onConn: onConn,
	}
	tl := newTestListener(o)
	tl.start()

	opts := &connectionOptions{
		remoteAddress: tl.addr.(*net.TCPAddr),
	}

	conn, err := newConnection(opts)
	if err != nil {
		t.Error(err)
	}

	if err := conn.connect(); err != nil {
		t.Error("unexpected error in connect", err)
	} else {
		tl.stop()
		cmd := &PingCommand{}
		if err := conn.execute(cmd); err != nil {
			if operr, ok := err.(*net.OpError); ok {
				t.Log("op error", operr, operr.Op)
			} else if err == io.EOF {
				t.Log("saw EOF")
			} else {
				t.Errorf("expected to see net.OpError or io.EOF, but got '%s' (type: %v)", err.Error(), reflect.TypeOf(err))
			}
		} else {
			t.Error("expected error in execute")
		}
	}
}

func TestConnectionTimeout(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp4", "10.255.255.1:65535")
	if err != nil {
		t.Error(err.Error())
	}

	opts := &connectionOptions{
		remoteAddress:  addr,
		connectTimeout: time.Millisecond * 150,
	}

	if conn, err := newConnection(opts); err == nil {
		if err := conn.connect(); err == nil {
			t.Error("expected to see timeout error")
		} else {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				t.Log("timeout error", neterr)
			} else if operr, ok := err.(*net.OpError); ok {
				t.Log("op error", operr)
			} else {
				t.Errorf("expected to see timeout error, but got '%s' (type: %v)", err.Error(), reflect.TypeOf(err))
			}
		}
	} else {
		t.Error(err)
	}
}

func TestConnectionSuccess(t *testing.T) {
	o := &testListenerOpts{
		test: t,
		host: "127.0.0.1",
		port: 1340,
	}
	tl := newTestListener(o)
	defer tl.stop()
	tl.start()

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:1340")
	if err != nil {
		t.Error(err.Error())
	}

	opts := &connectionOptions{
		remoteAddress:  addr,
		connectTimeout: tenSeconds,
	}

	if conn, err := newConnection(opts); err == nil {
		if err := conn.connect(); err != nil {
			t.Error("unexpected error:", err)
		}
	} else {
		t.Error("unexpected error:", err)
	}
}
