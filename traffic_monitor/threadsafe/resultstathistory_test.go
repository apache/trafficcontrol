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
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/srvhttp"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"

	jsoniter "github.com/json-iterator/go"
)

func randResultStatHistory() ResultStatHistory {
	hist := NewResultStatHistory()

	num := 5
	for i := 0; i < num; i++ {
		hist.Store(tc.CacheName(test.RandStr()), randResultStatValHistory())
	}
	return hist
}

func randResultStatValHistory() ResultStatValHistory {
	a := NewResultStatValHistory()
	num := 5
	numSlice := 5
	for i := 0; i < num; i++ {
		cacheName := test.RandStr()
		vals := []tc.ResultStatVal{}
		for j := 0; j < numSlice; j++ {
			vals = append(vals, randResultStatVal())
		}
		a.Store(cacheName, vals)
	}
	return a
}

func randResultStatVal() tc.ResultStatVal {
	return tc.ResultStatVal{
		Val:  uint64(test.RandInt64()),
		Time: time.Now(),
		Span: uint64(test.RandInt64()),
	}
}

func randResultInfoHistory() cache.ResultInfoHistory {
	hist := cache.ResultInfoHistory{}

	num := 5
	infNum := 5
	for i := 0; i < num; i++ {
		cacheName := tc.CacheName(test.RandStr())
		for j := 0; j < infNum; j++ {
			hist[cacheName] = append(hist[cacheName], randResultInfo())
		}
	}
	return hist
}

func randResultInfo() cache.ResultInfo {
	return cache.ResultInfo{
		ID:          test.RandStr(),
		Error:       fmt.Errorf(test.RandStr()),
		Time:        time.Now(),
		RequestTime: time.Millisecond * time.Duration(test.RandInt()),
		Vitals:      randVitals(),
		PollID:      uint64(test.RandInt64()),
		Available:   test.RandBool(),
	}
}

func randVitals() cache.Vitals {
	return cache.Vitals{
		LoadAvg:    test.RandFloat64(),
		BytesOut:   test.RandUint64(),
		BytesIn:    test.RandUint64(),
		KbpsOut:    test.RandInt64(),
		MaxKbpsOut: test.RandInt64(),
	}
}

type DummyFilterNever struct {
}

func (DummyFilterNever) UseStat(string) bool {
	return false
}

func (DummyFilterNever) UseInterfaceStat(string) bool {
	return false
}

func (DummyFilterNever) UseCache(tc.CacheName) bool {
	return false
}

func (DummyFilterNever) WithinStatHistoryMax(uint64) bool {
	return false
}
func TestLegacyStatsMarshall(t *testing.T) {
	statHist := randResultStatHistory()
	infHist := randResultInfoHistory()
	filter := DummyFilterNever{}
	params := url.Values{}
	beforeStatsMarshall := time.Now()
	bytes, err := LegacyStatsMarshall(statHist, infHist, tc.CRStates{}, tc.TrafficMonitorConfigMap{}, cache.Kbpses{}, filter, params)
	afterStatsMarshall := time.Now()
	if err != nil {
		t.Fatalf("StatsMarshall return expected nil err, actual err: %v", err)
	}

	stats := tc.LegacyStats{}
	json := jsoniter.ConfigFastest // TODO make configurable
	if err := json.Unmarshal(bytes, &stats); err != nil {
		t.Fatalf("unmarshalling expected nil err, actual err: %v", err)
	}

	if stats.CommonAPIData.QueryParams != "" {
		t.Errorf(`unmarshalling stats.CommonAPIData.QueryParams expected "", actual %v`, stats.CommonAPIData.QueryParams)
	}

	statsDate, err := time.Parse(srvhttp.CommonAPIDataDateFormat, stats.CommonAPIData.DateStr)
	if err != nil {
		t.Errorf(`stats.CommonAPIData.DateStr expected format %v, actual %v`, srvhttp.CommonAPIDataDateFormat, stats.CommonAPIData.DateStr)
	}
	if beforeStatsMarshall.Truncate(time.Second).After(statsDate) || statsDate.Truncate(time.Second).After(afterStatsMarshall.Truncate(time.Second)) { // round to second, because CommonAPIDataDateFormat is second-precision
		t.Errorf(`unmarshalling stats.CommonAPIData.DateStr expected between %v and %v, actual %v`, beforeStatsMarshall, afterStatsMarshall, stats.CommonAPIData.DateStr)
	}
	if len(stats.Caches) > 0 {
		t.Errorf(`unmarshalling stats.Caches expected empty, actual %+v`, stats.Caches)
	}
}

