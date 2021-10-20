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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
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
		UpdateTestRegions(t)
		UpdateTestRegionsWithHeaders(t, header)
		GetTestRegions(t)
		GetTestRegionsIMSAfterChange(t, header)
		DeleteTestRegionsByName(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestRegionsWithHeaders(t, header)
	})
}

func UpdateTestRegionsWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Regions) > 0 {
		firstRegion := testData.Regions[0]
		// Retrieve the Region by region so we can get the id for the Update
		resp, _, err := TOSession.GetRegionByNameWithHdr(firstRegion.Name, header)
		if err != nil {
			t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
		}
		if len(resp) > 0 {
			remoteRegion := resp[0]
			remoteRegion.Name = "OFFLINE-TEST"
			_, reqInf, err := TOSession.UpdateRegionByIDWithHdr(remoteRegion.ID, remoteRegion, header)
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
		_, reqInf, err := TOSession.GetRegionByNameWithHdr(region.Name, header)
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
		_, reqInf, err := TOSession.GetRegionByNameWithHdr(region.Name, header)
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
		_, reqInf, err := TOSession.GetRegionByNameWithHdr(region.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
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
	var header http.Header
	var sortedList []string
	resp, _, err := TOSession.GetRegionsWithHdr(header)
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

func UpdateTestRegions(t *testing.T) {

	firstRegion := testData.Regions[0]
	// Retrieve the Region by region so we can get the id for the Update
	resp, _, err := TOSession.GetRegionByName(firstRegion.Name)
	if err != nil {
		t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
	}
	remoteRegion := resp[0]
	expectedRegion := "OFFLINE-TEST"
	remoteRegion.Name = expectedRegion
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateRegionByID(remoteRegion.ID, remoteRegion)
	if err != nil {
		t.Errorf("cannot UPDATE Region by id: %v - %v", err, alert)
	}

	// Retrieve the Region to check region got updated
	resp, _, err = TOSession.GetRegionByID(remoteRegion.ID)
	if err != nil {
		t.Errorf("cannot GET Region by region: %v - %v", firstRegion.Name, err)
	}
	respRegion := resp[0]
	if respRegion.Name != expectedRegion {
		t.Errorf("results do not match actual: %s, expected: %s", respRegion.Name, expectedRegion)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRegion.Name = firstRegion.Name
	alert, _, err = TOSession.UpdateRegionByID(remoteRegion.ID, remoteRegion)
	if err != nil {
		t.Errorf("cannot UPDATE Region by id: %v - %v", err, alert)
	}

}

func GetTestRegions(t *testing.T) {
	for _, region := range testData.Regions {
		resp, _, err := TOSession.GetRegionByName(region.Name)
		if err != nil {
			t.Errorf("cannot GET Region by region: %v - %v", err, resp)
		}
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
		regionResp, _, err := TOSession.GetRegionByName(region.Name)
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
		resp, _, err := TOSession.GetRegionByName(region.Name)
		if err != nil {
			t.Errorf("cannot GET Region by name: %v - %v", region.Name, err)
		}
		respRegion := resp[0]

		delResp, _, err := TOSession.DeleteRegionByID(respRegion.ID)
		if err != nil {
			t.Errorf("cannot DELETE Region by region: %v - %v", err, delResp)
		}

		// Retrieve the Region to see if it got deleted
		regionResp, _, err := TOSession.GetRegionByName(region.Name)
		if err != nil {
			t.Errorf("error deleting Region region: %s", err.Error())
		}
		if len(regionResp) > 0 {
			t.Errorf("expected Region : %s to be deleted", region.Name)
		}
	}
}
