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
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/toreq"
)

// TOData is the Traffic Ops data needed to generate configs.
// See each field for details on the data required.
// - If a field says 'must', the creation of TOData is guaranteed to do so, and users of the struct may rely on that.
// - If it says 'may', the creation may or may not do so, and therefore users of the struct must filter if they
//   require the potential fields to be omitted to generate correctly.
type TOData struct {
	// Servers must be all the servers from Traffic Ops. May include servers not on the current cdn.
	Servers []tc.Server

	// CacheGroups must be all cachegroups in Traffic Ops with Servers on the current server's cdn. May also include CacheGroups without servers on the current cdn.
	CacheGroups []tc.CacheGroupNullable

	// GlobalParams must be all Parameters in Traffic Ops on the tc.GlobalProfileName Profile. Must not include other parameters.
	GlobalParams []tc.Parameter

	// ScopeParams must be all Parameters in Traffic Ops with the name "scope". Must not include other Parameters.
	ScopeParams []tc.Parameter

	// ServerParams must be all Parameters on the Profile of the current server. Must not include other Parameters.
	ServerParams []tc.Parameter

	// CacheKeyParams must be all Parameters with the ConfigFile atscfg.CacheKeyParameterConfigFile.
	CacheKeyParams []tc.Parameter

	// ParentConfigParams must be all Parameters with the ConfigFile "parent.config.
	ParentConfigParams []tc.Parameter

	// DeliveryServices must include all Delivery Services on the current server's cdn, including those not assigned to the server. Must not include delivery services on other cdns.
	DeliveryServices []tc.DeliveryServiceNullable

	// DeliveryServiceServers must include all delivery service servers in Traffic Ops for all delivery services on the current cdn, including those not assigned to the current server.
	DeliveryServiceServers []tc.DeliveryServiceServer

	// Server must be the server we're fetching configs from
	Server tc.Server

	// TOToolName must be the Parameter named 'tm.toolname' on the tc.GlobalConfigFileName Profile.
	TOToolName string

	// TOToolName must be the Parameter named 'tm.url' on the tc.GlobalConfigFileName Profile.
	TOURL string

	// Jobs must be all Jobs on the server's CDN. May include jobs on other CDNs.
	Jobs []tc.Job

	// CDN must be the CDN of the server.
	CDN tc.CDN

	// DeliveryServiceRegexes must be all regexes on all delivery services on this server's cdn.
	DeliveryServiceRegexes []tc.DeliveryServiceRegexes

	// Profile must be the Profile of the server being requested.
	Profile tc.Profile

	// URISigningKeys must be a map of every delivery service which is URI Signed, to its keys.
	URISigningKeys map[tc.DeliveryServiceName][]byte

	// URLSigKeys must be a map of every delivery service which uses URL Sig, to its keys.
	URLSigKeys map[tc.DeliveryServiceName]tc.URLSigKeys

	// ServerCapabilities must be a map of all server IDs on this server's CDN, to a set of their capabilities. May also include servers from other cdns.
	ServerCapabilities map[int]map[atscfg.ServerCapability]struct{}

	// DSRequiredCapabilities must be a map of all delivery service IDs on this server's CDN, to a set of their required capabilities. Delivery Services with no required capabilities may not have an entry in the map.
	DSRequiredCapabilities map[int]map[atscfg.ServerCapability]struct{}
}

// TODO: validate all "profile scope" files are the server's profile.
//       If they ever weren't, we'll send bad data, because we're only getting the server's profile data.
//       Getting all data for all profiles in TOData isn't reasonable.

