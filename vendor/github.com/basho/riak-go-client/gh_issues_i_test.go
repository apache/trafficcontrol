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

// +build integration

package riak

import (
	"testing"
	"time"
)

func TestGitHubIssue17UpdateMulipleCountersInMapAtOnce(t *testing.T) {
	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	const bucketName = "github-17"
	var err error
	var cmd Command

	mapOp := &MapOperation{}
	mapOp.IncrementCounter("c1", 1)
	mapOp.IncrementCounter("c2", 2)
	mapOp.IncrementCounter("c3", 3)
	cmd, err = NewUpdateMapCommandBuilder().
		WithBucketType(testMapBucketType).
		WithBucket(bucketName).
		WithMapOperation(mapOp).
		WithReturnBody(true).
		WithTimeout(time.Second * 20).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	key := "unknown"
	if uc, ok := cmd.(*UpdateMapCommand); ok {
		if uc.Response == nil {
			t.Fatal("expected non-nil Response")
		}
		rsp := uc.Response
		if rsp.GeneratedKey == "" {
			t.Errorf("expected non-empty generated key")
		} else {
			key = rsp.GeneratedKey
		}
	} else {
		t.FailNow()
	}

	cmd, err = NewFetchMapCommandBuilder().
		WithBucketType(testMapBucketType).
		WithBucket(bucketName).
		WithKey(key).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if fc, ok := cmd.(*FetchMapCommand); ok {
		if fc.Response == nil {
			t.Fatal("expected non-nil Response")
		}
		rsp := fc.Response
		if actual, expected := rsp.Map.Counters["c1"], int64(1); actual != expected {
			t.Errorf("actual %v, expected %v", actual, expected)
		}
		if actual, expected := rsp.Map.Counters["c2"], int64(2); actual != expected {
			t.Errorf("actual %v, expected %v", actual, expected)
		}
		if actual, expected := rsp.Map.Counters["c3"], int64(3); actual != expected {
			t.Errorf("actual %v, expected %v", actual, expected)
		}
	} else {
		t.FailNow()
	}
}
