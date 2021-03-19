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
	"bytes"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/pborman/getopt/v2"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

var TSHome string = "/opt/trafficserver"
var TSConfigDir string = "/opt/trafficserver/etc/trafficserver"

const (
	StatusDir          = "/opt/ort/status"
	AtsTcConfig        = "/opt/ort/atstccfg"
	Chkconfig          = "/sbin/chkconfig"
	Service            = "/sbin/service"
	SystemCtl          = "/bin/systemctl"
	TmpBase            = "/tmp/ort"
	TrafficCtl         = "/bin/traffic_ctl"
	TrafficServerOwner = "ats"
)

type Mode int

const (
	BadAss     Mode = 0
	Report     Mode = 1
	Revalidate Mode = 2
	SyncDS     Mode = 3
)

func (m Mode) String() string {
	switch m {
	case 0:
		return "BadAss"
	case 1:
		return "Report"
	case 2:
		return "Revalidate"
	case 3:
		return "SyncDS"
	}
	return ""
}

type SvcManagement int

const (
	Unknown SvcManagement = 0
	SystemD SvcManagement = 1
	SystemV SvcManagement = 2 // legacy System V Init.
)

func (s SvcManagement) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case SystemD:
		return "SystemD"
	case SystemV:
		return "SystemV"
	}
	return "Unknown"
}

type Cfg struct {
	Dispersion          time.Duration
	LogLocationDebug    string
	LogLocationErr      string
	LogLocationInfo     string
	LogLocationWarn     string
	LoginDispersion     time.Duration
	CacheHostName       string
	SvcManagement       SvcManagement
	Retries             int
	RevalWaitTime       time.Duration
	ReverseProxyDisable bool
	RunMode             Mode
	SkipOSCheck         bool
	TOInsecure          bool
	TOTimeoutMS         time.Duration
	TOUser              string
	TOPass              string
	TOURL               string
	DNSLocalBind        bool
	WaitForParents      bool
	YumOptions          string
	// UseGit is whether to create and maintain a git repo of config changes.
	// Note this only applies to the ATS config directory inferred or set via the flag.
	//      It does not do anything for config files generated outside that location.
	UseGit UseGitFlag
}

type UseGitFlag string

const (
	UseGitAuto    = "auto"
	UseGitYes     = "yes"
	UseGitNo      = "no"
	UseGitInvalid = ""
)

func StrToUseGitFlag(str string) UseGitFlag {
	str = strings.ToLower(strings.TrimSpace(str))
	switch str {
	case UseGitAuto:
		fallthrough
	case UseGitYes:
		fallthrough
	case UseGitNo:
		return UseGitFlag(str)
	default:
		return UseGitInvalid
	}
}

func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationErr) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // event logging is not used.

func directoryExists(dir string) (bool, os.FileInfo) {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	return info.IsDir(), info
}

// derives the ATS Installation directory from
// the rpm config file list.
func GetTSPackageHome() string {
	var dir []string
	var output bytes.Buffer
	var tsHome string = ""
	var files []string

	cmd := exec.Command("/bin/rpm", "-q", "-c", "trafficserver")
	cmd.Stdout = &output
	err := cmd.Run()
	// on error or if the trafficserver rpm is not installed indicated
	// by a return code of '1', return an empty string.
	if err != nil || cmd.ProcessState.ExitCode() == 1 {
		return ""
	} else {
		responseStr := string(output.Bytes())
		files = strings.Split(responseStr, "\n")
	}

	if files != nil && err == nil { // trafficserver is installed, derive TSHome
		for ii := range files {
			line := strings.TrimSpace(files[ii])
			if strings.Contains(line, "etc/trafficserver") {
				dir = strings.Split(line, "/etc")
				break
			}
		}
		if dir != nil && len(dir) >= 0 {
			if ok, _ := directoryExists(dir[0]); ok == true {
				tsHome = dir[0]
			}
		}
	}
	return tsHome
}

