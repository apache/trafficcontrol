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

func getSeries(db string, q influx.Query, client *influx.Client) (*tc.TrafficStatsSeries, error) {
	msgs := []influx.Message{}
	s := tc.TrafficStatsSeries{}

	defer log.Debugf("Messages from summary query: %s", tc.MessagesToString(msgs))
	log.Debugf("InfluxDB series query: %+v", q)

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

		s = tc.TrafficStatsSeries{
			Name:    series.Name,
			Tags:    series.Tags,
			Values:  series.Values,
			Columns: series.Columns,
			Count:   uint(len(series.Values)),
		}

	} else {
		log.Debugf("InfluxDB series response: %+v", resp)
		return nil, errors.New("'results' missing or improper!")
	}

	if resp.Error() != nil {
		log.Debugf("response error, series object was %+v", s)
		return nil, resp.Error()
	}

	return &s, nil
}
