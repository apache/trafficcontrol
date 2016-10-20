package manager

import (
	"sync"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/cache"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/config"
	ds "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
)

type StatHistory map[enum.CacheName][]cache.Result

func copyStat(a []cache.Result) []cache.Result {
	b := make([]cache.Result, len(a), len(a))
	copy(b, a)
	return b
}

func (a StatHistory) Copy() StatHistory {
	b := StatHistory{}
	for k, v := range a {
		b[k] = copyStat(v)
	}
	return b
}

// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
type StatHistoryThreadsafe struct {
	statHistory *StatHistory
	m           *sync.RWMutex
	max         uint64
}

func (h StatHistoryThreadsafe) Max() uint64 {
	return h.max
}

func NewStatHistoryThreadsafe(maxHistory uint64) StatHistoryThreadsafe {
	h := StatHistory{}
	return StatHistoryThreadsafe{m: &sync.RWMutex{}, statHistory: &h, max: maxHistory}
}

// Get returns the StatHistory. Callers MUST NOT modify. If mutation is necessary, call StatHistory.Copy()
func (t *StatHistoryThreadsafe) Get() StatHistory {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.statHistory
}

// Set sets the internal StatHistory. This is only safe for one thread of execution. This MUST NOT be called from multiple threads.
func (t *StatHistoryThreadsafe) Set(v StatHistory) {
	t.m.Lock()
	*t.statHistory = v
	t.m.Unlock()
}

func pruneHistory(history []cache.Result, limit uint64) []cache.Result {
	if uint64(len(history)) > limit {
		history = history[1:]
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
	errorCount UintThreadsafe,
	cfg config.Config,
	monitorConfig TrafficMonitorConfigMapThreadsafe,
) (StatHistoryThreadsafe, DurationMapThreadsafe, LastStatsThreadsafe, DSStatsThreadsafe, UnpolledCachesThreadsafe) {
	statHistory := NewStatHistoryThreadsafe(cfg.MaxStatHistory)
	lastStatDurations := NewDurationMapThreadsafe()
	lastStatEndTimes := map[enum.CacheName]time.Time{}
	lastStats := NewLastStatsThreadsafe()
	dsStats := NewDSStatsThreadsafe()
	unpolledCaches := NewUnpolledCachesThreadsafe()
	tickInterval := cfg.StatFlushInterval
	go func() {

		<-cachesChanged // wait for the signal that localStates have been set
		unpolledCaches.SetNewCaches(getNewCaches(localStates, monitorConfig))

		for {
			var results []cache.Result
			results = append(results, <-cacheStatChan)
			tick := time.Tick(tickInterval)
		innerLoop:
			for {
				select {
				case <-cachesChanged:
					unpolledCaches.SetNewCaches(getNewCaches(localStates, monitorConfig))
				case <-tick:
					log.Warnf("StatHistoryManager flushing queued results\n")
					processStatResults(results, statHistory, combinedStates.Get(), lastStats, toData.Get(), errorCount, dsStats, lastStatEndTimes, lastStatDurations, unpolledCaches)
					break innerLoop
				default:
					select {
					case r := <-cacheStatChan:
						results = append(results, r)
					default:
						processStatResults(results, statHistory, combinedStates.Get(), lastStats, toData.Get(), errorCount, dsStats, lastStatEndTimes, lastStatDurations, unpolledCaches)
						break innerLoop
					}
				}
			}
		}
	}()
	return statHistory, lastStatDurations, lastStats, dsStats, unpolledCaches
}

// processStatResults processes the given results, creating and setting DSStats, LastStats, and other stats. Note this is NOT threadsafe, and MUST NOT be called from multiple threads.
func processStatResults(
	results []cache.Result,
	statHistoryThreadsafe StatHistoryThreadsafe,
	combinedStates peer.Crstates,
	lastStats LastStatsThreadsafe,
	toData todata.TOData,
	errorCount UintThreadsafe,
	dsStats DSStatsThreadsafe,
	lastStatEndTimes map[enum.CacheName]time.Time,
	lastStatDurationsThreadsafe DurationMapThreadsafe,
	unpolledCaches UnpolledCachesThreadsafe,
) {
	statHistory := statHistoryThreadsafe.Get().Copy()
	maxStats := statHistoryThreadsafe.Max()
	for _, result := range results {
		// TODO determine if we want to add results with errors, or just print the errors now and don't add them.
		statHistory[enum.CacheName(result.Id)] = pruneHistory(append(statHistory[enum.CacheName(result.Id)], result), maxStats)
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
		if lastStatStart, ok := lastStatEndTimes[enum.CacheName(result.Id)]; ok {
			d := time.Since(lastStatStart)
			lastStatDurations[enum.CacheName(result.Id)] = d
		}
		lastStatEndTimes[enum.CacheName(result.Id)] = endTime

		// log.Debugf("poll %v %v statfinish\n", result.PollID, endTime)
		result.PollFinished <- result.PollID
	}
	lastStatDurationsThreadsafe.Set(lastStatDurations)
	unpolledCaches.SetPolled(results, lastStats)
}
