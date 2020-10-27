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
	"errors"
	"fmt"
	"os"
	"sync"

	riak "github.com/basho/riak-go-client"
)

/*
   Code samples from:
   http://docs.basho.com/riak/latest/dev/using/2i/

   make sure the 'indexes' bucket-type is created using the leveldb backend
*/

func main() {
	riak.EnableDebugLogging = false

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
		Nodes:             nodes,
		ExecutionAttempts: 1,
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

	if err = insertingObjects(cluster); err != nil {
		ErrExit(err)
	}

	if err = queryingIndexes(cluster); err != nil {
		ErrExit(err)
	}

	if err = indexingObjects(cluster); err != nil {
		ErrExit(err)
	}

	if err = invalidFieldNames(cluster); err != nil {
		ErrExit(err)
	}

	if err = incorrectDataType(cluster); err != nil {
		ErrExit(err)
	}

	if err = queryingExactMatch(cluster); err != nil {
		ErrExit(err)
	}

	if err = queryingRange(cluster); err != nil {
		ErrExit(err)
	}

	if err = queryingRangeWithTerms(cluster); err != nil {
		ErrExit(err)
	}

	if err = queryingPagination(cluster); err != nil {
		ErrExit(err)
	}
}

func ErrExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func printIndexQueryResults(cmd riak.Command) {
	sciq := cmd.(*riak.SecondaryIndexQueryCommand)
	if sciq.Response == nil {
		fmt.Println("[DevUsing2i] print query results: no response")
		return
	}

	if sciq.Response.Results == nil {
		fmt.Println("[DevUsing2i] print query results: no results")
		return
	}

	for _, r := range sciq.Response.Results {
		fmt.Println("[DevUsing2i] index key:", string(r.IndexKey), "object key:", string(r.ObjectKey))
	}
}

func insertingObjects(cluster *riak.Cluster) error {
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		BucketType:      "indexes",
		Bucket:          "users",
		Key:             "john_smith",
		Value:           []byte("…user data…"),
	}

	obj.AddToIndex("twitter_bin", "jsmith123")
	obj.AddToIndex("email_bin", "jsmith@basho.com")

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithContent(obj).
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	return nil
}

func queryingIndexes(cluster *riak.Cluster) error {
	cmd, err := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("users").
		WithIndexName("twitter_bin").
		WithIndexKey("jsmith123").
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	printIndexQueryResults(cmd)
	return nil
}

func indexingObjects(cluster *riak.Cluster) error {
	o1 := &riak.Object{
		Key:   "larry",
		Value: []byte("My name is Larry"),
	}
	o1.AddToIndex("field1_bin", "val1")
	o1.AddToIntIndex("field2_int", 1001)

	o2 := &riak.Object{
		Key:   "moe",
		Value: []byte("My name is Moe"),
	}
	o2.AddToIndex("Field1_bin", "val2")
	o2.AddToIntIndex("Field2_int", 1002)

	o3 := &riak.Object{
		Key:   "curly",
		Value: []byte("My name is Curly"),
	}
	o3.AddToIndex("FIELD1_BIN", "val3")
	o3.AddToIntIndex("FIELD2_INT", 1003)

	o4 := &riak.Object{
		Key:   "veronica",
		Value: []byte("My name is Veronica"),
	}
	o4.AddToIndex("FIELD1_bin", "val4")
	o4.AddToIndex("FIELD1_bin", "val4")
	o4.AddToIndex("FIELD1_bin", "val4a")
	o4.AddToIndex("FIELD1_bin", "val4b")
	o4.AddToIntIndex("FIELD2_int", 1004)
	o4.AddToIntIndex("FIELD2_int", 1005)
	o4.AddToIntIndex("FIELD2_int", 1006)
	o4.AddToIntIndex("FIELD2_int", 1004)
	o4.AddToIntIndex("FIELD2_int", 1004)
	o4.AddToIntIndex("FIELD2_int", 1007)

	objs := [...]*riak.Object{o1, o2, o3, o4}

	wg := &sync.WaitGroup{}
	for _, obj := range objs {
		obj.ContentType = "text/plain"
		obj.Charset = "utf-8"
		obj.ContentEncoding = "utf-8"

		cmd, err := riak.NewStoreValueCommandBuilder().
			WithBucketType("indexes").
			WithBucket("people").
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

func invalidFieldNames(cluster *riak.Cluster) error {
	cmd, err := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("users").
		WithIndexName("field2_foo").
		WithIndexKey("jsmith123").
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println("[DevUsing2i] field name error:", err)
	} else {
		return errors.New("[DevUsing2i] expected an error!")
	}

	return nil
}

func incorrectDataType(cluster *riak.Cluster) error {
	obj := &riak.Object{
		BucketType:      "indexes",
		Bucket:          "people",
		Key:             "larry",
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("My name is Larry"),
	}
	obj.AddToIndex("field2_int", "bar")

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithContent(obj).
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		fmt.Println("[DevUsing2i] index data type error:", err)
	} else {
		return errors.New("[DevUsing2i] expected an error!")
	}

	return nil
}

