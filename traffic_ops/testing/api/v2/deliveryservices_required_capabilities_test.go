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
	"fmt"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
)

func TestDeliveryServicesRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, DeliveryServices, DeliveryServicesRequiredCapabilities}, func() {
		InvalidDeliveryServicesRequiredCapabilityAddition(t)
		GetTestDeliveryServicesRequiredCapabilities(t)
	})
}

func GetTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	data := testData.DeliveryServicesRequiredCapabilities
	ds1 := helperGetDeliveryServiceID(t, data[0].XMLID)

	testCases := []struct {
		description string
		capability  tc.DeliveryServicesRequiredCapability
		expected    int
	}{
		{
			description: "get all deliveryservices required capabilities",
			expected:    len(testData.DeliveryServicesRequiredCapabilities),
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by deliveryServiceID: %d", *ds1),
			capability: tc.DeliveryServicesRequiredCapability{
				DeliveryServiceID: ds1,
			},
			expected: 1,
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by xmlID: %s", *data[0].XMLID),
			capability: tc.DeliveryServicesRequiredCapability{
				XMLID: data[0].XMLID,
			},
			expected: 1,
		},
		{
			description: fmt.Sprintf("get all deliveryservices required capabilities by requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServicesRequiredCapability{
				RequiredCapability: data[0].RequiredCapability,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.RequiredCapability)
			if err != nil {
				t.Fatalf("%s; got err= %v; expected err= nil", tc.description, err)
			}
			if len(capabilities) != tc.expected {
				t.Errorf("got %d; expected %d required capabilities assigned to deliveryservices", len(capabilities), tc.expected)
			}
		})
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
	capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(nil, util.StrPtr("ds1"), nil)
	if err != nil {
		t.Fatalf("cannot GET delivery service required capabilities: %v", err)
	}
	if len(capabilities) == 0 {
		t.Fatal("delivery service ds1 needs at least one capability required")
	}

	// First assign current capabilities to edge server so we can assign it to the DS
	servers, _, err := TOSession.GetServerByHostName("atlanta-edge-01")
	if err != nil {
		t.Fatalf("cannot GET Server by hostname: %v", err)
	}
	if len(servers) < 1 {
		t.Fatal("need at least one server to test invalid ds required capability assignment")
	}

	dsID := capabilities[0].DeliveryServiceID
	sID := servers[0].ID
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
	_, err = TOSession.CreateDeliveryServiceServers(*dsID, []int{sID}, false)
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
	capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilities(nil, nil, nil)
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
	ds, _, err := TOSession.GetDeliveryServiceByXMLIDNullable(*xmlID)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 1 {
		t.Fatalf("cannot GET deliveyservice by xml id: %v. Response did not include record.", *xmlID)
	}
	return ds[0].ID
}
