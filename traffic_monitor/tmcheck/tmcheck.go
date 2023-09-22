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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/datareq"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"

	jsoniter "github.com/json-iterator/go"
)

const RequestTimeout = time.Second * time.Duration(30)

const TrafficMonitorCRStatesPath = "/publish/CrStates"
const TrafficMonitorDSStatsPath = "/publish/DsStats"
const TrafficMonitorConfigDocPath = "/publish/ConfigDoc"
const TrafficMonitorStatsPath = "/publish/Stats"

func getClient() *http.Client {
	return &http.Client{
		Timeout: RequestTimeout,
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
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}

	configDoc := TrafficMonitorConfigDoc{}
	json := jsoniter.ConfigFastest
	if err := json.Unmarshal(respBytes, &configDoc); err != nil {
		return "", fmt.Errorf("unmarshalling: %v", err)
	}
	return configDoc.CDN, nil
}

// GetCRStates gets the CRStates from the given Traffic Monitor.
func GetCRStates(uri string) (*tc.CRStates, error) {
	resp, err := getClient().Get(uri)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}

	states := tc.CRStates{}
	json := jsoniter.ConfigFastest
	if err := json.Unmarshal(respBytes, &states); err != nil {
		return nil, fmt.Errorf("unmarshalling: %v", err)
	}
	return &states, nil
}

// GetDSStats gets the DSStats from the given Traffic Monitor.
func GetDSStats(uri string) (*dsdata.StatsOld, error) {
	resp, err := getClient().Get(uri)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}

	dsStats := dsdata.StatsOld{}
	json := jsoniter.ConfigFastest
	if err := json.Unmarshal(respBytes, &dsStats); err != nil {
		return nil, fmt.Errorf("unmarshalling: %v", err)
	}
	return &dsStats, nil
}

// GetStats gets the stats from the given Traffic Monitor.
func GetStats(uri string) (*datareq.Stats, error) {
	resp, err := getClient().Get(uri)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading reply from %v: %v\n", uri, err)
	}

	stats := datareq.JSONStats{}
	json := jsoniter.ConfigFastest
	if err := json.Unmarshal(respBytes, &stats); err != nil {
		return nil, fmt.Errorf("unmarshalling: %v", err)
	}
	return &stats.Stats, nil
}

type ValidatorFunc func(
	tmURI string,
	toClient *to.Session,
	interval time.Duration,
	grace time.Duration,
	onErr func(error),
	onResumeSuccess func(),
	onCheck func(error),
)

type AllValidatorFunc func(
	toClient *to.Session,
	interval time.Duration,
	includeOffline bool,
	grace time.Duration,
	onErr func(tc.TrafficMonitorName, error),
	onResumeSuccess func(tc.TrafficMonitorName),
	onCheck func(tc.TrafficMonitorName, error),
)

// CRStatesOfflineValidator is designed to be run as a goroutine, and does not return. It continously validates every `interval`, and calls `onErr` on failure, `onResumeSuccess` when a failure ceases, and `onCheck` on every poll.
func Validator(
	tmURI string,
	toClient *to.Session,
	interval time.Duration,
	grace time.Duration,
	onErr func(error),
	onResumeSuccess func(),
	onCheck func(error),
	validator func(tmURI string, toClient *to.Session) error,
) {
	invalid := false
	invalidStart := time.Time{}
	for {
		err := validator(tmURI, toClient)

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
	CRConfig *tc.CRConfig
	Err      error
}

func GetMonitors(toClient *to.Session, includeOffline bool) ([]tc.ServerV4, error) {
	trafficMonitorType := "RASCAL"
	query := url.Values{}
	query.Set("type", trafficMonitorType)
	response, _, err := toClient.GetServers(to.RequestOptions{QueryParameters: query})
	servers := response.Response
	if err != nil {
		return nil, fmt.Errorf("getting monitors from Traffic Ops: %v", err)
	}

	if !includeOffline {
		servers = FilterOfflines(servers)
	}
	return servers, nil
}

func AllValidator(
	toClient *to.Session,
	interval time.Duration,
	includeOffline bool,
	grace time.Duration,
	onErr func(tc.TrafficMonitorName, error),
	onResumeSuccess func(tc.TrafficMonitorName),
	onCheck func(tc.TrafficMonitorName, error),
	validator func(toClient *to.Session, includeOffline bool) (map[tc.TrafficMonitorName]error, error),
) {
	invalid := map[tc.TrafficMonitorName]bool{}
	invalidStart := map[tc.TrafficMonitorName]time.Time{}
	metaFail := false
	for {
		tmErrs, err := validator(toClient, includeOffline)
		if err != nil {
			onErr("", fmt.Errorf("Error validating monitors: %v", err))
			time.Sleep(interval)
			metaFail = true
		} else if metaFail {
			onResumeSuccess("")
			metaFail = false
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
func FilterOfflines(servers []tc.ServerV4) []tc.ServerV4 {
	onlineServers := []tc.ServerV4{}
	for _, server := range servers {
		status := tc.CacheStatusFromString(*server.Status)
		if status != tc.CacheStatusOnline && status != tc.CacheStatusReported {
			continue
		}
		onlineServers = append(onlineServers, server)
	}
	return onlineServers
}

func GetCDNs(servers []tc.ServerV4) map[tc.CDNName]struct{} {
	cdns := map[tc.CDNName]struct{}{}
	for _, server := range servers {
		cdns[tc.CDNName(*server.CDNName)] = struct{}{}
	}
	return cdns
}

func GetCRConfigs(cdns map[tc.CDNName]struct{}, toClient *to.Session) map[tc.CDNName]CRConfigOrError {
	crConfigs := map[tc.CDNName]CRConfigOrError{}
	for cdn, _ := range cdns {
		response, _, err := toClient.GetCRConfig(string(cdn), to.RequestOptions{})
		if err != nil {
			crConfigs[cdn] = CRConfigOrError{Err: fmt.Errorf("getting CRConfig: %v", err)}
			continue
		}

		crConfig := response.Response
		crConfigs[cdn] = CRConfigOrError{CRConfig: &crConfig}
	}
	return crConfigs
}
