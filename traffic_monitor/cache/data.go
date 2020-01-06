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
	"fmt"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"

	"github.com/json-iterator/go"
)

// AvailableStatusReported is the status string returned by caches set to "reported" in Traffic Ops.
// TODO put somewhere more generic
const AvailableStatusReported = "REPORTED"

// CacheAvailableStatus is the available status of the given cache. It includes a boolean available/unavailable flag, and a descriptive string.
type AvailableStatus struct {
	Available bool
	Status    string
	Why       string
	// UnavailableStat is the stat whose threshold made the cache unavailable. If this is the empty string, the cache is unavailable for a non-threshold reason. This exists so a poller (health, stat) won't mark an unavailable cache as available if the stat whose threshold was reached isn't available on that poller.
	UnavailableStat string
	// Poller is the name of the poller which set this available status
	Poller string
}

// CacheAvailableStatuses is the available status of each cache.
type AvailableStatuses map[tc.CacheName]AvailableStatus

// Copy copies this CacheAvailableStatuses. It does not modify, and thus is safe for multiple reader goroutines.
func (a AvailableStatuses) Copy() AvailableStatuses {
	b := AvailableStatuses(map[tc.CacheName]AvailableStatus{})
	for k, v := range a {
		b[k] = v
	}
	return b
}

// ResultHistory is a map of cache names, to an array of result history from each cache.
type ResultHistory map[tc.CacheName][]Result

func copyResult(a []Result) []Result {
	b := make([]Result, len(a), len(a))
	copy(b, a)
	return b
}

// Copy copies returns a deep copy of this ResultHistory
func (a ResultHistory) Copy() ResultHistory {
	b := ResultHistory{}
	for k, v := range a {
		b[k] = copyResult(v)
	}
	return b
}

// ResultStatVal is the value of an individual stat returned from a poll. Time is the time this stat was returned.
// Span is the number of polls this stat has been the same. For example, if History is set to 100, and the last 50 polls had the same value for this stat (but none of the previous 50 were the same), this stat's map value slice will actually contain 51 entries, and the first entry will have the value, the time of the last poll, and a Span of 50. Assuming the poll time is every 8 seconds, users will then know, looking at the Span, that the value was unchanged for the last 50*8=400 seconds.
// JSON values are all strings, for the TM1.0 /publish/CacheStats API.
type ResultStatVal struct {
	Val  interface{} `json:"value"`
	Time time.Time   `json:"time"`
	Span uint64      `json:"span"`
}

func (t *ResultStatVal) MarshalJSON() ([]byte, error) {
	v := struct {
		Val  string `json:"value"`
		Time int64  `json:"time"`
		Span uint64 `json:"span"`
	}{
		Val:  fmt.Sprintf("%v", t.Val),
		Time: t.Time.UnixNano() / 1000000, // ms since the epoch
		Span: t.Span,
	}
	json := jsoniter.ConfigFastest // TODO make configurable
	return json.Marshal(&v)
}

func (t *ResultStatVal) UnmarshalJSON(data []byte) error {
	v := struct {
		Val  string `json:"value"`
		Time int64  `json:"time"`
		Span uint64 `json:"span"`
	}{}
	json := jsoniter.ConfigFastest // TODO make configurable
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	t.Time = time.Unix(0, v.Time*1000000)
	t.Val = v.Val
	t.Span = v.Span
	return nil
}

func pruneStats(history []ResultStatVal, limit uint64) []ResultStatVal {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}

// TODO determine if anything ever needs more than the latest, and if not, change ResultInfo to not be a slice.
type ResultInfoHistory map[tc.CacheName][]ResultInfo

// ResultInfo contains all the non-stat result info. This includes the cache ID, any errors, the time of the poll, the request time duration, Astats System (Vitals), Poll ID, and Availability.
type ResultInfo struct {
	ID          tc.CacheName
	Error       error
	Time        time.Time
	RequestTime time.Duration
	Vitals      Vitals
	System      AstatsSystem
	PollID      uint64
	Available   bool
}

func ToInfo(r Result) ResultInfo {
	return ResultInfo{
		ID:          r.ID,
		Error:       r.Error,
		Time:        r.Time,
		RequestTime: r.RequestTime,
		Vitals:      r.Vitals,
		PollID:      r.PollID,
		Available:   r.Available,
		System:      r.Astats.System,
	}
}

func copyResultInfos(a []ResultInfo) []ResultInfo {
	b := make([]ResultInfo, len(a), len(a))
	copy(b, a)
	return b
}

func (a ResultInfoHistory) Copy() ResultInfoHistory {
	b := ResultInfoHistory{}
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
	a[r.ID] = pruneInfos(append([]ResultInfo{ToInfo(r)}, a[r.ID]...), limit)
}

// Kbpses is the kbps values of each cache.
type Kbpses map[tc.CacheName]int64

func (a Kbpses) Copy() Kbpses {
	b := Kbpses{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

func (a Kbpses) AddMax(r Result) {
	a[r.ID] = r.PrecomputedData.MaxKbps
}
