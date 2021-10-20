package v2

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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, DeliveryServices, Servers}, func() {
		UpdateTestServers(t)
		GetTestServers(t)
		GetTestServersDetails(t)
	})
}

func CreateTestServers(t *testing.T) {

	var name string

	// loop through servers, assign FKs and create
	for _, server := range testData.Servers {
		resp, _, err := TOSession.CreateServer(server)
		t.Log("Response: ", server.HostName, " ", resp)
		if err != nil {
			t.Errorf("could not CREATE servers: %v", err)
		} else {
			name = server.HostName
		}
	}

	if name == "" {
		t.Fatal("Failed to create at least one valid server to use as a template")
	}

	servers, _, err := TOSession.GetServerByHostName(name)
	if err != nil {
		t.Fatalf("Failed to get server by hostname '%s': %v", name, err)
	}
	if len(servers) < 1 {
		t.Fatalf("Failed to get at least one server by hostname '%s'", name)
	}

	testServer := servers[0]

	testServer.DomainName = "test"
	testServer.HostName = "test"
	testServer.ID = 0
	testServer.IPIsService = true
	testServer.IPAddress = ""
	testServer.IP6IsService = false
	testServer.IP6Address = "::1/64"

	if alerts, _, err := TOSession.CreateServer(testServer); err == nil {
		t.Error("Didn't get expected error trying to create a server with a blank service address")
	} else {
		t.Logf("Received expected error creating a server with a blank service address: %v", err)
		t.Logf("Alerts: %v", alerts)
	}

}

func GetTestServers(t *testing.T) {

	for _, server := range testData.Servers {
		resp, _, err := TOSession.GetServerByHostName(server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server by name: %v - %v", err, resp)
		}
	}
}

func GetTestServersDetails(t *testing.T) {

	for _, server := range testData.Servers {
		resp, _, err := TOSession.GetServerDetailsByHostName(server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server Details by name: %v - %v", err, resp)
		}
	}
}

func UpdateTestServers(t *testing.T) {

	firstServer := testData.Servers[0]
	hostName := firstServer.HostName
	// Retrieve the server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServerByHostName(hostName)

	if err != nil {
		t.Errorf("cannot GET Server by hostname: %v - %v", firstServer.HostName, err)
	}
	remoteServer := resp[0]
	originalHostname := resp[0].HostName
	originalXMPIDD := resp[0].XMPPID
	updatedServerInterface := "bond1"
	updatedServerRack := "RR 119.03"
	updatedHostName := "atl-edge-01"
	updatedXMPPID := "change-it"

	// update rack, interfaceName and hostName values on server
	remoteServer.InterfaceName = updatedServerInterface
	remoteServer.Rack = updatedServerRack
	remoteServer.HostName = updatedHostName

	var alert tc.Alerts
	alert, _, err = TOSession.UpdateServerByID(remoteServer.ID, remoteServer)
	if err != nil {
		t.Errorf("cannot UPDATE Server by hostname: %v - %v", err, alert)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err = TOSession.GetServerByID(remoteServer.ID)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v", remoteServer.HostName, err)
	}

	respServer := resp[0]
	if respServer.InterfaceName != updatedServerInterface || respServer.Rack != updatedServerRack {
		t.Errorf("results do not match actual: %s, expected: %s", respServer.InterfaceName, updatedServerInterface)
		t.Errorf("results do not match actual: %s, expected: %s", respServer.Rack, updatedServerRack)
	}

	//Check change in hostname with no change to xmppid
	if originalHostname == respServer.HostName && originalXMPIDD == respServer.XMPPID {
		t.Errorf("HostName didn't change. Expected: #{updatedHostName}, actual: #{originalHostname}")
	}

	//Check to verify XMPPID never gets updated
	remoteServer.XMPPID = updatedXMPPID
	al, _, err := TOSession.UpdateServerByID(remoteServer.ID, remoteServer)
	t.Logf("Alerts from changing server XMPPID: %v", al)
	if err == nil {
		t.Error("Didn't get expected error changing server XMPPID")
	} else {
		t.Logf("Got expected error trying to change server XMPPID: %v", err)
	}

	//Change back hostname and xmppid to its original name for other tests to pass
	remoteServer.HostName = originalHostname
	remoteServer.XMPPID = originalXMPIDD
	alerts, _, err := TOSession.UpdateServerByID(remoteServer.ID, remoteServer)
	if err != nil {
		t.Fatalf("cannot UPDATE Server by ID %d (hostname '%s'): %v - %v", remoteServer.ID, hostName, err, alerts)
	}
	resp, _, err = TOSession.GetServerByID(remoteServer.ID)
	if err != nil {
		t.Errorf("cannot GET Server by hostName: %v - %v", originalHostname, err)
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
		if t.ID != remoteServer.TypeID {
			remoteServer.TypeID = t.ID
			break
		}
	}

	// Assign server to DS
	_, err = TOSession.CreateDeliveryServiceServers(*dses[0].ID, []int{remoteServer.ID}, true)
	if err != nil {
		t.Fatalf("POST delivery service servers: %v", err)
	}

	// Attempt Update - should fail
	_, _, err = TOSession.UpdateServerByID(remoteServer.ID, remoteServer)
	if err == nil {
		t.Errorf("expected error when updating Server Type of a server assigned to DSes")
	}

}

func DeleteTestServers(t *testing.T) {

	for _, server := range testData.Servers {
		resp, _, err := TOSession.GetServerByHostName(server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server by hostname: %v - %v", server.HostName, err)
		}
		if len(resp) > 0 {
			respServer := resp[0]

			dses, _, err := TOSession.GetServerIDDeliveryServices(respServer.ID)
			for _, ds := range dses {
				if ds.ID == nil {
					t.Logf("got DS assigned to server #%d that has no ID - deletion may fail!", respServer.ID)
					continue
				}
				if ds.Active == nil || *ds.Active {
					setInactive(t, *ds.ID)
				}
			}

			delResp, _, err := TOSession.DeleteServerByID(respServer.ID)
			if err != nil {
				t.Errorf("cannot DELETE Server by ID: %v - %v", err, delResp)
			}

			// Retrieve the Server to see if it got deleted
			serv, _, err := TOSession.GetServerByHostName(server.HostName)
			if err != nil {
				t.Errorf("error deleting Server hostname: %s", err.Error())
			}
			if len(serv) > 0 {
				t.Errorf("expected Server hostname: %s to be deleted", server.HostName)
			}
		}
	}
}
