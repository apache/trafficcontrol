package manager

import (
	"encoding/json"
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/http_server"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	towrap "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
)

type JSONEvents struct {
	Events []Event `json:"events"`
}

type CacheState struct {
	Value bool `json:"value"`
}

type ApiPeerStates struct {
	Peers map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState `json:"peers"`
}

// TODO make fields nullable, so error fields can be omitted, letting API callers still get updates for unerrored fields
type CacheStatus struct {
	Type                  *string  `json:"type,omitempty"`
	Status                *string  `json:"status,omitempty"`
	LoadAverage           *float64 `json:"load_average,omitempty"`
	QueryTimeMilliseconds *int64   `json:"query_time_ms,omitempty"`
	BandwidthKbps         *float64 `json:"bandwidth_kbps,omitempty"`
	ConnectionCount       *int64   `json:"connection_count,omitempty"`
}

func StartDataRequestManager(dr <-chan http_server.DataRequest, opsConfig OpsConfigThreadsafe, toSession towrap.ITrafficOpsSession, localStates peer.CRStatesThreadsafe, peerStates peer.CRStatesPeersThreadsafe, combinedStates peer.CRStatesThreadsafe, statHistory StatHistoryThreadsafe, dsStats DSStatsThreadsafe, events EventsThreadsafe, staticAppData StaticAppData, healthPollInterval time.Duration, lastHealthDurations DurationMapThreadsafe, fetchCount UintThreadsafe, healthIteration UintThreadsafe, errorCount UintThreadsafe, toData todata.TODataThreadsafe, localCacheStatus CacheAvailableStatusThreadsafe, lastKbpsStats StatsLastKbpsThreadsafe) {
	go dataRequestManagerListen(dr, opsConfig, toSession, localStates, peerStates, combinedStates, statHistory, dsStats, events, staticAppData, healthPollInterval, lastHealthDurations, fetchCount, healthIteration, errorCount, toData, localCacheStatus, lastKbpsStats)
}

func dataRequestManagerListen(dr <-chan http_server.DataRequest, opsConfig OpsConfigThreadsafe, toSession towrap.ITrafficOpsSession, localStates peer.CRStatesThreadsafe, peerStates peer.CRStatesPeersThreadsafe, combinedStates peer.CRStatesThreadsafe, statHistory StatHistoryThreadsafe, dsStats DSStatsThreadsafe, events EventsThreadsafe, staticAppData StaticAppData, healthPollInterval time.Duration, lastHealthDurations DurationMapThreadsafe, fetchCount UintThreadsafe, healthIteration UintThreadsafe, errorCount UintThreadsafe, toData todata.TODataThreadsafe, localCacheStatus CacheAvailableStatusThreadsafe, lastKbpsStats StatsLastKbpsThreadsafe) {
	for {
		select {
		case req := <-dr:
			defer close(req.Response)

			var body []byte
			var err error

			switch req.Type {
			case http_server.TRConfig:
				cdnName := opsConfig.Get().CdnName
				if toSession == nil {
					err = fmt.Errorf("Unable to connect to Traffic Ops")
				} else if cdnName == "" {
					err = fmt.Errorf("No CDN Configured")
				} else {
					body, err = toSession.CRConfigRaw(cdnName)
				}
				if err != nil {
					err = fmt.Errorf("TR Config: %v", err)
				}
			case http_server.TRStateDerived:
				body, err = peer.CrstatesMarshall(combinedStates.Get())
				if err != nil {
					err = fmt.Errorf("TR State (derived): %v", err)
				}
			case http_server.TRStateSelf:
				body, err = peer.CrstatesMarshall(localStates.Get())
				if err != nil {
					err = fmt.Errorf("TR State (self): %v", err)
				}
			case http_server.CacheStats:
				// TODO: add support for ?hc=N query param, stats=, wildcard, individual caches
				// add pp and date to the json:
				/*
					pp: "0=[my-ats-edge-cache-1], hc=[1]",
					date: "Thu Oct 09 20:28:36 UTC 2014"
				*/
				params := req.Parameters
				hc := 1
				if _, exists := params["hc"]; exists {
					v, err := strconv.Atoi(params["hc"][0])
					if err == nil {
						hc = v
					}
				}
				body, err = cache.StatsMarshall(statHistory.Get(), hc)
				if err != nil {
					err = fmt.Errorf("CacheStats: %v", err)
				}
			case http_server.DSStats:
				body, err = json.Marshal(ds.StatsJSON(dsStats.Get())) // TODO marshall beforehand, for performance? (test to see how often requests are made)
				if err != nil {
					err = fmt.Errorf("DsStats: %v", err)
				}
			case http_server.EventLog:
				body, err = json.Marshal(JSONEvents{Events: events.Get()})
				if err != nil {
					err = fmt.Errorf("EventLog: %v", err)
				}
			case http_server.PeerStates:
				body, err = json.Marshal(createApiPeerStates(peerStates.Get()))
			case http_server.StatSummary:
				body = []byte("TODO implement")
			case http_server.Stats:
				body, err = getStats(staticAppData, healthPollInterval, lastHealthDurations.Get(), fetchCount.Get(), healthIteration.Get(), errorCount.Get())
				if err != nil {
					err = fmt.Errorf("Stats: %v", err)
				}
			case http_server.ConfigDoc:
				opsConfigCopy := opsConfig.Get()
				// if the password is blank, leave it blank, so callers can see it's missing.
				if opsConfigCopy.Password != "" {
					opsConfigCopy.Password = "*****"
				}
				body, err = json.Marshal(opsConfigCopy)
				if err != nil {
					err = fmt.Errorf("Config Doc: %v", err)
				}
			case http_server.APICacheCount: // TODO determine if this should use peerStates
				body = []byte(strconv.Itoa(len(localStates.Get().Caches)))
			case http_server.APICacheAvailableCount:
				body = []byte(strconv.Itoa(cacheAvailableCount(localStates.Get().Caches)))
			case http_server.APICacheDownCount:
				body = []byte(strconv.Itoa(cacheDownCount(localStates.Get().Caches)))
			case http_server.APIVersion:
				s := "traffic_monitor-" + staticAppData.Version + "."
				if len(staticAppData.GitRevision) > 6 {
					s += staticAppData.GitRevision[:6]
				} else {
					s += staticAppData.GitRevision
				}
				body = []byte(s)
			case http_server.APITrafficOpsURI:
				body = []byte(opsConfig.Get().Url)
			case http_server.APICacheStates:
				body, err = json.Marshal(createCacheStatuses(toData.Get().ServerTypes, statHistory.Get(), lastHealthDurations.Get(), localStates.Get().Caches, lastKbpsStats.Get(), localCacheStatus))
			case http_server.APIBandwidthKbps:
				serverTypes := toData.Get().ServerTypes
				kbpsStats := lastKbpsStats.Get()
				sum := float64(0.0)
				for cache, data := range kbpsStats.Caches {
					if serverTypes[cache] != enum.CacheTypeEdge {
						continue
					}
					sum += data.Kbps
				}
				body = []byte(fmt.Sprintf("%f", sum))
			default:
				err = fmt.Errorf("Unknown Request Type: %v", req.Type)
			}

			if err != nil {
				errorCount.Inc()
				log.Errorf("Request Error: %v\n", err)
			} else {
				req.Response <- body
			}
		}
	}
}

