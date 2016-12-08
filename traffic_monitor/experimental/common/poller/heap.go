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
	"sync"
	"time"
)

type HeapPollInfo struct {
	Info HTTPPollInfo
	Next time.Time
}

// Heap implements a Heap from Introduction to Algorithms (Cormen et al). A Heap allows fase access of the maximum object, in this case the latest Next time, and O(log(n)) insert. This Heap is specifically designed to be used as a Priority Queue.
type Heap struct {
	m        sync.Mutex
	info     []HeapPollInfo
	PollerID int64
}

func left(i int) int {
	return 2*i + 1
}

func right(i int) int {
	return 2*i + 2
}

// TODO benchmark directly replacing this, to see if Go inlines the function call
func parent(i int) int {
	return (i - 1) / 2
}

func (h *Heap) heapify(i int) {
	l := left(i)
	r := right(i)
	var largest int
	if l < len(h.info) && h.info[i].Next.After(h.info[l].Next) {
		largest = l
	} else {
		largest = i
	}

	if r < len(h.info) && h.info[largest].Next.After(h.info[r].Next) {
		largest = r
	}

	if largest != i {
		h.info[i], h.info[largest] = h.info[largest], h.info[i]
		h.heapify(largest)
	}
}

func (h *Heap) increaseKey(i int, key HeapPollInfo) {
	if h.info[i].Next.After(key.Next) {
		panic("Poll.Heap.increaseKey got key smaller than index")
	}

	h.info[i] = key

	for i > 0 && h.info[parent(i)].Next.After(h.info[i].Next) {
		h.info[i], h.info[parent(i)] = h.info[parent(i)], h.info[i]
		i = parent(i)
	}
}

// Pop gets the latest time from the heap. Implements Algorithms HEAP-EXTRACT-MAX.
// Returns the info with the latest time, and false if the heap is empty.
func (h *Heap) Pop() (HeapPollInfo, bool) {
	h.m.Lock()
	defer h.m.Unlock()
	if len(h.info) == 0 {
		return HeapPollInfo{}, false
	}
	max := h.info[0]
	h.info[0] = h.info[len(h.info)-1]
	h.info = h.info[:len(h.info)-1]
	h.heapify(0)
	if max.Info.ID == "odol-atsec-jac-04" {
		fmt.Printf("httpPoll %v Heap.Pop id %v next %v\n", h.PollerID, max.Info.ID, max.Next)
	}
	return max, true
}

// Pop gets the latest time from the heap. Implements Algorithms MAX-HEAP-INSERT.
func (h *Heap) Push(key HeapPollInfo) {
	h.m.Lock()
	defer h.m.Unlock()
	if key.Info.ID == "odol-atsec-jac-04" {
		fmt.Printf("httpPoll %v Heap.Push id %v next %v\n", h.PollerID, key.Info.ID, key.Next)
	}
	h.info = append(h.info, HeapPollInfo{Next: time.Unix(1<<63-1, 0)})
	h.increaseKey(len(h.info)-1, key)
}
