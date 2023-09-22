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

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"

	influx "github.com/influxdata/influxdb/client/v2"
)

// TrafficStatsDurationPattern reflects all the possible durations that can be requested via the
// /deliveryservice_stats endpoint.
var TrafficStatsDurationPattern = regexp.MustCompile(`^\d+[mhdw]$`)

// DurationLiteralToSeconds returns the number of seconds to which an InfluxQL duration literal is
// equivalent. For invalid objects, it returns -1 - otherwise it will always, of course be > 0.
func DurationLiteralToSeconds(d string) (int64, error) {
	if strings.HasSuffix(d, "m") {
		v, err := strconv.ParseInt(strings.Split(d, "m")[0], 10, 64)
		return v * 60, err
	}

	if strings.HasSuffix(d, "h") {
		v, err := strconv.ParseInt(strings.Split(d, "h")[0], 10, 64)
		return v * 3600, err
	}

	if strings.HasSuffix(d, "d") {
		v, err := strconv.ParseInt(strings.Split(d, "d")[0], 10, 64)
		return v * 86400, err
	}

	if strings.HasSuffix(d, "w") {
		v, err := strconv.ParseInt(strings.Split(d, "w")[0], 10, 64)
		return v * 604800, err
	}

	return -1, errors.New("Invalid duration literal, no recognized suffix")
}

// TrafficStatsOrderable encodes what columns by which the data returned from a Traffic Stats query
// may be ordered.
type TrafficStatsOrderable string

// TimeOrder indicates an ordering by time at which the measurement was taken.
const TimeOrder TrafficStatsOrderable = "time"

// OrderableFromString parses the passed string and returns the corresponding value as a pointer to
// a TrafficStatsOrderable - or nil if the value was invalid.
func OrderableFromString(v string) *TrafficStatsOrderable {
	var o TrafficStatsOrderable
	switch v {
	case "time":
		o = TimeOrder
	default:
		return nil
	}
	return &o
}

// TrafficStatsExclude encodes what parts of a response to a request to a "Traffic Stats" endpoint
// of the TO API may be omitted.
type TrafficStatsExclude string

const (
	// ExcludeSeries can be used to omit the data series from a response.
	ExcludeSeries TrafficStatsExclude = "series"

	// ExcludeSummary can be used to omit the summary series from a response.
	ExcludeSummary TrafficStatsExclude = "summary"

	// ExcludeInvalid can be used if the the key that you want to exclude fails
	// validation.
	ExcludeInvalid TrafficStatsExclude = "INVALID"
)

// ExcludeFromString parses the passed string and returns the corresponding value as a TrafficStatsExclude.
func ExcludeFromString(v string) TrafficStatsExclude {
	switch v {
	case "series":
		return ExcludeSeries
	case "summary":
		return ExcludeSummary
	default:
		return ExcludeInvalid
	}
}

// TrafficStatsConfig represents the configuration of a request made to Traffic Stats. This is
// typically constructed by parsing a request body submitted to Traffic Ops.
type TrafficStatsConfig struct {
	End            time.Time
	ExcludeSeries  bool
	ExcludeSummary bool
	Interval       string
	Limit          *uint64
	MetricType     string
	Offset         *uint64
	OrderBy        *TrafficStatsOrderable
	Start          time.Time
	Unix           bool
}

// TrafficDSStatsConfig represents the configuration of a request made to
// Traffic Stats for Delivery Service stats.
type TrafficDSStatsConfig struct {
	DeliveryService string
	TrafficStatsConfig
}

// TrafficCacheStatsConfig represents the configuration of a request made to
// Traffic Stats for cache server stats.
type TrafficCacheStatsConfig struct {
	CDN string
	TrafficStatsConfig
}

// OffsetString is a stupid, dirty hack to try to convince Influx to not
// give back data that's outside of the range in a WHERE clause. It doesn't
// work, but it helps.
// (https://github.com/influxdata/influxdb/issues/8010)
func (c *TrafficStatsConfig) OffsetString() string {
	iSecs, err := DurationLiteralToSeconds(c.Interval)
	if err != nil {
		// TODO: Don't log in library function - return error if there is one.
		log.Errorf("Parsing duration literal: %v", err)
		return "0s"
	}
	return fmt.Sprintf("%ds", int64(c.Start.Sub(time.Unix(0, 0))/time.Second)%iSecs)
}

// TrafficDSStatsResponse represents a response from the
// /deliveryservice_stats "Traffic Stats" endpoints.
type TrafficDSStatsResponse struct {
	// Series holds the actual data - it is NOT in general the same as a github.com/influxdata/influxdb1-client/models.Row
	Series *TrafficStatsSeries `json:"series,omitempty"`
	// Summary contains summary statistics of the data in Series
	Summary *TrafficDSStatsSummary `json:"summary,omitempty"`
}

