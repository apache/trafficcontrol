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
	tc "github.com/apache/trafficcontrol/lib/go-tc"
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
	if len(testData.Coordinates) > 0 {
		firstCoord := testData.Coordinates[0]
		resp, _, err := TOSession.GetCoordinateByName(firstCoord.Name, header)
		if err != nil {
			t.Errorf("cannot GET Coordinate by name: %v - %v", firstCoord.Name, err)
		}
		if len(resp) > 0 {
			coord := resp[0]
			expectedLat := 12.34
			coord.Latitude = expectedLat

			_, reqInf, err := TOSession.UpdateCoordinate(coord.ID, coord, header)
			if err == nil {
				t.Errorf("Expected error about precondition failed, but got none")
			}
			if reqInf.StatusCode != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		} else {
			t.Errorf("No coordinates available")
		}
	} else {
		t.Errorf("No Coordinates available to update")
	}
}

func GetTestCoordinatesIMSAfterChange(t *testing.T, header http.Header) {
	for _, coord := range testData.Coordinates {
		_, reqInf, err := TOSession.GetCoordinateByName(coord.Name, header)
		if err != nil {
			t.Fatalf("could not GET coordinates: %v", err)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, coord := range testData.Coordinates {
		_, reqInf, err := TOSession.GetCoordinateByName(coord.Name, header)
		if err != nil {
			t.Fatalf("could not GET coordinates: %v", err)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCoordinatesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, coord := range testData.Coordinates {
		_, reqInf, err := TOSession.GetCoordinateByName(coord.Name, header)
		if err != nil {
			t.Fatalf("No error expected, but got: %v", err)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestCoordinates(t *testing.T) {
	for _, coord := range testData.Coordinates {
		_, _, err := TOSession.CreateCoordinate(coord)
		if err != nil {
			t.Errorf("could not CREATE coordinates: %v", err)
		}
	}
}

func GetTestCoordinates(t *testing.T) {
	for _, coord := range testData.Coordinates {
		resp, _, err := TOSession.GetCoordinateByName(coord.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Coordinate: %v - %v", err, resp)
		}
	}
}

func SortTestCoordinates(t *testing.T) {
	var sortedList []string
	resp, _, err := TOSession.GetCoordinates(nil, nil)
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

func SortTestCoordinatesDesc(t *testing.T) {

	respAsc, _, err1 := TOSession.GetCoordinates(nil, nil)
	params := url.Values{}
	params.Set("sortOrder", "desc")
	respDesc, _, err2 := TOSession.GetCoordinates(params, nil)

	if err1 != nil {
		t.Errorf("Expected no error, but got error in Coordinates Ascending %v", err1)
	}
	if err2 != nil {
		t.Errorf("Expected no error, but got error in Coordinates Descending %v", err2)
	}

	if len(respAsc) == len(respDesc) {
		if len(respAsc) > 0 && len(respDesc) > 0 {
			// reverse the descending-sorted response and compare it to the ascending-sorted one
			for start, end := 0, len(respDesc)-1; start < end; start, end = start+1, end-1 {
				respDesc[start], respDesc[end] = respDesc[end], respDesc[start]
			}
			if respDesc[0].Name != "" && respAsc[0].Name != "" {
				if !reflect.DeepEqual(respDesc[0].Name, respAsc[0].Name) {
					t.Errorf("Coordinates responses are not equal after reversal: %s - %s", respDesc[0].Name, respAsc[0].Name)
				}
			}
		} else {
			t.Errorf("No Response returned from GET Coordinates using SortOrder")
		}
	} else {
		t.Fatalf("Coordinates response length are not equal Asc: %d Desc: %d", len(respAsc), len(respDesc))
	}
}

func UpdateTestCoordinates(t *testing.T) {
	if len(testData.Coordinates) > 0 {
		firstCoord := testData.Coordinates[0]
		resp, _, err := TOSession.GetCoordinateByName(firstCoord.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Coordinate by name: %v - %v", firstCoord.Name, err)
		}
		if len(resp) > 0 {
			coord := resp[0]
			expectedLat := 12.34
			coord.Latitude = expectedLat

			var alert tc.Alerts
			alert, _, err = TOSession.UpdateCoordinate(coord.ID, coord, nil)
			if err != nil {
				t.Errorf("cannot UPDATE Coordinate by id: %v - %v", err, alert)
			}

			// Retrieve the Coordinate to check Coordinate name got updated
			resp, _, err = TOSession.GetCoordinateByID(coord.ID, nil)
			if err != nil {
				t.Errorf("cannot GET Coordinate by name: '$%s', %v", firstCoord.Name, err)
			}
			if len(resp) > 0 {
				coord = resp[0]
				if coord.Latitude != expectedLat {
					t.Errorf("results do not match actual: %s, expected: %f", coord.Name, expectedLat)
				}
			} else {
				t.Errorf("Can't retrieve coordinates to check the updated value")
			}
		} else {
			t.Errorf("No Coordinates available to update")
		}
	} else {
		t.Errorf("No Coordinates available to update")
	}
}

func DeleteTestCoordinates(t *testing.T) {
	for _, coord := range testData.Coordinates {
		// Retrieve the Coordinate by name so we can get the id for the Update
		resp, _, err := TOSession.GetCoordinateByName(coord.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Coordinate by name: %v - %v", coord.Name, err)
		}
		if len(resp) > 0 {
			respCoord := resp[0]
			_, _, err := TOSession.DeleteCoordinate(respCoord.ID)
			if err != nil {
				t.Errorf("cannot DELETE Coordinate by name: '%s' %v", respCoord.Name, err)
			}
			// Retrieve the Coordinate to see if it got deleted
			coords, _, err := TOSession.GetCoordinateByName(coord.Name, nil)
			if err != nil {
				t.Errorf("error deleting Coordinate name: %s", err.Error())
			}
			if len(coords) > 0 {
				t.Errorf("expected Coordinate name: %s to be deleted", coord.Name)
			}
		} else {
			t.Errorf("No Coordinates available to delete")
		}
	}
}

func GetTestCoordinatesByInvalidId(t *testing.T) {
	coordinatesResp, _, err := TOSession.GetCoordinateByID(10000, nil)
	if err != nil {
		t.Errorf("Error!! Getting Coordinates by Invalid ID %v", err)
	}
	if len(coordinatesResp) >= 1 {
		t.Errorf("Error!! Invalid ID shouldn't have any response %v Error %v", coordinatesResp, err)
	}
}

func GetTestCoordiantesByInvalidName(t *testing.T) {
	coordinatesResp, _, err := TOSession.GetCoordinateByName("abcd", nil)
	if err != nil {
		t.Errorf("Error!! Getting Coordinates by Invalid Name %v", err)
	}
	if len(coordinatesResp) >= 1 {
		t.Errorf("Error!! Invalid Name shouldn't have any response %v Error %v", coordinatesResp, err)
	}
}

func GetTestPaginationSupportCoordinates(t *testing.T) {

	qparams := url.Values{}
	qparams.Set("orderby", "id")
	coordinates, _, err := TOSession.GetCoordinates(qparams, nil)
	if err != nil {
		t.Fatalf("cannot GET Coordinates: %v", err)
	}

	if len(coordinates) > 0 {
		qparams = url.Values{}
		qparams.Set("orderby", "id")
		qparams.Set("limit", "1")
		coordinatesWithLimit, _, err := TOSession.GetCoordinates(qparams, nil)
		if err == nil {
			if !reflect.DeepEqual(coordinates[:1], coordinatesWithLimit) {
				t.Error("expected GET Coordinates with limit = 1 to return first result")
			}
		} else {
			t.Error("Error in getting coordinates by limit")
		}

		if len(coordinates) > 1 {
			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("offset", "1")
			coordinatesWithOffset, _, err := TOSession.GetCoordinates(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(coordinates[1:2], coordinatesWithOffset) {
					t.Error("expected GET Coordinates with limit = 1, offset = 1 to return second result")
				}
			} else {
				t.Error("Error in getting coordinates by limit and offset")
			}

			qparams = url.Values{}
			qparams.Set("orderby", "id")
			qparams.Set("limit", "1")
			qparams.Set("page", "2")
			coordinatesWithPage, _, err := TOSession.GetCoordinates(qparams, nil)
			if err == nil {
				if !reflect.DeepEqual(coordinates[1:2], coordinatesWithPage) {
					t.Error("expected GET Coordinates with limit = 1, page = 2 to return second result")
				}
			} else {
				t.Error("Error in getting coordinates by limit and page")
			}
		} else {
			t.Errorf("only one Coordinates found, so offset functionality can't test")
		}
	} else {
		t.Errorf("No Coordinates found to check pagination")
	}

	qparams = url.Values{}
	qparams.Set("limit", "-2")
	_, _, err = TOSession.GetCoordinates(qparams, nil)
	if err == nil {
		t.Error("expected GET Coordinates to return an error when limit is not bigger than -1")
	} else if !strings.Contains(err.Error(), "must be bigger than -1") {
		t.Errorf("expected GET Coordinates to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("offset", "0")
	_, _, err = TOSession.GetCoordinates(qparams, nil)
	if err == nil {
		t.Error("expected GET Coordinates to return an error when offset is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Coordinates to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("page", "0")
	_, _, err = TOSession.GetCoordinates(qparams, nil)
	if err == nil {
		t.Error("expected GET Coordinates to return an error when page is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET Coordinates to return an error for page is not a positive integer, actual error: %v" + err.Error())
	}
}

func CreateTestCoordinatesWithInvalidName(t *testing.T) {
	if len(testData.Coordinates) > 0 {
		firstCoordinates := testData.Coordinates[0]
		firstCoordinates.Name = ""
		_, reqInf, err := TOSession.CreateCoordinate(firstCoordinates)
		if reqInf.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
		}
		if err == nil {
			t.Errorf("Getting Coordinates by Invalid Name")
		}
	} else {
		t.Errorf("No Coordinates available to fetch")
	}
}

func CreateTestCoordinatesWithInvalidLatitude(t *testing.T) {
	if len(testData.Coordinates) > 0 {
		firstCoordinates := testData.Coordinates[0]
		firstCoordinates.Latitude = 20000
		_, reqInf, err := TOSession.CreateCoordinate(firstCoordinates)
		if reqInf.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
		}
		if err == nil {
			t.Errorf("Getting Coordinates by Invalid Latitude")
		}
	} else {
		t.Errorf("No Coordinates available to fetch")
	}
}

func CreateTestCoordinatesWithInvalidLogitude(t *testing.T) {
	if len(testData.Coordinates) > 0 {
		firstCoordinates := testData.Coordinates[0]
		firstCoordinates.Longitude = 20000
		_, reqInf, err := TOSession.CreateCoordinate(firstCoordinates)
		if reqInf.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected 400 status code, got %v", reqInf.StatusCode)
		}
		if err == nil {
			t.Errorf("Getting Coordinates by Invalid Longitude ")
		}
	} else {
		t.Errorf("No Coordinates available to fetch")
	}
}

func UpdateTestCoordinatesByInvalidId(t *testing.T) {
	if len(testData.Coordinates) > 0 {
		firstCoord := testData.Coordinates[0]
		resp, reqInf, err := TOSession.GetCoordinateByName(firstCoord.Name, nil)
		if err != nil {
			t.Errorf("cannot GET Coordinate by name: %v - %v", firstCoord.Name, err)
		}
		if len(resp) > 0 {
			coord := resp[0]
			expectedLat := 12.34
			coord.Latitude = expectedLat

			var alert tc.Alerts
			alert, reqInf, err = TOSession.UpdateCoordinate(10000, coord, nil)
			if err == nil {
				t.Errorf("Updating Coordinate by invalid id: %v", alert)
			}
			if reqInf.StatusCode != http.StatusNotFound {
				t.Fatalf("Expected 404 status code, got %v", reqInf.StatusCode)
			}
		} else {
			t.Errorf("No coordinates available to update")
		}
	} else {
		t.Errorf("No Coordinates available to update")
	}
}

func DeleteTestCoordinatesByInvalidId(t *testing.T) {
	_, reqInf, err := TOSession.DeleteCoordinate(12345)
	if err == nil {
		t.Errorf("Deleting coordinate by invalid Id")
	}
	if reqInf.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 status code, got %v", reqInf.StatusCode)
	}
}
