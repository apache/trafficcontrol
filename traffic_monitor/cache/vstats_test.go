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
	"reflect"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

var vstatsData = `{
	"proc.net.dev": "eth0:47907832129 14601260    0    0    0     0          0   790726 728207677726 10210700052    0    0    0     0       0          0",
	"proc.loadavg": "0.30 0.12 0.21 803/863 1421",
	"not_available": false,
	"inf_speed": 70000,
	"stats": {}
}
`

func TestVstatsParse(t *testing.T) {
	reader := strings.NewReader(vstatsData)
	systemStats, statistics, err := vstatsParse("test", reader, nil)
	if err != nil {
		t.Errorf("got error %s", err)
	}
	// stats not implemented yet
	if len(statistics) != 0 {
		t.Errorf("expected statistics to be empty found %v", statistics)
	}
	load := Loadavg{One: 0.3, Five: 0.12, Fifteen: 0.21, CurrentProcesses: 803, TotalProcesses: 863, LatestPID: 1421}
	inf := Interface{Speed: 70000, BytesOut: 728207677726, BytesIn: 47907832129}
	if systemStats.Loadavg != load {
		t.Errorf("got %v want %v", systemStats.Loadavg, load)
	}
	if len(systemStats.Interfaces) != 1 {
		t.Errorf("expected 1 interface got %v", systemStats.Interfaces)
	}
	if systemStats.Interfaces["eth0"] != inf {
		t.Errorf("got %v want %v", systemStats.Interfaces["eth0"], inf)
	}
	if systemStats.NotAvailable {
		t.Errorf("expected NotAvailable to be false")
	}
}

func TestVstatsPrecompute(t *testing.T) {
	to := todata.New()
	to.DeliveryServiceRegexes.DirectMatches["infra.origin.ciab.test"] = "demo1"
	stats := Statistics{
		Interfaces: map[string]Interface{
			"eth0": {
				Speed:    1000,
				BytesIn:  2000,
				BytesOut: 3000,
			},
			"eth1": {
				Speed:    2000,
				BytesIn:  3000,
				BytesOut: 4000,
			},
		},
	}
	counters := map[string]interface{}{
		"infra.origin.ciab.test.200":     float64(5),
		"infra.origin.ciab.test.201":     float64(5),
		"infra.origin.ciab.test.300":     float64(6),
		"infra.origin.ciab.test.404":     float64(3),
		"not-found.origin.ciab.test.404": float64(10), // should not affect other stats and return error
	}
	dsDtats := map[string]*DSStat{
		"demo1": {Status2xx: 10, Status3xx: 6, Status4xx: 3},
	}
	precomputedData := vstatsPrecompute("cache", *to, stats, counters)

	if !reflect.DeepEqual(precomputedData.DeliveryServiceStats, dsDtats) {
		t.Errorf("expected %v got %v", dsDtats, precomputedData.DeliveryServiceStats)
	}
	if precomputedData.OutBytes != 7000 {
		t.Errorf("expected 7000 got %d", precomputedData.OutBytes)
	}
	if precomputedData.MaxKbps != 2000000 {
		t.Errorf("expected 2000000 got %d", precomputedData.MaxKbps)
	}
	if len(precomputedData.Errors) != 1 {
		t.Errorf("expected one error got %v", precomputedData.Errors)
	}

}
