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
	"github.com/apache/trafficcontrol/lib/go-tc/enum"
	"io"

	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

//
// To create a new Stats Type, for a custom caching proxy with its own stats format:
//
// 1. Create a file for your type in this directory and package, `traffic_monitor/cache/`
// 2. Create Parse and Precompute functions in your file, with the signature of `StatsTypeParser` and `StatsTypePrecomputer`
// 3. In your file, add `func init(){AddStatsType(myTypeParser, myTypePrecomputer})`
//
// Your Parser should take the raw bytes from the `io.Reader` and populate the raw stats from them. For maximum compatibility, the names of these should be of the same form as Apache Traffic Server's `stats_over_http`, of the form "plugin.remap_stats.delivery-service-fqdn.com.in_bytes" et cetera. Traffic Control _may_ work with custom stat names, but we don't currently guarantee it.
//
// Your Precomputer should take the Stats and System information your Parser created, and populate the PrecomputedData. It is essential that all PrecomputedData fields are populated, especially `DeliveryServiceStats`, as they are used for cache and delivery service availability and threshold computation. If PrecomputedData is not properly and fully populated, the cache's availability will not be properly computed.
//
// Note this function is not called for Health polls, only Stat polls. Your Cache should have two separate stats endpoints: a small light endpoint returning only system stats and used to quickly verify reachability, and a large endpoint with all stats. If your cache does not have two stat endpoints, you may use your large stat endpoint for the Health poll, and configure the Health poll interval to be arbitrarily slow.
//
// Note the PrecomputedData `Reporting` and `Time` fields are the exception: they do not need to be set, and will be forcibly overridden by the Handler after your Precomputer function returns.
//
// Note your stats functions SHOULD NOT reuse functions from other stats types, even if they are similar, or have identical helper functions. This is a case where "duplicate" code is acceptable, because it's not conceptually duplicate. You don't want your stat parsers to break if the similar stats format you reuse code from changes.
//

const DefaultStatsType = "astats"

// CacheStatsTypeDecoder is a pair of functions registered for decoding a particular Stats type, for parsing stats, and creating precomputed data
type StatsTypeDecoder struct {
	Parse      StatsTypeParser
	Precompute StatsTypePrecomputer
}

// StatsTypeParser takes the bytes returned from the cache's stats endpoint, along with the cache name, and returns the map of raw stats (whose names must be strings, and values may be any primitive type but MUST be float64 if they are used by a Parameter Threshold) and System information.
type StatsTypeParser func(cache enum.CacheName, r io.Reader) (error, map[string]interface{}, AstatsSystem)

// StatsTypePrecomputer takes the cache name, the time the given stats were received, the Traffic Ops data, and the raw stats and system information created by Parse, and returns the PrecomputedData. Note this will only be called for Stats polls, not Health polls. Note errors should be returned in PrecomputedData.Errors
//
type StatsTypePrecomputer func(cache enum.CacheName, toData todata.TOData, stats map[string]interface{}, system AstatsSystem) PrecomputedData

// StatsTypeDecoders holds the functions for parsing cache stats. This is not const, because Go doesn't allow constant maps. This is populated on startup, and MUST NOT be modified after startup.
var StatsTypeDecoders = map[string]StatsTypeDecoder{}

func AddStatsType(typeName string, parser StatsTypeParser, precomputer StatsTypePrecomputer) {
	StatsTypeDecoders[typeName] = StatsTypeDecoder{Parse: parser, Precompute: precomputer}
}
