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
	"time"

	rpbRiakDT "github.com/basho/riak-go-client/rpb/riak_dt"
	rpbRiakKV "github.com/basho/riak-go-client/rpb/riak_kv"
	proto "github.com/golang/protobuf/proto"
)

// UpdateCounter
// DtUpdateReq
// DtUpdateResp

// UpdateCounterCommand is used to increment or decrement a counter data type in Riak KV
type UpdateCounterCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *UpdateCounterResponse
	protobuf proto.Message
	isLegacy bool
}

// Name identifies this command
func (cmd *UpdateCounterCommand) Name() string {
	return cmd.getName("UpdateCounter")
}

func (cmd *UpdateCounterCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *UpdateCounterCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		// For legacy counters, the response may be different
		if rpbDtUpdateResp, is_DtUpdateResp := msg.(*rpbRiakDT.DtUpdateResp); is_DtUpdateResp && !cmd.isLegacy {
			cmd.Response = &UpdateCounterResponse{
				GeneratedKey: string(rpbDtUpdateResp.GetKey()),
				CounterValue: rpbDtUpdateResp.GetCounterValue(),
			}
		} else if rpbCounterUpdateResp, is_RpbCounterUpdateResp := msg.(*rpbRiakKV.RpbCounterUpdateResp); is_RpbCounterUpdateResp && cmd.isLegacy {
			cmd.Response = &UpdateCounterResponse{
				CounterValue: rpbCounterUpdateResp.GetValue(),
			}
		} else {
			return fmt.Errorf("[UpdateCounterCommand] could not convert %v to DtUpdateResp / RpbCounterUpdateResp, isLegacy: %v", reflect.TypeOf(msg), cmd.isLegacy)
		}
	}
	return nil
}

func (cmd *UpdateCounterCommand) getRequestCode() byte {
	if cmd.isLegacy {
		return rpbCode_RpbCounterUpdateReq
	}
	return rpbCode_DtUpdateReq
}

func (cmd *UpdateCounterCommand) getResponseCode() byte {
	if cmd.isLegacy {
		return rpbCode_RpbCounterUpdateResp
	}
	return rpbCode_DtUpdateResp
}

func (cmd *UpdateCounterCommand) getResponseProtobufMessage() proto.Message {
	if cmd.isLegacy {
		return &rpbRiakKV.RpbCounterUpdateResp{}
	}
	return &rpbRiakDT.DtUpdateResp{}
}

// UpdateCounterResponse is the object containing the response
type UpdateCounterResponse struct {
	GeneratedKey string
	CounterValue int64
}

// UpdateCounterCommandBuilder type is required for creating new instances of UpdateCounterCommand
//
//	command, err := NewUpdateCounterCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		WithIncrement(1).
//		Build()
type UpdateCounterCommandBuilder struct {
	bucketType string
	bucket     string
	key        string
	increment  int64
	w          uint32
	dw         uint32
	pw         uint32
	returnBody bool
	timeout    time.Duration
}

// NewUpdateCounterCommandBuilder is a factory function for generating the command builder struct
func NewUpdateCounterCommandBuilder() *UpdateCounterCommandBuilder {
	return &UpdateCounterCommandBuilder{}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *UpdateCounterCommandBuilder) WithBucketType(bucketType string) *UpdateCounterCommandBuilder {
	builder.bucketType = bucketType
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *UpdateCounterCommandBuilder) WithBucket(bucket string) *UpdateCounterCommandBuilder {
	builder.bucket = bucket
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *UpdateCounterCommandBuilder) WithKey(key string) *UpdateCounterCommandBuilder {
	builder.key = key
	return builder
}

// WithIncrement defines the increment the Counter value is to be increased / decreased by
func (builder *UpdateCounterCommandBuilder) WithIncrement(increment int64) *UpdateCounterCommandBuilder {
	builder.increment = increment
	return builder
}

// WithW sets the number of nodes that must report back a successful write in order for then
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateCounterCommandBuilder) WithW(w uint32) *UpdateCounterCommandBuilder {
	builder.w = w
	return builder
}

// WithPw sets the number of primary nodes (N) that must report back a successful write in order for
// the command operation to be considered a success by Riak.  If ommitted, the bucket default is
// used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateCounterCommandBuilder) WithPw(pw uint32) *UpdateCounterCommandBuilder {
	builder.pw = pw
	return builder
}

// WithDw (durable writes) sets the number of nodes that must report back a successful write to
// backend storage in order for the command operation to be considered a success by Riak
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateCounterCommandBuilder) WithDw(dw uint32) *UpdateCounterCommandBuilder {
	builder.dw = dw
	return builder
}

// WithReturnBody sets Riak to return the value within its response after completing the write
// operation
func (builder *UpdateCounterCommandBuilder) WithReturnBody(returnBody bool) *UpdateCounterCommandBuilder {
	builder.returnBody = returnBody
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *UpdateCounterCommandBuilder) WithTimeout(timeout time.Duration) *UpdateCounterCommandBuilder {
	builder.timeout = timeout
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *UpdateCounterCommandBuilder) Build() (Command, error) {
	var isLegacy = false
	var timeout time.Duration
	var protobuf proto.Message = nil
	if builder.bucketType == defaultBucketType && builder.returnBody == true {
		isLegacy = true
		rpbCounterUpdateReq := &rpbRiakKV.RpbCounterUpdateReq{
			Amount:      &builder.increment,
			W:           &builder.w,
			Dw:          &builder.dw,
			Pw:          &builder.pw,
			Returnvalue: &builder.returnBody,
		}
		// NB: strings must be handled this way to ensure that nil slices
		// are in the PB msg, rather than 0-len ones
		if builder.bucket != "" {
			rpbCounterUpdateReq.Bucket = []byte(builder.bucket)
		}
		if builder.key != "" {
			rpbCounterUpdateReq.Key = []byte(builder.key)
		}
		protobuf = rpbCounterUpdateReq
	} else {
		timeout = builder.timeout
		timeoutMilliseconds := uint32(builder.timeout / time.Millisecond)
		dtUpdateReq := &rpbRiakDT.DtUpdateReq{
			W:          &builder.w,
			Dw:         &builder.dw,
			Pw:         &builder.pw,
			ReturnBody: &builder.returnBody,
			Timeout:    &timeoutMilliseconds,
			Op: &rpbRiakDT.DtOp{
				CounterOp: &rpbRiakDT.CounterOp{
					Increment: &builder.increment,
				},
			},
		}
		// NB: strings must be handled this way to ensure that nil slices
		// are in the PB msg, rather than 0-len ones
		if builder.bucketType != "" {
			dtUpdateReq.Type = []byte(builder.bucketType)
		}
		if builder.bucket != "" {
			dtUpdateReq.Bucket = []byte(builder.bucket)
		}
		if builder.key != "" {
			dtUpdateReq.Key = []byte(builder.key)
		}
		protobuf = dtUpdateReq
	}

	if err := validateLocatable(protobuf); err != nil {
		return nil, err
	}

	return &UpdateCounterCommand{
		timeoutImpl: timeoutImpl{
			timeout: timeout,
		},
		protobuf: protobuf,
		isLegacy: isLegacy,
	}, nil
}

