package manager

import (
	"fmt"
	"sync"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/handler"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/poller"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/config"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/http_server"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	towrap "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/Comcast/traffic_control/traffic_ops/client"
)

// This could be made lock-free, if the performance was necessary
type OpsConfigThreadsafe struct {
	opsConfig *handler.OpsConfig
	m         *sync.RWMutex
}

func NewOpsConfigThreadsafe() OpsConfigThreadsafe {
	return OpsConfigThreadsafe{m: &sync.RWMutex{}, opsConfig: &handler.OpsConfig{}}
}

func (o *OpsConfigThreadsafe) Get() handler.OpsConfig {
	o.m.RLock()
	defer o.m.RUnlock()
	return *o.opsConfig
}

func (o *OpsConfigThreadsafe) Set(newOpsConfig handler.OpsConfig) {
	o.m.Lock()
	*o.opsConfig = newOpsConfig
	o.m.Unlock()
}

// Note the OpsConfigManager is in charge of the httpServer, because ops config changes trigger server changes. If other things needed to trigger server restarts, the server could be put in its own goroutine with signal channels
func StartOpsConfigManager(
	opsConfigFile string,
	toSession towrap.ITrafficOpsSession,
	toData todata.TODataThreadsafe,
	opsConfigChangeSubscribers []chan<- handler.OpsConfig,
	toChangeSubscribers []chan<- towrap.ITrafficOpsSession,
	localStates peer.CRStatesThreadsafe,
	peerStates peer.CRStatesPeersThreadsafe,
	combinedStates peer.CRStatesThreadsafe,
	statHistory StatHistoryThreadsafe,
	lastStats LastStatsThreadsafe,
	dsStats DSStatsThreadsafe,
	events EventsThreadsafe,
	staticAppData StaticAppData,
	healthPollInterval time.Duration,
	lastHealthDurations DurationMapThreadsafe,
	fetchCount UintThreadsafe,
	healthIteration UintThreadsafe,
	errorCount UintThreadsafe,
	localCacheStatus CacheAvailableStatusThreadsafe,
	unpolledCaches UnpolledCachesThreadsafe,
	cfg config.Config,
) OpsConfigThreadsafe {

	opsConfigFileChannel := make(chan interface{})
	opsConfigFilePoller := poller.FilePoller{
		File:          opsConfigFile,
		ResultChannel: opsConfigFileChannel,
	}

	opsConfigChannel := make(chan handler.OpsConfig)
	opsConfigFileHandler := handler.OpsConfigFileHandler{
		ResultChannel:    opsConfigFilePoller.ResultChannel,
		OpsConfigChannel: opsConfigChannel,
	}

	go opsConfigFileHandler.Listen()
	go opsConfigFilePoller.Poll()

	opsConfig := NewOpsConfigThreadsafe()

	// TODO remove change subscribers, give Threadsafes directly to the things that need them. If they only set vars, and don't actually do work on change.
	go func() {
		httpServer := http_server.Server{}

		for {
			select {
			case newOpsConfig := <-opsConfigChannel:
				var err error
				opsConfig.Set(newOpsConfig)

				listenAddress := ":80" // default

				if newOpsConfig.HttpListener != "" {
					listenAddress = newOpsConfig.HttpListener
				}

				handleErr := func(err error) {
					errorCount.Inc()
					log.Errorf("OpsConfigManager: %v\n", err)
				}

				err = httpServer.Run(func(req http_server.DataRequest) ([]byte, int) {
					return DataRequest(
						req,
						opsConfig,
						toSession,
						localStates,
						peerStates,
						combinedStates,
						statHistory,
						dsStats,
						events,
						staticAppData,
						healthPollInterval,
						lastHealthDurations,
						fetchCount,
						healthIteration,
						errorCount,
						toData,
						localCacheStatus,
						lastStats,
						unpolledCaches,
					)
				}, listenAddress, cfg.ServeReadTimeout, cfg.ServeWriteTimeout)
				if err != nil {
					handleErr(fmt.Errorf("MonitorConfigPoller: error creating HTTP server: %s\n", err))
					continue
				}

				realToSession, err := to.Login(newOpsConfig.Url, newOpsConfig.Username, newOpsConfig.Password, newOpsConfig.Insecure)
				if err != nil {
					handleErr(fmt.Errorf("MonitorConfigPoller: error instantiating Session with traffic_ops: %s\n", err))
					continue
				}
				toSession.Set(realToSession)

				if err := toData.Fetch(toSession, newOpsConfig.CdnName); err != nil {
					handleErr(fmt.Errorf("Error getting Traffic Ops data: %v\n", err))
					continue
				}

				// These must be in a goroutine, because the monitorConfigPoller tick sends to a channel this select listens for. Thus, if we block on sends to the monitorConfigPoller, we have a livelock race condition.
				// More generically, we're using goroutines as an infinite chan buffer, to avoid potential livelocks
				for _, subscriber := range opsConfigChangeSubscribers {
					go func(s chan<- handler.OpsConfig) { s <- newOpsConfig }(subscriber)
				}
				for _, subscriber := range toChangeSubscribers {
					go func(s chan<- towrap.ITrafficOpsSession) { s <- toSession }(subscriber)
				}
			}
		}
	}()

	return opsConfig
}
