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

package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	riak "github.com/basho/riak-go-client"
)

var stopping bool = false
var nodeCount uint16 = 1
var minConnections uint16 = 10
var maxConnections uint16 = 256

var key uint64
var data []byte

var slog = log.New(os.Stdout, "", log.LstdFlags)
var elog = log.New(os.Stderr, "", log.LstdFlags)

func LogInfo(source, format string, v ...interface{}) {
	slog.Printf(fmt.Sprintf("[INFO] %s %s", source, format), v...)
}

func LogDebug(source, format string, v ...interface{}) {
	slog.Printf(fmt.Sprintf("[DEBUG] %s %s", source, format), v...)
}

func LogError(source, format string, v ...interface{}) {
	elog.Printf(fmt.Sprintf("[DEBUG] %s %s", source, format), v...)
}

func LogErr(source string, err error) {
	elog.Println("[ERROR]", source, err)
}

func ErrExit(err error) {
	LogErr("[GH-47]", err)
	os.Exit(1)
}

func init() {
	c := 524288
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		ErrExit(err)
	}
	data = []byte(hex.EncodeToString(b)) // NB: bytes of utf-8 encoding
}

func keepAlive(c *riak.Cluster, sc chan struct{}) {
	tck := time.NewTicker(time.Second * 1)
	dc := make(chan riak.Command, minConnections)

	defer func() {
		tck.Stop()
		close(dc)
	}()

	LogDebug("[GH-47/KeepAlive]", "Starting keepalive process")
	defer LogDebug("[GH-47/KeepAlive]", "Stopped keepalive process")

	for !stopping {
		select {
		case <-sc:
			LogDebug("[GH-47/KeepAlive]", "Stopping keepalive process")
			stopping = true
			break
		case pc := <-dc:
			LogDebug("[GH-47/KeepAlive]", "%v completed", pc.Name())
		case t := <-tck.C:
			LogDebug("[GH-47/KeepAlive]", "Running keepalive at %v", t)
			for i := uint16(0); i < minConnections; i++ {
				if stopping {
					break
				}
				go func() {
					cmd := &riak.PingCommand{}
					if err := c.Execute(cmd); err != nil {
						LogErr("[GH-47/KeepAlive]", err)
					}
				}()
				/*
					a := &riak.Async{
						Command: &riak.PingCommand{},
						Done: dc,
					}
					if err := c.ExecuteAsync(a); err != nil {
						LogErr("[GH-47/KeepAlive]", err)
					}
				*/
			}
		}
	}
}

func storeData(c *riak.Cluster, sc chan struct{}) {
	tck := time.NewTicker(time.Millisecond * 125)
	dc := make(chan riak.Command, minConnections)

	defer func() {
		tck.Stop()
		close(dc)
	}()

	LogDebug("[GH-47/StoreData]", "Starting worker process")
	defer LogDebug("[GH-47/StoreData]", "Stopped worker process")

	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           data,
	}

	for !stopping {
		select {
		case <-sc:
			LogDebug("[GH-47/StoreData]", "Stopping worker process")
			stopping = true
			break
		case cmd := <-dc:
			LogDebug("[GH-47/StoreData]", "%v completed", cmd.Name())
		case <-tck.C:
			for i := uint16(0); i < minConnections; i++ {
				if stopping {
					break
				}
				svc, err := riak.NewStoreValueCommandBuilder().
					WithBucket("gh-47").
					WithKey(strconv.FormatUint(key, 10)).
					WithContent(obj).
					Build()
				if err != nil {
					ErrExit(err)
				}
				a := &riak.Async{
					Command: svc,
					Done:    dc,
				}
				if err = c.ExecuteAsync(a); err != nil {
					LogErr("[GH-47/StoreData]", err)
				}
				key++
			}
		}
	}
}

func listKeys(c *riak.Cluster, sc chan struct{}) {
	tck := time.NewTicker(time.Millisecond * 500)
	dc := make(chan riak.Command, minConnections)

	defer func() {
		tck.Stop()
		close(dc)
	}()

	LogDebug("[GH-47/ListKeys]", "Starting worker process")
	defer LogDebug("[GH-47/ListKeys]", "Stopped worker process")

	for !stopping {
		select {
		case <-sc:
			LogDebug("[GH-47/ListKeys]", "Stopping worker process")
			stopping = true
			break
		case cmd := <-dc:
			lk := cmd.(*riak.ListKeysCommand)
			kc := 0
			if lk.Response != nil {
				kc = len(lk.Response.Keys)
			}
			LogDebug("[GH-47/ListKeys]", "%v completed, keys: %d", cmd.Name(), kc)
		case <-tck.C:
			svc, err := riak.NewListKeysCommandBuilder().
				WithBucket("gh-47").
				WithStreaming(false).
				Build()
			if err != nil {
				ErrExit(err)
			}
			a := &riak.Async{
				Command: svc,
				Done:    dc,
			}
			if err = c.ExecuteAsync(a); err != nil {
				LogErr("[GH-47/ListKeys]", err)
			}
		}
	}
}

func main() {
	riak.EnableDebugLogging = true

	sc := make(chan struct{})

	if len(os.Args) > 1 {
		if i, err := strconv.Atoi(os.Args[1]); err != nil {
			ErrExit(err)
		} else {
			nodeCount = uint16(i)
		}
	}
	LogInfo("[GH-47]", "Node count: %v", nodeCount)

	nodes := make([]*riak.Node, nodeCount)
	port := 10017

	for i := uint16(0); i < nodeCount; i++ {
		var node *riak.Node
		var err error
		nodeOpts := &riak.NodeOptions{
			MinConnections: minConnections,
			MaxConnections: maxConnections,
			RemoteAddress:  fmt.Sprintf("riak-test:%d", port),
		}
		if node, err = riak.NewNode(nodeOpts); err != nil {
			ErrExit(err)
		}
		if node == nil {
			ErrExit(errors.New("node is nil"))
		}
		LogDebug("[GH-47]", "node: %v", node)
		nodes[i] = node
		port += 10
	}

	opts := &riak.ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: 3,
	}

	c, err := riak.NewCluster(opts)
	if err != nil {
		ErrExit(err)
	}

	if err = c.Start(); err != nil {
		ErrExit(err)
	}

	go keepAlive(c, sc)
	go storeData(c, sc)
	go listKeys(c, sc)

	defer func() {
		stopping = true
		close(sc)
		if err = c.Stop(); err != nil {
			LogErr("[GH-47]", err)
		}
	}()

	LogInfo("[GH-47]", "HIT ANY KEY TO STOP")
	bio := bufio.NewReader(os.Stdin)
	_, _, rerr := bio.ReadLine()
	if rerr != nil {
		LogErr("[GH-47]", rerr)
	}
}