// FetchCounter
// DtFetchReq
// DtFetchResp

// FetchCounterCommand fetches a counter CRDT from Riak
type FetchCounterCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *FetchCounterResponse
	protobuf *rpbRiakDT.DtFetchReq
}

// Name identifies this command
func (cmd *FetchCounterCommand) Name() string {
	return cmd.getName("FetchCounter")
}

func (cmd *FetchCounterCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchCounterCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtFetchResp, ok := msg.(*rpbRiakDT.DtFetchResp); ok {
			response := &FetchCounterResponse{}
			rpbValue := rpbDtFetchResp.GetValue()
			if rpbValue == nil {
				response.IsNotFound = true
			} else {
				response.CounterValue = rpbValue.GetCounterValue()
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[FetchCounterCommand] could not convert %v to DtFetchResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchCounterCommand) getRequestCode() byte {
	return rpbCode_DtFetchReq
}

func (cmd *FetchCounterCommand) getResponseCode() byte {
	return rpbCode_DtFetchResp
}

func (cmd *FetchCounterCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtFetchResp{}
}

// FetchCounterResponse contains the response data for a FetchCounterCommand
type FetchCounterResponse struct {
	IsNotFound   bool
	CounterValue int64
}

// FetchCounterCommandBuilder type is required for creating new instances of FetchCounterCommand
//
//	command, err := NewFetchCounterCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		Build()
type FetchCounterCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakDT.DtFetchReq
}

// NewFetchCounterCommandBuilder is a factory function for generating the command builder struct
func NewFetchCounterCommandBuilder() *FetchCounterCommandBuilder {
	return &FetchCounterCommandBuilder{protobuf: &rpbRiakDT.DtFetchReq{}}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *FetchCounterCommandBuilder) WithBucketType(bucketType string) *FetchCounterCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *FetchCounterCommandBuilder) WithBucket(bucket string) *FetchCounterCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *FetchCounterCommandBuilder) WithKey(key string) *FetchCounterCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithR sets the number of nodes that must report back a successful read in order for the
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchCounterCommandBuilder) WithR(r uint32) *FetchCounterCommandBuilder {
	builder.protobuf.R = &r
	return builder
}

// WithPr sets the number of primary nodes (N) that must be read from in order for the command
// operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchCounterCommandBuilder) WithPr(pr uint32) *FetchCounterCommandBuilder {
	builder.protobuf.Pr = &pr
	return builder
}

// WithNotFoundOk sets notfound_ok, whether to treat notfounds as successful reads for the purposes
// of R
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchCounterCommandBuilder) WithNotFoundOk(notFoundOk bool) *FetchCounterCommandBuilder {
	builder.protobuf.NotfoundOk = &notFoundOk
	return builder
}

// WithBasicQuorum sets basic_quorum, whether to return early in some failure cases (eg. when r=1
// and you get 2 errors and a success basic_quorum=true would return an error)
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchCounterCommandBuilder) WithBasicQuorum(basicQuorum bool) *FetchCounterCommandBuilder {
	builder.protobuf.BasicQuorum = &basicQuorum
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *FetchCounterCommandBuilder) WithTimeout(timeout time.Duration) *FetchCounterCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchCounterCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &FetchCounterCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// UpdateSet
// DtUpdateReq
// DtUpdateResp

// UpdateSetCommand stores or updates a set CRDT in Riak
type UpdateSetCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *UpdateSetResponse
	protobuf *rpbRiakDT.DtUpdateReq
}

// Name identifies this command
func (cmd *UpdateSetCommand) Name() string {
	return cmd.getName("UpdateSet")
}

func (cmd *UpdateSetCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *UpdateSetCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtUpdateResp, ok := msg.(*rpbRiakDT.DtUpdateResp); ok {
			response := &UpdateSetResponse{
				GeneratedKey: string(rpbDtUpdateResp.GetKey()),
				Context:      rpbDtUpdateResp.GetContext(),
				SetValue:     rpbDtUpdateResp.GetSetValue(),
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[UpdateSetCommand] could not convert %v to DtUpdateResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *UpdateSetCommand) getRequestCode() byte {
	return rpbCode_DtUpdateReq
}

func (cmd *UpdateSetCommand) getResponseCode() byte {
	return rpbCode_DtUpdateResp
}

func (cmd *UpdateSetCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtUpdateResp{}
}

// UpdateSetResponse contains the response data for a UpdateSetCommand
type UpdateSetResponse struct {
	GeneratedKey string
	Context      []byte
	SetValue     [][]byte
}

// UpdateSetCommandBuilder type is required for creating new instances of UpdateSetCommand
//
//	adds := [][]byte{
//		[]byte("a1"),
//		[]byte("a2"),
//		[]byte("a3"),
//		[]byte("a4"),
//	}
//
//	command, err := NewUpdateSetCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		WithContext(setContext).
//		WithAdditions(adds).
//		Build()
type UpdateSetCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakDT.DtUpdateReq
}

// NewUpdateSetCommandBuilder is a factory function for generating the command builder struct
func NewUpdateSetCommandBuilder() *UpdateSetCommandBuilder {
	return &UpdateSetCommandBuilder{
		protobuf: &rpbRiakDT.DtUpdateReq{
			Op: &rpbRiakDT.DtOp{
				SetOp: &rpbRiakDT.SetOp{},
			},
		},
	}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *UpdateSetCommandBuilder) WithBucketType(bucketType string) *UpdateSetCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *UpdateSetCommandBuilder) WithBucket(bucket string) *UpdateSetCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *UpdateSetCommandBuilder) WithKey(key string) *UpdateSetCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithContext sets the causal context needed to identify the state of the set when removing elements
