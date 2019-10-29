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
)

func TestCacheGroupsDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, DeliveryServices, CacheGroupsDeliveryServices}, func() {})
}

const TestEdgeServerCacheGroupName = "cachegroup1" // TODO this is the name hard-coded in the create servers test; change to be dynamic

func CreateTestCachegroupsDeliveryServices(t *testing.T) {
	dss, _, err := TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceServers: %v\n", err)
	}
	if len(dss.Response) > 0 {
		t.Errorf("cannot test cachegroups delivery services: expected no initial delivery service servers, actual %v\n", len(dss.Response))
	}

	dses, _, err := TOSession.GetDeliveryServices()
	if err != nil {
		t.Errorf("cannot GET DeliveryServices: %v - %v\n", err, dses)
	}

	clientCGs, _, err := TOSession.GetCacheGroupByName(TestEdgeServerCacheGroupName)
	if err != nil {
		t.Errorf("getting cachegroup: " + err.Error() + "\n")
	}
	if len(clientCGs) != 1 {
		t.Errorf("getting cachegroup expected 1, got %v\n", len(clientCGs))
	}
	clientCG := clientCGs[0]

	cgID := clientCG.ID

	dsIDs := []int64{}
	for _, ds := range dses {
		dsIDs = append(dsIDs, int64(ds.ID))
	}

	resp, _, err := TOSession.SetCachegroupDeliveryServices(cgID, dsIDs)
	if err != nil {
		t.Errorf("setting cachegroup delivery services returned error: %v\n", err)
	}
	if len(resp.Response.ServerNames) == 0 {
		t.Errorf("setting cachegroup delivery services returned success, but no servers set\n")
	}

	// Note this second post of the same cg-dses specifically tests a previous bug, where the query failed if any servers with location parameters were already assigned, due to an fk violation. See https://github.com/apache/trafficcontrol/pull/3199
	resp, _, err = TOSession.SetCachegroupDeliveryServices(cgID, dsIDs)
	if err != nil {
		t.Errorf("setting cachegroup delivery services returned error: %v\n", err)
	}
	if len(resp.Response.ServerNames) == 0 {
		t.Errorf("setting cachegroup delivery services returned success, but no servers set\n")
	}

	for _, serverName := range resp.Response.ServerNames {
		servers, _, err := TOSession.GetServerByHostName(string(serverName))
		if err != nil {
			t.Errorf("getting server: " + err.Error())
		}
		if len(servers) != 1 {
			t.Errorf("getting servers: expected 1 got %v\n", len(servers))
		}
		server := servers[0]
		serverID := server.ID

		serverDSes, _, err := TOSession.GetDeliveryServicesByServer(serverID)

		for _, dsID := range dsIDs {
			found := false
			for _, serverDS := range serverDSes {
				if serverDS.ID == int(dsID) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("post succeeded, but didn't assign delivery service %v to server\n", dsID)
			}
		}
	}
}

func DeleteTestCachegroupsDeliveryServices(t *testing.T) {
	dss, _, err := TOSession.GetDeliveryServiceServersN(1000000)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceServers: %v\n", err)
	}
	for _, ds := range dss.Response {
		_, _, err := TOSession.DeleteDeliveryServiceServer(*ds.DeliveryService, *ds.Server)
		if err != nil {
			t.Errorf("deleting delivery service servers: " + err.Error() + "\n")
		}
	}

	dss, _, err = TOSession.GetDeliveryServiceServers()
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceServers: %v\n", err)
	}
	if len(dss.Response) > 0 {
		t.Errorf("deleting delivery service servers: delete succeeded, expected empty subsequent get, actual %v\n", len(dss.Response))
	}
}
