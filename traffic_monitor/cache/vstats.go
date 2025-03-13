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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	registerDecoder("vstats", vstatsParse, vstatsPrecompute)
}

// Vstats holds Varnish cache statistics
type Vstats struct {
	ProcNetDev   string                 `json:"proc.net.dev"`
	ProcLoadAvg  string                 `json:"proc.loadavg"`
	NotAvailable bool                   `json:"not_available"`
	InfSpeed     int64                  `json:"inf_speed"`
	Stats        map[string]interface{} `json:"stats"`
}

func vstatsParse(cacheName string, r io.Reader, _ interface{}) (Statistics, map[string]interface{}, error) {
	var stats Statistics

	if r == nil {
		log.Warnf("%s handler got nil reader", cacheName)
		return stats, nil, errors.New("handler got nil reader")
	}

	var vstats Vstats
	json := jsoniter.ConfigFastest

	if err := json.NewDecoder(r).Decode(&vstats); err != nil {
		return stats, nil, fmt.Errorf("failed to decode reader data: %w", err)
	}
	if err := stats.AddInterfaceFromRawLine(vstats.ProcNetDev); err != nil {
		return stats, nil, fmt.Errorf("failed to add interface data %s: %w", vstats.ProcNetDev, err)
	}

	if loadAvg, err := LoadavgFromRawLine(vstats.ProcLoadAvg); err != nil {
		return stats, nil, fmt.Errorf("failed to read average load data %s: %w", vstats.ProcLoadAvg, err)
	} else {
		stats.Loadavg = loadAvg
	}

	stats.NotAvailable = vstats.NotAvailable
	inf := stats.Interfaces["eth0"]
	inf.Speed = vstats.InfSpeed
	stats.Interfaces["eth0"] = inf

	return stats, vstats.Stats, nil
}

func vstatsPrecompute(cacheName string, data todata.TOData, stats Statistics, miscStats map[string]interface{}) PrecomputedData {
	dsStats := make(map[string]*DSStat)
	var precomputed PrecomputedData
	precomputed.OutBytes = 0
	precomputed.MaxKbps = 0
	for _, iface := range stats.Interfaces {
		precomputed.OutBytes += iface.BytesOut
		kbps := iface.Speed * 1000
		if kbps > precomputed.MaxKbps {
			precomputed.MaxKbps = kbps
		}
	}

	for name, value := range miscStats {
		parts := strings.Split(name, ".")
		subsubdomain := parts[0]
		subdomain := parts[1]
		domain := strings.Join(parts[2:len(parts)-1], ".")

		ds, ok := data.DeliveryServiceRegexes.DeliveryService(domain, subdomain, subsubdomain)
		if !ok {
			precomputed.Errors = append(
				precomputed.Errors,
				fmt.Errorf("no Delivery Service match for '%s.%s.%s'", subsubdomain, subdomain, domain),
			)
			continue
		}
		if ds == "" {
			precomputed.Errors = append(
				precomputed.Errors,
				fmt.Errorf("empty Delivery Service fqdn '%s.%s.%s'", subsubdomain, subdomain, domain),
			)
			continue
		}

		dsName := string(ds)

		vstatsProcessCounter(dsStats, dsName, parts[len(parts)-1], value)
	}
	precomputed.DeliveryServiceStats = dsStats

	return precomputed
}

func vstatsProcessCounter(dsStats map[string]*DSStat, dsName, category string, value interface{}) error {
	if stat, ok := dsStats[dsName]; stat == nil || !ok {
		dsStats[dsName] = new(DSStat)
	}
	parsedValue, ok := value.(float64)
	if !ok {
		// only float counters are used now
		return fmt.Errorf("expected counter value of type float got type: %T, for value: %v", value, value)
	}
	statusCode, err := strconv.ParseInt(category, 10, 64)

	if err == nil {
		if statusCode >= 200 && statusCode < 300 {
			dsStats[dsName].Status2xx += uint64(parsedValue)
		} else if statusCode >= 300 && statusCode < 400 {
			dsStats[dsName].Status3xx += uint64(parsedValue)
		} else if statusCode >= 400 && statusCode < 500 {
			dsStats[dsName].Status4xx += uint64(parsedValue)
		} else if statusCode >= 500 && statusCode < 600 {
			dsStats[dsName].Status5xx += uint64(parsedValue)
		}
	}
	return nil
}
