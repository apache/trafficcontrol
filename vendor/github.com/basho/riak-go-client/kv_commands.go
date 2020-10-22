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
	"fmt"
	"reflect"
	"strconv"
	"time"

	rpbRiakKV "github.com/basho/riak-go-client/rpb/riak_kv"
	proto "github.com/golang/protobuf/proto"
)

// FetchValue
// RpbGetReq
// RpbGetResp

// ConflictResolver is an interface to handle sibling conflicts for a key
type ConflictResolver interface {
	Resolve([]*Object) []*Object
}

// FetchValueCommand is used to fetch / get a value from Riak KV
type FetchValueCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *FetchValueResponse
	protobuf *rpbRiakKV.RpbGetReq
	resolver ConflictResolver
}

// Name identifies this command
func (cmd *FetchValueCommand) Name() string {
	return cmd.getName("FetchValue")
}

func (cmd *FetchValueCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchValueCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg == nil {
		cmd.Response = &FetchValueResponse{
			IsNotFound:  true,
			IsUnchanged: false,
		}
	} else {
		if rpbGetResp, ok := msg.(*rpbRiakKV.RpbGetResp); ok {
			vclock := rpbGetResp.GetVclock()
			response := &FetchValueResponse{
				VClock:      vclock,
				IsUnchanged: rpbGetResp.GetUnchanged(),
				IsNotFound:  false,
			}

			if pbContent := rpbGetResp.GetContent(); pbContent == nil || len(pbContent) == 0 {
				object := &Object{
					IsTombstone: true,
					BucketType:  string(cmd.protobuf.Type),
					Bucket:      string(cmd.protobuf.Bucket),
					Key:         string(cmd.protobuf.Key),
				}
				response.Values = []*Object{object}
			} else {
				response.Values = make([]*Object, len(pbContent))
				for i, content := range pbContent {
					ro, err := fromRpbContent(content)
					if err != nil {
						return err
					}
					ro.VClock = vclock
					ro.BucketType = string(cmd.protobuf.Type)
					ro.Bucket = string(cmd.protobuf.Bucket)
					ro.Key = string(cmd.protobuf.Key)
					response.Values[i] = ro
				}
				if cmd.resolver != nil {
					response.Values = cmd.resolver.Resolve(response.Values)
				}
			}

			cmd.Response = response
		} else {
			return fmt.Errorf("[FetchValueCommand] could not convert %v to RpbGetResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchValueCommand) getRequestCode() byte {
	return rpbCode_RpbGetReq
}

func (cmd *FetchValueCommand) getResponseCode() byte {
	return rpbCode_RpbGetResp
}

func (cmd *FetchValueCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbGetResp{}
}

// FetchValueResponse contains the response data for a FetchValueCommand
type FetchValueResponse struct {
	IsNotFound  bool
	IsUnchanged bool
	VClock      []byte
	Values      []*Object
}

// FetchValueCommandBuilder type is required for creating new instances of FetchValueCommand
//
//	command, err := NewFetchValueCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		Build()
type FetchValueCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakKV.RpbGetReq
	resolver ConflictResolver
}

// NewFetchValueCommandBuilder is a factory function for generating the command builder struct
func NewFetchValueCommandBuilder() *FetchValueCommandBuilder {
	builder := &FetchValueCommandBuilder{protobuf: &rpbRiakKV.RpbGetReq{}}
	return builder
}

// WithConflictResolver builds the command object with a user defined ConflictResolver for handling conflicting key values
func (builder *FetchValueCommandBuilder) WithConflictResolver(resolver ConflictResolver) *FetchValueCommandBuilder {
	builder.resolver = resolver
	return builder
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *FetchValueCommandBuilder) WithBucketType(bucketType string) *FetchValueCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *FetchValueCommandBuilder) WithBucket(bucket string) *FetchValueCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *FetchValueCommandBuilder) WithKey(key string) *FetchValueCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithR sets the number of nodes that must report back a successful read in order for the
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchValueCommandBuilder) WithR(r uint32) *FetchValueCommandBuilder {
	builder.protobuf.R = &r
	return builder
}

// WithPr sets the number of primary nodes (N) that must be read from in order for the command
// operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchValueCommandBuilder) WithPr(pr uint32) *FetchValueCommandBuilder {
	builder.protobuf.Pr = &pr
	return builder
}

// WithNVal sets the number of times this command operation is replicated in the Cluster. If
// ommitted, the ring default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchValueCommandBuilder) WithNVal(nval uint32) *FetchValueCommandBuilder {
	builder.protobuf.NVal = &nval
	return builder
}

// WithBasicQuorum sets basic_quorum, whether to return early in some failure cases (eg. when r=1
// and you get 2 errors and a success basic_quorum=true would return an error)
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchValueCommandBuilder) WithBasicQuorum(basicQuorum bool) *FetchValueCommandBuilder {
	builder.protobuf.BasicQuorum = &basicQuorum
	return builder
}

// WithNotFoundOk sets notfound_ok, whether to treat notfounds as successful reads for the purposes
// of R
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchValueCommandBuilder) WithNotFoundOk(notFoundOk bool) *FetchValueCommandBuilder {
	builder.protobuf.NotfoundOk = &notFoundOk
	return builder
}