func (builder *UpdateSetCommandBuilder) WithContext(context []byte) *UpdateSetCommandBuilder {
	builder.protobuf.Context = context
	return builder
}

// WithAdditions sets the set elements to be added to the CRDT set via this update operation
func (builder *UpdateSetCommandBuilder) WithAdditions(adds ...[]byte) *UpdateSetCommandBuilder {
	opAdds := builder.protobuf.Op.SetOp.Adds
	opAdds = append(opAdds, adds...)
	builder.protobuf.Op.SetOp.Adds = opAdds
	return builder
}

// WithRemovals sets the set elements to be removed from the CRDT set via this update operation
func (builder *UpdateSetCommandBuilder) WithRemovals(removals ...[]byte) *UpdateSetCommandBuilder {
	opRemoves := builder.protobuf.Op.SetOp.Removes
	opRemoves = append(opRemoves, removals...)
	builder.protobuf.Op.SetOp.Removes = opRemoves
	return builder
}

// WithW sets the number of nodes that must report back a successful write in order for then
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateSetCommandBuilder) WithW(w uint32) *UpdateSetCommandBuilder {
	builder.protobuf.W = &w
	return builder
}

// WithPw sets the number of primary nodes (N) that must report back a successful write in order for
// the command operation to be considered a success by Riak.  If ommitted, the bucket default is
// used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateSetCommandBuilder) WithPw(pw uint32) *UpdateSetCommandBuilder {
	builder.protobuf.Pw = &pw
	return builder
}

// WithDw (durable writes) sets the number of nodes that must report back a successful write to
// backend storage in order for the command operation to be considered a success by Riak. If
// ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateSetCommandBuilder) WithDw(dw uint32) *UpdateSetCommandBuilder {
	builder.protobuf.Dw = &dw
	return builder
}

// WithReturnBody sets Riak to return the value within its response after completing the write
// operation
func (builder *UpdateSetCommandBuilder) WithReturnBody(returnBody bool) *UpdateSetCommandBuilder {
	builder.protobuf.ReturnBody = &returnBody
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *UpdateSetCommandBuilder) WithTimeout(timeout time.Duration) *UpdateSetCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *UpdateSetCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &UpdateSetCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// UpdateGSet
// DtUpdateReq
// DtUpdateResp

// UpdateGSetCommand stores or updates a set CRDT in Riak
type UpdateGSetCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *UpdateGSetResponse
	protobuf *rpbRiakDT.DtUpdateReq
}

// Name identifies this command
func (cmd *UpdateGSetCommand) Name() string {
	return cmd.getName("UpdateGSet")
}

func (cmd *UpdateGSetCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *UpdateGSetCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtUpdateResp, ok := msg.(*rpbRiakDT.DtUpdateResp); ok {
			response := &UpdateGSetResponse{
				GeneratedKey: string(rpbDtUpdateResp.GetKey()),
				Context:      rpbDtUpdateResp.GetContext(),
				GSetValue:    rpbDtUpdateResp.GetGsetValue(),
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[UpdateGSetCommand] could not convert %v to DtUpdateResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *UpdateGSetCommand) getRequestCode() byte {
	return rpbCode_DtUpdateReq
}

func (cmd *UpdateGSetCommand) getResponseCode() byte {
	return rpbCode_DtUpdateResp
}

func (cmd *UpdateGSetCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtUpdateResp{}
}

// UpdateGSetResponse contains the response data for a UpdateGSetCommand
type UpdateGSetResponse struct {
	GeneratedKey string
	Context      []byte
	GSetValue    [][]byte
}

// UpdateGSetCommandBuilder type is required for creating new instances of UpdateGSetCommand
//
//	adds := [][]byte{
//		[]byte("a1"),
//		[]byte("a2"),
//		[]byte("a3"),
//		[]byte("a4"),
//	}
//
//	command := NewUpdateGSetCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		WithContext(setContext).
//		WithAdditions(adds).
//		Build()
type UpdateGSetCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakDT.DtUpdateReq
}

// NewUpdateGSetCommandBuilder is a factory function for generating the command builder struct
func NewUpdateGSetCommandBuilder() *UpdateGSetCommandBuilder {
	return &UpdateGSetCommandBuilder{
		protobuf: &rpbRiakDT.DtUpdateReq{
			Op: &rpbRiakDT.DtOp{
				GsetOp: &rpbRiakDT.GSetOp{},
			},
		},
	}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *UpdateGSetCommandBuilder) WithBucketType(bucketType string) *UpdateGSetCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *UpdateGSetCommandBuilder) WithBucket(bucket string) *UpdateGSetCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *UpdateGSetCommandBuilder) WithKey(key string) *UpdateGSetCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithContext sets the causal context needed to identify the state of the set when removing elements
func (builder *UpdateGSetCommandBuilder) WithContext(context []byte) *UpdateGSetCommandBuilder {
	builder.protobuf.Context = context
	return builder
}

// WithAdditions sets the set elements to be added to the CRDT set via this update operation
func (builder *UpdateGSetCommandBuilder) WithAdditions(adds ...[]byte) *UpdateGSetCommandBuilder {
	opAdds := builder.protobuf.Op.GsetOp.Adds
	opAdds = append(opAdds, adds...)
	builder.protobuf.Op.GsetOp.Adds = opAdds
	return builder
}

// WithW sets the number of nodes that must report back a successful write in order for then
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateGSetCommandBuilder) WithW(w uint32) *UpdateGSetCommandBuilder {
	builder.protobuf.W = &w
	return builder
}

// WithPw sets the number of primary nodes (N) that must report back a successful write in order for
// the command operation to be considered a success by Riak.  If ommitted, the bucket default is
// used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateGSetCommandBuilder) WithPw(pw uint32) *UpdateGSetCommandBuilder {
	builder.protobuf.Pw = &pw
	return builder
}

// WithDw (durable writes) sets the number of nodes that must report back a successful write to
// backend storage in order for the command operation to be considered a success by Riak. If
// ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateGSetCommandBuilder) WithDw(dw uint32) *UpdateGSetCommandBuilder {
	builder.protobuf.Dw = &dw
	return builder
}

// WithReturnBody sets Riak to return the value within its response after completing the write
// operation
func (builder *UpdateGSetCommandBuilder) WithReturnBody(returnBody bool) *UpdateGSetCommandBuilder {
	builder.protobuf.ReturnBody = &returnBody
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *UpdateGSetCommandBuilder) WithTimeout(timeout time.Duration) *UpdateGSetCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *UpdateGSetCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &UpdateGSetCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// FetchSet
// DtFetchReq
// DtFetchResp

