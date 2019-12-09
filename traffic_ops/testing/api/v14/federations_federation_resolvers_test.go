package v14

import (
	"testing"
)

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
/*

func TestFederationFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServices, CDNFederations, FederationResolvers, FederationsFederationResolvers}, func() {
		GetTestFederationFederationResolvers(t)
	})
}

func GetTestFederationFederationResolvers(t *testing.T) {
	tdlen := len(testData.Federations)
	if tdlen < 1 {
		t.Fatal("no federations to test")
	}

	fs, _, err := TOSession.AllFederations()
	if err != nil {
		t.Errorf("Unexpected error getting Federations: %v", err)
	}
	if len(fs) != tdlen {
		t.Fatalf("Wrong number of Federation from GET, want %d got %d", tdlen, len(fs))
	}
	for _, id := range fedIDs {
		data, _, err := TOSession.GetCDNFederationsByID("foo", id)
		if err != nil {
			t.Errorf("could not GET federations: " + err.Error())
		}
		bytes, _ := json.Marshal(data)
		log.Debugf("GET Response: %s\n", bytes)
	}

	tdlen = len(testData.FederationResolvers)
	if tdlen < 1 {
		t.Fatal("no federation resolvers to test")
	}

	frs, _, err := TOSession.GetFederationResolvers()
	if err != nil {
		t.Errorf("Unexpected error getting Federation Resolvers: %v", err)
	}
	if len(frs) != tdlen {
		t.Fatalf("Wrong number of Federation Resolvers from GET, want %d got %d", tdlen, len(frs))
	}
}
*/

func CreateTestFederationsFederationResolvers(t *testing.T) {
	/*var frIDs []int
	var fedIDs []int

	// GET ALL
	feds, _, err := TOSession.AllFederations()
	fmt.Printf("\nfeds:%v, err:%v", feds, err)
	for _, f := range testData.Federations {
		n := *f.CName
		cf, _, err := TOSession.GetCDNFederationsByName(n)
		if err != nil {
			t.Fatalf("could not get cdn federation by name %s, error:%s", n, err.Error())
		}
		id := *cf.Response[0].ID
		fedIDs = append(fedIDs, id)
	}

	for _, fr := range testData.FederationResolvers {
		frIDs = append(frIDs, int(*fr.ID))
	}

	testCases := []struct {
		description string
		fedID       int
		resolverIDs []int
	}{
		{
			description: fmt.Sprintf("assign resolver(s):%v to federation:%v", frIDs, fedIDs[0]),
			fedID:       fedIDs[0],
			resolverIDs: frIDs,
		},
	}

	// Assign all FederationsFederationResolvers listed in `tc-fixtures.json`.
	for _, td := range testData.FederationsFederationResolvers {
		fmt.Printf("td:%v", td)
		t.Run(fmt.Sprintf("assign a federation resolver to a federation; federation: %d, federation_resolver: %d", td.Federation, td.FederationResolver), func(t *testing.T) {
			_, _, err := TOSession.AssignFederationsFederationResolver(td.Federation, []int{td.FederationResolver}, true)
			if err != nil {
				t.Fatalf(err.Error())
			}
		})
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, _, err := TOSession.AssignFederationsFederationResolver(tc.fedID, tc.resolverIDs, true)
			if err == nil {
				t.Fatalf("%s; expected err", tc.description)
			}
		})
	}
	*/
}

func DeleteTestFederationsFederationResolvers(t *testing.T) {
	// Get FFRs to delete them
	/*
		var ids []int

		for _, f := range testData.Federations {
			_, _, err := TOSession.GetFederationsFederationResolversByID(*f.ID)
			if err != nil {
				t.Fatalf(err.Error())
			}
			ids = append(ids, *f.ID)
		}

		testCases := []struct {
			description string
			fedID       int
			err         string
		}{
			{
				description: fmt.Sprintf("delete a federation_federation_resolver from fed %d", ids[0]),
				fedID:       ids[0],
				err:         "",
			},
			{
				description: fmt.Sprintf("delete a federation_federation_resolver from fed %d", ids[1]),
				fedID:       ids[1],
				err:         "",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				_, _, err := TOSession.AssignFederationsFederationResolver(tc.fedID, []int{}, true)
				if err != nil && !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("%s; got err= %s; expected err= %s", tc.description, err, tc.err)
				}
			})
		}
	*/
}
