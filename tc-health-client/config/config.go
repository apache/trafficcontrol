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
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/tc-health-client/util"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"

	"github.com/pborman/getopt/v2"
)

var userAgent = "tc-health-client/1.0"
var tmPollingInterval time.Duration
var toRequestTimeout time.Duration
var toSession *toclient.Session = nil

const (
	DefaultPollStateJSONLog         = "/var/log/trafficcontrol/poll-state.json"
	DefaultConfigFile               = "/etc/trafficcontrol/tc-health-client.json"
	DefaultLogDirectory             = "/var/log/trafficcontrol"
	DefaultLogFile                  = "tc-health-client.log"
	DefaultTOLoginDispersionFactor  = 90
	DefaultTrafficServerConfigDir   = "/opt/trafficserver/etc/trafficserver"
	DefaultTrafficServerBinDir      = "/opt/trafficserver/bin"
	DefaultUnavailablePollThreshold = 2
	DefaultMarkupPollThreshold      = 1
)

type Cfg struct {
	CDNName                  string          `json:"cdn-name"`
	EnableActiveMarkdowns    bool            `json:"enable-active-markdowns"`
	ReasonCode               string          `json:"reason-code"`
	TOCredentialFile         string          `json:"to-credential-file"`
	TORequestTimeOutSeconds  string          `json:"to-request-timeout-seconds"`
	TOPass                   string          `json:"to-pass"`
	TOUrl                    string          `json:"to-url"`
	TOUser                   string          `json:"to-user"`
	TmProxyURL               string          `json:"tm-proxy-url"`
	TmPollIntervalSeconds    string          `json:"tm-poll-interval-seconds"`
	TOLoginDispersionFactor  int             `json:"to-login-dispersion-factor"`
	UnavailablePollThreshold int             `json:"unavailable-poll-threshold"`
	MarkUpPollThreshold      int             `json:"markup-poll-threshold"`
	TrafficServerConfigDir   string          `json:"trafficserver-config-dir"`
	TrafficServerBinDir      string          `json:"trafficserver-bin-dir"`
	PollStateJSONLog         string          `json:"poll-state-json-log"`
	EnablePollStateLog       bool            `json:"enable-poll-state-log"`
	TrafficMonitors          map[string]bool `json:"trafficmonitors,omitempty"`
	HealthClientConfigFile   util.ConfigFile
	CredentialFile           util.ConfigFile
	ParsedProxyURL           *url.URL
}

type LogCfg struct {
	LogLocationErr   string
	LogLocationDebug string
	LogLocationInfo  string
	LogLocationWarn  string
}

func (lcfg LogCfg) ErrorLog() log.LogLocation   { return log.LogLocation(lcfg.LogLocationErr) }
func (lcfg LogCfg) WarningLog() log.LogLocation { return log.LogLocation(lcfg.LogLocationWarn) }
func (lcfg LogCfg) InfoLog() log.LogLocation    { return log.LogLocation(lcfg.LogLocationInfo) }
func (lcfg LogCfg) DebugLog() log.LogLocation   { return log.LogLocation(lcfg.LogLocationDebug) }
func (lcfg LogCfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // not used

/**
 * ReadCredentials
 *
 * cfg - the existing config
 * updating - when true, existing credentials may be updated from the credential file
 */
func ReadCredentials(cfg *Cfg, updating bool) error {
	if cfg.TOCredentialFile == "" {
		return nil
	}

	fn := cfg.CredentialFile

	// verify that we have credentials or can read them from the credential file
	if fn.Filename == "" {
		if cfg.TOPass == "" || cfg.TOUser == "" || cfg.TOUrl == "" {
			return fmt.Errorf("cannot continue, no TO credentials or TO URL have been specified, check configs")
		} else {
			return nil
		}
	}

	// You should not configure a credential file and the credentials simultaneously in the health client
	// config file.  Either use an external credential file or put the credentials in the health client
	// config.  Precedence is given to credentials in the health client config file.
	if !updating && (cfg.TOPass != "" && cfg.TOUser != "" && cfg.TOUrl != "") {
		log.Warnf("credentials are defined in the %s file, will not override them with those in the %s file", cfg.HealthClientConfigFile.Filename, cfg.CredentialFile.Filename)
		cfg.CredentialFile.LastModifyTime = math.MaxInt64
		return nil
	}

	f, err := os.Open(fn.Filename)
	if err != nil {
		return errors.New("failed to open + " + fn.Filename + " :" + err.Error())
	}
	defer f.Close()

	var to_pass_found = false
	var to_url_found = false
	var to_user_found = false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, " ")
		for _, v := range fields {
			if strings.HasPrefix(v, "TO_") {
				sf := strings.Split(v, "=")
				if len(sf) == 2 {
					if sf[0] == "TO_URL" {
						// parse the url after trimming off any surrounding double quotes
						cfg.TOUrl = strings.Trim(sf[1], "\"")
						to_url_found = true
					}
					if sf[0] == "TO_USER" {
						// set the TOUser after trimming off any surrounding quotes.
						cfg.TOUser = strings.Trim(sf[1], "\"")
						to_user_found = true
					}
					// set the TOPass after trimming off any surrounding quotes.
					if sf[0] == "TO_PASS" {
						cfg.TOPass = strings.Trim(sf[1], "\"")
						to_pass_found = true
					}
				}
			}
		}
	}
	if !to_url_found && !to_user_found && !to_pass_found {
		return errors.New("failed to retrieve one or more TrafficOps credentails")
	}

	modTime, err := util.GetFileModificationTime(fn.Filename)
	if err != nil {
		return fmt.Errorf("could not stat %s: %w", fn.Filename, err)
	}
	cfg.CredentialFile.LastModifyTime = modTime

	return nil
}

