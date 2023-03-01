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
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v7/lib/go-tc"
)

func (r *TCData) CreateTestDeliveryServices(t *testing.T) {
	pl := tc.Parameter{
		ConfigFile: "remap.config",
		Name:       "location",
		Value:      "/remap/config/location/parameter/",
	}
	_, _, err := TOSession.CreateParameter(pl)
	if err != nil {
		t.Errorf("cannot create parameter: %v", err)
	}
	for _, ds := range r.TestData.DeliveryServices {
		_, _, err = TOSession.CreateDeliveryServiceV30(ds)
		if err != nil {
			t.Errorf("could not CREATE delivery service '%s': %v", *ds.XMLID, err)
		}
	}
}

func (r *TCData) DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		t.Errorf("cannot GET deliveryservices: %v", err)
	}
	for _, testDS := range r.TestData.DeliveryServices {
		var ds tc.DeliveryServiceNullableV30
		found := false
		for _, realDS := range dses {
			if realDS.XMLID != nil && *realDS.XMLID == *testDS.XMLID {
				ds = realDS
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DeliveryService not found in Traffic Ops: %s", *ds.XMLID)
			continue
		}

		delResp, err := TOSession.DeleteDeliveryService(strconv.Itoa(*ds.ID))
		if err != nil {
			t.Errorf("cannot DELETE DeliveryService by ID: %v - %v", err, delResp)
			continue
		}

		// Retrieve the Server to see if it got deleted
		params := url.Values{}
		params.Set("id", strconv.Itoa(*ds.ID))
		foundDS, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, params)
		if err != nil {
			t.Errorf("Unexpected error deleting Delivery Service '%s': %v", *ds.XMLID, err)
		}
		if len(foundDS) > 0 {
			t.Errorf("expected Delivery Service: %s to be deleted, but %d exist with same ID (#%d)", *ds.XMLID, len(foundDS), *ds.ID)
		}
	}

	// clean up parameter created in CreateTestDeliveryServices()
	params, _, err := TOSession.GetParameterByNameAndConfigFile("location", "remap.config")
	for _, param := range params {
		deleted, _, err := TOSession.DeleteParameterByID(param.ID)
		if err != nil {
			t.Errorf("cannot DELETE parameter by ID (%d): %v - %v", param.ID, err, deleted)
		}
	}
}
