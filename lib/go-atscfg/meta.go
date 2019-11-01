package atscfg

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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

type ConfigProfileParams struct {
	FileNameOnDisk string
	Location       string
	URL            string
	APIURI         string
}

// APIVersion is the Traffic Ops API version for config fiels.
// This is used to generate the meta config, which has API paths.
// Note the version in the meta config is not used by the atstccfg generator, which isn't actually an API.
// TODO change the config system to not use old API paths, and remove this.
const APIVersion = "1.4"

func MakeMetaConfig(
	serverHostName tc.CacheName,
	server *ServerInfo,
	tmURL string, // global tm.url Parameter
	tmReverseProxyURL string, // global tm.rev_proxy.url Parameter
	locationParams map[string]ConfigProfileParams, // map[configFile]params; 'location' and 'URL' Parameters on serverHostName's Profile
	uriSignedDSes []tc.DeliveryServiceName,
	scopeParams map[string]string, // map[configFileName]scopeParam
) string {
	if tmURL == "" {
		log.Errorln("ats.GetConfigMetadata: global tm.url parameter missing or empty! Setting empty in meta config!")
	}

	atsData := tc.ATSConfigMetaData{
		Info: tc.ATSConfigMetaDataInfo{
			ProfileID:         int(server.ProfileID),
			TOReverseProxyURL: tmReverseProxyURL,
			TOURL:             tmURL,
			ServerIPv4:        server.IP,
			ServerPort:        server.Port,
			ServerName:        server.HostName,
			CDNID:             server.CDNID,
			CDNName:           string(server.CDN),
			ServerID:          server.ID,
			ProfileName:       server.ProfileName,
		},
		ConfigFiles: []tc.ATSConfigMetaDataConfigFile{},
	}

	if locationParams["remap.config"].Location != "" {
		configLocation := locationParams["remap.config"].Location
		for _, ds := range uriSignedDSes {
			cfgName := "uri_signing_" + string(ds) + ".config"
			// If there's already a parameter for it, don't clobber it. The user may wish to override the location.
			if _, ok := locationParams[cfgName]; !ok {
				p := locationParams[cfgName]
				p.FileNameOnDisk = cfgName
				p.Location = configLocation
				locationParams[cfgName] = p
			}
		}
	}

	for cfgFile, cfgParams := range locationParams {
		atsCfg := tc.ATSConfigMetaDataConfigFile{
			FileNameOnDisk: cfgParams.FileNameOnDisk,
			Location:       cfgParams.Location,
		}

		scope := getServerScope(cfgFile, server.Type, scopeParams)

		if cfgParams.URL != "" {
			scope = tc.ATSConfigMetaDataConfigFileScopeCDNs
			atsCfg.URL = cfgParams.URL
		} else {
			scopeID := ""
			if scope == tc.ATSConfigMetaDataConfigFileScopeCDNs {
				scopeID = string(server.CDN)
			} else if scope == tc.ATSConfigMetaDataConfigFileScopeProfiles {
				scopeID = server.ProfileName
			} else { // ATSConfigMetaDataConfigFileScopeServers
				scopeID = server.HostName
			}
			atsCfg.APIURI = "/api/" + APIVersion + "/" + string(scope) + "/" + scopeID + "/configfiles/ats/" + cfgFile
		}

		atsCfg.Scope = string(scope)

		atsData.ConfigFiles = append(atsData.ConfigFiles, atsCfg)
	}

	bts, err := json.Marshal(atsData)
	if err != nil {
		// should never happen
		log.Errorln("marshalling chkconfig NameVersions: " + err.Error())
		bts = []byte("error encoding to json, see log for details")
	}
	return string(bts)
}

func getServerScope(cfgFile string, serverType string, scopeParams map[string]string) tc.ATSConfigMetaDataConfigFileScope {
	switch {
	case cfgFile == "cache.config" && tc.CacheTypeFromString(serverType) == tc.CacheTypeMid:
		return tc.ATSConfigMetaDataConfigFileScopeServers
	default:
		return getScope(cfgFile, scopeParams)
	}
}

const DefaultScope = tc.ATSConfigMetaDataConfigFileScopeServers

// getScope returns the ATSConfigMetaDataConfigFileScope for the given config file, and potentially the given server. If the config is not a Server scope, i.e. was part of an endpoint which does not include a server name or id, the server may be nil.
func getScope(cfgFile string, scopeParams map[string]string) tc.ATSConfigMetaDataConfigFileScope {
	switch {
	case cfgFile == "ip_allow.config",
		cfgFile == "parent.config",
		cfgFile == "hosting.config",
		cfgFile == "packages",
		cfgFile == "chkconfig",
		cfgFile == "remap.config",
		strings.HasPrefix(cfgFile, "to_ext_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeServers
	case cfgFile == "12M_facts",
		cfgFile == "50-ats.rules",
		cfgFile == "astats.config",
		cfgFile == "cache.config",
		cfgFile == "drop_qstring.config",
		cfgFile == "logs_xml.config",
		cfgFile == "logging.config",
		cfgFile == "plugin.config",
		cfgFile == "records.config",
		cfgFile == "storage.config",
		cfgFile == "volume.config",
		cfgFile == "sysctl.conf",
		strings.HasPrefix(cfgFile, "url_sig_") && strings.HasSuffix(cfgFile, ".config"),
		strings.HasPrefix(cfgFile, "uri_signing_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeProfiles
	case cfgFile == "bg_fetch.config",
		cfgFile == "regex_revalidate.config",
		cfgFile == "ssl_multicert.config",
		strings.HasPrefix(cfgFile, "cacheurl") && strings.HasSuffix(cfgFile, ".config"),
		strings.HasPrefix(cfgFile, "hdr_rw_") && strings.HasSuffix(cfgFile, ".config"),
		strings.HasPrefix(cfgFile, "regex_remap_") && strings.HasSuffix(cfgFile, ".config"),
		strings.HasPrefix(cfgFile, "set_dscp_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeCDNs
	}

	scope, ok := scopeParams[cfgFile]
	if !ok {
		scope = string(DefaultScope)
	}
	return tc.ATSConfigMetaDataConfigFileScope(scope)
}