func GetConfig() (Cfg, error, bool) {
	var err error
	var configFile string
	var logLocationErr = log.LogLocationStderr
	var logLocationDebug = log.LogLocationNull
	var logLocationInfo = log.LogLocationNull
	var logLocationWarn = log.LogLocationNull

	configFilePtr := getopt.StringLong("config-file", 'f', DefaultConfigFile, "full path to the json config file")
	logdirPtr := getopt.StringLong("logging-dir", 'l', DefaultLogDirectory, "directory location for log files")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	verbosePtr := getopt.CounterLong("verbose", 'v', `Log verbosity. Logging is output to stderr. By default, errors are logged. To log warnings, pass '-v'. To log info, pass '-vv', debug pass '-vvv'`)

	getopt.Parse()

	if configFilePtr != nil {
		configFile = *configFilePtr
	} else {
		configFile = DefaultConfigFile
	}

	var logfile string

	logfile = filepath.Join(*logdirPtr, DefaultLogFile)

	logLocationErr = logfile

	if *verbosePtr == 1 {
		logLocationWarn = logfile
	} else if *verbosePtr == 2 {
		logLocationInfo = logfile
		logLocationWarn = logfile
	} else if *verbosePtr == 3 {
		logLocationInfo = logfile
		logLocationWarn = logfile
		logLocationDebug = logfile
	}

	if help := *helpPtr; help == true {
		Usage()
		return Cfg{}, nil, true
	}

	lcfg := LogCfg{
		LogLocationDebug: logLocationDebug,
		LogLocationErr:   logLocationErr,
		LogLocationInfo:  logLocationInfo,
		LogLocationWarn:  logLocationWarn,
	}

	if err := log.InitCfg(&lcfg); err != nil {
		return Cfg{}, errors.New("initializing loggers: " + err.Error() + "\n"), false
	}

	cf := util.ConfigFile{
		Filename:       configFile,
		LastModifyTime: 0,
	}

	cfg := Cfg{
		HealthClientConfigFile: cf,
		CredentialFile:         util.ConfigFile{},
	}

	if _, err = LoadConfig(&cfg); err != nil {
		return Cfg{}, errors.New(err.Error() + "\n"), false
	}

	if err = ReadCredentials(&cfg, false); err != nil {
		return cfg, err, false
	}

	dispersion := GetTOLoginDispersion(cfg.TOLoginDispersionFactor)
	log.Infof("waiting %v seconds before logging into TrafficOps", dispersion.Seconds())
	time.Sleep(dispersion)
	err = GetTrafficMonitors(&cfg)
	if err != nil {
		return cfg, err, false
	}

	return cfg, nil, false
}

func GetTrafficMonitors(cfg *Cfg) error {
	qry := &url.Values{}
	qry.Add("type", "RASCAL")
	qry.Add("status", "ONLINE")

	// login to traffic ops.
	if toSession == nil {
		session, _, err := toclient.LoginWithAgent(cfg.TOUrl, cfg.TOUser, cfg.TOPass, true, userAgent, false, GetRequestTimeout())
		if err != nil {
			return fmt.Errorf("could not establish a TrafficOps session: %w", err)
		} else {
			toSession = session
		}
	}
	srvs, _, err := toSession.GetServersWithHdr(qry, nil)
	if err != nil {
		// next time we'll login again and get a new session.
		toSession = nil
		return errors.New("error fetching Trafficmonitor server list: " + err.Error())
	}

	cfg.TrafficMonitors = make(map[string]bool, 0)
	for _, v := range srvs.Response {
		if *v.CDNName == cfg.CDNName && *v.Status == "ONLINE" {
			hostname := *v.HostName + "." + *v.DomainName
			cfg.TrafficMonitors[hostname] = true
		}
	}

	return nil
}

func GetTMPollingInterval() time.Duration {
	return tmPollingInterval
}

