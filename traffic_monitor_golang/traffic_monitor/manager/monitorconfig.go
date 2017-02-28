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
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// StartMonitorConfigManager runs the monitor config manager goroutine, and returns the threadsafe data which it sets.
func StartMonitorConfigManager(
	monitorConfigPollChan <-chan to.TrafficMonitorConfigMap,
	localStates peer.CRStatesThreadsafe,
	statURLSubscriber chan<- poller.HttpPollerConfig,
	healthURLSubscriber chan<- poller.HttpPollerConfig,
	peerURLSubscriber chan<- poller.HttpPollerConfig,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData config.StaticAppData,
) threadsafe.TrafficMonitorConfigMap {
	monitorConfig := threadsafe.NewTrafficMonitorConfigMap()
	go monitorConfigListen(monitorConfig,
		monitorConfigPollChan,
		localStates,
		statURLSubscriber,
		healthURLSubscriber,
		peerURLSubscriber,
		cachesChangeSubscriber,
		cfg,
		staticAppData,
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

// PollIntervalRatio is the ratio of the configuration interval to poll. The configured intervals are 'target' times, so we actually poll at some small fraction less, in attempt to make the actual poll marginally less than the target.
const PollIntervalRatio = float64(0.97) // TODO make config?

// getPollIntervals reads the Traffic Ops Client monitorConfig structure, and parses and returns the health, peer, and stat poll intervals
func getHealthPeerStatPollIntervals(monitorConfig to.TrafficMonitorConfigMap, cfg config.Config, logMissingHeartbeatParam bool) (time.Duration, time.Duration, time.Duration, error) {
	peerPollIntervalI, peerPollIntervalExists := monitorConfig.Config["peers.polling.interval"]
	if !peerPollIntervalExists {
		return 0, 0, 0, fmt.Errorf("Traffic Ops Monitor config missing 'peers.polling.interval', not setting config changes.\n")
	}
	peerPollIntervalInt, peerPollIntervalIsInt := peerPollIntervalI.(float64)
	if !peerPollIntervalIsInt {
		return 0, 0, 0, fmt.Errorf("Traffic Ops Monitor config 'peers.polling.interval' value '%v' type %T is not an integer, not setting config changes.\n", peerPollIntervalI, peerPollIntervalI)
	}
	peerPollInterval := trafficOpsPeerPollIntervalToDuration(int(peerPollIntervalInt))

	statPollIntervalI, statPollIntervalExists := monitorConfig.Config["health.polling.interval"]
	if !statPollIntervalExists {
		return 0, 0, 0, fmt.Errorf("Traffic Ops Monitor config missing 'health.polling.interval', not setting config changes.\n")
	}
	statPollIntervalInt, statPollIntervalIsInt := statPollIntervalI.(float64)
	if !statPollIntervalIsInt {
		return 0, 0, 0, fmt.Errorf("Traffic Ops Monitor config 'health.polling.interval' value '%v' type %T is not an integer, not setting config changes.\n", statPollIntervalI, statPollIntervalI)
	}
	statPollInterval := trafficOpsStatPollIntervalToDuration(int(statPollIntervalInt))

	healthPollIntervalI, healthPollIntervalExists := monitorConfig.Config["heartbeat.polling.interval"]
	healthPollIntervalInt, healthPollIntervalIsInt := healthPollIntervalI.(float64)
	if !healthPollIntervalExists {
		if logMissingHeartbeatParam {
			log.Warnln("Traffic Ops Monitor config missing 'heartbeat.polling.interval', using health for heartbeat.")
		}
		healthPollIntervalInt = statPollIntervalInt
	} else if !healthPollIntervalIsInt {
		log.Warnf("Traffic Ops Monitor config 'heartbeat.polling.interval' value '%v' type %T is not an integer, using health for heartbeat\n", statPollIntervalI, statPollIntervalI)
		healthPollIntervalInt = statPollIntervalInt
	}
	healthPollInterval := trafficOpsHealthPollIntervalToDuration(int(healthPollIntervalInt))

	healthPollInterval = time.Duration(float64(healthPollInterval) * PollIntervalRatio)
	peerPollInterval = time.Duration(float64(peerPollInterval) * PollIntervalRatio)
	statPollInterval = time.Duration(float64(statPollInterval) * PollIntervalRatio)
	return healthPollInterval, peerPollInterval, statPollInterval, nil
}

// TODO timing, and determine if the case, or its internal `for`, should be put in a goroutine
// TODO determine if subscribers take action on change, and change to mutexed objects if not.
func monitorConfigListen(
	monitorConfigTS threadsafe.TrafficMonitorConfigMap,
	monitorConfigPollChan <-chan to.TrafficMonitorConfigMap,
	localStates peer.CRStatesThreadsafe,
	statURLSubscriber chan<- poller.HttpPollerConfig,
	healthURLSubscriber chan<- poller.HttpPollerConfig,
	peerURLSubscriber chan<- poller.HttpPollerConfig,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData config.StaticAppData,
) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("MonitorConfigManager panic: %v\n", err)
		} else {
			log.Errorf("MonitorConfigManager failed without panic\n")
		}
		os.Exit(1) // The Monitor can't run without a MonitorConfigManager
	}()

	logMissingHeartbeatParam := true

	for monitorConfig := range monitorConfigPollChan {
		monitorConfigTS.Set(monitorConfig)
		healthURLs := map[string]poller.PollConfig{}
		statURLs := map[string]poller.PollConfig{}
		peerURLs := map[string]poller.PollConfig{}
		caches := map[string]string{}

		healthPollInterval, peerPollInterval, statPollInterval, err := getHealthPeerStatPollIntervals(monitorConfig, cfg, logMissingHeartbeatParam)
		logMissingHeartbeatParam = false // only log the heartbeat parameter missing once

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
		}

		statURLSubscriber <- poller.HttpPollerConfig{Urls: statURLs, Interval: statPollInterval}
		healthURLSubscriber <- poller.HttpPollerConfig{Urls: healthURLs, Interval: healthPollInterval}
		peerURLSubscriber <- poller.HttpPollerConfig{Urls: peerURLs, Interval: peerPollInterval}

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