// FetchSetCommand fetches a set CRDT from Riak
type FetchSetCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *FetchSetResponse
	protobuf *rpbRiakDT.DtFetchReq
}

// Name identifies this command
func (cmd *FetchSetCommand) Name() string {
	return cmd.getName("FetchSet")
}

func (cmd *FetchSetCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchSetCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtFetchResp, ok := msg.(*rpbRiakDT.DtFetchResp); ok {
			response := &FetchSetResponse{
				Context: rpbDtFetchResp.GetContext(),
			}
			rpbValue := rpbDtFetchResp.GetValue()
			if rpbValue == nil {
				response.IsNotFound = true
			} else {
				rpbType := rpbDtFetchResp.GetType()
				switch rpbType {
				case rpbRiakDT.DtFetchResp_SET:
					response.SetValue = rpbValue.GetSetValue()
				case rpbRiakDT.DtFetchResp_GSET:
					response.SetValue = rpbValue.GetGsetValue()
				}
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[FetchSetCommand] could not convert %v to DtFetchResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchSetCommand) getRequestCode() byte {
	return rpbCode_DtFetchReq
}

func (cmd *FetchSetCommand) getResponseCode() byte {
	return rpbCode_DtFetchResp
}

func (cmd *FetchSetCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtFetchResp{}
}

// FetchSetResponse contains the response data for a FetchSetCommand
type FetchSetResponse struct {
	IsNotFound bool
	Context    []byte
	SetValue   [][]byte
}

// FetchSetCommandBuilder type is required for creating new instances of FetchSetCommand
//
//	command, err := NewFetchSetCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		Build()
type FetchSetCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakDT.DtFetchReq
}

// NewFetchSetCommandBuilder is a factory function for generating the command builder struct
func NewFetchSetCommandBuilder() *FetchSetCommandBuilder {
	return &FetchSetCommandBuilder{protobuf: &rpbRiakDT.DtFetchReq{}}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *FetchSetCommandBuilder) WithBucketType(bucketType string) *FetchSetCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *FetchSetCommandBuilder) WithBucket(bucket string) *FetchSetCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *FetchSetCommandBuilder) WithKey(key string) *FetchSetCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithR sets the number of nodes that must report back a successful read in order for the
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchSetCommandBuilder) WithR(r uint32) *FetchSetCommandBuilder {
	builder.protobuf.R = &r
	return builder
}

// WithPr sets the number of primary nodes (N) that must be read from in order for the command
// operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchSetCommandBuilder) WithPr(pr uint32) *FetchSetCommandBuilder {
	builder.protobuf.Pr = &pr
	return builder
}

// WithNotFoundOk sets notfound_ok, whether to treat notfounds as successful reads for the purposes
// of R
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchSetCommandBuilder) WithNotFoundOk(notFoundOk bool) *FetchSetCommandBuilder {
	builder.protobuf.NotfoundOk = &notFoundOk
	return builder
}

// WithBasicQuorum sets basic_quorum, whether to return early in some failure cases (eg. when r=1
// and you get 2 errors and a success basic_quorum=true would return an error)
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchSetCommandBuilder) WithBasicQuorum(basicQuorum bool) *FetchSetCommandBuilder {
	builder.protobuf.BasicQuorum = &basicQuorum
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *FetchSetCommandBuilder) WithTimeout(timeout time.Duration) *FetchSetCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchSetCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &FetchSetCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// UpdateMap
// DtUpdateReq
// DtUpdateResp

// UpdateMapCommand updates a map CRDT in Riak
type UpdateMapCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *UpdateMapResponse
	op       *MapOperation
	protobuf *rpbRiakDT.DtUpdateReq
}

// Name identifies this command
func (cmd *UpdateMapCommand) Name() string {
	return cmd.getName("UpdateMap")
}

func (cmd *UpdateMapCommand) constructPbRequest() (proto.Message, error) {
	pbMapOp := &rpbRiakDT.MapOp{}
	populate(cmd.op, pbMapOp)

	cmd.protobuf.Op = &rpbRiakDT.DtOp{
		MapOp: pbMapOp,
	}
	return cmd.protobuf, nil
}

