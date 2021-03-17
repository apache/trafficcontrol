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
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestServerCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{ServerCapabilities}, func() {
		GetTestServerCapabilities(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		rfcTime := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, rfcTime)
		header.Set(rfc.IfUnmodifiedSince, rfcTime)
		SortTestServerCapabilities(t)
		UpdateTestServerCapabilities(t)
		UpdateTestServerCapabilitiesWithHeaders(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestServerCapabilitiesWithHeaders(t, header)
		ValidationTestServerCapabilities(t)
	})
}

func CreateTestServerCapabilities(t *testing.T) {

	for _, sc := range testData.ServerCapabilities {
		resp, _, err := TOSession.CreateServerCapability(sc)
		if err != nil {
			t.Errorf("could not CREATE server capability: %v", err)
		}
		t.Log("Response: ", resp)
	}

}

func SortTestServerCapabilities(t *testing.T) {
	var header http.Header
	var sortedList []string
	resp, _, err := TOSession.GetServerCapabilitiesWithHdr(header)
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

func GetTestServerCapabilities(t *testing.T) {

	for _, sc := range testData.ServerCapabilities {
		resp, _, err := TOSession.GetServerCapability(sc.Name)
		if err != nil {
			t.Errorf("cannot GET server capability: %v - %v", err, resp)
		} else if resp == nil {
			t.Error("GET server capability expected non-nil response")
		}
	}

	resp, _, err := TOSession.GetServerCapabilities()
	if err != nil {
		t.Errorf("cannot GET server capabilities: %v", err)
	}
	if len(resp) != len(testData.ServerCapabilities) {
		t.Errorf("expected to GET %d server capabilities, actual: %d", len(testData.ServerCapabilities), len(resp))
	}
}

func UpdateTestServerCapabilitiesWithHeaders(t *testing.T, header http.Header) {
	resp, _, err := TOSession.GetServerCapabilitiesWithHdr(header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if len(resp) == 0 {
		t.Fatal("no server capability in response, quitting")
	}
	originalName := resp[0].Name
	newSCName := "sc-test"
	resp[0].Name = newSCName

	_, reqInf, err := TOSession.UpdateServerCapability(originalName, &resp[0], header)
	if err == nil {
		t.Errorf("Expected error about Precondition Failed, got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func ValidationTestServerCapabilities(t *testing.T) {
	_, _, err := TOSession.CreateServerCapability(tc.ServerCapability{Name: "b@dname"})
	if err == nil {
		t.Error("expected POST with invalid name to return an error, actual: nil")
	}
}

func UpdateTestServerCapabilities(t *testing.T) {
	var header http.Header

	// Get server capability name and edit it to a new name
	resp, _, err := TOSession.GetServerCapabilitiesWithHdr(header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if len(resp) == 0 {
		t.Fatalf("no server capability in response, quitting")
	}
	origName := resp[0].Name
	newSCName := "sc-test"
	resp[0].Name = newSCName

	// Update server capability with new name
	updateResponse, _, err := TOSession.UpdateServerCapability(origName, &resp[0], nil)
	if err != nil {
		t.Errorf("cannot PUT server capability: %v - %v", err, updateResponse)
	}

	// Get updated name
	getResp, _, err := TOSession.GetServerCapabilityWithHdr(newSCName, header)
	if err != nil {
		t.Fatalf("Expected no error, but %v", err.Error())
	}
	if getResp == nil {
		t.Fatalf("no server capability in response, quitting")
	}
	if getResp.Name != newSCName {
		t.Errorf("failed to update server capability name, expected: %v but got: %v", newSCName, updateResponse.Name)
	}

	// Set everything back as it was for further testing.
	resp[0].Name = origName
	r, _, err := TOSession.UpdateServerCapability(newSCName, &resp[0], nil)
	if err != nil {
		t.Errorf("cannot PUT seerver capability: %v - %v", err, r)
	}
}

func DeleteTestServerCapabilities(t *testing.T) {

	for _, sc := range testData.ServerCapabilities {
		delResp, _, err := TOSession.DeleteServerCapability(sc.Name)
		if err != nil {
			t.Errorf("cannot DELETE server capability: %v - %v", err, delResp)
		}

		serverCapability, _, err := TOSession.GetServerCapability(sc.Name)
		if err == nil {
			t.Errorf("expected error trying to GET deleted server capability: %s, actual: nil", sc.Name)
		}
		if serverCapability != nil {
			t.Errorf("expected nil trying to GET deleted server capability: %s, actual: non-nil", sc.Name)
		}
	}
}
