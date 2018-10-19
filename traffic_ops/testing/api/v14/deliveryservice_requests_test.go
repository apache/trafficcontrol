package v14

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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
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

	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestDeliveryServiceRequests(t)
	GetTestDeliveryServiceRequests(t)
	UpdateTestDeliveryServiceRequests(t)
	DeleteTestDeliveryServiceRequests(t)

	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

func CreateTestDeliveryServiceRequests(t *testing.T) {
	log.Debugln("CreateTestDeliveryServiceRequests")

	// Attach CDNs
	cdn := testData.CDNs[0]
	resp, _, err := TOSession.GetCDNByName(cdn.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: %v - %v\n", cdn.Name, err)
	}
	respCDN := resp[0]

	// Attach Type
	typ := testData.DeliveryServiceRequests[dsrGood].DeliveryService.Type.String()
	respTypes, _, err := TOSession.GetTypeByName(typ)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v\n", typ, err)
	}
	respTyp := respTypes[0]

	dsr := testData.DeliveryServiceRequests[dsrGood]
	dsr.DeliveryService.CDNID = respCDN.ID
	dsr.DeliveryService.TypeID = respTyp.ID
	respDSR, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	log.Debugln("Response: ", respDSR)
	if err != nil {
		t.Errorf("could not CREATE DeliveryServiceRequests: %v\n", err)
	}

}

func TestDeliveryServiceRequestRequired(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	dsr := testData.DeliveryServiceRequests[dsrRequired]
	alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	if err != nil {
		t.Errorf("Error occurred %v", err)
	}

	if len(alerts.Alerts) == 0 {
		t.Errorf("Expected: validation error alerts, actual: %+v", alerts)
	}
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

func TestDeliveryServiceRequestRules(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)
	routingName := strings.Repeat("X", 1) + "." + strings.Repeat("X", 48)
	// Test the xmlId length and form
	XMLID := "X " + strings.Repeat("X", 46)
	displayName := strings.Repeat("X", 49)

	dsr := testData.DeliveryServiceRequests[dsrGood]
	dsr.DeliveryService.DisplayName = displayName
	dsr.DeliveryService.RoutingName = routingName
	dsr.DeliveryService.XMLID = XMLID

	// Attach Types
	typ := testData.Types[3]
	rt, _, err := TOSession.GetTypeByName(typ.Name)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v\n", typ.Name, err)
	}
	respType := rt[0]

	// Attach CDNs
	cdn := testData.CDNs[3]
	resp, _, err := TOSession.GetCDNByName(cdn.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: %v - %v\n", cdn.Name, err)
	}
	respCDN := resp[0]
	dsr.DeliveryService.TypeID = respType.ID
	dsr.DeliveryService.CDNID = respCDN.ID

	alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	if err != nil {
		t.Errorf("Error occurred %v", err)
	}
	if len(alerts.Alerts) == 0 {
		t.Errorf("Expected: validation error alerts, actual: %+v", alerts)
	}
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

func TestDeliveryServiceRequestTypeFields(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestParameters(t)

	dsr := testData.DeliveryServiceRequests[dsrBadTenant]

	// Attach Types
	typ := testData.Types[3]
	rt, _, err := TOSession.GetTypeByName(typ.Name)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v\n", typ.Name, err)
	}
	respType := rt[0]

	// Attach CDNs
	cdn := testData.CDNs[3]
	resp, _, err := TOSession.GetCDNByName(cdn.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: %v - %v\n", cdn.Name, err)
	}
	respCDN := resp[0]
	dsr.DeliveryService.TypeID = respType.ID
	dsr.DeliveryService.CDNID = respCDN.ID

	alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	if err != nil {
		t.Errorf("Error occurred %v", err)
	}

	expected := []string{
		"deliveryservice_request was created.",
		//"'xmlId' the length must be between 1 and 48",
	}

	utils.Compare(t, expected, alerts.ToStrings())

	dsrs, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if len(dsrs) != 1 {
		t.Errorf("expected 1 deliveryservice_request with XMLID %s;  got %d", dsr.DeliveryService.XMLID, len(dsrs))
	}
	alert, _, err := TOSession.DeleteDeliveryServiceRequestByID(dsrs[0].ID)
	if err != nil {
		t.Errorf("cannot DELETE DeliveryServiceRequest by id: %d - %v - %v\n", dsrs[0].ID, err, alert)
	}

	DeleteTestParameters(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

func TestDeliveryServiceRequestBad(t *testing.T) {
	CreateTestCDNs(t)
	CreateTestTypes(t)
	// try to create non-draft/submitted
	src := testData.DeliveryServiceRequests[dsrDraft]
	s, err := tc.RequestStatusFromString("pending")
	if err != nil {
		t.Errorf(`unable to create Status from string "pending"`)
	}
	src.Status = s

	// Attach Types
	typ := testData.Types[3]
	rt, _, err := TOSession.GetTypeByName(typ.Name)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v\n", typ.Name, err)
	}
	respType := rt[0]

	// Attach CDNs
	cdn := testData.CDNs[3]
	resp, _, err := TOSession.GetCDNByName(cdn.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: %v - %v\n", cdn.Name, err)
	}
	respCDN := resp[0]
	src.DeliveryService.TypeID = respType.ID
	src.DeliveryService.CDNID = respCDN.ID

	alerts, _, err := TOSession.CreateDeliveryServiceRequest(src)
	if err != nil {
		t.Errorf("Error creating DeliveryServiceRequest %v", err)
	}
	expected := []string{
		`'status' invalid transition from draft to pending`,
	}
	utils.Compare(t, expected, alerts.ToStrings())
	DeleteTestTypes(t)
	DeleteTestCDNs(t)
}

