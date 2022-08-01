package v4

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
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestServerServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, ServerCapabilities, ServerServerCapabilities, DeliveryServicesRequiredCapabilities}, func() {
		SortTestServerServerCapabilities(t)
		GetTestServerServerCapabilitiesIMS(t)
		GetTestServerServerCapabilities(t)
		GetDeliveryServiceServersWithCapabilities(t)
		UpdateTestServerServerCapabilities(t)
		GetTestPaginationSupportSsc(t)
		AssignMultipleTestServerCapabilities(t)
		DeleteTestServerServerCapabilityWithInvalidData(t)
		DeleteTestServerServerCapabilitiesForTopologiesValidation(t)
	})
}

func GetTestServerServerCapabilitiesIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	rfcTime := futureTime.Format(time.RFC1123)
	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, rfcTime)
	resp, reqInf, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestServerServerCapabilities(t *testing.T) {
	// Valid POSTs

	// loop through server ServerCapabilities, assign FKs and create
	opts := client.NewRequestOptions()
	for _, ssc := range testData.ServerServerCapabilities {
		if ssc.Server == nil || ssc.ServerCapability == nil {
			t.Fatalf("server-server-capability structure had nil server and/or Capability")
		}
		opts.QueryParameters.Set("hostName", *ssc.Server)
		resp, _, err := TOSession.GetServers(opts)
		if err != nil {
			t.Fatalf("cannot get Servers filtered by Host Name '%s': %v - alerts: %+v", *ssc.Server, err, resp.Alerts)
		}
		servResp := resp.Response
		if len(servResp) != 1 {
			t.Fatalf("cannot GET Server by hostname: %v. Response did not include record.", *ssc.Server)
		}
		server := servResp[0]
		ssc.ServerID = server.ID
		createResp, _, err := TOSession.CreateServerServerCapability(ssc, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not associate Capability '%s' with server '%s': %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, err, createResp.Alerts)
		}
	}

	// Invalid POSTs
	if len(testData.ServerServerCapabilities) < 1 {
		t.Fatal("Need at least one server/Capability relationship to test creating server/Capability associations")
	}
	ssc := testData.ServerServerCapabilities[0]

	// Attempt to assign already assigned server capability
	_, _, err := TOSession.CreateServerServerCapability(ssc, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error when assigning a already assigned server capability")
	}

	// Attempt to assign a server capability with no ID
	sscNilID := tc.ServerServerCapability{
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilID, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error when assigning a server capability without a server ID")
	}

	// Attempt to assign a server capability with no server capability
	sscNilCapability := tc.ServerServerCapability{
		ServerID: ssc.ServerID,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscNilCapability, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server without a server capability")
	}

	// Attempt to assign a server capability with invalid server capability
	sscInvalidCapability := tc.ServerServerCapability{
		ServerID:         ssc.ServerID,
		ServerCapability: util.StrPtr("bogus"),
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidCapability, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error when assigning a non existent server capability to a server")
	}

	// Attempt to assign a server capability with invalid server capability
	sscInvalidID := tc.ServerServerCapability{
		ServerID:         util.IntPtr(-1),
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidID, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a non existent server ID")
	}

	// Attempt to assign a server capability to a non MID/EDGE server
	// TODO: DON'T hard-code server hostnames!
	opts.QueryParameters.Set("hostName", "trafficvault")
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot get Servers filtered by hostname 'trafficvault': %v - alerts: %+v", err, resp.Alerts)
	}
	servers := resp.Response
	if len(servers) < 1 {
		t.Fatal("need at least one server to test invalid server type assignment")
	}

	sscInvalidType := tc.ServerServerCapability{
		ServerID:         servers[0].ID,
		ServerCapability: ssc.ServerCapability,
	}
	_, _, err = TOSession.CreateServerServerCapability(sscInvalidType, client.RequestOptions{})
	if err == nil {
		t.Error("expected to receive error when assigning a server capability to a server with incorrect type")
	}
}

