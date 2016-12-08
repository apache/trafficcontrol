package poller

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	h := &Heap{}

	num := 100
	for i := 0; i < num; i++ {
		h.Push(HeapPollInfo{
			Info: HTTPPollInfo{
				Interval: time.Second * time.Duration(8),
				ID:       fmt.Sprintf("%v", i),
			},
			Next: time.Now().Add(time.Second * time.Duration(i)), // time.Duration((i%2)*-1)
		})
	}

	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if val.Info.ID != fmt.Sprintf("%v", i) {
			t.Errorf("expected pop ID %v got %v next %v", i, val.Info.ID, val.Next)
		}
	}
}

func TestHeapRandom(t *testing.T) {
	h := &Heap{}

	num := 10
	for i := 0; i < num; i++ {
		h.Push(HeapPollInfo{
			Info: HTTPPollInfo{
				Interval: time.Second * time.Duration(8),
				ID:       fmt.Sprintf("%v", i),
			},
			Next: time.Now().Add(time.Duration(rand.Int63())),
		})
	}

	previousTime := time.Now()
	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}
}

func TestHeapRandomPopping(t *testing.T) {
	h := &Heap{}

	randInfo := func(id int) HeapPollInfo {
		return HeapPollInfo{
			Info: HTTPPollInfo{
				Interval: time.Second * time.Duration(8),
				ID:       fmt.Sprintf("%v", id),
			},
			Next: time.Now().Add(time.Duration(rand.Int63())),
		}
	}

	num := 10
	for i := 0; i < num; i++ {
		h.Push(randInfo(i))
	}

	previousTime := time.Now()
	for i := 0; i < num/2; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		h.Push(randInfo(i))
	}
	val, ok := h.Pop()
	if !ok {
		t.Errorf("expected pop, got empty heap")
	} else {
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		h.Push(randInfo(i))
	}
	val, ok = h.Pop()
	if !ok {
		t.Errorf("expected pop, got empty heap")
	} else {
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	for i := 0; i < num/2-2; i++ { // -2 for the two we manually popped in order to get the max
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	val, ok = h.Pop()
	if ok {
		t.Errorf("expected empty, got %+v", val)
	}
}
