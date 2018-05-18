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

var serverIdForIpsTest int

func TestIPs(t *testing.T) {

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
	CreateTestIPs(t)
	UpdateTestIPs(t)
	GetTestIPs(t)
	DeleteTestIPs(t)
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

func GetServerIdForIpTest(t *testing.T) int {
	// GET server "atlanta-edge-01"
	resp, _, err := TOSession.GetServerByHostName("atlanta-edge-01")
	if err != nil {
		t.Errorf("cannot GET server by name: atlanta-edge-01 - %v\n", err)
	}
	server := resp[0]

	return server.ID

}

func CreateTestIPs(t *testing.T) {

	serverIdForIpsTest := GetServerIdForIpTest(t)

	// GET IP_SECONDARY type
	respTypes, _, err := TOSession.GetTypeByName("IP_SECONDARY")
	if err != nil {
		t.Errorf("cannot GET Type by name: IP_SECONDARY - %v\n", err)
	}
	respType := respTypes[0]

	// GET interface on the server
	respIntfs, _, err := TOSession.GetInterfacesByServer(serverIdForIpsTest)
	if err != nil {
		t.Errorf("cannot GET Interfaces by server: %v - %v\n", serverIdForIpsTest, err)
	}
	respIntf := respIntfs[0]

	for _, ip := range testData.IPs {
		ip.TypeID = respType.ID
		ip.ServerID = serverIdForIpsTest
		ip.InterfaceID = respIntf.ID
		resp, _, err := TOSession.CreateIP(serverIdForIpsTest, ip)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE secondary IPs: %v\n", err)
		}
	}

}

func UpdateTestIPs(t *testing.T) {

	serverIdForIpsTest := GetServerIdForIpTest(t)

	// Retrieve an IP
	resp, _, err := TOSession.GetIPsByServer(serverIdForIpsTest)
	if err != nil {
		t.Errorf("cannot GET IPs by server: %v - %v\n", err, resp)
	}
	ip := resp[0]
	expectedIPGateway := "192.168.1.1"
	ip.IPGateway = expectedIPGateway
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateIPs(ip.ID, ip)
	if err != nil {
		t.Errorf("cannot UPDATE IP by id: %v - %v\n", err, alert)
	}

	// Retrieve the IP to check IP gateway got updated
	resp, _, err = TOSession.GetIPsById(ip.ID)
	if err != nil {
		t.Errorf("cannot GET IPs by server: %v - %v\n", err, resp)
	}
	ip = resp[0]
	if ip.IPGateway != expectedIPGateway {
		t.Errorf("results do not match actual: %s, expected: %s\n", ip.IPGateway, expectedIPGateway)
	}

}

func GetTestIPs(t *testing.T) {

	serverIdForIpsTest := GetServerIdForIpTest(t)

	resp, _, err := TOSession.GetIPsByServer(serverIdForIpsTest)
	if err != nil {
		t.Errorf("cannot GET IPs by server: %v - %v\n", err, resp)
	}

}

func DeleteTestIPs(t *testing.T) {

	serverIdForIpsTest := GetServerIdForIpTest(t)

	// Retrieve a Secondary IP
	resp, _, err := TOSession.GetIPsByServer(serverIdForIpsTest)
	if err != nil {
		t.Errorf("cannot GET IPs by server: %v - %v\n", serverIdForIpsTest, err)
	}
	if len(resp) > 0 {
		index := -1
		for i, _ := range resp {
			if resp[i].Type == "IP_SECONDARY" {
				index = i
				break
			}
		}
		if index == -1 {
			t.Errorf("cannot GET Secondary IP by server: %v - %v\n", serverIdForIpsTest, "no secondary IP")
		}
		ip := resp[index]

		_, _, err := TOSession.DeleteIPs(ip.ID)
		if err != nil {
			t.Errorf("cannot DELETE Secondary IP by id: '%s' %v\n", ip.ID, err)
		}

		// Retrieve the Secondary IP to see if it got deleted
		resp, _, err := TOSession.GetIPsById(ip.ID)
		if err != nil {
			t.Errorf("cannot GET Secondary IPs by server: %v - %v\n", err, resp)
		}
		if len(resp) > 0 {
			t.Errorf("expected Secondary IP by id %s is not deleted\n", ip.ID)
		}
	}
}