func SortTestServerServerCapabilities(t *testing.T) {
	resp, _, err := TOSession.GetServerServerCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}

	sortedList := make([]string, 0, len(resp.Response))
	for _, ssc := range resp.Response {
		if ssc.Server == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server")
			continue
		}
		sortedList = append(sortedList, *ssc.Server)
	}

	if !sort.StringsAreSorted(sortedList) {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func GetTestServerServerCapabilities(t *testing.T) {
	// Get All Server Capabilities
	sscs, _, err := TOSession.GetServerServerCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot server capability/server relationships: %v - alerts: %+v", err, sscs.Alerts)
	}
	if len(sscs.Response) != len(testData.ServerServerCapabilities) {
		t.Errorf("expect %v server capabilities assigned to servers received %v ", len(testData.ServerServerCapabilities), len(sscs.Response))
	}

	opts := client.NewRequestOptions()
	for _, ssc := range sscs.Response {
		if ssc.Server == nil || ssc.ServerID == nil || ssc.ServerCapability == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server and/or server ID and/or Capability")
			continue
		}
		// Get assigned Server Capabilities by server id
		opts.QueryParameters.Set("serverId", strconv.Itoa(*ssc.ServerID))
		sscs, _, err := TOSession.GetServerServerCapabilities(opts)
		opts.QueryParameters.Del("serverId")
		if err != nil {
			t.Fatalf("cannot get Capabilities assigned to server #%d: %v - alerts: %+v", *ssc.ServerID, err, sscs.Alerts)
		}
		for _, s := range sscs.Response {
			if s.ServerID == nil {
				t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server ID")
			} else if *s.ServerID != *ssc.ServerID {
				t.Errorf("GET server server capabilities by serverID returned non-matching server ID: %d", *s.ServerID)
			}
		}

		// Get assigned Server Capabilities by host name
		opts.QueryParameters.Set("serverHostName", *ssc.Server)
		sscs, _, err = TOSession.GetServerServerCapabilities(opts)
		opts.QueryParameters.Del("serverHostName")
		if err != nil {
			t.Fatalf("cannot get Capabilities assigned server '%s': %v alerts: %+v", *ssc.Server, err, sscs.Alerts)
		}
		for _, s := range sscs.Response {
			if s.Server == nil {
				t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server")
			} else if *s.Server != *ssc.Server {
				t.Errorf("GET server server capabilities by serverHostName returned non-matching server hostname: %s", *s.Server)
			}
		}

		// Get assigned Server Capabilities by server capability
		opts.QueryParameters.Set("serverCapability", *ssc.ServerCapability)
		sscs, _, err = TOSession.GetServerServerCapabilities(opts)
		opts.QueryParameters.Del("serverCapability")
		if err != nil {
			t.Fatalf("cannot get Capability/server associations for Capability '%s': %v alerts: %+v", *ssc.ServerCapability, err, sscs.Alerts)
		}
		for _, s := range sscs.Response {
			if s.ServerCapability == nil {
				t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined Capability")
			} else if *s.ServerCapability != *ssc.ServerCapability {
				t.Errorf("GET server server capabilities by server capability returned non-matching server capability: %s", *s.ServerCapability)
			}
		}
	}
}

