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

package tmcheck

import (
	"fmt"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

const QueryIntervalMax = time.Duration(10) * time.Second

// ValidateQueryInterval validates the given monitor has an acceptable Query Interval 95th percentile.
func ValidateQueryInterval(tmURI string, toClient *to.Session) error {
	stats, err := GetStats(tmURI + TrafficMonitorStatsPath)
	if err != nil {
		return fmt.Errorf("getting Stats: %v", err)
	}
	queryInterval := time.Duration(stats.QueryInterval95thPercentile) * time.Millisecond

	if queryInterval > QueryIntervalMax {
		return fmt.Errorf("Query Interval 95th Percentile %v greater than max %v", queryInterval, QueryIntervalMax)
	}
	return nil
}

// QueryIntervalValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll.
func QueryIntervalValidator(
	tmURI string,
	toClient *to.Session,
	interval time.Duration,
	grace time.Duration,
	onErr func(error),
	onResumeSuccess func(),
	onCheck func(error),
) {
	Validator(tmURI, toClient, interval, grace, onErr, onResumeSuccess, onCheck, ValidateQueryInterval)
}

// AllMonitorsQueryIntervalValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll. Note the error passed to `onErr` may be a general validation error not associated with any monitor, in which case the passed `tc.TrafficMonitorName` will be empty.
func AllMonitorsQueryIntervalValidator(
	toClient *to.Session,
	interval time.Duration,
	includeOffline bool,
	grace time.Duration,
	onErr func(tc.TrafficMonitorName, error),
	onResumeSuccess func(tc.TrafficMonitorName),
	onCheck func(tc.TrafficMonitorName, error),
) {
	AllValidator(toClient, interval, includeOffline, grace, onErr, onResumeSuccess, onCheck, ValidateAllMonitorsQueryInterval)
}

// ValidateAllMonitorsQueryInterval validates, for all monitors in the given Traffic Ops, an acceptable query interval 95th percentile.
func ValidateAllMonitorsQueryInterval(toClient *to.Session, includeOffline bool) (map[tc.TrafficMonitorName]error, error) {
	servers, err := GetMonitors(toClient, includeOffline)
	if err != nil {
		return nil, err
	}

	errs := map[tc.TrafficMonitorName]error{}
	for _, server := range servers {
		uri := fmt.Sprintf("http://%s.%s", *server.HostName, *server.DomainName)
		errs[tc.TrafficMonitorName(*server.HostName)] = ValidateQueryInterval(uri, toClient)
	}
	return errs, nil
}
