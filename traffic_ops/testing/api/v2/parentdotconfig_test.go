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
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestParentDotConfig(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices}, func() {
		defer DeleteTestDeliveryServiceServersCreated(t)
		CreateTestDeliveryServiceServers(t)
		GetTestParentDotConfig(t)
	})
}

func GetTestParentDotConfig(t *testing.T) {
	dsServers, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Fatalf("GET delivery service servers: %v", err)
	}
	if len(dsServers.Response) == 0 {
		t.Fatal("GET delivery service servers: no servers found")
	}

	dss := dsServers.Response[0]

	if dss.Server == nil {
		t.Fatal("GET delivery service servers: returned nil server")
	}
	if dss.DeliveryService == nil {
		t.Fatal("GET delivery service servers: returned nil ds")
	}

	ds, _, err := TOSession.GetDeliveryService(strconv.Itoa(*dss.DeliveryService))
	if err != nil {
		t.Fatalf("Getting ds %+v: %v", *dss.DeliveryService, err)
	}
	if ds == nil {
		t.Fatalf("Getting ds %+v: got nil response", *dss.DeliveryService)
	}
	if ds.OrgServerFQDN == "" {
		t.Fatalf("Getting ds %+v: got empty ds.OrgServerFQDN", *dss.DeliveryService)
	}

	originURI, err := url.Parse(ds.OrgServerFQDN)
	if err != nil {
		t.Fatalf("Getting ds %+v: ds.OrgServerFQDN '%v' failed to parse as a URL: %v", *dss.DeliveryService, ds.OrgServerFQDN, err)
	}
	originHost := originURI.Hostname()

	parentDotConfig, _, err := TOSession.GetATSServerConfig(*dss.Server, "parent.config")
	if err != nil {
		t.Fatalf("Getting server %v config parent.config: %v", *dss.Server, err)
	}

	if !strings.Contains(parentDotConfig, originHost) {
		t.Errorf("expected: parent.config to contain delivery service origin FQDN '%v' host '%v', actual:\n'''\n%+v\n'''", ds.OrgServerFQDN, originHost, parentDotConfig)
	}
}

func CreateTestDeliveryServiceServers(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) < 1 {
		t.Error("GET DeliveryServices returned no dses, must have at least 1 to test ds-servers")
	}

	servers, _, err := TOSession.GetServers()
	if err != nil {
		t.Errorf("cannot GET Servers: %v", err)
	}
	if len(servers) < 1 {
		t.Error("GET Servers returned no servers, must have at least 1 to test ds-servers")
	}

	for _, ds := range dses {
		serverIDs := make([]int, 0, len(servers))
		for _, server := range servers {
			if server.Type == "EDGE" && server.CDNName == ds.CDNName {
				serverIDs = append(serverIDs, server.ID)
			}
		}

		if len(serverIDs) > 0 {
			_, err = TOSession.CreateDeliveryServiceServers(ds.ID, serverIDs, true)
			if err != nil {
				t.Errorf("POST delivery service servers: %v", err)
			}
		}
	}
}

// DeleteTestDeliveryServiceServersCreated deletes the dss assignments created by CreateTestDeliveryServiceServers.
func DeleteTestDeliveryServiceServersCreated(t *testing.T) {
	// You gotta do this because TOSession.GetDeliveryServiceServers doesn't fetch the complete response.......
	dssLen := len(testData.Servers) * len(testData.DeliveryServices)
	dsServers, _, err := TOSession.GetDeliveryServiceServersN(dssLen)
	if err != nil {
		t.Fatalf("GET delivery service servers: %v", err)
	}

	for _, dss := range dsServers.Response {
		if dss.DeliveryService == nil {
			t.Error("Found ds-to-server assignment with nil Delivery Service")
			continue
		}
		if dss.Server == nil {
			t.Error("Found ds-to-server assignment with nil Server")
			continue
		}

		_, _, err := TOSession.DeleteDeliveryServiceServer(*dss.DeliveryService, *dss.Server)
		if err != nil {
			t.Errorf("Failed to remove assignment of server #%d to DS #%d: %v", *dss.Server, *dss.DeliveryService, err)
		}
	}

	dsServers, _, err = TOSession.GetDeliveryServiceServersN(dssLen)
	if err != nil {
		t.Fatalf("GET delivery service servers: %v", err)
	}

	for _, dss := range dsServers.Response {
		if dss.DeliveryService == nil {
			t.Error("Found ds-to-server assignment (after supposed deletion) with nil DeliveryService")
			continue
		}
		if dss.Server == nil {
			t.Error("Found ds-to-server assignment (after supposed deletion) with nil Server")
			continue
		}

		t.Errorf("Found ds-to-server assignment {DSID: %d, Server: %d} after deletion", *dss.DeliveryService, *dss.Server)
	}
}
