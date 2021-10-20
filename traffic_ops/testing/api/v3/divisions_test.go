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
		UpdateTestDivisions(t)
		UpdateTestDivisionsWithHeaders(t, header)
		GetTestDivisionsIMSAfterChange(t, header)
		GetTestDivisions(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestDivisionsWithHeaders(t, header)
	})
}

func UpdateTestDivisionsWithHeaders(t *testing.T, header http.Header) {
	firstDivision := testData.Divisions[0]
	// Retrieve the Division by division so we can get the id for the Update
	resp, _, err := TOSession.GetDivisionByNameWithHdr(firstDivision.Name, header)
	if err != nil {
		t.Errorf("cannot GET Division by division: %v - %v", firstDivision.Name, err)
	}
	if len(resp) > 0 {
		remoteDivision := resp[0]
		expectedDivision := "division-test"
		remoteDivision.Name = expectedDivision

		_, reqInf, err := TOSession.UpdateDivisionByIDWithHdr(remoteDivision.ID, remoteDivision, header)
		if err == nil {
			t.Errorf("Expected error about precondition failed, but got none")
		}
		if reqInf.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestDivisionsIMSAfterChange(t *testing.T, header http.Header) {
	for _, division := range testData.Divisions {
		_, reqInf, err := TOSession.GetDivisionByNameWithHdr(division.Name, header)
		if err != nil {
			t.Fatalf("could not GET divisions: %v", err)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, division := range testData.Divisions {
		_, reqInf, err := TOSession.GetDivisionByNameWithHdr(division.Name, header)
		if err != nil {
			t.Fatalf("could not GET divisions: %v", err)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestDivisionsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, division := range testData.Divisions {
		_, reqInf, err := TOSession.GetDivisionByNameWithHdr(division.Name, header)
		if err != nil {
			t.Fatalf("could not GET divisions: %v", err)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func TryToDeleteDivision(t *testing.T) {
	division := testData.Divisions[0]

	resp, _, err := TOSession.GetDivisionByName(division.Name)
	if err != nil {
		t.Errorf("cannot GET Division by name: %v - %v", division.Name, err)
	}
	division = resp[0]
	_, _, err = TOSession.DeleteDivisionByID(division.ID)

	if err == nil {
		t.Error("should not be able to delete a division prematurely")
		return
	}

	if strings.Contains(err.Error(), "Resource not found.") {
		t.Errorf("division with name %v does not exist", division.Name)
		return
	}

	if strings.Contains(err.Error(), "cannot delete division because it is being used by a region") {
		return
	}

	t.Errorf("unexpected error occured: %v", err)
}

func CreateTestDivisions(t *testing.T) {
	for _, division := range testData.Divisions {
		resp, _, err := TOSession.CreateDivision(division)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE division: %v", err)
		}
	}
}

func SortTestDivisions(t *testing.T) {
	var header http.Header
	var sortedList []string
	resp, _, err := TOSession.GetDivisionsWithHdr(header)
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

func UpdateTestDivisions(t *testing.T) {

	firstDivision := testData.Divisions[0]
	// Retrieve the Division by division so we can get the id for the Update
	resp, _, err := TOSession.GetDivisionByName(firstDivision.Name)
	if err != nil {
		t.Errorf("cannot GET Division by division: %v - %v", firstDivision.Name, err)
	}
	remoteDivision := resp[0]
	expectedDivision := "division-test"
	remoteDivision.Name = expectedDivision
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDivisionByID(remoteDivision.ID, remoteDivision)
	if err != nil {
		t.Errorf("cannot UPDATE Division by id: %v - %v", err, alert)
	}

	// Retrieve the Division to check division got updated
	resp, _, err = TOSession.GetDivisionByID(remoteDivision.ID)
	if err != nil {
		t.Errorf("cannot GET Division by division: %v - %v", firstDivision.Name, err)
	}
	if len(resp) > 0 {
		respDivision := resp[0]
		if respDivision.Name != expectedDivision {
			t.Errorf("results do not match actual: %s, expected: %s", respDivision.Name, expectedDivision)
		}

		// Set the name back to the fixture value so we can delete it after
		remoteDivision.Name = firstDivision.Name
		alert, _, err = TOSession.UpdateDivisionByID(remoteDivision.ID, remoteDivision)
		if err != nil {
			t.Errorf("cannot UPDATE Division by id: %v - %v", err, alert)
		}
	}
}

func GetTestDivisions(t *testing.T) {
	for _, division := range testData.Divisions {
		resp, _, err := TOSession.GetDivisionByName(division.Name)
		if err != nil {
			t.Errorf("cannot GET Division by division: %v - %v", err, resp)
		}
	}
}

func DeleteTestDivisions(t *testing.T) {

	for _, division := range testData.Divisions {
		// Retrieve the Division by name so we can get the id
		resp, _, err := TOSession.GetDivisionByName(division.Name)
		if err != nil {
			t.Errorf("cannot GET Division by name: %v - %v", division.Name, err)
		}
		respDivision := resp[0]

		delResp, _, err := TOSession.DeleteDivisionByID(respDivision.ID)
		if err != nil {
			t.Errorf("cannot DELETE Division by division: %v - %v", err, delResp)
		}

		// Retrieve the Division to see if it got deleted
		divisionResp, _, err := TOSession.GetDivisionByName(division.Name)
		if err != nil {
			t.Errorf("error deleting Division division: %s", err.Error())
		}
		if len(divisionResp) > 0 {
			t.Errorf("expected Division : %s to be deleted", division.Name)
		}
	}
}
