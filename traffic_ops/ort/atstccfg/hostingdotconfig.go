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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const ServerHostingDotConfigMidIncludeInactive = false
const ServerHostingDotConfigEdgeIncludeInactive = true

func GetConfigFileServerHostingDotConfig(cfg TCCfg, serverNameOrID string) (string, error) {
	// TODO TOAPI add /servers?cdn=1 query param
	servers, err := GetServers(cfg)
	if err != nil {
		return "", errors.New("getting servers: " + err.Error())
	}

	server := tc.Server{ID: atscfg.InvalidID}
	if serverID, err := strconv.Atoi(serverNameOrID); err == nil {
		for _, toServer := range servers {
			if toServer.ID == serverID {
				server = toServer
				break
			}
		}
	} else {
		serverName := serverNameOrID
		for _, toServer := range servers {
			if toServer.HostName == serverName {
				server = toServer
				break
			}
		}
	}
	if server.ID == atscfg.InvalidID {
		return "", errors.New("server '" + serverNameOrID + " not found in servers")
	}

	serverName := tc.CacheName(server.HostName)

	toToolName, toURL, err := GetTOToolNameAndURLFromTO(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	profileParams, err := GetProfileParameters(cfg, server.Profile)
	if err != nil {
		return "", errors.New("getting profile '" + server.Profile + "' parameters: " + err.Error())
	}
	if len(profileParams) == 0 {
		// The TO endpoint behind toclient.GetParametersByProfileName returns an empty object with a 200, if the Profile doesn't exist.
		// So we act as though we got a 404 if there are no params, to make ORT behave correctly.
		return "", ErrNotFound
	}

	fileParams := map[string]string{}
	for _, param := range profileParams {
		if param.ConfigFile != atscfg.HostingConfigParamConfigFile {
			continue
		}
		if val, ok := fileParams[param.Name]; ok {
			log.Errorln("hosting config parameter name '" + param.Name + "' got multiple values - using '" + val + "'")
			continue
		}
		fileParams[param.Name] = param.Value
	}

	dses, err := GetCDNDeliveryServices(cfg, server.CDNID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	cdnServers := map[tc.CacheName]tc.Server{}
	for _, sv := range servers {
		if sv.CDNID != server.CDNID {
			continue
		}
		cdnServers[tc.CacheName(sv.HostName)] = sv
	}

	serverIDs := []int{}
	for _, sv := range cdnServers {
		serverIDs = append(serverIDs, sv.ID)
	}

	dsIDs := []int{}
	for _, ds := range dses {
		if ds.ID != nil {
			dsIDs = append(dsIDs, *ds.ID)
		}
	}

	dsServers, err := GetDeliveryServiceServers(cfg, dsIDs, serverIDs)
	if err != nil {
		return "", errors.New("getting delivery service servers: " + err.Error())
	}

	dsServerMap := map[int]map[int]struct{}{} // set[dsID][serverID]
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			return "", errors.New("deliveryserviceservers returned dss with nil values")
		}
		if _, ok := dsServerMap[*dss.DeliveryService]; !ok {
			dsServerMap[*dss.DeliveryService] = map[int]struct{}{}
		}
		dsServerMap[*dss.DeliveryService][*dss.Server] = struct{}{}
	}

	hostingDSes := map[tc.DeliveryServiceName]tc.DeliveryServiceNullable{}
	for _, ds := range dses {
		if ds.Active == nil || ds.Type == nil || ds.XMLID == nil || ds.CDNID == nil || ds.ID == nil || ds.OrgServerFQDN == nil {
			// some DSes have nil origins. I think MSO? TODO: verify
			continue
		}

		if !ServerHostingDotConfigMidIncludeInactive && !*ds.Active {
			continue
		}
		if *ds.CDNID != server.CDNID {
			continue
		}

		if strings.HasPrefix(server.Type, tc.MidTypePrefix) {
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

			if _, ok := dsServerMap[*ds.ID][server.ID]; !ok {
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

	txt := atscfg.MakeHostingDotConfig(serverName, toToolName, toURL, fileParams, origins)
	return txt, nil
}
