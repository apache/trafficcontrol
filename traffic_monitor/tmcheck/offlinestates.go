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

// ValidateOfflineStates validates that no OFFLINE or ADMIN_DOWN caches in the given Traffic Ops' CRConfig are marked Available in the given Traffic Monitor's CRStates.
func ValidateOfflineStates(tmURI string, toClient *to.Session) error {
	cdn, err := GetCDN(tmURI)
	if err != nil {
		return fmt.Errorf("getting CDN from Traffic Monitor: %v", err)
	}
	return ValidateOfflineStatesWithCDN(tmURI, cdn, toClient)
}

// ValidateOfflineStatesWithCDN validates per ValidateOfflineStates, but saves an additional query if the Traffic Monitor's CDN is known.
func ValidateOfflineStatesWithCDN(tmURI string, tmCDN string, toClient *to.Session) error {
	response, _, err := toClient.GetCRConfig(tmCDN, to.RequestOptions{})
	if err != nil {
		return fmt.Errorf("getting CRConfig: %v", err)
	}

	crConfig := response.Response
	return ValidateOfflineStatesWithCRConfig(tmURI, &crConfig, toClient)
}

// ValidateOfflineStatesWithCRConfig validates per ValidateOfflineStates, but saves querying the CRconfig if it's already fetched.
func ValidateOfflineStatesWithCRConfig(tmURI string, crConfig *tc.CRConfig, toClient *to.Session) error {
	crStates, err := GetCRStates(tmURI + TrafficMonitorCRStatesPath)
	if err != nil {
		return fmt.Errorf("getting CRStates: %v", err)
	}

	return ValidateCRStates(crStates, crConfig)
}

// ValidateCRStates validates that no OFFLINE or ADMIN_DOWN caches in the given CRConfig are marked Available in the given CRStates.
func ValidateCRStates(crstates *tc.CRStates, crconfig *tc.CRConfig) error {
	for cacheName, cacheInfo := range crconfig.ContentServers {
		status := tc.CacheStatusFromString(string(*cacheInfo.ServerStatus))
		if status != tc.CacheStatusAdminDown && status != tc.CacheStatusOffline {
			continue
		}

		available, ok := crstates.Caches[tc.CacheName(cacheName)]
		if !ok {
			return fmt.Errorf("Cache %v in CRConfig but not CRStates", cacheName)
		}

		if available.IsAvailable {
			return fmt.Errorf("Cache %v is %v in CRConfig, but available in CRStates", cacheName, status)
		}

	}
	return nil
}

// CRStatesOfflineValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll.
func CRStatesOfflineValidator(
	tmURI string,
	toClient *to.Session,
	interval time.Duration,
	grace time.Duration,
	onErr func(error),
	onResumeSuccess func(),
	onCheck func(error),
) {
	Validator(tmURI, toClient, interval, grace, onErr, onResumeSuccess, onCheck, ValidateOfflineStates)
}

// AllMonitorsCRStatesOfflineValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll. Note the error passed to `onErr` may be a general validation error not associated with any monitor, in which case the passed `tc.TrafficMonitorName` will be empty.
func AllMonitorsCRStatesOfflineValidator(
	toClient *to.Session,
	interval time.Duration,
	includeOffline bool,
	grace time.Duration,
	onErr func(tc.TrafficMonitorName, error),
	onResumeSuccess func(tc.TrafficMonitorName),
	onCheck func(tc.TrafficMonitorName, error),
) {
	AllValidator(toClient, interval, includeOffline, grace, onErr, onResumeSuccess, onCheck, ValidateAllMonitorsOfflineStates)
}

// ValidateOfflineStates validates that no OFFLINE or ADMIN_DOWN caches in the given Traffic Ops' CRConfig are marked Available in the given Traffic Monitor's CRStates.
func ValidateAllMonitorsOfflineStates(toClient *to.Session, includeOffline bool) (map[tc.TrafficMonitorName]error, error) {
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
		errs[tc.TrafficMonitorName(*server.HostName)] = ValidateOfflineStatesWithCRConfig(uri, crConfig.CRConfig, toClient)
	}
	return errs, nil
}
