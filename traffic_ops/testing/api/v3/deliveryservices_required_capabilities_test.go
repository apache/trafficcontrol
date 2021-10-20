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
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
)

func TestDeliveryServicesRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, Topologies, DeliveryServices, DeliveryServicesRequiredCapabilities}, func() {
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
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, ServerServerCapabilitiesForTopologies, Topologies, DeliveryServices, TopologyBasedDeliveryServiceRequiredCapabilities}, func() {
		GetTestDeliveryServicesRequiredCapabilities(t)
		OriginAssignTopologyBasedDeliveryServiceWithRequiredCapabilities(t)
	})
}

func OriginAssignTopologyBasedDeliveryServiceWithRequiredCapabilities(t *testing.T) {
	resp, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr("ds-top-req-cap2", nil)
	if err != nil {
		t.Errorf("getting delivery service by xml ID: %v", err.Error())
	}
	if len(resp) != 1 {
		t.Fatalf("expected to get only one delivery service in the response, but got %d", len(resp))
	}
	if resp[0].ID == nil {
		t.Fatalf("no ID in the resulting delivery service")
	}
	dsID := *resp[0].ID
	params := url.Values{}
	_, _, err = TOSession.AssignServersToDeliveryService([]string{"denver-mso-org-01"}, "ds-top-req-cap2")
	if err != nil {
		t.Errorf("assigning ORG server to ds-top delivery service: %v", err.Error())
	}
	params.Add("dsId", strconv.Itoa(dsID))
	params.Add("type", tc.OriginTypeName)
	responseServers, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("getting servers for ds-top-req-cap2 delivery service: %v", err.Error())
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
	ds1 := helperGetDeliveryServiceID(t, data[0].XMLID)

	testCases := []struct {
		description string
		capability  tc.DeliveryServicesRequiredCapability
	}{
		{
			description: "get all deliveryservices required capabilities",
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by deliveryServiceID: %d", *ds1),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID: ds1,
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by xmlID: %s", *data[0].XMLID),
			capability: tc.DeliveryServicesRequiredCapability{
				XMLID: data[0].XMLID,
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				RequiredCapability: data[0].RequiredCapability,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, reqInf, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.RequiredCapability, header)
			if err != nil {
				t.Fatalf("Expected no error, but got %v", err.Error())
			}
			if reqInf.StatusCode != http.StatusOK {
				t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
			}
		})
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, reqInf, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.RequiredCapability, header)
			if err != nil {
				t.Fatalf("Expected no error, but got %v", err.Error())
			}
			if reqInf.StatusCode != http.StatusNotModified {
				t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
			}
		})
	}
}

func GetTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	data := testData.DeliveryServicesRequiredCapabilities
	ds1 := helperGetDeliveryServiceID(t, data[0].XMLID)

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
					if *dsrc.RequiredCapability != *dsRequiredCapability.RequiredCapability {
						t.Errorf("expected: all delivery service required capabilities to equal %s, actual: found %s", *dsRequiredCapability.RequiredCapability, *dsrc.RequiredCapability)
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.RequiredCapability, nil)
			if err != nil {
				t.Fatalf("%s; got err= %v; expected err= nil", tc.description, err)
			}
			tc.expectFunc(tc.capability, capabilities)
		})
	}
}

func GetTestDeliveryServicesRequiredCapabilitiesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	data := testData.DeliveryServicesRequiredCapabilities
	ds1 := helperGetDeliveryServiceID(t, data[0].XMLID)

	testCases := []struct {
		description string
		capability  tc.DeliveryServicesRequiredCapability
	}{
		{
			description: "get all deliveryservices required capabilities",
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by deliveryServiceID: %d", *ds1),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID: ds1,
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by xmlID: %s", *data[0].XMLID),
			capability: tc.DeliveryServicesRequiredCapability{
				XMLID: data[0].XMLID,
			},
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				RequiredCapability: data[0].RequiredCapability,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, reqInf, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.RequiredCapability, header)
			if err != nil {
				t.Fatalf("Expected no error, but got %v", err.Error())
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

		_, _, err := TOSession.CreateDeliveryServicesRequiredCapability(c)
		if err != nil {
			t.Fatalf("cannot create delivery service required capability: %v", err)
		}
	}

	invalid := tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  helperGetDeliveryServiceID(t, util.StrPtr("ds-top-req-cap")),
		RequiredCapability: util.StrPtr("asdf"),
	}
	_, reqInf, err := TOSession.CreateDeliveryServicesRequiredCapability(invalid)
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

			_, _, err := TOSession.CreateDeliveryServicesRequiredCapability(cap)
			if err != nil {
				t.Fatalf(err.Error())
			}
		})
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, _, err := TOSession.CreateDeliveryServicesRequiredCapability(tc.capability)
			if err == nil {
				t.Fatalf("%s; expected err", tc.description)
			}
		})
	}
}

