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
	"os"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/tc-health-client/config"
	"github.com/apache/trafficcontrol/tc-health-client/tmagent"
)

const (
	Success      = 0
	ConfigError  = 166
	RunTimeError = 167
)

func main() {
	cfg, err, helpflag := config.GetConfig()
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(ConfigError)
	} else {
		log.Infoln("Startup complete, the config has been loaded")
	}
	if helpflag { // user used --help option
		os.Exit(Success)
	}

	log.Infof("Polling interval: %d\n", config.GetTMPollingInterval())
	tmInfo, err := tmagent.NewParentInfo(cfg)
	if err != nil {
		log.Errorf("startup could not initialize parent info, check that trafficserver is running: %s\n", err.Error())
		os.Exit(RunTimeError)
	}

	tmInfo.PollAndUpdateCacheStatus()
}
