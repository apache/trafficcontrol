package tcdata

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
	"net/http"
	"net/url"
	"testing"

	"github.com/apache/trafficcontrol/v7/lib/go-tc"
	"github.com/apache/trafficcontrol/v7/lib/go-util"
)

func (r *TCData) CreateTestServerServerCapabilities(t *testing.T) {
	// Valid POSTs

	// loop through server ServerCapabilities, assign FKs and create
	params := url.Values{}
	for _, ssc := range r.TestData.ServerServerCapabilities {
		if ssc.Server == nil {
			t.Fatalf("server-server-capability structure had nil server")
		}
		params.Set("hostName", *ssc.Server)
		resp, _, err := TOSession.GetServersWithHdr(&params, nil)
		if err != nil {
			t.Fatalf("cannot GET Server by hostname '%s': %v - %v", *ssc.Server, err, resp.Alerts)
		}
		servResp := resp.Response
		if len(servResp) != 1 {
			t.Fatalf("cannot GET Server by hostname: %s. Response did not include record.", *ssc.Server)
		}
		server := servResp[0]
		ssc.ServerID = server.ID
		createResp, _, err := TOSession.CreateServerServerCapability(ssc)
		if err != nil {
			t.Errorf("could not POST the server capability %s to server %s: %v", *ssc.ServerCapability, *ssc.Server, err)
		}
		t.Logf("Response: %s %+v", *ssc.Server, createResp)
	}

	// Invalid POSTs

	ssc := r.TestData.ServerServerCapabilities[0]

	// Attempt to assign already assigned server capability
	_, _, err := TOSession.CreateServerServerCapability(ssc)
	if err == nil {
		t.Error("expected to receive error when assigning a already assigned server capability")
	}

	// Attempt to assign a server capability with no ID
	sscNilID := tc.ServerServerCapability{
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability without a server ID")
	}

	// Attempt to assign a server capability with no server capability
	sscNilCapability := tc.ServerServerCapability{
		ServerID: ssc.ServerID,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server without a server capability")
	}

	// Attempt to assign a server capability with invalid server capability
	sscInvalidCapability := tc.ServerServerCapability{
		ServerID:         ssc.ServerID,
		ServerCapability: util.StrPtr("bogus"),
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidCapability)
	if err == nil {
		t.Error("expected to receive error when assigning a non existent server capability to a server")
	}

	// Attempt to assign a server capability with invalid server capability
	sscInvalidID := tc.ServerServerCapability{
		ServerID:         util.IntPtr(-1),
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidID)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a non existent server ID")
	}

	// Attempt to assign a server capability to a non MID/EDGE server
	// TODO: DON'T hard-code server hostnames!
	params.Set("hostName", "trafficvault")
	resp, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("cannot GET Server by hostname 'trafficvault': %v - alerts: %+v", err, resp.Alerts)
	}
	servers := resp.Response
	if len(servers) < 1 {
		t.Fatal("need at least one server to test invalid server type assignment")
	}

	sscInvalidType := tc.ServerServerCapability{
		ServerID:         servers[0].ID,
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidType)
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server with incorrect type")
	}
}

