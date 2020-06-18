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
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"net/http"
	"sort"
	"testing"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestPhysLocations(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Parameters, Divisions, Regions, PhysLocations}, func() {
		GetTestPhysLocationsIMS(t)
		GetDefaultSortPhysLocationsTest(t)
		GetSortPhysLocationsTest(t)
		currentTime := time.Now().Add(-1 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		UpdateTestPhysLocations(t)
		GetTestPhysLocations(t)
		GetTestPhysLocationsIMSAfterChange(t, header)
	})
}

func GetTestPhysLocationsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0,0,1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, cdn := range testData.PhysLocations {
		_, reqInf, err := TOSession.GetPhysLocationByName(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}

}

func GetTestPhysLocationsIMSAfterChange(t *testing.T, header http.Header) {
	for _, cdn := range testData.PhysLocations {
		_, reqInf, err := TOSession.GetPhysLocationByName(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, cdn := range testData.PhysLocations {
		_, reqInf, err := TOSession.GetPhysLocationByName(cdn.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestPhysLocations(t *testing.T) {
	for _, pl := range testData.PhysLocations {
		resp, _, err := TOSession.CreatePhysLocation(pl)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE physlocations: %v", err)
		}
	}

}

func UpdateTestPhysLocations(t *testing.T) {

	firstPhysLocation := testData.PhysLocations[0]
	// Retrieve the PhysLocation by name so we can get the id for the Update
	resp, _, err := TOSession.GetPhysLocationByName(firstPhysLocation.Name, nil)
	if err != nil {
		t.Errorf("cannot GET PhysLocation by name: '%s', %v", firstPhysLocation.Name, err)
	}
	remotePhysLocation := resp[0]
	expectedPhysLocationCity := "city1"
	remotePhysLocation.City = expectedPhysLocationCity
	var alert tc.Alerts
	alert, _, err = TOSession.UpdatePhysLocationByID(remotePhysLocation.ID, remotePhysLocation)
	if err != nil {
		t.Errorf("cannot UPDATE PhysLocation by id: %v - %v", err, alert)
	}

	// Retrieve the PhysLocation to check PhysLocation name got updated
	resp, _, err = TOSession.GetPhysLocationByID(remotePhysLocation.ID, nil)
	if err != nil {
		t.Errorf("cannot GET PhysLocation by name: '$%s', %v", firstPhysLocation.Name, err)
	}
	respPhysLocation := resp[0]
	if respPhysLocation.City != expectedPhysLocationCity {
		t.Errorf("results do not match actual: %s, expected: %s", respPhysLocation.City, expectedPhysLocationCity)
	}

}

func GetTestPhysLocations(t *testing.T) {

	for _, cdn := range testData.PhysLocations {
		resp, _, err := TOSession.GetPhysLocationByName(cdn.Name, nil)
		if err != nil {
			t.Errorf("cannot GET PhysLocation by name: %v - %v", err, resp)
		}
	}

}

func GetSortPhysLocationsTest(t *testing.T) {
	resp, _, err := TOSession.GetPhysLocations(map[string]string{"orderby": "id"}, nil)
	if err != nil {
		t.Error(err.Error())
	}
	sorted := sort.SliceIsSorted(resp, func(i, j int) bool {
		return resp[i].ID < resp[j].ID
	})
	if !sorted {
		t.Error("expected response to be sorted by id")
	}
}

func GetDefaultSortPhysLocationsTest(t *testing.T) {
	resp, _, err := TOSession.GetPhysLocations(nil, nil)
	if err != nil {
		t.Error(err.Error())
	}
	sorted := sort.SliceIsSorted(resp, func(i, j int) bool {
		return resp[i].Name < resp[j].Name
	})
	if !sorted {
		t.Error("expected response to be sorted by name")
	}
}

func DeleteTestPhysLocations(t *testing.T) {

	for _, cdn := range testData.PhysLocations {
		// Retrieve the PhysLocation by name so we can get the id for the Update
		resp, _, err := TOSession.GetPhysLocationByName(cdn.Name, nil)
		if err != nil {
			t.Errorf("cannot GET PhysLocation by name: %v - %v", cdn.Name, err)
		}
		if len(resp) > 0 {
			respPhysLocation := resp[0]

			_, _, err := TOSession.DeletePhysLocationByID(respPhysLocation.ID)
			if err != nil {
				t.Errorf("cannot DELETE PhysLocation by name: '%s' %v", respPhysLocation.Name, err)
			}

			// Retrieve the PhysLocation to see if it got deleted
			cdns, _, err := TOSession.GetPhysLocationByName(cdn.Name, nil)
			if err != nil {
				t.Errorf("error deleting PhysLocation name: %s", err.Error())
			}
			if len(cdns) > 0 {
				t.Errorf("expected PhysLocation name: %s to be deleted", cdn.Name)
			}
		}
	}
}
