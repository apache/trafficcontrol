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

package integration

import (
	"encoding/json"
	"fmt"
	"testing"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

func TestTrafficMonitorConfig(t *testing.T) {
	cdn, err := GetCdn()
	if err != nil {
		t.Errorf("Could not get CDN, error was: %v\n", err)
	}
	uri := fmt.Sprintf("/api/1.2/cdns/%s/configs/monitoring.json", cdn.Name)
	resp, err := Request(*to, "GET", uri, nil)
	if err != nil {
		t.Errorf("Could not get %s reponse was: %v\n", uri, err)
		t.FailNow()
	}

	defer resp.Body.Close()
	var apiTMConfigRes traffic_ops.TMConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiTMConfigRes); err != nil {
		t.Errorf("Could not decode Traffic Monitor Config response.  Error is: %v\n", err)
		t.FailNow()
	}
	apiTMConfig := apiTMConfigRes.Response

	clientTMConfig, err := to.TrafficMonitorConfig(cdn.Name)
	if err != nil {
		t.Errorf("Could not get Traffic Monitor Config from client.  Error is: %v\n", err)
		t.FailNow()
	}

	if len(apiTMConfig.CacheGroups) != len(clientTMConfig.CacheGroups) {
		t.Errorf("Length of Traffic Monitor config cachegroups do not match! Expected %v, got %v\n", len(apiTMConfig.CacheGroups), len(clientTMConfig.CacheGroups))
	}

	for _, apiCg := range apiTMConfig.CacheGroups {
		match := false
		for _, clientCg := range clientTMConfig.CacheGroups {
			if apiCg == clientCg {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a cachegroup matching %+v\n", apiCg)
		}
	}

	if len(apiTMConfig.DeliveryServices) != len(clientTMConfig.DeliveryServices) {
		t.Errorf("Length of Traffic Monitor config deliveryserivces do not match! Expected %v, got %v\n", len(apiTMConfig.DeliveryServices), len(clientTMConfig.DeliveryServices))
	}

	for _, apiDs := range apiTMConfig.DeliveryServices {
		match := false
		for _, clientDs := range clientTMConfig.DeliveryServices {
			if apiDs == clientDs {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a Deliveryservice matching %+v\n", apiDs)
		}
	}

	if len(apiTMConfig.Profiles) != len(clientTMConfig.Profiles) {
		t.Errorf("Length of Traffic Monitor config profiles do not match! Expected %v, got %v\n", len(apiTMConfig.Profiles), len(clientTMConfig.Profiles))
	}

	for _, apiProfile := range apiTMConfig.Profiles {
		match := false
		for _, clientProfile := range clientTMConfig.Profiles {
			if apiProfile == clientProfile {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a Profile matching %+v\n", apiProfile)
		}
	}

	if len(apiTMConfig.TrafficMonitors) != len(clientTMConfig.TrafficMonitors) {
		t.Errorf("Length of Traffic Monitor config traffic monitors does not match! Expected %v, got %v\n", len(apiTMConfig.TrafficMonitors), len(clientTMConfig.TrafficMonitors))
	}

	for _, apiTM := range apiTMConfig.TrafficMonitors {
		match := false
		for _, clientTM := range clientTMConfig.TrafficMonitors {
			if apiTM == clientTM {
				match = true
			}
		}
		if !match {
			t.Errorf("Did not get a Traffic Monitor matching %+v\n", apiTM)
		}
	}

	if len(apiTMConfig.TrafficServers) != len(clientTMConfig.TrafficServers) {
		t.Errorf("Length of Traffic Monitor config traffic servers does not match! Expected %v, got %v\n", len(apiTMConfig.TrafficServers), len(clientTMConfig.TrafficServers))
	}

	for _, apiTS := range apiTMConfig.TrafficServers {
		match := false
		for _, clientTS := range clientTMConfig.TrafficServers {
			if apiTS.HostName == clientTS.HostName {
				match = true
				if apiTS.CacheGroup != clientTS.CacheGroup {
					t.Errorf("Cachegroup -- Expected %v, got %v\n", apiTS.CacheGroup, clientTS.CacheGroup)
				}
				if len(apiTS.DeliveryServices) != len(clientTS.DeliveryServices) {
					t.Errorf("len DeliveryServices -- Expected %v, got %v\n", len(apiTS.DeliveryServices), len(clientTS.DeliveryServices))
				}
				for _, apiDS := range apiTS.DeliveryServices {
					dsMatch := false
					for _, clientDS := range clientTS.DeliveryServices {
						if apiDS.Xmlid == clientDS.Xmlid && len(apiDS.Remaps) == len(clientDS.Remaps) {
							dsMatch = true
						}
					}
					if !dsMatch {
						t.Errorf("Could not finding a matching DS for %v\n", apiDS.Xmlid)
					}
				}
				if apiTS.FQDN != clientTS.FQDN {
					t.Errorf("FQDN -- Expected %v, got %v\n", apiTS.FQDN, clientTS.FQDN)
				}
				if apiTS.HashID != clientTS.HashID {
					t.Errorf("HashID -- Expected %v, got %v\n", apiTS.HashID, clientTS.HashID)
				}
				if apiTS.IP != clientTS.IP {
					t.Errorf("IP -- Expected %v, got %v\n", apiTS.IP, clientTS.IP)
				}
				if apiTS.IP6 != clientTS.IP6 {
					t.Errorf("IP6 -- Expected %v, got %v\n", apiTS.IP6, clientTS.IP6)
				}
				if apiTS.InterfaceName != clientTS.InterfaceName {
					t.Errorf("Interface Name -- Expected %v, got %v\n", apiTS.InterfaceName, clientTS.InterfaceName)
				}
				if apiTS.Port != clientTS.Port {
					t.Errorf("Port -- Expected %v, got %v\n", apiTS.Port, clientTS.Port)
				}
				if apiTS.Profile != clientTS.Profile {
					t.Errorf("Profile -- Expected %v, got %v\n", apiTS.Profile, clientTS.Profile)
				}
				if apiTS.Status != clientTS.Status {
					t.Errorf("Status -- Expected %v, got %v\n", apiTS.Status, clientTS.Status)
				}
				if apiTS.Type != clientTS.Type {
					t.Errorf("Type -- Expected %v, got %v\n", apiTS.Type, clientTS.Type)
				}
			}
		}
		if !match {
			t.Errorf("Did not get a Traffic Server matching %+v\n", apiTS)
		}
	}
}