func (r *TCData) DeleteTestServerServerCapabilities(t *testing.T) {
	// Get Server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilitiesWithHdr(nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil")
	}

	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET delivery services: %v", err)
	}
	dsIDtoDS := make(map[int]tc.DeliveryServiceNullableV30, len(dses))
	for _, ds := range dses {
		dsIDtoDS[*ds.ID] = ds
	}

	// Assign servers to DSes that have the capability required
	// Used to make sure we block server server_capability DELETE in that case
	dsServers := []tc.DeliveryServiceServer{}
	assignedServers := make(map[int]bool)
	for _, ssc := range sscs {

		dsReqCapResp, _, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(nil, nil, ssc.ServerCapability, nil)
		if err != nil {
			t.Fatalf("cannot GET delivery service required capabilities: %v", err)
		}
		if len(dsReqCapResp) == 0 {
			// capability is not required by any delivery service
			continue
		}
		var dsReqCap tc.DeliveryServicesRequiredCapability
		for _, dsrc := range dsReqCapResp {
			if dsIDtoDS[*dsrc.DeliveryServiceID].Topology == nil {
				dsReqCap = dsrc
				break
			}
		}
		if dsReqCap.DeliveryServiceID == nil {
			// didn't find a non-topology-based dsReqCap for this ssc
			continue
		}

		// Assign server to ds
		_, _, err = TOSession.CreateDeliveryServiceServers(*dsReqCap.DeliveryServiceID, []int{*ssc.ServerID}, false)
		if err != nil {
			t.Fatalf("cannot CREATE server delivery service assignment: %v", err)
		}
		dsServers = append(dsServers, tc.DeliveryServiceServer{
			Server:          ssc.ServerID,
			DeliveryService: dsReqCap.DeliveryServiceID,
		})
		assignedServers[*ssc.ServerID] = true
	}
	if len(dsServers) == 0 {
		t.Fatalf("test requires at least one server with a capability that is required by at least one delivery service")
	}

	// Delete should fail as their delivery services now require the capabilities
	for _, ssc := range sscs {
		if assignedServers[*ssc.ServerID] {
			_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
			if err == nil {
				t.Fatalf("should have gotten error when using DELETE on the server capability %s from server %s as it is required by associated dses", *ssc.ServerCapability, *ssc.Server)
			}
		}
	}

	for _, dsServer := range dsServers {
		_, _, err := TOSession.DeleteDeliveryServiceServer(*dsServer.DeliveryService, *dsServer.Server)
		if err != nil {
			t.Fatalf("could not DELETE the server %d from ds %d: %v", *dsServer.Server, *dsServer.DeliveryService, err)
		}
	}

	// Remove the requirement so we can actually delete them

	for _, ssc := range sscs {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err != nil {
			t.Errorf("could not DELETE the server capability %s from server %s: %v", *ssc.ServerCapability, *ssc.Server, err)
		}
	}

}

func DeleteTestServerServerCapabilitiesForTopologiesValidation(t *testing.T) {
	// dtrc-edge-01 and dtrc-edge-02 (capabilities = ram, disk) are assigned to
	// ds-top-req-cap (topology = top-for-ds-req; required capabilities = ram, disk) and
	// ds-top-req-cap2 (topology = top-for-ds-req2; required capabilities = ram)
	var edge1 tc.ServerV30
	var edge2 tc.ServerV30

	servers, _, err := TOSession.GetServersWithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET servers: %v", err)
	}
	for _, s := range servers.Response {
		if *s.HostName == "dtrc-edge-01" {
			edge1 = s
		}
		if *s.HostName == "dtrc-edge-02" {
			edge2 = s
		}
	}
	if edge1.HostName == nil || edge2.HostName == nil {
		t.Fatalf("expected servers with hostName dtrc-edge-01 and dtrc-edge-02")
	}

	// delete should succeed because dtrc-edge-02 still has the required capabilities
	// for ds-top-req-cap and ds-top-req-cap2 within the cachegroup
	_, _, err = TOSession.DeleteServerServerCapability(*edge1.ID, "ram")
	if err != nil {
		t.Fatalf("when deleting server server capability, expected: nil error, actual: %v", err)
	}

	// delete should fail because dtrc-edge-02 is the last server in the cachegroup that
	// has ds-top-req-cap's required capabilities
	_, reqInf, err := TOSession.DeleteServerServerCapability(*edge2.ID, "ram")
	if err == nil {
		t.Fatalf("when deleting server server capability, expected: error, actual: nil")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("when deleting server server capability, expected status code: %d, actual: %d", http.StatusBadRequest, reqInf.StatusCode)
	}

	// delete should fail because dtrc-edge-02 is the last server in the cachegroup that
	// has ds-top-req-cap's required capabilities
	_, r, err := TOSession.DeleteServerServerCapability(*edge2.ID, "disk")
	if err == nil {
		t.Fatalf("when deleting required server server capability, expected: error, actual: nil")
	}
	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("when deleting required server server capability, expected status code: %d, actual: %d", http.StatusBadRequest, reqInf.StatusCode)
	}

	// delete should succeed because dtrc-edge-02 still has the required capabilities
	// for ds-top-req-cap and ds-top-req-cap2 within the cachegroup
	_, _, err = TOSession.DeleteServerServerCapability(*edge1.ID, "disk")
	if err != nil {
		t.Fatalf("when deleting server server capability, expected: nil error, actual: %v", err)
	}
}

func DeleteTestServerServerCapabilitiesForTopologies(t *testing.T) {
	// Get Server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilitiesWithHdr(nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("cannot GET server capabilities assigned to servers: %v", err)
	}
	if sscs == nil {
		t.Fatal("returned server capabilities assigned to servers was nil")
	}

	for _, ssc := range sscs {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err != nil {
			t.Errorf("could not DELETE the server capability %s from server %s: %v", *ssc.ServerCapability, *ssc.Server, err)
		}
	}
}
