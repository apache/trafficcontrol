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
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestOrigins(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Users, Topologies, DeliveryServices, Coordinates, Origins}, func() {
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestOrigins(t)
		UpdateTestOriginsWithHeaders(t, header)
		GetTestOrigins(t)
		NotFoundDeleteTest(t)
		OriginTenancyTest(t)
		VerifyPaginationSupport(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestOriginsWithHeaders(t, header)
	})
}

func UpdateTestOriginsWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Origins) < 1 {
		t.Fatal("Need at least one Origin to test updating Origins with an HTTP header")
	}
	firstOrigin := testData.Origins[0]
	if firstOrigin.Name == nil {
		t.Fatalf("couldn't get the name of test origin server")
	}

	// Retrieve the origin by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstOrigin.Name)
	resp, _, err := TOSession.GetOrigins(opts)
	if err != nil {
		t.Errorf("cannot get Origin '%s'': %v - alerts: %+v", *firstOrigin.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Origin to exist with name '%s', found: %d", *firstOrigin.Name, len(resp.Response))
	}

	remoteOrigin := resp.Response[0]
	if remoteOrigin.ID == nil {
		t.Fatal("couldn't get the ID of the response origin server")
	}
	updatedPort := 4321
	updatedFQDN := "updated.example.com"

	// update port and FQDN values on origin
	remoteOrigin.Port = &updatedPort
	remoteOrigin.FQDN = &updatedFQDN
	opts.QueryParameters.Del("name")
	opts.Header = header
	_, reqInf, err := TOSession.UpdateOrigin(*remoteOrigin.ID, remoteOrigin, opts)
	if err == nil {
		t.Errorf("Expected error about precondition failed, but got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func CreateTestOrigins(t *testing.T) {
	// loop through origins, assign FKs and create
	for _, origin := range testData.Origins {
		resp, _, err := TOSession.CreateOrigin(origin, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Origins: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func NotFoundDeleteTest(t *testing.T) {
	resp, _, err := TOSession.DeleteOrigin(2020, client.RequestOptions{})
	if err == nil {
		t.Fatal("deleting origin with what should be a non-existent id - expected: error, actual: nil error")
	}

	found := false
	for _, alert := range resp.Alerts {
		if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, "not found") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("deleted origin with what should be a non-existent id - expected: 'not found' error-level alert, actual: %v - alerts: %+v", err, resp.Alerts)
	}
}

func GetTestOrigins(t *testing.T) {
	resp, _, err := TOSession.GetOrigins(client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot get Origins: %v - alerts: %+v", err, resp.Alerts)
	}

	opts := client.NewRequestOptions()
	for _, origin := range testData.Origins {
		opts.QueryParameters.Set("name", *origin.Name)
		resp, _, err := TOSession.GetOrigins(opts)
		if err != nil {
			t.Errorf("cannot get Origin by name: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func UpdateTestOrigins(t *testing.T) {
	if len(testData.Origins) < 1 {
		t.Fatal("Need at least one Origin to test updating Origins")
	}
	firstOrigin := testData.Origins[0]
	if firstOrigin.Name == nil {
		t.Fatal("Found an Origin in the testing data with null or undefined name")
	}
	foName := *firstOrigin.Name

	// Retrieve the origin by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", foName)
	resp, _, err := TOSession.GetOrigins(opts)
	if err != nil {
		t.Errorf("cannot get Origin '%s': %v - alerts: %+v", foName, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Origin to exist with name '%s', found: %d", foName, len(resp.Response))
	}

	remoteOrigin := resp.Response[0]
	if remoteOrigin.ID == nil {
		t.Fatalf("Traffic Ops responded with a representation of Origin '%s' that had null or undefined ID", foName)
	}
	updatedPort := 4321
	updatedFQDN := "updated.example.com"

	// update port and FQDN values on origin
	remoteOrigin.Port = &updatedPort
	remoteOrigin.FQDN = &updatedFQDN
	updResp, _, err := TOSession.UpdateOrigin(*remoteOrigin.ID, remoteOrigin, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Origin '%s' (#%d): %v - %v", foName, *remoteOrigin.ID, err, updResp.Alerts)
	}

	// Retrieve the origin to check port and FQDN values were updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(*remoteOrigin.ID))
	resp, _, err = TOSession.GetOrigins(opts)
	if err != nil {
		t.Errorf("cannot get Origin #%d ('%s'): %v - alerts: %+v", *remoteOrigin.ID, foName, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Origin to exist with ID %d, found: %d", *remoteOrigin.ID, len(resp.Response))
	}
	respOrigin := resp.Response[0]

	if respOrigin.Port == nil {
		t.Errorf("results do not match, actual: null or undefined - expected: %d", updatedPort)
	} else if *respOrigin.Port != updatedPort {
		t.Errorf("results do not match actual: %d, expected: %d", *respOrigin.Port, updatedPort)
	}
	if respOrigin.FQDN == nil {
		t.Errorf("results do not match, actual: null or undefined, expected: '%s'", updatedFQDN)
	} else if *respOrigin.FQDN != updatedFQDN {
		t.Errorf("results do not match actual: %s, expected: %s", *respOrigin.FQDN, updatedFQDN)
	}
}

func OriginTenancyTest(t *testing.T) {
	origins, _, err := TOSession.GetOrigins(client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot get Origins: %v - alerts: %+v", err, origins.Alerts)
	}
	if len(origins.Response) < 1 {
		t.Fatal("Need at least one Origin to exist in Traffic Ops to test Tenancy for Origins")
	}
	// This ID check specifically needs to be a fatal condition, despite also being an error below,
	// because we explicitly dereference the ID of the 0th Origin in this slice later on.
	if origins.Response[0].ID == nil || origins.Response[0].Name == nil {
		t.Fatal("Traffic Ops returned a representation for an Origin with null or undefined ID and/or Name")
	}

	var tenant3Origin tc.Origin
	foundTenant3Origin := false
	for _, o := range origins.Response {
		if o.FQDN == nil || o.ID == nil {
			t.Error("Traffic Ops responded with a representation of an Origin with null or undefined FQDN and/or ID")
			continue
		}
		if *o.FQDN == "origin.ds3.example.net" {
			tenant3Origin = o
			foundTenant3Origin = true
		}
	}
	if !foundTenant3Origin {
		t.Error("expected to find origin with tenant 'tenant3' and fqdn 'origin.ds3.example.net'")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant4TOClient, _, err := client.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v3-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with tenant4user: %v", err)
	}

	originsReadableByTenant4, _, err := tenant4TOClient.GetOrigins(client.RequestOptions{})
	if err != nil {
		t.Errorf("tenant4user cannot get Origins: %v - alerts: %+v", err, originsReadableByTenant4.Alerts)
	}

	// assert that tenant4user cannot read origins outside of its tenant
	for _, origin := range originsReadableByTenant4.Response {
		if origin.FQDN == nil {
			t.Error("Traffic Ops returned a representation of an Origin with null or undefined FQDN")
		} else if *origin.FQDN == "origin.ds3.example.net" {
			t.Error("expected tenant4 to be unable to read origins from tenant 3")
		}
	}

	// assert that tenant4user cannot update tenant3user's origin
	if _, _, err = tenant4TOClient.UpdateOrigin(*tenant3Origin.ID, tenant3Origin, client.RequestOptions{}); err == nil {
		t.Error("expected tenant4user to be unable to update tenant3's origin")
	}

	// assert that tenant4user cannot delete an origin outside of its tenant
	if _, _, err = tenant4TOClient.DeleteOrigin(*origins.Response[0].ID, client.RequestOptions{}); err == nil {
		t.Errorf("expected tenant4user to be unable to delete an origin outside of its tenant (origin %s)", *origins.Response[0].Name)
	}

	// assert that tenant4user cannot create origins outside of its tenant
	tenant3Origin.FQDN = util.StrPtr("origin.tenancy.test.example.com")
	if _, _, err = tenant4TOClient.CreateOrigin(tenant3Origin, client.RequestOptions{}); err == nil {
		t.Error("expected tenant4user to be unable to create an origin outside of its tenant")
	}
}

func alertsHaveError(alerts []tc.Alert, err string) bool {
	for _, alert := range alerts {
		if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, err) {
			return true
		}
	}
	return false
}

func VerifyPaginationSupport(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetOrigins(opts)
	if err != nil {
		t.Fatalf("cannot get Origins: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) < 3 {
		t.Fatalf("Need at least 3 Origins in Traffic Ops to test pagination")
	}
	origins := resp.Response

	opts.QueryParameters.Set("limit", "1")
	originsWithLimit, _, err := TOSession.GetOrigins(opts)
	if err != nil {
		t.Fatalf("cannot get Origins with limit: %v - alerts: %+v", err, originsWithLimit.Alerts)
	}
	if !reflect.DeepEqual(origins[:1], originsWithLimit.Response) {
		t.Error("expected GET origins with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("offset", "1")
	originsWithOffset, _, err := TOSession.GetOrigins(opts)
	if err != nil {
		t.Fatalf("cannot get Origins with offset: %v - alerts: %+v", err, originsWithOffset.Alerts)
	}
	if !reflect.DeepEqual(origins[1:2], originsWithOffset.Response) {
		t.Error("expected GET origins with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Del("offset")
	opts.QueryParameters.Set("page", "2")
	originsWithPage, _, err := TOSession.GetOrigins(opts)
	if err != nil {
		t.Fatalf("cannot get Origins with page: %v - alerts: %+v", err, originsWithPage.Alerts)
	}
	if !reflect.DeepEqual(origins[1:2], originsWithPage.Response) {
		t.Error("expected GET origins with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters.Del("page")
	opts.QueryParameters.Del("orderby")
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetOrigins(opts)
	if err == nil {
		t.Error("expected GET origins to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET origins to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetOrigins(opts)
	if err == nil {
		t.Error("expected GET origins to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET origins to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	opts.QueryParameters.Del("offset")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetOrigins(opts)
	if err == nil {
		t.Error("expected GET origins to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET origins to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}

func DeleteTestOrigins(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, origin := range testData.Origins {
		if origin.Name == nil {
			t.Error("Found an Origin in the testing data with null or undefined name")
			continue
		}

		opts.QueryParameters.Set("name", *origin.Name)
		resp, _, err := TOSession.GetOrigins(opts)
		if err != nil {
			t.Errorf("cannot get Origin '%s': %v - alerts: %+v", *origin.Name, err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			respOrigin := resp.Response[0]
			if respOrigin.ID == nil {
				t.Error("Traffic Ops returned a representation for an Origin that has null or undefined ID")
				continue
			}
			delResp, _, err := TOSession.DeleteOrigin(*respOrigin.ID, client.RequestOptions{})
			if err != nil {
				t.Errorf("cannot DELETE Origin by ID: %v - %v", err, delResp)
			}

			// Retrieve the Origin to see if it got deleted
			org, _, err := TOSession.GetOrigins(opts)
			if err != nil {
				t.Errorf("error fetching Origin '%s' after deletion: %v - alerts: %+v", *origin.Name, err, org.Alerts)
			}
			if len(org.Response) > 0 {
				t.Errorf("expected Origin '%s' to be deleted, but it was found in Traffic Ops", *origin.Name)
			}
		}
	}
}
