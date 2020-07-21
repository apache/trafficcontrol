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
	"errors"
	"path/filepath"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

type ConfigProfileParams struct {
	FileNameOnDisk string
	Location       string
	URL            string
}

// APIVersion is the Traffic Ops API version for config fiels.
// This is used to generate the meta config, which has API paths.
// Note the version in the meta config is not used by the atstccfg generator, which isn't actually an API.
// TODO change the config system to not use old API paths, and remove this.
const APIVersion = "2.0"

// requiredFiles is a constant (because Go doesn't allow const slices).
// Note these are not exhaustive. This is only used to error if these are missing.
// The presence of these is no guarantee the location Parameters are complete and correct.
func requiredFiles() []string {
	return []string{
		"cache.config",
		"hosting.config",
		"ip_allow.config",
		"parent.config",
		"plugin.config",
		"records.config",
		"remap.config",
		"storage.config",
		"volume.config",
	}
}

func MakeMetaConfig(
	serverHostName tc.CacheName,
	server *ServerInfo,
	tmURL string, // global tm.url Parameter
	tmReverseProxyURL string, // global tm.rev_proxy.url Parameter
	locationParams map[string]ConfigProfileParams, // map[configFile]params; 'location' and 'URL' Parameters on serverHostName's Profile
	uriSignedDSes []tc.DeliveryServiceName,
	scopeParams map[string]string, // map[configFileName]scopeParam
	dsNames map[tc.DeliveryServiceName]struct{},
) string {
	// dses are only used when configDir is not empty
	dses := map[tc.DeliveryServiceName]tc.DeliveryServiceNullable{}
	for dsName, _ := range dsNames {
		dses[dsName] = tc.DeliveryServiceNullable{}
	}
	configDir := ""
	atsData, err := MakeMetaObj(serverHostName, server, tmURL, tmReverseProxyURL, locationParams, uriSignedDSes, scopeParams, dses, configDir)
	if err != nil {
		return "error creating meta config: " + err.Error()
	}
	bts, err := json.Marshal(atsData)
	if err != nil {
		// should never happen
		log.Errorln("marshalling meta config: " + err.Error())
		bts = []byte("error encoding to json, see log for details")
	}
	return string(bts)
}

