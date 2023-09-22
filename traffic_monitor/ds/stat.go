package ds

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
	"errors"
	"fmt"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/health"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

// TODO remove 'ds' and 'stat' from names

func setStaticData(dsStats *dsdata.Stats, dsServers map[tc.DeliveryServiceName][]tc.CacheName) {
	for ds, stat := range dsStats.DeliveryService {
		stat.CommonStats.CachesConfiguredNum.Value = int64(len(dsServers[ds]))
	}
}

func addAvailableData(dsStats *dsdata.Stats, crStates tc.CRStates, serverCachegroups map[tc.CacheName]tc.CacheGroupName, serverDs map[tc.CacheName][]tc.DeliveryServiceName, serverTypes map[tc.CacheName]tc.CacheType, precomputed map[tc.CacheName]cache.PrecomputedData) {
	for cache, available := range crStates.Caches {
		cacheGroup, ok := serverCachegroups[cache]
		if !ok {
			log.Infof("CreateStats not adding availability data for '%s': not found in Cachegroups\n", cache)
			continue
		}
		cacheType, ok := serverTypes[cache]
		if !ok {
			log.Infof("CreateStats not adding availability data for '%s': not found in Server Types\n", cache)
			continue
		}
		deliveryServices, ok := serverDs[cache]
		if !ok && cacheType != tc.CacheTypeMid {
			log.Infof("CreateStats not adding availability data for '%s': not found in DeliveryServices\n", cache)
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
				if cacheGroupStats == nil {
					cacheGroupStats = &dsdata.StatCacheStats{} // TODO sync.Pool?
					stat.CacheGroups[cacheGroup] = cacheGroupStats
				}
				cacheGroupStats.IsAvailable.Value = true
				stat.TotalStats.IsAvailable.Value = true
				typeStats := stat.Types[cacheType]
				if typeStats == nil {
					typeStats = &dsdata.StatCacheStats{} // TODO sync.Pool?
					stat.Types[cacheType] = cacheGroupStats
				}
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
		}
	}

	// TODO move to its own func?
	for dsName := range crStates.DeliveryService {
		stat, ok := dsStats.DeliveryService[dsName]
		if !ok {
			log.Infof("CreateStats not adding disabledLocations for '%s': not found in Stats\n", dsName)
			continue // TODO log warning? Error?
		}
		stat.CommonStats.CachesDisabled = make([]string, 0)
	}
}

// BytesPerKilobit is the number of bytes in a kilobit.
const BytesPerKilobit = 125

// Adds the new stat to lastData.
// Note this mutates lastData, adding the new stat to it.
// Also note that lastData may be mutated, even if an error occurs. Specifically, if the new stat is less than the last stat, it will still be set, so that the per-second stats will be properly computed on the next poll.
func addLastStat(lastData *dsdata.LastStatData, newStat int64, newStatTime time.Time) error {
	if lastData == nil {
		return errors.New("nil lastData")
	}

	if lastData.Time == newStatTime {
		return nil // TODO fix callers to not pass the same stat twice
	}

	if newStat < lastData.Stat {
		// if a new stat comes in lower than current, assume rollover, set the 'last stat' to the new one, but leave PerSec what it was (not negative).
		err := fmt.Errorf("new stat '%d'@'%s' value less than last stat '%d'@'%s'", newStat, newStatTime.Format(time.RFC3339Nano), lastData.Stat, lastData.Time.Format(time.RFC3339Nano))
		lastData.Stat = newStat
		lastData.Time = newStatTime
		return err
	}

	if newStatTime.Before(lastData.Time) {
		return fmt.Errorf("new stat '%d'@'%s' time less than last stat '%d'@'%s'", newStat, newStatTime.Format(time.RFC3339Nano), lastData.Stat, lastData.Time.Format(time.RFC3339Nano))
	}

	if lastData.Stat != 0 {
		lastData.PerSec = float64(newStat-lastData.Stat) / newStatTime.Sub(lastData.Time).Seconds()
	}

	lastData.Stat = newStat
	lastData.Time = newStatTime
	return nil
}

