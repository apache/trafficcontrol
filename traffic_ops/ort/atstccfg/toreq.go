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
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func GetProfile(cfg TCCfg, profileID int) (tc.Profile, error) {
	profile := tc.Profile{}
	err := GetCachedJSON(cfg, "profile_"+strconv.Itoa(profileID)+".json", &profile, func(obj interface{}) error {
		toProfiles, reqInf, err := (*cfg.TOClient).GetProfileByID(profileID)
		if err != nil {
			return errors.New("getting profile '" + strconv.Itoa(profileID) + "' from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		if len(toProfiles) != 1 {
			return errors.New("getting profile '" + strconv.Itoa(profileID) + "'from Traffic Ops '" + MaybeIPStr(reqInf) + "': expected 1 Profile, got " + strconv.Itoa(len(toProfiles)))
		}

		profile := obj.(*tc.Profile)
		*profile = toProfiles[0]
		return nil
	})
	if err != nil {
		return tc.Profile{}, errors.New("getting profile '" + strconv.Itoa(profileID) + "': " + err.Error())
	}
	return profile, nil
}

func GetProfileByName(cfg TCCfg, profileName string) (tc.Profile, error) {
	profile := tc.Profile{}
	err := GetCachedJSON(cfg, "profile_"+profileName+".json", &profile, func(obj interface{}) error {
		toProfiles, reqInf, err := (*cfg.TOClient).GetProfileByName(profileName)
		if err != nil {
			return errors.New("getting profile '" + profileName + "' from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		if len(toProfiles) != 1 {
			return errors.New("getting profile '" + profileName + "'from Traffic Ops '" + MaybeIPStr(reqInf) + "': expected 1 Profile, got " + strconv.Itoa(len(toProfiles)))
		}

		profile := obj.(*tc.Profile)
		*profile = toProfiles[0]
		return nil
	})
	if err != nil {
		return tc.Profile{}, errors.New("getting profile '" + profileName + "': " + err.Error())
	}
	return profile, nil
}

func GetProfileParameters(cfg TCCfg, profileName string) ([]tc.Parameter, error) {
	profileParameters := []tc.Parameter{}
	err := GetCachedJSON(cfg, "profile_"+profileName+"_parameters.json", &profileParameters, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParametersByProfileName(profileName)
		if err != nil {
			return errors.New("getting profile '" + profileName + "' parameters from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		return nil
	})
	if err != nil {
		return nil, errors.New("getting profile '" + profileName + "' parameters: " + err.Error())
	}
	return profileParameters, nil
}

func GetGlobalParameters(cfg TCCfg) ([]tc.Parameter, error) {
	globalParams := []tc.Parameter{}
	err := GetCachedJSON(cfg, "profile_global_parameters.json", &globalParams, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParametersByProfileName(GlobalProfileName)
		if err != nil {
			return errors.New("getting global profile '" + GlobalProfileName + "' parameters from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		return nil
	})
	if err != nil {
		return nil, errors.New("getting global profile '" + GlobalProfileName + "' parameters: " + err.Error())
	}
	return globalParams, nil
}

func GetTOToolNameAndURL(globalParams []tc.Parameter) (string, string) {
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
	// TODO error here? Perl doesn't.
	if toToolName == "" {
		log.Warnln("Global Parameter tm.toolname not found, config may not be constructed properly!")
	}
	if toURL == "" {
		log.Warnln("Global Parameter tm.url not found, config may not be constructed properly!")
	}
	return toToolName, toURL
}

func GetTOToolNameAndURLFromTO(cfg TCCfg) (string, string, error) {
	globalParams, err := GetGlobalParameters(cfg)
	if err != nil {
		return "", "", err
	}
	toToolName, toURL := GetTOToolNameAndURL(globalParams)
	return toToolName, toURL, nil
}

func GetServers(cfg TCCfg) ([]tc.Server, error) {
	servers := []tc.Server{}
	err := GetCachedJSON(cfg, "servers.json", &servers, func(obj interface{}) error {
		toServers, reqInf, err := (*cfg.TOClient).GetServers()
		if err != nil {
			return errors.New("getting servers from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		servers := obj.(*[]tc.Server)
		*servers = toServers
		return nil
	})
	if err != nil {
		return nil, errors.New("getting servers: " + err.Error())
	}
	return servers, nil
}

func GetCacheGroups(cfg TCCfg) ([]tc.CacheGroupNullable, error) {
	cacheGroups := []tc.CacheGroupNullable{}
	err := GetCachedJSON(cfg, "cachegroups.json", &cacheGroups, func(obj interface{}) error {
		toCacheGroups, reqInf, err := (*cfg.TOClient).GetCacheGroupsNullable()
		if err != nil {
			return errors.New("getting cachegroups from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		cacheGroups := obj.(*[]tc.CacheGroupNullable)
		*cacheGroups = toCacheGroups
		return nil
	})
	if err != nil {
		return nil, errors.New("getting cachegroups: " + err.Error())
	}
	return cacheGroups, nil
}

func GetDeliveryServiceServers(cfg TCCfg, dsIDs []int, serverIDs []int) ([]tc.DeliveryServiceServer, error) {
	sortIDsInHash := true
	serverIDsStr := ""
	if len(serverIDs) > 0 {
		serverIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(serverIDs, sortIDsInHash)))
	}
	dsIDsStr := ""
	if len(dsIDs) > 0 {
		dsIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(dsIDs, sortIDsInHash)))
	}

	dsServers := []tc.DeliveryServiceServer{}
	err := GetCachedJSON(cfg, "deliveryservice_servers_s"+serverIDsStr+"_d_"+dsIDsStr+".json", &dsServers, func(obj interface{}) error {
		const noLimit = 999999 // TODO add "no limit" param to DSS endpoint
		toDSS, reqInf, err := (*cfg.TOClient).GetDeliveryServiceServersWithLimits(noLimit, dsIDs, serverIDs)
		if err != nil {
			return errors.New("getting delivery service servers from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}

		serverIDsMap := map[int]struct{}{}
		for _, id := range serverIDs {
			serverIDsMap[id] = struct{}{}
		}

		dsIDsMap := map[int]struct{}{}
		for _, id := range dsIDs {
			dsIDsMap[id] = struct{}{}
		}

		// Older TO's may ignore the server ID list, so we need to filter them out manually to be sure.
		filteredDSServers := []tc.DeliveryServiceServer{}
		for _, dsServer := range toDSS.Response {
			if dsServer.Server == nil || dsServer.DeliveryService == nil {
				continue // TODO warn? error?
			}
			if len(serverIDsMap) > 0 {
				if _, ok := serverIDsMap[*dsServer.Server]; !ok {
					continue
				}
			}
			if len(dsIDsMap) > 0 {
				if _, ok := dsIDsMap[*dsServer.DeliveryService]; !ok {
					continue
				}
			}
			filteredDSServers = append(filteredDSServers, dsServer)
		}

		dss := obj.(*[]tc.DeliveryServiceServer)
		*dss = filteredDSServers
		return nil
	})
	if err != nil {
		return nil, errors.New("getting delivery service servers: " + err.Error())
	}

	return dsServers, nil
}

func GetServerProfileParameters(cfg TCCfg, profileName string) ([]tc.Parameter, error) {
	serverProfileParameters := []tc.Parameter{}
	err := GetCachedJSON(cfg, "profile_"+profileName+"_parameters.json", &serverProfileParameters, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParametersByProfileName(profileName)
		if err != nil {
			return errors.New("getting server profile '" + profileName + "' parameters from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		return nil
	})
	if err != nil {
		return nil, errors.New("getting server profile '" + profileName + "' parameters: " + err.Error())
	}
	return serverProfileParameters, nil
}

func GetCDNDeliveryServices(cfg TCCfg, cdnID int) ([]tc.DeliveryServiceNullable, error) {
	deliveryServices := []tc.DeliveryServiceNullable{}
	err := GetCachedJSON(cfg, "cdn_"+strconv.Itoa(cdnID)+"_deliveryservices"+".json", &deliveryServices, func(obj interface{}) error {
		toDSes, reqInf, err := (*cfg.TOClient).GetDeliveryServicesByCDNID(cdnID)
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		dses := obj.(*[]tc.DeliveryServiceNullable)
		*dses = toDSes
		return nil
	})
	if err != nil {
		return nil, errors.New("getting delivery services: " + err.Error())
	}
	return deliveryServices, nil
}

func GetConfigFileParameters(cfg TCCfg, configFile string) ([]tc.Parameter, error) {
	params := []tc.Parameter{}
	err := GetCachedJSON(cfg, "config_file_"+configFile+"_parameters"+".json", &params, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParameterByConfigFile(configFile)
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		return nil
	})
	if err != nil {
		return nil, errors.New("getting parent.config parameters: " + err.Error())
	}
	return params, nil
}

func GetCDN(cfg TCCfg, cdnName tc.CDNName) (tc.CDN, error) {
	cdn := tc.CDN{}
	err := GetCachedJSON(cfg, "cdn_"+string(cdnName)+".json", &cdn, func(obj interface{}) error {
		toCDNs, reqInf, err := (*cfg.TOClient).GetCDNByName(string(cdnName))
		if err != nil {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		if len(toCDNs) != 1 {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf) + "': expected 1 CDN, got " + strconv.Itoa(len(toCDNs)))
		}
		cdn := obj.(*tc.CDN)
		*cdn = toCDNs[0]
		return nil
	})
	if err != nil {
		return tc.CDN{}, errors.New("getting cdn: " + err.Error())
	}
	return cdn, nil
}

func GetCDNByID(cfg TCCfg, cdnID int) (tc.CDN, error) {
	cdn := tc.CDN{}
	err := GetCachedJSON(cfg, "cdn_id_"+strconv.Itoa(cdnID)+".json", &cdn, func(obj interface{}) error {
		toCDNs, reqInf, err := (*cfg.TOClient).GetCDNByID(cdnID)
		if err != nil {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		if len(toCDNs) != 1 {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf) + "': expected 1 CDN, got " + strconv.Itoa(len(toCDNs)))
		}
		cdn := obj.(*tc.CDN)
		*cdn = toCDNs[0]
		return nil
	})
	if err != nil {
		return tc.CDN{}, errors.New("getting cdn: " + err.Error())
	}
	return cdn, nil
}

func GetURLSigKeys(cfg TCCfg, dsName string) (tc.URLSigKeys, error) {
	keys := tc.URLSigKeys{}
	err := GetCachedJSON(cfg, "urlsigkeys_"+string(dsName)+".json", &keys, func(obj interface{}) error {
		toKeys, reqInf, err := (*cfg.TOClient).GetDeliveryServiceURLSigKeys(dsName)
		if err != nil {
			return errors.New("getting url sig keys from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		keys := obj.(*tc.URLSigKeys)
		*keys = toKeys
		return nil
	})
	if err != nil {
		return tc.URLSigKeys{}, errors.New("getting url sig keys: " + err.Error())
	}
	return keys, nil
}

func GetURISigningKeys(cfg TCCfg, dsName string) ([]byte, error) {
	keys := []byte{}
	err := GetCachedJSON(cfg, "urisigningkeys_"+string(dsName)+".json", &keys, func(obj interface{}) error {
		toKeys, reqInf, err := (*cfg.TOClient).GetDeliveryServiceURISigningKeys(dsName)
		if err != nil {
			return errors.New("getting url sig keys from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}

		keys := obj.(*[]byte)
		*keys = toKeys
		return nil
	})
	if err != nil {
		return []byte{}, errors.New("getting url sig keys: " + err.Error())
	}
	return keys, nil
}

func GetParametersByName(cfg TCCfg, paramName string) ([]tc.Parameter, error) {
	params := []tc.Parameter{}
	err := GetCachedJSON(cfg, "parameters_name_"+paramName+".json", &params, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParameterByName(paramName)
		if err != nil {
			return errors.New("getting parameters name '" + paramName + "' from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		return nil
	})
	if err != nil {
		return nil, errors.New("getting params name '" + paramName + "': " + err.Error())
	}
	return params, nil
}

func GetDeliveryServiceRegexes(cfg TCCfg) ([]tc.DeliveryServiceRegexes, error) {
	regexes := []tc.DeliveryServiceRegexes{}
	err := GetCachedJSON(cfg, "ds_regexes.json", &regexes, func(obj interface{}) error {
		toRegexes, reqInf, err := (*cfg.TOClient).GetDeliveryServiceRegexes()
		if err != nil {
			return errors.New("getting ds regexes from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		regexes := obj.(*[]tc.DeliveryServiceRegexes)
		*regexes = toRegexes
		return nil
	})
	if err != nil {
		return nil, errors.New("getting ds regexes: " + err.Error())
	}
	return regexes, nil
}

func GetJobs(cfg TCCfg) ([]tc.Job, error) {
	jobs := []tc.Job{}
	err := GetCachedJSON(cfg, "jobs.json", &jobs, func(obj interface{}) error {
		toJobs, reqInf, err := (*cfg.TOClient).GetJobs(nil, nil)
		if err != nil {
			return errors.New("getting jobs from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		jobs := obj.(*[]tc.Job)
		*jobs = toJobs
		return nil
	})
	if err != nil {
		return nil, errors.New("getting jobs: " + err.Error())
	}
	return jobs, nil
}
func GetServerCapabilitiesByID(cfg TCCfg, serverIDs []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	serverIDsStr := ""
	if len(serverIDs) > 0 {
		sortIDsInHash := true
		serverIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(serverIDs, sortIDsInHash)))
	}

	serverCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	err := GetCachedJSON(cfg, "server_capabilities_s_"+serverIDsStr+".json", &serverCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toServerCaps, reqInf, err := (*cfg.TOClient).GetServerServerCapabilities(nil, nil, nil)
		if err != nil {
			return errors.New("getting server caps from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		serverCaps := obj.(*map[int]map[atscfg.ServerCapability]struct{})

		for _, sc := range toServerCaps {
			if sc.ServerID == nil {
				log.Errorln("Traffic Ops returned Server Capability with nil server id! Skipping!")
			}
			if sc.ServerCapability == nil {
				log.Errorln("Traffic Ops returned Server Capability with nil capability! Skipping!")
			}
			if _, ok := (*serverCaps)[*sc.ServerID]; !ok {
				(*serverCaps)[*sc.ServerID] = map[atscfg.ServerCapability]struct{}{}
			}
			(*serverCaps)[*sc.ServerID][atscfg.ServerCapability(*sc.ServerCapability)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("getting server server capabilities: " + err.Error())
	}
	return serverCaps, nil
}

func GetDeliveryServiceRequiredCapabilitiesByID(cfg TCCfg, dsIDs []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	dsIDsStr := ""
	if len(dsIDs) > 0 {
		sortIDsInHash := true
		dsIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(dsIDs, sortIDsInHash)))
	}

	dsCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	err := GetCachedJSON(cfg, "ds_capabilities_d_"+dsIDsStr+".json", &dsCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toDSCaps, reqInf, err := (*cfg.TOClient).GetDeliveryServicesRequiredCapabilities(nil, nil, nil)
		if err != nil {
			return errors.New("getting ds caps from Traffic Ops '" + MaybeIPStr(reqInf) + "': " + err.Error())
		}
		dsCaps := obj.(*map[int]map[atscfg.ServerCapability]struct{})

		for _, sc := range toDSCaps {
			if sc.DeliveryServiceID == nil {
				log.Errorln("Traffic Ops returned Delivery Service Capability with nil ds id! Skipping!")
			}
			if sc.RequiredCapability == nil {
				log.Errorln("Traffic Ops returned Delivery Service Capability with nil capability! Skipping!")
			}
			if (*dsCaps)[*sc.DeliveryServiceID] == nil {
				(*dsCaps)[*sc.DeliveryServiceID] = map[atscfg.ServerCapability]struct{}{}
			}
			(*dsCaps)[*sc.DeliveryServiceID][atscfg.ServerCapability(*sc.RequiredCapability)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("getting ds server capabilities: " + err.Error())
	}
	return dsCaps, nil
}