// WithIfModified tells Riak to only return the object if the vclock in Riak differs from what is
// provided
func (builder *FetchValueCommandBuilder) WithIfModified(ifModified []byte) *FetchValueCommandBuilder {
	builder.protobuf.IfModified = ifModified
	return builder
}

// WithHeadOnly returns only the meta data for the value, useful when objects contain large amounts
// of data
func (builder *FetchValueCommandBuilder) WithHeadOnly(headOnly bool) *FetchValueCommandBuilder {
	builder.protobuf.Head = &headOnly
	return builder
}

// WithReturnDeletedVClock sets the command to return a Tombstone if any our found for the key across
// all of the vnodes
func (builder *FetchValueCommandBuilder) WithReturnDeletedVClock(returnDeletedVClock bool) *FetchValueCommandBuilder {
	builder.protobuf.Deletedvclock = &returnDeletedVClock
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *FetchValueCommandBuilder) WithTimeout(timeout time.Duration) *FetchValueCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// WithSloppyQuorum sets the sloppy_quorum for this Command
//
// See http://docs.basho.com/riak/latest/theory/concepts/Eventual-Consistency/
func (builder *FetchValueCommandBuilder) WithSloppyQuorum(sloppyQuorum bool) *FetchValueCommandBuilder {
	builder.protobuf.SloppyQuorum = &sloppyQuorum
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchValueCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &FetchValueCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
		resolver: builder.resolver,
	}, nil
}

// StoreValue
// RpbPutReq
// RpbPutResp

// StoreValueCommand used to store a value from Riak KV.
type StoreValueCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *StoreValueResponse
	value    *Object
	protobuf *rpbRiakKV.RpbPutReq
	resolver ConflictResolver
}

// Name identifies this command
func (cmd *StoreValueCommand) Name() string {
	return cmd.getName("StoreValue")
}

func (cmd *StoreValueCommand) constructPbRequest() (msg proto.Message, err error) {
	value := cmd.value

	// Some properties of the value override options
	setProtobufFromValue(cmd.protobuf, cmd.value)

	cmd.protobuf.Content, err = toRpbContent(value)
	if err != nil {
		return
	}

	msg = cmd.protobuf
	return
}

func (cmd *StoreValueCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg == nil {
		cmd.Response = &StoreValueResponse{}
	} else {
		if rpbPutResp, ok := msg.(*rpbRiakKV.RpbPutResp); ok {
			var responseKey string
			if responseKeyBytes := rpbPutResp.GetKey(); responseKeyBytes != nil && len(responseKeyBytes) > 0 {
				responseKey = string(responseKeyBytes)
			}

			vclock := rpbPutResp.GetVclock()
			response := &StoreValueResponse{
				VClock:       vclock,
				GeneratedKey: responseKey,
			}

			if pbContent := rpbPutResp.GetContent(); pbContent != nil && len(pbContent) > 0 {
				response.Values = make([]*Object, len(pbContent))
				for i, content := range pbContent {
					ro, err := fromRpbContent(content)
					if err != nil {
						return err
					}

					ro.VClock = vclock
					ro.BucketType = string(cmd.protobuf.Type)
					ro.Bucket = string(cmd.protobuf.Bucket)
					if responseKey == "" {
						ro.Key = string(cmd.protobuf.Key)
					} else {
						ro.Key = responseKey
					}
					response.Values[i] = ro
				}
				if cmd.resolver != nil {
					response.Values = cmd.resolver.Resolve(response.Values)
				}
			}

			cmd.Response = response
		} else {
			return fmt.Errorf("[StoreValueCommand] could not convert %v to RpbPutResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *StoreValueCommand) getRequestCode() byte {
	return rpbCode_RpbPutReq
}

func (cmd *StoreValueCommand) getResponseCode() byte {
	return rpbCode_RpbPutResp
}

func (cmd *StoreValueCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbPutResp{}
}

func setProtobufFromValue(pb *rpbRiakKV.RpbPutReq, value *Object) {
	if value.VClock != nil {
		pb.Vclock = value.VClock
	}
	if value.BucketType != "" {
		pb.Type = []byte(value.BucketType)
	}
	if value.Bucket != "" {
		pb.Bucket = []byte(value.Bucket)
	}
	if value.Key != "" {
		pb.Key = []byte(value.Key)
	}
}

// StoreValueResponse contains the response data for a StoreValueCommand
type StoreValueResponse struct {
	GeneratedKey string
	VClock       []byte
	Values       []*Object
}

// StoreValueCommandBuilder type is required for creating new instances of StoreValueCommand
//
//	command, err := NewStoreValueCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		Build()
type StoreValueCommandBuilder struct {
	value    *Object
	timeout  time.Duration
	protobuf *rpbRiakKV.RpbPutReq
	resolver ConflictResolver
}

// NewStoreValueCommandBuilder is a factory function for generating the command builder struct
func NewStoreValueCommandBuilder() *StoreValueCommandBuilder {
	builder := &StoreValueCommandBuilder{protobuf: &rpbRiakKV.RpbPutReq{}}
	return builder
}

// WithConflictResolver sets the ConflictResolver that should be used when sibling conflicts are found
// for this operation
func (builder *StoreValueCommandBuilder) WithConflictResolver(resolver ConflictResolver) *StoreValueCommandBuilder {
	builder.resolver = resolver
	return builder
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *StoreValueCommandBuilder) WithBucketType(bucketType string) *StoreValueCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *StoreValueCommandBuilder) WithBucket(bucket string) *StoreValueCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *StoreValueCommandBuilder) WithKey(key string) *StoreValueCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithVClock sets the vclock for the object to be stored, providing causal context for conflicts
func (builder *StoreValueCommandBuilder) WithVClock(vclock []byte) *StoreValueCommandBuilder {
	builder.protobuf.Vclock = vclock
	return builder
}

// WithContent sets the object / value to be stored at the specified key
func (builder *StoreValueCommandBuilder) WithContent(object *Object) *StoreValueCommandBuilder {
	setProtobufFromValue(builder.protobuf, object)
	builder.value = object
	return builder
}

// WithW sets the number of nodes that must report back a successful write in order for then
// command operation to be considered a success by Riak
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *StoreValueCommandBuilder) WithW(w uint32) *StoreValueCommandBuilder {
	builder.protobuf.W = &w
	return builder
}

