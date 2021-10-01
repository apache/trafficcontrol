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
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDeliveryServicesRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, Topologies, ServiceCategories, DeliveryServices, DeliveryServicesRequiredCapabilities}, func() {
		GetTestDeliveryServicesRequiredCapabilitiesIMS(t)
		InvalidDeliveryServicesRequiredCapabilityAddition(t)
		GetTestDeliveryServicesRequiredCapabilities(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		GetTestDeliveryServicesRequiredCapabilitiesIMSAfterChange(t, header)
	})
}

func TestTopologyBasedDeliveryServicesRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, ServerServerCapabilitiesForTopologies, Topologies, ServiceCategories, DeliveryServices, TopologyBasedDeliveryServiceRequiredCapabilities}, func() {
		GetTestDeliveryServicesRequiredCapabilities(t)
		OriginAssignTopologyBasedDeliveryServiceWithRequiredCapabilities(t)
	})
}

func OriginAssignTopologyBasedDeliveryServiceWithRequiredCapabilities(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", "ds-top-req-cap2")
	resp, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Errorf("getting Delivery Service 'ds-top-req-cap2': %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("expected to get only one Delivery Service with XMLID 'ds-top-req-cap2', but got %d", len(resp.Response))
	}
	if resp.Response[0].ID == nil {
		t.Fatal("no ID in the resulting delivery service")
	}
	dsID := *resp.Response[0].ID
	opts.QueryParameters = url.Values{}
	alerts, _, err := TOSession.AssignServersToDeliveryService([]string{"denver-mso-org-01"}, "ds-top-req-cap2", client.RequestOptions{})
	if err != nil {
		t.Errorf("assigning server 'denver-mso-org-01' to Delivery Service 'ds-top-req-cap2': %v - alerts: %+v", err, alerts)
	}
	opts.QueryParameters.Set("dsId", strconv.Itoa(dsID))
	opts.QueryParameters.Set("type", tc.OriginTypeName)
	responseServers, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("getting servers for the 'ds-top-req-cap2' Delivery Service: %v - alerts: %+v", err, responseServers.Alerts)
	}
	if len(responseServers.Response) != 1 {
		t.Fatalf("expected just one ORG server in the response, but got %d", len(responseServers.Response))
	}
	if responseServers.Response[0].HostName == nil {
		t.Fatal("expected a valid host name for the resulting ORG server, but got nothing")
	}
	if *responseServers.Response[0].HostName != "denver-mso-org-01" {
		t.Errorf("expected host name of the resulting ORG server to be %v, but got %v", "denver-mso-org-01", *responseServers.Response[0].HostName)
	}
}

func GetTestDeliveryServicesRequiredCapabilitiesIMSAfterChange(t *testing.T, header http.Header) {
	data := testData.DeliveryServicesRequiredCapabilities
	if len(data) < 1 {
		t.Fatal("Need at least one Delivery Service Required Capability to test IMS updates to Delivery Service Required Capabilities")
	}
	if data[0].XMLID == nil {
		t.Fatal("Found a Delivery Service Required Capability in the testing data with null or undefined XMLID")
	}
	if data[0].RequiredCapability == nil {
		t.Fatal("Found a Delivery Service Required Capability in the testing data with null or undefined Required Capability")
	}
	xmlid := *data[0].XMLID
	cap := *data[0].RequiredCapability
	ds1 := helperGetDeliveryServiceID(t, &xmlid)
	if ds1 == nil {
		t.Fatalf("Failed to get ID for Delivery Service '%s'", xmlid)
	}

	testCases := []struct {
		description string
		params      url.Values
	}{
		{
			description: "get all deliveryservices required capabilities",
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by deliveryServiceID: %d", *ds1),
			params: url.Values{
				"deliveryServiceID": {strconv.Itoa(*ds1)},
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by xmlID: %s", xmlid),
			params: url.Values{
				"xmlID": {xmlid},
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by requiredCapability: %s", cap),
			params: url.Values{
				"requiredCapability": {cap},
			},
		},
	}

	opts := client.NewRequestOptions()
	opts.Header = header
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			opts.QueryParameters = tc.params
			resp, reqInf, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
			if err != nil {
				t.Errorf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
			}
			if reqInf.StatusCode != http.StatusOK {
				t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
			}
		})
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)

	opts.Header.Set(rfc.IfModifiedSince, timeStr)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			opts.QueryParameters = tc.params
			resp, reqInf, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
			if err != nil {
				t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
			}
			if reqInf.StatusCode != http.StatusNotModified {
				t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
			}
		})
	}
}

func GetTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	data := testData.DeliveryServicesRequiredCapabilities
	if len(data) < 1 {
		t.Fatal("Need at least one Delivery Service Required Capability to test IMS updates to Delivery Service Required Capabilities")
	}
	if data[0].XMLID == nil {
		t.Fatal("Found a Delivery Service Required Capability in the testing data with null or undefined XMLID")
	}
	if data[0].RequiredCapability == nil {
		t.Fatal("Found a Delivery Service Required Capability in the testing data with null or undefined Required Capability")
	}
	ds1 := helperGetDeliveryServiceID(t, data[0].XMLID)
	if ds1 == nil {
		t.Fatalf("Failed to get ID for Delivery Service '%s'", *data[0].XMLID)
	}

	testCases := []struct {
		description string
		capability  tc.DeliveryServicesRequiredCapability
		expectFunc  func(tc.DeliveryServicesRequiredCapability, []tc.DeliveryServicesRequiredCapability)
	}{
		{
			description: "get all deliveryservices required capabilities",
			expectFunc: func(expect tc.DeliveryServicesRequiredCapability, actual []tc.DeliveryServicesRequiredCapability) {
				if len(actual) != len(testData.DeliveryServicesRequiredCapabilities) {
					t.Errorf("expected length: %d, actual: %d", len(testData.DeliveryServicesRequiredCapabilities), len(actual))
				}
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by deliveryServiceID: %d", *ds1),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID: ds1,
			},
			expectFunc: func(dsRequiredCapability tc.DeliveryServicesRequiredCapability, dsReqCaps []tc.DeliveryServicesRequiredCapability) {
				for _, dsrc := range dsReqCaps {
					if dsrc.DeliveryServiceID == nil {
						t.Error("Traffic Ops returned a representation for a Delivery Service/Required Capability relationship with null or undefined Delivery Service ID")
						continue
					}
					if *dsrc.DeliveryServiceID != *dsRequiredCapability.DeliveryServiceID {
						t.Errorf("expected: all delivery service IDs to equal %d, actual: found %d", *dsRequiredCapability.DeliveryServiceID, *dsrc.DeliveryServiceID)
					}
				}
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by xmlID: %s", *data[0].XMLID),
			capability: tc.DeliveryServicesRequiredCapability{
				XMLID: data[0].XMLID,
			},
			expectFunc: func(dsRequiredCapability tc.DeliveryServicesRequiredCapability, dsReqCaps []tc.DeliveryServicesRequiredCapability) {
				for _, dsrc := range dsReqCaps {
					if dsrc.XMLID == nil {
						t.Error("Traffic Ops returned a representation for a Delivery Service/Required Capability relationship with null or undefined XMLID")
						continue
					}
					if *dsrc.XMLID != *dsRequiredCapability.XMLID {
						t.Errorf("expected: all delivery service XMLIDs to equal %s, actual: found %s", *dsRequiredCapability.XMLID, *dsrc.XMLID)
					}
				}
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				RequiredCapability: data[0].RequiredCapability,
			},
			expectFunc: func(dsRequiredCapability tc.DeliveryServicesRequiredCapability, dsReqCaps []tc.DeliveryServicesRequiredCapability) {
				for _, dsrc := range dsReqCaps {
					if dsrc.RequiredCapability == nil {
						t.Error("Traffic Ops returned a representation for a Delivery Service/Required Capability relationship with null or undefined required Capability")
						continue
					}
					if *dsrc.RequiredCapability != *dsRequiredCapability.RequiredCapability {
						t.Errorf("expected: all delivery service required capabilities to equal %s, actual: found %s", *dsRequiredCapability.RequiredCapability, *dsrc.RequiredCapability)
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			opts := client.NewRequestOptions()
			if tc.capability.XMLID != nil {
				opts.QueryParameters.Set("xmlID", *tc.capability.XMLID)
			}
			if tc.capability.RequiredCapability != nil {
				opts.QueryParameters.Set("requiredCapability", *tc.capability.RequiredCapability)
			}
			if tc.capability.DeliveryServiceID != nil {
				opts.QueryParameters.Set("deliveryServiceID", strconv.Itoa(*tc.capability.DeliveryServiceID))
			}
			capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
			if err != nil {
				t.Errorf("Unexpected error requesting Delivery Service Required Capabilities: %v - alerts: %+v", err, capabilities.Alerts)
			}
			tc.expectFunc(tc.capability, capabilities.Response)
		})
	}
}

func GetTestDeliveryServicesRequiredCapabilitiesIMS(t *testing.T) {
	data := testData.DeliveryServicesRequiredCapabilities
	if len(data) < 1 {
		t.Fatal("Need at least one Delivery Service Required Capability to test IMS updates to Delivery Service Required Capabilities")
	}
	if data[0].XMLID == nil {
		t.Fatal("Found a Delivery Service Required Capability in the testing data with null or undefined XMLID")
	}
	if data[0].RequiredCapability == nil {
		t.Fatal("Found a Delivery Service Required Capability in the testing data with null or undefined Required Capability")
	}
	xmlid := *data[0].XMLID
	cap := *data[0].RequiredCapability
	ds1 := helperGetDeliveryServiceID(t, &xmlid)
	if ds1 == nil {
		t.Fatalf("Failed to get ID for Delivery Service '%s'", xmlid)
	}

	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)

	testCases := []struct {
		description string
		params      url.Values
	}{
		{
			description: "get all deliveryservices required capabilities",
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by deliveryServiceID: %d", *ds1),
			params: url.Values{
				"deliveryServiceID": {strconv.Itoa(*ds1)},
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by xmlID: %s", xmlid),
			params: url.Values{
				"xmlID": {xmlid},
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by requiredCapability: %s", cap),
			params: url.Values{
				"requiredCapability": {cap},
			},
		},
	}

	for _, tc := range testCases {
		opts.QueryParameters = tc.params
		t.Run(tc.description, func(t *testing.T) {
			resp, reqInf, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
			if err != nil {
				t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
			}
			if reqInf.StatusCode != http.StatusNotModified {
				t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
			}
		})
	}
}

func CreateTestTopologyBasedDeliveryServicesRequiredCapabilities(t *testing.T) {
	for _, td := range testData.TopologyBasedDeliveryServicesRequiredCapabilities {

		c := tc.DeliveryServicesRequiredCapability{
			DeliveryServiceID:  helperGetDeliveryServiceID(t, td.XMLID),
			RequiredCapability: td.RequiredCapability,
		}

		alerts, _, err := TOSession.CreateDeliveryServicesRequiredCapability(c, client.RequestOptions{})
		if err != nil {
			t.Fatalf("cannot create delivery service required capability: %v - %+v", err, alerts.Alerts)
		}
	}

	invalid := tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  helperGetDeliveryServiceID(t, util.StrPtr("ds-top-req-cap")),
		RequiredCapability: util.StrPtr("asdf"),
	}
	_, reqInf, err := TOSession.CreateDeliveryServicesRequiredCapability(invalid, client.RequestOptions{})
	if err == nil {
		t.Fatal("when adding delivery service required capability to a delivery service with a topology that " +
			"doesn't have cachegroups with at least one server with the required capabilities - expected: error, actual: nil")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("when adding delivery service required capability to a delivery service with a topology that "+
			"doesn't have cachegroups with at least one server with the required capabilities - expected status code: "+
			"%d, actual: %d", http.StatusBadRequest, reqInf.StatusCode)
	}
}

func CreateTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	data := testData.DeliveryServicesRequiredCapabilities
	if len(data) == 0 {
		t.Fatal("there must be at least one test ds required capability defined")
	}
	ds1 := helperGetDeliveryServiceID(t, data[0].XMLID)
	amDS := helperGetDeliveryServiceID(t, util.StrPtr("anymap-ds"))
	testCases := []struct {
		description string
		capability  tc.DeliveryServicesRequiredCapability
	}{
		{
			description: fmt.Sprintf("re-assign a deliveryservice to a required capability; deliveryServiceID: %d, requiredCapability: %s", *ds1, *data[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  ds1,
				RequiredCapability: data[0].RequiredCapability,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with no delivery service id; deliveryServiceID: 0, requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				RequiredCapability: data[0].RequiredCapability,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with no requiredCapability; deliveryServiceID: %d, requiredCapability: 0", *ds1),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID: ds1,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with an invalid required capability; deliveryServiceID: %d, requiredCapability: bogus", *ds1),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  ds1,
				RequiredCapability: util.StrPtr("bogus"),
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with an invalid delivery service id; deliveryServiceID: -1, requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  util.IntPtr(-1),
				RequiredCapability: data[0].RequiredCapability,
			},
		},
		{
			description: "assign a deliveryservice to a required capability with an invalid deliveryservice type",
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  amDS,
				RequiredCapability: data[0].RequiredCapability,
			},
		},
	}

	// Assign all required capability to delivery services listed in `tc-fixtures.json`.
	for _, td := range testData.DeliveryServicesRequiredCapabilities {
		var dsID int
		if td.DeliveryServiceID != nil {
			dsID = *td.DeliveryServiceID
		}

		var capability string
		if td.RequiredCapability != nil {
			capability = *td.RequiredCapability
		}

		t.Run(fmt.Sprintf("assign a deliveryservice to a required capability; deliveryServiceID: %d, requiredCapability: %s", dsID, capability), func(t *testing.T) {
			cap := tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  helperGetDeliveryServiceID(t, td.XMLID),
				RequiredCapability: td.RequiredCapability,
			}

			resp, _, err := TOSession.CreateDeliveryServicesRequiredCapability(cap, client.RequestOptions{})
			if err != nil {
				t.Fatalf("Unexpected error creating a Delivery Service/Required Capability relationship: %v - alerts: %+v", err, resp.Alerts)
			}
		})
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, _, err := TOSession.CreateDeliveryServicesRequiredCapability(tc.capability, client.RequestOptions{})
			if err == nil {
				t.Fatalf("%s; expected err", tc.description)
			}
		})
	}
}

