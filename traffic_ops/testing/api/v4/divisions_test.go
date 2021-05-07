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
	firstDivision := testData.Divisions[0]
	// Retrieve the Division by division so we can get the id for the Update
	resp, _, err := TOSession.GetDivisionByName(firstDivision.Name, header)
	if err != nil {
		t.Errorf("cannot GET Division by division: %v - %v", firstDivision.Name, err)
	}
	if len(resp) > 0 {
		remoteDivision := resp[0]
		expectedDivision := "division-test"
		remoteDivision.Name = expectedDivision

		_, reqInf, err := TOSession.UpdateDivision(remoteDivision.ID, remoteDivision, header)
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
		_, reqInf, err := TOSession.GetDivisionByName(division.Name, header)
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
		_, reqInf, err := TOSession.GetDivisionByName(division.Name, header)
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
		_, reqInf, err := TOSession.GetDivisionByName(division.Name, header)
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

	resp, _, err := TOSession.GetDivisionByName(division.Name, nil)
	if err != nil {
		t.Errorf("cannot GET Division by name: %v - %v", division.Name, err)
	}
	division = resp[0]
	_, _, err = TOSession.DeleteDivision(division.ID)

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
	var sortedList []string
	resp, _, err := TOSession.GetDivisions(nil, nil)
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

func SortTestDivisionDesc(t *testing.T) {

	respAsc, _, err1 := TOSession.GetDivisions(nil, nil)
	params := url.Values{}
	params.Set("sortOrder", "desc")
	respDesc, _, err2 := TOSession.GetDivisions(params, nil)

	if err1 != nil {
		t.Errorf("Expected no error, but got error in Division Ascending %v", err1)
	}
	if err2 != nil {
		t.Errorf("Expected no error, but got error in Division Descending %v", err2)
	}

	if len(respAsc) == len(respDesc) {
		if len(respAsc) > 0 && len(respDesc) > 0 {
			// reverse the descending-sorted response and compare it to the ascending-sorted one
			for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
				respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
			}
			if respDesc[0].Name != "" && respAsc[0].Name != "" {
				if !reflect.DeepEqual(respDesc[0].Name, respAsc[0].Name) {
					t.Errorf("Division responses are not equal after reversal: Asc: %s - Desc: %s", respDesc[0].Name, respAsc[0].Name)
				}
			}
		} else {
			t.Errorf("No Response returned from GET Division using SortOrder")
		}
	} else {
		t.Fatalf("Division response length are not equal Asc: %d Desc: %d", len(respAsc), len(respDesc))
	}
}

func UpdateTestDivisions(t *testing.T) {

	firstDivision := testData.Divisions[0]
	// Retrieve the Division by division so we can get the id for the Update
	resp, _, err := TOSession.GetDivisionByName(firstDivision.Name, nil)
	if err != nil {
		t.Errorf("cannot GET Division by division: %v - %v", firstDivision.Name, err)
	}
	remoteDivision := resp[0]
	expectedDivision := "division-test"
	remoteDivision.Name = expectedDivision
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateDivision(remoteDivision.ID, remoteDivision, nil)
	if err != nil {
		t.Errorf("cannot UPDATE Division by id: %v - %v", err, alert)
	}

	// Retrieve the Division to check division got updated
	resp, _, err = TOSession.GetDivisionByID(remoteDivision.ID, nil)
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
		alert, _, err = TOSession.UpdateDivision(remoteDivision.ID, remoteDivision, nil)
		if err != nil {
			t.Errorf("cannot UPDATE Division by id: %v - %v", err, alert)
		}
	}
}

func GetTestDivisions(t *testing.T) {
	for _, division := range testData.Divisions {
		resp, _, err := TOSession.GetDivisionByName(division.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Division by division: %v - %v", err, resp)
		}
	}
}

