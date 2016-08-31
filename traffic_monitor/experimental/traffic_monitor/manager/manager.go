package manager

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/fetcher"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/handler"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/poller"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/http_server"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	towrap "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	//	to "github.com/Comcast/traffic_control/traffic_ops/client"
	"github.com/davecheney/gmx"
)

const (
	defaultCacheHealthPollingInterval   time.Duration = 6 * time.Second
	defaultCacheStatPollingInterval     time.Duration = 6 * time.Second
	defaultMonitorConfigPollingInterval time.Duration = 5 * time.Second
	defaultHttpTimeout                  time.Duration = 2 * time.Second
	defaultPeerPollingInterval          time.Duration = 5 * time.Second
)

type StaticAppData struct {
	StartTime      time.Time
	GitRevision    string
	FreeMemoryMB   uint64
	Version        string
	WorkingDir     string
	Name           string
	BuildTimestamp string
}

//
// Kicks off the pollers and handlers
//
func Start(opsConfigFile string, staticAppData StaticAppData) {
	toSession := towrap.ITrafficOpsSession(nil)

	counters := fetcher.Counters{
		Success: gmx.NewCounter("fetchSuccess"),
		Fail:    gmx.NewCounter("fetchFail"),
		Pending: gmx.NewGauge("fetchPending"),
	}

	// TODO investigate whether a unique client per cache to be polled is faster
	sharedClient := &http.Client{
		Timeout:   defaultHttpTimeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	localStates := peer.NewCRStatesThreadsafe()     // this is the local state as discoverer by this traffic_monitor
	peerStates := peer.NewCRStatesPeersThreadsafe() // each peer's last state is saved in this map
	fetchCount := NewUintThreadsafe()               // note this is the number of individual caches fetched from, not the number of times all the caches were polled.
	healthIteration := NewUintThreadsafe()
	errorCount := NewUintThreadsafe()

	toData := todata.NewThreadsafe()
	dr := make(chan http_server.DataRequest)

	cacheHealthHandler := cache.NewHandler()
	cacheHealthPoller := poller.NewHTTP(defaultCacheHealthPollingInterval, true, sharedClient, counters, cacheHealthHandler)
	cacheStatHandler := cache.NewPrecomputeHandler(toData, peerStates) // TODO figure out if this is necessary, with the CacheHealthPoller
	cacheStatPoller := poller.NewHTTP(defaultCacheStatPollingInterval, false, sharedClient, counters, cacheStatHandler)
	monitorConfigPoller := poller.NewMonitorConfig(defaultMonitorConfigPollingInterval)
	peerHandler := peer.NewHandler()
	peerPoller := poller.NewHTTP(defaultPeerPollingInterval, false, sharedClient, counters, peerHandler)

	go monitorConfigPoller.Poll()
	go cacheHealthPoller.Poll()
	go cacheStatPoller.Poll()
	go peerPoller.Poll()

	opsConfig := StartOpsConfigManager(
		opsConfigFile,
		dr,
		toSession,
		toData,
		[]chan<- handler.OpsConfig{monitorConfigPoller.OpsConfigChannel},
		[]chan<- towrap.ITrafficOpsSession{monitorConfigPoller.SessionChannel})

	monitorConfig := StartMonitorConfigManager(
		monitorConfigPoller.ConfigChannel,
		localStates,
		cacheStatPoller.ConfigChannel,
		cacheHealthPoller.ConfigChannel,
		peerPoller.ConfigChannel)
	combinedStates := StartPeerManager(peerHandler.ResultChannel, localStates, peerStates)
	statHistory, _, lastKbpsStats, dsStats := StartStatHistoryManager(cacheStatHandler.ResultChannel, combinedStates, toData, errorCount)
	lastHealthDurations, events, localCacheStatus := StartHealthResultManager(
		cacheHealthHandler.ResultChannel,
		toData,
		localStates,
		statHistory,
		monitorConfig,
		peerStates,
		combinedStates,
		fetchCount,
		errorCount)
	StartDataRequestManager(
		dr,
		opsConfig,
		toSession,
		localStates,
		peerStates,
		combinedStates,
		statHistory,
		dsStats,
		events,
		staticAppData,
		cacheHealthPoller.Config.Interval,
		lastHealthDurations,
		fetchCount,
		healthIteration,
		errorCount,
		toData,
		localCacheStatus,
		lastKbpsStats)

	healthTickListener(cacheHealthPoller.TickChan, healthIteration)
}

// healthTickListener listens for health ticks, and writes to the health iteration variable. Does not return.
func healthTickListener(cacheHealthTick <-chan uint64, healthIteration UintThreadsafe) {
	for {
		select {
		case i := <-cacheHealthTick:
			healthIteration.Set(i)
		}
	}
}