// addLastStats adds the new stats to the lastData.
// Note lastData is mutated, with the new stats added to it.
func addLastStats(lastData *dsdata.LastStatsData, newStats *dsdata.StatCacheStats, newStatsTime time.Time) error {
	return util.JoinErrs([]error{
		addLastStat(&lastData.Bytes, newStats.OutBytes.Value, newStatsTime),
		addLastStat(&lastData.Status2xx, newStats.Status2xx.Value, newStatsTime),
		addLastStat(&lastData.Status3xx, newStats.Status3xx.Value, newStatsTime),
		addLastStat(&lastData.Status4xx, newStats.Status4xx.Value, newStatsTime),
		addLastStat(&lastData.Status5xx, newStats.Status5xx.Value, newStatsTime),
	})
}

// addLastStatsToStatCacheStats adds the given LastStatsData to the given StatCacheStats.
// Note s is mutated, with l being added to it.
func addLastStatsToStatCacheStats(s *dsdata.StatCacheStats, l *dsdata.LastStatsData) {
	if s == nil {
		log.Errorln("ds.addLastStatsToStatCacheStats got nil StatCacheStats - skipping!")
		return
	}
	if l == nil {
		log.Errorln("ds.addLastStatsToStatCacheStats got nil LastStatsData - skipping!")
		return
	}
	s.Kbps.Value = l.Bytes.PerSec / BytesPerKilobit
	s.Tps2xx.Value = l.Status2xx.PerSec
	s.Tps3xx.Value = l.Status3xx.PerSec
	s.Tps4xx.Value = l.Status4xx.PerSec
	s.Tps5xx.Value = l.Status5xx.PerSec
	s.TpsTotal.Value = s.Tps2xx.Value + s.Tps3xx.Value + s.Tps4xx.Value + s.Tps5xx.Value
}

// addLastDSStatTotals takes a LastDSStat with only raw `Caches` data, and calculates and sets the `CacheGroups`, `Type`, and `Total` data, and returns the augmented structure.
// Note lastStat is mutated, with the calculated values being set in it.
func addLastDSStatTotals(lastStat *dsdata.LastDSStat, cachesReporting map[tc.CacheName]bool, serverCachegroups map[tc.CacheName]tc.CacheGroupName, serverTypes map[tc.CacheName]tc.CacheType) {
	cacheGroups := map[tc.CacheGroupName]*dsdata.LastStatsData{}
	cacheTypes := map[tc.CacheType]*dsdata.LastStatsData{}
	total := dsdata.LastStatsData{}
	for cacheName, cacheStats := range lastStat.Caches {
		if !cachesReporting[cacheName] {
			continue
		}

		if cacheGroup, ok := serverCachegroups[cacheName]; ok {
			cgStat := cacheGroups[cacheGroup]
			if cgStat == nil {
				cgStat = &dsdata.LastStatsData{}
				cacheGroups[cacheGroup] = cgStat
			}
			cgStat.Sum(cacheStats)
		} else {
			log.Warnf("while computing delivery service data, cache %v not in cachegroups\n", cacheName)
		}

		if cacheType, ok := serverTypes[cacheName]; ok {
			cacheTypeStat := cacheTypes[cacheType]
			if cacheTypeStat == nil {
				cacheTypeStat = &dsdata.LastStatsData{}
				cacheTypes[cacheType] = cacheTypeStat
			}
			cacheTypeStat.Sum(cacheStats)
		} else {
			log.Warnf("while computing delivery service data, cache %v not in types\n", cacheName)
		}
		total.Sum(cacheStats)
	}
	lastStat.CacheGroups = cacheGroups
	lastStat.Type = cacheTypes
	lastStat.Total = total
}

