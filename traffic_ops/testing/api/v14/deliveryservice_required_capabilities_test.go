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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestGetDeliveryServiceRequiredCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, RequiredCapabilities, DeliveryServices, DeliveryServiceRequiredCapabilities}, func() {
		GetTestDeliveryServiceRequiredCapabilities(t)
	})
}

func GetTestDeliveryServiceRequiredCapabilities(t *testing.T) {
	data := testData.DeliveryServiceRequiredCapabilities
	ds1 := helperGetDeliveryServiceID(t, data[0])

	testCases := []struct {
		description string
		capability  tc.DeliveryServiceRequiredCapability
		expected    int
	}{
		{
			description: "get all deliveryservice required capabilities",
			expected:    len(testData.DeliveryServiceRequiredCapabilities),
		},
		{
			description: fmt.Sprintf("get all deliveryservice required capabilities by deliveryServiceID: %d", *ds1),
			capability: tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID: ds1,
			},
			expected: 1,
		},
		{
			description: fmt.Sprintf("get all deliveryservice required capabilities by xmlID: %s", *data[0].XMLID),
			capability: tc.DeliveryServiceRequiredCapability{
				XMLID: data[0].XMLID,
			},
			expected: 1,
		},
		{
			description: fmt.Sprintf("get all deliveryservice required capabilities by requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServiceRequiredCapability{
				RequiredCapability: data[0].RequiredCapability,
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			capabilities, _, err := TOSession.GetDeliveryServiceRequiredCapabilities(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.RequiredCapability)
			if err != nil {
				t.Fatalf("%s; got err= %v; expected err= nil", tc.description, err)
			}
			if len(capabilities) != tc.expected {
				t.Errorf("got %d; expected %d required capabilities assigned to deliveryservices", len(capabilities), tc.expected)
			}
		})
	}
}

func CreateTestDeliveryServiceRequiredCapabilities(t *testing.T) {
	data := testData.DeliveryServiceRequiredCapabilities
	ds1 := helperGetDeliveryServiceID(t, data[0])

	testCases := []struct {
		description string
		capability  tc.DeliveryServiceRequiredCapability
	}{
		{
			description: fmt.Sprintf("re-assign a deliveryservice to a required capability; deliveryServiceID: %d, requiredCapability: %s", *ds1, *data[0].RequiredCapability),
			capability: tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID:  ds1,
				RequiredCapability: data[0].RequiredCapability,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with no delivery service id; deliveryServiceID: 0, requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServiceRequiredCapability{
				RequiredCapability: data[0].RequiredCapability,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with no requiredCapability; deliveryServiceID: %d, requiredCapability: 0", *ds1),
			capability: tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID: ds1,
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with an invalid required capability; deliveryServiceID: %d, requiredCapability: bogus", *ds1),
			capability: tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID:  ds1,
				RequiredCapability: util.StrPtr("bogus"),
			},
		},
		{
			description: fmt.Sprintf("assign a deliveryservice to a required capability with an invalid deliver service id; deliveryServiceID: -1, requiredCapability: %s", *data[0].RequiredCapability),
			capability: tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID:  util.IntPtr(-1),
				RequiredCapability: data[0].RequiredCapability,
			},
		},
	}

	// Assign all required capability to delivery services listed in `tc-fixtrues.json`.
	for _, td := range testData.DeliveryServiceRequiredCapabilities {
		var dsID int
		if td.DeliveryServiceID != nil {
			dsID = *td.DeliveryServiceID
		}

		var capability string
		if td.RequiredCapability != nil {
			capability = *td.RequiredCapability
		}

		t.Run(fmt.Sprintf("assign a deliveryservice to a required capability; deliveryServiceID: %d, requiredCapability: %s", dsID, capability), func(t *testing.T) {
			cap := tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID:  helperGetDeliveryServiceID(t, td),
				RequiredCapability: td.RequiredCapability,
			}

			_, _, err := TOSession.CreateDeliveryServiceRequiredCapability(cap)
			if err != nil {
				t.Fatalf(err.Error())
			}
		})
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, _, err := TOSession.CreateDeliveryServiceRequiredCapability(tc.capability)
			if err == nil {
				t.Fatalf("%s; expected err", tc.description)
			}
		})
	}

}

func DeleteTestDeliveryServiceRequiredCapabilities(t *testing.T) {
	// Get Required Capabilities to delete them
	capabilities, _, err := TOSession.GetDeliveryServiceRequiredCapabilities(nil, nil, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	type testCase struct {
		description string
		capability  tc.DeliveryServiceRequiredCapability
		err         string
	}

	testCases := []testCase{
		testCase{
			description: fmt.Sprintf("delete a deliveryservice required capability with an invalid delivery service id; deliveryServiceID: -1, requiredCapability: %s", *capabilities[0].RequiredCapability),
			capability: tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID:  util.IntPtr(-1),
				RequiredCapability: capabilities[0].RequiredCapability,
			},
			err: "no deliveryservice.RequiredCapability with that key found",
		},
		testCase{
			description: fmt.Sprintf("delete a deliveryservice required capability with an invalid required capability; deliveryServiceID: %d, requiredCapability: bogus", *capabilities[0].DeliveryServiceID),
			capability: tc.DeliveryServiceRequiredCapability{
				DeliveryServiceID:  capabilities[0].DeliveryServiceID,
				RequiredCapability: util.StrPtr("bogus"),
			},
			err: "no deliveryservice.RequiredCapability with that key found",
		},
	}

	for _, c := range capabilities {
		t := testCase{
			description: fmt.Sprintf("delete a deliveryservice required capability; deliveryServiceID: %d, requiredCapability: %s", *c.DeliveryServiceID, *c.RequiredCapability),
			capability:  c,
		}
		testCases = append(testCases, t)
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			_, _, err := TOSession.DeleteDeliveryServiceRequiredCapability(*c.capability.DeliveryServiceID, *c.capability.RequiredCapability)
			if err != nil && !strings.Contains(err.Error(), c.err) {
				t.Fatalf("%s; got err= %s; expected err= %s", c.description, err, c.err)
			}
		})
	}
}

func helperGetDeliveryServiceID(t *testing.T, capability tc.DeliveryServiceRequiredCapability) *int {
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
