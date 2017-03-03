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

// package tmcheck contains utility functions for validating a Traffic Monitor is acting correctly.
package tmcheck

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/crconfig"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

const RequestTimeout = time.Second * time.Duration(30)

const TrafficMonitorCRStatesPath = "/publish/CrStates"
const TrafficMonitorConfigDocPath = "/publish/ConfigDoc"

func getClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   RequestTimeout,
	}
}

// TrafficMonitorConfigDoc represents the JSON returned by Traffic Monitor's ConfigDoc endpoint. This currently only contains the CDN member, as needed by this library.
type TrafficMonitorConfigDoc struct {
	CDN string `json:"cdnName"`
}

// GetCDN gets the CDN of the given Traffic Monitor.
func GetCDN(uri string) (string, error) {
	resp, err := getClient().Get(uri + TrafficMonitorConfigDocPath)
	if err != nil {
		return "", fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}

	configDoc := TrafficMonitorConfigDoc{}
	if err := json.Unmarshal(respBytes, &configDoc); err != nil {
		return "", fmt.Errorf("unmarshalling: %v", err)
	}
	return configDoc.CDN, nil
}

// GetCRStates gets the CRStates from the given Traffic Monitor.
func GetCRStates(uri string) (*peer.Crstates, error) {
	resp, err := getClient().Get(uri)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}

	states := peer.Crstates{}
	if err := json.Unmarshal(respBytes, &states); err != nil {
		return nil, fmt.Errorf("unmarshalling: %v", err)
	}
	return &states, nil
}

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
	crConfigBytes, err := toClient.CRConfigRaw(tmCDN)
	if err != nil {
		return fmt.Errorf("getting CRConfig: %v", err)
	}

	crConfig := crconfig.CRConfig{}
	if err := json.Unmarshal(crConfigBytes, &crConfig); err != nil {
		return fmt.Errorf("unmarshalling CRConfig JSON: %v", err)
	}

	return ValidateOfflineStatesWithCRConfig(tmURI, &crConfig, toClient)
}

// ValidateOfflineStatesWithCRConfig validates per ValidateOfflineStates, but saves querying the CRconfig if it's already fetched.
func ValidateOfflineStatesWithCRConfig(tmURI string, crConfig *crconfig.CRConfig, toClient *to.Session) error {
	crStates, err := GetCRStates(tmURI + TrafficMonitorCRStatesPath)
	if err != nil {
		return fmt.Errorf("getting CRStates: %v", err)
	}

	return ValidateCRStates(crStates, crConfig)
}

