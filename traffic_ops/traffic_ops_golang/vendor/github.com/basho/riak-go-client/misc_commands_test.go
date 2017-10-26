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
	"reflect"
	"testing"

	rpbRiak "github.com/basho/riak-go-client/rpb/riak"
)

func buildRpbGetBucketResp() *rpbRiak.RpbGetBucketResp {
	trueVal := true
	uint32val := uint32(9)
	replMode := rpbRiak.RpbBucketProps_REALTIME

	rpbModFun := &rpbRiak.RpbModFun{
		Module:   []byte("module_name"),
		Function: []byte("function_name"),
	}

	rpbCommitHook := &rpbRiak.RpbCommitHook{
		Name:   []byte("hook_name"),
		Modfun: rpbModFun,
	}

	rpbBucketProps := &rpbRiak.RpbBucketProps{
		NVal:          &uint32val,
		AllowMult:     &trueVal,
		LastWriteWins: &trueVal,
		HasPrecommit:  &trueVal,
		HasPostcommit: &trueVal,
		OldVclock:     &uint32val,
		YoungVclock:   &uint32val,
		BigVclock:     &uint32val,
		SmallVclock:   &uint32val,
		R:             &uint32val,
		Pr:            &uint32val,
		W:             &uint32val,
		Pw:            &uint32val,
		Dw:            &uint32val,
		Rw:            &uint32val,
		BasicQuorum:   &trueVal,
		NotfoundOk:    &trueVal,
		Search:        &trueVal,
		Consistent:    &trueVal,
		Repl:          &replMode,
		Backend:       []byte("backend"),
		SearchIndex:   []byte("index"),
		Datatype:      []byte("datatype"),
		HllPrecision:  &uint32val,
	}

	rpbBucketProps.Precommit = []*rpbRiak.RpbCommitHook{rpbCommitHook}
	rpbBucketProps.Postcommit = []*rpbRiak.RpbCommitHook{rpbCommitHook}
	rpbBucketProps.ChashKeyfun = rpbModFun
	rpbBucketProps.Linkfun = rpbModFun

	rpbGetBucketResp := &rpbRiak.RpbGetBucketResp{
		Props: rpbBucketProps,
	}

	return rpbGetBucketResp
}

