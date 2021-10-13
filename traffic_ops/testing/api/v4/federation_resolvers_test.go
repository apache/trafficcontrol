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
	"strconv"
	"testing"
	"time"
	"net/url"
	"reflect"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestFederationResolvers(t *testing.T) {
	WithObjs(t, []TCObj{Types, FederationResolvers}, func() {
		GetTestFederationResolversIMS(t)
		GetTestFederationResolvers(t)
		GetTestPaginationSupportFedResolver(t)
	})
}
func GetTestFederationResolversIMS(t *testing.T) {
	var tdlen = len(testData.FederationResolvers)
	if tdlen < 1 {
		t.Fatal("no federation resolvers test data")
	}

	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	resp, reqInf, err := TOSession.GetFederationResolvers(opts)
	if err != nil {
		t.Fatalf("could not get Federation resolvers: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestFederationResolvers(t *testing.T) {
	var tdlen = len(testData.FederationResolvers)
	if tdlen < 1 {
		t.Fatal("no federation resolvers test data")
	}

	frs, _, err := TOSession.GetFederationResolvers(client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting Federation Resolvers: %v - alerts: %+v", err, frs.Alerts)
	}
	if len(frs.Response) != tdlen {
		t.Fatalf("Wrong number of Federation Resolvers from GET, want %d got %d", tdlen, len(frs.Response))
	}

	var testFr = frs.Response[0]
	if testFr.ID == nil || testFr.Type == nil || testFr.IPAddress == nil {
		t.Fatalf("Malformed federation resolver: %+v", testFr)
	}

	_ = t.Run("Get Federation Resolvers by ID", getFRByIDTest(testFr))
	_ = t.Run("Get Federation Resolvers by IPAddress", getFRByIPTest(testFr))
	_ = t.Run("Get Federation Resolvers by Type", getFRByTypeTest(testFr))
}

func getFRByIDTest(testFr tc.FederationResolver) func(*testing.T) {
	return func(t *testing.T) {
		if testFr.ID == nil {
			t.Fatal("Federation Resolver has nil ID")
		}
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.FormatUint(uint64(*testFr.ID), 10))
		fr, _, err := TOSession.GetFederationResolvers(opts)
		if err != nil {
			t.Fatalf("Unexpected error getting Federation Resolver by ID %d: %v - alerts: %+v", *testFr.ID, err, fr.Alerts)
		}
		if len(fr.Response) != 1 {
			t.Fatalf("Expected exactly one Federation Resolver to exist with ID %d, found: %d", *testFr.ID, len(fr.Response))
		}

		cmpr(testFr, fr.Response[0], t)

	}
}

func getFRByIPTest(testFr tc.FederationResolver) func(*testing.T) {
	return func(t *testing.T) {
		if testFr.IPAddress == nil {
			t.Fatal("Federation Resolver has nil IP Address")
		}
		ip := *testFr.IPAddress
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("ipAddress", ip)
		frs, _, err := TOSession.GetFederationResolvers(opts)
		if err != nil {
			t.Fatalf("Unexpected error getting Federation Resolver by IP %s: %v - alerts: %+v", ip, err, frs.Alerts)
		}

		if len(frs.Response) != 1 {
			t.Fatalf("Expected exactly one Federation Resolver with IP address '%s', got: %d", ip, len(frs.Response))
		}

		cmpr(testFr, frs.Response[0], t)

	}
}

func getFRByTypeTest(testFr tc.FederationResolver) func(*testing.T) {
	return func(t *testing.T) {
		if testFr.Type == nil {
			t.Fatal("Federation Resolver has nil Type")
		}
		typ := *testFr.Type
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("type", typ)
		frs, _, err := TOSession.GetFederationResolvers(opts)
		if err != nil {
			t.Fatalf("Unexpected error getting Federation Resolvers by Type %s: %v - alerts: %+v", typ, err, frs.Alerts)
		}

		if len(frs.Response) < 1 {
			t.Errorf("Expected at least one Federation Resolver by Type '%s' to exist, but none did", typ)
		}

		for _, fr := range frs.Response {
			if fr.ID == nil {
				t.Error("Retrieved Federation Resolver has nil ID")
			}
			if fr.IPAddress == nil {
				t.Error("Retrieved Federation Resolver has nil IPAddress")
			}
			if fr.Type == nil {
				t.Error("Retrieved Federation Resolver has nil Type")
			} else if *fr.Type != typ {
				t.Errorf("Retrieved Federation Resolver has incorrect Type; want '%s', got '%s'", typ, *fr.Type)
			}
		}
	}
}

func cmpr(testFr, apiFr tc.FederationResolver, t *testing.T) {
	if apiFr.ID == nil {
		t.Error("Retrieved Federation Resolver has nil ID")
	} else if *apiFr.ID != *testFr.ID {
		t.Errorf("Retrieved Federation Resolver has incorrect ID; want %d, got %d", *testFr.ID, *apiFr.ID)
	}

	if apiFr.IPAddress == nil {
		t.Error("Retrieved Federation Resolver has nil IP address")
	} else if *apiFr.IPAddress != *testFr.IPAddress {
		t.Errorf("Retrieved Federation Resolver has incorrect IPAddress; want %s, got %s", *testFr.IPAddress, *apiFr.IPAddress)
	}

	if apiFr.Type == nil {
		t.Error("Retrieved Federation Resolver has nil Type")
	} else if *apiFr.Type != *testFr.Type {
		t.Errorf("Retrieved Federation Resolver has incorrect Type; want %s, got %s", *testFr.Type, *apiFr.Type)
	}
}

func CreateTestFederationResolvers(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, fr := range testData.FederationResolvers {
		if fr.Type == nil {
			t.Fatal("testData Federation Resolver has nil Type")
		}

		opts.QueryParameters.Set("name", *fr.Type)
		tid, _, err := TOSession.GetTypes(opts)
		if err != nil {
			t.Fatalf("Couldn't get an ID for Type '%s': %v - alerts: %+v", *fr.Type, err, tid.Alerts)
		}
		if len(tid.Response) != 1 {
			t.Fatalf("Expected exactly one Type by name %s, got %d", *fr.Type, len(tid.Response))
		}

		fr.TypeID = util.UIntPtr(uint(tid.Response[0].ID))

		alerts, _, err := TOSession.CreateFederationResolver(fr, client.RequestOptions{})
		if err != nil {
			t.Fatalf("failed to create Federation Resolver %+v: %v - alerts: %+v", fr, err, alerts.Alerts)
		}
		for _, a := range alerts.Alerts.Alerts {
			if a.Level == tc.ErrorLevel.String() {
				t.Errorf("Unexpected error-level alert creating a federation resolver: %s", a.Text)
			}
		}
	}

	var invalidFR tc.FederationResolver
	response, _, err := TOSession.CreateFederationResolver(invalidFR, client.RequestOptions{})
	if err == nil {
		t.Error("Expected an error creating a bad Federation Resolver, but didn't get one")
	}
	for _, a := range response.Alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success creating a bad Federation Resolver: %s", a.Text)
		}
	}

	invalidFR.TypeID = util.UIntPtr(1)
	invalidFR.IPAddress = util.StrPtr("not a valid IP address")
	response, _, err = TOSession.CreateFederationResolver(invalidFR, client.RequestOptions{})
	if err == nil {
		t.Error("Expected an error creating a bad Federation Resolver, but didn't get one")
	}
	for _, a := range response.Alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success creating a bad Federation Resolver: %s", a.Text)
		}
	}
}

