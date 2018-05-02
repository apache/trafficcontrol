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
   http://docs.basho.com/riak/latest/dev/using/search/

   make sure these bucket-types are created:
   'animals', 'quotes', 'sports', 'cars', 'users', 'n_val_of_5'

   Simple example:
   for bt in animals sports quotes cars users n_val_of_5
   do
   	riak-admin bucket-type create $bt
   	riak-admin bucket-type activate $bt
   done
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

	fetchRufus(cluster)
	storeRufus(cluster)
	storeStray(cluster)
	fetchRufusWithR(cluster)
	storeQuote(cluster)
	storeAndUpdateSport(cluster)
	storeCar(cluster)
	storeUserThenDelete(cluster)
	fetchBucketProps(cluster)
}

func storeRufus(cluster *riak.Cluster) {
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("WOOF!"),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("animals").
		WithBucket("dogs").
		WithKey("rufus").
		WithContent(obj).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	svc := cmd.(*riak.StoreValueCommand)
	rsp := svc.Response
	fmt.Println(rsp.VClock)
}

func storeQuote(cluster *riak.Cluster) {
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("I have nothing to declare but my genius"),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("quotes").
		WithBucket("oscar_wilde").
		WithKey("genius").
		WithContent(obj).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	svc := cmd.(*riak.StoreValueCommand)
	rsp := svc.Response
	fmt.Println(rsp.VClock)
}

func storeCar(cluster *riak.Cluster) {
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("vroom"),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("cars").
		WithBucket("dodge").
		WithKey("viper").
		WithW(3).
		WithContent(obj).
		WithReturnBody(true).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func storeUserThenDelete(cluster *riak.Cluster) {
	obj := &riak.Object{
		ContentType:     "application/json",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("{'user':'data'}"),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("users").
		WithBucket("random_user_keys").
		WithContent(obj).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	svc := cmd.(*riak.StoreValueCommand)
	rsp := svc.Response
	fmt.Printf("Generated key: %v\n", rsp.GeneratedKey)

	cmd, err = riak.NewDeleteValueCommandBuilder().
		WithBucketType("users").
		WithBucket("random_user_keys").
		WithKey(rsp.GeneratedKey).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func storeAndUpdateSport(cluster *riak.Cluster) {
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("Washington Generals"),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("sports").
		WithBucket("nba").
		WithKey("champion").
		WithContent(obj).
		WithReturnBody(true).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	svc := cmd.(*riak.StoreValueCommand)
	rsp := svc.Response
	obj = rsp.Values[0]
	obj.Value = []byte("Harlem Globetrotters")

	cmd, err = riak.NewStoreValueCommandBuilder().
		WithBucketType("sports").
		WithBucket("nba").
		WithKey("champion").
		WithContent(obj).
		WithReturnBody(true).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	svc = cmd.(*riak.StoreValueCommand)
	rsp = svc.Response
	obj = rsp.Values[0]
	fmt.Printf("champion: %v\n", string(obj.Value))
}

func fetchRufus(cluster *riak.Cluster) {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucketType("animals").
		WithBucket("dogs").
		WithKey("rufus").
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response
	fmt.Println(rsp.IsNotFound)
}

func storeStray(cluster *riak.Cluster) {
	// There is not a default content type for an object, so it must be defined. For JSON content
	// types, please use 'application/json'
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("Found in alley."),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("animals").
		WithBucket("dogs").
		WithContent(obj).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	svc := cmd.(*riak.StoreValueCommand)
	rsp := svc.Response
	fmt.Println(rsp.GeneratedKey)
}

// https://godoc.org/github.com/basho/riak-go-client#FetchValueCommandBuilder.WithR
func fetchRufusWithR(cluster *riak.Cluster) {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucketType("animals").
		WithBucket("dogs").
		WithKey("rufus").
		WithR(3).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response
	fmt.Println(rsp.IsNotFound)
}

func fetchChampion(cluster *riak.Cluster) {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucketType("sports").
		WithBucket("nba").
		WithKey("champion").
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response
	var obj *riak.Object
	if len(rsp.Values) > 0 {
		obj = rsp.Values[0]
	} else {
		obj = &riak.Object{
			ContentType:     "text/plain",
			Charset:         "utf-8",
			ContentEncoding: "utf-8",
			Value:           nil,
		}
	}

	obj.Value = []byte("Harlem Globetrotters")

	cmd, err = riak.NewStoreValueCommandBuilder().
		WithBucketType("sports").
		WithBucket("nba").
		WithContent(obj).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func fetchBucketProps(cluster *riak.Cluster) {
	cmd, err := riak.NewFetchBucketPropsCommandBuilder().
		WithBucketType("n_val_of_5").
		WithBucket("any_bucket_name").
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	fbc := cmd.(*riak.FetchBucketPropsCommand)
	fmt.Println("bucket props:", fbc.Response)
}
