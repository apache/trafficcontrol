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
	"github.com/apache/trafficcontrol/traffic_ops_ort/atstccfg/config"
)

var scopeConfigFileFuncs = map[string]func(toData *config.TOData, fileName string) (string, string, string, error){
	"cdns":     GetConfigFileCDN,
	"servers":  GetConfigFileServer,
	"profiles": GetConfigFileProfile,
}

// GetConfigFile returns the text of the generated config file, the MIME Content Type of the config file, and any error.
func GetConfigFile(toData *config.TOData, fileInfo tc.ATSConfigMetaDataConfigFile) (string, string, string, error) {
	start := time.Now()
	defer func() {
		log.Infof("GetConfigFile %v took %v\n", fileInfo.FileNameOnDisk, time.Since(start).Round(time.Millisecond))
	}()
	log.Infoln("GetConfigFile scope '" + fileInfo.Scope + "' fileName '" + fileInfo.FileNameOnDisk + "'")
	if scopeConfigFileFunc, ok := scopeConfigFileFuncs[fileInfo.Scope]; ok {
		return scopeConfigFileFunc(toData, fileInfo.FileNameOnDisk)
	}
	return "", "", "", errors.New("unknown config file '" + fileInfo.FileNameOnDisk + "'")
}

type ConfigFilePrefixSuffixFunc struct {
	Prefix string
	Suffix string
	Func   func(toData *config.TOData, fileName string) (string, string, string, error)
}

func GetConfigFileCDN(toData *config.TOData, fileName string) (string, string, string, error) {
	log.Infoln("GetConfigFileCDN cdn '" + toData.Server.CDNName + "' fileName '" + fileName + "'")

	txt := ""
	contentType := ""
	lineComment := ""
	err := error(nil)
	if getCfgFunc, ok := CDNConfigFileFuncs()[fileName]; ok {
		txt, contentType, lineComment, err = getCfgFunc(toData)
	} else {
		for _, prefixSuffixFunc := range ConfigFileCDNPrefixSuffixFuncs {
			if strings.HasPrefix(fileName, prefixSuffixFunc.Prefix) && strings.HasSuffix(fileName, prefixSuffixFunc.Suffix) && len(fileName) > len(prefixSuffixFunc.Prefix)+len(prefixSuffixFunc.Suffix) {
				txt, contentType, lineComment, err = prefixSuffixFunc.Func(toData, fileName)
				break
			}
		}
	}

	if err == nil && txt == "" {
		err = config.ErrNotFound
	}

	if err != nil {
		return "", "", "", err
	}
	return txt, contentType, lineComment, nil
}

func GetConfigFileProfile(toData *config.TOData, fileName string) (string, string, string, error) {
	log.Infoln("GetConfigFileProfile profile '" + toData.Server.Profile + "' fileName '" + fileName + "'")

	txt := ""
	contentType := ""
	lineComment := ""
	err := error(nil)
	if getCfgFunc, ok := ProfileConfigFileFuncs()[fileName]; ok {
		txt, contentType, lineComment, err = getCfgFunc(toData)
	} else if strings.HasPrefix(fileName, "url_sig_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("url_sig_")+len(".config") {
		txt, contentType, lineComment, err = GetConfigFileProfileURLSigConfig(toData, fileName)
	} else if strings.HasPrefix(fileName, "uri_signing_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("uri_signing")+len(".config") {
		txt, contentType, lineComment, err = GetConfigFileProfileURISigningConfig(toData, fileName)
	} else {
		txt, contentType, lineComment, err = GetConfigFileProfileUnknownConfig(toData, fileName)
	}

	if err != nil {
		return "", "", "", err
	}
	return txt, contentType, lineComment, nil
}

// ConfigFileFuncs returns a map[scope][configFile]configFileFunc.
func ConfigFileFuncs() map[string]map[string]func(toData *config.TOData) (string, string, string, error) {
	return map[string]map[string]func(toData *config.TOData) (string, string, string, error){
		"cdns":     CDNConfigFileFuncs(),
		"servers":  ServerConfigFileFuncs(),
		"profiles": ProfileConfigFileFuncs(),
	}
}

func CDNConfigFileFuncs() map[string]func(toData *config.TOData) (string, string, string, error) {
	return map[string]func(toData *config.TOData) (string, string, string, error){
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

func ProfileConfigFileFuncs() map[string]func(toData *config.TOData) (string, string, string, error) {
	return map[string]func(toData *config.TOData) (string, string, string, error){
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

func ServerConfigFileFuncs() map[string]func(toData *config.TOData) (string, string, string, error) {
	return map[string]func(toData *config.TOData) (string, string, string, error){
		"parent.config":   GetConfigFileServerParentDotConfig,
		"remap.config":    GetConfigFileServerRemapDotConfig,
		"cache.config":    GetConfigFileServerCacheDotConfig,
		"ip_allow.config": GetConfigFileServerIPAllowDotConfig,
		"hosting.config":  GetConfigFileServerHostingDotConfig,
		"packages":        GetConfigFileServerPackages,
		"chkconfig":       GetConfigFileServerChkconfig,
	}
}

func GetConfigFileServer(toData *config.TOData, fileName string) (string, string, string, error) {
	log.Infoln("GetConfigFileServer server '" + toData.Server.HostName + "' fileName '" + fileName + "'")
	txt := ""
	contentType := ""
	lineComment := ""
	err := error(nil)
	if getCfgFunc, ok := ServerConfigFileFuncs()[fileName]; ok {
		txt, contentType, lineComment, err = getCfgFunc(toData)
	} else {
		txt, contentType, lineComment, err = GetConfigFileServerUnknownConfig(toData, fileName)
	}
	if err != nil {
		return "", "", "", err
	}
	return txt, contentType, lineComment, nil
}
