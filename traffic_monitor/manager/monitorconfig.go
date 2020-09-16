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
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/traffic_monitor/config"
	"github.com/apache/trafficcontrol/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/traffic_monitor/poller"
	"github.com/apache/trafficcontrol/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
	"github.com/apache/trafficcontrol/traffic_monitor/towrap"
)

type PollIntervals struct {
	Health            time.Duration
	HealthNoKeepAlive bool
	Peer              time.Duration
	PeerNoKeepAlive   bool
	Stat              time.Duration
	StatNoKeepAlive   bool
	TO                time.Duration
}

// getPollIntervals reads the Traffic Ops Client monitorConfig structure, and parses and returns the health, peer, stat, and TrafficOps poll intervals
func getIntervals(monitorConfig tc.TrafficMonitorConfigMap, cfg config.Config, logMissingParams bool) (PollIntervals, error) {
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

	getNoKeepAlive := func(param string) bool {
		keepAliveI, keepAliveExists := monitorConfig.Config[param]
		keepAliveStr, keepAliveIsStr := keepAliveI.(string)
		return keepAliveExists && keepAliveIsStr && !strings.HasPrefix(strings.ToLower(keepAliveStr), "t")
	}
	intervals.PeerNoKeepAlive = getNoKeepAlive("peer.polling.keepalive")
	intervals.HealthNoKeepAlive = getNoKeepAlive("health.polling.keepalive")
	intervals.StatNoKeepAlive = getNoKeepAlive("stat.polling.keepalive")

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
	statURLSubscriber chan<- poller.CachePollerConfig,
	healthURLSubscriber chan<- poller.CachePollerConfig,
	peerURLSubscriber chan<- poller.CachePollerConfig,
	toIntervalSubscriber chan<- time.Duration,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData config.StaticAppData,
	toSession towrap.TrafficOpsSessionThreadsafe,
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

const DefaultHealthConnectionTimeout = time.Second * 2

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
	statURLSubscriber chan<- poller.CachePollerConfig,
	healthURLSubscriber chan<- poller.CachePollerConfig,
	peerURLSubscriber chan<- poller.CachePollerConfig,
	toIntervalSubscriber chan<- time.Duration,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData config.StaticAppData,
	toSession towrap.TrafficOpsSessionThreadsafe,
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
		if err := toData.Update(toSession, cdn); err != nil {
			log.Errorln("Updating Traffic Ops Data: " + err.Error())
		}

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
			caches[srv.HostName] = srv.ServerStatus

			cacheName := tc.CacheName(srv.HostName)

			srvStatus := tc.CacheStatusFromString(srv.ServerStatus)
			if srvStatus == tc.CacheStatusOnline {
				localStates.AddCache(cacheName, tc.IsAvailable{IsAvailable: true})
				continue
			}
			if srvStatus == tc.CacheStatusOffline {
				continue
			}
			// seed states with available = false until our polling cycle picks up a result
			if _, exists := localStates.GetCache(cacheName); !exists {
				localStates.AddCache(cacheName, tc.IsAvailable{IsAvailable: false})
			}

			pollURLStr := monitorConfig.Profile[srv.Profile].Parameters.HealthPollingURL
			if pollURLStr == "" {
				log.Errorf("monitor config server %v profile %v has no polling URL; can't poll", srv.HostName, srv.Profile)
				continue
			}

			format := monitorConfig.Profile[srv.Profile].Parameters.HealthPollingFormat
			if format == "" {
				format = cache.DefaultStatsType
				log.Infof("health.polling.format for '%v' is empty, using default '%v'", srv.HostName, format)
			}

			pollType := monitorConfig.Profile[srv.Profile].Parameters.HealthPollingType
			if pollType == "" {
				pollType = poller.DefaultPollerType
				log.Infof("health.polling.type for '%v' is empty, using default '%v'", srv.HostName, pollType)
			}

			pollURL4Str, pollURL6Str := createServerHealthPollURLs(pollURLStr, srv)

			connTimeout := trafficOpsHealthConnectionTimeoutToDuration(monitorConfig.Profile[srv.Profile].Parameters.HealthConnectionTimeout)
			if connTimeout == 0 {
				connTimeout = DefaultHealthConnectionTimeout
				log.Warnln("profile " + srv.Profile + " health.connection.timeout Parameter is missing or zero, using default " + DefaultHealthConnectionTimeout.String())
			}

			healthURLs[srv.HostName] = poller.PollConfig{URL: pollURL4Str, URLv6: pollURL6Str, Host: srv.FQDN, Timeout: connTimeout, Format: format, PollType: pollType}

			statURL4 := createServerStatPollURL(pollURL4Str)
			statURL6 := createServerStatPollURL(pollURL6Str)
			statURLs[srv.HostName] = poller.PollConfig{URL: statURL4, URLv6: statURL6, Host: srv.FQDN, Timeout: connTimeout, Format: format, PollType: pollType}
		}

		peerSet := map[tc.TrafficMonitorName]struct{}{}
		for _, srv := range monitorConfig.TrafficMonitor {
			if srv.HostName == staticAppData.Hostname {
				continue
			}
			if tc.CacheStatusFromString(srv.ServerStatus) != tc.CacheStatusOnline {
				continue
			}
			// TODO: the URL should be config driven. -jse
			url4 := fmt.Sprintf("http://%s:%d/publish/CrStates?raw", srv.IP, srv.Port)
			url6 := fmt.Sprintf("http://[%s]:%d/publish/CrStates?raw", ipv6CIDRStrToAddr(srv.IP6), srv.Port)
			peerURLs[srv.HostName] = poller.PollConfig{URL: url4, URLv6: url6, Host: srv.FQDN} // TODO determine timeout.
			peerSet[tc.TrafficMonitorName(srv.HostName)] = struct{}{}
		}

		statURLSubscriber <- poller.CachePollerConfig{Urls: statURLs, PollingProtocol: cfg.CachePollingProtocol, Interval: intervals.Stat, NoKeepAlive: intervals.StatNoKeepAlive}
		healthURLSubscriber <- poller.CachePollerConfig{Urls: healthURLs, PollingProtocol: cfg.CachePollingProtocol, Interval: intervals.Health, NoKeepAlive: intervals.HealthNoKeepAlive}
		peerURLSubscriber <- poller.CachePollerConfig{Urls: peerURLs, PollingProtocol: cfg.PeerPollingProtocol, Interval: intervals.Peer, NoKeepAlive: intervals.PeerNoKeepAlive}
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
			if _, exists := localStates.GetDeliveryService(tc.DeliveryServiceName(ds.XMLID)); !exists {
				localStates.SetDeliveryService(tc.DeliveryServiceName(ds.XMLID), tc.CRStatesDeliveryService{IsAvailable: false, DisabledLocations: []tc.CacheGroupName{}}) // important to initialize DisabledLocations, so JSON is `[]` not `null`
			}
		}
		for ds := range localStates.GetDeliveryServices() {
			if _, exists := monitorConfig.DeliveryService[string(ds)]; !exists {
				localStates.DeleteDeliveryService(ds)
			}
		}
	}
}

