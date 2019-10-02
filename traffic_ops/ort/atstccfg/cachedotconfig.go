package main

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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetConfigFileProfileCacheDotConfig(cfg TCCfg, profileNameOrID string) (string, error) {
	profileName, err := GetProfileNameFromProfileNameOrID(cfg, profileNameOrID)
	if err != nil {
		return "", errors.New("getting profile name from '" + profileNameOrID + "': " + err.Error())
	}

	toToolName, toURL, err := GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	servers, err := GetServers(cfg)
	if err != nil {
		return "", errors.New("getting servers: " + err.Error())
	}

	profileServerIDs := []int{}
	profileServerIDsMap := map[int]struct{}{}
	profileServers := []tc.Server{}
	for _, sv := range servers {
		if sv.Profile != profileName {
			continue
		}
		profileServers = append(profileServers, sv)
		profileServerIDs = append(profileServerIDs, sv.ID)
		profileServerIDsMap[sv.ID] = struct{}{}
	}

	dsServers, err := GetDeliveryServiceServers(cfg, nil, profileServerIDs)
	if err != nil {
		return "", errors.New("getting parent.config cachegroup parent server delivery service servers: " + err.Error())
	}

	profile, err := GetProfileByName(cfg, profileName)
	if err != nil {
		return "", errors.New("getting profile '" + profileNameOrID + "': " + err.Error())
	}

	dses, err := GetCDNDeliveryServices(cfg, profile.CDNID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	dsIDs := map[int]struct{}{}
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn? err?
		}
		if _, ok := profileServerIDsMap[*dss.Server]; !ok {
			continue
		}
		dsIDs[*dss.DeliveryService] = struct{}{}
	}

	profileDSes := []atscfg.ProfileDS{}
	for _, ds := range dses {
		if ds.ID == nil || ds.Type == nil || ds.OrgServerFQDN == nil {
			continue // TODO warn? err?
		}
		if *ds.Type == tc.DSTypeInvalid {
			continue // TODO warn? err?
		}
		if *ds.OrgServerFQDN == "" {
			continue // TODO warn? err?
		}
		if _, ok := dsIDs[*ds.ID]; !ok {
			continue
		}
		origin := *ds.OrgServerFQDN
		profileDSes = append(profileDSes, atscfg.ProfileDS{Type: *ds.Type, OriginFQDN: &origin})
	}
	return atscfg.MakeCacheDotConfig(profileName, profileDSes, toToolName, toURL), nil
}
