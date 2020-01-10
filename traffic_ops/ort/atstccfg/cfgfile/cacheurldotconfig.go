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
	"github.com/apache/trafficcontrol/lib/go-tc/enum"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/toreq"
)

func GetConfigFileCDNCacheURL(cfg config.TCCfg, cdnNameOrID string, fileName string) (string, error) {
	cdnName, err := toreq.GetCDNNameFromCDNNameOrID(cfg, cdnNameOrID)
	if err != nil {
		return "", errors.New("getting CDN name from '" + cdnNameOrID + "': " + err.Error())
	}

	toToolName, toURL, err := toreq.GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	cdn, err := toreq.GetCDN(cfg, cdnName)
	if err != nil {
		return "", errors.New("getting cdn '" + string(cdnName) + "': " + err.Error())
	}

	dses, err := toreq.GetCDNDeliveryServices(cfg, cdn.ID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	dsIDs := []int{}
	for _, ds := range dses {
		if ds.ID != nil {
			dsIDs = append(dsIDs, *ds.ID)
		}
	}

	dss, err := toreq.GetDeliveryServiceServers(cfg, dsIDs, nil)
	if err != nil {
		return "", errors.New("getting delivery service servers: " + err.Error())
	}

	dssMap := map[int][]int{} // map[dsID]serverID
	for _, dss := range dss {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn?
		}
		dssMap[*dss.DeliveryService] = append(dssMap[*dss.DeliveryService], *dss.Server)
	}

	dsesWithServers := []tc.DeliveryServiceNullable{}
	for _, ds := range dses {
		if ds.ID == nil {
			continue // TODO warn
		}
		// ANY_MAP and STEERING DSes don't have origins, and thus can't be put into the cacheurl config.
		if ds.Type != nil && (*ds.Type == enum.DSTypeAnyMap || *ds.Type == enum.DSTypeSteering) {
			continue
		}
		if len(dssMap[*ds.ID]) == 0 {
			continue
		}
		dsesWithServers = append(dsesWithServers, ds)
	}

	cfgDSes := atscfg.DeliveryServicesToCacheURLDSes(dsesWithServers)

	txt := atscfg.MakeCacheURLDotConfig(cdnName, toToolName, toURL, fileName, cfgDSes)
	return txt, nil
}

func GetConfigFileCDNCacheURLPlain(cfg config.TCCfg, cdnNameOrID string) (string, error) {
	return GetConfigFileCDNCacheURL(cfg, cdnNameOrID, "cacheurl.config")
}
