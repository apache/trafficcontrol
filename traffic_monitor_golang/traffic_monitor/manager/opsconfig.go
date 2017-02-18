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
	"fmt"
	"sync"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/handler"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/poller"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/config"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/health"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/srvhttp"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// OpsConfigThreadsafe provides safe access for multiple reader goroutines and a single writer to a stored OpsConfig object.
// This could be made lock-free, if the performance was necessary
type OpsConfigThreadsafe struct {
	opsConfig *handler.OpsConfig
	m         *sync.RWMutex
}

// NewOpsConfigThreadsafe returns a new single-writer-multiple-reader OpsConfig
func NewOpsConfigThreadsafe() OpsConfigThreadsafe {
	return OpsConfigThreadsafe{m: &sync.RWMutex{}, opsConfig: &handler.OpsConfig{}}
}

// Get gets the internal OpsConfig object. This MUST NOT be modified. If modification is necessary, copy the object.
func (o *OpsConfigThreadsafe) Get() handler.OpsConfig {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.opsConfig
}

// Set sets the internal OpsConfig object. This MUST NOT be called from multiple goroutines.
func (o *OpsConfigThreadsafe) Set(newOpsConfig handler.OpsConfig) {
	o.m.Lock()
	*o.opsConfig = newOpsConfig
	o.m.Unlock()
}

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
	staticAppData StaticAppData,
	healthPollInterval time.Duration,
	lastHealthDurations DurationMapThreadsafe,
	fetchCount threadsafe.Uint,
	healthIteration threadsafe.Uint,
	errorCount threadsafe.Uint,
	localCacheStatus threadsafe.CacheAvailableStatus,
	unpolledCaches threadsafe.UnpolledCaches,
	monitorConfig TrafficMonitorConfigMapThreadsafe,
	cfg config.Config,
) OpsConfigThreadsafe {

	opsConfigFileChannel := make(chan interface{})
	opsConfigFilePoller := poller.FilePoller{
		File:          opsConfigFile,
		ResultChannel: opsConfigFileChannel,
	}

	opsConfigChannel := make(chan handler.OpsConfig)
	opsConfigFileHandler := handler.OpsConfigFileHandler{
		ResultChannel:    opsConfigFilePoller.ResultChannel,
		OpsConfigChannel: opsConfigChannel,
	}

	go opsConfigFileHandler.Listen()
	go opsConfigFilePoller.Poll()

	opsConfig := NewOpsConfigThreadsafe()

	// TODO remove change subscribers, give Threadsafes directly to the things that need them. If they only set vars, and don't actually do work on change.
	go func() {
		httpServer := srvhttp.Server{}

		for newOpsConfig := range opsConfigChannel {
			var err error
			opsConfig.Set(newOpsConfig)

			listenAddress := ":80" // default

			if newOpsConfig.HttpListener != "" {
				listenAddress = newOpsConfig.HttpListener
			}

			handleErr := func(err error) {
				errorCount.Inc()
				log.Errorf("OpsConfigManager: %v\n", err)
			}

			endpoints := MakeDispatchMap(
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
				continue
			}

			// TODO config? parameter?
			useCache := false
			trafficOpsRequestTimeout := time.Second * time.Duration(10)

			realToSession, err := to.LoginWithAgent(newOpsConfig.Url, newOpsConfig.Username, newOpsConfig.Password, newOpsConfig.Insecure, staticAppData.UserAgent, useCache, trafficOpsRequestTimeout)
			if err != nil {
				handleErr(fmt.Errorf("MonitorConfigPoller: error instantiating Session with traffic_ops: %s\n", err))
				continue
			}
			toSession.Set(realToSession)

			if err := toData.Fetch(toSession, newOpsConfig.CdnName); err != nil {
				handleErr(fmt.Errorf("Error getting Traffic Ops data: %v\n", err))
				continue
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
	}()

	return opsConfig
}
