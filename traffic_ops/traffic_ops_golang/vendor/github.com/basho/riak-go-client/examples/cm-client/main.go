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
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	riak "github.com/basho/riak-go-client"
	util "github.com/lukebakken/goutil"
)

var stopping bool = false
var nodeCount uint16 = 5
var minConnections uint16 = 10
var maxConnections uint16 = 256

var key int = 1
var data []byte

var fetchDataInterval time.Duration = time.Millisecond * 100
var storeDataInterval time.Duration = time.Millisecond * 100

func init() {
	c := 256
	b := make([]byte, c)
	_, err := crand.Read(b)
	if err != nil {
		util.ErrExit(err)
	}
	data = []byte(hex.EncodeToString(b)) // NB: bytes of utf-8 encoding

	rand.Seed(time.Now().Unix())
}

func main() {
	riak.EnableDebugLogging = true

	sc := make(chan struct{})

	port := 10017
	nodes := make([]*riak.Node, nodeCount)
	for i := uint16(0); i < nodeCount; i++ {
		addr := fmt.Sprintf("riak-test:%d", port)
		nodeOpts := &riak.NodeOptions{
			MinConnections: minConnections,
			MaxConnections: maxConnections,
			RemoteAddress:  addr,
		}
		if node, nerr := riak.NewNode(nodeOpts); nerr != nil {
			util.ErrExit(nerr)
		} else {
			if node == nil {
				util.ErrExit(errors.New("node is nil"))
			}
			util.LogDebug("[cm-client]", "node: %v", node)
			nodes[i] = node
		}
		port += 10
	}

	opts := &riak.ClusterOptions{
		Nodes:             nodes,
		ExecutionAttempts: 3,
	}

	c, err := riak.NewCluster(opts)
	if err != nil {
		util.ErrExit(err)
	}

	if serr := c.Start(); serr != nil {
		util.ErrExit(serr)
	}

	// ping
	ping := &riak.PingCommand{}
	if perr := c.Execute(ping); perr != nil {
		util.ErrExit(perr)
	} else {
		fmt.Println("ping passed")
	}

	sdc := make(chan riak.Command, minConnections)
	fdc := make(chan riak.Command, minConnections)

	defer func() {
		stopping = true
		close(sc)
		if serr := c.Stop(); serr != nil {
			util.ErrExit(serr)
		}
		close(sdc)
		close(fdc)
	}()

	go storeData(c, sc, sdc)
	go fetchData(c, sc, fdc)

	util.LogInfo("[cm-client]", "HIT ANY KEY TO STOP")
	bio := bufio.NewReader(os.Stdin)
	_, _, rerr := bio.ReadLine()
	if rerr != nil {
		util.LogErr("[GH-47]", rerr)
	}
}

func fetchData(c *riak.Cluster, sc chan struct{}, dc chan riak.Command) {
	tck := time.NewTicker(fetchDataInterval)
	defer func() {
		tck.Stop()
	}()

	util.LogDebug("[cm-client/FetchData]", "Starting worker process")
	defer util.LogDebug("[cm-client/FetchData]", "Stopped worker process")

	for !stopping {
		select {
		case <-sc:
			util.LogDebug("[cm-client/FetchData]", "Stopping worker process")
			stopping = true
			break
		case cmd := <-dc:
			util.LogDebug("[cm-client/FetchData]", "%v completed", cmd.Name())
		case <-tck.C:
			for i := uint16(0); i < (minConnections * nodeCount); i++ {
				if stopping {
					break
				}
				rkey := rand.Intn(key)
				svc, err := riak.NewFetchValueCommandBuilder().
					WithBucket("chaos-monkey").
					WithKey(strconv.Itoa(rkey)).
					Build()
				if err != nil {
					util.ErrExit(err)
				}
				a := &riak.Async{
					Command: svc,
					Done:    dc,
				}
				if err = c.ExecuteAsync(a); err != nil {
					util.LogErr("[cm-client/FetchData]", err)
				}
			}
		}
	}
}

func storeData(c *riak.Cluster, sc chan struct{}, dc chan riak.Command) {
	tck := time.NewTicker(storeDataInterval)
	defer func() {
		tck.Stop()
	}()

	util.LogDebug("[cm-client/StoreData]", "Starting worker process")
	defer util.LogDebug("[cm-client/StoreData]", "Stopped worker process")

	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           data,
	}

	for !stopping {
		select {
		case <-sc:
			util.LogDebug("[cm-client/StoreData]", "Stopping worker process")
			stopping = true
			break
		case cmd := <-dc:
			util.LogDebug("[cm-client/StoreData]", "%v completed", cmd.Name())
		case <-tck.C:
			for i := uint16(0); i < (minConnections * nodeCount); i++ {
				if stopping {
					break
				}
				svc, err := riak.NewStoreValueCommandBuilder().
					WithBucket("chaos-monkey").
					WithKey(strconv.Itoa(key)).
					WithContent(obj).
					Build()
				if err != nil {
					util.ErrExit(err)
				}
				a := &riak.Async{
					Command: svc,
					Done:    dc,
				}
				if err = c.ExecuteAsync(a); err != nil {
					util.LogErr("[cm-client/StoreData]", err)
				}
				key++
			}
		}
	}
}