// WithDw (durable writes) sets the number of nodes that must report back a successful write to
// backend storage in order for the command operation to be considered a success by Riak. If
// ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *StoreValueCommandBuilder) WithDw(dw uint32) *StoreValueCommandBuilder {
	builder.protobuf.Dw = &dw
	return builder
}

// WithPw sets the number of primary nodes (N) that must report back a successful write in order for
// the command operation to be considered a success by Riak.  If ommitted, the bucket default is
// used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *StoreValueCommandBuilder) WithPw(pw uint32) *StoreValueCommandBuilder {
	builder.protobuf.Pw = &pw
	return builder
}

// WithNVal sets the number of times this command operation is replicated in the Cluster. If
// ommitted, the ring default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *StoreValueCommandBuilder) WithNVal(nval uint32) *StoreValueCommandBuilder {
	builder.protobuf.NVal = &nval
	return builder
}

// WithReturnBody sets Riak to return the value within its response after completing the write
// operation
func (builder *StoreValueCommandBuilder) WithReturnBody(returnBody bool) *StoreValueCommandBuilder {
	builder.protobuf.ReturnBody = &returnBody
	return builder
}

// WithIfNotModified tells Riak to only update the object in Riak if the vclock provided matches the
// one currently in Riak
func (builder *StoreValueCommandBuilder) WithIfNotModified(ifNotModified bool) *StoreValueCommandBuilder {
	builder.protobuf.IfNotModified = &ifNotModified
	return builder
}

// WithIfNoneMatch tells Riak to store the object only if it does not already exist in the database
func (builder *StoreValueCommandBuilder) WithIfNoneMatch(ifNoneMatch bool) *StoreValueCommandBuilder {
	builder.protobuf.IfNoneMatch = &ifNoneMatch
	return builder
}

// WithReturnHead returns only the meta data for the value, useful when objects contain large amounts
// of data
func (builder *StoreValueCommandBuilder) WithReturnHead(returnHead bool) *StoreValueCommandBuilder {
	builder.protobuf.ReturnHead = &returnHead
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *StoreValueCommandBuilder) WithTimeout(timeout time.Duration) *StoreValueCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// WithAsis sets the asis option
// Please note, this is an advanced feature, only use with caution
func (builder *StoreValueCommandBuilder) WithAsis(asis bool) *StoreValueCommandBuilder {
	builder.protobuf.Asis = &asis
	return builder
}

// WithSloppyQuorum sets the sloppy_quorum for this Command
// Please note, this is an advanced feature, only use with caution
//
// See http://docs.basho.com/riak/latest/theory/concepts/Eventual-Consistency/
func (builder *StoreValueCommandBuilder) WithSloppyQuorum(sloppyQuorum bool) *StoreValueCommandBuilder {
	builder.protobuf.SloppyQuorum = &sloppyQuorum
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *StoreValueCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &StoreValueCommand{
		value: builder.value,
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
		resolver: builder.resolver}, nil
}

// DeleteValue
// RpbDelReq
// RpbDelResp

// DeleteValueCommand is used to delete a value from Riak KV.
type DeleteValueCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response bool
	protobuf *rpbRiakKV.RpbDelReq
}

// Name identifies this command
func (cmd *DeleteValueCommand) Name() string {
	return cmd.getName("DeleteValue")
}

func (cmd *DeleteValueCommand) constructPbRequest() (msg proto.Message, err error) {
	msg = cmd.protobuf
	return
}

func (cmd *DeleteValueCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	cmd.Response = true
	return nil
}

func (cmd *DeleteValueCommand) getRequestCode() byte {
	return rpbCode_RpbDelReq
}

func (cmd *DeleteValueCommand) getResponseCode() byte {
	return rpbCode_RpbDelResp
}

func (cmd *DeleteValueCommand) getResponseProtobufMessage() proto.Message {
	return nil
}

// DeleteValueCommandBuilder type is required for creating new instances of DeleteValueCommand
//
//	deleteValue := NewDeleteValueCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		WithVClock(vclock).
//		Build()
type DeleteValueCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakKV.RpbDelReq
}

// NewDeleteValueCommandBuilder is a factory function for generating the command builder struct
func NewDeleteValueCommandBuilder() *DeleteValueCommandBuilder {
	builder := &DeleteValueCommandBuilder{protobuf: &rpbRiakKV.RpbDelReq{}}
	return builder
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *DeleteValueCommandBuilder) WithBucketType(bucketType string) *DeleteValueCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *DeleteValueCommandBuilder) WithBucket(bucket string) *DeleteValueCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *DeleteValueCommandBuilder) WithKey(key string) *DeleteValueCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithVClock sets the vector clock.
//
// If not set siblings may be created depending on bucket properties.
func (builder *DeleteValueCommandBuilder) WithVClock(vclock []byte) *DeleteValueCommandBuilder {
	builder.protobuf.Vclock = vclock
	return builder
}

