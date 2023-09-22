package trafficstats

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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	influx "github.com/influxdata/influxdb/client/v2"
)

func getSummary(db string, q influx.Query, client *influx.Client) (*tc.TrafficStatsSummary, error) {
	s := tc.TrafficStatsSummary{}
	msgs := []influx.Message{}
	log.Debugf("InfluxDB SummaryQuery: %+v", q)

	defer log.Debugf("Messages from summary query: %s", tc.MessagesToString(msgs))

	resp, err := (*client).Query(q)
	if err != nil {
		return nil, err
	}

	if resp.Results != nil && len(resp.Results) == 1 {
		r := resp.Results[0]
		if r.Messages != nil {
			for _, m := range r.Messages {
				if m != nil {
					msgs = append(msgs, *m)
				}
			}
		}

		// Perl API handled no series as a non error
		// Matching implementation
		if len(r.Series) == 0 {
			return nil, nil
		}

		if len(r.Series) != 1 {
			return nil, fmt.Errorf("Improper number of series: %d", len(r.Series))
		}

		series := r.Series[0]

		if len(series.Values) != 1 {
			return nil, fmt.Errorf("Improper number of returned rows: %d", len(r.Series[0].Values))
		}

		vals := series.Values[0]
		if len(vals) != 8 || len(series.Columns) != 8 {
			return nil, fmt.Errorf("Improper number of returned values in row: %d (%d cols)", len(vals), len(series.Columns))
		}

		mappedValues := map[string]interface{}{}
		for i, v := range vals {
			mappedValues[series.Columns[i]] = v
		}

		var err error
		if s.Average, err = extractFloat64("average", mappedValues); err != nil {
			return nil, err
		}

		if s.Count, err = extractUInt("count", mappedValues); err != nil {
			return nil, err
		}

		if s.FifthPercentile, err = extractFloat64("fifthPercentile", mappedValues); err != nil {
			return nil, err
		}

		if s.Max, err = extractFloat64("max", mappedValues); err != nil {
			return nil, err
		}

		if s.Min, err = extractFloat64("min", mappedValues); err != nil {
			return nil, err
		}

		if s.Max, err = extractFloat64("max", mappedValues); err != nil {
			return nil, err
		}

		if s.NinetyEighthPercentile, err = extractFloat64("ninetyEighthPercentile", mappedValues); err != nil {
			return nil, err
		}

		if s.NinetyFifthPercentile, err = extractFloat64("ninetyFifthPercentile", mappedValues); err != nil {
			return nil, err
		}

	} else {
		log.Debugf("InfluxDB summary response: %+v", resp)
		return nil, errors.New("'results' missing or improper")
	}

	if resp.Error() != nil {
		log.Debugf("response error, summary object was: %+v", s)
		return nil, resp.Error()
	}

	return &s, nil
}
