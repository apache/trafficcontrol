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
	"net/url"
	"reflect"
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
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations}, func() {
		SortTestCDNFederations(t)
		SortTestCDNFederationsDesc(t)
		UpdateTestCDNFederations(t)
		GetTestCDNFederations(t)
		GetTestCDNFederationsIMS(t)
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
		GetTestPaginationSupportCdnFederation(t)
	})
}

func UpdateTestCDNFederationsWithHeaders(t *testing.T, h http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = h
	for _, id := range fedIDs {
		opts.QueryParameters.Set("id", strconv.Itoa(id))
		fed, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
		if err != nil {
			t.Errorf("cannot GET federation by id: %v - alerts: %+v", err, fed.Alerts)
		}
		if len(fed.Response) > 0 {
			expectedCName := "new.cname."
			fed.Response[0].CName = &expectedCName
			_, reqInf, err := TOSession.UpdateCDNFederation(fed.Response[0], "cdn1", id, opts)
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
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles, Tenants, CacheGroups, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations, FederationResolvers}, func() {
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
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", "ds1")
		resp, _, err := TOSession.GetDeliveryServices(opts)
		if err != nil {
			t.Errorf("could not get delivery service by xml ID: %v", err)
		}
		if len(resp.Response) != 1 {
			t.Fatalf("expected one response for delivery service, but got %d", len(resp.Response))
		}
		data, _, err := TOSession.CreateCDNFederation(f, "cdn1", client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create CDN Federations: %v - alerts: %+v", err, data.Alerts)
		}

		// need to save the ids, otherwise the other tests won't be able to reference the federations
		if data.Response.ID == nil {
			t.Error("Federation id is nil after posting")
		} else {
			fedIDs = append(fedIDs, *data.Response.ID)
			_, _, err = TOSession.CreateFederationDeliveryServices(*data.Response.ID, []int{*resp.Response[0].ID}, false, client.NewRequestOptions())
			if err != nil {
				t.Errorf("could not create federation delivery service: %v", err)
			}
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
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "cname")
	resp, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
	if err != nil {
		t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
	}
	for i := range resp.Response {
		if resp.Response[i].CName == nil {
			t.Fatalf("Federation resolver CName is nil, so sorting can't be tested")
		}
		sortedList = append(sortedList, *resp.Response[i].CName)
	}

	// Check if list was sorted
	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their CName: %v", sortedList)
	}

	// Delete the newly created federation
	resp1, _, err := TOSession.DeleteCDNFederation("cdn1", id, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot delete federation #%d: %v - alerts: %+v", id, err, resp1.Alerts)
	}
}

func SortTestCDNFederationsDesc(t *testing.T) {

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
	if err != nil {
		t.Fatalf("Expected no error, but got error in CDN Federation default ordering %v - alerts: %+v", err, resp.Alerts)
	}
	respAsc := resp.Response
	if len(respAsc) < 1 {
		t.Fatal("Need at least one CDN Federation in Traffic Ops to test CDN Federation sort ordering")
	}
	opts.QueryParameters.Set("sortOrder", "desc")
	resp, _, err = TOSession.GetCDNFederationsByName("cdn1", opts)
	if err != nil {
		t.Errorf("Expected no error, but got error in CDN Federation with Descending ordering: %v - alerts: %+v", err, resp.Alerts)
	}
	respDesc := resp.Response
	if len(respDesc) < 1 {
		t.Fatal("Need at least one CDN Federation in Traffic Ops to test CDN Federation sort ordering")
	}
	if len(respAsc) != len(respDesc) {
		t.Fatalf("Traffic Ops returned %d CDN Federation using default sort order, but %d CDN Federation when sort order was explicitly set to descending", len(respAsc), len(respDesc))
	}
	for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
		respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
	}
	if respDesc[0].ID == nil || respAsc[0].ID == nil {
		t.Fatalf("Response ID is nil in CDN Test federation")
	}
	if *respDesc[0].ID != *respAsc[0].ID {
		t.Errorf("CDN Federation responses are not equal after reversal: Asc: %d - Desc: %d", *respDesc[0].ID, *respAsc[0].ID)
	}
}

