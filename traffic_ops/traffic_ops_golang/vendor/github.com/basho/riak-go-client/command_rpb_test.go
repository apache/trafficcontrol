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
	rpbRiak "github.com/basho/riak-go-client/rpb/riak"
	rpbRiakDT "github.com/basho/riak-go-client/rpb/riak_dt"
	rpbRiakKV "github.com/basho/riak-go-client/rpb/riak_kv"
	rpbRiakSCH "github.com/basho/riak-go-client/rpb/riak_search"
	rpbRiakYZ "github.com/basho/riak-go-client/rpb/riak_yokozuna"
	proto "github.com/golang/protobuf/proto"
	"reflect"
	"testing"
)

func TestEnsureCorrectRequestAndResponseCodes(t *testing.T) {
	var cmd Command
	var msg proto.Message
	// Misc commands
	// Ping
	cmd = &PingCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbPingReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbPingResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if cmd.getResponseProtobufMessage() != nil {
		t.Error("want nil response protobuf message")
	}
	// GetServerInfo
	cmd = &GetServerInfoCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbGetServerInfoReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbGetServerInfoResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiak.RpbGetServerInfoResp); !ok {
		t.Errorf("error casting %v to RpbGetServerInfoResp", reflect.TypeOf(msg))
	}
	// StartTls
	cmd = &startTlsCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbStartTls; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbStartTls; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if cmd.getResponseProtobufMessage() != nil {
		t.Error("want nil response protobuf message")
	}
	// Auth
	cmd = &authCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbAuthReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbAuthResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if cmd.getResponseProtobufMessage() != nil {
		t.Error("want nil response protobuf message")
	}
	// FetchBucketTypeProps
	cmd = &FetchBucketTypePropsCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbGetBucketTypeReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbGetBucketResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiak.RpbGetBucketResp); !ok {
		t.Errorf("error casting %v to RpbGetBucketResp", reflect.TypeOf(msg))
	}
	// FetchBucketProps
	cmd = &FetchBucketPropsCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbGetBucketReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbGetBucketResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiak.RpbGetBucketResp); !ok {
		t.Errorf("error casting %v to RpbGetBucketResp", reflect.TypeOf(msg))
	}
	// StoreBucketTypeProps
	cmd = &StoreBucketTypePropsCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbSetBucketTypeReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbSetBucketResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if msg != nil {
		t.Error("want nil response protobuf message")
	}
	// StoreBucketProps
	cmd = &StoreBucketPropsCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbSetBucketReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbSetBucketResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if msg != nil {
		t.Error("want nil response protobuf message")
	}

	// KV commands
	// FetchValue
	cmd = &FetchValueCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbGetReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbGetResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakKV.RpbGetResp); !ok {
		t.Errorf("error casting %v to RpbGetResp", reflect.TypeOf(msg))
	}
	// StoreValue
	cmd = &StoreValueCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbPutReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbPutResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakKV.RpbPutResp); !ok {
		t.Errorf("error casting %v to RpbPutResp", reflect.TypeOf(msg))
	}
	// DeleteValue
	cmd = &DeleteValueCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbDelReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbDelResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if msg != nil {
		t.Error("want nil response protobuf message")
	}
	// ListBuckets
	cmd = &ListBucketsCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbListBucketsReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbListBucketsResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakKV.RpbListBucketsResp); !ok {
		t.Errorf("error casting %v to RpbListBucketsResp", reflect.TypeOf(msg))
	}
	// ListKeys
	cmd = &ListKeysCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbListKeysReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbListKeysResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakKV.RpbListKeysResp); !ok {
		t.Errorf("error casting %v to RpbListKeysResp", reflect.TypeOf(msg))
	}
	// FetchPreflist
	cmd = &FetchPreflistCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbGetBucketKeyPreflistReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbGetBucketKeyPreflistResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakKV.RpbGetBucketKeyPreflistResp); !ok {
		t.Errorf("error casting %v to RpbGetBucketKeyPreflistResp", reflect.TypeOf(msg))
	}
	// SecondaryIndexQuery
	cmd = &SecondaryIndexQueryCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbIndexReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbIndexResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakKV.RpbIndexResp); !ok {
		t.Errorf("error casting %v to RpbIndexResp", reflect.TypeOf(msg))
	}
	// MapReduce
	cmd = &MapReduceCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbMapRedReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbMapRedResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakKV.RpbMapRedResp); !ok {
		t.Errorf("error casting %v to RpbMapRedResp", reflect.TypeOf(msg))
	}

	// YZ commands
	// StoreIndex
	cmd = &StoreIndexCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbYokozunaIndexPutReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbPutResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if cmd.getResponseProtobufMessage() != nil {
		t.Error("want nil response protobuf message")
	}
	// FetchIndex
	cmd = &FetchIndexCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbYokozunaIndexGetReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbYokozunaIndexGetResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakYZ.RpbYokozunaIndexGetResp); !ok {
		t.Errorf("error casting %v to RpbYokozunaIndexGetResp", reflect.TypeOf(msg))
	}
	// DeleteIndex
	cmd = &DeleteIndexCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbYokozunaIndexDeleteReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbDelResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if cmd.getResponseProtobufMessage() != nil {
		t.Error("want nil response protobuf message")
	}
	// StoreSchema
	cmd = &StoreSchemaCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbYokozunaSchemaPutReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbPutResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if cmd.getResponseProtobufMessage() != nil {
		t.Error("want nil response protobuf message")
	}
	// FetchSchema
	cmd = &FetchSchemaCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbYokozunaSchemaGetReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbYokozunaSchemaGetResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakYZ.RpbYokozunaSchemaGetResp); !ok {
		t.Errorf("error casting %v to RpbYokozunaSchemaGetResp", reflect.TypeOf(msg))
	}
	// Search
	cmd = &SearchCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_RpbSearchQueryReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_RpbSearchQueryResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakSCH.RpbSearchQueryResp); !ok {
		t.Errorf("error casting %v to RpbSearchQueryResp", reflect.TypeOf(msg))
	}

	// CRDT commands
	// UpdateCounter
	cmd = &UpdateCounterCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_DtUpdateReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_DtUpdateResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakDT.DtUpdateResp); !ok {
		t.Errorf("error casting %v to DtUpdateResp", reflect.TypeOf(msg))
	}
	// FetchCounter
	cmd = &FetchCounterCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_DtFetchReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_DtFetchResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakDT.DtFetchResp); !ok {
		t.Errorf("error casting %v to DtFetchResp", reflect.TypeOf(msg))
	}
	// UpdateSet
	cmd = &UpdateSetCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_DtUpdateReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_DtUpdateResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakDT.DtUpdateResp); !ok {
		t.Errorf("error casting %v to DtUpdateResp", reflect.TypeOf(msg))
	}
	// FetchSet
	cmd = &FetchSetCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_DtFetchReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_DtFetchResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakDT.DtFetchResp); !ok {
		t.Errorf("error casting %v to DtFetchResp", reflect.TypeOf(msg))
	}
	// UpdateMap
	cmd = &UpdateMapCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_DtUpdateReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_DtUpdateResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakDT.DtUpdateResp); !ok {
		t.Errorf("error casting %v to DtUpdateResp", reflect.TypeOf(msg))
	}
	// FetchMap
	cmd = &FetchMapCommand{}
	if got, want := cmd.getRequestCode(), rpbCode_DtFetchReq; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := cmd.getResponseCode(), rpbCode_DtFetchResp; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	msg = cmd.getResponseProtobufMessage()
	if _, ok := msg.(*rpbRiakDT.DtFetchResp); !ok {
		t.Errorf("error casting %v to DtFetchResp", reflect.TypeOf(msg))
	}
}