// AddMetaObjConfigDir takes the Meta Object generated from TO data, and the ATS config directory
// and prepends the config directory to all relative paths,
// and creates MetaObj entries for all required config files which have no location parameter.
// If configDir is empty and any location Parameters have relative paths, returns an error.
func AddMetaObjConfigDir(
	metaObj tc.ATSConfigMetaData,
	configDir string,
	serverHostName tc.CacheName,
	server *ServerInfo,
	tmURL string, // global tm.url Parameter
	tmReverseProxyURL string, // global tm.rev_proxy.url Parameter
	locationParams map[string]ConfigProfileParams, // map[configFile]params; 'location' and 'URL' Parameters on serverHostName's Profile
	uriSignedDSes []tc.DeliveryServiceName,
	scopeParams map[string]string, // map[configFileName]scopeParam
	dses map[tc.DeliveryServiceName]tc.DeliveryServiceNullable,
) (tc.ATSConfigMetaData, error) {

	// Note there may be multiple files with the same name in different directories.
	configFilesM := map[string][]tc.ATSConfigMetaDataConfigFile{} // map[fileShortName]tc.ATSConfigMetaDataConfigFile
	for _, fi := range metaObj.ConfigFiles {
		configFilesM[fi.FileNameOnDisk] = append(configFilesM[fi.FileNameOnDisk], fi)
	}

	// add all strictly required files, all of which should be in the base config directory.
	// If they don't exist, create them.
	// If they exist with a relative path, prepend configDir.
	// If any exist with a relative path, or don't exist, and configDir is empty, return an error.
	for _, fileName := range requiredFiles() {
		if _, ok := configFilesM[fileName]; ok {
			continue
		}
		if configDir == "" {
			return metaObj, errors.New("required file '" + fileName + "' has no location Parameter, and ATS config directory not found.")
		}
		configFilesM[fileName] = []tc.ATSConfigMetaDataConfigFile{{
			FileNameOnDisk: fileName,
			Location:       configDir,
			Scope:          string(getServerScope(fileName, server.Type, nil)),
		}}
	}

	for fileName, fis := range configFilesM {
		newFis := []tc.ATSConfigMetaDataConfigFile{}
		for _, fi := range fis {
			if !filepath.IsAbs(fi.Location) {
				if configDir == "" {
					return metaObj, errors.New("file '" + fileName + "' has location Parameter with relative path '" + fi.Location + "', but ATS config directory was not found.")
				}
				absPath := filepath.Join(configDir, fi.Location)
				fi.Location = absPath
			}
			newFis = append(newFis, fi)
		}
		configFilesM[fileName] = newFis
	}

	for _, ds := range dses {
		if ds.XMLID == nil {
			log.Errorln("meta config generation got Delivery Service with nil XMLID - not considering!")
			continue
		}
		err := error(nil)
		// Note we log errors, but don't return them.
		// If an individual DS has an error, we don't want to break the rest of the CDN.
		if (ds.EdgeHeaderRewrite != nil || ds.MaxOriginConnections != nil) &&
			strings.HasPrefix(server.Type, tc.EdgeTypePrefix) {
			fileName := "hdr_rw_" + *ds.XMLID + ".config"
			scope := tc.ATSConfigMetaDataConfigFileScopeCDNs
			if configFilesM, err = ensureConfigFile(configFilesM, fileName, configDir, scope); err != nil {
				log.Errorln("meta config generation: " + err.Error())
			}
		}
		if (ds.MidHeaderRewrite != nil || ds.MaxOriginConnections != nil) &&
			ds.Type != nil && ds.Type.UsesMidCache() &&
			strings.HasPrefix(server.Type, tc.MidTypePrefix) {
			fileName := "hdr_rw_mid_" + *ds.XMLID + ".config"
			scope := tc.ATSConfigMetaDataConfigFileScopeCDNs
			if configFilesM, err = ensureConfigFile(configFilesM, fileName, configDir, scope); err != nil {
				log.Errorln("meta config generation: " + err.Error())
			}
		}
		if ds.RegexRemap != nil {
			configFile := "regex_remap_" + *ds.XMLID + ".config"
			scope := tc.ATSConfigMetaDataConfigFileScopeCDNs
			if configFilesM, err = ensureConfigFile(configFilesM, configFile, configDir, scope); err != nil {
				log.Errorln("meta config generation: " + err.Error())
			}
		}
		if ds.CacheURL != nil {
			configFile := "cacheurl_" + *ds.XMLID + ".config"
			scope := tc.ATSConfigMetaDataConfigFileScopeCDNs
			if configFilesM, err = ensureConfigFile(configFilesM, configFile, configDir, scope); err != nil {
				log.Errorln("meta config generation: " + err.Error())
			}
		}
		if ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == tc.SigningAlgorithmURLSig {
			configFile := "url_sig_" + *ds.XMLID + ".config"
			scope := tc.ATSConfigMetaDataConfigFileScopeProfiles
			if configFilesM, err = ensureConfigFile(configFilesM, configFile, configDir, scope); err != nil {
				log.Errorln("meta config generation: " + err.Error())
			}
		}
		if ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == tc.SigningAlgorithmURISigning {
			configFile := "uri_signing_" + *ds.XMLID + ".config"
			scope := tc.ATSConfigMetaDataConfigFileScopeProfiles
			if configFilesM, err = ensureConfigFile(configFilesM, configFile, configDir, scope); err != nil {
				log.Errorln("meta config generation: " + err.Error())
			}
		}
	}
	// TODO add location params for ds ensure garbage here

	newFiles := []tc.ATSConfigMetaDataConfigFile{}
	for _, fis := range configFilesM {
		for _, fi := range fis {
			newFiles = append(newFiles, fi)
		}
	}
	metaObj.ConfigFiles = newFiles
	return metaObj, nil
}