func UpdateTestCDNFederations(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, id := range fedIDs {
		opts.QueryParameters.Set("id", strconv.Itoa(id))
		fed, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
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
		resp, _, err := TOSession.UpdateCDNFederation(fed.Response[0], "cdn1", id, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot update federation by id: %v - alerts: %+v", err, resp.Alerts)
		}

		resp2, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
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
		data, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
		if err != nil {
			t.Errorf("could not get federations: %v - alerts: %+v", err, data.Alerts)
		}
	}
}

func GetTestCDNFederationsIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	fmtFutureTime := futureTime.Format(time.RFC1123)
	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, fmtFutureTime)
	for _, id := range fedIDs {
		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, reqInf, err := TOSession.GetCDNFederationsByName("cdn1", opts)
		if err != nil {
			t.Errorf("could not get federations: %v - alerts: %+v", err, data.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}

	pastTime := time.Now().AddDate(0, 0, -1)
	fmtPastTime := pastTime.Format(time.RFC1123)
	opts = client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, fmtPastTime)
	for _, id := range fedIDs {
		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, reqInf, err := TOSession.GetCDNFederationsByName("cdn1", opts)
		if err != nil {
			t.Errorf("could not get federations: %v - alerts: %+v", err, data.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
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
		resp, _, err := TOSession.DeleteCDNFederation("cdn1", id, opts)
		if err != nil {
			t.Errorf("cannot delete federation #%d: %v - alerts: %+v", id, err, resp.Alerts)
		}

		opts.QueryParameters.Set("id", strconv.Itoa(id))
		data, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
		if len(data.Response) != 0 {
			t.Error("expected federation to be deleted")
		}
		opts.QueryParameters.Del("id")
	}
	fedIDs = nil // reset the global variable for the next test
}

func GetTestPaginationSupportCdnFederation(t *testing.T) {

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
	if err != nil {
		t.Fatalf("cannot Get CDN Federation by name: %v - alerts: %+v", err, resp.Alerts)
	}
	cdnFederation := resp.Response
	if len(cdnFederation) < 3 {
		t.Fatalf("Need at least 3 CDN Federation by name in Traffic Ops to test pagination support, found: %d", len(cdnFederation))
	}

	opts.QueryParameters.Set("limit", "1")
	cdnFederationWithLimit, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
	if err != nil {
		t.Fatalf("cannot Get CDN Federation by name with Limit: %v - alerts: %+v", err, cdnFederationWithLimit.Alerts)
	}
	if !reflect.DeepEqual(cdnFederation[:1], cdnFederationWithLimit.Response) {
		t.Error("expected GET CDN Federation by name with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("offset", "1")
	cdnFederationWithOffset, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
	if err != nil {
		t.Fatalf("cannot Get CDN Federation by name with Limit and Offset: %v - alerts: %+v", err, cdnFederationWithOffset.Alerts)
	}
	if !reflect.DeepEqual(cdnFederation[1:2], cdnFederationWithOffset.Response) {
		t.Error("expected GET CDN Federation by name with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Del("offset")
	opts.QueryParameters.Set("page", "2")
	cdnFederationWithPage, _, err := TOSession.GetCDNFederationsByName("cdn1", opts)
	if err != nil {
		t.Fatalf("cannot Get CDN Federation by name with Limit and Page: %v - alerts: %+v", err, cdnFederationWithPage.Alerts)
	}
	if !reflect.DeepEqual(cdnFederation[1:2], cdnFederationWithPage.Response) {
		t.Error("expected GET CDN Federation by name with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, reqInf, err := TOSession.GetCDNFederationsByName("cdn1", opts)
	if err == nil {
		t.Error("expected GET CDN Federation by name to return an error when limit is not bigger than -1")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, reqInf, err = TOSession.GetCDNFederationsByName("cdn1", opts)
	if err == nil {
		t.Error("expected GET CDN Federation by name to return an error when offset is not a positive integer")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, reqInf, err = TOSession.GetCDNFederationsByName("cdn1", opts)
	if err == nil {
		t.Error("expected GET CDN Federation by name to return an error when page is not a positive integer")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}
