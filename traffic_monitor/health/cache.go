package health

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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

// Used as a dummy value for evaluating threshold stats (which don't use real
// combined states).
var dummyCombinedState = tc.IsAvailable{}

// AvailableStr is used to describe the state of a cache server that is
// available to serve traffic.
const AvailableStr = "available"

// UnavailableStr is used to describe the state of a cache server that is not
// available to serve traffic.
const UnavailableStr = "unavailable"

// GetVitals Gets the vitals to decide health on in the right format
func GetVitals(newResult *cache.Result, prevResult *cache.Result, mc *tc.TrafficMonitorConfigMap) {
	if newResult.Error != nil {
		log.Errorf("cache_health.GetVitals() called with an errored Result!")
		return
	}

	if mc == nil {
		log.Errorf("TrafficMonitorConfigMap must not be nil")
		return
	}

	if newResult.InterfaceVitals == nil {
		newResult.InterfaceVitals = map[string]cache.Vitals{}
	}

	// proc.loadavg -- we're using the 1 minute average (!?)
	newResult.Vitals.LoadAvg = newResult.Statistics.Loadavg.One

	ts, exists := mc.TrafficServer[newResult.ID]
	if !exists {
		log.Errorf("cache server not found in config map for cache: %s", newResult.ID)
		return
	}
	if ts.Interfaces == nil {
		log.Warnf("no interfaces reported in config map for cache: %s", newResult.ID)
		return
	}

	var monitoredInterfaces []tc.ServerInterfaceInfo
	for _, srvrIfaceInfo := range mc.TrafficServer[newResult.ID].Interfaces {
		if srvrIfaceInfo.Monitor {
			monitoredInterfaces = append(monitoredInterfaces, srvrIfaceInfo)
		}
	}

	if len(monitoredInterfaces) == 0 {
		log.Warnf("no interfaces selected to be monitored for %v", newResult.ID)
		return
	}

	for _, monitoredInterface := range monitoredInterfaces {
		ifaceName := monitoredInterface.Name
		iface, exists := newResult.Interfaces()[ifaceName]
		if !exists {
			// monitored interface doesn't exist in Result interfaces, skip
			log.Warnf("monitored interface %v does not exist in cache %v", ifaceName, newResult.ID)
			continue
		}

		ifaceVitals := cache.Vitals{
			BytesIn:    iface.BytesIn,
			BytesOut:   iface.BytesOut,
			MaxKbpsOut: iface.Speed * 1000,
		}

		if prevResult != nil && prevResult.InterfaceVitals != nil && prevResult.InterfaceVitals[ifaceName].BytesOut != 0 {
			elapsedTimeInSecs := float64(newResult.Time.UnixNano()-prevResult.Time.UnixNano()) / 1000000000
			ifaceVitals.KbpsOut = int64(float64((ifaceVitals.BytesOut-prevResult.InterfaceVitals[ifaceName].BytesOut)*8/1000) / elapsedTimeInSecs)
		}
		newResult.InterfaceVitals[ifaceName] = ifaceVitals

		// Overflow possible
		newResult.Vitals.BytesOut += iface.BytesOut
		newResult.Vitals.BytesIn += iface.BytesIn
		// TODO JvD: Should we really be running this code every second for every cache polled????? I don't think so.
		newResult.Vitals.MaxKbpsOut += iface.Speed * 1000
	}

	if prevResult != nil && prevResult.Vitals.BytesOut != 0 {
		elapsedTimeInSecs := float64(newResult.Time.UnixNano()-prevResult.Time.UnixNano()) / 1000000000
		newResult.Vitals.KbpsOut = int64(float64((newResult.Vitals.BytesOut-prevResult.Vitals.BytesOut)*8/1000) / elapsedTimeInSecs)
	}

}

