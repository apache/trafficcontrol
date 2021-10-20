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
)

func TestAssignments(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Tenants, DeliveryServices}, func() {
		AssignTestDeliveryService(t)
		AssignIncorrectTestDeliveryService(t)
	})
}

func AssignTestDeliveryService(t *testing.T) {
	rs, _, err := TOSession.GetServerByHostName(testData.Servers[0].HostName)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v", err)
	} else if len(rs) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	firstServer := rs[0]

	rd, _, err := TOSession.GetDeliveryServiceByXMLIDNullable(*testData.DeliveryServices[0].XMLID)
	if err != nil {
		t.Fatalf("Failed to fetch DS information: %v", err)
	} else if len(rd) == 0 {
		t.Fatalf("Failed to fetch DS information: No results returned!")
	}
	firstDS := rd[0]

	if firstDS.ID == nil {
		t.Fatal("Fetch DS information returned unknown ID")
	}
	alerts, _, err := TOSession.AssignDeliveryServiceIDsToServerID(firstServer.ID, []int{*firstDS.ID}, true)
	if err != nil {
		t.Errorf("Couldn't assign DS '%+v' to server '%+v': %v (alerts: %v)", firstDS, firstServer, err, alerts)
	}
	t.Logf("alerts: %+v", alerts)

	response, _, err := TOSession.GetServerIDDeliveryServices(firstServer.ID)
	t.Logf("response: %+v", response)
	if err != nil {
		t.Fatalf("Couldn't get Delivery Services assigned to Server '%+v': %v", firstServer, err)
	}

	var found bool
	for _, ds := range response {
		if ds.ID != nil && *ds.ID == *firstDS.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf(`Server/DS assignment not found after "successful" assignment!`)
	}
}

func AssignIncorrectTestDeliveryService(t *testing.T) {
	var server *tc.Server
	for _, s := range testData.Servers {
		if s.CDNName == "cdn2" {
			server = &s
			break
		}
	}
	if server == nil {
		t.Fatalf("Couldn't find a server in CDN 'cdn2'!")
	}

	rs, _, err := TOSession.GetServerByHostName(server.HostName)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v", err)
	} else if len(rs) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	server = &rs[0]

	rd, _, err := TOSession.GetDeliveryServiceByXMLIDNullable(*testData.DeliveryServices[0].XMLID)
	if err != nil {
		t.Fatalf("Failed to fetch DS information: %v", err)
	} else if len(rd) == 0 {
		t.Fatalf("Failed to fetch DS information: No results returned!")
	}
	firstDS := rd[0]

	if firstDS.ID == nil {
		t.Fatal("Fetch DS information returned unknown ID")
	}
	alerts, _, err := TOSession.AssignDeliveryServiceIDsToServerID(server.ID, []int{*firstDS.ID}, false)
	if err == nil {
		t.Errorf("Expected bad assignment to fail, but it didn't! (alerts: %v)", alerts)
	}

	response, _, err := TOSession.GetServerIDDeliveryServices(server.ID)
	t.Logf("response: %+v", response)
	if err != nil {
		t.Fatalf("Couldn't get Delivery Services assigned to Server '%+v': %v", *server, err)
	}

	var found bool
	for _, ds := range response {

		if ds.ID != nil && *ds.ID == *firstDS.ID {
			found = true
			break
		}
	}

	if found {
		t.Errorf(`Invalid Server/DS assignment was created!`)
	}
}
