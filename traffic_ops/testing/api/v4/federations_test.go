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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestFederations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, CDNFederations}, func() {
		PostDeleteTestFederationsDeliveryServices(t)
		GetTestFederations(t)
		GetTestFederationsIMS(t)
		AddFederationResolversForCurrentUserTest(t)
		ReplaceFederationResolversForCurrentUserTest(t)
		RemoveFederationResolversForCurrentUserTest(t)
		GetTestPaginationSupportFedDeliveryServices(t)
		SortTestFederationDs(t)
		SortTestFederationsDsDesc(t)
	})
}

func GetTestFederationsIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	fmtFutureTime := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, fmtFutureTime)
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	resp, reqInf, err := TOSession.AllFederations(opts)
	if err != nil {
		t.Fatalf("No error expected, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}

	pastTime := time.Now().AddDate(0, 0, -1)
	fmtPastTime := pastTime.Format(time.RFC1123)

	opts = client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, fmtPastTime)
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	resp, reqInf, err = TOSession.AllFederations(opts)
	if err != nil {
		t.Fatalf("No error expected, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestFederations(t *testing.T) {
	if len(testData.Federations) == 0 {
		t.Error("no federations test data")
	}

	feds, _, err := TOSession.AllFederations(client.RequestOptions{})
	if err != nil {
		t.Errorf("getting federations: %v - alerts: %+v", err, feds.Alerts)
	}

	if len(feds.Response) < 1 {
		t.Errorf("expected atleast 1 federation, but got none")
	}
	fed := feds.Response[0]

	if len(fed.Mappings) < 1 {
		t.Fatal("federation mappings expected > 1, actual: 0")
	}

	mapping := fed.Mappings[0]
	if mapping.CName == nil {
		t.Fatal("federation mapping expected cname, actual: nil")
	}
	if mapping.TTL == nil {
		t.Fatal("federation mapping expected ttl, actual: nil")
	}

	matched := false
	for _, testFed := range testData.Federations {
		if testFed.CName == nil {
			t.Error("test federation missing cname!")
			continue
		}
		if testFed.TTL == nil {
			t.Error("test federation missing ttl!")
			continue
		}

		if *mapping.CName != *testFed.CName {
			continue
		}
		matched = true

		if *mapping.TTL != *testFed.TTL {
			t.Errorf("federation mapping ttl expected: %v, actual: %v", *testFed.TTL, *mapping.TTL)
		}
	}
	if !matched {
		t.Errorf("federation mapping expected to match test data, actual: cname %v not in test data", *mapping.CName)
	}
}

func createFederationToDeliveryServiceAssociation() (int, tc.DeliveryServiceV4, tc.DeliveryServiceV4, error) {
	var ds tc.DeliveryServiceV4
	var ds1 tc.DeliveryServiceV4
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		return -1, ds, ds1, fmt.Errorf("cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 2 {
		return -1, ds, ds1, errors.New("no delivery services, must have at least 2 Delivery Services to test federations deliveryservices")
	}
	ds = dses.Response[0]
	ds1 = dses.Response[1]
	if ds.ID == nil || ds1.ID == nil {
		return -1, ds, ds1, errors.New("Traffic Ops returned at least one representation of a Delivery Service that had a null or undefined ID")
	}

	if len(fedIDs) == 0 {
		return -1, ds, ds1, errors.New("no federations, must have at least 1 federation to test federations deliveryservices")
	}
	fedID := fedIDs[0]

	alerts, _, err := TOSession.CreateFederationDeliveryServices(fedID, []int{*ds.ID, *ds1.ID}, true, client.RequestOptions{})
	if err != nil {
		err = fmt.Errorf("creating federations delivery services: %v - alerts: %+v", err, alerts.Alerts)
	}

	return fedID, ds, ds1, err

}

func createFederationToMultipleDeliveryServiceAssociation() (int, tc.DeliveryServiceV4, tc.DeliveryServiceV4, tc.DeliveryServiceV4, error) {
	var ds tc.DeliveryServiceV4
	var ds1 tc.DeliveryServiceV4
	var ds2 tc.DeliveryServiceV4
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	if err != nil {
		return -1, ds, ds1, ds2, fmt.Errorf("cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)
	}
	if len(dses.Response) < 3 {
		return -1, ds, ds1, ds2, errors.New("no delivery services, must have at least 3 Delivery Services to test federations deliveryservices")
	}
	ds = dses.Response[0]
	ds1 = dses.Response[1]
	ds2 = dses.Response[2]
	if ds.ID == nil || ds1.ID == nil || ds2.ID == nil {
		return -1, ds, ds1, ds2, errors.New("Traffic Ops returned at least one representation of a Delivery Service that had a null or undefined ID")
	}

	if len(fedIDs) == 0 {
		return -1, ds, ds1, ds2, errors.New("no federations, must have at least 1 federation to test federations deliveryservices")
	}
	fedID := fedIDs[1]

	alerts, _, err := TOSession.CreateFederationDeliveryServices(fedID, []int{*ds.ID, *ds1.ID, *ds2.ID}, true, client.RequestOptions{})
	if err != nil {
		err = fmt.Errorf("creating federations delivery services: %v - alerts: %+v", err, alerts.Alerts)
	}

	return fedID, ds, ds1, ds2, err

}

func PostDeleteTestFederationsDeliveryServices(t *testing.T) {
	fedID, ds, ds1, err := createFederationToDeliveryServiceAssociation()
	if err != nil {
		t.Fatal(err.Error())
	}

	// Test get created Federation Delivery Services
	fedDSes, _, err := TOSession.GetFederationDeliveryServices(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Federation #%d Delivery Services: %v - alerts: %+v", fedID, err, fedDSes.Alerts)
	}
	if len(fedDSes.Response) != 2 {
		t.Fatalf("two Federation DeliveryService expected for Federation %d, %d was returned", fedID, len(fedDSes.Response))
	}

	// Delete one of the Delivery Services from the Federation
	alerts, _, err := TOSession.DeleteFederationDeliveryService(fedID, *ds.ID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot Delete Federation #%d Delivery Service #%d: %v - alerts: %+v", fedID, ds.ID, err, alerts.Alerts)
	}

	// Make sure it is deleted

	// Test get created Federation Delivery Services
	fedDSes, _, err = TOSession.GetFederationDeliveryServices(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot get Federation Delivery Services: %v - alerts: %+v", err, fedDSes.Alerts)
	}
	if len(fedDSes.Response) != 1 {
		t.Fatalf("one Federation Delivery Service expected for Federation #%d, %d was returned", fedID, len(fedDSes.Response))
	}

	// Attempt to delete the last one which should fail as you cannot remove the last
	_, _, err = TOSession.DeleteFederationDeliveryService(fedID, *ds1.ID, client.RequestOptions{})
	if err == nil {
		t.Fatal("expected to receive error from attempting to delete last Delivery Service from a Federation")
	}
}

func RemoveFederationResolversForCurrentUserTest(t *testing.T) {
	if len(testData.Federations) < 1 {
		t.Fatal("No test Federations, deleting resolvers cannot be tested!")
	}

	alerts, _, err := TOSession.DeleteFederationResolverMappingsForCurrentUser(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error deleting Federation Resolvers for current user: %v - alerts: %+v", err, alerts)
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.ErrorLevel.String() {
			t.Errorf("Unexpected error-level alert from deleting Federation Resolvers for current user: %s", a.Text)
		}
	}

	// Now try deleting Federation Resolvers when there are none.
	_, _, err = TOSession.DeleteFederationResolverMappingsForCurrentUser(client.RequestOptions{})
	if err != nil {
		t.Logf("Received expected error deleting Federation Resolvers for current user: %v", err)
	} else {
		t.Error("Expected an error deleting zero Federation Resolvers, but didn't get one.")
	}
}

func AddFederationResolversForCurrentUserTest(t *testing.T) {
	fedID, ds, ds1, err := createFederationToDeliveryServiceAssociation()
	if err != nil {
		t.Fatal(err.Error())
	}

	// need to assign myself the federation to set its mappings
	me, _, err := TOSession.GetUserCurrent(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Couldn't figure out who I am: %v - alerts: %+v", err, me.Alerts)
	}
	if me.Response.ID == nil {
		t.Fatal("Current user has no ID, cannot continue.")
	}

	alerts, _, err := TOSession.CreateFederationUsers(fedID, []int{*me.Response.ID}, false, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Failed to assign Federation to current user: %v - alerts: %+v", err, alerts.Alerts)
	}

	mappings := tc.DeliveryServiceFederationResolverMappingRequest{
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: *ds.XMLID,
			Mappings: tc.ResolverMapping{
				Resolve4: []string{"0.0.0.0"},
				Resolve6: []string{"::1"},
			},
		},
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: *ds1.XMLID,
			Mappings: tc.ResolverMapping{
				Resolve4: []string{"1.2.3.4/28"},
				Resolve6: []string{"1234::/110"},
			},
		},
	}

	alerts, _, err = TOSession.AddFederationResolverMappingsForCurrentUser(mappings, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error adding Federation Resolver mappings for the current user: %v - alerts: %+v", err, alerts.Alerts)
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.ErrorLevel.String() {
			t.Errorf("Unexpected error-level alert from adding Federation Resolver mappings for the current user: %s", a.Text)
		}
	}

	mappings = tc.DeliveryServiceFederationResolverMappingRequest{
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: "aoeuhtns",
			Mappings: tc.ResolverMapping{
				Resolve4: []string{},
				Resolve6: []string{"dead::beef", "f1d0::f00d/82"},
			},
		},
	}

	alerts, _, err = TOSession.AddFederationResolverMappingsForCurrentUser(mappings, client.RequestOptions{})
	if err == nil {
		t.Fatal("Expected error adding Federation Resolver mappings for the current user, but didn't get one")
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success from adding Federation Resolver mappings for the current user: %s", a.Text)
		}
	}
}

