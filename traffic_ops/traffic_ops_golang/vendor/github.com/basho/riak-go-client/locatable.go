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
	proto "github.com/golang/protobuf/proto"
)

type rpbLocatable interface {
	GetType() []byte
	SetType(bt []byte) // NB: bt == bucket type
	BucketIsRequired() bool
	GetBucket() []byte
	KeyIsRequired() bool
	GetKey() []byte
}

func validateLocatable(msg proto.Message) error {
	l := msg.(rpbLocatable)
	if l.BucketIsRequired() {
		if bucket := l.GetBucket(); len(bucket) == 0 {
			return ErrBucketRequired
		}
	}
	if l.KeyIsRequired() {
		if key := l.GetKey(); len(key) == 0 {
			return ErrKeyRequired
		}
	}
	if bucketType := l.GetType(); len(bucketType) == 0 {
		l.SetType([]byte(defaultBucketType))
	}
	return nil
}
