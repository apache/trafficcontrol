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

func TestEnqueueDequeueCommandsConcurrently(t *testing.T) {
	queueSize := uint16(64)
	queue := newQueue(queueSize)

	w := &sync.WaitGroup{}
	for i := uint16(0); i < queueSize; i++ {
		w.Add(1)
		go func() {
			cmd := &PingCommand{}
			async := &Async{
				Command: cmd,
			}
			if err := queue.enqueue(async); err != nil {
				t.Error(err)
			}
			w.Done()
		}()
	}

	w.Wait()

	cmd := &PingCommand{}
	async := &Async{
		Command: cmd,
	}
	if err := queue.enqueue(async); err == nil {
		t.Error("expected non-nil err when enqueueing one more command than max")
	}
	if expected, actual := false, queue.isEmpty(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	w = &sync.WaitGroup{}
	for i := uint16(0); i < queueSize; i++ {
		w.Add(1)
		go func() {
			cmd, err := queue.dequeue()
			if cmd == nil {
				t.Error("expected non-nil cmd")
			}
			if err != nil {
				t.Error("expected nil err")
			}
			w.Done()
		}()
	}

	w.Wait()

	if expected, actual := true, queue.isEmpty(); expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	queue.destroy()

	_, err := queue.dequeue()
	if err == nil {
		t.Error("expected non-nil err")
	}
}