func UpdateTestServerServerCapabilities(t *testing.T) {
	// Get server capability name and edit it to a new name
	resp, _, err := TOSession.GetServerCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("no server capability in response, quitting")
	}
	originalName := resp.Response[0].Name
	newSCName := "sc-test"
	resp.Response[0].Name = newSCName

	// Get all servers related to original sever capability name
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("serverCapability", originalName)
	servOrigResp, _, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get Capability/server associations for Capability '%s': %v alerts: %+v", originalName, err, servOrigResp.Alerts)
	}
	if len(servOrigResp.Response) == 0 {
		t.Fatalf("no servers associated with server capability name: %v", originalName)
	}
	mapOrigServ := make(map[string]string, len(servOrigResp.Response))
	for _, s := range servOrigResp.Response {
		if s.Server == nil || s.ServerCapability == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server and/or Capability")
			continue
		}
		mapOrigServ[*s.Server] = *s.ServerCapability
	}

	// Update server capability with new name
	updateResponse, _, err := TOSession.UpdateServerCapability(originalName, resp.Response[0], client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Server Capability: %v - alerts: %+v", err, updateResponse.Alerts)
	}

	//To check whether the primary key change trickled down to server table
	opts.QueryParameters.Set("serverCapability", newSCName)
	servUpdatedResp, _, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get Capability/server associations for Capability '%s': %v alerts: %+v", newSCName, err, servUpdatedResp.Alerts)
	}
	if len(servUpdatedResp.Response) == 0 {
		t.Fatalf("no server associated with server capability '%s'", newSCName)
	}
	if len(servOrigResp.Response) != len(servUpdatedResp.Response) {
		t.Fatalf("length of servers for a given server capability name is different, expected: %s-%d, got: %s-%d", originalName, len(servOrigResp.Response), newSCName, len(servUpdatedResp.Response))
	}
	for _, s := range servUpdatedResp.Response {
		if s.ServerCapability == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server")
			continue
		}
		if newSCName != *s.ServerCapability {
			t.Errorf("GET server server capabilities by server capability returned non-matching server capability: %s", *s.ServerCapability)
		}
		_, ok := mapOrigServ[*s.Server]
		if !ok {
			t.Fatalf("server capability name change didn't trickle to server: %v", *s.Server)
		}
	}

	// Set everything back as it was for further testing.
	resp.Response[0].Name = originalName
	r, _, err := TOSession.UpdateServerCapability(newSCName, resp.Response[0], client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Server Capability: %v - alerts: %+v", err, r.Alerts)
	}
}

func AssignMultipleTestServerCapabilities(t *testing.T) {
	//Get list of server capabilities
	resp, _, err := TOSession.GetServerCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("no server capability in response, quitting")
	}

	originalName := resp.Response[0].Name
	var multipleSCs []string
	multipleSCs = append(multipleSCs, resp.Response[1].Name, resp.Response[2].Name)

	// Get all servers related to original sever capability name
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("serverCapability", originalName)
	servOrigResp, _, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get Capability/server associations for Capability '%s': %v alerts: %+v", originalName, err, servOrigResp.Alerts)
	}
	if len(servOrigResp.Response) == 0 {
		t.Fatalf("no servers associated with server capability name: %v", originalName)
	}
	origServerID := *servOrigResp.Response[3].ServerID

	msc := tc.MultipleServerCapabilities{
		ServerCapability: multipleSCs,
		ServerID:         &origServerID,
	}
	_, reqInf, err := TOSession.AssignMultipleServerCapability(msc, client.NewRequestOptions(), origServerID)
	if err != nil {
		t.Fatalf("unable to assign update multiple server capabilities to server: %d, error: %v", origServerID, err.Error())
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected:%v, got:%v", http.StatusOK, reqInf.StatusCode)
	}
}