func InvalidDeliveryServicesRequiredCapabilityAddition(t *testing.T) {
	// Tests that a capability cannot be made required if the DS's services do not have it assigned

	// Get Delivery Capability for a DS
	capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(nil, util.StrPtr("ds1"), nil, nil)
	if err != nil {
		t.Fatalf("cannot GET delivery service required capabilities: %v", err)
	}
	if len(capabilities) == 0 {
		t.Fatal("delivery service ds1 needs at least one capability required")
	}

	// First assign current capabilities to edge server so we can assign it to the DS
	// TODO: DON'T hard-code hostnames!
	params := url.Values{}
	params.Add("hostName", "atlanta-edge-01")
	resp, _, err := TOSession.GetServersWithHdr(&params, nil)
	if err != nil {
		t.Fatalf("cannot GET Server by hostname: %v", err)
	}
	servers := resp.Response
	if len(servers) < 1 {
		t.Fatal("need at least one server to test invalid ds required capability assignment")
	}

	if servers[0].ID == nil {
		t.Fatal("server 'atlanta-edge-01' had nil ID")
	}

	dsID := capabilities[0].DeliveryServiceID
	sID := *servers[0].ID
	serverCaps := []tc.ServerServerCapability{}

	for _, cap := range capabilities {
		sCap := tc.ServerServerCapability{
			ServerID:         &sID,
			ServerCapability: cap.RequiredCapability,
		}
		_, _, err := TOSession.CreateServerServerCapability(sCap)
		if err != nil {
			t.Errorf("could not POST the server capability %v to server %v: %v", *cap.RequiredCapability, sID, err)
		}
		serverCaps = append(serverCaps, sCap)
	}

	// Assign server to ds
	_, _, err = TOSession.CreateDeliveryServiceServers(*dsID, []int{sID}, false)
	if err != nil {
		t.Fatalf("cannot CREATE server delivery service assignement: %v", err)
	}

	// Create new bogus server capability
	_, _, err = TOSession.CreateServerCapability(tc.ServerCapability{
		Name: "newcap",
	})
	if err != nil {
		t.Fatalf("cannot CREATE newcap server capability: %v", err)
	}

	// Attempt to assign to DS should fail
	_, _, err = TOSession.CreateDeliveryServicesRequiredCapability(tc.DeliveryServicesRequiredCapability{
		DeliveryServiceID:  dsID,
		RequiredCapability: util.StrPtr("newcap"),
	})
	if err == nil {
		t.Fatal("expected error requiring a capability that is not associated on the delivery service's servers")
	}

	// Disassociate server from DS
	setInactive(t, *dsID)
	_, _, err = TOSession.DeleteDeliveryServiceServer(*dsID, sID)
	if err != nil {
		t.Fatalf("could not DELETE the server %v from ds %v: %v", sID, *dsID, err)
	}

	// Remove server capabilities from server
	for _, ssc := range serverCaps {
		_, _, err := TOSession.DeleteServerServerCapability(*ssc.ServerID, *ssc.ServerCapability)
		if err != nil {
			t.Errorf("could not DELETE the server capability %v from server %v: %v", *ssc.ServerCapability, *ssc.Server, err)
		}
	}

	// Delete server capability
	_, _, err = TOSession.DeleteServerCapability("newcap")
	if err != nil {
		t.Fatalf("cannot DELETE newcap server capability: %v", err)
	}

}

func DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Get Required Capabilities to delete them
	capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(nil, nil, nil, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(capabilities) < 1 {
		t.Fatal("no delivery services returned")
	}

	type testCase struct {
		description string
		capability  tc.DeliveryServicesRequiredCapability
		err         string
	}

	testCases := []testCase{
		testCase{
			description: fmt.Sprintf("delete a deliveryservices required capability with an invalid delivery service id; deliveryServiceID: -1, requiredCapability: %s", *capabilities[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  util.IntPtr(-1),
				RequiredCapability: capabilities[0].RequiredCapability,
			},
			err: "no deliveryservice.RequiredCapability with that key found",
		},
		testCase{
			description: fmt.Sprintf("delete a deliveryservices required capability with an invalid required capability; deliveryServiceID: %d, requiredCapability: bogus", *capabilities[0].DeliveryServiceID),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID:  capabilities[0].DeliveryServiceID,
				RequiredCapability: util.StrPtr("bogus"),
			},
			err: "no deliveryservice.RequiredCapability with that key found",
		},
	}

	for _, c := range capabilities {
		t := testCase{
			description: fmt.Sprintf("delete a deliveryservices required capability; deliveryServiceID: %d, requiredCapability: %s", *c.DeliveryServiceID, *c.RequiredCapability),
			capability:  c,
		}
		testCases = append(testCases, t)
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			_, _, err := TOSession.DeleteDeliveryServicesRequiredCapability(*c.capability.DeliveryServiceID, *c.capability.RequiredCapability)
			if err != nil && !strings.Contains(err.Error(), c.err) {
				t.Fatalf("%s; got err= %s; expected err= %s", c.description, err, c.err)
			}
		})
	}
}

func helperGetDeliveryServiceID(t *testing.T, xmlID *string) *int {
	t.Helper()
	if xmlID == nil {
		t.Fatal("xml id must not be nil")
	}
	ds, _, err := TOSession.GetDeliveryServiceByXMLIDNullableWithHdr(*xmlID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 1 {
		t.Fatalf("cannot GET deliveyservice by xml id: %v. Response did not include record.", *xmlID)
	}
	return ds[0].ID
}
