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

/*
Package riak provides the interfaces needed to interact with Riak KV using
Protocol Buffers. Riak KV is a distributed key-value datastore designed to be
fault tolerant, scalable, and flexible.

Currently, this library was written for and tested against Riak KV 2.0+.

TL;DR;

	import (
		"fmt"
		"os"
		riak "github.com/basho/riak-go-client"
	)

	func main() {
		nodeOpts := &riak.NodeOptions{
			RemoteAddress: "127.0.0.1:8087",
		}

		var node *riak.Node
		var err error
		if node, err = riak.NewNode(nodeOpts); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		nodes := []*riak.Node{node}
		opts := &riak.ClusterOptions{
			Nodes: nodes,
		}

		cluster, err := riak.NewCluster(opts)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		defer func() {
			if err := cluster.Stop(); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}()

		if err := cluster.Start(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		obj := &riak.Object{
			ContentType:     "text/plain",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Value:           []byte("this is a value in Riak"),
		}

		cmd, err := riak.NewStoreValueCommandBuilder().
			WithBucket("testBucketName").
			WithContent(obj).
			Build()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if err := cluster.Execute(cmd); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		svc := cmd.(*riak.StoreValueCommand)
		rsp := svc.Response
		fmt.Println(rsp.GeneratedKey)
	}
*/
package riak
