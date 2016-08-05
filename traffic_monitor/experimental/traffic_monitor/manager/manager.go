package manager

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/fetcher"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/handler"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/poller"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/cache"
	ds "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/deliveryservice"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/health"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/http_server"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	todata "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopsdata"
	towrap "github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/trafficopswrapper"
	to "github.com/Comcast/traffic_control/traffic_ops/client"
	"github.com/davecheney/gmx"
)

const (
	defaultCacheHealthPollingInterval   time.Duration = 1 * time.Second
	defaultCacheStatPollingInterval     time.Duration = 5 * time.Second
	defaultMonitorConfigPollingInterval time.Duration = 5 * time.Second
	defaultHttpTimeout                  time.Duration = 2 * time.Second
	defaultPeerPollingInterval          time.Duration = 5 * time.Second
)

//const maxHistory = (60 / pollingInterval) * 5
const defaultMaxHistory = 5

const maxEvents = 200

type StaticAppData struct {
	StartTime      time.Time
	GitRevision    string
	FreeMemoryMB   uint64
	Version        string
	WorkingDir     string
	Name           string
	BuildTimestamp string
}

// TODO put somewhere more generic
const CacheAvailableStatusReported = "REPORTED"

type CacheAvailableStatus struct {
	Available bool
	Status    string
}

