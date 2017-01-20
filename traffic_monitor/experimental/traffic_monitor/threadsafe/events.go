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
	"sync"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/health"
)

// Events provides safe access for multiple goroutines readers and a single writer to a stored Events slice.
type Events struct {
	events    *[]health.Event
	m         *sync.RWMutex
	nextIndex *uint64
	max       uint64
}

func copyEvents(a []health.Event) []health.Event {
	b := make([]health.Event, len(a), len(a))
	copy(b, a)
	return b
}

// NewEvents creates a new single-writer-multiple-reader Threadsafe object
func NewEvents(maxEvents uint64) Events {
	i := uint64(0)
	return Events{m: &sync.RWMutex{}, events: &[]health.Event{}, nextIndex: &i, max: maxEvents}
}

// Get returns the internal slice of Events for reading. This MUST NOT be modified. If modification is necessary, copy the slice.
func (o *Events) Get() []health.Event {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.events
}

// Add adds the given event. This is threadsafe for one writer, multiple readers. This MUST NOT be called by multiple threads, as it non-atomically fetches and adds.
func (o *Events) Add(e health.Event) {
	// host="hostname", type=EDGE, available=true, msg="REPORTED - available"
	log.Eventf(e.Time, "host=\"%s\", type=%s, available=%t, msg=\"%s\"", e.Hostname, e.Type, e.Available, e.Description)
	o.m.Lock() // TODO test removing
	events := copyEvents(*o.events)
	e.Index = *o.nextIndex
	events = append([]health.Event{e}, events...)
	if len(events) > int(o.max) {
		events = (events)[:o.max-1]
	}
	// o.m.Lock()
	*o.events = events
	*o.nextIndex++
	o.m.Unlock()
}
