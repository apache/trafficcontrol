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
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/util"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/cache"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// GetVitals Gets the vitals to decide health on in the right format
func GetVitals(newResult *cache.Result, prevResult *cache.Result, mc *to.TrafficMonitorConfigMap) {
	if newResult.Error != nil {
		log.Errorf("cache_health.GetVitals() called with an errored Result!")
		return
	}
	// proc.loadavg -- we're using the 1 minute average (!?)
	// value looks like: "0.20 0.07 0.07 1/967 29536" (without the quotes)
	loadAverages := strings.Fields(newResult.Astats.System.ProcLoadavg)
	if len(loadAverages) > 0 {
		oneMinAvg, err := strconv.ParseFloat(loadAverages[0], 64)
		if err != nil {
			setErr(newResult, fmt.Errorf("Error converting load average string '%s': %v", newResult.Astats.System.ProcLoadavg, err))
			return
		}
		newResult.Vitals.LoadAvg = oneMinAvg
	} else {
		setErr(newResult, fmt.Errorf("Can't make sense of '%s' as a load average for %s", newResult.Astats.System.ProcLoadavg, newResult.ID))
		return
	}

	// proc.net.dev -- need to compare to prevSample
	// value looks like
	// "bond0:8495786321839 31960528603    0    0    0     0          0   2349716 143283576747316 101104535041    0    0    0     0       0          0"
	// (without the quotes)
	parts := strings.Split(newResult.Astats.System.ProcNetDev, ":")
	if len(parts) > 1 {
		numbers := strings.Fields(parts[1])
		var err error
		newResult.Vitals.BytesOut, err = strconv.ParseInt(numbers[8], 10, 64)
		if err != nil {
			setErr(newResult, fmt.Errorf("Error converting BytesOut from procnetdev: %v", err))
			return
		}
		newResult.Vitals.BytesIn, err = strconv.ParseInt(numbers[0], 10, 64)
		if err != nil {
			setErr(newResult, fmt.Errorf("Error converting BytesIn from procnetdev: %v", err))
			return
		}
		if prevResult != nil && prevResult.Vitals.BytesOut != 0 {
			elapsedTimeInSecs := float64(newResult.Time.UnixNano()-prevResult.Time.UnixNano()) / 1000000000
			newResult.Vitals.KbpsOut = int64(float64(((newResult.Vitals.BytesOut - prevResult.Vitals.BytesOut) * 8 / 1000)) / elapsedTimeInSecs)
		} else {
			// log.Infoln("prevResult == nil for id " + newResult.Id + ". Hope we're just starting up?")
		}
	} else {
		setErr(newResult, fmt.Errorf("Error parsing procnetdev: no fields found"))
		return
	}

	// inf.speed -- value looks like "10000" (without the quotes) so it is in Mbps.
	// TODO JvD: Should we really be running this code every second for every cache polled????? I don't think so.
	interfaceBandwidth := newResult.Astats.System.InfSpeed
	newResult.Vitals.MaxKbpsOut = int64(interfaceBandwidth) * 1000

	// log.Infoln(newResult.Id, "BytesOut", newResult.Vitals.BytesOut, "BytesIn", newResult.Vitals.BytesIn, "Kbps", newResult.Vitals.KbpsOut, "max", newResult.Vitals.MaxKbpsOut)
}

