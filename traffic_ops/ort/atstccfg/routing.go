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
	"errors"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/cfgfile"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

var scopeConfigFileFuncs = map[string]func(toData *cfgfile.TOData, fileName string) (string, int, error){
	"cdns":     GetConfigFileCDN,
	"servers":  GetConfigFileServer,
	"profiles": GetConfigFileProfile,
}

func GetConfigFile(toData *cfgfile.TOData, fileInfo tc.ATSConfigMetaDataConfigFile) (string, int, error) {
	path := fileInfo.APIURI
	// TODO remove the URL path parsing. It's a legacy from when config files were endpoints in the meta config.
	// We should replace it with actually calling the right file and name directly.
	start := time.Now()
	defer func() {
		log.Infof("GetConfigFile %v took %v\n", path, time.Since(start).Round(time.Millisecond))
	}()

	pathParts := strings.Split(path, "/")
	if len(pathParts) < 8 {
		return "", 0, errors.New("unknown config file '" + path + "'")
	}
	scope := pathParts[3]
	resource := pathParts[4]
	fileName := pathParts[7]

	log.Infoln("GetConfigFile scope '" + scope + "' resource '" + resource + "' fileName '" + fileName + "'")

	if scopeConfigFileFunc, ok := scopeConfigFileFuncs[scope]; ok {
		return scopeConfigFileFunc(toData, fileName)
	}

	return "", 0, errors.New("unknown config file '" + fileInfo.APIURI + "'")
}

type ConfigFilePrefixSuffixFunc struct {
	Prefix string
	Suffix string
	Func   func(toData *cfgfile.TOData, fileName string) (string, error)
}

func GetConfigFileCDN(toData *cfgfile.TOData, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileCDN cdn '" + toData.Server.CDNName + "' fileName '" + fileName + "'")

	txt := ""
	err := error(nil)
	if getCfgFunc, ok := CDNConfigFileFuncs()[fileName]; ok {
		txt, err = getCfgFunc(toData)
	} else {
		for _, prefixSuffixFunc := range ConfigFileCDNPrefixSuffixFuncs {
			if strings.HasPrefix(fileName, prefixSuffixFunc.Prefix) && strings.HasSuffix(fileName, prefixSuffixFunc.Suffix) && len(fileName) > len(prefixSuffixFunc.Prefix)+len(prefixSuffixFunc.Suffix) {
				txt, err = prefixSuffixFunc.Func(toData, fileName)
				break
			}
		}
	}

	if err == nil && txt == "" {
		err = config.ErrNotFound
	}

	if err != nil {
		code := config.ExitCodeErrGeneric
		if err == config.ErrNotFound {
			code = config.ExitCodeNotFound
		} else if err == config.ErrBadRequest {
			code = config.ExitCodeBadRequest
		}
		return "", code, err
	}
	return txt, config.ExitCodeSuccess, nil
}

func GetConfigFileProfile(toData *cfgfile.TOData, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileProfile profile '" + toData.Server.Profile + "' fileName '" + fileName + "'")

	txt := ""
	err := error(nil)
	if getCfgFunc, ok := ProfileConfigFileFuncs()[fileName]; ok {
		txt, err = getCfgFunc(toData)
	} else if strings.HasPrefix(fileName, "url_sig_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("url_sig_")+len(".config") {
		txt, err = cfgfile.GetConfigFileProfileURLSigConfig(toData, fileName)
	} else if strings.HasPrefix(fileName, "uri_signing_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("uri_signing")+len(".config") {
		txt, err = cfgfile.GetConfigFileProfileURISigningConfig(toData, fileName)
	} else {
		txt, err = cfgfile.GetConfigFileProfileUnknownConfig(toData, fileName)
	}

	if err != nil {
		code := config.ExitCodeErrGeneric
		if err == config.ErrNotFound {
			code = config.ExitCodeNotFound
		} else if err == config.ErrBadRequest {
			code = config.ExitCodeBadRequest
		}
		return "", code, err
	}
	return txt, config.ExitCodeSuccess, nil
}

