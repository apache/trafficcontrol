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

package riak

import (
	"reflect"
	"testing"
	"time"

	rpbRiak "github.com/basho/riak-go-client/rpb/riak"
	rpbRiakSCH "github.com/basho/riak-go-client/rpb/riak_search"
	rpbRiakYZ "github.com/basho/riak-go-client/rpb/riak_yokozuna"
)

// StoreIndex
// RpbYokozunaIndexPutReq

func TestBuildRpbYokozunaIndexPutReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewStoreIndexCommandBuilder().
		WithIndexName("indexName").
		WithSchemaName("indexName_schema").
		WithNVal(5).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakYZ.RpbYokozunaIndexPutReq); ok {
		index := req.Index
		if expected, actual := "indexName", string(index.GetName()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "indexName_schema", string(index.GetSchema()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(5), index.GetNVal(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbYokozunaIndexPutReq", ok, reflect.TypeOf(protobuf))
	}
}

// FetchIndex
// RpbYokozunaIndexGetReq
// RpbYokozunaIndexGetResp

func TestBuildRpbYokozunaIndexGetReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewFetchIndexCommandBuilder().
		WithIndexName("indexName")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakYZ.RpbYokozunaIndexGetReq); ok {
		if expected, actual := "indexName", string(req.GetName()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbYokozunaIndexGetReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbYokozunaIndexGetRespCorrectly(t *testing.T) {
	var nval uint32 = 9
	indexes := make([]*rpbRiakYZ.RpbYokozunaIndex, 1)
	indexes[0] = &rpbRiakYZ.RpbYokozunaIndex{
		Name:   []byte("indexName"),
		Schema: []byte("_yz_default"),
		NVal:   &nval,
	}
	resp := &rpbRiakYZ.RpbYokozunaIndexGetResp{Index: indexes}
	builder := NewFetchIndexCommandBuilder().WithIndexName("indexName")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cmd.onSuccess(resp); err != nil {
		t.Fatal(err.Error())
	} else {
		if fcmd, ok := cmd.(*FetchIndexCommand); ok {
			if expected, actual := 1, len(fcmd.Response); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			idx := fcmd.Response[0]
			if expected, actual := "indexName", idx.Name; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "_yz_default", idx.Schema; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := nval, idx.NVal; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("ok: %v - could not convert %v to FetchIndexCommand", ok, reflect.TypeOf(cmd))
		}
	}
}

// DeleteIndex
// RpbYokozunaIndexDeleteReq

func TestBuildRpbYokozunaIndexDeleteReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewDeleteIndexCommandBuilder().
		WithIndexName("indexName")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakYZ.RpbYokozunaIndexDeleteReq); ok {
		if expected, actual := "indexName", string(req.GetName()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbYokozunaIndexDeleteReq", ok, reflect.TypeOf(protobuf))
	}
}

// StoreSchema
// RpbYokozunaSchemaPutReq

func TestBuildRpbYokozunaSchemaPutReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewStoreSchemaCommandBuilder().
		WithSchemaName("schemaName").
		WithSchema("schema_xml")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakYZ.RpbYokozunaSchemaPutReq); ok {
		schema := req.Schema
		if expected, actual := "schemaName", string(schema.GetName()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "schema_xml", string(schema.GetContent()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbYokozunaSchemaPutReq", ok, reflect.TypeOf(protobuf))
	}
}

// FetchSchema
// RpbYokozunaSchemaGetReq
// RpbYokozunaSchemaGetResp

func TestBuildRpbYokozunaSchemaGetReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewFetchSchemaCommandBuilder().
		WithSchemaName("schemaName")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakYZ.RpbYokozunaSchemaGetReq); ok {
		if expected, actual := "schemaName", string(req.GetName()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbYokozunaSchemaGetReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbYokozunaSchemaGetRespCorrectly(t *testing.T) {
	schema := &rpbRiakYZ.RpbYokozunaSchema{
		Name:    []byte("schemaName"),
		Content: []byte("schema_xml"),
	}
	resp := &rpbRiakYZ.RpbYokozunaSchemaGetResp{Schema: schema}
	builder := NewFetchSchemaCommandBuilder().WithSchemaName("schemaName")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cmd.onSuccess(resp); err != nil {
		t.Fatal(err.Error())
	} else {
		if fcmd, ok := cmd.(*FetchSchemaCommand); ok {
			schema := fcmd.Response
			if expected, actual := "schemaName", schema.Name; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "schema_xml", schema.Content; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("ok: %v - could not convert %v to FetchSchemaCommand", ok, reflect.TypeOf(cmd))
		}
	}
}

// Search
// pbSearchQueryReq
// RpbSearchQueryResp

func TestBuildRpbSearchQueryReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewSearchCommandBuilder().
		WithIndexName("indexName").
		WithQuery("*:*").
		WithNumRows(128).
		WithStart(2).
		WithSortField("sortField").
		WithFilterQuery("filterQuery").
		WithDefaultField("defaultField").
		WithDefaultOperation("and").
		WithReturnFields("field1", "field2").
		WithPresort("score")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); ok {
		t.Errorf("got %v, want cmd %s to NOT implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakSCH.RpbSearchQueryReq); ok {
		if expected, actual := "indexName", string(req.GetIndex()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "*:*", string(req.GetQ()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(128), req.GetRows(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(2), req.GetStart(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "sortField", string(req.GetSort()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "filterQuery", string(req.GetFilter()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "defaultField", string(req.GetDf()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "and", string(req.GetOp()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		rf := req.GetFl()
		if expected, actual := 2, len(rf); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "field1", string(rf[0]); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "field2", string(rf[1]); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "score", string(req.GetPresort()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbSearchQueryReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbSearchQueryRespCorrectly(t *testing.T) {
	resp := &rpbRiakSCH.RpbSearchQueryResp{
		Docs: make([]*rpbRiakSCH.RpbSearchDoc, 1),
	}

	p1 := &rpbRiak.RpbPair{
		Key:   []byte("leader_b"),
		Value: []byte("true"),
	}
	p2 := &rpbRiak.RpbPair{
		Key:   []byte("age_i"),
		Value: []byte("30"),
	}
	p3 := &rpbRiak.RpbPair{
		Key:   []byte("_yz_id"),
		Value: []byte("id"),
	}
	p4 := &rpbRiak.RpbPair{
		Key:   []byte("nullValue"),
		Value: []byte("null"),
	}
	p5 := &rpbRiak.RpbPair{
		Key:   []byte("array"),
		Value: []byte("val_0"),
	}
	p6 := &rpbRiak.RpbPair{
		Key:   []byte("array"),
		Value: []byte("val_1"),
	}
	p7 := &rpbRiak.RpbPair{
		Key:   []byte("_yz_rk"),
		Value: []byte("key"),
	}
	p8 := &rpbRiak.RpbPair{
		Key:   []byte("_yz_rt"),
		Value: []byte("bucket_type"),
	}
	p9 := &rpbRiak.RpbPair{
		Key:   []byte("_yz_rb"),
		Value: []byte("bucket"),
	}
	p10 := &rpbRiak.RpbPair{
		Key:   []byte("score"),
		Value: []byte("2.23"),
	}
	doc := &rpbRiakSCH.RpbSearchDoc{
		Fields: []*rpbRiak.RpbPair{p1, p2, p3, p4, p5, p6, p7, p8, p9, p10},
	}
	maxScore := float32(1.123)
	numFound := uint32(1)
	resp.Docs[0] = doc
	resp.MaxScore = &maxScore
	resp.NumFound = &numFound

	builder := NewSearchCommandBuilder().WithIndexName("index").WithQuery("some solr query")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cmd.onSuccess(resp); err != nil {
		t.Fatal(err.Error())
	} else {
		if scmd, ok := cmd.(*SearchCommand); ok {
			r := scmd.Response // SearchResponse
			if expected, actual := uint32(1), r.NumFound; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := float32(1.123), r.MaxScore; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := 1, len(r.Docs); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			doc := r.Docs[0]
			if expected, actual := "true", doc.Fields["leader_b"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "30", doc.Fields["age_i"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "null", doc.Fields["nullValue"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "id", doc.Fields["_yz_id"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "id", doc.Id; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "key", doc.Fields["_yz_rk"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "key", doc.Key; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "bucket", doc.Fields["_yz_rb"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "bucket", doc.Bucket; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "bucket_type", doc.Fields["_yz_rt"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "bucket_type", doc.BucketType; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "val_0", doc.Fields["array"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "val_1", doc.Fields["array"][1]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "2.23", doc.Fields["score"][0]; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "2.23", doc.Score; expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("ok: %v - could not convert %v to SearchCommand", ok, reflect.TypeOf(cmd))
		}
	}
}