// addDSPerSecStats calculates and adds the per-second delivery service stats to
// both the Stats and LastStats structures. Note this mutates both dsStats and
// lastStats, adding the per-second stats to them.
func addDSPerSecStats(lastStats *dsdata.LastStats, dsName tc.DeliveryServiceName, stat *dsdata.Stat, serverCachegroups map[tc.CacheName]tc.CacheGroupName, serverTypes map[tc.CacheName]tc.CacheType, mc tc.TrafficMonitorConfigMap, events health.ThreadsafeEvents, precomputed map[tc.CacheName]cache.PrecomputedData, states peer.CRStatesThreadsafe) {
	lastStat, lastStatExists := lastStats.DeliveryServices[dsName]
	if !lastStatExists {
		lastStat = &dsdata.LastDSStat{
			CacheGroups: map[tc.CacheGroupName]*dsdata.LastStatsData{},
			Type:        map[tc.CacheType]*dsdata.LastStatsData{},
			Caches:      map[tc.CacheName]*dsdata.LastStatsData{},
		}
		lastStats.DeliveryServices[dsName] = lastStat
	}
	for cacheName, cacheStats := range stat.Caches {
		if cacheStats == nil {
			log.Errorln("ds.addDSPerSecStats - stat.Caches[" + cacheName + "] exists, but unexpected nil! Setting new!")
			stat.Caches[cacheName] = &dsdata.StatCacheStats{}
		}

		if lastStat.Caches[cacheName] == nil {
			lastStat.Caches[cacheName] = &dsdata.LastStatsData{}
		}
		if _, ok := precomputed[cacheName]; ok {
			if err := addLastStats(lastStat.Caches[cacheName], cacheStats, precomputed[cacheName].Time); err != nil {
				log.Warnf("%s adding per-second stats for cache %s: %s", dsName.String(), cacheName.String(), err.Error())
				continue
			}
		}
		cacheStats.Kbps.Value = lastStat.Caches[cacheName].Bytes.PerSec / BytesPerKilobit
	}

	addLastDSStatTotals(lastStat, stat.CommonStats.CachesReporting, serverCachegroups, serverTypes)

	for cacheGroup, cacheGroupStat := range lastStat.CacheGroups {
		if cacheGroupStat == nil {
			log.Errorln("ds.addDSPerSecStats - lastStats.DeliveryServices[" + string(dsName) + "].CacheGroups[" + string(cacheGroup) + "] exists, but unexpected nil! Setting new!")
			lastStat.CacheGroups[cacheGroup] = &dsdata.LastStatsData{}
		}
		if stat.CacheGroups[cacheGroup] == nil {
			stat.CacheGroups[cacheGroup] = &dsdata.StatCacheStats{}
		}
		addLastStatsToStatCacheStats(stat.CacheGroups[cacheGroup], cacheGroupStat)
	}
	for cacheType, cacheTypeStat := range lastStat.Type {
		if cacheTypeStat == nil {
			log.Errorln("ds.addDSPerSecStats - lastStat.DeliveryServices[" + string(dsName) + "].Type[" + string(cacheType) + "] exists, but unexpected nil! Setting new!")
			lastStat.Type[cacheType] = &dsdata.LastStatsData{}
		}
		if stat.Types[cacheType] == nil {
			stat.Types[cacheType] = &dsdata.StatCacheStats{}
		}
		addLastStatsToStatCacheStats(stat.Types[cacheType], cacheTypeStat)
	}
	addLastStatsToStatCacheStats(&stat.TotalStats, &lastStat.Total)

	dsErr := getDSErr(dsName, stat.TotalStats, mc)
	if dsErr != nil {
		stat.CommonStats.IsAvailable.Value = false
		stat.CommonStats.IsHealthy.Value = false
		stat.CommonStats.ErrorStr.Value = dsErr.Error()
	}
	//it's ok to ignore the 'ok' return here.  If the DS doesn't exist, an empty struct will be returned and we can use it.
	dsState, _ := states.GetDeliveryService(dsName)
	dsState.IsAvailable = stat.CommonStats.IsAvailable.Value
	states.SetDeliveryService(dsName, dsState) // TODO sync.Map? Determine if slow.

	getEvent := func(desc string) health.Event {
		// TODO sync.Pool?
		return health.Event{
			Time:        health.Time(time.Now()),
			Description: desc,
			Name:        dsName.String(),
			Hostname:    dsName.String(),
			Type:        health.DeliveryServiceEventType,
			Available:   stat.CommonStats.IsAvailable.Value,
		}
	}
	if stat.CommonStats.IsAvailable.Value == false && lastStat.Available == true {
		eventDesc := "Unavailable"
		if dsErr != nil {
			eventDesc = eventDesc + " err: " + dsErr.Error()
		}
		events.Add(getEvent(eventDesc)) // TODO change events.Add to not allocate new memory, after the limit is reached.
	} else if stat.CommonStats.IsAvailable.Value == true && lastStat.Available == false {
		events.Add(getEvent("Available caches"))
	}

	lastStat.Available = stat.CommonStats.IsAvailable.Value
}