// ConfigFileFuncs returns a map[scope][configFile]configFileFunc.
func ConfigFileFuncs() map[string]map[string]func(toData *cfgfile.TOData) (string, error) {
	return map[string]map[string]func(toData *cfgfile.TOData) (string, error){
		"cdns":     CDNConfigFileFuncs(),
		"servers":  ServerConfigFileFuncs(),
		"profiles": ProfileConfigFileFuncs(),
	}
}

func CDNConfigFileFuncs() map[string]func(toData *cfgfile.TOData) (string, error) {
	return map[string]func(toData *cfgfile.TOData) (string, error){
		"regex_revalidate.config": cfgfile.GetConfigFileCDNRegexRevalidateDotConfig,
		"bg_fetch.config":         cfgfile.GetConfigFileCDNBGFetchDotConfig,
		"ssl_multicert.config":    cfgfile.GetConfigFileCDNSSLMultiCertDotConfig,
		"cacheurl.config":         cfgfile.GetConfigFileCDNCacheURLPlain,
	}
}

var ConfigFileCDNPrefixSuffixFuncs = []ConfigFilePrefixSuffixFunc{
	{"hdr_rw_mid_", ".config", cfgfile.GetConfigFileCDNHeaderRewriteMid},
	{"hdr_rw_", ".config", cfgfile.GetConfigFileCDNHeaderRewrite},
	{"cacheurl", ".config", cfgfile.GetConfigFileCDNCacheURL},
	{"regex_remap_", ".config", cfgfile.GetConfigFileCDNRegexRemap},
	{"set_dscp_", ".config", cfgfile.GetConfigFileCDNSetDSCP},
}

func ProfileConfigFileFuncs() map[string]func(toData *cfgfile.TOData) (string, error) {
	return map[string]func(toData *cfgfile.TOData) (string, error){
		"12M_facts":           cfgfile.GetConfigFileProfile12MFacts,
		"50-ats.rules":        cfgfile.GetConfigFileProfileATSDotRules,
		"astats.config":       cfgfile.GetConfigFileProfileAstatsDotConfig,
		"cache.config":        cfgfile.GetConfigFileProfileCacheDotConfig,
		"drop_qstring.config": cfgfile.GetConfigFileProfileDropQStringDotConfig,
		"logging.config":      cfgfile.GetConfigFileProfileLoggingDotConfig,
		"logging.yaml":        cfgfile.GetConfigFileProfileLoggingDotYAML,
		"logs_xml.config":     cfgfile.GetConfigFileProfileLogsXMLDotConfig,
		"plugin.config":       cfgfile.GetConfigFileProfilePluginDotConfig,
		"records.config":      cfgfile.GetConfigFileProfileRecordsDotConfig,
		"storage.config":      cfgfile.GetConfigFileProfileStorageDotConfig,
		"sysctl.conf":         cfgfile.GetConfigFileProfileSysCtlDotConf,
		"volume.config":       cfgfile.GetConfigFileProfileVolumeDotConfig,
	}
}

func ServerConfigFileFuncs() map[string]func(toData *cfgfile.TOData) (string, error) {
	return map[string]func(toData *cfgfile.TOData) (string, error){
		"parent.config":   cfgfile.GetConfigFileServerParentDotConfig,
		"remap.config":    cfgfile.GetConfigFileServerRemapDotConfig,
		"cache.config":    cfgfile.GetConfigFileServerCacheDotConfig,
		"ip_allow.config": cfgfile.GetConfigFileServerIPAllowDotConfig,
		"hosting.config":  cfgfile.GetConfigFileServerHostingDotConfig,
		"packages":        cfgfile.GetConfigFileServerPackages,
		"chkconfig":       cfgfile.GetConfigFileServerChkconfig,
	}
}

func GetConfigFileServer(toData *cfgfile.TOData, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileServer server '" + toData.Server.HostName + "' fileName '" + fileName + "'")
	txt := ""
	err := error(nil)
	if getCfgFunc, ok := ServerConfigFileFuncs()[fileName]; ok {
		txt, err = getCfgFunc(toData)
	} else {
		txt, err = cfgfile.GetConfigFileServerUnknownConfig(toData, fileName)
	}
	if err != nil {
		return "", config.ExitCodeErrGeneric, err
	}
	return txt, config.ExitCodeSuccess, nil
}
