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
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestAssignments(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Tenants, Topologies, DeliveryServices}, func() {
		AssignTestDeliveryService(t)
		AssignIncorrectTestDeliveryService(t)
		AssignTopologyBasedDeliveryService(t)
		OriginAssignTopologyBasedDeliveryService(t)
	})
}

func AssignTestDeliveryService(t *testing.T) {
	if len(testData.Servers) < 1 {
		t.Fatal("Need at least one test server to test delivery service assignment")
	}

	server := testData.Servers[0]
	if server.HostName == nil {
		t.Fatalf("First server had nil hostname: %+v", server)
	}

	params := url.Values{}
	params.Add("hostName", *server.HostName)

	rs, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v", err)
	} else if len(rs.Response) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	firstServer := rs.Response[0]
	if firstServer.ID == nil {
		t.Fatalf("Server '%s' had nil ID", *server.HostName)
	}

	rd, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(*testData.DeliveryServices[0].XMLID, nil)
	if err != nil {
		t.Fatalf("Failed to fetch DS information: %v", err)
	} else if len(rd) == 0 {
		t.Fatalf("Failed to fetch DS information: No results returned!")
	}
	firstDS := rd[0]

	if firstDS.ID == nil {
		t.Fatal("Fetch DS information returned unknown ID")
	}
	alerts, _, err := TOSession.AssignDeliveryServiceIDsToServerID(*firstServer.ID, []int{*firstDS.ID}, false)
	if err != nil {
		t.Errorf("Couldn't assign DS '%+v' to server '%+v': %v (alerts: %v)", firstDS, firstServer, err, alerts)
	}
	t.Logf("alerts: %+v", alerts)

	response, _, err := TOSession.GetServerIDDeliveryServicesWithHdr(*firstServer.ID, nil)
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

	currentTime := time.Now().UTC().Add(5 * time.Second)
	time := currentTime.Format(time.RFC1123)
	var header http.Header
	header = make(map[string][]string)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, _ := TOSession.GetServerIDDeliveryServicesWithHdr(*firstServer.ID, header)
	if reqInf.StatusCode != http.StatusNotModified {
		t.Errorf("Expected a status code of 304, got %v", reqInf.StatusCode)
	}
}

func AssignIncorrectTestDeliveryService(t *testing.T) {
	var server *tc.ServerV30
	for _, s := range testData.Servers {
		if s.CDNName != nil && *s.CDNName == "cdn2" {
			server = &s
			break
		}
	}
	if server == nil {
		t.Fatal("Couldn't find a server in CDN 'cdn2'!")
	}
	if server.HostName == nil {
		t.Fatalf("Server found with nil hostname: %+v", *server)
	}
	hostname := *server.HostName

	params := url.Values{}
	params.Add("hostName", hostname)
	rs, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v - %v", err, rs.Alerts)
	} else if len(rs.Response) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	server = &rs.Response[0]
	if server.ID == nil {
		t.Fatalf("Server '%s' has nil ID", hostname)
	}

	rd, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(*testData.DeliveryServices[0].XMLID, nil)
	if err != nil {
		t.Fatalf("Failed to fetch DS information: %v", err)
	} else if len(rd) == 0 {
		t.Fatalf("Failed to fetch DS information: No results returned!")
	}
	firstDS := rd[0]

	if firstDS.ID == nil {
		t.Fatal("Fetch DS information returned unknown ID")
	}
	alerts, _, err := TOSession.AssignDeliveryServiceIDsToServerID(*server.ID, []int{*firstDS.ID}, false)
	if err == nil {
		t.Errorf("Expected bad assignment to fail, but it didn't! (alerts: %v)", alerts)
	}

	response, _, err := TOSession.GetServerIDDeliveryServicesWithHdr(*server.ID, nil)
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