func InvalidDeliveryServicesRequiredCapabilityAddition(t *testing.T) {
	// Tests that a capability cannot be made required if the DS's services do not have it assigned

	// Get Delivery Capability for a DS
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlID", "ds1")
	capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(opts)
	if err != nil {
		t.Fatalf("cannot get Delivery Service Required Capabilities: %v - alerts: %+v", err, capabilities.Alerts)
	}
	if len(capabilities.Response) == 0 {
		t.Fatal("delivery service ds1 needs at least one capability required")
	}
	dsID := capabilities.Response[0].DeliveryServiceID
	if dsID == nil {
		t.Fatal("Traffic Ops returned a representation for a Delivery Service/Required Capability relationship with null or undefined Delivery Service ID")
	}

	// First assign current capabilities to edge server so we can assign it to the DS
	// TODO: DON'T hard-code hostnames!
	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("hostName", "atlanta-edge-01")
	resp, _, err := TOSession.GetServers(opts)
	if err != nil {
		t.Fatalf("cannot get Server by Host Name 'atlanta-edge-01': %v - alerts: %+v", err, resp.Alerts)
	}
	servers := resp.Response
	if len(servers) < 1 {
		t.Fatal("need at least one server to test invalid ds required capability assignment")
	}

	if servers[0].ID == nil {
		t.Fatal("server 'atlanta-edge-01' had nil ID")
	}

	sID := *servers[0].ID
	serverCaps := []tc.ServerServerCapability{}

	for _, cap := range capabilities.Response {
		if cap.RequiredCapability == nil {
			t.Errorf("Traffic Ops returned a representation for a Delivery Service/Required Capability relationship with null or undefined required Capability")
			continue
		}
		sCap := tc.ServerServerCapability{
			ServerID:         &sID,
			ServerCapability: cap.RequiredCapability,
		}
		resp, _, err := TOSession.CreateServerServerCapability(sCap, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not associate Capability %s with server #%d: %v - alerts: %+v", *cap.RequiredCapability, sID, err, resp.Alerts)
		}
		serverCaps = append(serverCaps, sCap)
	}

	// Assign server to ds
	alerts, _, err := TOSession.CreateDeliveryServiceServers(*dsID, []int{sID}, false, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error assigning server #%d to Delivery Service #%d: %v - alerts: %+v", sID, *dsID, err, alerts.Alerts)
	}

	// Create new bogus server capability
	scResp, _, err := TOSession.CreateServerCapability(tc.ServerCapability{
		Name: "newcap",
	}, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot create 'newcap' Server Capability: %v - alerts: %+v", err, scResp.Alerts)
	}

	// Attempt to assign to DS should fail
	_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  dsID,
		RequiredCapability: util.StrPtr("newcap"),
	}, client.RequestOptions{})
	if err == nil {
		t.Error("expected error requiring a capability that is not associated on the delivery service's servers")
	}

	// Disassociate server from DS
	setInactive(t, *dsID)
	deleteResp, _, err := TOSession.DeleteDeliveryServiceServer(*dsID, sID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("could not remove server #%d from Delivery Service #%d: %v - alerts: %+v", sID, *dsID, err, deleteResp.Alerts)
	}

	// Remove server capabilities from server
	for _, ssc := range serverCaps {
		resp, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not remove Capability '%s' from server #%d: %v - alerts: %+v", *ssc.ServerCapability, *ssc.ServerID, err, resp.Alerts)
		}
	}

	// Delete server capability
	deleteAlerts, _, err := TOSession.DeleteServerCapability("newcap", client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot delete 'newcap' Server Capability: %v - alerts: %+v", err, deleteAlerts.Alerts)
	}

}

func DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Get Required Capabilities to delete them
	capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error getting Delivery Service/Required Capability relationships: %v - alerts: %+v", err, capabilities.Alerts)
	}
	if len(capabilities.Response) < 1 {
		t.Fatal("no Delivery Service/Required Capability relationships returned")
	}
	cap := capabilities.Response[0]
	if cap.DeliveryServiceID == nil || cap.RequiredCapability == nil {
		t.Fatal("Traffic Ops returned a representation of a Delivery Service/Required Capability relationship with null or undefined required Capability and/or Delivery Service ID")
	}

	type testCase struct {
		description string
		dsID        int
		capability  string
		err         string
	}

	testCases := []testCase{
		{
			description: fmt.Sprintf("delete a deliveryservices required capability with an invalid delivery service id; deliveryServiceID: -1, requiredCapability: %s", *cap.RequiredCapability),
			dsID:        -1,
			capability:  *cap.RequiredCapability,
			err:         "no deliveryservice.RequiredCapability with that key found",
		},
		{
			description: fmt.Sprintf("delete a deliveryservices required capability with an invalid required capability; deliveryServiceID: %d, requiredCapability: bogus", *cap.DeliveryServiceID),
			dsID:        *cap.DeliveryServiceID,
			capability:  "bogus",
			err:         "no deliveryservice.RequiredCapability with that key found",
		},
	}

	for _, c := range capabilities.Response {
		if c.DeliveryServiceID == nil || c.RequiredCapability == nil {
			t.Error("Traffic Ops returned a representation of a Delivery Service/Required Capability relationship with null or undefined required Capability and/or Delivery Service ID")
			continue
		}
		t := testCase{
			description: fmt.Sprintf("delete a deliveryservices required capability; deliveryServiceID: %d, requiredCapability: %s", *c.DeliveryServiceID, *c.RequiredCapability),
			capability:  *c.RequiredCapability,
			dsID:        *c.DeliveryServiceID,
		}
		testCases = append(testCases, t)
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			alerts, _, err := TOSession.DeleteDeliveryServicesRequiredCapability(c.dsID, c.capability, client.RequestOptions{})
			if err != nil {
				if c.err != "" {
					found := false
					for _, alert := range alerts.Alerts {
						if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, c.err) {
							found = true
							continue
						}
					}
					if !found {
						t.Errorf("Expected to find an error-level alert containing the text '%s', but it was not found - alerts: %+v", c.err, alerts.Alerts)
					}
				} else {
					t.Errorf("Unexpected error deleting a relationship between a Delivery Service and a Capability it requires: %v - alerts: %+v", err, alerts.Alerts)
				}
			} else if c.err != "" {
				t.Errorf("Expected deletion to fail with reason '%s' but it succeeded", c.err)
			}
		})
	}
}

func helperGetDeliveryServiceID(t *testing.T, xmlID *string) *int {
	t.Helper()
	if xmlID == nil {
		t.Error("xml id must not be nil")
		return nil
	}
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("xmlId", *xmlID)
	ds, _, err := TOSession.GetDeliveryServices(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Delivery Services filtered by XMLID '%s': %v - alerts: %+v", *xmlID, err, ds.Alerts)
		return nil
	}
	if len(ds.Response) != 1 {
		t.Errorf("Expected exactly one Delivery Service to have XMLID '%s', found: %d", *xmlID, len(ds.Response))
		return nil
	}
	return ds.Response[0].ID
}
