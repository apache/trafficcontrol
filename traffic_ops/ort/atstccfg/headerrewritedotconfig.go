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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetConfigFileCDNHeaderRewrite(cfg TCCfg, cdnNameOrID string, fileName string) (string, error) {

	cdnName, err := GetCDNNameFromCDNNameOrID(cfg, cdnNameOrID)
	if err != nil {
		return "", errors.New("getting CDN name from '" + cdnNameOrID + "': " + err.Error())
	}

	toToolName, toURL, err := GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	cdn, err := GetCDN(cfg, cdnName)
	if err != nil {
		return "", errors.New("getting cdn '" + string(cdnName) + "': " + err.Error())
	}

	dsName := strings.TrimSuffix(strings.TrimPrefix(fileName, atscfg.HeaderRewritePrefix), atscfg.ConfigSuffix) // TODO verify prefix and suffix? Perl doesn't

	// TODO determine if it's faster overall for ORT to get all DSes (which we already have for other configs), or to get the single DS we need
	dses, err := GetCDNDeliveryServices(cfg, cdn.ID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}
	tcDS := tc.DeliveryServiceNullable{}
	for _, ds := range dses {
		if ds.XMLID == nil || *ds.XMLID != dsName {
			continue
		}
		tcDS = ds
		break
	}
	if tcDS.ID == nil {
		return "", errors.New("ds '" + dsName + "' not found")
	}

	if tcDS.CDNName == nil {
		return "", errors.New("ds '" + dsName + "' missing cdn")
	}

	cfgDS, err := atscfg.HeaderRewriteDSFromDS(&tcDS)
	if err != nil {
		return "", errors.New("converting ds to config ds: " + err.Error())
	}

	dsServers, err := GetDeliveryServiceServers(cfg, []int{cfgDS.ID}, nil)
	if err != nil {
		return "", errors.New("getting delivery service servers: " + err.Error())
	}

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

	servers, err := GetServers(cfg)
	if err != nil {
		return "", errors.New("getting servers: " + err.Error())
	}

	assignedEdges := []atscfg.HeaderRewriteServer{}
	for _, server := range servers {
		if server.CDNName != *tcDS.CDNName {
			continue
		}
		if _, ok := dsServerIDs[server.ID]; !ok {
			continue
		}
		cfgServer, err := atscfg.HeaderRewriteServerFromServerNotNullable(server)
		if err != nil {
			continue // TODO warn?
		}
		assignedEdges = append(assignedEdges, cfgServer)
	}

	return atscfg.MakeHeaderRewriteDotConfig(cdnName, toToolName, toURL, cfgDS, assignedEdges), nil
}
