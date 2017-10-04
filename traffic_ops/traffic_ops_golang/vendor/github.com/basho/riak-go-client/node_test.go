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

func TestCreateNodeWithOptions(t *testing.T) {
	builder := &PingCommandBuilder{}
	opts := &NodeOptions{
		RemoteAddress:       "8.8.8.8:1234",
		MinConnections:      2,
		MaxConnections:      2048,
		IdleTimeout:         tenSeconds,
		ConnectTimeout:      tenSeconds,
		RequestTimeout:      tenSeconds,
		HealthCheckInterval: tenSeconds,
		HealthCheckBuilder:  builder,
		TempNetErrorRetries: 16,
	}
	node, err := NewNode(opts)
	if err != nil {
		t.Error(err.Error())
	}
	if expected, actual := "8.8.8.8:1234|0|0", node.String(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := nodeCreated, node.getState(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if node.addr.Port != 1234 {
		t.Errorf("expected port 1234, got: %s", string(node.addr.Port))
	}
	if node.addr.Zone != "" {
		t.Errorf("expected empty zone, got: %s", string(node.addr.Zone))
	}
	var testIP = net.ParseIP("8.8.8.8")
	if !node.addr.IP.Equal(testIP) {
		t.Errorf("expected %v, got: %v", testIP, node.addr.IP)
	}
	if expected, actual := node.cm.minConnections, opts.MinConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.maxConnections, opts.MaxConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.idleTimeout, opts.IdleTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.connectTimeout, opts.ConnectTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.requestTimeout, opts.RequestTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if got, want := node.cm.tempNetErrorRetries, opts.TempNetErrorRetries; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if expected, actual := node.healthCheckInterval, opts.HealthCheckInterval; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := builder, node.healthCheckBuilder; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
}

func TestEnsureNodeValuesWithZeroValOptions(t *testing.T) {
	opts := &NodeOptions{
		TempNetErrorRetries: 8,
	}
	node, err := NewNode(opts)
	if err != nil {
		t.Error(err.Error())
	}
	exp := fmt.Sprintf("%s|0|0", defaultRemoteAddress)
	if expected, actual := exp, node.String(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := nodeCreated, node.getState(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if node.addr.Port != int(defaultRemotePort) {
		t.Errorf("expected port %v, got: %v", defaultRemotePort, node.addr.Port)
	}
	if expected, actual := node.cm.minConnections, defaultMinConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.maxConnections, defaultMaxConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.idleTimeout, defaultIdleTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.connectTimeout, defaultConnectTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := node.cm.requestTimeout, defaultRequestTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if got, want := node.cm.tempNetErrorRetries, opts.TempNetErrorRetries; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if expected, actual := node.healthCheckInterval, defaultHealthCheckInterval; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
}

func TestEnsureDefaultNodeValues(t *testing.T) {
	node, err := NewNode(nil)
	if err != nil {
		t.Error(err.Error())
	}
	if node.addr.Port != 8087 {
		t.Errorf("expected port 8087, got: %v", string(node.addr.Port))
	}
	if node.addr.Zone != "" {
		t.Errorf("expected empty zone, got: %v", string(node.addr.Zone))
	}
	if !node.addr.IP.Equal(localhost) {
		t.Errorf("expected %v, got: %v", localhost, node.addr.IP)
	}
	if expected, actual := defaultMinConnections, node.cm.minConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := defaultMaxConnections, node.cm.maxConnections; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := defaultIdleTimeout, node.cm.idleTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := defaultConnectTimeout, node.cm.connectTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := defaultRequestTimeout, node.cm.requestTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if expected, actual := defaultConnectTimeout, node.cm.connectTimeout; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if got, want := node.cm.tempNetErrorRetries, defaultTempNetErrorRetries; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if expected, actual := defaultHealthCheckInterval, node.healthCheckInterval; expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
	if node.healthCheckBuilder != nil {
		t.Errorf("expected nil, got: %v", node.healthCheckBuilder)
	}
	if expected, actual := nodeCreated, node.getState(); expected != actual {
		t.Errorf("expected %v, got: %v", expected, actual)
	}
}