// EvalCache returns whether the given cache should be marked available, a string describing why, and which stat exceeded a threshold. The `stats` may be nil, for pollers which don't poll stats.
// The availability of EvalCache MAY NOT be used to directly set the cache's local availability, because the threshold stats may not be part of the poller which produced the result. Rather, if the cache was previously unavailable from a threshold, it must be verified that threshold stat is in the results before setting the cache to available.
// TODO change to return a `cache.AvailableStatus`
func EvalCache(result cache.ResultInfo, resultStats cache.ResultStatValHistory, mc *to.TrafficMonitorConfigMap) (bool, string, string) {
	serverInfo, ok := mc.TrafficServer[string(result.ID)]
	if !ok {
		log.Errorf("Cache %v missing from from Traffic Ops Monitor Config - treating as OFFLINE\n", result.ID)
		return false, "ERROR - server missing in Traffic Ops monitor config", ""
	}
	serverProfile, ok := mc.Profile[serverInfo.Profile]
	if !ok {
		log.Errorf("Cache %v profile %v missing from from Traffic Ops Monitor Config - treating as OFFLINE\n", result.ID, serverInfo.Profile)
		return false, "ERROR - server profile missing in Traffic Ops monitor config", ""
	}

	status := enum.CacheStatusFromString(serverInfo.Status)
	if status == enum.CacheStatusInvalid {
		log.Errorf("Cache %v got invalid status from Traffic Ops '%v' - treating as Reported\n", result.ID, serverInfo.Status)
	}

	availability := "available"
	if !result.Available {
		availability = "unavailable"
	}

	switch {
	case status == enum.CacheStatusInvalid:
		log.Errorf("Cache %v got invalid status from Traffic Ops '%v' - treating as OFFLINE\n", result.ID, serverInfo.Status)
		return false, eventDesc(status, availability+"; invalid status"), ""
	case status == enum.CacheStatusAdminDown:
		return false, eventDesc(status, availability), ""
	case status == enum.CacheStatusOffline:
		log.Errorf("Cache %v set to offline, but still polled\n", result.ID)
		return false, eventDesc(status, availability), ""
	case status == enum.CacheStatusOnline:
		return true, eventDesc(status, availability), ""
	case result.Error != nil:
		return false, eventDesc(status, fmt.Sprintf("%v", result.Error)), ""
	}

	computedStats := cache.ComputedStats()

	for stat, threshold := range serverProfile.Parameters.Thresholds {
		resultStat := interface{}(nil)
		if computedStatF, ok := computedStats[stat]; ok {
			dummyCombinedstate := peer.IsAvailable{} // the only stats which use combinedState are things like isAvailable, which don't make sense to ever be thresholds.
			resultStat = computedStatF(result, serverInfo, serverProfile, dummyCombinedstate)
		} else {
			if resultStats == nil {
				continue
			}
			resultStatHistory, ok := resultStats[stat]
			if !ok {
				continue
			}
			if len(resultStatHistory) < 1 {
				continue
			}
			resultStat = resultStatHistory[0].Val
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

	return result.Available, eventDesc(status, availability), ""
}

// CalcAvailability calculates the availability of the cache, from the given result. Availability is stored in `localCacheStatus` and `localStates`, and if the status changed an event is added to `events`. statResultHistory may be nil, for pollers which don't poll stats.
// TODO add enum for poller names?
func CalcAvailability(results []cache.Result, pollerName string, statResultHistory cache.ResultStatHistory, mc to.TrafficMonitorConfigMap, toData todata.TOData, localCacheStatusThreadsafe threadsafe.CacheAvailableStatus, localStates peer.CRStatesThreadsafe, events ThreadsafeEvents) {
	localCacheStatuses := localCacheStatusThreadsafe.Get().Copy()
	for _, result := range results {
		statResults := cache.ResultStatValHistory(nil)
		if statResultHistory != nil {
			statResults = statResultHistory[result.ID]
		}

		isAvailable, whyAvailable, unavailableStat := EvalCache(cache.ToInfo(result), statResults, &mc)

		// if the cache is now Available, and was previously unavailable due to a threshold, make sure this poller contains the stat which exceeded the threshold.
		if previousStatus, hasPreviousStatus := localCacheStatuses[result.ID]; isAvailable && hasPreviousStatus && !previousStatus.Available && previousStatus.UnavailableStat != "" {
			if !result.HasStat(previousStatus.UnavailableStat) {
				return
			}
		}
		localCacheStatuses[result.ID] = cache.AvailableStatus{
			Available:       isAvailable,
			Status:          mc.TrafficServer[string(result.ID)].Status,
			Why:             whyAvailable,
			UnavailableStat: unavailableStat,
			Poller:          pollerName,
		} // TODO move within localStates?

		if available, ok := localStates.GetCache(result.ID); !ok || available.IsAvailable != isAvailable {
			log.Infof("Changing state for %s was: %t now: %t because %s poller: %v error: %v", result.ID, available.IsAvailable, isAvailable, whyAvailable, pollerName, result.Error)
			events.Add(Event{Time: Time(time.Now()), Description: whyAvailable + " (" + pollerName + ")", Name: string(result.ID), Hostname: string(result.ID), Type: toData.ServerTypes[result.ID].String(), Available: isAvailable})
		}

		localStates.SetCache(result.ID, peer.IsAvailable{IsAvailable: isAvailable})
	}
	calculateDeliveryServiceState(toData.DeliveryServiceServers, localStates)
	localCacheStatusThreadsafe.Set(localCacheStatuses)
}

func setErr(newResult *cache.Result, err error) {
	newResult.Error = err
	newResult.Available = false
}

// ExceedsThresholdMsg returns a human-readable message for why the given value exceeds the threshold. It does NOT check whether the value actually exceeds the threshold; call `InThreshold` to check first.
func exceedsThresholdMsg(stat string, threshold to.HealthThreshold, val float64) string {
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

func inThreshold(threshold to.HealthThreshold, val float64) bool {
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

func eventDesc(status enum.CacheStatus, message string) string {
	return fmt.Sprintf("%s - %s", status, message)
}

//calculateDeliveryServiceState calculates the state of delivery services from the new cache state data `cacheState` and the CRConfig data `deliveryServiceServers` and puts the calculated state in the outparam `deliveryServiceStates`
func calculateDeliveryServiceState(deliveryServiceServers map[enum.DeliveryServiceName][]enum.CacheName, states peer.CRStatesThreadsafe) {
	deliveryServices := states.GetDeliveryServices()
	for deliveryServiceName, deliveryServiceState := range deliveryServices {
		if _, ok := deliveryServiceServers[deliveryServiceName]; !ok {
			log.Infof("CRConfig does not have delivery service %s, but traffic monitor poller does; skipping\n", deliveryServiceName)
			continue
		}
		deliveryServiceState.DisabledLocations = []enum.CacheName{} // it's important this isn't nil, so it serialises to the JSON `[]` instead of `null`
		for _, server := range deliveryServiceServers[deliveryServiceName] {
				deliveryServiceState.DisabledLocations = append(deliveryServiceState.DisabledLocations, server)
			}
		}
		states.SetDeliveryService(deliveryServiceName, deliveryServiceState)
	}
}
