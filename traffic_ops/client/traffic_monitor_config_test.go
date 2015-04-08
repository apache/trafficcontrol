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
	"fmt"
	// "os"
	"io/ioutil"
	"strings"
	"testing"
)

func TestTrafficMonitorConfig(t *testing.T) {

	fmt.Println("Running Traffic Monitor Config Tests")

	files, _ := ioutil.ReadDir("./testdata")
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "TrafficMonitorConfig.json") {
			t.Logf("Skipping test for %v, doesn't end in TrafficMonitorConfig.json, or is unreadable.", f.Name())
			continue
		}
		text, err := ioutil.ReadFile("./testdata/" + f.Name())
		if err != nil {
			t.Logf("Skipping test for %v: %v", f.Name(), err.Error())
			continue
		}
		fmt.Printf("Testing %v...\n", f.Name())

		trafficMonitorConfig, err := trafficMonitorConfigUnmarshall(text)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(len(trafficMonitorConfig.TrafficServers), "TrafficServers found")
		for _, trafficServer := range trafficMonitorConfig.TrafficServers {
			t.Logf(" Traffic Server | FQDN: %65v | IP: %15v |", trafficServer.Fqdn, trafficServer.Ip)
		}

		t.Log(len(trafficMonitorConfig.CacheGroups), "CacheGroups found")
		for _, cacheGroup := range trafficMonitorConfig.CacheGroups {
			t.Logf(" Cache Group | Name: %20v | Coordinates (%10v, %10v) |", cacheGroup.Name, cacheGroup.Coordinates.Longitude, cacheGroup.Coordinates.Latitude)
		}

		t.Log(len(trafficMonitorConfig.TrafficMonitors), "TrafficMonitors found")
		for _, trafficMonitor := range trafficMonitorConfig.TrafficMonitors {
			t.Logf(" Traffic Monitor | HostName: %30v | IP: %15v |", trafficMonitor.HostName, trafficMonitor.Ip)
		}

		t.Log(len(trafficMonitorConfig.DeliveryServices), "DeliveryServices found")
		for _, deliveryService := range trafficMonitorConfig.DeliveryServices {
			t.Logf(" Delivery Service | xmlId: %20v | totalTpsThreshold: %10v | totalKbpsThreshold: %10v | status: %10v |", deliveryService.XmlId, deliveryService.TotalTpsThreshold, deliveryService.TotalKbpsThreshold, deliveryService.Status)
		}

		t.Log(len(trafficMonitorConfig.Config), "Config settings found")
		for parameterKey, parameterVal := range trafficMonitorConfig.Config {
			t.Logf(" Parameter | Name: %30v | Value: %80v |", parameterKey, parameterVal)
		}

		t.Log(len(trafficMonitorConfig.Profiles), "Profiles found")
		for _, profile := range trafficMonitorConfig.Profiles {
			t.Logf("Profile   | Name: %10v | Type: %30v | ", profile.Name, profile.Type)
			pp := profile.Parameters
			t.Logf("Profile Parameters | Timeout: %10v |  Polling URL: %10v, Avail BW: %10v | Load Avg: %10v, Query Time: %10v | History Count: %10v | Min Free KBps: %10v", pp.HealthConnectionTimeout, pp.HealthPollingUrl, pp.HealthThresholdAvailableBandwidthInKbps, pp.HealthThresholdLoadAvg, pp.HealthThresholdQueryTime, pp.HistoryCount, pp.MinFreeKbps)
		}

		trafficMonitorConfigMap := trafficMonitorTransformToMap(trafficMonitorConfig)

		t.Log(len(trafficMonitorConfigMap.TrafficServer), "TrafficServers foundn in map")
		for _, trafficServer := range trafficMonitorConfigMap.TrafficServer {
			t.Logf(" Traffic Server | FQDN: %65v | IP: %15v |", trafficServer.Fqdn, trafficServer.Ip)
		}

		t.Log(len(trafficMonitorConfigMap.CacheGroup), "CacheGroups found in map")
		for _, cacheGroup := range trafficMonitorConfigMap.CacheGroup {
			t.Logf(" Cache Group | Name: %20v | Coordinates (%10v, %10v) |", cacheGroup.Name, cacheGroup.Coordinates.Longitude, cacheGroup.Coordinates.Latitude)
		}

		t.Log(len(trafficMonitorConfigMap.TrafficMonitor), "TrafficMonitors found in map")
		for _, trafficMonitor := range trafficMonitorConfigMap.TrafficMonitor {
			t.Logf(" Traffic Monitor | HostName: %30v | IP: %15v |", trafficMonitor.HostName, trafficMonitor.Ip)
		}

		t.Log(len(trafficMonitorConfigMap.DeliveryService), "DeliveryServices found in map")
		for _, deliveryService := range trafficMonitorConfigMap.DeliveryService {
			t.Logf(" Delivery Service | xmlId: %20v | totalTpsThreshold: %10v | totalKbpsThreshold: %10v | status: %10v |", deliveryService.XmlId, deliveryService.TotalTpsThreshold, deliveryService.TotalKbpsThreshold, deliveryService.Status)
		}

		t.Log(len(trafficMonitorConfigMap.Config), "Config settings found in map")
		for parameterKey, parameterVal := range trafficMonitorConfigMap.Config {
			t.Logf(" Parameter | Name: %30v | Value: %80v |", parameterKey, parameterVal)
		}

		t.Log(len(trafficMonitorConfigMap.Profile), "Profiles found in map")
		for _, profile := range trafficMonitorConfigMap.Profile {
			t.Logf(" Profile   | Name: %10v | Type: %30v |", profile.Name, profile.Type)
			pp := profile.Parameters
			t.Logf("Profile Parameters | Timeout: %10v |  Polling URL: %10v, Avail BW: %10v | Load Avg: %10v, Query Time: %10v | History Count: %10v | Min Free KBps: %10v", pp.HealthConnectionTimeout, pp.HealthPollingUrl, pp.HealthThresholdAvailableBandwidthInKbps, pp.HealthThresholdLoadAvg, pp.HealthThresholdQueryTime, pp.HistoryCount, pp.MinFreeKbps)
		}

	}
}