func GetCfg() (Cfg, error) {
	var err error

	dispersionPtr := getopt.IntLong("dispersion", 'D', 300, "[seconds] wait a random number of seconds between 0 and [seconds] before starting, default 300")
	loginDispersionPtr := getopt.IntLong("login-dispersion", 'l', 0, "[seconds] wait a random number of seconds between 0 and [seconds] before login to traffic ops, default 0")
	logLocationDebugPtr := getopt.StringLong("log-location-debug", 'd', "", "Where to log debugs. May be a file path, stdout, stderr, or null, default ''")
	logLocationErrorPtr := getopt.StringLong("log-location-error", 'e', "stderr", "Where to log errors. May be a file path, stdout, stderr, or null, default stderr")
	logLocationInfoPtr := getopt.StringLong("log-location-info", 'i', "stderr", "Where to log info. May be a file path, stdout, stderr, or null, default stderr")
	logLocationWarnPtr := getopt.StringLong("log-location-warning", 'w', "stderr", "Where to log warnings. May be a file path, stdout, stderr, or null, default stderr")
	cacheHostNamePtr := getopt.StringLong("cache-host-name", 'H', "", "Host name of the cache to generate config for. Must be the server host name in Traffic Ops, not a URL, and not the FQDN")
	retriesPtr := getopt.IntLong("num-retries", 'r', 3, "[number] retry connection to Traffic Ops URL [number] times, default is 3")
	revalWaitTimePtr := getopt.IntLong("reval-wait-time", 'T', 60, "[seconds] wait a random number of seconds between 0 and [seconds] before revlidation, default is 60")
	reverseProxyDisablePtr := getopt.BoolLong("reverse-proxy-disable", 'p', "[false | true] bypass the reverse proxy even if one has been configured default is false")
	runModePtr := getopt.StringLong("run-mode", 'm', "report", "[badass | report | revalidate | syncds] run mode, default is 'report'")
	skipOSCheckPtr := getopt.BoolLong("skip-os-check", 's', "[false | true] skip os check, default is false")
	toInsecurePtr := getopt.BoolLong("traffic-ops-insecure", 'I', "[true | false] ignore certificate errors from Traffic Ops")
	toTimeoutMSPtr := getopt.IntLong("traffic-ops-timeout-milliseconds", 't', 30000, "Timeout in milli-seconds for Traffic Ops requests, default is 30000")
	toURLPtr := getopt.StringLong("traffic-ops-url", 'u', "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with the environment variable TO_URL")
	toUserPtr := getopt.StringLong("traffic-ops-user", 'U', "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER")
	toPassPtr := getopt.StringLong("traffic-ops-password", 'P', "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS")
	tsHomePtr := getopt.StringLong("trafficserver-home", 'R', "", "Trafficserver Package directory. May also be set with the environment variable TS_HOME")
	waitForParentsPtr := getopt.BoolLong("wait-for-parents", 'W', "[true | false] do not update if parent_pending = 1 in the update json. default is false, wait for parents")
	dnsLocalBindPtr := getopt.BoolLong("dns-local-bind", 'b', "[true | false] whether to use the server's Service Addresses to set the ATS DNS local bind address")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	useGitStr := getopt.StringLong("git", 'g', "auto", "Create and use a git repo in the config directory. Options are yes, no, and auto. If yes, create and use. If auto, use if it exist. Default is auto.")
	getopt.Parse()

	dispersion := time.Second * time.Duration(*dispersionPtr)
	loginDispersion := time.Second * time.Duration(*loginDispersionPtr)
	logLocationDebug := *logLocationDebugPtr
	logLocationError := *logLocationErrorPtr
	logLocationInfo := *logLocationInfoPtr
	logLocationWarn := *logLocationWarnPtr

	var cacheHostName string
	if len(*cacheHostNamePtr) > 0 {
		cacheHostName = *cacheHostNamePtr
	} else {
		cacheHostName, err = os.Hostname()
		if err != nil {
			return Cfg{}, errors.New("Could not get the hostname from the O.S., please supply a hostname: " + err.Error())
		}
	}

	useGit := StrToUseGitFlag(*useGitStr)

	if useGit == UseGitInvalid {
		return Cfg{}, errors.New("Invalid git flag '" + *useGitStr + "'. Valid options are yes, no, auto.")
	}

	retries := *retriesPtr
	revalWaitTime := time.Second * time.Duration(*revalWaitTimePtr)
	reverseProxyDisable := *reverseProxyDisablePtr
	skipOsCheck := *skipOSCheckPtr
	toInsecure := *toInsecurePtr
	toTimeoutMS := time.Millisecond * time.Duration(*toTimeoutMSPtr)
	toURL := *toURLPtr
	toUser := *toUserPtr
	toPass := *toPassPtr
	waitForParents := *waitForParentsPtr
	dnsLocalBind := *dnsLocalBindPtr
	help := *helpPtr

	if help {
		Usage()
		return Cfg{}, nil
	}

	runModeStr := strings.ToUpper(*runModePtr)
	runMode := Mode(Report)
	switch runModeStr {
	case "REPORT":
		runMode = Report
	case "BADASS":
		runMode = BadAss
	case "SYNCDS":
		runMode = SyncDS
	case "REVALIDATE":
		runMode = Revalidate
	default:
		Usage()
		return Cfg{}, errors.New(runModeStr + " is an invalid mode.")
	}

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

	// set TSHome
	var tsHome = ""
	if *tsHomePtr != "" {
		tsHome = *tsHomePtr
		fmt.Printf("set TSHome from command line: '%s'\n\n", TSHome)
	}
	if *tsHomePtr == "" { // evironment or rpm check.
		tsHome = os.Getenv("TS_HOME") // check for the environment variable.
		if tsHome != "" {
			fmt.Printf("set TSHome from TS_HOME environment variable '%s'\n", TSHome)
		} else { // finally check using the config file listing from the rpm package.
			tsHome = GetTSPackageHome()
			if tsHome != "" {
				fmt.Printf("set TSHome from the RPM config file  list '%s'\n", tsHome)
			} else {
				fmt.Printf("no override for TSHome was found, using the configured default: '%s'\n", TSHome)
			}
		}
	}
	if tsHome != "" {
		TSHome = tsHome
		TSConfigDir = tsHome + "/etc/trafficserver"
		fmt.Printf("TSHome: %s, TSConfigDir: %s\n", TSHome, TSConfigDir)
	}

	usageStr := "basic usage: t3c  --traffic-ops-url=myurl --traffic-ops-user=myuser --traffic-ops-password=mypass --cache-host-name=my-cache"
	if strings.TrimSpace(toURL) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-url or TO_URL environment variable. " + usageStr)
	}
	if strings.TrimSpace(toUser) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-user or TO_USER environment variable. " + usageStr)
	}
	if strings.TrimSpace(toPass) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-password or TO_PASS environment variable. " + usageStr)
	}
	if strings.TrimSpace(cacheHostName) == "" {
		return Cfg{}, errors.New("Missing required argument --cache-host-name. " + usageStr)
	}

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return Cfg{}, errors.New("parsing Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	} else if err = validateURL(toURLParsed); err != nil {
		return Cfg{}, errors.New("invalid Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	}

	svcManagement := getOSSvcManagement()
	yumOptions := os.Getenv("YUM_OPTIONS")

	cfg := Cfg{
		Dispersion:          dispersion,
		LogLocationDebug:    logLocationDebug,
		LogLocationErr:      logLocationError,
		LogLocationInfo:     logLocationInfo,
		LogLocationWarn:     logLocationWarn,
		LoginDispersion:     loginDispersion,
		CacheHostName:       cacheHostName,
		SvcManagement:       svcManagement,
		Retries:             retries,
		RevalWaitTime:       revalWaitTime,
		ReverseProxyDisable: reverseProxyDisable,
		RunMode:             runMode,
		SkipOSCheck:         skipOsCheck,
		TOInsecure:          toInsecure,
		TOTimeoutMS:         toTimeoutMS,
		TOUser:              toUser,
		TOPass:              toPass,
		TOURL:               toURL,
		DNSLocalBind:        dnsLocalBind,
		WaitForParents:      waitForParents,
		YumOptions:          yumOptions,
		UseGit:              useGit,
	}

	if err = log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}

	printConfig(cfg)

	return cfg, nil
}

func validateURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil url")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme expected 'http' or 'https', actual '" + u.Scheme + "'")
	}
	if strings.TrimSpace(u.Host) == "" {
		return errors.New("no host")
	}
	return nil
}

func isCommandAvailable(name string) bool {
	command := name + " --version"
	cmd := exec.Command("/bin/sh", "-c", command)
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func getOSSvcManagement() SvcManagement {
	var _svcManager SvcManagement

	if isCommandAvailable(SystemCtl) {
		_svcManager = SystemD
	} else if isCommandAvailable(Service) {
		_svcManager = SystemV
	}
	if !isCommandAvailable(Chkconfig) {
		return Unknown
	}

	// we have what we need
	return _svcManager
}

func printConfig(cfg Cfg) {
	log.Debugf("Dispersion: %d\n", cfg.Dispersion)
	log.Debugf("LogLocationDebug: %s\n", cfg.LogLocationDebug)
	log.Debugf("LogLocationErr: %s\n", cfg.LogLocationErr)
	log.Debugf("LogLocationInfo: %s\n", cfg.LogLocationInfo)
	log.Debugf("LogLocationWarn: %s\n", cfg.LogLocationWarn)
	log.Debugf("LoginDispersion: %d\n", cfg.LoginDispersion)
	log.Debugf("CacheHostName: %s\n", cfg.CacheHostName)
	log.Debugf("SvcManagement: %s\n", cfg.SvcManagement)
	log.Debugf("Retries: %d\n", cfg.Retries)
	log.Debugf("RevalWaitTime: %d\n", cfg.RevalWaitTime)
	log.Debugf("ReverseProxyDisable: %t\n", cfg.ReverseProxyDisable)
	log.Debugf("RunMode: %s\n", cfg.RunMode)
	log.Debugf("SkipOSCheck: %t\n", cfg.SkipOSCheck)
	log.Debugf("TOInsecure: %t\n", cfg.TOInsecure)
	log.Debugf("TOTimeoutMS: %d\n", cfg.TOTimeoutMS)
	log.Debugf("TOUser: %s\n", cfg.TOUser)
	log.Debugf("TOPass: Pass len: '%d'\n", len(cfg.TOPass))
	log.Debugf("TOURL: %s\n", cfg.TOURL)
	log.Debugf("TSHome: %s\n", TSHome)
	log.Debugf("WaitForParents: %t\n", cfg.WaitForParents)
	log.Debugf("YumOptions: %s\n", cfg.YumOptions)
}

func Usage() {
	fmt.Println("Usage: t3c [options]")
	fmt.Println("\t[options]:")
	fmt.Println("\t  --dispersion=[time in seconds] | -D, [time in seconds] wait a random number between 0 and <time in seconds> before starting, default = 300s")
	fmt.Println("\t  --login-dispersion=[time in seconds] | -l, [time in seconds] wait a random number between 0 and <time in seconds> befor login, default = 0")
	fmt.Println("\t  --log-location-debug=[value] | -d [value], Where to log debugs. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --log-location-error=[value] | -e [value], Where to log errors. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --log-location-info=[value] | -i [value], Where to log info. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --log-location-warning=[value] | -w [value], Where to log warnings. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --run-mode=[mode] | -m [mode] where mode is one of [ report | badass | syncds | revalidate ], default = report")
	fmt.Println("\t  --cache-hostname=[hostname] | -H [hostname], Host name of the cache to generate config for. Must be the server host name in Traffic Ops, not a URL, and not the FQDN")
	fmt.Println("\t  --num-retries=[number] | -r [number], retry connection to Traffic Ops URL [number] times, default is 3")
	fmt.Println("\t  --reval-wait-time=[seconds] | -T [seconds] wait a random number of seconds between 0 and [seconds] before revlidation, default is 60")
	fmt.Println("\t  --rev-proxy-disable=[true|false] | -p [true|false] bypass the reverse proxy even if one has been configured, default = false")
	fmt.Println("\t  --skip-os-check=[true|false] | -s [true | false] bypass the check for a supported CentOS version. default = false")
	fmt.Println("\t  --traffic-ops-insecure=[true|false] -I [true | false] Whether to ignore HTTPS certificate errors from Traffic Ops. It is HIGHLY RECOMMENDED to never use this in a production environment, but only for debugging, default = false")
	fmt.Println("\t  --traffic-ops-timeout-milliseconds=[milliseconds] | -t [milliseconds] the Traffic Ops request timeout in milliseconds. Default = 30000 (30 seconds)")
	fmt.Println("\t  --traffic-ops-url=[url] | -u [url], Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with the environment variable TO_URL")
	fmt.Println("\t  --traffic-ops-user=[username] | -U [username], Traffic Ops username. Required. May also be set with the environment variable TO_USER")
	fmt.Println("\t  --traffic-ops-password=[password] | -P [password], Traffic Ops password. Required. May also be set with the environment variable TO_PASS")
	fmt.Println("\t  --trafficserver-home=[value] | -R [value], Trafficserver Package directory. May also be set with the environment variable TS_HOME")
	fmt.Println("\t  --wait-for-parents | -W [true | false] do not update if parent_pending = 1 in the update json. default = true, wait for parents\n")
	fmt.Println("\t  --help | -h, Print usage information and exit")
}
