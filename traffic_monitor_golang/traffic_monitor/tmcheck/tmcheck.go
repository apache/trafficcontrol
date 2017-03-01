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

	crConfigBytes, err := toClient.CRConfigRaw(cdn)
	if err != nil {
		return fmt.Errorf("getting CRConfig: %v", err)
	}

	crConfig := crconfig.CRConfig{}
	if err := json.Unmarshal(crConfigBytes, &crConfig); err != nil {
		return fmt.Errorf("unmarshalling CRConfig JSON: %v", err)
	}

	crStates, err := GetCRStates(tmURI + TrafficMonitorCRStatesPath)
	if err != nil {
		return fmt.Errorf("getting CRStates: %v", err)
	}

	return ValidateCRStates(crStates, &crConfig)
}

// ValidateCRStates validates that no OFFLINE or ADMIN_DOWN caches in the given CRConfig are marked Available in the given CRStates.
func ValidateCRStates(crstates *peer.Crstates, crconfig *crconfig.CRConfig) error {
	for cacheName, cacheInfo := range crconfig.ContentServers {
		status := enum.CacheStatusFromString(string(*cacheInfo.Status))
		if status != enum.CacheStatusOffline || status != enum.CacheStatusOffline {
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

// Validator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll.
func Validator(
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
