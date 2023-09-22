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
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/poller"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

func TestHandlerPrecompute(t *testing.T) {
	if NewHandler().Precompute() {
		t.Errorf("expected NewHandler().Precompute() false, actual true")
	}
	if !NewPrecomputeHandler(todata.NewThreadsafe()).Precompute() {
		t.Errorf("expected NewPrecomputeHandler().Precompute() true, actual false")
	}
}

type DummyFilterNever struct {
}

func (f DummyFilterNever) UseStat(name string) bool {
	return false
}

func (f DummyFilterNever) UseCache(name tc.CacheName) bool {
	return false
}

func (f DummyFilterNever) WithinStatHistoryMax(i int) bool {
	return false
}

var AggregatedStats = Vitals{
	KbpsOut:    1500000,
	MaxKbpsOut: 2500000,
	LoadAvg:    100000,
}

var StatToDesiredValue = map[string]float64{
	"availableBandwidthInKbps": 1000000,
	"availableBandwidthInMbps": 1000,
	"bandwidth":                1500000,
	"kbps":                     1500000,
	"gbps":                     1.5,
	"loadavg":                  100000,
	"maxKbps":                  2500000,
}

func TestComputeAggregateStats(t *testing.T) {
	computedStats := ComputedStats()

	for stat, want := range StatToDesiredValue {
		got := computedStats[stat](ResultInfo{Vitals: AggregatedStats}, tc.TrafficServer{}, tc.TMProfile{}, tc.IsAvailable{})
		if numGot, ok := util.ToNumeric(got); ok {
			if numGot != want {
				t.Errorf("ComputedStats[\"%v\"] return %v instead of %v", stat, got, want)
			}
		} else {
			t.Errorf(`ComputedStats["%s"] returned non-numeric value: %v`, stat, got)
		}
	}
}

func TestComputeStatGbps(t *testing.T) {
	serverInfo := tc.TrafficServer{}
	serverProfile := tc.TMProfile{}
	combinedState := tc.IsAvailable{}
	computedStats := ComputedStats()

	for stat, want := range StatToDesiredValue {
		got := computedStats[stat](ResultInfo{Vitals: AggregatedStats}, serverInfo, serverProfile, combinedState)
		if numGot, ok := util.ToNumeric(got); ok && numGot != want {
			t.Errorf("ComputedStats[\"%v\"] return %v instead of %v", stat, got, want)
		}
	}
}

func TestParseAndDecode(t *testing.T) {
	file, err := ioutil.ReadFile("stats_over_http.json")
	if err != nil {
		t.Fatal(err)
	}

	pl := &poller.HTTPPollCtx{HTTPHeader: http.Header{}}
	ctx := interface{}(pl)
	ctx.(*poller.HTTPPollCtx).HTTPHeader.Set(rfc.ContentType, rfc.ApplicationJSON)

	decoder, err := GetDecoder("stats_over_http")
	if err != nil {
		t.Errorf("decoder error, expected: nil, got: %v", err)
	}

	_, miscStats, err := decoder.Parse("1", bytes.NewReader(file), ctx)
	if err != nil {
		t.Errorf("decoder parse error, expected: nil, got: %v", err)
	}

	if len(miscStats) < 1 {
		t.Errorf("empty miscStats structure")
	}

	if val, ok := miscStats["plugin.system_stats.timestamp_ms"]; ok {
		valType := reflect.TypeOf(val)
		if valType.Kind() != reflect.String {
			t.Errorf("type mismatch, expected: string, got:%s", valType)
		}
		val1, _ := parseNumericStat(val)
		if val1 != uint64(1684784877939) {
			t.Errorf("unable to read `plugin.system_stats.timestamp_ms`")
		}
	}
	if val, ok := miscStats["plugin.system_stats.timestamp_ms_float64"]; ok {
		valType := reflect.TypeOf(val)
		if valType.Kind() != reflect.Float64 {
			t.Errorf("type mismatch, expected: float64, got:%s", valType)
		}
		val1, _ := parseNumericStat(val)
		if val1 != uint64(1684784877939) {
			t.Errorf("unable to read `plugin.system_stats.timestamp_ms_float64`")
		}
	}
}