func GetTOData(cfg config.TCCfg) (*TOData, error) {
	// TODO TOAPI add /servers?cdn=1 query param
	servers, err := toreq.GetServers(cfg)
	if err != nil {
		return nil, errors.New("getting servers: " + err.Error())
	}

	server := tc.Server{ID: atscfg.InvalidID}
	for _, toServer := range servers {
		if toServer.HostName == cfg.CacheHostName {
			server = toServer
			break
		}
	}
	if server.ID == atscfg.InvalidID {
		return nil, errors.New("server '" + cfg.CacheHostName + " not found in servers")
	}

	cacheGroups, err := toreq.GetCacheGroups(cfg)
	if err != nil {
		return nil, errors.New("getting cachegroups: " + err.Error())
	}

	// TODO test whether it's faster to get all params and then filter
	globalParams, err := toreq.GetGlobalParameters(cfg)
	if err != nil {
		return nil, errors.New("getting global parameters: " + err.Error())
	}

	scopeParams, err := toreq.GetParametersByName(cfg, "scope")
	if err != nil {
		return nil, errors.New("getting scope parameters: " + err.Error())
	}

	// TODO check if len == 0 and return an error here?
	//      Basically ever config file does, so we might as well do it here and reduce duplicate code.
	//      It shouldn't ever be 0 anyway, because that should only happen if the profile doesn't exist, and a server's profile should always exist.
	serverParams, err := toreq.GetServerProfileParameters(cfg, server.Profile)
	if err != nil {
		return nil, errors.New("getting server profile '" + server.Profile + "' parameters: " + err.Error())
	}
	if len(serverParams) == 0 {
		return nil, errors.New("getting server profile '" + server.Profile + "' parameters: no parameters (profile not found?)")
	}

	deliveryServices, err := toreq.GetCDNDeliveryServices(cfg, server.CDNID)
	if err != nil {
		return nil, errors.New("getting delivery services: " + err.Error())
	}

	dsServers, err := toreq.GetDeliveryServiceServers(cfg, nil, nil)
	if err != nil {
		return nil, errors.New("getting delivery service servers: " + err.Error())
	}

	toToolName, toURL := toreq.GetTOToolNameAndURL(globalParams)

	jobs, err := toreq.GetJobs(cfg) // TODO add cdn query param to jobs endpoint
	if err != nil {
		return nil, errors.New("getting jobs: " + err.Error())
	}

	cdn, err := toreq.GetCDN(cfg, tc.CDNName(server.CDNName))
	if err != nil {
		return nil, errors.New("getting cdn '" + server.CDNName + "': " + err.Error())
	}

	dsRegexes, err := toreq.GetDeliveryServiceRegexes(cfg)
	if err != nil {
		return nil, errors.New("getting delivery service regexes: " + err.Error())
	}

	cacheKeyParams, err := toreq.GetConfigFileParameters(cfg, atscfg.CacheKeyParameterConfigFile)
	if err != nil {
		return nil, errors.New("getting cache key parameters: " + err.Error())
	}

	parentConfigParams, err := toreq.GetConfigFileParameters(cfg, "parent.config") // TODO make const in lib/go-atscfg
	if err != nil {
		return nil, errors.New("getting parent.config parameters: " + err.Error())
	}

	// TODO verify used
	profile, err := toreq.GetProfileByName(cfg, server.Profile)
	if err != nil {
		return nil, errors.New("getting profile '" + server.Profile + "': " + err.Error())
	}

	uriSigningKeys := map[tc.DeliveryServiceName][]byte{}
	for _, ds := range deliveryServices {
		if ds.XMLID == nil {
			continue // TODO warn?
		}
		// TODO read meta config gen, and only include servers which are included in the meta (assigned to edge or all for mids? read the meta gen to find out)
		if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURISigning {
			continue
		}
		keys, err := toreq.GetURISigningKeys(cfg, *ds.XMLID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				log.Errorln("Delivery service '" + *ds.XMLID + "' is uri_signing, but keys were not found! Skipping!")
				continue
			} else {
				return nil, errors.New("getting uri signing keys for ds '" + *ds.XMLID + "': " + err.Error())
			}
		}
		uriSigningKeys[tc.DeliveryServiceName(*ds.XMLID)] = keys
	}

	urlSigKeys := map[tc.DeliveryServiceName]tc.URLSigKeys{}
	for _, ds := range deliveryServices {
		if ds.XMLID == nil {
			continue // TODO warn?
		}
		// TODO read meta config gen, and only include servers which are included in the meta (assigned to edge or all for mids? read the meta gen to find out)
		if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURLSig {
			continue
		}
		keys, err := toreq.GetURLSigKeys(cfg, *ds.XMLID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				log.Errorln("Delivery service '" + *ds.XMLID + "' is url_sig, but keys were not found! Skipping!: " + err.Error())
				continue
			} else {
				return nil, errors.New("getting url sig keys for ds '" + *ds.XMLID + "': " + err.Error())
			}
		}
		urlSigKeys[tc.DeliveryServiceName(*ds.XMLID)] = keys
	}

	serverCapabilities, err := toreq.GetServerCapabilitiesByID(cfg, nil) // TODO change to not take a param; it doesn't use it to request TO anyway.
	if err != nil {
		log.Errorln("Server Capabilities error, skipping!")
		// return nil, errors.New("getting server caps from Traffic Ops: " + err.Error())
	}

	dsRequiredCapabilities, err := toreq.GetDeliveryServiceRequiredCapabilitiesByID(cfg, nil)
	if err != nil {
		log.Errorln("DS Required Capabilities error, skipping!")
		// return nil, errors.New("getting DS required capabilities: " + err.Error())
	}

	return &TOData{
		Servers:                servers,
		CacheGroups:            cacheGroups,
		GlobalParams:           globalParams,
		ScopeParams:            scopeParams,
		ServerParams:           serverParams,
		CacheKeyParams:         cacheKeyParams,
		ParentConfigParams:     parentConfigParams,
		DeliveryServices:       deliveryServices,
		DeliveryServiceServers: dsServers,
		Server:                 server,
		TOToolName:             toToolName,
		TOURL:                  toURL,
		Jobs:                   jobs,
		CDN:                    cdn,
		DeliveryServiceRegexes: dsRegexes,
		Profile:                profile,
		URISigningKeys:         uriSigningKeys,
		URLSigKeys:             urlSigKeys,
		ServerCapabilities:     serverCapabilities,
		DSRequiredCapabilities: dsRequiredCapabilities,
	}, nil
}

