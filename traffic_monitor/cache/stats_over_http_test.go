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
	"os"
	"testing"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/poller"
)

func TestStatsOverHTTPParse(t *testing.T) {
	fd, err := os.Open("stats_over_http.json")
	if err != nil {
		t.Fatal(err)
	}

	pl := &poller.HTTPPollCtx{HTTPHeader: http.Header{}}
	ctx := interface{}(pl)
	ctx.(*poller.HTTPPollCtx).HTTPHeader.Set("Content-Type", "text/json")

	stats, misc, err := statsOverHTTPParse("test", fd, ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Floating-Point arithmetic...
	if misc["plugin.remap_stats.edge-cache-0.delivery.service.zero.in_bytes"] != float64(296727207) {
		t.Errorf("Expected 296727207 for remap_stats edge-cache in_bytes, got %.10f", misc["plugin.remap_stats.edge-cache-0.delivery.service.zero.in_bytes"])
	}
	if stats.Loadavg.One <= 0.092773437 || stats.Loadavg.One >= 0.092773439 {
		t.Errorf("Incorrect one-minute loadavg, expected roughly 0.092773438, got '%.10f'", stats.Loadavg.One)
	}
	if stats.Loadavg.Five <= 0.25439453 || stats.Loadavg.Five >= 0.254394532 {
		t.Errorf("Incorrect five-minute loadavg, expected roughly 0.254394531, got %.10f", stats.Loadavg.Five)
	}
	if stats.Loadavg.Fifteen <= 0.639160155 || stats.Loadavg.Fifteen >= 639160157 {
		t.Errorf("Incorrect fifteen-minute loadavg, expected roughly 0.639160156, got %.10f", stats.Loadavg.Fifteen)
	}
	if stats.Loadavg.TotalProcesses != 803 {
		t.Errorf("Incorrect current_processes, expected 803, got %d", stats.Loadavg.CurrentProcesses)
	}

	if len(stats.Interfaces) != 1 {
		t.Errorf("Expected exactly one interface, got %d", len(stats.Interfaces))
		if len(stats.Interfaces) < 1 {
			t.FailNow()
		}
	}

	found := false
	for name, iface := range stats.Interfaces {
		if name != "docker0" {
			t.Errorf("Found unexpected network interface '%s'", name)
			continue
		}
		found = true

		if iface.Speed != 70000 {
			t.Errorf("Incorrect interface speed, expected 70000, got %d", iface.Speed)
		}
		if iface.BytesIn != 4363732 {
			t.Errorf("Incorrect interface rx_bytes, expected 4363732, got %d", iface.BytesIn)
		}
		if iface.BytesOut != 237634637 {
			t.Errorf("Incorrect interface tx_bytes, expceted 237634637, got %d", iface.BytesOut)
		}
	}
	if !found {
		t.Error("Didn't find the expected 'docker0' network interface")
	}

}

func TestStatsOverHTTPParseCSV(t *testing.T) {
	fd, err := os.Open("stats_over_http.csv")
	if err != nil {
		t.Fatal(err)
	}

	pl := &poller.HTTPPollCtx{HTTPHeader: http.Header{}}
	ctx := interface{}(pl)
	ctx.(*poller.HTTPPollCtx).HTTPHeader.Set("Content-Type", "text/csv")

	stats, misc, err := statsOverHTTPParse("test", fd, ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Floating-Point arithmetic...
	if misc["plugin.remap_stats.edge-cache-0.delivery.service.zero.in_bytes"] != float64(296727207) {
		t.Errorf("Expected 296727207 for remap_stats edge-cache in_bytes, got %.10f", misc["plugin.remap_stats.edge-cache-0.delivery.service.zero.in_bytes"])
	}

	if stats.Loadavg.One <= 0.092773437 || stats.Loadavg.One >= 0.092773439 {
		t.Errorf("Incorrect one-minute loadavg, expected roughly 0.092773438, got '%.10f'", stats.Loadavg.One)
	}
	if stats.Loadavg.Five <= 0.25439453 || stats.Loadavg.Five >= 0.254394532 {
		t.Errorf("Incorrect five-minute loadavg, expected roughly 0.254394531, got %.10f", stats.Loadavg.Five)
	}
	if stats.Loadavg.Fifteen <= 0.639160155 || stats.Loadavg.Fifteen >= 639160157 {
		t.Errorf("Incorrect fifteen-minute loadavg, expected roughly 0.639160156, got %.10f", stats.Loadavg.Fifteen)
	}
	if stats.Loadavg.TotalProcesses != 803 {
		t.Errorf("Incorrect current_processes, expected 803, got %d", stats.Loadavg.TotalProcesses)
	}

	if len(stats.Interfaces) != 1 {
		t.Errorf("Expected exactly one interface, got %d", len(stats.Interfaces))
		if len(stats.Interfaces) < 1 {
			t.FailNow()
		}
	}

	found := false
	for name, iface := range stats.Interfaces {
		if name != "docker0" {
			t.Errorf("Found unexpected network interface '%s'", name)
			continue
		}
		found = true

		if iface.Speed != 70000 {
			t.Errorf("Incorrect interface speed, expected 70000, got %d", iface.Speed)
		}
		if iface.BytesIn != 4363732 {
			t.Errorf("Incorrect interface rx_bytes, expected 4363732, got %d", iface.BytesIn)
		}
		if iface.BytesOut != 237634637 {
			t.Errorf("Incorrect interface tx_bytes, expceted 237634637, got %d", iface.BytesOut)
		}
	}
	if !found {
		t.Error("Didn't find the expected 'docker0' network interface")
	}

}

func BenchmarkStatsJson(b *testing.B) {
	file, err := ioutil.ReadFile("stats_over_http.json")
	if err != nil {
		b.Fatal(err)
	}

	pl := &poller.HTTPPollCtx{HTTPHeader: http.Header{}}
	ctx := interface{}(pl)
	ctx.(*poller.HTTPPollCtx).HTTPHeader.Set("Content-Type", "text/json")
	// Reset benchmark timer to not include reading the file
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := statsOverHTTPParse("test", bytes.NewReader(file), ctx)

		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkStatsCSV(b *testing.B) {
	file, err := ioutil.ReadFile("stats_over_http.csv")
	if err != nil {
		b.Fatal(err)
	}

	pl := &poller.HTTPPollCtx{HTTPHeader: http.Header{}}
	ctx := interface{}(pl)
	ctx.(*poller.HTTPPollCtx).HTTPHeader.Set("Content-Type", "text/csv")
	// Reset benchmark timer to not include reading the file
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := statsOverHTTPParse("test", bytes.NewReader(file), ctx)

		if err != nil {
			b.Error(err)
		}
	}
}