// ensureConfigFile ensures files contains the given fileName. If so, returns files unmodified.
// If not, if configDir is empty, returns an error.
// If not, and configDir is nonempty, creates the given file, configDir location, and scope, and returns files.
func ensureConfigFile(files map[string][]tc.ATSConfigMetaDataConfigFile, fileName string, configDir string, scope tc.ATSConfigMetaDataConfigFileScope) (map[string][]tc.ATSConfigMetaDataConfigFile, error) {
	if _, ok := files[fileName]; ok {
		return files, nil
	}
	if configDir == "" {
		return files, errors.New("required file '" + fileName + "' has no location Parameter, and ATS config directory not found.")
	}
	files[fileName] = []tc.ATSConfigMetaDataConfigFile{{
		FileNameOnDisk: fileName,
		Location:       configDir,
		Scope:          string(scope),
	}}
	return files, nil
}

func MakeMetaObj(
	serverHostName tc.CacheName,
	server *ServerInfo,
	tmURL string, // global tm.url Parameter
	tmReverseProxyURL string, // global tm.rev_proxy.url Parameter
	locationParams map[string]ConfigProfileParams, // map[configFile]params; 'location' and 'URL' Parameters on serverHostName's Profile
	uriSignedDSes []tc.DeliveryServiceName,
	scopeParams map[string]string, // map[configFileName]scopeParam
	dses map[tc.DeliveryServiceName]tc.DeliveryServiceNullable,
	configDir string,
) (tc.ATSConfigMetaData, error) {
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

locationParamsFor:
	for cfgFile, cfgParams := range locationParams {
		if strings.HasSuffix(cfgFile, ".config") {
			dsConfigFilePrefixes := []string{
				"hdr_rw_mid_", // must come before hdr_rw_, to avoid thinking we have a "hdr_rw_" with a ds of "mid_x"
				"hdr_rw_",
				"regex_remap_",
				"url_sig_",
				"uri_signing_",
			}
		prefixFor:
			for _, prefix := range dsConfigFilePrefixes {
				if strings.HasPrefix(cfgFile, prefix) {
					dsName := strings.TrimSuffix(strings.TrimPrefix(cfgFile, prefix), ".config")
					if _, ok := dses[tc.DeliveryServiceName(dsName)]; !ok {
						log.Warnln("Server Profile had 'location' Parameter '" + cfgFile + "', but delivery Service '" + dsName + "' is not assigned to this Server! Not including in meta config!")
						continue locationParamsFor
					}
					break prefixFor // if it has a prefix, don't check the next prefix. This is important for hdr_rw_mid_, which will match hdr_rw_ and result in a "ds name" of "mid_x" if we don't continue here.
				}
			}
		}

		atsCfg := tc.ATSConfigMetaDataConfigFile{
			FileNameOnDisk: cfgParams.FileNameOnDisk,
			Location:       cfgParams.Location,
		}

		scope := getServerScope(cfgFile, server.Type, scopeParams)

		if cfgParams.URL != "" {
			// TODO this is legacy, from when a custom URL could be set via Parameters.
			//      verify nobody is relying on it in a production system, and remove.
			scope = tc.ATSConfigMetaDataConfigFileScopeCDNs
		}

		atsCfg.Scope = string(scope)

		atsData.ConfigFiles = append(atsData.ConfigFiles, atsCfg)
	}

	return AddMetaObjConfigDir(atsData, configDir, serverHostName, server, tmURL, tmReverseProxyURL, locationParams, uriSignedDSes, scopeParams, dses)
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
		cfgFile == "logging.yaml",
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
		cfgFile == SSLMultiCertConfigFileName,
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
