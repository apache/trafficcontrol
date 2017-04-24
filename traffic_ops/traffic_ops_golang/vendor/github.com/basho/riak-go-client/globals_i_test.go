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

// +build integration timeseries integration_hll

package riak

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"sync"
	"testing"

	rpb_riak "github.com/basho/riak-go-client/rpb/riak"
	proto "github.com/golang/protobuf/proto"
)

func integrationTestsBuildCluster() *Cluster {
	var cluster *Cluster
	var err error
	nodeOpts := &NodeOptions{
		RemoteAddress: getRiakAddress(),
	}
	var node *Node
	node, err = NewNode(nodeOpts)
	if err != nil {
		panic(fmt.Sprintf("error building integration test node object: %s", err.Error()))
	}
	if node == nil {
		panic("NewNode returned nil!")
	}
	nodes := []*Node{node}
	opts := &ClusterOptions{
		Nodes: nodes,
	}
	cluster, err = NewCluster(opts)
	if err != nil {
		panic(fmt.Sprintf("error building integration test cluster object: %s", err.Error()))
	}
	if err = cluster.Start(); err != nil {
		panic(fmt.Sprintf("error starting integration test cluster object: %s", err.Error()))
	}
	return cluster
}

type testListenerOpts struct {
	test   *testing.T
	host   string
	port   uint16
	onConn func(c net.Conn) bool
}

type testListener struct {
	test   *testing.T
	host   string
	port   uint16
	addr   net.Addr
	onConn func(c net.Conn) bool
	ln     net.Listener
}

func newTestListener(o *testListenerOpts) *testListener {
	if o.test == nil {
		panic("testing object is required")
	}
	if o.host == "" {
		o.host = "127.0.0.1"
	}
	if o.onConn == nil {
		o.onConn = func(c net.Conn) bool {
			if readWriteResp(o.test, c, false) {
				return false // connection is not done
			}
			return true // connection is done
		}
	}
	t := &testListener{
		test:   o.test,
		host:   o.host,
		port:   o.port,
		onConn: o.onConn,
	}
	if t.port > 0 {
		addrstr := net.JoinHostPort(t.host, strconv.Itoa(int(t.port)))
		if addr, err := net.ResolveTCPAddr("tcp4", addrstr); err != nil {
			t.test.Fatal(err)
		} else {
			t.addr = addr
		}
	}
	return t
}

func (t *testListener) start() {
	if t.test == nil {
		panic("testing object is required")
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	addr := net.JoinHostPort(t.host, strconv.Itoa(int(t.port)))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.test.Fatal(err)
	} else {
		t.ln = ln
		t.addr = ln.Addr()
		tcpaddr := t.addr.(*net.TCPAddr)
		t.port = uint16(tcpaddr.Port)
	}

	go func() {
		wg.Done()
		logDebug("[testListener]", "(%v) started", t.addr)
		for {
			c, err := t.ln.Accept()
			if err != nil {
				if _, ok := err.(*net.OpError); !ok {
					t.test.Log(err)
				}
				return
			}
			go func() {
				for {
					if t.onConn(c) {
						break
					}
				}
			}()
		}
	}()

	wg.Wait()
	return
}

func (t *testListener) stop() {
	if t.ln == nil {
		logDebugln("[testListener]", "never started!")
	} else {
		if err := t.ln.Close(); err != nil {
			t.test.Error(err)
		}
		logDebug("[testListener]", "(%v) stopped", t.addr)
	}
}

func readWriteResp(t *testing.T, c net.Conn, shouldClose bool) (success bool) {
	success = false
	var err error
	var msgCode byte

	if msgCode, err = readClientMessage(c); err != nil {
		if err == io.EOF {
			c.Close()
		} else {
			logErr("[testListener]", err)
			t.Error(err)
		}
		success = false
		return
	}

	var data []byte
	switch msgCode {
	case rpbCode_RpbPingReq:
		data = buildRiakMessage(rpbCode_RpbPingResp, nil)
	case rpbCode_RpbGetServerInfoReq:
		data, err = buildGetServerInfoResp()
	default:
		msg := fmt.Sprintf("unknown msg code: %v", msgCode)
		data, err = buildRiakError(msg)
	}

	if err != nil {
		t.Error(err)
		success = false
	}

	count, err := c.Write(data)
	if err == nil {
		success = true
	} else {
		t.Error(err)
		success = false
	}

	if count != len(data) {
		t.Errorf("expected to write %v bytes, wrote %v bytes", len(data), count)
		success = false
	}

	if shouldClose {
		c.Close()
	}

	return
}

// TODO this is copied from connection.go and should be shared
func readClientMessage(c net.Conn) (msgCode byte, err error) {
	var sizeBuf []byte = make([]byte, 4)
	var count int = 0
	if count, err = io.ReadFull(c, sizeBuf); err == nil && count == 4 {
		messageLength := binary.BigEndian.Uint32(sizeBuf)
		data := make([]byte, messageLength)
		count, err = io.ReadFull(c, data)
		if err != nil {
			return
		} else if uint32(count) != messageLength {
			err = fmt.Errorf("[readClientMessage] message length: %d, only read: %d", messageLength, count)
		}
		msgCode = data[0]
	} else {
		if err != io.EOF {
			err = errors.New(fmt.Sprintf("[readClientMessage] error reading command size into sizeBuf: count %d, err %s, errtype %v", count, err, reflect.TypeOf(err)))
		}
	}
	return
}

func handleClientMessageWithRiakError(t *testing.T, c net.Conn, msgCount uint16, respChan chan bool) {
	defer func() {
		if err := c.Close(); err != nil {
			t.Error(err)
		}
	}()

	for i := 0; i < int(msgCount); i++ {
		if _, err := readClientMessage(c); err != nil {
			t.Error(err)
		}

		data, err := buildRiakError("this is an error")
		if err != nil {
			t.Error(err)
		}

		count, err := c.Write(data)
		if err != nil {
			t.Error(err)
		}
		if count != len(data) {
			t.Errorf("expected to write %v bytes, wrote %v bytes", len(data), count)
		}
		if respChan != nil {
			respChan <- true
		}
	}
}

func buildGetServerInfoResp() ([]byte, error) {
	n := bytes.NewBufferString("golang-test")
	v := bytes.NewBufferString("9.9.9")
	rpb := &rpb_riak.RpbGetServerInfoResp{
		Node:          n.Bytes(),
		ServerVersion: v.Bytes(),
	}
	if encoded, err := proto.Marshal(rpb); err != nil {
		return nil, err
	} else {
		data := buildRiakMessage(rpbCode_RpbGetServerInfoResp, encoded)
		return data, nil
	}
}

func buildRiakError(errmsg string) ([]byte, error) {
	var errcode uint32 = 1
	emsg := bytes.NewBufferString(errmsg)
	rpbErr := &rpb_riak.RpbErrorResp{
		Errcode: &errcode,
		Errmsg:  emsg.Bytes(),
	}
	if encoded, err := proto.Marshal(rpbErr); err != nil {
		return nil, err
	} else {
		data := buildRiakMessage(rpbCode_RpbErrorResp, encoded)
		return data, nil
	}
}