func DeleteTestDivisions(t *testing.T) {

	for _, division := range testData.Divisions {
		// Retrieve the Division by name so we can get the id
		resp, _, err := TOSession.GetDivisionByName(division.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Division by name: %v - %v", division.Name, err)
		}
		respDivision := resp[0]

		delResp, _, err := TOSession.DeleteDivision(respDivision.ID)
		if err != nil {
			t.Errorf("cannot DELETE Division by division: %v - %v", err, delResp)
		}

		// Retrieve the Division to see if it got deleted
		divisionResp, _, err := TOSession.GetDivisionByName(division.Name, nil)
		if err != nil {
			t.Errorf("error deleting Division division: %s", err.Error())
		}
		if len(divisionResp) > 0 {
			t.Errorf("expected Division : %s to be deleted", division.Name)
		}
	}
}

func GetTestPaginationSupportDivision(t *testing.T) {
	qparams := url.Values{}
	qparams.Set("orderby", "id")
	divisions, _, err := TOSession.GetDivisions(qparams, nil)
	if err != nil {
		t.Fatalf("cannot GET Divisions: %v", err)
	}

	if len(divisions) > 0 {
		qparams = url.Values{}
		qparams.Set("orderby", "id")
		qparams.Set("limit", "1")
		divisionsWithLimit, _, err := TOSession.GetDivisions(qparams, nil)
		if err == nil {
			if !reflect.DeepEqual(divisions[:1], divisionsWithLimit) {
				t.Error("expected GET Divisions with limit = 1 to return first result")
			}
		} else {
			t.Error("Error in getting division by limit")
		}
		if len(divisions) > 1 {
			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("offset", "1")
			divisionsWithOffset, _, err := TOSession.GetDivisions(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(divisions[1:2], divisionsWithOffset) {
					t.Error("expected GET Divisions with limit = 1, offset = 1 to return second result")
				}
			} else {
				t.Error("Error in getting division by limit and offset")
			}

			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("page", "2")
			divisionsWithPage, _, err := TOSession.GetDivisions(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(divisions[1:2], divisionsWithPage) {
					t.Error("expected GET Divisions with limit = 1, page = 2 to return second result")
				}
			} else {
				t.Error("Error in getting divisions by limit and page")
			}
		} else {
			t.Errorf("only one division found, so offset functionality can't test")
		}
	} else {
		t.Errorf("No division found to check pagination")
	}

	qparams = url.Values{}
	qparams.Set("limit", "-2")
	_, _, err = TOSession.GetDivisions(qparams, nil)
	if err == nil {
		t.Error("expected GET Divisions to return an error when limit is not bigger than -1")
	} else if !strings.Contains(err.Error(), "must be bigger than -1") {
		t.Errorf("expected GET Divisions to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("offset", "0")
	_, _, err = TOSession.GetDivisions(qparams, nil)
	if err == nil {
		t.Error("expected GET Divisions to return an error when offset is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Divisions to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("page", "0")
	_, _, err = TOSession.GetDivisions(qparams, nil)
	if err == nil {
		t.Error("expected GET Divisions to return an error when page is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Divisions to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}

func GetDivisionByInvalidId(t *testing.T) {
	resp, _, err := TOSession.GetDivisionByID(10000, nil)
	if err != nil {
		t.Errorf("Error!! Getting Division by Invalid ID %v", err)
	}
	if len(resp) >= 1 {
		t.Errorf("Error!! Invalid ID shouldn't have any response %v Error %v", resp, err)
	}
}

func GetDivisionByInvalidName(t *testing.T) {
	resp, _, err := TOSession.GetDivisionByName("abcd", nil)
	if err != nil {
		t.Errorf("Getting Division by Invalid Name %v", err)
	}
	if len(resp) >= 1 {
		t.Errorf("Invalid Name shouldn't have any response %v Error %v", resp, err)
	}
}

func DeleteTestDivisionsByInvalidId(t *testing.T) {
	delResp, _, err := TOSession.DeleteDivision(10000)
	if err == nil {
		t.Errorf("cannot DELETE Division by Invalid ID: %v - %v", err, delResp)
	}
}
