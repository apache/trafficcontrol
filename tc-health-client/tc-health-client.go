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
	"context"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/sys/unix"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/tc-health-client/config"
	"github.com/apache/trafficcontrol/v8/tc-health-client/tmagent"
	"github.com/apache/trafficcontrol/v8/tc-health-client/util"
)

const (
	Success      = 0
	ConfigError  = 166
	RunTimeError = 167
	PidFile      = "/run/tc-health-client.pid"
)

// the BuildTimestamp and Version are set via ld flags
// when the RPM is built, see build/build_rpm.sh
var (
	BuildTimestamp = ""
	Version        = ""
)

func main() {
	rand.Seed(time.Now().UnixNano()) // TODO make deterministic. Seed hostname?

	cfg, err, helpflag := config.GetConfig()
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(ConfigError)
	}
	cfgPtr := config.NewCfgPtr(cfg)

	if helpflag { // user used --help option
		os.Exit(Success)
	}

	log.Infof("Polling interval: %v seconds\n", cfg.TMPollingInterval.Seconds())
	tmInfo, err := tmagent.NewParentInfo(cfgPtr)
	if err != nil {
		log.Errorf("startup could not initialize parent info, check that trafficserver is running: %s\n", err.Error())
		os.Exit(RunTimeError)
	}

	if err := tmInfo.GetTOData(cfg); err != nil {
		log.Errorln("startup could not get data from Traffic Ops: " + err.Error())
		os.Exit(RunTimeError)
	}

	pid := os.Getpid()
	err = os.WriteFile(PidFile, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Errorf("could not write the process id to %s: %s", PidFile, err.Error())
		os.Exit(RunTimeError)
	}

	log.Infof("startup complete, version: %s, built: %s\n", Version, BuildTimestamp)

	// TODO: make configurable
	cfg.ParentHealthLogLocation = filepath.Join(config.DefaultLogDirectory, "parent-health.log")

	switch cfg.ParentHealthLogLocation {
	case "":
		fallthrough
	case "null":
		tmInfo.ParentHealthLog = util.NopWriter{}
		log.Infoln("Logging parent health to /dev/null")
	case "stdout":
		tmInfo.ParentHealthLog = util.MakeNoCloser(os.Stdout)
		log.Infoln("Logging parent health to out")
	case "stderr":
		tmInfo.ParentHealthLog = util.MakeNoCloser(os.Stderr)
		log.Infoln("Logging parent health to stderr")
	default:
		logFi, err := os.OpenFile(cfg.ParentHealthLogLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Errorf("Opening parent health log file '" + cfg.ParentHealthLogLocation + "': " + err.Error())
			os.Exit(RunTimeError)
		}
		log.Infoln("Logging parent health to file '" + cfg.ParentHealthLogLocation + "'")
		tmInfo.ParentHealthLog = logFi
	}

	// start goroutines

	log.Infoln("Starting Internal Services")

	log.Infof("Health Methods: %+v\n", cfg.HealthMethods)
	log.Infof("Markdown Methods: %+v\n", cfg.HealthMethods)
	log.Infof("Num Health Workers: %+v\n", cfg.NumHealthWorkers)
	log.Infof("Parent Health L4 Poll MS: %+v\n", cfg.ParentHealthL4PollMS)
	log.Infof("Parent Health L7 Poll MS: %+v\n", cfg.ParentHealthL7PollMS)
	log.Infof("Parent Health Sv Poll MS: %+v\n", cfg.ParentHealthServicePollMS)

	parentHealthL4PollInterval := time.Duration(cfg.ParentHealthL4PollMS) * time.Millisecond
	parentHealthL7PollInterval := time.Duration(cfg.ParentHealthL7PollMS) * time.Millisecond
	parentHealthServicePollInterval := time.Duration(cfg.ParentHealthServicePollMS) * time.Millisecond
	markdownMinInterval := time.Duration(*cfg.MarkdownMinIntervalMS) * time.Millisecond

	markdownSvc := tmagent.StartMarkdownService(tmInfo, markdownMinInterval)

	parentHealthL4DoneChan := make(chan<- struct{}, 1)
	if _, use := tmInfo.HealthMethods[config.HealthMethodParentL4]; use {
		parentHealthL4DoneChan, err = tmagent.StartParentHealthPoller(tmInfo, cfg.NumHealthWorkers, parentHealthL4PollInterval, tmagent.ParentHealthPollTypeL4, markdownSvc.UpdateHealth)
		if err != nil {
			log.Errorln("starting parent health poller l4 failed: " + err.Error())
			os.Exit(RunTimeError)
		}
	}

	parentHealthL7DoneChan := make(chan<- struct{}, 1)
	if _, use := tmInfo.HealthMethods[config.HealthMethodParentL7]; use {
		parentHealthL7DoneChan, err = tmagent.StartParentHealthPoller(tmInfo, cfg.NumHealthWorkers, parentHealthL7PollInterval, tmagent.ParentHealthPollTypeL7, markdownSvc.UpdateHealth)
		if err != nil {
			log.Errorln("starting parent health poller l7 failed: " + err.Error())
			os.Exit(RunTimeError)
		}
	}

	tmHealthDoneChan := make(chan<- struct{}, 1)
	if _, use := tmInfo.HealthMethods[config.HealthMethodTrafficMonitor]; use {
		tmHealthDoneChan = tmagent.StartTrafficMonitorHealthPoller(tmInfo, markdownSvc.UpdateHealth)
	}

	parentServiceHealthDoneChan := make(chan<- struct{}, 1)
	if _, use := tmInfo.HealthMethods[config.HealthMethodParentService]; use {
		parentServiceHealthDoneChan = tmagent.StartParentServiceHealthPoller(tmInfo, cfg.NumHealthWorkers, parentHealthServicePollInterval, markdownSvc.UpdateHealth)
	}

	const parentHealthServiceTimeout = time.Second * 30 // TODO make configurable

	parentHealthSvc := &tmagent.ParentHealthServer{}
	if cfg.ParentHealthServicePort > 0 {
		parentHealthSvc = tmagent.StartParentHealthService(tmInfo, cfg.ParentHealthServicePort, parentHealthServiceTimeout)
	}

	terminate := func() {
		log.Infoln("Received SIGTERM signal, terminating")
		parentHealthSvc.Shutdown(context.Background()) // TODO timeout context?
		parentHealthL4DoneChan <- struct{}{}
		parentHealthL7DoneChan <- struct{}{}
		parentServiceHealthDoneChan <- struct{}{}
		tmHealthDoneChan <- struct{}{}
		markdownSvc.Shutdown <- struct{}{}
		if err := tmInfo.ParentHealthLog.Close(); err != nil {
			log.Errorln("Closing Parent Health Log: " + err.Error())
		}
		os.Exit(0)
	}

	go signalReloader(unix.SIGHUP, func() { reloadConfig(tmInfo) })
	signalReloader(unix.SIGTERM, terminate)
}

