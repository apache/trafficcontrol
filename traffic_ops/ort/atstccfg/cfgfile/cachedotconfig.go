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

func GetConfigFileProfileCacheDotConfig(toData *TOData) (string, string, error) {
	profileServerIDsMap := map[int]struct{}{}
	for _, sv := range toData.Servers {
		if sv.Profile != toData.Server.Profile {
			continue
		}
		profileServerIDsMap[sv.ID] = struct{}{}
	}

	dsServers := FilterDSS(toData.DeliveryServiceServers, nil, profileServerIDsMap)

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
	for _, ds := range toData.DeliveryServices {
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

	return atscfg.MakeCacheDotConfig(toData.Server.Profile, profileDSes, toData.TOToolName, toData.TOURL), atscfg.ContentTypeCacheDotConfig, nil
}