func AssignTopologyBasedDeliveryService(t *testing.T) {
	var server *tc.ServerV30
	for _, s := range testData.Servers {
		if s.CDNName != nil && *s.CDNName == "cdn1" && s.Type == string(tc.CacheTypeEdge) {
			server = &s
			break
		}
	}
	if server == nil || server.HostName == nil {
		t.Fatalf("Couldn't find an EDGE server in CDN 'cdn1'!")
	}

	params := url.Values{}
	params.Add("hostName", *server.HostName)
	rs, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v", err)
	} else if len(rs.Response) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	server = &rs.Response[0]
	if server.ID == nil {
		t.Fatal("Server had nil ID")
	}

	rd, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr("ds-top", nil)
	if err != nil {
		t.Fatalf("Failed to fetch DS information: %v", err)
	} else if len(rd) == 0 {
		t.Fatalf("Failed to fetch DS information: No results returned!")
	}
	firstDS := rd[0]

	if firstDS.ID == nil {
		t.Fatal("Fetch DS information returned unknown ID")
	}
	alerts, reqInf, err := TOSession.AssignDeliveryServiceIDsToServerID(*server.ID, []int{*firstDS.ID}, false)
	if err != nil {
		t.Errorf("Expected assignment to succeed, but it didn't! (alerts: %v)", alerts)
	}
	if reqInf.StatusCode >= http.StatusBadRequest {
		t.Fatalf("assigning Topology-based delivery service to server - expected: non-error status code, actual: %d", reqInf.StatusCode)
	}

	response, _, err := TOSession.GetServerIDDeliveryServicesWithHdr(*server.ID, nil)
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

	if !found {
		t.Errorf(`Valid Server/DS assignment was not created!`)
	}
}

func OriginAssignTopologyBasedDeliveryService(t *testing.T) {
	params := url.Values{}
	params.Add("hostName", "denver-mso-org-01")
	rs, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Failed to fetch server information: %v", err)
	} else if len(rs.Response) == 0 {
		t.Fatalf("Failed to fetch server information: No results returned!")
	}
	origin := &rs.Response[0]
	if origin.ID == nil {
		t.Fatal("Server had nil ID")
	}

	rd, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr("ds-top-req-cap", nil)
	if err != nil {
		t.Fatalf("Failed to fetch DS information: %v", err)
	} else if len(rd) == 0 {
		t.Fatalf("Failed to fetch DS information: No results returned!")
	}
	firstDS := rd[0]
	if firstDS.ID == nil {
		t.Fatal("Fetch DS information returned unknown ID")
	}

	// invalid assignment: ORG server cachegroup does not belong to the topology
	alerts, reqInf, err := TOSession.AssignDeliveryServiceIDsToServerID(*origin.ID, []int{*firstDS.ID}, true)
	if err == nil {
		t.Errorf("Expected assigning ORG server to topology-based delivery service where the ORG server does not belong to the topology to fail, but it didn't! (alerts: %v)", alerts)
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("assigning Topology-based delivery service to ORG server that does not belong to the topology - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}

	// valid assignment ORG server cachegroup belongs to the topology
	rd, _, err = TOSession.GetDeliveryServiceByXMLIDNullableWithHdr("ds-top", nil)
	if err != nil {
		t.Fatalf("Failed to fetch DS information: %v", err)
	} else if len(rd) == 0 {
		t.Fatalf("Failed to fetch DS information: No results returned!")
	}
	firstDS = rd[0]
	if firstDS.ID == nil {
		t.Fatal("Fetch DS information returned unknown ID")
	}

	alerts, reqInf, err = TOSession.AssignDeliveryServiceIDsToServerID(*origin.ID, []int{*firstDS.ID}, true)
	if err != nil {
		t.Errorf("Expected assigning ORG server to topology-based delivery service where the ORG server belongs to the topology to succeed, but it didn't! (alerts: %v, err: %v)", alerts, err)
	}
	if reqInf.StatusCode < http.StatusOK || reqInf.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("assigning Topology-based delivery service to ORG server that belongs to the topology - expected: 200-level status code, actual: %d", reqInf.StatusCode)
	}

	response, _, err := TOSession.GetServerIDDeliveryServicesWithHdr(*origin.ID, nil)
	if err != nil {
		t.Fatalf("Couldn't get Delivery Services assigned to Server '%+v': %v", *origin, err)
	}
	var found bool
	for _, ds := range response {

		if ds.ID != nil && *ds.ID == *firstDS.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf(`Valid Server/DS assignment was not created!`)
	}
}
