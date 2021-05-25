package config

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
	"os"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/pborman/getopt/v2"
)

type Cfg struct {
	CommandArgs            []string
	LogLocationDebug       string
	LogLocationError       string
	LogLocationInfo        string
	TrafficServerConfigDir string
	TrafficServerPluginDir string
	FilesAdding            map[string]struct{}
}

var (
	defaultATSConfigDir = "/opt/trafficserver/etc/trafficserver"
	defaultATSPluginDir = "/opt/trafficserver/libexec/trafficserver"
)

func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationError) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(log.LogLocationNull) } // warn logging is not used.
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // event logging is not used.

// Usage() writes command line options and usage to 'stderr'
func Usage() {
	getopt.PrintUsage(os.Stderr)
	os.Exit(0)
}

// InitConfig() intializes the configuration variables and loggers.
func InitConfig() (Cfg, error) {

	logLocationDebugPtr := getopt.StringLong("log-location-debug", 'd', "", "Where to log debugs. May be a file path, stdout, stderr")
	logLocationErrorPtr := getopt.StringLong("log-location-error", 'e', "stderr", "Where to log errors. May be a file path, stdout, stderr")
	logLocationInfoPtr := getopt.StringLong("log-location-info", 'i', "stderr", "Where to log infos. May be a file path, stdout, stderr")
	atsConfigDirPtr := getopt.StringLong("trafficserver-config-dir", 'c', defaultATSConfigDir, "directory where ATS config files are stored.")
	atsPluginDirPtr := getopt.StringLong("trafficserver-plugin-dir", 'p', defaultATSPluginDir, "directory where ATS plugins are stored.")
	filesAdding := getopt.StringLong("files-adding", 'f', "", "comma-delimited list of file names being added, to not fail to verify if they don't already exist.")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	getopt.Parse()

	if *helpPtr == true {
		Usage()
	}

	filesAddingSet := map[string]struct{}{}
	for _, fileAdding := range strings.Split(*filesAdding, ",") {
		fileAdding := strings.TrimSpace(fileAdding)
		if fileAdding == "" {
			continue
		}
		filesAddingSet[fileAdding] = struct{}{}
	}

	cfg := Cfg{
		CommandArgs:            getopt.Args(),
		LogLocationDebug:       *logLocationDebugPtr,
		LogLocationError:       *logLocationErrorPtr,
		LogLocationInfo:        *logLocationInfoPtr,
		TrafficServerConfigDir: *atsConfigDirPtr,
		TrafficServerPluginDir: *atsPluginDirPtr,
		FilesAdding:            filesAddingSet,
	}

	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("initializing loggers: " + err.Error())
	}

	return cfg, nil
}
