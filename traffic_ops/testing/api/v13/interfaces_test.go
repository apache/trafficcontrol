package v13

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

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestInterfaces(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestServers(t)
	CreateTestInterfaces(t)
	UpdateTestInterfaces(t)
	GetTestInterfaces(t)
	DeleteTestInterfaces(t)
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

func GetServerIdForInterfacesTest(t *testing.T) int {
	// GET server "atlanta-edge-01"
	resp, _, err := TOSession.GetServerByHostName("atlanta-edge-01")
	if err != nil {
		t.Errorf("cannot GET server by name: atlanta-edge-01 - %v\n", err)
	}
	server := resp[0]

	return server.ID
}

func CreateTestInterfaces(t *testing.T) {

	serverIdForInterfacesTest := GetServerIdForInterfacesTest(t)

	for _, intf := range testData.Interfaces {
		intf.ServerID = serverIdForInterfacesTest
		resp, _, err := TOSession.CreateInterface(serverIdForInterfacesTest, intf)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE Interfaces: %v\n", err)
		}
	}

}

func UpdateTestInterfaces(t *testing.T) {

	serverIdForInterfacesTest := GetServerIdForInterfacesTest(t)

	// Retrieve an Interface
	resp, _, err := TOSession.GetInterfacesByServer(serverIdForInterfacesTest)
	if err != nil {
		t.Errorf("cannot GET Interfaces by server: %v - %v\n", err, resp)
	}
	intf := resp[0]
	expectedInterfaceName := "ens32"
	intf.InterfaceName = expectedInterfaceName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateInterfaces(intf.ID, intf)
	if err != nil {
		t.Errorf("cannot UPDATE Interface by id: %v - %v\n", err, alert)
	}

	// Retrieve the Interface to check InterfaceName got updated
	resp, _, err = TOSession.GetInterfacesById(intf.ID)
	if err != nil {
		t.Errorf("cannot GET Interfaces by server: %v - %v\n", err, resp)
	}
	intf = resp[0]
	if intf.InterfaceName != expectedInterfaceName {
		t.Errorf("results do not match actual: %s, expected: %s\n", intf.InterfaceName, expectedInterfaceName)
	}

}

func GetTestInterfaces(t *testing.T) {

	serverIdForInterfacesTest := GetServerIdForInterfacesTest(t)

	resp, _, err := TOSession.GetInterfacesByServer(serverIdForInterfacesTest)
	if err != nil {
		t.Errorf("cannot GET Interfaces by server: %v - %v\n", err, resp)
	}
}

func DeleteTestInterfaces(t *testing.T) {

	serverIdForInterfacesTest := GetServerIdForInterfacesTest(t)

	// Retrieve the Server
	Servers, _, err := TOSession.GetServerByID(serverIdForInterfacesTest)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v\n", serverIdForInterfacesTest, err)
	}
	server := Servers[0]

	// Retrieve the Interface
	resp, _, err := TOSession.GetInterfacesByServer(serverIdForInterfacesTest)
	if err != nil {
		t.Errorf("cannot GET Interfaces by server: %v - %v\n", err, resp)
	}
	if len(resp) > 0 {
		index := -1
		for i, _ := range resp {
			if resp[i].InterfaceName != server.InterfaceName {
				index = i
				break
			}
		}
		if index == -1 {
			t.Errorf("cannot GET new Interface by server: %v - %v\n", serverIdForIpsTest, "no new Interface")
		}
		intf := resp[index]

		_, _, err := TOSession.DeleteInterfaces(intf.ID)
		if err != nil {
			t.Errorf("cannot DELETE Interface by id: '%s' %v\n", intf.ID, err)
		}

		// Retrieve the Interface to see if it got deleted
		resp, _, err := TOSession.GetInterfacesById(intf.ID)
		if err != nil {
			t.Errorf("cannot GET Interfaces by server: %v - %v\n", err, resp)
		}
		if len(resp) > 0 {
			t.Errorf("expected Interface by id %s is not deleted\n", intf.ID)
		}
	}
}
