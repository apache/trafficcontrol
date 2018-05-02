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
	"sync"
	"testing"
)

func TestReadFromEmptyQueue(t *testing.T) {
	q := newQueue(1)
	v, err := q.dequeue()
	if err != nil {
		t.Error("expected nil error when reading from empty queue")
	}
	if v != nil {
		t.Error("expected nil value when reading from empty queue")
	}
	if expected, actual := uint16(0), q.count(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestIterateEmptyQueue(t *testing.T) {
	count := uint16(128)
	q := newQueue(count)
	executed := false
	var f = func(val interface{}) (bool, bool) {
		executed = true
		if val == nil {
			return true, true
		} else {
			return false, true
		}
	}
	err := q.iterate(f)
	if err != nil {
		t.Error("expected nil error when iterating queue")
	}
	if expected, actual := false, executed; expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	if expected, actual := uint16(0), q.count(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestConcurrentIterateQueue(t *testing.T) {
	count := uint16(2)
	wg := &sync.WaitGroup{}
	q := newQueue(count + 2) // make room for 666
	for i := uint16(0); i < count; i++ {
		q.enqueue(i)
	}

	wg_inner := &sync.WaitGroup{}
	for i := uint16(0); i < count; i++ {
		wg.Add(1)
		go func() {
			var f = func(val interface{}) (bool, bool) {
				wg_inner.Add(1)
				go func() {
					q.enqueue(666)
					wg_inner.Done()
				}()
				return false, true
			}
			err := q.iterate(f)
			if err != nil {
				t.Error("expected nil error when iterating queue")
			}
			wg.Done()
		}()
	}

	wg.Wait()
	wg_inner.Wait()

	if expected, actual := uint16(4), q.count(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
