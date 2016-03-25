/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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

// TrafficRouterConfig is the json unmarshalled without any changes
// note all structs are local to this file _except_ the TrafficRouterConfig struct.
type TrafficRouterConfig struct {
	TrafficServers   []TrafficServer        `json:"trafficServers"`
	TrafficMonitors  []trafficMonitor       `json:"trafficMonitors"`
	TrafficRouters   []trafficRouter        `json:"trafficRouters"`
	CacheGroups      []cacheGroup           `json:"cacheGroups"`
	DeliveryServices []deliveryService      `json:"deliveryServices"`
	Stats            map[string]interface{} `json:"stats"`
	Config           map[string]interface{} `json:"config"`
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
	Fqdn             string              `json:"fqdn"`
	InterfaceName    string              `json:"interfaceName"`
	Type             string              `json:"type"`
	HashID           string              `json:"hashId"`
	DeliveryServices []tsdeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
}

type tsdeliveryService struct {
	Xmlid  string   `json:"xmlId"`
	Remaps []string `json:"remaps"`
}

type trafficMonitor struct {
	Port     int    `json:"port"`
	IP6      string `json:"ip6"`
	IP       string `json:"ip"`
	HostName string `json:"hostName"`
	Fqdn     string `json:"fqdn"`
	Profile  string `json:"profile"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

type trafficRouter struct {
	Port     int    `json:"port"`
	IP6      string `json:"ip6"`
	IP       string `json:"ip"`
	Fqdn     string `json:"fqdn"`
	Profile  string `json:"profile"`
	Location string `json:"location"`
	Status   string `json:"status"`
	APIPort  int    `json:"apiPort"`
}

// !!! Note the lowercase!!! this is local to this file, there's a CacheGroup definition in cachegroup.go!
type cacheGroup struct {
	Name        string      `json:"name"`
	Coordinates coordinates `json:"coordinates"`
}

type coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TODO JvD: move to deliveryservice.go ??
type deliveryService struct {
	XMLID             string            `json:"xmlId"`
	MissLocation      missLocation      `json:"missLocation"`
	CoverageZoneOnly  bool              `json:"coverageZoneOnly"`
	MatchSets         []matchSet        `json:"matchSets"`
	TTL               int               `json:"ttl"`
	TTLs              ttls              `json:"ttls"`
	BypassDestination bypassDestination `json:"bypassDestination"`
	StatcDNSEntries   []staticDNS       `json:"statitDnsEntries"`
	Soa               soa               `json:"soa"`
}

type missLocation struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type matchSet struct {
	Protocol  string      `json:"protocol"`
	MatchList []matchList `json:"matchList"`
}

type matchList struct {
	Regex     string `json:"regex"`
	MatchType string `json:"matchType"`
}

type bypassDestination struct {
	Fqdn string `json:"fqdn"`
	Type string `json:"type"`
	Port int    `json:"Port"`
}

type ttls struct {
	Arecord    int `json:"A"`
	SoaRecord  int `json:"SOA"`
	NsRecord   int `json:"NS"`
	AaaaRecord int `json:"AAAA"`
}

type staticDNS struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type soa struct {
	Admin   string `json:"admin"`
	Retry   int    `json:"retry"`
	Minimum int    `json:"minimum"`
	Refresh int    `json:"refresh"`
	Expire  int    `json:"expire"`
}

// TrafficRouterConfigMap ...
type TrafficRouterConfigMap struct {
	TrafficServer   map[string]TrafficServer
	TrafficMonitor  map[string]trafficMonitor
	TrafficRouter   map[string]trafficRouter
	CacheGroup      map[string]cacheGroup
	DeliveryService map[string]deliveryService
	Config          map[string]interface{}
	Stat            map[string]interface{}
}

// TrafficRouterConfigMap gets a bunch of maps
func (to *Session) TrafficRouterConfigMap(cdn string) (*TrafficRouterConfigMap, error) {
	trConfig, err := to.TrafficRouterConfig(cdn)
	if err != nil {
		return nil, err
	}
	trConfigMap := trTransformToMap(*trConfig)
	return &trConfigMap, nil
}

// TrafficRouterConfigRaw ...
func (to *Session) TrafficRouterConfigRaw(cdn string) ([]byte, error) {
	url := fmt.Sprintf("/api/1.1/configs/routing/%s.json", cdn)
	body, err := to.getBytesWithTTL(url, tmPollingInterval)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// TrafficRouterConfig gets the json arrays
func (to *Session) TrafficRouterConfig(cdn string) (*TrafficRouterConfig, error) {
	body, err := to.TrafficRouterConfigRaw(cdn)
	if err != nil {
		return nil, err
	}
	trConfig, err := trUnmarshall(body)
	if err != nil {
		return nil, err
	}
	return trConfig, nil
}

// in a seperate function for unit testing with files
func trUnmarshall(body []byte) (*TrafficRouterConfig, error) {
	var trConfig TrafficRouterConfig
	if err := json.Unmarshal(body, &trConfig); err != nil {
		return nil, err
	}
	return &trConfig, nil
}

func trTransformToMap(trConfig TrafficRouterConfig) TrafficRouterConfigMap {
	var trConfigMap TrafficRouterConfigMap
	trConfigMap.TrafficServer = make(map[string]TrafficServer)
	trConfigMap.TrafficMonitor = make(map[string]trafficMonitor)
	trConfigMap.TrafficRouter = make(map[string]trafficRouter)
	trConfigMap.CacheGroup = make(map[string]cacheGroup)
	trConfigMap.DeliveryService = make(map[string]deliveryService)
	trConfigMap.Config = make(map[string]interface{})
	trConfigMap.Stat = make(map[string]interface{})

	for _, trServer := range trConfig.TrafficServers {
		trConfigMap.TrafficServer[trServer.HostName] = trServer
	}
	for _, trMonitor := range trConfig.TrafficMonitors {
		trConfigMap.TrafficMonitor[trMonitor.HostName] = trMonitor
	}
	for _, trServer := range trConfig.TrafficServers {
		trConfigMap.TrafficServer[trServer.HostName] = trServer
	}
	for _, trRouter := range trConfig.TrafficRouters {
		trConfigMap.TrafficRouter[trRouter.Fqdn] = trRouter
	}
	for _, trCacheGroup := range trConfig.CacheGroups {
		trConfigMap.CacheGroup[trCacheGroup.Name] = trCacheGroup
	}
	for _, trDeliveryService := range trConfig.DeliveryServices {
		trConfigMap.DeliveryService[trDeliveryService.XMLID] = trDeliveryService
	}
	for trSettingKey, trSettingVal := range trConfig.Config {
		trConfigMap.Config[trSettingKey] = trSettingVal
	}
	for trStatKey, trStatVal := range trConfig.Stats {
		trConfigMap.Stat[trStatKey] = trStatVal
	}
	return trConfigMap
}
