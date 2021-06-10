package v4

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
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDivisions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {
		GetTestDivisionsIMS(t)
		TryToDeleteDivision(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestDivisions(t)
		SortTestDivisionDesc(t)
		UpdateTestDivisions(t)
		UpdateTestDivisionsWithHeaders(t, header)
		GetTestDivisionsIMSAfterChange(t, header)
		GetTestDivisions(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestDivisionsWithHeaders(t, header)
		GetTestPaginationSupportDivision(t)
		GetDivisionByInvalidId(t)
		GetDivisionByInvalidName(t)
		DeleteTestDivisionsByInvalidId(t)
	})
}

func UpdateTestDivisionsWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Divisions) < 1 {
		t.Error("Need at least one Division to test updating a Division with HTTP headers")
		return
	}
	firstDivision := testData.Divisions[0]

	// Retrieve the Division by division so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.Header = header
	opts.QueryParameters.Set("name", firstDivision.Name)
	resp, _, err := TOSession.GetDivisions(opts)
	if err != nil {
		t.Errorf("cannot get Division '%s': %v - alerts: %+v", firstDivision.Name, err, resp.Alerts)
	}
	if len(resp.Response) > 0 {
		remoteDivision := resp.Response[0]
		expectedDivision := "division-test"
		remoteDivision.Name = expectedDivision

		opts.QueryParameters.Del("name")
		_, reqInf, err := TOSession.UpdateDivision(remoteDivision.ID, remoteDivision, opts)
		if err == nil {
			t.Errorf("Expected error about precondition failed, but got none")
		}
		if reqInf.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestDivisionsIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, division := range testData.Divisions {
		opts.QueryParameters.Set("name", division.Name)
		resp, reqInf, err := TOSession.GetDivisions(opts)
		if err != nil {
			t.Errorf("could not get Divisions: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}

	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts.Header = make(map[string][]string)
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, division := range testData.Divisions {
		opts.QueryParameters.Set("name", division.Name)
		resp, reqInf, err := TOSession.GetDivisions(opts)
		if err != nil {
			t.Errorf("could not get Divisions: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestDivisionsIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, division := range testData.Divisions {
		opts.QueryParameters.Set("name", division.Name)
		resp, reqInf, err := TOSession.GetDivisions(opts)
		if err != nil {
			t.Errorf("could not get Divisions: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func TryToDeleteDivision(t *testing.T) {
	if len(testData.Divisions) < 1 {
		t.Fatal("Need at least one Division to attempt to delete Divisions")
	}
	division := testData.Divisions[0]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", division.Name)

	resp, _, err := TOSession.GetDivisions(opts)
	if err != nil {
		t.Errorf("cannot get Division '%s': %v - alerts %+v", division.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Division to exist with name '%s', found: %d", division.Name, len(resp.Response))
	}
	division = resp.Response[0]

	alerts, _, err := TOSession.DeleteDivision(division.ID, client.RequestOptions{})
	if err == nil {
		t.Fatal("should not be able to delete a Division prematurely")
	}

	found := false
	for _, alert := range alerts.Alerts {
		if strings.Contains(alert.Text, "Resource not found.") {
			t.Errorf("Division with name '%s' does not exist", division.Name)
		}
		if alert.Level == tc.ErrorLevel.String() && strings.Contains(alert.Text, "cannot delete division because it is being used by a region") {
			found = true
		}
	}
	if !found {
		t.Errorf("unexpected error occured: %v - alerts: %+v", err, alerts)
	}

}

func CreateTestDivisions(t *testing.T) {
	for _, division := range testData.Divisions {
		resp, _, err := TOSession.CreateDivision(division, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Division '%s': %v - alerts: %+v", division.Name, err, resp.Alerts)
		}
	}
}

func SortTestDivisions(t *testing.T) {
	resp, _, err := TOSession.GetDivisions(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}

	sortedList := make([]string, 0, len(resp.Response))
	for _, division := range resp.Response {
		sortedList = append(sortedList, division.Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func SortTestDivisionDesc(t *testing.T) {
	respAsc, _, err := TOSession.GetDivisions(client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected no error, but got error in Division Ascending: %v - alerts: %+v", err, respAsc.Alerts)
	}
	if len(respAsc.Response) < 1 {
		t.Fatal("Need at least one Division in Traffic Ops to test default vs explicit sort order")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("sortOrder", "desc")
	respDesc, _, err := TOSession.GetDivisions(opts)
	if err != nil {
		t.Errorf("Expected no error, but got error in Division Descending: %v - alerts: %+v", err, respDesc.Alerts)
	}
	if len(respDesc.Response) < 1 {
		t.Fatal("Traffic Ops returned at least one Division using default sort order, but zero Divisions using explicit 'desc' sort order")
	}

	if len(respAsc.Response) != len(respDesc.Response) {
		t.Fatalf("Ascending sort response lists %d Divisions, descending lists %d; these should/must be the same", len(respAsc.Response), len(respDesc.Response))
	}

	// TODO: check that the whole thing is sorted, not just the first/last elements?
	// TODO: verify more than one in each response - list of length one is trivially sorted both ascending and descending

	// reverse the descending-sorted response and compare it to the ascending-sorted one
	for start, end := 0, len(respDesc.Response)-1; start < end; start, end = start+1, end-1 {
		respDesc.Response[start], respDesc.Response[end] = respDesc.Response[end], respDesc.Response[start]
	}
	if respDesc.Response[0].Name != respAsc.Response[0].Name {
		t.Errorf("Division responses are not equal after reversal: %s - %s", respDesc.Response[0].Name, respAsc.Response[0].Name)
	}
}

func UpdateTestDivisions(t *testing.T) {
	if len(testData.Divisions) < 1 {
		t.Fatal("Need at least one Division to test updating Divisions")
	}
	firstDivision := testData.Divisions[0]

	// Retrieve the Division by division so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstDivision.Name)
	resp, _, err := TOSession.GetDivisions(opts)
	if err != nil {
		t.Errorf("cannot get Division '%s': %v - alerts: %+v", firstDivision.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Division to exist with name '%s', found: %d", firstDivision.Name, len(resp.Response))
	}
	remoteDivision := resp.Response[0]
	expectedDivision := "division-test"
	remoteDivision.Name = expectedDivision

	alert, _, err := TOSession.UpdateDivision(remoteDivision.ID, remoteDivision, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Division '%s' (#%d): %v - alerts: %+v", firstDivision.Name, remoteDivision.ID, err, alert.Alerts)
	}

	// Retrieve the Division to check division got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(remoteDivision.ID))
	resp, _, err = TOSession.GetDivisions(opts)
	if err != nil {
		t.Errorf("cannot get Division #%d: %v - alerts: %+v", remoteDivision.ID, err, resp.Alerts)
	}
	if len(resp.Response) > 0 {
		respDivision := resp.Response[0]
		if respDivision.Name != expectedDivision {
			t.Errorf("results do not match actual: %s, expected: %s", respDivision.Name, expectedDivision)
		}

		// Set the name back to the fixture value so we can delete it after
		remoteDivision.Name = firstDivision.Name
		alert, _, err = TOSession.UpdateDivision(remoteDivision.ID, remoteDivision, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot update Division #%d: %v - alerts: %+v", remoteDivision.ID, err, alert.Alerts)
		}
	}
}

func GetTestDivisions(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, division := range testData.Divisions {
		opts.QueryParameters.Set("name", division.Name)
		resp, _, err := TOSession.GetDivisions(opts)
		if err != nil {
			t.Errorf("cannot get Division '%s': %v - alerts: %+v", division.Name, err, resp)
		}
	}
}

func DeleteTestDivisions(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, division := range testData.Divisions {
		// Retrieve the Division by name so we can get the id
		opts.QueryParameters.Set("name", division.Name)
		resp, _, err := TOSession.GetDivisions(opts)
		if err != nil {
			t.Errorf("cannot get Division '%s': %v - alerts: %+v", division.Name, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Division to exist with the name '%s', found: %d", division.Name, len(resp.Response))
			continue
		}
		respDivision := resp.Response[0]

		delResp, _, err := TOSession.DeleteDivision(respDivision.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot DELETE Division '%s' (#%d): %v - alerts: %+v", division.Name, respDivision.ID, err, delResp.Alerts)
		}

		// Retrieve the Division to see if it got deleted
		resp, _, err = TOSession.GetDivisions(opts)
		if err != nil {
			t.Errorf("error fetching Division '%s' after deletion: %v - alerts: %+v", division.Name, err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			t.Errorf("expected Division : %s to be deleted, but it was returned by Traffic Ops", division.Name)
		}
	}
}

func GetTestPaginationSupportDivision(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetDivisions(opts)
	if err != nil {
		t.Fatalf("cannot get Divisions: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) < 2 {
		t.Fatalf("Need at least 2 Divisions to test Division pagination, only found: %d", len(resp.Response))
	}
	divisions := resp.Response

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	divisionsWithLimit, _, _ := TOSession.GetDivisions(opts)
	if !reflect.DeepEqual(divisions[:1], divisionsWithLimit.Response) {
		t.Error("expected GET Divisions with limit = 1 to return first result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "1")
	divisionsWithOffset, _, _ := TOSession.GetDivisions(opts)
	if !reflect.DeepEqual(divisions[1:2], divisionsWithOffset.Response) {
		t.Error("expected GET Divisions with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "2")
	divisionsWithPage, _, _ := TOSession.GetDivisions(opts)
	if !reflect.DeepEqual(divisions[1:2], divisionsWithPage.Response) {
		t.Error("expected GET Divisions with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetDivisions(opts)
	if err == nil {
		t.Error("expected GET Divisions to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET Divisions to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetDivisions(opts)
	if err == nil {
		t.Error("expected GET Divisions to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Divisions to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetDivisions(opts)
	if err == nil {
		t.Error("expected GET Divisions to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Divisions to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}

func GetDivisionByInvalidId(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("id", "10000")
	resp, _, err := TOSession.GetDivisions(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Division by presumably invalid ID (10000): %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) >= 1 {
		t.Errorf("Expected to find exactly zero Divisions with presumably invalid ID (10000), found: %d", len(resp.Response))
	}
}

func GetDivisionByInvalidName(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "abcd")
	resp, _, err := TOSession.GetDivisions(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Division by presumably invalid name ('abcd'): %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) >= 1 {
		t.Errorf("Expected to find exactly zero Divisions with presumably invalid name ('abcd'), found: %d", len(resp.Response))
	}
}

func DeleteTestDivisionsByInvalidId(t *testing.T) {
	delResp, _, err := TOSession.DeleteDivision(10000, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected an error deleting Division with presumably invalid ID (10000), didn't get one - alerts: %+v", delResp.Alerts)
	}
}
