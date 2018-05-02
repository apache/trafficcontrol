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
)

// Update and Reset bucket properties

func TestSetAndResetBucketProperties(t *testing.T) {
	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	const bucket = "set-reset-bucket-props"
	orig_nval := uint32(3)
	new_nval := uint32(9)

	var err error
	var cmd Command

	cmd, err = NewStoreBucketPropsCommandBuilder().
		WithBucket(bucket).
		WithNVal(new_nval).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if sc, ok := cmd.(*StoreBucketPropsCommand); ok {
		if got, want := sc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}

	cmd, err = NewFetchBucketPropsCommandBuilder().
		WithBucket(bucket).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if fc, ok := cmd.(*FetchBucketPropsCommand); ok {
		if got, want := fc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := fc.Response.NVal, new_nval; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}

	cmd, err = NewResetBucketCommandBuilder().
		WithBucket(bucket).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if rc, ok := cmd.(*ResetBucketCommand); ok {
		if got, want := rc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}

	cmd, err = NewFetchBucketPropsCommandBuilder().
		WithBucket(bucket).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if fc, ok := cmd.(*FetchBucketPropsCommand); ok {
		if got, want := fc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := fc.Response.NVal, orig_nval; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}
}

// Update and Reset bucket type properties

func TestSetAndResetBucketTypeProperties(t *testing.T) {
	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	const bucketType = "plain"
	orig_nval := uint32(3)
	new_nval := uint32(4)

	var err error
	var cmd Command

	cmd, err = NewStoreBucketTypePropsCommandBuilder().
		WithBucketType(bucketType).
		WithNVal(new_nval).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if sc, ok := cmd.(*StoreBucketTypePropsCommand); ok {
		if got, want := sc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}

	cmd, err = NewFetchBucketTypePropsCommandBuilder().
		WithBucketType(bucketType).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if fc, ok := cmd.(*FetchBucketTypePropsCommand); ok {
		if got, want := fc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := fc.Response.NVal, new_nval; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}

	cmd, err = NewStoreBucketTypePropsCommandBuilder().
		WithBucketType(bucketType).
		WithNVal(orig_nval).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if rc, ok := cmd.(*StoreBucketTypePropsCommand); ok {
		if got, want := rc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}

	cmd, err = NewFetchBucketTypePropsCommandBuilder().
		WithBucketType(bucketType).
		Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if fc, ok := cmd.(*FetchBucketTypePropsCommand); ok {
		if got, want := fc.Success(), true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := fc.Response.NVal, orig_nval; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.FailNow()
	}
}
