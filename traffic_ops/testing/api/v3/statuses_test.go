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
	tc "github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestStatuses(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Statuses}, func() {
		GetTestStatusesIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestStatuses(t)
		UpdateTestStatuses(t)
		UpdateTestStatusesWithHeaders(t, header)
		GetTestStatuses(t)
		GetTestStatusesIMSAfterChange(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestStatusesWithHeaders(t, header)
	})
}

func UpdateTestStatusesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Statuses) > 0 {
		firstStatus := testData.Statuses[0]
		if firstStatus.Name == nil {
			t.Fatal("cannot update test statuses: first test data status must have a name")
		}

		// Retrieve the Status by name so we can get the id for the Update
		resp, _, err := TOSession.GetStatusByNameWithHdr(*firstStatus.Name, header)
		if err != nil {
			t.Errorf("cannot GET Status by name: %v - %v", firstStatus.Name, err)
		}
		if len(resp) > 0 {
			remoteStatus := resp[0]
			expectedStatusDesc := "new description"
			remoteStatus.Description = expectedStatusDesc
			_, reqInf, err := TOSession.UpdateStatusByIDWithHdr(remoteStatus.ID, remoteStatus, header)
			if err == nil {
				t.Errorf("Expected error about precondition failed, but got none")
			}
			if reqInf.StatusCode != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		}
	}
}

func GetTestStatusesIMSAfterChange(t *testing.T, header http.Header) {
	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}
		_, reqInf, err := TOSession.GetStatusByNameWithHdr(*status.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}
		_, reqInf, err := TOSession.GetStatusByNameWithHdr(*status.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestStatusesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}
		_, reqInf, err := TOSession.GetStatusByNameWithHdr(*status.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestStatuses(t *testing.T) {

	for _, status := range testData.Statuses {
		resp, _, err := TOSession.CreateStatusNullable(status)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE types: %v", err)
		}
	}

}

func SortTestStatuses(t *testing.T) {
	var header http.Header
	var sortedList []string
	resp, _, err := TOSession.GetStatusesWithHdr(header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	for i, _ := range resp {
		sortedList = append(sortedList, resp[i].Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestStatuses(t *testing.T) {

	firstStatus := testData.Statuses[0]
	if firstStatus.Name == nil {
		t.Fatal("cannot update test statuses: first test data status must have a name")
	}

	// Retrieve the Status by name so we can get the id for the Update
	resp, _, err := TOSession.GetStatusByName(*firstStatus.Name)
	if err != nil {
		t.Errorf("cannot GET Status by name: %v - %v", firstStatus.Name, err)
	}
	remoteStatus := resp[0]
	expectedStatusDesc := "new description"
	remoteStatus.Description = expectedStatusDesc
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateStatusByID(remoteStatus.ID, remoteStatus)
	if err != nil {
		t.Errorf("cannot UPDATE Status by id: %v - %v", err, alert)
	}

	// Retrieve the Status to check Status name got updated
	resp, _, err = TOSession.GetStatusByID(remoteStatus.ID)
	if err != nil {
		t.Errorf("cannot GET Status by ID: %v - %v", firstStatus.Description, err)
	}
	respStatus := resp[0]
	if respStatus.Description != expectedStatusDesc {
		t.Errorf("results do not match actual: %s, expected: %s", respStatus.Name, expectedStatusDesc)
	}

}

func GetTestStatuses(t *testing.T) {

	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestStatuses(t *testing.T) {

	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}

		// Retrieve the Status by name so we can get the id for the Update
		resp, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("cannot GET Status by name: %v - %v", status.Name, err)
		}
		respStatus := resp[0]

		delResp, _, err := TOSession.DeleteStatusByID(respStatus.ID)
		if err != nil {
			t.Errorf("cannot DELETE Status by name: %v - %v", err, delResp)
		}

		// Retrieve the Status to see if it got deleted
		types, _, err := TOSession.GetStatusByName(*status.Name)
		if err != nil {
			t.Errorf("error deleting Status name: %s", err.Error())
		}
		if len(types) > 0 {
			t.Errorf("expected Status name: %s to be deleted", *status.Name)
		}
	}
}
