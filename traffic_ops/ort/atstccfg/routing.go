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
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/cfgfile"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/toreq"
)

var scopeConfigFileFuncs = map[string]func(cfg config.TCCfg, resource string, fileName string) (string, int, error){
	"cdns":     GetConfigFileCDN,
	"servers":  GetConfigFileServer,
	"profiles": GetConfigFileProfile,
}

func GetConfigFile(cfg config.TCCfg) (string, int, error) {
	pathParts := strings.Split(cfg.TOURL.Path, "/")

	if len(pathParts) == 7 && pathParts[1] == `api` && pathParts[3] == `servers` && pathParts[5] == `configfiles` && pathParts[6] == `ats` {
		// "/api/1.x/servers/name/configfiles/ats" is the "meta" config route, which lists all the other configs for this server.
		server := pathParts[4]
		log.Infoln("GetConfigFile is meta config request for server '" + server + "'; generating")
		txt, err := cfgfile.GetConfigFileMeta(cfg, server)
		if err != nil {
			if err == config.ErrNotFound {
				return "", config.ExitCodeNotFound, err
			} else if err == config.ErrBadRequest {
				return "", config.ExitCodeBadRequest, err
			} else {
				return "", config.ExitCodeErrGeneric, err
			}
		}
		return txt, config.ExitCodeSuccess, nil
	}

	if len(pathParts) < 8 {
		log.Infoln("GetConfigFile pathParts < 7, calling TO")
		return GetConfigFileFromTrafficOps(cfg)
	}
	scope := pathParts[3]
	resource := pathParts[4]
	fileName := pathParts[7]

	log.Infoln("GetConfigFile scope '" + scope + "' resource '" + resource + "' fileName '" + fileName + "'")

	if scopeConfigFileFunc, ok := scopeConfigFileFuncs[scope]; ok {
		return scopeConfigFileFunc(cfg, resource, fileName)
	}

	log.Infoln("GetConfigFile unknown scope, calling TO")
	return GetConfigFileFromTrafficOps(cfg)
}

type ConfigFilePrefixSuffixFunc struct {
	Prefix string
	Suffix string
	Func   func(cfg config.TCCfg, resource string, fileName string) (string, error)
}

func GetConfigFileCDN(cfg config.TCCfg, cdnNameOrID string, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileCDN cdn '" + cdnNameOrID + "' fileName '" + fileName + "'")

	txt := ""
	err := error(nil)
	if getCfgFunc, ok := CDNConfigFileFuncs()[fileName]; ok {
		txt, err = getCfgFunc(cfg, cdnNameOrID)
	} else {
		for _, prefixSuffixFunc := range ConfigFileCDNPrefixSuffixFuncs {
			if strings.HasPrefix(fileName, prefixSuffixFunc.Prefix) && strings.HasSuffix(fileName, prefixSuffixFunc.Suffix) && len(fileName) > len(prefixSuffixFunc.Prefix)+len(prefixSuffixFunc.Suffix) {
				txt, err = prefixSuffixFunc.Func(cfg, cdnNameOrID, fileName)
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

func GetConfigFileProfile(cfg config.TCCfg, profileNameOrID string, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileProfile profile '" + profileNameOrID + "' fileName '" + fileName + "'")

	txt := ""
	err := error(nil)
	if getCfgFunc, ok := ProfileConfigFileFuncs()[fileName]; ok {
		txt, err = getCfgFunc(cfg, profileNameOrID)
	} else if strings.HasPrefix(fileName, "url_sig_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("url_sig_")+len(".config") {
		txt, err = cfgfile.GetConfigFileProfileURLSigConfig(cfg, profileNameOrID, fileName)
	} else if strings.HasPrefix(fileName, "uri_signing_") && strings.HasSuffix(fileName, ".config") && len(fileName) > len("uri_signing")+len(".config") {
		txt, err = cfgfile.GetConfigFileProfileURISigningConfig(cfg, profileNameOrID, fileName)
	} else {
		txt, err = cfgfile.GetConfigFileProfileUnknownConfig(cfg, profileNameOrID, fileName)
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
func ConfigFileFuncs() map[string]map[string]func(cfg config.TCCfg, serverNameOrID string) (string, error) {
	return map[string]map[string]func(cfg config.TCCfg, serverNameOrID string) (string, error){
		"cdns":     CDNConfigFileFuncs(),
		"servers":  ServerConfigFileFuncs(),
		"profiles": ProfileConfigFileFuncs(),
	}
}

func CDNConfigFileFuncs() map[string]func(cfg config.TCCfg, cdnNameOrID string) (string, error) {
	return map[string]func(cfg config.TCCfg, cdnNameOrID string) (string, error){
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

func ProfileConfigFileFuncs() map[string]func(cfg config.TCCfg, serverNameOrID string) (string, error) {
	return map[string]func(cfg config.TCCfg, serverNameOrID string) (string, error){
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

func ServerConfigFileFuncs() map[string]func(cfg config.TCCfg, serverNameOrID string) (string, error) {
	return map[string]func(cfg config.TCCfg, serverNameOrID string) (string, error){
		"parent.config":   cfgfile.GetConfigFileServerParentDotConfig,
		"remap.config":    cfgfile.GetConfigFileServerRemapDotConfig,
		"cache.config":    cfgfile.GetConfigFileServerCacheDotConfig,
		"ip_allow.config": cfgfile.GetConfigFileServerIPAllowDotConfig,
		"hosting.config":  cfgfile.GetConfigFileServerHostingDotConfig,
		"packages":        cfgfile.GetConfigFileServerPackages,
		"chkconfig":       cfgfile.GetConfigFileServerChkconfig,
	}
}

func GetConfigFileServer(cfg config.TCCfg, serverNameOrID string, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileServer server '" + serverNameOrID + "' fileName '" + fileName + "'")
	txt := ""
	err := error(nil)
	if getCfgFunc, ok := ServerConfigFileFuncs()[fileName]; ok {
		txt, err = getCfgFunc(cfg, serverNameOrID)
	} else {
		txt, err = cfgfile.GetConfigFileServerUnknownConfig(cfg, serverNameOrID, fileName)
	}
	if err != nil {
		return "", config.ExitCodeErrGeneric, err
	}
	return txt, config.ExitCodeSuccess, nil
}

func GetConfigFileFromTrafficOps(cfg config.TCCfg) (string, int, error) {
	path := cfg.TOURL.Path
	if cfg.TOURL.RawQuery != "" {
		path += "?" + cfg.TOURL.RawQuery
	}
	log.Infoln("GetConfigFile path '" + path + "' not generated locally, requesting from Traffic Ops")
	log.Infoln("GetConfigFile url '" + cfg.TOURL.String() + "'")

	body, code, err := toreq.TrafficOpsRequest(cfg, http.MethodGet, cfg.TOURL.String(), nil)
	if err != nil {
		return "", code, errors.New("Requesting path '" + path + "': " + err.Error())
	}

	toreq.WriteCookiesToFile(toreq.CookiesToString((*cfg.TOClient).Client.Jar.Cookies(cfg.TOURL)), cfg.TempDir)

	return string(body), HTTPCodeToExitCode(code), nil
}
