package main

/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	influx "github.com/influxdata/influxdb/client/v2"
)

func TestCalcCacheValuesWithInvalidValue(t *testing.T) {
	stats := make(map[string][]tc.ResultStatVal)
	caches := make(map[tc.CacheName]map[string][]tc.ResultStatVal)
	cacheMap := make(map[string]tc.ServerV5)
	resultSatsVal := []tc.ResultStatVal{
		{
			Span: 0,
			Time: time.Now(),
			Val:  "invalid stat",
		},
	}
	stats["maxKbps"] = resultSatsVal
	caches["cache1"] = stats
	cacheMap["cache1"] = tc.ServerV5{}
	config := StartupConfig{
		BpsChan: make(chan influx.BatchPoints),
	}
	legacyStats := tc.LegacyStats{
		CommonAPIData: tc.CommonAPIData{},
		Caches:        caches,
	}
	data, err := json.Marshal(legacyStats)
	if err != nil {
		t.Fatalf("couldn't marshal struct %v", caches)
	}
	go calcCacheValues(data, "cdn", 0, cacheMap, config)
	result := <-config.BpsChan
	if len(result.Points()) == 0 {
		t.Fatalf("expected one point in the result, got none")
	}
	fields, err := result.Points()[0].Fields()
	if err != nil {
		t.Fatalf("couldn't read the fields of the result: %v", err.Error())
	}
	if val, ok := fields["value"]; !ok {
		t.Fatalf("couldn't find a 'value' field")
	} else {
		if val.(float64) != 0.0 {
			t.Errorf("expected invalid stat to result in a value of 0.0, but got %v instead", val.(float64))
		}
	}
}

func TestCalcCacheValuesWithEmptyInterface(t *testing.T) {
	stats := make(map[string][]tc.ResultStatVal)
	caches := make(map[tc.CacheName]map[string][]tc.ResultStatVal)
	cacheMap := make(map[string]tc.ServerV5)
	resultSatsVal := []tc.ResultStatVal{
		{
			Span: 0,
			Time: time.Now(),
		},
	}
	stats["maxKbps"] = resultSatsVal
	caches["cache1"] = stats
	cacheMap["cache1"] = tc.ServerV5{}
	config := StartupConfig{
		BpsChan: make(chan influx.BatchPoints),
	}
	legacyStats := tc.LegacyStats{
		CommonAPIData: tc.CommonAPIData{},
		Caches:        caches,
	}
	data, err := json.Marshal(legacyStats)
	if err != nil {
		t.Fatalf("couldn't marshal struct %v", caches)
	}
	go calcCacheValues(data, "cdn", 0, cacheMap, config)
	result := <-config.BpsChan
	if len(result.Points()) == 0 {
		t.Fatalf("expected one point in the result, got none")
	}
	fields, err := result.Points()[0].Fields()
	if err != nil {
		t.Fatalf("couldn't read the fields of the result: %v", err.Error())
	}
	if val, ok := fields["value"]; !ok {
		t.Fatalf("couldn't find a 'value' field")
	} else {
		if val.(float64) != 0.0 {
			t.Errorf("expected invalid stat to result in a value of 0.0, but got %v instead", val.(float64))
		}
	}
}

func TestCalcCacheValues(t *testing.T) {
	stats := make(map[string][]tc.ResultStatVal)
	caches := make(map[tc.CacheName]map[string][]tc.ResultStatVal)
	cacheMap := make(map[string]tc.ServerV5)
	resultSatsVal := []tc.ResultStatVal{
		{
			Span: 0,
			Time: time.Now(),
			Val:  "25.56",
		},
	}
	stats["maxKbps"] = resultSatsVal
	caches["cache1"] = stats
	cacheMap["cache1"] = tc.ServerV5{}
	config := StartupConfig{
		BpsChan: make(chan influx.BatchPoints),
	}
	legacyStats := tc.LegacyStats{
		CommonAPIData: tc.CommonAPIData{},
		Caches:        caches,
	}
	data, err := json.Marshal(legacyStats)
	if err != nil {
		t.Fatalf("couldn't marshal struct %v", caches)
	}
	go calcCacheValues(data, "cdn", 0, cacheMap, config)
	result := <-config.BpsChan
	if len(result.Points()) == 0 {
		t.Fatalf("expected one point in the result, got none")
	}
	fields, err := result.Points()[0].Fields()
	if err != nil {
		t.Fatalf("couldn't read the fields of the result: %v", err.Error())
	}
	if val, ok := fields["value"]; !ok {
		t.Fatalf("couldn't find a 'value' field")
	} else {
		if val.(float64) != 25.56 {
			t.Errorf("expected invalid stat to result in a value of 0.0, but got %v instead", val.(float64))
		}
	}
}