// TestDeliveryServiceRequestWorkflow tests that transitions of Status are
func TestDeliveryServiceRequestWorkflow(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)
	// test empty request table
	dsrs, _, err := TOSession.GetDeliveryServiceRequests()
	if err != nil {
		t.Errorf("Error getting empty list of DeliveryServiceRequests %v++", err)
	}
	if dsrs == nil {
		t.Errorf("Expected empty DeliveryServiceRequest slice -- got nil")
	}
	if len(dsrs) != 0 {
		t.Errorf("Expected no entries in DeliveryServiceRequest slice -- got %d", len(dsrs))
	}

	// Create a draft request
	src := testData.DeliveryServiceRequests[dsrDraft]

	// Attach Types
	typ := testData.Types[3]
	rt, _, err := TOSession.GetTypeByName(typ.Name)
	if err != nil {
		t.Errorf("cannot GET Type by name: %v - %v\n", typ.Name, err)
	}
	respType := rt[0]

	// Attach CDNs
	cdn := testData.CDNs[3]
	resp, _, err := TOSession.GetCDNByName(cdn.Name)
	if err != nil {
		t.Errorf("cannot GET CDN by name: %v - %v\n", cdn.Name, err)
	}
	respCDN := resp[0]
	src.DeliveryService.TypeID = respType.ID
	src.DeliveryService.CDNID = respCDN.ID

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

	dsrs, _, err = TOSession.GetDeliveryServiceRequestByXMLID(`test-transitions`)
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
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

func updateDeliveryServiceRequestStatus(t *testing.T, dsr tc.DeliveryServiceRequest, newstate string) (tc.Alerts, tc.DeliveryServiceRequest) {
	ID := dsr.ID
	dsr.Status = tc.RequestStatus("submitted")

	alerts, _, err := TOSession.UpdateDeliveryServiceRequestByID(ID, dsr)
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

func GetTestDeliveryServiceRequests(t *testing.T) {
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by XMLID: %v - %v\n", err, resp)
	}
}

func UpdateTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Fatalf("cannot GET DeliveryServiceRequest by name: %v - %v\n", dsr.DeliveryService.XMLID, err)
	}
	if len(resp) == 0 {
		t.Fatal("Length of GET DeliveryServiceRequest is 0")
	}
	respDSR := resp[0]
	expDisplayName := "new display name"
	respDSR.DeliveryService.DisplayName = expDisplayName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDeliveryServiceRequestByID(respDSR.ID, respDSR)
	log.Debugln("Response: ", alert)
	if err != nil {
		t.Errorf("cannot UPDATE DeliveryServiceRequest by id: %v - %v\n", err, alert)
		return
	}

	// Retrieve the DeliveryServiceRequest to check DeliveryServiceRequest name got updated
	resp, _, err = TOSession.GetDeliveryServiceRequestByID(respDSR.ID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v\n", respDSR.ID, err)
	} else {
		respDSR = resp[0]
		if respDSR.DeliveryService.DisplayName != expDisplayName {
			t.Errorf("results do not match actual: %s, expected: %s\n", respDSR.DeliveryService.DisplayName, expDisplayName)
		}
	}

}

func DeleteTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[dsrGood]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by id: %v - %v\n", dsr.DeliveryService.XMLID, err)
	}
	respDSR := resp[0]
	alert, _, err := TOSession.DeleteDeliveryServiceRequestByID(respDSR.ID)
	log.Debugln("Response: ", alert)
	if err != nil {
		t.Errorf("cannot DELETE DeliveryServiceRequest by id: %d - %v - %v\n", respDSR.ID, err, alert)
	}

	// Retrieve the DeliveryServiceRequest to see if it got deleted
	dsrs, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("error deleting DeliveryServiceRequest name: %s\n", err.Error())
	}
	if len(dsrs) > 0 {
		t.Errorf("expected DeliveryServiceRequest XMLID: %s to be deleted\n", dsr.DeliveryService.XMLID)
	}
}
