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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"net/http"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
)

func TestASN(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, ASN}, func() {
		GetTestASNsIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		SortTestASNs(t)
		UpdateTestASNs(t)
		GetTestASNs(t)
		GetTestASNsIMSAfterChange(t, header)
	})
}

func GetTestASNsIMSAfterChange(t *testing.T, header http.Header) {
	for _, asn := range testData.ASNs {
		_, reqInf, err := TOSession.GetASNByASN(asn.ASN)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, asn := range testData.ASNs {
		_, reqInf, err := TOSession.GetASNByASN(asn.ASN)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestASNsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	for _, asn := range testData.ASNs {
		futureTime := time.Now().AddDate(0, 0, 1)
		time := futureTime.Format(time.RFC1123)
		header.Set(rfc.IfModifiedSince, time)
		_, reqInf, err := TOSession.GetASNByID(asn.ID)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestASNs(t *testing.T) {

	for _, asn := range testData.ASNs {
		resp, _, err := TOSession.CreateASN(asn)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE ASNs: %v", err)
		}
	}

}

func SortTestASNs(t *testing.T) {
	var sortedList []string
	resp, _, err := TOSession.GetASNs()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	for i, _ := range resp {
		sortedList = append(sortedList, strconv.Itoa(resp[i].ASN))
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestASNs(t *testing.T) {

	firstASN := testData.ASNs[0]
	// Retrieve the ASN by name so we can get the id for the Update
	resp, _, err := TOSession.GetASNByASN(firstASN.ASN)
	if err != nil {
		t.Errorf("cannot GET ASN by name: '%v', %v", firstASN.ASN, err)
	}
	remoteASN := resp[0]
	remoteASN.ASN = 7777
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateASNByID(resp[0].ID, remoteASN)
	if err != nil {
		t.Errorf("cannot UPDATE ASN by id: %v - %v", err, alert)
	}

	// Retrieve the ASN to check ASN name got updated
	resp, _, err = TOSession.GetASNByID(remoteASN.ID)
	if err != nil {
		t.Errorf("cannot GET ANS by number: '$%v', %v", firstASN.ASN, err)
	}
	respASN := resp[0]
	if respASN.ASN != remoteASN.ASN {
		t.Errorf("results do not match actual: %v, expected: %v", respASN.ASN, remoteASN.ASN)
	}

	//Revert back to original ASN number for further functions to work correctly
	alert, _, err = TOSession.UpdateASNByID(resp[0].ID, firstASN)
	if err != nil {
		t.Errorf("cannot UPDATE ASN by id: %v - %v", err, alert)
	}
}

func GetTestASNs(t *testing.T) {

	for _, asn := range testData.ASNs {
		resp, _, err := TOSession.GetASNByASN(asn.ASN)
		if err != nil {
			t.Errorf("cannot GET ASN by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestASNs(t *testing.T) {

	for _, asn := range testData.ASNs {
		// Retrieve the ASN by name so we can get the id for the Update
		resp, _, err := TOSession.GetASNByASN(asn.ASN)
		if err != nil {
			t.Errorf("cannot GET ASN by number: %v - %v", asn.ASN, err)
		}
		if len(resp) > 0 {
			respASN := resp[0]

			_, _, err := TOSession.DeleteASNByASN(respASN.ID)
			if err != nil {
				t.Errorf("cannot DELETE ASN by ASN number: '%v' %v", respASN.ASN, err)
			}

			// Retrieve the ASN to see if it got deleted
			asns, _, err := TOSession.GetASNByASN(asn.ASN)
			if err != nil {
				t.Errorf("error deleting ASN number: %s", err.Error())
			}
			if len(asns) > 0 {
				t.Errorf("expected ASN number: %v to be deleted", asn.ASN)
			}
		}
	}
}
