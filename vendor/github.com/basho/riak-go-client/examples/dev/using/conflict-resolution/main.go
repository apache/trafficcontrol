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
	"fmt"

	riak "github.com/basho/riak-go-client"
)

/*
    Code samples from:
    http://docs.basho.com/riak/latest/dev/using/conflict-resolution/golang/

	make sure this bucket-type is created:
	siblings

	riak-admin bucket-type create siblings '{"props":{"allow_mult":true}}'
	riak-admin bucket-type activate siblings
*/

func main() {
	//riak.EnableDebugLogging = true

	nodeOpts := &riak.NodeOptions{
		RemoteAddress: "riak-test:10017",
	}

	var node *riak.Node
	var err error
	if node, err = riak.NewNode(nodeOpts); err != nil {
		fmt.Println(err.Error())
	}

	nodes := []*riak.Node{node}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}

	cluster, err := riak.NewCluster(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if err = cluster.Stop(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	if err = cluster.Start(); err != nil {
		fmt.Println(err.Error())
	}

	// ping
	ping := &riak.PingCommand{}
	if err = cluster.Execute(ping); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ping passed")
	}

	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("Ren"),
	}

	storeObject(cluster, obj)
	storeObject(cluster, obj)
	readSiblings(cluster)

	resolveViaOverwrite(cluster)
	readSiblings(cluster)

	storeObject(cluster, obj)
	storeObject(cluster, obj)
	readSiblings(cluster)

	resolveChoosingFirst(cluster)
	readSiblings(cluster)

	storeObject(cluster, obj)
	storeObject(cluster, obj)
	readSiblings(cluster)

	resolveUsingResolver(cluster)
	readSiblings(cluster)
}

func storeObject(cluster *riak.Cluster, obj *riak.Object) error {
	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("nickelodeon").
		WithKey("best_character").
		WithContent(obj).
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}

func readSiblings(cluster *riak.Cluster) error {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("nickelodeon").
		WithKey("best_character").
		Build()
	if err != nil {
		return err
	}

	err = cluster.Execute(cmd)
	if err != nil {
		return err
	}

	fcmd := cmd.(*riak.FetchValueCommand)
	fmt.Printf("[DevUsingConflictRes] nickelodeon/best_character has '%v' siblings\n", len(fcmd.Response.Values))

	return nil
}

func resolveViaOverwrite(cluster *riak.Cluster) error {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("nickelodeon").
		WithKey("best_character").
		Build()
	if err != nil {
		return err
	}

	err = cluster.Execute(cmd)
	if err != nil {
		return err
	}

	fcmd := cmd.(*riak.FetchValueCommand)
	obj := fcmd.Response.Values[0]
	// This overwrites the value and provides the canonical one
	obj.Value = []byte("Stimpy")

	return storeObject(cluster, obj)
}

func resolveChoosingFirst(cluster *riak.Cluster) error {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("nickelodeon").
		WithKey("best_character").
		Build()
	if err != nil {
		return err
	}

	err = cluster.Execute(cmd)
	if err != nil {
		return err
	}

	fcmd := cmd.(*riak.FetchValueCommand)
	obj := fcmd.Response.Values[0]

	return storeObject(cluster, obj)
}

type FirstSiblingResolver struct {
}

func (cr *FirstSiblingResolver) Resolve(objs []*riak.Object) []*riak.Object {
	// return the first one
	return []*riak.Object{
		objs[0],
	}
}

func resolveUsingResolver(cluster *riak.Cluster) error {
	// Note: a more sophisticated resolver would
	// look into the objects to pick one, or perhaps
	// present the list to a user to choose
	cr := &FirstSiblingResolver{}

	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("nickelodeon").
		WithKey("best_character").
		WithConflictResolver(cr).
		Build()
	if err != nil {
		return err
	}

	err = cluster.Execute(cmd)
	if err != nil {
		return err
	}

	fcmd := cmd.(*riak.FetchValueCommand)

	// Test that the resolver just returned one riak.Object
	vlen := len(fcmd.Response.Values)
	if vlen != 1 {
		return fmt.Errorf("expected 1 object, got %v", vlen)
	}

	obj := fcmd.Response.Values[0]
	return storeObject(cluster, obj)
}
