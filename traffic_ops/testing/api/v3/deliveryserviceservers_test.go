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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
)

func TestDeliveryServiceServers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices}, func() {
		DeleteTestDeliveryServiceServers(t)
		AssignServersToTopologyBasedDeliveryService(t)
		AssignOriginsToTopologyBasedDeliveryServices(t)
		TryToRemoveLastServerInDeliveryService(t)
		AssignServersToNonTopologyBasedDeliveryServiceThatUsesMidTier(t)
		GetTestDSSIMS(t)
	})
}

func TestDeliveryServiceServersWithRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, Topologies, DeliveryServices, DeliveryServicesRequiredCapabilities, ServerServerCapabilities}, func() {
		CreateTestDeliveryServiceServersWithRequiredCapabilities(t)
		CreateTestMSODSServerWithReqCap(t)
	})
}

const dssaTestingXMLID = "test-ds-server-assignments"

func TryToRemoveLastServerInDeliveryService(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(dssaTestingXMLID, nil)
	if err != nil {
		t.Fatalf("Unexpected error trying to get Delivery service with XMLID '%s': %v", dssaTestingXMLID, err)
	}
	if len(dses) != 1 {
		t.Fatalf("Expected exactly one Delivery service with XMLID '%s', got: %d", dssaTestingXMLID, len(dses))
	}
	ds := dses[0]
	if ds.ID == nil {
		t.Fatalf("Delivery Service '%s' has no ID", dssaTestingXMLID)
	}

	statuses, _, err := TOSession.GetStatusesWithHdr(nil)
	if err != nil {
		t.Fatalf("Could not fetch Statuses: %v", err)
	}
	if len(statuses) < 1 {
		t.Fatal("Need at least one Status")
	}

	var badStatusID int
	found := false
	for _, status := range statuses {
		if status.Name != tc.CacheStatusOnline.String() && status.Name != tc.CacheStatusReported.String() {
			badStatusID = status.ID
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Need at least one status with a name other than '%s' or '%s'", tc.CacheStatusOnline, tc.CacheStatusReported)
	}

	// TODO: this isn't sufficient to define a single server, so there might be
	// a better way to do this.
	params := url.Values{}
	params.Set("hostName", dssaTestingXMLID)
	params.Set("domainName", dssaTestingXMLID)
	servers, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Unexpected error fetching server '%s.%s': %v", dssaTestingXMLID, dssaTestingXMLID, err)
	}
	if len(servers.Response) != 1 {
		t.Fatalf("Expected exactly one server with FQDN '%s.%s', got: %d", dssaTestingXMLID, dssaTestingXMLID, len(servers.Response))
	}
	server := servers.Response[0]
	if server.ID == nil {
		t.Fatal("Server had null/undefined ID after creation")
	}

	_, _, err = TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*server.ID}, true)
	if err != nil {
		t.Fatalf("Failed to assign server to Delivery Service: %v", err)
	}

	_, _, err = TOSession.CreateDeliveryServiceServers(*ds.ID, []int{}, true)
	if err == nil {
		t.Error("Didn't get expected error trying to remove the only server assigned to a Delivery Service")
	} else {
		t.Logf("Got expected error trying to remove the only server assigned to a Delivery Service: %v", err)
	}

	_, _, err = TOSession.DeleteDeliveryServiceServer(*ds.ID, *server.ID)
	if err == nil {
		t.Error("Didn't get expected error trying to remove the only server assigned to a Delivery Service")
	} else {
		t.Logf("Got expected error trying to remove the only server assigned to a Delivery Service: %v", err)
	}

	alerts, _, err := TOSession.DeleteServerByID(*server.ID)
	t.Logf("Alerts from deleting server: %s", strings.Join(alerts.ToStrings(), ", "))
	if err == nil {
		t.Error("Didn't get expected error trying to delete the only server assigned to a Delivery Service")
	} else {
		t.Logf("Got expected error trying to delete the only server assigned to a Delivery Service: %v", err)
	}

	alerts, _, err = TOSession.AssignDeliveryServiceIDsToServerID(*server.ID, []int{}, true)
	t.Logf("Alerts from removing Delivery Service from server: %s", strings.Join(alerts.ToStrings(), ", "))
	if err == nil {
		t.Error("Didn't get expected error trying to remove a Delivery Service from the only server to which it is assigned")
	} else {
		t.Logf("Got expected error trying to remove a Delivery Service from the only server to which it is assigned: %v", err)
	}

	server.StatusID = &badStatusID
	putRequest := tc.ServerPutStatus{
		Status:        util.JSONNameOrIDStr{ID: &badStatusID},
		OfflineReason: util.StrPtr("test"),
	}
	alertsPtr, _, err := TOSession.UpdateServerStatus(*server.ID, putRequest)
	if alertsPtr != nil {
		t.Logf("Alerts from updating server status: %s", strings.Join(alertsPtr.ToStrings(), ", "))
	}
	if err == nil {
		t.Error("Didn't get expected error trying to put server into a bad state when it's the only one assigned to a Delivery Service")
	} else {
		t.Logf("Got expected error trying to put server into a bad state when it's the only one assigned to a Delivery Service: %v", err)
	}

	alerts, _, err = TOSession.UpdateServerByIDWithHdr(*server.ID, server, nil)
	t.Logf("Alerts from updating server status: %s", strings.Join(alerts.ToStrings(), ", "))
	if err == nil {
		t.Error("Didn't get expected error trying to put server into a bad state when it's the only one assigned to a Delivery Service")
	} else {
		t.Logf("Got expected error trying to put server into a bad state when it's the only one assigned to a Delivery Service: %v", err)
	}

	server.HostName = util.StrPtr(dssaTestingXMLID + "-quest")
	server.ID = nil
	if len(server.Interfaces) == 0 {
		t.Fatal("no interfaces in this server, quitting")
	}
	interfaces := make([]tc.ServerInterfaceInfo, 0)
	ipAddresses := make([]tc.ServerIPAddress, 0)
	gateway := "1.2.3.4"
	ipAddresses = append(ipAddresses, tc.ServerIPAddress{
		Address:        "1.1.1.1",
		Gateway:        &gateway,
		ServiceAddress: true,
	})
	interfaces = append(interfaces, tc.ServerInterfaceInfo{
		IPAddresses:  ipAddresses,
		MaxBandwidth: server.Interfaces[0].MaxBandwidth,
		Monitor:      false,
		MTU:          server.Interfaces[0].MTU,
		Name:         server.Interfaces[0].Name,
	})
	server.Interfaces = interfaces
	alerts, _, err = TOSession.CreateServerWithHdr(server, nil)
	if err != nil {
		t.Fatalf("Failed to create server: %v - alerts: %s", err, strings.Join(alerts.ToStrings(), ", "))
	}
	params.Set("hostName", *server.HostName)
	servers, _, err = TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("Could not fetch server after creation: %v", err)
	}
	if len(servers.Response) != 1 {
		t.Fatalf("Expected exactly 1 server with hostname '%s'; got: %d", *server.HostName, len(servers.Response))
	}
	server = servers.Response[0]
	if server.ID == nil {
		t.Fatal("Server had null/undefined ID after creation")
	}

	_, _, err = TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*server.ID}, true)
	if err == nil {
		t.Error("Didn't get expected error trying to replace the last server assigned to a Delivery Service with a server in a bad state")
	} else {
		t.Logf("Got expected error trying to replace the last server assigned to a Delivery Service with a server in a bad state: %v", err)
	}

	// Cleanup
	setInactive(t, *ds.ID)
	alerts, _, err = TOSession.DeleteServerByID(*server.ID)
	t.Logf("Alerts from deleting a server: %s", strings.Join(alerts.ToStrings(), ", "))
	if err != nil {
		t.Errorf("Failed to delete server: %v", err)
	}
}