// WithR sets the number of nodes that must report back a successful read in order for the
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *DeleteValueCommandBuilder) WithR(r uint32) *DeleteValueCommandBuilder {
	builder.protobuf.R = &r
	return builder
}

// WithW sets the number of nodes that must report back a successful write in order for then
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *DeleteValueCommandBuilder) WithW(w uint32) *DeleteValueCommandBuilder {
	builder.protobuf.W = &w
	return builder
}

// WithPr sets the number of primary nodes (N) that must be read from in order for the command
// operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *DeleteValueCommandBuilder) WithPr(pr uint32) *DeleteValueCommandBuilder {
	builder.protobuf.Pr = &pr
	return builder
}

// WithPw sets the number of primary nodes (N) that must report back a successful write in order for
// the command operation to be considered a success by Riak.  If ommitted, the bucket default is
// used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *DeleteValueCommandBuilder) WithPw(pw uint32) *DeleteValueCommandBuilder {
	builder.protobuf.Pw = &pw
	return builder
}

// WithDw (durable writes) sets the number of nodes that must report back a successful write to
// backend storage in order for the command operation to be considered a success by Riak. If
// ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *DeleteValueCommandBuilder) WithDw(dw uint32) *DeleteValueCommandBuilder {
	builder.protobuf.Dw = &dw
	return builder
}

// WithRw (delete quorum) sets the number of nodes that must report back a successful delete to
// backend storage in order for the command operation to be considered a success by Riak. It
// represents the read and write operations that are completed internal to Riak to complete a delete.
// If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *DeleteValueCommandBuilder) WithRw(rw uint32) *DeleteValueCommandBuilder {
	builder.protobuf.Rw = &rw
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *DeleteValueCommandBuilder) WithTimeout(timeout time.Duration) *DeleteValueCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *DeleteValueCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &DeleteValueCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// ListBuckets
// RpbListBucketsReq
// RpbListBucketsResp

// ListBucketsCommand is used to list buckets in a bucket type
type ListBucketsCommand struct {
	commandImpl
	listingImpl
	Response *ListBucketsResponse
	protobuf *rpbRiakKV.RpbListBucketsReq
	callback func(buckets []string) error
	done     bool
}

// Name identifies this command
func (cmd *ListBucketsCommand) Name() string {
	return cmd.getName("ListBuckets")
}

func (cmd *ListBucketsCommand) isDone() bool {
	if cmd.protobuf.GetStream() {
		return cmd.done
	}

	return true
}

func (cmd *ListBucketsCommand) constructPbRequest() (msg proto.Message, err error) {
	msg = cmd.protobuf
	return
}