func DeleteTestServerServerCapabilities(t *testing.T) {
	// Get Server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get server/Capability associations: %v - alerts: %+v", err, sscs.Alerts)
	}

	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	dsIDtoDS := make(map[int]tc.DeliveryServiceV4, len(dses.Response))
	for _, ds := range dses.Response {
		if ds.ID == nil {
			t.Error("Traffic Ops responded with a representation of a Delivery Service that had null or undefined ID")
			continue
		}
		dsIDtoDS[*ds.ID] = ds
	}

	// Assign servers to DSes that have the capability required
	// Used to make sure we block server server_capability DELETE in that case
	dsServers := []tc.DeliveryServiceServer{}
	assignedServers := make(map[int]bool)
	opts := client.NewRequestOptions()
	for _, ssc := range sscs.Response {
		if ssc.ServerCapability == nil {
			t.Error("Traffic Ops returned a representation of a Server/Capability relationship with null or undefined Capability")
			continue
		}
		opts.QueryParameters.Set("requiredCapability", *ssc.ServerCapability)

		dsReqCapResp, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
		if err != nil {
			t.Fatalf("Unexpected error retrieving relationships between Delivery Services and the Capabilities they require filtered by Capability '%s': %v - alerts: %+v", *ssc.ServerCapability, err, dsReqCapResp.Alerts)
		}
		if len(dsReqCapResp.Response) == 0 {
			// capability is not required by any delivery service
			continue
		}
		var dsReqCap tc.DeliveryServicesRequiredCapability
		for _, dsrc := range dsReqCapResp.Response {
			if dsrc.DeliveryServiceID == nil {
				t.Error("Traffic Ops returned a representation of a Delivery Service/Required Capability relationship with null or undefined Delivery Service ID")
				continue
			}
			if ds, ok := dsIDtoDS[*dsrc.DeliveryServiceID]; ok {
				if ds.Topology == nil {
					dsReqCap = dsrc
					break
				}
			} else {
				t.Errorf("Traffic Ops reports that Delivery Service #%d requires one or more capabilities, but also does not report that said Delivery Service exists", *dsrc.DeliveryServiceID)
			}
		}
		if dsReqCap.DeliveryServiceID == nil {
			// didn't find a non-topology-based dsReqCap for this ssc
			continue
		}

		// Assign server to ds
		assignResp, _, err := TOSession.CreateDeliveryServiceServers(*dsReqCap.DeliveryServiceID, []int{*ssc.ServerID}, false, client.RequestOptions{})
		if err != nil {
			t.Fatalf("Unexpected error retrieving server-to-Delivery-Service assignments: %v - alerts: %+v", err, assignResp.Alerts)
		}
		dsServers = append(dsServers, tc.DeliveryServiceServer{
			Server:          ssc.ServerID,
			DeliveryService: dsReqCap.DeliveryServiceID,
		})
		assignedServers[*ssc.ServerID] = true
	}

	// Delete should fail as their delivery services now require the capabilities
	for _, ssc := range sscs.Response {
		if ssc.ServerID == nil || ssc.ServerCapability == nil || ssc.Server == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server and/or server ID and/or Capability")
			continue
		}
		if assignedServers[*ssc.ServerID] {
			_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability, client.RequestOptions{})
			if err == nil {
				t.Fatalf("should have gotten error when removing Capability '%s' from server '%s' (#%d) as it is required by associated Delivery Services", *ssc.ServerCapability, *ssc.Server, *ssc.ServerID)
			}
		}
	}

	for _, dsServer := range dsServers {
		setInactive(t, *dsServer.DeliveryService)
		alerts, _, err := TOSession.DeleteDeliveryServiceServer(*dsServer.DeliveryService, *dsServer.Server, client.RequestOptions{})
		if err != nil {
			t.Fatalf("could not remove server #%d from Delivery Service #%d: %v - alerts: %+v", *dsServer.Server, *dsServer.DeliveryService, err, alerts.Alerts)
		}
	}

	// Remove the requirement so we can actually delete them

	for _, ssc := range sscs.Response {
		if ssc.ServerID == nil || ssc.ServerCapability == nil || ssc.Server == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server and/or server ID and/or Capability")
			continue
		}
		alerts, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not remove Capability '%s' from server '%s' (#%d): %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, *ssc.ServerID, err, alerts.Alerts)
		}
	}
}

