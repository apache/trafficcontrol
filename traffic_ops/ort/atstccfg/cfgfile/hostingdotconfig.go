package cfgfile

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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

const ServerHostingDotConfigMidIncludeInactive = false
const ServerHostingDotConfigEdgeIncludeInactive = true

func GetConfigFileServerHostingDotConfig(toData *config.TOData) (string, string, error) {
	fileParams := ParamsToMap(FilterParams(toData.ServerParams, atscfg.HostingConfigParamConfigFile, "", "", ""))

	cdnServers := map[tc.CacheName]tc.Server{}
	for _, sv := range toData.Servers {
		if sv.CDNID != toData.Server.CDNID {
			continue
		}
		cdnServers[tc.CacheName(sv.HostName)] = sv
	}

	serverIDs := map[int]struct{}{}
	for _, sv := range cdnServers {
		serverIDs[sv.ID] = struct{}{}
	}

	dsIDs := map[int]struct{}{}
	for _, ds := range toData.DeliveryServices {
		if ds.ID != nil {
			dsIDs[*ds.ID] = struct{}{}
		}
	}

	dsServers := FilterDSS(toData.DeliveryServiceServers, dsIDs, serverIDs)

	dsServerMap := map[int]map[int]struct{}{} // set[dsID][serverID]
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			return "", "", errors.New("deliveryserviceservers returned dss with nil values")
		}
		if _, ok := dsServerMap[*dss.DeliveryService]; !ok {
			dsServerMap[*dss.DeliveryService] = map[int]struct{}{}
		}
		dsServerMap[*dss.DeliveryService][*dss.Server] = struct{}{}
	}

	hostingDSes := map[tc.DeliveryServiceName]tc.DeliveryServiceNullable{}

	isMid := strings.HasPrefix(toData.Server.Type, tc.MidTypePrefix)

	for _, ds := range toData.DeliveryServices {
		if ds.Active == nil || ds.Type == nil || ds.XMLID == nil || ds.CDNID == nil || ds.ID == nil || ds.OrgServerFQDN == nil {
			// some DSes have nil origins. I think MSO? TODO: verify
			continue
		}

		if !*ds.Active && ((!isMid && !ServerHostingDotConfigEdgeIncludeInactive) || (isMid && !ServerHostingDotConfigMidIncludeInactive)) {
			continue
		}

		if *ds.CDNID != toData.Server.CDNID {
			continue
		}

		if isMid {
			if !strings.HasSuffix(string(*ds.Type), tc.DSTypeLiveNationalSuffix) {
				continue
			}

			// mids: include all DSes with at least one server assigned
			if len(dsServerMap[*ds.ID]) == 0 {
				continue
			}
		} else {
			if !strings.HasSuffix(string(*ds.Type), tc.DSTypeLiveNationalSuffix) && !strings.HasSuffix(string(*ds.Type), tc.DSTypeLiveSuffix) {
				continue
			}

			// edges: only include DSes assigned to this edge
			if dsServerMap[*ds.ID] == nil {
				continue
			}

			if _, ok := dsServerMap[*ds.ID][toData.Server.ID]; !ok {
				continue
			}
		}

		hostingDSes[tc.DeliveryServiceName(*ds.XMLID)] = ds
	}

	originSet := map[string]struct{}{}
	for _, ds := range hostingDSes {
		originSet[*ds.OrgServerFQDN] = struct{}{}
	}
	origins := []string{}
	for origin, _ := range originSet {
		origins = append(origins, origin)
	}

	return atscfg.MakeHostingDotConfig(tc.CacheName(toData.Server.HostName), toData.TOToolName, toData.TOURL, fileParams, origins), atscfg.ContentTypeHostingDotConfig, nil
}
