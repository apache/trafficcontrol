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

func TestTRConfig(t *testing.T) {

	t.Log("Running TRConfig Tests")

	files, _ := ioutil.ReadDir("./testdata")
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "TRConfig.json") {
			t.Logf("Skipping test for %v, doesn't end in TRConfig, or is unreadable.", f.Name())
			continue
		}

		text, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s", f.Name()))
		if err != nil {
			t.Logf("Skipping test for %v: %v", f.Name(), err.Error())
			continue
		}
		fmt.Printf("Testing %v...\n", f.Name())

		trConfig, err := trUnmarshall(text)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(len(trConfig.TrafficServers), "TrafficServers found")
		for _, tServer := range trConfig.TrafficServers {
			t.Logf("  %v -> %v (%v remaps)", tServer.HostName, tServer.Ip, len(tServer.DeliveryServices))
		}

		t.Log(len(trConfig.CacheGroups), "CacheGroups found")
		for _, tLoc := range trConfig.CacheGroups {
			t.Logf("  %v -> (%v, %v)", tLoc.Name, tLoc.Coordinates.Longitude, tLoc.Coordinates.Latitude)
		}

		t.Log(len(trConfig.TrafficMonitors), "TrafficMonitors found")
		for _, tMon := range trConfig.TrafficMonitors {
			t.Logf("  %v -> %v", tMon.HostName, tMon.Ip)
		}

		t.Log(len(trConfig.DeliveryServices), "DeliveryServices found")
		for _, tDeliveryService := range trConfig.DeliveryServices {
			t.Logf("  %v -> %v MatchSets", tDeliveryService.XmlId, len(tDeliveryService.MatchSets))
		}

		t.Log(len(trConfig.Config), "Config settings  found")
		for cKey, cVal := range trConfig.Config {
			t.Logf("  %v -> %v", cKey, cVal)
		}

		trConfigMap := trTransformToMap(trConfig)

		t.Log(len(trConfigMap.TrafficServer), "TrafficServers found in Map")
		for tServerName, tServer := range trConfigMap.TrafficServer {
			t.Logf("   Map: %v -> %v (%v remaps)", tServerName, tServer.Ip, len(tServer.DeliveryServices))
		}

		t.Log(len(trConfigMap.CacheGroup), "CacheGroups found in Map")
		for _, tLoc := range trConfigMap.CacheGroup {
			t.Logf("  %v -> (%v, %v)", tLoc.Name, tLoc.Coordinates.Longitude, tLoc.Coordinates.Latitude)
		}

		t.Log(len(trConfigMap.TrafficMonitor), "TrafficMonitors found in Map")
		for _, tMon := range trConfigMap.TrafficMonitor {
			t.Logf("  %v -> %v", tMon.HostName, tMon.Ip)
		}

		t.Log(len(trConfigMap.DeliveryService), "DeliveryServices found in Map")
		for _, tDeliveryService := range trConfigMap.DeliveryService {
			t.Logf("  %v -> %v MatchSets", tDeliveryService.XmlId, len(tDeliveryService.MatchSets))
		}

		t.Log(len(trConfigMap.Config), "Config settings  found in Map")
		for sName, sValue := range trConfigMap.Config {
			t.Logf("  %v -> %v", sName, sValue)
		}

		t.Log(len(trConfigMap.Stat), "Stats found in Map")
		for sName, sValue := range trConfigMap.Stat {
			t.Logf("  %v -> %v", sName, sValue)
		}

	}

}
