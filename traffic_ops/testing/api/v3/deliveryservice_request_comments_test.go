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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestDeliveryServiceRequestComments(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Tenants, DeliveryServiceRequests, DeliveryServiceRequestComments}, func() {
		GetTestDeliveryServiceRequestCommentsIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfUnmodifiedSince, time)
		header.Set(rfc.IfModifiedSince, time)
		SortTestDeliveryServiceRequestComments(t)
		UpdateTestDeliveryServiceRequestComments(t)
		UpdateTestDeliveryServiceRequestCommentsWithHeaders(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestDeliveryServiceRequestCommentsWithHeaders(t, header)
		GetTestDeliveryServiceRequestComments(t)
		GetTestDeliveryServiceRequestCommentsIMSAfterChange(t, header)
	})
}

func UpdateTestDeliveryServiceRequestCommentsWithHeaders(t *testing.T, header http.Header) {
	comments, _, _ := TOSession.GetDeliveryServiceRequestCommentsWithHdr(header)

	if len(comments) > 0 {
		firstComment := comments[0]
		newFirstCommentValue := "new comment value"
		firstComment.Value = newFirstCommentValue

		_, reqInf, err := TOSession.UpdateDeliveryServiceRequestCommentByIDWithHdr(firstComment.ID, firstComment, header)
		if err == nil {
			t.Errorf("expected precondition failed error, but got none")
		}
		if reqInf.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestDeliveryServiceRequestCommentsIMSAfterChange(t *testing.T, header http.Header) {
	_, reqInf, err := TOSession.GetDeliveryServiceRequestCommentsWithHdr(header)
	if err != nil {
		t.Fatalf("could not GET delivery service request comments: %v", err)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, err = TOSession.GetDeliveryServiceRequestCommentsWithHdr(header)
	if err != nil {
		t.Fatalf("could not GET delivery service request comments: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestDeliveryServiceRequestComments(t *testing.T) {

	// Retrieve a delivery service request by xmlId so we can get the ID needed to create a dsr comment
	dsr := testData.DeliveryServiceRequests[0].DeliveryService

	resp, _, err := TOSession.GetDeliveryServiceRequestByXMLID(dsr.XMLID)
	if err != nil {
		t.Errorf("cannot GET delivery service request by xml id: %v - %v", dsr.XMLID, err)
	}
	if len(resp) != 1 {
		t.Errorf("found %d delivery service request by xml id, expected %d: %s", len(resp), 1, dsr.XMLID)
	} else {
		respDSR := resp[0]

		for _, comment := range testData.DeliveryServiceRequestComments {
			comment.DeliveryServiceRequestID = respDSR.ID
			resp, _, err := TOSession.CreateDeliveryServiceRequestComment(comment)
			if err != nil {
				t.Errorf("could not CREATE delivery service request comment: %v - %v", err, resp)
			}
		}
	}

}

func SortTestDeliveryServiceRequestComments(t *testing.T) {
	var header http.Header
	var sortedList []string
	resp, _, err := TOSession.GetDeliveryServiceRequestCommentsWithHdr(header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	for i, _ := range resp {
		sortedList = append(sortedList, resp[i].XMLID)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
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
		t.Errorf("cannot UPDATE delivery service request comment by id: %v - %v", err, alert)
	}

	// Retrieve the delivery service request comment to check that the value got updated
	resp, _, err := TOSession.GetDeliveryServiceRequestCommentByID(firstComment.ID)
	if err != nil {
		t.Errorf("cannot GET delivery service request comment by id: '$%d', %v", firstComment.ID, err)
	}
	respDSRC := resp[0]
	if respDSRC.Value != newFirstCommentValue {
		t.Errorf("results do not match actual: %s, expected: %s", respDSRC.Value, newFirstCommentValue)
	}

}

func GetTestDeliveryServiceRequestCommentsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, err := TOSession.GetDeliveryServiceRequestCommentsWithHdr(header)
	if err != nil {
		t.Fatalf("could not GET delivery service request comments: %v", err)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestDeliveryServiceRequestComments(t *testing.T) {

	comments, _, _ := TOSession.GetDeliveryServiceRequestComments()

	for _, comment := range comments {
		resp, _, err := TOSession.GetDeliveryServiceRequestCommentByID(comment.ID)
		if err != nil {
			t.Errorf("cannot GET delivery service request comment by id: %v - %v", err, resp)
		}
	}
}

func DeleteTestDeliveryServiceRequestComments(t *testing.T) {

	comments, _, _ := TOSession.GetDeliveryServiceRequestComments()

	for _, comment := range comments {
		_, _, err := TOSession.DeleteDeliveryServiceRequestCommentByID(comment.ID)
		if err != nil {
			t.Errorf("cannot DELETE delivery service request comment by id: '%d' %v", comment.ID, err)
		}

		// Retrieve the delivery service request comment to see if it got deleted
		comments, _, err := TOSession.GetDeliveryServiceRequestCommentByID(comment.ID)
		if err != nil {
			t.Errorf("error deleting delivery service request comment: %s", err.Error())
		}
		if len(comments) > 0 {
			t.Errorf("expected delivery service request comment: %d to be deleted", comment.ID)
		}
	}
}
