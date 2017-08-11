package manager

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
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/sys/unix"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/handler"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/config"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/datareq"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/health"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/srvhttp"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// StartOpsConfigManager starts the ops config manager goroutine, returning the (threadsafe) variables which it sets.
// Note the OpsConfigManager is in charge of the httpServer, because ops config changes trigger server changes. If other things needed to trigger server restarts, the server could be put in its own goroutine with signal channels
func StartOpsConfigManager(
	opsConfigFile string,
	toSession towrap.ITrafficOpsSession,
	toData todata.TODataThreadsafe,
	opsConfigChangeSubscribers []chan<- handler.OpsConfig,
	toChangeSubscribers []chan<- towrap.ITrafficOpsSession,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	statInfoHistory threadsafe.ResultInfoHistory,
	statResultHistory threadsafe.ResultStatHistory,
	statMaxKbpses threadsafe.CacheKbpses,
	healthHistory threadsafe.ResultHistory,
	lastStats threadsafe.LastStats,
	dsStats threadsafe.DSStatsReader,
	events health.ThreadsafeEvents,
	staticAppData config.StaticAppData,
	healthPollInterval time.Duration,
	lastHealthDurations threadsafe.DurationMap,
	fetchCount threadsafe.Uint,
	healthIteration threadsafe.Uint,
	errorCount threadsafe.Uint,
	localCacheStatus threadsafe.CacheAvailableStatus,
	unpolledCaches threadsafe.UnpolledCaches,
	monitorConfig threadsafe.TrafficMonitorConfigMap,
	cfg config.Config,
) (threadsafe.OpsConfig, error) {

	handleErr := func(err error) {
		errorCount.Inc()
		log.Errorf("OpsConfigManager: %v\n", err)
	}

	httpServer := srvhttp.Server{}
	opsConfig := threadsafe.NewOpsConfig()

	// TODO remove change subscribers, give Threadsafes directly to the things that need them. If they only set vars, and don't actually do work on change.
	onChange := func(bytes []byte, err error) {
		if err != nil {
			handleErr(err)
			return
		}

		newOpsConfig := handler.OpsConfig{}
		if err = json.Unmarshal(bytes, &newOpsConfig); err != nil {
			handleErr(fmt.Errorf("Could not unmarshal Ops Config JSON: %s\n", err))
			return
		}

		opsConfig.Set(newOpsConfig)

		listenAddress := ":80" // default

		if newOpsConfig.HttpListener != "" {
			listenAddress = newOpsConfig.HttpListener
		}

		endpoints := datareq.MakeDispatchMap(
			opsConfig,
			toSession,
			localStates,
			peerStates,
			combinedStates,
			statInfoHistory,
			statResultHistory,
			statMaxKbpses,
			healthHistory,
			dsStats,
			events,
			staticAppData,
			healthPollInterval,
			lastHealthDurations,
			fetchCount,
			healthIteration,
			errorCount,
			toData,
			localCacheStatus,
			lastStats,
			unpolledCaches,
			monitorConfig,
		)
		err = httpServer.Run(endpoints, listenAddress, cfg.ServeReadTimeout, cfg.ServeWriteTimeout, cfg.StaticFileDir)
		if err != nil {
			handleErr(fmt.Errorf("MonitorConfigPoller: error creating HTTP server: %s\n", err))
			return
		}

		// TODO config? parameter?
		useCache := false
		trafficOpsRequestTimeout := time.Second * time.Duration(10)

		realToSession, err := to.LoginWithAgent(newOpsConfig.Url, newOpsConfig.Username, newOpsConfig.Password, newOpsConfig.Insecure, staticAppData.UserAgent, useCache, trafficOpsRequestTimeout)
		if err != nil {
			handleErr(fmt.Errorf("MonitorConfigPoller: error instantiating Session with traffic_ops: %s\n", err))
			return
		}
		toSession.Set(realToSession)

		if cdn, err := getMonitorCDN(realToSession, staticAppData.Hostname); err != nil {
			handleErr(fmt.Errorf("getting CDN name from Traffic Ops, using config CDN '%s': %s\n", newOpsConfig.CdnName, err))
		} else {
			if newOpsConfig.CdnName != "" && newOpsConfig.CdnName != cdn {
				log.Warnf("%s Traffic Ops CDN '%s' doesn't match config CDN '%s' - using Traffic Ops CDN\n", staticAppData.Hostname, cdn, newOpsConfig.CdnName)
			}
			newOpsConfig.CdnName = cdn
		}

		if err := toData.Fetch(toSession, newOpsConfig.CdnName); err != nil {
			handleErr(fmt.Errorf("Error getting Traffic Ops data: %v\n", err))
			return
		}

		// These must be in a goroutine, because the monitorConfigPoller tick sends to a channel this select listens for. Thus, if we block on sends to the monitorConfigPoller, we have a livelock race condition.
		// More generically, we're using goroutines as an infinite chan buffer, to avoid potential livelocks
		for _, subscriber := range opsConfigChangeSubscribers {
			go func(s chan<- handler.OpsConfig) { s <- newOpsConfig }(subscriber)
		}
		for _, subscriber := range toChangeSubscribers {
			go func(s chan<- towrap.ITrafficOpsSession) { s <- toSession }(subscriber)
		}
	}

	bytes, err := ioutil.ReadFile(opsConfigFile)
	if err != nil {
		return opsConfig, err
	}
	onChange(bytes, err)

	startSignalFileReloader(opsConfigFile, unix.SIGHUP, onChange)

	return opsConfig, nil
}

// getMonitorCDN returns the CDN of a given Traffic Monitor.
// TODO change to get by name, when Traffic Ops supports querying a single server.
func getMonitorCDN(toc *to.Session, monitorHostname string) (string, error) {
	servers, err := toc.Servers()
	if err != nil {
		return "", fmt.Errorf("getting monitor %s CDN: %v", monitorHostname, err)
	}

	for _, server := range servers {
		if server.HostName != monitorHostname {
			continue
		}
		return server.CDNName, nil
	}
	return "", fmt.Errorf("no monitor named %v found in Traffic Ops", monitorHostname)
}