func validateRpbBucketPropsForSetCommand(t *testing.T, rpb *rpbRiak.RpbBucketProps) {
	uint32val := uint32(9)

	if got, want := rpb.GetNVal(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetAllowMult(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetLastWriteWins(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetOldVclock(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetYoungVclock(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetBigVclock(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetSmallVclock(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetR(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetPr(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetW(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetPw(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetDw(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetRw(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetBasicQuorum(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetNotfoundOk(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetSearch(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.GetBackend()), "backend"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.GetSearchIndex()), "index"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.ChashKeyfun.Module), "module_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.ChashKeyfun.Function), "function_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.Precommit[0].Name), "hook_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.Precommit[0].Modfun.Module), "module_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.Precommit[0].Modfun.Function), "function_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.Postcommit[0].Name), "hook_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.Postcommit[0].Modfun.Module), "module_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(rpb.Postcommit[0].Modfun.Function), "function_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := rpb.GetHllPrecision(), uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func validateFetchBucketPropsResponse(t *testing.T, r *FetchBucketPropsResponse) {
	uint32val := uint32(9)
	replMode := rpbRiak.RpbBucketProps_REALTIME

	if r == nil {
		t.Fatal("want non-nil response")
	}
	if got, want := r.NVal, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.AllowMult, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.LastWriteWins, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.HasPrecommit, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.HasPostcommit, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.OldVClock, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.YoungVClock, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.BigVClock, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.SmallVClock, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.R, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.Pr, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.W, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.Pw, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.Dw, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.Rw, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.BasicQuorum, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.NotFoundOk, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.Search, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.Consistent, true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := int32(r.Repl), int32(replMode); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.Backend, "backend"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.SearchIndex, "index"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.DataType, "datatype"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.PreCommit[0].Name, "hook_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.PreCommit[0].ModFun.Module, "module_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.PreCommit[0].ModFun.Function, "function_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.PostCommit[0].Name, "hook_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.PostCommit[0].ModFun.Module, "module_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.PostCommit[0].ModFun.Function, "function_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.ChashKeyFun.Module, "module_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.ChashKeyFun.Function, "function_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.LinkFun.Module, "module_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.LinkFun.Function, "function_name"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := r.HllPrecision, uint32val; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// FetchBucketTypeProps

func TestBuildRpbGetBucketTypeReqCorrectlyViaBuilder(t *testing.T) {
	bt := "bucket_type"
	builder := NewFetchBucketTypePropsCommandBuilder().WithBucketType(bt)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.Fatal("protobuf is nil")
	}
	if req, ok := protobuf.(*rpbRiak.RpbGetBucketTypeReq); ok {
		if got, want := string(req.GetType()), bt; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiak.RpbGetBucketTypeReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbGetBucketRespCorrectlyForBucketType(t *testing.T) {
	builder := NewFetchBucketTypePropsCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	err = cmd.onSuccess(buildRpbGetBucketResp())
	if err != nil {
		t.Fatal(err.Error())
	}

	if got, want := cmd.Success(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if fetchBucketTypePropsCommand, ok := cmd.(*FetchBucketTypePropsCommand); ok {
		if fetchBucketTypePropsCommand.Response == nil {
			t.Fatal("unexpected nil object")
		}
		if got, want := fetchBucketTypePropsCommand.success, true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		validateFetchBucketPropsResponse(t, fetchBucketTypePropsCommand.Response)
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchBucketTypePropsCommand", ok, reflect.TypeOf(cmd))
	}
}

// FetchBucketProps

func TestBuildRpbGetBucketReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewFetchBucketPropsCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.Fatal("protobuf is nil")
	}
	if req, ok := protobuf.(*rpbRiak.RpbGetBucketReq); ok {
		if got, want := string(req.GetType()), "bucket_type"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := string(req.GetBucket()), "bucket_name"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiak.RpbGetBucketReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbGetBucketRespCorrectly(t *testing.T) {
	builder := NewFetchBucketPropsCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	err = cmd.onSuccess(buildRpbGetBucketResp())
	if err != nil {
		t.Fatal(err.Error())
	}

	if got, want := cmd.Success(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if fetchBucketPropsCommand, ok := cmd.(*FetchBucketPropsCommand); ok {
		if fetchBucketPropsCommand.Response == nil {
			t.Fatal("unexpected nil object")
		}
		if got, want := fetchBucketPropsCommand.success, true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		validateFetchBucketPropsResponse(t, fetchBucketPropsCommand.Response)
	} else {
		t.Errorf("ok: %v - could not convert %v to *FetchBucketPropsCommand", ok, reflect.TypeOf(cmd))
	}
}

// StoreBucketTypeProps

func TestBuildRpbSetBucketTypeReqCorrectlyViaBuilder(t *testing.T) {
	trueVal := true
	uint32val := uint32(9)

	modFun := &ModFun{
		Module:   "module_name",
		Function: "function_name",
	}
	hook := &CommitHook{
		Name:   "hook_name",
		ModFun: modFun,
	}

	builder := NewStoreBucketTypePropsCommandBuilder().
		WithBucketType("bucket_type").
		WithNVal(uint32val).
		WithAllowMult(trueVal).
		WithLastWriteWins(trueVal).
		WithOldVClock(uint32val).
		WithYoungVClock(uint32val).
		WithBigVClock(uint32val).
		WithSmallVClock(uint32val).
		WithR(uint32val).
		WithPr(uint32val).
		WithW(uint32val).
		WithPw(uint32val).
		WithDw(uint32val).
		WithRw(uint32val).
		WithBasicQuorum(trueVal).
		WithNotFoundOk(trueVal).
		WithSearch(trueVal).
		WithBackend("backend").
		WithSearchIndex("index").
		AddPreCommit(hook).
		AddPostCommit(hook).
		WithChashKeyFun(modFun).
		WithHllPrecision(uint32val)

	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.Fatal("protobuf is nil")
	}
	if req, ok := protobuf.(*rpbRiak.RpbSetBucketTypeReq); ok {
		if got, want := string(req.GetType()), "bucket_type"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		validateRpbBucketPropsForSetCommand(t, req.Props)
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiak.RpbSetBucketTypeReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbStoreBucketTypeRespCorrectly(t *testing.T) {
	builder := NewStoreBucketTypePropsCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	err = cmd.onSuccess(nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	if got, want := cmd.Success(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if storeBucketTypePropsCommand, ok := cmd.(*StoreBucketTypePropsCommand); ok {
		if got, want := storeBucketTypePropsCommand.success, true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *StoreBucketTypePropsCommand", ok, reflect.TypeOf(cmd))
	}
}

// StoreBucketProps

func TestBuildRpbSetBucketReqCorrectlyViaBuilder(t *testing.T) {
	trueVal := true
	uint32val := uint32(9)

	modFun := &ModFun{
		Module:   "module_name",
		Function: "function_name",
	}
	hook := &CommitHook{
		Name:   "hook_name",
		ModFun: modFun,
	}

	builder := NewStoreBucketPropsCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		WithNVal(uint32val).
		WithAllowMult(trueVal).
		WithLastWriteWins(trueVal).
		WithOldVClock(uint32val).
		WithYoungVClock(uint32val).
		WithBigVClock(uint32val).
		WithSmallVClock(uint32val).
		WithR(uint32val).
		WithPr(uint32val).
		WithW(uint32val).
		WithPw(uint32val).
		WithDw(uint32val).
		WithRw(uint32val).
		WithBasicQuorum(trueVal).
		WithNotFoundOk(trueVal).
		WithSearch(trueVal).
		WithBackend("backend").
		WithSearchIndex("index").
		AddPreCommit(hook).
		AddPostCommit(hook).
		WithChashKeyFun(modFun).
		WithHllPrecision(uint32val)

	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.Fatal("protobuf is nil")
	}
	if req, ok := protobuf.(*rpbRiak.RpbSetBucketReq); ok {
		if got, want := string(req.GetType()), "bucket_type"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := string(req.GetBucket()), "bucket_name"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		validateRpbBucketPropsForSetCommand(t, req.Props)
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiak.RpbSetBucketReq", ok, reflect.TypeOf(protobuf))
	}
}

func TestParseRpbStoreBucketRespCorrectly(t *testing.T) {
	builder := NewStoreBucketPropsCommandBuilder()
	cmd, err := builder.
		WithBucketType("bucket_type").
		WithBucket("bucket_name").
		Build()
	if err != nil {
		t.Error(err.Error())
	}

	err = cmd.onSuccess(nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	if got, want := cmd.Success(), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if storeBucketPropsCommand, ok := cmd.(*StoreBucketPropsCommand); ok {
		if got, want := storeBucketPropsCommand.success, true; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *StoreBucketPropsCommand", ok, reflect.TypeOf(cmd))
	}
}

// ResetBucket

func TestBuildRpbResetBucketReqCorrectlyViaBuilder(t *testing.T) {
	builder := NewResetBucketCommandBuilder().
		WithBucketType("bucket_type").
		WithBucket("bucket_name")
	cmd, err := builder.Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, ok := cmd.(retryableCommand); !ok {
		t.Errorf("got %v, want cmd %s to implement retryableCommand", ok, reflect.TypeOf(cmd))
	}

	protobuf, err := cmd.constructPbRequest()
	if err != nil {
		t.Fatal(err.Error())
	}
	if protobuf == nil {
		t.Fatal("protobuf is nil")
	}
	if req, ok := protobuf.(*rpbRiak.RpbResetBucketReq); ok {
		if got, want := string(req.GetType()), "bucket_type"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := string(req.GetBucket()), "bucket_name"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	} else {
		t.Errorf("ok: %v - could not convert %v to *rpbRiak.RpbResetBucketReq", ok, reflect.TypeOf(protobuf))
	}
}
