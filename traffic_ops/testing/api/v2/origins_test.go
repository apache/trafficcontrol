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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/lib/go-util"
	toclient "github.com/apache/trafficcontrol/v6/traffic_ops/v2-client"
)

func TestOrigins(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Users, DeliveryServices, Coordinates, Origins}, func() {
		UpdateTestOrigins(t)
		GetTestOrigins(t)
		NotFoundDeleteTest(t)
		OriginTenancyTest(t)
		VerifyPaginationSupport(t)
	})
}

func CreateTestOrigins(t *testing.T) {
	// loop through origins, assign FKs and create
	for _, origin := range testData.Origins {
		_, _, err := TOSession.CreateOrigin(origin)
		if err != nil {
			t.Errorf("could not CREATE origins: %v", err)
		}
	}
}

func NotFoundDeleteTest(t *testing.T) {
	_, _, err := TOSession.DeleteOriginByID(2020)
	if !strings.Contains(err.Error(), "not found") {
		t.Error("deleted origin with what should be a non-existent id")
	}
}

func GetTestOrigins(t *testing.T) {
	_, _, err := TOSession.GetOrigins()
	if err != nil {
		t.Errorf("cannot GET origins: %v", err)
	}

	for _, origin := range testData.Origins {
		resp, _, err := TOSession.GetOriginByName(*origin.Name)
		if err != nil {
			t.Errorf("cannot GET Origin by name: %v - %v", err, resp)
		}
	}
}

func UpdateTestOrigins(t *testing.T) {
	firstOrigin := testData.Origins[0]
	// Retrieve the origin by name so we can get the id for the Update
	resp, _, err := TOSession.GetOriginByName(*firstOrigin.Name)
	if err != nil {
		t.Errorf("cannot GET origin by name: %v - %v", *firstOrigin.Name, err)
	}
	remoteOrigin := resp[0]
	updatedPort := 4321
	updatedFQDN := "updated.example.com"

	// update port and FQDN values on origin
	remoteOrigin.Port = &updatedPort
	remoteOrigin.FQDN = &updatedFQDN
	updResp, _, err := TOSession.UpdateOriginByID(*remoteOrigin.ID, remoteOrigin)
	if err != nil {
		t.Errorf("cannot UPDATE Origin by name: %v - %v", err, updResp.Alerts)
	}

	// Retrieve the origin to check port and FQDN values were updated
	resp, _, err = TOSession.GetOriginByID(*remoteOrigin.ID)
	if err != nil {
		t.Errorf("cannot GET Origin by ID: %v - %v", *remoteOrigin.Name, err)
	}

	respOrigin := resp[0]
	if *respOrigin.Port != updatedPort {
		t.Errorf("results do not match actual: %d, expected: %d", *respOrigin.Port, updatedPort)
	}
	if *respOrigin.FQDN != updatedFQDN {
		t.Errorf("results do not match actual: %s, expected: %s", *respOrigin.FQDN, updatedFQDN)
	}
}

