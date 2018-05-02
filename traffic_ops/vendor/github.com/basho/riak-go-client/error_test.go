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
	"testing"

	rpb_riak "github.com/basho/riak-go-client/rpb/riak"
)

func TestBuildRiakErrorFromRpbErrorResp(t *testing.T) {
	var errcode uint32 = 1
	errmsg := bytes.NewBufferString("this is an error")
	rpbErr := &rpb_riak.RpbErrorResp{
		Errcode: &errcode,
		Errmsg:  errmsg.Bytes(),
	}
	err := newRiakError(rpbErr)
	if riakError, ok := err.(RiakError); ok == true {
		if expected, actual := errcode, riakError.Errcode; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "this is an error", riakError.Errmsg; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "RiakError|1|this is an error", riakError.Error(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.Error("error in type conversion")
	}
}
