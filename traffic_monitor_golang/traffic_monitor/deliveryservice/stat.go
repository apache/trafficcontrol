package deliveryservice

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
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/util"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/cache"
	dsdata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/deliveryservicedata"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/health"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/peer"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/trafficopsdata"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

// TODO remove 'ds' and 'stat' from names

func setStaticData(dsStats dsdata.Stats, dsServers map[enum.DeliveryServiceName][]enum.CacheName) dsdata.Stats {
	for ds, stat := range dsStats.DeliveryService {
		stat.CommonStats.CachesConfiguredNum.Value = int64(len(dsServers[ds]))
		dsStats.DeliveryService[ds] = stat // TODO consider changing dsStats.DeliveryService[ds] to a pointer so this kind of thing isn't necessary; possibly more performant, as well
	}
	return dsStats
}

func addAvailableData(dsStats dsdata.Stats, crStates peer.Crstates, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverDs map[enum.CacheName][]enum.DeliveryServiceName, serverTypes map[enum.CacheName]enum.CacheType, precomputed map[enum.CacheName]cache.PrecomputedData) (dsdata.Stats, error) {
	for cache, available := range crStates.Caches {
		cacheGroup, ok := serverCachegroups[cache]
		if !ok {
			log.Infof("CreateStats not adding availability data for '%s': not found in Cachegroups\n", cache)
			continue
		}
		deliveryServices, ok := serverDs[cache]
		if !ok {
			log.Infof("CreateStats not adding availability data for '%s': not found in DeliveryServices\n", cache)
			continue
		}
		cacheType, ok := serverTypes[cache]
		if !ok {
			log.Infof("CreateStats not adding availability data for '%s': not found in Server Types\n", cache)
			continue
		}

		for _, deliveryService := range deliveryServices {
			if deliveryService == "" {
				log.Errorln("EMPTY addAvailableData DS") // various bugs in other functions can cause this - this will help identify and debug them.
				continue
			}

			stat, ok := dsStats.DeliveryService[deliveryService]
			if !ok {
				log.Infof("CreateStats not adding availability data for '%s': not found in Stats\n", cache)
				continue // TODO log warning? Error?
			}

			if available.IsAvailable {
				stat.CommonStats.IsAvailable.Value = true
				stat.CommonStats.IsHealthy.Value = true
				stat.CommonStats.CachesAvailableNum.Value++
				cacheGroupStats := stat.CacheGroups[cacheGroup]
				cacheGroupStats.IsAvailable.Value = true
				stat.CacheGroups[cacheGroup] = cacheGroupStats
				stat.TotalStats.IsAvailable.Value = true
				typeStats := stat.Types[cacheType]
				typeStats.IsAvailable.Value = true
				stat.Types[cacheType] = typeStats
			}

			// TODO fix nested ifs
			if pc, ok := precomputed[cache]; ok {
				if pc.Reporting {
					stat.CommonStats.CachesReporting[cache] = true
				} else {
					log.Debugf("no reporting %v %v\n", cache, deliveryService)
				}
			} else {
				log.Debugf("no result for %v %v\n", cache, deliveryService)
			}

			dsStats.DeliveryService[deliveryService] = stat // TODO Necessary? Remove?
		}
	}

	// TODO move to its own func?
	for dsName, ds := range crStates.Deliveryservice {
		stat, ok := dsStats.DeliveryService[dsName]
		if !ok {
			log.Warnf("CreateStats not adding disabledLocations for '%s': not found in Stats\n", dsName)
			continue // TODO log warning? Error?
		}

		// TODO determine if a deep copy is necessary
		stat.CommonStats.CachesDisabled = make([]string, len(ds.DisabledLocations), len(ds.DisabledLocations))
		for i, v := range ds.DisabledLocations {
			stat.CommonStats.CachesDisabled[i] = string(v)
		}
		dsStats.DeliveryService[dsName] = stat // TODO Necessary? Remove?
	}

	return dsStats, nil
}

func newLastDSStat() dsdata.LastDSStat {
	return dsdata.LastDSStat{
		CacheGroups: map[enum.CacheGroupName]dsdata.LastStatsData{},
		Type:        map[enum.CacheType]dsdata.LastStatsData{},
		Caches:      map[enum.CacheName]dsdata.LastStatsData{},
	}
}

// BytesPerKilobit is the number of bytes in a kilobit.
const BytesPerKilobit = 125

