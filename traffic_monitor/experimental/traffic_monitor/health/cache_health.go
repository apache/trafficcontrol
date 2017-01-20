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

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/cache"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func setError(newResult *cache.Result, err error) {
	newResult.Error = err
	newResult.Available = false
}

// GetVitals Gets the vitals to decide health on in the right format
func GetVitals(newResult *cache.Result, prevResult *cache.Result, mc *traffic_ops.TrafficMonitorConfigMap) {
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
			setError(newResult, fmt.Errorf("Error converting load average string '%s': %v", newResult.Astats.System.ProcLoadavg, err))
			return
		}
		newResult.Vitals.LoadAvg = oneMinAvg
	} else {
		setError(newResult, fmt.Errorf("Can't make sense of '%s' as a load average for %s", newResult.Astats.System.ProcLoadavg, newResult.ID))
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
			setError(newResult, fmt.Errorf("Error converting BytesOut from procnetdev: %v", err))
			return
		}
		newResult.Vitals.BytesIn, err = strconv.ParseInt(numbers[0], 10, 64)
		if err != nil {
			setError(newResult, fmt.Errorf("Error converting BytesIn from procnetdev: %v", err))
			return
		}
		if prevResult != nil && prevResult.Vitals.BytesOut != 0 {
			elapsedTimeInSecs := float64(newResult.Time.UnixNano()-prevResult.Time.UnixNano()) / 1000000000
			newResult.Vitals.KbpsOut = int64(float64(((newResult.Vitals.BytesOut - prevResult.Vitals.BytesOut) * 8 / 1000)) / elapsedTimeInSecs)
		} else {
			// log.Infoln("prevResult == nil for id " + newResult.Id + ". Hope we're just starting up?")
		}
	} else {
		setError(newResult, fmt.Errorf("Error parsing procnetdev: no fields found"))
		return
	}

	// inf.speed -- value looks like "10000" (without the quotes) so it is in Mbps.
	// TODO JvD: Should we really be running this code every second for every cache polled????? I don't think so.
	interfaceBandwidth := newResult.Astats.System.InfSpeed
	newResult.Vitals.MaxKbpsOut = int64(interfaceBandwidth)*1000 - mc.Profile[mc.TrafficServer[string(newResult.ID)].Profile].Parameters.MinFreeKbps

	// log.Infoln(newResult.Id, "BytesOut", newResult.Vitals.BytesOut, "BytesIn", newResult.Vitals.BytesIn, "Kbps", newResult.Vitals.KbpsOut, "max", newResult.Vitals.MaxKbpsOut)
}

// getKbpsThreshold returns the numeric kbps threshold, from the Traffic Ops string value. If there is a parse error, it logs a warning and returns the max floating point number, signifying no limit
// TODO add float64 to Traffic Ops Client interface
func getKbpsThreshold(threshStr string) (int64, bool) {
	if len(threshStr) == 0 {
		return 0, false
	}
	if threshStr[0] == '>' {
		threshStr = threshStr[1:]
	}
	thresh, err := strconv.ParseInt(threshStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return thresh, true
}

// TODO add time.Duration to Traffic Ops Client interface
func getQueryThreshold(threshInt int64) (time.Duration, bool) {
	if threshInt == 0 {
		return 0, false
	}
	return time.Duration(threshInt) * time.Millisecond, true
}

func cacheCapacityKbps(result cache.Result) int64 {
	kbpsInMbps := int64(1000)
	return int64(result.Astats.System.InfSpeed) * kbpsInMbps
}

// EvalCache returns whether the given cache should be marked available, and a string describing why
func EvalCache(result cache.Result, mc *traffic_ops.TrafficMonitorConfigMap) (bool, string) {
	toServer := mc.TrafficServer[string(result.ID)]
	status := enum.CacheStatusFromString(toServer.Status)
	params := mc.Profile[toServer.Profile].Parameters
	kbpsThreshold, hasKbpsThreshold := getKbpsThreshold(params.HealthThresholdAvailableBandwidthInKbps)
	queryTimeThreshold, hasQueryTimeThreshold := getQueryThreshold(int64(params.HealthThresholdQueryTime))

	availability := "available"
	if !result.Available {
		availability = "unavailable"
	}

	switch {
	case status == enum.CacheStatusInvalid:
		log.Errorf("Cache %v got invalid status from Traffic Ops '%v' - treating as OFFLINE\n", result.ID, toServer.Status)
		return false, getEventDescription(status, availability+"; invalid status")
	case status == enum.CacheStatusAdminDown:
		return false, getEventDescription(status, availability)
	case status == enum.CacheStatusOffline:
		log.Errorf("Cache %v set to OFFLINE, but still polled\n", result.ID)
		return false, getEventDescription(status, availability)
	case status == enum.CacheStatusOnline:
		return true, getEventDescription(status, availability)
	case result.Error != nil:
		return false, getEventDescription(status, fmt.Sprintf("%v", result.Error))
	case result.Vitals.LoadAvg > params.HealthThresholdLoadAvg && params.HealthThresholdLoadAvg != 0:
		return false, getEventDescription(status, fmt.Sprintf("loadavg too high (%.5f > %.5f)", result.Vitals.LoadAvg, params.HealthThresholdLoadAvg))
	case hasKbpsThreshold && cacheCapacityKbps(result)-result.Vitals.KbpsOut < kbpsThreshold:
		return false, getEventDescription(status, fmt.Sprintf("availableBandwidthInKbps too low (%d < %d)", cacheCapacityKbps(result)-result.Vitals.KbpsOut, kbpsThreshold))
	case hasQueryTimeThreshold && result.RequestTime > queryTimeThreshold:
		return false, getEventDescription(status, fmt.Sprintf("queryTime too high (%.5f > %.5f)", float64(result.RequestTime.Nanoseconds())/1e6, float64(queryTimeThreshold.Nanoseconds())/1e6))
	default:
		return result.Available, getEventDescription(status, availability)
	}
}

func getEventDescription(status enum.CacheStatus, message string) string {
	return fmt.Sprintf("%s - %s", status, message)
}
