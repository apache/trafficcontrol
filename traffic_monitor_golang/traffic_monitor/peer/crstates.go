package peer

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
	"encoding/json"
	"sync"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
)

// Crstates includes availability data for caches and delivery services, as gathered and aggregated by this Traffic Monitor. It is designed to be served at an API endpoint primarily for Traffic Routers (Content Router) to consume.
// TODO rename to `CRStates`
type Crstates struct {
	Caches          map[enum.CacheName]IsAvailable               `json:"caches"`
	Deliveryservice map[enum.DeliveryServiceName]Deliveryservice `json:"deliveryServices"`
}

// NewCrstates creates a new CR states object, initializing pointer members.
func NewCrstates() Crstates {
	return Crstates{
		Caches:          map[enum.CacheName]IsAvailable{},
		Deliveryservice: map[enum.DeliveryServiceName]Deliveryservice{},
	}
}

// Copy creates a deep copy of this object. It does not mutate, and is thus safe for multiple goroutines.
func (a Crstates) Copy() Crstates {
	b := NewCrstates()
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	for k, v := range a.Deliveryservice {
		b.Deliveryservice[k] = v
	}
	return b
}

// CopyDeliveryservices creates a deep copy of the delivery service availability data.. It does not mutate, and is thus safe for multiple goroutines.
func (a Crstates) CopyDeliveryservices() map[enum.DeliveryServiceName]Deliveryservice {
	b := map[enum.DeliveryServiceName]Deliveryservice{}
	for k, v := range a.Deliveryservice {
		b[k] = v
	}
	return b
}

// CopyCaches creates a deep copy of the cache availability data.. It does not mutate, and is thus safe for multiple goroutines.
func (a Crstates) CopyCaches() map[enum.CacheName]IsAvailable {
	b := map[enum.CacheName]IsAvailable{}
	for k, v := range a.Caches {
		b[k] = v
	}
	return b
}

// IsAvailable contains whether the given cache or delivery service is available. It is designed for JSON serialization, namely in the Traffic Monitor 1.0 API.
type IsAvailable struct {
	IsAvailable bool `json:"isAvailable"`
}

// Deliveryservice contains data about the availability of a particular delivery service, and which caches in that delivery service have been marked as unavailable.
type Deliveryservice struct {
	DisabledLocations []enum.CacheGroupName `json:"disabledLocations"`
	IsAvailable       bool                  `json:"isAvailable"`
}

// CrstatesUnMarshall takes bytes of a JSON string, and unmarshals them into a Crstates object.
func CrstatesUnMarshall(body []byte) (Crstates, error) {
	var crStates Crstates
	err := json.Unmarshal(body, &crStates)
	return crStates, err
}

// CrstatesMarshall serializes the given Crstates into bytes.
func CrstatesMarshall(states Crstates) ([]byte, error) {
	return json.Marshal(states)
}

// CRStatesThreadsafe provides safe access for multiple goroutines to read a single Crstates object, with a single goroutine writer.
// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
type CRStatesThreadsafe struct {
	crStates *Crstates
	m        *sync.RWMutex
}

// NewCRStatesThreadsafe creates a new CRStatesThreadsafe object safe for multiple goroutine readers and a single writer.
func NewCRStatesThreadsafe() CRStatesThreadsafe {
	crs := NewCrstates()
	return CRStatesThreadsafe{m: &sync.RWMutex{}, crStates: &crs}
}

// Get returns the internal Crstates object for reading.
func (t *CRStatesThreadsafe) Get() Crstates {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.Copy()
}

// GetDeliveryServices returns the internal Crstates delivery services map for reading.
func (t *CRStatesThreadsafe) GetDeliveryServices() map[enum.DeliveryServiceName]Deliveryservice {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyDeliveryservices()
}

// GetCache returns the availability data of the given cache. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetCache(name enum.CacheName) (available IsAvailable, ok bool) {
	t.m.RLock()
	available, ok = t.crStates.Caches[name]
	t.m.RUnlock()
	return
}

// GetCaches returns the availability data of all caches. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetCaches() map[enum.CacheName]IsAvailable {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyCaches()
}

// GetDeliveryService returns the availability data of the given delivery service. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetDeliveryService(name enum.DeliveryServiceName) (ds Deliveryservice, ok bool) {
	t.m.RLock()
	ds, ok = t.crStates.Deliveryservice[name]
	t.m.RUnlock()
	return
}

