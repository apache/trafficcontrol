package manager

import (
	"fmt"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/common/poller"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/log"
	"github.com/Comcast/traffic_control/traffic_monitor/experimental/traffic_monitor/peer"
	to "github.com/Comcast/traffic_control/traffic_ops/client"
	"strings"
	"sync"
)

func copyTMConfig(a to.TrafficMonitorConfigMap) to.TrafficMonitorConfigMap {
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

// TrafficMonitorConfigMap ...
type TrafficMonitorConfigMap struct {
}

type TrafficMonitorConfigMapThreadsafe struct {
	monitorConfig *to.TrafficMonitorConfigMap
	m             *sync.Mutex
}

func NewTrafficMonitorConfigMapThreadsafe() TrafficMonitorConfigMapThreadsafe {
	return TrafficMonitorConfigMapThreadsafe{monitorConfig: &to.TrafficMonitorConfigMap{}, m: &sync.Mutex{}}
}

func (t *TrafficMonitorConfigMapThreadsafe) Get() to.TrafficMonitorConfigMap {
	t.m.Lock()
	defer func() {
		t.m.Unlock()
	}()
	return copyTMConfig(*t.monitorConfig)
}

func (t *TrafficMonitorConfigMapThreadsafe) Set(newMonitorConfig to.TrafficMonitorConfigMap) {
	t.m.Lock()
	*t.monitorConfig = copyTMConfig(newMonitorConfig)
	t.m.Unlock()
}

func StartMonitorConfigManager(monitorConfigPollChan <-chan to.TrafficMonitorConfigMap, localStates peer.CRStatesThreadsafe, statUrlSubscriber chan<- poller.HttpPollerConfig, healthUrlSubscriber chan<- poller.HttpPollerConfig, peerUrlSubscriber chan<- poller.HttpPollerConfig) TrafficMonitorConfigMapThreadsafe {
	monitorConfig := NewTrafficMonitorConfigMapThreadsafe()
	go monitorConfigListen(monitorConfig, monitorConfigPollChan, localStates, statUrlSubscriber, healthUrlSubscriber, peerUrlSubscriber)
	return monitorConfig
}

// TODO timing, and determine if the case, or its internal `for`, should be put in a goroutine
// TODO determine if subscribers take action on change, and change to mutexed objects if not.
func monitorConfigListen(monitorConfigTS TrafficMonitorConfigMapThreadsafe, monitorConfigPollChan <-chan to.TrafficMonitorConfigMap, localStates peer.CRStatesThreadsafe, statUrlSubscriber chan<- poller.HttpPollerConfig, healthUrlSubscriber chan<- poller.HttpPollerConfig, peerUrlSubscriber chan<- poller.HttpPollerConfig) {
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

				if srv.Status == "ONLINE" {
					localStates.SetCache(srv.HostName, peer.IsAvailable{IsAvailable: true})
					continue
				}
				if srv.Status == "OFFLINE" {
					localStates.SetCache(srv.HostName, peer.IsAvailable{IsAvailable: false})
					continue
				}
				// seed states with available = false until our polling cycle picks up a result
				if _, exists := localStates.Get().Caches[srv.HostName]; !exists {
					localStates.SetCache(srv.HostName, peer.IsAvailable{IsAvailable: false})
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

			statUrlSubscriber <- poller.HttpPollerConfig{Urls: statUrls, Interval: defaultCacheStatPollingInterval}
			healthUrlSubscriber <- poller.HttpPollerConfig{Urls: healthUrls, Interval: defaultCacheHealthPollingInterval}
			peerUrlSubscriber <- poller.HttpPollerConfig{Urls: peerUrls, Interval: defaultPeerPollingInterval}

			for k := range localStates.GetCaches() {
				if _, exists := monitorConfig.TrafficServer[k]; !exists {
					log.Warnf("Removing %s from localStates", k)
					localStates.DeleteCache(k)
				}
			}

			addStateDeliveryServices(monitorConfig, localStates.Get().Deliveryservice)
		}
	}
}

// addStateDeliveryServices adds delivery services in `mc` as keys in `deliveryServices`, with empty Deliveryservice values.
// TODO add disabledLocations
func addStateDeliveryServices(mc to.TrafficMonitorConfigMap, deliveryServices map[string]peer.Deliveryservice) {
	for _, ds := range mc.DeliveryService {
		// since caches default to unavailable, also default DS false
		deliveryServices[ds.XMLID] = peer.Deliveryservice{}
	}
}