func DeleteTestFederationResolvers(t *testing.T) {
	frs, _, err := TOSession.GetFederationResolvers(client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting Federation Resolvers: %v - alerts: %+v", err, frs.Alerts)
	}
	if len(frs.Response) < 1 {
		t.Fatal("Found no Federation Resolvers to delete")
	}
	for _, fr := range frs.Response {
		if fr.ID == nil {
			t.Fatalf("Malformed Federation Resolver: %+v", fr)
		}
		alerts, _, err := TOSession.DeleteFederationResolver(*fr.ID, client.RequestOptions{})
		if err != nil {
			t.Fatalf("failed to delete Federation Resolver %+v: %v - alerts: %+v", fr, err, alerts.Alerts)
		}
		for _, a := range alerts.Alerts.Alerts {
			if a.Level == tc.ErrorLevel.String() {
				t.Errorf("Unexpected error-level alert deleting a federation resolver: %s", a.Text)
			}
		}
	}

	alerts, _, err := TOSession.DeleteFederationResolver(0, client.RequestOptions{})
	if err == nil {
		t.Error("Expected an error deleting a non-existent Federation Resolver, but didn't get one")
	}
	for _, a := range alerts.Alerts.Alerts {
		if a.Level == tc.SuccessLevel.String() {
			t.Errorf("Unexpected success deleting a non-existent Federation Resolver: %s", a.Text)
		}
	}

}

func GetTestPaginationSupportFedResolver(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetFederationResolvers(opts)
	if err != nil {
		t.Fatalf("cannot Get FederationResolver: %v - alerts: %+v", err, resp.Alerts)
	}
	fedResolver := resp.Response
	if len(fedResolver) < 3 {
		t.Fatalf("Need at least 3 FederationResolver in Traffic Ops to test pagination support, found: %d", len(fedResolver))
	}

	opts.QueryParameters.Set("limit", "1")
	fedResolverWithLimit, _, err := TOSession.GetFederationResolvers(opts)
	if err != nil {
		t.Fatalf("cannot Get FederationResolver with Limit: %v - alerts: %+v", err, fedResolverWithLimit.Alerts)
	}
	if !reflect.DeepEqual(fedResolver[:1], fedResolverWithLimit.Response) {
		t.Error("expected GET FederationResolver with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("offset", "1")
	fedResolverWithOffset, _, err := TOSession.GetFederationResolvers(opts)
	if err != nil {
		t.Fatalf("cannot Get FederationResolver with Limit and Offset: %v - alerts: %+v", err, fedResolverWithOffset.Alerts)
	}
	if !reflect.DeepEqual(fedResolver[1:2], fedResolverWithOffset.Response) {
		t.Error("expected GET FederationResolver with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Del("offset")
	opts.QueryParameters.Set("page", "2")
	fedResolverWithPage, _, err := TOSession.GetFederationResolvers(opts)
	if err != nil {
		t.Fatalf("cannot Get FederationResolver with Limit and Page: %v - alerts: %+v", err, fedResolverWithPage.Alerts)
	}
	if !reflect.DeepEqual(fedResolver[1:2], fedResolverWithPage.Response) {
		t.Error("expected GET FederationResolver with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetFederationResolvers(opts)
	if err == nil {
		t.Error("expected GET FederationResolver to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET FederationResolver to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetFederationResolvers(opts)
	if err == nil {
		t.Error("expected GET FederationResolver to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET FederationResolver to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetFederationResolvers(opts)
	if err == nil {
		t.Error("expected GET FederationResolver to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET FederationResolver to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}
