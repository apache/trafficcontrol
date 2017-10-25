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
	"fmt"
	"sync"
)

type state byte

type stateful interface {
	fmt.Stringer
	setStateDesc(desc ...string)
	isCurrentState(st state) (rv bool)
	isStateLessThan(st state) (rv bool)
	setState(st state)
	getState() (st state)
	stateCheck(allowed ...state) (err error)
}

type stateData struct {
	sync.RWMutex
	stateVal     state
	stateDesc    []string
	setStateFunc func(sd *stateData, st state)
}

var defaultSetStateFunc = func(sd *stateData, st state) {
	sd.stateVal = st
}

func (s *stateData) initStateData(desc ...string) {
	s.stateDesc = desc
	s.setStateFunc = defaultSetStateFunc
}

func (s *stateData) String() string {
	stateIdx := int(s.stateVal)
	if len(s.stateDesc) > stateIdx {
		return s.stateDesc[stateIdx]
	} else {
		return fmt.Sprintf("STATE_%v", stateIdx)
	}
}

func (s *stateData) isCurrentState(st state) bool {
	s.RLock()
	defer s.RUnlock()
	return s.stateVal == st
}

func (s *stateData) isStateLessThan(st state) bool {
	s.RLock()
	defer s.RUnlock()
	return s.stateVal < st
}

func (s *stateData) getState() state {
	s.RLock()
	defer s.RUnlock()
	return s.stateVal
}

func (s *stateData) setState(st state) {
	s.Lock()
	defer s.Unlock()
	s.setStateFunc(s, st)
}

func (s *stateData) stateCheck(allowed ...state) error {
	s.RLock()
	defer s.RUnlock()
	stateAllowed := false
	for _, st := range allowed {
		if s.stateVal == st {
			stateAllowed = true
			break
		}
	}
	if !stateAllowed {
		return fmt.Errorf("Illegal State - required %v: current: %v", allowed, s.stateVal)
	}
	return nil
}