// EvalCacheWithStatusInfo evaluates whether the given cache should be marked
// available, taking the server's configured Status into account as well as its
// polling information.
func EvalCacheWithStatusInfo(result cache.ResultInfo, mc *tc.TrafficMonitorConfigMap, status tc.CacheStatus, serverStatus string) (bool, string, string) {
	availability := AvailableStr
	if !result.Available {
		availability = UnavailableStr
	}
	switch {
	case status == tc.CacheStatusInvalid:
		log.Errorf("Cache %v got invalid status from Traffic Ops '%v' - treating as OFFLINE\n", result.ID, serverStatus)
		return false, eventDesc(status, availability+"; invalid status"), ""
	case status == tc.CacheStatusAdminDown:
		return false, eventDesc(status, availability), ""
	case status == tc.CacheStatusOffline:
		log.Errorf("Cache %v set to offline, but still polled\n", result.ID)
		return false, eventDesc(status, availability), ""
	case status == tc.CacheStatusOnline:
		return true, eventDesc(status, availability), ""
	case result.Error != nil:
		return false, eventDesc(status, fmt.Sprintf("%v", result.Error)), ""
	case result.Statistics.NotAvailable == true:
		return false, eventDesc(status, fmt.Sprintf("system.notAvailable == %v", result.Statistics.NotAvailable)), ""
	}
	return result.Available, eventDesc(status, availability), ""
}

// EvalInterface returns whether the given interface should be marked
// available, a boolean of whether the result was over IPv4 (false means it
// was IPv6), a string describing why, and which stat exceeded a threshold. The
// `stats` may be nil, for pollers which don't poll stats. The availability of
// EvalCache MAY NOT be used to directly set the cache's local availability,
// because the threshold stats may not be part of the poller which produced the
// result. Rather, if the cache was previously unavailable from a threshold, it
// must be verified that threshold stat is in the results before setting the
// cache to available. The resultStats may be nil, and if so, won't be checked
// for thresholds. For example, the Health poller doesn't have Stats.
// TODO change to return a `cache.AvailableStatus`
func EvalInterface(infVitals map[string]cache.Vitals, inf tc.ServerInterfaceInfo) (bool, string) {
	if !inf.Monitor {
		return true, ""
	}

	vitals, ok := infVitals[inf.Name]
	if !ok {
		return false, "not found in polled data"
	}

	if inf.MaxBandwidth == nil {
		return true, ""
	}

	if *inf.MaxBandwidth < uint64(vitals.KbpsOut) {
		return false, "maximum bandwidth exceeded"
	}

	return true, ""
}

// EvalAggregate calculates the availability of a cache server as an aggregate
// of server metrics and metrics of its network interfaces.
func EvalAggregate(result cache.ResultInfo, resultStats *threadsafe.ResultStatValHistory, mc *tc.TrafficMonitorConfigMap) (bool, string, string) {
	serverInfo, ok := mc.TrafficServer[string(result.ID)]
	if !ok {
		log.Errorf("Cache %v missing from from Traffic Ops Monitor Config - treating as OFFLINE\n", result.ID)
		return false, "ERROR - server missing in Traffic Ops monitor config", ""
	}
	status := tc.CacheStatusFromString(serverInfo.ServerStatus)
	if status == tc.CacheStatusOnline {
		// return here first, even though EvalCacheWithStatus checks online, because we later assume that if EvalCacheWithStatus returns true, to return false if thresholds are exceeded; but, if the cache is ONLINE, we don't want to check thresholds.
		return true, eventDesc(status, AvailableStr), ""
	}

	profile, ok := mc.Profile[serverInfo.Profile]
	if !ok {
		log.Errorf("Profile '%v' for cache server '%v' missing from monitoring configuration - treating as OFFLINE", serverInfo.Profile, result.ID)
		return false, "ERROR - server profile missing in Traffic Ops monitor config", ""
	}

	avail, eventDescVal, eventMsg := EvalCacheWithStatusInfo(result, mc, status, serverInfo.ServerStatus)
	if !avail {
		return avail, eventDescVal, eventMsg
	}

	computedStats := cache.ComputedStats()

	for stat, threshold := range profile.Parameters.Thresholds {
		resultStat := interface{}(nil)
		computedStatF, ok := computedStats[stat]
		if !ok {
			if resultStats == nil {
				continue
			}
			resultStatHistory := resultStats.Load(stat)
			if len(resultStatHistory) == 0 {
				continue
			}
			resultStat = resultStatHistory[0].Val
		} else {
			resultStat = computedStatF(result, serverInfo, profile, dummyCombinedState)
		}

		resultStatNum, ok := util.ToNumeric(resultStat)
		if !ok {
			log.Errorf("health.EvalCache threshold stat %s was not a number: %v", stat, resultStat)
			continue
		}

		if !inThreshold(threshold, resultStatNum) {
			return false, eventDesc(status, exceedsThresholdMsg(stat, threshold, resultStatNum)), stat
		}
	}

	return avail, eventDescVal, eventMsg
}

