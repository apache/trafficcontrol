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
	"errors"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
)

func TestDeliveryServiceServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		DeleteTestDeliveryServiceServers(t)
	})
}

func TestDeliveryServiceServersWithRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, DeliveryServices, DeliveryServicesRequiredCapabilities, ServerServerCapabilities}, func() {
		CreateTestDeliveryServiceServersWithRequiredCapabilities(t)
		CreateTestMSODSServerWithReqCap(t)
	})
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
			err:         errors.New("cannot be assigned to this delivery service"),
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
			servers, _, err := TOSession.GetServerByHostName(ctc.serverName)
			if err != nil {
				t.Fatalf("cannot GET Server by hostname: %v", err)
			}
			server := servers[0]

			_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(ctc.capability)
			if err != nil {
				t.Fatalf("*POST delivery service required capability: %v", err)
			}

			ctc.ssc.ServerID = &server.ID
			_, _, err = TOSession.CreateServerServerCapability(ctc.ssc)
			if err != nil {
				t.Fatalf("could not POST the server capability %v to server %v: %v", *ctc.ssc.ServerCapability, *ctc.ssc.Server, err)
			}

			_, got := TOSession.CreateDeliveryServiceServers(*ctc.capability.DeliveryServiceID, []int{server.ID}, true)
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
	dsReqCap, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(nil, util.StrPtr("msods1"), nil)
	if err != nil {
		t.Fatalf("GET delivery service required capabilites: %v", err)
	}

	if len(dsReqCap) == 0 {
		t.Fatal("no delivery service required capabilites found for ds msods1")
	}

	// Associate origin server to msods1 even though it does not have req cap

	servers, _, err := TOSession.GetServerByHostName("denver-mso-org-01")
	if err != nil {
		t.Fatalf("GET server denver-mso-org-01: %v", err)
	}
	if len(servers) != 1 {
		t.Fatal("expected 1 server with hostname denver-mso-org-01")
	}

	s := servers[0]

	// Make sure server has no caps to ensure test correctness
	sccs, _, err := TOSession.GetServerServerCapabilities(&s.ID, nil, nil)
	if err != nil {
		t.Fatalf("GET server server capabilities for denver-mso-org-01: %v", err)
	}
	if len(sccs) != 0 {
		t.Fatal("expected 0 server server capabilities for server denver-mso-org-01")
	}

	// Is origin included in eligible servers even though it doesnt have required capability
	eServers, _, err := TOSession.GetDeliveryServicesEligible(*dsReqCap[0].DeliveryServiceID)
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

	if _, err = TOSession.CreateDeliveryServiceServers(*dsReqCap[0].DeliveryServiceID, []int{s.ID}, true); err != nil {
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

	_, err := TOSession.CreateDeliveryServiceServers(*ds.ID, []int{server.ID}, true)
	if err != nil {
		t.Errorf("POST delivery service servers: %v", err)
	}

	dsServers, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("GET delivery service servers: %v", err)
	}

	found := false
	for _, dss := range dsServers.Response {
		if *dss.DeliveryService == *ds.ID && *dss.Server == server.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("POST delivery service servers returned success, but ds-server not in GET")
	}

	setInactive(t, *ds.ID)
	if _, _, err := TOSession.DeleteDeliveryServiceServer(*ds.ID, server.ID); err != nil {
		t.Errorf("DELETE delivery service server: %v", err)
	}

	dsServers, _, err = TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("GET delivery service servers: %v", err)
	}

	found = false
	for _, dss := range dsServers.Response {
		if *dss.DeliveryService == *ds.ID && *dss.Server == server.ID {
			found = true
			break
		}
	}
	if found {
		t.Error("DELETE delivery service servers returned success, but still in GET")
	}
}

func getServerAndDSofSameCDN(t *testing.T) (tc.DeliveryServiceNullable, tc.Server) {
	dses, _, err := TOSession.GetDeliveryServicesNullable()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) < 1 {
		t.Fatal("GET DeliveryServices returned no dses, must have at least 1 to test ds-servers")
	}

	servers, _, err := TOSession.GetServers()
	if err != nil {
		t.Fatalf("cannot GET Servers: %v", err)
	}
	if len(servers) < 1 {
		t.Fatal("GET Servers returned no dses, must have at least 1 to test ds-servers")
	}

	for _, ds := range dses {
		for _, s := range servers {
			if ds.CDNName != nil && *ds.CDNName == s.CDNName {
				return ds, s
			}
		}
	}
	t.Fatal("expected at least one delivery service and server in the same CDN")

	return tc.DeliveryServiceNullable{}, tc.Server{}
}
