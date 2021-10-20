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
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-rfc"
	"github.com/apache/trafficcontrol/v6/lib/go-tc"
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
	params := url.Values{}
	for _, asn := range testData.ASNs {
		params.Add("asn", strconv.Itoa(asn.ASN))
		_, reqInf, err := TOSession.GetASNsWithHeader(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
		params.Del("asn")
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	for _, asn := range testData.ASNs {
		params.Add("asn", strconv.Itoa(asn.ASN))
		_, reqInf, err := TOSession.GetASNsWithHeader(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
		params.Del("asn")
	}
}

func GetTestASNsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	params := url.Values{}
	for _, asn := range testData.ASNs {
		params.Add("asn", strconv.Itoa(asn.ASN))
		futureTime := time.Now().AddDate(0, 0, 1)
		time := futureTime.Format(time.RFC1123)
		header.Set(rfc.IfModifiedSince, time)
		_, reqInf, err := TOSession.GetASNsWithHeader(&params, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
		params.Del("asn")
	}
}

func CreateTestASNs(t *testing.T) {
	var header http.Header
	resp, _, err := TOSession.GetCacheGroupNullableByNameWithHdr(*testData.CacheGroups[0].Name, header)
	if err != nil {
		t.Fatalf("unable to get cachgroup ID: %v", err)
	}
	for _, asn := range testData.ASNs {
		asn.CachegroupID = *resp[0].ID
		resp, _, err := TOSession.CreateASN(asn)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE ASNs: %v", err)
		}
	}

}

func SortTestASNs(t *testing.T) {
	var header http.Header
	var sortedList []string
	params := url.Values{}
	resp, _, err := TOSession.GetASNsWithHeader(&params, header)
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
	var header http.Header
	firstASN := testData.ASNs[0]
	params := url.Values{}
	params.Add("asn", strconv.Itoa(firstASN.ASN))
	// Retrieve the ASN by name so we can get the id for the Update
	resp, _, err := TOSession.GetASNsWithHeader(&params, header)
	if err != nil {
		t.Errorf("cannot GET ASN by name: '%v', %v", firstASN.ASN, err)
	}
	remoteASN := resp[0]
	remoteASN.ASN = 7777
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateASNByID(remoteASN.ID, remoteASN)
	if err != nil {
		t.Errorf("cannot UPDATE ASN by id: %v - %v", err, alert)
	}

	// Retrieve the ASN to check ASN name got updated
	params.Del("asn")
	params.Add("id", strconv.Itoa(remoteASN.ID))
	resp, _, err = TOSession.GetASNsWithHeader(&params, header)
	if err != nil {
		t.Errorf("cannot GET ANS by number: '$%v', %v", firstASN.ASN, err)
	}
	respASN := resp[0]
	if respASN.ASN != remoteASN.ASN {
		t.Errorf("results do not match actual: %v, expected: %v", respASN.ASN, remoteASN.ASN)
	}

	//Revert back to original ASN number for further functions to work correctly
	respASN.ASN = firstASN.ASN
	alert, _, err = TOSession.UpdateASNByID(respASN.ID, respASN)
	if err != nil {
		t.Errorf("cannot UPDATE ASN by id: %v - %v", err, alert)
	}
}

func GetTestASNs(t *testing.T) {

	var header http.Header
	params := url.Values{}
	for _, asn := range testData.ASNs {
		params.Add("asn", strconv.Itoa(asn.ASN))
		resp, _, err := TOSession.GetASNsWithHeader(&params, header)
		if err != nil {
			t.Errorf("cannot GET ASN by name: %v - %v", err, resp)
		}
		params.Del("asn")
	}
}

func DeleteTestASNs(t *testing.T) {

	var header http.Header
	params := url.Values{}
	for _, asn := range testData.ASNs {
		params.Add("asn", strconv.Itoa(asn.ASN))
		// Retrieve the ASN by name so we can get the id for the Update
		resp, _, err := TOSession.GetASNsWithHeader(&params, header)
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
			asns, _, err := TOSession.GetASNsWithHeader(&params, header)
			if err != nil {
				t.Errorf("error deleting ASN number: %s", err.Error())
			}
			if len(asns) > 0 {
				t.Errorf("expected ASN number: %v to be deleted", asn.ASN)
			}
		}
		params.Del("asn")
	}
}
