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
	"strconv"
	"strings"
	"testing"
)

func TestCacheGroupsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, CacheGroupsDeliveryServices}, func() {})
}

// TODO this is the name hard-coded in the create servers test; change to be dynamic
// TODO this test assumes that a CDN named "cdn1" exists, has at least one Delivery Service, and also
// assumes that ALL SERVERS IN "cachegroup3" ARE EDGE-TIER CACHE SERVERS IN "cdn1". If that EVER changes,
// this WILL break.
const TestEdgeServerCacheGroupName = "cachegroup3"

func CreateTestCachegroupsDeliveryServices(t *testing.T) {
	dss, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Fatalf("cannot GET DeliveryServiceServers: %v", err)
	}
	if len(dss.Response) > 0 {
		t.Fatalf("cannot test cachegroups delivery services: expected no initial delivery service servers, actual %v", len(dss.Response))
	}

	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v - %v", err, dses)
	}

	clientCGs, _, err := TOSession.GetCacheGroupNullableByName(TestEdgeServerCacheGroupName)
	if err != nil {
		t.Fatalf("getting cachegroup: %v", err)
	}
	if len(clientCGs) != 1 {
		t.Fatalf("getting cachegroup expected 1, got %v", len(clientCGs))
	}

	clientCG := clientCGs[0]

	if clientCG.ID == nil {
		t.Fatalf("Cachegroup has a nil ID")
	}
	cgID := *clientCG.ID

	dsIDs := []int{}
	topologyDsIDs := []int{}
	for _, ds := range dses {
		if *ds.CDNName == "cdn1" && ds.Topology == nil {
			dsIDs = append(dsIDs, *ds.ID)
		} else if *ds.CDNName == "cdn1" && ds.Topology != nil {
			topologyDsIDs = append(topologyDsIDs, *ds.ID)
		}
	}
	if len(dsIDs) < 1 {
		t.Fatal("No Delivery Services found in CDN 'cdn1', cannot continue.")
	}

	if len(topologyDsIDs) < 1 {
		t.Fatal("No Topology-based Delivery Services found in CDN 'cdn1', cannot continue.")
	}

	_, reqInf, err := TOSession.SetCachegroupDeliveryServices(cgID, topologyDsIDs)
	if err == nil {
		t.Fatal("assigning Topology-based delivery service to cachegroup - expected: error, actual: nil")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("assigning Topology-based delivery service to cachegroup - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}

	resp, _, err := TOSession.SetCachegroupDeliveryServices(cgID, dsIDs)
	if err != nil {
		t.Fatalf("setting cachegroup delivery services returned error: %v", err)
	}
	if len(resp.Response.ServerNames) == 0 {
		t.Fatal("setting cachegroup delivery services returned success, but no servers set")
	}

	// Note this second post of the same cg-dses specifically tests a previous bug, where the query
	// failed if any servers with location parameters were already assigned, due to a foreign key
	// violation. See https://github.com/apache/trafficcontrol/pull/3199
	resp, _, err = TOSession.SetCachegroupDeliveryServices(cgID, dsIDs)
	if err != nil {
		t.Fatalf("setting cachegroup delivery services returned error: %v", err)
	}
	if len(resp.Response.ServerNames) == 0 {
		t.Fatal("setting cachegroup delivery services returned success, but no servers set")
	}

	params := url.Values{}
	for _, serverName := range resp.Response.ServerNames {
		params.Set("hostName", string(serverName))
		resp, _, err := TOSession.GetServers(&params, nil)
		if err != nil {
			t.Fatalf("getting server: %v", err)
		}
		servers := resp
		if len(servers) != 1 {
			t.Fatalf("getting servers: expected 1 got %v", len(servers))
		}
		server := servers[0]
		serverID := server.ID

		if serverID == nil {
			t.Fatalf("got a nil server ID in response, quitting")
		}
		serverDSes, _, err := TOSession.GetDeliveryServicesByServer(*serverID)

		for _, dsID := range dsIDs {
			found := false
			for _, serverDS := range serverDSes {
				if *serverDS.ID == int(dsID) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("post succeeded, but didn't assign delivery service %v to server", dsID)
			}
		}
	}
}

func setInactive(t *testing.T, dsID int) {
	ds, _, err := TOSession.GetDeliveryServiceNullableWithHdr(strconv.Itoa(dsID), nil)
	if err != nil {
		t.Errorf("Failed to fetch details for Delivery Service #%d", dsID)
		return
	}
	if ds == nil {
		t.Errorf("Got null or undefined Delivery Service for #%d", dsID)
		return
	}
	if ds.Active == nil {
		t.Errorf("Deliver Service #%d had null or undefined 'active'", dsID)
		ds.Active = new(bool)
	}
	if *ds.Active {
		*ds.Active = false
		_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(dsID, *ds, nil)
		if err != nil {
			t.Errorf("Failed to set Delivery Service #%d to inactive: %v", dsID, err)
		}
	}
}

func DeleteTestCachegroupsDeliveryServices(t *testing.T) {
	dss, _, err := TOSession.GetDeliveryServiceServersNWithHdr(1000000, nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceServers: %v", err)
	}
	for _, ds := range dss.Response {
		if ds.DeliveryService == nil {
			t.Error("Got deliveryserviceserver with no Delivery Service")
			continue
		}
		if ds.Server == nil {
			t.Error("Got deliveryserviceserver with no server")
			continue
		}

		setInactive(t, *ds.DeliveryService)

		alerts, _, err := TOSession.DeleteDeliveryServiceServer(*ds.DeliveryService, *ds.Server)
		t.Logf("Alerts from deleting deliveryserviceserver: %s", strings.Join(alerts.ToStrings(), ", "))
		if err != nil {
			t.Errorf("deleting delivery service servers: %v", err)
		}
	}

	dss, _, err = TOSession.GetDeliveryServiceServersWithHdr(nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceServers: %v", err)
	}
	if len(dss.Response) > 0 {
		t.Errorf("deleting delivery service servers: delete succeeded, expected empty subsequent get, actual %v", len(dss.Response))
	}
}