func TestStatsMarshall(t *testing.T) {
	statHist := randResultStatHistory()
	infHist := randResultInfoHistory()
	filter := DummyFilterNever{}
	params := url.Values{}
	beforeStatsMarshall := time.Now()
	bytes, err := StatsMarshall(statHist, infHist, tc.CRStates{}, tc.TrafficMonitorConfigMap{}, cache.Kbpses{}, filter, params)
	afterStatsMarshall := time.Now()
	if err != nil {
		t.Fatalf("StatsMarshall return expected nil err, actual err: %v", err)
	}

	stats := tc.Stats{}
	json := jsoniter.ConfigFastest // TODO make configurable
	if err := json.Unmarshal(bytes, &stats); err != nil {
		t.Fatalf("unmarshalling expected nil err, actual err: %v", err)
	}

	if stats.CommonAPIData.QueryParams != "" {
		t.Errorf(`unmarshalling stats.CommonAPIData.QueryParams expected "", actual %v`, stats.CommonAPIData.QueryParams)
	}

	statsDate, err := time.Parse(srvhttp.CommonAPIDataDateFormat, stats.CommonAPIData.DateStr)
	if err != nil {
		t.Errorf(`stats.CommonAPIData.DateStr expected format %v, actual %v`, srvhttp.CommonAPIDataDateFormat, stats.CommonAPIData.DateStr)
	}
	if beforeStatsMarshall.Truncate(time.Second).After(statsDate) || statsDate.Truncate(time.Second).After(afterStatsMarshall.Truncate(time.Second)) { // round to second, because CommonAPIDataDateFormat is second-precision
		t.Errorf(`unmarshalling stats.CommonAPIData.DateStr expected between %v and %v, actual %v`, beforeStatsMarshall, afterStatsMarshall, stats.CommonAPIData.DateStr)
	}
	if len(stats.Caches) > 0 {
		t.Errorf(`unmarshalling stats.Caches expected empty, actual %+v`, stats.Caches)
	}
}

func TestSystemComputedStats(t *testing.T) {
	stats := cache.ComputedStats()

	for stat, function := range stats {
		if strings.HasPrefix(stat, "system.") {
			computedStat := function(cache.ResultInfo{}, tc.TrafficServer{}, tc.TMProfile{}, tc.IsAvailable{})
			_, err := newStatEqual([]tc.ResultStatVal{{Val: float64(0)}}, computedStat)
			if err != nil {
				t.Errorf("expected no errors from newStatEqual: %s", err)
			}
		}
	}
}

func TestCompareAndAppendStatForInterface(t *testing.T) {
	var errs strings.Builder
	var limit uint64 = 1
	stat := interfaceStat{
		InterfaceName: "test",
		Stat:          uint64(5),
		StatName:      "test",
		Time:          time.Now(),
	}

	history := compareAndAppendStatForInterface(nil, errs, limit, stat)
	if errs.Len() > 0 {
		t.Errorf("Unexpected errors comparing previously non-existent interface stat: %s", errs.String())
	}
	if len(history) == 0 {
		t.Fatal("Empty history after comparing previously non-existent interface stat")
	}
	if len(history) > 1 {
		t.Fatalf("Too many stats returned from comparing previously non-existent interface stat: %d", len(history))
	}

	result := history[0]
	if result.Span != 1 {
		t.Errorf("Incorrect span comparing previously non-existent interface stat; want: 1, got: %d", result.Span)
	}

	if result.Time != stat.Time {
		t.Errorf("Incorrect time comparing previously non-existent interface stat; want: %v, got: %v", stat.Time, result.Time)
	}

	if v, ok := result.Val.(uint64); !ok {
		t.Errorf("Incorrect value type from comparing previously non-existent interface stat; want: uint64, got: %T", result.Val)
	} else if v != stat.Stat {
		t.Errorf("Incorrect value from comparing previously non-existent interface stat; want: %d, got: %d", stat.Stat, v)
	}

	errs.Reset()

	history = compareAndAppendStatForInterface(history, errs, limit, stat)
	if errs.Len() > 0 {
		t.Errorf("Unexpected errors comparing previously non-existent interface stat: %s", errs.String())
	}
	if len(history) == 0 {
		t.Fatal("Empty history after comparing previously non-existent interface stat")
	}
	if len(history) > 1 {
		t.Fatalf("Too many stats returned from comparing previously non-existent interface stat: %d", len(history))
	}

	result = history[0]
	if result.Span != 2 {
		t.Errorf("Incorrect span comparing previously non-existent interface stat; want: 2, got: %d", result.Span)
	}

	if result.Time != stat.Time {
		t.Errorf("Incorrect time comparing previously non-existent interface stat; want: %v, got: %v", stat.Time, result.Time)
	}

	if v, ok := result.Val.(uint64); !ok {
		t.Errorf("Incorrect value type from comparing previously non-existent interface stat; want: uint64, got: %T", result.Val)
	} else if v != stat.Stat {
		t.Errorf("Incorrect value from comparing previously non-existent interface stat; want: %d, got: %d", stat.Stat, v)
	}

	errs.Reset()
	stat.Stat = uint64(6)

	history = compareAndAppendStatForInterface(history, errs, limit, stat)
	if errs.Len() > 0 {
		t.Errorf("Unexpected errors comparing previously non-existent interface stat: %s", errs.String())
	}
	if len(history) == 0 {
		t.Fatal("Empty history after comparing previously non-existent interface stat")
	}
	if len(history) > 1 {
		t.Fatalf("Too many stats returned from comparing previously non-existent interface stat: %d", len(history))
	}

	result = history[0]
	if result.Span != 1 {
		t.Errorf("Incorrect span comparing previously non-existent interface stat; want: 1, got: %d", result.Span)
	}

	if result.Time != stat.Time {
		t.Errorf("Incorrect time comparing previously non-existent interface stat; want: %v, got: %v", stat.Time, result.Time)
	}

	if v, ok := result.Val.(uint64); !ok {
		t.Errorf("Incorrect value type from comparing previously non-existent interface stat; want: uint64, got: %T", result.Val)
	} else if v != stat.Stat {
		t.Errorf("Incorrect value from comparing previously non-existent interface stat; want: %d, got: %d", stat.Stat, v)
	}
}