func AssignServersToTopologyBasedDeliveryService(t *testing.T) {
	params := url.Values{}
	params.Set("xmlId", "ds-top")
	ds, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET delivery service 'ds-top': %s", err.Error())
	}
	if len(ds) != 1 {
		t.Fatalf("expected one delivery service: 'ds-top', actual: %v", len(ds))
	}
	if ds[0].Topology == nil {
		t.Fatal("expected delivery service: 'ds-top' to have a non-nil Topology, actual: nil")
	}
	serversResp, _, err := TOSession.GetServersWithHdr(nil, nil)
	servers := []tc.ServerV30{}
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
		if s.CDNID != nil && *s.CDNID == *ds[0].CDNID && s.Type == tc.CacheTypeEdge.String() {
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

func AssignOriginsToTopologyBasedDeliveryServices(t *testing.T) {
	// attempt to assign ORG server to a topology-based DS while the ORG server's cachegroup doesn't belong to the topology
	params := url.Values{}
	params.Add("hostName", "denver-mso-org-01")
	resp, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("unable to GET server: %v", err)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("GET server expected length: 1, actual: %d", len(resp.Response))
	}
	orgServer := resp.Response[0]
	_, reqInf, err := TOSession.AssignServersToDeliveryService([]string{*orgServer.HostName}, "ds-top-req-cap")
	if err == nil {
		t.Fatal("assigning ORG server to topology-based delivery service while the ORG server's cachegroup does not belong to the topology - expected: error, actual: nil error")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("assigning ORG server to topology-based delivery service while the ORG server's cachegroup does not belong to the topology - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}
	params = url.Values{}
	params.Set("xmlId", "ds-top-req-cap")
	ds, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET delivery service 'ds-top-req-cap': %s", err.Error())
	}
	if len(ds) != 1 {
		t.Fatalf("expected one delivery service: 'ds-top-req-cap', actual: %v", len(ds))
	}
	if ds[0].Topology == nil {
		t.Fatal("expected delivery service: 'ds-top-req-cap' to have a non-nil Topology, actual: nil")
	}
	_, reqInf, err = TOSession.CreateDeliveryServiceServers(*ds[0].ID, []int{*orgServer.ID}, false)
	if err == nil {
		t.Fatal("creating deliveryserviceserver assignment for ORG server to topology-based delivery service while the ORG server's cachegroup does not belong to the topology - expected: error, actual: nil error")
	}
	if reqInf.StatusCode < http.StatusBadRequest || reqInf.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("creating deliveryserviceserver assignment for ORG server to topology-based delivery service while the ORG server's cachegroup does not belong to the topology - expected: 400-level status code, actual: %d", reqInf.StatusCode)
	}

	// attempt to assign ORG server to a topology-based DS while the ORG server's cachegroup belongs to the topology
	_, reqInf, err = TOSession.AssignServersToDeliveryService([]string{*orgServer.HostName}, "ds-top")
	if err != nil {
		t.Fatalf("assigning ORG server to topology-based delivery service while the ORG server's cachegroup belongs to the topology - expected: no error, actual: %v", err)
	}
	if reqInf.StatusCode < http.StatusOK || reqInf.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("assigning ORG server to topology-based delivery service while the ORG server's cachegroup belongs to the topology - expected: 200-level status code, actual: %d", reqInf.StatusCode)
	}
	params = url.Values{}
	params.Set("xmlId", "ds-top")
	ds, _, err = TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET delivery service 'ds-top': %s", err.Error())
	}
	if len(ds) != 1 {
		t.Fatalf("expected one delivery service: 'ds-top', actual: %v", len(ds))
	}
	if ds[0].Topology == nil {
		t.Fatal("expected delivery service: 'ds-top' to have a non-nil Topology, actual: nil")
	}
	_, reqInf, err = TOSession.CreateDeliveryServiceServers(*ds[0].ID, []int{*orgServer.ID}, true)
	if err != nil {
		t.Fatalf("creating deliveryserviceserver assignment for ORG server to topology-based delivery service while the ORG server's cachegroup belongs to the topology - expected: no error, actual: %v", err)
	}
	if reqInf.StatusCode < http.StatusOK || reqInf.StatusCode >= http.StatusMultipleChoices {
		t.Fatalf("creating deliveryserviceserver assignment for ORG server to topology-based delivery service while the ORG server's cachegroup belongs to the topology - expected: 200-level status code, actual: %d", reqInf.StatusCode)
	}
}