func (cmd *UpdateMapCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtUpdateResp, ok := msg.(*rpbRiakDT.DtUpdateResp); ok {
			response := &UpdateMapResponse{
				GeneratedKey: string(rpbDtUpdateResp.GetKey()),
				Context:      rpbDtUpdateResp.GetContext(),
				Map:          parsePbResponse(rpbDtUpdateResp.GetMapValue()),
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[UpdateMapCommand] could not convert %v to DtUpdateResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *UpdateMapCommand) getRequestCode() byte {
	return rpbCode_DtUpdateReq
}

func (cmd *UpdateMapCommand) getResponseCode() byte {
	return rpbCode_DtUpdateResp
}

func (cmd *UpdateMapCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtUpdateResp{}
}

func addMapUpdate(pbMapOp *rpbRiakDT.MapOp, update *rpbRiakDT.MapUpdate) {
	pbMapOp.Updates = append(pbMapOp.Updates, update)
}

func addMapRemove(pbMapOp *rpbRiakDT.MapOp, field *rpbRiakDT.MapField) {
	pbMapOp.Removes = append(pbMapOp.Removes, field)
}

func populate(mapOp *MapOperation, pbMapOp *rpbRiakDT.MapOp) {
	if mapOp.hasRemoves(false) {
		for name := range mapOp.removeCounters {
			field := &rpbRiakDT.MapField{
				Name: []byte(name),
				Type: rpbRiakDT.MapField_COUNTER.Enum(),
			}
			addMapRemove(pbMapOp, field)
		}
		for name := range mapOp.removeSets {
			field := &rpbRiakDT.MapField{
				Name: []byte(name),
				Type: rpbRiakDT.MapField_SET.Enum(),
			}
			addMapRemove(pbMapOp, field)
		}
		for name := range mapOp.removeMaps {
			field := &rpbRiakDT.MapField{
				Name: []byte(name),
				Type: rpbRiakDT.MapField_MAP.Enum(),
			}
			addMapRemove(pbMapOp, field)
		}
		for name := range mapOp.removeRegisters {
			field := &rpbRiakDT.MapField{
				Name: []byte(name),
				Type: rpbRiakDT.MapField_REGISTER.Enum(),
			}
			addMapRemove(pbMapOp, field)
		}
		for name := range mapOp.removeFlags {
			field := &rpbRiakDT.MapField{
				Name: []byte(name),
				Type: rpbRiakDT.MapField_FLAG.Enum(),
			}
			addMapRemove(pbMapOp, field)
		}
	}

	for name, increment := range mapOp.incrementCounters {
		i := increment
		field := &rpbRiakDT.MapField{
			Name: []byte(name),
			Type: rpbRiakDT.MapField_COUNTER.Enum(),
		}
		counterOp := &rpbRiakDT.CounterOp{
			Increment: &i,
		}
		update := &rpbRiakDT.MapUpdate{
			Field:     field,
			CounterOp: counterOp,
		}
		addMapUpdate(pbMapOp, update)
	}
	for name, adds := range mapOp.addToSets {
		field := &rpbRiakDT.MapField{
			Name: []byte(name),
			Type: rpbRiakDT.MapField_SET.Enum(),
		}
		setOp := &rpbRiakDT.SetOp{
			Adds: make([][]byte, len(adds)),
		}
		for i, add := range adds {
			setOp.Adds[i] = add
		}
		update := &rpbRiakDT.MapUpdate{
			Field: field,
			SetOp: setOp,
		}
		addMapUpdate(pbMapOp, update)
	}
	for name, removes := range mapOp.removeFromSets {
		field := &rpbRiakDT.MapField{
			Name: []byte(name),
			Type: rpbRiakDT.MapField_SET.Enum(),
		}
		setOp := &rpbRiakDT.SetOp{
			Removes: make([][]byte, len(removes)),
		}
		for i, remove := range removes {
			setOp.Removes[i] = remove
		}
		update := &rpbRiakDT.MapUpdate{
			Field: field,
			SetOp: setOp,
		}
		addMapUpdate(pbMapOp, update)
	}
	for name, register := range mapOp.registersToSet {
		field := &rpbRiakDT.MapField{
			Name: []byte(name),
			Type: rpbRiakDT.MapField_REGISTER.Enum(),
		}
		update := &rpbRiakDT.MapUpdate{
			Field:      field,
			RegisterOp: register,
		}
		addMapUpdate(pbMapOp, update)
	}
	for name, flag := range mapOp.flagsToSet {
		field := &rpbRiakDT.MapField{
			Name: []byte(name),
			Type: rpbRiakDT.MapField_FLAG.Enum(),
		}
		var flagOp rpbRiakDT.MapUpdate_FlagOp
		if flag {
			flagOp = rpbRiakDT.MapUpdate_ENABLE
		} else {
			flagOp = rpbRiakDT.MapUpdate_DISABLE
		}
		update := &rpbRiakDT.MapUpdate{
			Field:  field,
			FlagOp: flagOp.Enum(),
		}
		addMapUpdate(pbMapOp, update)
	}
	for name, mapOp := range mapOp.maps {
		field := &rpbRiakDT.MapField{
			Name: []byte(name),
			Type: rpbRiakDT.MapField_MAP.Enum(),
		}
		nestedMapOp := &rpbRiakDT.MapOp{}
		populate(mapOp, nestedMapOp)
		update := &rpbRiakDT.MapUpdate{
			Field: field,
			MapOp: nestedMapOp,
		}
		addMapUpdate(pbMapOp, update)
	}
}

// MapOperation contains the instructions to send to Riak what updates to the Map you want to complete
type MapOperation struct {
	incrementCounters map[string]int64
	removeCounters    map[string]bool

	addToSets      map[string][][]byte
	removeFromSets map[string][][]byte
	removeSets     map[string]bool

	registersToSet  map[string][]byte
	removeRegisters map[string]bool

	flagsToSet  map[string]bool
	removeFlags map[string]bool

	maps       map[string]*MapOperation
	removeMaps map[string]bool
}

// IncrementCounter increments a child counter CRDT of the map at the specified key
func (mapOp *MapOperation) IncrementCounter(key string, increment int64) *MapOperation {
	if mapOp.removeCounters != nil {
		delete(mapOp.removeCounters, key)
	}
	if mapOp.incrementCounters == nil {
		mapOp.incrementCounters = make(map[string]int64)
	}
	mapOp.incrementCounters[key] += increment
	return mapOp
}

// RemoveCounter removes a child counter CRDT from the map at the specified key
func (mapOp *MapOperation) RemoveCounter(key string) *MapOperation {
	if mapOp.incrementCounters != nil {
		delete(mapOp.incrementCounters, key)
	}
	if mapOp.removeCounters == nil {
		mapOp.removeCounters = make(map[string]bool)
	}
	mapOp.removeCounters[key] = true
	return mapOp
}

// AddToSet adds an element to the child set CRDT of the map at the specified key
func (mapOp *MapOperation) AddToSet(key string, value []byte) *MapOperation {
	if mapOp.removeSets != nil {
		delete(mapOp.removeSets, key)
	}
	if mapOp.addToSets == nil {
		mapOp.addToSets = make(map[string][][]byte)
	}
	mapOp.addToSets[key] = append(mapOp.addToSets[key], value)
	return mapOp
}

// RemoveFromSet removes elements from the child set CRDT of the map at the specified key
func (mapOp *MapOperation) RemoveFromSet(key string, value []byte) *MapOperation {
	if mapOp.removeSets != nil {
		delete(mapOp.removeSets, key)
	}
	if mapOp.removeFromSets == nil {
		mapOp.removeFromSets = make(map[string][][]byte)
	}
	mapOp.removeFromSets[key] = append(mapOp.removeFromSets[key], value)
	return mapOp
}

// RemoveSet removes the child set CRDT from the map
func (mapOp *MapOperation) RemoveSet(key string) *MapOperation {
	if mapOp.addToSets != nil {
		delete(mapOp.addToSets, key)
	}
	if mapOp.removeFromSets != nil {
		delete(mapOp.removeFromSets, key)
	}
	if mapOp.removeSets == nil {
		mapOp.removeSets = make(map[string]bool)
	}
	mapOp.removeSets[key] = true
	return mapOp
}

// SetRegister sets a register CRDT on the map with the provided value
func (mapOp *MapOperation) SetRegister(key string, value []byte) *MapOperation {
	if mapOp.removeRegisters != nil {
		delete(mapOp.removeRegisters, key)
	}
	if mapOp.registersToSet == nil {
		mapOp.registersToSet = make(map[string][]byte)
	}
	mapOp.registersToSet[key] = value
	return mapOp
}

// RemoveRegister removes a register CRDT from the map
func (mapOp *MapOperation) RemoveRegister(key string) *MapOperation {
	if mapOp.registersToSet != nil {
		delete(mapOp.registersToSet, key)
	}
	if mapOp.removeRegisters == nil {
		mapOp.removeRegisters = make(map[string]bool)
	}
	mapOp.removeRegisters[key] = true
	return mapOp
}

// SetFlag sets a flag CRDT on the map
func (mapOp *MapOperation) SetFlag(key string, value bool) *MapOperation {
	if mapOp.removeFlags != nil {
		delete(mapOp.removeFlags, key)
	}
	if mapOp.flagsToSet == nil {
		mapOp.flagsToSet = make(map[string]bool)
	}
	mapOp.flagsToSet[key] = value
	return mapOp
}

// RemoveFlag removes a flag CRDT from the map
func (mapOp *MapOperation) RemoveFlag(key string) *MapOperation {
	if mapOp.flagsToSet != nil {
		delete(mapOp.flagsToSet, key)
	}
	if mapOp.removeFlags == nil {
		mapOp.removeFlags = make(map[string]bool)
	}
	mapOp.removeFlags[key] = true
	return mapOp
}

// Map returns a nested map operation for manipulation
func (mapOp *MapOperation) Map(key string) *MapOperation {
	if mapOp.removeMaps != nil {
		delete(mapOp.removeMaps, key)
	}
	if mapOp.maps == nil {
		mapOp.maps = make(map[string]*MapOperation)
	}

	innerMapOp, ok := mapOp.maps[key]
	if ok {
		return innerMapOp
	}

	innerMapOp = &MapOperation{}
	mapOp.maps[key] = innerMapOp
	return innerMapOp
}

// RemoveMap removes a nested map from the map
func (mapOp *MapOperation) RemoveMap(key string) *MapOperation {
	if mapOp.maps != nil {
		delete(mapOp.maps, key)
	}
	if mapOp.removeMaps == nil {
		mapOp.removeMaps = make(map[string]bool)
	}
	mapOp.removeMaps[key] = true
	return mapOp
}

func (mapOp *MapOperation) hasRemoves(includeRemoveFromSets bool) bool {
	nestedHaveRemoves := false
	for _, m := range mapOp.maps {
		if m.hasRemoves(false) {
			nestedHaveRemoves = true
			break
		}
	}

	rv := nestedHaveRemoves ||
		len(mapOp.removeCounters) > 0 ||
		len(mapOp.removeSets) > 0 ||
		len(mapOp.removeRegisters) > 0 ||
		len(mapOp.removeFlags) > 0 ||
		len(mapOp.removeMaps) > 0

	if includeRemoveFromSets {
		rv = rv || len(mapOp.removeFromSets) > 0
	}

	return rv
}

func parsePbResponse(pbMapEntries []*rpbRiakDT.MapEntry) *Map {
	m := &Map{}
	for _, mapEntry := range pbMapEntries {
		mapField := mapEntry.GetField()
		key := string(mapField.GetName())
		switch mapField.GetType() {
		case rpbRiakDT.MapField_COUNTER:
			if m.Counters == nil {
				m.Counters = make(map[string]int64)
			}
			m.Counters[key] = mapEntry.GetCounterValue()
		case rpbRiakDT.MapField_SET:
			if m.Sets == nil {
				m.Sets = make(map[string][][]byte)
			}
			m.Sets[key] = mapEntry.SetValue
		case rpbRiakDT.MapField_REGISTER:
			if m.Registers == nil {
				m.Registers = make(map[string][]byte)
			}
			m.Registers[key] = mapEntry.GetRegisterValue()
		case rpbRiakDT.MapField_FLAG:
			if m.Flags == nil {
				m.Flags = make(map[string]bool)
			}
			m.Flags[key] = mapEntry.GetFlagValue()
		case rpbRiakDT.MapField_MAP:
			if m.Maps == nil {
				m.Maps = make(map[string]*Map)
			}
			m.Maps[key] = parsePbResponse(mapEntry.MapValue)
		}
	}
	return m
}

// Map object represents the Riak Map object and is returned within the Response objects for both
// UpdateMapCommand and FetchMapCommand
type Map struct {
	Counters  map[string]int64
	Sets      map[string][][]byte
	Registers map[string][]byte
	Flags     map[string]bool
	Maps      map[string]*Map
}

// UpdateMapResponse contains the response data for a UpdateMapCommand
type UpdateMapResponse struct {
	GeneratedKey string
	Context      []byte
	Map          *Map
}

// UpdateMapCommandBuilder type is required for creating new instances of UpdateMapCommand
//
//	mapOp := &MapOperation{}
//	mapOp.SetRegister("register_1", []byte("register_value_1"))
//
//	command, err := NewUpdateMapCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		WithMapOperation(mapOp).
//		Build()
type UpdateMapCommandBuilder struct {
	mapOperation *MapOperation
	timeout      time.Duration
	protobuf     *rpbRiakDT.DtUpdateReq
}

// NewUpdateMapCommandBuilder is a factory function for generating the command builder struct
func NewUpdateMapCommandBuilder() *UpdateMapCommandBuilder {
	return &UpdateMapCommandBuilder{protobuf: &rpbRiakDT.DtUpdateReq{}}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *UpdateMapCommandBuilder) WithBucketType(bucketType string) *UpdateMapCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *UpdateMapCommandBuilder) WithBucket(bucket string) *UpdateMapCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *UpdateMapCommandBuilder) WithKey(key string) *UpdateMapCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithContext sets the causal context needed to identify the state of the map when removing elements
