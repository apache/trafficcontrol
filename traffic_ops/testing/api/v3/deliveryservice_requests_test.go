package v3

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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	dsrGood      = 0
	dsrBadTenant = 1
	dsrRequired  = 2
	dsrDraft     = 3
)

func deleteDSR(id int, t *testing.T) {
	alerts, _, err := TOSession.DeleteDeliveryServiceRequestByID(id)
	if err != nil {
		t.Errorf("Cleaning up DSR #%d: %v - alerts: %v", id, err, alerts)
	}
}

func TestDeliveryServiceRequests(t *testing.T) {
	ReloadFixtures() // resets IDs
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServiceRequests}, func() {
		GetTestDeliveryServiceRequestsIMS(t)
		GetTestDeliveryServiceRequests(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		UpdateTestDeliveryServiceRequests(t)
		GetTestDeliveryServiceRequestsIMSAfterChange(t, header)
	})
}

func GetTestDeliveryServiceRequestsIMSAfterChange(t *testing.T, header http.Header) {
	// dsr := testData.DeliveryServiceRequests[dsrGood]
	_, reqInf, err := TOSession.GetDeliveryServiceRequestsV30(header, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	_, reqInf, err = TOSession.GetDeliveryServiceRequestsV30(header, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestDeliveryServiceRequests(t *testing.T) {
	t.Log("CreateTestDeliveryServiceRequests")

	dsr := testData.DeliveryServiceRequests[dsrGood]
	respDSR, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
	t.Log("Response: ", respDSR)
	if err != nil {
		t.Errorf("could not CREATE DeliveryServiceRequests: %v", err)
	}

}

// verifies requirements of creating a "delete" DSR
func TestCreateDeletionRequest(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{Types, CDNs, Tenants, CacheGroups, Topologies, DeliveryServices, Parameters}, func() {
		dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
		if err != nil {
			t.Fatalf("Failed to retrieve Delivery Services: %v", err)
		}
		if len(dses) < 1 {
			t.Fatalf("Failed to retrieve Delivery Services: no DSes returned from Traffic Ops")
		}
		ds := dses[0]
		if ds.ID == nil {
			t.Fatal("DS chosen for testing has no ID")
		}

		dsr := tc.DeliveryServiceRequestV30{
			ChangeType: tc.DSRChangeTypeDelete,
			Requested:  &ds,
			Status:     tc.RequestStatusDraft,
		}

		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err == nil {
			t.Error("Didn't get expected error creating a 'delete' DSR with no Original")
		} else {
			t.Logf("Received expected error creating 'delete' DSR with no Original: %v", err)
		}

		found := false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.SuccessLevel.String() {
				t.Errorf("Received unexpected success-level alert creating 'delete' DSR with no Original: %s", alert.Text)
			} else if alert.Level == tc.ErrorLevel.String() {
				t.Logf("Received expected error-level alert creating 'delete' DSR with no Original: %s", alert.Text)
				found = true
			}
		}
		if !found {
			t.Error("Didn't find expected error-level alert creating 'delete' DSR with no Original")
		}

		dsr.Original = &ds
		dsr.Requested = nil
		alerts, inf, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err != nil {
			t.Fatalf("Failed to create 'delete' DSR: %v", err)
		}
		t.Logf("alerts: %+v", alerts)
		t.Logf("status code: %d %s", inf.StatusCode, http.StatusText(inf.StatusCode))

		dsrs, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, nil)
		t.Logf("Fetched DSRs: %+v", dsrs)
		if err != nil {
			t.Errorf("Failed to get DSRs after creation: %v", err)
		} else if len(dsrs) < 1 {
			t.Error("Expected a DSR to exist after one was created, but it didn't")
		} else if dsrs[0].ID == nil {
			t.Error("DSR had no ID after creation")
		} else if alerts, _, err = TOSession.DeleteDeliveryServiceRequestByID(*dsrs[0].ID); err != nil {
			t.Errorf("Failed to delete DSR: %v - alerts: %v", err, alerts)
		}
	})
}

