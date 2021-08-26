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
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	client "github.com/apache/trafficcontrol/traffic_ops/v5-client"
)

func TestRegions(t *testing.T) {
	WithObjs(t, []TCObj{Parameters, Divisions, Regions}, func() {
		GetTestRegionsIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestRegions(t)
		SortTestRegionsDesc(t)
		UpdateTestRegions(t)
		UpdateTestRegionsWithHeaders(t, header)
		GetTestRegions(t)
		GetTestRegionsIMSAfterChange(t, header)
		DeleteTestRegionsByName(t)
		VerifyPaginationSupportRegion(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestRegionsWithHeaders(t, header)
		DeleteTestRegionsByInvalidId(t)
		DeleteTestRegionsByInvalidName(t)
		GetTestRegionByInvalidId(t)
		GetTestRegionByInvalidName(t)
		GetTestRegionByDivision(t)
		GetTestRegionByInvalidDivision(t)
		CreateTestRegionsInvalidDivision(t)
	})
}

func UpdateTestRegionsWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Regions) < 1 {
		t.Fatal("Need at least one Region to test updating Regions with HTTP headers")
	}
	firstRegion := testData.Regions[0]

	// Retrieve the Region by region so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstRegion.Name)
	opts.Header = header
	resp, _, err := TOSession.GetRegions(opts)
	if err != nil {
		t.Errorf("cannot get Region '%s' by name: %v - alerts: %+v", firstRegion.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Region to exist with name '%s', found: %d", firstRegion.Name, len(resp.Response))
	}

	remoteRegion := resp.Response[0]
	remoteRegion.Name = "OFFLINE-TEST"

	opts.QueryParameters.Del("name")
	_, reqInf, err := TOSession.UpdateRegion(remoteRegion.ID, remoteRegion, opts)
	if err == nil {
		t.Errorf("Expected error about precondition failed, but got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestRegionsIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, region := range testData.Regions {
		opts.QueryParameters.Set("name", region.Name)
		resp, reqInf, err := TOSession.GetRegions(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestRegionsIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, region := range testData.Regions {
		opts.QueryParameters.Set("name", region.Name)
		resp, reqInf, err := TOSession.GetRegions(opts)
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

	for _, region := range testData.Regions {
		opts.QueryParameters.Set("name", region.Name)
		resp, reqInf, err := TOSession.GetRegions(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestRegions(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, region := range testData.Regions {
		opts.QueryParameters.Set("name", region.Name)
		resp, _, err := TOSession.GetRegions(opts)
		if err != nil {
			t.Errorf("cannot get Region '%s' by name: %v - alerts: %+v", region.Name, err, resp.Alerts)
		}
	}
}

func CreateTestRegions(t *testing.T) {
	for _, region := range testData.Regions {
		resp, _, err := TOSession.CreateRegion(region, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Region '%s': %v - alerts: %+v", region.Name, err, resp.Alerts)
		}
	}
}

func SortTestRegions(t *testing.T) {
	resp, _, err := TOSession.GetRegions(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}

	sortedList := make([]string, 0, len(resp.Response))
	for _, region := range resp.Response {
		sortedList = append(sortedList, region.Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if !res {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func SortTestRegionsDesc(t *testing.T) {
	resp, _, err := TOSession.GetRegions(client.RequestOptions{})
	if err != nil {
		t.Errorf("Expected no error, but got error in Regions with default ordering: %v - alerts: %+v", err, resp.Alerts)
	}
	respAsc := resp.Response
	if len(respAsc) < 1 {
		t.Fatal("Need at least one Region in Traffic Ops to test Regions sort ordering")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("sortOrder", "desc")
	resp, _, err = TOSession.GetRegions(opts)
	if err != nil {
		t.Errorf("Expected no error, but got error in Regions with Descending ordering: %v - alerts: %+v", err, resp.Alerts)
	}
	respDesc := resp.Response
	if len(respDesc) < 1 {
		t.Fatal("Need at least one Region in Traffic Ops to test Regions sort ordering")
	}

	if len(respAsc) != len(respDesc) {
		t.Fatalf("Traffic Ops returned %d Regions using default sort order, but %d Regions when sort order was explicitly set to descending", len(respAsc), len(respDesc))
	}

	// reverse the descending-sorted response and compare it to the ascending-sorted one
	// TODO ensure at least two in each slice? A list of length one is
	// trivially sorted both ascending and descending.
	for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
		respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
	}
	if respDesc[0].Name != respAsc[0].Name {
		t.Errorf("Regions responses are not equal after reversal: Asc: %s - Desc: %s", respDesc[0].Name, respAsc[0].Name)
	}
}

func UpdateTestRegions(t *testing.T) {
	if len(testData.Regions) < 1 {
		t.Fatal("Need at least one Region to test updating a Region")
	}
	firstRegion := testData.Regions[0]

	// Retrieve the Region by region so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstRegion.Name)
	resp, _, err := TOSession.GetRegions(opts)
	if err != nil {
		t.Errorf("cannot get Region '%s' by name: %v - alerts: %+v", firstRegion.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Region to exist with name '%s', found: %d", firstRegion.Name, len(resp.Response))
	}

	remoteRegion := resp.Response[0]
	expectedRegion := "OFFLINE-TEST"
	remoteRegion.Name = expectedRegion

	alert, _, err := TOSession.UpdateRegion(remoteRegion.ID, remoteRegion, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Region: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Region to check region got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(remoteRegion.ID))
	resp, _, err = TOSession.GetRegions(opts)
	if err != nil {
		t.Errorf("cannot get Region '%s' (#%d) by ID: %v - alerts: %+v", firstRegion.Name, remoteRegion.ID, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Region to exist with ID %d, found: %d", remoteRegion.ID, len(resp.Response))
	}
	respRegion := resp.Response[0]
	if respRegion.Name != expectedRegion {
		t.Errorf("results do not match actual: %s, expected: %s", respRegion.Name, expectedRegion)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRegion.Name = firstRegion.Name
	alert, _, err = TOSession.UpdateRegion(remoteRegion.ID, remoteRegion, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Region: %v - alerts: %+v", err, alert.Alerts)
	}
}

func VerifyPaginationSupportRegion(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetRegions(opts)
	if err != nil {
		t.Fatalf("cannot get Regions: %v - alerts: %+v", err, resp.Alerts)
	}
	regions := resp.Response
	if len(regions) < 2 {
		t.Fatalf("Need at least 2 Regions in Traffic Ops to test pagination support, found: %d", len(regions))
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	regionsWithLimit, _, err := TOSession.GetRegions(opts)
	if err == nil {
		if !reflect.DeepEqual(regions[:1], regionsWithLimit.Response) {
			t.Error("expected GET Regions with limit = 1 to return first result")
		}
	} else {
		t.Error("Error in getting regions by limit")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "1")
	regionsWithOffset, _, err := TOSession.GetRegions(opts)
	if err == nil {
		if !reflect.DeepEqual(regions[1:2], regionsWithOffset.Response) {
			t.Error("expected GET Regions with limit = 1, offset = 1 to return second result")
		}
	} else {
		t.Error("Error in getting regions by limit and offset")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "2")
	regionsWithPage, _, err := TOSession.GetRegions(opts)
	if err == nil {
		if !reflect.DeepEqual(regions[1:2], regionsWithPage.Response) {
			t.Error("expected GET Regions with limit = 1, page = 2 to return second result")
		}
	} else {
		t.Error("Error in getting regions by limit and page")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetRegions(opts)
	if err == nil {
		t.Error("expected GET Regions to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET Regions to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetRegions(opts)
	if err == nil {
		t.Error("expected GET Regions to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Regions to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetRegions(opts)
	if err == nil {
		t.Error("expected GET Regions to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET Regions to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}

func DeleteTestRegionsByName(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, region := range testData.Regions {
		delResp, _, err := TOSession.DeleteRegion(region.Name, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Region '%s': %v - alerts: %+v", region.Name, err, delResp.Alerts)
		}

		opts.QueryParameters.Set("limit", "1")
		deleteLog, _, err := TOSession.GetLogs(opts)
		opts.QueryParameters.Del("limit")
		if err != nil {
			t.Fatalf("unable to get latest audit log entry")
		}
		if len(deleteLog.Response) != 1 {
			t.Fatalf("log entry length - expected: 1, actual: %d", len(deleteLog.Response))
		}
		if deleteLog.Response[0].Message == nil {
			t.Fatal("Traffic Ops returned a representation for a log entry with null or undefined message")
		}
		if !strings.Contains(*deleteLog.Response[0].Message, region.Name) {
			t.Errorf("region deletion audit log entry - expected: message containing region name '%s', actual: %s", region.Name, *deleteLog.Response[0].Message)
		}

		// Retrieve the Region to see if it got deleted
		opts.QueryParameters.Set("name", region.Name)
		regionResp, _, err := TOSession.GetRegions(opts)
		opts.QueryParameters.Del("name")
		if err != nil {
			t.Errorf("error deleting Region '%s': %v - alerts: %+v", region.Name, err, regionResp.Alerts)
		}
		if len(regionResp.Response) > 0 {
			t.Errorf("expected Region '%s' to be deleted, but it was found in Traffic Ops", region.Name)
		}
	}
	CreateTestRegions(t)
}

func DeleteTestRegions(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, region := range testData.Regions {
		// Retrieve the Region by name so we can get the id
		opts.QueryParameters.Set("name", region.Name)
		resp, _, err := TOSession.GetRegions(opts)
		opts.QueryParameters.Del("name")
		if err != nil {
			t.Errorf("cannot get Region '%s' by name: %v - alerts: %+v", region.Name, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Region to exist with name '%s', found: %d", region.Name, len(resp.Response))
		}
		respRegion := resp.Response[0]

		opts.QueryParameters.Set("id", strconv.Itoa(respRegion.ID))
		delResp, _, err := TOSession.DeleteRegion("", opts)
		opts.QueryParameters.Del("id")
		if err != nil {
			t.Errorf("cannot delete Region: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Region to see if it got deleted
		opts.QueryParameters.Set("name", region.Name)
		regionResp, _, err := TOSession.GetRegions(opts)
		if err != nil {
			t.Errorf("error fetching Region '%s' after deletion: %v - alerts: %+v", region.Name, err, regionResp.Alerts)
		}
		if len(regionResp.Response) > 0 {
			t.Errorf("expected Region '%s' to be deleted, but it was found in Traffic Ops", region.Name)
		}
	}
}

func DeleteTestRegionsByInvalidId(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("id", "10000")
	delResp, _, err := TOSession.DeleteRegion("", opts)
	if err == nil {
		t.Errorf("cannot delete Regions by invalid ID: %v - alerts: %+v", err, delResp.Alerts)
	}
}

func DeleteTestRegionsByInvalidName(t *testing.T) {
	delResp, _, err := TOSession.DeleteRegion("invalid", client.RequestOptions{})
	if err == nil {
		t.Errorf("cannot delete Regions by invalid name: %v - alerts: %+v", err, delResp.Alerts)
	}
}

func GetTestRegionByInvalidId(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("id", "10000")
	regionResp, _, err := TOSession.GetRegions(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Regions filtered by presumably non-existent ID: %v - alerts: %+v", err, regionResp.Alerts)
	}
	if len(regionResp.Response) >= 1 {
		t.Errorf("Didn't expect to find any Regions with presumably non-existent ID, found: %d", len(regionResp.Response))
	}
}

func GetTestRegionByInvalidName(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", "abcd")
	regionResp, _, err := TOSession.GetRegions(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Regions filtered by presumably non-existent name: %v - alerts: %+v", err, regionResp.Alerts)
	}
	if len(regionResp.Response) >= 1 {
		t.Errorf("Didn't expect to find any Regions with presumably non-existent name, found: %d", len(regionResp.Response))
	}
}

func GetTestRegionByDivision(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, region := range testData.Regions {
		opts.QueryParameters.Set("name", region.DivisionName)
		resp, _, err := TOSession.GetDivisions(opts)
		opts.QueryParameters.Del("name")
		if err != nil {
			t.Errorf("cannot get Division '%s' by name: %v - alerts: %+v", region.DivisionName, err, resp.Alerts)
		}
		if len(resp.Response) != 1 {
			t.Errorf("Expected exactly one Division to exist with name '%s', found: %d", region.DivisionName, len(resp.Response))
			continue
		}
		respDivision := resp.Response[0]

		opts.QueryParameters.Set("division", strconv.Itoa(respDivision.ID))
		regionsResp, reqInf, err := TOSession.GetRegions(opts)
		opts.QueryParameters.Del("division")
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, regionsResp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestRegionByInvalidDivision(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("division", "100000")
	regionResp, _, err := TOSession.GetRegions(opts)
	if err != nil {
		t.Errorf("Unexpected error getting Regions filtered by presumably non-existent Division ID: %v - alerts: %+v", err, regionResp.Alerts)
	}
	if len(regionResp.Response) >= 1 {
		t.Errorf("Didn't expect to find any Regions in presumably non-existent Division, found: %d", len(regionResp.Response))
	}
}

func CreateTestRegionsInvalidDivision(t *testing.T) {
	if len(testData.Regions) < 1 {
		t.Fatal("Need at least one Region to test creating an invalid Region")
	}
	firstRegion := testData.Regions[0]
	firstRegion.Division = 100
	firstRegion.Name = "abcd"
	_, _, err := TOSession.CreateRegion(firstRegion, client.RequestOptions{})
	if err == nil {
		t.Error("Expected an error creating a presumably invalid Region (name: 'abcd', Division ID: 100), but didn't get one")
	}
}
