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

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/tc-health-client/util"

	"github.com/pborman/getopt/v2"
)

// userAgent is the UA used by this service in HTTP requests.
// TODO dynamically add Version from RPM build.
const userAgent = "tc-health-client/1.0"

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

// NewCfgPtr is a convenience func for NewAtomicPtr(cfg).
func NewCfgPtr(cfg *Cfg) *util.AtomicPtr[Cfg] {
	return util.NewAtomicPtr(cfg)
}

type HealthMethod string

const (
	HealthMethodTrafficMonitor = HealthMethod("traffic-monitor")
	HealthMethodParentL4       = HealthMethod("parent-l4")
	HealthMethodParentL7       = HealthMethod("parent-l7")
	HealthMethodParentService  = HealthMethod("parent-service")
)

var DefaultHealthMethods = []HealthMethod{HealthMethodTrafficMonitor}
var DefaultMarkdownMethods = []HealthMethod{HealthMethodTrafficMonitor}

type Cfg struct {
	// TODO fix to not take CDN name, so it can't be wrong.
	//      If this is wrong, the health client runs but just can't find tons of stuff,
	//      resulting in strange not-obvious errors.
	// Infer from server somehow (config header comment? t3c metadata? hostname -s?)
	CDNName                  string `json:"cdn-name"`
	EnableActiveMarkdowns    bool   `json:"enable-active-markdowns"`
	ReasonCode               string `json:"reason-code"`
	TOCredentialFile         string `json:"to-credential-file"`
	TORequestTimeOutSeconds  string `json:"to-request-timeout-seconds"`
	TOPass                   string `json:"to-pass"`
	TOUrl                    string `json:"to-url"`
	TOUser                   string `json:"to-user"`
	TmProxyURL               string `json:"tm-proxy-url"`
	TmPollIntervalSeconds    string `json:"tm-poll-interval-seconds"`
	TOLoginDispersionFactor  int    `json:"to-login-dispersion-factor"`
	UnavailablePollThreshold int    `json:"unavailable-poll-threshold"`
	MarkUpPollThreshold      int    `json:"markup-poll-threshold"`
	TrafficServerConfigDir   string `json:"trafficserver-config-dir"`
	TrafficServerBinDir      string `json:"trafficserver-bin-dir"`
	PollStateJSONLog         string `json:"poll-state-json-log"`
	EnablePollStateLog       bool   `json:"enable-poll-state-log"`
	HealthClientConfigFile   util.ConfigFile
	CredentialFile           util.ConfigFile
	ParsedProxyURL           *url.URL

	MarkdownMinIntervalMS     *uint64 `json:"markdown-min-interval-ms"`
	ParentHealthL4PollMS      uint64  `json:"parent-health-l4-poll-ms"`
	ParentHealthL7PollMS      uint64  `json:"parent-health-l7-poll-ms"`
	ParentHealthServicePollMS uint64  `json:"parent-health-service-poll-ms"`

	// ParentHealthServicePort is the port to serve parent health on. To disable serving parent health, set to 0.
	ParentHealthServicePort int `json:"parent-health-service-port"`

	// ParentHealthLogLocation may be a file path, stdout, stderr, or null (the empty string is equivalent to null)
	ParentHealthLogLocation string `json:"parent-health-log-location"`

	// HealthMethods are the types of health to poll. If omitted, Traffic Monitor health will be used.
	HealthMethods *[]HealthMethod `json:"health-methods"`

	// MarkdownMethods are the types of health to consider when marking down parents.
	// This may be empty or omitted, to never marking down parents.
	MarkdownMethods *[]HealthMethod `json:"markdown-methods"`

	// NumHealthWorkers is the number of worker microthreads (goroutines) per health poll method.
	// Note this only applies to Parent L4, Parent L7, and Parent Service health; the Traffic
	// Monitor health poll is a single HTTP request and thus doesn't need workers.
	NumHealthWorkers int `json:"num-health-workers"`

	// Monitor peers inside strategies.yaml file
	MonitorStrategiesPeers bool `json:"monitor-strategies-peers"`

	TMPollingInterval time.Duration
	TORequestTimeout  time.Duration

	UserAgent string
	// HostName is this host's short hostname.
	// If empty, the system's hostname will be used.
	HostName string
}

func (cfg *Cfg) Clone() *Cfg {
	newCfg := *cfg
	newCfg.ParsedProxyURL = nil
	if cfg.ParsedProxyURL != nil {
		urlCopy := *cfg.ParsedProxyURL
		newCfg.ParsedProxyURL = &urlCopy
	}

	if cfg.HealthMethods != nil {
		healthMethods := make([]HealthMethod, len(*cfg.HealthMethods), len(*cfg.HealthMethods))
		copy(healthMethods, *cfg.HealthMethods)
		newCfg.HealthMethods = &healthMethods
	} else {
		newCfg.HealthMethods = &[]HealthMethod{}
	}

	if cfg.MarkdownMethods != nil {
		markdownMethods := make([]HealthMethod, len(*cfg.MarkdownMethods), len(*cfg.MarkdownMethods))
		copy(markdownMethods, *cfg.MarkdownMethods)
		newCfg.MarkdownMethods = &markdownMethods
	} else {
		newCfg.MarkdownMethods = &[]HealthMethod{}
	}

	return &newCfg
}