func TestSetHealthURLs(t *testing.T) {
	cfg := StartupConfig{
		StatusToMon: tc.CacheStatusOffline.String(),
	}
	runningCfg := RunningConfig{}
	runningCfg.CacheMap = map[string]tc.ServerV5{
		"tm1": {
			CDN:        "foo",
			DomainName: "example.org",
			HostName:   "tm1",
			Status:     tc.CacheStatusOffline.String(),
			TCPPort:    util.Ptr(8080),
			Type:       tc.MonitorTypeName,
		},
		"tm2": {
			CDN:        "bar",
			DomainName: "example.org",
			HostName:   "tm2",
			Status:     tc.CacheStatusOffline.String(),
			TCPPort:    util.Ptr(8080),
			Type:       tc.MonitorTypeName,
		},
		"tm3": {
			CDN:        "foo",
			DomainName: "example.org",
			HostName:   "tm3",
			Status:     tc.CacheStatusOffline.String(),
			TCPPort:    util.Ptr(8080),
			Type:       tc.MonitorTypeName,
		},
		"tm4": {
			CDN:        "bar",
			DomainName: "example.org",
			HostName:   "tm4",
			Status:     tc.CacheStatusOffline.String(),
			TCPPort:    util.Ptr(8080),
			Type:       tc.MonitorTypeName,
		},
		"tm5": {
			CDN:        "foo",
			DomainName: "example.org",
			HostName:   "tm5",
			Status:     tc.CacheStatusOnline.String(),
			TCPPort:    util.Ptr(8080),
			Type:       tc.MonitorTypeName,
		},
		"tm6": {
			CDN:        "bar",
			DomainName: "example.org",
			HostName:   "tm6",
			Status:     tc.CacheStatusOnline.String(),
			TCPPort:    util.Ptr(8080),
			Type:       tc.MonitorTypeName,
		},
		"cache1": {
			CDN:        "foo",
			DomainName: "example.org",
			HostName:   "cache1",
			Status:     tc.CacheStatusOffline.String(),
			TCPPort:    util.Ptr(8080),
			Type:       tc.CacheTypeEdge.String(),
		},
	}
	setHealthURLs(cfg, &runningCfg, "/cacheStats", "/dsStats")
	expected := map[string]map[string][]string{
		"foo": {
			"CacheStats": []string{
				"http://tm1.example.org:8080/cacheStats",
				"http://tm3.example.org:8080/cacheStats",
			},
			"DsStats": []string{
				"http://tm1.example.org:8080/dsStats",
				"http://tm3.example.org:8080/dsStats",
			},
		},
		"bar": {
			"CacheStats": []string{
				"http://tm2.example.org:8080/cacheStats",
				"http://tm4.example.org:8080/cacheStats",
			},
			"DsStats": []string{
				"http://tm2.example.org:8080/dsStats",
				"http://tm4.example.org:8080/dsStats",
			},
		},
	}
	for _, toBeSorted := range []map[string]map[string][]string{
		expected,
		runningCfg.HealthUrls,
	} {
		for _, healthURLs := range toBeSorted {
			for _, urls := range healthURLs {
				sort.Strings(urls)
			}
		}
	}
	if !reflect.DeepEqual(expected, runningCfg.HealthUrls) {
		t.Errorf("expected: %+v, actual: %+v", expected, runningCfg.HealthUrls)
	}
}
