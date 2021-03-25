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
	"net"
	"strings"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

const TrafficOpsProxyParameterName = `tm.rev_proxy.url`

// GetTOData gets all the data from Traffic Ops needed to generate config.
// Returns the data, the addresses of all Traffic Ops' requested, and any error.
func GetTOData(cfg config.TCCfg) (*config.TOData, []net.Addr, error) {
	start := time.Now()
	defer func() { log.Infof("GetTOData took %v\n", time.Since(start)) }()

	toIPs := &sync.Map{} // each Traffic Ops request could get a different IP, so track them all
	toData := &config.TOData{}

	globalParams, toAddr, err := cfg.TOClient.GetGlobalParameters()
	if err != nil {
		return nil, nil, errors.New("getting global parameters: " + err.Error())
	}
	toIPs.Store(toAddr, nil)
	toData.GlobalParams = globalParams

	if !cfg.DisableProxy {
		toProxyURLStr := ""
		for _, param := range globalParams {
			if param.Name == TrafficOpsProxyParameterName {
				toProxyURLStr = param.Value
				break
			}
		}
		if toProxyURLStr != "" {
			realTOURL := cfg.TOClient.C.URL
			cfg.TOClient.C.URL = toProxyURLStr
			log.Infoln("using Traffic Ops proxy '" + toProxyURLStr + "'")
			if _, _, err := cfg.TOClient.C.GetCDNs(); err != nil {
				log.Warnln("Traffic Ops proxy '" + toProxyURLStr + "' failed to get CDNs, falling back to real Traffic Ops")
				cfg.TOClient.C.URL = realTOURL
			}
		} else {
			log.Infoln("Traffic Ops proxy enabled, but GLOBAL Parameter '" + TrafficOpsProxyParameterName + "' missing or empty, not using proxy")
		}
	} else {
		log.Infoln("Traffic Ops proxy is disabled, not checking or using GLOBAL Parameter '" + TrafficOpsProxyParameterName)
	}

	serversF := func() error {
		defer func(start time.Time) { log.Infof("serversF took %v\n", time.Since(start)) }(time.Now())
		// TODO TOAPI add /servers?cdn=1 query param
		servers, toAddr, err := cfg.TOClient.GetServers()
		if err != nil {
			return errors.New("getting servers: " + err.Error())
		}
		toData.Servers = servers
		toIPs.Store(toAddr, nil)

		server := &atscfg.Server{}
		for _, toServer := range servers {
			if toServer.HostName != nil && *toServer.HostName == cfg.CacheHostName {
				server = &toServer
				break
			}
		}
		if server.ID == nil {
			return errors.New("server '" + cfg.CacheHostName + " not found in servers")
		} else if server.CDNName == nil {
			return errors.New("server '" + cfg.CacheHostName + " missing CDNName")
		} else if server.CDNID == nil {
			return errors.New("server '" + cfg.CacheHostName + " missing CDNID")
		} else if server.Profile == nil {
			return errors.New("server '" + cfg.CacheHostName + " missing Profile")
		}

		toData.Server = server

		sslF := func() error {
			defer func(start time.Time) { log.Infof("sslF took %v\n", time.Since(start)) }(time.Now())
			keys, toAddr, err := cfg.TOClient.GetCDNSSLKeys(tc.CDNName(*server.CDNName))
			if err != nil {
				return errors.New("getting cdn '" + *server.CDNName + "': " + err.Error())
			}
			toData.SSLKeys = keys
			toIPs.Store(toAddr, nil)
			return nil
		}
		dsF := func() error {
			defer func(start time.Time) { log.Infof("dsF took %v\n", time.Since(start)) }(time.Now())

			dses, toAddr, err := cfg.TOClient.GetCDNDeliveryServices(*server.CDNID)
			if err != nil {
				return errors.New("getting delivery services: " + err.Error())
			}
			toData.DeliveryServices = dses
			toIPs.Store(toAddr, nil)

			// TODO uncomment when MSO Origins are changed to not use DSS, to avoid the DSS call if it isn't necessary
			// allDSesHaveTopologies := true
			// for _, ds := range toData.DeliveryServices {
			// 	if ds.CDNID == nil || *ds.CDNID != *server.CDNID {
			// 		continue
			// 	}
			// 	if ds.Topology == nil {
			// 		allDSesHaveTopologies = false
			// 		break
			// 	}
			// }

			dssF := func() error {
				defer func(start time.Time) { log.Infof("dssF took %v\n", time.Since(start)) }(time.Now())
				dss, toAddr, err := cfg.TOClient.GetDeliveryServiceServers(nil, nil)
				if err != nil {
					return errors.New("getting delivery service servers: " + err.Error())
				}
				toData.DeliveryServiceServers = dss
				toIPs.Store(toAddr, nil)
				return nil
			}

			uriSignKeysF := func() error {
				defer func(start time.Time) { log.Infof("uriF took %v\n", time.Since(start)) }(time.Now())
				uriSigningKeys := map[tc.DeliveryServiceName][]byte{}
				for _, ds := range dses {
					if ds.XMLID == nil {
						continue // TODO warn?
					}
					// TODO read meta config gen, and only include servers which are included in the meta (assigned to edge or all for mids? read the meta gen to find out)
					if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURISigning {
						continue
					}
					keys, toAddr, err := cfg.TOClient.GetURISigningKeys(*ds.XMLID)
					if err != nil {
						if strings.Contains(strings.ToLower(err.Error()), "not found") {
							log.Errorln("Delivery service '" + *ds.XMLID + "' is uri_signing, but keys were not found! Skipping!")
							continue
						} else {
							return errors.New("getting uri signing keys for ds '" + *ds.XMLID + "': " + err.Error())
						}
					}
					toIPs.Store(toAddr, nil)
					uriSigningKeys[tc.DeliveryServiceName(*ds.XMLID)] = keys
				}
				toData.URISigningKeys = uriSigningKeys
				return nil
			}

			urlSigKeysF := func() error {
				defer func(start time.Time) { log.Infof("urlF took %v\n", time.Since(start)) }(time.Now())
				urlSigKeys := map[tc.DeliveryServiceName]tc.URLSigKeys{}
				for _, ds := range dses {
					if ds.XMLID == nil {
						continue // TODO warn?
					}
					// TODO read meta config gen, and only include servers which are included in the meta (assigned to edge or all for mids? read the meta gen to find out)
					if ds.SigningAlgorithm == nil || *ds.SigningAlgorithm != tc.SigningAlgorithmURLSig {
						continue
					}
					keys, toAddr, err := cfg.TOClient.GetURLSigKeys(*ds.XMLID)
					if err != nil {
						if strings.Contains(strings.ToLower(err.Error()), "not found") {
							log.Errorln("Delivery service '" + *ds.XMLID + "' is url_sig, but keys were not found! Skipping!: " + err.Error())
							continue
						} else {
							return errors.New("getting url sig keys for ds '" + *ds.XMLID + "': " + err.Error())
						}
					}
					toIPs.Store(toAddr, nil)
					urlSigKeys[tc.DeliveryServiceName(*ds.XMLID)] = keys
				}
				toData.URLSigKeys = urlSigKeys
				return nil
			}

			fs := []func() error{}
			if !cfg.RevalOnly {
				fs = append([]func() error{uriSignKeysF, urlSigKeysF}, fs...) // skip keys for reval-only, which doesn't need them
			}
			if !cfg.RevalOnly { // TODO when MSO Origins are changed to not use DSS, we can add `&& !allDSesHaveTopologies` here, to avoid the DSS call if it isn't necessary
				// skip DSS if reval-only (which doesn't need DSS)
				fs = append([]func() error{dssF}, fs...)
			}

			return util.JoinErrs(runParallel(fs))
		}
		serverParamsF := func() error {
			defer func(start time.Time) { log.Infof("serverParamsF took %v\n", time.Since(start)) }(time.Now())
			params, toAddr, err := cfg.TOClient.GetServerProfileParameters(*server.Profile)
			if err != nil {
				return errors.New("getting server profile '" + *server.Profile + "' parameters: " + err.Error())
			} else if len(params) == 0 {
				return errors.New("getting server profile '" + *server.Profile + "' parameters: no parameters (profile not found?)")
			}
			toData.ServerParams = params
			toIPs.Store(toAddr, nil)
			return nil
		}
		cdnF := func() error {
			defer func(start time.Time) { log.Infof("cdnF took %v\n", time.Since(start)) }(time.Now())
			cdn, toAddr, err := cfg.TOClient.GetCDN(tc.CDNName(*server.CDNName))
			if err != nil {
				return errors.New("getting cdn '" + *server.CDNName + "': " + err.Error())
			}
			toData.CDN = &cdn
			toIPs.Store(toAddr, nil)
			return nil
		}
		profileF := func() error {
			defer func(start time.Time) { log.Infof("profileF took %v\n", time.Since(start)) }(time.Now())
			profile, toAddr, err := cfg.TOClient.GetProfileByName(*server.Profile)
			if err != nil {
				return errors.New("getting profile '" + *server.Profile + "': " + err.Error())
			}
			toData.Profile = profile
			toIPs.Store(toAddr, nil)
			return nil
		}
		fs := []func() error{dsF, serverParamsF, cdnF, profileF}
		if !cfg.RevalOnly {
			fs = append([]func() error{sslF}, fs...) // skip ssl keys for reval only, which doesn't need them
		}
		return util.JoinErrs(runParallel(fs))
	}

	cgF := func() error {
		defer func(start time.Time) { log.Infof("cfF took %v\n", time.Since(start)) }(time.Now())
		cacheGroups, toAddr, err := cfg.TOClient.GetCacheGroups()
		if err != nil {
			return errors.New("getting cachegroups: " + err.Error())
		}
		toData.CacheGroups = cacheGroups
		toIPs.Store(toAddr, nil)
		return nil
	}
	jobsF := func() error {
		defer func(start time.Time) { log.Infof("jobsF took %v\n", time.Since(start)) }(time.Now())
		jobs, toAddr, err := cfg.TOClient.GetJobs() // TODO add cdn query param to jobs endpoint
		if err != nil {
			return errors.New("getting jobs: " + err.Error())
		}
		toData.Jobs = jobs
		toIPs.Store(toAddr, nil)
		return nil
	}
	capsF := func() error {
		defer func(start time.Time) { log.Infof("capsF took %v\n", time.Since(start)) }(time.Now())
		caps, toAddr, err := cfg.TOClient.GetServerCapabilitiesByID(nil) // TODO change to not take a param; it doesn't use it to request TO anyway.
		if err != nil {
			log.Errorln("Server Capabilities error, skipping!")
			// return errors.New("getting server caps from Traffic Ops: " + err.Error())
		} else {
			toData.ServerCapabilities = caps
			toIPs.Store(toAddr, nil)
		}
		return nil
	}
	dsCapsF := func() error {
		defer func(start time.Time) { log.Infof("dscapsF took %v\n", time.Since(start)) }(time.Now())
		caps, toAddr, err := cfg.TOClient.GetDeliveryServiceRequiredCapabilitiesByID(nil)
		if err != nil {
			log.Errorln("DS Required Capabilities error, skipping!")
			// return errors.New("getting DS required capabilities: " + err.Error())
		} else {
			toData.DSRequiredCapabilities = caps
			toIPs.Store(toAddr, nil)
		}
		return nil
	}
	dsrF := func() error {
		defer func(start time.Time) { log.Infof("dsrF took %v\n", time.Since(start)) }(time.Now())
		dsr, toAddr, err := cfg.TOClient.GetDeliveryServiceRegexes()
		if err != nil {
			return errors.New("getting delivery service regexes: " + err.Error())
		}
		toIPs.Store(toAddr, nil)
		toData.DeliveryServiceRegexes = dsr
		return nil
	}
	cacheKeyParamsF := func() error {
		defer func(start time.Time) { log.Infof("cacheKeyParamsF took %v\n", time.Since(start)) }(time.Now())
		params, toAddr, err := cfg.TOClient.GetConfigFileParameters(atscfg.CacheKeyParameterConfigFile)
		if err != nil {
			return errors.New("getting cache key parameters: " + err.Error())
		}
		toIPs.Store(toAddr, nil)
		toData.CacheKeyParams = params
		return nil
	}
	parentConfigParamsF := func() error {
		defer func(start time.Time) { log.Infof("parentConfigParamsF took %v\n", time.Since(start)) }(time.Now())
		parentConfigParams, toAddr, err := cfg.TOClient.GetConfigFileParameters("parent.config") // TODO make const in lib/go-atscfg
		if err != nil {
			return errors.New("getting parent.config parameters: " + err.Error())
		}
		toIPs.Store(toAddr, nil)
		toData.ParentConfigParams = parentConfigParams
		return nil
	}

	topologiesF := func() error {
		defer func(start time.Time) { log.Infof("topologiesF took %v\n", time.Since(start)) }(time.Now())
		topologies, toAddr, err := cfg.TOClient.GetTopologies()
		if err != nil {
			return errors.New("getting topologies: " + err.Error())
		}
		toIPs.Store(toAddr, nil)
		toData.Topologies = topologies
		return nil
	}

	fs := []func() error{serversF, cgF, jobsF}
	if !cfg.RevalOnly {
		// skip data not needed for reval, if we're reval-only
		fs = append([]func() error{dsrF, cacheKeyParamsF, parentConfigParamsF, capsF, dsCapsF, topologiesF}, fs...)
	}
	errs := runParallel(fs)

	toIPArr := []net.Addr{}
	toIPs.Range(func(key, val interface{}) bool {
		toIPArr = append(toIPArr, key.(net.Addr))
		return true
	})

	return toData, toIPArr, util.JoinErrs(errs)
}