// ValidateCRStates validates that no OFFLINE or ADMIN_DOWN caches in the given CRConfig are marked Available in the given CRStates.
func ValidateCRStates(crstates *peer.Crstates, crconfig *crconfig.CRConfig) error {
	for cacheName, cacheInfo := range crconfig.ContentServers {
		status := enum.CacheStatusFromString(string(*cacheInfo.Status))
		if status != enum.CacheStatusAdminDown || status != enum.CacheStatusOffline {
			continue
		}

		available, ok := crstates.Caches[enum.CacheName(cacheName)]
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
	invalid := false
	invalidStart := time.Time{}
	for {
		err := ValidateOfflineStates(tmURI, toClient)

		if err != nil && !invalid {
			invalid = true
			invalidStart = time.Now()
		}

		if err != nil {
			invalidSpan := time.Now().Sub(invalidStart)
			if invalidSpan > grace {
				onErr(fmt.Errorf("invalid state for %v: %v\n", invalidSpan, err))
			}
		}

		if err == nil && invalid {
			onResumeSuccess()
			invalid = false
		}

		onCheck(err)

		time.Sleep(interval)
	}
}

// CRConfigOrError contains a CRConfig or an error. Union types? Monads? What are those?
type CRConfigOrError struct {
	CRConfig *crconfig.CRConfig
	Err      error
}

// ValidateOfflineStates validates that no OFFLINE or ADMIN_DOWN caches in the given Traffic Ops' CRConfig are marked Available in the given Traffic Monitor's CRStates.
func ValidateAllMonitorsOfflineStates(toClient *to.Session, includeOffline bool) (map[enum.TrafficMonitorName]error, error) {
	trafficMonitorType := "RASCAL"
	monitorTypeQuery := map[string][]string{"type": []string{trafficMonitorType}}
	servers, err := toClient.ServersByType(monitorTypeQuery)
	if err != nil {
		return nil, fmt.Errorf("getting monitors from Traffic Ops: %v", err)
	}

	if !includeOffline {
		servers = FilterOfflines(servers)
	}

	crConfigs := GetCRConfigs(GetCDNs(servers), toClient)

	errs := map[enum.TrafficMonitorName]error{}
	for _, server := range servers {
		crConfig := crConfigs[enum.CDNName(server.CDNName)]
		if err := crConfig.Err; err != nil {
			errs[enum.TrafficMonitorName(server.HostName)] = fmt.Errorf("getting CRConfig: %v", err)
			continue
		}

		uri := fmt.Sprintf("http://%s.%s", server.HostName, server.DomainName)
		errs[enum.TrafficMonitorName(server.HostName)] = ValidateOfflineStatesWithCRConfig(uri, crConfig.CRConfig, toClient)
	}
	return errs, nil
}

// AllMonitorsCRStatesOfflineValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll. Note the error passed to `onErr` may be a general validation error not associated with any monitor, in which case the passed `enum.TrafficMonitorName` will be empty.
func AllMonitorsCRStatesOfflineValidator(
	toClient *to.Session,
	interval time.Duration,
	includeOffline bool,
	grace time.Duration,
	onErr func(enum.TrafficMonitorName, error),
	onResumeSuccess func(enum.TrafficMonitorName),
	onCheck func(enum.TrafficMonitorName, error),
) {
	invalid := map[enum.TrafficMonitorName]bool{}
	invalidStart := map[enum.TrafficMonitorName]time.Time{}
	for {
		tmErrs, err := ValidateAllMonitorsOfflineStates(toClient, includeOffline)
		if err != nil {
			onErr("", fmt.Errorf("Error validating monitors: %v", err))
			time.Sleep(interval)
		}

		for name, err := range tmErrs {
			if err != nil && !invalid[name] {
				invalid[name] = true
				invalidStart[name] = time.Now()
			}

			if err != nil {
				invalidSpan := time.Now().Sub(invalidStart[name])
				if invalidSpan > grace {
					onErr(name, fmt.Errorf("invalid state for %v: %v\n", invalidSpan, err))
				}
			}

			onCheck(name, err)
		}

		for tm, tmInvalid := range invalid {
			if _, ok := tmErrs[tm]; tmInvalid && !ok {
				onResumeSuccess(tm)
				invalid[tm] = false
			}
		}

		time.Sleep(interval)
	}
}

// FilterOfflines returns only servers which are REPORTED or ONLINE
func FilterOfflines(servers []to.Server) []to.Server {
	onlineServers := []to.Server{}
	for _, server := range servers {
		status := enum.CacheStatusFromString(server.Status)
		if status != enum.CacheStatusOnline && status != enum.CacheStatusReported {
			continue
		}
		onlineServers = append(onlineServers, server)
	}
	return onlineServers
}

func GetCDNs(servers []to.Server) map[enum.CDNName]struct{} {
	cdns := map[enum.CDNName]struct{}{}
	for _, server := range servers {
		cdns[enum.CDNName(server.CDNName)] = struct{}{}
	}
	return cdns
}

func GetCRConfigs(cdns map[enum.CDNName]struct{}, toClient *to.Session) map[enum.CDNName]CRConfigOrError {
	crConfigs := map[enum.CDNName]CRConfigOrError{}
	for cdn, _ := range cdns {
		crConfigBytes, err := toClient.CRConfigRaw(string(cdn))
		if err != nil {
			crConfigs[cdn] = CRConfigOrError{Err: fmt.Errorf("getting CRConfig: %v", err)}
			continue
		}

		crConfig := crconfig.CRConfig{}
		if err := json.Unmarshal(crConfigBytes, &crConfig); err != nil {
			crConfigs[cdn] = CRConfigOrError{Err: fmt.Errorf("unmarshalling CRConfig JSON: %v", err)}
		}

		crConfigs[cdn] = CRConfigOrError{CRConfig: &crConfig}
	}
	return crConfigs
}
