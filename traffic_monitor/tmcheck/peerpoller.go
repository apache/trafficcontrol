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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
	"time"
)

const PeerPollMax = time.Duration(10) * time.Second

func GetOldestPolledPeerTime(uri string) (time.Duration, error) {
	stats, err := GetStats(uri + TrafficMonitorStatsPath)
	if err != nil {
		return time.Duration(0), fmt.Errorf("getting stats: %v", err)
	}

	oldestPolledPeerTime := time.Duration(stats.OldestPolledPeerMs) * time.Millisecond

	return oldestPolledPeerTime, nil
}

func ValidatePeerPoller(uri string) error {
	lastPollTime, err := GetOldestPolledPeerTime(uri)
	if err != nil {
		return fmt.Errorf("failed to get oldest peer time: %v", err)
	}
	if lastPollTime > PeerPollMax {
		return fmt.Errorf("Peer poller is dead, last poll was %v ago", lastPollTime)
	}
	return nil
}

func ValidateAllPeerPollers(toClient *to.Session, includeOffline bool) (map[tc.TrafficMonitorName]error, error) {
	servers, err := GetMonitors(toClient, includeOffline)
	if err != nil {
		return nil, err
	}
	errs := map[tc.TrafficMonitorName]error{}
	for _, server := range servers {
		uri := fmt.Sprintf("http://%s.%s", *server.HostName, *server.DomainName)
		errs[tc.TrafficMonitorName(*server.HostName)] = ValidatePeerPoller(uri)
	}
	return errs, nil
}

func PeerPollersValidator(
	tmURI string,
	toClient *to.Session,
	interval time.Duration,
	grace time.Duration,
	onErr func(error),
	onResumeSuccess func(),
	onCheck func(error),
) {
	wrapValidatePeerPoller := func(uri string, _ *to.Session) error { return ValidatePeerPoller(uri) }
	Validator(tmURI, toClient, interval, grace, onErr, onResumeSuccess, onCheck, wrapValidatePeerPoller)
}

func PeerPollersAllValidator(
	toClient *to.Session,
	interval time.Duration,
	includeOffline bool,
	grace time.Duration,
	onErr func(tc.TrafficMonitorName, error),
	onResumeSuccess func(tc.TrafficMonitorName),
	onCheck func(tc.TrafficMonitorName, error),
) {
	AllValidator(toClient, interval, includeOffline, grace, onErr, onResumeSuccess, onCheck, ValidateAllPeerPollers)
}