func AssignServersToNonTopologyBasedDeliveryServiceThatUsesMidTier(t *testing.T) {
	params := url.Values{}
	params.Set("xmlId", "ds1")
	dsWithMid, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
	if err != nil {
		t.Fatalf("cannot GET delivery service 'ds1': %s", err.Error())
	}
	if len(dsWithMid) != 1 {
		t.Fatalf("expected one delivery service: 'ds1', actual: %v", len(dsWithMid))
	}
	if dsWithMid[0].Topology != nil {
		t.Fatal("expected delivery service: 'ds1' to have a nil Topology, actual: non-nil")
	}
	serversResp, _, err := TOSession.GetServersWithHdr(nil, nil)
	if err != nil {
		t.Fatalf("unable to fetch all servers: %v", err)
	}
	serversIds := []int{}
	for _, s := range serversResp.Response {
		if s.CDNID != nil && *s.CDNID == *dsWithMid[0].CDNID && s.Type == tc.CacheTypeEdge.String() {
			serversIds = append(serversIds, *s.ID)
		}
	}
	if len(serversIds) < 1 {
		t.Fatalf("expected: at least one EDGE in cdn %s, actual: 0", *dsWithMid[0].CDNName)
	}

	_, _, err = TOSession.CreateDeliveryServiceServers(*dsWithMid[0].ID, serversIds, true)
	if err != nil {
		t.Fatalf("unable to create delivery service server associations: %v", err)
	}

	params = url.Values{"dsId": []string{strconv.Itoa(*dsWithMid[0].ID)}}
	dsServersResp, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("unable to fetch delivery service servers: %v", err)
	}
	dsServerIds := []int{}
	for _, dss := range dsServersResp.Response {
		dsServerIds = append(dsServerIds, *dss.ID)
	}
	if len(dsServerIds) <= len(serversIds) {
		t.Fatalf("delivery service servers (%d) expected to exceed directly assigned servers (%d) to account for implicitly assigned mid servers", len(dsServerIds), len(serversIds))
	}

	for _, dss := range dsServersResp.Response {
		if dss.CDNID != nil && *dss.CDNID != *dsWithMid[0].CDNID {
			t.Fatalf("a server for another cdn was returned for this delivery service")
		}
	}
}

