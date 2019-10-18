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

func TestGetDeliveryServiceServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, ServerCapabilities, DeliveryServices, DeliveryServiceServerCapabilities}, func() {
		testCases := []struct {
			description string
			capability  tc.DeliveryServiceServerCapability
			expected    int
			err         string
		}{
			{
				description: "get all deliveryservice server capabilities",
				expected:    len(testData.DeliveryServiceServerCapabilities),
			},
			{
				description: fmt.Sprintf("get all deliveryservice server capabilities by delivery service id"),
				capability: tc.DeliveryServiceServerCapability{
					DeliveryServiceID: getDSID(t, testData.DeliveryServiceServerCapabilities[0]),
				},
				expected: 1,
			},
			{
				description: fmt.Sprintf("get all deliveryservice server capabilities by xml id"),
				capability: tc.DeliveryServiceServerCapability{
					XMLID: testData.DeliveryServiceServerCapabilities[0].XMLID,
				},
				expected: 1,
			},
			{
				description: fmt.Sprintf("get all deliveryservice server capabilities by server capability"),
				capability: tc.DeliveryServiceServerCapability{
					ServerCapability: testData.DeliveryServiceServerCapabilities[0].ServerCapability,
				},
				expected: 1,
			},
		}

		for i, tc := range testCases {
			t.Run(fmt.Sprintf("case %d: %s", i, tc.description), func(t *testing.T) {
				capabilities, _, err := TOSession.GetDeliveryServiceServerCapabilities(tc.capability.DeliveryServiceID, tc.capability.XMLID, tc.capability.ServerCapability)
				if err != nil {
					if !strings.Contains(err.Error(), tc.err) {
						t.Fatalf("case %d: %s; got err= %v; expected err = %v", i, tc.description, err, tc.err)
					}
				}
				if len(capabilities) != tc.expected {
					t.Errorf("got %d; expected %d server capabilities assigned to deliveryservices", len(capabilities), tc.expected)
				}
			})
		}
	})
}

func CreateTestDeliveryServiceServerCapabilities(t *testing.T) {
	testCases := []struct {
		description string
		capability  tc.DeliveryServiceServerCapability
		err         string
	}{
		{
			description: "re-assign a deliveryservice server capability",
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: getDSID(t, testData.DeliveryServiceServerCapabilities[0]),
				ServerCapability:  testData.DeliveryServiceServerCapabilities[0].ServerCapability,
			},
			err: fmt.Sprintf("deliveryservice_server_capability deliveryservice_id, server_capability '%d, foo' already exists", *getDSID(t, testData.DeliveryServiceServerCapabilities[0])),
		},
		{
			description: "assign a deliveryservice server capability with no deliverys service id",
			capability: tc.DeliveryServiceServerCapability{
				ServerCapability: testData.DeliveryServiceServerCapabilities[0].ServerCapability,
			},
			err: "'deliveryServiceID' cannot be blank",
		},
		{
			description: "assign a deliveryservice server capability with no server capability",
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: getDSID(t, testData.DeliveryServiceServerCapabilities[0]),
			},
			err: "'serverCapability' cannot be blank",
		},
		{
			description: "assign a deliveryservice server capability with an invalid server capability",
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: getDSID(t, testData.DeliveryServiceServerCapabilities[0]),
				ServerCapability:  util.StrPtr("bogus"),
			},
			err: "server_capability not found",
		},
		{
			description: "assign a deliveryservice server capability with an invalid deliver service id",
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: util.IntPtr(-1),
				ServerCapability:  testData.DeliveryServiceServerCapabilities[0].ServerCapability,
			},
			err: "deliveryservice not found",
		},
	}

	// Assign all server capability to delivery services listed in `tc-fixtrues.json`
	for i, td := range testData.DeliveryServiceServerCapabilities {
		t.Run(fmt.Sprintf("case %d: assign a deliveryservice server capability: %s", i, *td.XMLID), func(t *testing.T) {
			cap := tc.DeliveryServiceServerCapability{
				DeliveryServiceID: getDSID(t, td),
				ServerCapability:  td.ServerCapability,
			}

			_, _, err := TOSession.CreateDeliveryServiceServerCapability(cap)
			if err != nil {
				t.Fatalf(err.Error())
			}
		})
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.description), func(t *testing.T) {
			_, _, err := TOSession.CreateDeliveryServiceServerCapability(tc.capability)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("case %d: %s; got err= %v; expected err = %v", i, tc.description, err, tc.err)
				}
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
		err         string
	}

	testCases := []testCase{
		testCase{
			description: "delete a deliveryservice server capability with an invalid delivery service id",
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: util.IntPtr(-1),
				ServerCapability:  capabilities[0].ServerCapability,
			},
		},
		testCase{
			description: "delete a deliveryservice server capability with an invalid server capability",
			capability: tc.DeliveryServiceServerCapability{
				DeliveryServiceID: capabilities[0].DeliveryServiceID,
				ServerCapability:  util.StrPtr("bogus"),
			},
		},
	}

	for _, c := range capabilities {
		t := testCase{
			description: "delete a deliveryservice server capability",
			capability:  c,
		}
		testCases = append(testCases, t)
	}

	for i, c := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, c.description), func(t *testing.T) {
			_, _, err := TOSession.DeleteDeliveryServiceServerCapability(*c.capability.DeliveryServiceID, *c.capability.ServerCapability)
			if err != nil {
				if !strings.Contains(err.Error(), c.err) {
					t.Fatalf("case %d: %s; got err= %v; expected err = %v", i, c.description, err, c.err)
				}
			}
		})
	}
}

func getDSID(t *testing.T, capability tc.DeliveryServiceServerCapability) *int {
	ds, _, err := TOSession.GetDeliveryServiceByXMLID(*capability.XMLID)
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) != 1 {
		t.Fatalf("cannot GET deliveyservice by xml id: %v. Response did not include record.\n", *capability.XMLID)
	}
	return &ds[0].ID
}
