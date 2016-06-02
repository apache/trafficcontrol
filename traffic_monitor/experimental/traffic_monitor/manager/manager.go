package manager

import (
	"crypto/tls"
	"fmt"
	"net/http"
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

//
// Kicks off the pollers and handlers
//
func Start(opsConfigFile string) {
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
	cacheHealthPoller := poller.HttpPoller{
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
	localStates := peer.Crstates{Caches: make(map[string]peer.IsAvailable)}    // this is the local state as discoverer by this traffic_monitor
	peerStates := make(map[string]peer.Crstates)                               // each peer's last state is saved in this map
	combinedStates := peer.Crstates{Caches: make(map[string]peer.IsAvailable)} // this is the result of combining the localStates and all the peerStates using the var ??

	for {
		select {
		case req := <-dr:
			var body []byte
			var err error

			if req.T == http_server.TR_CONFIG && toSession != nil && opsConfig.CdnName != "" {
				body, err = toSession.CRConfigRaw(opsConfig.CdnName)
			} else if req.T == http_server.TR_STATE_DERIVED {
				body, err = peer.CrStatesMarshall(combinedStates)
			} else if req.T == http_server.TR_STATE_SELF {
				body, err = peer.CrStatesMarshall(localStates)
			} else if req.T == http_server.CACHE_STATS {
				// TODO: add support for ?hc=N query param, stats=, wildcard, individual caches
				// add pp and date to the json:
				/*
					pp: "0=[my-ats-edge-cache-1], hc=[1]",
					date: "Thu Oct 09 20:28:36 UTC 2014"
				*/

				params := req.Parameters
				hc := 1
				_, exists := params["hc"]

				if exists {
					v, err := strconv.Atoi(params["hc"][0])

					if err == nil {
						hc = v
					}
				}

				body, err = cache.StatsMarshall(statHistory, hc)
			} else if req.T == http_server.DS_STATS {
			} else if req.T == http_server.EVENT_LOG {
			} else if req.T == http_server.PEER_STATES {
			} else if req.T == http_server.STAT_SUMMARY {
			} else if req.T == http_server.STATS {
			}

			if err == nil {
				req.C <- body
			}

			close(req.C)
		case oc := <-opsConfigFileHandler.OpsConfigChannel:
			var err error
			opsConfig = oc

			listenAddress := ":80" // default

			if opsConfig.HttpListener != "" {
				listenAddress = opsConfig.HttpListener
			}

			go http_server.Run(dr, listenAddress)

			toSession, err = traffic_ops.Login(opsConfig.Url, opsConfig.Username, opsConfig.Password, opsConfig.Insecure)

			if err != nil {
				fmt.Printf("MonitorConfigPoller: error instantiating Session with traffic_ops: %s\n", err)
			} else {
				monitorConfigPoller.OpsConfigChannel <- opsConfig // this is needed for cdnName
				monitorConfigPoller.SessionChannel <- toSession
			}
		case mc := <-monitorConfigPoller.ConfigChannel:
			monitorConfig = mc

			healthUrls := make(map[string]string)
			statUrls := make(map[string]string)
			peerUrls := make(map[string]string)
			caches := make(map[string]string)

			for _, srv := range monitorConfig.TrafficServer {
				caches[srv.HostName] = srv.Status

				if srv.Status == "ONLINE" {
					localStates.Caches[srv.HostName] = peer.IsAvailable{IsAvailable: true}
					continue
				} else if srv.Status == "OFFLINE" {
					localStates.Caches[srv.HostName] = peer.IsAvailable{IsAvailable: false}
					continue
				} else {
					_, exists := localStates.Caches[srv.HostName]

					// seed states with available = false until our polling cycle picks up a result
					if !exists {
						localStates.Caches[srv.HostName] = peer.IsAvailable{IsAvailable: false}
					}
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
		case healthResult := <-cacheHealthChannel:
			var prevResult cache.Result
			if len(healthHistory[healthResult.Id]) != 0 {
				prevResult = healthHistory[healthResult.Id][len(healthHistory[healthResult.Id])-1].(cache.Result)
			}
			health.GetVitals(&healthResult, &prevResult, &monitorConfig)
			healthHistory[healthResult.Id] = pruneHistory(append(healthHistory[healthResult.Id], healthResult), defaultMaxHistory)
			isAvailable, whyAvailable := health.EvalCache(healthResult, &monitorConfig)
			if localStates.Caches[healthResult.Id].IsAvailable != isAvailable {
				fmt.Println("Changing state for", healthResult.Id, " was:", prevResult.Available, " is now:", isAvailable, " because:", whyAvailable, " errors:", healthResult.Errors)
			}
			localStates.Caches[healthResult.Id] = peer.IsAvailable{IsAvailable: isAvailable}
		case stats := <-cacheStatChannel:
			statHistory[stats.Id] = pruneHistory(append(statHistory[stats.Id], stats), defaultMaxHistory)
		case crStatesResult := <-peerChannel:
			peerStates[crStatesResult.Id] = crStatesResult.PeerStats
			combinedStates = combineCrStates(peerStates, localStates)
		}
	}
}

// TODO JvD: add deliveryservice stuff
func combineCrStates(peerStates map[string]peer.Crstates, localStates peer.Crstates) peer.Crstates {

	combinedStates := peer.Crstates{Caches: make(map[string]peer.IsAvailable)}
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

	return combinedStates
}

func pruneHistory(history []interface{}, limit int) []interface{} {
	if len(history) > limit {
		history = history[1:]
	}

	return history
}
