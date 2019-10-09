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
	"errors"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"

	flag "github.com/ogier/pflag"
)

type Cfg struct {
	CacheFileMaxAge     time.Duration
	LogLocationErr      string
	LogLocationInfo     string
	LogLocationWarn     string
	NumRetries          int
	TempDir             string
	TOInsecure          bool
	TOPass              string
	TOTimeout           time.Duration
	TOURL               *url.URL
	TOUser              string
	PrintGeneratedFiles bool
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
	noCachePtr := flag.BoolP("no-cache", "n", false, "Whether not to use existing cache files. Optional. Cache files will still be created, existing ones just won't be used.")
	numRetriesPtr := flag.IntP("num-retries", "r", 5, "The number of times to retry getting a file if it fails.")
	logLocationErrPtr := flag.StringP("log-location-error", "e", "stderr", "Where to log errors. May be a file path, stdout, stderr, or null.")
	logLocationWarnPtr := flag.StringP("log-location-warning", "w", "stderr", "Where to log warnings. May be a file path, stdout, stderr, or null.")
	logLocationInfoPtr := flag.StringP("log-location-info", "i", "stderr", "Where to log information messages. May be a file path, stdout, stderr, or null.")
	printGeneratedFilesPtr := flag.BoolP("print-generated-files", "g", false, "Whether to print a list of files which are generated (and not proxied to Traffic Ops).")
	toInsecurePtr := flag.BoolP("traffic-ops-insecure", "s", false, "Whether to ignore HTTPS certificate errors from Traffic Ops. It is HIGHLY RECOMMENDED to never use this in a production environment, but only for debugging.")
	toTimeoutMSPtr := flag.IntP("traffic-ops-timeout-milliseconds", "t", 10000, "Timeout in seconds for Traffic Ops requests.")
	cacheFileMaxAgeSecondsPtr := flag.IntP("cache-file-max-age-seconds", "a", 60, "Maximum age to use cached files.")
	flag.Parse()

	if *printGeneratedFilesPtr {
		return Cfg{PrintGeneratedFiles: true}, nil
	}

	toURL := *toURLPtr
	toUser := *toUserPtr
	toPass := *toPassPtr
	noCache := *noCachePtr
	numRetries := *numRetriesPtr
	logLocationErr := *logLocationErrPtr
	logLocationWarn := *logLocationWarnPtr
	logLocationInfo := *logLocationInfoPtr
	toInsecure := *toInsecurePtr
	toTimeout := time.Millisecond * time.Duration(*toTimeoutMSPtr)
	cacheFileMaxAge := time.Second * time.Duration(*cacheFileMaxAgeSecondsPtr)

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

	if strings.TrimSpace(toURL) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-url or TO_URL environment variable. Usage: ./" + AppName + " --traffic-ops-url myurl --traffic-ops-user myuser --traffic-ops-password mypass")
	}
	if strings.TrimSpace(toUser) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-user or TO_USER environment variable. Usage: ./" + AppName + " --traffic-ops-url myurl --traffic-ops-user myuser --traffic-ops-password mypass")
	}
	if strings.TrimSpace(toPass) == "" {
		return Cfg{}, errors.New("Missing required argument --traffic-ops-password or TO_PASS environment variable. Usage: ./" + AppName + " --traffic-ops-url myurl --traffic-ops-user myuser --traffic-ops-password mypass")
	}

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return Cfg{}, errors.New("parsing Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	} else if err := ValidateURL(toURLParsed); err != nil {
		return Cfg{}, errors.New("invalid Traffic Ops URL from " + urlSourceStr + " '" + toURL + "': " + err.Error())
	}

	tmpDir := os.TempDir()
	tmpDir = filepath.Join(tmpDir, TempSubdir)

	cfg := Cfg{
		CacheFileMaxAge: cacheFileMaxAge,
		LogLocationErr:  logLocationErr,
		LogLocationWarn: logLocationWarn,
		LogLocationInfo: logLocationInfo,
		NumRetries:      numRetries,
		TempDir:         tmpDir,
		TOInsecure:      toInsecure,
		TOPass:          toPass,
		TOTimeout:       toTimeout,
		TOURL:           toURLParsed,
		TOUser:          toUser,
	}

	if err := log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}

	if noCache {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Errorln("deleting cache directory '" + tmpDir + "': " + err.Error())
		}
	}

	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return Cfg{}, errors.New("creating temp directory '" + tmpDir + "': " + err.Error())
	}
	if err := ValidateDirWriteable(tmpDir); err != nil {
		return Cfg{}, errors.New("validating temp directory is writeable '" + tmpDir + "': " + err.Error())
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
