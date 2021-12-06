package tcdata

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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func (r *TCData) CreateTestTopologyBasedDeliveryServicesRequiredCapabilities(t *testing.T) {
	for _, td := range r.TestData.TopologyBasedDeliveryServicesRequiredCapabilities {

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

func (r *TCData) CreateTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	data := r.TestData.DeliveryServicesRequiredCapabilities
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
	for _, td := range r.TestData.DeliveryServicesRequiredCapabilities {
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
				t.Fatalf("failed to associate a capability with a Delivery Service: %v", err)
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

func (r *TCData) DeleteTestDeliveryServicesRequiredCapabilities(t *testing.T) {
	// Get Required Capabilities to delete them
	capabilities, _, err := TOSession.GetDeliveryServicesRequiredCapabilitiesWithHdr(nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to fetch associations between Capabilities and the Delivery Services that require them: %v", err)
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
		t.Fatalf("cannot GET deliveyservice by xml id: %s, response did not include record", *xmlID)
	}
	return ds[0].ID
}