func GetTOLoginDispersion(dispersionFactor int) time.Duration {
	dispersionSeconds := uint64(tmPollingInterval.Seconds()) * uint64(dispersionFactor)
	hostName, err := os.Hostname()
	if err != nil {
		log.Errorf("the OS hostname is not set, cannot continue: %s", err.Error())
		os.Exit(1)
	}
	md5hash := md5.Sum([]byte(hostName))
	sl := md5hash[0:8]
	disp := (binary.BigEndian.Uint64(sl) % dispersionSeconds)
	if disp < uint64(tmPollingInterval.Seconds()*2) {
		disp = ((disp * 2) + uint64(tmPollingInterval.Seconds()))
	}
	return time.Duration(disp) * time.Second
}

func GetRequestTimeout() time.Duration {
	return toRequestTimeout
}

func LoadConfig(cfg *Cfg) (bool, error) {
	updated := false
	configFile := cfg.HealthClientConfigFile.Filename
	modTime, err := util.GetFileModificationTime(configFile)
	if err != nil {
		return updated, errors.New(err.Error())
	}

	if modTime > cfg.HealthClientConfigFile.LastModifyTime {
		log.Infoln("Loading a new config file.")
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return updated, errors.New(err.Error())
		}
		err = json.Unmarshal(content, cfg)
		if err != nil {
			return updated, fmt.Errorf("config parsing failed: %w", err)
		}
		tmPollingInterval, err = time.ParseDuration(cfg.TmPollIntervalSeconds)
		if err != nil {
			return updated, errors.New("parsing TMPollingIntervalSeconds: " + err.Error())
		}
		if cfg.TOLoginDispersionFactor == 0 {
			cfg.TOLoginDispersionFactor = DefaultTOLoginDispersionFactor
		}
		toRequestTimeout, err = time.ParseDuration(cfg.TORequestTimeOutSeconds)
		if err != nil {
			return updated, errors.New("parsing TORequestTimeOutSeconds: " + err.Error())
		}
		if cfg.ReasonCode != "active" && cfg.ReasonCode != "local" {
			return updated, errors.New("invalid reason-code: " + cfg.ReasonCode + ", valid reason codes are 'active' or 'local'")
		}
		if cfg.TrafficServerConfigDir == "" {
			cfg.TrafficServerConfigDir = DefaultTrafficServerConfigDir
		}
		if cfg.TrafficServerBinDir == "" {
			cfg.TrafficServerBinDir = DefaultTrafficServerBinDir
		}
		if cfg.UnavailablePollThreshold == 0 {
			cfg.UnavailablePollThreshold = DefaultUnavailablePollThreshold
		}
		if cfg.PollStateJSONLog == "" {
			cfg.PollStateJSONLog = DefaultPollStateJSONLog
		}

		cfg.HealthClientConfigFile.LastModifyTime = modTime

		if cfg.TOCredentialFile != "" {
			cfg.CredentialFile.Filename = cfg.TOCredentialFile
		}

		// if tm-proxy-url is set in the config, verify the proxy
		// url
		if cfg.TmProxyURL != "" {
			if cfg.ParsedProxyURL, err = url.Parse(cfg.TmProxyURL); err != nil {
				cfg.ParsedProxyURL = nil
				return false, errors.New("parsing TmProxyUrl: " + err.Error())
			}
			if cfg.ParsedProxyURL.Port() == "" {
				cfg.ParsedProxyURL = nil
				return false, errors.New("TmProxyUrl invalid port specified")
			}
			log.Infof("TM queries will use the proxy: %s", cfg.TmProxyURL)
		} else {
			cfg.ParsedProxyURL = nil
		}
		updated = true
	}
	return updated, nil
}

func UpdateConfig(cfg *Cfg, newCfg *Cfg) {
	cfg.CDNName = newCfg.CDNName
	cfg.EnableActiveMarkdowns = newCfg.EnableActiveMarkdowns
	cfg.ReasonCode = newCfg.ReasonCode
	cfg.TOCredentialFile = newCfg.TOCredentialFile
	cfg.TORequestTimeOutSeconds = newCfg.TORequestTimeOutSeconds
	cfg.TOPass = newCfg.TOPass
	cfg.TOUrl = newCfg.TOUrl
	cfg.TOUser = newCfg.TOUser
	cfg.TmPollIntervalSeconds = newCfg.TmPollIntervalSeconds
	cfg.TOLoginDispersionFactor = newCfg.TOLoginDispersionFactor
	if cfg.TOLoginDispersionFactor == 0 {
		cfg.TOLoginDispersionFactor = DefaultTOLoginDispersionFactor
	}
	cfg.UnavailablePollThreshold = newCfg.UnavailablePollThreshold
	cfg.TrafficServerConfigDir = newCfg.TrafficServerConfigDir
	cfg.TrafficServerBinDir = newCfg.TrafficServerBinDir
	cfg.TrafficMonitors = newCfg.TrafficMonitors
	cfg.HealthClientConfigFile = newCfg.HealthClientConfigFile
	cfg.PollStateJSONLog = newCfg.PollStateJSONLog
	cfg.EnablePollStateLog = newCfg.EnablePollStateLog
}

func Usage() {
	getopt.PrintUsage(os.Stdout)
}
