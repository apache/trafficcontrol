// t3c-generate is a tool for generating configuration files server-side on ATC cache servers. See README.md for usage.

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
	"sort"
	"strings"

	"github.com/apache/trafficcontrol/cache-config/t3c-generate/cfgfile"
	"github.com/apache/trafficcontrol/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/cache-config/t3c-generate/getdata"
	"github.com/apache/trafficcontrol/cache-config/t3c-generate/plugin"
	"github.com/apache/trafficcontrol/cache-config/t3c-generate/toreq"
	"github.com/apache/trafficcontrol/lib/go-log"
)

func main() {
	cfg, err := config.GetCfg()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting config: "+err.Error()+"\n")
		os.Exit(config.ExitCodeErrGeneric)
	}

	if cfg.ListPlugins {
		fmt.Println(strings.Join(plugin.List(), "\n"))
		os.Exit(0)
	}

	// Because logs will be appended, we want a "start" message, so individual runs are easily distinguishable.
	// log the "start" message to each distinct logger.
	startMsg := "Starting t3c-generate"
	log.Infoln(startMsg)
	if cfg.WarningLog() != cfg.ErrorLog() {
		log.Warnln(startMsg)
	}
	if cfg.InfoLog() != cfg.WarningLog() && cfg.InfoLog() != cfg.ErrorLog() {
		log.Infoln(startMsg)
	}

	plugins := plugin.Get(cfg)
	plugins.OnStartup(plugin.StartupData{Cfg: cfg})

	toClient, err := toreq.New(cfg.TOURL, cfg.TOUser, cfg.TOPass, cfg.TOInsecure, cfg.TOTimeout, config.UserAgent)
	if err != nil {
		log.Errorln(err)
		os.Exit(config.ExitCodeErrGeneric)
	}

	if toClient.FellBack() {
		log.Warnln("Traffic Ops does not support the latest version supported by this app! Falling back to previous major Traffic Ops API version!")
	}

	tccfg := config.TCCfg{Cfg: cfg, TOClient: toClient}

	if tccfg.GetData != "" {
		if err := getdata.WriteData(tccfg); err != nil {
			log.Errorln("writing data: " + err.Error())
			os.Exit(config.ExitCodeErrGeneric)
		}
		os.Exit(config.ExitCodeSuccess)
	}

	if tccfg.SetRevalStatus != "" || tccfg.SetQueueStatus != "" {
		if err := getdata.SetQueueRevalStatuses(tccfg); err != nil {
			log.Errorln("writing queue and reval statuses: " + err.Error())
			os.Exit(config.ExitCodeErrGeneric)
		}
		os.Exit(config.ExitCodeSuccess)
	}

	toData, toIPs, err := cfgfile.GetTOData(tccfg)
	if err != nil {
		log.Errorln("getting data from traffic ops: " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	configs, err := cfgfile.GetAllConfigs(toData, config.UserAgent, toIPs, tccfg)
	if err != nil {
		log.Errorln("Getting config for'" + cfg.CacheHostName + "': " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	modifyFilesData := plugin.ModifyFilesData{Cfg: tccfg, TOData: toData, Files: configs}
	configs = plugins.ModifyFiles(modifyFilesData)

	sort.Sort(config.ATSConfigFiles(configs))

	if err := cfgfile.WriteConfigs(configs, os.Stdout); err != nil {
		log.Errorln("Writing configs for '" + cfg.CacheHostName + "': " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	os.Exit(config.ExitCodeSuccess)
}
