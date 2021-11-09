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

func TestPhysLocations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters, Divisions, Regions, PhysLocations}, func() {
		GetTestPhysLocationsIMS(t)
		GetDefaultSortPhysLocationsTest(t)
		GetSortPhysLocationsTest(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestPhysLocations(t)
		UpdateTestPhysLocations(t)
		UpdateTestPhysLocationsWithHeaders(t, header)
		GetTestPhysLocations(t)
		GetTestPhysLocationsIMSAfterChange(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestPhysLocationsWithHeaders(t, header)
		GetTestPaginationSupportPhysLocation(t)
		CreatePhysLocationWithMismatchedRegionNameAndID(t)
	})
}

func UpdateTestPhysLocationsWithHeaders(t *testing.T, header http.Header) {
	if len(testData.PhysLocations) < 1 {
		t.Fatal("Need at least one Physical Location to test updating Physical Locations, with an HTTP header")
	}
	firstPhysLocation := testData.PhysLocations[0]

	// Retrieve the PhysLocation by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.Header = header
	opts.QueryParameters.Set("name", firstPhysLocation.Name)
	resp, _, err := TOSession.GetPhysLocations(opts)
	if err != nil {
		t.Errorf("cannot get Physical Location by name '%s': %v - alerts: %+v", firstPhysLocation.Name, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected exactly one Physical Location to exist with name '%s', found: %d", firstPhysLocation.Name, len(resp.Response))
	}

	remotePhysLocation := resp.Response[0]
	expectedPhysLocationCity := "city1"
	remotePhysLocation.City = expectedPhysLocationCity
	opts.QueryParameters.Del("name")
	_, reqInf, err := TOSession.UpdatePhysLocation(remotePhysLocation.ID, remotePhysLocation, opts)
	if err == nil {
		t.Errorf("Expected error about precondition failed, but got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestPhysLocationsIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, physLoc := range testData.PhysLocations {
		opts.QueryParameters.Set("name", physLoc.Name)
		resp, reqInf, err := TOSession.GetPhysLocations(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}

}

func GetTestPhysLocationsIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, physLoc := range testData.PhysLocations {
		opts.QueryParameters.Set("name", physLoc.Name)
		resp, reqInf, err := TOSession.GetPhysLocations(opts)
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

	for _, physLoc := range testData.PhysLocations {
		opts.QueryParameters.Set("name", physLoc.Name)
		resp, reqInf, err := TOSession.GetPhysLocations(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestPhysLocations(t *testing.T) {
	for _, pl := range testData.PhysLocations {
		resp, _, err := TOSession.CreatePhysLocation(pl, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Physical Location '%s': %v - alerts: %+v", pl.Name, err, resp.Alerts)
		}
	}

}

func SortTestPhysLocations(t *testing.T) {
	var sortedList []string
	resp, _, err := TOSession.GetPhysLocations(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	for _, pl := range resp.Response {
		sortedList = append(sortedList, pl.Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestPhysLocations(t *testing.T) {
	if len(testData.PhysLocations) < 1 {
		t.Fatal("Need at least one Physical Location to test updating Physical Locations")
	}
	firstPhysLocation := testData.PhysLocations[0]

	// Retrieve the PhysLocation by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstPhysLocation.Name)
	resp, _, err := TOSession.GetPhysLocations(opts)
	if err != nil {
		t.Errorf("cannot get Physical Location by name '%s': %v - alerts: %+v", firstPhysLocation.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Physical Location to exist with name '%s', found: %d", firstPhysLocation.Name, len(resp.Response))
	}

	remotePhysLocation := resp.Response[0]
	expectedPhysLocationCity := "city1"
	remotePhysLocation.City = expectedPhysLocationCity

	alerts, _, err := TOSession.UpdatePhysLocation(remotePhysLocation.ID, remotePhysLocation, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Physical Location: %v - alerts: %+v", err, alerts.Alerts)
	}

	// Retrieve the PhysLocation to check PhysLocation name got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(remotePhysLocation.ID))
	resp, _, err = TOSession.GetPhysLocations(opts)
	if err != nil {
		t.Errorf("cannot Physical Location '%s' (#%d) by ID: %v - alerts: %+v", firstPhysLocation.Name, remotePhysLocation.ID, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Physical Location to exist with ID %d, found: %d", remotePhysLocation.ID, len(resp.Response))
	}
	respPhysLocation := resp.Response[0]
	if respPhysLocation.City != expectedPhysLocationCity {
		t.Errorf("results do not match actual: %s, expected: %s", respPhysLocation.City, expectedPhysLocationCity)
	}

}

func GetTestPhysLocations(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, pl := range testData.PhysLocations {
		opts.QueryParameters.Set("name", pl.Name)
		resp, _, err := TOSession.GetPhysLocations(opts)
		if err != nil {
			t.Errorf("cannot get Physical Location '%s': %v - alerts: %+v", pl.Name, err, resp.Alerts)
		}
	}
}

func GetSortPhysLocationsTest(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetPhysLocations(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Physical Locations ordered by ID: %v - alerts: %+v", err, resp.Alerts)
	}
	sorted := sort.SliceIsSorted(resp.Response, func(i, j int) bool {
		return resp.Response[i].ID < resp.Response[j].ID
	})
	if !sorted {
		t.Error("expected response to be sorted by id")
	}
}

func GetDefaultSortPhysLocationsTest(t *testing.T) {
	resp, _, err := TOSession.GetPhysLocations(client.RequestOptions{})
	if err != nil {
		t.Errorf("Unexpected error getting Physical Locations: %v - alerts: %+v", err, resp.Alerts)
	}
	sorted := sort.SliceIsSorted(resp.Response, func(i, j int) bool {
		return resp.Response[i].Name < resp.Response[j].Name
	})
	if !sorted {
		t.Error("expected response to be sorted by name")
	}
}

func DeleteTestPhysLocations(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, pl := range testData.PhysLocations {
		// Retrieve the PhysLocation by name so we can get the id for the Update
		opts.QueryParameters.Set("name", pl.Name)
		resp, _, err := TOSession.GetPhysLocations(opts)
		if err != nil {
			t.Errorf("cannot get Physical Location by name '%s': %v - alerts: %+v", pl.Name, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Physical Location to exist with name '%s', found: %d", pl.Name, len(resp.Response))
			continue
		}

		respPhysLocation := resp.Response[0]

		alerts, _, err := TOSession.DeletePhysLocation(respPhysLocation.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("Unexpected error deleting Physical Location '%s' (#%d): %v - alerts: %+v", respPhysLocation.Name, respPhysLocation.ID, err, alerts.Alerts)
		}

		// Retrieve the PhysLocation to see if it got deleted
		resp, _, err = TOSession.GetPhysLocations(opts)
		if err != nil {
			t.Errorf("error getting Physical Location '%s' after deletion: %v - alerts: %+v", pl.Name, err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			t.Errorf("expected Physical Location '%s' to be deleted, but it was found in Traffic Ops", pl.Name)
		}
	}
}

func CreatePhysLocationWithMismatchedRegionNameAndID(t *testing.T) {
	resp, _, err := TOSession.GetRegions(client.NewRequestOptions())
	if err != nil {
		t.Errorf("Unexpected error getting regions: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) < 2 {
		t.Fatalf("expected at least two regions to be returned, but got none")
	}
	physLocation := tc.PhysLocation{
		Address:    "100 blah lane",
		City:       "foo",
		Comments:   "comment",
		Email:      "bar@foobar.com",
		Name:       "testPhysicalLocation",
		Phone:      "111-222-3333",
		RegionID:   resp.Response[0].ID,
		RegionName: resp.Response[1].Name,
		ShortName:  "testLocation1",
		State:      "CO",
		Zip:        "80602",
	}
	_, _, err = TOSession.CreatePhysLocation(physLocation, client.NewRequestOptions())
	if err == nil {
		t.Fatalf("expected an error about mismatched region name and ID, but got nothing")
	}

	physLocation.RegionName = resp.Response[0].Name
	_, _, err = TOSession.CreatePhysLocation(physLocation, client.NewRequestOptions())
	if err != nil {
		t.Fatalf("expected no error while creating phys location, but got %v", err)
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "testPhysicalLocation")
	response, _, err := TOSession.GetPhysLocations(opts)
	if err != nil {
		t.Fatalf("cannot get Physical Location by name 'testPhysicalLocation': %v - alerts: %+v", err, resp.Alerts)
	}
	if len(response.Response) != 1 {
		t.Fatalf("Expected exactly one Physical Location to exist with name 'testPhysicalLocation', found: %d", len(resp.Response))
	}

	_, _, err = TOSession.DeletePhysLocation(response.Response[0].ID, client.NewRequestOptions())
	if err != nil {
		t.Errorf("error deleteing physical location 'testPhysicalLocation': %v", err)
	}
}

func GetTestPaginationSupportPhysLocation(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetPhysLocations(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Physical Locations: %v - alerts: %+v", err, resp.Alerts)
	}
	physlocations := resp.Response

	if len(physlocations) > 0 {
		opts.QueryParameters = url.Values{}
		opts.QueryParameters.Set("orderby", "id")
		opts.QueryParameters.Set("limit", "1")
		physlocationsWithLimit, _, err := TOSession.GetPhysLocations(opts)
		if err == nil {
			if !reflect.DeepEqual(physlocations[:1], physlocationsWithLimit.Response) {
				t.Error("expected GET PhysLocation with limit = 1 to return first result")
			}
		} else {
			t.Errorf("Unexpected error getting Physical Locations with a limit: %v - alerts: %+v", err, physlocationsWithLimit.Alerts)
		}
		if len(physlocations) > 1 {
			opts.QueryParameters = url.Values{}
			opts.QueryParameters.Set("orderby", "id")
			opts.QueryParameters.Set("limit", "1")
			opts.QueryParameters.Set("offset", "1")
			physlocationsWithOffset, _, err := TOSession.GetPhysLocations(opts)
			if err == nil {
				if !reflect.DeepEqual(physlocations[1:2], physlocationsWithOffset.Response) {
					t.Error("expected GET PhysLocation with limit = 1, offset = 1 to return second result")
				}
			} else {
				t.Errorf("Unexpected error getting Physical Locations with a limit and an offset: %v - alerts: %+v", err, physlocationsWithOffset.Alerts)
			}

			opts.QueryParameters = url.Values{}
			opts.QueryParameters.Set("orderby", "id")
			opts.QueryParameters.Set("limit", "1")
			opts.QueryParameters.Set("page", "2")
			physlocationsWithPage, _, err := TOSession.GetPhysLocations(opts)
			if err == nil {
				if !reflect.DeepEqual(physlocations[1:2], physlocationsWithPage.Response) {
					t.Error("expected GET PhysLocation with limit = 1, page = 2 to return second result")
				}
			} else {
				t.Errorf("Unexpected error getting Physical Locations with a limit and a page: %v - alerts: %+v", err, physlocationsWithPage.Alerts)
			}
		} else {
			t.Errorf("only one PhysLocation found, so offset functionality can't test")
		}
	} else {
		t.Errorf("No PhysLocation found to check pagination")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetPhysLocations(opts)
	if err == nil {
		t.Error("expected GET PhysLocation to return an error when limit is not bigger than -1")
	} else if !strings.Contains(err.Error(), "must be bigger than -1") {
		t.Errorf("expected GET PhysLocation to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetPhysLocations(opts)
	if err == nil {
		t.Error("expected GET PhysLocation to return an error when offset is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET PhysLocation to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetPhysLocations(opts)
	if err == nil {
		t.Error("expected GET PhysLocation to return an error when page is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET PhysLocation to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}
