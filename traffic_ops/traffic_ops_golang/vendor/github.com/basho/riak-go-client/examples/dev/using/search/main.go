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
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	riak "github.com/basho/riak-go-client"
)

/*
   Code samples from:
   http://docs.basho.com/riak/latest/dev/using/search/

   make sure these bucket-types are created:
   'animals', 'quotes', 'sports', 'cars', 'users', 'n_val_of_5'
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

	if err = storeIndex(cluster); err != nil {
		ErrExit(err)
	}

	if err = storeBucketProperties(cluster); err != nil {
		ErrExit(err)
	}

	if err = storeObjects(cluster); err != nil {
		ErrExit(err)
	}

	time.Sleep(time.Millisecond * 1250)

	if err = doSearchRequest(cluster); err != nil {
		ErrExit(err)
	}

	if err = doAgeSearchRequest(cluster); err != nil {
		ErrExit(err)
	}

	if err = doAndSearchRequest(cluster); err != nil {
		ErrExit(err)
	}

	if err = doPaginatedSearchRequest(cluster); err != nil {
		ErrExit(err)
	}

	if err = deleteIndex(cluster); err != nil {
		ErrExit(err)
	}
}

func ErrExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func storeIndex(cluster *riak.Cluster) error {
	cmd, err := riak.NewStoreIndexCommandBuilder().
		WithIndexName("famous").
		WithSchemaName("_yz_default").
		WithTimeout(time.Second * 30).
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}

func storeBucketProperties(cluster *riak.Cluster) error {
	cmd, err := riak.NewStoreBucketPropsCommandBuilder().
		WithBucketType("animals").
		WithBucket("cats").
		WithSearchIndex("famous").
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}

func storeObjects(cluster *riak.Cluster) error {
	o1 := &riak.Object{
		Key:   "liono",
		Value: []byte("{\"name_s\":\"Lion-o\",\"age_i\":30,\"leader_b\":true}"),
	}
	o2 := &riak.Object{
		Key:   "cheetara",
		Value: []byte("{\"name_s\":\"Cheetara\",\"age_i\":30,\"leader_b\":false}"),
	}
	o3 := &riak.Object{
		Key:   "snarf",
		Value: []byte("{\"name_s\":\"Snarf\",\"age_i\":43,\"leader_b\":false}"),
	}
	o4 := &riak.Object{
		Key:   "panthro",
		Value: []byte("{\"name_s\":\"Panthro\",\"age_i\":36,\"leader_b\":false}"),
	}

	objs := [...]*riak.Object{o1, o2, o3, o4}

	wg := &sync.WaitGroup{}
	for _, obj := range objs {
		obj.ContentType = "application/json"
		obj.Charset = "utf-8"
		obj.ContentEncoding = "utf-8"

		cmd, err := riak.NewStoreValueCommandBuilder().
			WithBucketType("animals").
			WithBucket("cats").
			WithContent(obj).
			Build()
		if err != nil {
			return err
		}

		args := &riak.Async{
			Command: cmd,
			Wait:    wg,
		}
		if err = cluster.ExecuteAsync(args); err != nil {
			return err
		}
	}

	wg.Wait()

	return nil
}

func printDocs(cmd riak.Command, desc string) error {
	sc := cmd.(*riak.SearchCommand)
	if json, jerr := json.MarshalIndent(sc.Response.Docs, "", "  "); jerr != nil {
		return jerr
	} else {
		fmt.Println("------------------------------------------------------------------------")
		fmt.Println(desc)
		fmt.Println(string(json))
	}
	return nil
}

func doSearchRequest(cluster *riak.Cluster) error {
	cmd, err := riak.NewSearchCommandBuilder().
		WithIndexName("famous").
		WithQuery("name_s:Lion*").
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	if err = printDocs(cmd, "Search Request Documents:"); err != nil {
		return err
	}

	sc := cmd.(*riak.SearchCommand)
	doc := sc.Response.Docs[0] // NB: SearchDoc struct type

	cmd, err = riak.NewFetchValueCommandBuilder().
		WithBucketType(doc.BucketType).
		WithBucket(doc.Bucket).
		WithKey(doc.Key).
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	fc := cmd.(*riak.FetchValueCommand)
	if json, jerr := json.MarshalIndent(fc.Response, "", "  "); jerr != nil {
		return jerr
	} else {
		fmt.Println(string(json))
	}

	return nil
}

func doAgeSearchRequest(cluster *riak.Cluster) error {
	cmd, err := riak.NewSearchCommandBuilder().
		WithIndexName("famous").
		WithQuery("age_i:[30 TO *]").
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	return printDocs(cmd, "Age Search Documents:")
}

func doAndSearchRequest(cluster *riak.Cluster) error {
	cmd, err := riak.NewSearchCommandBuilder().
		WithIndexName("famous").
		WithQuery("leader_b:true AND age_i:[30 TO *]").
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	return printDocs(cmd, "AND Search Documents:")
}

func doPaginatedSearchRequest(cluster *riak.Cluster) error {
	rowsPerPage := uint32(2)
	page := uint32(2)
	start := rowsPerPage * (page - uint32(1))

	cmd, err := riak.NewSearchCommandBuilder().
		WithIndexName("famous").
		WithQuery("*:*").
		WithStart(start).
		WithNumRows(rowsPerPage).
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	return printDocs(cmd, "Paginated Search Documents:")
}

func deleteIndex(cluster *riak.Cluster) error {
	cmd, err := riak.NewStoreBucketPropsCommandBuilder().
		WithBucketType("animals").
		WithBucket("cats").
		WithSearchIndex("_dont_index_").
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	cmd, err = riak.NewDeleteIndexCommandBuilder().
		WithIndexName("famous").
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}