// verifies requirements of creating a "create" DSR
func TestCreateCreationRequest(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{Types, CDNs, Tenants}, func() {
		if len(testData.DeliveryServices) < 1 {
			t.Fatal("Need at least one testing DS for testing")
		}
		ds := testData.DeliveryServices[0]

		dsr := tc.DeliveryServiceRequestV30{
			ChangeType: tc.DSRChangeTypeCreate,
			Original:   &ds,
			Status:     tc.RequestStatusDraft,
		}

		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err == nil {
			t.Error("Didn't get expected error creating a 'create' DSR with no Requested")
		} else {
			t.Logf("Received expected error creating 'create' DSR with no Requested: %v", err)
		}

		found := false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.SuccessLevel.String() {
				t.Errorf("Received unexpected success-level alert creating 'create' DSR with no Requested: %s", alert.Text)
			} else if alert.Level == tc.ErrorLevel.String() {
				t.Logf("Received expected error-level alert creating 'create' DSR with no Requested: %s", alert.Text)
				found = true
			}
		}
		if !found {
			t.Error("Didn't find expected error-level alert creating 'create' DSR with no Requested")
		}

		dsr.Requested = &ds
		dsr.Original = nil
		alerts, _, err = TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err != nil {
			t.Errorf("Failed to create 'create' DSR: %v", err)
		} else if dsrs, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, nil); err != nil {
			t.Errorf("Failed to get DSRs after creation: %v", err)
		} else if len(dsrs) < 1 {
			t.Error("Expected a DSR to exist after one was created, but it didn't")
		} else if dsrs[0].ID == nil {
			t.Error("DSR had no ID after creation")
		} else if alerts, _, err = TOSession.DeleteDeliveryServiceRequestByID(*dsrs[0].ID); err != nil {
			t.Errorf("Failed to delete DSR: %v - alerts: %v", err, alerts)
		}
	})
}

func TestDeliveryServiceRequestRequired(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		dsr := testData.DeliveryServiceRequests[dsrRequired]
		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err == nil {
			t.Fatal("Expected an error creating a DSR missing required fields, but didn't get one")
		}
		t.Logf("Received expected error creating DSR missing required fields %v", err)

		if len(alerts.Alerts) == 0 {
			t.Errorf("Expected: validation error alerts, actual: %+v", alerts)
		}
	})
}

func TestDeliveryServiceRequestGetAssignee(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		if len(testData.DeliveryServiceRequests) < 1 {
			t.Fatal("Need at least one DSR for testing")
		}
		dsr := testData.DeliveryServiceRequests[0]
		me, _, err := TOSession.GetUserCurrent()
		if err != nil {
			t.Fatalf("Fetching current user: %v", err)
		}
		if me.UserName == nil {
			t.Fatal("Current user has no username")
		}
		if me.ID == nil {
			t.Fatal("Current user has no ID")
		}
		dsr.Assignee = me.UserName
		dsr.AssigneeID = me.ID
		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err != nil {
			t.Fatalf("Creating DSR: %v - %v", err, alerts)
		}

		dsrs, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, nil)
		if err != nil {
			t.Fatalf("Fetching DSRs: %v", err)
		}
		if len(dsrs) < 1 {
			t.Fatal("No DSRs returned after creating one")
		}
		d := dsrs[0]
		if len(dsrs) > 1 {
			t.Errorf("Too many DSRs returned after creating only one: %d", len(dsrs))
			t.Logf("Testing will proceed with DSR: %v", d)
		}
		if d.ID == nil {
			t.Fatal("Got DSR with no ID")
		}

		assignee, _, err := TOSession.GetDeliveryServiceRequestAssignment(*d.ID, nil)
		if err != nil {
			t.Errorf("Error fetching DSR assignee: %v", err)
		}
		if assignee == nil {
			t.Fatal("DSR had no assignee after assigning")
		}
		if *assignee != *me.UserName {
			t.Fatalf("Incorrect assignee after assignment; want: '%s', got: '%s'", *me.UserName, *assignee)
		}

	})
}

func TestDeliveryServiceRequestRules(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		routingName := strings.Repeat("X", 1) + "." + strings.Repeat("X", 48)
		// Test the xmlId length and form
		XMLID := "X " + strings.Repeat("X", 46)
		displayName := strings.Repeat("X", 49)

		dsr := testData.DeliveryServiceRequests[dsrGood]
		dsr.Requested.DisplayName = &displayName
		dsr.Requested.RoutingName = &routingName
		dsr.Requested.XMLID = &XMLID

		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err == nil {
			t.Error("Expected creating DSR with fields that fail validation to fail, but it didn't")
		} else {
			t.Logf("Received expected error creating DSR with fields that fail validation: %v", err)
		}
		if len(alerts.Alerts) == 0 {
			t.Errorf("Expected: validation error alerts, actual: %+v", alerts)
		}
	})
}