func (cmd *ListBucketsCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg == nil {
		cmd.done = true
		cmd.Response = &ListBucketsResponse{}
	} else {
		if rpbListBucketsResp, ok := msg.(*rpbRiakKV.RpbListBucketsResp); ok {
			cmd.done = rpbListBucketsResp.GetDone()
			response := cmd.Response
			if response == nil {
				response = &ListBucketsResponse{}
				cmd.Response = response
			}
			if rpbListBucketsResp.GetBuckets() != nil {
				buckets := make([]string, len(rpbListBucketsResp.GetBuckets()))
				for i, bucket := range rpbListBucketsResp.GetBuckets() {
					buckets[i] = string(bucket)
				}
				if cmd.protobuf.GetStream() {
					if cmd.callback == nil {
						panic("ListBucketsCommand requires a callback when streaming.")
					} else {
						if err := cmd.callback(buckets); err != nil {
							cmd.Response = nil
							return err
						}
					}
				} else {
					response.Buckets = append(response.Buckets, buckets...)
				}
			}
		} else {
			cmd.done = true
			return fmt.Errorf("[ListBucketsCommand] could not convert %v to RpbListBucketsResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *ListBucketsCommand) getRequestCode() byte {
	return rpbCode_RpbListBucketsReq
}

func (cmd *ListBucketsCommand) getResponseCode() byte {
	return rpbCode_RpbListBucketsResp
}

func (cmd *ListBucketsCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbListBucketsResp{}
}

// ListBucketsResponse contains the response data for a ListBucketsCommand
type ListBucketsResponse struct {
	Buckets []string
}

// ListBucketsCommandBuilder type is required for creating new instances of ListBucketsCommand
//
//	cb := func(buckets []string) error {
//		// Do something with the result
//		return nil
//	}
//	cmd, err := NewListBucketsCommandBuilder().
//		WithBucketType("myBucketType").
//		WithStreaming(true).
//		WithCallback(cb).
//		Build()
type ListBucketsCommandBuilder struct {
	allowListing bool
	callback     func(buckets []string) error
	protobuf     *rpbRiakKV.RpbListBucketsReq
}

// NewListBucketsCommandBuilder is a factory function for generating the command builder struct
func NewListBucketsCommandBuilder() *ListBucketsCommandBuilder {
	builder := &ListBucketsCommandBuilder{protobuf: &rpbRiakKV.RpbListBucketsReq{}}
	return builder
}

// WithAllowListing will allow this command to be built and execute
func (builder *ListBucketsCommandBuilder) WithAllowListing() *ListBucketsCommandBuilder {
	builder.allowListing = true
	return builder
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *ListBucketsCommandBuilder) WithBucketType(bucketType string) *ListBucketsCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithStreaming sets the command to provide a streamed response
//
// If true, a callback must be provided via WithCallback()
func (builder *ListBucketsCommandBuilder) WithStreaming(streaming bool) *ListBucketsCommandBuilder {
	builder.protobuf.Stream = &streaming
	return builder
}

// WithCallback sets the callback to be used when handling a streaming response
//
// Requires WithStreaming(true)
func (builder *ListBucketsCommandBuilder) WithCallback(callback func([]string) error) *ListBucketsCommandBuilder {
	builder.callback = callback
	return builder
}

// WithTimeout sets a timeout in milliseconds to be used for this command operation
func (builder *ListBucketsCommandBuilder) WithTimeout(timeout time.Duration) *ListBucketsCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *ListBucketsCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	if builder.protobuf.GetStream() && builder.callback == nil {
		return nil, newClientError("ListBucketsCommand requires a callback when streaming.", nil)
	}
	if !builder.allowListing {
		return nil, ErrListingDisabled
	}
	return &ListBucketsCommand{
		listingImpl: listingImpl{
			allowListing: builder.allowListing,
		},
		protobuf: builder.protobuf,
		callback: builder.callback}, nil
}

// ListKeys
// RpbListKeysReq
// RpbListKeysResp

// ListKeysCommand is used to fetch a list of keys within a bucket from Riak KV
type ListKeysCommand struct {
	commandImpl
	timeoutImpl
	listingImpl
	Response  *ListKeysResponse
	protobuf  *rpbRiakKV.RpbListKeysReq
	streaming bool
	callback  func(keys []string) error
	done      bool
}

// Name identifies this command
func (cmd *ListKeysCommand) Name() string {
	return cmd.getName("ListKeys")
}

func (cmd *ListKeysCommand) isDone() bool {
	// NB: RpbListKeysReq is *always* streaming so no need to take
	// cmd.streaming into account here, unlike RpbListBucketsReq
	return cmd.done
}

func (cmd *ListKeysCommand) constructPbRequest() (msg proto.Message, err error) {
	msg = cmd.protobuf
	return
}

func (cmd *ListKeysCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg == nil {
		cmd.done = true
		cmd.Response = &ListKeysResponse{}
	} else {
		if rpbListKeysResp, ok := msg.(*rpbRiakKV.RpbListKeysResp); ok {
			cmd.done = rpbListKeysResp.GetDone()
			response := cmd.Response
			if response == nil {
				response = &ListKeysResponse{}
				cmd.Response = response
			}
			if rpbListKeysResp.GetKeys() != nil {
				keys := make([]string, len(rpbListKeysResp.GetKeys()))
				for i, key := range rpbListKeysResp.GetKeys() {
					keys[i] = string(key)
				}
				if cmd.streaming {
					if cmd.callback == nil {
						panic("ListKeysCommand requires a callback when streaming.")
					} else {
						if err := cmd.callback(keys); err != nil {
							cmd.Response = nil
							return err
						}
					}
				} else {
					response.Keys = append(response.Keys, keys...)
				}
			}
		} else {
			cmd.done = true
			return fmt.Errorf("[ListKeysCommand] could not convert %v to RpbListKeysResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *ListKeysCommand) getRequestCode() byte {
	return rpbCode_RpbListKeysReq
}

func (cmd *ListKeysCommand) getResponseCode() byte {
	return rpbCode_RpbListKeysResp
}

func (cmd *ListKeysCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbListKeysResp{}
}

// ListKeysResponse contains the response data for a ListKeysCommand
type ListKeysResponse struct {
	Keys []string
}

// ListKeysCommandBuilder type is required for creating new instances of ListKeysCommand
//
//	cb := func(buckets []string) error {
//		// Do something with the result
//		return nil
//	}
//	cmd, err := NewListKeysCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithStreaming(true).
//		WithCallback(cb).
//		Build()
type ListKeysCommandBuilder struct {
	allowListing bool
	timeout      time.Duration
	protobuf     *rpbRiakKV.RpbListKeysReq
	streaming    bool
	callback     func(buckets []string) error
}

// NewListKeysCommandBuilder is a factory function for generating the command builder struct
func NewListKeysCommandBuilder() *ListKeysCommandBuilder {
	builder := &ListKeysCommandBuilder{protobuf: &rpbRiakKV.RpbListKeysReq{}}
	return builder
}

// WithAllowListing will allow this command to be built and execute
func (builder *ListKeysCommandBuilder) WithAllowListing() *ListKeysCommandBuilder {
	builder.allowListing = true
	return builder
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *ListKeysCommandBuilder) WithBucketType(bucketType string) *ListKeysCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *ListKeysCommandBuilder) WithBucket(bucket string) *ListKeysCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithStreaming sets the command to provide a streamed response
//
// If true, a callback must be provided via WithCallback()
func (builder *ListKeysCommandBuilder) WithStreaming(streaming bool) *ListKeysCommandBuilder {
	builder.streaming = streaming
	return builder
}

// WithCallback sets the callback to be used when handling a streaming response
//
// Requires WithStreaming(true)
func (builder *ListKeysCommandBuilder) WithCallback(callback func([]string) error) *ListKeysCommandBuilder {
	builder.callback = callback
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *ListKeysCommandBuilder) WithTimeout(timeout time.Duration) *ListKeysCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *ListKeysCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	if builder.streaming && builder.callback == nil {
		return nil, newClientError("ListKeysCommand requires a callback when streaming.", nil)
	}
	if !builder.allowListing {
		return nil, ErrListingDisabled
	}
	return &ListKeysCommand{
		listingImpl: listingImpl{
			allowListing: builder.allowListing,
		},
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf:  builder.protobuf,
		streaming: builder.streaming,
		callback:  builder.callback,
	}, nil
}

// FetchPreflist
// RpbGetBucketKeyPreflistReq
// RpbGetBucketKeyPreflistResp

// FetchPreflistCommand is used to fetch the preference list for a key from Riak KV
type FetchPreflistCommand struct {
	commandImpl
	retryableCommandImpl
	Response *FetchPreflistResponse
	protobuf *rpbRiakKV.RpbGetBucketKeyPreflistReq
}

// Name identifies this command
func (cmd *FetchPreflistCommand) Name() string {
	return cmd.getName("FetchPreflist")
}

func (cmd *FetchPreflistCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchPreflistCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg == nil {
		cmd.Response = &FetchPreflistResponse{}
	} else {
		if rpbGetBucketKeyPreflistResp, ok := msg.(*rpbRiakKV.RpbGetBucketKeyPreflistResp); ok {
			response := &FetchPreflistResponse{}
			if rpbGetBucketKeyPreflistResp.GetPreflist() != nil {
				rpbPreflist := rpbGetBucketKeyPreflistResp.GetPreflist()
				response.Preflist = make([]*PreflistItem, len(rpbPreflist))
				for i, rpbItem := range rpbPreflist {
					response.Preflist[i] = &PreflistItem{
						Partition: rpbItem.GetPartition(),
						Node:      string(rpbItem.GetNode()),
						Primary:   rpbItem.GetPrimary(),
					}
				}
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[FetchPreflistCommand] could not convert %v to RpbGetBucketKeyPreflistResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchPreflistCommand) getRequestCode() byte {
	return rpbCode_RpbGetBucketKeyPreflistReq
}

func (cmd *FetchPreflistCommand) getResponseCode() byte {
	return rpbCode_RpbGetBucketKeyPreflistResp
}

func (cmd *FetchPreflistCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbGetBucketKeyPreflistResp{}
}

// PreflistItem represents an individual result from the FetchPreflistResponse result set
type PreflistItem struct {
	Partition int64
	Node      string
	Primary   bool
}

// FetchPreflistResponse contains the response data for a FetchPreflistCommand
type FetchPreflistResponse struct {
	Preflist []*PreflistItem
}

// FetchPreflistCommandBuilder type is required for creating new instances of FetchPreflistCommand
//
//	preflist, err := NewFetchPreflistCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		Build()
type FetchPreflistCommandBuilder struct {
	protobuf *rpbRiakKV.RpbGetBucketKeyPreflistReq
}

// NewFetchPreflistCommandBuilder is a factory function for generating the command builder struct
func NewFetchPreflistCommandBuilder() *FetchPreflistCommandBuilder {
	builder := &FetchPreflistCommandBuilder{protobuf: &rpbRiakKV.RpbGetBucketKeyPreflistReq{}}
	return builder
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *FetchPreflistCommandBuilder) WithBucketType(bucketType string) *FetchPreflistCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *FetchPreflistCommandBuilder) WithBucket(bucket string) *FetchPreflistCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *FetchPreflistCommandBuilder) WithKey(key string) *FetchPreflistCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchPreflistCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &FetchPreflistCommand{protobuf: builder.protobuf}, nil
}

// SecondaryIndexQuery
// RpbGetBucketKeyPreflistReq
// RpbGetBucketKeyPreflistResp

// SecondaryIndexQueryCommand is used to query for keys from Riak KV using secondary indexes
type SecondaryIndexQueryCommand struct {
	commandImpl
	timeoutImpl
	Response *SecondaryIndexQueryResponse
	protobuf *rpbRiakKV.RpbIndexReq
	callback func([]*SecondaryIndexQueryResult) error
	done     bool
}

func (cmd *SecondaryIndexQueryCommand) isDone() bool {
	if cmd.protobuf.GetStream() {
		return cmd.done
	}

	return true
}

// Name identifies this command
func (cmd *SecondaryIndexQueryCommand) Name() string {
	return cmd.getName("SecondaryIndexQuery")
}

func (cmd *SecondaryIndexQueryCommand) constructPbRequest() (proto.Message, error) {
	if cmd.protobuf.GetKey() != nil {
		cmd.protobuf.Qtype = rpbRiakKV.RpbIndexReq_eq.Enum()
	} else {
		cmd.protobuf.Qtype = rpbRiakKV.RpbIndexReq_range.Enum()
	}
	return cmd.protobuf, nil
}

func (cmd *SecondaryIndexQueryCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg == nil {
		cmd.Response = &SecondaryIndexQueryResponse{}
		cmd.done = true
	} else {
		if rpbIndexResp, ok := msg.(*rpbRiakKV.RpbIndexResp); ok {
			cmd.done = rpbIndexResp.GetDone()
			response := cmd.Response
			if response == nil {
				response = &SecondaryIndexQueryResponse{}
				cmd.Response = response
			}

			response.Continuation = rpbIndexResp.GetContinuation()

			var results []*SecondaryIndexQueryResult
			rpbIndexRespResultsLen := len(rpbIndexResp.GetResults())
			if rpbIndexRespResultsLen > 0 {
				// Index keys and object keys were returned
				results = make([]*SecondaryIndexQueryResult, rpbIndexRespResultsLen)
				for i, rpbIndexResult := range rpbIndexResp.GetResults() {
					results[i] = &SecondaryIndexQueryResult{
						IndexKey:  rpbIndexResult.Key,
						ObjectKey: rpbIndexResult.Value,
					}
				}
			} else {
				// Only object keys were returned
				var key []byte
				if cmd.protobuf.GetReturnTerms() {
					// this is only possible if this was a single key query
					key = cmd.protobuf.GetKey()
				}
				rpbIndexRespKeys := rpbIndexResp.GetKeys()
				results = make([]*SecondaryIndexQueryResult, len(rpbIndexRespKeys))
				for i, rpbIndexKey := range rpbIndexRespKeys {
					results[i] = &SecondaryIndexQueryResult{
						IndexKey:  key,
						ObjectKey: rpbIndexKey,
					}
				}
			}

			if cmd.protobuf.GetStream() {
				if cmd.callback == nil {
					panic("SecondaryIndexQueryCommand requires a callback when streaming.")
				} else {
					if err := cmd.callback(results); err != nil {
						cmd.Response = nil
						return err
					}
				}
			} else {
				response.Results = append(response.Results, results...)
			}
		} else {
			cmd.done = true
			return fmt.Errorf("[SecondaryIndexQueryCommand] could not convert %v to RpbIndexResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *SecondaryIndexQueryCommand) getRequestCode() byte {
	return rpbCode_RpbIndexReq
}

func (cmd *SecondaryIndexQueryCommand) getResponseCode() byte {
	return rpbCode_RpbIndexResp
}

func (cmd *SecondaryIndexQueryCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbIndexResp{}
}

// SecondaryIndexQueryResult represents an individual result of the SecondaryIndexQueryResponse
// result set
type SecondaryIndexQueryResult struct {
	IndexKey  []byte
	ObjectKey []byte
}

// SecondaryIndexQueryResponse contains the response data for a SecondaryIndexQueryCommand
type SecondaryIndexQueryResponse struct {
	Results      []*SecondaryIndexQueryResult
	Continuation []byte
}

// SecondaryIndexQueryCommandBuilder type is required for creating new instances of SecondaryIndexQueryCommand
//
//	command, err := NewSecondaryIndexQueryCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithIndexName("myIndexName").
//		WithIndexKey("myIndexKey").
//		WithIntIndexKey(1234).
//		Build()
type SecondaryIndexQueryCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakKV.RpbIndexReq
	callback func([]*SecondaryIndexQueryResult) error
}

// NewSecondaryIndexQueryCommandBuilder is a factory function for generating the command builder struct
func NewSecondaryIndexQueryCommandBuilder() *SecondaryIndexQueryCommandBuilder {
	builder := &SecondaryIndexQueryCommandBuilder{protobuf: &rpbRiakKV.RpbIndexReq{}}
	return builder
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *SecondaryIndexQueryCommandBuilder) WithBucketType(bucketType string) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *SecondaryIndexQueryCommandBuilder) WithBucket(bucket string) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithIndexName sets the index to use for the command
func (builder *SecondaryIndexQueryCommandBuilder) WithIndexName(indexName string) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.Index = []byte(indexName)
	return builder
}

// WithRange sets the range of index values to return
func (builder *SecondaryIndexQueryCommandBuilder) WithRange(min string, max string) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.RangeMin = []byte(min)
	builder.protobuf.RangeMax = []byte(max)
	return builder
}

// WithIntRange sets the range of integer type index values to return, useful when you want 1,3,5,11
// and not 1,11,3,5
func (builder *SecondaryIndexQueryCommandBuilder) WithIntRange(min int64, max int64) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.RangeMin = []byte(strconv.FormatInt(min, 10))
	builder.protobuf.RangeMax = []byte(strconv.FormatInt(max, 10))
	return builder
}

// WithIndexKey defines the index to search against
func (builder *SecondaryIndexQueryCommandBuilder) WithIndexKey(key string) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithIntIndexKey defines the integer index to search against
func (builder *SecondaryIndexQueryCommandBuilder) WithIntIndexKey(key int) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.Key = []byte(strconv.Itoa(key))
	return builder
}

