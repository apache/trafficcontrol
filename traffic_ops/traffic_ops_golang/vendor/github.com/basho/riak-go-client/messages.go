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

// Convert from csv with:
// %s/\(\d\+\),\([^,]\+\),.*/const rpbCode_\2 byte = \1/
const rpbCode_RpbErrorResp byte = 0
const rpbCode_RpbPingReq byte = 1
const rpbCode_RpbPingResp byte = 2
const rpbCode_RpbGetClientIdReq byte = 3
const rpbCode_RpbGetClientIdResp byte = 4
const rpbCode_RpbSetClientIdReq byte = 5
const rpbCode_RpbSetClientIdResp byte = 6
const rpbCode_RpbGetServerInfoReq byte = 7
const rpbCode_RpbGetServerInfoResp byte = 8
const rpbCode_RpbGetReq byte = 9
const rpbCode_RpbGetResp byte = 10
const rpbCode_RpbPutReq byte = 11
const rpbCode_RpbPutResp byte = 12
const rpbCode_RpbDelReq byte = 13
const rpbCode_RpbDelResp byte = 14
const rpbCode_RpbListBucketsReq byte = 15
const rpbCode_RpbListBucketsResp byte = 16
const rpbCode_RpbListKeysReq byte = 17
const rpbCode_RpbListKeysResp byte = 18
const rpbCode_RpbGetBucketReq byte = 19
const rpbCode_RpbGetBucketResp byte = 20
const rpbCode_RpbSetBucketReq byte = 21
const rpbCode_RpbSetBucketResp byte = 22
const rpbCode_RpbMapRedReq byte = 23
const rpbCode_RpbMapRedResp byte = 24
const rpbCode_RpbIndexReq byte = 25
const rpbCode_RpbIndexResp byte = 26
const rpbCode_RpbSearchQueryReq byte = 27
const rpbCode_RpbSearchQueryResp byte = 28
const rpbCode_RpbResetBucketReq byte = 29
const rpbCode_RpbResetBucketResp byte = 30
const rpbCode_RpbGetBucketTypeReq byte = 31
const rpbCode_RpbSetBucketTypeReq byte = 32
const rpbCode_RpbGetBucketKeyPreflistReq byte = 33
const rpbCode_RpbGetBucketKeyPreflistResp byte = 34
const rpbCode_RpbCSBucketReq byte = 40
const rpbCode_RpbCSBucketResp byte = 41
const rpbCode_RpbCounterUpdateReq byte = 50
const rpbCode_RpbCounterUpdateResp byte = 51
const rpbCode_RpbCounterGetReq byte = 52
const rpbCode_RpbCounterGetResp byte = 53
const rpbCode_RpbYokozunaIndexGetReq byte = 54
const rpbCode_RpbYokozunaIndexGetResp byte = 55
const rpbCode_RpbYokozunaIndexPutReq byte = 56
const rpbCode_RpbYokozunaIndexDeleteReq byte = 57
const rpbCode_RpbYokozunaSchemaGetReq byte = 58
const rpbCode_RpbYokozunaSchemaGetResp byte = 59
const rpbCode_RpbYokozunaSchemaPutReq byte = 60
const rpbCode_DtFetchReq byte = 80
const rpbCode_DtFetchResp byte = 81
const rpbCode_DtUpdateReq byte = 82
const rpbCode_DtUpdateResp byte = 83
const rpbCode_TsQueryReq byte = 90
const rpbCode_TsQueryResp byte = 91
const rpbCode_TsPutReq byte = 92
const rpbCode_TsPutResp byte = 93
const rpbCode_TsDelReq byte = 94
const rpbCode_TsDelResp byte = 95
const rpbCode_TsGetReq byte = 96
const rpbCode_TsGetResp byte = 97
const rpbCode_TsListKeysReq byte = 98
const rpbCode_TsListKeysResp byte = 99
const rpbCode_RpbAuthReq byte = 253
const rpbCode_RpbAuthResp byte = 254
const rpbCode_RpbStartTls byte = 255
