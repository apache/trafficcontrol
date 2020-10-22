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

package riak_kv

// RpbGetReq

func (m *RpbGetReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbGetReq) BucketIsRequired() bool {
	return true
}

func (m *RpbGetReq) KeyIsRequired() bool {
	return true
}

// RpbPutReq

func (m *RpbPutReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbPutReq) BucketIsRequired() bool {
	return true
}

func (m *RpbPutReq) KeyIsRequired() bool {
	return false
}

// RpbDelReq

func (m *RpbDelReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbDelReq) BucketIsRequired() bool {
	return true
}

func (m *RpbDelReq) KeyIsRequired() bool {
	return true
}

// RpbListBucketsReq

func (m *RpbListBucketsReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbListBucketsReq) BucketIsRequired() bool {
	return false
}

func (m *RpbListBucketsReq) GetBucket() []byte {
	return nil
}

func (m *RpbListBucketsReq) KeyIsRequired() bool {
	return false
}

func (m *RpbListBucketsReq) GetKey() []byte {
	return nil
}

// RpbListKeysReq

func (m *RpbListKeysReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbListKeysReq) BucketIsRequired() bool {
	return true
}

func (m *RpbListKeysReq) KeyIsRequired() bool {
	return false
}

func (m *RpbListKeysReq) GetKey() []byte {
	return nil
}

// RpbGetBucketKeyPreflistReq

func (m *RpbGetBucketKeyPreflistReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbGetBucketKeyPreflistReq) BucketIsRequired() bool {
	return true
}

func (m *RpbGetBucketKeyPreflistReq) KeyIsRequired() bool {
	return true
}

// RpbIndexReq

func (m *RpbIndexReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbIndexReq) BucketIsRequired() bool {
	return true
}

func (m *RpbIndexReq) KeyIsRequired() bool {
	return false
}

// RpbCounterUpdateReq

func (m *RpbCounterUpdateReq) SetType(bt []byte) {
}

func (m *RpbCounterUpdateReq) GetType() []byte {
	return nil
}

func (m *RpbCounterUpdateReq) BucketIsRequired() bool {
	return true
}

func (m *RpbCounterUpdateReq) KeyIsRequired() bool {
	return true
}
