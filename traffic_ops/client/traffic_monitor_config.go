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
	"strconv"
	"strings"
)

// TMConfigResponse ...
type TMConfigResponse struct {
	Version  string               `json:"version"`
	Response TrafficMonitorConfig `json:"response"`
}

// TrafficMonitorConfig ...
type TrafficMonitorConfig struct {
	TrafficServers   []TrafficServer        `json:"trafficServers,omitempty"`
	CacheGroups      []TMCacheGroup         `json:"cacheGroups,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
	TrafficMonitors  []TrafficMonitor       `json:"trafficMonitors,omitempty"`
	DeliveryServices []TMDeliveryService    `json:"deliveryServices,omitempty"`
	Profiles         []TMProfile            `json:"profiles,omitempty"`
}

// TrafficMonitorConfigMap ...
type TrafficMonitorConfigMap struct {
	TrafficServer   map[string]TrafficServer
	CacheGroup      map[string]TMCacheGroup
	Config          map[string]interface{}
	TrafficMonitor  map[string]TrafficMonitor
	DeliveryService map[string]TMDeliveryService
	Profile         map[string]TMProfile
}

// TrafficMonitor ...
type TrafficMonitor struct {
	Port     int    `json:"port"`
	IP6      string `json:"ip6"`
	IP       string `json:"ip"`
	HostName string `json:"hostName"`
	FQDN     string `json:"fqdn"`
	Profile  string `json:"profile"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

// TMDeliveryService ...
type TMDeliveryService struct {
	XMLID              string `json:"xmlId"`
	TotalTPSThreshold  int64  `json:"TotalTpsThreshold"`
	Status             string `json:"status"`
	TotalKbpsThreshold int64  `json:"TotalKbpsThreshold"`
}

// TMProfile ...
type TMProfile struct {
	Parameters TMParameters `json:"parameters"`
	Name       string       `json:"name"`
	Type       string       `json:"type"`
}

// TMParameters ...
type TMParameters struct {
	HealthConnectionTimeout                 int     `json:"health.connection.timeout"`
	HealthPollingURL                        string  `json:"health.polling.url"`
	HealthThresholdQueryTime                int     `json:"health.threshold.queryTime"`
	HistoryCount                            int     `json:"history.count"`
	HealthThresholdAvailableBandwidthInKbps string  `json:"health.threshold.availableBandwidthInKbps"`
	HealthThresholdLoadAvg                  float64 `json:"health.threshold.loadavg,string"`
	MinFreeKbps                             int64
}

// TrafficMonitorConfigMap ...
func (to *Session) TrafficMonitorConfigMap(cdn string) (*TrafficMonitorConfigMap, error) {
	tmConfig, err := to.TrafficMonitorConfig(cdn)
	if err != nil {
		return nil, err
	}
	tmConfigMap, err := trafficMonitorTransformToMap(tmConfig)
	if err != nil {
		return nil, err
	}
	return tmConfigMap, nil
}

// TrafficMonitorConfig ...
func (to *Session) TrafficMonitorConfig(cdn string) (*TrafficMonitorConfig, error) {
	url := fmt.Sprintf("/api/1.2/cdns/%s/configs/monitoring.json", cdn)
	resp, err := to.request(url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data TMConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data.Response, nil
}

func trafficMonitorTransformToMap(tmConfig *TrafficMonitorConfig) (*TrafficMonitorConfigMap, error) {
	var tm TrafficMonitorConfigMap

	tm.TrafficServer = make(map[string]TrafficServer)
	tm.CacheGroup = make(map[string]TMCacheGroup)
	tm.Config = make(map[string]interface{})
	tm.TrafficMonitor = make(map[string]TrafficMonitor)
	tm.DeliveryService = make(map[string]TMDeliveryService)
	tm.Profile = make(map[string]TMProfile)

	for _, trafficServer := range tmConfig.TrafficServers {
		tm.TrafficServer[trafficServer.HostName] = trafficServer
	}

	for _, cacheGroup := range tmConfig.CacheGroups {
		tm.CacheGroup[cacheGroup.Name] = cacheGroup
	}

	for parameterKey, parameterVal := range tmConfig.Config {
		tm.Config[parameterKey] = parameterVal
	}

	for _, trafficMonitor := range tmConfig.TrafficMonitors {
		tm.TrafficMonitor[trafficMonitor.HostName] = trafficMonitor
	}

	for _, deliveryService := range tmConfig.DeliveryServices {
		tm.DeliveryService[deliveryService.XMLID] = deliveryService
	}

	for _, profile := range tmConfig.Profiles {
		bwThresholdString := profile.Parameters.HealthThresholdAvailableBandwidthInKbps
		if strings.HasPrefix(bwThresholdString, ">") {
			var err error
			profile.Parameters.MinFreeKbps, err = strconv.ParseInt(bwThresholdString[1:len(bwThresholdString)], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		tm.Profile[profile.Name] = profile
	}

	return &tm, nil
}