func TestDeliveryServiceRequestBadStatusOnCreation(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		src := testData.DeliveryServiceRequests[dsrDraft]

		for _, s := range []tc.RequestStatus{tc.RequestStatusPending, tc.RequestStatusComplete, tc.RequestStatusRejected} {
			src.Status = s

			alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(src, nil)
			if err == nil {
				t.Error("Expected an error creating a bad DSR, but didn't get one")
			} else {
				t.Logf("Received expected error creating DSR: %v", err)
			}

			found := false
			for _, alert := range alerts.Alerts {
				if alert.Level == tc.SuccessLevel.String() {
					t.Errorf("Unexpected success creating bad DSR: %v", alert.Text)
				} else if alert.Level == tc.ErrorLevel.String() {
					t.Logf("Received expected error-level alert creating bad DSR: %v", alert.Text)
					found = true
				}
			}

			if !found {
				t.Error("Didn't find an error alert when creating a bad DSR")
			}
		}
	})
}

// TestDeliveryServiceRequestWorkflow tests that the rules governing transitions
// of Status are properly enforced.
func TestDeliveryServiceRequestWorkflow(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{Types, CDNs, Tenants, CacheGroups, Topologies, DeliveryServices, Parameters}, func() {
		// test empty request table
		dsrs, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, nil)
		if err != nil {
			t.Errorf("Error getting empty list of DeliveryServiceRequests %v++", err)
		}
		if dsrs == nil {
			t.Error("Expected empty DeliveryServiceRequest slice -- got nil")
		}
		if len(dsrs) != 0 {
			t.Errorf("Expected no entries in DeliveryServiceRequest slice -- got %d", len(dsrs))
		}

		// Create a draft request
		src := testData.DeliveryServiceRequests[dsrDraft]
		src.SetXMLID()

		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(src, nil)
		if err != nil {
			t.Errorf("Error creating DeliveryServiceRequest %v - %v", err, alerts)
		}

		found := false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.ErrorLevel.String() {
				t.Errorf("Unexpected error-level alert creating DSR: %s", alert.Text)
			} else if alert.Level == tc.SuccessLevel.String() {
				t.Logf("Received expected success-level alert creating DSR: %s", alert.Text)
				found = true
			}
		}
		if !found {
			t.Error("Didn't find success-level alert after creating a DSR")
		}

		// Create a duplicate request -- should fail because xmlId is the same
		alerts, _, err = TOSession.CreateDeliveryServiceRequestV30(src, nil)
		if err == nil {
			t.Error("Expected an error creating duplicate request - didn't get one")
		} else {
			t.Logf("Received expected error creating Delivery Service Request: %v", err)
		}

		found = false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.SuccessLevel.String() {
				t.Errorf("Unexpected success message creating duplicate DSR: %v", alert.Text)
			} else if alert.Level == tc.ErrorLevel.String() {
				found = true
			}
		}
		if !found {
			t.Error("Didn't find expected error-level alert when creating duplicate DSR")
		}

		params := url.Values{}
		params.Set("xmlId", src.XMLID)
		dsrs, _, err = TOSession.GetDeliveryServiceRequestsV30(nil, params)
		if len(dsrs) != 1 {
			t.Errorf("Expected 1 deliveryServiceRequest -- got %d", len(dsrs))
			if len(dsrs) == 0 {
				t.Fatal("Cannot proceed")
			}
		}

		alerts, dsr := updateDeliveryServiceRequestStatus(t, dsrs[0], "submitted")
		found = false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.ErrorLevel.String() {
				t.Errorf("Unexpected error-level alert: %s", alert.Text)
			} else if alert.Level == tc.SuccessLevel.String() && strings.Contains(alert.Text, "updated") {
				found = true
			}
		}

		if !found {
			t.Error("Didn't find success-level alert after updating")
		}

		if dsr.Status != tc.RequestStatus("submitted") {
			t.Errorf("expected status=submitted,  got %s", string(dsr.Status))
		}
	})
}