// getProcessAvailableTuple gets a function to process an availability tuple
// based on the protocol used.
func getProcessAvailableTuple(protocol config.PollingProtocol) func(cache.AvailableTuple, tc.TrafficServer) bool {
	switch protocol {
	case config.IPv4Only:
		return func(tuple cache.AvailableTuple, _ tc.TrafficServer) bool {
			return tuple.IPv4
		}
	case config.IPv6Only:
		return func(tuple cache.AvailableTuple, _ tc.TrafficServer) bool {
			return tuple.IPv6
		}
	case config.Both:
		return func(tuple cache.AvailableTuple, serverInfo tc.TrafficServer) bool {
			if serverInfo.IPv4() == "" {
				return tuple.IPv6
			} else if serverInfo.IPv6() == "" {
				return tuple.IPv4
			}
			return tuple.IPv4 || tuple.IPv6
		}
	default:
		log.Errorf("received an unknown Polling Protocol: %s", protocol)
	}
	return func(cache.AvailableTuple, tc.TrafficServer) bool { return false }
}

// CalcAvailability calculates the availability of each cache in results.
// statResultHistory may be nil, in which case stats won't be used to calculate
// availability.
func CalcAvailability(
	results []cache.Result,
	pollerName string,
	statResultHistory *threadsafe.ResultStatHistory,
	mc tc.TrafficMonitorConfigMap,
	toData todata.TOData,
	localCacheStatusThreadsafe threadsafe.CacheAvailableStatus,
	localStates peer.CRStatesThreadsafe,
	events ThreadsafeEvents,
	protocol config.PollingProtocol,
) {
	localCacheStatuses := localCacheStatusThreadsafe.Get().Copy()
	var statResultsVal *threadsafe.CacheStatHistory
	processAvailableTuple := getProcessAvailableTuple(protocol)

	for _, result := range results {
		if statResultHistory != nil {
			t := statResultHistory.LoadOrStore(result.ID)
			statResultsVal = &t
		}
		serverInfo, ok := mc.TrafficServer[result.ID]
		if !ok {
			log.Errorf("Cache %v missing from from Traffic Ops Monitor Config - treating as OFFLINE\n", result.ID)
		}

		availStatus := cache.AvailableStatus{
			LastCheckedIPv4:    result.UsingIPv4,
			ProcessedAvailable: true,
			Poller:             pollerName,
			Status:             serverInfo.ServerStatus,
		}

		lastStatus, ok := localCacheStatuses[result.ID]
		if ok {
			if result.UsingIPv4 {
				availStatus.Available.IPv4 = true
				availStatus.Available.IPv6 = serverInfo.IPv6() != "" && lastStatus.Available.IPv6
			} else {
				availStatus.Available.IPv6 = true
				availStatus.Available.IPv4 = serverInfo.IPv4() != "" && lastStatus.Available.IPv4
			}
		}

		reasons := []string{}
		resultInfo := cache.ToInfo(result)
		for _, inf := range serverInfo.Interfaces {
			if !inf.Monitor {
				continue
			}

			available, why := EvalInterface(resultInfo.InterfaceVitals, inf)
			if result.UsingIPv4 {
				availStatus.Available.IPv4 = availStatus.Available.IPv4 && available
			} else {
				availStatus.Available.IPv6 = availStatus.Available.IPv6 && available
			}

			if why != "" {
				reasons = append(reasons, inf.Name+": "+why)
			}
		}

		var aggIsAvailable bool
		var aggWhyAvailable string
		var aggUnavailableStat string

		if statResultsVal != nil {
			aggIsAvailable, aggWhyAvailable, aggUnavailableStat = EvalAggregate(cache.ToInfo(result), &statResultsVal.Stats, &mc)
		} else {
			aggIsAvailable, aggWhyAvailable, aggUnavailableStat = EvalAggregate(cache.ToInfo(result), nil, &mc)
		}

		if result.UsingIPv4 {
			availStatus.Available.IPv4 = availStatus.Available.IPv4 && aggIsAvailable
		} else {
			availStatus.Available.IPv6 = availStatus.Available.IPv6 && aggIsAvailable
		}

		availStatus.ProcessedAvailable = processAvailableTuple(availStatus.Available, serverInfo)

		if aggWhyAvailable != "" {
			reasons = append([]string{aggWhyAvailable}, reasons...)
		}
		availStatus.Why = strings.Join(reasons, "; ")
		if aggUnavailableStat != "" {
			availStatus.UnavailableStat = aggUnavailableStat
		}

		localStates.SetCache(tc.CacheName(result.ID), tc.IsAvailable{
			IsAvailable:    availStatus.ProcessedAvailable,
			Ipv4Available:  availStatus.Available.IPv4,
			Ipv6Available:  availStatus.Available.IPv6,
			DirectlyPolled: true, // we know this cache was directly polled because otherwise we wouldn't have a cache.Result for it
			Status:         availStatus.Why,
			LastPoll:       result.Time,
		})

		if available, ok := localStates.GetCache(tc.CacheName(result.ID)); !ok || available.IsAvailable != lastStatus.ProcessedAvailable {
			protocol := "IPv4"
			if !availStatus.LastCheckedIPv4 {
				protocol = "IPv6"
			}
			log.Infof("Changing state for %s was: %t now: %t because %s poller: %v on protocol %v error: %v",
				result.ID, available.IsAvailable, availStatus.ProcessedAvailable, availStatus.Why, pollerName, protocol, result.Error)

			event := Event{
				Time:          Time(time.Now()),
				Description:   "Protocol (" + protocol + ") " + availStatus.Why + " (" + pollerName + ") ",
				Name:          result.ID,
				Hostname:      result.ID,
				Type:          toData.ServerTypes[tc.CacheName(result.ID)].String(),
				Available:     availStatus.ProcessedAvailable,
				IPv4Available: availStatus.Available.IPv4,
				IPv6Available: availStatus.Available.IPv6,
			}
			events.Add(event)
		}

		localCacheStatuses[result.ID] = availStatus
	}
	calculateDeliveryServiceState(localStates)
	localCacheStatusThreadsafe.Set(localCacheStatuses)
}

