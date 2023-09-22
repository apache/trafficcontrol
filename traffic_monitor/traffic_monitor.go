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
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/manager"
)

// GitRevision is the git revision of the app. The app SHOULD always be built with this set via the `-X` flag.
var GitRevision = "No Git Revision Specified. Please build with '-X main.GitRevision=${git rev-parse HEAD}'"

// BuildTimestamp is the time the app was built. The app SHOULD always be built with this set via the `-X` flag.
var BuildTimestamp = "No Build Timestamp Specified. Please build with '-X main.BuildTimestamp=`date +'%Y-%M-%dT%H:%M:%S'`"

func InitAccessCfg(cfg config.Config) error {
	accessW, err := config.GetAccessLogWriter(cfg)
	if err != nil {
		return err
	}
	log.InitAccess(accessW)
	return nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	staticData, err := config.GetStaticAppData(Version, GitRevision, BuildTimestamp)
	if err != nil {
		fmt.Printf("Error starting service: failed to get static app data: %v\n", err)
		os.Exit(1)
	}

	opsConfigFile := flag.String("opsCfg", "", "The traffic ops config file")
	configFileName := flag.String("config", "", "The Traffic Monitor config file path")
	flag.Parse()

	if *opsConfigFile == "" {
		fmt.Println("Error starting service: The --opsCfg argument is required")
		os.Exit(1)
	}

	// TODO add hot reloading (like opsConfigFile)?
	cfg, err := config.Load(*configFileName)
	if err != nil {
		fmt.Printf("Error starting service: failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := log.InitCfg(cfg); err != nil {
		fmt.Printf("Error starting service: failed to create log writers: %v\n", err)
		os.Exit(1)
	}

	if err := InitAccessCfg(cfg); err != nil {
		fmt.Printf("Error starting service: failed to create access log writer: %v\n", err)
		os.Exit(1)
	}

	if cfg.ShortHostnameOverride != "" {
		staticData.Hostname = cfg.ShortHostnameOverride
	}

	rand.Seed(time.Now().UnixNano())
	log.Infof("Starting with config %+v\n", cfg)

	err = manager.Start(*opsConfigFile, cfg, staticData, *configFileName)
	if err != nil {
		fmt.Printf("Error starting service: failed to start managers: %v\n", err)
		os.Exit(1)
	}
}
