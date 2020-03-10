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
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"

	flag "github.com/ogier/pflag"
)

const AppName = "atstccfg"
const Version = "0.2"
const UserAgent = AppName + "/" + Version

const ExitCodeSuccess = 0
const ExitCodeErrGeneric = 1
const ExitCodeNotFound = 104
const ExitCodeBadRequest = 100

var ErrNotFound = errors.New("not found")
var ErrBadRequest = errors.New("bad request")

type Cfg struct {
	CacheHostName       string
	ListPlugins         bool
	LogLocationErr      string
	LogLocationInfo     string
	LogLocationWarn     string
	NumRetries          int
	OutputDir           string
	PrintGeneratedFiles bool
	TOInsecure          bool
	TOPass              string
	TOTimeout           time.Duration
	TOURL               *url.URL
	TOUser              string
	GetData             string
	SetQueueStatus      string
	SetRevalStatus      string
	RevalOnly           bool
}

type TCCfg struct {
	Cfg
	TOClient **toclient.Session
}

func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationErr) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // atstccfg doesn't use the debug logger, use Info instead.
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // atstccfg doesn't use the event logger.

// GetCfg gets the application configuration, from arguments and environment variables.
// Note if PrintGeneratedFiles is configured, the config will be returned with PrintGeneratedFiles true and all other values set to their defaults. This is because other values may have requirements and return errors, where if PrintGeneratedFiles is set by the user, no other setting should be considered.
func GetCfg() (Cfg, error) {
	toURLPtr := flag.StringP("traffic-ops-url", "u", "", "Traffic Ops URL. Must be the full URL, including the scheme. Required. May also be set with the environment variable TO_URL.")
	toUserPtr := flag.StringP("traffic-ops-user", "U", "", "Traffic Ops username. Required. May also be set with the environment variable TO_USER.")
	toPassPtr := flag.StringP("traffic-ops-password", "P", "", "Traffic Ops password. Required. May also be set with the environment variable TO_PASS.")
	numRetriesPtr := flag.IntP("num-retries", "r", 5, "The number of times to retry getting a file if it fails.")
	logLocationErrPtr := flag.StringP("log-location-error", "e", "stderr", "Where to log errors. May be a file path, stdout, stderr, or null.")
	logLocationWarnPtr := flag.StringP("log-location-warning", "w", "stderr", "Where to log warnings. May be a file path, stdout, stderr, or null.")
	logLocationInfoPtr := flag.StringP("log-location-info", "i", "stderr", "Where to log information messages. May be a file path, stdout, stderr, or null.")
	printGeneratedFilesPtr := flag.BoolP("print-generated-files", "g", false, "Whether to print a list of files which are generated (and not proxied to Traffic Ops).")
	toInsecurePtr := flag.BoolP("traffic-ops-insecure", "s", false, "Whether to ignore HTTPS certificate errors from Traffic Ops. It is HIGHLY RECOMMENDED to never use this in a production environment, but only for debugging.")
	toTimeoutMSPtr := flag.IntP("traffic-ops-timeout-milliseconds", "t", 60000, "Timeout in seconds for Traffic Ops requests.")
	versionPtr := flag.BoolP("version", "v", false, "Print version information and exit.")
	listPluginsPtr := flag.BoolP("list-plugins", "l", false, "Print the list of plugins.")
	helpPtr := flag.BoolP("help", "h", false, "Print usage information and exit")
	outputDirPtr := flag.StringP("output-directory", "o", "", "Directory to output config files to.")
	cacheHostNamePtr := flag.StringP("cache-host-name", "n", "", "Host name of the cache to generate config for. Must be the server host name in Traffic Ops, not a URL, and not the FQDN")
	getDataPtr := flag.StringP("get-data", "d", "", "non-config-file Traffic Ops Data to get. Valid values are update-status, packages, chkconfig, system-info, and statuses")
	setQueueStatusPtr := flag.StringP("set-queue-status", "q", "", "POSTs to Traffic Ops setting the queue status of the server. Must be 'true' or 'false'. Requires --set-reval-status also be set")
	setRevalStatusPtr := flag.StringP("set-reval-status", "a", "", "POSTs to Traffic Ops setting the revaliate status of the server. Must be 'true' or 'false'. Requires --set-queue-status also be set")
	revalOnlyPtr := flag.BoolP("revalidate-only", "y", false, "Whether to exclude files not named 'regex_revalidate.config'")

	flag.Parse()

	if *versionPtr {
		fmt.Println(AppName + " v" + Version)
		os.Exit(0)
	} else if *helpPtr {
		flag.PrintDefaults()
		os.Exit(0)
	} else if *printGeneratedFilesPtr {
		return Cfg{PrintGeneratedFiles: true}, nil
	} else if *listPluginsPtr {
		return Cfg{ListPlugins: true}, nil
	}

	toURL := *toURLPtr
	toUser := *toUserPtr
	toPass := *toPassPtr
	numRetries := *numRetriesPtr
	logLocationErr := *logLocationErrPtr
	logLocationWarn := *logLocationWarnPtr
	logLocationInfo := *logLocationInfoPtr
	toInsecure := *toInsecurePtr
	toTimeout := time.Millisecond * time.Duration(*toTimeoutMSPtr)
	listPlugins := *listPluginsPtr
	outputDir := *outputDirPtr
	cacheHostName := *cacheHostNamePtr
	getData := *getDataPtr
	setQueueStatus := *setQueueStatusPtr
	setRevalStatus := *setRevalStatusPtr
	revalOnly := *revalOnlyPtr

	urlSourceStr := "argument" // for error messages
	if toURL == "" {
		urlSourceStr = "environment variable"
		toURL = os.Getenv("TO_URL")
	}
	if toUser == "" {
		toUser = os.Getenv("TO_USER")
	}
	if toPass == "" {
		toPass = os.Getenv("TO_PASS")
	}

	usageStr := "Usage: ./" + AppName + " --traffic-ops-url=myurl --traffic-ops-user=myuser --traffic-ops-password=mypass --cache-host-name=my-cache --output-directory=/opt/trafficserver/etc/trafficserver-temp/"
	if strings.TrimSpace(toURL) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-url or TO_URL environment variable. " + usageStr)
	}
	if strings.TrimSpace(toUser) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-user or TO_USER environment variable. " + usageStr)
	}
	if strings.TrimSpace(toPass) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-password or TO_PASS environment variable. " + usageStr)
	}
	// if strings.TrimSpace(outputDir) == "" {
	// 	return Cfg{}, errors.New("Missing required argument --output-directory. If you wish to use the current directory, pass '.'. " + usageStr)
	// }
	if strings.TrimSpace(cacheHostName) == "" {
		return Cfg{}, errors.New("Missing required argument --cache-host-name. " + usageStr)
	}

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return Cfg{}, errors.New("parsing Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	} else if err := ValidateURL(toURLParsed); err != nil {
		return Cfg{}, errors.New("invalid Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	}

	cfg := Cfg{
		LogLocationErr:  logLocationErr,
		LogLocationWarn: logLocationWarn,
		LogLocationInfo: logLocationInfo,
		NumRetries:      numRetries,
		TOInsecure:      toInsecure,
		TOPass:          toPass,
		TOTimeout:       toTimeout,
		TOURL:           toURLParsed,
		TOUser:          toUser,
		ListPlugins:     listPlugins,
		OutputDir:       outputDir,
		CacheHostName:   cacheHostName,
		GetData:         getData,
		SetRevalStatus:  setRevalStatus,
		SetQueueStatus:  setQueueStatus,
		RevalOnly:       revalOnly,
	}

	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}

	if err := ValidateDirWriteable(outputDir); err != nil {
		return Cfg{}, errors.New("validating output directory is writeable '" + outputDir + "': " + err.Error())
	}

	return cfg, nil
}

func ValidateURL(u *url.URL) error {
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

func ValidateDirWriteable(dir string) error {
	testFileName := "testwrite.txt"
	testFilePath := filepath.Join(dir, testFileName)
	if err := os.RemoveAll(testFilePath); err != nil {
		// TODO don't log? This can be normal
		log.Infoln("error removing temp test file '" + testFilePath + "' (ok if it didn't exist): " + err.Error())
	}

	fl, err := os.Create(testFilePath)
	if err != nil {
		return errors.New("creating temp test file '" + testFilePath + "': " + err.Error())
	}
	defer fl.Close()

	if _, err := fl.WriteString("test"); err != nil {
		return errors.New("writing to temp test file '" + testFilePath + "': " + err.Error())
	}

	return nil
}

func RetryBackoffSeconds(currentRetry int) int {
	// TODO make configurable?
	return int(math.Pow(2.0, float64(currentRetry)))
}
