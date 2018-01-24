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
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
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

func GetTestDeliveryServiceRequests(t *testing.T) {

	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID("xmlId")
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
