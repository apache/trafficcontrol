package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
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
	go func() {
		for {
			select {
			case stat := <-cacheStatChan:
				statHistory.Add(stat)

				now := time.Now()

				var err error
				createStatsCopyStatHistory := statHistory.Get()
				createStatsCopyCombinedStates := combinedStates.Get()
				createStatsCopyLastKbpsStats := lastKbpsStats.Get()
				toDataCopy := toData.Get()

				//				for _, healthResult := range results {
				fmt.Printf("DEBUG poll %v %v CreateStats start\n", stat.PollID, time.Now())
				//				}

				newDsStats, newLastKbpsStats, err := ds.CreateStats(createStatsCopyStatHistory, toDataCopy, createStatsCopyCombinedStates, createStatsCopyLastKbpsStats, now)

				//				for _, healthResult := range results {
				fmt.Printf("DEBUG poll %v %v CreateStats end\n", stat.PollID, time.Now())
				//				}

				if err != nil {
					errorCount.Inc()
					fmt.Printf("ERROR getting deliveryservice: %v\n", err)
				} else {
					dsStats.Set(newDsStats)
					lastKbpsStats.Set(newLastKbpsStats)
				}

				// for _, healthResult := range results {
				if lastStatStart, ok := lastStatEndTimes[enum.CacheName(stat.Id)]; ok {
					d := time.Since(lastStatStart)
					lastStatDurations.Set(enum.CacheName(stat.Id), d)
				}
				lastStatEndTimes[enum.CacheName(stat.Id)] = now

				fmt.Printf("DEBUG poll %v %v statfinish\n", stat.PollID, time.Now())
				stat.PollFinished <- stat.PollID
				//				}

			}
		}
	}()
	return statHistory, lastStatDurations, lastKbpsStats, dsStats
}
