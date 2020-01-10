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
	"github.com/apache/trafficcontrol/lib/go-tc/enum"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

// CRStatesThreadsafe provides safe access for multiple goroutines to read a single Crstates object, with a single goroutine writer.
// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and DeliveryService maps?
type CRStatesThreadsafe struct {
	crStates *tc.CRStates
	m        *sync.RWMutex
}

// NewCRStatesThreadsafe creates a new CRStatesThreadsafe object safe for multiple goroutine readers and a single writer.
func NewCRStatesThreadsafe() CRStatesThreadsafe {
	crs := tc.NewCRStates()
	return CRStatesThreadsafe{m: &sync.RWMutex{}, crStates: &crs}
}

// Get returns the internal Crstates object for reading.
func (t *CRStatesThreadsafe) Get() tc.CRStates {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.Copy()
}

// GetDeliveryServices returns the internal Crstates delivery services map for reading.
func (t *CRStatesThreadsafe) GetDeliveryServices() map[enum.DeliveryServiceName]tc.CRStatesDeliveryService {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyDeliveryServices()
}

// GetCache returns the availability data of the given cache. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetCache(name enum.CacheName) (available tc.IsAvailable, ok bool) {
	t.m.RLock()
	available, ok = t.crStates.Caches[name]
	t.m.RUnlock()
	return
}

// GetCaches returns the availability data of all caches. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetCaches() map[enum.CacheName]tc.IsAvailable {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyCaches()
}

// GetDeliveryService returns the availability data of the given delivery service. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetDeliveryService(name enum.DeliveryServiceName) (ds tc.CRStatesDeliveryService, ok bool) {
	t.m.RLock()
	ds, ok = t.crStates.DeliveryService[name]
	t.m.RUnlock()
	return
}

// SetCache sets the internal availability data for a particular cache. It does NOT set data if the cache doesn't already exist. By adding newly received caches with `AddCache`, this allows easily avoiding a race condition when an in-flight poller tries to set a cache which has been removed.
func (t *CRStatesThreadsafe) SetCache(cacheName enum.CacheName, available tc.IsAvailable) {
	t.m.Lock()
	if _, ok := t.crStates.Caches[cacheName]; ok {
		t.crStates.Caches[cacheName] = available
	}
	t.m.Unlock()
}

// AddCache adds the internal availability data for a particular cache.
func (t *CRStatesThreadsafe) AddCache(cacheName enum.CacheName, available tc.IsAvailable) {
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
func (t *CRStatesThreadsafe) SetDeliveryService(name enum.DeliveryServiceName, ds tc.CRStatesDeliveryService) {
	t.m.Lock()
	t.crStates.DeliveryService[name] = ds
	t.m.Unlock()
}

// DeleteDeliveryService deletes the given delivery service from the internal data. This MUST NOT be called by multiple goroutines.
func (t *CRStatesThreadsafe) DeleteDeliveryService(name enum.DeliveryServiceName) {
	t.m.Lock()
	delete(t.crStates.DeliveryService, name)
	t.m.Unlock()
}

// CRStatesPeersThreadsafe provides safe access for multiple goroutines to read a map of Traffic Monitor peers to their returned Crstates, with a single goroutine writer.
// This could be made lock-free, if the performance was necessary
type CRStatesPeersThreadsafe struct {
	crStates   map[enum.TrafficMonitorName]tc.CRStates
	peerStates map[enum.TrafficMonitorName]bool
	peerTimes  map[enum.TrafficMonitorName]time.Time
	peerOnline map[enum.TrafficMonitorName]bool
	timeout    *time.Duration
	m          *sync.RWMutex
}

// NewCRStatesPeersThreadsafe creates a new CRStatesPeers object safe for multiple goroutine readers and a single writer.
func NewCRStatesPeersThreadsafe() CRStatesPeersThreadsafe {
	timeout := time.Hour // default to a large timeout
	return CRStatesPeersThreadsafe{
		m:          &sync.RWMutex{},
		timeout:    &timeout,
		peerOnline: map[enum.TrafficMonitorName]bool{},
		crStates:   map[enum.TrafficMonitorName]tc.CRStates{},
		peerStates: map[enum.TrafficMonitorName]bool{},
		peerTimes:  map[enum.TrafficMonitorName]time.Time{},
	}
}

func (t *CRStatesPeersThreadsafe) SetTimeout(timeout time.Duration) {
	t.m.Lock()
	defer t.m.Unlock()
	*t.timeout = timeout
}

func (t *CRStatesPeersThreadsafe) SetPeers(newPeers map[enum.TrafficMonitorName]struct{}) {
	t.m.Lock()
	defer t.m.Unlock()
	for peer, _ := range t.crStates {
		_, ok := newPeers[peer]
		t.peerOnline[peer] = ok
	}
}

// GetCrstates returns the internal Traffic Monitor peer Crstates data. This MUST NOT be modified.
func (t *CRStatesPeersThreadsafe) GetCrstates() map[enum.TrafficMonitorName]tc.CRStates {
	t.m.RLock()
	m := map[enum.TrafficMonitorName]tc.CRStates{}
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

func copyPeerAvailable(a map[enum.TrafficMonitorName]bool) map[enum.TrafficMonitorName]bool {
	m := make(map[enum.TrafficMonitorName]bool, len(a))
	for k, v := range a {
		m[k] = v
	}
	return m
}

// GetPeerAvailability returns the state of the given peer
func (t *CRStatesPeersThreadsafe) GetPeerAvailability(peer enum.TrafficMonitorName) bool {
	t.m.RLock()
	availability := t.peerStates[peer] && t.peerOnline[peer] && time.Since(t.peerTimes[peer]) < *t.timeout
	t.m.RUnlock()
	return availability
}

// GetPeersOnline return a map of peers which are marked ONLINE in the latest CRConfig from Traffic Ops. This is NOT guaranteed to actually _contain_ all OFFLINE monitors returned by other functions, such as `GetPeerAvailability` and `GetQueryTimes`, but bool defaults to false, so the value of any key is guaranteed to be correct.
func (t *CRStatesPeersThreadsafe) GetPeersOnline() map[enum.TrafficMonitorName]bool {
	t.m.RLock()
	defer t.m.RUnlock()
	return copyPeerAvailable(t.peerOnline)
}

// GetQueryTimes returns the last query time of all peers
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
