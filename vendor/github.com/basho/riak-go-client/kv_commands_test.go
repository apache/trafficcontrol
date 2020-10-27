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
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	rpbRiak "github.com/basho/riak-go-client/rpb/riak"
	rpbRiakKV "github.com/basho/riak-go-client/rpb/riak_kv"
	proto "github.com/golang/protobuf/proto"
)

type testConflictResolver struct {
}

func (cr *testConflictResolver) Resolve(objs []*Object) []*Object {
	// return the first one
	return []*Object{
		objs[0],
	}
}

var resolver = &testConflictResolver{}

// FetchValue

func TestBuildRpbGetReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewFetchValueCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("key").
		WithR(3).
		WithPr(1).
		WithBasicQuorum(true).
		WithNotFoundOk(true).
		WithIfModified(vclockBytes).
		WithHeadOnly(true).
		WithReturnDeletedVClock(true).
		WithTimeout(time.Second * 20).
		WithSloppyQuorum(true).
		WithNVal(4)
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
	validateRpbGetReq(t, protobuf)
}

func validateRpbGetReq(t *testing.T, protobuf proto.Message) {
	if req, ok := protobuf.(*rpbRiakKV.RpbGetReq); ok {
		if expected, actual := "bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "bucket_name", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "key", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(3), req.GetR(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(1), req.GetPr(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetNotfoundOk(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 0, bytes.Compare(vclockBytes, req.GetIfModified()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetHead(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetDeletedvclock(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
		if expected, actual := true, req.GetSloppyQuorum(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := uint32(4), req.GetNVal(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbGetReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestBuildRpbGetReqCorrectlyWithDefaults(t *testing.T) {
	builder := NewFetchValueCommandBuilder().
		WithBucket("bucket_name").
		WithKey("key")
	cmd, err := builder.Build()

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.FailNow()
	}
	if req, ok := protobuf.(*rpbRiakKV.RpbGetReq); ok {
		if expected, actual := "default", string(req.Type); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "bucket_name", string(req.Bucket); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "key", string(req.Key); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if req.R != nil {
			t.Errorf("expected nil value")
		}
		if req.Pr != nil {
			t.Errorf("expected nil value")
		}
		if req.NotfoundOk != nil {
			t.Error("expected nil value")
		}
		if req.IfModified != nil {
			t.Errorf("expected nil value")
		}
		if req.Head != nil {
			t.Error("expected nil value")
		}
		if req.Deletedvclock != nil {
			t.Error("expected nil value")
		}
		if req.Timeout != nil {
			t.Errorf("expected nil value")
		}
		if req.SloppyQuorum != nil {
			t.Error("expected nil value")
		}
		if req.NVal != nil {
			t.Errorf("expected nil value")
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbGetReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbGetRespCorrectly(t *testing.T) {
	rpbContent := generateTestRpbContent("this is a value", "application/json")

	rpbGetResp := &rpbRiakKV.RpbGetResp{
		Content: []*rpbRiakKV.RpbContent{rpbContent},
		Vclock:  vclock.Bytes(),
	}

	builder := NewFetchValueCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("key").
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	cmd.onSuccess(rpbGetResp)
	if expected, actual := true, cmd.Success(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	if fetchValueCommand, ok := cmd.(*FetchValueCommand); ok {
		if fetchValueCommand.Response == nil {
			t.Fatal("unexpected nil object")
		}
		if expected, actual := true, fetchValueCommand.success; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 1, len(fetchValueCommand.Response.Values); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		riakObject := fetchValueCommand.Response.Values[0]
		if riakObject == nil {
			t.Fatal("unexpected nil object")
		}
		if expected, actual := "bucket_type", riakObject.BucketType; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "bucket_name", riakObject.Bucket; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "key", riakObject.Key; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "application/json", riakObject.ContentType; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "utf-8", riakObject.Charset; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "utf-8", riakObject.ContentEncoding; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "test-vtag", riakObject.VTag; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := time.Unix(1234, 123456789), riakObject.LastModified; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, riakObject.HasIndexes(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, riakObject.HasIndexes(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "golang@basho.com", riakObject.Indexes["email_bin"][0]; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, riakObject.HasUserMeta(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "golang@basho.com", riakObject.Indexes["email_bin"][0]; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "frazzle@basho.com", riakObject.Indexes["email_bin"][1]; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "metaKey1", riakObject.UserMeta[0].Key; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "metaValue1", riakObject.UserMeta[0].Value; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "b0", riakObject.Links[0].Bucket; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "k0", riakObject.Links[0].Key; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "t0", riakObject.Links[0].Tag; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "b1", riakObject.Links[1].Bucket; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "k1", riakObject.Links[1].Key; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "t1", riakObject.Links[1].Tag; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "vclock123456789", string(riakObject.VClock); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchValueCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestParseRpbGetRespWithSiblingsCorrectly(t *testing.T) {
	rpb1 := generateTestRpbContent("value_1", "text/plain")
	rpb2 := generateTestRpbContent("value_2", "text/plain")

	rpbGetResp := &rpbRiakKV.RpbGetResp{
		Content: []*rpbRiakKV.RpbContent{rpb1, rpb2},
		Vclock:  vclock.Bytes(),
	}

	builder := NewFetchValueCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("key").
		WithConflictResolver(resolver).
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	cmd.onSuccess(rpbGetResp)
	if actual, expected := cmd.Success(), true; actual != expected {
		t.Errorf("got %v, expected %v", actual, expected)
	}

	if fetchValueCommand, ok := cmd.(*FetchValueCommand); ok {
		if fetchValueCommand.Response == nil {
			t.Fatal("unexpected nil object")
		}
		if actual, expected := len(fetchValueCommand.Response.Values), 1; actual != expected {
			t.Errorf("got %v, expected %v", actual, expected)
		}
		ro := fetchValueCommand.Response.Values[0]
		if ro == nil {
			t.Fatal("unexpected nil object")
		}
		if actual, expected := string(ro.Value), "value_1"; actual != expected {
			t.Errorf("got %v, expected %v", actual, expected)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchValueCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestParseRpbGetRespWithoutContentCorrectly(t *testing.T) {
	builder := NewFetchValueCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("key").
		Build()
	if err != nil {
		t.Error(err.Error())
	}
	cmd.onSuccess(nil)
	if expected, actual := true, cmd.Success(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if fetchValueCommand, ok := cmd.(*FetchValueCommand); ok {
		if fetchValueCommand.Response == nil {
			t.Fatal("unexpected nil object")
		}
		if expected, actual := true, fetchValueCommand.Response.IsNotFound; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchValueCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestValidationOfRpbGetReqViaBuilder(t *testing.T) {
	// validate that Bucket is required
	builder := NewFetchValueCommandBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is required
	builder = NewFetchValueCommandBuilder()
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrKeyRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func generateTestRpbContent(value string, contentType string) (rpbContent *rpbRiakKV.RpbContent) {
	lastMod := uint32(1234)
	lastModUsecs := uint32(123456789)
	deleted := false

	rpbContent = &rpbRiakKV.RpbContent{
		Value:           []byte(value),
		ContentType:     []byte(contentType),
		Charset:         []byte("utf-8"),
		ContentEncoding: []byte("utf-8"),
		Vtag:            []byte("test-vtag"),
		Links:           make([]*rpbRiakKV.RpbLink, 2),
		LastMod:         &lastMod,
		LastModUsecs:    &lastModUsecs,
		Usermeta:        make([]*rpbRiak.RpbPair, 2),
		Indexes:         make([]*rpbRiak.RpbPair, 3),
		Deleted:         &deleted,
	}

	rpbContent.Links[0] = &rpbRiakKV.RpbLink{
		Bucket: []byte("b0"),
		Key:    []byte("k0"),
		Tag:    []byte("t0"),
	}
	rpbContent.Links[1] = &rpbRiakKV.RpbLink{
		Bucket: []byte("b1"),
		Key:    []byte("k1"),
		Tag:    []byte("t1"),
	}

	rpbContent.Usermeta[0] = &rpbRiak.RpbPair{
		Key:   []byte("metaKey1"),
		Value: []byte("metaValue1"),
	}
	rpbContent.Usermeta[1] = &rpbRiak.RpbPair{
		Key:   []byte("metaKey2"),
		Value: []byte("metaValue2"),
	}

	rpbContent.Indexes[0] = &rpbRiak.RpbPair{
		Key:   []byte("email_bin"),
		Value: []byte("golang@basho.com"),
	}
	rpbContent.Indexes[1] = &rpbRiak.RpbPair{
		Key:   []byte("email_bin"),
		Value: []byte("frazzle@basho.com"),
	}
	rpbContent.Indexes[2] = &rpbRiak.RpbPair{
		Key:   []byte("phone_bin"),
		Value: []byte("15551234567"),
	}

	return rpbContent
}

// StoreValue

func TestValidationOfRpbPutReqViaBuilder(t *testing.T) {
	// validate that Bucket is required
	builder := NewStoreValueCommandBuilder()
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is NOT required
	builder = NewStoreValueCommandBuilder()
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err != nil {
		t.Fatal("expected nil err since PUT requests can generate keys")
	}
}

func TestBuildRpbPutReqCorrectlyViaBuilder(t *testing.T) {
	value := "this is a value"
	userMeta := []*Pair{
		{"metaKey1", "metaValue1"},
		{"metaKey2", "metaValue2"},
	}
	links := []*Link{
		{"b0", "k0", "t0"},
		{"b1", "k1", "t1"},
	}
	ro := &Object{
		ContentType:     "application/json",
		ContentEncoding: "gzip",
		Charset:         "utf-8",
		UserMeta:        userMeta,
		Links:           links,
		Value:           []byte(value),
	}
	ro.AddToIndex("email_bin", "golang@basho.com")
	ro.AddToIndex("email_bin", "frazzle@basho.com")

	key := "key"
	builder := NewStoreValueCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey(key).
		WithW(3).
		WithPw(1).
		WithDw(2).
		WithNVal(3).
		WithVClock(vclockBytes).
		WithReturnHead(true).
		WithReturnBody(true).
		WithIfNotModified(true).
		WithIfNoneMatch(true).
		WithAsis(true).
		WithSloppyQuorum(true).
		WithTimeout(time.Second * 20).
		WithContent(ro)
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

	if req, ok := protobuf.(*rpbRiakKV.RpbPutReq); ok {
		if expected, actual := "bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "bucket_name", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := key, string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(3), req.GetW(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(1), req.GetPw(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(2), req.GetDw(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(3), req.GetNVal(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 0, bytes.Compare(vclockBytes, req.GetVclock()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetReturnHead(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetReturnBody(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetIfNotModified(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetIfNoneMatch(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetAsis(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, req.GetSloppyQuorum(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
		content := req.GetContent()
		if content == nil {
			t.Fatal("expected non-nil content")
		} else {
			if expected, actual := value, string(content.GetValue()); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "application/json", string(content.GetContentType()); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "gzip", string(content.GetContentEncoding()); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "utf-8", string(content.GetCharset()); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			indexes := content.GetIndexes()
			if expected, actual := 2, len(indexes); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "email_bin", string(indexes[0].Key); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "golang@basho.com", string(indexes[0].Value); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "email_bin", string(indexes[1].Key); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "frazzle@basho.com", string(indexes[1].Value); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			usermeta := content.GetUsermeta()
			if expected, actual := 2, len(usermeta); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "metaKey1", string(usermeta[0].Key); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "metaValue1", string(usermeta[0].Value); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "metaKey2", string(usermeta[1].Key); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "metaValue2", string(usermeta[1].Value); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			links := content.GetLinks()
			if expected, actual := 2, len(links); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "b0", string(links[0].Bucket); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "k0", string(links[0].Key); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "t0", string(links[0].Tag); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "b1", string(links[1].Bucket); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "k1", string(links[1].Key); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
			if expected, actual := "t1", string(links[1].Tag); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbPutReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestBuildRpbPutReqUsingObjectValues(t *testing.T) {
	objectVclock := "object_vclock"
	ro := &Object{
		BucketType: "object_bucket_type",
		Bucket:     "object_bucket_name",
		Key:        "object_key",
		VClock:     []byte(objectVclock),
	}

	builder := NewStoreValueCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("key").
		WithVClock([]byte("vclock")).
		WithContent(ro)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}

	if req, ok := protobuf.(*rpbRiakKV.RpbPutReq); ok {
		if expected, actual := "object_bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "object_bucket_name", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "object_key", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "object_vclock", string(req.GetVclock()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbPutReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbPutRespCorrectly(t *testing.T) {
	rpbContent := &rpbRiakKV.RpbContent{
		Value:           []byte("this is a value"),
		ContentType:     []byte("text/plain"),
		ContentEncoding: []byte("ascii"),
		Charset:         []byte("ascii"),
		Links: []*rpbRiakKV.RpbLink{
			{[]byte("b0"), []byte("k0"), []byte("t0"), nil},
			{[]byte("b1"), []byte("k1"), []byte("t1"), nil},
		},
		Usermeta: []*rpbRiak.RpbPair{
			{[]byte("metaKey0"), []byte("metaValue0"), nil},
			{[]byte("metaKey1"), []byte("metaValue1"), nil},
		},
		Indexes: []*rpbRiak.RpbPair{
			{[]byte("email_bin"), []byte("golang@basho.com"), nil},
			{[]byte("email_bin"), []byte("frazzle@basho.com"), nil},
			{[]byte("test_int"), []byte("1"), nil},
			{[]byte("test_int"), []byte("2"), nil},
		},
	}

	rpbPutResp := &rpbRiakKV.RpbPutResp{
		Content: []*rpbRiakKV.RpbContent{rpbContent},
		Vclock:  vclock.Bytes(),
		Key:     []byte("generated_riak_key"),
	}

	builder := NewStoreValueCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("ignored_key").
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	cmd.onSuccess(rpbPutResp)
	if expected, actual := true, cmd.Success(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	if storeValueCommand, ok := cmd.(*StoreValueCommand); ok {
		if storeValueCommand.Response == nil {
			t.Fatal("unexpected nil object")
		}
		if expected, actual := true, storeValueCommand.success; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 1, len(storeValueCommand.Response.Values); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		ro := storeValueCommand.Response.Values[0]
		if ro == nil {
			t.Fatal("unexpected nil object")
		}
		if expected, actual := "bucket_type", ro.BucketType; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "bucket_name", ro.Bucket; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "generated_riak_key", ro.Key; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "text/plain", ro.ContentType; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "ascii", ro.Charset; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "ascii", ro.ContentEncoding; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, ro.HasLinks(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 2, len(ro.Links); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		for i, link := range ro.Links {
			bucket := fmt.Sprintf("b%d", i)
			if expected, actual := bucket, string(link.Bucket); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
			key := fmt.Sprintf("k%d", i)
			if expected, actual := key, string(link.Key); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
			tag := fmt.Sprintf("t%d", i)
			if expected, actual := tag, string(link.Tag); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
		}
		if expected, actual := true, ro.HasUserMeta(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 2, len(ro.UserMeta); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		for i, meta := range ro.UserMeta {
			key := fmt.Sprintf("metaKey%d", i)
			if expected, actual := key, string(meta.Key); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
			value := fmt.Sprintf("metaValue%d", i)
			if expected, actual := value, string(meta.Value); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
		}
		if expected, actual := true, ro.HasIndexes(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 2, len(ro.Indexes); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "golang@basho.com", ro.Indexes["email_bin"][0]; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "frazzle@basho.com", ro.Indexes["email_bin"][1]; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "1", ro.Indexes["test_int"][0]; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "2", ro.Indexes["test_int"][1]; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchValueCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestParseRpbPutRespWithSiblingsCorrectly(t *testing.T) {
	rpb1 := &rpbRiakKV.RpbContent{
		Value:           []byte("value_1"),
		ContentType:     []byte("text/plain"),
		ContentEncoding: []byte("ascii"),
		Charset:         []byte("ascii"),
	}
	rpb2 := &rpbRiakKV.RpbContent{
		Value:           []byte("value_2"),
		ContentType:     []byte("text/plain"),
		ContentEncoding: []byte("ascii"),
		Charset:         []byte("ascii"),
	}

	rpbPutResp := &rpbRiakKV.RpbPutResp{
		Content: []*rpbRiakKV.RpbContent{rpb1, rpb2},
		Vclock:  vclock.Bytes(),
		Key:     []byte("generated_riak_key"),
	}

	builder := NewStoreValueCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("ignored_key").
		WithConflictResolver(resolver).
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	cmd.onSuccess(rpbPutResp)
	if actual, expected := cmd.Success(), true; actual != expected {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	if storeValueCommand, ok := cmd.(*StoreValueCommand); ok {
		if actual, expected := len(storeValueCommand.Response.Values), 1; actual != expected {
			t.Errorf("got %v, expected %v", actual, expected)
		}
		ro := storeValueCommand.Response.Values[0]
		if ro == nil {
			t.Fatal("unexpected nil object")
		}
		if actual, expected := string(ro.Value), "value_1"; actual != expected {
			t.Errorf("got %v, expected %v", actual, expected)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchValueCommand", ok, reflect.TypeOf(cmd))
	}
}

// DeleteValue

func TestBuildRpbDelReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewDeleteValueCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("key").
		WithR(1).
		WithPr(2).
		WithW(3).
		WithPw(4).
		WithDw(5).
		WithRw(6).
		WithVClock(vclockBytes).
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

	if req, ok := protobuf.(*rpbRiakKV.RpbDelReq); ok {
		if expected, actual := "bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "bucket_name", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "key", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(1), req.GetR(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(2), req.GetPr(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(3), req.GetW(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(4), req.GetPw(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(5), req.GetDw(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(6), req.GetRw(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 0, bytes.Compare(vclockBytes, req.GetVclock()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbDelReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestValidationOfRpbDelReqViaBuilder(t *testing.T) {
	builder := NewDeleteValueCommandBuilder()
	// validate that Bucket is required
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is required
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrKeyRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

// ListBuckets
func TestListBucketsErrorsViaBuilder(t *testing.T) {
	var streamingCallback = func(buckets []string) error { return nil }
	builder := NewListBucketsCommandBuilder().
		WithBucketType("bucket_type").
		WithStreaming(true).
		WithCallback(streamingCallback).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err == nil {
		t.Errorf("expected cmd %s to error when building if WithAllowListing not called!", reflect.TypeOf(cmd))
	}
}

func TestBuildRpbListBucketsReqCorrectlyViaBuilder(t *testing.T) {
	var streamingCallback = func(buckets []string) error { return nil }
	builder := NewListBucketsCommandBuilder().
		WithAllowListing().
		WithBucketType("bucket_type").
		WithStreaming(true).
		WithCallback(streamingCallback).
		WithTimeout(time.Second * 20)
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

	if req, ok := protobuf.(*rpbRiakKV.RpbListBucketsReq); ok {
		if expected, actual := "bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, req.GetStream(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbDelReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestMultipleRpbListBucketsRespValuesNonStreaming(t *testing.T) {
	builder := NewListBucketsCommandBuilder().
		WithAllowListing().
		WithBucketType("bucket_type").
		WithStreaming(false)

	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbListBucketsResp := &rpbRiakKV.RpbListBucketsResp{}
		buckets := make([][]byte, 5)
		for j := 0; j < 5; j++ {
			buckets[j] = []byte("bucket")
		}
		rpbListBucketsResp.Buckets = buckets
		if i == 19 {
			rpbListBucketsResp.Done = &done
		}

		cmd.onSuccess(rpbListBucketsResp)
	}

	if listBucketsCommand, ok := cmd.(*ListBucketsCommand); ok {
		response := listBucketsCommand.Response
		if expected, actual := 100, len(response.Buckets); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *ListBucketsCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestMultipleRpbListBucketsRespValuesWithStreaming(t *testing.T) {
	count := 0
	timesCalled := 0
	var streamingCallback = func(buckets []string) error {
		timesCalled++
		count += len(buckets)
		return nil
	}

	builder := NewListBucketsCommandBuilder().
		WithAllowListing().
		WithBucketType("bucket_type").
		WithStreaming(true).
		WithCallback(streamingCallback)

	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbListBucketsResp := &rpbRiakKV.RpbListBucketsResp{}
		buckets := make([][]byte, 5)
		for j := 0; j < 5; j++ {
			buckets[j] = []byte("bucket")
		}
		rpbListBucketsResp.Buckets = buckets
		if i == 19 {
			rpbListBucketsResp.Done = &done
		}

		cmd.onSuccess(rpbListBucketsResp)
	}

	if expected, actual := 20, timesCalled; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if expected, actual := 100, count; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestValidationOfRpbListBucketsReqViaBuilder(t *testing.T) {
	builder := NewListBucketsCommandBuilder().
		WithAllowListing()
	// validate that Bucket and Key are NOT required
	// and that type is "default"
	var err error
	var cmd Command
	var protobuf proto.Message
	cmd, err = builder.Build()
	if err == nil {
		protobuf, err = cmd.constructPbRequest()
		if err != nil {
			t.Fatal(err.Error())
		}
		if req, ok := protobuf.(*rpbRiakKV.RpbListBucketsReq); ok {
			if expected, actual := "default", string(req.GetType()); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
		} else {
			t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbDelReq", ok, reflect.TypeOf(protobuf))
		}
	} else {
		t.Fatal("expected nil err")
	}
}

// ListKeys
func TestListKeysErrorsViaBuilder(t *testing.T) {
	var streamingCallback = func(buckets []string) error { return nil }
	builder := NewListKeysCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket").
		WithStreaming(true).
		WithCallback(streamingCallback).
		WithTimeout(time.Second * 20)
	cmd, err := builder.Build()
	if err == nil {
		t.Errorf("expected cmd %s to error when building if WithAllowListing not called!", reflect.TypeOf(cmd))
	}
}

func TestBuildRpbListKeysReqCorrectlyViaBuilder(t *testing.T) {
	var streamingCallback = func(buckets []string) error { return nil }
	builder := NewListKeysCommandBuilder().
		WithAllowListing().
		WithBucketType("bucket_type").
		WithBucket("bucket").
		WithStreaming(true).
		WithCallback(streamingCallback).
		WithTimeout(time.Second * 20)
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

	if req, ok := protobuf.(*rpbRiakKV.RpbListKeysReq); ok {
		if expected, actual := "bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "bucket", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbDelReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestMultipleRpbListKeysRespValuesNonStreaming(t *testing.T) {
	builder := NewListKeysCommandBuilder().
		WithAllowListing().
		WithBucketType("bucket_type").
		WithBucket("bucket").
		WithStreaming(false)

	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbListKeysResp := &rpbRiakKV.RpbListKeysResp{}
		buckets := make([][]byte, 5)
		for j := 0; j < 5; j++ {
			buckets[j] = []byte("bucket")
		}
		rpbListKeysResp.Keys = buckets
		if i == 19 {
			rpbListKeysResp.Done = &done
		}

		cmd.onSuccess(rpbListKeysResp)
	}

	if listKeysCommand, ok := cmd.(*ListKeysCommand); ok {
		response := listKeysCommand.Response
		if expected, actual := 100, len(response.Keys); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *ListKeysCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestMultipleRpbListKeysRespValuesWithStreaming(t *testing.T) {
	count := 0
	timesCalled := 0
	var streamingCallback = func(buckets []string) error {
		timesCalled++
		count += len(buckets)
		return nil
	}

	builder := NewListKeysCommandBuilder().
		WithAllowListing().
		WithBucketType("bucket_type").
		WithBucket("bucket").
		WithStreaming(true).
		WithCallback(streamingCallback)

	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbListKeysResp := &rpbRiakKV.RpbListKeysResp{}
		buckets := make([][]byte, 5)
		for j := 0; j < 5; j++ {
			buckets[j] = []byte("bucket")
		}
		rpbListKeysResp.Keys = buckets
		if i == 19 {
			rpbListKeysResp.Done = &done
		}

		cmd.onSuccess(rpbListKeysResp)
	}

	if expected, actual := 20, timesCalled; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if expected, actual := 100, count; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestValidationOfRpbListKeysReqViaBuilder(t *testing.T) {
	builder := NewListKeysCommandBuilder().
		WithAllowListing().
		WithBucket("bucket")
	// validate that Key is NOT required
	// and that type is "default"
	var err error
	var cmd Command
	var protobuf proto.Message
	cmd, err = builder.Build()
	if err == nil {
		protobuf, err = cmd.constructPbRequest()
		if err != nil {
			t.Fatal(err.Error())
		}
		if req, ok := protobuf.(*rpbRiakKV.RpbListKeysReq); ok {
			if expected, actual := "default", string(req.GetType()); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
			if expected, actual := "bucket", string(req.GetBucket()); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
		} else {
			t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbDelReq", ok, reflect.TypeOf(protobuf))
		}
	} else {
		t.Fatal("expected nil err")
	}
}

// FetchPreflist

func TestBuildRpbGetBucketKeyPreflistReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewFetchPreflistCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithKey("key")
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

	if req, ok := protobuf.(*rpbRiakKV.RpbGetBucketKeyPreflistReq); ok {
		if expected, actual := "bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "bucket_name", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "key", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbDelReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestValidationOfRpbGetBucketKeyPreflistReqViaBuilder(t *testing.T) {
	builder := NewFetchPreflistCommandBuilder()
	// validate that Bucket is required
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is required
	builder.WithBucket("bucket_name")
	_, err = builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrKeyRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

// SecondaryIndexQuery

func TestBuildRpbIndexReqCorrectlyViaBuilder(t *testing.T) {
	// should error due to no index query data
	builder := NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name")
	cmd, err := builder.Build()
	if err == nil {
		t.Fatal("expected error")
	} else {
		if expected, actual := "ClientError|either WithIndexKey or WithRange are required", err.Error(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	}

	continuationBytes := bytes.NewBufferString("continuation_1234").Bytes()
	var cb = func(results []*SecondaryIndexQueryResult) error {
		return nil
	}

	builder = NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithIndexName("email_bin").
		WithIndexKey("golang@basho.com").
		WithReturnKeyAndIndex(true).
		WithCallback(cb).
		WithStreaming(true).
		WithPaginationSort(true).
		WithMaxResults(1024).
		WithContinuation(continuationBytes).
		WithTermRegex("^yomama").
		WithTimeout(time.Second * 20)
	cmd, err = builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); ok {
		t.Errorf("got %v, want cmd %s to NOT implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	var protobuf proto.Message
	protobuf, err = cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if req, ok := protobuf.(*rpbRiakKV.RpbIndexReq); ok {
		if expected, actual := "bucket_type", string(req.GetType()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "bucket_name", string(req.GetBucket()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "email_bin", string(req.GetIndex()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "golang@basho.com", string(req.GetKey()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, req.GetReturnTerms(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, req.GetStream(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := true, req.GetPaginationSort(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := uint32(1024), req.GetMaxResults(); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := 0, bytes.Compare(continuationBytes, req.GetContinuation()); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "^yomama", string(req.GetTermRegex()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		validateTimeout(t, time.Second*20, req.GetTimeout())
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbIndexReq", ok, reflect.TypeOf(protobuf))
	}

	builder = NewSecondaryIndexQueryCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithIndexName("email_int").
		WithIntRange(100, 200)
	cmd, err = builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	protobuf, err = cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if req, ok := protobuf.(*rpbRiakKV.RpbIndexReq); ok {
		if expected, actual := "100", string(req.GetRangeMin()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "200", string(req.GetRangeMax()); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbIndexReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestMultipleRpbIndexRespWithObjectKeysCorrectly(t *testing.T) {
	builder := NewSecondaryIndexQueryCommandBuilder().
		WithBucket("bucket").
		WithIndexName("id_bin").
		WithStreaming(false).
		WithIntRange(0, 50)
	cmd, err := builder.Build()
	if err != nil {
		t.Error(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbIndexResp := &rpbRiakKV.RpbIndexResp{
			Keys: make([][]byte, 5),
		}
		for j := 0; j < 5; j++ {
			rpbIndexResp.Keys[j] = []byte("object_key")
		}
		if i == 19 {
			rpbIndexResp.Done = &done
			rpbIndexResp.Continuation = []byte("1234")
		}
		cmd.onSuccess(rpbIndexResp)
	}
	if expected, actual := true, cmd.Success(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if siq, ok := cmd.(*SecondaryIndexQueryCommand); ok {
		if expected, actual := true, siq.done; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		rsp := siq.Response
		if rsp == nil {
			t.Fatal("expected non-nil Response")
		}
		if expected, actual := 100, len(rsp.Results); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if rsp.Results[0].IndexKey != nil {
			t.Error("expected nil IndexKey value")
		}
		if expected, actual := "object_key", string(rsp.Results[0].ObjectKey); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *SecondaryIndexQueryCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestMultipleRpbIndexRespWithTermKeyPairsCorrectly(t *testing.T) {
	builder := NewSecondaryIndexQueryCommandBuilder().
		WithBucket("bucket").
		WithIndexName("id_bin").
		WithStreaming(false).
		WithReturnKeyAndIndex(true).
		WithIntRange(0, 50)
	cmd, err := builder.Build()
	if err != nil {
		t.Error(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbIndexResp := &rpbRiakKV.RpbIndexResp{}
		for j := 0; j < 5; j++ {
			pair := &rpbRiak.RpbPair{
				Key:   []byte("index_key"),
				Value: []byte("object_key"),
			}
			rpbIndexResp.Results = append(rpbIndexResp.Results, pair)
		}
		if i == 19 {
			rpbIndexResp.Done = &done
			rpbIndexResp.Continuation = []byte("1234")
		}
		cmd.onSuccess(rpbIndexResp)
	}
	if expected, actual := true, cmd.Success(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if siq, ok := cmd.(*SecondaryIndexQueryCommand); ok {
		if expected, actual := true, siq.done; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		rsp := siq.Response
		if rsp == nil {
			t.Fatal("expected non-nil Response")
		}
		if expected, actual := 100, len(rsp.Results); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "index_key", string(rsp.Results[0].IndexKey); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "object_key", string(rsp.Results[0].ObjectKey); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *SecondaryIndexQueryCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestMultipleRpbIndexRespWithObjectKeysCorrectlyStreams(t *testing.T) {
	count := 0
	timesCalled := 0
	var cb = func(results []*SecondaryIndexQueryResult) error {
		timesCalled++
		count += len(results)
		if results[0].IndexKey != nil {
			t.Error("expected nil IndexKey value")
		}
		if expected, actual := "object_key", string(results[0].ObjectKey); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		return nil
	}

	builder := NewSecondaryIndexQueryCommandBuilder().
		WithBucket("bucket").
		WithIndexName("id_bin").
		WithStreaming(true).
		WithCallback(cb).
		WithIntRange(0, 50)
	cmd, err := builder.Build()
	if err != nil {
		t.Error(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbIndexResp := &rpbRiakKV.RpbIndexResp{
			Keys: make([][]byte, 5),
		}
		for j := 0; j < 5; j++ {
			rpbIndexResp.Keys[j] = []byte("object_key")
		}
		if i == 19 {
			rpbIndexResp.Done = &done
			rpbIndexResp.Continuation = []byte("1234")
		}
		cmd.onSuccess(rpbIndexResp)
	}

	if expected, actual := 20, timesCalled; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if expected, actual := 100, count; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if siq, ok := cmd.(*SecondaryIndexQueryCommand); ok {
		if expected, actual := true, siq.done; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		rsp := siq.Response
		if rsp == nil {
			t.Fatal("expected non-nil Response")
		}
		if expected, actual := "1234", string(rsp.Continuation); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *SecondaryIndexQueryCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestMultipleRpbIndexRespWithTermKeyPairsCorrectlyStreams(t *testing.T) {
	count := 0
	timesCalled := 0
	var cb = func(results []*SecondaryIndexQueryResult) error {
		timesCalled++
		count += len(results)
		if expected, actual := "index_key", string(results[0].IndexKey); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		if expected, actual := "object_key", string(results[0].ObjectKey); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		return nil
	}

	builder := NewSecondaryIndexQueryCommandBuilder().
		WithBucket("bucket").
		WithIndexName("id_int").
		WithCallback(cb).
		WithStreaming(true).
		WithReturnKeyAndIndex(true).
		WithIntRange(0, 50)
	cmd, err := builder.Build()
	if err != nil {
		t.Error(err.Error())
	}

	done := true
	for i := 0; i < 20; i++ {
		rpbIndexResp := &rpbRiakKV.RpbIndexResp{}
		for j := 0; j < 5; j++ {
			pair := &rpbRiak.RpbPair{
				Key:   []byte("index_key"),
				Value: []byte("object_key"),
			}
			rpbIndexResp.Results = append(rpbIndexResp.Results, pair)
		}
		if i == 19 {
			rpbIndexResp.Done = &done
			rpbIndexResp.Continuation = []byte("1234")
		}
		cmd.onSuccess(rpbIndexResp)
	}

	if expected, actual := 20, timesCalled; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if expected, actual := 100, count; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if siq, ok := cmd.(*SecondaryIndexQueryCommand); ok {
		if expected, actual := true, siq.done; expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
		rsp := siq.Response
		if rsp == nil {
			t.Fatal("expected non-nil Response")
		}
		if expected, actual := "1234", string(rsp.Continuation); expected != actual {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *SecondaryIndexQueryCommand", ok, reflect.TypeOf(cmd))
	}
}

func TestValidationOfRpbIndexReqViaBuilder(t *testing.T) {
	builder := NewSecondaryIndexQueryCommandBuilder()
	// validate that Bucket is required
	_, err := builder.Build()
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if expected, actual := ErrBucketRequired.Error(), err.Error(); expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}

	// validate that Key is NOT required
	builder.WithBucket("bucket_name")
	builder.WithRange("frazzle", "pop")
	_, err = builder.Build()
	if err != nil {
		t.Fatal("expected nil err")
	}
}

// MapReduce

func TestBuildRpbMapRedReqCorrectlyViaBuilder(t *testing.T) {
	query := "{\"inputs\":\"goog\",\"query\":[{\"map\":{\"language\":\"javascript\",\"source\":\"function(value, keyData, arg) { var data = Riak.mapValuesJson(value)[0]; if(data.High && parseFloat(data.High) > 600.00) return [value.key];else return [];}\",\"keep\":true}}]}"
	var err error
	var mr Command
	var protobuf proto.Message
	if mr, err = NewMapReduceCommandBuilder().WithQuery(query).WithStreaming(false).Build(); err == nil {
		if protobuf, err = mr.constructPbRequest(); err == nil {
			if req, ok := protobuf.(*rpbRiakKV.RpbMapRedReq); ok {
				if expected, actual := query, string(req.GetRequest()); expected != actual {
					t.Errorf("expected %v, actual %v", expected, actual)
				}
				if expected, actual := "application/json", string(req.GetContentType()); expected != actual {
					t.Errorf("expected %v, actual %v", expected, actual)
				}
			} else {
				t.Errorf("ok: %v - could not convert %v to *rpbRiakKV.RpbMapRedReq", ok, reflect.TypeOf(protobuf))
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Error(err.Error())
	}
}

func TestParseRpbMapRedRespCorrectly(t *testing.T) {
	done := true
	rspJSON := "[{\"the\": 8}]"
	rpbResponse := []byte(rspJSON)
	if cmd, err := NewMapReduceCommandBuilder().WithQuery("some query").Build(); err == nil {
		for i := 1; i <= 10; i++ {
			phase := uint32(i)
			rpbMapRedResp := &rpbRiakKV.RpbMapRedResp{
				Phase:    &phase,
				Response: rpbResponse,
			}
			if i == 10 {
				rpbMapRedResp.Done = &done
			}
			cmd.onSuccess(rpbMapRedResp)
		}
		if mr, ok := cmd.(*MapReduceCommand); ok {
			rsp := mr.Response
			if expected, actual := 10, len(rsp); expected != actual {
				t.Errorf("expected %v, actual %v", expected, actual)
			}
			if expected, actual := 0, bytes.Compare(rpbResponse, rsp[0]); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		} else {
			t.Errorf("Could not convert %v to *MapReduceCommand", reflect.TypeOf(cmd))
		}
	} else {
		t.Error(err.Error())
	}
}

func TestParseRpbMapRedRespCorrectlyWithStreaming(t *testing.T) {
	done := true
	rspJSON := "[{\"the\": 8}]"
	rpbResponse := []byte(rspJSON)

	count := 0
	var cb = func(response []byte) error {
		count++
		if expected, actual := 0, bytes.Compare(rpbResponse, response); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		return nil
	}

	if cmd, err := NewMapReduceCommandBuilder().
		WithQuery("some query").
		WithCallback(cb).
		WithStreaming(true).
		Build(); err == nil {

		if _, ok := cmd.(retryableCommand); ok {
			t.Errorf("got %v, want cmd %s to NOT implement retryableCommand", ok, reflect.TypeOf(cmd))
		}

		for i := 1; i <= 10; i++ {
			phase := uint32(i)
			rpbMapRedResp := &rpbRiakKV.RpbMapRedResp{
				Phase:    &phase,
				Response: rpbResponse,
			}
			if i == 10 {
				rpbMapRedResp.Done = &done
			}
			cmd.onSuccess(rpbMapRedResp)
			if mr, ok := cmd.(*MapReduceCommand); ok {
				if mr.Response != nil {
					t.Error("expected nil results")
				}
			} else {
				t.Errorf("Could not convert %v to *MapReduceCommand", reflect.TypeOf(cmd))
			}
		}
		if expected, actual := 10, count; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Error(err.Error())
	}
}
