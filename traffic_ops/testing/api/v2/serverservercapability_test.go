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
	"github.com/apache/trafficcontrol/v6/lib/go-util"
)

func TestServerServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices, ServerCapabilities, DeliveryServicesRequiredCapabilities, ServerServerCapabilities}, func() {
		GetTestServerServerCapabilities(t)
	})
}

func CreateTestServerServerCapabilities(t *testing.T) {
	// Valid POSTs

	// loop through server ServerCapabilities, assign FKs and create
	for _, ssc := range testData.ServerServerCapabilities {
		servResp, _, err := TOSession.GetServerByHostName(*ssc.Server)
		if err != nil {
			t.Fatalf("cannot GET Server by hostname: %v - %v", *ssc.Server, err)
		}
		if len(servResp) != 1 {
			t.Fatalf("cannot GET Server by hostname: %v. Response did not include record.", *ssc.Server)
		}
		server := servResp[0]
		ssc.ServerID = &server.ID
		resp, _, err := TOSession.CreateServerServerCapability(ssc)
		if err != nil {
			t.Errorf("could not POST the server capability %v to server %v: %v", *ssc.ServerCapability, *ssc.Server, err)
		}
		t.Log("Response: ", server.HostName, " ", resp)
	}

	// Invalid POSTs

	ssc := testData.ServerServerCapabilities[0]

	// Attempt to assign already assigned server capability
	_, _, err := TOSession.CreateServerServerCapability(ssc)
	if err == nil {
		t.Error("expected to receive error when assigning a already assigned server capability\n")
	}

	// Attempt to assign a server capability with no ID
	sscNilID := tc.ServerServerCapability{
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability without a server ID\n")
	}

	// Attempt to assign a server capability with no server capability
	sscNilCapability := tc.ServerServerCapability{
		ServerID: ssc.ServerID,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server without a server capability\n")
	}

	// Attempt to assign a server capability with invalid server capability
	sscInvalidCapability := tc.ServerServerCapability{
		ServerID:         ssc.ServerID,
		ServerCapability: util.StrPtr("bogus"),
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a non existent server capability to a server\n")
	}

	// Attempt to assign a server capability with invalid server capability
	sscInvalidID := tc.ServerServerCapability{
		ServerID:         util.IntPtr(-1),
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a non existent server ID\n")
	}

	// Attempt to assign a server capability to a non MID/EDGE server
	servers, _, err := TOSession.GetServerByHostName("riak")
	if err != nil {
		t.Fatalf("cannot GET Server by hostname: %v - %v", *ssc.Server, err)
	}
	if len(servers) < 1 {
		t.Fatal("need at least one server to test invalid server type assignment")
	}

	sscInvalidType := tc.ServerServerCapability{
		ServerID:         &servers[0].ID,
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidType)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server with incorrect type\n")
	}
}

func GetTestServerServerCapabilities(t *testing.T) {
	// Get All Server Capabilities
	sscs, _, err := TOSession.GetServerServerCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil\n")
	}
	if len(sscs) != len(testData.ServerServerCapabilities) {
		t.Errorf("expect %v server capabilities assigned to servers received %v ", len(testData.ServerServerCapabilities), len(sscs))
	}

	checkResp := func(t *testing.T, sscs []tc.ServerServerCapability) {
		if sscs == nil {
			t.Fatal("returned server capabilities assigned to servers was nil\n")
		}
		if len(sscs) != 1 {
			t.Errorf("expect 1 server capabilities assigned to server received %v ", len(sscs))
		}
	}

	for _, ssc := range sscs {
		// Get assigned Server Capabilities by server id
		sscs, _, err := TOSession.GetServerServerCapabilities(ssc.ServerID, nil, nil)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server ID %v: %v", *ssc.ServerID, err)
		}
		checkResp(t, sscs)
		// Get assigned Server Capabilities by host name
		sscs, _, err = TOSession.GetServerServerCapabilities(nil, ssc.Server, nil)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server host name %v: %v", *ssc.Server, err)
		}
		checkResp(t, sscs)

		// Get assigned Server Capabilities by server capability
		sscs, _, err = TOSession.GetServerServerCapabilities(nil, nil, ssc.ServerCapability)
		if err != nil {
			t.Fatalf("cannot GET server capabilities assigned to servers by server capability %v: %v", *ssc.ServerCapability, err)
		}
		checkResp(t, sscs)
	}
}

func DeleteTestServerServerCapabilities(t *testing.T) {
	// Get Server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil\n")
	}

	// Assign servers to DSes that have the capability required
	// Used to make sure we block server server_capability DELETE in that case
	dsServers := []tc.DeliveryServiceServer{}
	for _, ssc := range sscs {

		dsReqCapResp, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(nil, nil, ssc.ServerCapability)
		if err != nil {
			t.Fatalf("cannot GET delivery service required capabilities: %v", err)
		}
		if len(dsReqCapResp) == 0 {
			t.Fatalf("at least one delivery service needs the capability %v required", *ssc.ServerCapability)
		}
		dsReqCap := dsReqCapResp[0]

		// Assign server to ds
		_, err = TOSession.CreateDeliveryServiceServers(*dsReqCap.DeliveryServiceID, []int{*ssc.ServerID}, false)
		if err != nil {
			t.Fatalf("cannot CREATE server delivery service assignment: %v", err)
		}
		dsServers = append(dsServers, tc.DeliveryServiceServer{
			Server:          ssc.ServerID,
			DeliveryService: dsReqCap.DeliveryServiceID,
		})
	}

	// Delete should fail as their delivery services now require the capabilities
	for _, ssc := range sscs {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err == nil {
			t.Fatalf("should have gotten error when using DELETE on the server capability %v from server %v as it is required by associated dses", *ssc.ServerCapability, *ssc.Server)
		}
	}

	for _, dsServer := range dsServers {
		if dsServer.DeliveryService == nil {
			t.Error("nil DeliveryService property")
			continue
		}
		if dsServer.Server == nil {
			t.Error("nil Server property")
			continue
		}
		setInactive(t, *dsServer.DeliveryService)
		_, _, err := TOSession.DeleteDeliveryServiceServer(*dsServer.DeliveryService, *dsServer.Server)
		if err != nil {
			t.Fatalf("could not DELETE the server %v from ds %v: %v", *dsServer.Server, *dsServer.DeliveryService, err)
		}
	}

	// Remove the requirement so we can actually delete them

	for _, ssc := range sscs {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err != nil {
			t.Errorf("could not DELETE the server capability %v from server %v: %v", *ssc.ServerCapability, *ssc.Server, err)
		}
	}

}