// latestBytes returns the most recent OutBytes from the given cache results,
// and the time of that result. It assumes zero results are not valid, but
// nonzero results with errors are valid.
func latestBytes(p cache.PrecomputedData) (uint64, time.Time, error) {
	if p.OutBytes == 0 {
		return 0, time.Time{}, fmt.Errorf("no valid results")
	}
	return p.OutBytes, p.Time, nil
}

// addCachePerSecStats calculates the cache per-second stats, adds them to LastStats.
// Note this mutates lastStats, adding the calculated per-second stats to it.
func addCachePerSecStats(lastStats *dsdata.LastStats, cacheName tc.CacheName, precomputed cache.PrecomputedData) {
	outBytes, outBytesTime, err := latestBytes(precomputed) // it's ok if `latestBytes` returns 0s with an error, `addLastStat` will refrain from setting it (unless the previous calculation was nonzero, in which case it will error appropriately).
	if err != nil {
		log.Warnf("while computing delivery service data for cache %v: %v\n", cacheName, err)
	}
	lastStat, ok := lastStats.Caches[cacheName] // if lastStats.Caches[cacheName] doesn't exist, it will be zero-constructed, and `addLastStat` will refrain from setting the PerSec for zero LastStats
	if !ok {
		lastStat = &dsdata.LastStatsData{}
		lastStats.Caches[cacheName] = lastStat
	}
	if err = addLastStat(&lastStat.Bytes, int64(outBytes), outBytesTime); err != nil {
		log.Warnf("while computing delivery service data for cache %v: %v\n", cacheName, err)
	}
}

// addPerSecStats adds Kbps fields to the NewStats, based on the previous
// out_bytes in the oldStats, and the time difference.
//
// Traffic Server only updates its data every N seconds. So, often we get a new
// Stats with the same OutBytes as the previous one, so, we must record the last
// changed value, and the time it changed. Then, if the new OutBytes is
// different from the previous, we set the (new - old) / lastChangedTime as the
// KBPS, and update the recorded LastChangedTime and LastChangedValue
//
// TODO handle ATS byte rolling (when the `out_bytes` overflows back to 0)
//
// Note this mutates both dsStats and lastStats, adding the per-second stats to
// them.
func addPerSecStats(precomputed map[tc.CacheName]cache.PrecomputedData, dsStats *dsdata.Stats, lastStats *dsdata.LastStats, serverCachegroups map[tc.CacheName]tc.CacheGroupName, serverTypes map[tc.CacheName]tc.CacheType, mc tc.TrafficMonitorConfigMap, events health.ThreadsafeEvents, states peer.CRStatesThreadsafe) {
	for dsName, stat := range dsStats.DeliveryService {
		addDSPerSecStats(lastStats, dsName, stat, serverCachegroups, serverTypes, mc, events, precomputed, states)
	}
	for cacheName, precomputedData := range precomputed {
		addCachePerSecStats(lastStats, cacheName, precomputedData)
	}
}

