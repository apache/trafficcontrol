package cgsrch

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
	"errors"
	"sync"

	"github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/quadtree"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// CGSearcher is the interface that wraps the Nearest method.
//
// Nearest searches its list of objects and returns the closest to the given point.
type CGSearcher interface {
	Nearest(top float64, left float64) (quadtree.DataT, bool)
}

func Create(crc *tc.CRConfig) (CGSearcher, error) {
	// TODO change to Searcher interface?
	if crc == nil {
		return nil, errors.New("CRConfig is nil")
	}
	qt := quadtree.New()
	for cg, ll := range crc.EdgeLocations {
		qt.Insert(quadtree.DataT{Lat: ll.Lat, Lon: ll.Lon, Obj: quadtree.ObjT(cg)})
	}
	return qt, nil
}

// ThsT is the Threadsafe type used by this package. ThsT should usually be a pointer or an interface which holds a pointer.
type ThsT CGSearcher

// Ths provides threadsafe access to a ThsT
type Ths struct {
	v *ThsT
	m *sync.RWMutex
}

// NewThs creates a new Threadsafe Ths container.
func NewThs() Ths {
	v := ThsT(nil)
	return Ths{m: &sync.RWMutex{}, v: &v}
}

// Set sets the given object in the threadsafe container. The given object MUST NOT be modified after calling this.
func (t Ths) Set(v ThsT) {
	t.m.Lock()
	defer t.m.Unlock()
	*t.v = v
}

// Get returns the object held by the threadsafe container. The object MUST NOT be modified.
func (t Ths) Get() ThsT {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.v
}
