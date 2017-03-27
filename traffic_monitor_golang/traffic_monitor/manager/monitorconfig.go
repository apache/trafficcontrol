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
	"os"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/poller"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/config"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
	towrap "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopswrapper"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

type PollIntervals struct {
	Health time.Duration
	Peer   time.Duration
	Stat   time.Duration
	TO     time.Duration
}

// getPollIntervals reads the Traffic Ops Client monitorConfig structure, and parses and returns the health, peer, stat, and TrafficOps poll intervals
func getIntervals(monitorConfig to.TrafficMonitorConfigMap, cfg config.Config, logMissingParams bool) (PollIntervals, error) {
	intervals := PollIntervals{}
	peerPollIntervalI, peerPollIntervalExists := monitorConfig.Config["peers.polling.interval"]
	if !peerPollIntervalExists {
		return PollIntervals{}, fmt.Errorf("Traffic Ops Monitor config missing 'peers.polling.interval', not setting config changes.\n")
	}
	peerPollIntervalInt, peerPollIntervalIsInt := peerPollIntervalI.(float64)
	if !peerPollIntervalIsInt {
		return PollIntervals{}, fmt.Errorf("Traffic Ops Monitor config 'peers.polling.interval' value '%v' type %T is not an integer, not setting config changes.\n", peerPollIntervalI, peerPollIntervalI)
	}
	intervals.Peer = trafficOpsPeerPollIntervalToDuration(int(peerPollIntervalInt))

	statPollIntervalI, statPollIntervalExists := monitorConfig.Config["health.polling.interval"]
	if !statPollIntervalExists {
		return PollIntervals{}, fmt.Errorf("Traffic Ops Monitor config missing 'health.polling.interval', not setting config changes.\n")
	}
	statPollIntervalInt, statPollIntervalIsInt := statPollIntervalI.(float64)
	if !statPollIntervalIsInt {
		return PollIntervals{}, fmt.Errorf("Traffic Ops Monitor config 'health.polling.interval' value '%v' type %T is not an integer, not setting config changes.\n", statPollIntervalI, statPollIntervalI)
	}
	intervals.Stat = trafficOpsStatPollIntervalToDuration(int(statPollIntervalInt))

	healthPollIntervalI, healthPollIntervalExists := monitorConfig.Config["heartbeat.polling.interval"]
	healthPollIntervalInt, healthPollIntervalIsInt := healthPollIntervalI.(float64)
	if !healthPollIntervalExists {
		if logMissingParams {
			log.Warnln("Traffic Ops Monitor config missing 'heartbeat.polling.interval', using health for heartbeat.")
		}
		healthPollIntervalInt = statPollIntervalInt
	} else if !healthPollIntervalIsInt {
		log.Warnf("Traffic Ops Monitor config 'heartbeat.polling.interval' value '%v' type %T is not an integer, using health for heartbeat\n", statPollIntervalI, statPollIntervalI)
		healthPollIntervalInt = statPollIntervalInt
	}
	intervals.Health = trafficOpsHealthPollIntervalToDuration(int(healthPollIntervalInt))

	toPollIntervalI, toPollIntervalExists := monitorConfig.Config["tm.polling.interval"]
	toPollIntervalInt, toPollIntervalIsInt := toPollIntervalI.(float64)
	intervals.TO = cfg.MonitorConfigPollingInterval
	if !toPollIntervalExists {
		if logMissingParams {
			log.Warnf("Traffic Ops Monitor config missing 'tm.polling.interval', using config value '%v'\n", cfg.MonitorConfigPollingInterval)
		}
	} else if !toPollIntervalIsInt {
		log.Warnf("Traffic Ops Monitor config 'tm.polling.interval' value '%v' type %T is not an integer, using config value '%v'\n", toPollIntervalI, toPollIntervalI, cfg.MonitorConfigPollingInterval)
	} else {
		intervals.TO = trafficOpsTOPollIntervalToDuration(int(toPollIntervalInt))
	}

	multiplyByRatio := func(i time.Duration) time.Duration {
		return time.Duration(float64(i) * PollIntervalRatio)
	}

	intervals.TO = multiplyByRatio(intervals.TO)
	intervals.Health = multiplyByRatio(intervals.Health)
	intervals.Peer = multiplyByRatio(intervals.Peer)
	intervals.Stat = multiplyByRatio(intervals.Stat)
	return intervals, nil
}

