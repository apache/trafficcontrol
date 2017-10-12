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
	"net"
	"testing"
)

func TestCreateConnection(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8098")
	if err != nil {
		t.Error(err.Error())
	}
	opts := &connectionOptions{
		remoteAddress:       addr,
		connectTimeout:      tenSeconds,
		requestTimeout:      tenSeconds,
		tempNetErrorRetries: 10,
	}
	var conn *connection
	if conn, err = newConnection(opts); err == nil {
		if conn.addr.Port != 8098 {
			t.Errorf("expected port 8098, got: %s", string(conn.addr.Port))
		}
		if conn.addr.Zone != "" {
			t.Errorf("expected empty zone, got: %s", string(conn.addr.Zone))
		}
		if !conn.addr.IP.Equal(localhost) {
			t.Errorf("expected %v, got: %v", localhost, conn.addr.IP)
		}
		if conn.connectTimeout != tenSeconds {
			t.Errorf("expected %v, got: %v", tenSeconds, conn.connectTimeout)
		}
		if conn.requestTimeout != tenSeconds {
			t.Errorf("expected %v, got: %v", tenSeconds, conn.requestTimeout)
		}
		if expected, actual := false, conn.inFlight; expected != actual {
			t.Errorf("expected %v, got: %v", expected, actual)
		}
		if got, want := conn.tempNetErrorRetries, opts.tempNetErrorRetries; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.Error(err.Error())
	}
}

func TestCreateConnectionWithBadAddress(t *testing.T) {
	_, err := net.ResolveTCPAddr("tcp4", "123456.89.9813948.19328419348:80983r6")
	if err == nil {
		t.Error("expected error")
	}
}

func TestCreateConnectionRequiresOptions(t *testing.T) {
	if _, err := newConnection(nil); err == nil {
		t.Error("expected error when creating Connection without options")
	}
}

func TestEnsureDefaultConnectionValues(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8087")
	if err != nil {
		t.Error(err.Error())
	}
	opts := &connectionOptions{remoteAddress: addr}
	var conn *connection
	if conn, err = newConnection(opts); err == nil {
		if conn.addr.Port != 8087 {
			t.Errorf("expected port 8087, got: %s", string(conn.addr.Port))
		}
		if conn.addr.Zone != "" {
			t.Errorf("expected empty zone, got: %s", string(conn.addr.Zone))
		}
		if !conn.addr.IP.Equal(localhost) {
			t.Errorf("expected %v, got: %v", localhost, conn.addr.IP)
		}
		if conn.connectTimeout != defaultConnectTimeout {
			t.Errorf("expected %v, got: %v", defaultConnectTimeout, conn.connectTimeout)
		}
		if conn.requestTimeout != defaultRequestTimeout {
			t.Errorf("expected %v, got: %v", defaultRequestTimeout, conn.requestTimeout)
		}
		if expected, actual := false, conn.inFlight; expected != actual {
			t.Errorf("expected %v, got: %v", expected, actual)
		}
	} else {
		t.Error(err.Error())
	}
}