func ReplaceFederationResolversForCurrentUserTest(t *testing.T) {
	fedID, ds, ds1, err := createFederationToDeliveryServiceAssociation()
	if err != nil {
		t.Fatal(err.Error())
	}

	// need to assign myself the federation to set its mappings
	me, _, err := TOSession.GetUserCurrent(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Couldn't figure out who I am: %v - alerts: %+v", err, me.Alerts)
	}
	if me.Response.ID == nil {
		t.Fatal("Current user has no ID, cannot continue.")
	}

	fedUsers, _, err := TOSession.GetFederationUsers(fedID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("unexpected error getting federation users: %v", err)
	}
	foundFedUser := false
	for _, fedUser := range fedUsers.Response {
		if *fedUser.ID == *me.Response.ID {
			foundFedUser = true
			break
		}
	}
	if !foundFedUser {
		alerts, _, err := TOSession.CreateFederationUsers(fedID, []int{*me.Response.ID}, false, client.RequestOptions{})
		if err != nil {
			t.Fatalf("Failed to assign Federation to current user: %v - alerts: %+v", err, alerts.Alerts)
		}
	}
	expectedResolve4 := []string{"192.0.2.0/25", "192.0.2.128/25"}
	expectedResolve6 := []string{"2001:db8::/33", "2001:db8:8000::/33"}
	sort.Strings(expectedResolve4)
	sort.Strings(expectedResolve6)

	mappings := tc.DeliveryServiceFederationResolverMappingRequest{
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: *ds.XMLID,
			Mappings: tc.ResolverMapping{
				Resolve4: expectedResolve4,
				Resolve6: expectedResolve6,
			},
		},
		// for the purpose of this test, it's important that at least two different mappings have the same resolvers
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: *ds1.XMLID,
			Mappings: tc.ResolverMapping{
				Resolve4: expectedResolve4,
				Resolve6: expectedResolve6,
			},
		},
	}

	alerts, _, err := TOSession.ReplaceFederationResolverMappingsForCurrentUser(mappings, client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error replacing Federation Resolver mappings for the current user: %v - alerts: %+v", err, alerts.Alerts)
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.ErrorLevel.String() {
			t.Errorf("Unexpected error-level alert from replacing Federation Resolver mappings for the current user: %s", a.Text)
		}
	}

	feds, _, err := TOSession.Federations(client.RequestOptions{})
	if err != nil {
		t.Fatalf("unexpected error getting federations: %v", err)
	}
	for _, fed := range feds.Response {
		if fed.DeliveryService.String() == *ds.XMLID || fed.DeliveryService.String() == *ds1.XMLID {
			if len(fed.Mappings) != 1 {
				t.Fatalf("expected 1 mapping, got %d", len(fed.Mappings))
			}
			sort.Strings(fed.Mappings[0].Resolve4)
			sort.Strings(fed.Mappings[0].Resolve6)
			if !reflect.DeepEqual(expectedResolve4, fed.Mappings[0].Resolve4) {
				t.Errorf("checking federation resolver mappings, expected: %+v, actual: %+v", expectedResolve4, fed.Mappings[0].Resolve4)
			}
			if !reflect.DeepEqual(expectedResolve6, fed.Mappings[0].Resolve6) {
				t.Errorf("checking federation resolver mappings, expected: %+v, actual: %+v", expectedResolve6, fed.Mappings[0].Resolve6)
			}
		}
	}

	mappings = tc.DeliveryServiceFederationResolverMappingRequest{
		tc.DeliveryServiceFederationResolverMapping{
			DeliveryService: "aoeuhtns",
			Mappings: tc.ResolverMapping{
				Resolve4: []string{},
				Resolve6: []string{"dead::beef", "f1d0::f00d/82"},
			},
		},
	}

	alerts, _, err = TOSession.ReplaceFederationResolverMappingsForCurrentUser(mappings, client.RequestOptions{})
	if err == nil {
		t.Fatal("Expected error replacing Federation Resolver mappings for the current user, but didn't get one")
	}
	for _, a := range alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success from replacing Federation Resolver mappings for the current user: %s", a.Text)
		}
	}
}

