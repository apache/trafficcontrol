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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/cache-config/t3c-generate/torequtil"
	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

func (cl *TOClient) GetProfileByName(profileName string, reqHdr http.Header) (tc.Profile, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetProfileByName(profileName)
	}

	profile := tc.Profile{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "profile_"+profileName, &profile, func(obj interface{}) error {
		toProfiles, toReqInf, err := cl.C.GetProfileByNameWithHdr(profileName, reqHdr)
		if err != nil {
			return errors.New("getting profile '" + profileName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if toReqInf.StatusCode != http.StatusNotModified {
			if len(toProfiles) != 1 {
				return errors.New("getting profile '" + profileName + "'from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': expected 1 Profile, got " + strconv.Itoa(len(toProfiles)))
			}

			profile := obj.(*tc.Profile)
			*profile = toProfiles[0]
		}
		reqInf = toReqInf
		return nil
	})

	if err != nil {
		return tc.Profile{}, reqInf, errors.New("getting profile '" + profileName + "': " + err.Error())
	}
	return profile, reqInf, nil
}

func (cl *TOClient) GetGlobalParameters(reqHdr http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetGlobalParameters()
	}

	globalParams := []tc.Parameter{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "profile_global_parameters", &globalParams, func(obj interface{}) error {
		toParams, toReqInf, err := cl.C.GetParametersByProfileNameWithHdr(tc.GlobalProfileName, reqHdr)
		if err != nil {
			return errors.New("getting global profile '" + tc.GlobalProfileName + "' parameters from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting global profile '" + tc.GlobalProfileName + "' parameters: " + err.Error())
	}
	return globalParams, reqInf, nil
}

func (cl *TOClient) GetServers(reqHdr http.Header) ([]atscfg.Server, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetServers()
	}

	servers := []atscfg.Server{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "servers", &servers, func(obj interface{}) error {
		toServers, toReqInf, err := cl.C.GetServersWithHdr(nil, reqHdr)
		if err != nil {
			return errors.New("getting servers from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		servers := obj.(*[]atscfg.Server)
		*servers, err = serversToLatest(toServers)
		if err != nil {
			return errors.New("upgrading servers from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting servers: " + err.Error())
	}
	return servers, reqInf, nil
}

func (cl *TOClient) GetServerByHostName(serverHostName string, reqHdr http.Header) (*atscfg.Server, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetServerByHostName(serverHostName)
	}

	server := atscfg.Server{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "server-name-"+serverHostName, &server, func(obj interface{}) error {
		params := &url.Values{}
		params.Add("hostName", serverHostName)
		toServers, toReqInf, err := cl.C.GetServersWithHdr(params, reqHdr)
		if err != nil {
			return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if toReqInf.StatusCode != http.StatusNotModified {
			if len(toServers.Response) < 1 {
				return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': no servers returned")
			}
			asv, err := serverToLatest(&toServers.Response[0])
			if err != nil {
				return errors.New("converting server to latest version: " + err.Error())
			}
			server := obj.(*atscfg.Server)
			*server = *asv
		}
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting server name '" + serverHostName + "': " + err.Error())
	}
	return &server, reqInf, nil
}

func (cl *TOClient) GetCacheGroups(reqHdr http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetCacheGroups()
	}

	cacheGroups := []tc.CacheGroupNullable{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cachegroups", &cacheGroups, func(obj interface{}) error {
		toCacheGroups, toReqInf, err := cl.C.GetCacheGroupsNullableWithHdr(reqHdr)
		if err != nil {
			return errors.New("getting cachegroups from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		cacheGroups := obj.(*[]tc.CacheGroupNullable)
		*cacheGroups = toCacheGroups
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting cachegroups: " + err.Error())
	}
	return cacheGroups, reqInf, nil
}

// DeliveryServiceServersAlwaysGetAll indicates whether to always get all delivery service servers from Traffic Ops, and cache all in a file (but still return to the caller only the objects they requested).
// This exists and is currently true, because with an ORT run, it's typically more efficient to get them all in a single request, and re-use that cache; than for every config file to get and cache its own unique set.
// If your use case is more efficient to only get the needed objects, for example if you're frequently requesting one file, set this false to get and cache the specific needed delivery services and servers.
const DeliveryServiceServersAlwaysGetAll = true

func (cl *TOClient) GetDeliveryServiceServers(dsIDs []int, serverIDs []int, reqHdr http.Header) ([]tc.DeliveryServiceServer, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetDeliveryServiceServers(dsIDs, serverIDs)
	}

	const sortIDsInHash = true
	reqInf := toclientlib.ReqInf{}
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
	err := torequtil.GetRetry(cl.NumRetries, "deliveryservice_servers_s"+serverIDsStr+"_d_"+dsIDsStr, &dsServers, func(obj interface{}) error {
		const noLimit = 999999 // TODO add "no limit" param to DSS endpoint
		toDSS, toReqInf, err := cl.C.GetDeliveryServiceServersWithLimitsWithHdr(noLimit, dsIDsToFetch, sIDsToFetch, reqHdr)
		if err != nil {
			return errors.New("getting delivery service servers from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dss := obj.(*[]tc.DeliveryServiceServer)
		*dss = toDSS.Response
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting delivery service servers: " + err.Error())
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

	return filteredDSServers, reqInf, nil
}

func (cl *TOClient) GetServerProfileParameters(profileName string, reqHdr http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetServerProfileParameters(profileName)
	}

	serverProfileParameters := []tc.Parameter{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "profile_"+profileName+"_parameters", &serverProfileParameters, func(obj interface{}) error {
		toParams, toReqInf, err := cl.C.GetParametersByProfileNameWithHdr(profileName, reqHdr)
		if err != nil {
			return errors.New("getting server profile '" + profileName + "' parameters from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting server profile '" + profileName + "' parameters: " + err.Error())
	}
	return serverProfileParameters, reqInf, nil
}

// GetCDNDeliveryServices returns the data, the Traffic Ops address, and any error.
func (cl *TOClient) GetCDNDeliveryServices(cdnID int, reqHdr http.Header) ([]atscfg.DeliveryService, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetCDNDeliveryServices(cdnID)
	}

	deliveryServices := []atscfg.DeliveryService{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cdn_"+strconv.Itoa(cdnID)+"_deliveryservices", &deliveryServices, func(obj interface{}) error {
		params := url.Values{}
		params.Set("cdn", strconv.Itoa(cdnID))
		toDSes, toReqInf, err := cl.C.GetDeliveryServicesV30WithHdr(reqHdr, params)
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dses := obj.(*[]atscfg.DeliveryService)
		*dses = dsesToLatest(toDSes)
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting delivery services: " + err.Error())
	}
	return deliveryServices, reqInf, nil
}

// GetTopologies returns the data, the Traffic Ops address, and any error.
func (cl *TOClient) GetTopologies(reqHdr http.Header) ([]tc.Topology, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetTopologies()
	}

	topologies := []tc.Topology{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "topologies", &topologies, func(obj interface{}) error {
		toTopologies, toReqInf, err := cl.C.GetTopologiesWithHdr(reqHdr)
		if err != nil {
			return errors.New("getting topologies from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		topologies := obj.(*[]tc.Topology)
		*topologies = toTopologies
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting topologies: " + err.Error())
	}
	return topologies, reqInf, nil
}

func (cl *TOClient) GetConfigFileParameters(configFile string, reqHdr http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetConfigFileParameters(configFile)
	}

	params := []tc.Parameter{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "config_file_"+configFile+"_parameters", &params, func(obj interface{}) error {
		toParams, toReqInf, err := cl.C.GetParameterByConfigFileWithHdr(configFile, reqHdr)
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting parent.config parameters: " + err.Error())
	}
	return params, reqInf, nil
}

func (cl *TOClient) GetCDN(cdnName tc.CDNName, reqHdr http.Header) (tc.CDN, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetCDN(cdnName)
	}

	cdn := tc.CDN{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cdn_"+string(cdnName), &cdn, func(obj interface{}) error {
		toCDNs, toReqInf, err := cl.C.GetCDNByNameWithHdr(string(cdnName), reqHdr)
		if err != nil {
			return errors.New("getting cdn from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if toReqInf.StatusCode != http.StatusNotModified {
			if len(toCDNs) != 1 {
				return errors.New("getting cdn from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': expected 1 CDN, got " + strconv.Itoa(len(toCDNs)))
			}
			cdn := obj.(*tc.CDN)
			*cdn = toCDNs[0]
		}
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return tc.CDN{}, reqInf, errors.New("getting cdn: " + err.Error())
	}
	return cdn, reqInf, nil
}

func (cl *TOClient) GetURLSigKeys(dsName string, reqHdr http.Header) (tc.URLSigKeys, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetURLSigKeys(dsName)
	}

	keys := tc.URLSigKeys{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "urlsigkeys_"+string(dsName), &keys, func(obj interface{}) error {
		toKeys, toReqInf, err := cl.C.GetDeliveryServiceURLSigKeysWithHdr(dsName, reqHdr)
		if err != nil {
			return errors.New("getting url sig keys from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		keys := obj.(*tc.URLSigKeys)
		*keys = toKeys
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return tc.URLSigKeys{}, reqInf, errors.New("getting url sig keys: " + err.Error())
	}
	return keys, reqInf, nil
}

func (cl *TOClient) GetURISigningKeys(dsName string, reqHdr http.Header) ([]byte, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetURISigningKeys(dsName)
	}

	keys := []byte{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "urisigningkeys_"+string(dsName), &keys, func(obj interface{}) error {
		toKeys, toReqInf, err := cl.C.GetDeliveryServiceURISigningKeysWithHdr(dsName, reqHdr)
		if err != nil {
			return errors.New("getting url sig keys from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}

		keys := obj.(*[]byte)
		*keys = toKeys
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return []byte{}, reqInf, errors.New("getting url sig keys: " + err.Error())
	}
	return keys, reqInf, nil
}

func (cl *TOClient) GetParametersByName(paramName string, reqHdr http.Header) ([]tc.Parameter, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetParametersByName(paramName)
	}

	params := []tc.Parameter{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "parameters_name_"+paramName, &params, func(obj interface{}) error {
		toParams, toReqInf, err := cl.C.GetParameterByNameWithHdr(paramName, reqHdr)
		if err != nil {
			return errors.New("getting parameters name '" + paramName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.Parameter)
		*params = toParams
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting params name '" + paramName + "': " + err.Error())
	}
	return params, reqInf, nil
}

func (cl *TOClient) GetDeliveryServiceRegexes(reqHdr http.Header) ([]tc.DeliveryServiceRegexes, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetDeliveryServiceRegexes()
	}

	regexes := []tc.DeliveryServiceRegexes{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "ds_regexes", &regexes, func(obj interface{}) error {
		toRegexes, toReqInf, err := cl.C.GetDeliveryServiceRegexesWithHdr(reqHdr)
		if err != nil {
			return errors.New("getting ds regexes from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		regexes := obj.(*[]tc.DeliveryServiceRegexes)
		*regexes = toRegexes
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting ds regexes: " + err.Error())
	}
	return regexes, reqInf, nil
}

func (cl *TOClient) GetJobs(reqHdr http.Header) ([]tc.InvalidationJob, toclientlib.ReqInf, error) {
	if cl.C == nil {
		oldJobs, inf, err := cl.Old.GetJobs()
		jobs, err := atscfg.JobsToInvalidationJobs(oldJobs)
		if err != nil {
			return nil, inf, errors.New("converting old []tc.Job to []tc.InvalidationJob: " + err.Error())
		}
		return jobs, inf, err
	}

	jobs := []tc.InvalidationJob{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "jobs", &jobs, func(obj interface{}) error {
		toJobs, toReqInf, err := cl.C.GetInvalidationJobsWithHdr(nil, nil, reqHdr)
		if err != nil {
			return errors.New("getting jobs from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		jobs := obj.(*[]tc.InvalidationJob)
		*jobs = toJobs
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting jobs: " + err.Error())
	}
	return jobs, reqInf, nil
}

func (cl *TOClient) GetServerCapabilitiesByID(serverIDs []int, reqHdr http.Header) (map[int]map[atscfg.ServerCapability]struct{}, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetServerCapabilitiesByID(serverIDs)
	}

	serverIDsStr := ""
	if len(serverIDs) > 0 {
		sortIDsInHash := true
		serverIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(serverIDs, sortIDsInHash)))
	}

	serverCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "server_capabilities_s_"+serverIDsStr, &serverCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toServerCaps, toReqInf, err := cl.C.GetServerServerCapabilitiesWithHdr(nil, nil, nil, reqHdr)
		if err != nil {
			return errors.New("getting server caps from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting server server capabilities: " + err.Error())
	}
	return serverCaps, reqInf, nil
}

func (cl *TOClient) GetDeliveryServiceRequiredCapabilitiesByID(dsIDs []int, reqHdr http.Header) (map[int]map[atscfg.ServerCapability]struct{}, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetDeliveryServiceRequiredCapabilitiesByID(dsIDs)
	}

	dsIDsStr := ""
	if len(dsIDs) > 0 {
		sortIDsInHash := true
		dsIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(dsIDs, sortIDsInHash)))
	}

	dsCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "ds_capabilities_d_"+dsIDsStr, &dsCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toDSCaps, toReqInf, err := cl.C.GetDeliveryServicesRequiredCapabilitiesWithHdr(nil, nil, nil, reqHdr)
		if err != nil {
			return errors.New("getting ds caps from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
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
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting ds server capabilities: " + err.Error())
	}
	return dsCaps, reqInf, nil
}

func (cl *TOClient) GetCDNSSLKeys(cdnName tc.CDNName, reqHdr http.Header) ([]tc.CDNSSLKeys, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetCDNSSLKeys(cdnName)
	}

	keys := []tc.CDNSSLKeys{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cdn_sslkeys_"+string(cdnName), &keys, func(obj interface{}) error {
		toKeys, toReqInf, err := cl.C.GetCDNSSLKeysWithHdr(string(cdnName), reqHdr)
		if err != nil {
			return errors.New("getting cdn ssl keys from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		keys := obj.(*[]tc.CDNSSLKeys)
		*keys = toKeys
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return []tc.CDNSSLKeys{}, reqInf, errors.New("getting cdn ssl keys: " + err.Error())
	}
	return keys, reqInf, nil
}

func (cl *TOClient) GetStatuses(reqHdr http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetStatuses()
	}

	statuses := []tc.Status{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "statuses", &statuses, func(obj interface{}) error {
		toStatus, toReqInf, err := cl.C.GetStatusesWithHdr(reqHdr)
		if err != nil {
			return errors.New("getting server update status from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		status := obj.(*[]tc.Status)
		*status = toStatus
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting server update status: " + err.Error())
	}
	return statuses, reqInf, nil
}

// GetServerUpdateStatus returns the data, the Traffic Ops address, and any error.
func (cl *TOClient) GetServerUpdateStatus(cacheHostName tc.CacheName, reqHdr http.Header) (tc.ServerUpdateStatus, toclientlib.ReqInf, error) {
	if cl.C == nil {
		return cl.Old.GetServerUpdateStatus(cacheHostName)
	}

	status := tc.ServerUpdateStatus{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "server_update_status_"+string(cacheHostName), &status, func(obj interface{}) error {
		toStatus, toReqInf, err := cl.C.GetServerUpdateStatusWithHdr(string(cacheHostName), reqHdr)
		if err != nil {
			return errors.New("getting server update status from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		status := obj.(*tc.ServerUpdateStatus)
		*status = toStatus
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return tc.ServerUpdateStatus{}, reqInf, errors.New("getting server update status: " + err.Error())
	}
	return status, reqInf, nil
}
