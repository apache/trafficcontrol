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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetConfigFileMeta(cfg TCCfg, serverNameOrID string) (string, error) {
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

	serverHostName := server.HostName

	cacheGroups, err := GetCacheGroups(cfg)
	if err != nil {
		return "", errors.New("getting cachegroups: " + err.Error())
	}

	cgMap := map[string]tc.CacheGroupNullable{}
	for _, cg := range cacheGroups {
		if cg.Name == nil {
			return "", errors.New("got cachegroup with nil name!'")
		}
		cgMap[*cg.Name] = cg
	}

	serverCG, ok := cgMap[server.Cachegroup]
	if !ok {
		return "", errors.New("server '" + serverNameOrID + "' cachegroup '" + server.Cachegroup + "' not found in CacheGroups")
	}

	parentCGID := -1
	parentCGType := ""
	if serverCG.ParentName != nil && *serverCG.ParentName != "" {
		parentCG, ok := cgMap[*serverCG.ParentName]
		if !ok {
			return "", errors.New("server '" + serverNameOrID + "' cachegroup '" + server.Cachegroup + "' parent '" + *serverCG.ParentName + "' not found in CacheGroups")
		}
		if parentCG.ID == nil {
			return "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		parentCGID = *parentCG.ID

		if parentCG.Type == nil {
			return "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}
		parentCGType = *parentCG.Type
	}

	secondaryParentCGID := -1
	secondaryParentCGType := ""
	if serverCG.SecondaryParentName != nil && *serverCG.SecondaryParentName != "" {
		parentCG, ok := cgMap[*serverCG.SecondaryParentName]
		if !ok {
			return "", errors.New("server '" + serverNameOrID + "' cachegroup '" + server.Cachegroup + "' secondary parent '" + *serverCG.SecondaryParentName + "' not found in CacheGroups")
		}

		if parentCG.ID == nil {
			return "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		secondaryParentCGID = *parentCG.ID
		if parentCG.Type == nil {
			return "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}

		secondaryParentCGType = *parentCG.Type
	}

	serverInfo := atscfg.ServerInfo{
		CacheGroupID:                  server.CachegroupID,
		CDN:                           tc.CDNName(server.CDNName),
		CDNID:                         server.CDNID,
		DomainName:                    server.DomainName,
		HostName:                      server.HostName,
		ID:                            server.ID,
		IP:                            server.IPAddress,
		ParentCacheGroupID:            parentCGID,
		ParentCacheGroupType:          parentCGType,
		ProfileID:                     atscfg.ProfileID(server.ProfileID),
		ProfileName:                   server.Profile,
		Port:                          server.TCPPort,
		SecondaryParentCacheGroupID:   secondaryParentCGID,
		SecondaryParentCacheGroupType: secondaryParentCGType,
		Type:                          server.Type,
	}

	globalParams, err := GetGlobalParameters(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	toReverseProxyURL := ""
	toURL := ""
	for _, param := range globalParams {
		if param.Name == "tm.rev_proxy.url" {
			toReverseProxyURL = param.Value
		} else if param.Name == "tm.url" {
			toURL = param.Value
		}
		if toReverseProxyURL != "" && toURL != "" {
			break
		}
	}

	scopeParamsRaw, err := GetParametersByName(cfg, "scope")
	if err != nil {
		return "", errors.New("getting scope parameters: " + err.Error())
	}

	scopeParams := map[string]string{}
	for _, param := range scopeParamsRaw {
		scopeParams[param.ConfigFile] = param.Value
	}

	serverProfileParameters, err := GetServerProfileParameters(cfg, server.Profile)
	if err != nil {
		return "", errors.New("getting server profile '" + server.Profile + "' parameters: " + err.Error())
	}

	locationParams := map[string]atscfg.ConfigProfileParams{}
	for _, param := range serverProfileParameters {
		if param.Name == "location" {
			p := locationParams[param.ConfigFile]
			p.FileNameOnDisk = param.ConfigFile
			p.Location = param.Value
			locationParams[param.ConfigFile] = p
		} else if param.Name == "URL" {
			p := locationParams[param.ConfigFile]
			p.URL = param.Value
			locationParams[param.ConfigFile] = p
		}
	}

	deliveryServices, err := GetCDNDeliveryServices(cfg, server.CDNID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	dsIDs := []int{}
	for _, ds := range deliveryServices {
		if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURISigning {
			continue
		}
		if ds.ID == nil {
			// TODO log error?
			continue
		}
		dsIDs = append(dsIDs, *ds.ID)
	}

	serverIDs := []int{server.ID}

	dsServers, err := GetDeliveryServiceServers(cfg, dsIDs, serverIDs)
	if err != nil {
		return "", errors.New("getting meta config delivery service servers: " + err.Error())
	}

	dssMap := map[int]struct{}{} // set of map[dsID]. We know we only asked for our own server, so we don't care about the servers returned.
	for _, dss := range dsServers {
		if dss.DeliveryService == nil {
			continue // TODO log?
		}
		dssMap[*dss.DeliveryService] = struct{}{}
	}

	uriSignedDSes := []tc.DeliveryServiceName{}
	for _, ds := range deliveryServices {
		if ds.ID == nil {
			continue
		}
		if ds.XMLID == nil {
			continue // TODO log?
		}
		if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURISigning {
			continue
		}
		if _, ok := dssMap[*ds.ID]; !ok {
			continue
		}
		uriSignedDSes = append(uriSignedDSes, tc.DeliveryServiceName(*ds.XMLID))
	}

	return atscfg.MakeMetaConfig(tc.CacheName(serverHostName), &serverInfo, toURL, toReverseProxyURL, locationParams, uriSignedDSes, scopeParams), nil
}