// StartMonitorConfigManager runs the monitor config manager goroutine, and returns the threadsafe data which it sets.
func StartMonitorConfigManager(
	monitorConfigPollChan <-chan poller.MonitorCfg,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	statURLSubscriber chan<- poller.HttpPollerConfig,
	healthURLSubscriber chan<- poller.HttpPollerConfig,
	peerURLSubscriber chan<- poller.HttpPollerConfig,
	toIntervalSubscriber chan<- time.Duration,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData config.StaticAppData,
	toSession towrap.ITrafficOpsSession,
	toData todata.TODataThreadsafe,
) threadsafe.TrafficMonitorConfigMap {
	monitorConfig := threadsafe.NewTrafficMonitorConfigMap()
	go monitorConfigListen(monitorConfig,
		monitorConfigPollChan,
		localStates,
		peerStates,
		statURLSubscriber,
		healthURLSubscriber,
		peerURLSubscriber,
		toIntervalSubscriber,
		cachesChangeSubscriber,
		cfg,
		staticAppData,
		toSession,
		toData,
	)
	return monitorConfig
}

// trafficOpsHealthConnectionTimeoutToDuration takes the int from Traffic Ops, which is in milliseconds, and returns a time.Duration
// TODO change Traffic Ops Client API to a time.Duration
func trafficOpsHealthConnectionTimeoutToDuration(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

// trafficOpsPeerPollIntervalToDuration takes the int from Traffic Ops, which is in milliseconds, and returns a time.Duration
// TODO change Traffic Ops Client API to a time.Duration
func trafficOpsPeerPollIntervalToDuration(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

// trafficOpsStatPollIntervalToDuration takes the int from Traffic Ops, which is in milliseconds, and returns a time.Duration
// TODO change Traffic Ops Client API to a time.Duration
func trafficOpsStatPollIntervalToDuration(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

// trafficOpsHealthPollIntervalToDuration takes the int from Traffic Ops, which is in milliseconds, and returns a time.Duration
// TODO change Traffic Ops Client API to a time.Duration
func trafficOpsHealthPollIntervalToDuration(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

// trafficOpsTOPollIntervalToDuration takes the int from Traffic Ops, which is in milliseconds, and returns a time.Duration
// TODO change Traffic Ops Client API to a time.Duration
func trafficOpsTOPollIntervalToDuration(t int) time.Duration {
	return time.Duration(t) * time.Millisecond
}

// PollIntervalRatio is the ratio of the configuration interval to poll. The configured intervals are 'target' times, so we actually poll at some small fraction less, in attempt to make the actual poll marginally less than the target.
const PollIntervalRatio = float64(0.97) // TODO make config?

// TODO timing, and determine if the case, or its internal `for`, should be put in a goroutine
// TODO determine if subscribers take action on change, and change to mutexed objects if not.
func monitorConfigListen(
	monitorConfigTS threadsafe.TrafficMonitorConfigMap,
	monitorConfigPollChan <-chan poller.MonitorCfg,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	statURLSubscriber chan<- poller.HttpPollerConfig,
	healthURLSubscriber chan<- poller.HttpPollerConfig,
	peerURLSubscriber chan<- poller.HttpPollerConfig,
	toIntervalSubscriber chan<- time.Duration,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData config.StaticAppData,
	toSession towrap.ITrafficOpsSession,
	toData todata.TODataThreadsafe,
) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("MonitorConfigManager panic: %v\n", err)
		} else {
			log.Errorf("MonitorConfigManager failed without panic\n")
		}
		os.Exit(1) // The Monitor can't run without a MonitorConfigManager
	}()

	logMissingIntervalParams := true

	for pollerMonitorCfg := range monitorConfigPollChan {
		monitorConfig := pollerMonitorCfg.Cfg
		cdn := pollerMonitorCfg.CDN
		monitorConfigTS.Set(monitorConfig)
		toData.Update(toSession, cdn)

		healthURLs := map[string]poller.PollConfig{}
		statURLs := map[string]poller.PollConfig{}
		peerURLs := map[string]poller.PollConfig{}
		caches := map[string]string{}

		intervals, err := getIntervals(monitorConfig, cfg, logMissingIntervalParams)
		logMissingIntervalParams = false // only log missing parameters once
		if err != nil {
			log.Errorf("monitor config error getting polling intervals, can't poll: %v", err)
			continue
		}

		for _, srv := range monitorConfig.TrafficServer {
			caches[srv.HostName] = srv.Status

			cacheName := enum.CacheName(srv.HostName)

			srvStatus := enum.CacheStatusFromString(srv.Status)
			if srvStatus == enum.CacheStatusOnline {
				localStates.AddCache(cacheName, peer.IsAvailable{IsAvailable: true})
				continue
			}
			if srvStatus == enum.CacheStatusOffline {
				continue
			}
			// seed states with available = false until our polling cycle picks up a result
			if _, exists := localStates.GetCache(cacheName); !exists {
				localStates.AddCache(cacheName, peer.IsAvailable{IsAvailable: false})
			}

			url := monitorConfig.Profile[srv.Profile].Parameters.HealthPollingURL
			if url == "" {
				log.Errorf("monitor config server %v profile %v has no polling URL; can't poll", srv.HostName, srv.Profile)
				continue
			}
			r := strings.NewReplacer(
				"${hostname}", srv.IP,
				"${interface_name}", srv.InterfaceName,
				"application=plugin.remap", "application=system",
				"application=", "application=system",
			)
			url = r.Replace(url)

			connTimeout := trafficOpsHealthConnectionTimeoutToDuration(monitorConfig.Profile[srv.Profile].Parameters.HealthConnectionTimeout)
			healthURLs[srv.HostName] = poller.PollConfig{URL: url, Host: srv.FQDN, Timeout: connTimeout}
			r = strings.NewReplacer("application=system", "application=")
			statURL := r.Replace(url)
			statURLs[srv.HostName] = poller.PollConfig{URL: statURL, Host: srv.FQDN, Timeout: connTimeout}
		}

		peerSet := map[enum.TrafficMonitorName]struct{}{}
		for _, srv := range monitorConfig.TrafficMonitor {
			if srv.HostName == staticAppData.Hostname {
				continue
			}
			if enum.CacheStatusFromString(srv.Status) != enum.CacheStatusOnline {
				continue
			}
			// TODO: the URL should be config driven. -jse
			url := fmt.Sprintf("http://%s:%d/publish/CrStates?raw", srv.IP, srv.Port)
			peerURLs[srv.HostName] = poller.PollConfig{URL: url, Host: srv.FQDN} // TODO determine timeout.
			peerSet[enum.TrafficMonitorName(srv.HostName)] = struct{}{}
		}

		statURLSubscriber <- poller.HttpPollerConfig{Urls: statURLs, Interval: intervals.Stat}
		healthURLSubscriber <- poller.HttpPollerConfig{Urls: healthURLs, Interval: intervals.Health}
		peerURLSubscriber <- poller.HttpPollerConfig{Urls: peerURLs, Interval: intervals.Peer}
		toIntervalSubscriber <- intervals.TO
		peerStates.SetTimeout((intervals.Peer + cfg.HTTPTimeout) * 2)
		peerStates.SetPeers(peerSet)

		for cacheName := range localStates.GetCaches() {
			if _, exists := monitorConfig.TrafficServer[string(cacheName)]; !exists {
				log.Warnf("Removing %s from localStates", cacheName)
				localStates.DeleteCache(cacheName)
			}
		}

		if len(healthURLs) == 0 {
			log.Errorf("No REPORTED caches exist in Traffic Ops, nothing to poll.")
		}

		cachesChangeSubscriber <- struct{}{}

		// TODO because there are multiple writers to localStates.DeliveryService, there is a race condition, where MonitorConfig (this func) and HealthResultManager could write at the same time, and the HealthResultManager could overwrite a delivery service addition or deletion here. Probably the simplest and most performant fix would be a lock-free algorithm using atomic compare-and-swaps.
		for _, ds := range monitorConfig.DeliveryService {
			// since caches default to unavailable, also default DS false
			if _, exists := localStates.GetDeliveryService(enum.DeliveryServiceName(ds.XMLID)); !exists {
				localStates.SetDeliveryService(enum.DeliveryServiceName(ds.XMLID), peer.Deliveryservice{IsAvailable: false, DisabledLocations: []enum.CacheGroupName{}}) // important to initialize DisabledLocations, so JSON is `[]` not `null`
			}
		}
		for ds := range localStates.GetDeliveryServices() {
			if _, exists := monitorConfig.DeliveryService[string(ds)]; !exists {
				localStates.DeleteDeliveryService(ds)
			}
		}
	}
}
