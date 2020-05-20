package v1

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

import (
	"testing"
)

func TestATSConfigs(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		defer DeleteTestDeliveryServiceServersCreated(t)
		CreateTestDeliveryServiceServers(t)
		GetTestATSConfigs(t)
	})
}

func GetTestATSConfigs(t *testing.T) {
	if len(testData.Servers) < 1 {
		t.Fatal("cannot GET Server: no test data")
	}
	testServer := testData.Servers[0]

	serverList, _, err := TOSession.GetServerByHostName(testServer.HostName)
	if err != nil {
		t.Fatalf("cannot GET Server: %v", err)
	}
	if len(serverList) < 1 {
		t.Fatalf("cannot GET Server '" + testServer.HostName + "', returned no servers\n")
	}
	server := serverList[0]

	_, _, err = TOSession.GetATSServerConfigList(server.ID)
	if err != nil {
		t.Fatalf("Getting server '" + server.HostName + "' config list: " + err.Error() + "\n")
	}

	_, _, err = TOSession.GetATSServerConfigListByName(server.HostName)
	if err != nil {
		t.Fatalf("Getting server by name '" + server.HostName + "' config list: " + err.Error() + "\n")
	}

	serverConfigs := []string{
		"hosting.config",
		"packages",
		"chkconfig",
	}
	for _, serverConfig := range serverConfigs {
		if _, _, err = TOSession.GetATSServerConfig(server.ID, serverConfig); err != nil {
			t.Errorf("Getting server '" + server.HostName + "' config '" + serverConfig + "': " + err.Error() + "\n")
		}
		if _, _, err = TOSession.GetATSServerConfigByName(server.HostName, serverConfig); err != nil {
			t.Errorf("Getting server by name '" + server.HostName + "' config '" + serverConfig + "': " + err.Error() + "\n")
		}
	}

	profileConfigs := []string{
		"12M_facts",
		"50-ats.rules",
		"astats.config",
		"cache.config",
		"drop_qstring.config",
		"logging.config",
		"logging.yaml",
		"logs_xml.config",
		"plugin.config",
		"records.config",
		"storage.config",
		"sysctl.conf",
		"volume.config",
	}
	for _, profileConfig := range profileConfigs {
		if _, _, err = TOSession.GetATSProfileConfig(server.ProfileID, profileConfig); err != nil {
			t.Errorf("Getting profile '" + server.Profile + "' config '" + profileConfig + "': " + err.Error() + "\n")
		}
		if _, _, err = TOSession.GetATSProfileConfigByName(server.Profile, profileConfig); err != nil {
			t.Errorf("Getting profile by name '" + server.Profile + "' config '" + profileConfig + "': " + err.Error() + "\n")
		}
	}

	cdnConfigs := []string{
		"regex_revalidate.config",
		"bg_fetch.config",
		"set_dscp_42.config",
		"ssl_multicert.config",
	}
	for _, cdnConfig := range cdnConfigs {
		if _, _, err = TOSession.GetATSCDNConfig(server.CDNID, cdnConfig); err != nil {
			t.Errorf("Getting cdn '" + server.CDNName + "' config '" + cdnConfig + "': " + err.Error() + "\n")
		}
		if _, _, err = TOSession.GetATSCDNConfigByName(server.CDNName, cdnConfig); err != nil {
			t.Errorf("Getting cdn by name '" + server.CDNName + "' config '" + cdnConfig + "': " + err.Error() + "\n")
		}
	}
}



func CreateTestDeliveryServiceServers(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) < 1 {
		t.Error("GET DeliveryServices returned no dses, must have at least 1 to test ds-servers")
	}

	servers, _, err := TOSession.GetServers()
	if err != nil {
		t.Errorf("cannot GET Servers: %v", err)
	}
	if len(servers) < 1 {
		t.Error("GET Servers returned no servers, must have at least 1 to test ds-servers")
	}

	for _, ds := range dses {
		serverIDs := make([]int, 0, len(servers))
		for _, server := range servers {
			if server.Type == "EDGE" && server.CDNName == ds.CDNName {
				serverIDs = append(serverIDs, server.ID)
			}
		}

		if len(serverIDs) > 0 {
			_, err = TOSession.CreateDeliveryServiceServers(ds.ID, serverIDs, true)
			if err != nil {
				t.Errorf("POST delivery service servers: %v", err)
			}
		}
	}
}

// DeleteTestDeliveryServiceServersCreated deletes the dss assignments created by CreateTestDeliveryServiceServers.
func DeleteTestDeliveryServiceServersCreated(t *testing.T) {
	// You gotta do this because TOSession.GetDeliveryServiceServers doesn't fetch the complete response.......
	dssLen := len(testData.Servers) * len(testData.DeliveryServices)
	dsServers, _, err := TOSession.GetDeliveryServiceServersN(dssLen)
	if err != nil {
		t.Fatalf("GET delivery service servers: %v", err)
	}

	for _, dss := range dsServers.Response {
		if dss.DeliveryService == nil {
			t.Error("Found ds-to-server assignment with nil Delivery Service")
			continue
		}
		if dss.Server == nil {
			t.Error("Found ds-to-server assignment with nil Server")
			continue
		}

		_, _, err := TOSession.DeleteDeliveryServiceServer(*dss.DeliveryService, *dss.Server)
		if err != nil {
			t.Errorf("Failed to remove assignment of server #%d to DS #%d: %v", *dss.Server, *dss.DeliveryService, err)
		}
	}

	dsServers, _, err = TOSession.GetDeliveryServiceServersN(dssLen)
	if err != nil {
		t.Fatalf("GET delivery service servers: %v", err)
	}

	for _, dss := range dsServers.Response {
		if dss.DeliveryService == nil {
			t.Error("Found ds-to-server assignment (after supposed deletion) with nil DeliveryService")
			continue
		}
		if dss.Server == nil {
			t.Error("Found ds-to-server assignment (after supposed deletion) with nil Server")
			continue
		}

		t.Errorf("Found ds-to-server assignment {DSID: %d, Server: %d} after deletion", *dss.DeliveryService, *dss.Server)
	}
}