func GetTestPaginationSupportFedDeliveryServices(t *testing.T) {

	fedID, _, _, _, err := createFederationToMultipleDeliveryServiceAssociation()
	if err != nil {
		t.Fatal(err.Error())
	}
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "dsID")
	resp, _, err := TOSession.GetFederationDeliveryServices(fedID, opts)
	if err != nil {
		t.Fatalf("cannot get Federation #%d Delivery Services: %v - alerts: %+v", fedID, err, resp.Alerts)
	}
	fedDs := resp.Response
	if len(fedDs) < 3 {
		t.Fatalf("Need at least 3 Federation Delivery services in Traffic Ops to test pagination support, found: %d", len(fedDs))
	}

	opts.QueryParameters.Set("limit", "1")
	fedDsWithLimit, _, err := TOSession.GetFederationDeliveryServices(fedID, opts)
	if err != nil {
		t.Fatalf("cannot Get Federation deliveryservices with Limit: %v - alerts: %+v", err, fedDsWithLimit.Alerts)
	}
	if !reflect.DeepEqual(fedDs[:1], fedDsWithLimit.Response) {
		t.Error("expected GET Federation deliveryservices with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("offset", "1")
	fedDsWithOffset, _, err := TOSession.GetFederationDeliveryServices(fedID, opts)
	if err != nil {
		t.Fatalf("cannot Get Federation deliveryservices with Limit and Offset: %v - alerts: %+v", err, fedDsWithOffset.Alerts)
	}
	if !reflect.DeepEqual(fedDs[1:2], fedDsWithOffset.Response) {
		t.Error("expected GET Federation deliveryservices with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Del("offset")
	opts.QueryParameters.Set("page", "2")
	fedDsWithPage, _, err := TOSession.GetFederationDeliveryServices(fedID, opts)
	if err != nil {
		t.Fatalf("cannot Get Federation deliveryservices with Limit and Page: %v - alerts: %+v", err, fedDsWithPage.Alerts)
	}
	if !reflect.DeepEqual(fedDs[1:2], fedDsWithPage.Response) {
		t.Error("expected GET Federation deliveryservices with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, reqInf, err := TOSession.GetFederationDeliveryServices(fedID, opts)
	if err == nil {
		t.Error("expected GET Federation deliveryservices to return an error when limit is not bigger than -1")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, reqInf, err = TOSession.GetFederationDeliveryServices(fedID, opts)
	if err == nil {
		t.Error("expected GET Federation deliveryservices to return an error when offset is not a positive integer")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, reqInf, err = TOSession.GetFederationDeliveryServices(fedID, opts)
	if err == nil {
		t.Error("expected GET Federation deliveryservices to return an error when page is not a positive integer")
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
}

func SortTestFederationDs(t *testing.T) {
	var sortedList []int
	if len(fedIDs) == 0 {
		t.Fatalf("no federations, must have at least 1 federation to test federations deliveryservices")
	}
	fedID := fedIDs[1]
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "dsID")
	resp, _, err := TOSession.GetFederationDeliveryServices(fedID, opts)
	if err != nil {
		t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
	}
	for _, fedDs := range resp.Response {
		if fedDs.ID == nil {
			t.Fatalf("Federation Deliveryservices ID is nil, so sorting can't be tested")
		}
		sortedList = append(sortedList, *fedDs.ID)
	}
	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their DS ID: %v", sortedList)
	}
}

func SortTestFederationsDsDesc(t *testing.T) {

	if len(fedIDs) == 0 {
		t.Fatalf("no federations, must have at least 1 federation to test federations deliveryservices")
	}
	fedID := fedIDs[1]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "dsID")
	resp, _, err := TOSession.GetFederationDeliveryServices(fedID, opts)
	if err != nil {
		t.Fatalf("Expected no error, but got error in Federation DS default ordering %v - alerts: %+v", err, resp.Alerts)
	}
	respAsc := resp.Response
	if len(respAsc) < 1 {
		t.Fatal("Need at least one Federation DS in Traffic Ops to test Federation DS sort ordering")
	}
	opts.QueryParameters.Set("sortOrder", "desc")
	resp, _, err = TOSession.GetFederationDeliveryServices(fedID, opts)
	if err != nil {
		t.Errorf("Expected no error, but got error in Federation DS with Descending ordering: %v - alerts: %+v", err, resp.Alerts)
	}
	respDesc := resp.Response
	if len(respDesc) < 1 {
		t.Fatal("Need at least one Federation DS in Traffic Ops to test Federation DS sort ordering")
	}
	if len(respAsc) != len(respDesc) {
		t.Fatalf("Traffic Ops returned %d Federation DS using default sort order, but %d Federation DS when sort order was explicitly set to descending", len(respAsc), len(respDesc))
	}
	for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
		respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
	}
	if respDesc[0].ID == nil || respAsc[0].ID == nil {
		t.Fatalf("Response ID is nil in federation deliveryservices")
	}
	if *respDesc[0].ID != *respAsc[0].ID {
		t.Errorf("Federation DS responses are not equal after reversal: Asc: %d - Desc: %d", *respDesc[0].ID, *respAsc[0].ID)
	}
}