func (builder *UpdateMapCommandBuilder) WithContext(context []byte) *UpdateMapCommandBuilder {
	builder.protobuf.Context = context
	return builder
}

// WithMapOperation provides the details of what is supposed to be updated on the map
func (builder *UpdateMapCommandBuilder) WithMapOperation(mapOperation *MapOperation) *UpdateMapCommandBuilder {
	builder.mapOperation = mapOperation
	return builder
}

// WithW sets the number of nodes that must report back a successful write in order for then
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateMapCommandBuilder) WithW(w uint32) *UpdateMapCommandBuilder {
	builder.protobuf.W = &w
	return builder
}

// WithPw sets the number of primary nodes (N) that must report back a successful write in order for
// the command operation to be considered a success by Riak.  If ommitted, the bucket default is
// used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateMapCommandBuilder) WithPw(pw uint32) *UpdateMapCommandBuilder {
	builder.protobuf.Pw = &pw
	return builder
}

// WithDw (durable writes) sets the number of nodes that must report back a successful write to
// backend storage in order for the command operation to be considered a success by Riak. If
// ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateMapCommandBuilder) WithDw(dw uint32) *UpdateMapCommandBuilder {
	builder.protobuf.Dw = &dw
	return builder
}

// WithReturnBody sets Riak to return the value within its response after completing the write
// operation
func (builder *UpdateMapCommandBuilder) WithReturnBody(returnBody bool) *UpdateMapCommandBuilder {
	builder.protobuf.ReturnBody = &returnBody
	return builder
}

