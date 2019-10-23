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
	"fmt"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestGetDeliveryServiceServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, DeliveryServices, DeliveryServiceServerCapabilities}, func() {
		GetTestDeliveryServiceServerCapabilities(t)
	})
}

func GetTestDeliveryServiceServerCapabilities(t *testing.T) {
	data := testData.DeliveryServiceServerCapabilities
	ds1 := helperGetDeliveryServiceID(t, data[0])

	testCases := []struct {
		description string
		capability  tc.DeliveryServiceServerCapability
		expected    int
	}{
		{
			description: "get all deliveryservice server capabilities",
			expected:    len(testData.DeliveryServiceServerCapabilities),
		},
		{
			description: fmt.Sprintf("get all deliveryservice server capabilities by deliveryServiceID: %d", *ds1),
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: ds1,
			},
			expected: 1,
		},
		{
			description: fmt.Sprintf("get all deliveryservice server capabilities by xmlID: %s", *data[0].XMLID),
			capability: tc.DeliveryServiceServerCapability{
				XMLID: data[0].XMLID,
			},
			expected: 1,
		},
		{
			description: fmt.Sprintf("get all deliveryservice server capabilities by serverCapability: %s", *data[0].ServerCapability),
			capability: tc.DeliveryServiceServerCapability{
				ServerCapability: data[0].ServerCapability,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			capabilities, _, err := TOSession.GetDeliveryServiceServerCapabilities(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.ServerCapability)
			if err != nil {
				t.Fatalf("%s; got err= %v; expected err nil", tc.description, err)
			}
			if len(capabilities) != tc.expected {
				t.Errorf("got %d; expected %d server capabilities assigned to deliveryservices", len(capabilities), tc.expected)
			}
		})
	}
}

func CreateTestDeliveryServiceServerCapabilities(t *testing.T) {
	data := testData.DeliveryServiceServerCapabilities
	ds1 := helperGetDeliveryServiceID(t, data[0])

	testCases := []struct {
		description string
		capability  tc.DeliveryServiceServerCapability
	}{
		{
			description: fmt.Sprintf("re-assign a deliveryservice to a server capability; deliveryServiceID: %d, serverCapability: %s", *ds1, *data[0].ServerCapability),
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: ds1,
				ServerCapability:  data[0].ServerCapability,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a server capability with no delivery service id; deliveryServiceID: 0, serverCapability: %s", *data[0].ServerCapability),
			capability: tc.DeliveryServiceServerCapability{
				ServerCapability: data[0].ServerCapability,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a server capability with no serverCapability; deliveryServiceID: %d, serverCapability: 0", *ds1),
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: ds1,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a server capability with an invalid server capability; deliveryServiceID: %d, serverCapability: bogus", *ds1),
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: ds1,
				ServerCapability:  util.StrPtr("bogus"),
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a server capability with an invalid deliver service id; deliveryServiceID: -1, serverCapability: %s", *data[0].ServerCapability),
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: util.IntPtr(-1),
				ServerCapability:  data[0].ServerCapability,
			},
		},
	}

	// Assign all server capability to delivery services listed in `tc-fixtrues.json`.
	for _, td := range testData.DeliveryServiceServerCapabilities {
		var dsID int
		if td.DeliveryServiceID != nil {
			dsID = *td.DeliveryServiceID
		}

		var capability string
		if td.ServerCapability != nil {
			capability = *td.ServerCapability
		}

		t.Run(fmt.Sprintf("assign a deliveryservice to a server capability; deliveryServiceID: %d, serverCapability: %s", dsID, capability), func(t *testing.T) {
			cap := tc.DeliveryServiceServerCapability{
				DeliveryServiceID: helperGetDeliveryServiceID(t, td),
				ServerCapability:  td.ServerCapability,
			}

			_, _, err := TOSession.CreateDeliveryServiceServerCapability(cap)
			if err != nil {
				t.Fatalf(err.Error())
			}
		})
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, _, err := TOSession.CreateDeliveryServiceServerCapability(tc.capability)
			if err == nil {
				t.Fatalf("%s; expected error", tc.description)
			}
		})
	}

}

func DeleteTestDeliveryServiceServerCapabilities(t *testing.T) {
	// Get Server Capabilities to delete them
	capabilities, _, err := TOSession.GetDeliveryServiceServerCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	type testCase struct {
		description string
		capability  tc.DeliveryServiceServerCapability
	}

	testCases := []testCase{
		testCase{
			description: fmt.Sprintf("delete a deliveryservice server capability with an invalid delivery service id; deliveryServiceID: -1, serverCapability: %s", *capabilities[0].ServerCapability),
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: util.IntPtr(-1),
				ServerCapability:  capabilities[0].ServerCapability,
			},
		},
		testCase{
			description: fmt.Sprintf("delete a deliveryservice server capability with an invalid server capability; deliveryServiceID: %d, serverCapability: bogus", *capabilities[0].DeliveryServiceID),
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: capabilities[0].DeliveryServiceID,
				ServerCapability:  util.StrPtr("bogus"),
			},
		},
	}

	for _, c := range capabilities {
		t := testCase{
			description: fmt.Sprintf("delete a deliveryservice server capability; deliveryServiceID: %d, serverCapability: %s", *c.DeliveryServiceID, *c.ServerCapability),
			capability:  c,
		}
		testCases = append(testCases, t)
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			_, _, err := TOSession.DeleteDeliveryServiceServerCapability(*c.capability.DeliveryServiceID, *c.capability.ServerCapability)
			if err == nil {
				t.Fatalf("%s; expected err", c.description)
			}
		})
	}
}

func helperGetDeliveryServiceID(t *testing.T, capability tc.DeliveryServiceServerCapability) *int {
	t.Helper()
	ds, _, err := TOSession.GetDeliveryServiceByXMLID(*capability.XMLID)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) != 1 {
		t.Fatalf("cannot GET deliveyservice by xml id: %v. Response did not include record.\n", *capability.XMLID)
	}
	return &ds[0].ID
}