// WithReturnKeyAndIndex set to true, the result set will include both index keys and object keys
func (builder *SecondaryIndexQueryCommandBuilder) WithReturnKeyAndIndex(val bool) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.ReturnTerms = &val
	return builder
}

// WithStreaming sets the command to provide a streamed response
//
// If true, a callback must be provided via WithCallback()
func (builder *SecondaryIndexQueryCommandBuilder) WithStreaming(streaming bool) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.Stream = &streaming
	return builder
}

// WithCallback sets the callback to be used when handling a streaming response
//
// Requires WithStreaming(true)
func (builder *SecondaryIndexQueryCommandBuilder) WithCallback(callback func([]*SecondaryIndexQueryResult) error) *SecondaryIndexQueryCommandBuilder {
	builder.callback = callback
	return builder
}

// WithPaginationSort set to true, the results of a non-paginated query will return sorted from Riak
func (builder *SecondaryIndexQueryCommandBuilder) WithPaginationSort(paginationSort bool) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.PaginationSort = &paginationSort
	return builder
}

// WithMaxResults sets the maximum number of values to return in the result set
func (builder *SecondaryIndexQueryCommandBuilder) WithMaxResults(maxResults uint32) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.MaxResults = &maxResults
	return builder
}

// WithContinuation sets the position at which the result set should continue from, value can be
// found within the result set of the previous page for the same query
func (builder *SecondaryIndexQueryCommandBuilder) WithContinuation(cont []byte) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.Continuation = cont
	return builder
}