func createCacheStatuses(cacheTypes map[enum.CacheName]enum.CacheType, statHistory map[enum.CacheName][]cache.Result, lastHealthDurations map[enum.CacheName]time.Duration, cacheStates map[string]peer.IsAvailable, lastKbpsStats ds.StatsLastKbps, localCacheStatusThreadsafe CacheAvailableStatusThreadsafe) map[enum.CacheName]CacheStatus {
	conns := createCacheConnections(statHistory)
	statii := map[enum.CacheName]CacheStatus{}
	localCacheStatus := localCacheStatusThreadsafe.Get()

	for cacheName, cacheType := range cacheTypes {
		cacheStatHistory, ok := statHistory[cacheName]
		if !ok {
			log.Warnf("createCacheStatuses stat history missing cache %s\n", cacheName)
			continue
		}

		if len(cacheStatHistory) < 1 {
			log.Warnf("createCacheStatuses stat history empty for cache %s\n", cacheName)
			continue
		}

		log.Debugf("createCacheStatuses NOT empty for cache %s\n", cacheName)

		var loadAverage *float64
		procLoadAvg := cacheStatHistory[0].Astats.System.ProcLoadavg
		if procLoadAvg != "" {
			firstSpace := strings.IndexRune(procLoadAvg, ' ')
			if firstSpace == -1 {
				log.Warnf("WARNING unexpected proc.loadavg '%s' for cache %s\n", procLoadAvg, cacheName)
			} else {
				loadAverageVal, err := strconv.ParseFloat(procLoadAvg[:firstSpace], 64)
				if err != nil {
					log.Warnf("proc.loadavg doesn't contain a float prefix '%s' for cache %s\n", procLoadAvg, cacheName)
				} else {
					loadAverage = &loadAverageVal
				}
			}
		}

		var queryTime *int64
		queryTimeVal, ok := lastHealthDurations[cacheName]
		if !ok {
			log.Warnf("cache not in last health durations cache %s\n", cacheName)
		} else {
			queryTimeInt := int64(queryTimeVal / time.Millisecond)
			queryTime = &queryTimeInt
		}

		var kbps *float64
		kbpsVal, ok := lastKbpsStats.Caches[enum.CacheName(cacheName)]
		if !ok {
			log.Warnf("cache not in last kbps cache %s\n", cacheName)
		} else {
			kbps = &kbpsVal.Kbps
		}

		var connections *int64
		connectionsVal, ok := conns[enum.CacheName(cacheName)]
		if !ok {
			log.Warnf("cache not in connections %s\n", cacheName)
		} else {
			connections = &connectionsVal
		}

		var status *string
		statusVal, ok := localCacheStatus[enum.CacheName(cacheName)]
		if !ok {
			log.Warnf("cache not in statuses %s\n", cacheName)
		} else {
			statusString := statusVal.Status + " - "
			if localCacheStatus[enum.CacheName(cacheName)].Available {
				statusString += "available"
			} else {
				statusString += "unavailable"
			}
			status = &statusString
		}

		cacheTypeStr := string(cacheType)
		statii[enum.CacheName(cacheName)] = CacheStatus{Type: &cacheTypeStr, LoadAverage: loadAverage, QueryTimeMilliseconds: queryTime, BandwidthKbps: kbps, ConnectionCount: connections, Status: status}
	}
	return statii
}

