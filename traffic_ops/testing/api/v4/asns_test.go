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
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestASN(t *testing.T) {
	WithObjs(t, []TCObj{Types, CacheGroups, ASN}, func() {
		GetTestASNsIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		opts := client.NewRequestOptions()
		opts.Header.Set(rfc.IfModifiedSince, time)
		SortTestASNs(t)
		UpdateTestASNs(t)
		GetTestASNs(t)
		GetTestASNsIMSAfterChange(t, opts)
	})
}

func GetTestASNsIMSAfterChange(t *testing.T, opts client.RequestOptions) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	for _, asn := range testData.ASNs {
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		_, reqInf, err := TOSession.GetASNs(opts)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)
	for _, asn := range testData.ASNs {
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		_, reqInf, err := TOSession.GetASNs(opts)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestASNsIMS(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, asn := range testData.ASNs {
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		futureTime := time.Now().AddDate(0, 0, 1)
		time := futureTime.Format(time.RFC1123)
		opts.Header.Set(rfc.IfModifiedSince, time)
		_, reqInf, err := TOSession.GetASNs(opts)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestASNs(t *testing.T) {
	if len(testData.CacheGroups) < 1 {
		t.Fatal("Need at least one Cache Group to test creating ASNs")
	}
	cg := testData.CacheGroups[0]
	if cg.Name == nil {
		t.Fatal("Cache Group found in the test data with null or undefined name")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *cg.Name)
	resp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("unable to get cachgroup ID: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Cache Group with Name '%s', got: %d", *cg.Name, len(resp.Response))
	}
	if resp.Response[0].ID == nil {
		t.Fatalf("Cache Group '%s' had no ID in Traffic Ops response", *cg.Name)
	}
	id := *resp.Response[0].ID
	for _, asn := range testData.ASNs {
		asn.CachegroupID = id
		resp, _, err := TOSession.CreateASN(asn, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create ASN: %v - alerts: %+v", err, resp)
		}
	}

}

func SortTestASNs(t *testing.T) {
	var sortedList []string
	resp, _, err := TOSession.GetASNs(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if len(resp.Response) < 2 {
		t.Fatal("Cannot test sort order with less than 2 ASNs")
	}
	for _, asn := range resp.Response {
		sortedList = append(sortedList, strconv.Itoa(asn.ASN))
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestASNs(t *testing.T) {
	if len(testData.ASNs) < 1 {
		t.Fatal("Need at least one ASN to test updating ASNs")
	}
	firstASN := testData.ASNs[0]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("asn", strconv.Itoa(firstASN.ASN))

	resp, _, err := TOSession.GetASNs(opts)
	if err != nil {
		t.Fatalf("cannot get ASN by ASN %d: %v - alerts: %+v", firstASN.ASN, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected ASN %d to exist, but Traffic Ops returned no such ASN", firstASN.ASN)
	}

	remoteASN := resp.Response[0]
	remoteASN.ASN = 7777
	alert, _, err := TOSession.UpdateASN(remoteASN.ID, remoteASN, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot update ASN by id: %v - alerts: %+v", err, alert)
	}

	opts.QueryParameters.Del("asn")
	opts.QueryParameters.Set("id", strconv.Itoa(remoteASN.ID))
	resp, _, err = TOSession.GetASNs(opts)
	if err != nil {
		t.Errorf("cannot get ANS by ID %d: %v - alerts: %+v", firstASN.ASN, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected ASN with ID %d to exist after update, but Traffic Ops returned no such ASN", remoteASN.ID)
	}
	respASN := resp.Response[0]
	if respASN.ASN != remoteASN.ASN {
		t.Errorf("results do not match actual: %v, expected: %v", respASN.ASN, remoteASN.ASN)
	}

	//Revert back to original ASN number for further functions to work correctly
	respASN.ASN = firstASN.ASN
	alert, _, err = TOSession.UpdateASN(respASN.ID, respASN, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update ASN by id: %v - alerts: %+v", err, alert)
	}
}

func GetTestASNs(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, asn := range testData.ASNs {
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		resp, _, err := TOSession.GetASNs(opts)
		if err != nil {
			t.Errorf("cannot get ASN by asn: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

func DeleteTestASNs(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, asn := range testData.ASNs {
		opts.QueryParameters.Set("asn", strconv.Itoa(asn.ASN))
		resp, _, err := TOSession.GetASNs(opts)
		if err != nil {
			t.Errorf("cannot get ASN %d: %v - alerts: %+v", asn.ASN, err, resp.Alerts)
			continue
		}
		if len(resp.Response) < 1 {
			t.Errorf("ASN %d existed in the test data, but not in Traffic Ops", asn.ASN)
			continue
		}

		respASN := resp.Response[0]

		alerts, _, err := TOSession.DeleteASN(respASN.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete ASN %d: %v - alerts: %+v", respASN.ASN, err, alerts)
		}

		// Retrieve the ASN to see if it got deleted
		asns, _, err := TOSession.GetASNs(opts)
		if err != nil {
			t.Errorf("error trying to fetch ASN after deletion: %v - alerts: %+v", err, asns.Alerts)
		}
		if len(asns.Response) > 0 {
			t.Errorf("expected ASN %d to be deleted, but it was found in Traffic Ops's response", asn.ASN)
		}
	}
}
