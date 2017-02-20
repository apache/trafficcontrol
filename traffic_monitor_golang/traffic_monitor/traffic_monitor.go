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
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/instrumentation"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/config"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/manager"
	_ "github.com/davecheney/gmx"
)

// GitRevision is the git revision of the app. The app SHOULD always be built with this set via the `-X` flag.
var GitRevision = "No Git Revision Specified. Please build with '-X main.GitRevision=${git rev-parse HEAD}'"

// BuildTimestamp is the time the app was built. The app SHOULD always be built with this set via the `-X` flag.
var BuildTimestamp = "No Build Timestamp Specified. Please build with '-X main.BuildTimestamp=`date +'%Y-%M-%dT%H:%M:%S'`"

// getHostNameWithoutDomain returns the machine hostname, without domain information.
// Modified from http://stackoverflow.com/a/34331660/292623
func getHostNameWithoutDomain() (string, error) {
	cmd := exec.Command("/bin/hostname", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	hostname := out.String()
	if len(hostname) < 1 {
		return "", fmt.Errorf("OS returned empty hostname")
	}
	hostname = hostname[:len(hostname)-1] // removing EOL
	return hostname, nil
}

// getStaticAppData returns app data available at start time.
// This should be called immediately, as it includes calculating when the app was started.
func getStaticAppData() (manager.StaticAppData, error) {
	var d manager.StaticAppData
	var err error
	d.StartTime = time.Now()
	d.GitRevision = GitRevision
	d.FreeMemoryMB = math.MaxUint64 // TODO remove if/when nothing needs this
	d.Version = Version
	if d.WorkingDir, err = os.Getwd(); err != nil {
		return manager.StaticAppData{}, err
	}
	d.Name = os.Args[0]
	d.BuildTimestamp = BuildTimestamp
	if d.Hostname, err = getHostNameWithoutDomain(); err != nil {
		return manager.StaticAppData{}, err
	}

	d.UserAgent = fmt.Sprintf("%s/%s", filepath.Base(d.Name), d.Version)
	return d, nil
}

func getLogWriter(location string) (io.Writer, error) {
	switch location {
	case config.LogLocationStdout:
		return os.Stdout, nil
	case config.LogLocationStderr:
		return os.Stderr, nil
	case config.LogLocationNull:
		return ioutil.Discard, nil
	default:
		return os.OpenFile(location, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
}
func getLogWriters(eventLoc, errLoc, warnLoc, infoLoc, debugLoc string) (io.Writer, io.Writer, io.Writer, io.Writer, io.Writer, error) {
	eventW, err := getLogWriter(eventLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log event writer %v: %v", eventLoc, err)
	}
	errW, err := getLogWriter(errLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log error writer %v: %v", errLoc, err)
	}
	warnW, err := getLogWriter(warnLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log warning writer %v: %v", warnLoc, err)
	}
	infoW, err := getLogWriter(infoLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log info writer %v: %v", infoLoc, err)
	}
	debugW, err := getLogWriter(debugLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log debug writer %v: %v", debugLoc, err)
	}
	return eventW, errW, warnW, infoW, debugW, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	staticData, err := getStaticAppData()
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

	eventW, errW, warnW, infoW, debugW, err := getLogWriters(cfg.LogLocationEvent, cfg.LogLocationError, cfg.LogLocationWarning, cfg.LogLocationInfo, cfg.LogLocationDebug)
	if err != nil {
		fmt.Printf("Error starting service: failed to create log writers: %v\n", err)
		os.Exit(1)
	}
	log.Init(eventW, errW, warnW, infoW, debugW)

	log.Infof("Starting with config %+v\n", cfg)

	manager.Start(*opsConfigFile, cfg, staticData)
}
