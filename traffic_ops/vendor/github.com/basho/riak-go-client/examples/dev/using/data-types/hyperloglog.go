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
	"os"
	"time"

	"errors"
	riak "github.com/basho/riak-go-client"
)

/*
   Code samples from:
   http://docs.basho.com/riak/kv/latest/developing/data-types/hyperloglogs/

   make sure these bucket-types are created:

   riak_admin bucket-type create hlls '{"props":{"datatype":"hll"}}'
*/

func main() {
	//riak.EnableDebugLogging = true

	nodeOpts := &riak.NodeOptions{
		RemoteAddress:  "riak-test:10017",
		RequestTimeout: time.Second * 60,
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

	if err = assertHyperloglogStartsEmpty(cluster); err != nil {
		ErrExit(err)
	}

	if err = updateHyperloglog(cluster); err != nil {
		ErrExit(err)
	}

	if err = fetchHyperloglog(cluster); err != nil {
		ErrExit(err)
	}
}

func ErrExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func assertHyperloglogStartsEmpty(cluster *riak.Cluster) error {
	var resp *riak.FetchHllResponse

	builder := riak.NewFetchHllCommandBuilder()
	cmd, err := builder.WithBucketType("hlls").
		WithBucket("hello").
		WithKey("darkness").
		Build()
	if err != nil {
		return err
	}
	if err = cluster.Execute(cmd); err != nil {
		return err
	}
	if fc, ok := cmd.(*riak.FetchHllCommand); ok {
		if fc.Response == nil {
			return errors.New("expected non-nil Response")
		}
		resp = fc.Response
	}

	fmt.Println("Hyperloglog cardinality: ", resp.Cardinality)
	fmt.Println("Hyperloglog isNotFound: ", resp.IsNotFound)
	return nil
}

func updateHyperloglog(cluster *riak.Cluster) error {
	adds := [][]byte{
		[]byte("Jokes"),
		[]byte("Are"),
		[]byte("Better"),
		[]byte("Explained"),
		[]byte("Jokes"),
	}

	builder := riak.NewUpdateHllCommandBuilder()
	cmd, err := builder.WithBucketType("hlls").
		WithBucket("hello").
		WithKey("darkness").
		WithAdditions(adds...).
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}

func fetchHyperloglog(cluster *riak.Cluster) error {
	var resp *riak.FetchHllResponse

	builder := riak.NewFetchHllCommandBuilder()
	cmd, err := builder.WithBucketType("hlls").
		WithBucket("hello").
		WithKey("darkness").
		Build()
	if err != nil {
		return err
	}
	if err = cluster.Execute(cmd); err != nil {
		return err
	}
	if fc, ok := cmd.(*riak.FetchHllCommand); ok {
		if fc.Response == nil {
			return errors.New("expected non-nil Response")
		}
		resp = fc.Response
	}

	// We added "Jokes" twice, but, remember, the algorithm only counts the
	// unique elements we've added to the data structure.
	fmt.Println("Hyperloglog cardinality: ", resp.Cardinality)
	return nil
}
