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
	"sort"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func GetConfigFileServerRemapDotConfig(toData *TOData) (string, string, error) {
	// TODO TOAPI add /servers?cdn=1 query param

	atsVersionParam := ""
	for _, param := range toData.ServerParams {
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
		return "", "", errors.New("getting ATS major version from version parameter (profile '" + toData.Server.Profile + "' configFile 'package' name 'trafficserver'): " + err.Error())
	}

	dsIDs := map[int]struct{}{}
	for _, ds := range toData.DeliveryServices {
		if ds.ID == nil {
			// TODO log error?
			continue
		}
		dsIDs[*ds.ID] = struct{}{}
	}

	isMid := strings.HasPrefix(toData.Server.Type, string(tc.CacheTypeMid))

	serverIDs := map[int]struct{}{}
	if !isMid {
		// mids use all servers, so pass empty=all. Edges only use this current server
		serverIDs[toData.Server.ID] = struct{}{}
	}

	dsServers := FilterDSS(toData.DeliveryServiceServers, dsIDs, serverIDs)

	dssMap := map[int]map[int]struct{}{} // set of map[dsID][serverID]
	for _, dss := range dsServers {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO log?
		}
		if dssMap[*dss.DeliveryService] == nil {
			dssMap[*dss.DeliveryService] = map[int]struct{}{}
		}
		dssMap[*dss.DeliveryService][*dss.Server] = struct{}{}
	}

	useInactive := false
	if !isMid {
		// mids get inactive DSes, edges don't. This is how it's always behaved, not necessarily how it should.
		useInactive = true
	}

	filteredDSes := []tc.DeliveryServiceNullable{}
	for _, ds := range toData.DeliveryServices {
		if ds.ID == nil {
			continue // TODO log?
		}
		if ds.Active == nil {
			continue // TODO log?
		}
		if _, ok := dssMap[*ds.ID]; !ok {
			continue
		}
		if !useInactive && !*ds.Active {
			continue
		}
		filteredDSes = append(filteredDSes, ds)
	}

	dsRegexMap := map[tc.DeliveryServiceName][]tc.DeliveryServiceRegex{}
	for _, dsRegex := range toData.DeliveryServiceRegexes {
		sort.Sort(DeliveryServiceRegexesSortByTypeThenSetNum(dsRegex.Regexes))
		dsRegexMap[tc.DeliveryServiceName(dsRegex.DSName)] = dsRegex.Regexes
	}

	remapConfigDSData := []atscfg.RemapConfigDSData{}
	for _, ds := range filteredDSes {
		if ds.ID == nil || ds.Type == nil || ds.XMLID == nil || ds.DSCP == nil || ds.Active == nil {
			continue // TODO log error?
		}
		// TODO sort by DS ID? the old Perl query does, but it shouldn't be necessary, except for determinism.
		// TODO warn if no regexes?
		for _, dsRegex := range dsRegexMap[tc.DeliveryServiceName(*ds.XMLID)] {
			remapConfigDSData = append(remapConfigDSData, atscfg.RemapConfigDSData{
				ID:                       *ds.ID,
				Type:                     *ds.Type,
				OriginFQDN:               ds.OrgServerFQDN,
				MidHeaderRewrite:         ds.MidHeaderRewrite,
				CacheURL:                 ds.CacheURL,
				RangeRequestHandling:     ds.RangeRequestHandling,
				RemapText:                ds.RemapText,
				EdgeHeaderRewrite:        ds.EdgeHeaderRewrite,
				SigningAlgorithm:         ds.SigningAlgorithm,
				Name:                     *ds.XMLID,
				QStringIgnore:            ds.QStringIgnore,
				RegexRemap:               ds.RegexRemap,
				FQPacingRate:             ds.FQPacingRate,
				DSCP:                     *ds.DSCP,
				RoutingName:              ds.RoutingName,
				Pattern:                  util.StrPtr(dsRegex.Pattern),
				RegexType:                util.StrPtr(dsRegex.Type),
				Domain:                   util.StrPtr(toData.CDN.DomainName), // note this is intentionally the CDN domain, not the DS or Server Domain. Must be the remap domain.
				OriginShield:             ds.OriginShield,
				ProfileID:                ds.ProfileID,
				Protocol:                 ds.Protocol,
				AnonymousBlockingEnabled: ds.AnonymousBlockingEnabled,
				Active:                   *ds.Active,
			})
		}
	}

	serverPackageParamData := map[string]string{}
	for _, param := range toData.ServerParams {
		if param.ConfigFile != "package" { // TODO put in const
			continue
		}

		if param.Name == "location" { // TODO put in const
			continue
		}

		paramName := param.Name
		// some files have multiple lines with the same key... handle that with param id.
		if _, ok := serverPackageParamData[param.Name]; ok {
			paramName += "__" + strconv.Itoa(param.ID)
		}
		paramValue := param.Value
		if paramValue == "STRING __HOSTNAME__" {
			paramValue = toData.Server.HostName + "." + toData.Server.DomainName // TODO strings.Replace to replace all anywhere, instead of just an exact match?
		}
		serverPackageParamData[paramName] = paramValue
	}

	cacheURLParams := map[string]string{}
	for _, param := range toData.ServerParams {
		if param.ConfigFile != atscfg.CacheURLParameterConfigFile {
			continue
		}
		if existingVal, ok := cacheURLParams[param.Name]; ok {
			log.Warnln("generating remap.config: server profile '" + toData.Server.Profile + "' cacheurl.config has multiple parameters for '" + param.Name + "' - using '" + existingVal + "' and ignoring the rest!")
			continue
		}
		cacheURLParams[param.Name] = param.Value
	}

	cacheKeyParamsWithProfiles, err := TCParamsToParamsWithProfiles(toData.CacheKeyParams)
	if err != nil {
		return "", "", errors.New("decoding cache key parameter profiles: " + err.Error())
	}

	cacheKeyParamsWithProfilesMap := ParameterWithProfilesToMap(cacheKeyParamsWithProfiles)

	dsProfileNamesToIDs := map[string]int{}
	for _, ds := range filteredDSes {
		if ds.ProfileID == nil || ds.ProfileName == nil {
			continue // TODO log
		}
		dsProfileNamesToIDs[*ds.ProfileName] = *ds.ProfileID
	}

	dsProfilesCacheKeyConfigParams := map[int]map[string]string{}
	for _, param := range cacheKeyParamsWithProfilesMap {
		for dsProfileName, dsProfileID := range dsProfileNamesToIDs {
			if _, ok := param.ProfileNames[dsProfileName]; ok {
				if _, ok := dsProfilesCacheKeyConfigParams[dsProfileID]; !ok {
					dsProfilesCacheKeyConfigParams[dsProfileID] = map[string]string{}
				}
				if _, ok := dsProfilesCacheKeyConfigParams[dsProfileID][param.Name]; ok {
					// TODO warn
					continue
				}
				dsProfilesCacheKeyConfigParams[dsProfileID][param.Name] = param.Value
			}
		}
	}

	// TODO get dses first, so we can get the profile names-to-IDs without fetching all profiles

	// TODO put parentcg logic in func, to remove duplication with parent.config

	cgMap := map[string]tc.CacheGroupNullable{}
	for _, cg := range toData.CacheGroups {
		if cg.Name == nil {
			return "", "", errors.New("got cachegroup with nil name!'")
		}
		cgMap[*cg.Name] = cg
	}

	serverCG, ok := cgMap[toData.Server.Cachegroup]
	if !ok {
		return "", "", errors.New("server '" + toData.Server.HostName + "' cachegroup '" + toData.Server.Cachegroup + "' not found in CacheGroups")
	}

	parentCGID := -1
	parentCGType := ""
	if serverCG.ParentName != nil && *serverCG.ParentName != "" {
		parentCG, ok := cgMap[*serverCG.ParentName]
		if !ok {
			return "", "", errors.New("server '" + toData.Server.HostName + "' cachegroup '" + toData.Server.Cachegroup + "' parent '" + *serverCG.ParentName + "' not found in CacheGroups")
		}
		if parentCG.ID == nil {
			return "", "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		parentCGID = *parentCG.ID

		if parentCG.Type == nil {
			return "", "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}
		parentCGType = *parentCG.Type
	}

	secondaryParentCGID := -1
	secondaryParentCGType := ""
	if serverCG.SecondaryParentName != nil && *serverCG.SecondaryParentName != "" {
		parentCG, ok := cgMap[*serverCG.SecondaryParentName]
		if !ok {
			return "", "", errors.New("server '" + toData.Server.HostName + "' cachegroup '" + toData.Server.Cachegroup + "' secondary parent '" + *serverCG.SecondaryParentName + "' not found in CacheGroups")
		}

		if parentCG.ID == nil {
			return "", "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil ID!'")
		}
		secondaryParentCGID = *parentCG.ID
		if parentCG.Type == nil {
			return "", "", errors.New("got cachegroup '" + *parentCG.Name + "' with nil Type!'")
		}

		secondaryParentCGType = *parentCG.Type
	}

	serverInfo := &atscfg.ServerInfo{
		CacheGroupID:                  toData.Server.CachegroupID,
		CDN:                           tc.CDNName(toData.Server.CDNName),
		CDNID:                         toData.Server.CDNID,
		DomainName:                    toData.CDN.DomainName, // note this is intentionally the CDN domain, not the server domain. It's what's remapped to.
		HostName:                      toData.Server.HostName,
		ID:                            toData.Server.ID,
		IP:                            toData.Server.IPAddress,
		ParentCacheGroupID:            parentCGID,
		ParentCacheGroupType:          parentCGType,
		ProfileID:                     atscfg.ProfileID(toData.Server.ProfileID),
		ProfileName:                   toData.Server.Profile,
		Port:                          toData.Server.TCPPort,
		HTTPSPort:                     toData.Server.HTTPSPort,
		SecondaryParentCacheGroupID:   secondaryParentCGID,
		SecondaryParentCacheGroupType: secondaryParentCGType,
		Type:                          toData.Server.Type,
	}
	return atscfg.MakeRemapDotConfig(tc.CacheName(toData.Server.HostName), toData.TOToolName, toData.TOURL, atsMajorVer, cacheURLParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapConfigDSData), atscfg.ContentTypeRemapDotConfig, nil
}

type DeliveryServiceRegexesSortByTypeThenSetNum []tc.DeliveryServiceRegex

func (r DeliveryServiceRegexesSortByTypeThenSetNum) Len() int { return len(r) }
func (r DeliveryServiceRegexesSortByTypeThenSetNum) Less(i, j int) bool {
	if rc := strings.Compare(r[i].Type, r[j].Type); rc != 0 {
		return rc < 0
	}
	return r[i].SetNumber < r[j].SetNumber
}
func (r DeliveryServiceRegexesSortByTypeThenSetNum) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
