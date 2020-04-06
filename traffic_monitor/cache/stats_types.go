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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
)

const DefaultStatsType = "astats"

// CacheStatsTypeDecoder is a pair of functions registered for decoding a
// particular Stats type, for parsing stats, and creating precomputed data
type StatsTypeDecoder struct {
	Parse      StatsTypeParser
	Precompute StatsTypePrecomputer
}

// StatsDecoder is a pair of functions registered for decoding statistics
// of a particular format, and parsing that data and precomputing related
// data, respectively.
type StatsDecoder struct {
	Parse StatisticsParser
	Precompute StatisticsPrecomputer
}

// StatisticsParser is a function that parses raw input data for a given
// statistics format and returns the meaningful statistics.
// In addition to the decoded statistics, the decoder should also return
// whatever miscellaneous data was in the payload but not represented by
// the properties of a Statistics object, so that it can be used in later
// calculations if necessary.
type StatisticsParser func(string, io.Reader) (Statistics, map[string]interface{}, error)

// StatisticsPrecomputer is a function that "pre-computes" some statistics
// beyond the basic ones covered by a Statistics object.
// Precomputers aren't called until a statistics poll is done, whereas basic
// Statistics are calculated even for Health polls.
type StatisticsPrecomputer func(string, todata.TOData, Statistics, map[string]interface{}) PrecomputedData

// StatsTypeParser takes the bytes returned from the cache's stats endpoint,
// along with the cache name, and returns the map of raw stats (whose names
// must be strings, and values may be any primitive type but MUST be float64
// if they are used by a Parameter Threshold) and System information.
type StatsTypeParser func(cache tc.CacheName, r io.Reader) (error, map[string]interface{}, AstatsSystem)

// StatsTypePrecomputer takes the cache name, the time the given stats were
// received, the Traffic Ops data, and the raw stats and system information
//created by Parse, and returns the PrecomputedData. Note this will only be
// called for Stats polls, not Health polls. Note errors should be returned
// in PrecomputedData.Errors
type StatsTypePrecomputer func(cache tc.CacheName, toData todata.TOData, stats map[string]interface{}, system AstatsSystem) PrecomputedData

// StatsTypeDecoders holds the functions for parsing cache stats. This is
// not const, because Go doesn't allow constant maps. This is populated
// on startup, and MUST NOT be modified after startup.
var StatsTypeDecoders = map[string]StatsTypeDecoder{}

func AddStatsType(typeName string, parser StatsTypeParser, precomputer StatsTypePrecomputer) {
	StatsTypeDecoders[typeName] = StatsTypeDecoder{Parse: parser, Precompute: precomputer}
}

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
