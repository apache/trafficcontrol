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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/v6/lib/go-tc"
)

const (
	dsrGood      = 0
	dsrBadTenant = 1
	dsrRequired  = 2
	dsrDraft     = 3
)

func TestDeliveryServiceRequests(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServiceRequests}, func() {
		GetTestDeliveryServiceRequestsIMS(t)
		GetTestDeliveryServiceRequests(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestDeliveryServiceRequests(t)
		UpdateTestDeliveryServiceRequestsWithHeaders(t, header)
		GetTestDeliveryServiceRequestsIMSAfterChange(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestDeliveryServiceRequestsWithHeaders(t, header)
	})
}

func UpdateTestDeliveryServiceRequestsWithHeaders(t *testing.T, header http.Header) {
	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLIDWithHdr(dsr.DeliveryService.XMLID, header)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v", dsr.DeliveryService.XMLID, err)
	}
	if len(resp) == 0 {
		t.Fatal("Length of GET DeliveryServiceRequest is 0")
	}
	respDSR := resp[0]
	respDSR.DeliveryService.DisplayName = "new display name"

	_, reqInf, err := TOSession.UpdateDeliveryServiceRequestByIDWithHdr(respDSR.ID, respDSR, header)
	if err == nil {
		t.Errorf("Expected error about precondition failed, but got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServiceRequestsIMSAfterChange(t *testing.T, header http.Header) {
	dsr := testData.DeliveryServiceRequests[dsrGood]
	_, reqInf, err := TOSession.GetDeliveryServiceRequestByXMLIDWithHdr(dsr.DeliveryService.XMLID, header)
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
	_, reqInf, err = TOSession.GetDeliveryServiceRequestByXMLIDWithHdr(dsr.DeliveryService.XMLID, header)
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
	respDSR, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	t.Log("Response: ", respDSR)
	if err != nil {
		t.Errorf("could not CREATE DeliveryServiceRequests: %v", err)
	}

}

func TestDeliveryServiceRequestRequired(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		dsr := testData.DeliveryServiceRequests[dsrRequired]
		_, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
		if err == nil {
			t.Error("expected: validation error, actual: nil")
		}
	})
}

func TestDeliveryServiceRequestRules(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		routingName := strings.Repeat("X", 1) + "." + strings.Repeat("X", 48)
		// Test the xmlId length and form
		XMLID := "X " + strings.Repeat("X", 46)
		displayName := strings.Repeat("X", 49)

		dsr := testData.DeliveryServiceRequests[dsrGood]
		dsr.DeliveryService.DisplayName = displayName
		dsr.DeliveryService.RoutingName = routingName
		dsr.DeliveryService.XMLID = XMLID

		_, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
		if err == nil {
			t.Error("expected: validation error, actual: nil")
		}
	})
}

func TestDeliveryServiceRequestTypeFields(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters}, func() {
		dsr := testData.DeliveryServiceRequests[dsrBadTenant]

		alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
		if err != nil {
			t.Errorf("Error occurred %v", err)
		}

		found := false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.ErrorLevel.String() {
				t.Errorf("Expected only succuss-level alerts creating a DSR, got error-level alert: %s", alert.Text)
			} else if alert.Level == tc.SuccessLevel.String() {
				t.Logf("Got expected alert creating a DSR: %s", alert.Text)
				found = true
			}
		}
		if !found {
			t.Errorf("Expected a success-level alert creating a DSR, got none")
		}

		dsrs, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
		if len(dsrs) != 1 {
			t.Errorf("expected 1 deliveryservice_request with XMLID %s;  got %d", dsr.DeliveryService.XMLID, len(dsrs))
		}
		alert, _, err := TOSession.DeleteDeliveryServiceRequestByID(dsrs[0].ID)
		if err != nil {
			t.Errorf("cannot DELETE DeliveryServiceRequest by id: %d - %v - %v", dsrs[0].ID, err, alert)
		}
	})
}

func TestDeliveryServiceRequestBad(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		// try to create non-draft/submitted
		src := testData.DeliveryServiceRequests[dsrDraft]
		s, err := tc.RequestStatusFromString("pending")
		if err != nil {
			t.Errorf(`unable to create Status from string "pending"`)
		}
		src.Status = s

		_, _, err = TOSession.CreateDeliveryServiceRequest(src)
		if err == nil {
			t.Fatal("expected: validation error, actual: nil")
		}
	})
}

