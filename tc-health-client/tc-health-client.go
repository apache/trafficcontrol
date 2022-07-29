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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/tc-health-client/config"
	"github.com/apache/trafficcontrol/tc-health-client/tmagent"
)

const (
	Success      = 0
	ConfigError  = 166
	RunTimeError = 167
	PidFile      = "/run/tc-health-client.pid"
)

var (
	// BuildTimestamp is set via ld flags when the RPM is built. See build/build_rpm.sh.
	BuildTimestamp = ""
	// Version is set via ld flags when the RPM is built. See build/build_rpm.sh.
	Version = ""
)

func main() {
	cfg, err, helpflag := config.GetConfig()
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(ConfigError)
	}

	if helpflag { // user used --help option
		os.Exit(Success)
	}

	log.Infof("Polling interval: %v seconds\n", config.GetTMPollingInterval().Seconds())
	tmInfo, err := tmagent.NewParentInfo(cfg)
	if err != nil {
		log.Errorf("startup could not initialize parent info, check that trafficserver is running: %s\n", err.Error())
		os.Exit(RunTimeError)
	}

	pid := os.Getpid()
	err = os.WriteFile(PidFile, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Errorf("could not write the process id to %s: %s", PidFile, err.Error())
		os.Exit(RunTimeError)
	}

	log.Infof("startup complete, version: %s, built: %s\n", Version, BuildTimestamp)

	tmInfo.PollAndUpdateCacheStatus()
}
