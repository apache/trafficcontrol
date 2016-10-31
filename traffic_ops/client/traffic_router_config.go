/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package client

import (
	"encoding/json"
	"fmt"
)

// TRConfigResponse ...
type TRConfigResponse struct {
	Version  string              `json:"version"`
	Response TrafficRouterConfig `json:"response"`
}

// TrafficRouterConfig is the json unmarshalled without any changes
// note all structs are local to this file _except_ the TrafficRouterConfig struct.
type TrafficRouterConfig struct {
	TrafficServers   []TrafficServer        `json:"trafficServers,omitempty"`
	TrafficMonitors  []TrafficMonitor       `json:"trafficMonitors,omitempty"`
	TrafficRouters   []TrafficRouter        `json:"trafficRouters,omitempty"`
	CacheGroups      []TMCacheGroup         `json:"cacheGroups,omitempty"`
	DeliveryServices []TRDeliveryService    `json:"deliveryServices,omitempty"`
	Stats            map[string]interface{} `json:"stats,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
}

// TrafficRouterConfigMap ...
type TrafficRouterConfigMap struct {
	TrafficServer   map[string]TrafficServer
	TrafficMonitor  map[string]TrafficMonitor
	TrafficRouter   map[string]TrafficRouter
	CacheGroup      map[string]TMCacheGroup
	DeliveryService map[string]TRDeliveryService
	Config          map[string]interface{}
	Stat            map[string]interface{}
}

// TrafficServer ...
type TrafficServer struct {
	Profile          string              `json:"profile"`
	IP               string              `json:"ip"`
	Status           string              `json:"status"`
	CacheGroup       string              `json:"cacheGroup"`
	IP6              string              `json:"ip6"`
	Port             int                 `json:"port"`
	HostName         string              `json:"hostName"`
	FQDN             string              `json:"fqdn"`
	InterfaceName    string              `json:"interfaceName"`
	Type             string              `json:"type"`
	HashID           string              `json:"hashId"`
	DeliveryServices []tsdeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
}

type tsdeliveryService struct {
	Xmlid  string   `json:"xmlId"`
	Remaps []string `json:"remaps"`
}

// TrafficRouter ...
type TrafficRouter struct {
	Port     int    `json:"port"`
	IP6      string `json:"ip6"`
	IP       string `json:"ip"`
	FQDN     string `json:"fqdn"`
	Profile  string `json:"profile"`
	Location string `json:"location"`
	Status   string `json:"status"`
	APIPort  int    `json:"apiPort"`
}

// TMCacheGroup ...
// !!! Note the lowercase!!! this is local to this file, there's a CacheGroup definition in cachegroup.go!
type TMCacheGroup struct {
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
}

// Coordinates ...
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TRDeliveryService ...
// TODO JvD: move to deliveryservice.go ??
type TRDeliveryService struct {
	XMLID             string            `json:"xmlId"`
	Domains           []string          `json:"domains"`
	MissLocation      MissLocation      `json:"missCoordinates"`
	CoverageZoneOnly  bool              `json:"coverageZoneOnly"`
	MatchSets         []MatchSet        `json:"matchSets"`
	TTL               int               `json:"ttl"`
	TTLs              TTLS              `json:"ttls"`
	BypassDestination BypassDestination `json:"bypassDestination"`
	StatcDNSEntries   []StaticDNS       `json:"statitDnsEntries"`
	Soa               SOA               `json:"soa"`
}

// MissLocation ...
type MissLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitudef"`
}

// MatchSet ...
type MatchSet struct {
	Protocol  string      `json:"protocol"`
	MatchList []MatchList `json:"matchList"`
}

// MatchList ...
type MatchList struct {
	Regex     string `json:"regex"`
	MatchType string `json:"matchType"`
}

// BypassDestination ...
type BypassDestination struct {
	FQDN string `json:"fqdn"`
	Type string `json:"type"`
	Port int    `json:"Port"`
}

