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
	"encoding/json"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v6/lib/go-log"
)

var fedIDs []int

func TestCDNFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServices, CDNFederations}, func() {
		UpdateTestCDNFederations(t)
		GetTestCDNFederations(t)
	})
}

func TestFederationFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServices, CDNFederations, FederationResolvers}, func() {
		AssignTestFederationFederationResolvers(t)
		GetTestFederationFederationResolvers(t)
	})
}

func CreateTestCDNFederations(t *testing.T) {

	// Every federation is associated with a cdn
	for i, f := range testData.Federations {

		// CDNs test data and Federations test data are not naturally parallel
		if i >= len(testData.CDNs) {
			break
		}

		data, _, err := TOSession.CreateCDNFederationByName(f, "cdn1")
		if err != nil {
			t.Errorf("could not POST federations: " + err.Error())
		}
		bytes, _ := json.Marshal(data)
		t.Logf("POST Response: %s\n", bytes)

		// need to save the ids, otherwise the other tests won't be able to reference the federations
		if data.Response.ID == nil {
			t.Error("Federation id is nil after posting")
		} else {
			fedIDs = append(fedIDs, *data.Response.ID)
			resp, _, err := TOSession.GetDeliveryServiceByXMLIDNullable("ds1")
			if err != nil {
				t.Errorf("could not get delivery service by xml ID: %v", err)
			}
			if len(resp) != 1 {
				t.Fatalf("expected one response for delivery service, but got %d", len(resp))
			}
			_, err = TOSession.CreateFederationDeliveryServices(*data.Response.ID, []int{*resp[0].ID}, false)
			if err != nil {
				t.Errorf("could not create federation delivery service: %v", err)
			}
		}
	}
}

func UpdateTestCDNFederations(t *testing.T) {

	for _, id := range fedIDs {
		fed, _, err := TOSession.GetCDNFederationsByID("cdn1", id)
		if err != nil {
			t.Errorf("cannot GET federation by id: %v", err)
		}

		expectedCName := "new.cname."
		fed.Response[0].CName = &expectedCName
		resp, _, err := TOSession.UpdateCDNFederationsByID(fed.Response[0], "cdn1", id)
		if err != nil {
			t.Errorf("cannot PUT federation by id: %v", err)
		}
		bytes, _ := json.Marshal(resp)
		t.Logf("PUT Response: %s\n", bytes)

		resp2, _, err := TOSession.GetCDNFederationsByID("cdn1", id)
		if err != nil {
			t.Errorf("cannot GET federation by id after PUT: %v", err)
		}
		bytes, _ = json.Marshal(resp2)
		t.Logf("GET Response: %s\n", bytes)

		if resp2.Response[0].CName == nil {
			log.Errorln("CName is nil after updating")
		} else if *resp2.Response[0].CName != expectedCName {
			t.Errorf("results do not match actual: %s, expected: %s", *resp2.Response[0].CName, expectedCName)
		}

	}
}

func GetTestCDNFederations(t *testing.T) {

	// TOSession.GetCDNFederationsByName can't be tested until
	// POST /federations/{{id}}/deliveryservices has been
	// created. (DELETE cdns/:name/federations/:id may need to
	// clean up fedIDs connection?)

	for _, id := range fedIDs {
		data, _, err := TOSession.GetCDNFederationsByID("cdn1", id)
		if err != nil {
			t.Errorf("could not GET federations: " + err.Error())
		}
		bytes, _ := json.Marshal(data)
		t.Logf("GET Response: %s\n", bytes)
	}
}

func AssignTestFederationFederationResolvers(t *testing.T) {
	// Setup
	if len(fedIDs) < 2 {
		t.Fatal("not enough federations to test")
	}

	frCnt := len(testData.FederationResolvers)
	if frCnt < 2 {
		t.Fatal("not enough federation resolvers to test")
	}

	frs, _, err := TOSession.GetFederationResolvers()
	if err != nil {
		t.Fatalf("Unexpected error getting Federation Resolvers: %v", err)
	}
	if len(frs) != frCnt {
		t.Fatalf("Wrong number of Federation Resolvers from GET, want %d got %d", frCnt, len(frs))
	}

	frIDs := make([]int, 0, len(frs))
	for _, fr := range frs {
		frIDs = append(frIDs, int(*fr.ID))
	}

	// Test Cases
	testCases := []struct {
		description string
		fedID       int
		resolverIDs []int
		replace     bool
		err         string
	}{
		{
			description: "Successfully assign one federation_resolver to a federation",
			fedID:       fedIDs[0],
			resolverIDs: frIDs[0:0],
			replace:     false,
			err:         "",
		},
		{
			description: "Successfully assign multiple federation_resolver to a federation",
			fedID:       fedIDs[0],
			resolverIDs: frIDs[1:frCnt],
			replace:     false,
			err:         "",
		},
		{
			description: "Successfully replace all federation_resolver for a federation",
			fedID:       fedIDs[0],
			resolverIDs: frIDs[0:frCnt],
			replace:     true,
			err:         "",
		},
		{
			description: "Fail to assign federation_resolver to a federation when federation does not exist",
			fedID:       -1,
			resolverIDs: frIDs[0:0],
			replace:     false,
			err:         "no such Federation",
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			_, _, err := TOSession.AssignFederationFederationResolver(c.fedID, c.resolverIDs, c.replace)

			if err != nil && !strings.Contains(err.Error(), c.err) {
				t.Fatalf("error: expected error result %v, want: %v", err, c.err)
			}
		})
	}

}

func GetTestFederationFederationResolvers(t *testing.T) {
	if len(fedIDs) < 2 {
		t.Fatal("not enough federations to test")
	}

	testCases := []struct {
		description string
		fedID       int
		hasRecords  bool
	}{
		{
			description: "successfully get federation_federation_resolvers for a federation with some",
			fedID:       fedIDs[0],
			hasRecords:  true,
		},
		{
			description: "successfully get federation_federation_resolvers for a federation without any",
			fedID:       fedIDs[1],
			hasRecords:  false,
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			resp, _, err := TOSession.GetFederationFederationResolversByID(c.fedID)
			if err != nil {
				t.Fatalf("Error getting federation federation resolvers by federation id: %d, err: %s", c.fedID, err.Error())
			}
			if len(resp.Response) == 0 && c.hasRecords {
				t.Fatalf("expected federation of ID %d to have associated federation resolvers, but had 0", c.fedID)
			}
		})
	}
}

func DeleteTestCDNFederations(t *testing.T) {

	for _, id := range fedIDs {
		resp, _, err := TOSession.DeleteCDNFederationByID("cdn1", id)
		if err != nil {
			t.Errorf("cannot DELETE federation by id: '%d' %v", id, err)
		}
		bytes, err := json.Marshal(resp)
		t.Logf("DELETE Response: %s\n", bytes)

		data, _, err := TOSession.GetCDNFederationsByID("cdn1", id)
		if len(data.Response) != 0 {
			t.Error("expected federation to be deleted")
		}
	}
	fedIDs = nil // reset the global variable for the next test
}
