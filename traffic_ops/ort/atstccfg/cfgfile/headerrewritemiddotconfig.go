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
)

func GetConfigFileCDNHeaderRewriteMid(toData *TOData, fileName string) (string, string, error) {
	dsName := strings.TrimSuffix(strings.TrimPrefix(fileName, atscfg.HeaderRewriteMidPrefix), atscfg.ConfigSuffix) // TODO verify prefix and suffix? Perl doesn't

	tcDS := tc.DeliveryServiceNullable{}
	for _, ds := range toData.DeliveryServices {
		if ds.XMLID == nil || *ds.XMLID != dsName {
			continue
		}
		tcDS = ds
		break
	}
	if tcDS.ID == nil {
		return "", "", errors.New("ds '" + dsName + "' not found")
	}

	if tcDS.CDNName == nil {
		return "", "", errors.New("ds '" + dsName + "' missing cdn")
	}

	cfgDS, err := atscfg.HeaderRewriteDSFromDS(&tcDS)
	if err != nil {
		return "", "", errors.New("converting ds to config ds: " + err.Error())
	}

	dsServers := FilterDSS(toData.DeliveryServiceServers, map[int]struct{}{cfgDS.ID: {}}, nil)

	dsServerIDs := map[int]struct{}{}
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn?
		}
		if *dss.DeliveryService != *tcDS.ID {
			continue
		}
		dsServerIDs[*dss.Server] = struct{}{}
	}

	serverCGs := map[tc.CacheGroupName]struct{}{}
	for _, sv := range toData.Servers {
		if sv.CDNName != toData.Server.CDNName {
			continue
		}
		if tc.CacheStatus(sv.Status) != tc.CacheStatusReported && tc.CacheStatus(sv.Status) != tc.CacheStatusOnline {
			continue
		}
		serverCGs[tc.CacheGroupName(sv.Cachegroup)] = struct{}{}
	}

	parentCGs := map[string]struct{}{} // names of cachegroups which are parent cachegroups of the cachegroup of any edge assigned to the given DS
	for _, cg := range toData.CacheGroups {
		if cg.Name == nil {
			continue // TODO warn?
		}
		if cg.ParentName == nil {
			continue // TODO warn?
		}
		if _, ok := serverCGs[tc.CacheGroupName(*cg.Name)]; !ok {
			continue
		}
		parentCGs[*cg.ParentName] = struct{}{}
	}

	parentServers := []tc.Server{}
	for _, sv := range toData.Servers {
		if _, ok := parentCGs[sv.Cachegroup]; !ok {
			continue
		}
		parentServers = append(parentServers, sv)
	}

	assignedMids := []atscfg.HeaderRewriteServer{}
	for _, server := range toData.Servers {
		if server.CDNName != *tcDS.CDNName {
			continue
		}
		if _, ok := parentCGs[server.Cachegroup]; !ok {
			continue
		}
		cfgServer, err := atscfg.HeaderRewriteServerFromServerNotNullable(server)
		if err != nil {
			continue // TODO warn?
		}
		assignedMids = append(assignedMids, cfgServer)
	}

	return atscfg.MakeHeaderRewriteMidDotConfig(tc.CDNName(toData.Server.CDNName), toData.TOToolName, toData.TOURL, cfgDS, assignedMids), atscfg.ContentTypeHeaderRewriteDotConfig, nil
}
