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

package riak_dt

// DtUpdateReq

func (m *DtUpdateReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *DtUpdateReq) BucketIsRequired() bool {
	return true
}

func (m *DtUpdateReq) KeyIsRequired() bool {
	return false
}

// DtFetchReq

func (m *DtFetchReq) SetType(bt []byte) {
	m.Type = bt
}

func (m *DtFetchReq) BucketIsRequired() bool {
	return true
}

func (m *DtFetchReq) KeyIsRequired() bool {
	return true
}
