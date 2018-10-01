package v14

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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestServers(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestServers(t)
	UpdateTestServers(t)
	GetTestServers(t)
	DeleteTestServers(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

func CreateTestServers(t *testing.T) {

	// GET EDGE1 profile
	respProfiles, _, err := TOSession.GetProfileByName("EDGE1")
	if err != nil {
		t.Errorf("cannot GET Profiles - %v\n", err)
	}
	respProfile := respProfiles[0]

	// GET ONLINE status
	respStatuses, _, err := TOSession.GetStatusByName("ONLINE")
	if err != nil {
		t.Errorf("cannot GET Status by name: ONLINE - %v\n", err)
	}
	respStatus := respStatuses[0]

	// GET Denver physlocation
	respPhysLocations, _, err := TOSession.GetPhysLocationByName("Denver")
	if err != nil {
		t.Errorf("cannot GET PhysLocation by name: Denver - %v\n", err)
	}
	respPhysLocation := respPhysLocations[0]

	// GET cachegroup1 cachegroup
	respCacheGroups, _, err := TOSession.GetCacheGroupNullableByName("cachegroup1")
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: cachegroup1 - %v\n", err)
	}
	respCacheGroup := respCacheGroups[0]

	// loop through servers, assign FKs and create
	for _, server := range testData.Servers {
		// GET EDGE type
		respTypes, _, err := TOSession.GetTypeByName(server.Type)
		if err != nil {
			t.Errorf("cannot GET Division by name: EDGE - %v\n", err)
		}
		respType := respTypes[0]

		server.CDNID = respProfile.CDNID
		server.ProfileID = respProfile.ID
		server.TypeID = respType.ID
		server.StatusID = respStatus.ID
		server.PhysLocationID = respPhysLocation.ID
		server.CachegroupID = *respCacheGroup.ID

		resp, _, err := TOSession.CreateServer(server)
		log.Debugln("Response: ", server.HostName, " ", resp)
		if err != nil {
			t.Errorf("could not CREATE servers: %v\n", err)
		}
	}

}

func GetTestServers(t *testing.T) {

	for _, server := range testData.Servers {
		resp, _, err := TOSession.GetServerByHostName(server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server by name: %v - %v\n", err, resp)
		}
	}
}

func UpdateTestServers(t *testing.T) {

	firstServer := testData.Servers[0]
	hostName := firstServer.HostName
	// Retrieve the server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServerByHostName(hostName)

	if err != nil {
		t.Errorf("cannot GET Server by hostname: %v - %v\n", firstServer.HostName, err)
	}
	remoteServer := resp[0]
	updatedServerInterface := "bond1"
	updatedServerRack := "RR 119.03"

	// update rack and interfaceName values on server
	remoteServer.InterfaceName = updatedServerInterface
	remoteServer.Rack = updatedServerRack
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateServerByID(remoteServer.ID, remoteServer)
	if err != nil {
		t.Errorf("cannot UPDATE Server by hostname: %v - %v\n", err, alert)
	}

	// Retrieve the server to check rack and interfaceName values were updated
	resp, _, err = TOSession.GetServerByID(remoteServer.ID)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v\n", remoteServer.HostName, err)
	}

	respServer := resp[0]
	if respServer.InterfaceName != updatedServerInterface || respServer.Rack != updatedServerRack {
		t.Errorf("results do not match actual: %s, expected: %s\n", respServer.InterfaceName, updatedServerInterface)
		t.Errorf("results do not match actual: %s, expected: %s\n", respServer.Rack, updatedServerRack)
	}

}

func DeleteTestServers(t *testing.T) {

	for _, server := range testData.Servers {
		resp, _, err := TOSession.GetServerByHostName(server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server by hostname: %v - %v\n", server.HostName, err)
		}
		if len(resp) > 0 {
			respServer := resp[0]

			delResp, _, err := TOSession.DeleteServerByID(respServer.ID)
			if err != nil {
				t.Errorf("cannot DELETE Server by ID: %v - %v\n", err, delResp)
			}

			// Retrieve the Server to see if it got deleted
			serv, _, err := TOSession.GetServerByHostName(server.HostName)
			if err != nil {
				t.Errorf("error deleting Server hostname: %s\n", err.Error())
			}
			if len(serv) > 0 {
				t.Errorf("expected Server hostname: %s to be deleted\n", server.HostName)
			}
		}
	}
}
