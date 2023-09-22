// Package cache contains definitions for mechanisms used to extract health
// and statistics data from cache-server-provided data. The most commonly
// used format is the “stats_over_http” format provided by the plugin of the
// same name for Apache Traffic Server, followed closely by “astats”  which
// is the legacy format used by older versions of Apache Traffic Control.
//
// # Creating A New Stats Type
//
// To create a new Stats Type, for a custom caching proxy with its own stats
// format:
//
//  1. Create a file for your type in the traffic_monitor/cache directory and
//     package, `github.com/apache/trafficcontrol/v8/traffic_monitor/cache/`
//  2. Create Parse and (optionally) Precompute functions in your file, with the
//     signature of `StatisticsParser` and `StatisticsPrecomputer`, respectively
//  3. In your file's special `init` func, call `registerDecoder` with your two
//     functions to register the new format. The name of the format MUST be
//     unique!
//  4. To apply the new parsing format to a cache server, set its Profile's
//     “health.polling.format“ Parameter's Value to the name of the desired
//     format.
//
// Your Parser should take the raw bytes from the `io.Reader` and populate the
// raw stats from them. It needs to provide (nearly) all of the data in a
// Statistics structure. Specifically, the available statistics MUST include:
//
//   - One-minute "loadavg" value for the cache server. The others are optional,
//     as we only use the one-minute value for health checks.
//   - At least one network interface (which will be considered the one used for
//     routing, and if multiple "monitored" network interfaces are configured for
//     the cache server in Traffic Ops they MUST all be present) and specifically
//     its name, “speed”, and bytes in and out. Parsers SHOULD return an error
//     if at least one interface cannot be found in the payload data.
//   - If your format does not directly indicate if the cache server is available
//     then NotAvailable should just be set to “false”.
//
// All other statistics (e.g. Delivery Service stats) should be returned in the
// map of statistic names to their values.
//
// Your Precomputer should take the Statistics and other miscellaneous stats
// that your Parser created, and populate the PrecomputedData. It is essential
// that all PrecomputedData fields are populated, especially
// `DeliveryServiceStats`, as they are used for cache and Delivery Service
// availability and threshold computation. If PrecomputedData is not properly
// and fully populated, the cache's availability will not be properly computed.
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
// configure the Health poll interval to be arbitrarily slow. These are
// controlled by the “health.polling.url' Parameter in Traffic Ops.
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

import (
	"fmt"
	"io"

	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

const DefaultStatsType = "astats"

// StatsDecoder is a pair of functions registered for decoding statistics
// of a particular format, and parsing that data and precomputing related
// data, respectively.
type StatsDecoder struct {
	Parse      StatisticsParser
	Precompute StatisticsPrecomputer
}

// StatisticsParser is a function that parses raw input data for a given
// statistics format and returns the meaningful statistics.
// In addition to the decoded statistics, the decoder should also return
// whatever miscellaneous data was in the payload but not represented by
// the properties of a Statistics object, so that it can be used in later
// calculations if necessary.
type StatisticsParser func(string, io.Reader, interface{}) (Statistics, map[string]interface{}, error)

// StatisticsPrecomputer is a function that "pre-computes" some statistics
// beyond the basic ones covered by a Statistics object.
// Precomputers aren't called until a statistics poll is done, whereas basic
// Statistics are calculated even for Health polls.
type StatisticsPrecomputer func(string, todata.TOData, Statistics, map[string]interface{}) PrecomputedData

var statDecoders = map[string]StatsDecoder{}

// GetDecoder gets a decoder for the given statistic format. Returns an error
// if no parser for the given format exists.
func GetDecoder(format string) (StatsDecoder, error) {
	if decoder, ok := statDecoders[format]; ok {
		return decoder, nil
	}
	return StatsDecoder{}, fmt.Errorf("No decoder registered for format '%s'", format)
}

func registerDecoder(format string, parser StatisticsParser, precomputer StatisticsPrecomputer) {
	statDecoders[format] = StatsDecoder{Parse: parser, Precompute: precomputer}
}
