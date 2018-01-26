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

package v13

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/testing/api/utils"
)

const GOOD_DSR = 0
const REQUIRED_DSR = 2

func TestDeliveryServiceRequests(t *testing.T) {

	CreateTestDeliveryServiceRequests(t)
	GetTestDeliveryServiceRequests(t)
	UpdateTestDeliveryServiceRequests(t)
	DeleteTestDeliveryServiceRequests(t)

}

func CreateTestDeliveryServiceRequests(t *testing.T) {
	fmt.Printf("CreateTestDeliveryServiceRequests\n")

	dsr := testData.DeliveryServiceRequests[GOOD_DSR]
	resp, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	log.Debugln("Response: ", resp)
	if err != nil {
		t.Errorf("could not CREATE DeliveryServiceRequests: %v\n", err)
	}

}

func TestDeliveryServiceRequestRequired(t *testing.T) {

	dsr := testData.DeliveryServiceRequests[REQUIRED_DSR]
	alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	if err != nil {
		t.Errorf("Error occurred %v", err)
	}

	expected := []string{
		"'active' cannot be blank",
		"'cdnId' cannot be blank",
		"'dscp' cannot be blank",
		"'geoLimit' cannot be blank",
		"'geoProvider' cannot be blank",
		"'infoUrl' must be a valid URL",
		"'initialDispersion' must be greater than zero",
		"'logsEnabled' cannot be blank",
		"'orgServerFqdn' must be a valid URL",
		"'regionalGeoBlocking' cannot be blank",
		"'routingName' must be a valid hostname",
		"'typeId' cannot be blank",
		"'xmlId' cannot contain spaces",
	}

	utils.Compare(t, expected, alerts.ToStrings())

}

func TestDeliveryServiceRequestRules(t *testing.T) {
	fmt.Printf("TestDeliveryServiceRequestRules\n")

	routingName := strings.Repeat("X", 1) + "." + strings.Repeat("X", 48)
	// Test the xmlId length
	XMLID := strings.Repeat("X", 48)
	displayName := strings.Repeat("X", 49)

	dsr := testData.DeliveryServiceRequests[GOOD_DSR]
	dsr.DeliveryService.DisplayName = displayName
	dsr.DeliveryService.RoutingName = routingName
	dsr.DeliveryService.XMLID = XMLID
	alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	if err != nil {
		t.Errorf("Error occurred %v", err)
	}

	expected := []string{
		"'routingName' cannot contain periods",
		"'xmlId' cannot contain spaces",
	}

	utils.Compare(t, expected, alerts.ToStrings())

}

func TestDeliveryServiceRequestTypeFields(t *testing.T) {
	fmt.Printf("TestDeliveryServiceRequestTypeFields\n")

	dsr := testData.DeliveryServiceRequests[GOOD_DSR]
	alerts, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
	if err != nil {
		t.Errorf("Error occurred %v", err)
	}

	expected := []string{
		"not authorized on this tenant",
		//"'xmlId' the length must be between 1 and 48",
	}

	utils.Compare(t, expected, alerts.ToStrings())

}

func GetTestDeliveryServiceRequests(t *testing.T) {
	fmt.Printf("GetTestDeliveryServiceRequests\n")

	dsr := testData.DeliveryServiceRequests[GOOD_DSR]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by XMLID: %v - %v\n", err, resp)
	}
}

func UpdateTestDeliveryServiceRequests(t *testing.T) {

	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	dsr := testData.DeliveryServiceRequests[GOOD_DSR]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v\n", dsr.DeliveryService.XMLID, err)
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
	dsr := testData.DeliveryServiceRequests[GOOD_DSR]
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
		t.Errorf("expected DeliveryServiceRequest name: %s to be deleted\n", dsr.DeliveryService.XMLID)
	}
}
