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
	"os"
	"runtime"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/cache"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/config"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/ds"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/health"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/peer"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/threadsafe"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

func pruneHistory(history []cache.Result, limit uint64) []cache.Result {
	if uint64(len(history)) > limit {
		history = history[:limit-1]
	}
	return history
}

func getNewCaches(localStates peer.CRStatesThreadsafe, monitorConfigTS threadsafe.TrafficMonitorConfigMap) map[tc.CacheName]bool {
	monitorConfig := monitorConfigTS.Get()
	caches := map[tc.CacheName]bool{}
	for cacheName, a := range localStates.GetCaches() {
		// ONLINE and OFFLINE caches are not polled.
		if ts, ok := monitorConfig.TrafficServer[string(cacheName)]; !ok || ts.ServerStatus == string(tc.CacheStatusOnline) || ts.ServerStatus == string(tc.CacheStatusOffline) {
			continue
		}
		caches[cacheName] = a.DirectlyPolled
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
	cfg config.Config,
	monitorConfig threadsafe.TrafficMonitorConfigMap,
	events health.ThreadsafeEvents,
	combineState func(),
) (threadsafe.ResultInfoHistory, threadsafe.ResultStatHistory, threadsafe.CacheKbpses, threadsafe.DurationMap, threadsafe.LastStats, threadsafe.DSStatsReader, threadsafe.UnpolledCaches, threadsafe.CacheAvailableStatus) {
	statInfoHistory := threadsafe.NewResultInfoHistory()
	statResultHistory := threadsafe.NewResultStatHistory()
	statMaxKbpses := threadsafe.NewCacheKbpses()
	lastStatDurations := threadsafe.NewDurationMap()
	lastStatEndTimes := map[tc.CacheName]time.Time{}
	lastStats := threadsafe.NewLastStats()
	dsStats := threadsafe.NewDSStats()
	statUnpolledCaches := threadsafe.NewUnpolledCaches()
	localCacheStatus := threadsafe.NewCacheAvailableStatus()

	precomputedData := map[tc.CacheName]cache.PrecomputedData{}

	lastResults := map[tc.CacheName]cache.Result{}

	haveCachesChanged := func() bool {
		select {
		case <-cachesChanged:
			return true
		default:
			return false
		}
	}

	process := func(results []cache.Result) {
		if haveCachesChanged() {
			statUnpolledCaches.SetNewCaches(getNewCaches(localStates, monitorConfig))
		}
		processStatResults(results, statInfoHistory, statResultHistory, statMaxKbpses, combinedStates, lastStats, toData.Get(), dsStats, lastStatEndTimes, lastStatDurations, statUnpolledCaches, monitorConfig.Get(), precomputedData, lastResults, localStates, events, localCacheStatus, combineState, cfg.CachePollingProtocol)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("StatHistoryManager panic: %v\n", err)
			} else {
				log.Errorln("StatHistoryManager failed without panic")
			}
			log.Errorf("%s\n", stacktrace())
			os.Exit(1) // The monitor can't run without a stat processor
		}()

		flushTimer := time.NewTimer(cfg.StatFlushInterval)
		// Note! bufferTimer MAY be uninitialized! If there is no cfg.StatBufferInterval, the timer WILL NOT be created with time.NewTimer(), and thus is NOT initialized, and MUST NOT have functions called, such as timer.Stop()! Those functions WILL panic.
		bufferTimer := &time.Timer{}
		bufferFakeChan := make(chan time.Time, 1) // fake chan, if there is no stat buffer interval. Unused, if cfg.StatBufferInterval != nil. Buffer 1, so don't need a separate goroutine to write.
		if cfg.StatBufferInterval == 0 {
			// if there is no stat buffer interval, make a timer which has already expired.
			bufferFakeChan <- time.Now()
			bufferTimer.C = bufferFakeChan
		} else {
			bufferTimer = time.NewTimer(cfg.StatBufferInterval)
		}

		// resetBufferTimer resets the Buffer timer. It MUST have expired and been read.
		// If the buffer loop is changed to allow finishing without being expired and read, this MUST be changed to stop and drain the channel (with a select/default, if it's possible to expire but not be read (like flush is now). Otherwise, it will deadlock and/or leak resources.
		resetBufferTimer := func() {
			if cfg.StatBufferInterval == 0 {
				bufferFakeChan <- time.Now()
			} else {
				bufferTimer.Reset(cfg.StatBufferInterval)
			}
		}

		// resetFlushTimer resets the Flush timer. It may or may not have been read or expired.
		resetFlushTimer := func() {
			if !flushTimer.Stop() {
				select { // need to select/default because we don't know whether the flush timer was read
				case <-flushTimer.C:
				default:
				}
			}
			flushTimer.Reset(cfg.StatFlushInterval)
		}

		// There are 2 timers: the Buffer, and the Flush.
		// The Buffer says "never process until this much time has elapsed"
		// The Flush says "if we're continuously getting new stats, with no break, and this much time has elasped, go ahead and process anyway to prevent starvation"
		//
		// So, we continuously read from the stat channel, until Buffer has elasped. Even if the channel is empty, wait and keep trying to read.
		// Then, we move to State 2: continuously read from the stat channel, while there are things to read. If at any point there's nothing more to read, then process. Otherwise, if there are always thing to read, then after Flush time has elapsed, then go ahead and read anyway, and go to State 1.
		//
		// Note that either the Buffer or Flush may be configured to be 0.
		// If the Buffer is 0, we immediately move to phase 2: process as fast as we can, only flush to prevent starvation. This optimizes the Monitor for getting health as quickly as possible, at the cost of CPU. (Having a buffer itself puts CPU above getting health results quickly, and the buffer interval is a factor of that)
		// If the Flush is 0, then the Monitor will process every Buffer interval, regardless whether results are still coming in. This attempts to optimize for stability, attempting to ensure a poll every (Buffer + Poll Time) interval. Note this attempt may fail, and in particular, if the Monitor is unable to keep up with the given poll time and buffer, it will continuously back up. For this reason, setting a Flush of 0 is not recommended.
		//
		// Note the Flush and Buffer times are cumulative. That is, the total "maximum time a cache can be unhealthy before we know" is the Poll+Flush+Buffer. So, the buffer time elapses, then we start a new flush interval. They don't run concurrently.

		results := []cache.Result{}

		// flush loop - breaks after processing - processes when there are no pending results, or the flush time elapses.
		flushLoop := func() {
			log.Infof("StatHistoryManager starting flushLoop with %+v results\n", len(results))
			resetFlushTimer()
			for {
				select {
				case <-flushTimer.C: // first, make sure the flushTimer hasn't expired, by itself (because GO selects aren't ordered, so it could starve if we were reading <-cacheStatChan at the same level
					log.Infof("StatHistoryManager flushLoop: flush timer fired, processing %+v results\n", len(results))
					process(results)
					return
				default: // flushTimer hadn't expired: read cacheStatChan at the same level now.
					// This extra level is necessary, because Go selects aren't ordered, so even after the Flush timer expires, the "case" could still never get hit,
					// if there were continuously results from <-cacheStatChan at the same level.
					select {
					case r := <-cacheStatChan:
						results = append(results, r)
						// we're still processing as much as possible, and flushing, don't break to the outer Buffer loop, until we process.
					default:
						log.Infof("StatHistoryManager flushLoop: stat chan is empty, processing %+v results\n", len(results))
						// Buffer expired (above), and the cacheStatChan is empty, so process
						process(results)
						return
					}
				}
			}
		}

		results = []cache.Result{}
		// no point doing anything, until we read at least one stat. If stat polling is disabled, this blocks forever
		results = append(results, <-cacheStatChan)

		// buffer loop - never breaks - calls flushLoop to actually process, when the buffer time elapses.
		for {
			// select only the bufferTimer first, to make sure it doesn't starve.
			select {
			case <-bufferTimer.C:
				// buffer expired, move to State 2 (Flush)
				flushLoop()
				log.Infof("StatHistoryManager bufferLoop exiting flush loop, resetting buffer timer\n")
				resetBufferTimer()
				results = []cache.Result{}
				results = append(results, <-cacheStatChan) // no point doing anything, until we read at least one stat.
			default:
				// buffer time hadn't elapsed, so we know we aren't starving. Go ahead and read the stat chan + buffer now.
				select {
				case r := <-cacheStatChan:
					results = append(results, r)
				case <-bufferTimer.C: // TODO protect against bufferTimer starvation
					// buffer expired, move to State 2 (Flush): process until there's nothing to process, or the Flush elapses.
					flushLoop()
					log.Infof("StatHistoryManager bufferLoop (within stat select) exiting flush loop, resetting buffer timer\n")
					resetBufferTimer()
					results = []cache.Result{}
					results = append(results, <-cacheStatChan) // no point doing anything, until we read at least one stat.
				}
			}
		}
	}()
	return statInfoHistory, statResultHistory, statMaxKbpses, lastStatDurations, lastStats, &dsStats, statUnpolledCaches, localCacheStatus
}

