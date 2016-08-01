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
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/health"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/http_server"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	traffic_ops "github.com/Comcast/traffic_control/traffic_ops/client"
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
	var toSession *traffic_ops.Session

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

	sessionChannel := make(chan *traffic_ops.Session)
	monitorConfigChannel := make(chan traffic_ops.TrafficMonitorConfigMap)
	monitorOpsConfigChannel := make(chan handler.OpsConfig)
	monitorConfigPoller := poller.MonitorConfigPoller{
		Interval:         defaultMonitorConfigPollingInterval,
		SessionChannel:   sessionChannel,
		ConfigChannel:    monitorConfigChannel,
		OpsConfigChannel: monitorOpsConfigChannel,
	}

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

	go opsConfigFileHandler.Listen()
	go opsConfigFilePoller.Poll()
	go monitorConfigPoller.Poll()
	go cacheHealthPoller.Poll()
	go cacheStatPoller.Poll()
	go peerPoller.Poll()

	dr := make(chan http_server.DataRequest)

	healthHistory := map[string][]interface{}{}
	statHistory := map[string][]interface{}{}

	var opsConfig handler.OpsConfig
	var monitorConfig traffic_ops.TrafficMonitorConfigMap
	localStates := peer.NewCRStates()        // this is the local state as discoverer by this traffic_monitor
	peerStates := map[string]peer.Crstates{} // each peer's last state is saved in this map
	combinedStates := peer.NewCRStates()     // this is the result of combining the localStates and all the peerStates using the var ??

	deliveryServiceServers := map[string][]string{}
	serverDeliveryServices := map[string]string{}
	serverTypes := map[string]ds.StatCacheType{}
	deliveryServiceTypes := map[string]ds.StatType{}
	deliveryServiceRegexes := map[string][]string{}
	serverCachegroups := map[string]string{}

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
	localCacheStatus := map[ds.CacheName]CacheAvailableStatus{}

	for {
		select {
		case req := <-dr:
			defer close(req.Response)

			var body []byte
			var err error

			switch req.Type {
			case http_server.TRConfig:
				if toSession == nil {
					err = fmt.Errorf("Unable to connect to Traffic Ops")
				} else if opsConfig.CdnName == "" {
					err = fmt.Errorf("No CDN Configured")
				} else {
					body, err = toSession.CRConfigRaw(opsConfig.CdnName)
				}
				if err != nil {
					err = fmt.Errorf("TR Config: %v", err)
				}
			case http_server.TRStateDerived:
				body, err = peer.CrStatesMarshall(combinedStates)
				if err != nil {
					err = fmt.Errorf("TR State (derived): %v", err)
				}
			case http_server.TRStateSelf:
				body, err = peer.CrStatesMarshall(localStates)
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
				body, err = json.Marshal(createApiPeerStates(peerStates))
			case http_server.StatSummary:
				body = []byte("TODO implement")
			case http_server.Stats:
				body, err = getStats(staticAppData, cacheHealthPoller.Config.Interval, lastHealthDurations, fetchCount, healthIteration, errorCount)
				if err != nil {
					err = fmt.Errorf("Stats: %v", err)
				}
			case http_server.ConfigDoc:
				opsConfigCopy := opsConfig
				// if the password is blank, leave it blank, so callers can see it's missing.
				if opsConfigCopy.Password != "" {
					opsConfigCopy.Password = "*****"
				}
				body, err = json.Marshal(opsConfigCopy)
				if err != nil {
					err = fmt.Errorf("Config Doc: %v", err)
				}
			case http_server.APICacheCount: // TODO determine if this should use peerStates
				body = []byte(strconv.Itoa(len(localStates.Caches)))
			case http_server.APICacheAvailableCount:
				body = []byte(strconv.Itoa(cacheAvailableCount(localStates.Caches)))
			case http_server.APICacheDownCount:
				body = []byte(strconv.Itoa(cacheDownCount(localStates.Caches)))
			case http_server.APIVersion:
				s := "traffic_monitor-" + staticAppData.Version + "."
				if len(staticAppData.GitRevision) > 6 {
					s += staticAppData.GitRevision[:6]
				} else {
					s += staticAppData.GitRevision
				}
				body = []byte(s)
			case http_server.APITrafficOpsURI:
				body = []byte(opsConfig.Url)
			case http_server.APICacheStates:
				body, err = json.Marshal(createCacheStatuses(serverTypes, statHistory, lastHealthDurations, localStates.Caches, lastKbpsStats, localCacheStatus))
			default:
				err = fmt.Errorf("Unknown Request Type: %v", req.Type)
			}

			if err != nil {
				errorCount++
				log.Printf("ERROR Request Error: %v\n", err)
			} else {
				req.Response <- body
			}
		case oc := <-opsConfigFileHandler.OpsConfigChannel:
			var err error
			opsConfig = oc

			listenAddress := ":80" // default

			if opsConfig.HttpListener != "" {
				listenAddress = opsConfig.HttpListener
			}

			handleErr := func(err error) {
				errorCount++
				log.Printf("%v\n", err)
			}

			err = http_server.Run(dr, listenAddress)
			if err != nil {
				handleErr(fmt.Errorf("MonitorConfigPoller: error creating HTTP server: %s\n", err))
				continue
			}

			toSession, err = traffic_ops.Login(opsConfig.Url, opsConfig.Username, opsConfig.Password, opsConfig.Insecure)
			if err != nil {
				handleErr(fmt.Errorf("MonitorConfigPoller: error instantiating Session with traffic_ops: %s\n", err))
				continue
			}

			deliveryServiceServers, serverDeliveryServices, err = getDeliveryServiceServers(toSession, opsConfig.CdnName)
			if err != nil {
				handleErr(fmt.Errorf("Error getting delivery service servers from Traffic Ops: %v\n", err))
				continue
			}

			deliveryServiceTypes, err = getDeliveryServiceTypes(toSession, opsConfig.CdnName)
			if err != nil {
				handleErr(fmt.Errorf("Error getting delivery service types from Traffic Ops: %v\n", err))
				continue
			}

			deliveryServiceRegexes, err = getDeliveryServiceRegexes(toSession, opsConfig.CdnName)
			if err != nil {
				handleErr(fmt.Errorf("Error getting delivery service regexes from Traffic Ops: %v\n", err))
				continue
			}

			serverCachegroups, err = getServerCachegroups(toSession, opsConfig.CdnName)
			if err != nil {
				handleErr(fmt.Errorf("Error getting server cachegroups from Traffic Ops: %v\n", err))
				continue
			}

			serverTypes, err = getServerTypes(toSession, opsConfig.CdnName)
			if err != nil {
				handleErr(fmt.Errorf("Error getting server types from Traffic Ops: %v\n", err))
				continue
			}

			// This must be in a goroutine, because the monitorConfigPoller tick sends to a channel this select listens for. Thus, if we block on sends to the monitorConfigPoller, we have a livelock race condition.
			go func() {
				monitorConfigPoller.OpsConfigChannel <- opsConfig // this is needed for cdnName
				monitorConfigPoller.SessionChannel <- toSession
			}()
		case monitorConfig = <-monitorConfigPoller.ConfigChannel:
			healthUrls := map[string]string{}
			statUrls := map[string]string{}
			peerUrls := map[string]string{}
			caches := map[string]string{}

			for _, srv := range monitorConfig.TrafficServer {
				caches[srv.HostName] = srv.Status

				if srv.Status == "ONLINE" {
					localStates.Caches[srv.HostName] = peer.IsAvailable{IsAvailable: true}
					continue
				}
				if srv.Status == "OFFLINE" {
					localStates.Caches[srv.HostName] = peer.IsAvailable{IsAvailable: false}
					continue
				}
				// seed states with available = false until our polling cycle picks up a result
				if _, exists := localStates.Caches[srv.HostName]; !exists {
					localStates.Caches[srv.HostName] = peer.IsAvailable{IsAvailable: false}
				}

				url := monitorConfig.Profile[srv.Profile].Parameters.HealthPollingURL
				r := strings.NewReplacer(
					"${hostname}", srv.FQDN,
					"${interface_name}", srv.InterfaceName,
					"application=system", "application=plugin.remap",
					"application=", "application=plugin.remap",
				)
				url = r.Replace(url)
				healthUrls[srv.HostName] = url
				r = strings.NewReplacer("application=plugin.remap", "application=")
				url = r.Replace(url)
				statUrls[srv.HostName] = url
			}

			for _, srv := range monitorConfig.TrafficMonitor {
				if srv.Status != "ONLINE" {
					continue
				}
				// TODO: the URL should be config driven. -jse
				url := fmt.Sprintf("http://%s:%d/publish/CrStates?raw", srv.IP, srv.Port)
				peerUrls[srv.HostName] = url
			}

			cacheStatPoller.ConfigChannel <- poller.HttpPollerConfig{Urls: statUrls, Interval: defaultCacheStatPollingInterval}
			cacheHealthPoller.ConfigChannel <- poller.HttpPollerConfig{Urls: healthUrls, Interval: defaultCacheHealthPollingInterval}
			peerPoller.ConfigChannel <- poller.HttpPollerConfig{Urls: peerUrls, Interval: defaultPeerPollingInterval}

			for k := range localStates.Caches {
				_, exists := monitorConfig.TrafficServer[k]

				if !exists {
					fmt.Printf("Warning: removing %s from localStates", k)
					delete(localStates.Caches, k)
				}
			}

			addStateDeliveryServices(monitorConfig, localStates.Deliveryservice)
		case i := <-cacheHealthTick:
			healthIteration = i
		case healthResult := <-cacheHealthChannel:
			fetchCount++
			var prevResult cache.Result
			if len(healthHistory[healthResult.Id]) != 0 {
				prevResult = healthHistory[healthResult.Id][len(healthHistory[healthResult.Id])-1].(cache.Result)
			}
			health.GetVitals(&healthResult, &prevResult, &monitorConfig)
			healthHistory[healthResult.Id] = pruneHistory(append(healthHistory[healthResult.Id], healthResult), defaultMaxHistory)
			isAvailable, whyAvailable := health.EvalCache(healthResult, &monitorConfig)
			if localStates.Caches[healthResult.Id].IsAvailable != isAvailable {
				fmt.Println("Changing state for", healthResult.Id, " was:", prevResult.Available, " is now:", isAvailable, " because:", whyAvailable, " errors:", healthResult.Errors)
				e := Event{Index: eventIndex, Time: time.Now().Unix(), Description: whyAvailable, Name: healthResult.Id, Hostname: healthResult.Id, Type: serverTypes[healthResult.Id].String(), Available: isAvailable}
				events = append([]Event{e}, events...)
				if len(events) > maxEvents {
					events = events[:maxEvents-1]
				}
				eventIndex++
			}

			localCacheStatus[ds.CacheName(healthResult.Id)] = CacheAvailableStatus{Available: isAvailable, Status: monitorConfig.TrafficServer[healthResult.Id].Status} // TODO move within localStates
			localStates.Caches[healthResult.Id] = peer.IsAvailable{IsAvailable: isAvailable}
			calculateDeliveryServiceState(deliveryServiceServers, localStates.Caches, localStates.Deliveryservice)

			// TODO determine if we should combineCrStates() here

			now := time.Now()

			var err error
			dsStats, lastKbpsStats, err = ds.CreateStats(statHistory, deliveryServiceServers, serverDeliveryServices, deliveryServiceTypes, deliveryServiceRegexes, serverCachegroups, serverTypes, combinedStates, lastKbpsStats, now)
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
			peerStates[crStatesResult.Id] = crStatesResult.PeerStats
			combinedStates = combineCrStates(peerStates, localStates)
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

func createCacheStatuses(cacheTypes map[string]ds.StatCacheType, statHistory map[string][]interface{}, lastHealthDurations map[string]time.Duration, cacheStates map[string]peer.IsAvailable, lastKbpsStats ds.StatsLastKbps, localCacheStatus map[ds.CacheName]CacheAvailableStatus) map[ds.CacheName]CacheStatus {
	conns := createCacheConnections(statHistory)
	statii := map[ds.CacheName]CacheStatus{}
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
		kbpsVal, ok := lastKbpsStats.Caches[ds.CacheName(cacheName)]
		if !ok {
			log.Printf("WARNING DEBUGQ cache not in last kbps cache %s\n", cacheName)
		} else {
			kbps = &kbpsVal.Kbps
		}

		var connections *int64
		connectionsVal, ok := conns[ds.CacheName(cacheName)]
		if !ok {
			log.Printf("WARNING DEBUGQ cache not in connections %s\n", cacheName)
		} else {
			connections = &connectionsVal
		}

		var status *string
		statusVal, ok := localCacheStatus[ds.CacheName(cacheName)]
		if !ok {
			log.Printf("WARNING DEBUGQ cache not in statuses %s\n", cacheName)
		} else {
			statusString := statusVal.Status + " - "
			if localCacheStatus[ds.CacheName(cacheName)].Available {
				statusString += "available"
			} else {
				statusString += "unavailable"
			}
			status = &statusString
		}

		cacheTypeStr := string(cacheType)
		statii[ds.CacheName(cacheName)] = CacheStatus{Type: &cacheTypeStr, LoadAverage: loadAverage, QueryTimeMilliseconds: queryTime, BandwidthKbps: kbps, ConnectionCount: connections, Status: status}
	}
	return statii
}

// TODO: run these in goroutines
func createCacheHealthStatuses(statHistory map[string][]interface{}) map[ds.CacheName]string {
	statuses := map[ds.CacheName]string{}
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

			statuses[ds.CacheName(server)] = v
		}
	}
	return statuses
}

