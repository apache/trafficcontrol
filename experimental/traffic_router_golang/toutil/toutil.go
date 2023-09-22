package toutil

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
 *
 */

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v3-client"
)

// MonitorAvailableStatuses returns a set of the status names that indicate a
// Traffic Monitor instance is "available".
func MonitorAvailableStatuses() map[string]struct{} {
	return map[string]struct{}{
		"ONLINE":   {},
		"REPORTED": {},
	}
}

// GetMonitorURIs returns a list of URLs for the root of all Traffic Monitor
// instances in the given CDN - or an error if one occurred.
func GetMonitorURIs(toc *client.Session, cdn string) ([]string, error) {
	params := url.Values{}
	params.Set("type", tc.MonitorTypeName)
	servers, reqInf, err := toc.GetServersWithHdr(&params, nil)
	if err != nil {
		return nil, errors.New("getting servers by type from '" + reqInf.RemoteAddr.String() + "':" + err.Error())
	}
	availableStatuses := MonitorAvailableStatuses()

	monitors := []string{}
	for _, server := range servers.Response {
		if server.CDNName == nil || server.HostName == nil || server.DomainName == nil || server.Status == nil {
			continue
		}
		if *server.CDNName != cdn {
			continue
		}
		if _, ok := availableStatuses[*server.Status]; !ok {
			continue
		}
		m := "http://" + *server.HostName + "." + *server.DomainName
		if server.TCPPort != nil && *server.TCPPort > 0 && *server.TCPPort != 80 {
			m += ":" + strconv.Itoa(*server.TCPPort)
		}
		monitors = append(monitors, m)
	}
	return monitors, nil
}
