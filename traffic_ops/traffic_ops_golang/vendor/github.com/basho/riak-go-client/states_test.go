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
	"testing"
)

type testStateData struct {
	stateData
}

const (
	STATE_ONE state = iota
	STATE_TWO
	STATE_THREE
	STATE_FOUR

	OTHER_STATE_ONE
	OTHER_STATE_TWO
	OTHER_STATE_THREE
)

func TestStateConsts(t *testing.T) {
	data1 := &testStateData{}
	data1.initStateData("STATE_ONE")
	data1.setState(STATE_ONE)

	data2 := &testStateData{}
	data2.initStateData("OTHER_STATE_ONE")
	data2.setState(OTHER_STATE_ONE)

	if s1, s2 := data1.getState(), data2.getState(); s1 == s2 {
		t.Errorf("whoops, %v equals %v", s1, s2)
	}
}

func TestStateData(t *testing.T) {
	data := &testStateData{}
	data.initStateData("STATE_TWO")
	data.setState(STATE_TWO)

	if expected, actual := true, data.isCurrentState(STATE_TWO); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	if expected, actual := false, data.isCurrentState(STATE_ONE); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestAllowedState(t *testing.T) {
	data := &testStateData{}
	data.initStateData("STATE_TWO")
	data.setState(STATE_TWO)

	if err := data.stateCheck(STATE_ONE, STATE_THREE); err == nil {
		t.Errorf("expected non-nil error, got %v", err)
	}
}

func TestStateDesc(t *testing.T) {
	data := &testStateData{}
	data.initStateData("STATE_ONE", "STATE_TWO", "STATE_THREE")

	data.setState(STATE_ONE)
	if expected, actual := "STATE_ONE", data.String(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	data.setState(STATE_TWO)
	if expected, actual := "STATE_TWO", data.String(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	data.setState(STATE_THREE)
	if expected, actual := "STATE_THREE", data.String(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestStateDescUnknown(t *testing.T) {
	data := &testStateData{}
	data.initStateData("STATE_ONE", "STATE_TWO", "STATE_THREE")
	data.setState(STATE_FOUR)

	if expected, actual := "STATE_3", data.String(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
