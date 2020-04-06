// Package cache contains definitions for mechanisms used to extract health
// and statistics data from cache-server-provided data. The most commonly
// used format is the ``stats_over_http'' format provided by the plugin of the
// same name for Apache Traffic Server, followed closely by ``astats''  which
// is the legacy format used by older versions of Apache Traffic Control.
//
// Creating A New Stats Type
//
// To create a new Stats Type, for a custom caching proxy with its own stats
// format:
//
// 1. Create a file for your type in the traffic_monitor/cache directory and
//    package, `github.com/apache/trafficcontrol/traffic_monitor/cache/`
// 2. Create Parse and (optionally) Precompute functions in your file, with the
//     signature of `StatsTypeParser` and `StatsTypePrecomputer`
// 3. In your file, add
//    `func init(){AddStatsType(myTypeParser, myTypePrecomputer})`
//
// Your Parser should take the raw bytes from the `io.Reader` and populate the
// raw stats from them. For maximum compatibility, the names of these should be
// of the same form as Apache Traffic Server's `stats_over_http`, of the form
// "plugin.remap_stats.delivery-service-fqdn.com.in_bytes" et cetera. Traffic
// Control _may_ work with custom stat names, but we don't currently guarantee
// it.
//
// Your Precomputer should take the Stats and System information your Parser
// created, and populate the PrecomputedData. It is essential that all
// PrecomputedData fields are populated, especially `DeliveryServiceStats`,
// as they are used for cache and delivery service availability and threshold
// computation. If PrecomputedData is not properly and fully populated, the
// cache's availability will not be properly computed.
//
// Note the PrecomputedData `Reporting` and `Time` fields are the exception:
// they do not need to be set, and will be forcibly overridden by the Handler
// after your Precomputer function returns.
//
// Note these functions will not be called for Health polls, only Stat polls.
// Your Cache should have two separate stats endpoints: a small light endpoint
// returning only system stats and used to quickly verify reachability, and a
// large endpoint with all stats. If your cache does not have two stat
// endpoints, you may use your large stat endpoint for the Health poll, and
// configure the Health poll interval to be arbitrarily slow.
//
// Note your stats functions SHOULD NOT reuse functions from other stats types,
// even if they are similar, or have identical helper functions. This is a case
// where "duplicate" code is acceptable, because it's not conceptually
// duplicate. You don't want your stat parsers to break if the similar stats
// format you reuse code from changes.
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


func ExampleLoadavgFromRawLine() {
	loadavg, err := LoadavgFromRawLine("0.30 0.12 0.21 1/863 1421")
	fmt.Println(err)
	fmt.Printf("%.2f %.2f %.2f %d/%d %d", loadavg.One, loadavg.Five, loadavg.Fifteen, loadavg.CurrentProcesses, loadavg.TotalProcesses, loadavg.LatestPID)
	// Output: <nil>
	// 0.30 0.12 0.21 1/863 1421
}

func ExampleAddInterfaceFromRawLine() {
	var s Statistics
	raw := "eth0:47907832129 14601260    0    0    0     0          0   790726 728207677726 10210700052    0    0    0     0       0          0"

	if err := s.AddInterfaceFromRawLine(raw); err != nil {
		fmt.Println(err)
		return
	}

	iface, ok := s.Interfaces["eth0"]
	if !ok {
		fmt.Printf("Error, no 'eth0' interface!\n%+v", s.Interfaces)
		return
	}
	fmt.Printf("eth0: {BytesOut: %d, BytesIn: %d}", iface.BytesOut, iface.BytesIn)
	// Output: eth0: {BytesOut: 728207677726, BytesIn: 47907832129}
}