// runParallel runs all funcs in fs in parallel goroutines, and returns after all funcs have returned.
// Returns a slice of the errors returned by each func. The order of the errors will not be the same as the order of fs.
// All funcs in fs must be safe to run in parallel.
func runParallel(fs []func() error) []error {
	errs := []error{}
	doneChan := make(chan error, len(fs))
	for _, fPtr := range fs {
		f := fPtr // because functions are pointers, f will change in the loop. Need create a new variable here to close around.
		go func() { doneChan <- f() }()
	}
	for i := 0; i < len(fs); i++ {
		errs = append(errs, <-doneChan)
	}
	return errs
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

// ParamsToMap converts a []tc.Parameter to a map[paramName]paramValue.
// If multiple params have the same value, the first one in params will be used an an error will be logged.
// See ParamArrToMultiMap.
func ParamsToMap(params []tc.Parameter) map[string]string {
	mp := map[string]string{}
	for _, param := range params {
		if val, ok := mp[param.Name]; ok {
			if val < param.Value {
				log.Errorln("config generation got multiple parameters for name '" + param.Name + "' - ignoring '" + param.Value + "'")
				continue
			} else {
				log.Errorln("config generation got multiple parameters for name '" + param.Name + "' - ignoring '" + val + "'")
			}
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

// logWarnings writes all strings in warnings to the warning log, with the context prefix.
// If warnings is empty, no log is written.
func logWarnings(context string, warnings []string) {
	for _, warn := range warnings {
		log.Warnln(context + warn)
	}
}
