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
	"net/url"
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/srvhttp"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/todata"
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

func TestStatsMarshall(t *testing.T) {
	hist := randResultStatHistory()
	filter := DummyFilterNever{}
	params := url.Values{}
	info := randResultInfoHistory()
	combinedStates := tc.CRStates{}
	monitorConfig := tc.TrafficMonitorConfigMap{}
	statMaxKbpses := randKbpses()

	beforeStatsMarshall := time.Now()
	bytes, err := StatsMarshall(hist, info, combinedStates, monitorConfig, statMaxKbpses, filter, params)
	afterStatsMarshall := time.Now()

	if err != nil {
		t.Fatalf("StatsMarshall return expected nil err, actual err: %v", err)
	}
	// if len(bytes) > 0 {
	// 	t.Errorf("expected empty bytes, actual: %v", string(bytes))
	// }

	stats := Stats{}
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
	if beforeStatsMarshall.Truncate(time.Second).After(statsDate) || statsDate.After(afterStatsMarshall.Truncate(time.Second).Add(time.Second)) { // +/- 1 second, because CommonAPIDataDateFormat is second-precision
		t.Errorf(`unmarshalling stats.CommonAPIData.DateStr expected between %v and %v, actual (%v) %v`, beforeStatsMarshall, afterStatsMarshall, statsDate, stats.CommonAPIData.DateStr)
	}
	if len(stats.Caches) > 0 {
		t.Errorf(`unmarshalling stats.Caches expected empty, actual %+v`, stats.Caches)
	}
}