func updateDeliveryServiceRequestStatus(t *testing.T, dsr tc.DeliveryServiceRequestV30, newstate string) (tc.Alerts, tc.DeliveryServiceRequestV30) {
	if dsr.ID == nil {
		t.Error("Cannot update DSR with no ID")
		return tc.Alerts{}, tc.DeliveryServiceRequestV30{}
	}

	ID := *dsr.ID
	dsr.Status = tc.RequestStatus("submitted")

	alerts, _, err := TOSession.UpdateDeliveryServiceRequest(ID, dsr, nil)
	if err != nil {
		t.Errorf("Error updating deliveryservice_request: %v", err)
		return alerts, dsr
	}

	params := url.Values{}
	params.Set("id", strconv.Itoa(ID))
	d, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, params)
	if err != nil {
		t.Errorf("Error updating deliveryservice_request %d: %v", ID, err)
		return alerts, dsr
	}

	if len(d) != 1 {
		t.Errorf("Expected 1 deliveryservice_request, got %d", len(d))
	}
	return alerts, d[0]
}

func GetTestDeliveryServiceRequestsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	// dsr := testData.DeliveryServiceRequests[dsrGood]
	_, reqInf, err := TOSession.GetDeliveryServiceRequestsV30(header, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServiceRequests(t *testing.T) {
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(*dsr.Requested.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by XMLID: %v - %v", err, resp)
	}
}

// verifies that 'original' is a moving target until a DSR is closed
func TestOriginalNotFixedUntilClosed(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{Types, CDNs, Tenants, CacheGroups, Topologies, DeliveryServices, Parameters}, func() {
		dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
		if err != nil {
			t.Fatalf("Failed to retrieve Delivery Services: %v", err)
		}
		if len(dses) < 1 {
			t.Fatalf("Failed to retrieve Delivery Services: no DSes returned from Traffic Ops")
		}
		ds := dses[0]
		if ds.ID == nil {
			t.Fatal("DS chosen for testing has no ID")
		}
		if ds.Active == nil {
			t.Fatal("DS chosen for testing has no Active property")
		}

		originalActive := *ds.Active

		*ds.Active = !*ds.Active

		me, _, err := TOSession.GetUserCurrentWithHdr(nil)
		if err != nil {
			t.Fatalf("Failed to get current user from Traffic Ops: %v", err)
		}
		if me == nil {
			t.Fatal("Failed to get current user from Traffic Ops: user was nil")
		}
		if me.UserName == nil {
			t.Fatal("Current user has no username")
		}
		if me.ID == nil {
			t.Fatal("Current user has no ID")
		}

		dsr := tc.DeliveryServiceRequestV30{
			Author:     *me.UserName,
			ChangeType: tc.DSRChangeTypeUpdate,
			Requested:  &ds,
			Status:     tc.RequestStatusSubmitted,
		}

		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err != nil {
			t.Fatalf("Failed to create initial DSR: %v - alerts: %v", err, alerts)
		}

		dsrs, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, url.Values{"author": []string{*me.UserName}})
		if err != nil {
			t.Fatalf("Failed to fetch newly created DSR: %v", err)
		}
		if len(dsrs) != 1 {
			t.Fatalf("Expected exactly one dsr to exist after creating only one; got: %d", len(dsrs))
		}
		dsr = dsrs[0]
		if dsr.ID == nil {
			t.Fatal("Test DSR had nil ID after creation")
		}
		defer deleteDSR(*dsr.ID, t)
		if dsr.Requested == nil {
			t.Fatal("Test DSR had no 'requested' field after creation")
		}
		if dsr.Original == nil {
			t.Fatal("Test DSR had no 'original' field in GET response")
		}
		if dsr.Original.Active == nil {
			t.Fatal("Test DSR's original had no 'active' field in GET response")
		}

		if *dsr.Original.Active != originalActive {
			t.Errorf("Incorrect original.active in GET response; want: %t, got: %t", *ds.Active, *dsr.Original.Active)
		}

		_, _, err = TOSession.UpdateDeliveryServiceV30(*ds.ID, ds)
		if err != nil {
			t.Fatalf("Failed to modify Delivery Service: %v", err)
		}

		dsrs, _, err = TOSession.GetDeliveryServiceRequestsV30(nil, url.Values{"author": []string{*me.UserName}})
		if err != nil {
			t.Fatalf("Failed to fetch DSR after original was modified: %v", err)
		}
		if len(dsrs) != 1 {
			t.Fatalf("Expected exactly one dsr to exist after creating only one; got: %d", len(dsrs))
		}
		dsr = dsrs[0]
		if dsr.ID == nil {
			t.Fatal("Test DSR had nil ID after original was modified")
		}
		if dsr.Requested == nil {
			t.Fatal("Test DSR had no 'requested' field after original was modified")
		}
		if dsr.Original == nil {
			t.Fatal("Test DSR had no 'original' field in GET response")
		}
		if dsr.Original.Active == nil {
			t.Fatal("Test DSR's original had no 'active' field in GET response")
		}

		if *dsr.Original.Active == originalActive {
			t.Errorf("Incorrect original.active in GET response; want: %t, got: %t", !originalActive, *dsr.Original.Active)
		}
	})
}