func stacktrace() []byte {
	initialBufSize := 1024
	buf := make([]byte, initialBufSize)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, len(buf)*2)
	}
}

// processStatResults processes the given results, creating and setting DSStats, LastStats, and other stats. Note this is NOT threadsafe, and MUST NOT be called from multiple threads.
func processStatResults(
	results []cache.Result,
	statInfoHistoryThreadsafe threadsafe.ResultInfoHistory,
	statResultHistoryThreadsafe threadsafe.ResultStatHistory,
	statMaxKbpsesThreadsafe threadsafe.CacheKbpses,
	combinedStatesThreadsafe peer.CRStatesThreadsafe,
	lastStats threadsafe.LastStats,
	toData todata.TOData,
	dsStats threadsafe.DSStats,
	lastStatEndTimes map[tc.CacheName]time.Time,
	lastStatDurationsThreadsafe threadsafe.DurationMap,
	statUnpolledCaches threadsafe.UnpolledCaches,
	mc tc.TrafficMonitorConfigMap,
	precomputedData map[tc.CacheName]cache.PrecomputedData,
	lastResults map[tc.CacheName]cache.Result,
	localStates peer.CRStatesThreadsafe,
	events health.ThreadsafeEvents,
	localCacheStatusThreadsafe threadsafe.CacheAvailableStatus,
	combineState func(),
	pollingProtocol config.PollingProtocol,
) {
	if len(results) == 0 {
		return
	}
	combinedStates := combinedStatesThreadsafe.Get()
	defer func() {
		for _, r := range results {
			r.PollFinished <- r.PollID
		}
	}()

	// setting the statHistory could be put in a goroutine concurrent with `ds.CreateStats`, if it were slow
	statInfoHistory := statInfoHistoryThreadsafe.Get().Copy()
	statMaxKbpses := statMaxKbpsesThreadsafe.Get().Copy()

	for i, result := range results {
		maxStats := uint64(mc.Profile[mc.TrafficServer[string(result.ID)].Profile].Parameters.HistoryCount)
		if maxStats < 1 {
			log.Infof("processStatResults got history count %v for %v, setting to 1\n", maxStats, result.ID)
			maxStats = 1
		}

		// TODO determine if we want to add results with errors, or just print the errors now and don't add them.
		if lastResult, ok := lastResults[tc.CacheName(result.ID)]; ok && result.Error == nil {
			health.GetVitals(&result, &lastResult, &mc) // TODO precompute
			if result.Error == nil {
				results[i] = result
			} else {
				log.Errorf("stat poll getting vitals for %v: %v\n", result.ID, result.Error)
			}
		}
		statInfoHistory.Add(result, maxStats)
		if err := statResultHistoryThreadsafe.Add(result, maxStats); err != nil {
			log.Errorf("Adding result from %v: %v\n", result.ID, err)
		}
		// Don't add errored maxes or precomputed DSStats
		if result.Error == nil {
			// max and precomputed always contain the latest result from each cache
			statMaxKbpses[result.ID] = uint64(result.PrecomputedData.MaxKbps)
			// if we failed to compute the OutBytes, keep the outbytes of the last result.
			if result.PrecomputedData.OutBytes == 0 {
				result.PrecomputedData.OutBytes = precomputedData[tc.CacheName(result.ID)].OutBytes
			}
			precomputedData[tc.CacheName(result.ID)] = result.PrecomputedData

		}
		lastResults[tc.CacheName(result.ID)] = result
	}
	statInfoHistoryThreadsafe.Set(statInfoHistory)
	statMaxKbpsesThreadsafe.Set(statMaxKbpses)

	lastStatsVal := lastStats.Get()
	lastStatsCopy := lastStatsVal.Copy()
	newDsStats := ds.CreateStats(precomputedData, toData, combinedStates, lastStatsCopy, mc, events, localStates)

	dsStats.Set(*newDsStats)
	lastStats.Set(*lastStatsCopy)

	pollerName := "stat"
	health.CalcAvailability(results, pollerName, &statResultHistoryThreadsafe, mc, toData, localCacheStatusThreadsafe, localStates, events, pollingProtocol)
	combineState()

	endTime := time.Now()
	lastStatDurations := threadsafe.CopyDurationMap(lastStatDurationsThreadsafe.Get())
	for _, result := range results {
		if lastStatStart, ok := lastStatEndTimes[tc.CacheName(result.ID)]; ok {
			d := time.Since(lastStatStart)
			lastStatDurations[tc.CacheName(result.ID)] = d
		}
		lastStatEndTimes[tc.CacheName(result.ID)] = endTime
	}
	lastStatDurationsThreadsafe.Set(lastStatDurations)
	statUnpolledCaches.SetPolled(results, lastStats.Get())
}
