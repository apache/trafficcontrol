package threadsafe

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
	"sync/atomic"
)

// Uint provides safe access for multiple goroutines readers and a single writer to a stored uint.
type Uint struct {
	val *uint64
}

// NewUint returns a new single-writer-multiple-reader threadsafe uint
func NewUint() Uint {
	v := uint64(0)
	return Uint{val: &v}
}

// Get gets the internal uint. This is safe for multiple readers
func (u *Uint) Get() uint64 {
	return atomic.LoadUint64(u.val)
}

// Set sets the internal uint. This MUST NOT be called by multiple goroutines.
func (u *Uint) Set(v uint64) {
	atomic.StoreUint64(u.val, v)
}

// Inc increments the internal uint64.
// TODO make sure everything using this uses the value it returns, not a separate Get
func (u *Uint) Inc() uint64 {
	return atomic.AddUint64(u.val, 1)
}
