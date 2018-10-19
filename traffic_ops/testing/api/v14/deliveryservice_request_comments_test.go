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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestDeliveryServiceRequestComments(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestDeliveryServiceRequests(t)
	CreateTestDeliveryServiceRequestComments(t)
	UpdateTestDeliveryServiceRequestComments(t)
	GetTestDeliveryServiceRequestComments(t)
	DeleteTestDeliveryServiceRequestComments(t)
	DeleteTestDeliveryServiceRequests(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

func CreateTestDeliveryServiceRequestComments(t *testing.T) {

	// Retrieve a delivery service request by xmlId so we can get the ID needed to create a dsr comment
	dsr := testData.DeliveryServiceRequests[0].DeliveryService

	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.XMLID)
	if err != nil {
		t.Errorf("cannot GET delivery service request by xml id: %v - %v\n", dsr.XMLID, err)
	}
	respDSR := resp[0]

	for _, comment := range testData.DeliveryServiceRequestComments {
		comment.DeliveryServiceRequestID = respDSR.ID
		resp, _, err := TOSession.CreateDeliveryServiceRequestComment(comment)
		if err != nil {
			t.Errorf("could not CREATE delivery service request comment: %v - %v\n", err, resp)
		}
	}

}

func UpdateTestDeliveryServiceRequestComments(t *testing.T) {

	comments, _, err := TOSession.GetDeliveryServiceRequestComments()

	firstComment := comments[0]
	newFirstCommentValue := "new comment value"
	firstComment.Value = newFirstCommentValue

	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDeliveryServiceRequestCommentByID(firstComment.ID, firstComment)
	if err != nil {
		t.Errorf("cannot UPDATE delivery service request comment by id: %v - %v\n", err, alert)
	}

	// Retrieve the delivery service request comment to check that the value got updated
	resp, _, err := TOSession.GetDeliveryServiceRequestCommentByID(firstComment.ID)
	if err != nil {
		t.Errorf("cannot GET delivery service request comment by id: '$%d', %v\n", firstComment.ID, err)
	}
	respDSRC := resp[0]
	if respDSRC.Value != newFirstCommentValue {
		t.Errorf("results do not match actual: %s, expected: %s\n", respDSRC.Value, newFirstCommentValue)
	}

}

func GetTestDeliveryServiceRequestComments(t *testing.T) {

	comments, _, _ := TOSession.GetDeliveryServiceRequestComments()

	for _, comment := range comments {
		resp, _, err := TOSession.GetDeliveryServiceRequestCommentByID(comment.ID)
		if err != nil {
			t.Errorf("cannot GET delivery service request comment by id: %v - %v\n", err, resp)
		}
	}
}

func DeleteTestDeliveryServiceRequestComments(t *testing.T) {

	comments, _, _ := TOSession.GetDeliveryServiceRequestComments()

	for _, comment := range comments {
		_, _, err := TOSession.DeleteDeliveryServiceRequestCommentByID(comment.ID)
		if err != nil {
			t.Errorf("cannot DELETE delivery service request comment by id: '%d' %v\n", comment.ID, err)
		}

		// Retrieve the delivery service request comment to see if it got deleted
		comments, _, err := TOSession.GetDeliveryServiceRequestCommentByID(comment.ID)
		if err != nil {
			t.Errorf("error deleting delivery service request comment: %s\n", err.Error())
		}
		if len(comments) > 0 {
			t.Errorf("expected delivery service request comment: %d to be deleted\n", comment.ID)
		}
	}
}
