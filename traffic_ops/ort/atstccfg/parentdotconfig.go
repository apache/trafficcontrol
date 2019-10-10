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
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const GlobalProfileName = "GLOBAL"

func GetConfigFileServerParentDotConfig(cfg TCCfg, serverNameOrID string) (string, error) {
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

	parentCacheGroups := map[string]struct{}{}
	if serverInfo.IsTopLevelCache() {
		log.Infoln("This cache Is Top Level!")
		for _, cg := range cacheGroups {
			if cg.Type == nil {
				return "", errors.New("cachegroup type is nil!")
			}
			if cg.Name == nil {
				return "", errors.New("cachegroup type is nil!")
			}

			if *cg.Type != tc.CacheGroupOriginTypeName {
				continue
			}
			parentCacheGroups[*cg.Name] = struct{}{}
		}
	} else {
		if server.Cachegroup == "" {
			return "", errors.New("server cachegroup is nil!")
		}
		for _, cg := range cacheGroups {
			if cg.Type == nil {
				return "", errors.New("cachegroup type is nil!")
			}
			if cg.Name == nil {
				return "", errors.New("cachegroup type is nil!")
			}

			if *cg.Name == server.Cachegroup {
				if cg.ParentName != nil && *cg.ParentName != "" {
					parentCacheGroups[*cg.ParentName] = struct{}{}
				}
				if cg.SecondaryParentName != nil && *cg.SecondaryParentName != "" {
					parentCacheGroups[*cg.SecondaryParentName] = struct{}{}
				}
				break
			}
		}
	}

	cgServers := map[int]tc.Server{} // map[serverID]server
	for _, sv := range servers {
		if sv.CDNName != server.CDNName {
			continue
		}
		if _, ok := parentCacheGroups[sv.Cachegroup]; !ok {
			continue
		}
		if sv.Type != tc.OriginTypeName &&
			!strings.HasPrefix(sv.Type, tc.EdgeTypePrefix) &&
			!strings.HasPrefix(sv.Type, tc.MidTypePrefix) {
			continue
		}
		if sv.Status != string(tc.CacheStatusReported) && sv.Status != string(tc.CacheStatusOnline) {
			continue
		}
		cgServers[sv.ID] = sv
	}

	cgServerIDs := []int{}
	for serverID, _ := range cgServers {
		cgServerIDs = append(cgServerIDs, serverID)
	}
	cgServerIDs = append(cgServerIDs, server.ID)

	serverCapabilities, err := GetServerCapabilitiesByID(cfg, cgServerIDs)
	if err != nil {
		return "", errors.New("getting server capabilities: " + err.Error())
	}

	cgDSServers, err := GetDeliveryServiceServers(cfg, nil, cgServerIDs)
	if err != nil {
		return "", errors.New("getting parent.config cachegroup parent server delivery service servers: " + err.Error())
	}

	parentServerDSes := map[int]map[int]struct{}{} // map[serverID][dsID] // cgServerDSes
	for _, dss := range cgDSServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			return "", errors.New("getting parent.config cachegroup parent server delivery service servers: got dss with nil members!" + err.Error())
		}
		if parentServerDSes[*dss.Server] == nil {
			parentServerDSes[*dss.Server] = map[int]struct{}{}
		}
		parentServerDSes[*dss.Server][*dss.DeliveryService] = struct{}{}
	}

	serverProfileParameters, err := GetServerProfileParameters(cfg, server.Profile)
	if err != nil {
		return "", errors.New("getting server profile '" + server.Profile + "' parameters: " + err.Error())
	}

	atsVersionParam := ""
	for _, param := range serverProfileParameters {
		if param.ConfigFile != "package" || param.Name != "trafficserver" {
			continue
		}
		atsVersionParam = param.Value
		break
	}
	if atsVersionParam == "" {
		atsVersionParam = atscfg.DefaultATSVersion
	}

	atsMajorVer, err := atscfg.GetATSMajorVersionFromATSVersion(atsVersionParam)
	if err != nil {
		return "", errors.New("getting ATS major version from version parameter (profile '" + server.Profile + "' configFile 'package' name 'trafficserver'): " + err.Error())
	}

	globalParams, err := GetGlobalParameters(cfg)
	if err != nil {
		return "", errors.New("getting global parameters: " + err.Error())
	}

	toToolName := ""
	toURL := ""
	for _, param := range globalParams {
		if param.Name == "tm.toolname" {
			toToolName = param.Value
		} else if param.Name == "tm.url" {
			toURL = param.Value
		}
		if toToolName != "" && toURL != "" {
			break
		}
	}

	deliveryServices, err := GetCDNDeliveryServices(cfg, server.CDNID)
	if err != nil {
		return "", errors.New("getting delivery services: " + err.Error())
	}

	parentConfigParams, err := GetConfigFileParameters(cfg, "parent.config")
	if err != nil {
		return "", errors.New("getting parent.config parameters: " + err.Error())
	}

	parentConfigParamsWithProfiles, err := TCParamsToParamsWithProfiles(parentConfigParams)
	if err != nil {
		return "", errors.New("unmarshalling parent.config parameters profiles: " + err.Error())
	}

	// this is an optimization, to avoid looping over all params, for every DS. Instead, we loop over all params only once, and put them in a profile map.
	profileParentConfigParams := map[string]map[string]string{} // map[profileName][paramName]paramVal
	for _, param := range parentConfigParamsWithProfiles {
		for _, profile := range param.ProfileNames {
			if _, ok := profileParentConfigParams[profile]; !ok {
				profileParentConfigParams[profile] = map[string]string{}
			}
			profileParentConfigParams[profile][param.Name] = param.Value
		}
	}

	serverParams := map[string]string{}
	if serverInfo.ProfileName != "" { // TODO warn/error if false? Servers requires profiles.
		for name, val := range profileParentConfigParams[serverInfo.ProfileName] {
			if name == atscfg.ParentConfigParamQStringHandling ||
				name == atscfg.ParentConfigParamAlgorithm ||
				name == atscfg.ParentConfigParamQString {
				serverParams[name] = val
			}
		}
	}

	cdn, err := GetCDN(cfg, serverInfo.CDN)
	if err != nil {
		return "", errors.New("getting cdn '" + string(serverInfo.CDN) + "': " + err.Error())
	}

	serverCDNDomain := cdn.DomainName

	parentConfigServerCacheProfileParams := map[string]atscfg.ProfileCache{} // map[profileName]ProfileCache
	for _, cgServer := range cgServers {
		profileCache, ok := parentConfigServerCacheProfileParams[cgServer.Profile]
		if !ok {
			profileCache = atscfg.DefaultProfileCache()
		}
		params, ok := profileParentConfigParams[cgServer.Profile]
		if !ok {
			parentConfigServerCacheProfileParams[cgServer.Profile] = profileCache
			continue
		}
		for name, val := range params {
			switch name {
			case atscfg.ParentConfigCacheParamWeight:
				// f, err := strconv.ParseFloat(param.Val, 64)
				// if err != nil {
				// 	log.Errorln("parent.config generation: weight param is not a float, skipping! : " + err.Error())
				// } else {
				// 	profileCache.Weight = f
				// }
				// TODO validate float?
				profileCache.Weight = val
			case atscfg.ParentConfigCacheParamPort:
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Errorln("parent.config generation: port param is not an integer, skipping! : " + err.Error())
				} else {
					profileCache.Port = int(i)
				}
			case atscfg.ParentConfigCacheParamUseIP:
				profileCache.UseIP = val == "1"
			case atscfg.ParentConfigCacheParamRank:
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Errorln("parent.config generation: rank param is not an integer, skipping! : " + err.Error())
				} else {
					profileCache.Rank = int(i)
				}
			case atscfg.ParentConfigCacheParamNotAParent:
				profileCache.NotAParent = val != "false"
			}
		}
		parentConfigServerCacheProfileParams[cgServer.Profile] = profileCache
	}

	dsIDMap := map[int]tc.DeliveryServiceNullable{}
	for _, ds := range deliveryServices {
		if ds.ID == nil {
			log.Errorln("delivery services got nil ID!")
			os.Exit(1)
		}

		if !ds.Type.IsHTTP() && !ds.Type.IsDNS() {
			continue // skip ANY_MAP, STEERING, etc
		}

		dsIDMap[*ds.ID] = ds
	}

	allDSMap := map[int]tc.DeliveryServiceNullable{} // all DSes for this server, NOT all dses in TO
	for _, dsIDs := range parentServerDSes {
		for dsID, _ := range dsIDs {
			if _, ok := dsIDMap[dsID]; !ok {
				// this is normal if the TO was too old to understand our /deliveryserviceserver?servers= query param
				// In which case, the DSS will include DSes from other CDNs, which aren't in the dsIDMap
				// If the server was new enough to respect the params, this should never happen.
				// log.Warnln("getting delivery services: parent server DS %v not in dsIDMap\n", dsID)
				continue
			}
			if _, ok := allDSMap[dsID]; !ok {
				allDSMap[dsID] = dsIDMap[dsID]
			}
		}
	}

	allDSes := []int{}
	for ds, _ := range allDSMap {
		allDSes = append(allDSes, int(ds))
	}

	dsRequiredCapabilities, err := GetDeliveryServiceRequiredCapabilitiesByID(cfg, allDSes)
	if err != nil {
		return "", errors.New("getting DS required capabilities: " + err.Error())
	}

	parentConfigDSes := []atscfg.ParentConfigDSTopLevel{}
	for _, tcDS := range deliveryServices {
		if tcDS.ID == nil {
			continue // TODO warn?
		}

		if !serverInfo.IsTopLevelCache() {
			if _, ok := parentServerDSes[server.ID][*tcDS.ID]; !ok {
				continue // skip DSes not assigned to this server.
			}
		}

		if !tcDS.Type.IsHTTP() && !tcDS.Type.IsDNS() {
			continue // skip ANY_MAP, STEERING, etc
		}
		if tcDS.XMLID == nil || *tcDS.XMLID == "" {
			log.Errorln("got delivery service with no XMLID! Skipping!")
			continue
		}
		if tcDS.OrgServerFQDN == nil || *tcDS.OrgServerFQDN == "" {
			log.Errorln("ds  '" + *tcDS.XMLID + "' has no origin server! Skipping!")
			continue
		}

		xmlID := tc.DeliveryServiceName(*tcDS.XMLID)
		originFQDN := *tcDS.OrgServerFQDN
		qStringIgnore := 0
		multiSiteOrigin := false
		originShield := ""
		dsType := tc.DSTypeFromString("")
		if tcDS.QStringIgnore != nil {
			qStringIgnore = *tcDS.QStringIgnore
		}
		if tcDS.MultiSiteOrigin != nil {
			multiSiteOrigin = *tcDS.MultiSiteOrigin
		}
		if tcDS.OriginShield != nil {
			originShield = *tcDS.OriginShield
		}
		if tcDS.Type != nil {
			dsType = *tcDS.Type
		}

		ds := atscfg.ParentConfigDSTopLevel{
			ParentConfigDS: atscfg.ParentConfigDS{
				Name:            xmlID,
				QStringIgnore:   tc.QStringIgnore(qStringIgnore),
				OriginFQDN:      originFQDN,
				MultiSiteOrigin: multiSiteOrigin,
				OriginShield:    originShield,
				Type:            dsType,
			},
		}

		ds.MSOAlgorithm = atscfg.ParentConfigDSParamDefaultMSOAlgorithm
		ds.MSOParentRetry = atscfg.ParentConfigDSParamDefaultMSOParentRetry
		ds.MSOUnavailableServerRetryResponses = atscfg.ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses
		ds.MSOMaxSimpleRetries = atscfg.ParentConfigDSParamDefaultMaxSimpleRetries
		ds.MSOMaxUnavailableServerRetries = atscfg.ParentConfigDSParamDefaultMaxUnavailableServerRetries

		if tcDS.ProfileName != nil && *tcDS.ProfileName != "" {
			if dsParams, ok := profileParentConfigParams[*tcDS.ProfileName]; ok {
				ds.QStringHandling = dsParams[atscfg.ParentConfigParamQStringHandling] // may be blank, no default
				if v, ok := dsParams[atscfg.ParentConfigParamMSOAlgorithm]; ok && strings.TrimSpace(v) != "" {
					ds.MSOAlgorithm = v
				}
				if v, ok := dsParams[atscfg.ParentConfigParamMSOParentRetry]; ok {
					ds.MSOParentRetry = v
				}
				if v, ok := dsParams[atscfg.ParentConfigParamUnavailableServerRetryResponses]; ok {
					ds.MSOUnavailableServerRetryResponses = v
				}
				if v, ok := dsParams[atscfg.ParentConfigParamMaxSimpleRetries]; ok {
					ds.MSOMaxSimpleRetries = v
				}
				if v, ok := dsParams[atscfg.ParentConfigParamMaxUnavailableServerRetries]; ok {
					ds.MSOMaxUnavailableServerRetries = v
				}
			}
		}

		ds.RequiredCapabilities = dsRequiredCapabilities[*tcDS.ID]

		parentConfigDSes = append(parentConfigDSes, ds)
	}

	log.Infof("len(parentServerDSes) %v!\n", len(parentServerDSes))
	log.Infof("len(dsIDMap) %v!\n", len(dsIDMap))
	log.Infof("len(allDSMap) %v!\n", len(allDSMap))

	dsOrigins, err := GetDSOrigins(allDSMap)
	if err != nil {
		log.Errorln("getting delivery service origins: " + err.Error())
		os.Exit(1)
	}

	log.Infof("len(dsOrigins) %v!\n", len(dsOrigins))

	profileParams := parentConfigServerCacheProfileParams

	originServers := map[atscfg.OriginHost][]atscfg.CGServer{}  // "deliveryServices" in Perl
	profileCaches := map[atscfg.ProfileID]atscfg.ProfileCache{} // map[profileID]ProfileCache

	for _, cgServer := range cgServers {
		realCGServer := atscfg.CGServer{
			ServerID:     atscfg.ServerID(cgServer.ID),
			ServerHost:   cgServer.HostName,
			ServerIP:     cgServer.IPAddress,
			ServerPort:   cgServer.TCPPort,
			CacheGroupID: cgServer.CachegroupID,
			Status:       cgServer.StatusID,
			Type:         cgServer.TypeID,
			ProfileID:    atscfg.ProfileID(cgServer.ProfileID),
			CDN:          cgServer.CDNID,
			TypeName:     cgServer.Type,
			Domain:       cgServer.DomainName,
			Capabilities: serverCapabilities[cgServer.ID],
		}

		if cgServer.Type == tc.OriginTypeName {
			for dsID, _ := range parentServerDSes[cgServer.ID] { // map[serverID][]dsID
				orgURI := dsOrigins[dsID]
				if orgURI == nil {
					// log.Warnln("ds %v has no origins! Skipping!\n", dsID) // TODO determine if this is normal
					continue
				}
				if atscfg.HasRequiredCapabilities(serverCapabilities[int(cgServer.ID)], dsRequiredCapabilities[dsID]) {
					orgHost := atscfg.OriginHost(orgURI.Host)
					originServers[orgHost] = append(originServers[orgHost], realCGServer)
				} else {
					log.Errorf("ds %v server %v missing required caps, skipping!\n", dsID, orgURI.Host)
				}
			}
		} else {
			originServers[atscfg.DeliveryServicesAllParentsKey] = append(originServers[atscfg.DeliveryServicesAllParentsKey], realCGServer)
		}

		if _, profileCachesHasProfile := profileCaches[realCGServer.ProfileID]; !profileCachesHasProfile {
			if profileCache, profileParamsHasProfile := profileParams[cgServer.Profile]; !profileParamsHasProfile {
				log.Warnf("cachegroup has server with profile %+v but that profile has no parameters\n", cgServer.ProfileID)
				profileCaches[realCGServer.ProfileID] = atscfg.DefaultProfileCache()
			} else {
				profileCaches[realCGServer.ProfileID] = profileCache
			}
		}
	}

	parentInfos := atscfg.MakeParentInfo(&serverInfo, serverCDNDomain, profileCaches, originServers)

	return atscfg.MakeParentDotConfig(&serverInfo, atsMajorVer, toToolName, toURL, parentConfigDSes, serverParams, parentInfos), nil
}

