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
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetConfigFileCDNCacheURL(toData *TOData, fileName string) (string, string, error) {
	dsIDs := map[int]struct{}{}
	for _, ds := range toData.DeliveryServices {
		if ds.ID != nil {
			dsIDs[*ds.ID] = struct{}{}
		}
	}

	dss := FilterDSS(toData.DeliveryServiceServers, dsIDs, nil)

	dssMap := map[int][]int{} // map[dsID]serverID
	for _, dss := range dss {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn?
		}
		dssMap[*dss.DeliveryService] = append(dssMap[*dss.DeliveryService], *dss.Server)
	}

	dsesWithServers := []tc.DeliveryServiceNullable{}
	for _, ds := range toData.DeliveryServices {
		if ds.ID == nil {
			continue // TODO warn
		}
		// ANY_MAP and STEERING DSes don't have origins, and thus can't be put into the cacheurl config.
		if ds.Type != nil && (*ds.Type == tc.DSTypeAnyMap || *ds.Type == tc.DSTypeSteering) {
			continue
		}
		if len(dssMap[*ds.ID]) == 0 {
			continue
		}
		dsesWithServers = append(dsesWithServers, ds)
	}

	cfgDSes := atscfg.DeliveryServicesToCacheURLDSes(dsesWithServers)

	return atscfg.MakeCacheURLDotConfig(tc.CDNName(toData.Server.CDNName), toData.TOToolName, toData.TOURL, fileName, cfgDSes), atscfg.ContentTypeCacheURLDotConfig, nil
}

func GetConfigFileCDNCacheURLPlain(toData *TOData) (string, string, error) {
	return GetConfigFileCDNCacheURL(toData, "cacheurl.config")
}
