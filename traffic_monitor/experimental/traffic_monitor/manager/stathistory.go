package manager

import (
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"sync"
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

func StartStatHistoryManager(cacheStatChan <-chan cache.Result) StatHistoryThreadsafe {
	statHistory := NewStatHistoryThreadsafe()
	go func() {
		for {
			select {
			case stat := <-cacheStatChan:
				statHistory.Add(stat)
			}
		}
	}()
	return statHistory
}
