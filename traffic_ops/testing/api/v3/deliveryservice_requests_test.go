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
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"net/http"
	"strings"
	"testing"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
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
		UpdateTestDeliveryServiceRequests(t)
		GetTestDeliveryServiceRequestsIMSAfterChange(t, header)
	})
}

func GetTestDeliveryServiceRequestsIMSAfterChange(t *testing.T, header http.Header) {
	dsr := testData.DeliveryServiceRequests[dsrGood]
	_, reqInf, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, header)
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
	_, reqInf, err = TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, header)
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
		alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
		if err != nil {
			t.Errorf("Error occurred %v", err)
		}

		if len(alerts.Alerts) == 0 {
			t.Errorf("Expected: validation error alerts, actual: %+v", alerts)
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

		alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
		if err != nil {
			t.Errorf("Error occurred %v", err)
		}
		if len(alerts.Alerts) == 0 {
			t.Errorf("Expected: validation error alerts, actual: %+v", alerts)
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

		expected := []string{
			"deliveryservice_request was created.",
			//"'xmlId' the length must be between 1 and 48",
		}

		utils.Compare(t, expected, alerts.ToStrings())

		dsrs, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, nil)
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

		alerts, _, err := TOSession.CreateDeliveryServiceRequest(src)
		if err != nil {
			t.Errorf("Error creating DeliveryServiceRequest %v", err)
		}
		expected := []string{
			`'status' invalid transition from draft to pending`,
		}
		utils.Compare(t, expected, alerts.ToStrings())
	})
}

// TestDeliveryServiceRequestWorkflow tests that transitions of Status are
func TestDeliveryServiceRequestWorkflow(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants}, func() {
		// test empty request table
		dsrs, _, err := TOSession.GetDeliveryServiceRequests(nil)
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

		expected := []string{`deliveryservice_request was created.`}
		utils.Compare(t, expected, alerts.ToStrings())

		// Create a duplicate request -- should fail because xmlId is the same
		alerts, _, err = TOSession.CreateDeliveryServiceRequest(src)
		if err != nil {
			t.Errorf("Error creating DeliveryServiceRequest %v", err)
		}

		expected = []string{`An active request exists for delivery service 'test-transitions'`}
		utils.Compare(t, expected, alerts.ToStrings())

		dsrs, _, err = TOSession.GetDeliveryServiceRequestByXMLID(`test-transitions`, nil)
		if len(dsrs) != 1 {
			t.Errorf("Expected 1 deliveryServiceRequest -- got %d", len(dsrs))
			if len(dsrs) == 0 {
				t.Fatal("Cannot proceed")
			}
		}

		alerts, dsr := updateDeliveryServiceRequestStatus(t, dsrs[0], "submitted")

		expected = []string{
			"deliveryservice_request was updated.",
		}

		utils.Compare(t, expected, alerts.ToStrings())
		if dsr.Status != tc.RequestStatus("submitted") {
			t.Errorf("expected status=submitted,  got %s", string(dsr.Status))
		}
	})
}

func updateDeliveryServiceRequestStatus(t *testing.T, dsr tc.DeliveryServiceRequest, newstate string) (tc.Alerts, tc.DeliveryServiceRequest) {
	ID := dsr.ID
	dsr.Status = tc.RequestStatus("submitted")

	alerts, _, err := TOSession.UpdateDeliveryServiceRequestByID(ID, dsr)
	if err != nil {
		t.Errorf("Error updating deliveryservice_request: %v", err)
		return alerts, dsr
	}

	d, _, err := TOSession.GetDeliveryServiceRequestByID(ID, nil)
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
	futureTime := time.Now().AddDate(0,0,1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	dsr := testData.DeliveryServiceRequests[dsrGood]
	_, reqInf, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServiceRequests(t *testing.T) {
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, nil)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by XMLID: %v - %v", err, resp)
	}
}

func UpdateTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, nil)
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
	resp, _, err = TOSession.GetDeliveryServiceRequestByID(respDSR.ID, nil)
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
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, nil)
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
	dsrs, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID, nil)
	if err != nil {
		t.Errorf("error deleting DeliveryServiceRequest name: %s", err.Error())
	}
	if len(dsrs) > 0 {
		t.Errorf("expected DeliveryServiceRequest XMLID: %s to be deleted", dsr.DeliveryService.XMLID)
	}
}
