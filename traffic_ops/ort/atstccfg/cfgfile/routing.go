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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

var scopeConfigFileFuncs = map[string]func(toData *config.TOData, fileName string) (string, string, error){
	"cdns":     GetConfigFileCDN,
	"servers":  GetConfigFileServer,
	"profiles": GetConfigFileProfile,
}

// GetConfigFile returns the text of the generated config file, the MIME Content Type of the config file, and any error.
func GetConfigFile(toData *config.TOData, fileInfo tc.ATSConfigMetaDataConfigFile) (string, string, error) {
	path := fileInfo.APIURI
	// TODO remove the URL path parsing. It's a legacy from when config files were endpoints in the meta config.
	// We should replace it with actually calling the right file and name directly.
	start := time.Now()
	defer func() {
		log.Infof("GetConfigFile %v took %v\n", path, time.Since(start).Round(time.Millisecond))
	}()

	pathParts := strings.Split(path, "/")
	if len(pathParts) < 8 {
		return "", "", errors.New("unknown config file '" + path + "'")
	}
	scope := pathParts[3]
	resource := pathParts[4]
	fileName := pathParts[7]

	log.Infoln("GetConfigFile scope '" + scope + "' resource '" + resource + "' fileName '" + fileName + "'")

	if scopeConfigFileFunc, ok := scopeConfigFileFuncs[scope]; ok {
		return scopeConfigFileFunc(toData, fileName)
	}

	return "", "", errors.New("unknown config file '" + fileInfo.APIURI + "'")
}

type ConfigFilePrefixSuffixFunc struct {
	Prefix string
	Suffix string
	Func   func(toData *config.TOData, fileName string) (string, string, error)
}

func GetConfigFileCDN(toData *config.TOData, fileName string) (string, string, error) {
	log.Infoln("GetConfigFileCDN cdn '" + toData.Server.CDNName + "' fileName '" + fileName + "'")

	txt := ""
	contentType := ""
	err := error(nil)
	if getCfgFunc, ok := CDNConfigFileFuncs()[fileName]; ok {
		txt, contentType, err = getCfgFunc(toData)
	} else {
		for _, prefixSuffixFunc := range ConfigFileCDNPrefixSuffixFuncs {
			if strings.HasPrefix(fileName, prefixSuffixFunc.Prefix) && strings.HasSuffix(fileName, prefixSuffixFunc.Suffix) && len(fileName) > len(prefixSuffixFunc.Prefix)+len(prefixSuffixFunc.Suffix) {
				txt, contentType, err = prefixSuffixFunc.Func(toData, fileName)
				break
			}
		}
	}

	if err == nil && txt == "" {
		err = config.ErrNotFound
	}

	if err != nil {
		return "", "", err
	}
	return txt, contentType, nil
}

func GetConfigFileProfile(toData *config.TOData, fileName string) (string, string, error) {
	log.Infoln("GetConfigFileProfile profile '" + toData.Server.Profile + "' fileName '" + fileName + "'")

	txt := ""
	contentType := ""
	err := error(nil)
	if getCfgFunc, ok := ProfileConfigFileFuncs()[fileName]; ok {
		txt, contentType, err = getCfgFunc(toData)
	} else if strings.HasPrefix(fileName, "url_sig_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("url_sig_")+len(".config") {
		txt, contentType, err = GetConfigFileProfileURLSigConfig(toData, fileName)
	} else if strings.HasPrefix(fileName, "uri_signing_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("uri_signing")+len(".config") {
		txt, contentType, err = GetConfigFileProfileURISigningConfig(toData, fileName)
	} else {
		txt, contentType, err = GetConfigFileProfileUnknownConfig(toData, fileName)
	}

	if err != nil {
		return "", "", err
	}
	return txt, contentType, nil
}

// ConfigFileFuncs returns a map[scope][configFile]configFileFunc.
func ConfigFileFuncs() map[string]map[string]func(toData *config.TOData) (string, string, error) {
	return map[string]map[string]func(toData *config.TOData) (string, string, error){
		"cdns":     CDNConfigFileFuncs(),
		"servers":  ServerConfigFileFuncs(),
		"profiles": ProfileConfigFileFuncs(),
	}
}

func CDNConfigFileFuncs() map[string]func(toData *config.TOData) (string, string, error) {
	return map[string]func(toData *config.TOData) (string, string, error){
		"regex_revalidate.config": GetConfigFileCDNRegexRevalidateDotConfig,
		"bg_fetch.config":         GetConfigFileCDNBGFetchDotConfig,
		"ssl_multicert.config":    GetConfigFileCDNSSLMultiCertDotConfig,
		"cacheurl.config":         GetConfigFileCDNCacheURLPlain,
	}
}

var ConfigFileCDNPrefixSuffixFuncs = []ConfigFilePrefixSuffixFunc{
	{"hdr_rw_mid_", ".config", GetConfigFileCDNHeaderRewriteMid},
	{"hdr_rw_", ".config", GetConfigFileCDNHeaderRewrite},
	{"cacheurl", ".config", GetConfigFileCDNCacheURL},
	{"regex_remap_", ".config", GetConfigFileCDNRegexRemap},
	{"set_dscp_", ".config", GetConfigFileCDNSetDSCP},
}

func ProfileConfigFileFuncs() map[string]func(toData *config.TOData) (string, string, error) {
	return map[string]func(toData *config.TOData) (string, string, error){
		"12M_facts":           GetConfigFileProfile12MFacts,
		"50-ats.rules":        GetConfigFileProfileATSDotRules,
		"astats.config":       GetConfigFileProfileAstatsDotConfig,
		"cache.config":        GetConfigFileProfileCacheDotConfig,
		"drop_qstring.config": GetConfigFileProfileDropQStringDotConfig,
		"logging.config":      GetConfigFileProfileLoggingDotConfig,
		"logging.yaml":        GetConfigFileProfileLoggingDotYAML,
		"logs_xml.config":     GetConfigFileProfileLogsXMLDotConfig,
		"plugin.config":       GetConfigFileProfilePluginDotConfig,
		"records.config":      GetConfigFileProfileRecordsDotConfig,
		"storage.config":      GetConfigFileProfileStorageDotConfig,
		"sysctl.conf":         GetConfigFileProfileSysCtlDotConf,
		"volume.config":       GetConfigFileProfileVolumeDotConfig,
	}
}

func ServerConfigFileFuncs() map[string]func(toData *config.TOData) (string, string, error) {
	return map[string]func(toData *config.TOData) (string, string, error){
		"parent.config":   GetConfigFileServerParentDotConfig,
		"remap.config":    GetConfigFileServerRemapDotConfig,
		"cache.config":    GetConfigFileServerCacheDotConfig,
		"ip_allow.config": GetConfigFileServerIPAllowDotConfig,
		"hosting.config":  GetConfigFileServerHostingDotConfig,
		"packages":        GetConfigFileServerPackages,
		"chkconfig":       GetConfigFileServerChkconfig,
	}
}

func GetConfigFileServer(toData *config.TOData, fileName string) (string, string, error) {
	log.Infoln("GetConfigFileServer server '" + toData.Server.HostName + "' fileName '" + fileName + "'")
	txt := ""
	contentType := ""
	err := error(nil)
	if getCfgFunc, ok := ServerConfigFileFuncs()[fileName]; ok {
		txt, contentType, err = getCfgFunc(toData)
	} else {
		txt, contentType, err = GetConfigFileServerUnknownConfig(toData, fileName)
	}
	if err != nil {
		return "", "", err
	}
	return txt, contentType, nil
}
