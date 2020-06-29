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
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestDeliveryServiceServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		DeleteTestDeliveryServiceServers(t)
		AssignServersToTopologyBasedDeliveryService(t)
	})
}

func TestDeliveryServiceServersWithRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, Topologies, DeliveryServices, DeliveryServicesRequiredCapabilities, ServerServerCapabilities}, func() {
		CreateTestDeliveryServiceServersWithRequiredCapabilities(t)
		CreateTestMSODSServerWithReqCap(t)
	})
}

func AssignServersToTopologyBasedDeliveryService(t *testing.T) {
	ds, _, err := TOSession.GetDeliveryServiceByXMLIDNullable("ds-top", nil)
	if err != nil {
		t.Fatalf("cannot GET delivery service 'ds-top': %s", err.Error())
	}
	if len(ds) != 1 {
		t.Fatalf("expected one delivery service: 'ds-top', actual: %v", ds)
	}
	if ds[0].Topology == nil {
		t.Fatal("expected delivery service: 'ds-top' to have a non-nil Topology, actual: nil")
	}
	serversResp, _, err := TOSession.GetServers(nil, nil)
	servers := []tc.ServerNullable{}
	for _, s := range serversResp.Response {
		if s.CDNID != nil && *s.CDNID == *ds[0].CDNID && s.Type == tc.CacheTypeEdge.String() {
			servers = append(servers, s)
		}
	}
	if len(servers) < 1 {
		t.Fatalf("expected: at least one EDGE in cdn %s, actual: %v", *ds[0].CDNName, servers)
	}
	serverNames := []string{}
	for _, s := range servers {
		if s.CDNID != nil && *s.CDNID == *ds[0].CDNID && s.Type == tc.CacheTypeEdge.String() && s.HostName != nil {
			serverNames = append(serverNames, *s.HostName)
		} else {
			t.Fatalf("expected only EDGE servers in cdn '%s', actual: %v", *ds[0].CDNName, servers)
		}
	}
	_, reqInf, err := TOSession.AssignServersToDeliveryService(serverNames, "ds-top")
	if err == nil {
		t.Fatal("assigning servers to topology-based delivery service - expected: error, actual: nil error")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("assigning servers to topology-based delivery service - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}

	_, reqInf, err = TOSession.CreateDeliveryServiceServers(*ds[0].ID, []int{*servers[0].ID}, false)
	if err == nil {
		t.Fatal("creating deliveryserviceserver assignment for topology-based delivery service - expected: error, actual: nil error")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("creating deliveryserviceserver assignment for topology-based delivery service - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}
}

func CreateTestDeliveryServiceServersWithRequiredCapabilities(t *testing.T) {
	sscs := testData.ServerServerCapabilities

	testCases := []struct {
		ds          string
		serverName  string
		ssc         tc.ServerServerCapability
		description string
		err         error
		capability  tc.DeliveryServicesRequiredCapability
	}{
		{
			serverName:  "atlanta-edge-01",
			description: "missing requirements for server -> DS assignment",
			err:         errors.New(`Caching server cannot be assigned to this delivery service without having the required delivery service capabilities`),
			ssc:         sscs[0],
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  helperGetDeliveryServiceID(t, util.StrPtr("ds-test-minor-versions")),
				RequiredCapability: sscs[1].ServerCapability,
			},
		},
		{
			serverName:  "atlanta-mid-01",
			description: "successful server -> DS assignment",
			err:         nil,
			ssc:         sscs[1],
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  helperGetDeliveryServiceID(t, util.StrPtr("ds3")),
				RequiredCapability: sscs[1].ServerCapability,
			},
		},
	}

	for _, ctc := range testCases {
		t.Run(ctc.description, func(t *testing.T) {
			params := url.Values{}
			params.Add("hostName", ctc.serverName)
			resp, _, err := TOSession.GetServers(&params, nil)
			if err != nil {
				t.Fatalf("cannot GET Server by hostname: %v", err)
			}
			servers := resp.Response
			server := servers[0]
			if server.ID == nil {
				t.Fatalf("server %s had nil ID", ctc.serverName)
			}

			_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(ctc.capability)
			if err != nil {
				t.Fatalf("*POST delivery service required capability: %v", err)
			}

			ctc.ssc.ServerID = server.ID
			_, _, err = TOSession.CreateServerServerCapability(ctc.ssc)
			if err != nil {
				t.Fatalf("could not POST the server capability %v to server %v: %v", *ctc.ssc.ServerCapability, *ctc.ssc.Server, err)
			}

			_, _, got := TOSession.CreateDeliveryServiceServers(*ctc.capability.DeliveryServiceID, []int{*server.ID}, true)
			if (ctc.err == nil && got != nil) || (ctc.err != nil && !strings.Contains(got.Error(), ctc.err.Error())) {
				t.Fatalf("expected ctc.err to contain %v, got %v", ctc.err, got)
			}

			_, _, err = TOSession.DeleteDeliveryServicesRequiredCapability(*ctc.capability.DeliveryServiceID, *ctc.capability.RequiredCapability)
			if err != nil {
				t.Fatalf("*DELETE delivery service required capability: %v", err)
			}
		})
	}
}

func CreateTestMSODSServerWithReqCap(t *testing.T) {
	dsReqCap, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(nil, util.StrPtr("msods1"), nil, nil)
	if err != nil {
		t.Fatalf("GET delivery service required capabilites: %v", err)
	}

	if len(dsReqCap) == 0 {
		t.Fatal("no delivery service required capabilites found for ds msods1")
	}

	// Associate origin server to msods1 even though it does not have req cap
	// TODO: DON'T hard-code server hostnames!
	params := url.Values{}
	params.Add("hostName", "denver-mso-org-01")
	resp, _, err := TOSession.GetServers(&params, nil)
	if err != nil {
		t.Fatalf("GET server denver-mso-org-01: %v", err)
	}
	servers := resp.Response
	if len(servers) != 1 {
		t.Fatal("expected 1 server with hostname denver-mso-org-01")
	}

	s := servers[0]
	if s.ID == nil {
		t.Fatal("server 'denver-mso-org-01' had nil ID")
	}

	// Make sure server has no caps to ensure test correctness
	sccs, _, err := TOSession.GetServerServerCapabilities(s.ID, nil, nil, nil)
	if err != nil {
		t.Fatalf("GET server server capabilities for denver-mso-org-01: %v", err)
	}
	if len(sccs) != 0 {
		t.Fatal("expected 0 server server capabilities for server denver-mso-org-01")
	}

	// Is origin included in eligible servers even though it doesnt have required capability
	eServers, _, err := TOSession.GetDeliveryServicesEligible(*dsReqCap[0].DeliveryServiceID, nil)
	if err != nil {
		t.Fatalf("GET delivery service msods1 eligible servers: %v", err)
	}
	found := false
	for _, es := range eServers {
		if es.HostName != nil && *es.HostName == "denver-mso-org-01" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find origin server denver-mso-org-01 to be in eligible server return even though it is missing a required capability")
	}

	if _, _, err = TOSession.CreateDeliveryServiceServers(*dsReqCap[0].DeliveryServiceID, []int{*s.ID}, true); err != nil {
		t.Fatalf("POST delivery service origin servers without capabilities: %v", err)
	}

	// Create new bogus server capability
	if _, _, err = TOSession.CreateServerCapability(tc.ServerCapability{Name: "newfun"}); err != nil {
		t.Fatalf("cannot CREATE newfun server capability: %v", err)
	}

	// Attempt to assign to DS should not fail
	if _, _, err = TOSession.CreateDeliveryServicesRequiredCapability(tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  dsReqCap[0].DeliveryServiceID,
		RequiredCapability: util.StrPtr("newfun"),
	}); err != nil {
		t.Fatalf("POST required capability newfun to ds msods1: %v", err)
	}

	// Remove required capablity
	if _, _, err = TOSession.DeleteDeliveryServicesRequiredCapability(*dsReqCap[0].DeliveryServiceID, "newfun"); err != nil {
		t.Fatalf("DELETE delivery service required capability: %v", err)
	}

	// Delete server capability
	if _, _, err = TOSession.DeleteServerCapability("newfun"); err != nil {
		t.Fatalf("DELETE newfun server capability: %v", err)
	}
}

func DeleteTestDeliveryServiceServers(t *testing.T) {
	ds, server := getServerAndDSofSameCDN(t)
	if server.ID == nil {
		t.Fatalf("Got a server with a nil ID: %+v", server)
	}
	if ds.ID == nil {
		t.Fatalf("Got a delivery service with a nil ID %+v", ds)
	}

	_, _, err := TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*server.ID}, true)
	if err != nil {
		t.Errorf("POST delivery service servers: %v", err)
	}

	dsServers, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("GET delivery service servers: %v", err)
	}

	found := false
	for _, dss := range dsServers.Response {
		if dss.DeliveryService != nil && *dss.DeliveryService == *ds.ID && dss.Server != nil && *dss.Server == *server.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("POST delivery service servers returned success, but ds-server not in GET")
	}

	if _, _, err := TOSession.DeleteDeliveryServiceServer(*ds.ID, *server.ID); err != nil {
		t.Errorf("DELETE delivery service server: %v", err)
	}

	dsServers, _, err = TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("GET delivery service servers: %v", err)
	}

	found = false
	for _, dss := range dsServers.Response {
		if dss.DeliveryService != nil && *dss.DeliveryService == *ds.ID && dss.Server != nil && *dss.Server == *server.ID {
			found = true
			break
		}
	}
	if found {
		t.Error("DELETE delivery service servers returned success, but still in GET")
	}
}

func getServerAndDSofSameCDN(t *testing.T) (tc.DeliveryServiceNullable, tc.ServerNullable) {
	dses, _, err := TOSession.GetDeliveryServicesNullable(nil)
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) < 1 {
		t.Fatal("GET DeliveryServices returned no dses, must have at least 1 to test ds-servers")
	}

	resp, _, err := TOSession.GetServers(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET Servers: %v", err)
	}
	servers := resp.Response
	if len(servers) < 1 {
		t.Fatal("GET Servers returned no dses, must have at least 1 to test ds-servers")
	}

	for _, ds := range dses {
		for _, s := range servers {
			if ds.CDNName != nil && s.CDNName != nil && *ds.CDNName == *s.CDNName {
				return ds, s
			}
		}
	}
	t.Fatal("expected at least one delivery service and server in the same CDN")

	return tc.DeliveryServiceNullable{}, tc.ServerNullable{}
}
