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

// RpbGetBucketTypeReq

func (m *RpbGetBucketTypeReq) GetKey() []byte {
	return nil
}

func (m *RpbGetBucketTypeReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbGetBucketTypeReq) BucketIsRequired() bool {
	return false
}

func (m *RpbGetBucketTypeReq) GetBucket() []byte {
	return nil
}

func (m *RpbGetBucketTypeReq) KeyIsRequired() bool {
	return false
}

// RpbGetBucketReq

func (m *RpbGetBucketReq) GetKey() []byte {
	return nil
}

func (m *RpbGetBucketReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbGetBucketReq) BucketIsRequired() bool {
	return true
}

func (m *RpbGetBucketReq) KeyIsRequired() bool {
	return false
}

// RpbSetBucketTypeReq

func (m *RpbSetBucketTypeReq) GetKey() []byte {
	return nil
}

func (m *RpbSetBucketTypeReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbSetBucketTypeReq) BucketIsRequired() bool {
	return false
}

func (m *RpbSetBucketTypeReq) GetBucket() []byte {
	return nil
}

func (m *RpbSetBucketTypeReq) KeyIsRequired() bool {
	return false
}

// RpbSetBucketReq

func (m *RpbSetBucketReq) GetKey() []byte {
	return nil
}

func (m *RpbSetBucketReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbSetBucketReq) BucketIsRequired() bool {
	return true
}

func (m *RpbSetBucketReq) KeyIsRequired() bool {
	return false
}

// RpbResetBucketReq

func (m *RpbResetBucketReq) GetKey() []byte {
	return nil
}

func (m *RpbResetBucketReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *RpbResetBucketReq) BucketIsRequired() bool {
	return true
}

func (m *RpbResetBucketReq) KeyIsRequired() bool {
	return false
}
