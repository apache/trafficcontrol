package v4

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
	"net/http"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

var fedIDs []int

func TestCDNFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, DeliveryServices, CDNFederations}, func() {
		SortTestCDNFederations(t)
		UpdateTestCDNFederations(t)
		GetTestCDNFederations(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestCDNFederationsWithHeaders(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestCDNFederationsWithHeaders(t, header)
	})
}

func UpdateTestCDNFederationsWithHeaders(t *testing.T, h http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = h
	for _, id := range fedIDs {
		opts.QueryParameters.Set("id", strconv.Itoa(id))
		fed, _, err := TOSession.GetCDNFederationsByName("foo", opts)
		if err != nil {
			t.Errorf("cannot GET federation by id: %v - alerts: %+v", err, fed.Alerts)
		}
		if len(fed.Response) > 0 {
			expectedCName := "new.cname."
			fed.Response[0].CName = &expectedCName
			_, reqInf, err := TOSession.UpdateCDNFederation(fed.Response[0], "foo", id, opts)
			if err == nil {
				t.Errorf("Expected an error saying precondition failed, but got none")
			}
			if reqInf.StatusCode != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		}
	}
}

func TestFederationFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, DeliveryServices, CDNFederations, FederationResolvers}, func() {
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

		data, _, err := TOSession.CreateCDNFederation(f, testData.CDNs[i].Name, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create CDN Federations: %v - alerts: %+v", err, data.Alerts)
		}

		// need to save the ids, otherwise the other tests won't be able to reference the federations
		if data.Response.ID == nil {
			t.Error("Federation id is nil after posting")
		} else {
			fedIDs = append(fedIDs, *data.Response.ID)
		}
	}
}

// This test will not work unless a given CDN has more than one federation associated with it.
func SortTestCDNFederations(t *testing.T) {
	var sortedList []string

	//Create a new federation under the same CDN
	cname := "bar.foo."
	ttl := 50
	description := "test"
	f := tc.CDNFederation{ID: nil, CName: &cname, TTL: &ttl, Description: &description, LastUpdated: nil}
	data, _, err := TOSession.CreateCDNFederation(f, "cdn1", client.RequestOptions{})
	if err != nil {
		t.Errorf("could not create CDN Federations: %v - alerts: %+v", err, data.Alerts)
	}
	id := *data.Response.ID

	//Get list of federations for one type of cdn
	resp, _, err := TOSession.GetCDNFederationsByName("cdn1", client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
	}
	for i := range resp.Response {
		sortedList = append(sortedList, *resp.Response[i].CName)
	}

	// Check if list was sorted
	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}

	// Delete the newly created federation
	resp1, _, err := TOSession.DeleteCDNFederation("cdn1", id, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot delete federation #%d: %v - alerts: %+v", id, err, resp1.Alerts)
	}
}

func UpdateTestCDNFederations(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, id := range fedIDs {
		opts.QueryParameters.Set("id", strconv.Itoa(id))
		fed, _, err := TOSession.GetCDNFederationsByName("foo", opts)
		if err != nil {
			t.Errorf("cannot get federation by id: %v - alerts: %+v", err, fed.Alerts)
			continue
		}
		if len(fed.Response) < 1 {
			t.Error("Response from Traffic Ops included no Federations")
			continue
		}

		expectedCName := "new.cname."
		fed.Response[0].CName = &expectedCName
		resp, _, err := TOSession.UpdateCDNFederation(fed.Response[0], "foo", id, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot update federation by id: %v - alerts: %+v", err, resp.Alerts)
		}

		resp2, _, err := TOSession.GetCDNFederationsByName("foo", opts)
		if err != nil {
			t.Errorf("cannot get federation by id after update: %v - alerts: %+v", err, resp2.Alerts)
		}
		if len(resp2.Response) < 1 {
			t.Error("cannot get federation by id after update: Federation not found in Traffic Ops response")
		}

		if resp2.Response[0].CName == nil {
			t.Error("CName is nil after updating")
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

	opts := client.NewRequestOptions()
	for _, id := range fedIDs {
		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, _, err := TOSession.GetCDNFederationsByName("foo", opts)
		if err != nil {
			t.Errorf("could not get federations: %v - alerts: %+v", err, data.Alerts)
		}
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

	frs, _, err := TOSession.GetFederationResolvers(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error getting Federation Resolvers: %v - alerts: %+v", err, frs.Alerts)
	}
	if len(frs.Response) != frCnt {
		t.Fatalf("Wrong number of Federation Resolvers from GET, want %d got %d", frCnt, len(frs.Response))
	}

	frIDs := make([]int, 0, len(frs.Response))
	for _, fr := range frs.Response {
		if fr.ID == nil {
			// Because we need 'frCnt' resolver ids, this is not recoverable
			t.Fatal("Traffic Ops returned a representation for a Federation Resolver with null or undefined ID")
		}
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
			resp, _, err := TOSession.AssignFederationFederationResolver(c.fedID, c.resolverIDs, c.replace, client.RequestOptions{})
			if err != nil {
				if c.err == "" {
					t.Errorf("Unexpected error assigning Federation Resolvers to Federation: %v - alerts: %+v", err, resp.Alerts)
					return
				}
				found := false
				for _, alert := range resp.Alerts.Alerts {
					if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, c.err) {
						found = true
					}
				}
				if !found {
					t.Errorf("Expected to find an error-level alert relating to '%s', but didn't - err: %v - alerts: %+v", c.err, err, resp.Alerts)
				}
			} else if c.err != "" {
				t.Errorf("Expected to get an error assigning Federation Resolvers to Federation, but didn't - alerts: %+v", resp.Alerts)
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
			resp, _, err := TOSession.GetFederationFederationResolvers(c.fedID, client.RequestOptions{})
			if err != nil {
				t.Fatalf("Error getting federation federation resolvers by federation id %d: %v - alerts: %+v", c.fedID, err, resp.Alerts)
			}
			if len(resp.Response) == 0 && c.hasRecords {
				t.Fatalf("expected federation of ID %d to have associated federation resolvers, but had 0", c.fedID)
			}
		})
	}
}

func DeleteTestCDNFederations(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, id := range fedIDs {
		resp, _, err := TOSession.DeleteCDNFederation("foo", id, opts)
		if err != nil {
			t.Errorf("cannot delete federation #%d: %v - alerts: %+v", id, err, resp.Alerts)
		}

		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, _, err := TOSession.GetCDNFederationsByName("foo", opts)
		if len(data.Response) != 0 {
			t.Error("expected federation to be deleted")
		}
		opts.QueryParameters.Del("id")
	}
	fedIDs = nil // reset the global variable for the next test
}
