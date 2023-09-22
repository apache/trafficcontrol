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
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

const defaultMapCapacity = 8

// CRStatesThreadsafe provides safe access for multiple goroutines to read a single Crstates object, with a single goroutine writer.
// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and DeliveryService maps?
type CRStatesThreadsafe struct {
	crStates *tc.CRStates
	m        *sync.RWMutex
}

// NewCRStatesThreadsafe creates a new CRStatesThreadsafe object safe for multiple goroutine readers and a single writer.
func NewCRStatesThreadsafe() CRStatesThreadsafe {
	crs := tc.NewCRStates(defaultMapCapacity, defaultMapCapacity)
	return CRStatesThreadsafe{m: &sync.RWMutex{}, crStates: &crs}
}

// Get returns the internal Crstates object for reading.
func (t *CRStatesThreadsafe) Get() tc.CRStates {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.Copy()
}

// GetDeliveryServices returns the internal Crstates delivery services map for reading.
func (t *CRStatesThreadsafe) GetDeliveryServices() map[tc.DeliveryServiceName]tc.CRStatesDeliveryService {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyDeliveryServices()
}

// GetCache returns the availability data of the given cache. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetCache(name tc.CacheName) (available tc.IsAvailable, ok bool) {
	t.m.RLock()
	available, ok = t.crStates.Caches[name]
	t.m.RUnlock()
	return
}

// GetCaches returns the availability data of all caches. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetCaches() map[tc.CacheName]tc.IsAvailable {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.crStates.CopyCaches()
}

// GetDeliveryService returns the availability data of the given delivery service. This does not mutate, and is thus safe for multiple goroutines to call.
func (t *CRStatesThreadsafe) GetDeliveryService(name tc.DeliveryServiceName) (ds tc.CRStatesDeliveryService, ok bool) {
	t.m.RLock()
	ds, ok = t.crStates.DeliveryService[name]
	t.m.RUnlock()
	return
}

// SetCache sets the internal availability data for a particular cache. It does NOT set data if the cache doesn't already exist. By adding newly received caches with `AddCache`, this allows easily avoiding a race condition when an in-flight poller tries to set a cache which has been removed.
func (t *CRStatesThreadsafe) SetCache(cacheName tc.CacheName, available tc.IsAvailable) {
	t.m.Lock()
	if _, ok := t.crStates.Caches[cacheName]; ok {
		t.crStates.Caches[cacheName] = available
	}
	t.m.Unlock()
}

// AddCache adds the internal availability data for a particular cache.
func (t *CRStatesThreadsafe) AddCache(cacheName tc.CacheName, available tc.IsAvailable) {
	t.m.Lock()
	t.crStates.Caches[cacheName] = available
	t.m.Unlock()
}

// DeleteCache deletes the given cache from the internal data.
func (t *CRStatesThreadsafe) DeleteCache(name tc.CacheName) {
	t.m.Lock()
	delete(t.crStates.Caches, name)
	t.m.Unlock()
}

// SetDeliveryService sets the availability data for the given delivery service.
func (t *CRStatesThreadsafe) SetDeliveryService(name tc.DeliveryServiceName, ds tc.CRStatesDeliveryService) {
	t.m.Lock()
	t.crStates.DeliveryService[name] = ds
	t.m.Unlock()
}

// DeleteDeliveryService deletes the given delivery service from the internal data. This MUST NOT be called by multiple goroutines.
func (t *CRStatesThreadsafe) DeleteDeliveryService(name tc.DeliveryServiceName) {
	t.m.Lock()
	delete(t.crStates.DeliveryService, name)
	t.m.Unlock()
}

// CRStatesPeersThreadsafe provides safe access for multiple goroutines to read a map of Traffic Monitor peers to their returned Crstates, with a single goroutine writer.
// This could be made lock-free, if the performance was necessary
type CRStatesPeersThreadsafe struct {
	crStates   map[tc.TrafficMonitorName]tc.CRStates
	peerStates map[tc.TrafficMonitorName]bool
	peerTimes  map[tc.TrafficMonitorName]time.Time
	peerOnline map[tc.TrafficMonitorName]bool
	peerCount  *int
	quorumMin  *int
	timeout    *time.Duration
	m          *sync.RWMutex
}

