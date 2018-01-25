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
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/testing/api/utils"
)

func TestDeliveryServiceRequests(t *testing.T) {

	CreateTestDeliveryServiceRequests(t)
	GetTestDeliveryServiceRequests(t)
	/*
		UpdateTestDeliveryServiceRequests(t)
		DeleteTestDeliveryServiceRequests(t)
	*/

}

func CreateTestDeliveryServiceRequests(t *testing.T) {

	for _, dsr := range testData.DeliveryServiceRequests {
		resp, _, err := TOSession.CreateDeliveryServiceRequest(dsr)
		log.Debugln("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE DeliveryServiceRequests: %v\n", err)
		}
	}

}

func TestBadDeliveryServiceCreateRequests(t *testing.T) {

	routingName := strings.Repeat("X", 1) + "." + strings.Repeat("X", 48)
	// Test the xmlId length
	XMLID := strings.Repeat("X", 1) + " " + strings.Repeat("X", 48)
	displayName := strings.Repeat("X", 49)

	dsr := testData.DeliveryServiceRequests[2]
	dsr.DeliveryService.DisplayName = displayName
	dsr.DeliveryService.RoutingName = routingName
	dsr.DeliveryService.XMLID = XMLID
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
		"'routingName' cannot contain periods",
		"'typeId' cannot be blank",
		"'xmlId' cannot contain spaces",
		//"'xmlId' the length must be between 1 and 48",
	}

	alertsStrs := alerts.ToStrings()
	sort.Strings(alertsStrs)
	expectedFmt, _ := json.MarshalIndent(expected, "", "  ")
	errorsFmt, _ := json.MarshalIndent(alertsStrs, "", "  ")

	// Compare both directions
	for _, s := range alertsStrs {
		if !utils.FindNeedle(s, expected) {
			t.Errorf("\nExpected %s and \n Actual %v must match exactly", string(expectedFmt), string(errorsFmt))
			break
		}
	}

	// Compare both directions
	for _, s := range expected {
		if !utils.FindNeedle(s, alertsStrs) {
			t.Errorf("\nExpected %s and \n Actual %v must match exactly", string(expectedFmt), string(errorsFmt))
			break
		}
	}

}

func GetTestDeliveryServiceRequests(t *testing.T) {

	dsr := testData.DeliveryServiceRequests[0]
	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.DeliveryService.XMLID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v\n", err, resp)
	}
}

/*
func UpdateTestDeliveryServiceRequests(t *testing.T) {

	firstDeliveryServiceRequest := testData.DeliveryServiceRequests[0]
	// Retrieve the DeliveryServiceRequest by name so we can get the id for the Update
	resp, _, err := TOSession.GetDeliveryServiceRequestByName(firstDeliveryServiceRequest.Name)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v\n", firstDeliveryServiceRequest.Name, err)
	}
	remoteDeliveryServiceRequest := resp[0]
	expectedDeliveryServiceRequestName := "testDSR1"
	remoteDeliveryServiceRequest.Name = expectedDeliveryServiceRequestName
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDeliveryServiceRequestByID(remoteDeliveryServiceRequest.ID, remoteDeliveryServiceRequest)
	if err != nil {
		t.Errorf("cannot UPDATE DeliveryServiceRequest by id: %v - %v\n", err, alert)
	}

	// Retrieve the DeliveryServiceRequest to check DeliveryServiceRequest name got updated
	resp, _, err = TOSession.GetDeliveryServiceRequestByID(remoteDeliveryServiceRequest.ID)
	if err != nil {
		t.Errorf("cannot GET DeliveryServiceRequest by name: %v - %v\n", firstDeliveryServiceRequest.Name, err)
	}
	respDeliveryServiceRequest := resp[0]
	if respDeliveryServiceRequest.Name != expectedDeliveryServiceRequestName {
		t.Errorf("results do not match actual: %s, expected: %s\n", respDeliveryServiceRequest.Name, expectedDeliveryServiceRequestName)
	}

}

func DeleteTestDeliveryServiceRequests(t *testing.T) {

	secondDeliveryServiceRequest := testData.DeliveryServiceRequests[1]
	resp, _, err := TOSession.DeleteDeliveryServiceRequestByName(secondDeliveryServiceRequest.Name)
	if err != nil {
		t.Errorf("cannot DELETE DeliveryServiceRequest by name: %v - %v\n", err, resp)
	}

	// Retrieve the DeliveryServiceRequest to see if it got deleted
	dsrs, _, err := TOSession.GetDeliveryServiceRequestByName(secondDeliveryServiceRequest.Name)
	if err != nil {
		t.Errorf("error deleting DeliveryServiceRequest name: %s\n", err.Error())
	}
	if len(dsrs) > 0 {
		t.Errorf("expected DeliveryServiceRequest name: %s to be deleted\n", secondDeliveryServiceRequest.Name)
	}
}
*/
