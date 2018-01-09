package crconfigregex

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
 *
 */

import (
	"sync"
)

// Ths provides threadsafe access to a ThsT pointer. Note the object itself is not safe for multiple access, and must not be mutated, either by the original owner after calling Set, or by future users who call Get. If you need to mutate, perform a deep copy.
type Ths struct {
	v *ThsT
	m *sync.RWMutex
}

func NewThs() Ths {
	v := ThsT(nil)
	return Ths{m: &sync.RWMutex{}, v: &v}
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