// SetCache sets the internal availability data for a particular cache. It does NOT set data if the cache doesn't already exist. By adding newly received caches with `AddCache`, this allows easily avoiding a race condition when an in-flight poller tries to set a cache which has been removed.
func (t *CRStatesThreadsafe) SetCache(cacheName enum.CacheName, available IsAvailable) {
	t.m.Lock()
	if _, ok := t.crStates.Caches[cacheName]; ok {
		t.crStates.Caches[cacheName] = available
	}
	t.m.Unlock()
}

// AddCache adds the internal availability data for a particular cache.
func (t *CRStatesThreadsafe) AddCache(cacheName enum.CacheName, available IsAvailable) {
	t.m.Lock()
	t.crStates.Caches[cacheName] = available
	t.m.Unlock()
}

// DeleteCache deletes the given cache from the internal data.
func (t *CRStatesThreadsafe) DeleteCache(name enum.CacheName) {
	t.m.Lock()
	delete(t.crStates.Caches, name)
	t.m.Unlock()
}

// SetDeliveryService sets the availability data for the given delivery service.
func (t *CRStatesThreadsafe) SetDeliveryService(name enum.DeliveryServiceName, ds Deliveryservice) {
	t.m.Lock()
	t.crStates.Deliveryservice[name] = ds
	t.m.Unlock()
}

// DeleteDeliveryService deletes the given delivery service from the internal data. This MUST NOT be called by multiple goroutines.
func (t *CRStatesThreadsafe) DeleteDeliveryService(name enum.DeliveryServiceName) {
	t.m.Lock()
	delete(t.crStates.Deliveryservice, name)
	t.m.Unlock()
}

// CRStatesPeersThreadsafe provides safe access for multiple goroutines to read a map of Traffic Monitor peers to their returned Crstates, with a single goroutine writer.
// This could be made lock-free, if the performance was necessary
type CRStatesPeersThreadsafe struct {
	crStates   map[enum.TrafficMonitorName]Crstates
	peerStates map[enum.TrafficMonitorName]bool
	peerTimes  map[enum.TrafficMonitorName]time.Time
	m          *sync.RWMutex
}

// NewCRStatesPeersThreadsafe creates a new CRStatesPeers object safe for multiple goroutine readers and a single writer.
func NewCRStatesPeersThreadsafe() CRStatesPeersThreadsafe {
	return CRStatesPeersThreadsafe{
		m:          &sync.RWMutex{},
		crStates:   map[enum.TrafficMonitorName]Crstates{},
		peerStates: map[enum.TrafficMonitorName]bool{},
		peerTimes:  map[enum.TrafficMonitorName]time.Time{},
	}
}

// GetCrstates returns the internal Traffic Monitor peer Crstates data. This MUST NOT be modified.
func (t *CRStatesPeersThreadsafe) GetCrstates() map[enum.TrafficMonitorName]Crstates {
	t.m.RLock()
	m := map[enum.TrafficMonitorName]Crstates{}
	for k, v := range t.crStates {
		m[k] = v.Copy()
	}
	t.m.RUnlock()
	return m
}

func copyPeerTimes(a map[enum.TrafficMonitorName]time.Time) map[enum.TrafficMonitorName]time.Time {
	m := make(map[enum.TrafficMonitorName]time.Time, len(a))
	for k, v := range a {
		m[k] = v
	}
	return m
}

// GetPeerAvailability returns the state of the given peer
func (t *CRStatesPeersThreadsafe) GetPeerAvailability(peer enum.TrafficMonitorName) bool {
	t.m.RLock()
	availability := t.peerStates[peer]
	t.m.RUnlock()
	return availability
}

// GetPeerAvailability returns the state of the given peer
func (t *CRStatesPeersThreadsafe) GetQueryTimes() map[enum.TrafficMonitorName]time.Time {
	t.m.RLock()
	defer t.m.RUnlock()
	return copyPeerTimes(t.peerTimes)
}

// HasAvailablePeers returns true if at least one peer is online
func (t *CRStatesPeersThreadsafe) HasAvailablePeers() bool {
	availablePeers := false

	t.m.RLock()

	for _, available := range t.peerStates {
		if available {
			availablePeers = true
			break
		}
	}

	t.m.RUnlock()

	return availablePeers
}

// Set sets the internal Traffic Monitor peer state and Crstates data. This MUST NOT be called by multiple goroutines.
func (t *CRStatesPeersThreadsafe) Set(result Result) {
	t.m.Lock()
	t.crStates[result.ID] = result.PeerStates
	t.peerStates[result.ID] = result.Available
	t.peerTimes[result.ID] = result.Time
	t.m.Unlock()
}
