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

// noop is a no-op parser designed to work with the the no-op poller,
// to report caches as healthy without actually polling them.

import (
	"io"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

func init() {
	registerDecoder("noop", noOpParse, noopPrecompute)
}

func noOpParse(string, io.Reader, interface{}) (Statistics, map[string]interface{}, error) {
	stats := Statistics{
		Loadavg: Loadavg{
			One:              0.1,
			Five:             0.05,
			Fifteen:          0.05,
			CurrentProcesses: 1,
			TotalProcesses:   1000,
			LatestPID:        30000,
		},
		Interfaces: map[string]Interface{
			"bond0": Interface{
				Speed:    20000,
				BytesIn:  10000,
				BytesOut: 100000,
			},
		},
	}
	return stats, map[string]interface{}{}, nil
}

func noopPrecompute(cache string, toData todata.TOData, stats Statistics, rawStats map[string]interface{}) PrecomputedData {
	return PrecomputedData{DeliveryServiceStats: map[string]*DSStat{}}
}
