package toreq

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
	"fmt"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

func GetProfile(cfg config.TCCfg, profileID int) (tc.Profile, error) {
	profile := tc.Profile{}
	err := GetCached(cfg, "profile_"+strconv.Itoa(profileID), &profile, func(obj interface{}) error {
		toProfiles, reqInf, err := (*cfg.TOClient).GetProfileByID(profileID)
		if err != nil {
			return errors.New("getting profile '" + strconv.Itoa(profileID) + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if len(toProfiles) != 1 {
			return errors.New("getting profile '" + strconv.Itoa(profileID) + "'from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': expected 1 Profile, got " + strconv.Itoa(len(toProfiles)))
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

func GetProfileByName(cfg config.TCCfg, profileName string) (tc.Profile, error) {
	profile := tc.Profile{}

	err := GetCached(cfg, "profile_"+profileName, &profile, func(obj interface{}) error {
		toProfiles, reqInf, err := (*cfg.TOClient).GetProfileByName(profileName)
		if err != nil {
			return errors.New("getting profile '" + profileName + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if len(toProfiles) != 1 {
			return errors.New("getting profile '" + profileName + "'from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': expected 1 Profile, got " + strconv.Itoa(len(toProfiles)))
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

func GetProfileParameters(cfg config.TCCfg, profileName string) ([]tc.Parameter, error) {
	profileParameters := []tc.Parameter{}
	err := GetCached(cfg, "profile_"+profileName+"_parameters", &profileParameters, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParametersByProfileName(profileName)
		if err != nil {
			return errors.New("getting profile '" + profileName + "' parameters from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetGlobalParameters(cfg config.TCCfg) ([]tc.Parameter, error) {
	globalParams := []tc.Parameter{}
	err := GetCached(cfg, "profile_global_parameters", &globalParams, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParametersByProfileName(tc.GlobalProfileName)
		if err != nil {
			return errors.New("getting global profile '" + tc.GlobalProfileName + "' parameters from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		return nil
	})
	if err != nil {
		return nil, errors.New("getting global profile '" + tc.GlobalProfileName + "' parameters: " + err.Error())
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

func GetTOToolNameAndURLFromTO(cfg config.TCCfg) (string, string, error) {
	globalParams, err := GetGlobalParameters(cfg)
	if err != nil {
		return "", "", err
	}
	toToolName, toURL := GetTOToolNameAndURL(globalParams)
	return toToolName, toURL, nil
}

func GetServers(cfg config.TCCfg) ([]tc.Server, error) {
	servers := []tc.Server{}
	err := GetCached(cfg, "servers", &servers, func(obj interface{}) error {
		toServers, reqInf, err := (*cfg.TOClient).GetServers()
		if err != nil {
			return errors.New("getting servers from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetServerByHostName(cfg config.TCCfg, serverHostName string) (tc.Server, error) {
	server := tc.Server{}
	err := GetCached(cfg, "server-name-"+serverHostName, &server, func(obj interface{}) error {
		toServers, reqInf, err := (*cfg.TOClient).GetServerByHostName(serverHostName)
		if err != nil {
			return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		} else if len(toServers) < 1 {
			return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': no servers returned")
		}
		server := obj.(*tc.Server)
		*server = toServers[0]
		return nil
	})
	if err != nil {
		return tc.Server{}, errors.New("getting server name '" + serverHostName + "': " + err.Error())
	}
	return server, nil
}

func GetServerByID(cfg config.TCCfg, serverID int) (tc.Server, error) {
	server := tc.Server{}
	err := GetCached(cfg, "server-id-"+strconv.Itoa(serverID), &server, func(obj interface{}) error {
		toServers, reqInf, err := (*cfg.TOClient).GetServerByID(serverID)
		if err != nil {
			return fmt.Errorf("getting server id %v from Traffic Ops '%v': %v", serverID, MaybeIPStr(reqInf.RemoteAddr), err)
		} else if len(toServers) < 1 {
			return fmt.Errorf("getting server id %v from Traffic Ops '%v': %v", serverID, MaybeIPStr(reqInf.RemoteAddr), "no servers returned")
		}
		server := obj.(*tc.Server)
		*server = toServers[0]
		return nil
	})
	if err != nil {
		return tc.Server{}, fmt.Errorf("getting server id %v: %v", serverID, err)
	}
	return server, nil
}

func GetCacheGroups(cfg config.TCCfg) ([]tc.CacheGroupNullable, error) {
	cacheGroups := []tc.CacheGroupNullable{}
	err := GetCached(cfg, "cachegroups", &cacheGroups, func(obj interface{}) error {
		toCacheGroups, reqInf, err := (*cfg.TOClient).GetCacheGroupsNullable()
		if err != nil {
			return errors.New("getting cachegroups from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

// DeliveryServiceServersAlwaysGetAll indicates whether to always get all delivery service servers from Traffic Ops, and cache all in a file (but still return to the caller only the objects they requested).
// This exists and is currently true, because with an ORT run, it's typically more efficient to get them all in a single request, and re-use that cache; than for every config file to get and cache its own unique set.
// If your use case is more efficient to only get the needed objects, for example if you're frequently requesting one file, set this false to get and cache the specific needed delivery services and servers.
const DeliveryServiceServersAlwaysGetAll = true

func GetDeliveryServiceServers(cfg config.TCCfg, dsIDs []int, serverIDs []int) ([]tc.DeliveryServiceServer, error) {
	const sortIDsInHash = true

	serverIDsStr := ""
	dsIDsStr := ""
	dsIDsToFetch := ([]int)(nil)
	sIDsToFetch := ([]int)(nil)
	if !DeliveryServiceServersAlwaysGetAll {
		if len(dsIDs) > 0 {
			dsIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(dsIDs, sortIDsInHash)))
		}
		if len(serverIDs) > 0 {
			serverIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(serverIDs, sortIDsInHash)))
		}
		dsIDsToFetch = dsIDs
		sIDsToFetch = serverIDs
	}

	dsServers := []tc.DeliveryServiceServer{}
	err := GetCached(cfg, "deliveryservice_servers_s"+serverIDsStr+"_d_"+dsIDsStr, &dsServers, func(obj interface{}) error {
		const noLimit = 999999 // TODO add "no limit" param to DSS endpoint
		toDSS, reqInf, err := (*cfg.TOClient).GetDeliveryServiceServersWithLimits(noLimit, dsIDsToFetch, sIDsToFetch)
		if err != nil {
			return errors.New("getting delivery service servers from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dss := obj.(*[]tc.DeliveryServiceServer)
		*dss = toDSS.Response
		return nil
	})
	if err != nil {
		return nil, errors.New("getting delivery service servers: " + err.Error())
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
	// Also, if DeliveryServiceServersAlwaysGetAll, we need to filter here anyway.
	filteredDSServers := []tc.DeliveryServiceServer{}
	for _, dsServer := range dsServers {
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

	return filteredDSServers, nil
}

func GetServerProfileParameters(cfg config.TCCfg, profileName string) ([]tc.Parameter, error) {
	serverProfileParameters := []tc.Parameter{}
	err := GetCached(cfg, "profile_"+profileName+"_parameters", &serverProfileParameters, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParametersByProfileName(profileName)
		if err != nil {
			return errors.New("getting server profile '" + profileName + "' parameters from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetCDNDeliveryServices(cfg config.TCCfg, cdnID int) ([]tc.DeliveryServiceNullable, error) {
	deliveryServices := []tc.DeliveryServiceNullable{}
	err := GetCached(cfg, "cdn_"+strconv.Itoa(cdnID)+"_deliveryservices", &deliveryServices, func(obj interface{}) error {
		toDSes, reqInf, err := (*cfg.TOClient).GetDeliveryServicesByCDNID(cdnID)
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetConfigFileParameters(cfg config.TCCfg, configFile string) ([]tc.Parameter, error) {
	params := []tc.Parameter{}
	err := GetCached(cfg, "config_file_"+configFile+"_parameters", &params, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParameterByConfigFile(configFile)
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetCDN(cfg config.TCCfg, cdnName tc.CDNName) (tc.CDN, error) {
	cdn := tc.CDN{}
	err := GetCached(cfg, "cdn_"+string(cdnName), &cdn, func(obj interface{}) error {
		toCDNs, reqInf, err := (*cfg.TOClient).GetCDNByName(string(cdnName))
		if err != nil {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if len(toCDNs) != 1 {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': expected 1 CDN, got " + strconv.Itoa(len(toCDNs)))
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

func GetCDNByID(cfg config.TCCfg, cdnID int) (tc.CDN, error) {
	cdn := tc.CDN{}
	err := GetCached(cfg, "cdn_id_"+strconv.Itoa(cdnID), &cdn, func(obj interface{}) error {
		toCDNs, reqInf, err := (*cfg.TOClient).GetCDNByID(cdnID)
		if err != nil {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if len(toCDNs) != 1 {
			return errors.New("getting cdn from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': expected 1 CDN, got " + strconv.Itoa(len(toCDNs)))
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

func GetURLSigKeys(cfg config.TCCfg, dsName string) (tc.URLSigKeys, error) {
	keys := tc.URLSigKeys{}
	err := GetCached(cfg, "urlsigkeys_"+string(dsName), &keys, func(obj interface{}) error {
		toKeys, reqInf, err := (*cfg.TOClient).GetDeliveryServiceURLSigKeys(dsName)
		if err != nil {
			return errors.New("getting url sig keys from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetURISigningKeys(cfg config.TCCfg, dsName string) ([]byte, error) {
	keys := []byte{}
	err := GetCached(cfg, "urisigningkeys_"+string(dsName), &keys, func(obj interface{}) error {
		toKeys, reqInf, err := (*cfg.TOClient).GetDeliveryServiceURISigningKeys(dsName)
		if err != nil {
			return errors.New("getting url sig keys from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetParametersByName(cfg config.TCCfg, paramName string) ([]tc.Parameter, error) {
	params := []tc.Parameter{}
	err := GetCached(cfg, "parameters_name_"+paramName, &params, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetParameterByName(paramName)
		if err != nil {
			return errors.New("getting parameters name '" + paramName + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetCacheGroupParameters(cfg config.TCCfg, cacheGroupID int) ([]tc.Parameter, error) {
	params := []tc.Parameter{}
	err := GetCached(cfg, "cachegroup_parameters_id_"+strconv.Itoa(cacheGroupID), &params, func(obj interface{}) error {
		toParams, reqInf, err := (*cfg.TOClient).GetCacheGroupParameters(cacheGroupID)
		if err != nil {
			return errors.New("getting cachegroup parameters id '" + strconv.Itoa(cacheGroupID) + "' from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		for _, cgParam := range toParams {
			*params = append(*params, tc.Parameter{
				ConfigFile:  cgParam.ConfigFile,
				ID:          cgParam.ID,
				LastUpdated: cgParam.LastUpdated,
				Name:        cgParam.Name,
				Secure:      cgParam.Secure,
				Value:       cgParam.Value,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("getting params cachegroup id '" + strconv.Itoa(cacheGroupID) + "': " + err.Error())
	}
	return params, nil
}

func GetDeliveryServiceRegexes(cfg config.TCCfg) ([]tc.DeliveryServiceRegexes, error) {
	regexes := []tc.DeliveryServiceRegexes{}
	err := GetCached(cfg, "ds_regexes", &regexes, func(obj interface{}) error {
		toRegexes, reqInf, err := (*cfg.TOClient).GetDeliveryServiceRegexes()
		if err != nil {
			return errors.New("getting ds regexes from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetJobs(cfg config.TCCfg) ([]tc.Job, error) {
	jobs := []tc.Job{}
	err := GetCached(cfg, "jobs", &jobs, func(obj interface{}) error {
		toJobs, reqInf, err := (*cfg.TOClient).GetJobs(nil, nil)
		if err != nil {
			return errors.New("getting jobs from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetServerCapabilitiesByID(cfg config.TCCfg, serverIDs []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	serverIDsStr := ""
	if len(serverIDs) > 0 {
		sortIDsInHash := true
		serverIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(serverIDs, sortIDsInHash)))
	}

	serverCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	err := GetCached(cfg, "server_capabilities_s_"+serverIDsStr, &serverCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toServerCaps, reqInf, err := (*cfg.TOClient).GetServerServerCapabilities(nil, nil, nil)
		if err != nil {
			return errors.New("getting server caps from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetDeliveryServiceRequiredCapabilitiesByID(cfg config.TCCfg, dsIDs []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	dsIDsStr := ""
	if len(dsIDs) > 0 {
		sortIDsInHash := true
		dsIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(dsIDs, sortIDsInHash)))
	}

	dsCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	err := GetCached(cfg, "ds_capabilities_d_"+dsIDsStr, &dsCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toDSCaps, reqInf, err := (*cfg.TOClient).GetDeliveryServicesRequiredCapabilities(nil, nil, nil)
		if err != nil {
			return errors.New("getting ds caps from Traffic Ops '" + MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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

func GetProfileNameFromProfileNameOrID(cfg config.TCCfg, profileNameOrID string) (string, error) {
	profileName := profileNameOrID
	if profileID, err := strconv.Atoi(profileNameOrID); err == nil {
		profile, err := GetProfile(cfg, profileID)
		if err != nil {
			return "", errors.New("getting profile '" + profileNameOrID + "': " + err.Error())
		}
		if profile.Name == "" {
			return "", errors.New("getting profile '" + profileNameOrID + "': got profile with empty name")
		}
		profileName = profile.Name
	}
	return profileName, nil
}

func GetCDNNameFromCDNNameOrID(cfg config.TCCfg, cdnNameOrID string) (tc.CDNName, error) {
	cdnName := cdnNameOrID
	if cdnID, err := strconv.Atoi(cdnNameOrID); err == nil {
		cdn, err := GetCDNByID(cfg, cdnID)
		if err != nil {
			return "", errors.New("getting cdn '" + cdnNameOrID + "': " + err.Error())
		}
		if cdn.Name == "" {
			return "", errors.New("getting cdn '" + cdnNameOrID + "': got cdn with empty name")
		}
		cdnName = cdn.Name
	}
	return tc.CDNName(cdnName), nil
}