// NewCRStatesPeersThreadsafe creates a new CRStatesPeers object safe for multiple goroutine readers and a single writer.
func NewCRStatesPeersThreadsafe(quorumMin int) CRStatesPeersThreadsafe {
	count := 0
	timeout := time.Hour // default to a large timeout
	return CRStatesPeersThreadsafe{
		m:          &sync.RWMutex{},
		timeout:    &timeout,
		peerOnline: map[tc.TrafficMonitorName]bool{},
		crStates:   map[tc.TrafficMonitorName]tc.CRStates{},
		peerStates: map[tc.TrafficMonitorName]bool{},
		peerTimes:  map[tc.TrafficMonitorName]time.Time{},
		peerCount:  &count,
		quorumMin:  &quorumMin,
	}
}

func (t *CRStatesPeersThreadsafe) SetTimeout(timeout time.Duration) {
	t.m.Lock()
	defer t.m.Unlock()
	*t.timeout = timeout
}

func (t *CRStatesPeersThreadsafe) SetPeers(newPeers map[tc.TrafficMonitorName]struct{}) {
	t.m.Lock()
	defer t.m.Unlock()

	peerCount := 0

	for peer, _ := range t.crStates {
		_, ok := newPeers[peer]
		t.peerOnline[peer] = ok

		if ok {
			peerCount++
		}
	}

	*t.peerCount = peerCount
}

// GetCrstates returns the internal Traffic Monitor peer Crstates data. This MUST NOT be modified.
func (t *CRStatesPeersThreadsafe) GetCrstates() map[tc.TrafficMonitorName]tc.CRStates {
	t.m.RLock()
	m := map[tc.TrafficMonitorName]tc.CRStates{}
	for k, v := range t.crStates {
		m[k] = v.Copy()
	}
	t.m.RUnlock()
	return m
}

func copyPeerTimes(a map[tc.TrafficMonitorName]time.Time) map[tc.TrafficMonitorName]time.Time {
	m := make(map[tc.TrafficMonitorName]time.Time, len(a))
	for k, v := range a {
		m[k] = v
	}
	return m
}

func copyPeerAvailable(a map[tc.TrafficMonitorName]bool) map[tc.TrafficMonitorName]bool {
	m := make(map[tc.TrafficMonitorName]bool, len(a))
	for k, v := range a {
		m[k] = v
	}
	return m
}

// GetPeerAvailability returns the state of the given peer
func (t *CRStatesPeersThreadsafe) GetPeerAvailability(peer tc.TrafficMonitorName) bool {
	t.m.RLock()
	availability := t.peerStates[peer] && t.peerOnline[peer] && time.Since(t.peerTimes[peer]) < *t.timeout
	t.m.RUnlock()
	return availability
}

type CRStatesPeersInfo struct {
	peerStates map[tc.TrafficMonitorName]bool
	peerOnline map[tc.TrafficMonitorName]bool
	peerTimes  map[tc.TrafficMonitorName]time.Time
	crStates   map[tc.TrafficMonitorName]tc.CRStates
	timeout    time.Duration
}

func (i *CRStatesPeersInfo) GetCrStates() map[tc.TrafficMonitorName]tc.CRStates {
	return i.crStates
}

func (i *CRStatesPeersInfo) GetPeerAvailability(peer tc.TrafficMonitorName) bool {
	return i.peerStates[peer] && i.peerOnline[peer] && time.Since(i.peerTimes[peer]) < i.timeout
}

func (i *CRStatesPeersInfo) HasAvailablePeers() bool {
	for _, available := range i.peerStates {
		if available {
			return true
		}
	}
	return false
}