func createCacheConnections(statHistory map[string][]interface{}) map[ds.CacheName]int64 {
	conns := map[ds.CacheName]int64{}
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

			conns[ds.CacheName(server)] = int64(v)
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
	Peers map[TrafficMonitorName]map[ds.CacheName][]CacheState `json:"peers"`
}

type CacheState struct {
	Value bool `json:"value"`
}

func createApiPeerStates(peerStates map[string]peer.Crstates) ApiPeerStates {
	apiPeerStates := ApiPeerStates{Peers: map[TrafficMonitorName]map[ds.CacheName][]CacheState{}}

	for peer, state := range peerStates {
		if _, ok := apiPeerStates.Peers[TrafficMonitorName(peer)]; !ok {
			apiPeerStates.Peers[TrafficMonitorName(peer)] = map[ds.CacheName][]CacheState{}
		}
		peerState := apiPeerStates.Peers[TrafficMonitorName(peer)]
		for cache, available := range state.Caches {
			peerState[ds.CacheName(cache)] = []CacheState{CacheState{Value: available.IsAvailable}}
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
func addStateDeliveryServices(mc traffic_ops.TrafficMonitorConfigMap, deliveryServices map[string]peer.Deliveryservice) {
	for _, ds := range mc.DeliveryService {
		// since caches default to unavailable, also default DS false
		deliveryServices[ds.XMLID] = peer.Deliveryservice{}
	}
}

// getDeliveryServiceServers gets the servers on each delivery services, for the given CDN, from Traffic Ops.
// Returns a map[deliveryService][]server, and a map[server]deliveryService
func getDeliveryServiceServers(to *traffic_ops.Session, cdn string) (map[string][]string, map[string]string, error) {
	dsServers := map[string][]string{}
	serverDs := map[string]string{}

	crcData, err := to.CRConfigRaw(cdn)
	if err != nil {
		return nil, nil, err
	}
	type CrConfig struct {
		ContentServers map[string]struct {
			DeliveryServices map[string][]string `json:"deliveryServices"`
		} `json:"contentServers"`
	}
	var crc CrConfig
	if err := json.Unmarshal(crcData, &crc); err != nil {
		return nil, nil, err
	}

	for serverName, serverData := range crc.ContentServers {
		for deliveryServiceName, _ := range serverData.DeliveryServices {
			dsServers[deliveryServiceName] = append(dsServers[deliveryServiceName], serverName)
			serverDs[serverName] = deliveryServiceName
		}
	}
	return dsServers, serverDs, nil
}

// getDeliveryServiceRegexes gets the regexes of each delivery service, for the given CDN, from Traffic Ops.
// Returns a map[deliveryService][]regex.
func getDeliveryServiceRegexes(to *traffic_ops.Session, cdn string) (map[string][]string, error) {
	dsRegexes := map[string][]string{}

	crcData, err := to.CRConfigRaw(cdn)
	if err != nil {
		return nil, err
	}
	type CrConfig struct {
		DeliveryServices map[string]struct {
			Matchsets []struct {
				MatchList []struct {
					Regex string `json:"regex"`
				} `json:"matchlist"`
			} `json:"matchsets"`
		} `json:"deliveryServices"`
	}
	var crc CrConfig
	if err := json.Unmarshal(crcData, &crc); err != nil {
		return nil, err
	}

	for dsName, dsData := range crc.DeliveryServices {
		if len(dsData.Matchsets) < 1 {
			return nil, fmt.Errorf("CRConfig missing regex for '%s'", dsName)
		}
		for _, matchset := range dsData.Matchsets {
			if len(matchset.MatchList) < 1 {
				return nil, fmt.Errorf("CRConfig missing Regex for '%s'", dsName)
			}
			dsRegexes[dsName] = append(dsRegexes[dsName], matchset.MatchList[0].Regex)
		}
	}
	return dsRegexes, nil
}

// getServerCachegroups gets the cachegroup of each ATS Edge+Mid Cache server, for the given CDN, from Traffic Ops.
// Returns a map[server]cachegroup.
func getServerCachegroups(to *traffic_ops.Session, cdn string) (map[string]string, error) {
	serverCachegroups := map[string]string{}

	crcData, err := to.CRConfigRaw(cdn)
	if err != nil {
		return nil, err
	}
	type CrConfig struct {
		ContentServers map[string]struct {
			CacheGroup string `json:"cacheGroup"`
		} `json:"contentServers"`
	}
	var crc CrConfig
	if err := json.Unmarshal(crcData, &crc); err != nil {
		return nil, err
	}

	for server, serverData := range crc.ContentServers {
		serverCachegroups[server] = serverData.CacheGroup
	}
	return serverCachegroups, nil
}

// getServerTypes gets the cache type of each ATS Edge+Mid Cache server, for the given CDN, from Traffic Ops.
func getServerTypes(to *traffic_ops.Session, cdn string) (map[string]ds.StatCacheType, error) {
	serverTypes := map[string]ds.StatCacheType{}

	crcData, err := to.CRConfigRaw(cdn)
	if err != nil {
		return nil, err
	}
	type CrConfig struct {
		ContentServers map[string]struct {
			Type string `json:"type"`
		} `json:"contentServers"`
	}
	var crc CrConfig
	if err := json.Unmarshal(crcData, &crc); err != nil {
		return nil, err
	}

	for server, serverData := range crc.ContentServers {
		t := ds.StatCacheTypeFromString(serverData.Type)
		if t == ds.StatCacheTypeInvalid {
			return nil, fmt.Errorf("getServerTypes CRConfig unknown type for '%s': '%s'", server, serverData.Type)
		}
		serverTypes[server] = t
	}
	return serverTypes, nil
}

func getDeliveryServiceTypes(to *traffic_ops.Session, cdn string) (map[string]ds.StatType, error) {
	dsTypes := map[string]ds.StatType{}

	crcData, err := to.CRConfigRaw(cdn)
	if err != nil {
		return nil, err
	}
	type CrConfig struct {
		DeliveryServices map[string]struct {
			Matchsets []struct {
				Protocol string `json:"protocol"`
			} `json:"matchsets"`
		} `json:"deliveryServices"`
	}
	var crc CrConfig
	if err := json.Unmarshal(crcData, &crc); err != nil {
		return nil, fmt.Errorf("Error unmarshalling CRConfig: %v", err)
	}

	for dsName, dsData := range crc.DeliveryServices {
		if len(dsData.Matchsets) < 1 {
			return nil, fmt.Errorf("CRConfig missing protocol for '%s'", dsName)
		}
		dsTypeStr := dsData.Matchsets[0].Protocol
		dsType := ds.StatTypeFromString(dsTypeStr)
		if dsType == ds.StatTypeInvalid {
			return nil, fmt.Errorf("CRConfig unknowng protocol for '%s': '%s'", dsName, dsTypeStr)
		}
		dsTypes[dsName] = dsType
	}
	return dsTypes, nil
}

// calculateDeliveryServiceState calculates the state of delivery services from the new cache state data `cacheState` and the CRConfig data `deliveryServiceServers` and puts the calculated state in the outparam `deliveryServiceStates`
func calculateDeliveryServiceState(deliveryServiceServers map[string][]string, cacheState map[string]peer.IsAvailable, deliveryServiceStates map[string]peer.Deliveryservice) {
	for deliveryServiceName, deliveryServiceState := range deliveryServiceStates {
		if _, ok := deliveryServiceServers[deliveryServiceName]; !ok {
			// log.Printf("ERROR CRConfig does not have delivery service %s, but traffic monitor poller does; skipping\n", deliveryServiceName)
			continue
		}
		deliveryServiceState.IsAvailable = false
		for _, server := range deliveryServiceServers[deliveryServiceName] {
			if cacheState[server].IsAvailable {
				deliveryServiceState.IsAvailable = true
			} else {
				deliveryServiceState.DisabledLocations = append(deliveryServiceState.DisabledLocations, server)
			}
		}
		deliveryServiceStates[deliveryServiceName] = deliveryServiceState
	}
}

// TODO JvD: add deliveryservice stuff
func combineCrStates(peerStates map[string]peer.Crstates, localStates peer.Crstates) peer.Crstates {
	combinedStates := peer.NewCRStates()
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
