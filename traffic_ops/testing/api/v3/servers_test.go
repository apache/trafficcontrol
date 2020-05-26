package v3

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
	"net/url"
	"testing"
)

func TestServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Topologies, DeliveryServices, Servers}, func() {
		UpdateTestServers(t)
		GetTestServersDetails(t)
		GetTestServers(t)
	})
}

func CreateTestServers(t *testing.T) {
	// loop through servers, assign FKs and create
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		resp, _, err := TOSession.CreateServer(server)
		t.Log("Response: ", server.HostName, " ", resp)
		if err != nil {
			t.Errorf("could not CREATE servers: %v", err)
		}
	}

}

func GetTestServers(t *testing.T) {
	serverCount := uint64(len(testData.Servers))

	params := url.Values{}
	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		params.Set("hostName", *server.HostName)
		_, alerts, count, _, err := TOSession.GetServers(&params)
		if err != nil {
			t.Errorf("cannot GET Server by name '%s': %v - %v", *server.HostName, err, alerts)
		} else if count != serverCount {
			t.Errorf("incorrect server count, expected: %d, actual: %d", serverCount, count)
		}
	}
}

func GetTestServersDetails(t *testing.T) {

	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}
		resp, _, err := TOSession.GetServerDetailsByHostName(*server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server Details by name: %v - %v", err, resp)
		}
	}
}

func UpdateTestServers(t *testing.T) {
	if len(testData.Servers) < 1 {
		t.Fatal("Need at least one server to test updating")
	}

	firstServer := testData.Servers[0]
	if firstServer.HostName == nil {
		t.Fatalf("First test server had nil hostname: %+v", firstServer)
	}

	hostName := *firstServer.HostName
	params := url.Values{}
	params.Add("hostName", hostName)

	// Retrieve the server by hostname so we can get the id for the Update
	resp, alerts, _, _, err := TOSession.GetServers(&params)
	if err != nil {
		t.Fatalf("cannot GET Server by hostname '%s': %v - %v", hostName, err, alerts)
	}
	if len(resp) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp))
		t.Logf("Testing will proceed with server: %+v", resp[0])
	}

	remoteServer := resp[0]
	if remoteServer.ID == nil {
		t.Fatalf("Got null ID for server '%s'", hostName)
	}

	infs := remoteServer.Interfaces
	if len(infs) < 1 {
		t.Fatalf("Expected server '%s' to have at least one network interface", hostName)
	}
	inf := infs[0]

	updatedServerInterface := "bond1"
	updatedServerRack := "RR 119.03"

	// update rack and interfaceName values on server
	inf.Name = updatedServerInterface
	infs[0] = inf
	remoteServer.Interfaces = infs
	remoteServer.Rack = &updatedServerRack

	alerts, _, err = TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", *remoteServer.ID, hostName, err, alerts)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, alerts, _, _, err = TOSession.GetServers(&params)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v", remoteServer.HostName, err)
	}
	if len(resp) < 1 {
		t.Fatalf("Expected at least one server to exist by hostname '%s'", hostName)
	}
	if len(resp) > 1 {
		t.Errorf("Expected exactly one server to exist by hostname '%s' - actual: %d", hostName, len(resp))
		t.Logf("Testing will proceed with server: %+v", resp[0])
	}

	respServer := resp[0]
	infs = respServer.Interfaces
	found := false
	for _, inf = range infs {
		if inf.Name == updatedServerInterface {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected server '%s' to have an interface named '%s' after update", hostName, updatedServerInterface)
		t.Logf("Actual interfaces: %+v", infs)
	}

	if respServer.Rack == nil {
		t.Errorf("results do not match actual: null, expected: '%s'", updatedServerRack)
	} else if *respServer.Rack != updatedServerRack {
		t.Errorf("results do not match actual: '%s', expected: '%s'", *respServer.Rack, updatedServerRack)
	}

	if remoteServer.TypeID == nil {
		t.Fatalf("Cannot test server type change update; server '%s' had nil type ID", hostName)
	}

	// Assign server to DS and then attempt to update to a different type
	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) < 1 {
		t.Fatal("GET DeliveryServices returned no dses, must have at least 1 to test invalid type server update")
	}

	serverTypes, _, err := TOSession.GetTypes("server")
	if err != nil {
		t.Fatalf("cannot GET Server Types: %v", err)
	}
	if len(serverTypes) < 2 {
		t.Fatal("GET Server Types returned less then 2 types, must have at least 2 to test invalid type server update")
	}
	for _, t := range serverTypes {
		if t.ID != *remoteServer.TypeID {
			remoteServer.TypeID = &t.ID
			break
		}
	}

	// Assign server to DS
	_, _, err = TOSession.CreateDeliveryServiceServers(*dses[0].ID, []int{*remoteServer.ID}, true)
	if err != nil {
		t.Fatalf("POST delivery service servers: %v", err)
	}

	// Attempt Update - should fail
	alerts, _, err = TOSession.UpdateServerByID(*remoteServer.ID, remoteServer)
	if err == nil {
		t.Errorf("expected error when updating Server Type of a server assigned to DSes")
	} else {
		t.Logf("type change update alerts: %+v", alerts)
	}

}

func DeleteTestServers(t *testing.T) {
	params := url.Values{}

	for _, server := range testData.Servers {
		if server.HostName == nil {
			t.Errorf("found server with nil hostname: %+v", server)
			continue
		}

		params.Set("hostName", *server.HostName)

		resp, alerts, _, _, err := TOSession.GetServers(&params)
		if err != nil {
			t.Errorf("cannot GET Server by hostname '%s': %v - %v", *server.HostName, err, alerts)
			continue
		}
		if len(resp) > 0 {
			if len(resp) > 1 {
				t.Errorf("Expected exactly one server by hostname '%s' - actual: %d", *server.HostName, len(resp))
				t.Logf("Testing will proceed with server: %+v", resp[0])
			}
			respServer := resp[0]

			if respServer.ID == nil {
				t.Errorf("Server '%s' had nil ID", *server.HostName)
				continue
			}

			delResp, _, err := TOSession.DeleteServerByID(*respServer.ID)
			if err != nil {
				t.Errorf("cannot DELETE Server by ID %d: %v - %v", *respServer.ID, err, delResp)
				continue
			}

			// Retrieve the Server to see if it got deleted
			serv, alerts, _, _, err := TOSession.GetServers(&params)
			if err != nil {
				t.Errorf("error deleting Server hostname '%s': %v - %v", *server.HostName, err, alerts)
			}
			if len(serv) > 0 {
				t.Errorf("expected Server hostname: %s to be deleted", *server.HostName)
			}
		}
	}
}