func createCacheConnections(statHistory map[enum.CacheName][]cache.Result) map[enum.CacheName]int64 {
	conns := map[enum.CacheName]int64{}
	for server, history := range statHistory {
		for _, result := range history {
			val, ok := result.Astats.Ats["proxy.process.http.current_client_connections"]
			if !ok {
				continue
			}

			v, ok := val.(float64)
			if !ok {
				continue
			}

			conns[server] = int64(v)
			break
		}
	}
	return conns
}

func cacheDownCount(caches map[string]peer.IsAvailable) int {
	count := 0
	for _, available := range caches {
		if !available.IsAvailable {
			count++
		}
	}
	return count
}

func cacheAvailableCount(caches map[string]peer.IsAvailable) int {
	return len(caches) - cacheDownCount(caches)
}

func createApiPeerStates(peerStates map[string]peer.Crstates) ApiPeerStates {
	apiPeerStates := ApiPeerStates{Peers: map[enum.TrafficMonitorName]map[enum.CacheName][]CacheState{}}

	for peer, state := range peerStates {
		if _, ok := apiPeerStates.Peers[enum.TrafficMonitorName(peer)]; !ok {
			apiPeerStates.Peers[enum.TrafficMonitorName(peer)] = map[enum.CacheName][]CacheState{}
		}
		peerState := apiPeerStates.Peers[enum.TrafficMonitorName(peer)]
		for cache, available := range state.Caches {
			peerState[enum.CacheName(cache)] = []CacheState{CacheState{Value: available.IsAvailable}}
		}
		apiPeerStates.Peers[enum.TrafficMonitorName(peer)] = peerState
	}
	return apiPeerStates
}

type Stats struct {
	MaxMemoryMB         uint64 `json:"Max Memory (MB)"`
	GitRevision         string `json:"git-revision"`
	ErrorCount          uint64 `json:"Error Count"`
	Uptime              uint64 `json:"uptime"`
	FreeMemoryMB        uint64 `json:"Free Memory (MB)"`
	TotalMemoryMB       uint64 `json:"Total Memory (MB)"`
	Version             string `json:"version"`
	DeployDir           string `json:"deploy-dir"`
	FetchCount          uint64 `json:"Fetch Count"`
	QueryIntervalDelta  int    `json:"Query Interval Delta"`
	IterationCount      uint64 `json:"Iteration Count"`
	Name                string `json:"name"`
	BuildTimestamp      string `json:"buildTimestamp"`
	QueryIntervalTarget int    `json:"Query Interval Target"`
	QueryIntervalActual int    `json:"Query Interval Actual"`
	SlowestCache        string `json:"Slowest Cache"`
	LastQueryInterval   int    `json:"Last Query Interval"`
}

func getLongestPoll(lastHealthTimes map[enum.CacheName]time.Duration) (enum.CacheName, time.Duration) {
	var longestCache enum.CacheName
	var longestTime time.Duration
	for cache, time := range lastHealthTimes {
		if time > longestTime {
			longestTime = time
			longestCache = cache
		}
	}
	return longestCache, longestTime
}

func getStats(staticAppData StaticAppData, pollingInterval time.Duration, lastHealthTimes map[enum.CacheName]time.Duration, fetchCount uint64, healthIteration uint64, errorCount uint64) ([]byte, error) {
	longestPollCache, longestPollTime := getLongestPoll(lastHealthTimes)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var s Stats
	s.MaxMemoryMB = memStats.TotalAlloc / (1024 * 1024)
	s.GitRevision = staticAppData.GitRevision
	s.ErrorCount = errorCount
	s.Uptime = uint64(time.Since(staticAppData.StartTime) / time.Second)
	s.FreeMemoryMB = staticAppData.FreeMemoryMB
	s.TotalMemoryMB = memStats.Alloc / (1024 * 1024) // TODO rename to "used memory" if/when nothing is using the JSON entry
	s.Version = staticAppData.Version
	s.DeployDir = staticAppData.WorkingDir
	s.FetchCount = fetchCount
	s.SlowestCache = string(longestPollCache)
	s.IterationCount = healthIteration
	s.Name = staticAppData.Name
	s.BuildTimestamp = staticAppData.BuildTimestamp
	s.QueryIntervalTarget = int(pollingInterval / time.Millisecond)
	s.QueryIntervalActual = int(longestPollTime / time.Millisecond)
	s.QueryIntervalDelta = s.QueryIntervalActual - s.QueryIntervalTarget
	s.LastQueryInterval = int(math.Max(float64(s.QueryIntervalActual), float64(s.QueryIntervalTarget)))

	return json.Marshal(s)
}