// TrafficStatsResponse represents the generic response from one of the "Traffic Stats endpoints" of the
// Traffic Ops API, e.g. `/cache_stats`.
type TrafficStatsResponse struct {
	// Series holds the actual data - it is NOT in general the same as a github.com/influxdata/influxdb1-client/models.Row
	Series *TrafficStatsSeries `json:"series,omitempty"`
	// Summary contains summary statistics of the data in Series
	Summary *TrafficStatsSummary `json:"summary,omitempty"`
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
}

// LegacyTrafficDSStatsSummary contains summary statistics for a data series
// for Delivery Service stats.
//
// Deprecated: This uses a confusing, incorrect name for the total number of
// kilobytes served, so new code should use TrafficDSStatsSummary instead.
type LegacyTrafficDSStatsSummary struct {
	TrafficStatsSummary
	// TotalBytes is the total number of kilobytes served when the "metric type" requested is "kbps"
	// (or actually just contains "kbps"). If this is not nil, TotalTransactions *should* always be
	// nil.
	TotalBytes *float64 `json:"totalBytes"`
	// Totaltransactions is the total number of transactions within the requested window. Whenever
	// the requested metric doesn't contain "kbps", it assumed to be some kind of transactions
	// measurement. In that case, this will not be nil - otherwise it will be nil. If this not nil,
	// TotalBytes *should* always be nil.
	TotalTransactions *float64 `json:"totalTransactions"`
}

// TrafficDSStatsSummary contains summary statistics for a data series for deliveryservice stats.
type TrafficDSStatsSummary struct {
	TrafficStatsSummary
	// TotalKiloBytes is the total number of kilobytes served when the "metric type" requested is "kbps"
	// (or actually just contains "kbps"). If this is not nil, TotalTransactions *should* always be
	// nil.
	TotalKiloBytes *float64 `json:"totalKiloBytes"`
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

// FormatTimestamps formats the timestamps contained in the Values array as RFC3339 strings.
// This returns an error if the data is not in the expected format.
func (s TrafficStatsSeries) FormatTimestamps() error {
	for i, v := range s.Values {
		if len(v) != 2 {
			return fmt.Errorf("Datapoint %d (%v) malformed", i, v)
		}

		switch v[0].(type) {
		case int64:
			s.Values[i][0] = time.Unix(0, v[0].(int64)).Format(time.RFC3339)
		case float64:
			s.Values[i][0] = time.Unix(0, int64(v[0].(float64))).Format(time.RFC3339)
		case json.Number:
			val, err := v[0].(json.Number).Int64()
			if err != nil {
				return fmt.Errorf("Datapoint %d (%v) malformed: %v", i, v, err)
			}
			s.Values[i][0] = time.Unix(0, val).Format(time.RFC3339)
		case string:
			t, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", v[0]))
			if err != nil {
				return fmt.Errorf("error converting timestamp %s into RFC3339 format: %v", v[0], err)
			}
			s.Values[i][0] = t.Format(time.RFC3339)
		default:
			return fmt.Errorf("Invalid type %T for datapoint %d (%v)", v[0], i, v)
		}
	}
	return nil
}

// MessagesToString converts a set of messages from an InfluxDB node into a
// single, print-able string.
func MessagesToString(msgs []influx.Message) string {
	if msgs == nil || len(msgs) == 0 {
		return ""
	}

	b := strings.Builder{}
	b.Write([]byte("Messages: ["))
	for _, m := range msgs {
		b.WriteString(m.Level)
		b.WriteRune(':')
		b.WriteString(m.Text)
		b.Write([]byte(", "))
	}
	b.WriteRune(']')
	return b.String()
}

// TrafficStatsCDNStats contains summary statistics for a given CDN.
type TrafficStatsCDNStats struct {
	Bandwidth    *float64 `json:"bandwidth"`
	Capacity     *float64 `json:"capacity"`
	CDN          string   `json:"cdn"`
	Connnections *float64 `json:"connections"`
}

// TrafficStatsTotalStats contains summary statistics across CDNs.
//
// This differs from TrafficStatsCDNStats in that it omits Capacity.
type TrafficStatsTotalStats struct {
	Bandwidth    *float64 `json:"bandwidth"`
	CDN          string   `json:"cdn"`
	Connnections *float64 `json:"connections"`
}

// TrafficStatsCDNStatsResponse contains response for getting current stats.
type TrafficStatsCDNStatsResponse struct {
	Response []TrafficStatsCDNsStats `json:"response"`
	Alerts
}

// CurrentStatsResponse is the type of a response from Traffic Ops
// to a request to its /current_stats endpoint.
type CurrentStatsResponse = TrafficStatsCDNStatsResponse

// TrafficStatsCDNsStats contains a list of CDN summary statistics.
type TrafficStatsCDNsStats struct {
	Stats []TrafficStatsCDNStats `json:"currentStats"`
}
