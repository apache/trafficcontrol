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
	"net/url"
	"os"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-update"

type Cfg struct {
	Debug            bool
	CommandArgs      []string
	LogLocationDebug string
	LogLocationError string
	LogLocationInfo  string
	LogLocationWarn  string
	LoginDispersion  time.Duration
	CacheHostName    string
	GetData          string
	ConfigApplyTime  *time.Time
	RevalApplyTime   *time.Time
	ConfigApplyBool  *bool
	RevalApplyBool   *bool
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
	const setConfigApplyTimeFlagName = "set-config-apply-time"
	configApplyTimeStringPtr := getopt.StringLong(setConfigApplyTimeFlagName, 'q', "", "[RFC3339Nano Timestamp] sets the server's config apply time")
	const setRevalApplyTimeFlagName = "set-reval-apply-time"
	revalApplyTimeStringPtr := getopt.StringLong(setRevalApplyTimeFlagName, 'a', "", "[RFC3339Nano Timestamp] sets the server's reval apply time")
	toInsecurePtr := getopt.BoolLong("traffic-ops-insecure", 'I', "[true | false] ignore certificate errors from Traffic Ops")
	toTimeoutMSPtr := getopt.IntLong("traffic-ops-timeout-milliseconds", 't', 30000, "Timeout in milli-seconds for Traffic Ops requests, default is 30000")
	toURLPtr := getopt.StringLong("traffic-ops-url", 'u', "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with     the environment variable TO_URL")
	toUserPtr := getopt.StringLong("traffic-ops-user", 'U', "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER")
	toPassPtr := getopt.StringLong("traffic-ops-password", 'P', "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS    ")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	versionPtr := getopt.BoolLong("version", 'V', "Print the version")
	verbosePtr := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv'. To omit error logging, see '-s'`)
	silentPtr := getopt.BoolLong("silent", 's', `Silent. Errors are not logged, and the 'verbose' flag is ignored. If a fatal error occurs, the return code will be non-zero but no text will be output to stderr`)

	// *** Compatability requirement until ATC (v7.0+) is deployed with the timestamp features
	const setConfigApplyBoolFlagName = "set-update-status"
	configApplyBoolFlag := getopt.BoolLong(setConfigApplyBoolFlagName, 'y', `[false or nonexistent] Set the Update Status to false for the server`)
	const setRevalApplyBoolFlagName = "set-reval-status"
	revalApplyBoolFlag := getopt.BoolLong(setRevalApplyBoolFlagName, 'z', `[false or nonexistent] Set the Reval Status to false for the server`)
	// ***

	getopt.Parse()

	if *helpPtr == true {
		Usage()
	} else if *versionPtr {
		cfg := &Cfg{Version: appVersion, GitRevision: gitRevision}
		fmt.Println(cfg.AppVersion())
		os.Exit(0)
	}

	// Verify at least one flag is passed
	if (!getopt.IsSet(setConfigApplyTimeFlagName) && !getopt.IsSet(setRevalApplyTimeFlagName)) &&
		(!getopt.IsSet(setConfigApplyBoolFlagName) && !getopt.IsSet(setRevalApplyBoolFlagName)) { // TODO: Remove once ATC (v7.0+) is deployed
		fmt.Printf("Must set either %s or %s. One is at least required.\n", setConfigApplyTimeFlagName, setRevalApplyTimeFlagName)
		os.Exit(0)
	}

	var configApplyTimePtr, revalApplyTimePtr *time.Time
	// Validate that it can be parsed to a valid timestamp
	if getopt.IsSet(setConfigApplyTimeFlagName) {
		parsed, err := time.Parse(time.RFC3339Nano, *configApplyTimeStringPtr)
		if err != nil {
			fmt.Printf("%s must be a valid RFC3339Nano timestamp", setConfigApplyTimeFlagName)
		}
		configApplyTimePtr = &parsed
	}
	if getopt.IsSet(setRevalApplyTimeFlagName) {
		parsed, err := time.Parse(time.RFC3339Nano, *revalApplyTimeStringPtr)
		if err != nil {
			fmt.Printf("%s must be a valid RFC3339Nano timestamp", setRevalApplyTimeFlagName)
		}
		revalApplyTimePtr = &parsed
	}

	// TODO: Remove once ATC (v7.0+) is deployed
	var configApplyBoolPtr, revalApplyBoolPtr *bool
	if getopt.IsSet(setConfigApplyBoolFlagName) {
		if *configApplyBoolFlag {
			fmt.Println("set-update-status must be false or nonexistent")
			os.Exit(0)
		}
		configApplyBoolPtr = configApplyBoolFlag
	} else {
		configApplyBoolPtr = nil
	}
	if getopt.IsSet(setRevalApplyBoolFlagName) {
		if *revalApplyBoolFlag {
			fmt.Println("set-reval-status must be false or nonexistent")
			os.Exit(0)
		}
		revalApplyBoolPtr = revalApplyBoolFlag
	} else {
		revalApplyBoolPtr = nil
	}

	// Booleans must trump time for backwards compatibility
	if configApplyTimePtr != nil {
		configApplyBoolPtr = nil
	}
	if revalApplyTimePtr != nil {
		revalApplyBoolPtr = nil
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
		ConfigApplyTime:  configApplyTimePtr,
		RevalApplyTime:   revalApplyTimePtr,
		ConfigApplyBool:  configApplyBoolPtr,
		RevalApplyBool:   revalApplyBoolPtr,
		TCCfg: t3cutil.TCCfg{
			CacheHostName: cacheHostName,
			GetData:       "update-status",
			TOInsecure:    *toInsecurePtr,
			TOTimeoutMS:   toTimeoutMS,
			TOUser:        toUser,
			TOPass:        toPass,
			TOURL:         toURLParsed,
		},
		Version:     appVersion,
		GitRevision: gitRevision,
	}

	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("initializing loggers: " + err.Error())
	}

	return cfg, nil
}
