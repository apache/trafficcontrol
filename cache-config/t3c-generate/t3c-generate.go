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
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/cfgfile"
	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/plugin"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

func main() {
	cfg, err := config.GetCfg(Version, GitRevision)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Getting config: "+err.Error()+"\n")
		os.Exit(config.ExitCodeErrGeneric)
	}

	if cfg.ListPlugins {
		log.Errorln(strings.Join(plugin.List(), "\n"))
		os.Exit(0)
	}

	// Because logs will be appended, we want a "start" message, so individual runs are easily distinguishable.
	// log the "start" message to each distinct logger.
	startMsg := "Starting t3c-generate"
	log.Errorln(startMsg)
	if cfg.WarningLog() != cfg.ErrorLog() {
		log.Warnln(startMsg)
	}
	if cfg.InfoLog() != cfg.WarningLog() && cfg.InfoLog() != cfg.ErrorLog() {
		log.Infoln(startMsg)
	}

	plugins := plugin.Get(cfg)
	plugins.OnStartup(plugin.StartupData{Cfg: cfg})

	log.Infoln("reading Traffic Ops data from stdin")

	toData := &t3cutil.ConfigData{}
	if err := json.NewDecoder(os.Stdin).Decode(toData); err != nil {
		log.Errorln("reading and parsing input Traffic Ops data: " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	if toData.Server == nil {
		log.Errorln("input had no server")
		os.Exit(config.ExitCodeErrGeneric)
	} else if toData.Server.HostName == "" {
		log.Errorln("input server had no host name")
		os.Exit(config.ExitCodeErrGeneric)
	}

	if cfg.Cache == "varnish" {
		configs, err := cfgfile.GetVarnishConfigs(toData, cfg)
		if err != nil {
			log.Errorln("Generating varnish config for'" + toData.Server.HostName + "': " + err.Error())
			os.Exit(config.ExitCodeErrGeneric)
		}
		err = cfgfile.WriteConfigs(configs, os.Stdout)
		if err != nil {
			log.Errorln("Writing configs for '" + toData.Server.HostName + "': " + err.Error())
			os.Exit(config.ExitCodeErrGeneric)
		}
		os.Exit(config.ExitCodeSuccess)
	}

	configs, err := cfgfile.GetAllConfigs(toData, cfg)
	if err != nil {
		log.Errorln("Getting config for'" + toData.Server.HostName + "': " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	modifyFilesData := plugin.ModifyFilesData{Cfg: cfg, TOData: toData, Files: configs}
	configs = plugins.ModifyFiles(modifyFilesData)

	sort.Sort(t3cutil.ATSConfigFiles(configs))

	if err := cfgfile.WriteConfigs(configs, os.Stdout); err != nil {
		log.Errorln("Writing configs for '" + toData.Server.HostName + "': " + err.Error())
		os.Exit(config.ExitCodeErrGeneric)
	}

	os.Exit(config.ExitCodeSuccess)
}