// WithTermRegex sets the regex pattern to filter the result set by
func (builder *SecondaryIndexQueryCommandBuilder) WithTermRegex(regex string) *SecondaryIndexQueryCommandBuilder {
	builder.protobuf.TermRegex = []byte(regex)
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *SecondaryIndexQueryCommandBuilder) WithTimeout(timeout time.Duration) *SecondaryIndexQueryCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *SecondaryIndexQueryCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	if builder.protobuf.GetKey() == nil &&
		(builder.protobuf.GetRangeMin() == nil || builder.protobuf.GetRangeMax() == nil) {
		return nil, newClientError("either WithIndexKey or WithRange are required", nil)
	}
	if builder.protobuf.GetStream() && builder.callback == nil {
		return nil, newClientError("SecondaryIndexQueryCommand requires a callback when streaming.", nil)
	}
	return &SecondaryIndexQueryCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
		callback: builder.callback,
	}, nil
}

// MapReduce
// RpbMapRedReq
// RpbMapRedResp

// MapReduceCommand is used to fetch keys or data from Riak KV using the MapReduce technique
type MapReduceCommand struct {
	commandImpl
	Response  [][]byte
	protobuf  *rpbRiakKV.RpbMapRedReq
	streaming bool
	callback  func(response []byte) error
	done      bool
}