//
// Kicks off the pollers and handlers
//
func Start(opsConfigFile string, staticAppData StaticAppData) {
	toSession := towrap.ITrafficOpsSession(nil)

	fetchSuccessCounter := gmx.NewCounter("fetchSuccess")
	fetchFailCounter := gmx.NewCounter("fetchFail")
	fetchPendingGauge := gmx.NewGauge("fetchPending")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	sharedClient := http.Client{
		Timeout:   defaultHttpTimeout,
		Transport: tr,
	}

	cacheHealthConfigChannel := make(chan poller.HttpPollerConfig)
	cacheHealthChannel := make(chan cache.Result)
	cacheHealthTick := make(chan uint64)
	cacheHealthPoller := poller.HttpPoller{
		TickChan:      cacheHealthTick,
		ConfigChannel: cacheHealthConfigChannel,
		Config: poller.HttpPollerConfig{
			Interval: defaultCacheHealthPollingInterval,
		},
		Fetcher: fetcher.HttpFetcher{
			Handler: cache.Handler{ResultChannel: cacheHealthChannel},
			Client:  sharedClient,
			Success: fetchSuccessCounter,
			Fail:    fetchFailCounter,
			Pending: fetchPendingGauge,
		},
	}

	cacheStatConfigChannel := make(chan poller.HttpPollerConfig)
	cacheStatChannel := make(chan cache.Result)
	cacheStatPoller := poller.HttpPoller{
		ConfigChannel: cacheStatConfigChannel,
		Config: poller.HttpPollerConfig{
			Interval: defaultCacheStatPollingInterval,
		},
		Fetcher: fetcher.HttpFetcher{
			Handler: cache.Handler{ResultChannel: cacheStatChannel},
			Client:  sharedClient,
			Success: fetchSuccessCounter,
			Fail:    fetchFailCounter,
			Pending: fetchPendingGauge,
		},
	}

	sessionChannel := make(chan towrap.ITrafficOpsSession)
	monitorConfigChannel := make(chan to.TrafficMonitorConfigMap)
	monitorOpsConfigChannel := make(chan handler.OpsConfig)
	monitorConfigPoller := poller.MonitorConfigPoller{
		Interval:         defaultMonitorConfigPollingInterval,
		SessionChannel:   sessionChannel,
		ConfigChannel:    monitorConfigChannel,
		OpsConfigChannel: monitorOpsConfigChannel,
	}

	peerConfigChannel := make(chan poller.HttpPollerConfig)
	peerChannel := make(chan peer.Result)
	peerPoller := poller.HttpPoller{
		ConfigChannel: peerConfigChannel,
		Config: poller.HttpPollerConfig{
			Interval: defaultPeerPollingInterval,
		},
		Fetcher: fetcher.HttpFetcher{
			Handler: peer.Handler{ResultChannel: peerChannel},
			Client:  sharedClient,
			Success: fetchSuccessCounter,
			Fail:    fetchFailCounter,
			Pending: fetchPendingGauge,
		},
	}

	go monitorConfigPoller.Poll()
	go cacheHealthPoller.Poll()
	go cacheStatPoller.Poll()
	go peerPoller.Poll()

	toData := todata.NewThreadsafe()
	dr := make(chan http_server.DataRequest)

	healthHistory := map[string][]interface{}{}
	statHistory := map[string][]interface{}{}

	opsConfig := StartOpsConfigManager(opsConfigFile, dr, toSession, toData, []chan<- handler.OpsConfig{monitorConfigPoller.OpsConfigChannel}, []chan<- towrap.ITrafficOpsSession{monitorConfigPoller.SessionChannel})

	localStates := NewCRStatesThreadsafe()     // this is the local state as discoverer by this traffic_monitor
	peerStates := NewCRStatesPeersThreadsafe() // each peer's last state is saved in this map
	combinedStates := NewCRStatesThreadsafe()  // this is the result of combining the localStates and all the peerStates using the var ??

	// TODO put stat data in a struct, for brevity
	lastHealthEndTimes := map[string]time.Time{}
	lastHealthDurations := map[string]time.Duration{}
	fetchCount := uint64(0) // note this is the number of individual caches fetched from, not the number of times all the caches were polled.
	healthIteration := uint64(0)
	errorCount := uint64(0)
	events := []Event{}
	eventIndex := uint64(0)
	dsStats := ds.NewStats()
	lastKbpsStats := ds.NewStatsLastKbps()
	localCacheStatus := map[enum.CacheName]CacheAvailableStatus{}

	monitorConfig := StartMonitorConfigManager(monitorConfigPoller.ConfigChannel, localStates, cacheStatPoller.ConfigChannel, cacheHealthPoller.ConfigChannel, peerPoller.ConfigChannel)

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
				body, err = cache.StatsMarshall(statHistory, hc)
				if err != nil {
					err = fmt.Errorf("CacheStats: %v", err)
				}
			case http_server.DSStats:
				body, err = json.Marshal(ds.StatsJSON(dsStats)) // TODO marshall beforehand, for performance? (test to see how often requests are made)
				if err != nil {
					err = fmt.Errorf("DsStats: %v", err)
				}
			case http_server.EventLog:
				body, err = json.Marshal(JSONEvents{Events: events})
				if err != nil {
					err = fmt.Errorf("EventLog: %v", err)
				}
			case http_server.PeerStates:
				body, err = json.Marshal(createApiPeerStates(peerStates.Get()))
			case http_server.StatSummary:
				body = []byte("TODO implement")
			case http_server.Stats:
				body, err = getStats(staticAppData, cacheHealthPoller.Config.Interval, lastHealthDurations, fetchCount, healthIteration, errorCount)
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
				body, err = json.Marshal(createCacheStatuses(toData.Get().ServerTypes, statHistory, lastHealthDurations, localStates.Get().Caches, lastKbpsStats, localCacheStatus))
			default:
				err = fmt.Errorf("Unknown Request Type: %v", req.Type)
			}

			if err != nil {
				errorCount++
				log.Printf("ERROR Request Error: %v\n", err)
			} else {
				req.Response <- body
			}
			// case monitorConfig
		case i := <-cacheHealthTick:
			healthIteration = i
		case healthResult := <-cacheHealthChannel:
			fetchCount++
			toDataCopy := toData.Get() // create a copy, so the same data used for all processing of this cache health result
			var prevResult cache.Result
			if len(healthHistory[healthResult.Id]) != 0 {
				prevResult = healthHistory[healthResult.Id][len(healthHistory[healthResult.Id])-1].(cache.Result)
			}
			monitorConfigCopy := monitorConfig.Get() // copy now, so all calculations are on the same data
			health.GetVitals(&healthResult, &prevResult, &monitorConfigCopy)
			healthHistory[healthResult.Id] = pruneHistory(append(healthHistory[healthResult.Id], healthResult), defaultMaxHistory)
			isAvailable, whyAvailable := health.EvalCache(healthResult, &monitorConfigCopy)
			if localStates.Get().Caches[healthResult.Id].IsAvailable != isAvailable {
				fmt.Println("Changing state for", healthResult.Id, " was:", prevResult.Available, " is now:", isAvailable, " because:", whyAvailable, " errors:", healthResult.Errors)
				e := Event{Index: eventIndex, Time: time.Now().Unix(), Description: whyAvailable, Name: healthResult.Id, Hostname: healthResult.Id, Type: toDataCopy.ServerTypes[healthResult.Id].String(), Available: isAvailable}
				events = append([]Event{e}, events...)
				if len(events) > maxEvents {
					events = events[:maxEvents-1]
				}
				eventIndex++
			}

			localCacheStatus[enum.CacheName(healthResult.Id)] = CacheAvailableStatus{Available: isAvailable, Status: monitorConfigCopy.TrafficServer[healthResult.Id].Status} // TODO move within localStates
			localStates.SetCache(healthResult.Id, peer.IsAvailable{IsAvailable: isAvailable})
			calculateDeliveryServiceState(toDataCopy.DeliveryServiceServers, localStates)

			// TODO determine if we should combineCrStates() here

			now := time.Now()

			var err error
			dsStats, lastKbpsStats, err = ds.CreateStats(statHistory, toDataCopy, combinedStates.Get(), lastKbpsStats, now)
			if err != nil {
				errorCount++
				log.Printf("ERROR getting deliveryservice: %v\n", err)
			}

			if lastHealthStart, ok := lastHealthEndTimes[healthResult.Id]; ok {
				lastHealthDurations[healthResult.Id] = time.Since(lastHealthStart)
			}
			lastHealthEndTimes[healthResult.Id] = now
			fmt.Printf("DEBUG health duration for %s: %v\n", healthResult.Id, lastHealthDurations[healthResult.Id])

			// if _, ok := queryIntervalStart[pollI]; !ok {
			// 	log.Printf("ERROR poll start index not found")
			// 	continue
			// }
			// lastQueryIntervalTime = time.Since(queryIntervalStart[pollI])
		case stats := <-cacheStatChannel:
			statHistory[stats.Id] = pruneHistory(append(statHistory[stats.Id], stats), defaultMaxHistory)
		case crStatesResult := <-peerChannel:
			peerStates.Set(crStatesResult.Id, crStatesResult.PeerStats)
			combinedStates.Set(combineCrStates(peerStates.Get(), localStates.Get()))
		}
	}
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