func OriginTenancyTest(t *testing.T) {
	origins, _, err := TOSession.GetOrigins()
	if err != nil {
		t.Errorf("cannot GET origins: %v", err)
	}
	tenant3Origin := tc.Origin{}
	foundTenant3Origin := false
	for _, o := range origins {
		if *o.FQDN == "origin.ds3.example.net" {
			tenant3Origin = o
			foundTenant3Origin = true
		}
	}
	if !foundTenant3Origin {
		t.Error("expected to find origin with tenant 'tenant3' and fqdn 'origin.ds3.example.net'")
	}

	toReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)
	tenant4TOClient, _, err := toclient.LoginWithAgent(TOSession.URL, "tenant4user", "pa$$word", true, "to-api-v2-client-tests/tenant4user", true, toReqTimeout)
	if err != nil {
		t.Fatalf("failed to log in with tenant4user: %v", err.Error())
	}

	originsReadableByTenant4, _, err := tenant4TOClient.GetOrigins()
	if err != nil {
		t.Error("tenant4user cannot GET origins")
	}

	// assert that tenant4user cannot read origins outside of its tenant
	for _, origin := range originsReadableByTenant4 {
		if *origin.FQDN == "origin.ds3.example.net" {
			t.Error("expected tenant4 to be unable to read origins from tenant 3")
		}
	}

	// assert that tenant4user cannot update tenant3user's origin
	if _, _, err = tenant4TOClient.UpdateOriginByID(*tenant3Origin.ID, tenant3Origin); err == nil {
		t.Error("expected tenant4user to be unable to update tenant3's origin")
	}

	// assert that tenant4user cannot delete an origin outside of its tenant
	if _, _, err = tenant4TOClient.DeleteOriginByID(*origins[0].ID); err == nil {
		t.Errorf("expected tenant4user to be unable to delete an origin outside of its tenant (origin %s)", *origins[0].Name)
	}

	// assert that tenant4user cannot create origins outside of its tenant
	tenant3Origin.FQDN = util.StrPtr("origin.tenancy.test.example.com")
	if _, _, err = tenant4TOClient.CreateOrigin(tenant3Origin); err == nil {
		t.Error("expected tenant4user to be unable to create an origin outside of its tenant")
	}
}

func VerifyPaginationSupport(t *testing.T) {
	origins, _, err := TOSession.GetOriginsByQueryParams("?orderby=id")
	if err != nil {
		t.Fatalf("cannot GET origins: %v", err)
	}

	originsWithLimit, _, err := TOSession.GetOriginsByQueryParams("?orderby=id&limit=1")
	if !reflect.DeepEqual(origins[:1], originsWithLimit) {
		t.Error("expected GET origins with limit = 1 to return first result")
	}

	originsWithOffset, _, err := TOSession.GetOriginsByQueryParams("?orderby=id&limit=1&offset=1")
	if !reflect.DeepEqual(origins[1:2], originsWithOffset) {
		t.Error("expected GET origins with limit = 1, offset = 1 to return second result")
	}

	originsWithPage, _, err := TOSession.GetOriginsByQueryParams("?orderby=id&limit=1&page=2")
	if !reflect.DeepEqual(origins[1:2], originsWithPage) {
		t.Error("expected GET origins with limit = 1, page = 2 to return second result")
	}

	_, _, err = TOSession.GetOriginsByQueryParams("?limit=-2")
	if err == nil {
		t.Error("expected GET origins to return an error when limit is not bigger than -1")
	} else if !strings.Contains(err.Error(), "must be bigger than -1") {
		t.Errorf("expected GET origins to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}
	_, _, err = TOSession.GetOriginsByQueryParams("?limit=1&offset=0")
	if err == nil {
		t.Error("expected GET origins to return an error when offset is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET origins to return an error for offset is not a positive integer, actual error: " + err.Error())
	}
	_, _, err = TOSession.GetOriginsByQueryParams("?limit=1&page=0")
	if err == nil {
		t.Error("expected GET origins to return an error when page is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET origins to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}

func DeleteTestOrigins(t *testing.T) {
	for _, origin := range testData.Origins {
		resp, _, err := TOSession.GetOriginByName(*origin.Name)
		if err != nil {
			t.Errorf("cannot GET Origin by name: %v - %v", *origin.Name, err)
		}
		if len(resp) > 0 {
			respOrigin := resp[0]

			delResp, _, err := TOSession.DeleteOriginByID(*respOrigin.ID)
			if err != nil {
				t.Errorf("cannot DELETE Origin by ID: %v - %v", err, delResp)
			}

			// Retrieve the Origin to see if it got deleted
			org, _, err := TOSession.GetOriginByName(*origin.Name)
			if err != nil {
				t.Errorf("error deleting Origin name: %s", err.Error())
			}
			if len(org) > 0 {
				t.Errorf("expected Origin name: %s to be deleted", *origin.Name)
			}
		}
	}
}