// TestDeliveryServiceRequestWorkflow tests that transitions of Status are
func TestDeliveryServiceRequestWorkflow(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		// test empty request table
		dsrs, _, err := TOSession.GetDeliveryServiceRequests()
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

		alerts, _, err := TOSession.CreateDeliveryServiceRequest(src)
		if err != nil {
			t.Errorf("Error creating DeliveryServiceRequest %v", err)
		}

		found := false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.ErrorLevel.String() {
				t.Errorf("Expected only succuss-level alerts creating a DSR, got error-level alert: %s", alert.Text)
			} else if alert.Level == tc.SuccessLevel.String() {
				t.Logf("Got expected alert creating a DSR: %s", alert.Text)
				found = true
			}
		}
		if !found {
			t.Errorf("Expected a success-level alert creating a DSR, got none")
		}

		// Create a duplicate request -- should fail because xmlId is the same
		alerts, _, err = TOSession.CreateDeliveryServiceRequest(src)
		if err == nil {
			t.Fatal("expected: validation error, actual: nil")
		}

		// TODO: the client needs to be fixed to return alerts on error. But
		// that's gotten harder to do "semantically" now that most of the logic
		// is shared between all API versions.
		// found = false
		// for _, alert := range alerts.Alerts {
		// 	if alert.Level == tc.SuccessLevel.String() {
		// 		t.Errorf("Expected only error-level alerts creating a duplicate DSR, got success-level alert: %s", alert.Text)
		// 	} else if alert.Level == tc.ErrorLevel.String() {
		// 		t.Logf("Got expected alert creating a duplicate DSR: %s", alert.Text)
		// 		found = true
		// 	}
		// }
		// if !found {
		// 	t.Errorf("Expected an error-level alert creating a duplicate DSR, got none")
		// }

		dsrs, _, err = TOSession.GetDeliveryServiceRequestByXMLID(`test-transitions`)
		if len(dsrs) != 1 {
			t.Errorf("Expected 1 deliveryServiceRequest -- got %d", len(dsrs))
			if len(dsrs) == 0 {
				t.Fatal("Cannot proceed")
			}
		}

		alerts, dsr := updateDeliveryServiceRequestStatus(t, dsrs[0], "submitted", nil)

		found = false
		for _, alert := range alerts.Alerts {
			if alert.Level == tc.ErrorLevel.String() {
				t.Errorf("Expected only succuss-level alerts updating a DSR, got error-level alert: %s", alert.Text)
			} else if alert.Level == tc.SuccessLevel.String() {
				t.Logf("Got expected alert updating a DSR: %s", alert.Text)
				found = true
			}
		}
		if !found {
			t.Errorf("Expected a success-level alert updating a DSR, got none")
		}

		if dsr.Status != tc.RequestStatus("submitted") {
			t.Errorf("expected status=submitted,  got %s", string(dsr.Status))
		}
	})
}

func updateDeliveryServiceRequestStatus(t *testing.T, dsr tc.DeliveryServiceRequest, newstate string, header http.Header) (tc.Alerts, tc.DeliveryServiceRequest) {
	ID := dsr.ID
	dsr.Status = tc.RequestStatus("submitted")

	alerts, _, err := TOSession.UpdateDeliveryServiceRequestByIDWithHdr(ID, dsr, header)
	if err != nil {
		t.Errorf("Error updating deliveryservice_request: %v", err)
		return alerts, dsr
	}

	d, _, err := TOSession.GetDeliveryServiceRequestByID(ID)
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
	dsr := testData.DeliveryServiceRequests[dsrGood]
	_, reqInf, err := TOSession.GetDeliveryServiceRequestByXMLIDWithHdr(dsr.DeliveryService.XMLID, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServiceRequests(t *testing.T) {
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by XMLID: %v - %v", err, resp)
	}
}

func UpdateTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v", dsr.DeliveryService.XMLID, err)
	}
	if len(resp) == 0 {
		t.Fatal("Length of GET DeliveryServiceRequest is 0")
	}
	respDSR := resp[0]
	expDisplayName := "new display name"
	respDSR.DeliveryService.DisplayName = expDisplayName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDeliveryServiceRequestByID(respDSR.ID, respDSR)
	t.Log("Response: ", alert)
	if err != nil {
		t.Errorf("cannot UPDATE DeliveryServiceRequest by id: %v - %v", err, alert)
		return
	}

	// Retrieve the DeliveryServiceRequest to check DeliveryServiceRequest name got updated
	resp, _, err = TOSession.GetDeliveryServiceRequestByID(respDSR.ID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v", respDSR.ID, err)
	} else {
		respDSR = resp[0]
		if respDSR.DeliveryService.DisplayName != expDisplayName {
			t.Errorf("results do not match actual: %s, expected: %s", respDSR.DeliveryService.DisplayName, expDisplayName)
		}
	}

}

func DeleteTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by id: %v - %v", dsr.DeliveryService.XMLID, err)
	}
	respDSR := resp[0]
	alert, _, err := TOSession.DeleteDeliveryServiceRequestByID(respDSR.ID)
	t.Log("Response: ", alert)
	if err != nil {
		t.Errorf("cannot DELETE DeliveryServiceRequest by id: %d - %v - %v", respDSR.ID, err, alert)
	}

	// Retrieve the DeliveryServiceRequest to see if it got deleted
	dsrs, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("error deleting DeliveryServiceRequest name: %s", err.Error())
	}
	if len(dsrs) > 0 {
		t.Errorf("expected DeliveryServiceRequest XMLID: %s to be deleted", dsr.DeliveryService.XMLID)
	}
}