// WithTimeout sets a timeout in milliseconds to be used for this command operation
func (builder *UpdateMapCommandBuilder) WithTimeout(timeout time.Duration) *UpdateMapCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *UpdateMapCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	if builder.mapOperation == nil {
		return nil, newClientError("UpdateMapCommandBuilder requires non-nil MapOperation. Use WithMapOperation()", nil)
	}
	if builder.mapOperation.hasRemoves(true) && builder.protobuf.GetContext() == nil {
		return nil, newClientError("When doing any removes a context must be provided.", nil)
	}
	return &UpdateMapCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
		op:       builder.mapOperation,
	}, nil
}

// FetchMap
// DtFetchReq
// DtFetchResp

// FetchMapCommand fetches a map CRDT from Riak
type FetchMapCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *FetchMapResponse
	protobuf *rpbRiakDT.DtFetchReq
}

// Name identifies this command
func (cmd *FetchMapCommand) Name() string {
	return cmd.getName("FetchMap")
}

func (cmd *FetchMapCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchMapCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtFetchResp, ok := msg.(*rpbRiakDT.DtFetchResp); ok {
			response := &FetchMapResponse{
				Context: rpbDtFetchResp.GetContext(),
			}
			rpbValue := rpbDtFetchResp.GetValue()
			if rpbValue == nil {
				response.IsNotFound = true
			} else {
				rpbMapValue := rpbValue.GetMapValue()
				if rpbMapValue == nil {
					response.IsNotFound = true
				} else {
					response.Map = parsePbResponse(rpbMapValue)
				}
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[FetchMapCommand] could not convert %v to DtFetchResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchMapCommand) getRequestCode() byte {
	return rpbCode_DtFetchReq
}

func (cmd *FetchMapCommand) getResponseCode() byte {
	return rpbCode_DtFetchResp
}

func (cmd *FetchMapCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtFetchResp{}
}

// FetchMapResponse contains the response data for a FetchMapCommand
type FetchMapResponse struct {
	IsNotFound bool
	Context    []byte
	Map        *Map
}

// FetchMapCommandBuilder type is required for creating new instances of FetchMapCommand
//
//	command, err := NewFetchMapCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		Build()
type FetchMapCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakDT.DtFetchReq
}

// NewFetchMapCommandBuilder is a factory function for generating the command builder struct
func NewFetchMapCommandBuilder() *FetchMapCommandBuilder {
	return &FetchMapCommandBuilder{protobuf: &rpbRiakDT.DtFetchReq{}}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *FetchMapCommandBuilder) WithBucketType(bucketType string) *FetchMapCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *FetchMapCommandBuilder) WithBucket(bucket string) *FetchMapCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *FetchMapCommandBuilder) WithKey(key string) *FetchMapCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithR sets the number of nodes that must report back a successful read in order for the
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchMapCommandBuilder) WithR(r uint32) *FetchMapCommandBuilder {
	builder.protobuf.R = &r
	return builder
}

// WithPr sets the number of primary nodes (N) that must be read from in order for the command
// operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchMapCommandBuilder) WithPr(pr uint32) *FetchMapCommandBuilder {
	builder.protobuf.Pr = &pr
	return builder
}

// WithNotFoundOk sets notfound_ok, whether to treat notfounds as successful reads for the purposes
// of R
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchMapCommandBuilder) WithNotFoundOk(notFoundOk bool) *FetchMapCommandBuilder {
	builder.protobuf.NotfoundOk = &notFoundOk
	return builder
}

// WithBasicQuorum sets basic_quorum, whether to return early in some failure cases (eg. when r=1
// and you get 2 errors and a success basic_quorum=true would return an error)
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchMapCommandBuilder) WithBasicQuorum(basicQuorum bool) *FetchMapCommandBuilder {
	builder.protobuf.BasicQuorum = &basicQuorum
	return builder
}

// WithTimeout sets a timeout in milliseconds to be used for this command operation
func (builder *FetchMapCommandBuilder) WithTimeout(timeout time.Duration) *FetchMapCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchMapCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &FetchMapCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// UpdateHll
// DtUpdateReq
// DtUpdateResp

// UpdateHllCommand stores or updates a set CRDT in Riak
type UpdateHllCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *UpdateHllResponse
	protobuf *rpbRiakDT.DtUpdateReq
}

// Name identifies this command
func (cmd *UpdateHllCommand) Name() string {
	return cmd.getName("UpdateHll")
}

func (cmd *UpdateHllCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *UpdateHllCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtUpdateResp, ok := msg.(*rpbRiakDT.DtUpdateResp); ok {
			response := &UpdateHllResponse{
				GeneratedKey: string(rpbDtUpdateResp.GetKey()),
				Cardinality:  rpbDtUpdateResp.GetHllValue(),
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[UpdateHllCommand] could not convert %v to DtUpdateResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *UpdateHllCommand) getRequestCode() byte {
	return rpbCode_DtUpdateReq
}

func (cmd *UpdateHllCommand) getResponseCode() byte {
	return rpbCode_DtUpdateResp
}

func (cmd *UpdateHllCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtUpdateResp{}
}

// UpdateHllResponse contains the response data for a UpdateHllCommand
type UpdateHllResponse struct {
	GeneratedKey string
	Cardinality  uint64
}

// UpdateHllCommandBuilder type is required for creating new instances of UpdateHllCommand
//
//	adds := [][]byte{
//		[]byte("a1"),
//		[]byte("a2"),
//		[]byte("a3"),
//		[]byte("a4"),
//	}
//
//	command, err := NewUpdateHllCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		WithAdditions(adds).
//		Build()
type UpdateHllCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakDT.DtUpdateReq
}

// NewUpdateHllCommandBuilder is a factory function for generating the command builder struct
func NewUpdateHllCommandBuilder() *UpdateHllCommandBuilder {
	return &UpdateHllCommandBuilder{
		protobuf: &rpbRiakDT.DtUpdateReq{
			Op: &rpbRiakDT.DtOp{
				HllOp: &rpbRiakDT.HllOp{},
			},
		},
	}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *UpdateHllCommandBuilder) WithBucketType(bucketType string) *UpdateHllCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *UpdateHllCommandBuilder) WithBucket(bucket string) *UpdateHllCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *UpdateHllCommandBuilder) WithKey(key string) *UpdateHllCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithAdditions sets the Hll elements to be added to the Hll Data Type via this update operation