func createCacheStatuses(cacheTypes map[string]enum.CacheType, statHistory map[string][]interface{}, lastHealthDurations map[string]time.Duration, cacheStates map[string]peer.IsAvailable, lastKbpsStats ds.StatsLastKbps, localCacheStatus map[enum.CacheName]CacheAvailableStatus) map[enum.CacheName]CacheStatus {
	conns := createCacheConnections(statHistory)
	statii := map[enum.CacheName]CacheStatus{}
	for cacheName, cacheType := range cacheTypes {

		cacheStatHistory, ok := statHistory[cacheName]
		if !ok {
			log.Printf("WARNING DEBUG6 createCacheStatuses stat history missing cache %s\n", cacheName)
			continue
		}

		if len(cacheStatHistory) < 1 {
			log.Printf("WARNING DEBUG6 createCacheStatuses stat history empty for cache %s\n", cacheName)
			continue
		}

		log.Printf("DEBUGQ createCacheStatuses NOT empty for cache %s\n", cacheName)

		var loadAverage *float64
		procLoadAvg := cacheStatHistory[0].(cache.Result).Astats.System.ProcLoadavg
		if procLoadAvg != "" {
			firstSpace := strings.IndexRune(procLoadAvg, ' ')
			if firstSpace == -1 {
				log.Printf("WARNING DEBUG6 unexpected proc.loadavg '%s' for cache %s\n", procLoadAvg, cacheName)
			} else {
				loadAverageVal, err := strconv.ParseFloat(procLoadAvg[:firstSpace], 64)
				if err != nil {
					log.Printf("WARNING proc.loadavg doesn't contain a float prefix '%s' for cache %s\n", procLoadAvg, cacheName)
				} else {
					loadAverage = &loadAverageVal
				}
			}
		}

		var queryTime *int64
		queryTimeVal, ok := lastHealthDurations[cacheName]
		if !ok {
			log.Printf("WARNING DEBUGQ cache not in last health durations cache %s\n", cacheName)
		} else {
			queryTimeInt := int64(queryTimeVal / time.Millisecond)
			queryTime = &queryTimeInt
		}

		var kbps *float64
		kbpsVal, ok := lastKbpsStats.Caches[enum.CacheName(cacheName)]
		if !ok {
			log.Printf("WARNING DEBUGQ cache not in last kbps cache %s\n", cacheName)
		} else {
			kbps = &kbpsVal.Kbps
		}

		var connections *int64
		connectionsVal, ok := conns[enum.CacheName(cacheName)]
		if !ok {
			log.Printf("WARNING DEBUGQ cache not in connections %s\n", cacheName)
		} else {
			connections = &connectionsVal
		}

		var status *string
		statusVal, ok := localCacheStatus[enum.CacheName(cacheName)]
		if !ok {
			log.Printf("WARNING DEBUGQ cache not in statuses %s\n", cacheName)
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

// TODO: run these in goroutines
func createCacheHealthStatuses(statHistory map[string][]interface{}) map[enum.CacheName]string {
	statuses := map[enum.CacheName]string{}
	for server, history := range statHistory {
		for _, iresult := range history {
			result, ok := iresult.(cache.Result)
			if !ok {
				fmt.Printf("ERROR DEBUG6 history contained unexpected result type %T\n", iresult)
				continue
			}

			val, ok := result.Astats.Ats["status"]
			if !ok {
				fmt.Printf("ERROR DEBUG8 status stat not found for %s\n", server)
				continue
			}
			fmt.Printf("ERROR DEBUG6 status stat WAS FOUND for %s\n", server)

			v, ok := val.(string)
			if !ok {
				fmt.Printf("ERROR status stat value expected string actual '%v' type %T", val, val)
				continue
			}

			statuses[enum.CacheName(server)] = v
		}
	}
	return statuses
}

func createCacheConnections(statHistory map[string][]interface{}) map[enum.CacheName]int64 {
	conns := map[enum.CacheName]int64{}
	for server, history := range statHistory {
		for _, iresult := range history {
			result, ok := iresult.(cache.Result)
			if !ok {
				fmt.Printf("ERROR DEBUG6 history contained unexpected result type %T\n", iresult)
				continue
			}

			val, ok := result.Astats.Ats["proxy.process.http.total_incoming_connections"]
			if !ok {
				fmt.Printf("ERROR DEBUG6 connections stat not found for %s\n", server)
				continue
			}

			v, ok := val.(float64)
			if !ok {
				fmt.Printf("ERROR connection stat value expected int actual '%v' type %T", val, val)
				continue
			}

			conns[enum.CacheName(server)] = int64(v)
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

type TrafficMonitorName string

type ApiPeerStates struct {
	Peers map[TrafficMonitorName]map[enum.CacheName][]CacheState `json:"peers"`
}

type CacheState struct {
	Value bool `json:"value"`
}

func createApiPeerStates(peerStates map[string]peer.Crstates) ApiPeerStates {
	apiPeerStates := ApiPeerStates{Peers: map[TrafficMonitorName]map[enum.CacheName][]CacheState{}}

	for peer, state := range peerStates {
		if _, ok := apiPeerStates.Peers[TrafficMonitorName(peer)]; !ok {
			apiPeerStates.Peers[TrafficMonitorName(peer)] = map[enum.CacheName][]CacheState{}
		}
		peerState := apiPeerStates.Peers[TrafficMonitorName(peer)]
		for cache, available := range state.Caches {
			peerState[enum.CacheName(cache)] = []CacheState{CacheState{Value: available.IsAvailable}}
		}
		apiPeerStates.Peers[TrafficMonitorName(peer)] = peerState
	}
	return apiPeerStates
}

type JSONEvents struct {
	Events []Event `json:"events"`
}

type Event struct {
	Index       uint64 `json:"index"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Hostname    string `json:"hostname"`
	Type        string `json:"type"`
	Available   bool   `json:"isAvailable"`
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

func getLongestPoll(lastHealthTimes map[string]time.Duration) (string, time.Duration) {
	var longestCache string
	var longestTime time.Duration
	for cache, time := range lastHealthTimes {
		if time > longestTime {
			longestTime = time
			longestCache = cache
		}
	}
	return longestCache, longestTime
}

func getStats(staticAppData StaticAppData, pollingInterval time.Duration, lastHealthTimes map[string]time.Duration, fetchCount uint64, healthIteration uint64, errorCount uint64) ([]byte, error) {
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
	s.SlowestCache = longestPollCache
	s.IterationCount = healthIteration
	s.Name = staticAppData.Name
	s.BuildTimestamp = staticAppData.BuildTimestamp
	s.QueryIntervalTarget = int(pollingInterval / time.Millisecond)
	s.QueryIntervalActual = int(longestPollTime / time.Millisecond)
	s.QueryIntervalDelta = s.QueryIntervalActual - s.QueryIntervalTarget
	s.LastQueryInterval = int(math.Max(float64(s.QueryIntervalActual), float64(s.QueryIntervalTarget)))

	return json.Marshal(s)
}

// addStateDeliveryServices adds delivery services in `mc` as keys in `deliveryServices`, with empty Deliveryservice values.
// TODO add disabledLocations
func addStateDeliveryServices(mc to.TrafficMonitorConfigMap, deliveryServices map[string]peer.Deliveryservice) {
	for _, ds := range mc.DeliveryService {
		// since caches default to unavailable, also default DS false
		deliveryServices[ds.XMLID] = peer.Deliveryservice{}
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

// TODO JvD: add deliveryservice stuff
func combineCrStates(peerStates map[string]peer.Crstates, localStates peer.Crstates) peer.Crstates {
	combinedStates := peer.NewCrstates()
	for cacheName, localCacheState := range localStates.Caches { // localStates gets pruned when servers are disabled, it's the source of truth
		downVotes := 0 // TODO JvD: change to use parameter when deciding to be optimistic or pessimistic.
		if localCacheState.IsAvailable {
			// fmt.Println(cacheName, " is available locally - setting to IsAvailable: true")
			combinedStates.Caches[cacheName] = peer.IsAvailable{IsAvailable: true} // we don't care about the peers, we got a "good one", and we're optimistic
		} else {
			downVotes++ // localStates says it's not happy
			for _, peerCrStates := range peerStates {
				if peerCrStates.Caches[cacheName].IsAvailable {
					// fmt.Println(cacheName, "- locally we think it's down, but", peerName, "says IsAvailable: ", peerCrStates.Caches[cacheName].IsAvailable, "trusting the peer.")
					combinedStates.Caches[cacheName] = peer.IsAvailable{IsAvailable: true} // we don't care about the peers, we got a "good one", and we're optimistic
					break                                                                  // one peer that thinks we're good is all we need.
				} else {
					// fmt.Println(cacheName, "- locally we think it's down, and", peerName, "says IsAvailable: ", peerCrStates.Caches[cacheName].IsAvailable, "down voting")
					downVotes++ // peerStates for this peer doesn't like it
				}
			}
		}
		if downVotes > len(peerStates) {
			// fmt.Println(cacheName, "-", downVotes, "down votes, setting to IsAvailable: false")
			combinedStates.Caches[cacheName] = peer.IsAvailable{IsAvailable: false}
		}
	}

	for deliveryServiceName, localDeliveryService := range localStates.Deliveryservice {
		deliveryService := peer.Deliveryservice{}
		if localDeliveryService.IsAvailable {
			deliveryService.IsAvailable = true
		}
		deliveryService.DisabledLocations = localDeliveryService.DisabledLocations

		for peerName, iPeerStates := range peerStates {
			peerDeliveryService, ok := iPeerStates.Deliveryservice[deliveryServiceName]
			if !ok {
				log.Printf("WARN local delivery service %s not found in peer %s\n", deliveryServiceName, peerName)
				continue
			}
			if peerDeliveryService.IsAvailable {
				deliveryService.IsAvailable = true
			}
			deliveryService.DisabledLocations = intersection(deliveryService.DisabledLocations, peerDeliveryService.DisabledLocations)
		}
		combinedStates.Deliveryservice[deliveryServiceName] = deliveryService
	}

	return combinedStates
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

func pruneHistory(history []interface{}, limit int) []interface{} {
	if len(history) > limit {
		history = history[1:]
	}

	return history
}
