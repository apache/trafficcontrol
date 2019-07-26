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
	"os"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
)

const AppName = "atstccfg"
const Version = "0.1"
const UserAgent = AppName + "/" + Version

const APIVersion = "1.2"
const TempSubdir = AppName + "_cache"
const TempCookieFileName = "cookies"
const TOCookieName = "mojolicious"

// TODO make the below configurable?
const TOInsecure = false
const TOTimeout = time.Second * 10
const CacheFileMaxAge = time.Minute

func main() {
	cfg, err := GetCfg()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting config: "+err.Error()+"\n")
		os.Exit(1)
	}

	log.Infoln("URL: '" + cfg.TOURL.String() + "' User: '" + cfg.TOUser + "' Pass len: '" + strconv.Itoa(len(cfg.TOPass)) + "'")
	log.Infoln("TempDir: '" + cfg.TempDir + "'")

	toFQDN := cfg.TOURL.Scheme + "://" + cfg.TOURL.Host
	log.Infoln("TO FQDN: '" + toFQDN + "'")
	log.Infoln("TO URL: '" + cfg.TOURL.String() + "'")

	toClient, err := GetClient(toFQDN, cfg.TOUser, cfg.TOPass, cfg.TempDir)
	if err != nil {
		log.Errorln("Logging in to Traffic Ops: " + err.Error())
		os.Exit(1)
	}

	cfgFile, code, err := GetConfigFile(&toClient, cfg)
	if err != nil {
		log.Errorln("Getting config file '" + cfg.TOURL.String() + "' from Traffic Ops: " + err.Error())
		if code == 0 {
			code = 1
		}
		os.Exit(code)
	}
	fmt.Println(cfgFile)
	os.Exit(0)
}
