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
	"encoding/json"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
)

var fedIDs []int

func TestCDNFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServices, CDNFederations}, func() {
		UpdateTestCDNFederations(t)
		GetTestCDNFederations(t)
	})
}

func CreateTestCDNFederations(t *testing.T) {

	// Every federation is associated with a cdn
	for i, f := range testData.Federations {

		// CDNs test data and Federations test data are not naturally parallel
		if i >= len(testData.CDNs) {
			break
		}

		data, _, err := TOSession.CreateCDNFederationByName(f, testData.CDNs[i].Name)
		if err != nil {
			t.Errorf("could not POST federations: " + err.Error())
		}
		bytes, _ := json.Marshal(data)
		log.Debugf("POST Response: %s\n", bytes)

		// need to save the ids, otherwise the other tests won't be able to reference the federations
		if data.Response.ID == nil {
			t.Error("Federation id is nil after posting")
		} else {
			fedIDs = append(fedIDs, *data.Response.ID)
		}
	}
}

func UpdateTestCDNFederations(t *testing.T) {

	for _, id := range fedIDs {
		fed, _, err := TOSession.GetCDNFederationsByID("foo", id)
		if err != nil {
			t.Errorf("cannot GET federation by id: %v", err)
		}

		expectedCName := "new.cname."
		fed.Response[0].CName = &expectedCName
		resp, _, err := TOSession.UpdateCDNFederationsByID(fed.Response[0], "foo", id)
		if err != nil {
			t.Errorf("cannot PUT federation by id: %v", err)
		}
		bytes, _ := json.Marshal(resp)
		log.Debugf("PUT Response: %s\n", bytes)

		resp2, _, err := TOSession.GetCDNFederationsByID("foo", id)
		if err != nil {
			t.Errorf("cannot GET federation by id after PUT: %v", err)
		}
		bytes, _ = json.Marshal(resp2)
		log.Debugf("GET Response: %s\n", bytes)

		if resp2.Response[0].CName == nil {
			log.Errorln("CName is nil after updating")
		} else if *resp2.Response[0].CName != expectedCName {
			t.Errorf("results do not match actual: %s, expected: %s\n", *resp2.Response[0].CName, expectedCName)
		}

	}
}

func GetTestCDNFederations(t *testing.T) {

	// TOSession.GetCDNFederationsByName can't be tested until
	// POST /api/1.2/federations/:id/deliveryservices has been
	// created. (DELETE cdns/:name/federations/:id may need to
	// clean up fedIDs connection?)

	for _, id := range fedIDs {
		data, _, err := TOSession.GetCDNFederationsByID("foo", id)
		if err != nil {
			t.Errorf("could not GET federations: " + err.Error())
		}
		bytes, _ := json.Marshal(data)
		log.Debugf("GET Response: %s\n", bytes)
	}
}

func DeleteTestCDNFederations(t *testing.T) {

	for _, id := range fedIDs {
		resp, _, err := TOSession.DeleteCDNFederationByID("foo", id)
		if err != nil {
			t.Errorf("cannot DELETE federation by id: '%d' %v\n", id, err)
		}
		bytes, err := json.Marshal(resp)
		log.Debugf("DELETE Response: %s\n", bytes)

		data, _, err := TOSession.GetCDNFederationsByID("foo", id)
		if len(data.Response) != 0 {
			t.Error("expected federation to be deleted")
		}
	}
	fedIDs = nil // reset the global variable for the next test
}