func queryingExactMatch(cluster *riak.Cluster) error {
	c1, c1err := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("people").
		WithIndexName("field1_bin").
		WithIndexKey("val1").
		Build()
	if c1err != nil {
		return c1err
	}

	c2, c2err := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("people").
		WithIndexName("field2_int").
		WithIntIndexKey(1001).
		Build()
	if c2err != nil {
		return c2err
	}

	wg := &sync.WaitGroup{}
	cmds := [...]riak.Command{c1, c2}

	for _, cmd := range cmds {
		args := &riak.Async{
			Command: cmd,
			Wait:    wg,
		}
		if err := cluster.ExecuteAsync(args); err != nil {
			return err
		}
	}

	wg.Wait()

	for _, cmd := range cmds {
		printIndexQueryResults(cmd)
	}

	return nil
}

func queryingRange(cluster *riak.Cluster) error {
	c1, c1err := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("people").
		WithIndexName("field1_bin").
		WithRange("val2", "val4").
		Build()
	if c1err != nil {
		return c1err
	}

	c2, c2err := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("people").
		WithIndexName("field2_int").
		WithIntRange(1002, 1004).
		Build()
	if c2err != nil {
		return c2err
	}

	wg := &sync.WaitGroup{}
	cmds := [...]riak.Command{c1, c2}

	for _, cmd := range cmds {
		args := &riak.Async{
			Command: cmd,
			Wait:    wg,
		}
		if err := cluster.ExecuteAsync(args); err != nil {
			return err
		}
	}

	wg.Wait()

	for _, cmd := range cmds {
		printIndexQueryResults(cmd)
	}

	return nil
}

func queryingRangeWithTerms(cluster *riak.Cluster) error {
	cmd, err := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("tweets").
		WithIndexName("hashtags_bin").
		WithRange("rock", "rocl").
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	printIndexQueryResults(cmd)
	return nil
}

func doPaginatedQuery(cluster *riak.Cluster, continuation []byte) error {
	builder := riak.NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("indexes").
		WithBucket("tweets").
		WithIndexName("hashtags_bin").
		WithRange("ri", "ru").
		WithMaxResults(5)

	if continuation != nil && len(continuation) > 0 {
		builder.WithContinuation(continuation)
	}

	cmd, err := builder.Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	printIndexQueryResults(cmd)

	sciq := cmd.(*riak.SecondaryIndexQueryCommand)
	if sciq.Response == nil {
		return errors.New("[DevUsing2i] expected response but did not get one")
	}

	rc := sciq.Response.Continuation
	if rc != nil && len(rc) > 0 {
		return doPaginatedQuery(cluster, sciq.Response.Continuation)
	}

	return nil
}

func queryingPagination(cluster *riak.Cluster) error {
	return doPaginatedQuery(cluster, nil)
}
