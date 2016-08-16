package manager

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/health"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
)

type DurationMapThreadsafe struct {
	durationMap map[enum.CacheName]time.Duration // TODO change string -> CacheName
	m           *sync.Mutex
}

func copyDurationMap(a map[enum.CacheName]time.Duration) map[enum.CacheName]time.Duration {
	b := map[enum.CacheName]time.Duration{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

func NewDurationMapThreadsafe() DurationMapThreadsafe {
	return DurationMapThreadsafe{m: &sync.Mutex{}, durationMap: map[enum.CacheName]time.Duration{}}
}

func (o DurationMapThreadsafe) Get() map[enum.CacheName]time.Duration {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	return copyDurationMap(o.durationMap)
}

func (o DurationMapThreadsafe) GetDuration(cacheName enum.CacheName) (time.Duration, bool) {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	duration, ok := o.durationMap[cacheName]
	return duration, ok
}

func (o DurationMapThreadsafe) Set(cacheName enum.CacheName, d time.Duration) {
	o.m.Lock()
	o.durationMap[cacheName] = d
	o.m.Unlock()
}

type TimeMapThreadsafe struct {
	timeMap map[enum.CacheName]time.Time // TODO change string -> CacheName
	m       *sync.Mutex
}

func copyTimeMap(a map[enum.CacheName]time.Time) map[enum.CacheName]time.Time {
	b := map[enum.CacheName]time.Time{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

func NewTimeMapThreadsafe() TimeMapThreadsafe {
	return TimeMapThreadsafe{m: &sync.Mutex{}, timeMap: map[enum.CacheName]time.Time{}}
}

func (o TimeMapThreadsafe) Get() map[enum.CacheName]time.Time {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	return copyTimeMap(o.timeMap)
}

func (o TimeMapThreadsafe) GetTime(cacheName enum.CacheName) (time.Time, bool) {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	time, ok := o.timeMap[cacheName]
	return time, ok
}

func (o TimeMapThreadsafe) Set(cacheName enum.CacheName, d time.Time) {
	o.m.Lock()
	o.timeMap[cacheName] = d
	o.m.Unlock()
}

type ResultsThreadsafe struct {
	r map[enum.CacheName][]cache.Result // TODO change string -> CacheName
	m *sync.Mutex
}

func NewResultsThreadsafe() ResultsThreadsafe {
	return ResultsThreadsafe{m: &sync.Mutex{}, r: map[enum.CacheName][]cache.Result{}}
}

func (o ResultsThreadsafe) Get(cacheName enum.CacheName) []cache.Result {
	o.m.Lock()
	defer func() {
		o.m.Unlock()
	}()
	return o.r[cacheName]
}

func (o ResultsThreadsafe) Set(cacheName enum.CacheName, newR []cache.Result) {
	o.m.Lock()
	o.r[cacheName] = newR
	o.m.Unlock()
}

// StartHealthResultManager starts the goroutine which listens for health results.
// Returns the last health durations, events, and the local cache statuses.
func StartHealthResultManager(cacheHealthChan <-chan cache.Result, toData todata.TODataThreadsafe, localStates CRStatesThreadsafe, statHistory StatHistoryThreadsafe, monitorConfig TrafficMonitorConfigMapThreadsafe, peerStates CRStatesPeersThreadsafe, combinedStates CRStatesThreadsafe, fetchCount UintThreadsafe, errorCount UintThreadsafe) (DurationMapThreadsafe, EventsThreadsafe, CacheAvailableStatusThreadsafe, DSStatsThreadsafe, StatsLastKbpsThreadsafe) {
	lastHealthDurations := NewDurationMapThreadsafe()
	events := NewEventsThreadsafe()
	localCacheStatus := NewCacheAvailableStatusThreadsafe()
	dsStats := NewDSStatsThreadsafe()
	lastKbpsStats := NewStatsLastKbpsThreadsafe()
	go healthResultManagerListen(cacheHealthChan, toData, localStates, lastHealthDurations, statHistory, monitorConfig, peerStates, combinedStates, fetchCount, errorCount, events, localCacheStatus, dsStats, lastKbpsStats)
	return lastHealthDurations, events, localCacheStatus, dsStats, lastKbpsStats
}

func healthResultManagerListen(cacheHealthChan <-chan cache.Result, toData todata.TODataThreadsafe, localStates CRStatesThreadsafe, lastHealthDurations DurationMapThreadsafe, statHistory StatHistoryThreadsafe, monitorConfig TrafficMonitorConfigMapThreadsafe, peerStates CRStatesPeersThreadsafe, combinedStates CRStatesThreadsafe, fetchCount UintThreadsafe, errorCount UintThreadsafe, events EventsThreadsafe, localCacheStatus CacheAvailableStatusThreadsafe, dsStats DSStatsThreadsafe, lastKbpsStats StatsLastKbpsThreadsafe) {
	lastHealthEndTimes := map[enum.CacheName]time.Time{}
	healthHistory := map[enum.CacheName][]cache.Result{}
	eventIndex := uint64(0) // TODO move to EventsThreadsafe.Add() ?
	// lastHealthEndTimes := NewTimeMapThreadsafe()
	// healthHistory := NewResultsThreadsafe()
	// eventIndex := NewUintThreadsafe()

	for {
		select {
		case healthResult := <-cacheHealthChan:
			// go func() {
			{
				fetchCount.Inc()
				toDataCopy := toData.Get() // create a copy, so the same data used for all processing of this cache health result
				var prevResult cache.Result
				healthResultHistory := healthHistory[enum.CacheName(healthResult.Id)]
				if len(healthResultHistory) != 0 {
					prevResult = healthResultHistory[len(healthResultHistory)-1]
				}

				monitorConfigCopy := monitorConfig.Get() // copy now, so all calculations are on the same data
				health.GetVitals(&healthResult, &prevResult, &monitorConfigCopy)
				healthHistory[enum.CacheName(healthResult.Id)] = pruneHistory(append(healthHistory[enum.CacheName(healthResult.Id)], healthResult), defaultMaxHistory)
				isAvailable, whyAvailable := health.EvalCache(healthResult, &monitorConfigCopy)
				if localStates.Get().Caches[healthResult.Id].IsAvailable != isAvailable {
					fmt.Println("Changing state for", healthResult.Id, " was:", prevResult.Available, " is now:", isAvailable, " because:", whyAvailable, " errors:", healthResult.Errors)
					e := Event{Index: eventIndex, Time: time.Now().Unix(), Description: whyAvailable, Name: healthResult.Id, Hostname: healthResult.Id, Type: toDataCopy.ServerTypes[enum.CacheName(healthResult.Id)].String(), Available: isAvailable}
					eventIndex++
					events.Add(e)
				}

				localCacheStatus.Set(enum.CacheName(healthResult.Id), CacheAvailableStatus{Available: isAvailable, Status: monitorConfigCopy.TrafficServer[healthResult.Id].Status}) // TODO move within localStates
				localStates.SetCache(healthResult.Id, peer.IsAvailable{IsAvailable: isAvailable})
				calculateDeliveryServiceState(toDataCopy.DeliveryServiceServers, localStates)

				// TODO determine if we should combineCrStates() here

				now := time.Now()

				var err error
				createStatsCopyStatHistory := statHistory.Get()
				createStatsCopyCombinedStates := combinedStates.Get()
				createStatsCopyLastKbpsStats := lastKbpsStats.Get()
				newDsStats, newLastKbpsStats, err := ds.CreateStats(createStatsCopyStatHistory, toDataCopy, createStatsCopyCombinedStates, createStatsCopyLastKbpsStats, now)
				if err != nil {
					errorCount.Inc()
					log.Printf("ERROR getting deliveryservice: %v\n", err)
				} else {
					dsStats.Set(newDsStats)
					lastKbpsStats.Set(newLastKbpsStats)
				}

				//				if lastHealthStart, ok := lastHealthEndTimes.GetTime(enum.CacheName(healthResult.Id)); ok {
				if lastHealthStart, ok := lastHealthEndTimes[enum.CacheName(healthResult.Id)]; ok {
					d := time.Since(lastHealthStart)
					lastHealthDurations.Set(enum.CacheName(healthResult.Id), d)
				}
				lastHealthEndTimes[enum.CacheName(healthResult.Id)] = now
				//}()
			}
		}
	}
}

// calculateDeliveryServiceState calculates the state of delivery services from the new cache state data `cacheState` and the CRConfig data `deliveryServiceServers` and puts the calculated state in the outparam `deliveryServiceStates`
func calculateDeliveryServiceState(deliveryServiceServers map[string][]string, states CRStatesThreadsafe) {
	deliveryServices := states.GetDeliveryServices()
	for deliveryServiceName, deliveryServiceState := range deliveryServices {
		if _, ok := deliveryServiceServers[deliveryServiceName]; !ok {
			// log.Printf("ERROR CRConfig does not have delivery service %s, but traffic monitor poller does; skipping\n", deliveryServiceName)
			continue
		}
		deliveryServiceState.IsAvailable = false
		deliveryServiceState.DisabledLocations = nil
		for _, server := range deliveryServiceServers[deliveryServiceName] {
			if states.GetCache(server).IsAvailable {
				deliveryServiceState.IsAvailable = true
			} else {
				deliveryServiceState.DisabledLocations = append(deliveryServiceState.DisabledLocations, server)
			}
		}
		deliveryServices[deliveryServiceName] = deliveryServiceState
	}
	states.SetDeliveryServices(deliveryServices)
}

// intersection returns strings in both a and b.
// Note this modifies a and b. Specifically, it sorts them. If that isn't acceptable, pass copies of your real data.
func intersection(a []string, b []string) []string {
	sort.Strings(a)
	sort.Strings(b)
	var c []string
	for _, s := range a {
		i := sort.SearchStrings(b, s)
		if i < len(b) && b[i] == s {
			c = append(c, s)
		}
	}
	return c
}