func (builder *UpdateHllCommandBuilder) WithAdditions(adds ...[]byte) *UpdateHllCommandBuilder {
	opAdds := builder.protobuf.Op.HllOp.Adds
	opAdds = append(opAdds, adds...)
	builder.protobuf.Op.HllOp.Adds = opAdds
	return builder
}

// WithW sets the number of nodes that must report back a successful write in order for then
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateHllCommandBuilder) WithW(w uint32) *UpdateHllCommandBuilder {
	builder.protobuf.W = &w
	return builder
}

// WithPw sets the number of primary nodes (N) that must report back a successful write in order for
// the command operation to be considered a success by Riak.  If ommitted, the bucket default is
// used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateHllCommandBuilder) WithPw(pw uint32) *UpdateHllCommandBuilder {
	builder.protobuf.Pw = &pw
	return builder
}

// WithDw (durable writes) sets the number of nodes that must report back a successful write to
// backend storage in order for the command operation to be considered a success by Riak. If
// ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *UpdateHllCommandBuilder) WithDw(dw uint32) *UpdateHllCommandBuilder {
	builder.protobuf.Dw = &dw
	return builder
}

// WithReturnBody sets Riak to return the value within its response after completing the write
// operation
func (builder *UpdateHllCommandBuilder) WithReturnBody(returnBody bool) *UpdateHllCommandBuilder {
	builder.protobuf.ReturnBody = &returnBody
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *UpdateHllCommandBuilder) WithTimeout(timeout time.Duration) *UpdateHllCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *UpdateHllCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &UpdateHllCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// FetchHll
// DtFetchReq
// DtFetchResp

// FetchHllCommand fetches an Hll Data Type from Riak
type FetchHllCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *FetchHllResponse
	protobuf *rpbRiakDT.DtFetchReq
}

// Name identifies this command
func (cmd *FetchHllCommand) Name() string {
	return cmd.getName("FetchHll")
}

func (cmd *FetchHllCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchHllCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbDtFetchResp, ok := msg.(*rpbRiakDT.DtFetchResp); ok {
			response := &FetchHllResponse{}
			rpbValue := rpbDtFetchResp.GetValue()
			if rpbValue == nil {
				response.IsNotFound = true
			} else {
				response.Cardinality = rpbValue.GetHllValue()
			}
			cmd.Response = response
		} else {
			return fmt.Errorf("[FetchHllCommand] could not convert %v to DtFetchResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchHllCommand) getRequestCode() byte {
	return rpbCode_DtFetchReq
}

func (cmd *FetchHllCommand) getResponseCode() byte {
	return rpbCode_DtFetchResp
}

func (cmd *FetchHllCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakDT.DtFetchResp{}
}

// FetchHllResponse contains the response data for a FetchHllCommand
type FetchHllResponse struct {
	IsNotFound  bool
	Cardinality uint64
}

// FetchHllCommandBuilder type is required for creating new instances of FetchHllCommand
//
//	command, err := NewFetchHllCommandBuilder().
//		WithBucketType("myBucketType").
//		WithBucket("myBucket").
//		WithKey("myKey").
//		Build()
type FetchHllCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakDT.DtFetchReq
}

// NewFetchHllCommandBuilder is a factory function for generating the command builder struct
func NewFetchHllCommandBuilder() *FetchHllCommandBuilder {
	return &FetchHllCommandBuilder{protobuf: &rpbRiakDT.DtFetchReq{}}
}

// WithBucketType sets the bucket-type to be used by the command. If omitted, 'default' is used
func (builder *FetchHllCommandBuilder) WithBucketType(bucketType string) *FetchHllCommandBuilder {
	builder.protobuf.Type = []byte(bucketType)
	return builder
}

// WithBucket sets the bucket to be used by the command
func (builder *FetchHllCommandBuilder) WithBucket(bucket string) *FetchHllCommandBuilder {
	builder.protobuf.Bucket = []byte(bucket)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *FetchHllCommandBuilder) WithKey(key string) *FetchHllCommandBuilder {
	builder.protobuf.Key = []byte(key)
	return builder
}

// WithR sets the number of nodes that must report back a successful read in order for the
// command operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchHllCommandBuilder) WithR(r uint32) *FetchHllCommandBuilder {
	builder.protobuf.R = &r
	return builder
}

// WithPr sets the number of primary nodes (N) that must be read from in order for the command
// operation to be considered a success by Riak. If ommitted, the bucket default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *FetchHllCommandBuilder) WithPr(pr uint32) *FetchHllCommandBuilder {
	builder.protobuf.Pr = &pr
	return builder
}

// WithNotFoundOk sets notfound_ok, whether to treat notfounds as successful reads for the purposes
// of R
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchHllCommandBuilder) WithNotFoundOk(notFoundOk bool) *FetchHllCommandBuilder {
	builder.protobuf.NotfoundOk = &notFoundOk
	return builder
}

// WithBasicQuorum sets basic_quorum, whether to return early in some failure cases (eg. when r=1
// and you get 2 errors and a success basic_quorum=true would return an error)
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-3/
func (builder *FetchHllCommandBuilder) WithBasicQuorum(basicQuorum bool) *FetchHllCommandBuilder {
	builder.protobuf.BasicQuorum = &basicQuorum
	return builder
}

// WithTimeout sets a timeout to be used for this command operation
func (builder *FetchHllCommandBuilder) WithTimeout(timeout time.Duration) *FetchHllCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchHllCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if err := validateLocatable(builder.protobuf); err != nil {
		return nil, err
	}
	return &FetchHllCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}
