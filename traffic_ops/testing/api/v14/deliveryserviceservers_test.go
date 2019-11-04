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
	"errors"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestDeliveryServiceServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		DeleteTestDeliveryServiceServers(t)
	})
}

func TestDeliveryServiceServersWithRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, DeliveryServices, DeliveryServicesRequiredCapabilities, ServerServerCapabilities}, func() {
		CreateTestDeliveryServiceServersWithRequiredCapabilities(t)
	})
}

func CreateTestDeliveryServiceServersWithRequiredCapabilities(t *testing.T) {
	dses, _ := getServersAndDSes(t)
	sscs := testData.ServerServerCapabilities

	testCases := []struct {
		ds          tc.DeliveryService
		serverName  string
		ssc         tc.ServerServerCapability
		description string
		err         error
		capability  tc.DeliveryServicesRequiredCapability
	}{
		{
			ds:          dses[1],
			serverName:  "atlanta-edge-01",
			description: "missing requirements for server -> DS assignment",
			err:         errors.New(`Caching server cannot be assigned to this delivery service without having the required delivery service capabilities`),
			ssc:         sscs[0],
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  &dses[1].ID,
				RequiredCapability: sscs[1].ServerCapability,
			},
		},
		{
			ds:          dses[0],
			serverName:  "atlanta-mid-01",
			description: "successful server -> DS assignment",
			err:         nil,
			ssc:         sscs[1],
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  &dses[0].ID,
				RequiredCapability: sscs[1].ServerCapability,
			},
		},
	}

	for _, ctc := range testCases {
		t.Run(ctc.description, func(t *testing.T) {

			servers, _, err := TOSession.GetServerByHostName(ctc.serverName)
			if err != nil {
				t.Fatalf("cannot GET Server by hostname: %v\n", err)
			}
			server := servers[0]

			_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(ctc.capability)
			if err != nil {
				t.Fatalf("*POST delivery service required capability: %v\n", err)
			}

			ctc.ssc.ServerID = &server.ID
			_, _, err = TOSession.CreateServerServerCapability(ctc.ssc)
			if err != nil {
				t.Fatalf("could not POST the server capability %v to server %v: %v\n", *ctc.ssc.ServerCapability, *ctc.ssc.Server, err)
			}

			_, got := TOSession.CreateDeliveryServiceServers(ctc.ds.ID, []int{server.ID}, true)
			if (ctc.err == nil && got != nil) || (ctc.err != nil && !strings.Contains(got.Error(), ctc.err.Error())) {
				t.Fatalf("expected ctc.err to contain %v, got %v\n", ctc.err, got)
			}

			_, _, err = TOSession.DeleteDeliveryServicesRequiredCapability(*ctc.capability.DeliveryServiceID, *ctc.capability.RequiredCapability)
			if err != nil {
				t.Fatalf("*DELETE delivery service required capability: %v\n", err)
			}
		})
	}
}

func DeleteTestDeliveryServiceServers(t *testing.T) {
	dses, servers := getServersAndDSes(t)
	ds, server := dses[0], servers[0]

	_, err := TOSession.CreateDeliveryServiceServers(ds.ID, []int{server.ID}, true)
	if err != nil {
		t.Errorf("POST delivery service servers: %v\n", err)
	}

	dsServers, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("GET delivery service servers: %v\n", err)
	}

	found := false
	for _, dss := range dsServers.Response {
		if *dss.DeliveryService == ds.ID && *dss.Server == server.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("POST delivery service servers returned success, but ds-server not in GET")
	}

	if _, _, err := TOSession.DeleteDeliveryServiceServer(ds.ID, server.ID); err != nil {
		t.Errorf("DELETE delivery service server: %v\n", err)
	}

	dsServers, _, err = TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("GET delivery service servers: %v\n", err)
	}

	found = false
	for _, dss := range dsServers.Response {
		if *dss.DeliveryService == ds.ID && *dss.Server == server.ID {
			found = true
			break
		}
	}
	if found {
		t.Errorf("DELETE delivery service servers returned success, but still in GET")
	}
}

func getServersAndDSes(t *testing.T) ([]tc.DeliveryService, []tc.Server) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v\n", err)
	}
	if len(dses) < 1 {
		t.Fatalf("GET DeliveryServices returned no dses, must have at least 1 to test ds-servers")
	}

	servers, _, err := TOSession.GetServers()
	if err != nil {
		t.Fatalf("cannot GET Servers: %v\n", err)
	}
	if len(servers) < 1 {
		t.Fatalf("GET Servers returned no dses, must have at least 1 to test ds-servers")
	}

	return dses, servers
}
