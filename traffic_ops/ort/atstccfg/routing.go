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
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

var scopeConfigFileFuncs = map[string]func(toClient **toclient.Session, cfg Cfg, resource string, fileName string) (string, int, error){
	"cdns":     GetConfigFileCDN,
	"servers":  GetConfigFileServer,
	"profiles": GetConfigFileProfile,
}

func GetConfigFile(toClient **toclient.Session, cfg Cfg) (string, int, error) {

	pathParts := strings.Split(cfg.TOURL.Path, "/")

	log.Infof("GetConfigFile pathParts %++v\n", pathParts)

	if len(pathParts) < 8 {
		log.Infoln("GetConfigFile pathParts < 7, calling TO")
		return GetConfigFileFromTrafficOps(toClient, cfg)
	}
	scope := pathParts[3]
	resource := pathParts[4]
	fileName := pathParts[7]

	log.Infoln("GetConfigFile scope '" + scope + "' resource '" + resource + "' fileName '" + fileName + "'")

	if scopeConfigFileFunc, ok := scopeConfigFileFuncs[scope]; ok {
		return scopeConfigFileFunc(toClient, cfg, resource, fileName)
	}

	log.Infoln("GetConfigFile unknown scope, calling TO")
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

func GetConfigFileCDN(toClient **toclient.Session, cfg Cfg, cdnNameOrID string, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileCDN cdn '" + cdnNameOrID + "' fileName '" + fileName + "'")
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

func GetConfigFileProfile(toClient **toclient.Session, cfg Cfg, profileNameOrID string, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileProfile profile '" + profileNameOrID + "' fileName '" + fileName + "'")
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

// ConfigFileFuncs returns a map[scope][configFile]configFileFunc.
func ConfigFileFuncs() map[string]map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error) {
	return map[string]map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error){
		"cdns":     CDNConfigFileFuncs(),
		"servers":  ServerConfigFileFuncs(),
		"profiles": ProfileConfigFileFuncs(),
	}
}

func CDNConfigFileFuncs() map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error) {
	return map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error){}
}

func ProfileConfigFileFuncs() map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error) {
	return map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error){}
}

func ServerConfigFileFuncs() map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error) {
	return map[string]func(toClient **toclient.Session, cfg Cfg, serverNameOrID string) (string, error){
		"parent.config": GetConfigFileServerParentDotConfig,
	}
}

func GetConfigFileServer(toClient **toclient.Session, cfg Cfg, serverNameOrID string, fileName string) (string, int, error) {
	log.Infoln("GetConfigFileServer server '" + serverNameOrID + "' fileName '" + fileName + "'")
	if getCfgFunc, ok := ServerConfigFileFuncs()[fileName]; ok {
		txt, err := getCfgFunc(toClient, cfg, serverNameOrID)
		if err != nil {
			return "", 1, err
		}
		return txt, 0, nil
	}
	return GetConfigFileFromTrafficOps(toClient, cfg)
}

func GetConfigFileFromTrafficOps(toClient **toclient.Session, cfg Cfg) (string, int, error) {
	path := cfg.TOURL.Path
	if cfg.TOURL.RawQuery != "" {
		path += "?" + cfg.TOURL.RawQuery
	}
	log.Infoln("GetConfigFile path '" + path + "' not generated locally, requesting from Traffic Ops")
	log.Infoln("GetConfigFile url '" + cfg.TOURL.String() + "'")

	body, code, err := TrafficOpsRequest(toClient, cfg, http.MethodGet, cfg.TOURL.String(), nil)
	if err != nil {
		return "", code, errors.New("Requesting path '" + path + "': " + err.Error())
	}

	WriteCookiesToFile(CookiesToString((*toClient).Client.Jar.Cookies(cfg.TOURL)), cfg.TempDir)

	return string(body), code, nil
}
