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

	healthHistory := make(map[string][]interface{})
	statHistory := make(map[string][]interface{})

	var opsConfig handler.OpsConfig
	var monitorConfig traffic_ops.TrafficMonitorConfigMap
	localStates := peer.Crstates{Caches: make(map[string]peer.IsAvailable), Deliveryservice: make(map[string]peer.Deliveryservice)}    // this is the local state as discoverer by this traffic_monitor
	peerStates := make(map[string]peer.Crstates)                                                                                       // each peer's last state is saved in this map
	combinedStates := peer.Crstates{Caches: make(map[string]peer.IsAvailable), Deliveryservice: make(map[string]peer.Deliveryservice)} // this is the result of combining the localStates and all the peerStates using the var ??

	deliveryServiceServers := map[string][]string{}
	serverTypes := map[string]string{}

	// TODO put stat data in a struct, for brevity
	lastHealthEndTimes := map[string]time.Time{}
	lastHealthDurations := map[string]time.Duration{}
	fetchCount := uint64(0) // note this is the number of individual caches fetched from, not the number of times all the caches were polled.
	healthIteration := uint64(0)
	errorCount := uint64(0)
	events := []Event{}
	eventIndex := uint64(0)
	for {
		select {
		case req := <-dr:
			defer close(req.C)

			var body []byte
			var err error

			switch req.T {
			case http_server.TR_CONFIG:
				if toSession != nil && opsConfig.CdnName != "" {
					body, err = toSession.CRConfigRaw(opsConfig.CdnName)
				}
			case http_server.TR_STATE_DERIVED:
				body, err = peer.CrStatesMarshall(combinedStates)
			case http_server.TR_STATE_SELF:
				body, err = peer.CrStatesMarshall(localStates)
			case http_server.CACHE_STATS:
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
			case http_server.DS_STATS:
				body = []byte("TODO implement")
			case http_server.EVENT_LOG:
				body, err = json.Marshal(JSONEvents{Events: events})
			case http_server.PEER_STATES:
				body = []byte("TODO implement")
			case http_server.STAT_SUMMARY:
				body = []byte("TODO implement")
			case http_server.STATS:
				body, err = getStats(staticAppData, cacheHealthPoller.Config.Interval, lastHealthDurations, fetchCount, healthIteration, errorCount)
				if err != nil {
					// TODO send error to client
					errorCount++
					log.Printf("ERROR getting stats %v\n", err)
					continue
				}
			case http_server.CONFIG_DOC:
				opsConfigCopy := opsConfig
				// if the password is blank, leave it blank, so callers can see it's missing.
				if opsConfigCopy.Password != "" {
					opsConfigCopy.Password = "*****"
				}
				body, err = json.Marshal(opsConfigCopy)
			default:
				body = []byte("TODO error message")
			}
			req.C <- body
		case oc := <-opsConfigFileHandler.OpsConfigChannel:
			var err error
			opsConfig = oc

			listenAddress := ":80" // default

			if opsConfig.HttpListener != "" {
				listenAddress = opsConfig.HttpListener
			}

			err = http_server.Run(dr, listenAddress)
			if err != nil {
				errorCount++
				log.Printf("MonitorConfigPoller: error creating HTTP server: %s\n", err)
				continue
			}

			toSession, err = traffic_ops.Login(opsConfig.Url, opsConfig.Username, opsConfig.Password, opsConfig.Insecure)
			if err != nil {
				errorCount++
				log.Printf("MonitorConfigPoller: error instantiating Session with traffic_ops: %s\n", err)
				continue
			}

			deliveryServiceServers, err = getDeliveryServiceServers(toSession, opsConfig.CdnName)
			if err != nil {
				errorCount++
				log.Printf("Error getting delivery service servers from Traffic Ops: %v\n", err)
				continue
			}

			serverTypes, err = getServerTypes(toSession, opsConfig.CdnName)
			if err != nil {
				errorCount++
				log.Printf("Error getting server types from Traffic Ops: %v\n", err)
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
				e := Event{
					Index:       eventIndex,
					Time:        time.Now().Unix(),
					Description: whyAvailable,
					Name:        healthResult.Id,
					Hostname:    healthResult.Id,
					Type:        serverTypes[healthResult.Id],
					Available:   isAvailable,
				}
				events = append([]Event{e}, events...)
				if len(events) > maxEvents {
					events = events[:maxEvents-1]
				}
				eventIndex++
			}
			localStates.Caches[healthResult.Id] = peer.IsAvailable{IsAvailable: isAvailable}
			calculateDeliveryServiceState(deliveryServiceServers, localStates.Caches, localStates.Deliveryservice)

			if lastHealthStart, ok := lastHealthEndTimes[healthResult.Id]; ok {
				lastHealthDurations[healthResult.Id] = time.Since(lastHealthStart)
			}
			lastHealthEndTimes[healthResult.Id] = time.Now()

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
	s.ErrorCount = errorCount // TODO implement
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

func getServerTypes(to *traffic_ops.Session, cdn string) (map[string]string, error) {
	// This is efficient (with getDeliveryServiceServers) because the traffic_ops client caches its result.
	// Were that not the case, these functions could be refactored to only call traffic_ops.Session.CRConfigRaw() once.

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

	serverTypes := map[string]string{}
	for serverName, serverData := range crc.ContentServers {
		serverTypes[serverName] = serverData.Type
	}
	return serverTypes, nil
}

// TODO add disabledLocations
func addStateDeliveryServices(mc traffic_ops.TrafficMonitorConfigMap, deliveryServices map[string]peer.Deliveryservice) {
	for _, ds := range mc.DeliveryService {
		// since caches default to unavailable, also default DS false
		deliveryServices[ds.XMLID] = peer.Deliveryservice{}
	}
}

// getDeliveryServiceServers returns a map[deliveryService][]server
func getDeliveryServiceServers(to *traffic_ops.Session, cdn string) (map[string][]string, error) {
	dsServers := map[string][]string{}

	crcData, err := to.CRConfigRaw(cdn)
	if err != nil {
		return nil, err
	}
	type CrConfig struct {
		ContentServers map[string]struct {
			DeliveryServices map[string][]string `json:"deliveryServices"`
		} `json:"contentServers"`
	}
	var crc CrConfig
	if err := json.Unmarshal(crcData, &crc); err != nil {
		return nil, err
	}

	for serverName, serverData := range crc.ContentServers {
		for deliveryServiceName, _ := range serverData.DeliveryServices {
			dsServers[deliveryServiceName] = append(dsServers[deliveryServiceName], serverName)
		}
	}
	return dsServers, nil
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
	combinedStates := peer.Crstates{Caches: make(map[string]peer.IsAvailable), Deliveryservice: make(map[string]peer.Deliveryservice)}
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
