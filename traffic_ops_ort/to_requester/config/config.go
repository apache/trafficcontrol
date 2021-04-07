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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops_ort/t3clib"
	"github.com/pborman/getopt/v2"
)

const AppName = "to_requester"
const Version = "0.1"
const UserAgent = AppName + "/" + Version

type Cfg struct {
	CommandArgs      []string
	LogLocationDebug string
	LogLocationError string
	LogLocationInfo  string
	LoginDispersion  time.Duration
	t3clib.TCCfg
}

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
	dispersionPtr := getopt.IntLong("login-dispersion", 'l', 0, "[seconds] wait a random number of seconds between 0 and [seconds] before login to traffic ops, default 0")
	cacheHostNamePtr := getopt.StringLong("cache-host-name", 'H', "", "Host name of the cache to generate config for. Must be the server host name in Traffic Ops, not a URL, and not the FQDN")
	getDataPtr := getopt.StringLong("get-data", 'D', "system-info", "non-config-file Traffic Ops Data to get. Valid values are update-status, packages, chkconfig, system-info, and statuses")
	toInsecurePtr := getopt.BoolLong("traffic-ops-insecure", 'I', "[true | false] ignore certificate errors from Traffic Ops")
	toTimeoutMSPtr := getopt.IntLong("traffic-ops-timeout-milliseconds", 't', 30000, "Timeout in milli-seconds for Traffic Ops requests, default is 30000")
	toURLPtr := getopt.StringLong("traffic-ops-url", 'u', "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with     the environment variable TO_URL")
	toUserPtr := getopt.StringLong("traffic-ops-user", 'U', "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER")
	toPassPtr := getopt.StringLong("traffic-ops-password", 'P', "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS    ")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	versionPtr := getopt.BoolLong("version", 'v', "Print the to_requester version")

	getopt.Parse()

	if *helpPtr == true {
		Usage()
	}
	if *versionPtr == true {
		fmt.Println(AppName + " v" + Version)
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
	} else if err := t3clib.ValidateURL(toURLParsed); err != nil {
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
		LogLocationDebug: *logLocationDebugPtr,
		LogLocationError: *logLocationErrorPtr,
		LogLocationInfo:  *logLocationInfoPtr,
		LoginDispersion:  dispersion,
		TCCfg: t3clib.TCCfg{
			CacheHostName: cacheHostName,
			GetData:       *getDataPtr,
			TOInsecure:    *toInsecurePtr,
			TOTimeoutMS:   toTimeoutMS,
			TOUser:        toUser,
			TOPass:        toPass,
			TOURL:         toURLParsed,
			UserAgent:     UserAgent,
		},
	}

	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("initializing loggers: " + err.Error())
	}

	return cfg, nil
}

func (cfg Cfg) PrintConfig() {
	fmt.Printf("CommandArgs: %s\n", cfg.CommandArgs)
	fmt.Printf("LogLocationDebug: %s\n", cfg.LogLocationDebug)
	fmt.Printf("LogLocationError: %s\n", cfg.LogLocationError)
	fmt.Printf("LogLocationInfo: %s\n", cfg.LogLocationInfo)
	fmt.Printf("LoginDispersion : %s\n", cfg.LoginDispersion)
	fmt.Printf("CacheHostName: %s\n", cfg.CacheHostName)
	fmt.Printf("TOInsecure: %s\n", cfg.TOInsecure)
	fmt.Printf("TOTimeoutMS: %s\n", cfg.TOTimeoutMS)
	fmt.Printf("TOUser: %s\n", cfg.TOUser)
	fmt.Printf("TOPass: xxxxxx\n")
	fmt.Printf("TOURL: %s\n", cfg.TOURL)
}
