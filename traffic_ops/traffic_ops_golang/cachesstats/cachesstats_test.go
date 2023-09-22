package cachesstats

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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestAddStatsInvalidString(t *testing.T) {
	cacheData := make([]CacheData, 0)
	cacheServerStats := make(map[string]tc.ServerStats, 0)
	stats := make(map[string][]tc.ResultStatVal, 0)
	resultStatVals := make([]tc.ResultStatVal, 0)
	resultStatVal := tc.ResultStatVal{
		Span: 0,
		Time: time.Now(),
		Val:  "invalid bandwidth",
	}
	resultStatVals = append(resultStatVals, resultStatVal)
	stats[tc.StatNameBandwidth] = resultStatVals

	resultStatVal = tc.ResultStatVal{
		Span: 0,
		Time: time.Now(),
		Val:  "invalid ats current connection stat",
	}
	resultStatVals = make([]tc.ResultStatVal, 0)
	resultStatVals = append(resultStatVals, resultStatVal)
	stats[ATSCurrentConnectionsStat] = resultStatVals

	cacheServerStats["hostName"] = tc.ServerStats{
		Interfaces: nil,
		Stats:      stats,
	}
	cacheStats := tc.Stats{
		CommonAPIData: tc.CommonAPIData{},
		Caches:        cacheServerStats,
	}
	data := CacheData{
		HostName:    "hostName",
		CacheGroup:  "cacheGroup",
		Status:      "ONLINE",
		Profile:     "profile",
		IP:          util.StrPtr("127.30.30.30"),
		Healthy:     true,
		KBPS:        0,
		Connections: 0,
	}
	cacheData = append(cacheData, data)
	result := addStats(cacheData, cacheStats, "url")
	if len(result) != 1 {
		t.Fatalf("expected a cache stat in the response, but got nothing")
	}
	if result[0].KBPS != 0 || result[0].Connections != 0 {
		t.Errorf("expected 0 KBPS and 0 connections, but got %v and %v respectively", result[0].KBPS, result[0].Connections)
	}
}

func TestAddStatsValidString(t *testing.T) {
	cacheData := make([]CacheData, 0)
	cacheServerStats := make(map[string]tc.ServerStats, 0)
	stats := make(map[string][]tc.ResultStatVal, 0)
	resultStatVals := make([]tc.ResultStatVal, 0)
	resultStatVal := tc.ResultStatVal{
		Span: 0,
		Time: time.Now(),
		Val:  "200",
	}
	resultStatVals = append(resultStatVals, resultStatVal)
	stats[tc.StatNameBandwidth] = resultStatVals

	resultStatVal = tc.ResultStatVal{
		Span: 0,
		Time: time.Now(),
		Val:  "100",
	}
	resultStatVals = make([]tc.ResultStatVal, 0)
	resultStatVals = append(resultStatVals, resultStatVal)
	stats[ATSCurrentConnectionsStat] = resultStatVals

	cacheServerStats["hostName"] = tc.ServerStats{
		Interfaces: nil,
		Stats:      stats,
	}
	cacheStats := tc.Stats{
		CommonAPIData: tc.CommonAPIData{},
		Caches:        cacheServerStats,
	}
	data := CacheData{
		HostName:    "hostName",
		CacheGroup:  "cacheGroup",
		Status:      "ONLINE",
		Profile:     "profile",
		IP:          util.StrPtr("127.30.30.30"),
		Healthy:     true,
		KBPS:        0,
		Connections: 0,
	}
	cacheData = append(cacheData, data)
	result := addStats(cacheData, cacheStats, "url")
	if len(result) != 1 {
		t.Fatalf("expected a cache stat in the response, but got nothing")
	}
	if result[0].KBPS != 200 || result[0].Connections != 100 {
		t.Errorf("expected 200 KBPS and 100 connections, but got %v and %v respectively", result[0].KBPS, result[0].Connections)
	}
}

func TestAddStatsEmptyJsonObject(t *testing.T) {
	cacheData := make([]CacheData, 0)
	cacheServerStats := make(map[string]tc.ServerStats, 0)
	stats := make(map[string][]tc.ResultStatVal, 0)
	resultStatVals := make([]tc.ResultStatVal, 0)
	type testJson struct{}
	var req testJson
	_ = json.Unmarshal([]byte(""), &req)
	resultStatVal := tc.ResultStatVal{
		Span: 0,
		Time: time.Now(),
		Val:  req,
	}
	resultStatVals = append(resultStatVals, resultStatVal)
	stats[tc.StatNameBandwidth] = resultStatVals

	resultStatVal = tc.ResultStatVal{
		Span: 0,
		Time: time.Now(),
		Val:  req,
	}
	resultStatVals = make([]tc.ResultStatVal, 0)
	resultStatVals = append(resultStatVals, resultStatVal)
	stats[ATSCurrentConnectionsStat] = resultStatVals

	cacheServerStats["hostName"] = tc.ServerStats{
		Interfaces: nil,
		Stats:      stats,
	}
	cacheStats := tc.Stats{
		CommonAPIData: tc.CommonAPIData{},
		Caches:        cacheServerStats,
	}
	data := CacheData{
		HostName:    "hostName",
		CacheGroup:  "cacheGroup",
		Status:      "ONLINE",
		Profile:     "profile",
		IP:          util.StrPtr("127.30.30.30"),
		Healthy:     true,
		KBPS:        0,
		Connections: 0,
	}
	cacheData = append(cacheData, data)
	result := addStats(cacheData, cacheStats, "url")
	if len(result) != 1 {
		t.Fatalf("expected a cache stat in the response, but got nothing")
	}
	if result[0].KBPS != 0 || result[0].Connections != 0 {
		t.Errorf("expected 0 KBPS and 0 connections, but got %v and %v respectively", result[0].KBPS, result[0].Connections)
	}
}