// Name identifies this command
func (cmd *MapReduceCommand) Name() string {
	return cmd.getName("MapReduce")
}

func (cmd *MapReduceCommand) isDone() bool {
	// NB: RpbMapRedReq is *always* streaming so no need to take
	// cmd.streaming into account here, unlike RpbListBucketsReq
	return cmd.done
}

func (cmd *MapReduceCommand) constructPbRequest() (msg proto.Message, err error) {
	msg = cmd.protobuf
	return
}

func (cmd *MapReduceCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg == nil {
		cmd.done = true
	} else {
		if rpbMapRedResp, ok := msg.(*rpbRiakKV.RpbMapRedResp); ok {
			cmd.done = rpbMapRedResp.GetDone()
			rpbMapRedRespData := rpbMapRedResp.GetResponse()
			if cmd.streaming {
				if cmd.callback == nil {
					panic("MapReduceCommand requires a callback when streaming.")
				} else {
					if err := cmd.callback(rpbMapRedRespData); err != nil {
						cmd.Response = nil
						return err
					}
				}
			} else {
				cmd.Response = append(cmd.Response, rpbMapRedRespData)
			}
		} else {
			cmd.done = true
			return fmt.Errorf("[MapReduceCommand] could not convert %v to RpbMapRedResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *MapReduceCommand) getRequestCode() byte {
	return rpbCode_RpbMapRedReq
}

func (cmd *MapReduceCommand) getResponseCode() byte {
	return rpbCode_RpbMapRedResp
}

func (cmd *MapReduceCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakKV.RpbMapRedResp{}
}

// MapReduceCommandBuilder type is required for creating new instances of MapReduceCommand
//
//	command, err := NewMapReduceCommandBuilder().
//		WithQuery("myMapReduceQuery").
//		Build()
type MapReduceCommandBuilder struct {
	protobuf  *rpbRiakKV.RpbMapRedReq
	streaming bool
	callback  func(response []byte) error
}

// NewMapReduceCommandBuilder is a factory function for generating the command builder struct
func NewMapReduceCommandBuilder() *MapReduceCommandBuilder {
	return &MapReduceCommandBuilder{
		protobuf: &rpbRiakKV.RpbMapRedReq{
			ContentType: []byte("application/json"),
		},
	}
}

// WithQuery sets the map reduce query to be executed on Riak
func (builder *MapReduceCommandBuilder) WithQuery(query string) *MapReduceCommandBuilder {
	builder.protobuf.Request = []byte(query)
	return builder
}

// WithStreaming sets the command to provide a streamed response
//
// If true, a callback must be provided via WithCallback()
func (builder *MapReduceCommandBuilder) WithStreaming(streaming bool) *MapReduceCommandBuilder {
	builder.streaming = streaming
	return builder
}

// WithCallback sets the callback to be used when handling a streaming response
//
// Requires WithStreaming(true)
func (builder *MapReduceCommandBuilder) WithCallback(callback func([]byte) error) *MapReduceCommandBuilder {
	builder.callback = callback
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *MapReduceCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if builder.streaming && builder.callback == nil {
		return nil, newClientError("MapReduceCommand requires a callback when streaming.", nil)
	}
	return &MapReduceCommand{
		protobuf:  builder.protobuf,
		streaming: builder.streaming,
		callback:  builder.callback,
	}, nil
}
