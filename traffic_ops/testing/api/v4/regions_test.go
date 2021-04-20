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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
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
	if len(testData.Regions) > 0 {
		firstRegion := testData.Regions[0]
		// Retrieve the Region by region so we can get the id for the Update
		resp, _, err := TOSession.GetRegionByName(firstRegion.Name, header)
		if err != nil {
			t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
		}
		if len(resp) > 0 {
			remoteRegion := resp[0]
			remoteRegion.Name = "OFFLINE-TEST"
			_, reqInf, err := TOSession.UpdateRegion(remoteRegion.ID, remoteRegion, header)
			if err == nil {
				t.Errorf("Expected error about precondition failed, but got none")
			}
			if reqInf.StatusCode != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		}
	}
}

func GetTestRegionsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, region := range testData.Regions {
		_, reqInf, err := TOSession.GetRegionByName(region.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestRegionsIMSAfterChange(t *testing.T, header http.Header) {
	for _, region := range testData.Regions {
		_, reqInf, err := TOSession.GetRegionByName(region.Name, header)
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
	for _, region := range testData.Regions {
		_, reqInf, err := TOSession.GetRegionByName(region.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestRegions(t *testing.T) {
	for _, region := range testData.Regions {
		resp, _, err := TOSession.GetRegionByName(region.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Region by region: %v - %v", err, resp)
		}
	}
}

func CreateTestRegions(t *testing.T) {

	for _, region := range testData.Regions {
		resp, _, err := TOSession.CreateRegion(region)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE region: %v", err)
		}
	}
}

func SortTestRegions(t *testing.T) {
	var sortedList []string
	resp, _, err := TOSession.GetRegions(nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	for i := range resp {
		sortedList = append(sortedList, resp[i].Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func SortTestRegionsDesc(t *testing.T) {

	respAsc, _, err1 := TOSession.GetRegions(nil, nil)
	params := url.Values{}
	params.Set("sortOrder", "desc")
	respDesc, _, err2 := TOSession.GetRegions(params, nil)

	if err1 != nil {
		t.Errorf("Expected no error, but got error in Regions Ascending %v", err1)
	}
	if err2 != nil {
		t.Errorf("Expected no error, but got error in Regions Descending %v", err2)
	}
	if len(respAsc) == len(respDesc) {
		if len(respAsc) > 0 && len(respDesc) > 0 {
			// reverse the descending-sorted response and compare it to the ascending-sorted one
			for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
				respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
			}
			if respDesc[0].Name != "" && respAsc[0].Name != "" {
				if !reflect.DeepEqual(respDesc[0].Name, respAsc[0].Name) {
					t.Errorf("Regions responses are not equal after reversal: Asc: %s - Desc: %s", respDesc[0].Name, respAsc[0].Name)
				}
			}
		} else {
			t.Errorf("No Response returned from GET Regions using SortOrder")
		}
	} else {
		t.Fatalf("Region response length are not equal Asc: %d Desc: %d", len(respAsc), len(respDesc))
	}
}

func UpdateTestRegions(t *testing.T) {

	firstRegion := testData.Regions[0]
	// Retrieve the Region by region so we can get the id for the Update
	resp, _, err := TOSession.GetRegionByName(firstRegion.Name, nil)
	if err != nil {
		t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
	}
	remoteRegion := resp[0]
	expectedRegion := "OFFLINE-TEST"
	remoteRegion.Name = expectedRegion
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateRegion(remoteRegion.ID, remoteRegion, nil)
	if err != nil {
		t.Errorf("cannot UPDATE Region by id: %v - %v", err, alert)
	}

	// Retrieve the Region to check region got updated
	resp, _, err = TOSession.GetRegionByID(remoteRegion.ID, nil)
	if err != nil {
		t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
	}
	respRegion := resp[0]
	if respRegion.Name != expectedRegion {
		t.Errorf("results do not match actual: %s, expected: %s", respRegion.Name, expectedRegion)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRegion.Name = firstRegion.Name
	alert, _, err = TOSession.UpdateRegion(remoteRegion.ID, remoteRegion, nil)
	if err != nil {
		t.Errorf("cannot UPDATE Region by id: %v - %v", err, alert)
	}
}

func VerifyPaginationSupportRegion(t *testing.T) {

	qparams := url.Values{}
	qparams.Set("orderby", "id")
	regions, _, err := TOSession.GetRegions(qparams, nil)
	if err != nil {
		t.Fatalf("cannot GET Regions: %v", err)
	}

	if len(regions) > 0 {
		qparams = url.Values{}
		qparams.Set("orderby", "id")
		qparams.Set("limit", "1")
		regionsWithLimit, _, err := TOSession.GetRegions(qparams, nil)
		if err == nil {
			if !reflect.DeepEqual(regions[:1], regionsWithLimit) {
				t.Error("expected GET Regions with limit = 1 to return first result")
			}
		} else {
			t.Error("Error in getting regions by limit")
		}

		if len(regions) > 1 {
			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("offset", "1")
			regionsWithOffset, _, err := TOSession.GetRegions(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(regions[1:2], regionsWithOffset) {
					t.Error("expected GET Regions with limit = 1, offset = 1 to return second result")
				}
			} else {
				t.Error("Error in getting regions by limit and offset")
			}

			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("page", "2")
			regionsWithPage, _, err := TOSession.GetRegions(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(regions[1:2], regionsWithPage) {
					t.Error("expected GET Regions with limit = 1, page = 2 to return second result")
				}
			} else {
				t.Error("Error in getting regions by limit and page")
			}
		} else {
			t.Errorf("only one region found, so offset functionality can't test")
		}
	} else {
		t.Errorf("No region found to check pagination")
	}

	qparams = url.Values{}
	qparams.Set("limit", "-2")
	_, _, err = TOSession.GetRegions(qparams, nil)
	if err == nil {
		t.Error("expected GET Regions to return an error when limit is not bigger than -1")
	} else if !strings.Contains(err.Error(), "must be bigger than -1") {
		t.Errorf("expected GET Regions to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("offset", "0")
	_, _, err = TOSession.GetRegions(qparams, nil)
	if err == nil {
		t.Error("expected GET Regions to return an error when offset is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Regions to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("page", "0")
	_, _, err = TOSession.GetRegions(qparams, nil)
	if err == nil {
		t.Error("expected GET Regions to return an error when page is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Regions to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}

func DeleteTestRegionsByName(t *testing.T) {
	for _, region := range testData.Regions {
		delResp, _, err := TOSession.DeleteRegion(nil, &region.Name)
		if err != nil {
			t.Errorf("cannot DELETE Region by name: %v - %v", err, delResp)
		}

		deleteLog, _, err := TOSession.GetLogsByLimit(1)
		if err != nil {
			t.Fatalf("unable to get latest audit log entry")
		}
		if len(deleteLog) != 1 {
			t.Fatalf("log entry length - expected: 1, actual: %d", len(deleteLog))
		}
		if !strings.Contains(*deleteLog[0].Message, region.Name) {
			t.Errorf("region deletion audit log entry - expected: message containing region name '%s', actual: %s", region.Name, *deleteLog[0].Message)
		}

		// Retrieve the Region to see if it got deleted
		regionResp, _, err := TOSession.GetRegionByName(region.Name, nil)
		if err != nil {
			t.Errorf("error deleting Region region: %s", err.Error())
		}
		if len(regionResp) > 0 {
			t.Errorf("expected Region : %s to be deleted", region.Name)
		}
	}
	CreateTestRegions(t)
}

func DeleteTestRegions(t *testing.T) {

	for _, region := range testData.Regions {
		// Retrieve the Region by name so we can get the id
		resp, _, err := TOSession.GetRegionByName(region.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Region by name: %v - %v", region.Name, err)
		}
		respRegion := resp[0]

		delResp, _, err := TOSession.DeleteRegion(&respRegion.ID, nil)
		if err != nil {
			t.Errorf("cannot DELETE Region by region: %v - %v", err, delResp)
		}

		// Retrieve the Region to see if it got deleted
		regionResp, _, err := TOSession.GetRegionByName(region.Name, nil)
		if err != nil {
			t.Errorf("error deleting Region region: %s", err.Error())
		}
		if len(regionResp) > 0 {
			t.Errorf("expected Region : %s to be deleted", region.Name)
		}
	}
}

func DeleteTestRegionsByInvalidId(t *testing.T) {
	i := 10000
	delResp, _, err := TOSession.DeleteRegion(&i, nil)
	if err == nil {
		t.Errorf("cannot DELETE Regions by Invalid ID: %v - %v", err, delResp)
	}
}

func DeleteTestRegionsByInvalidName(t *testing.T) {
	i := "invalid"
	delResp, _, err := TOSession.DeleteRegion(nil, &i)
	if err == nil {
		t.Errorf("cannot DELETE Regions by Invalid ID: %v - %v", err, delResp)
	}
}

func GetTestRegionByInvalidId(t *testing.T) {
	regionResp, _, err := TOSession.GetRegionByID(10000, nil)
	if err != nil {
		t.Errorf("Error!! Getting Region by Invalid ID %v", err)
	}
	if len(regionResp) >= 1 {
		t.Errorf("Error!! Invalid ID shouldn't have any response %v Error %v", regionResp, err)
	}
}

func GetTestRegionByInvalidName(t *testing.T) {
	regionResp, _, err := TOSession.GetRegionByName("abcd", nil)
	if err != nil {
		t.Errorf("Error!! Getting Region by Invalid Name %v", err)
	}
	if len(regionResp) >= 1 {
		t.Errorf("Error!! Invalid Name shouldn't have any response %v Error %v", regionResp, err)
	}
}

func GetTestRegionByDivision(t *testing.T) {
	for _, region := range testData.Regions {

		resp, _, err := TOSession.GetDivisionByName(region.DivisionName, nil)
		if err != nil {
			t.Errorf("cannot GET Division by name: %v - %v", region.DivisionName, err)
		}
		respDivision := resp[0]

		_, reqInf, err := TOSession.GetRegionByDivision(respDivision.ID, nil)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestRegionByInvalidDivision(t *testing.T) {
	regionResp, _, err := TOSession.GetRegionByDivision(100000, nil)
	if err != nil {
		t.Errorf("Getting Region by Invalid Divisions %v", err)
	}
	if len(regionResp) >= 1 {
		t.Errorf("Invalid Division shouldn't have any response %v Error %v", regionResp, err)
	}
}

func CreateTestRegionsInvalidDivision(t *testing.T) {

	firstRegion := testData.Regions[0]
	firstRegion.Division = 100
	firstRegion.Name = "abcd"
	_, _, err := TOSession.CreateRegion(firstRegion)
	if err == nil {
		t.Errorf("Expected division not found Error %v", err)
	}
}
