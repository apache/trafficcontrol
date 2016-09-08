package manager

import (
	"sync"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
)

//const maxHistory = (60 / pollingInterval) * 5
const defaultMaxHistory = 5 // TODO make config setting?

// This could be made lock-free, if the performance was necessary
// TODO add separate locks for Caches and Deliveryservice maps?
type StatHistoryThreadsafe struct {
	statHistory map[enum.CacheName][]cache.Result
	m           *sync.Mutex
}

func NewStatHistoryThreadsafe() StatHistoryThreadsafe {
	return StatHistoryThreadsafe{m: &sync.Mutex{}, statHistory: map[enum.CacheName][]cache.Result{}}
}

func (t *StatHistoryThreadsafe) GetStat(stat enum.CacheName) []cache.Result {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return copyStat(t.statHistory[stat])
}

func (t *StatHistoryThreadsafe) Get() map[enum.CacheName][]cache.Result {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return copyStats(t.statHistory)
}

func (t *StatHistoryThreadsafe) Add(stat cache.Result) {
	t.m.Lock()
	t.statHistory[enum.CacheName(stat.Id)] = pruneHistory(append(t.statHistory[enum.CacheName(stat.Id)], stat), defaultMaxHistory)
	t.m.Unlock()
}

func pruneHistory(history []cache.Result, limit int) []cache.Result {
	if len(history) > limit {
		history = history[1:]
	}
	return history
}

func copyStat(a []cache.Result) []cache.Result {
	b := make([]cache.Result, len(a), len(a))
	for i, v := range a {
		b[i] = v
	}
	return b
}

func copyStats(a map[enum.CacheName][]cache.Result) map[enum.CacheName][]cache.Result {
	b := map[enum.CacheName][]cache.Result{}
	for k, v := range a {
		b[k] = copyStat(v)
	}
	return b
}

// StartStatHistoryManager fetches the full statistics data from ATS Astats. This includes everything needed for all calculations, such as Delivery Services. This is expensive, though, and may be hard on ATS, so it should poll less often.
// For a fast 'is it alive' poll, use the Health Result Manager poll.
// Returns the stat history, the duration between the stat poll for each cache, the last Kbps data, and the calculated Delivery Service stats.
func StartStatHistoryManager(cacheStatChan <-chan cache.Result, combinedStates peer.CRStatesThreadsafe, toData todata.TODataThreadsafe, errorCount UintThreadsafe) (StatHistoryThreadsafe, DurationMapThreadsafe, StatsLastKbpsThreadsafe, DSStatsThreadsafe) {
	statHistory := NewStatHistoryThreadsafe()
	lastStatDurations := NewDurationMapThreadsafe()
	lastStatEndTimes := map[enum.CacheName]time.Time{}
	lastKbpsStats := NewStatsLastKbpsThreadsafe()
	dsStats := NewDSStatsThreadsafe()
	tickInterval := time.Millisecond * 200 // TODO make config setting
	go func() {
		for {
			var results []cache.Result
			results = append(results, <-cacheStatChan)
			tick := time.Tick(tickInterval)
		innerLoop:
			for {
				select {
				case <-tick:
					log.Warnf("StatHistoryManager flushing queued results\n")
					processStatResults(results, statHistory, combinedStates.Get(), lastKbpsStats, toData.Get(), errorCount, dsStats, lastStatEndTimes, lastStatDurations)
					break innerLoop
				default:
					select {
					case r := <-cacheStatChan:
						results = append(results, r)
					default:
						processStatResults(results, statHistory, combinedStates.Get(), lastKbpsStats, toData.Get(), errorCount, dsStats, lastStatEndTimes, lastStatDurations)
						break innerLoop
					}
				}
			}
		}
	}()
	return statHistory, lastStatDurations, lastKbpsStats, dsStats
}

func processStatResults(results []cache.Result, statHistory StatHistoryThreadsafe, combinedStates peer.Crstates, lastKbpsStats StatsLastKbpsThreadsafe, toData todata.TOData, errorCount UintThreadsafe, dsStats DSStatsThreadsafe, lastStatEndTimes map[enum.CacheName]time.Time, lastStatDurations DurationMapThreadsafe) {
	for _, result := range results {
		// TODO determine if we want to add results with errors, or just print the errors now and don't add them.
		statHistory.Add(result)
	}

	for _, result := range results {
		log.Debugf("poll %v %v CreateStats start\n", result.PollID, time.Now())
	}

	newDsStats, newLastKbpsStats, err := ds.CreateStats(statHistory.Get(), toData, combinedStates, lastKbpsStats.Get(), time.Now())

	for _, result := range results {
		log.Debugf("poll %v %v CreateStats end\n", result.PollID, time.Now())
	}

	if err != nil {
		errorCount.Inc()
		log.Errorf("getting deliveryservice: %v\n", err)
	} else {
		dsStats.Set(newDsStats)
		lastKbpsStats.Set(newLastKbpsStats)
	}

	endTime := time.Now()
	for _, result := range results {
		if lastStatStart, ok := lastStatEndTimes[enum.CacheName(result.Id)]; ok {
			d := time.Since(lastStatStart)
			lastStatDurations.Set(enum.CacheName(result.Id), d)
		}
		lastStatEndTimes[enum.CacheName(result.Id)] = endTime

		// log.Debugf("poll %v %v statfinish\n", result.PollID, endTime)
		result.PollFinished <- result.PollID
	}
}
