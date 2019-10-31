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
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"

	"golang.org/x/sys/unix"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/config"
	"github.com/apache/trafficcontrol/traffic_monitor/handler"
	"github.com/apache/trafficcontrol/traffic_monitor/health"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/poller"
	"github.com/apache/trafficcontrol/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
	"github.com/apache/trafficcontrol/traffic_monitor/towrap"
)

//
// Start starts the poller and handler goroutines
//
func Start(opsConfigFile string, cfg config.Config, appData config.StaticAppData, trafficMonitorConfigFileName string) error {
	toSession := towrap.ITrafficOpsSession(towrap.NewTrafficOpsSessionThreadsafe(nil, cfg.CRConfigHistoryCount, cfg))

	url6Regex := regexp.MustCompile(`/\d+`)

	localStates := peer.NewCRStatesThreadsafe() // this is the local state as discoverer by this traffic_monitor
	fetchCount := threadsafe.NewUint()          // note this is the number of individual caches fetched from, not the number of times all the caches were polled.
	healthIteration := threadsafe.NewUint()
	errorCount := threadsafe.NewUint()

	toData := todata.NewThreadsafe()

	cacheHealthHandler := cache.NewHandler()
	cacheHealthPoller := poller.NewCache(cfg.CacheHealthPollingInterval, true, cacheHealthHandler, cfg, appData, cfg.CachePollingProtocol)
	cacheStatHandler := cache.NewPrecomputeHandler(toData)
	cacheStatPoller := poller.NewCache(cfg.CacheStatPollingInterval, false, cacheStatHandler, cfg, appData, cfg.CachePollingProtocol)
	monitorConfigPoller := poller.NewMonitorConfig(cfg.MonitorConfigPollingInterval)
	peerHandler := peer.NewHandler()
	peerPoller := poller.NewCache(cfg.PeerPollingInterval, false, peerHandler, cfg, appData, cfg.PeerPollingProtocol)

	go monitorConfigPoller.Poll()
	go cacheHealthPoller.Poll()
	go cacheStatPoller.Poll()
	go peerPoller.Poll()

	events := health.NewThreadsafeEvents(cfg.MaxEvents)

	cachesChanged := make(chan struct{})
	peerStates := peer.NewCRStatesPeersThreadsafe(cfg.PeerOptimisticQuorumMin) // each peer's last state is saved in this map

	monitorConfig := StartMonitorConfigManager(
		monitorConfigPoller.ConfigChannel,
		localStates,
		peerStates,
		cacheStatPoller.ConfigChannel,
		cacheHealthPoller.ConfigChannel,
		peerPoller.ConfigChannel,
		monitorConfigPoller.IntervalChan,
		cachesChanged,
		cfg,
		appData,
		toSession,
		toData,
		url6Regex,
	)

	combinedStates, combineStateFunc := StartStateCombiner(events, peerStates, localStates, toData)

	StartPeerManager(
		peerHandler.ResultChannel,
		peerStates,
		events,
		combineStateFunc,
	)

	statInfoHistory, statResultHistory, statMaxKbpses, _, lastKbpsStats, dsStats, unpolledCaches, localCacheStatus := StartStatHistoryManager(
		cacheStatHandler.ResultChan(),
		localStates,
		combinedStates,
		toData,
		cachesChanged,
		errorCount,
		cfg,
		monitorConfig,
		events,
		combineStateFunc,
	)

	lastHealthDurations, healthHistory := StartHealthResultManager(
		cacheHealthHandler.ResultChan(),
		toData,
		localStates,
		monitorConfig,
		combinedStates,
		fetchCount,
		errorCount,
		cfg,
		events,
		localCacheStatus,
	)

	StartOpsConfigManager(
		opsConfigFile,
		toSession,
		toData,
		[]chan<- handler.OpsConfig{monitorConfigPoller.OpsConfigChannel},
		[]chan<- towrap.ITrafficOpsSession{monitorConfigPoller.SessionChannel},
		localStates,
		peerStates,
		combinedStates,
		statInfoHistory,
		statResultHistory,
		statMaxKbpses,
		healthHistory,
		lastKbpsStats,
		dsStats,
		events,
		appData,
		cacheHealthPoller.Config.Interval,
		lastHealthDurations,
		fetchCount,
		healthIteration,
		errorCount,
		localCacheStatus,
		unpolledCaches,
		monitorConfig,
		cfg,
	)

	if err := startMonitorConfigFilePoller(trafficMonitorConfigFileName); err != nil {
		return fmt.Errorf("starting monitor config file poller: %v", err)
	}

	healthTickListener(cacheHealthPoller.TickChan, healthIteration)
	return nil
}

// healthTickListener listens for health ticks, and writes to the health iteration variable. Does not return.
func healthTickListener(cacheHealthTick <-chan uint64, healthIteration threadsafe.Uint) {
	for i := range cacheHealthTick {
		healthIteration.Set(i)
	}
}

func startMonitorConfigFilePoller(filename string) error {
	onChange := func(bytes []byte, err error) {
		if err != nil {
			log.Errorf("monitor config file poll, polling file '%v': %v", filename, err)
			return
		}
		cfg, err := config.LoadBytes(bytes)
		if err != nil {
			log.Errorf("monitor config file poll, loading bytes '%v' from '%v': %v", string(bytes), filename, err)
			return
		}
		if err := log.InitCfg(cfg); err != nil {
			log.Errorf("monitor config file poll, getting log writers '%v': %v", filename, err)
			return
		}
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	onChange(bytes, nil)

	startSignalFileReloader(filename, unix.SIGHUP, onChange)
	return nil
}

// signalFileReloader starts a goroutine which, when the given signal is received, attempts to load the given file and calls the given function with its bytes or error. There is no way to stop the goroutine or stop listening for signals, thus this should not be called if it's ever necessary to stop handling or change the listened file. The initialRead parameter determines whether the given handler is called immediately with an attempted file read (without a signal).
func startSignalFileReloader(filename string, sig os.Signal, f func([]byte, error)) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, sig)
		for range c {
			f(ioutil.ReadFile(filename))
		}
	}()
}
