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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

func GetMeta(toData *config.TOData) (*tc.ATSConfigMetaData, error) {
	cgMap := map[string]tc.CacheGroupNullable{}
	for _, cg := range toData.CacheGroups {
		if cg.Name == nil {
			return nil, errors.New("got cachegroup with nil name!'")
		}
		cgMap[*cg.Name] = cg
	}

	serverCG, ok := cgMap[toData.Server.Cachegroup]
	if !ok {
		return nil, errors.New("server '" + toData.Server.HostName + "' cachegroup '" + toData.Server.Cachegroup + "' not found in CacheGroups")
	}

	parentCGID := -1
	parentCGType := ""
	if serverCG.ParentName != nil && *serverCG.ParentName != "" {
		parentCG, ok := cgMap[*serverCG.ParentName]
		if !ok {
			return nil, errors.New("server '" + toData.Server.HostName + "' cachegroup '" + toData.Server.Cachegroup + "' parent '" + *serverCG.ParentName + "' not found in CacheGroups")
		}
		if parentCG.ID == nil {
			return nil, errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		parentCGID = *parentCG.ID

		if parentCG.Type == nil {
			return nil, errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}
		parentCGType = *parentCG.Type
	}

	secondaryParentCGID := -1
	secondaryParentCGType := ""
	if serverCG.SecondaryParentName != nil && *serverCG.SecondaryParentName != "" {
		parentCG, ok := cgMap[*serverCG.SecondaryParentName]
		if !ok {
			return nil, errors.New("server '" + toData.Server.HostName + "' cachegroup '" + toData.Server.Cachegroup + "' secondary parent '" + *serverCG.SecondaryParentName + "' not found in CacheGroups")
		}

		if parentCG.ID == nil {
			return nil, errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		secondaryParentCGID = *parentCG.ID
		if parentCG.Type == nil {
			return nil, errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}

		secondaryParentCGType = *parentCG.Type
	}

	serverInfo := atscfg.ServerInfo{
		CacheGroupID:                  toData.Server.CachegroupID,
		CDN:                           tc.CDNName(toData.Server.CDNName),
		CDNID:                         toData.Server.CDNID,
		DomainName:                    toData.Server.DomainName,
		HostName:                      toData.Server.HostName,
		ID:                            toData.Server.ID,
		IP:                            toData.Server.IPAddress,
		ParentCacheGroupID:            parentCGID,
		ParentCacheGroupType:          parentCGType,
		ProfileID:                     atscfg.ProfileID(toData.Server.ProfileID),
		ProfileName:                   toData.Server.Profile,
		Port:                          toData.Server.TCPPort,
		SecondaryParentCacheGroupID:   secondaryParentCGID,
		SecondaryParentCacheGroupType: secondaryParentCGType,
		Type:                          toData.Server.Type,
	}

	toReverseProxyURL := ""
	toURL := ""
	for _, param := range toData.GlobalParams {
		if param.Name == "tm.rev_proxy.url" {
			toReverseProxyURL = param.Value
		} else if param.Name == "tm.url" {
			toURL = param.Value
		}
		if toReverseProxyURL != "" && toURL != "" {
			break
		}
	}

	scopeParams := ParamsToMap(toData.ScopeParams)

	locationParams := map[string]atscfg.ConfigProfileParams{}
	for _, param := range toData.ServerParams {
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

	dsNames := map[tc.DeliveryServiceName]struct{}{}
	if tc.CacheTypeFromString(toData.Server.Type) != tc.CacheTypeMid {
		dsIDs := map[int]struct{}{}
		for _, ds := range toData.DeliveryServices {
			if ds.ID == nil {
				// TODO log error?
				continue
			}
			dsIDs[*ds.ID] = struct{}{}
		}

		// TODO verify?
		//		serverIDs := []int{toData.Server.ID}

		dssMap := map[int]struct{}{}
		for _, dss := range toData.DeliveryServiceServers {
			if dss.Server == nil || dss.DeliveryService == nil {
				continue // TODO warn?
			}
			if *dss.Server != toData.Server.ID {
				continue
			}
			if _, ok := dsIDs[*dss.DeliveryService]; !ok {
				continue
			}
			dssMap[*dss.DeliveryService] = struct{}{}
		}

		for _, ds := range toData.DeliveryServices {
			if ds.ID == nil {
				continue
			}
			if ds.XMLID == nil {
				continue // TODO log?
			}
			if _, ok := dssMap[*ds.ID]; !ok {
				continue
			}
			dsNames[tc.DeliveryServiceName(*ds.XMLID)] = struct{}{}
		}
	} else {
		for _, ds := range toData.DeliveryServices {
			if ds.ID == nil {
				continue
			}
			if ds.XMLID == nil {
				continue // TODO log?
			}
			if ds.CDNID == nil || *ds.CDNID != toData.Server.CDNID {
				continue
			}
			dsNames[tc.DeliveryServiceName(*ds.XMLID)] = struct{}{}
		}
	}

	uriSignedDSes := []tc.DeliveryServiceName{}
	for _, ds := range toData.DeliveryServices {
		if ds.ID == nil {
			continue
		}
		if ds.XMLID == nil {
			continue // TODO log?
		}
		if _, ok := dsNames[tc.DeliveryServiceName(*ds.XMLID)]; !ok {
			continue
		}
		if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURISigning {
			continue
		}
		uriSignedDSes = append(uriSignedDSes, tc.DeliveryServiceName(*ds.XMLID))
	}

	metaObj := atscfg.MakeMetaObj(tc.CacheName(toData.Server.HostName), &serverInfo, toURL, toReverseProxyURL, locationParams, uriSignedDSes, scopeParams, dsNames)
	return &metaObj, nil
}
