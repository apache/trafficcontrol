package toreqold

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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil/toreq/torequtil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

/*func (cl *TOClient) GetProfileByName(profileName string) (tc.Profile, toclientlib.ReqInf, error) {
	profile := tc.Profile{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "profile_"+profileName, &profile, func(obj interface{}) error {
		toProfiles, toReqInf, err := cl.c.GetProfileByName(profileName)
		if err != nil {
			return errors.New("getting profile '" + profileName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		if len(toProfiles) != 1 {
			return errors.New("getting profile '" + profileName + "'from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': expected 1 Profile, got " + strconv.Itoa(len(toProfiles)))
		}

		profile := obj.(*tc.Profile)
		*profile = toProfiles[0]
		reqInf = toReqInf
		return nil
	})

	if err != nil {
		return tc.Profile{}, reqInf, errors.New("getting profile '" + profileName + "': " + err.Error())
	}
	return profile, reqInf, nil
}*/

func (cl *TOClient) GetGlobalParameters(reqHdr http.Header) ([]tc.ParameterV5, toclientlib.ReqInf, error) {
	globalParams := []tc.ParameterV5{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "profile_global_parameters", &globalParams, func(obj interface{}) error {
		toParams, toReqInf, err := cl.c.GetParametersByProfileName(tc.GlobalProfileName, *ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting global profile '" + tc.GlobalProfileName + "' parameters from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.ParameterV5)
		*params = parametersToLatest(toParams.Response)
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting global profile '" + tc.GlobalProfileName + "' parameters: " + err.Error())
	}
	return globalParams, reqInf, nil
}

func (cl *TOClient) GetServers(reqHdr http.Header) ([]atscfg.Server, toclientlib.ReqInf, error) {
	servers := []atscfg.Server{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "servers", &servers, func(obj interface{}) error {
		toServers, toReqInf, err := cl.c.GetServers(*ReqOpts(reqHdr))
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
	server := atscfg.Server{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "server-name-"+serverHostName, &server, func(obj interface{}) error {

		vals := &url.Values{}
		vals.Add("hostName", serverHostName)
		params := *ReqOpts(reqHdr)
		params.QueryParameters = *vals
		toServers, toReqInf, err := cl.c.GetServers(params)
		if err != nil {
			return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		} else if len(toServers.Response) < 1 {
			return errors.New("getting server name '" + serverHostName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': no servers returned")
		}
		asv, err := serverToLatest(toServers.Response[0])
		if err != nil {
			return errors.New("converting server to latest version: " + err.Error())
		}
		server := obj.(*atscfg.Server)
		*server = *asv
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting server name '" + serverHostName + "': " + err.Error())
	}
	return &server, reqInf, nil
}

func (cl *TOClient) GetCacheGroups(reqHdr http.Header) ([]tc.CacheGroupNullable, toclientlib.ReqInf, error) {
	cacheGroups := []tc.CacheGroupNullable{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cachegroups", &cacheGroups, func(obj interface{}) error {
		toCacheGroups, toReqInf, err := cl.c.GetCacheGroups(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting cachegroups from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		cacheGroups := obj.(*[]tc.CacheGroupNullable)
		*cacheGroups = toCacheGroups.Response
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

func (cl *TOClient) GetDeliveryServiceServers(dsIDs []int, serverIDs []int, cdnName string, reqHdr http.Header) ([]tc.DeliveryServiceServerV5, toclientlib.ReqInf, error) {
	const sortIDsInHash = true
	reqInf := toclientlib.ReqInf{}
	serverIDsStr := ""
	dsIDsStr := ""
	dsIDsToFetch := ([]int)(nil)
	sIDsToFetch := ([]int)(nil)
	if !DeliveryServiceServersAlwaysGetAll {
		if len(dsIDs) > 0 {
			dsIDsStr = base64.RawURLEncoding.EncodeToString(util.HashInts(dsIDs, sortIDsInHash))
		}
		if len(serverIDs) > 0 {
			serverIDsStr = base64.RawURLEncoding.EncodeToString(util.HashInts(serverIDs, sortIDsInHash))
		}
		dsIDsToFetch = dsIDs
		sIDsToFetch = serverIDs
	}

	dsServers := []tc.DeliveryServiceServerV5{}
	err := torequtil.GetRetry(cl.NumRetries, "deliveryservice_servers_s"+serverIDsStr+"_d_"+dsIDsStr, &dsServers, func(obj interface{}) error {
		dsIDStrs := []string{}
		for _, dsID := range dsIDsToFetch {
			dsIDStrs = append(dsIDStrs, strconv.Itoa(dsID))
		}

		serverIDStrs := []string{}
		for _, serverID := range sIDsToFetch {
			serverIDStrs = append(serverIDStrs, strconv.Itoa(serverID))
		}

		queryParams := url.Values{}
		queryParams.Set("limit", "999999") // TODO add "no limit" param to DSS endpoint
		queryParams.Set("cdn", cdnName)
		queryParams.Set("orderby", "") // prevent unnecessary sorting of the response
		if len(dsIDsToFetch) > 0 {
			queryParams.Set("deliveryserviceids", strings.Join(dsIDStrs, ","))
		}
		if len(sIDsToFetch) > 0 {
			queryParams.Set("serverids", strings.Join(serverIDStrs, ","))
		}
		opts := *ReqOpts(reqHdr)
		opts.QueryParameters = queryParams
		toDSS, toReqInf, err := cl.c.GetDeliveryServiceServers(opts)
		if err != nil {
			return errors.New("getting delivery service servers from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dss := obj.(*[]tc.DeliveryServiceServerV5)
		*dss = deliveryServiceServersToLatest(toDSS.Response)
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
	filteredDSServers := []tc.DeliveryServiceServerV5{}
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

func (cl *TOClient) GetServerProfileParameters(profileName string, reqHdr http.Header) ([]tc.ParameterV5, toclientlib.ReqInf, error) {
	serverProfileParameters := []tc.ParameterV5{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "profile_"+profileName+"_parameters", &serverProfileParameters, func(obj interface{}) error {
		toParams, toReqInf, err := cl.c.GetParametersByProfileName(profileName, *ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting server profile '" + profileName + "' parameters from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.ParameterV5)
		*params = parametersToLatest(toParams.Response)
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
	deliveryServices := []atscfg.DeliveryService{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cdn_"+strconv.Itoa(cdnID)+"_deliveryservices", &deliveryServices, func(obj interface{}) error {
		params := url.Values{}
		params.Set("cdn", strconv.Itoa(cdnID))
		opts := *ReqOpts(reqHdr)
		opts.QueryParameters = params
		toDSes, toReqInf, err := cl.c.GetDeliveryServices(opts)
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dses := obj.(*[]atscfg.DeliveryService)
		*dses = dsesToLatest(toDSes.Response)
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
	topologies := []tc.Topology{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "topologies", &topologies, func(obj interface{}) error {
		toTopologies, toReqInf, err := cl.c.GetTopologies(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting topologies from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		topologies := obj.(*[]tc.Topology)
		*topologies = toTopologies.Response
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting topologies: " + err.Error())
	}
	return topologies, reqInf, nil
}

func (cl *TOClient) GetConfigFileParameters(configFile string, reqHdr http.Header) ([]tc.ParameterV5, toclientlib.ReqInf, error) {
	params := []tc.ParameterV5{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "config_file_"+configFile+"_parameters", &params, func(obj interface{}) error {
		toParams, toReqInf, err := GetParametersByConfigFile(cl.c, configFile, ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting delivery services from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.ParameterV5)
		*params = parametersToLatest(toParams)
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting parent.config parameters: " + err.Error())
	}
	return params, reqInf, nil
}

func (cl *TOClient) GetCDN(cdnName tc.CDNName, reqHdr http.Header) (tc.CDN, toclientlib.ReqInf, error) {
	cdn := tc.CDN{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cdn_"+string(cdnName), &cdn, func(obj interface{}) error {
		toCDNs, toReqInf, err := GetCDNByName(cl.c, cdnName, ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting cdn from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		cdn := obj.(*tc.CDN)
		*cdn = toCDNs
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return tc.CDN{}, reqInf, errors.New("getting cdn: " + err.Error())
	}
	return cdn, reqInf, nil
}

func (cl *TOClient) GetCDNs(reqHdr http.Header) ([]tc.CDN, toclientlib.ReqInf, error) {
	cdns := []tc.CDN{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cdns", &cdns, func(obj interface{}) error {
		toCDNs, toReqInf, err := cl.c.GetCDNs(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting cdn from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		cdns := obj.(*[]tc.CDN)
		*cdns = toCDNs.Response
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return []tc.CDN{}, reqInf, errors.New("getting cdn: " + err.Error())
	}
	return cdns, reqInf, nil
}

func (cl *TOClient) GetURLSigKeys(dsName string, reqHdr http.Header) (tc.URLSigKeys, toclientlib.ReqInf, error) {
	keys := tc.URLSigKeys{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "urlsigkeys_"+string(dsName), &keys, func(obj interface{}) error {
		toKeys, toReqInf, err := GetDeliveryServiceURLSigKeys(cl.c, dsName, ReqOpts(reqHdr))
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
	keys := []byte{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "urisigningkeys_"+string(dsName), &keys, func(obj interface{}) error {
		toKeys, toReqInf, err := cl.c.GetDeliveryServiceURISigningKeys(dsName, *ReqOpts(reqHdr))
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

func (cl *TOClient) GetParametersByName(paramName string, reqHdr http.Header) ([]tc.ParameterV5, toclientlib.ReqInf, error) {
	params := []tc.ParameterV5{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "parameters_name_"+paramName, &params, func(obj interface{}) error {
		toParams, toReqInf, err := GetParametersByName(cl.c, paramName, ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting parameters name '" + paramName + "' from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		params := obj.(*[]tc.ParameterV5)
		*params = parametersToLatest(toParams)
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting params name '" + paramName + "': " + err.Error())
	}
	return params, reqInf, nil
}

func (cl *TOClient) GetDeliveryServiceRegexes(reqHdr http.Header) ([]tc.DeliveryServiceRegexes, toclientlib.ReqInf, error) {
	regexes := []tc.DeliveryServiceRegexes{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "ds_regexes", &regexes, func(obj interface{}) error {
		toRegexes, toReqInf, err := cl.c.GetDeliveryServiceRegexes(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting ds regexes from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		regexes := obj.(*[]tc.DeliveryServiceRegexes)
		*regexes = toRegexes.Response
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting ds regexes: " + err.Error())
	}
	return regexes, reqInf, nil
}

func (cl *TOClient) GetJobs(reqHdr http.Header) ([]tc.InvalidationJobV4, toclientlib.ReqInf, error) {
	jobs := []tc.InvalidationJobV4{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "jobs", &jobs, func(obj interface{}) error {
		toJobs, toReqInf, err := cl.c.GetInvalidationJobs(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting jobs from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		jobs := obj.(*[]tc.InvalidationJobV4)
		*jobs = toJobs.Response
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting jobs: " + err.Error())
	}
	return jobs, reqInf, nil
}

func (cl *TOClient) GetServerCapabilitiesByID(serverIDs []int, reqHdr http.Header) (map[int]map[atscfg.ServerCapability]struct{}, toclientlib.ReqInf, error) {
	serverIDsStr := ""
	if len(serverIDs) > 0 {
		sortIDsInHash := true
		serverIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(serverIDs, sortIDsInHash)))
	}

	serverCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "server_capabilities_s_"+serverIDsStr, &serverCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toServerCaps, toReqInf, err := cl.c.GetServerServerCapabilities(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting server caps from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		serverCaps := obj.(*map[int]map[atscfg.ServerCapability]struct{})

		for _, sc := range toServerCaps.Response {
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
	dsIDsStr := ""
	if len(dsIDs) > 0 {
		sortIDsInHash := true
		dsIDsStr = base64.RawURLEncoding.EncodeToString((util.HashInts(dsIDs, sortIDsInHash)))
	}

	dsCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "ds_capabilities_d_"+dsIDsStr, &dsCaps, func(obj interface{}) error {
		// TODO add list of IDs to API+Client
		toDSCaps, toReqInf, err := cl.c.GetDeliveryServicesRequiredCapabilities(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting ds caps from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		dsCaps := obj.(*map[int]map[atscfg.ServerCapability]struct{})

		for _, sc := range toDSCaps.Response {
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
	keys := []tc.CDNSSLKeys{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "cdn_sslkeys_"+string(cdnName), &keys, func(obj interface{}) error {
		toKeys, toReqInf, err := cl.c.GetCDNSSLKeys(string(cdnName), *ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting cdn ssl keys from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		keys := obj.(*[]tc.CDNSSLKeys)
		*keys = toKeys.Response
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return []tc.CDNSSLKeys{}, reqInf, errors.New("getting cdn ssl keys: " + err.Error())
	}
	return keys, reqInf, nil
}

func (cl *TOClient) GetStatuses(reqHdr http.Header) ([]tc.Status, toclientlib.ReqInf, error) {
	statuses := []tc.Status{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "statuses", &statuses, func(obj interface{}) error {
		toStatus, toReqInf, err := cl.c.GetStatuses(*ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting old server update status from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		status := obj.(*[]tc.Status)
		*status = toStatus.Response
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return nil, reqInf, errors.New("getting old server update status: " + err.Error())
	}
	return statuses, reqInf, nil
}

// GetServerUpdateStatus returns the data, the Traffic Ops address, and any error.
func (cl *TOClient) GetServerUpdateStatus(cacheHostName tc.CacheName, reqHdr http.Header) (atscfg.ServerUpdateStatus, toclientlib.ReqInf, error) {
	status := atscfg.ServerUpdateStatus{}
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "server_update_status_"+string(cacheHostName), &status, func(obj interface{}) error {
		toStatus, toReqInf, err := cl.c.GetServerUpdateStatus(string(cacheHostName), *ReqOpts(reqHdr))
		if err != nil {
			return errors.New("getting server update status from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		status := obj.(*atscfg.ServerUpdateStatus)
		*status = serverUpdateStatusToLatest(toStatus.Response)[0]
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return atscfg.ServerUpdateStatus{}, reqInf, errors.New("getting server update status: " + err.Error())
	}
	return status, reqInf, nil
}

// SetServerUpdateStatus sets the server's update status in Traffic Ops.
func (cl *TOClient) SetServerUpdateStatus(cacheHostName tc.CacheName, configApply, revalApply *time.Time) (toclientlib.ReqInf, error) {
	reqInf := toclientlib.ReqInf{}
	err := torequtil.GetRetry(cl.NumRetries, "server_update_status_"+string(cacheHostName), nil, func(obj interface{}) error {
		_, toReqInf, err := cl.c.SetUpdateServerStatusTimes(string(cacheHostName), configApply, revalApply, *ReqOpts(nil))
		if err != nil {
			return errors.New("setting server update status from Traffic Ops '" + torequtil.MaybeIPStr(reqInf.RemoteAddr) + "': " + err.Error())
		}
		reqInf = toReqInf
		return nil
	})
	if err != nil {
		return reqInf, errors.New("setting server update status: " + err.Error())
	}
	return reqInf, nil
}