// ExceedsThresholdMsg returns a human-readable message for why the given value exceeds the threshold. It does NOT check whether the value actually exceeds the threshold; call `InThreshold` to check first.
func exceedsThresholdMsg(stat string, threshold tc.HealthThreshold, val float64) string {
	switch threshold.Comparator {
	case "=":
		return fmt.Sprintf("%s not equal (%.2f != %.2f)", stat, val, threshold.Val)
	case ">":
		return fmt.Sprintf("%s too low (%.2f < %.2f)", stat, val, threshold.Val)
	case "<":
		return fmt.Sprintf("%s too high (%.2f > %.2f)", stat, val, threshold.Val)
	case ">=":
		return fmt.Sprintf("%s too low (%.2f <= %.2f)", stat, val, threshold.Val)
	case "<=":
		return fmt.Sprintf("%s too high (%.2f >= %.2f)", stat, val, threshold.Val)
	default:
		return fmt.Sprintf("ERROR: Invalid Threshold: %+v", threshold)
	}
}

func inThreshold(threshold tc.HealthThreshold, val float64) bool {
	switch threshold.Comparator {
	case "=":
		return val == threshold.Val
	case ">":
		return val > threshold.Val
	case "<":
		return val < threshold.Val
	case ">=":
		return val >= threshold.Val
	case "<=":
		return val <= threshold.Val
	default:
		log.Errorf("Invalid Threshold: %+v", threshold)
		return true // for safety, if a threshold somehow gets corrupted, don't start marking caches down.
	}
}

func eventDesc(status tc.CacheStatus, message string) string {
	return fmt.Sprintf("%s - %s", status, message)
}

// calculateDeliveryServiceState calculates the state of delivery services from the new cache state data `cacheState` and the CRConfig data `deliveryServiceServers` and puts the calculated state in the outparam `deliveryServiceStates`
func calculateDeliveryServiceState(states peer.CRStatesThreadsafe) {
	deliveryServices := states.GetDeliveryServices()
	for deliveryServiceName, deliveryServiceState := range deliveryServices {
		// NOTE: DisabledLocations is always empty, and it's important that it isn't nil, so it serialises to the JSON `[]` instead of `null`.
		// It's no longer populated due to it being an unnecessary optimization for Traffic Router, but the field is kept for compatibility.
		deliveryServiceState.DisabledLocations = []tc.CacheGroupName{}
		states.SetDeliveryService(deliveryServiceName, deliveryServiceState)
	}
}