func DeleteTestServerServerCapabilitiesForTopologiesValidation(t *testing.T) {
	// dtrc-edge-01 and dtrc-edge-02 (capabilities = ram, disk) are assigned to
	// ds-top-req-cap (topology = top-for-ds-req; required capabilities = ram, disk) and
	// ds-top-req-cap2 (topology = top-for-ds-req2; required capabilities = ram)
	var edge1 tc.ServerV4
	var edge2 tc.ServerV4

	servers, _, err := TOSession.GetServers(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get servers: %v - alerts: %+v", err, servers.Alerts)
	}
	for _, s := range servers.Response {
		if s.HostName == nil || s.ID == nil {
			t.Fatal("Traffic Ops returned a representation for a server with null or undefined ID and/or Host Name")
		}
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
	alerts, _, err := TOSession.DeleteServerServerCapability(*edge1.ID, "ram", client.RequestOptions{})
	if err != nil {
		t.Fatalf("when deleting server server capability, expected: nil error, actual: %v - alerts: %+v", err, alerts.Alerts)
	}

	// delete should fail because dtrc-edge-02 is the last server in the cachegroup that
	// has ds-top-req-cap's required capabilities
	_, reqInf, err := TOSession.DeleteServerServerCapability(*edge2.ID, "ram", client.RequestOptions{})
	if err == nil {
		t.Fatalf("when deleting server server capability, expected: error, actual: nil")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("when deleting server server capability, expected status code: %d, actual: %d", http.StatusBadRequest, reqInf.StatusCode)
	}

	// delete should fail because dtrc-edge-02 is the last server in the cachegroup that
	// has ds-top-req-cap's required capabilities
	_, r, err := TOSession.DeleteServerServerCapability(*edge2.ID, "disk", client.RequestOptions{})
	if err == nil {
		t.Fatalf("when deleting required server server capability, expected: error, actual: nil")
	}
	if r.StatusCode != http.StatusBadRequest {
		t.Errorf("when deleting required server server capability, expected status code: %d, actual: %d", http.StatusBadRequest, reqInf.StatusCode)
	}

	// delete should succeed because dtrc-edge-02 still has the required capabilities
	// for ds-top-req-cap and ds-top-req-cap2 within the cachegroup
	alerts, _, err = TOSession.DeleteServerServerCapability(*edge1.ID, "disk", client.RequestOptions{})
	if err != nil {
		t.Fatalf("when deleting server server capability, expected: nil error, actual: %v - alerts: %+v", err, alerts.Alerts)
	}
}

func DeleteTestServerServerCapabilitiesForTopologies(t *testing.T) {
	// Get Server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get server/Capability associations: %v - alerts: %+v", err, sscs.Alerts)
	}

	for _, ssc := range sscs.Response {
		if ssc.ServerID == nil || ssc.ServerCapability == nil || ssc.Server == nil {
			t.Error("Traffic Ops returned a representation of a relationship between a Server and one of its Capabilities with null or undefined server and/or server ID and/or Capability")
			continue
		}
		resp, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not remove Capability '%s' from server '%s': %v - alerts: %+v", *ssc.ServerCapability, *ssc.Server, err, resp.Alerts)
		}
	}
}

func GetDeliveryServiceServersWithCapabilities(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{QueryParameters: url.Values{"xmlId": []string{"ds4"}}})
	if err != nil {
		t.Fatalf("Failed to get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 1 {
		t.Fatal("Failed to get at least one Delivery Service")
	}

	ds := dses.Response[0]
	if ds.ID == nil {
		t.Fatal("Got Delivery Service with nil ID")
	}

	// Get an edge
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("hostName", "atlanta-edge-16")
	rs, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v - alerts: %+v", err, rs.Alerts)
	} else if len(rs.Response) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	edgeID := *rs.Response[0].ID

	// Get a MID
	opts.QueryParameters.Set("hostName", "atlanta-mid-02")
	rs, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v", err)
	} else if len(rs.Response) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	midID := *rs.Response[0].ID
	// assign edge and mid
	assignResp, _, err := TOSession.CreateDeliveryServiceServers(*ds.ID, []int{edgeID, midID}, true, client.RequestOptions{})
	if err != nil {
		t.Fatalf("expected no error while assigning servers to Delivery Service, but got: %v - alerts: %+v", err, assignResp.Alerts)
	}
	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Add("dsId", strconv.Itoa(*ds.ID))
	servers, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get server by Delivery Service ID: %v - alerts: %+v", err, servers.Alerts)
	}
	if len(servers.Response) != 2 {
		t.Fatalf("expected to get 2 servers for Delivery Service: %d, actual: %d", *ds.ID, len(servers.Response))
	}

	// now assign a capability
	reqCap := tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  ds.ID,
		RequiredCapability: util.StrPtr("blah"),
	}
	_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(reqCap, client.RequestOptions{})
	// this should fail because the mid doesn't have the reqd capability
	if err == nil {
		t.Fatalf("expected error creating DS reqd capability, but got nothing")
	}
	ssc := tc.ServerServerCapability{
		ServerID:         &midID,
		ServerCapability: util.StrPtr("blah"),
	}
	// assign the capability to the mid
	_, _, err = TOSession.CreateServerServerCapability(ssc, client.RequestOptions{})
	if err != nil {
		t.Fatalf("couldn't assign server capability to server with ID %d, err: %s", midID, err.Error())
	}
	resp, _, err := TOSession.CreateDeliveryServicesRequiredCapability(reqCap, client.RequestOptions{})
	// this should pass now because the mid has the reqd capability
	if err != nil {
		t.Fatalf("expected no error creating DS reqd capability, but got: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters.Set("dsId", strconv.Itoa(*ds.ID))
	servers, _, err = TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("Failed to get server by Delivery Service ID: %v", err)
	}
	if len(servers.Response) != 2 {
		t.Fatalf("expected to get 2 servers for Delivery Service: %d, actual: %d", *ds.ID, len(servers.Response))
	}
	alerts, _, err := TOSession.DeleteDeliveryServiceServer(*ds.ID, edgeID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error removing server #%d from Delivery Service #%d: %v - alerts: %+v", edgeID, *ds.ID, err, alerts.Alerts)
	}
	alerts, _, err = TOSession.DeleteDeliveryServiceServer(*ds.ID, midID, client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error removing server #%d from Delivery Service #%d: %v - alerts: %+v", midID, *ds.ID, err, alerts.Alerts)
	}
}