// TTLS ...
type TTLS struct {
	Arecord    int `json:"A"`
	SoaRecord  int `json:"SOA"`
	NsRecord   int `json:"NS"`
	AaaaRecord int `json:"AAAA"`
}

// StaticDNS ...
type StaticDNS struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

// SOA ...
type SOA struct {
	Admin   string `json:"admin"`
	Retry   int    `json:"retry"`
	Minimum int    `json:"minimum"`
	Refresh int    `json:"refresh"`
	Expire  int    `json:"expire"`
}

// TrafficRouterConfigMap Deprecated: use GetTrafficRouterConfigMap instead.
func (to *Session) TrafficRouterConfigMap(cdn string) (*TrafficRouterConfigMap, error) {
	cfg, _, err := to.GetTrafficRouterConfigMap(cdn)
	return cfg, err
}

// TrafficRouterConfigMap gets a bunch of maps
func (to *Session) GetTrafficRouterConfigMap(cdn string) (*TrafficRouterConfigMap, CacheHitStatus, error) {
	trConfig, cacheHitStatus, err := to.GetTrafficRouterConfig(cdn)
	if err != nil {
		return nil, CacheHitStatusInvalid, err
	}

	trConfigMap := TRTransformToMap(*trConfig)
	return &trConfigMap, cacheHitStatus, nil
}

// TrafficRouterConfig Deprecated: use GetTrafficRouterConfig instead.
func (to *Session) TrafficRouterConfig(cdn string) (*TrafficRouterConfig, error) {
	cfg, _, err := to.GetTrafficRouterConfig(cdn)
	return cfg, err
}

// GetTrafficRouterConfig gets the json arrays
func (to *Session) GetTrafficRouterConfig(cdn string) (*TrafficRouterConfig, CacheHitStatus, error) {
	url := fmt.Sprintf("/api/1.2/cdns/%s/configs/routing.json", cdn)
	body, cacheHitStatus, err := to.getBytesWithTTL(url, tmPollingInterval)
	if err != nil {
		return nil, CacheHitStatusInvalid, err
	}

	var data TRConfigResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, CacheHitStatusInvalid, err
	}
	return &data.Response, cacheHitStatus, nil
}

// TRTransformToMap ...
func TRTransformToMap(trConfig TrafficRouterConfig) TrafficRouterConfigMap {
	var tr TrafficRouterConfigMap
	tr.TrafficServer = make(map[string]TrafficServer)
	tr.TrafficMonitor = make(map[string]TrafficMonitor)
	tr.TrafficRouter = make(map[string]TrafficRouter)
	tr.CacheGroup = make(map[string]TMCacheGroup)
	tr.DeliveryService = make(map[string]TRDeliveryService)
	tr.Config = make(map[string]interface{})
	tr.Stat = make(map[string]interface{})

	for _, trServer := range trConfig.TrafficServers {
		tr.TrafficServer[trServer.HostName] = trServer
	}
	for _, trMonitor := range trConfig.TrafficMonitors {
		tr.TrafficMonitor[trMonitor.HostName] = trMonitor
	}
	for _, trServer := range trConfig.TrafficServers {
		tr.TrafficServer[trServer.HostName] = trServer
	}
	for _, trRouter := range trConfig.TrafficRouters {
		tr.TrafficRouter[trRouter.FQDN] = trRouter
	}
	for _, trCacheGroup := range trConfig.CacheGroups {
		tr.CacheGroup[trCacheGroup.Name] = trCacheGroup
	}
	for _, trDeliveryService := range trConfig.DeliveryServices {
		tr.DeliveryService[trDeliveryService.XMLID] = trDeliveryService
	}
	for trSettingKey, trSettingVal := range trConfig.Config {
		tr.Config[trSettingKey] = trSettingVal
	}
	for trStatKey, trStatVal := range trConfig.Stats {
		tr.Stat[trStatKey] = trStatVal
	}
	return tr
}