type LogCfg struct {
	LogLocationErr   string
	LogLocationDebug string
	LogLocationInfo  string
	LogLocationWarn  string
	LogLocationEvent string
}

func (lcfg LogCfg) ErrorLog() log.LogLocation   { return log.LogLocation(lcfg.LogLocationErr) }
func (lcfg LogCfg) WarningLog() log.LogLocation { return log.LogLocation(lcfg.LogLocationWarn) }
func (lcfg LogCfg) InfoLog() log.LogLocation    { return log.LogLocation(lcfg.LogLocationInfo) }
func (lcfg LogCfg) DebugLog() log.LogLocation   { return log.LogLocation(lcfg.LogLocationDebug) }
func (lcfg LogCfg) EventLog() log.LogLocation   { return log.LogLocation(lcfg.LogLocationEvent) }

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

	err := error(nil)
	cfg.TOUrl, cfg.TOUser, cfg.TOPass, err = getCredentialsFromFile(cfg.CredentialFile.Filename)
	if err != nil {
		return errors.New("reading credentials from file '" + fn.Filename + "' :" + err.Error())
	}

	if cfg.TOUrl == "" || cfg.TOUser == "" || cfg.TOPass == "" {
		return errors.New("failed to retrieve one or more TrafficOps credentails")
	}

	modTime, err := util.GetFileModificationTime(fn.Filename)
	if err != nil {
		return fmt.Errorf("could not stat %s: %w", fn.Filename, err)
	}
	cfg.CredentialFile.LastModifyTime = modTime

	return nil
}

func GetConfig() (*Cfg, error, bool) {
	var err error
	var configFile string
	var logLocationErr = log.LogLocationStderr
	var logLocationDebug = log.LogLocationNull
	var logLocationInfo = log.LogLocationFile
	var logLocationWarn = log.LogLocationNull
	var logLocationEvent = log.LogLocationNull

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
		logLocationEvent = logfile
		logLocationDebug = logfile
	}

	if help := *helpPtr; help == true {
		Usage()
		return nil, nil, true
	}

	lcfg := LogCfg{
		LogLocationDebug: logLocationDebug,
		LogLocationErr:   logLocationErr,
		LogLocationInfo:  logLocationInfo,
		LogLocationWarn:  logLocationWarn,
		LogLocationEvent: logLocationEvent,
	}

	if err := log.InitCfg(&lcfg); err != nil {
		return nil, errors.New("initializing loggers: " + err.Error() + "\n"), false
	}

	cf := util.ConfigFile{
		Filename:       configFile,
		LastModifyTime: 0,
	}

	cfg := &Cfg{
		HealthClientConfigFile: cf,
		CredentialFile:         util.ConfigFile{},
		UserAgent:              userAgent,
	}

	if _, err = LoadConfig(cfg); err != nil {
		return nil, errors.New(err.Error() + "\n"), false
	}

	if err = ReadCredentials(cfg, false); err != nil {
		return cfg, err, false
	}

	dispersion := GetTOLoginDispersion(cfg.TMPollingInterval, cfg.TOLoginDispersionFactor)
	log.Infof("waiting %v seconds before logging into TrafficOps", dispersion.Seconds())
	time.Sleep(dispersion)

	return cfg, nil, false
}

func GetTOLoginDispersion(pollingInterval time.Duration, dispersionFactor int) time.Duration {
	dispersionSeconds := uint64(pollingInterval.Seconds()) * uint64(dispersionFactor)
	hostName, err := os.Hostname()
	if err != nil {
		log.Errorf("the OS hostname is not set, cannot continue: %s", err.Error())
		os.Exit(1)
	}
	md5hash := md5.Sum([]byte(hostName))
	sl := md5hash[0:8]
	disp := (binary.BigEndian.Uint64(sl) % dispersionSeconds)
	if disp < uint64(pollingInterval.Seconds()*2) {
		disp = ((disp * 2) + uint64(pollingInterval.Seconds()))
	}
	return time.Duration(disp) * time.Second
}

const DefaultParentHealthL4PollMS = 30000
const DefaultParentHealthL7PollMS = 30000
const ParentHealthServicePollMS = 30000
const DefaultMarkdownMinIntervalMS = 5000

