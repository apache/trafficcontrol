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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// CacheAvailableStatusReported is the status string returned by caches set to "reported" in Traffic Ops.
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

// ResultStatHistory is a map[cache][statName]val
type ResultStatHistory map[tc.CacheName]ResultStatValHistory

type ResultStatValHistory map[string][]ResultStatVal

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
	return json.Marshal(&v)
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

// newStatEqual Returns whether the given stat is equal to the latest stat in history. If len(history)==0, this returns false without error. If the given stat is not a JSON primitive (string, number, bool), this returns an error. We explicitly refuse to compare arrays and objects, for performance.
func newStatEqual(history []ResultStatVal, stat interface{}) (bool, error) {
	if len(history) == 0 {
		return false, nil // if there's no history, it's "not equal", i.e. store this new history
	}
	switch stat.(type) {
	case string:
	case float64:
	case bool:
	default:
		return false, fmt.Errorf("incomparable stat type %T", stat)
	}
	switch history[0].Val.(type) {
	case string:
	case float64:
	case bool:
	default:
		return false, fmt.Errorf("incomparable history stat type %T", stat)
	}
	return stat == history[0].Val, nil
}

func (a ResultStatHistory) Add(r Result, limit uint64) error {
	errStrs := ""
	for statName, statVal := range r.Astats.Ats {
		statHistory := a[r.ID][statName]
		// If the new stat value is the same as the last, update the time and increment the span. Span is the number of polls the latest value has been the same, and hence the length of time it's been the same is span*pollInterval.
		if ok, err := newStatEqual(statHistory, statVal); err != nil {
			errStrs += "cannot add stat " + statName + ": " + err.Error() + "; "
		} else if ok {
			statHistory[0].Time = r.Time
			statHistory[0].Span++
		} else {
			resultVal := ResultStatVal{
				Val:  statVal,
				Time: r.Time,
				Span: 1,
			}
			statHistory = pruneStats(append([]ResultStatVal{resultVal}, statHistory...), limit)
		}
		if _, ok := a[r.ID]; !ok {
			a[r.ID] = ResultStatValHistory{}
		}
		a[r.ID][statName] = statHistory // TODO determine if necessary for the first conditional
	}

	if errStrs != "" {
		return errors.New("some stats could not be added: " + errStrs[:len(errStrs)-2])
	}
	return nil
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
