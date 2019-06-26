package fakesrvrdata

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
	"sync"
	"unsafe"
)

type MinMaxUint64 struct {
	Min uint64
	Max uint64
}

// Ths provides threadsafe access to a ThsT pointer. Note the object itself is not safe for multiple access, and must not be mutated, either by the original owner after calling Set, or by future users who call Get. If you need to mutate, perform a deep copy.
type Ths struct {
	v *ThsT
	m *sync.RWMutex
	// IncrementChan may be used to set the increments for a particular remap.
	// Note this is not synchronized with GetIncrementChan, so multiple writers calling GetIncrementChan and IncrmeentChan to get and set will race, unless they are externally synchronized.
	IncrementChan chan IncrementChanT
	// GetIncrementsChan may be used to get the current increments for all remaps.
	// The returned map must not be modified.
	// Note this is not synchronized with GetIncrementChan, so multiple writers calling GetIncrementChan and IncrmeentChan to get and set will race, unless they are externally synchronized.
	GetIncrementsChan chan map[string]BytesPerSec

	// DelayMS is the minimum and maximum delay to serve requests, in milliseconds.
	// Atomic - MUST be accessed with sync/atomic.LoadUintptr and sync/atomic.StoreUintptr.
	DelayMS *unsafe.Pointer
}

func NewThs() Ths {
	v := ThsT(nil)
	delayMSPtr := &MinMaxUint64{}
	delayMSUnsafePtr := unsafe.Pointer(delayMSPtr)
	return Ths{
		m:                 &sync.RWMutex{},
		v:                 &v,
		IncrementChan:     make(chan IncrementChanT, 10), // arbitrarily allow 10 writes before blocking. TODO document? config?
		GetIncrementsChan: make(chan map[string]BytesPerSec),
		DelayMS:           &delayMSUnsafePtr,
	}
}

func (t Ths) Set(v ThsT) {
	t.m.Lock()
	defer t.m.Unlock()
	*t.v = v
}

func (t Ths) Get() ThsT {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.v
}
