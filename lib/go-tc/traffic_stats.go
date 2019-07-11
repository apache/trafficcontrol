package tc

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

import "strings"

// This was supposed to be the "API version", but actually the plugin - this route used to be a
// plugin in Perl - always returned this static value
const TRAFFIC_STATS_VERSION = "1.2"

// Perl always returns source="TrafficStats", so we do too
const TRAFFIC_STATS_SOURCE = "TrafficStats"

// This type reflects all the possible durations that can be requested via the deliveryservice_stats endpoint
type TrafficStatsDuration string

const (
	InvalidDuration TrafficStatsDuration = ""
	OneMinute       TrafficStatsDuration = "1m"
	FiveMinutes     TrafficStatsDuration = "5m"
	OneHour         TrafficStatsDuration = "1h"
	SixHours        TrafficStatsDuration = "6h"
	OneDay          TrafficStatsDuration = "1d"
	OneWeek         TrafficStatsDuration = "1w"
	OneMonth        TrafficStatsDuration = "4w"
)

// Converts the given string into the appropriate TrafficStatsDuration if possible - returns the
// InvalidDuration constant if not.
func TrafficStatsDurationFromString(v string) TrafficStatsDuration {
	switch TrafficStatsDuration(strings.Trim(v, " ")) {
	case OneMinute:
		return OneMinute
	case FiveMinutes:
		return FiveMinutes
	case OneHour:
		return OneHour
	case SixHours:
		return SixHours
	case OneDay:
		return OneDay
	case OneWeek:
		return OneWeek
	case OneMonth:
		return OneMonth
	}
	return InvalidDuration
}

// For valid TrafficStatsDurations this returns the number of seconds to which it is equivalent.
// For invalid objects, it returns -1 - otherwise it will always, of course be > 0.
func (d TrafficStatsDuration) Seconds() int64 {
	switch d {
	case OneMinute:
		return 60
	case FiveMinutes:
		return 300
	case OneHour:
		return 3600
	case SixHours:
		return 21600
	case OneDay:
		return 86400
	case OneWeek:
		return 604800
	case OneMonth:
		return 2419200
	}

	return -1
}

// Represents a response from one of the "Traffic Stats endpoints" of the Traffic Ops API, e.g.
// `/deliveryservice_stats`.
type TrafficStatsResponse struct {
	// This holds the actual data - it is NOT in general the same as a github.com/influxdata/influxdb1-client/models.Row
	Series *TrafficStatsSeries `json:"series,omitempty"`
	// I believe this is supposed to name the "plugin" that provided the data - kept for compatibility
	// with the Perl version(s) of the "Traffic Stats endpoints".
	// Deprecated: this'll be removed or reworked to make more sense in the future
	Source string `json:"source"`
	// Contains summary statistics of the data in Series
	Summary *TrafficStatsSummary `json:"summary,omitempty"`
	// This is supposed to represent the API version - but actually the API just reports a static
	// number (TRAFFIC_STATS_VERSION).
	// Deprecated: this'll be removed or reworked to make more sense in the future
	Version string `json:"version"`
}

// Contains summary statistics for a data series.
type TrafficStatsSummary struct {
	// Calculated as an arithmetic mean
	Average float64 `json:"average"`
	// The total number of data points _except_ for any values that would appear as 'nil' in the
	// corresponding series
	Count                  uint    `json:"count"`
	FifthPercentile        float64 `json:"fifthPercentile"`
	Max                    float64 `json:"max"`
	Min                    float64 `json:"min"`
	NinetyEighthPercentile float64 `json:"ninetyEighthPercentile"`
	NinetyFifthPercentile  float64 `json:"ninetyFifthPercentile"`
	// This is the total number of bytes served when the "metric type" requested is "kbps" (or actually)
	// just contains "kbps"). If this is not nil, TotalTransactions *should* always be nil.
	TotalBytes *float64 `json:"totalBytes"`
	// Whenever the requested metric doesn't contain "kbps", it is assumed to be some kind of
	// transactions measurement. In that case, this indicates the total number of transactions
	// within the requested window.
	TotalTransactions *float64 `json:"totalTransactions"`
}

// This is the actual data returned by a request to a "Traffic Stats endpoint". Note that this
// section shouldn't appear if the "interval" query parameter was not specified.
type TrafficStatsSeries struct {
	// This is a list of column names. Each "row" in Values is ordered to match up with these column
	// names.
	Columns []string `json:"columns"`
	// The total number of returned data points. Should be the same as len(Values)
	Count uint `json:"count"`
	// The name of the InfluxDB database from which the data was retrieved
	Name string `json:"name"`
	// A set of InfluxDB tags associated with the requested database.
	Tags map[string]string `json:"tags"`
	// In general, each element of Values is a row of arbitrary data that can only really be
	// interpreted by inspecting Columns. In practice, however, each element is nearly always a
	// slice where the first element is an RFC3339 timestamp (as a string) and the second/final
	// element is a floating point number (or nil) indicating the value at that time.
	Values [][]interface{} `json:"values"`
}