func GetTestDSSIMS(t *testing.T) {
	const noLimit = 999999
	_, reqInf, err := TOSession.GetDeliveryServiceServersWithLimitsWithHdr(noLimit, nil, nil, nil)
	if err != nil {
		t.Errorf("deliveryserviceservers expected: no error, actual: %v", err)
	} else if reqInf.StatusCode != http.StatusOK {
		t.Errorf("expected deliveryserviceservers response code %v, actual %v", http.StatusOK, reqInf.StatusCode)
	}

	reqHdr := http.Header{}
	reqHdr.Set(rfc.IfModifiedSince, time.Now().UTC().Add(2*time.Second).Format(time.RFC1123))

	_, reqInf, err = TOSession.GetDeliveryServiceServersWithLimitsWithHdr(noLimit, nil, nil, reqHdr)
	if err != nil {
		t.Errorf("deliveryserviceservers IMS request expected: no error, actual: %v", err)
	} else if reqInf.StatusCode != http.StatusNotModified {
		t.Errorf("expected deliveryserviceservers IMS response code %v, actual %v", http.StatusNotModified, reqInf.StatusCode)
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
			params := url.Values{}
			params.Add("hostName", ctc.serverName)
			resp, _, err := TOSession.GetServersWithHdr(&params, nil)
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
	dsReqCap, _, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(nil, util.StrPtr("msods1"), nil, nil)
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
	resp, _, err := TOSession.GetServersWithHdr(&params, nil)
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
	sccs, _, err := TOSession.GetServerServerCapabilitiesWithHdr(s.ID, nil, nil, nil)
	if err != nil {
		t.Fatalf("GET server server capabilities for denver-mso-org-01: %v", err)
	}
	if len(sccs) != 0 {
		t.Fatal("expected 0 server server capabilities for server denver-mso-org-01")
	}

	// Is origin included in eligible servers even though it doesnt have required capability
	eServers, _, err := TOSession.GetDeliveryServicesEligibleWithHdr(*dsReqCap[0].DeliveryServiceID, nil)
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
	if ds.Active == nil {
		t.Fatalf("Got a Delivery Service with nil 'Active': %+v", ds)
	}

	_, _, err := TOSession.CreateDeliveryServiceServers(*ds.ID, []int{*server.ID}, true)
	if err != nil {
		t.Errorf("POST delivery service servers: %v", err)
	}

	dsServers, _, err := TOSession.GetDeliveryServiceServersWithHdr(nil)
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

	if *ds.Active {
		*ds.Active = false
		_, _, err = TOSession.UpdateDeliveryServiceV30WithHdr(*ds.ID, ds, nil)
		if err != nil {
			t.Errorf("Setting Delivery Service #%d to inactive", *ds.ID)
		}
	}

	if _, _, err := TOSession.DeleteDeliveryServiceServer(*ds.ID, *server.ID); err != nil {
		t.Errorf("DELETE delivery service server: %v", err)
	}

	dsServers, _, err = TOSession.GetDeliveryServiceServersWithHdr(nil)
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

func getServerAndDSofSameCDN(t *testing.T) (tc.DeliveryServiceNullableV30, tc.ServerV30) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Fatalf("cannot GET DeliveryServices: %v", err)
	}
	if len(dses) < 1 {
		t.Fatal("GET DeliveryServices returned no dses, must have at least 1 to test ds-servers")
	}

	resp, _, err := TOSession.GetServersWithHdr(nil, nil)
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

	return tc.DeliveryServiceNullableV30{}, tc.ServerV30{}
}

func CreateTestDeliveryServiceServerAssignments(t *testing.T) {
	for _, dss := range testData.DeliveryServiceServerAssignments {
		resp, _, err := TOSession.AssignServersToDeliveryService(dss.ServerNames, dss.XmlId)
		assert.NoError(t, err, "Could not create Delivery Service Server Assignments: %v - alerts: %+v", err, resp.Alerts)
	}
}
