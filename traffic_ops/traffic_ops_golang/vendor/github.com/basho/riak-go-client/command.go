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
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"

	proto "github.com/golang/protobuf/proto"
)

// Global private var used in debug mode to differentiate command names in debug output
var c uint64 = 0

// Interface implemented by Command types that can be re-tried
type retryableCommand interface {
	setLastNode(*Node)
	getLastNode() *Node
}

// Implementation of retryableCommand
type retryableCommandImpl struct {
	lastNode *Node
}

func (cmd *retryableCommandImpl) setLastNode(lastNode *Node) {
	if lastNode == nil {
		panic("[retryableCommandImpl] nil last node")
	}
	cmd.lastNode = lastNode
}

func (cmd *retryableCommandImpl) getLastNode() *Node {
	return cmd.lastNode
}

type commandImpl struct {
	error   error
	success bool
	name    string
}

func (cmd *commandImpl) Success() bool {
	return cmd.success == true
}

func (cmd *commandImpl) Error() error {
	return cmd.error
}

func (cmd *commandImpl) onError(err error) {
	cmd.success = false
	cmd.error = err
}

func (cmd *commandImpl) onRetry() {
	cmd.error = nil
}

func (cmd *commandImpl) getName(n string) string {
	if n == "" {
		panic("getName: n must not be empty")
	}
	if cmd.name == "" {
		if EnableDebugLogging == true {
			cmd.name = fmt.Sprintf("%s-%v", n, atomic.AddUint64(&c, 1))
		} else {
			cmd.name = n
		}
	}
	return cmd.name
}

// Interface implemented by Command types that can be streamed
type streamingCommand interface {
	isDone() bool
}

// Interface implemented by Command types that have a timeout
type timeoutCommand interface {
	getTimeout() time.Duration
}

type timeoutImpl struct {
	timeout time.Duration
}

func (cmd *timeoutImpl) getTimeout() time.Duration {
	return cmd.timeout
}

// Interface implemented by Commands that list data from Riak
type listingCommand interface {
	getAllowListing() bool
}

type listingImpl struct {
	allowListing bool
}

func (cmd *listingImpl) getAllowListing() bool {
	return cmd.allowListing
}

// CommandBuilder interface requires Build() method for generating the Command
// to be executed
type CommandBuilder interface {
	Build() (Command, error)
}

// Command interface enforces proper structure of a Command object
type Command interface {
	Name() string
	Success() bool
	Error() error
	getRequestCode() byte
	constructPbRequest() (proto.Message, error)
	onRetry()
	onError(error)
	onSuccess(proto.Message) error // NB: important for streaming commands to "do the right thing" here
	getResponseCode() byte
	getResponseProtobufMessage() proto.Message
}

func getRiakMessage(cmd Command) (msg []byte, err error) {
	requestCode := cmd.getRequestCode()
	if requestCode == 0 {
		panic(fmt.Sprintf("Must have non-zero value for getRequestCode(): %s", cmd.Name()))
	}

	var rpb proto.Message
	rpb, err = cmd.constructPbRequest()
	if err != nil {
		return
	}

	var bytes []byte
	if rpb != nil {
		bytes, err = proto.Marshal(rpb)
		if err != nil {
			return nil, err
		}
	}

	msg = buildRiakMessage(requestCode, bytes)
	return
}

func decodeRiakMessage(cmd Command, data []byte) (msg proto.Message, err error) {
	responseCode := cmd.getResponseCode()
	if responseCode == 0 {
		panic(fmt.Sprintf("Must have non-zero value for getResponseCode(): %s", cmd.Name()))
	}

	err = rpbValidateResp(data, responseCode)
	if err != nil {
		return
	}

	if len(data) > 1 {
		msg = cmd.getResponseProtobufMessage()
		if msg != nil {
			err = proto.Unmarshal(data[1:], msg)
		}
	}

	return
}

func buildRiakMessage(code byte, data []byte) []byte {
	buf := new(bytes.Buffer)
	// write total message length, including one byte for msg code
	binary.Write(buf, binary.BigEndian, uint32(len(data)+1))
	// write the message code
	binary.Write(buf, binary.BigEndian, byte(code))
	// write the protobuf data
	buf.Write(data)
	return buf.Bytes()
}