func GetTestPaginationSupportSsc(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get server server capabilities: %v - alerts: %+v", err, resp.Alerts)
	}
	ServerServerCapabilities := resp.Response
	if len(ServerServerCapabilities) < 3 {
		t.Fatalf("Need at least 3 server server capabilities in Traffic Ops to test pagination support, found: %d", len(ServerServerCapabilities))
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	serverServerCapWithLimit, _, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get server server capabilities by Order and Limit: %v - alerts: %+v", err, resp.Alerts)
	}
	if !reflect.DeepEqual(ServerServerCapabilities[:1], serverServerCapWithLimit.Response) {
		t.Error("expected GET Server Server Capabilities with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "1")
	serverServerCapWithOffset, _, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get server server capabilities by Order, Limit and Offset: %v - alerts: %+v", err, resp.Alerts)
	}
	if !reflect.DeepEqual(ServerServerCapabilities[1:2], serverServerCapWithOffset.Response) {
		t.Error("expected GET server server capabilities with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "2")
	serverServerCapWithPage, _, err := TOSession.GetServerServerCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get server server capabilities by Order, Limit and Page: %v - alerts: %+v", err, resp.Alerts)
	}
	if !reflect.DeepEqual(ServerServerCapabilities[1:2], serverServerCapWithPage.Response) {
		t.Error("expected GET Server Server Capabilities with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetServerServerCapabilities(opts)
	if err == nil {
		t.Error("expected GET Server Server Capabilities to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET Server Server Capabilities to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetServerServerCapabilities(opts)
	if err == nil {
		t.Error("expected GET Server Server Capabilities to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Server Server Capabilities to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetServerServerCapabilities(opts)
	if err == nil {
		t.Error("expected GET Server Server Capabilities to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Server Server Capabilities to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}

func DeleteTestServerServerCapabilityWithInvalidData(t *testing.T) {
	// Get Server server Capabilities to delete them
	sscs, _, err := TOSession.GetServerServerCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get server/Capability associations: %v - alerts: %+v", err, sscs.Alerts)
	}
	if len(sscs.Response) < 1 {
		t.Fatalf("No Server Server Capability available to test invalid scenario")
	}
	ssc := sscs.Response[0]
	if ssc.ServerID == nil {
		t.Fatal("Cache Group selected for testing had a null or undefined name")
	}

	//Delete Server Server Capability with Invalid Server Capability
	alerts, _, err := TOSession.DeleteServerServerCapability(*sscs.Response[0].ServerID, "abcd", client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected no server server_capability with that key found, actual: %v - alerts: %+v", err, alerts.Alerts)
	}

	//Missing Server Capability
	alerts, _, err = TOSession.DeleteServerServerCapability(*sscs.Response[0].ServerID, "", client.RequestOptions{})
	if err == nil {
		t.Fatalf("Expected missing key: serverCapability, actual: %v - alerts: %+v", err, alerts.Alerts)
	}
}