func addLastStat(lastData dsdata.LastStatData, newStat int64, newStatTime time.Time) (dsdata.LastStatData, error) {
	if lastData.Time == newStatTime {
		return lastData, nil // TODO fix callers to not pass the same stat twice
	}

	if newStat < lastData.Stat {
		// if a new stat comes in lower than current, assume rollover, set the 'last stat' to the new one, but leave PerSec what it was (not negative).
		lastData.Stat = newStat
		lastData.Time = newStatTime
		err := fmt.Errorf("new stat '%d'@'%v' value less than last stat '%d'@'%v'", newStat, newStatTime, lastData.Stat, lastData.Time)
		return lastData, err
	}

	if newStatTime.Before(lastData.Time) {
		return lastData, fmt.Errorf("new stat '%d'@'%v' time less than last stat '%d'@'%v'", newStat, newStatTime, lastData.Stat, lastData.Time)
	}

	if lastData.Stat != 0 {
		lastData.PerSec = float64(newStat-lastData.Stat) / newStatTime.Sub(lastData.Time).Seconds()
	}

	lastData.Stat = newStat
	lastData.Time = newStatTime
	return lastData, nil
}

func addLastStats(lastData dsdata.LastStatsData, newStats dsdata.StatCacheStats, newStatsTime time.Time) (dsdata.LastStatsData, error) {
	errs := []error{nil, nil, nil, nil, nil}
	lastData.Bytes, errs[0] = addLastStat(lastData.Bytes, newStats.OutBytes.Value, newStatsTime)
	lastData.Status2xx, errs[1] = addLastStat(lastData.Status2xx, newStats.Status2xx.Value, newStatsTime)
	lastData.Status3xx, errs[2] = addLastStat(lastData.Status3xx, newStats.Status3xx.Value, newStatsTime)
	lastData.Status4xx, errs[3] = addLastStat(lastData.Status4xx, newStats.Status4xx.Value, newStatsTime)
	lastData.Status5xx, errs[4] = addLastStat(lastData.Status5xx, newStats.Status5xx.Value, newStatsTime)
	return lastData, util.JoinErrors(errs)
}

func addLastStatsToStatCacheStats(s dsdata.StatCacheStats, l dsdata.LastStatsData) dsdata.StatCacheStats {
	s.Kbps.Value = l.Bytes.PerSec / BytesPerKilobit
	s.Tps2xx.Value = l.Status2xx.PerSec
	s.Tps3xx.Value = l.Status3xx.PerSec
	s.Tps4xx.Value = l.Status4xx.PerSec
	s.Tps5xx.Value = l.Status5xx.PerSec
	s.TpsTotal.Value = s.Tps2xx.Value + s.Tps3xx.Value + s.Tps4xx.Value + s.Tps5xx.Value
	return s
}

// addLastDSStatTotals takes a LastDSStat with only raw `Caches` data, and calculates and sets the `CacheGroups`, `Type`, and `Total` data, and returns the augmented structure.
func addLastDSStatTotals(lastStat dsdata.LastDSStat, cachesReporting map[enum.CacheName]bool, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverTypes map[enum.CacheName]enum.CacheType) dsdata.LastDSStat {
	cacheGroups := map[enum.CacheGroupName]dsdata.LastStatsData{}
	cacheTypes := map[enum.CacheType]dsdata.LastStatsData{}
	total := dsdata.LastStatsData{}
	for cacheName, cacheStats := range lastStat.Caches {
		if !cachesReporting[cacheName] {
			continue
		}

		if cacheGroup, ok := serverCachegroups[cacheName]; ok {
			cacheGroups[cacheGroup] = cacheGroups[cacheGroup].Sum(cacheStats)
		} else {
			log.Warnf("while computing delivery service data, cache %v not in cachegroups\n", cacheName)
		}

		if cacheType, ok := serverTypes[cacheName]; ok {
			cacheTypes[cacheType] = cacheTypes[cacheType].Sum(cacheStats)
		} else {
			log.Warnf("while computing delivery service data, cache %v not in types\n", cacheName)
		}
		total = total.Sum(cacheStats)
	}
	lastStat.CacheGroups = cacheGroups
	lastStat.Type = cacheTypes
	lastStat.Total = total
	return lastStat
}

