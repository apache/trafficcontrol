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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	tc "github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCoordinates(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Coordinates}, func() {
		GetTestCoordinatesIMS(t)
		GetTestCoordinates(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestCoordinates(t)
		SortTestCoordinatesDesc(t)
		UpdateTestCoordinates(t)
		UpdateTestCoordinatesWithHeaders(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestCoordinatesWithHeaders(t, header)
		GetTestCoordinatesIMSAfterChange(t, header)
		GetTestCoordinatesByInvalidId(t)
		GetTestCoordiantesByInvalidName(t)
		GetTestPaginationSupportCoordinates(t)
		CreateTestCoordinatesWithInvalidName(t)
		CreateTestCoordinatesWithInvalidLatitude(t)
		CreateTestCoordinatesWithInvalidLogitude(t)
		UpdateTestCoordinatesByInvalidId(t)
		DeleteTestCoordinatesByInvalidId(t)
	})
}

func UpdateTestCoordinatesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Coordinates) < 1 {
		t.Error("Need at least one Coordinate to test updating a Coordinate with an HTTP header")
		return
	}
	firstCoord := testData.Coordinates[0]

	opts := client.NewRequestOptions()
	opts.Header = header
	opts.QueryParameters.Set("name", firstCoord.Name)

	resp, _, err := TOSession.GetCoordinates(opts)
	if err != nil {
		t.Errorf("cannot get Coordinate '%s' from Traffic Ops: %v - alerts: %+v", firstCoord.Name, err, resp.Alerts)
	}
	if len(resp.Response) > 0 {
		coord := resp.Response[0]
		expectedLat := 12.34
		coord.Latitude = expectedLat

		opts.QueryParameters.Del("name")
		_, reqInf, err := TOSession.UpdateCoordinate(coord.ID, coord, opts)
		if err == nil {
			t.Errorf("Expected error about precondition failed, but got none")
		}
		if reqInf.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
		}
	} else {
		t.Errorf("No Coordinates available to update")
	}
}

func GetTestCoordinatesIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, coord := range testData.Coordinates {
		opts.QueryParameters.Set("name", coord.Name)
		resp, reqInf, err := TOSession.GetCoordinates(opts)
		if err != nil {
			t.Errorf("could not get Coordinate '%s' from Traffic Ops: %v - alerts: %+v", coord.Name, err, resp.Alerts)
			return
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
			return
		}
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)

	opts.Header.Set(rfc.IfModifiedSince, timeStr)

	for _, coord := range testData.Coordinates {
		opts.QueryParameters.Set("name", coord.Name)
		resp, reqInf, err := TOSession.GetCoordinates(opts)
		if err != nil {
			t.Fatalf("could not get Coordinate '%s' from Traffic Ops: %v - alerts: %+v", coord.Name, err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCoordinatesIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, coord := range testData.Coordinates {
		opts.QueryParameters.Set("name", coord.Name)
		resp, reqInf, err := TOSession.GetCoordinates(opts)
		if err != nil {
			t.Errorf("Unexpected error getting Coordinate '%s': %v - alerts: %+v", coord.Name, err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestCoordinates(t *testing.T) {
	for _, coord := range testData.Coordinates {
		resp, _, err := TOSession.CreateCoordinate(coord, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create coordinate: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func GetTestCoordinates(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, coord := range testData.Coordinates {
		opts.QueryParameters.Set("name", coord.Name)
		resp, _, err := TOSession.GetCoordinates(opts)
		if err != nil {
			t.Errorf("cannot get Coordinate '%s' from Traffic Ops: %v - alerts: %v", coord.Name, err, resp.Alerts)
		}
	}
}

func SortTestCoordinates(t *testing.T) {
	var sortedList []string
	resp, _, err := TOSession.GetCoordinates(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Unexpected error getting Coordinates from Traffic Ops: %v - alerts: %+v", err, resp.Alerts)
	}
	for _, coord := range resp.Response {
		sortedList = append(sortedList, coord.Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func SortTestCoordinatesDesc(t *testing.T) {

	resp, _, err := TOSession.GetCoordinates(client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting Coordinates with default sort order: %v - alerts: %+v", err, resp.Alerts)
	}
	respAsc := resp.Response
	if len(respAsc) < 1 {
		t.Fatal("Need at least one Coordinate in Traffic Ops to test sort ordering")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("sortOrder", "desc")
	resp, _, err = TOSession.GetCoordinates(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Coordinates with explicit descending sort order: %v - alerts: %+v", err, resp.Alerts)
	}
	respDesc := resp.Response
	if len(respDesc) < 1 {
		t.Fatal("Need at least one Coordinate in Traffic Ops to test sort ordering")
	}

	if len(respAsc) != len(respDesc) {
		t.Fatalf("Traffic Ops returned %d Coordinates using default sort order, but returned %d Coordinates using explicit descending sort order", len(respAsc), len(respDesc))
	}

	// reverse the descending-sorted response and compare it to the ascending-sorted one
	for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
		respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
	}
	if len(respDesc[0].Name) > 0 && len(respAsc[0].Name) > 0 {
		if respDesc[0].Name != respAsc[0].Name {
			t.Errorf("Coordinates responses are not equal after reversal: %s - %s", respDesc[0].Name, respAsc[0].Name)
		}
	} else {
		t.Errorf("Coordinates name shouldn't be empty while sorting the response")
	}
}

func UpdateTestCoordinates(t *testing.T) {
	if len(testData.Coordinates) < 1 {
		t.Fatal("Need at least one Coordinate to test updating Coordinates")
	}
	firstCoord := testData.Coordinates[0]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstCoord.Name)
	resp, _, err := TOSession.GetCoordinates(opts)
	if err != nil {
		t.Errorf("cannot get Coordinate '%s' by name: %v - alerts: %+v", firstCoord.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Coordinate to exist with name '%s', found: %d", firstCoord.Name, len(resp.Response))
	}
	coord := resp.Response[0]
	expectedLat := 12.34
	coord.Latitude = expectedLat

	alert, _, err := TOSession.UpdateCoordinate(coord.ID, coord, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Coordinate: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Coordinate to check Coordinate name got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(coord.ID))
	resp, _, err = TOSession.GetCoordinates(opts)
	if err != nil {
		t.Errorf("cannot get Coordinate '%s' by id: %v - alerts: %+v", firstCoord.Name, err, resp.Alerts)
	}
	if len(resp.Response) > 0 {
		coord = resp.Response[0]
		if coord.Latitude != expectedLat {
			t.Errorf("results do not match actual: %s, expected: %f", coord.Name, expectedLat)
		}
	} else {
		t.Errorf("Can't retrieve coordinates to check the updated value")
	}
}

func DeleteTestCoordinates(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, coord := range testData.Coordinates {
		// Retrieve the Coordinate by name so we can get the id for the Update
		opts.QueryParameters.Set("name", coord.Name)
		resp, _, err := TOSession.GetCoordinates(opts)
		if err != nil {
			t.Errorf("cannot get Coordinate '%s' from Traffic Ops: %v - alerts: %+v", coord.Name, err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			respCoord := resp.Response[0]
			delResp, _, err := TOSession.DeleteCoordinate(respCoord.ID, client.RequestOptions{})
			if err != nil {
				t.Errorf("cannot delete Coordinate '%s' (#%d): %v - alerts: %+v", respCoord.Name, respCoord.ID, err, delResp.Alerts)
			}
			// Retrieve the Coordinate to see if it got deleted
			coords, _, err := TOSession.GetCoordinates(opts)
			if err != nil {
				t.Errorf("Unexpected error fetching Coordinate '%s' (#%d) after deletion: %v - alerts: %+v", coord.Name, respCoord.ID, err, coords.Alerts)
			}
			if len(coords.Response) > 0 {
				t.Errorf("expected Coordinate '%s' (#%d) to be deleted, but found it in Traffic Ops after deletion", coord.Name, respCoord.ID)
			}
		} else {
			t.Errorf("No Coordinates available to delete")
		}
	}
}

func GetTestCoordinatesByInvalidId(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("id", "10000")
	coordinatesResp, _, err := TOSession.GetCoordinates(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Coordinates by presumably invalid ID: %v - alerts: %+v", err, coordinatesResp.Alerts)
	}
	if len(coordinatesResp.Response) >= 1 {
		t.Errorf("Didn't expect to find a Coordinate with a presumably invalid ID in Traffic Ops response: %v", coordinatesResp.Response)
	}
}

func GetTestCoordiantesByInvalidName(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "abcd")
	coordinatesResp, _, err := TOSession.GetCoordinates(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Coordinates filtered by presumably non-existent name: %v - alerts: %+v", err, coordinatesResp.Alerts)
	}
	if len(coordinatesResp.Response) >= 1 {
		t.Errorf("Didn't expect to find Coordinate with presumably non-existent name in Traffic Ops response: %v", coordinatesResp.Response)
	}
}

func GetTestPaginationSupportCoordinates(t *testing.T) {

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetCoordinates(opts)
	if err != nil {
		t.Fatalf("cannot get Coordinates: %v - alerts: %+v", err, resp.Alerts)
	}
	coordinates := resp.Response
	if len(coordinates) < 2 {
		t.Fatal("Need at least two Coordinates in Traffic Ops to test Coordinate pagination")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	coordinatesWithLimit, _, err := TOSession.GetCoordinates(opts)
	if err == nil {
		if !reflect.DeepEqual(coordinates[:1], coordinatesWithLimit.Response) {
			t.Error("expected GET Coordinates with limit = 1 to return first result")
		}
	} else {
		t.Error("Error in getting coordinates by limit")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "1")
	coordinatesWithOffset, _, err := TOSession.GetCoordinates(opts)
	if err == nil {
		if !reflect.DeepEqual(coordinates[1:2], coordinatesWithOffset.Response) {
			t.Error("expected GET Coordinates with limit = 1, offset = 1 to return second result")
		}
	} else {
		t.Error("Error in getting coordinates by limit and offset")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "2")
	coordinatesWithPage, _, err := TOSession.GetCoordinates(opts)
	if err == nil {
		if !reflect.DeepEqual(coordinates[1:2], coordinatesWithPage.Response) {
			t.Error("expected GET Coordinates with limit = 1, page = 2 to return second result")
		}
	} else {
		t.Error("Error in getting coordinates by limit and page")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetCoordinates(opts)
	if err == nil {
		t.Error("expected GET Coordinates to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET Coordinates to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetCoordinates(opts)
	if err == nil {
		t.Error("expected GET Coordinates to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Coordinates to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetCoordinates(opts)
	if err == nil {
		t.Error("expected GET Coordinates to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Coordinates to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}

func CreateTestCoordinatesWithInvalidName(t *testing.T) {
	if len(testData.Coordinates) < 1 {
		t.Fatal("No Coordinates available to fetch")
	}
	firstCoordinates := testData.Coordinates[0]
	firstCoordinates.Name = ""
	_, reqInf, err := TOSession.CreateCoordinate(firstCoordinates, client.RequestOptions{})
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("Expected an error creating a Coordinate with an invalid name, but didn't get one")
	}
}

func CreateTestCoordinatesWithInvalidLatitude(t *testing.T) {
	if len(testData.Coordinates) < 1 {
		t.Fatal("No Coordinates available to fetch")
	}
	firstCoordinates := testData.Coordinates[0]
	firstCoordinates.Latitude = 20000
	_, reqInf, err := TOSession.CreateCoordinate(firstCoordinates, client.RequestOptions{})
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("Expected an error creating a Coordinate with an invalid Latitude, but didn't get one")
	}
}

func CreateTestCoordinatesWithInvalidLogitude(t *testing.T) {
	if len(testData.Coordinates) < 1 {
		t.Fatal("No Coordinates available to fetch")
	}
	firstCoordinates := testData.Coordinates[0]
	firstCoordinates.Longitude = 20000
	_, reqInf, err := TOSession.CreateCoordinate(firstCoordinates, client.RequestOptions{})
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
	}
	if err == nil {
		t.Errorf("Expected an error creating a Coordinate with an invalid Latitude, but didn't get one")
	}
}

func UpdateTestCoordinatesByInvalidId(t *testing.T) {
	if len(testData.Coordinates) < 1 {
		t.Fatal("No Coordinates available to update")
	}
	firstCoord := testData.Coordinates[0]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstCoord.Name)
	resp, reqInf, err := TOSession.GetCoordinates(opts)
	if err != nil {
		t.Errorf("cannot get Coordinate '%s' by name: %v - alerts: %+v", firstCoord.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Coordinate to exist with name '%s', found: %d", firstCoord.Name, len(resp.Response))
	}
	coord := resp.Response[0]
	expectedLat := 12.34
	coord.Latitude = expectedLat

	var alert tc.Alerts
	alert, reqInf, err = TOSession.UpdateCoordinate(10000, coord, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected an error updating a Coordinate with a presumably non-existent ID, but didn't get one - alerts: %+v", alert.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 status code, got %v", reqInf.StatusCode)
	}
}

func DeleteTestCoordinatesByInvalidId(t *testing.T) {
	alerts, reqInf, err := TOSession.DeleteCoordinate(12345, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected an error deleting a Coordinate with a presumably non-existent ID, but didn't get one - alerts: %+v", alerts.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 status code, got %v", reqInf.StatusCode)
	}
}
