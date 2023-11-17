package cache

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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// AvailableTuple contains a boolean value to indicate whether IPv4 is
// available and a boolean value to indicate whether IPv6 is available.
type AvailableTuple struct {
	IPv4 bool
	IPv6 bool
}

// SetAvailablility sets two booleans to indicate whether IPv4 is available
// and whether IPv6 is available.
func (a *AvailableTuple) SetAvailability(usingIPv4 bool, isAvailable bool) {
	if usingIPv4 {
		a.IPv4 = isAvailable
	} else {
		a.IPv6 = isAvailable
	}
}

// AvailableStatus is the available status of the given cache. It includes
// a boolean available/unavailable flag, and a descriptive string.
type AvailableStatus struct {
	// Available indicates whether a Cache Server is available for various IP
	// protocol versions.
	Available          AvailableTuple
	ProcessedAvailable bool
	LastCheckedIPv4    bool
	// The name of the actual status the cache server has, as configured in
	// Traffic Ops.
	Status string
	// `Why` will contain the reason a cache server has been purposely marked
	// unavailable by a Traffic Ops operator, if indeed that has occurred.
	Why string
	// UnavailableStat is the stat whose threshold made the cache unavailable.
	// If this is the empty string, the cache is unavailable for a
	// non-threshold reason. This exists so a poller (health, stat) won't mark
	// an unavailable cache as available if the stat whose threshold was
	// reached isn't available on that poller.
	UnavailableStat string
	// Poller is the name of the poller which set this availability status.
	Poller string
}

// CacheAvailableStatuses is the available status of each cache.
type AvailableStatuses map[string]AvailableStatus

// Copy copies this CacheAvailableStatuses. It does not modify, and thus is
// safe for multiple reader goroutines.
func (a AvailableStatuses) Copy() AvailableStatuses {
	b := AvailableStatuses(make(map[string]AvailableStatus, len(a)))
	for cacheName, status := range a {
		b[cacheName] = status
	}
	return b
}

// ResultHistory is a map of cache names, to an array of result history from
// each cache server.
type ResultHistory map[tc.CacheName][]Result

func copyResult(a []Result) []Result {
	b := make([]Result, len(a), len(a))
	copy(b, a)
	return b
}

// Copy copies returns a deep copy of this ResultHistory.
func (a ResultHistory) Copy() ResultHistory {
	b := make(ResultHistory, len(a))
	for k, v := range a {
		b[k] = copyResult(v)
	}
	return b
}

func pruneStats(history []tc.ResultStatVal, limit uint64) []tc.ResultStatVal {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}

// TODO determine if anything ever needs more than the latest, and if not, change
// ResultInfo to not be a slice.
type ResultInfoHistory map[tc.CacheName][]ResultInfo

// ResultInfo contains all the non-stat result info. This includes the cache ID,
// any errors, the time of the poll, the request time duration, Astats System
// (Vitals), Poll ID, and Availability.
// TODO: Determine why this exists, it doesn't seem to differ from Result in any
// meaningful way.
type ResultInfo struct {
	Available       bool
	Error           error
	ID              string
	PollID          uint64
	RequestTime     time.Duration
	Statistics      Statistics
	Time            time.Time
	UsingIPv4       bool
	Vitals          Vitals
	InterfaceVitals map[string]Vitals
}

func ToInfo(r Result) ResultInfo {
	return ResultInfo{
		Available:       r.Available,
		Error:           r.Error,
		ID:              r.ID,
		PollID:          r.PollID,
		RequestTime:     r.RequestTime,
		Statistics:      r.Statistics,
		Time:            r.Time,
		UsingIPv4:       r.UsingIPv4,
		Vitals:          r.Vitals,
		InterfaceVitals: r.InterfaceVitals,
	}
}

func copyResultInfos(a []ResultInfo) []ResultInfo {
	b := make([]ResultInfo, len(a), len(a))
	copy(b, a)
	return b
}

func (a ResultInfoHistory) Copy() ResultInfoHistory {
	b := make(ResultInfoHistory, len(a))
	for k, v := range a {
		b[k] = copyResultInfos(v) // TODO determine if copy is necessary
	}
	return b
}

func pruneInfos(history []ResultInfo, limit uint64) []ResultInfo {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}

func (a ResultInfoHistory) Add(r Result, limit uint64) {
	a[tc.CacheName(r.ID)] = pruneInfos(append([]ResultInfo{ToInfo(r)}, a[tc.CacheName(r.ID)]...), limit)
}

// Kbpses is the kbps values of each interface of a cache server.
type Kbpses map[string]uint64

// Copy returns a deep copy of the Kbpses.
func (a Kbpses) Copy() Kbpses {
	b := make(Kbpses, len(a))
	for k, v := range a {
		b[k] = v
	}
	return b
}