// addDSPerSecStats calculates and adds the per-second delivery service stats to both the Stats and LastStats structures, and returns the augmented structures.
func addDSPerSecStats(dsName enum.DeliveryServiceName, stat dsdata.Stat, lastStats dsdata.LastStats, dsStats dsdata.Stats, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverTypes map[enum.CacheName]enum.CacheType, mc to.TrafficMonitorConfigMap, events health.ThreadsafeEvents, precomputed map[enum.CacheName]cache.PrecomputedData) (dsdata.Stats, dsdata.LastStats) {
	err := error(nil)
	lastStat, lastStatExists := lastStats.DeliveryServices[dsName]
	if !lastStatExists {
		lastStat = newLastDSStat()
	}

	for cacheName, cacheStats := range stat.Caches {
		if _, ok := precomputed[cacheName]; ok {
			lastStat.Caches[cacheName], err = addLastStats(lastStat.Caches[cacheName], cacheStats, precomputed[cacheName].Time)
			if err != nil {
				log.Warnf("%v adding per-second stats for cache %v: %v", dsName, cacheName, err)
				continue
			}
		}
		cacheStats.Kbps.Value = lastStat.Caches[cacheName].Bytes.PerSec / BytesPerKilobit
		stat.Caches[cacheName] = cacheStats
	}

	lastStat = addLastDSStatTotals(lastStat, stat.CommonStats.CachesReporting, serverCachegroups, serverTypes)

	for cacheGroup, cacheGroupStat := range lastStat.CacheGroups {
		stat.CacheGroups[cacheGroup] = addLastStatsToStatCacheStats(stat.CacheGroups[cacheGroup], cacheGroupStat)
	}
	for cacheType, cacheTypeStat := range lastStat.Type {
		stat.Types[cacheType] = addLastStatsToStatCacheStats(stat.Types[cacheType], cacheTypeStat)
	}
	stat.TotalStats = addLastStatsToStatCacheStats(stat.TotalStats, lastStat.Total)

	dsErr := getDSErr(dsName, stat.TotalStats, mc)
	if dsErr != nil {
		stat.CommonStats.IsAvailable.Value = false
		stat.CommonStats.IsHealthy.Value = false
		stat.CommonStats.ErrorStr.Value = err.Error()
	}

	getEvent := func(desc string) health.Event {
		return health.Event{
			Time:        health.Time(time.Now()),
			Description: desc,
			Name:        dsName.String(),
			Hostname:    dsName.String(),
			Type:        "Delivery Service",
			Available:   stat.CommonStats.IsAvailable.Value,
		}
	}
	if stat.CommonStats.IsAvailable.Value == false && lastStat.Available == true {
		events.Add(getEvent(dsErr.Error()))
	} else if stat.CommonStats.IsAvailable.Value == true && lastStat.Available == false {
		events.Add(getEvent("REPORTED - available"))
	}
	lastStat.Available = stat.CommonStats.IsAvailable.Value

	lastStats.DeliveryServices[dsName] = lastStat
	dsStats.DeliveryService[dsName] = stat
	return dsStats, lastStats
}

// latestBytes returns the most recent OutBytes from the given cache results, and the time of that result. It assumes zero results are not valid, but nonzero results with errors are valid.
func latestBytes(p cache.PrecomputedData) (int64, time.Time, error) {
	if p.OutBytes == 0 {
		return 0, time.Time{}, fmt.Errorf("no valid results")
	}
	return p.OutBytes, p.Time, nil
}

// addCachePerSecStats calculates the cache per-second stats, adds them to LastStats, and returns the augmented object.
func addCachePerSecStats(cacheName enum.CacheName, precomputed cache.PrecomputedData, lastStats dsdata.LastStats) dsdata.LastStats {
	outBytes, outBytesTime, err := latestBytes(precomputed) // it's ok if `latestBytes` returns 0s with an error, `addLastStat` will refrain from setting it (unless the previous calculation was nonzero, in which case it will error appropriately).
	if err != nil {
		log.Warnf("while computing delivery service data for cache %v: %v\n", cacheName, err)
	}
	lastStat := lastStats.Caches[cacheName] // if lastStats.Caches[cacheName] doesn't exist, it will be zero-constructed, and `addLastStat` will refrain from setting the PerSec for zero LastStats
	lastStat.Bytes, err = addLastStat(lastStat.Bytes, outBytes, outBytesTime)
	if err != nil {
		log.Warnf("while computing delivery service data for cache %v: %v\n", cacheName, err)
		return lastStats
	}
	lastStats.Caches[cacheName] = lastStat

	return lastStats
}