// GetDSOrigins takes a map[deliveryServiceID]DeliveryService, and returns a map[DeliveryServiceID]OriginURI.
func GetDSOrigins(dses map[int]tc.DeliveryServiceNullable) (map[int]*atscfg.OriginURI, error) {
	dsOrigins := map[int]*atscfg.OriginURI{}
	for _, ds := range dses {
		if ds.ID == nil {
			return nil, errors.New("ds has nil ID")
		}
		if ds.XMLID == nil {
			return nil, fmt.Errorf("ds id %v has nil XMLID", *ds.ID)
		}
		if ds.OrgServerFQDN == nil {
			log.Warnf("GetDSOrigins ds %v got nil OrgServerFQDN, skipping!\n", *ds.XMLID)
			continue
		}
		orgURL, err := url.Parse(*ds.OrgServerFQDN)
		if err != nil {
			return nil, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + err.Error())
		}
		if orgURL.Scheme == "" {
			return nil, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "missing scheme")
		}
		if orgURL.Host == "" {
			return nil, errors.New("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "missing scheme")
		}

		scheme := orgURL.Scheme
		host := orgURL.Hostname()
		port := orgURL.Port()
		if port == "" {
			if scheme == "http" {
				port = "80"
			} else if scheme == "https" {
				port = "443"
			} else {
				log.Warnln("parsing ds '" + *ds.XMLID + "' OrgServerFQDN '" + *ds.OrgServerFQDN + "': " + "unknown scheme '" + scheme + "' and no port, leaving port empty!")
			}
		}
		dsOrigins[*ds.ID] = &atscfg.OriginURI{Scheme: scheme, Host: host, Port: port}
	}
	return dsOrigins, nil
}
