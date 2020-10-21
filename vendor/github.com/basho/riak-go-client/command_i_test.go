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
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	var (
		addr *net.TCPAddr
		err  error
		conn *connection
	)
	addr, err = net.ResolveTCPAddr("tcp4", getRiakAddress())
	if err != nil {
		t.Error(err.Error())
	}
	opts := &connectionOptions{
		remoteAddress:  addr,
		connectTimeout: time.Second * 5,
		requestTimeout: time.Millisecond * 500,
	}
	if conn, err = newConnection(opts); err == nil {
		if err = conn.connect(); err == nil {
			cmd := &PingCommand{}
			if expected, actual := false, conn.inFlight; expected != actual {
				t.Errorf("expected %v, got: %v", expected, actual)
			}
			if err = conn.execute(cmd); err == nil {
				if cmd.Success() != true {
					t.Error("ping did not return true")
				}
			}
		}
	}
	if err != nil {
		t.Error(err.Error())
	}
}
