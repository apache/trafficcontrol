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
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/cache"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/config"
	ds "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func pruneHistory(history []cache.Result, limit uint64) []cache.Result {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}

func getNewCaches(localStates peer.CRStatesThreadsafe, monitorConfigTS TrafficMonitorConfigMapThreadsafe) map[enum.CacheName]struct{} {
	monitorConfig := monitorConfigTS.Get()
	caches := map[enum.CacheName]struct{}{}
	for cacheName := range localStates.GetCaches() {
		// ONLINE and OFFLINE caches are not polled.
		// TODO add a function IsPolled() which can be called by this and the monitorConfig func which sets the polling, to prevent updating in one place breaking the other.
		if ts, ok := monitorConfig.TrafficServer[string(cacheName)]; !ok || ts.Status == "ONLINE" || ts.Status == "OFFLINE" {
			continue
		}
		caches[cacheName] = struct{}{}
	}
	return caches
}

// StartStatHistoryManager fetches the full statistics data from ATS Astats. This includes everything needed for all calculations, such as Delivery Services. This is expensive, though, and may be hard on ATS, so it should poll less often.
// For a fast 'is it alive' poll, use the Health Result Manager poll.
// Returns the stat history, the duration between the stat poll for each cache, the last Kbps data, the calculated Delivery Service stats, and the unpolled caches list.
func StartStatHistoryManager(
	cacheStatChan <-chan cache.Result,
	localStates peer.CRStatesThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	toData todata.TODataThreadsafe,
	cachesChanged <-chan struct{},
	errorCount threadsafe.Uint,
	cfg config.Config,
	monitorConfig TrafficMonitorConfigMapThreadsafe,
) (threadsafe.ResultHistory, DurationMapThreadsafe, threadsafe.LastStats, threadsafe.DSStatsReader, threadsafe.UnpolledCaches) {
	statHistory := threadsafe.NewResultHistory()
	lastStatDurations := NewDurationMapThreadsafe()
	lastStatEndTimes := map[enum.CacheName]time.Time{}
	lastStats := threadsafe.NewLastStats()
	dsStats := threadsafe.NewDSStats()
	unpolledCaches := threadsafe.NewUnpolledCaches()
	tickInterval := cfg.StatFlushInterval
	go func() {
		var ticker *time.Ticker
		<-cachesChanged // wait for the signal that localStates have been set
		unpolledCaches.SetNewCaches(getNewCaches(localStates, monitorConfig))

		for {
			var results []cache.Result
			results = append(results, <-cacheStatChan)
			if ticker != nil {
				ticker.Stop()
			}
			ticker = time.NewTicker(tickInterval)
		innerLoop:
			for {
				select {
				case <-cachesChanged:
					unpolledCaches.SetNewCaches(getNewCaches(localStates, monitorConfig))
				case <-ticker.C:
					log.Warnf("StatHistoryManager flushing queued results\n")
					processStatResults(results, statHistory, combinedStates.Get(), lastStats, toData.Get(), errorCount, dsStats, lastStatEndTimes, lastStatDurations, unpolledCaches, monitorConfig.Get())
					break innerLoop
				default:
					select {
					case r := <-cacheStatChan:
						results = append(results, r)
					default:
						processStatResults(results, statHistory, combinedStates.Get(), lastStats, toData.Get(), errorCount, dsStats, lastStatEndTimes, lastStatDurations, unpolledCaches, monitorConfig.Get())
						break innerLoop
					}
				}
			}
		}
	}()
	return statHistory, lastStatDurations, lastStats, &dsStats, unpolledCaches
}

// processStatResults processes the given results, creating and setting DSStats, LastStats, and other stats. Note this is NOT threadsafe, and MUST NOT be called from multiple threads.
func processStatResults(
	results []cache.Result,
	statHistoryThreadsafe threadsafe.ResultHistory,
	combinedStates peer.Crstates,
	lastStats threadsafe.LastStats,
	toData todata.TOData,
	errorCount threadsafe.Uint,
	dsStats threadsafe.DSStats,
	lastStatEndTimes map[enum.CacheName]time.Time,
	lastStatDurationsThreadsafe DurationMapThreadsafe,
	unpolledCaches threadsafe.UnpolledCaches,
	mc to.TrafficMonitorConfigMap,
) {
	statHistory := statHistoryThreadsafe.Get().Copy()
	for _, result := range results {
		maxStats := uint64(mc.Profile[mc.TrafficServer[string(result.ID)].Profile].Parameters.HistoryCount)
		// TODO determine if we want to add results with errors, or just print the errors now and don't add them.
		statHistory[result.ID] = pruneHistory(append([]cache.Result{result}, statHistory[result.ID]...), maxStats)
	}
	statHistoryThreadsafe.Set(statHistory)

	for _, result := range results {
		log.Debugf("poll %v %v CreateStats start\n", result.PollID, time.Now())
	}

	newDsStats, newLastStats, err := ds.CreateStats(statHistory, toData, combinedStates, lastStats.Get().Copy(), time.Now())

	for _, result := range results {
		log.Debugf("poll %v %v CreateStats end\n", result.PollID, time.Now())
	}

	if err != nil {
		errorCount.Inc()
		log.Errorf("getting deliveryservice: %v\n", err)
	} else {
		dsStats.Set(newDsStats)
		lastStats.Set(newLastStats)
	}

	endTime := time.Now()
	lastStatDurations := lastStatDurationsThreadsafe.Get().Copy()
	for _, result := range results {
		if lastStatStart, ok := lastStatEndTimes[result.ID]; ok {
			d := time.Since(lastStatStart)
			lastStatDurations[result.ID] = d
		}
		lastStatEndTimes[result.ID] = endTime

		// log.Debugf("poll %v %v statfinish\n", result.PollID, endTime)
		result.PollFinished <- result.PollID
	}
	lastStatDurationsThreadsafe.Set(lastStatDurations)
	unpolledCaches.SetPolled(results, lastStats.Get())
}