// addPerSecStats adds Kbps fields to the NewStats, based on the previous out_bytes in the oldStats, and the time difference.
//
// Traffic Server only updates its data every N seconds. So, often we get a new Stats with the same OutBytes as the previous one,
// So, we must record the last changed value, and the time it changed. Then, if the new OutBytes is different from the previous,
// we set the (new - old) / lastChangedTime as the KBPS, and update the recorded LastChangedTime and LastChangedValue
//
// TODO handle ATS byte rolling (when the `out_bytes` overflows back to 0)
func addPerSecStats(precomputed map[enum.CacheName]cache.PrecomputedData, dsStats dsdata.Stats, lastStats dsdata.LastStats, serverCachegroups map[enum.CacheName]enum.CacheGroupName, serverTypes map[enum.CacheName]enum.CacheType, mc to.TrafficMonitorConfigMap, events health.ThreadsafeEvents) (dsdata.Stats, dsdata.LastStats) {
	for dsName, stat := range dsStats.DeliveryService {
		dsStats, lastStats = addDSPerSecStats(dsName, stat, lastStats, dsStats, serverCachegroups, serverTypes, mc, events, precomputed)
	}
	for cacheName, precomputedData := range precomputed {
		lastStats = addCachePerSecStats(cacheName, precomputedData, lastStats)
	}

	return dsStats, lastStats
}

// CreateStats aggregates and creates statistics from given precomputed stat history. It returns the created stats, information about these stats necessary for the next calculation, and any error.
func CreateStats(precomputed map[enum.CacheName]cache.PrecomputedData, toData todata.TOData, crStates peer.Crstates, lastStats dsdata.LastStats, now time.Time, mc to.TrafficMonitorConfigMap, events health.ThreadsafeEvents) (dsdata.Stats, dsdata.LastStats, error) {
	start := time.Now()
	dsStats := dsdata.NewStats()
	for deliveryService := range toData.DeliveryServiceServers {
		if deliveryService == "" {
			log.Errorf("EMPTY CreateStats deliveryService")
			continue
		}
		dsStats.DeliveryService[deliveryService] = *dsdata.NewStat()
	}
	dsStats = setStaticData(dsStats, toData.DeliveryServiceServers)
	var err error
	dsStats, err = addAvailableData(dsStats, crStates, toData.ServerCachegroups, toData.ServerDeliveryServices, toData.ServerTypes, precomputed) // TODO move after stat summarisation
	if err != nil {
		return dsStats, lastStats, fmt.Errorf("Error getting Cache availability data: %v", err)
	}

	for server, precomputedData := range precomputed {
		cachegroup, ok := toData.ServerCachegroups[server]
		if !ok {
			log.Warnf("server %s has no cachegroup, skipping\n", server)
			continue
		}
		serverType, ok := toData.ServerTypes[server]
		if !ok {
			log.Warnf("server %s not in CRConfig, skipping\n", server)
			continue
		}

		// TODO check result.PrecomputedData.Errors
		for ds, resultStat := range precomputedData.DeliveryServiceStats {
			if ds == "" {
				log.Errorf("EMPTY precomputed delivery service")
				continue
			}

			if _, ok := dsStats.DeliveryService[ds]; !ok {
				dsStats.DeliveryService[ds] = resultStat
				continue
			}
			httpDsStat := dsStats.DeliveryService[ds]
			httpDsStat.TotalStats = httpDsStat.TotalStats.Sum(resultStat.TotalStats)
			httpDsStat.CacheGroups[cachegroup] = httpDsStat.CacheGroups[cachegroup].Sum(resultStat.CacheGroups[cachegroup])
			httpDsStat.Types[serverType] = httpDsStat.Types[serverType].Sum(resultStat.Types[serverType])
			httpDsStat.Caches[server] = httpDsStat.Caches[server].Sum(resultStat.Caches[server])
			httpDsStat.CachesTimeReceived[server] = resultStat.CachesTimeReceived[server]
			httpDsStat.CommonStats = dsStats.DeliveryService[ds].CommonStats
			dsStats.DeliveryService[ds] = httpDsStat // TODO determine if necessary
		}
	}

	perSecStats, lastStats := addPerSecStats(precomputed, dsStats, lastStats, toData.ServerCachegroups, toData.ServerTypes, mc, events)
	log.Infof("CreateStats took %v\n", time.Since(start))
	perSecStats.Time = time.Now()
	return perSecStats, lastStats, nil
}

func getDSErr(dsName enum.DeliveryServiceName, dsStats dsdata.StatCacheStats, monitorConfig to.TrafficMonitorConfigMap) error {
	if tpsThreshold := monitorConfig.DeliveryService[dsName.String()].TotalTPSThreshold; tpsThreshold > 0 && dsStats.TpsTotal.Value > float64(tpsThreshold) {
		return fmt.Errorf("total.tps_total too high (%v > %v)", dsStats.TpsTotal.Value, tpsThreshold)
	}
	if kbpsThreshold := monitorConfig.DeliveryService[dsName.String()].TotalKbpsThreshold; kbpsThreshold > 0 && dsStats.Kbps.Value > float64(kbpsThreshold) {
		return fmt.Errorf("total.kbps too high (%v > %v)", dsStats.Kbps.Value, kbpsThreshold)
	}
	return nil
}