// verifies that closed DSRs cannot be updated - except to complete pending DSRs
func TestUpdateClosedDSR(t *testing.T) {
	ReloadFixtures()
	WithObjs(t, []TCObj{Types, CDNs, Tenants, CacheGroups, Topologies, DeliveryServices, Parameters}, func() {
		dses, _, err := TOSession.GetDeliveryServicesV30WithHdr(nil, nil)
		if err != nil {
			t.Fatalf("Failed to retrieve Delivery Services: %v", err)
		}
		if len(dses) < 1 {
			t.Fatalf("Failed to retrieve Delivery Services: no DSes returned from Traffic Ops")
		}
		ds := dses[0]
		if ds.ID == nil {
			t.Fatal("DS chosen for testing has no ID")
		}
		if ds.Active == nil {
			ds.Active = new(bool)
		}

		*ds.Active = !*ds.Active

		me, _, err := TOSession.GetUserCurrentWithHdr(nil)
		if err != nil {
			t.Fatalf("Failed to get current user from Traffic Ops: %v", err)
		}
		if me == nil {
			t.Fatal("Failed to get current user from Traffic Ops: user was nil")
		}
		if me.UserName == nil {
			t.Fatal("Current user has no username")
		}
		if me.ID == nil {
			t.Fatal("Current user has no ID")
		}

		dsr := tc.DeliveryServiceRequestV30{
			Author:     *me.UserName,
			ChangeType: tc.DSRChangeTypeUpdate,
			Requested:  &ds,
			Status:     tc.RequestStatusSubmitted,
		}

		alerts, _, err := TOSession.CreateDeliveryServiceRequestV30(dsr, nil)
		if err != nil {
			t.Fatalf("Failed to create initial DSR: %v - alerts: %v", err, alerts)
		}

		dsrs, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, url.Values{"author": []string{*me.UserName}})
		if err != nil {
			t.Fatalf("Failed to fetch newly created DSR: %v", err)
		}
		if len(dsrs) != 1 {
			t.Fatalf("Expected exactly one dsr to exist after creating only one; got: %d", len(dsrs))
		}
		dsr = dsrs[0]
		if dsr.ID == nil {
			t.Fatalf("Test DSR had nil ID after creation")
		}
		defer deleteDSR(*dsr.ID, t)

		alerts, _, err = TOSession.SetDeliveryServiceRequestStatus(*dsr.ID, tc.RequestStatusPending, nil)
		if err != nil {
			t.Fatalf("Failed to close test DSR: %v - alerts: %v", err, alerts)
		}

		*dsr.Requested.Active = !*dsr.Requested.Active
		alerts, _, err = TOSession.UpdateDeliveryServiceRequest(*dsr.ID, dsr, nil)
		if err == nil {
			t.Errorf("Didn't get expected error modifying a closed DSR")
		} else {
			t.Logf("Received expected error modifying closed DSR: %v", err)
		}

		found := false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.SuccessLevel.String() {
				t.Errorf("Found success message updating a closed DSR: %s", alert.Text)
			} else if alert.Level == tc.ErrorLevel.String() {
				t.Logf("Found expected error-level alert: %v", alert.Text)
				found = true
			}
		}

		if !found {
			t.Error("Didn't find expected error-level alert modiying a closed DSR")
		}

		alerts, _, err = TOSession.SetDeliveryServiceRequestStatus(*dsr.ID, tc.RequestStatusComplete, nil)
		if err != nil {
			t.Fatalf("Failed to complete the pending test DSR: %v - alerts: %v", err, alerts)
		}

		alerts, _, err = TOSession.UpdateDeliveryServiceRequest(*dsr.ID, dsr, nil)
		if err == nil {
			t.Errorf("Didn't get expected error modifying a closed DSR")
		} else {
			t.Logf("Received expected error modifying closed DSR: %v", err)
		}

		found = false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.SuccessLevel.String() {
				t.Errorf("Found success message updating a closed DSR: %s", alert.Text)
			} else if alert.Level == tc.ErrorLevel.String() {
				t.Logf("Found expected error-level alert: %v", alert.Text)
				found = true
			}
		}
		if !found {
			t.Error("Didn't find expected error-level alert modiying a closed DSR")
		}

		alerts, _, err = TOSession.SetDeliveryServiceRequestStatus(*dsr.ID, tc.RequestStatusRejected, nil)
		if err == nil {
			t.Error("Didn't get expected error trying to set the status of a complete DSR")
		} else {
			t.Logf("Received expected error setting the status of a complete DSR: %v", err)
		}

		found = false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.SuccessLevel.String() {
				t.Errorf("Found success message setting the status of a complete DSR: %s", alert.Text)
			} else if alert.Level == tc.ErrorLevel.String() {
				t.Logf("Found expected error-level alert: %v", alert.Text)
				found = true
			}
		}
		if !found {
			t.Errorf("Didn't find expected error-level alert setting the status of a complete DSR: %+v", alerts)
		}
	})
}

func UpdateTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	dsr.SetXMLID()
	params := url.Values{}
	params.Set("xmlId", dsr.XMLID)
	resp, _, err := TOSession.GetDeliveryServiceRequestsV30(nil, params)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by XMLID '%s': %v - %v", dsr.XMLID, *dsr.Requested.XMLID, err)
	}
	if len(resp) == 0 {
		t.Fatal("Length of GET DeliveryServiceRequest is 0")
	}
	respDSR := resp[0]
	if respDSR.Requested == nil {
		t.Fatalf("Got back DSR without 'requested' (changetype: '%s')", respDSR.ChangeType)
	}
	if respDSR.ID == nil {
		t.Fatal("Got back DSR without ID")
	}
	expDisplayName := "new display name"
	respDSR.Requested.DisplayName = &expDisplayName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDeliveryServiceRequest(*respDSR.ID, respDSR, nil)
	t.Log("Response: ", alert)
	if err != nil {
		t.Fatalf("cannot UPDATE DeliveryServiceRequest by id: %v - %v", err, alert)
	}

	// Retrieve the DeliveryServiceRequest to check DeliveryServiceRequest name got updated
	params.Del("xmlId")
	params.Set("id", strconv.Itoa(*respDSR.ID))
	resp, _, err = TOSession.GetDeliveryServiceRequestsV30(nil, params)
	if err != nil {
		t.Fatalf("cannot GET DeliveryServiceRequest by ID: %v - %v", respDSR.ID, err)
	}
	if len(resp) < 1 {
		t.Fatalf("No DSR by ID %d after updating that DSR", *respDSR.ID)
	}
	respDSR = resp[0]
	if respDSR.Requested == nil {
		t.Fatal("Got back DSR without 'requested' after update")
	}
	if respDSR.Requested.DisplayName == nil {
		t.Fatal("Got back DSR with null 'requested.displayName' after updating that field to non-null value")
	}
	if *respDSR.Requested.DisplayName != expDisplayName {
		t.Errorf("results do not match actual: '%s', expected: '%s'", *respDSR.Requested.DisplayName, expDisplayName)
	}
}

func DeleteTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(*dsr.Requested.XMLID)
	if err != nil || len(resp) < 1 {
		t.Fatalf("cannot GET DeliveryServiceRequest by XMLID: %v - %v", *dsr.Requested.XMLID, err)
	}
	respDSR := resp[0]
	alert, _, err := TOSession.DeleteDeliveryServiceRequestByID(respDSR.ID)
	t.Log("Response: ", alert)
	if err != nil {
		t.Errorf("cannot DELETE DeliveryServiceRequest by id: %d - %v - %v", respDSR.ID, err, alert)
	}

	// Retrieve the DeliveryServiceRequest to see if it got deleted
	dsrs, _, err := TOSession.GetDeliveryServiceRequestByXMLID(*dsr.Requested.XMLID)
	if err != nil {
		t.Errorf("error deleting DeliveryServiceRequest name: %s", err.Error())
	}
	if len(dsrs) > 0 {
		t.Errorf("expected DeliveryServiceRequest XMLID: %s to be deleted", *dsr.Requested.XMLID)
	}
}
