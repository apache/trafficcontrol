package manager

import (
	"fmt"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/common/poller"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/config"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/peer"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"strings"
	"sync"
)

func CopyTrafficMonitorConfigMap(a *to.TrafficMonitorConfigMap) to.TrafficMonitorConfigMap {
	b := to.TrafficMonitorConfigMap{}
	b.TrafficServer = map[string]to.TrafficServer{}
	b.CacheGroup = map[string]to.TMCacheGroup{}
	b.Config = map[string]interface{}{}
	b.TrafficMonitor = map[string]to.TrafficMonitor{}
	b.DeliveryService = map[string]to.TMDeliveryService{}
	b.Profile = map[string]to.TMProfile{}
	for k, v := range a.TrafficServer {
		b.TrafficServer[k] = v
	}
	for k, v := range a.CacheGroup {
		b.CacheGroup[k] = v
	}
	for k, v := range a.Config {
		b.Config[k] = v
	}
	for k, v := range a.TrafficMonitor {
		b.TrafficMonitor[k] = v
	}
	for k, v := range a.DeliveryService {
		b.DeliveryService[k] = v
	}
	for k, v := range a.Profile {
		b.Profile[k] = v
	}
	return b
}

type TrafficMonitorConfigMapThreadsafe struct {
	monitorConfig *to.TrafficMonitorConfigMap
	m             *sync.RWMutex
}

func NewTrafficMonitorConfigMapThreadsafe() TrafficMonitorConfigMapThreadsafe {
	return TrafficMonitorConfigMapThreadsafe{monitorConfig: &to.TrafficMonitorConfigMap{}, m: &sync.RWMutex{}}
}

// Get returns the TrafficMonitorConfigMap. Callers MUST NOT modify, it is not threadsafe for mutation. If mutation is necessary, call CopyTrafficMonitorConfigMap().
func (t *TrafficMonitorConfigMapThreadsafe) Get() to.TrafficMonitorConfigMap {
	t.m.RLock()
	defer t.m.RUnlock()
	return *t.monitorConfig
}

// Set sets the TrafficMonitorConfigMap. This is only safe for one writer. This MUST NOT be called by multiple threads.
func (t *TrafficMonitorConfigMapThreadsafe) Set(c to.TrafficMonitorConfigMap) {
	t.m.Lock()
	*t.monitorConfig = c
	t.m.Unlock()
}

func StartMonitorConfigManager(
	monitorConfigPollChan <-chan to.TrafficMonitorConfigMap,
	localStates peer.CRStatesThreadsafe,
	statUrlSubscriber chan<- poller.HttpPollerConfig,
	healthUrlSubscriber chan<- poller.HttpPollerConfig,
	peerUrlSubscriber chan<- poller.HttpPollerConfig,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData StaticAppData,
) TrafficMonitorConfigMapThreadsafe {
	monitorConfig := NewTrafficMonitorConfigMapThreadsafe()
	go monitorConfigListen(monitorConfig,
		monitorConfigPollChan,
		localStates,
		statUrlSubscriber,
		healthUrlSubscriber,
		peerUrlSubscriber,
		cachesChangeSubscriber,
		cfg,
		staticAppData,
	)
	return monitorConfig
}

// TODO timing, and determine if the case, or its internal `for`, should be put in a goroutine
// TODO determine if subscribers take action on change, and change to mutexed objects if not.
func monitorConfigListen(
	monitorConfigTS TrafficMonitorConfigMapThreadsafe,
	monitorConfigPollChan <-chan to.TrafficMonitorConfigMap,
	localStates peer.CRStatesThreadsafe,
	statUrlSubscriber chan<- poller.HttpPollerConfig,
	healthUrlSubscriber chan<- poller.HttpPollerConfig,
	peerUrlSubscriber chan<- poller.HttpPollerConfig,
	cachesChangeSubscriber chan<- struct{},
	cfg config.Config,
	staticAppData StaticAppData,
) {
	for {
		select {
		case monitorConfig := <-monitorConfigPollChan:
			monitorConfigTS.Set(monitorConfig)
			healthUrls := map[string]string{}
			statUrls := map[string]string{}
			peerUrls := map[string]string{}
			caches := map[string]string{}

			for _, srv := range monitorConfig.TrafficServer {
				caches[srv.HostName] = srv.Status

				cacheName := enum.CacheName(srv.HostName)

				if srv.Status == "ONLINE" {
					localStates.SetCache(cacheName, peer.IsAvailable{IsAvailable: true})
					continue
				}
				if srv.Status == "OFFLINE" {
					localStates.SetCache(cacheName, peer.IsAvailable{IsAvailable: false})
					continue
				}
				// seed states with available = false until our polling cycle picks up a result
				if _, exists := localStates.Get().Caches[cacheName]; !exists {
					localStates.SetCache(cacheName, peer.IsAvailable{IsAvailable: false})
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
				if srv.HostName == staticAppData.Hostname {
					continue
				}
				if srv.Status != "ONLINE" {
					continue
				}
				// TODO: the URL should be config driven. -jse
				url := fmt.Sprintf("http://%s:%d/publish/CrStates?raw", srv.IP, srv.Port)
				peerUrls[srv.HostName] = url
			}

			statUrlSubscriber <- poller.HttpPollerConfig{Urls: statUrls, Interval: cfg.CacheStatPollingInterval}
			healthUrlSubscriber <- poller.HttpPollerConfig{Urls: healthUrls, Interval: cfg.CacheHealthPollingInterval}
			peerUrlSubscriber <- poller.HttpPollerConfig{Urls: peerUrls, Interval: cfg.PeerPollingInterval}

			for cacheName := range localStates.GetCaches() {
				if _, exists := monitorConfig.TrafficServer[string(cacheName)]; !exists {
					log.Warnf("Removing %s from localStates", cacheName)
					localStates.DeleteCache(cacheName)
				}
			}

			cachesChangeSubscriber <- struct{}{}

			// TODO because there are multiple writers to localStates.DeliveryService, there is a race condition, where MonitorConfig (this func) and HealthResultManager could write at the same time, and the HealthResultManager could overwrite a delivery service addition or deletion here. Probably the simplest and most performant fix would be a lock-free algorithm using atomic compare-and-swaps.
			for _, ds := range monitorConfig.DeliveryService {
				// since caches default to unavailable, also default DS false
				if _, exists := localStates.Get().Deliveryservice[enum.DeliveryServiceName(ds.XMLID)]; !exists {
					localStates.SetDeliveryService(enum.DeliveryServiceName(ds.XMLID), peer.Deliveryservice{IsAvailable: false, DisabledLocations: []enum.CacheName{}}) // important to initialize DisabledLocations, so JSON is `[]` not `null`
				}
			}
			for ds, _ := range localStates.Get().Deliveryservice {
				if _, exists := monitorConfig.DeliveryService[string(ds)]; !exists {
					localStates.DeleteDeliveryService(ds)
				}
			}
		}
	}
}
