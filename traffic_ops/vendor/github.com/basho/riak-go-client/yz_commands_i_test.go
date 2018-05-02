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

// +build integration

package riak

import (
	"fmt"
	"testing"
	"time"
)

// FetchIndex
// StoreIndex

func TestStoreFetchAndDeleteAYokozunaIndex(t *testing.T) {
	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	var err error
	var cmd Command
	indexName := "indexName"
	sbuilder := NewStoreIndexCommandBuilder()
	cmd, err = sbuilder.WithIndexName(indexName).WithTimeout(time.Second * 30).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*StoreIndexCommand); ok {
		if expected, actual := true, scmd.Response; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}

	fbuilder := NewFetchIndexCommandBuilder()
	cmd, err = fbuilder.WithIndexName(indexName).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if fcmd, ok := cmd.(*FetchIndexCommand); ok {
		if fcmd.Response == nil {
			t.Errorf("expected non-nil Response")
		}
		idx := fcmd.Response[0]
		if expected, actual := indexName, idx.Name; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "_yz_default", idx.Schema; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(3), idx.NVal; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}

	dbuilder := NewDeleteIndexCommandBuilder()
	cmd, err = dbuilder.WithIndexName(indexName).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if dcmd, ok := cmd.(*DeleteIndexCommand); ok {
		if expected, actual := true, dcmd.Response; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// FetchSchema
// StoreSchema

func TestStoreFetchAndDeleteAYokozunaSchema(t *testing.T) {
	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	var err error
	var cmd Command
	defaultSchemaName := "_yz_default"
	schemaName := "schemaName"
	schemaXml := "dummy"

	fbuilder := NewFetchSchemaCommandBuilder()
	cmd, err = fbuilder.WithSchemaName(defaultSchemaName).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err)
	}

	fcmd := cmd.(*FetchSchemaCommand)
	if fcmd.Response == nil {
		t.Fatal("expected non-nil Response")
	}

	sch := fcmd.Response
	if expected, actual := defaultSchemaName, sch.Name; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	schemaXml = sch.Content

	sbuilder := NewStoreSchemaCommandBuilder()
	cmd, err = sbuilder.WithSchemaName(schemaName).WithSchema(schemaXml).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*StoreSchemaCommand); ok {
		if expected, actual := true, scmd.Response; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// Search

func TestSearchViaYokozunaIndex(t *testing.T) {
	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	var err error
	var cmd Command
	indexName := "myIndex"

	b1 := NewStoreIndexCommandBuilder()
	cmd, err = b1.WithIndexName(indexName).WithTimeout(time.Second * 20).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*StoreIndexCommand); ok {
		if expected, actual := true, scmd.Response; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}

	searchBucket := fmt.Sprintf("%s_search", testBucketName)
	b2 := NewStoreBucketPropsCommandBuilder()
	cmd, err = b2.WithBucket(searchBucket).WithSearchIndex(indexName).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}

	stuffToStore := [...]string{
		"{ \"content_s\":\"Alice was beginning to get very tired of sitting by her sister on the bank, and of having nothing to do: once or twice she had peeped into the book her sister was reading, but it had no pictures or conversations in it, 'and what is the use of a book,' thought Alice 'without pictures or conversation?'\"}",
		"{ \"content_s\":\"So she was considering in her own mind (as well as she could, for the hot day made her feel very sleepy and stupid), whether the pleasure of making a daisy-chain would be worth the trouble of getting up and picking the daisies, when suddenly a White Rabbit with pink eyes ran close by her.\", \"multi_ss\":[\"this\",\"that\"]}",
		"{ \"content_s\":\"The rabbit-hole went straight on like a tunnel for some way, and then dipped suddenly down, so suddenly that Alice had not a moment to think about stopping herself before she found herself falling down a very deep well.\"}",
	}

	for _, s := range stuffToStore {
		b3 := NewStoreValueCommandBuilder()
		obj := &Object{
			ContentType: "application/json", // NB: *very* important for extractor
			Value:       []byte(s),
		}
		cmd, err = b3.WithBucket(searchBucket).WithContent(obj).Build()
		if err != nil {
			t.Fatal(err.Error())
		}
		if err = cluster.Execute(cmd); err != nil {
			t.Fatal(err.Error())
		}
	}

	for count := 0; count < 10; count++ {
		time.Sleep(time.Millisecond * 500)
		b4 := NewSearchCommandBuilder()
		cmd, err = b4.WithIndexName(indexName).WithQuery("multi_ss:t*").Build()
		if err != nil {
			t.Fatal(err.Error())
		}
		if err = cluster.Execute(cmd); err != nil {
			t.Fatal(err.Error())
		}
		if scmd, ok := cmd.(*SearchCommand); ok {
			resp := scmd.Response
			if expected, actual := 1, len(resp.Docs); expected != actual {
				if count < 10 {
					t.Logf("expected %v, got %v - RETRYING", expected, actual)
				} else {
					t.Errorf("expected %v, got %v - DONE", expected, actual)
				}
			}
		} else {
			t.FailNow()
		}
	}
}
