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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-request"

type Cfg struct {
	CommandArgs      []string
	LogLocationDebug string
	LogLocationWarn  string
	LogLocationError string
	LogLocationInfo  string
	LoginDispersion  time.Duration
	t3cutil.TCCfg
	Version     string
	GitRevision string
}

func (cfg Cfg) AppVersion() string { return t3cutil.VersionStr(AppName, cfg.Version, cfg.GitRevision) }
func (cfg Cfg) UserAgent() string  { return t3cutil.UserAgentStr(AppName, cfg.Version, cfg.GitRevision) }

func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationError) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // event logging is not used.

// Usage() writes command line options and usage to 'stderr'
func Usage() {
	getopt.PrintUsage(os.Stderr)
	os.Exit(0)
}

// InitConfig() intializes the configuration variables and loggers.
func InitConfig(appVersion string, gitRevision string) (Cfg, error) {
	dispersionPtr := getopt.IntLong("login-dispersion", 'l', 0, "[seconds] wait a random number of seconds between 0 and [seconds] before login to traffic ops, default 0")
	cacheHostNamePtr := getopt.StringLong("cache-host-name", 'H', "", "Host name of the cache to generate config for. Must be the server host name in Traffic Ops, not a URL, and not the FQDN")
	getDataPtr := getopt.StringLong("get-data", 'D', "system-info", "non-config-file Traffic Ops Data to get. Valid values are update-status, packages, chkconfig, system-info, and statuses")
	toInsecurePtr := getopt.BoolLong("traffic-ops-insecure", 'I', "[true | false] ignore certificate errors from Traffic Ops")
	toTimeoutMSPtr := getopt.IntLong("traffic-ops-timeout-milliseconds", 't', 30000, "Timeout in milli-seconds for Traffic Ops requests, default is 30000")
	toURLPtr := getopt.StringLong("traffic-ops-url", 'u', "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with     the environment variable TO_URL")
	toUserPtr := getopt.StringLong("traffic-ops-user", 'U', "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER")
	revalOnlyPtr := getopt.BoolLong("reval-only", 'r', "[true | false] whether to only fetch data needed to revalidate, versus all config data. Only used if get-data is config")
	disableProxyPtr := getopt.BoolLong("traffic-ops-disable-proxy", 'p', "[true | false] whether to not use any configure Traffic Ops proxy parameter. Only used if get-data is config")
	toPassPtr := getopt.StringLong("traffic-ops-password", 'P', "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS    ")
	oldCfgPtr := getopt.StringLong("old-config", 'c', "", "Old config from a previous config request. Optional. May be a file path, or 'stdin' to read from stdin. Used to make conditional requests.")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	versionPtr := getopt.BoolLong("version", 'V', "Print the app version")
	verbosePtr := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv'. To omit error logging, see '-s'`)
	silentPtr := getopt.BoolLong("silent", 's', `Silent. Errors are not logged, and the 'verbose' flag is ignored. If a fatal error occurs, the return code will be non-zero but no text will be output to stderr`)

	getopt.Parse()

	if *helpPtr == true {
		Usage()
	} else if *versionPtr == true {
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

	dispersion := time.Second * time.Duration(*dispersionPtr)
	toTimeoutMS := time.Millisecond * time.Duration(*toTimeoutMSPtr)
	toURL := *toURLPtr
	toUser := *toUserPtr
	toPass := *toPassPtr

	urlSourceStr := "argument" // for error messages
	if toURL == "" {
		urlSourceStr = "environment variable"
		toURL = os.Getenv("TO_URL")
	}
	if toUser == "" {
		toUser = os.Getenv("TO_USER")
	}
	if *toPassPtr == "" {
		toPass = os.Getenv("TO_PASS")
	}

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return Cfg{}, errors.New("parsing Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	} else if err := t3cutil.ValidateURL(toURLParsed); err != nil {
		return Cfg{}, errors.New("invalid Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	}

	var cacheHostName string
	if len(*cacheHostNamePtr) > 0 {
		cacheHostName = *cacheHostNamePtr
	} else {
		cacheHostName, err = os.Hostname()
		if err != nil {
			return Cfg{}, errors.New("could not get the OS hostname, please supply a hostname: " + err.Error())
		}
	}

	cfg := Cfg{
		CommandArgs:      getopt.Args(),
		LogLocationDebug: logLocationDebug,
		LogLocationError: logLocationError,
		LogLocationInfo:  logLocationInfo,
		LogLocationWarn:  logLocationWarn,
		LoginDispersion:  dispersion,
		TCCfg: t3cutil.TCCfg{
			CacheHostName:  cacheHostName,
			GetData:        *getDataPtr,
			TOInsecure:     *toInsecurePtr,
			TOTimeoutMS:    toTimeoutMS,
			TOUser:         toUser,
			TOPass:         toPass,
			TOURL:          toURLParsed,
			RevalOnly:      *revalOnlyPtr,
			TODisableProxy: *disableProxyPtr,
			T3CVersion:     gitRevision,
		},
		Version:     appVersion,
		GitRevision: gitRevision,
	}

	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("initializing loggers: " + err.Error())
	}

	// load old config after initializing the loggers, because we want to log how long it takes
	oldCfg, err := LoadOldCfg(*oldCfgPtr)
	if err != nil {
		log.Warnf("loading old config failed, old config will not be used! Error: %v\n", err)
	} else {
		log.Infof("using old config for IMS requests")
		cfg.OldCfg = oldCfg
	}

	return cfg, nil
}

func (cfg Cfg) PrintConfig() {
	log.Debugf("CommandArgs: %s\n", cfg.CommandArgs)
	log.Debugf("LogLocationDebug: %s\n", cfg.LogLocationDebug)
	log.Debugf("LogLocationError: %s\n", cfg.LogLocationError)
	log.Debugf("LogLocationInfo: %s\n", cfg.LogLocationInfo)
	log.Debugf("LogLocationWarn: %s\n", cfg.LogLocationWarn)
	log.Debugf("LoginDispersion : %s\n", cfg.LoginDispersion)
	log.Debugf("CacheHostName: %s\n", cfg.CacheHostName)
	log.Debugf("TOInsecure: %v\n", cfg.TOInsecure)
	log.Debugf("TOTimeoutMS: %s\n", cfg.TOTimeoutMS)
	log.Debugf("TOUser: %s\n", cfg.TOUser)
	log.Debugf("TOPass: xxxxxx\n")
	log.Debugf("TOURL: %s\n", cfg.TOURL)
}

func LoadOldCfg(path string) (*t3cutil.ConfigData, error) {
	defer func(start time.Time) { log.Infof("loading old config took %v\n", time.Since(start)) }(time.Now())
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil // old config is optional.
	}

	if strings.ToLower(path) == "stdin" {
		cfg := &t3cutil.ConfigData{}
		if err := json.NewDecoder(os.Stdin).Decode(cfg); err != nil {
			return nil, errors.New("decoding old config from stdin: " + err.Error())
		}
		return cfg, nil
	}

	fi, err := os.Open(path)
	if err != nil {
		return nil, errors.New("opening old config file '" + path + "': " + err.Error())
	}
	defer fi.Close()

	cfg := &t3cutil.ConfigData{}
	if err := json.NewDecoder(fi).Decode(cfg); err != nil {
		return nil, errors.New("decoding old config file '" + path + "': " + err.Error())
	}

	return cfg, nil
}