// LoadConfig returns whether the config was updated and any error.
//
// Note the cfg may be modified, even if the returned updated is false.
// Users should create and pass a copy of Cfg if changes on error are not acceptable.
func LoadConfig(cfg *Cfg) (bool, error) {
	configFile := cfg.HealthClientConfigFile.Filename
	// Load default value for strategies.yaml
	cfg.MonitorStrategiesPeers = true
	modTime, err := util.GetFileModificationTime(configFile)
	if err != nil {
		return false, errors.New(err.Error())
	}

	if modTime <= cfg.HealthClientConfigFile.LastModifyTime {
		return false, nil
	}

	log.Infoln("Loading a new config file.")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return false, errors.New(err.Error())
	}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		return false, fmt.Errorf("config parsing failed: %w", err)
	}
	cfg.TMPollingInterval, err = time.ParseDuration(cfg.TmPollIntervalSeconds)
	if err != nil {
		return false, errors.New("parsing TMPollingIntervalSeconds: " + err.Error())
	}
	if cfg.TOLoginDispersionFactor == 0 {
		cfg.TOLoginDispersionFactor = DefaultTOLoginDispersionFactor
	}
	cfg.TORequestTimeout, err = time.ParseDuration(cfg.TORequestTimeOutSeconds)
	if err != nil {
		return false, errors.New("parsing TORequestTimeOutSeconds: " + err.Error())
	}
	if cfg.ReasonCode != "active" && cfg.ReasonCode != "local" {
		return false, errors.New("invalid reason-code: " + cfg.ReasonCode + ", valid reason codes are 'active' or 'local'")
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

	if cfg.HostName == "" {
		hostName, err := os.Hostname()
		if err != nil {
			return false, errors.New("No hostname configured, getting from OS: " + err.Error())
		}
		cfg.HostName = util.HostNameToShort(hostName)
	}

	if cfg.HealthMethods == nil {
		healthMethods := make([]HealthMethod, len(DefaultHealthMethods), len(DefaultHealthMethods))
		copy(healthMethods, DefaultHealthMethods)
		cfg.HealthMethods = &healthMethods
	}

	if cfg.MarkdownMethods == nil {
		markdownMethods := make([]HealthMethod, len(DefaultMarkdownMethods), len(DefaultMarkdownMethods))
		copy(markdownMethods, DefaultMarkdownMethods)
		cfg.MarkdownMethods = &markdownMethods
	}

	if cfg.NumHealthWorkers < 1 {
		cfg.NumHealthWorkers = 1
	}

	if cfg.MarkdownMinIntervalMS == nil {
		mi := uint64(DefaultMarkdownMinIntervalMS)
		cfg.MarkdownMinIntervalMS = &mi
	}
	if cfg.ParentHealthL4PollMS <= 0 {
		cfg.ParentHealthL4PollMS = DefaultParentHealthL4PollMS
	}
	if cfg.ParentHealthL7PollMS <= 0 {
		cfg.ParentHealthL7PollMS = DefaultParentHealthL7PollMS
	}
	if cfg.ParentHealthServicePollMS <= 0 {
		cfg.ParentHealthServicePollMS = ParentHealthServicePollMS
	}

	return true, nil
}

func Usage() {
	getopt.PrintUsage(os.Stdout)
}

// getCredentialsFromFile gets the TO URL, user, and password from an environment variable file.
// from environment variables declared in a credentials file bash script, if they exist.
//
// Returns the TO URL, user, password, and any error.
//
// Note this returns empty strings with no error if the file doesn't exist,
// or if any variables aren't declared.
func getCredentialsFromFile(filePath string) (string, string, string, error) {

	if inf, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", "", "", nil
	} else if inf.IsDir() {
		return "", "", "", errors.New("credentials path is a directory, must be a file")
	}

	// we execute sh and source the file to get the environment variables,
	// because it's easier and more accurate than writing our own sh env var parser.

	stdOut, stdErr, code := t3cutil.Do("sh", "-c", `(source "`+filePath+`" && printf "${TO_URL}\n")`)
	if code != 0 {
		return "", "", "", fmt.Errorf("getting credentials from file returned error code %v stderr '%v' stdout '%v'", code, string(stdErr), string(stdOut))
	}
	toURL := strings.TrimSpace(string(stdOut))

	stdOut, stdErr, code = t3cutil.Do("sh", "-c", `(source "`+filePath+`" && printf "${TO_USER}\n")`)
	if code != 0 {
		return "", "", "", fmt.Errorf("getting credentials from file returned error code %v stderr '%v' stdout '%v'", code, string(stdErr), string(stdOut))
	}
	toUser := strings.TrimSpace(string(stdOut))

	stdOut, stdErr, code = t3cutil.Do("sh", "-c", `(source "`+filePath+`" && printf "${TO_PASS}\n")`)
	if code != 0 {
		return "", "", "", fmt.Errorf("getting credentials from file returned error code %v stderr '%v' stdout '%v'", code, string(stdErr), string(stdOut))
	}
	toPass := strings.TrimSpace(string(stdOut))

	return toURL, toUser, toPass, nil
}
