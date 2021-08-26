package v5

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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
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
	if len(testData.Statuses) < 1 {
		t.Fatal("Need at least one Status to test updating a status with an HTTP header")
	}

	firstStatus := testData.Statuses[0]
	if firstStatus.Name == nil {
		t.Fatal("cannot update test statuses: first test data status must have a name")
	}

	// Retrieve the Status by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.Header = header
	opts.QueryParameters.Set("name", *firstStatus.Name)
	resp, _, err := TOSession.GetStatuses(opts)
	if err != nil {
		t.Errorf("cannot get Status by name '%s': %v - alerts %+v", *firstStatus.Name, err, resp.Alerts)
	}
	if len(resp.Response) > 0 {
		remoteStatus := resp.Response[0]
		expectedStatusDesc := "new description"
		remoteStatus.Description = expectedStatusDesc

		opts.QueryParameters.Del("name")
		_, reqInf, err := TOSession.UpdateStatus(remoteStatus.ID, remoteStatus, opts)
		if err == nil {
			t.Errorf("Expected error about precondition failed, but got none")
		}
		if reqInf.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("Expected status code 412, got %d", reqInf.StatusCode)
		}
	}
}

func GetTestStatusesIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}

		opts.QueryParameters.Set("name", *status.Name)
		resp, reqInf, err := TOSession.GetStatuses(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)

	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}
		opts.QueryParameters.Set("name", *status.Name)
		resp, reqInf, err := TOSession.GetStatuses(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestStatusesIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)

	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get test Statuses: test data Statuses must have names")
		}
		opts.QueryParameters.Set("name", *status.Name)
		resp, reqInf, err := TOSession.GetStatuses(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestStatuses(t *testing.T) {
	for _, status := range testData.Statuses {
		resp, _, err := TOSession.CreateStatus(status, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Status: %v - alerts: %+v", err, resp.Alerts)
		}
	}

}

func SortTestStatuses(t *testing.T) {
	resp, _, err := TOSession.GetStatuses(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}

	sortedList := make([]string, 0, len(resp.Response))
	for _, status := range resp.Response {
		sortedList = append(sortedList, status.Name)
	}

	if !sort.StringsAreSorted(sortedList) {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestStatuses(t *testing.T) {
	if len(testData.Statuses) < 1 {
		t.Fatal("Need at least one Status to test updating a Status")
	}
	firstStatus := testData.Statuses[0]
	if firstStatus.Name == nil {
		t.Fatal("cannot update test statuses: first test data status must have a name")
	}

	// Retrieve the Status by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstStatus.Name)
	resp, _, err := TOSession.GetStatuses(opts)
	if err != nil {
		t.Errorf("cannot get Status by name '%s': %v - alerts: %+v", *firstStatus.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Status to exist with name '%s', found: %d", *firstStatus.Name, len(resp.Response))
	}
	remoteStatus := resp.Response[0]
	expectedStatusDesc := "new description"
	remoteStatus.Description = expectedStatusDesc

	alert, _, err := TOSession.UpdateStatus(remoteStatus.ID, remoteStatus, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Status: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Status to check Status name got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(remoteStatus.ID))
	resp, _, err = TOSession.GetStatuses(opts)
	if err != nil {
		t.Errorf("cannot get Status '%s' by ID %d: %v - alerts: %+v", *firstStatus.Name, remoteStatus.ID, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Status to exist with ID %d, found: %d", remoteStatus.ID, len(resp.Response))
	}
	respStatus := resp.Response[0]
	if respStatus.Description != expectedStatusDesc {
		t.Errorf("results do not match actual: %s, expected: %s", respStatus.Name, expectedStatusDesc)
	}

}

func GetTestStatuses(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get ftest statuses: test data statuses must have names")
		}
		opts.QueryParameters.Set("name", *status.Name)
		resp, _, err := TOSession.GetStatuses(opts)
		if err != nil {
			t.Errorf("cannot get Status by name: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func DeleteTestStatuses(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, status := range testData.Statuses {
		if status.Name == nil {
			t.Fatal("cannot get test statuses: test data statuses must have names")
		}

		// Retrieve the Status by name so we can get the id for the Update
		opts.QueryParameters.Set("name", *status.Name)
		resp, _, err := TOSession.GetStatuses(opts)
		if err != nil {
			t.Errorf("cannot get Statuses filtered by name '%s': %v - alerts: %+v", *status.Name, err, resp.Alerts)
		}
		respStatus := resp.Response[0]

		delResp, _, err := TOSession.DeleteStatus(respStatus.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Status: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Status to see if it got deleted
		resp, _, err = TOSession.GetStatuses(opts)
		if err != nil {
			t.Errorf("Unexpected error getting Statuses filtered by name after deletion: %v - alerts: %+v", err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			t.Errorf("expected Status '%s' to be deleted, but it was found in Traffic Ops", *status.Name)
		}
	}
}
