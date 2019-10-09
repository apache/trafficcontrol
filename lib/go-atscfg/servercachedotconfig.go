package atscfg

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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

type ServerCacheConfigDS struct {
	OrgServerFQDN string
	Type          tc.DSType
}

func MakeServerCacheDotConfig(
	serverName tc.CacheName,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	dses map[tc.DeliveryServiceName]ServerCacheConfigDS,
) string {
	text := GenericHeaderComment(string(serverName), toToolName, toURL)

	seenOrigins := map[string]struct{}{}
	for _, ds := range dses {
		if ds.Type != tc.DSTypeHTTPNoCache {
			continue
		}
		if _, ok := seenOrigins[ds.OrgServerFQDN]; ok {
			continue
		}
		seenOrigins[ds.OrgServerFQDN] = struct{}{}

		originFQDN, originPort := GetOriginFQDNAndPort(ds.OrgServerFQDN)
		if originPort != nil {
			text += `dest_domain=` + originFQDN + ` port=` + strconv.Itoa(*originPort) + ` scheme=http action=never-cache` + "\n"
		} else {
			text += `dest_domain=` + originFQDN + ` scheme=http action=never-cache` + "\n"
		}
	}
	return text
}

// TODO unit test
func GetOriginFQDNAndPort(origin string) (string, *int) {
	origin = strings.TrimSpace(origin)
	origin = strings.Replace(origin, `https://`, ``, -1)
	origin = strings.Replace(origin, `http://`, ``, -1)

	// if the origin includes a path, strip it

	slashI := strings.Index(origin, `/`)
	if slashI != -1 {
		origin = origin[:slashI]
	}

	hostName := origin

	colonI := strings.Index(origin, ":")
	if colonI == -1 {
		return hostName, nil // no :, the origin must not include a port
	}
	portStr := origin[colonI+1:]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		// either the port isn't an integer, or the : we found was something else.
		// Return the origin, as if it didn't contain a port.
		return hostName, nil
	}

	hostName = origin[:colonI]
	return hostName, &port
}