func signalReloader(sig os.Signal, f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig)
	for range c {
		f()
	}
}

// reloadConfig reloads the config if necessary, and puts the new config in cfgGetter.
func reloadConfig(pi *tmagent.ParentInfo) {
	cfg := pi.Cfg.Get()

	newCfg := &config.Cfg{
		HealthClientConfigFile: cfg.HealthClientConfigFile,
		MonitorStrategiesPeers: true,
	}

	isNew, err := config.LoadConfig(newCfg)
	if err != nil {
		log.Errorf("error reading changed config file %s: %s\n", cfg.HealthClientConfigFile.Filename, err.Error())
		return
	}

	if isNew {
		if err = config.ReadCredentials(newCfg, false); err != nil {
			log.Errorf("could not load credentials for config updates, keeping the old config: %v", err.Error())
			return
		}
		if err = pi.GetTOData(newCfg); err != nil {
			log.Errorf("could not update the list of trafficmonitors, keeping the old config: %v", err.Error())
		} else {
			// TODO this was calling a custom copy func that wasn't copying:
			//      MarkUpPollThreshold, TmProxyURL, CredentialFile, ParsedProxyURL
			//      Verify that was a bug, and none of those need to not be updated
			pi.Cfg.Set(newCfg)
			log.Infoln("the configuration has been successfully updated")
		}
		return
	}

	// check for updates to the credentials file
	if cfg.CredentialFile.Filename != "" {
		modTime, err := util.GetFileModificationTime(cfg.CredentialFile.Filename)
		if err != nil {
			log.Errorf("could not stat the credential file '" + cfg.CredentialFile.Filename + "', considering modified! Error: " + err.Error())
			// There's no good option here. If we fail to get the modify time,
			// we can either end up never updating, or update every time.
			// Since this is only ever called on a user-initiated HUP, it's probably ok
			// to always refresh.
			modTime = cfg.CredentialFile.LastModifyTime + 1 // fake it as modified, so we update
		}
		if modTime > cfg.CredentialFile.LastModifyTime {
			log.Infoln("the credential file has changed, loading new credentials")
			newCfg := cfg.Clone()
			if err = config.ReadCredentials(newCfg, true); err != nil {
				log.Errorf("could not load credentials from the updated credential file: %s", newCfg.CredentialFile.Filename)
			} else {
				pi.Cfg.Set(newCfg)
			}
		}
	}
}
