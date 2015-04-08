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

import "encoding/json"

// TrafficRouterConfig is the json unmarshalled without any changes
// note all structs are local to this file _except_ the TrafficRouterConfig struct.
type TrafficRouterConfig struct {
	TrafficServers   []trafficServer        `json:"trafficServers"`
	TrafficMonitors  []trafficMonitor       `json:"trafficMonitors"`
	TrafficRouters   []trafficRouter        `json:"trafficRouters"`
	CacheGroups      []cacheGroup           `json:"cacheGroups"`
	DeliveryServices []deliveryService      `json:"deliveryServices"`
	Stats            map[string]interface{} `json:"stats"`
	Config           map[string]interface{} `json:"config"`
}

type trafficServer struct {
	Profile          string              `json:"profile"`
	Ip               string              `json:"ip"`
	Status           string              `json:"status"`
	CacheGroup       string              `json:"cacheGroup"`
	Ip6              string              `json:"ip6"`
	Port             int                 `json:"port"`
	HostName         string              `json:"hostName"`
	Fqdn             string              `json:"fqdn"`
	InterfaceName    string              `json:"interfaceName"`
	Type             string              `json:"type"`
	HashId           string              `json:"hashId"`
	DeliveryServices []tsdeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
}

type tsdeliveryService struct {
	Xmlid  string   `json:"xmlId"`
	Remaps []string `json:"remaps"`
}

type trafficMonitor struct {
	Port     int    `json:"port"`
	Ip6      string `json:"ip6"`
	Ip       string `json:"ip"`
	HostName string `json:"hostName"`
	Fqdn     string `json:"fqdn"`
	Profile  string `json:"profile"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

type trafficRouter struct {
	Port     int    `json:"port"`
	Ip6      string `json:"ip6"`
	Ip       string `json:"ip"`
	Fqdn     string `json:"fqdn"`
	Profile  string `json:"profile"`
	Location string `json:"location"`
	Status   string `json:"status"`
	ApiPort  int    `json:"apiPort"`
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
	XmlId             string            `json:"xmlId"`
	missLocation      missLocation      `json:"missLocation"`
	CoverageZoneOnly  bool              `json:"coverageZoneOnly"`
	MatchSets         []matchSet        `json:"matchSets"`
	Ttl               int               `json:"ttl"`
	Ttls              ttls              `json:"ttls"`
	BypassDestination bypassDestination `json:"bypassDestination"`
	StatcDnsEntries   []staticDns       `json:"statitDnsEntries"`
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

type staticDns struct {
	Value string `json:"value"`
	Ttl   int    `json:"ttl"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type soa struct {
	Admin   string `json:"admin"`
	retry   int    `json:"retry"`
	minimum int    `json:"minimum"`
	refresh int    `json:"refresh"`
	expire  int    `json:"expire"`
}

type TrafficRouterConfigMap struct {
	TrafficServer   map[string]trafficServer
	TrafficMonitor  map[string]trafficMonitor
	TrafficRouter   map[string]trafficRouter
	CacheGroup      map[string]cacheGroup
	DeliveryService map[string]deliveryService
	Config          map[string]interface{}
	Stat            map[string]interface{}
}

// get a  bunch of maps
func (to *Session) TrafficRouterConfigMap(cdn string) (TrafficRouterConfigMap, error) {
	trConfig, err := to.TrafficRouterConfig(cdn)
	trConfigMap := trTransformToMap(trConfig)
	return trConfigMap, err
}

func (to *Session) TrafficRouterConfigRaw(cdn string) ([]byte, error) {
	body, err := to.getBytesWithTTL("/api/1.1/configs/routing/"+cdn+".json", tmPollingInterval)
	return body, err
}

// get the json arrays
func (to *Session) TrafficRouterConfig(cdn string) (TrafficRouterConfig, error) {
	body, err := to.TrafficRouterConfigRaw(cdn)
	trConfig, err := trUnmarshall(body)
	return trConfig, err
}

// in a seperate function for unit testing with files
func trUnmarshall(body []byte) (TrafficRouterConfig, error) {
	var trConfig TrafficRouterConfig
	err := json.Unmarshal(body, &trConfig)
	return trConfig, err
}

func trTransformToMap(trConfig TrafficRouterConfig) TrafficRouterConfigMap {
	var trConfigMap TrafficRouterConfigMap
	trConfigMap.TrafficServer = make(map[string]trafficServer)
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
		trConfigMap.DeliveryService[trDeliveryService.XmlId] = trDeliveryService
	}
	for trSettingKey, trSettingVal := range trConfig.Config {
		trConfigMap.Config[trSettingKey] = trSettingVal
	}
	for trStatKey, trStatVal := range trConfig.Stats {
		trConfigMap.Stat[trStatKey] = trStatVal
	}
	return trConfigMap
}
