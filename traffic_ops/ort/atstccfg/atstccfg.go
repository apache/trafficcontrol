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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

const AppName = "atstccfg"
const Version = "0.1"
const UserAgent = AppName + "/" + Version

const APIVersion = "1.2"
const TempSubdir = AppName + "_cache"
const TempCookieFileName = "cookies"
const TOCookieName = "mojolicious"

const ExitCodeSuccess = 0
const ExitCodeErrGeneric = 1
const ExitCodeNotFound = 104
const ExitCodeBadRequest = 100

type TCCfg struct {
	Cfg
	TOClient **toclient.Session
}

func main() {
	cfg, err := GetCfg()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting config: "+err.Error()+"\n")
		os.Exit(ExitCodeErrGeneric)
	}

	if cfg.PrintGeneratedFiles {
		fmt.Println(strings.Join(GetGeneratedFilesList(), "\n"))
		os.Exit(ExitCodeSuccess)
	}

	log.Infoln("URL: '" + cfg.TOURL.String() + "' User: '" + cfg.TOUser + "' Pass len: '" + strconv.Itoa(len(cfg.TOPass)) + "'")
	log.Infoln("TempDir: '" + cfg.TempDir + "'")

	toFQDN := cfg.TOURL.Scheme + "://" + cfg.TOURL.Host
	log.Infoln("TO FQDN: '" + toFQDN + "'")
	log.Infoln("TO URL: '" + cfg.TOURL.String() + "'")

	toClient, err := GetClient(toFQDN, cfg.TOUser, cfg.TOPass, cfg.TempDir, cfg.CacheFileMaxAge, cfg.TOTimeout, cfg.TOInsecure)
	if err != nil {
		log.Errorln("Logging in to Traffic Ops: " + err.Error())
		os.Exit(ExitCodeErrGeneric)
	}

	tccfg := TCCfg{Cfg: cfg, TOClient: &toClient}

	cfgFile, code, err := GetConfigFile(tccfg)
	log.Infof("GetConfigFile returned %v %v\n", code, err)
	if err != nil {
		log.Errorln("Getting config file '" + cfg.TOURL.String() + "': " + err.Error())
		if code == 0 {
			code = ExitCodeErrGeneric
		}
		log.Infof("GetConfigFile exiting with code %v\n", code)
		os.Exit(code)
	}
	fmt.Println(cfgFile)
	os.Exit(ExitCodeSuccess)
}

func GetGeneratedFilesList() []string {
	names := []string{}
	for scope, fileFuncs := range ConfigFileFuncs() {
		for cfgFile, _ := range fileFuncs {
			names = append(names, scope+"/"+cfgFile)
		}
	}

	names = append(names, "profiles/url_sig_*.config")     // url_sig configs are generated, but not in the funcs because they're not a literal match
	names = append(names, "profiles/uri_signing_*.config") // uri_signing configs are generated, but not in the funcs because they're not a literal match
	names = append(names, "profiles/*")                    // unknown profiles configs are generated, a.k.a. "take-and-bake"

	return names
}

func HTTPCodeToExitCode(httpCode int) int {
	switch httpCode {
	case http.StatusNotFound:
		return ExitCodeNotFound
	}
	return ExitCodeErrGeneric
}