// CreateStats aggregates and creates statistics from given precomputed stat history. It returns the created stats, information about these stats necessary for the next calculation, and any error.
// Note lastStats is mutated, being set with the new last stats.
func CreateStats(precomputed map[tc.CacheName]cache.PrecomputedData, toData todata.TOData, crStates tc.CRStates, lastStats *dsdata.LastStats, mc tc.TrafficMonitorConfigMap, events health.ThreadsafeEvents, states peer.CRStatesThreadsafe) *dsdata.Stats {
	start := time.Now()
	dsStats := dsdata.NewStats(len(toData.DeliveryServiceServers)) // TODO sync.Pool?
	for deliveryService := range toData.DeliveryServiceServers {
		if deliveryService == "" {
			log.Errorf("EMPTY CreateStats deliveryService")
			continue
		}
		dsStats.DeliveryService[deliveryService] = dsdata.NewStat() // TODO sync.Pool?
	}
	setStaticData(dsStats, toData.DeliveryServiceServers)
	addAvailableData(dsStats, crStates, toData.ServerCachegroups, toData.ServerDeliveryServices, toData.ServerTypes, precomputed) // TODO move after stat summarisation

	for server, precomputedData := range precomputed {
		cachegroup, ok := toData.ServerCachegroups[server]
		if !ok {
			// this can happen if we have precomputed data for a cache but the cache has since been deleted from Traffic Ops
			log.Infof("server %s has no cachegroup, skipping", server)
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

			httpDsStat, hadHttpDsStat := dsStats.DeliveryService[tc.DeliveryServiceName(ds)]
			if !hadHttpDsStat {
				httpDsStat = dsdata.NewStat() // TODO sync.Pool?
				dsStats.DeliveryService[tc.DeliveryServiceName(ds)] = httpDsStat
			}

			httpDsStatCg := httpDsStat.CacheGroups[cachegroup]
			if httpDsStatCg == nil {
				httpDsStatCg = &dsdata.StatCacheStats{}
				httpDsStat.CacheGroups[cachegroup] = httpDsStatCg
			}

			httpDsStatType := httpDsStat.Types[serverType]
			if httpDsStatType == nil {
				httpDsStatType = &dsdata.StatCacheStats{}
				httpDsStat.Types[serverType] = httpDsStatType
			}

			httpDsStatCache := httpDsStat.Caches[server]
			if httpDsStatCache == nil {
				httpDsStatCache = &dsdata.StatCacheStats{}
				httpDsStat.Caches[server] = httpDsStatCache
			}

			SumDSAstats(&httpDsStat.TotalStats, resultStat)
			SumDSAstats(httpDsStatCg, resultStat)
			SumDSAstats(httpDsStatType, resultStat)
			SumDSAstats(httpDsStatCache, resultStat)
			httpDsStat.CommonStats = dsStats.DeliveryService[tc.DeliveryServiceName(ds)].CommonStats // TODO verify whether this should be a sum
		}
	}

	addPerSecStats(precomputed, dsStats, lastStats, toData.ServerCachegroups, toData.ServerTypes, mc, events, states)
	log.Infof("CreateStats took %v\n", time.Since(start))
	dsStats.Time = time.Now()
	return dsStats
}

func getDSErr(dsName tc.DeliveryServiceName, dsStats dsdata.StatCacheStats, monitorConfig tc.TrafficMonitorConfigMap) error {
	if tpsThreshold := monitorConfig.DeliveryService[dsName.String()].TotalTPSThreshold; tpsThreshold > 0 && dsStats.TpsTotal.Value > float64(tpsThreshold) {
		return fmt.Errorf("total.tps_total too high (%.2f > %v)", dsStats.TpsTotal.Value, tpsThreshold)
	}
	if kbpsThreshold := monitorConfig.DeliveryService[dsName.String()].TotalKbpsThreshold; kbpsThreshold > 0 && dsStats.Kbps.Value > float64(kbpsThreshold) {
		return fmt.Errorf("total.kbps too high (%.2f > %v)", dsStats.Kbps.Value, kbpsThreshold)
	}
	return nil
}

func SumDSAstats(ds *dsdata.StatCacheStats, cacheStat *cache.DSStat) {
	ds.OutBytes.Value += int64(cacheStat.OutBytes)
	ds.InBytes.Value += float64(cacheStat.InBytes)
	ds.Status2xx.Value += int64(cacheStat.Status2xx)
	ds.Status3xx.Value += int64(cacheStat.Status3xx)
	ds.Status4xx.Value += int64(cacheStat.Status4xx)
	ds.Status5xx.Value += int64(cacheStat.Status5xx)
}