// GetCRStatesPeersInfo returns a CRStatesPeersInfo which contains a copy of some
// internal fields of this CRStatesPeersThreadsafe. This all happens within
// a single, locked read transaction so that all fields are in sync with each other.
// This is to avoid making multiple read copies within the same goroutine for the same
// operation, which is very inefficient.
func (t *CRStatesPeersThreadsafe) GetCRStatesPeersInfo() CRStatesPeersInfo {
	t.m.RLock()
	defer t.m.RUnlock()
	info := CRStatesPeersInfo{
		peerStates: copyPeerAvailable(t.peerStates),
		peerOnline: copyPeerAvailable(t.peerOnline),
		peerTimes:  copyPeerTimes(t.peerTimes),
		crStates:   make(map[tc.TrafficMonitorName]tc.CRStates, len(t.crStates)),
		timeout:    *t.timeout,
	}
	for k, v := range t.crStates {
		info.crStates[k] = v.Copy()
	}
	return info
}

// GetPeersOnline return a map of peers which are marked ONLINE in the latest CRConfig from Traffic Ops. This is NOT guaranteed to actually _contain_ all OFFLINE monitors returned by other functions, such as `GetPeerAvailability` and `GetQueryTimes`, but bool defaults to false, so the value of any key is guaranteed to be correct.
func (t *CRStatesPeersThreadsafe) GetPeersOnline() map[tc.TrafficMonitorName]bool {
	t.m.RLock()
	defer t.m.RUnlock()
	return copyPeerAvailable(t.peerOnline)
}

// GetQueryTimes returns the last query time of all peers
func (t *CRStatesPeersThreadsafe) GetQueryTimes() map[tc.TrafficMonitorName]time.Time {
	t.m.RLock()
	defer t.m.RUnlock()
	return copyPeerTimes(t.peerTimes)
}

// numAvailablePeers is a private function to determine how many peers are currently available; callers must lock t
func (t *CRStatesPeersThreadsafe) numAvailablePeers() int {
	count := 0

	for _, available := range t.peerStates {
		if available {
			count++
		}
	}

	return count
}

// Set sets the internal Traffic Monitor peer state and Crstates data. This MUST NOT be called by multiple goroutines.
func (t *CRStatesPeersThreadsafe) Set(result Result) {
	t.m.Lock()
	t.crStates[result.ID] = result.PeerStates
	t.peerStates[result.ID] = result.Available
	t.peerTimes[result.ID] = result.Time
	t.m.Unlock()
}

// HasOptimisticQuorum returns true when the number of available peers is equal to or greater than the peer_optimistic_quorum_min setting.
func (t *CRStatesPeersThreadsafe) HasOptimisticQuorum() (bool, int, int, int) {
	t.m.RLock()
	defer t.m.RUnlock()

	available := t.numAvailablePeers()

	if available >= *t.quorumMin {
		return true, available, *t.peerCount, *t.quorumMin
	}

	return false, available, *t.peerCount, *t.quorumMin
}

// OptimisticQuorumEnabled returns true when peer_optimistic_quorum_min is set to a value greater than zero and the number of peers is greater than 1. Optimistic quorum requires a minimum of three Traffic Monitors; every individual monitor requires at least two peers to prevent a split-brain scenario that would be caused by having a single peer. If a single peer was legal (i.e.: two Traffic Monitors), neither peer would know which peer is reachable, and consequently both would serve 503s. This would force all Traffic Routers to use only their last-known state until the peering is restored, despite the fact that one of the two Traffic Monitors could still be reachable. A future enhancement could employ a heuristic to enable two monitors to determine whether they are offline independently by combining peer connectivity state with a calculation around the number of caches that are reachable, which might also include a rate of change in cache health state.
func (t *CRStatesPeersThreadsafe) OptimisticQuorumEnabled() bool {
	t.m.RLock()
	defer t.m.RUnlock()

	if *t.quorumMin > 0 && *t.peerCount > 1 {
		return true
	}

	return false
}
