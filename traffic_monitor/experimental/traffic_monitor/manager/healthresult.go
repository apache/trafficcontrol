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
	"sync"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/cache"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/config"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/health"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/threadsafe"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
)

// DurationMap represents a map of cache names to durations
type DurationMap map[enum.CacheName]time.Duration

// DurationMapThreadsafe wraps a DurationMap in an object safe for a single writer and multiple readers
type DurationMapThreadsafe struct {
	durationMap *DurationMap
	m           *sync.RWMutex
}

// Copy copies this duration map.
func (a DurationMap) Copy() DurationMap {
	b := DurationMap{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

// NewDurationMapThreadsafe returns a new DurationMapThreadsafe safe for multiple readers and a single writer goroutine.
func NewDurationMapThreadsafe() DurationMapThreadsafe {
	m := DurationMap{}
	return DurationMapThreadsafe{m: &sync.RWMutex{}, durationMap: &m}
}

// Get returns the duration map. Callers MUST NOT mutate. If mutation is necessary, call DurationMap.Copy().
func (o *DurationMapThreadsafe) Get() DurationMap {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.durationMap
}

// Set sets the internal duration map. This MUST NOT be called by multiple goroutines.
func (o *DurationMapThreadsafe) Set(d DurationMap) {
	o.m.Lock()
	*o.durationMap = d
	o.m.Unlock()
}

// StartHealthResultManager starts the goroutine which listens for health results.
// Note this polls the brief stat endpoint from ATS Astats, not the full stats.
// This poll should be quicker and less computationally expensive for ATS, but
// doesn't include all stat data needed for e.g. delivery service calculations.4
// Returns the last health durations, events, the local cache statuses, and the health result history.
func StartHealthResultManager(
	cacheHealthChan <-chan cache.Result,
	toData todata.TODataThreadsafe,
	localStates peer.CRStatesThreadsafe,
	monitorConfig TrafficMonitorConfigMapThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	fetchCount threadsafe.Uint,
	errorCount threadsafe.Uint,
	cfg config.Config,
) (DurationMapThreadsafe, threadsafe.Events, threadsafe.CacheAvailableStatus, threadsafe.ResultHistory) {
	lastHealthDurations := NewDurationMapThreadsafe()
	events := threadsafe.NewEvents(cfg.MaxEvents)
	localCacheStatus := threadsafe.NewCacheAvailableStatus()
	healthHistory := threadsafe.NewResultHistory()
	go healthResultManagerListen(
		cacheHealthChan,
		toData,
		localStates,
		lastHealthDurations,
		healthHistory,
		monitorConfig,
		peerStates,
		combinedStates,
		fetchCount,
		errorCount,
		events,
		localCacheStatus,
		cfg,
	)
	return lastHealthDurations, events, localCacheStatus, healthHistory
}

func healthResultManagerListen(
	cacheHealthChan <-chan cache.Result,
	toData todata.TODataThreadsafe,
	localStates peer.CRStatesThreadsafe,
	lastHealthDurations DurationMapThreadsafe,
	healthHistory threadsafe.ResultHistory,
	monitorConfig TrafficMonitorConfigMapThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	fetchCount threadsafe.Uint,
	errorCount threadsafe.Uint,
	events threadsafe.Events,
	localCacheStatus threadsafe.CacheAvailableStatus,
	cfg config.Config,
) {
	lastHealthEndTimes := map[enum.CacheName]time.Time{}
	// This reads at least 1 value from the cacheHealthChan. Then, we loop, and try to read from the channel some more. If there's nothing to read, we hit `default` and process. If there is stuff to read, we read it, then inner-loop trying to read more. If we're continuously reading and the channel is never empty, and we hit the tick time, process anyway even though the channel isn't empty, to prevent never processing (starvation).
	var ticker *time.Ticker
	for {
		var results []cache.Result
		results = append(results, <-cacheHealthChan)
		if ticker != nil {
			ticker.Stop()
		}
		ticker = time.NewTicker(cfg.HealthFlushInterval)
	innerLoop:
		for {
			select {
			case <-ticker.C:
				log.Warnf("Health Result Manager flushing queued results\n")
				processHealthResult(
					cacheHealthChan,
					toData,
					localStates,
					lastHealthDurations,
					monitorConfig,
					peerStates,
					combinedStates,
					fetchCount,
					errorCount,
					events,
					localCacheStatus,
					lastHealthEndTimes,
					healthHistory,
					results,
					cfg,
				)
				break innerLoop
			default:
				select {
				case r := <-cacheHealthChan:
					results = append(results, r)
				default:
					processHealthResult(
						cacheHealthChan,
						toData,
						localStates,
						lastHealthDurations,
						monitorConfig,
						peerStates,
						combinedStates,
						fetchCount,
						errorCount,
						events,
						localCacheStatus,
						lastHealthEndTimes,
						healthHistory,
						results,
						cfg,
					)
					break innerLoop
				}
			}
		}
	}
}

// processHealthResult processes the given health results, adding their stats to the CacheAvailableStatus. Note this is NOT threadsafe, because it non-atomically gets CacheAvailableStatuses, Events, LastHealthDurations and later updates them. This MUST NOT be called from multiple threads.
func processHealthResult(
	cacheHealthChan <-chan cache.Result,
	toData todata.TODataThreadsafe,
	localStates peer.CRStatesThreadsafe,
	lastHealthDurationsThreadsafe DurationMapThreadsafe,
	monitorConfig TrafficMonitorConfigMapThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	fetchCount threadsafe.Uint,
	errorCount threadsafe.Uint,
	events threadsafe.Events,
	localCacheStatusThreadsafe threadsafe.CacheAvailableStatus,
	lastHealthEndTimes map[enum.CacheName]time.Time,
	healthHistory threadsafe.ResultHistory,
	results []cache.Result,
	cfg config.Config,
) {
	if len(results) == 0 {
		return
	}
	toDataCopy := toData.Get() // create a copy, so the same data used for all processing of this cache health result
	localCacheStatus := localCacheStatusThreadsafe.Get().Copy()
	monitorConfigCopy := monitorConfig.Get()
	healthHistoryCopy := healthHistory.Get().Copy()
	for _, healthResult := range results {
		log.Debugf("poll %v %v healthresultman start\n", healthResult.PollID, time.Now())
		fetchCount.Inc()
		var prevResult cache.Result
		healthResultHistory := healthHistoryCopy[healthResult.ID]
		if len(healthResultHistory) != 0 {
			prevResult = healthResultHistory[len(healthResultHistory)-1]
		}

		if healthResult.Error == nil {
			health.GetVitals(&healthResult, &prevResult, &monitorConfigCopy)
		}

		maxHistory := uint64(monitorConfigCopy.Profile[monitorConfigCopy.TrafficServer[string(healthResult.ID)].Profile].Parameters.HistoryCount)
		if maxHistory < 1 {
			log.Warnf("processHealthResult got history count %v for %v, setting to 1\n", maxHistory, healthResult.ID)
			maxHistory = 1
		}

		healthHistoryCopy[healthResult.ID] = pruneHistory(append([]cache.Result{healthResult}, healthHistoryCopy[healthResult.ID]...), maxHistory)

		isAvailable, whyAvailable := health.EvalCache(healthResult, &monitorConfigCopy)
		if available, ok := localStates.GetCache(healthResult.ID); !ok || available.IsAvailable != isAvailable {
			log.Infof("Changing state for %s was: %t now: %t because %s error: %v", healthResult.ID, prevResult.Available, isAvailable, whyAvailable, healthResult.Error)
			events.Add(cache.Event{Time: time.Now().Unix(), Description: whyAvailable, Name: healthResult.ID, Hostname: healthResult.ID, Type: toDataCopy.ServerTypes[healthResult.ID].String(), Available: isAvailable})
		}

		localCacheStatus[healthResult.ID] = cache.AvailableStatus{
			Available: isAvailable,
			Status:    monitorConfigCopy.TrafficServer[string(healthResult.ID)].Status,
			Why:       whyAvailable,
		} // TODO move within localStates?
		localStates.SetCache(healthResult.ID, peer.IsAvailable{IsAvailable: isAvailable})
	}
	calculateDeliveryServiceState(toDataCopy.DeliveryServiceServers, localStates)
	healthHistory.Set(healthHistoryCopy)
	localCacheStatusThreadsafe.Set(localCacheStatus)
	// TODO determine if we should combineCrStates() here

	lastHealthDurations := lastHealthDurationsThreadsafe.Get().Copy()
	for _, healthResult := range results {
		if lastHealthStart, ok := lastHealthEndTimes[healthResult.ID]; ok {
			d := time.Since(lastHealthStart)
			lastHealthDurations[healthResult.ID] = d
		}
		lastHealthEndTimes[healthResult.ID] = time.Now()

		log.Debugf("poll %v %v finish\n", healthResult.PollID, time.Now())
		healthResult.PollFinished <- healthResult.PollID
	}
	lastHealthDurationsThreadsafe.Set(lastHealthDurations)
}

// calculateDeliveryServiceState calculates the state of delivery services from the new cache state data `cacheState` and the CRConfig data `deliveryServiceServers` and puts the calculated state in the outparam `deliveryServiceStates`
func calculateDeliveryServiceState(deliveryServiceServers map[enum.DeliveryServiceName][]enum.CacheName, states peer.CRStatesThreadsafe) {
	deliveryServices := states.GetDeliveryServices()
	for deliveryServiceName, deliveryServiceState := range deliveryServices {
		if _, ok := deliveryServiceServers[deliveryServiceName]; !ok {
			// log.Errorf("CRConfig does not have delivery service %s, but traffic monitor poller does; skipping\n", deliveryServiceName)
			continue
		}
		deliveryServiceState.IsAvailable = false
		deliveryServiceState.DisabledLocations = []enum.CacheName{} // it's important this isn't nil, so it serialises to the JSON `[]` instead of `null`
		for _, server := range deliveryServiceServers[deliveryServiceName] {
			if available, _ := states.GetCache(server); available.IsAvailable {
				deliveryServiceState.IsAvailable = true
			} else {
				deliveryServiceState.DisabledLocations = append(deliveryServiceState.DisabledLocations, server)
			}
		}
		states.SetDeliveryService(deliveryServiceName, deliveryServiceState)
	}
}
