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
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

// ValidateDSStates validates that all Delivery Services in the CRConfig exist in given Traffic Monitor's DSStats.
// Existence in DSStats is useful to verify, because "Available: false" in CRStates
func ValidateDSStats(tmURI string, toClient *to.Session) error {
	cdn, err := GetCDN(tmURI)
	if err != nil {
		return fmt.Errorf("getting CDN from Traffic Monitor: %v", err)
	}
	return ValidateDSStatsWithCDN(tmURI, cdn, toClient)
}

// ValidateOfflineStatesWithCDN validates per ValidateOfflineStates, but saves an additional query if the Traffic Monitor's CDN is known.
func ValidateDSStatsWithCDN(tmURI string, tmCDN string, toClient *to.Session) error {
	response, _, err := toClient.GetCRConfig(tmCDN, to.RequestOptions{})
	if err != nil {
		return fmt.Errorf("getting CRConfig: %v", err)
	}

	crConfig := response.Response
	return ValidateDSStatsWithCRConfig(tmURI, &crConfig, toClient)
}

// ValidateOfflineStatesWithCRConfig validates per ValidateOfflineStates, but saves querying the CRconfig if it's already fetched.
func ValidateDSStatsWithCRConfig(tmURI string, crConfig *tc.CRConfig, toClient *to.Session) error {
	dsStats, err := GetDSStats(tmURI + TrafficMonitorDSStatsPath)
	if err != nil {
		return fmt.Errorf("getting DSStats: %v", err)
	}

	return ValidateDSStatsData(dsStats, crConfig)
}

func hasCaches(dsName string, crconfig *tc.CRConfig) bool {
	for _, server := range crconfig.ContentServers {
		if _, ok := server.DeliveryServices[dsName]; ok {
			return true
		}
	}
	return false
}

// ValidateDSStatsData validates that all delivery services in the given CRConfig with caches assigned exist in the given DSStats.
func ValidateDSStatsData(dsStats *dsdata.StatsOld, crconfig *tc.CRConfig) error {
	for dsName, _ := range crconfig.DeliveryServices {
		if !hasCaches(dsName, crconfig) {
			continue
		}
		if _, ok := dsStats.DeliveryService[tc.DeliveryServiceName(dsName)]; !ok {
			return fmt.Errorf("Delivery Service %v in CRConfig but not DSStats", dsName)
		}
	}
	return nil
}

// DSStatsValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll.
func DSStatsValidator(
	tmURI string,
	toClient *to.Session,
	interval time.Duration,
	grace time.Duration,
	onErr func(error),
	onResumeSuccess func(),
	onCheck func(error),
) {
	Validator(tmURI, toClient, interval, grace, onErr, onResumeSuccess, onCheck, ValidateDSStats)
}

// AllMonitorsDSStatsValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll. Note the error passed to `onErr` may be a general validation error not associated with any monitor, in which case the passed `tc.TrafficMonitorName` will be empty.
func AllMonitorsDSStatsValidator(
	toClient *to.Session,
	interval time.Duration,
	includeOffline bool,
	grace time.Duration,
	onErr func(tc.TrafficMonitorName, error),
	onResumeSuccess func(tc.TrafficMonitorName),
	onCheck func(tc.TrafficMonitorName, error),
) {
	AllValidator(toClient, interval, includeOffline, grace, onErr, onResumeSuccess, onCheck, ValidateAllMonitorsDSStats)
}

// ValidateAllMonitorDSStats validates, for all monitors in the given Traffic Ops, DSStats contains all Delivery Services in the CRConfig.
func ValidateAllMonitorsDSStats(toClient *to.Session, includeOffline bool) (map[tc.TrafficMonitorName]error, error) {
	servers, err := GetMonitors(toClient, includeOffline)
	if err != nil {
		return nil, err
	}

	crConfigs := GetCRConfigs(GetCDNs(servers), toClient)

	errs := map[tc.TrafficMonitorName]error{}
	for _, server := range servers {
		crConfig := crConfigs[tc.CDNName(*server.CDNName)]
		if err := crConfig.Err; err != nil {
			errs[tc.TrafficMonitorName(*server.HostName)] = fmt.Errorf("getting CRConfig: %v", err)
			continue
		}

		uri := fmt.Sprintf("http://%s.%s", *server.HostName, *server.DomainName)
		errs[tc.TrafficMonitorName(*server.HostName)] = ValidateDSStatsWithCRConfig(uri, crConfig.CRConfig, toClient)
	}
	return errs, nil
}
