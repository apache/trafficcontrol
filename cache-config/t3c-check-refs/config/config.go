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
	"fmt"
	"os"
	"strings"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-check-refs"

// ArgFilesAddingInput is the value of the --files-adding flag to indicate to
// the input (file or stdin, depending whether a filename argument was passed)
// is a t3cutil.CheckRefsInputFileAndAdding JSON object.
const ArgFilesAddingInput = "input"

type Cfg struct {
	CommandArgs            []string
	LogLocationDebug       string
	LogLocationWarn        string
	LogLocationError       string
	LogLocationInfo        string
	TrafficServerConfigDir string
	TrafficServerPluginDir string
	FilesAdding            map[string]struct{}
	FilesAddingInput       bool
	Version                string
	GitRevision            string
}

var (
	defaultATSConfigDir = "/opt/trafficserver/etc/trafficserver"
	defaultATSPluginDir = "/opt/trafficserver/libexec/trafficserver"
)

func (cfg Cfg) AppVersion() string { return t3cutil.VersionStr(AppName, cfg.Version, cfg.GitRevision) }

func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationError) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) } // warn logging is not used.
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // event logging is not used.

// Usage() writes command line options and usage to 'stderr'
func Usage() {
	getopt.PrintUsage(os.Stderr)
	os.Exit(0)
}

// InitConfig() intializes the configuration variables and loggers.
func InitConfig(appVersion string, gitRevision string) (Cfg, error) {
	versionPtr := getopt.BoolLong("version", 'V', "Print version information and exit.")
	atsConfigDirPtr := getopt.StringLong("trafficserver-config-dir", 'c', defaultATSConfigDir, "directory where ATS config files are stored.")
	atsPluginDirPtr := getopt.StringLong("trafficserver-plugin-dir", 'p', defaultATSPluginDir, "directory where ATS plugins are stored.")
	filesAdding := getopt.StringLong("files-adding", 'f', "", "comma-delimited list of file names being added, to not fail to verify if they don't already exist.")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	verbosePtr := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv'. To omit error logging, see '-s'`)
	silentPtr := getopt.BoolLong("silent", 's', `Silent. Errors are not logged, and the 'verbose' flag is ignored. If a fatal error occurs, the return code will be non-zero but no text will be output to stderr`)

	getopt.Parse()

	if *helpPtr == true {
		Usage()
		os.Exit(0)
	} else if *versionPtr {
		cfg := &Cfg{Version: appVersion, GitRevision: gitRevision}
		fmt.Println(cfg.AppVersion())
		os.Exit(0)
	}

	logLocationError := log.LogLocationStderr
	logLocationWarn := log.LogLocationNull
	logLocationInfo := log.LogLocationNull
	logLocationDebug := log.LogLocationNull
	if *silentPtr {
		logLocationError = log.LogLocationNull
	} else {
		if *verbosePtr >= 1 {
			logLocationWarn = log.LogLocationStderr
		}
		if *verbosePtr >= 2 {
			logLocationInfo = log.LogLocationStderr
			logLocationDebug = log.LogLocationStderr // t3c only has 3 verbosity options: none (-s), error (default or --verbose=0), warning (-v), and info (-vv). Any code calling log.Debug is treated as Info.
		}
	}

	if *verbosePtr > 2 {
		return Cfg{}, errors.New("Too many verbose options. The maximum log verbosity level is 2 (-vv or --verbose=2) for errors (0), warnings (1), and info (2)")
	}

	filesAddingInput := false
	filesAddingSet := map[string]struct{}{}
	if strings.ToLower(strings.TrimSpace(*filesAdding)) == ArgFilesAddingInput {
		filesAddingInput = true
	} else {
		filesAddingSet = ArgToFilesAdding(*filesAdding)
	}

	cfg := Cfg{
		CommandArgs:            getopt.Args(),
		LogLocationDebug:       logLocationDebug,
		LogLocationError:       logLocationError,
		LogLocationInfo:        logLocationInfo,
		LogLocationWarn:        logLocationWarn,
		TrafficServerConfigDir: *atsConfigDirPtr,
		TrafficServerPluginDir: *atsPluginDirPtr,
		FilesAdding:            filesAddingSet,
		Version:                appVersion,
		GitRevision:            gitRevision,
		FilesAddingInput:       filesAddingInput,
	}

	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("initializing loggers: " + err.Error())
	}

	return cfg, nil
}

// ArgToFilesAdding converts a comma-delimited list of files being added to a set.
func ArgToFilesAdding(filesAddingVal string) map[string]struct{} {
	filesAddingSet := map[string]struct{}{}
	for _, fileAdding := range strings.Split(filesAddingVal, ",") {
		fileAdding := strings.TrimSpace(fileAdding)
		if fileAdding == "" {
			continue
		}
		filesAddingSet[fileAdding] = struct{}{}
	}
	return filesAddingSet
}
