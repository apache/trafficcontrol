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
	"reflect"
	"testing"
)

func TestNewClientWithPort(t *testing.T) {
	ports := []uint16{1234, 5678}
	for _, p := range ports {
		o := &testListenerOpts{
			test: t,
			host: "127.0.0.1",
			port: p,
		}
		tl := newTestListener(o)
		tl.start()
		defer tl.stop()
	}

	opts := &NewClientOptions{
		Port: 1234,
		RemoteAddresses: []string{
			"127.0.0.1",
			"127.0.0.1:5678",
			"127.0.0.1",
		},
	}
	c, err := NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
	if err != nil {
		t.Error(err)
	}
	if expected, actual := true, reflect.DeepEqual(addr, c.cluster.nodes[0].addr); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	addr, err = net.ResolveTCPAddr("tcp", "127.0.0.1:5678")
	if err != nil {
		t.Error(err)
	}
	if expected, actual := true, reflect.DeepEqual(addr, c.cluster.nodes[1].addr); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	addr, err = net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
	if err != nil {
		t.Error(err)
	}
	if expected, actual := true, reflect.DeepEqual(addr, c.cluster.nodes[2].addr); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}
