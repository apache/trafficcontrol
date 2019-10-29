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
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestServerServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices, CacheGroupsDeliveryServices, ServerCapabilities, ServerServerCapabilities}, func() {
		GetTestServerServerCapabilities(t)
	})
}

func CreateTestServerServerCapabilities(t *testing.T) {

	// Valid POSTs

	// loop through server ServerCapabilities, assign FKs and create
	for _, ssc := range testData.ServerServerCapabilities {
		servResp, _, err := TOSession.GetServerByHostName(*ssc.Server)
		if err != nil {
			t.Fatalf("cannot GET Server by hostname: %v - %v\n", *ssc.Server, err)
		}
		if len(servResp) != 1 {
			t.Fatalf("cannot GET Server by hostname: %v. Response did not include record.\n", *ssc.Server)
		}
		server := servResp[0]
		ssc.ServerID = &server.ID
		resp, _, err := TOSession.CreateServerServerCapability(ssc)
		if err != nil {
			t.Errorf("could not POST the server capability %v to server %v: %v\n", *ssc.ServerCapability, *ssc.Server, err)
		}
		log.Debugln("Response: ", server.HostName, " ", resp)
	}

	// Invalid POSTs

	ssc := testData.ServerServerCapabilities[0]

	// Attempt to assign already assigned server capability
	_, _, err := TOSession.CreateServerServerCapability(ssc)
	if err == nil {
		t.Error("expected to receive error when assigning a already assigned server capability\n")
	}

	// Attempt to assign an server capability with no ID
	sscNilID := tc.ServerServerCapability{
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability without a server ID\n")
	}

	// Attempt to assign an server capability with no server capability
	sscNilCapability := tc.ServerServerCapability{
		ServerID: ssc.ServerID,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server without a server capability\n")
	}

	// Attempt to assign an server capability with invalid server capability
	sscInvalidCapability := tc.ServerServerCapability{
		ServerID:         ssc.ServerID,
		ServerCapability: util.StrPtr("bogus"),
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a non existent server capability to a server\n")
	}

	// Attempt to assign an server capability with invalid server capability
	sscInvalidID := tc.ServerServerCapability{
		ServerID:         util.IntPtr(-1),
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a non existent server ID\n")
	}
}

func GetTestServerServerCapabilities(t *testing.T) {
	// Get All Server Capabilities
	sscs, _, err := TOSession.GetServerServerCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v\n", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil\n")
	}
	if len(sscs) != len(testData.ServerServerCapabilities) {
		t.Errorf("expect %v server capabilities assigned to servers received %v \n", len(testData.ServerServerCapabilities), len(sscs))
	}

	checkResp := func(t *testing.T, sscs []tc.ServerServerCapability) {
		if sscs == nil {
			t.Fatal("returned server capabilities assigned to servers was nil\n")
		}
		if len(sscs) != 1 {
			t.Errorf("expect 1 server capabilities assigned to server received %v \n", len(sscs))
		}
	}

	for _, ssc := range sscs {
		// Get assigned Server Capabilities by server id
		sscs, _, err := TOSession.GetServerServerCapabilities(ssc.ServerID, nil, nil)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server ID %v: %v\n", *ssc.ServerID, err)
		}
		checkResp(t, sscs)
		// Get assigned Server Capabilities by host name
		sscs, _, err = TOSession.GetServerServerCapabilities(nil, ssc.Server, nil)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server host name %v: %v\n", *ssc.Server, err)
		}
		checkResp(t, sscs)

		// Get assigned Server Capabilities by server capability
		sscs, _, err = TOSession.GetServerServerCapabilities(nil, nil, ssc.ServerCapability)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server capability %v: %v\n", *ssc.ServerCapability, err)
		}
		checkResp(t, sscs)
	}
}

func DeleteTestServerServerCapabilities(t *testing.T) {
	// Get Server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v\n", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil\n")
	}

	// Make the server capabilities required on the servers' delivery services
	// Used to make sure we block server server_capability DELETE in that case
	dsReqCaps := []tc.DeliveryServicesRequiredCapability{}

	for _, ssc := range sscs {
		dsServers, _, err := TOSession.GetDeliveryServiceServersWithLimits(1, []int{}, []int{*ssc.ServerID})
		if err != nil {
			t.Fatalf("cannot GET server delivery services assigned to servers: %v\n", err)
		}
		if len(dsServers.Response) == 0 {
			t.Fatal("servers must be assigned to delivery service")
		}
		dsReqCap := tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  dsServers.Response[0].DeliveryService,
			RequiredCapability: ssc.ServerCapability,
		}
		_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(dsReqCap)
		if err != nil {
			t.Fatalf("could not POST the server capability %v to ds %v: %v\n", *dsReqCap.RequiredCapability, *dsReqCap.DeliveryServiceID, err)
		}
		dsReqCaps = append(dsReqCaps, dsReqCap)
	}

	// Delete should fail as their delivery services now require the capabilities
	for _, ssc := range sscs {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err == nil {
			t.Fatalf("should have gotten error when using DELETE on the server capability %v from server %v as it is required by associated dses\n", *ssc.ServerCapability, *ssc.Server)
		}
	}

	for _, dsReqCap := range dsReqCaps {
		_, _, err := TOSession.DeleteDeliveryServicesRequiredCapability(*dsReqCap.DeliveryServiceID, *dsReqCap.RequiredCapability)
		if err != nil {
			t.Fatalf("could not DELETE the server capability %v from ds %v: %v\n", *dsReqCap.RequiredCapability, *dsReqCap.DeliveryServiceID, err)
		}
	}

	// Remove the requirement so we can actually delete them

	for _, ssc := range sscs {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err != nil {
			t.Errorf("could not DELETE the server capability %v from server %v: %v\n", *ssc.ServerCapability, *ssc.Server, err)
		}
	}

}