// createServerHealthPollURLs takes the template pollingURLStr, and replaces
// variables with data from srv, and returns the polling URL for srv.
//
// Note: `${hostname}` is replaced with the server's service IPv4 address (when
// possible) for IPv4 polls, and its IPv6 service address (when possible) for
// IPv6 polls - NOT the servers hostname!
func createServerHealthPollURLs(pollingURLStr string, srv tc.TrafficServer) (string, string) {
	lid, err := tc.InterfaceInfoToLegacyInterfaces(srv.Interfaces)
	if err != nil {
		log.Errorf("Failed to parse polling strings for cache server '%s': %v", srv.HostName, err)
		return "", ""
	}

	var infName string
	if lid.InterfaceName != nil {
		infName = *lid.InterfaceName
	}

	var pollingURL4Str string
	if lid.IPAddress != nil && *lid.IPAddress != "" {
		pollingURL4Str = strings.NewReplacer(
			"${hostname}", *lid.IPAddress,
			"${interface_name}", infName,
			"application=plugin.remap", "application=system",
			"application=", "application=system",
		).Replace(pollingURLStr)

		pollingURL4Str = insertPorts(pollingURL4Str, srv)
	}

	var pollingURL6Str string
	if lid.IP6Address != nil && *lid.IP6Address != "" {
		r := strings.NewReplacer(
			"${hostname}", "["+ipv6CIDRStrToAddr(*lid.IP6Address)+"]",
			"${interface_name}", infName,
			"application=plugin.remap", "application=system",
			"application=", "application=system",
		)

		pollingURL6Str = insertPorts(r.Replace(pollingURLStr), srv)
	}

	return pollingURL4Str, pollingURL6Str
}

func insertPorts(pollingURLStr string, srv tc.TrafficServer) string {
	if strings.HasPrefix(strings.ToLower(pollingURLStr), "https") {
		if srv.HTTPSPort != 0 {
			pollURL, err := url.Parse(pollingURLStr)
			if err != nil {
				log.Warnf("profile '%s' cache server '%s' polling URL '%s' failed to parse, may not be a valid URL! Using anyway, not using custom HTTPS Port %d!", srv.Profile, srv.FQDN, pollingURLStr, srv.HTTPSPort)
			} else if pollURL.Port() == "" { // if there's both an HTTPS Port and a port in the polling URL, the polling URL takes precedence
				pollURL.Host += ":" + strconv.Itoa(srv.HTTPSPort)
				pollingURLStr = pollURL.String()
			}
		}
	} else {
		if srv.Port != 0 {
			pollURL, err := url.Parse(pollingURLStr)
			if err != nil {
				log.Warnf("profile '%s' cache server '%s' polling URL '%s' failed to parse, may not be a valid URL! Using anyway, not using custom TCP Port %d!", srv.Profile, srv.FQDN, pollingURLStr, srv.Port)
			} else if pollURL.Port() == "" { // if there's both a TCP Port and a port in the polling URL, the polling URL takes precedence
				pollURL.Host += ":" + strconv.Itoa(srv.Port)
				pollingURLStr = pollURL.String()
			}
		}
	}
	return pollingURLStr
}

// createServerStatPollURL takes the health polling URL string, and modifies it to be the stat poll URL.
// Note this does not replace template variables with server values, healthPollURLStr must be the health URL for a given server, not a template.
func createServerStatPollURL(healthPollURLStr string) string {
	return strings.NewReplacer("application=system", "application=").Replace(healthPollURLStr)
}