func FilterDSS(dsses []tc.DeliveryServiceServer, dsIDs map[int]struct{}, serverIDs map[int]struct{}) []tc.DeliveryServiceServer {
	// TODO filter only DSes on this server's CDN? Does anything ever needs DSS cross-CDN? Surely not.
	//      Then, we can remove a bunch of config files that filter only DSes on the current cdn.
	filtered := []tc.DeliveryServiceServer{}
	for _, dss := range dsses {
		if dss.Server == nil || dss.DeliveryService == nil {
			continue // TODO warn?
		}
		if len(dsIDs) > 0 {
			if _, ok := dsIDs[*dss.DeliveryService]; !ok {
				continue
			}
		}
		if len(serverIDs) > 0 {
			if _, ok := serverIDs[*dss.Server]; !ok {
				continue
			}
		}
		filtered = append(filtered, dss)
	}
	return filtered
}

// TCParamsToParamsWithProfiles unmarshals the Profiles that the tc struct doesn't.
func TCParamsToParamsWithProfiles(tcParams []tc.Parameter) ([]ParameterWithProfiles, error) {
	params := make([]ParameterWithProfiles, 0, len(tcParams))
	for _, tcParam := range tcParams {
		param := ParameterWithProfiles{Parameter: tcParam}

		profiles := []string{}
		if err := json.Unmarshal(tcParam.Profiles, &profiles); err != nil {
			return nil, errors.New("unmarshalling JSON from parameter '" + strconv.Itoa(param.ID) + "': " + err.Error())
		}
		param.ProfileNames = profiles
		param.Profiles = nil
		params = append(params, param)
	}
	return params, nil
}

type ParameterWithProfiles struct {
	tc.Parameter
	ProfileNames []string
}

type ParameterWithProfilesMap struct {
	tc.Parameter
	ProfileNames map[string]struct{}
}

func ParameterWithProfilesToMap(tcParams []ParameterWithProfiles) []ParameterWithProfilesMap {
	params := []ParameterWithProfilesMap{}
	for _, tcParam := range tcParams {
		param := ParameterWithProfilesMap{Parameter: tcParam.Parameter, ProfileNames: map[string]struct{}{}}
		for _, profile := range tcParam.ProfileNames {
			param.ProfileNames[profile] = struct{}{}
		}
		params = append(params, param)
	}
	return params
}

// FilterParams filters params and returns only the parameters which match configFile, name, and value.
// If configFile, name, or value is the empty string, it is not filtered.
// Returns a slice of parameters.
func FilterParams(params []tc.Parameter, configFile string, name string, value string, omitName string) []tc.Parameter {
	filtered := []tc.Parameter{}
	for _, param := range params {
		if configFile != "" && param.ConfigFile != configFile {
			continue
		}
		if name != "" && param.Name != name {
			continue
		}
		if value != "" && param.Value != value {
			continue
		}
		if omitName != "" && param.Name == omitName {
			continue
		}
		filtered = append(filtered, param)
	}
	return filtered
}

// ParamArrToMap converts a []tc.Parameter to a map[paramName]paramValue.
// If multiple params have the same value, the first one in params will be used an an error will be logged.
// See ParamArrToMultiMap.
func ParamsToMap(params []tc.Parameter) map[string]string {
	mp := map[string]string{}
	for _, param := range params {
		if val, ok := mp[param.Name]; ok {
			log.Errorln("config generation got multiple parameters for name '" + param.Name + "' - using '" + val + "'")
			continue
		}
		mp[param.Name] = param.Value
	}
	return mp
}

// ParamArrToMultiMap converts a []tc.Parameter to a map[paramName][]paramValue.
func ParamsToMultiMap(params []tc.Parameter) map[string][]string {
	mp := map[string][]string{}
	for _, param := range params {
		mp[param.Name] = append(mp[param.Name], param.Value)
	}
	return mp
}

// type TOData struct {
// 	Servers                []tc.Server
// 	CacheGroups            []tc.CacheGroupNullable
// 	GlobalParams           []tc.Parameter
// 	ScopeParams            []tc.Parameter
// 	ServerParams           []tc.Parameter
// 	CacheKeyParams         []tc.Parameter
// 	ParentConfigParams     []tc.Parameter
// 	DeliveryServices       []tc.DeliveryServiceNullable
// 	DeliveryServiceServers []tc.DeliveryServiceServer
// 	Server                 tc.Server
// 	TOToolName             string
// 	TOURL                  string
// 	Jobs                   []tc.Job
// 	CDN                    tc.CDN
// 	DeliveryServiceRegexes []tc.DeliveryServiceRegexes
// 	Profile                tc.Profile
// 	URISigningKeys         map[tc.DeliveryServiceName][]byte
// 	URLSigKeys             map[tc.DeliveryServiceName]tc.URLSigKeys
// 	ServerCapabilities     map[int]map[atscfg.ServerCapability]struct{}
// 	DSRequiredCapabilities map[int]map[atscfg.ServerCapability]struct{}
// }
