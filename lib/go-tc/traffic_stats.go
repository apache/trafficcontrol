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

// TRAFFIC_STATS_VERSION was supposed to be the "API version", but actually the plugin (this route
// used to be a plugin in Perl) always returned this static value
const TRAFFIC_STATS_VERSION = "1.2"

// TRAFFIC_STATS_SOURCE is the value of the "source" field in an API response. Perl always returned
// source="TrafficStats", so we do too
const TRAFFIC_STATS_SOURCE = "TrafficStats"

// TrafficStatsDuration reflects all the possible durations that can be requested via the
// deliveryservice_stats endpoint
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

// TrafficStatsDurationFromString converts the given string into the appropriate
// TrafficStatsDuration if possible - returns the InvalidDuration constant if not.
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

// Seconds returns the number of seconds to which a TrafficStatsDuration is equivalent.
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

// TrafficStatsResponse represents a response from one of the "Traffic Stats endpoints" of the
// Traffic Ops API, e.g. `/deliveryservice_stats`.
type TrafficStatsResponse struct {
	// Series holds the actual data - it is NOT in general the same as a github.com/influxdata/influxdb1-client/models.Row
	Series *TrafficStatsSeries `json:"series,omitempty"`
	// Source has an unknown purpose. I believe this is supposed to name the "plugin" that provided
	// the data - kept for compatibility with the Perl version(s) of the "Traffic Stats endpoints".
	// Deprecated: this'll be removed or reworked to make more sense in the future
	Source string `json:"source"`
	// Summary contains summary statistics of the data in Series
	Summary *TrafficStatsSummary `json:"summary,omitempty"`
	// Version is supposed to represent the API version - but actually the API just reports a static
	// number (TRAFFIC_STATS_VERSION).
	// Deprecated: this'll be removed or reworked to make more sense in the future
	Version string `json:"version"`
}

// TrafficStatsSummary contains summary statistics for a data series.
type TrafficStatsSummary struct {
	// Average is calculated as an arithmetic mean
	Average float64 `json:"average"`
	// Count is the total number of data points _except_ for any values that would appear as 'nil'
	// in the corresponding series
	Count                  uint    `json:"count"`
	FifthPercentile        float64 `json:"fifthPercentile"`
	Max                    float64 `json:"max"`
	Min                    float64 `json:"min"`
	NinetyEighthPercentile float64 `json:"ninetyEighthPercentile"`
	NinetyFifthPercentile  float64 `json:"ninetyFifthPercentile"`
	// TotalBytes is the total number of bytes served when the "metric type" requested is "kbps"
	// (or actually just contains "kbps"). If this is not nil, TotalTransactions *should* always be
	// nil.
	TotalBytes *float64 `json:"totalBytes"`
	// Totaltransactions is the total number of transactions within the requested window. Whenever
	// the requested metric doesn't contain "kbps", it assumed to be some kind of transactions
	// measurement. In that case, this will not be nil - otherwise it will be nil. If this not nil,
	// TotalBytes *should* always be nil.
	TotalTransactions *float64 `json:"totalTransactions"`
}

// TrafficStatsSeries is the actual data returned by a request to a "Traffic Stats endpoint".
type TrafficStatsSeries struct {
	// Columns is a list of column names. Each "row" in Values is ordered to match up with these
	// column names.
	Columns []string `json:"columns"`
	// Count is the total number of returned data points. Should be the same as len(Values)
	Count uint `json:"count"`
	// Name is the name of the InfluxDB database from which the data was retrieved
	Name string `json:"name"`
	// Tags is a set of InfluxDB tags associated with the requested database.
	Tags map[string]string `json:"tags"`
	// Values is an array of rows of arbitrary data that can only really be interpreted by
	// inspecting Columns, in general. In practice, however, each element is nearly always a
	// slice where the first element is an RFC3339 timestamp (as a string) and the second/final
	// element is a floating point number (or nil) indicating the value at that time.
	Values [][]interface{} `json:"values"`
}
