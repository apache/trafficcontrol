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

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
)

// CacheAvailableStatusReported is the status string returned by caches set to "reported" in Traffic Ops.
// TODO put somewhere more generic
const AvailableStatusReported = "REPORTED"

// CacheAvailableStatus is the available status of the given cache. It includes a boolean available/unavailable flag, and a descriptive string.
type AvailableStatus struct {
	Available bool
	Status    string
	Why       string
}

// CacheAvailableStatuses is the available status of each cache.
type AvailableStatuses map[enum.CacheName]AvailableStatus

// Copy copies this CacheAvailableStatuses. It does not modify, and thus is safe for multiple reader goroutines.
func (a AvailableStatuses) Copy() AvailableStatuses {
	b := AvailableStatuses(map[enum.CacheName]AvailableStatus{})
	for k, v := range a {
		b[k] = v
	}
	return b
}

// Event represents an event change in aggregated data. For example, a cache being marked as unavailable.
type Event struct {
	Index       uint64         `json:"index"`
	Time        int64          `json:"time"`
	Description string         `json:"description"`
	Name        enum.CacheName `json:"name"`
	Hostname    enum.CacheName `json:"hostname"`
	Type        string         `json:"type"`
	Available   bool           `json:"isAvailable"`
}

// ResultHistory is a map of cache names, to an array of result history from each cache.
type ResultHistory map[enum.CacheName][]Result

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

// ResultStatHistory is a map[cache][statName]val
type ResultStatHistory map[enum.CacheName]ResultStatValHistory

type ResultStatValHistory map[string][]ResultStatVal

// ResultStatVal is the value of an individual stat returned from a poll. Time is the time this stat was returned.
// Span is the number of polls this stat has been the same. For example, if History is set to 100, and the last 50 polls had the same value for this stat (but none of the previous 50 were the same), this stat's map value slice will actually contain 51 entries, and the first entry will have the value, the time of the last poll, and a Span of 50. Assuming the poll time is every 8 seconds, users will then know, looking at the Span, that the value was unchanged for the last 50*8=400 seconds.
// JSON values are all strings, for the TM1.0 /publish/CacheStats API.
type ResultStatVal struct {
	Val  interface{} `json:"value,string"`
	Time TM1Time     `json:"time,string"`
	Span uint64      `json:"span,string"`
}

// TM1Time provides a custom MarshalJSON func to serialise a time.Time into milliseconds since the epoch, as served in Traffic Monitor 1.x APIs
// TODO move somewhere more generic (enum?)
type TM1Time time.Time

func (t *TM1Time) MarshalJSON() ([]byte, error) {
	NanosecondsPerMillisecond := int64(1000000)
	it := (*time.Time)(t).UnixNano() / NanosecondsPerMillisecond
	return []byte(fmt.Sprintf("%d", it)), nil
}

func copyResultStatVals(a []ResultStatVal) []ResultStatVal {
	b := make([]ResultStatVal, len(a), len(a))
	copy(b, a)
	return b
}

func copyResultStatValHistory(a ResultStatValHistory) ResultStatValHistory {
	b := ResultStatValHistory{}
	for k, v := range a {
		b[k] = copyResultStatVals(v) // TODO determine if necessary
	}
	return b
}

func (a ResultStatHistory) Copy() ResultStatHistory {
	b := ResultStatHistory{}
	for k, v := range a {
		b[k] = copyResultStatValHistory(v)
	}
	return b
}

func pruneStats(history []ResultStatVal, limit uint64) []ResultStatVal {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}

func (a ResultStatHistory) Add(r Result, limit uint64) {
	for statName, statVal := range r.Astats.Ats {
		statHistory := a[r.ID][statName]
		// If the new stat value is the same as the last, update the time and increment the span. Span is the number of polls the latest value has been the same, and hence the length of time it's been the same is span*pollInterval.
		if len(statHistory) > 0 && statHistory[0].Val == statVal {
			statHistory[0].Time = TM1Time(r.Time)
			statHistory[0].Span++
		} else {
			resultVal := ResultStatVal{
				Val:  statVal,
				Time: TM1Time(r.Time),
				Span: 1,
			}
			statHistory = pruneStats(append([]ResultStatVal{resultVal}, statHistory...), limit)
		}
		if _, ok := a[r.ID]; !ok {
			a[r.ID] = ResultStatValHistory{}
		}
		a[r.ID][statName] = statHistory // TODO determine if necessary for the first conditional
	}
}

type ResultInfoHistory map[enum.CacheName][]ResultInfo

// ResultInfo contains all the non-stat result info. This includes the cache ID, any errors, the time of the poll, the request time duration, Astats System (Vitals), Poll ID, and Availability.
type ResultInfo struct {
	ID          enum.CacheName
	Error       error
	Time        time.Time
	RequestTime time.Duration
	Vitals      Vitals
	PollID      uint64
	Available   bool
}

func toInfo(r Result) ResultInfo {
	return ResultInfo{
		ID:          r.ID,
		Error:       r.Error,
		Time:        r.Time,
		RequestTime: r.RequestTime,
		Vitals:      r.Vitals,
		PollID:      r.PollID,
		Available:   r.Available,
	}
}

func toInfos(rs []Result) []ResultInfo {
	infos := make([]ResultInfo, len(rs), len(rs))
	for i, r := range rs {
		infos[i] = toInfo(r)
	}
	return infos
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
	a[r.ID] = pruneInfos(append([]ResultInfo{toInfo(r)}, a[r.ID]...), limit)
}

// Kbpses is the kbps values of each cache.
type Kbpses map[enum.CacheName]int64

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
