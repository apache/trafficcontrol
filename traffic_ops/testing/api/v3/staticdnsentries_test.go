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
	"testing"
	"time"

	"reflect"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
)

func TestStaticDNSEntries(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, DeliveryServices, StaticDNSEntries}, func() {
		GetTestStaticDNSEntriesIMS(t)
		GetTestStaticDNSEntries(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		UpdateTestStaticDNSEntries(t)
		GetTestStaticDNSEntriesIMSAfterChange(t, header)
		UpdateTestStaticDNSEntriesInvalidAddress(t)
	})
}

func GetTestStaticDNSEntriesIMSAfterChange(t *testing.T, header http.Header) {
	for _, staticDNSEntry := range testData.StaticDNSEntries {
		_, reqInf, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host, header)
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
	for _, staticDNSEntry := range testData.StaticDNSEntries {
		_, reqInf, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestStaticDNSEntriesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0,0,1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)

	for _, staticDNSEntry := range testData.StaticDNSEntries {
		_, reqInf, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestStaticDNSEntries(t *testing.T) {
	for _, staticDNSEntry := range testData.StaticDNSEntries {
		resp, _, err := TOSession.CreateStaticDNSEntry(staticDNSEntry)
		t.Log("Response: ", resp)
		if err != nil {
			t.Errorf("could not CREATE staticDNSEntry: %v", err)
		}
	}

}

func UpdateTestStaticDNSEntries(t *testing.T) {

	firstStaticDNSEntry := testData.StaticDNSEntries[0]
	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	resp, _, err := TOSession.GetStaticDNSEntriesByHost(firstStaticDNSEntry.Host, nil)
	if err != nil {
		t.Errorf("cannot GET StaticDNSEntries by name: '%s', %v", firstStaticDNSEntry.Host, err)
	}
	remoteStaticDNSEntry := resp[0]
	expectedAddress := "192.168.0.2"
	remoteStaticDNSEntry.Address = expectedAddress
	var alert tc.Alerts
	var status int
	alert, _, status, err = TOSession.UpdateStaticDNSEntryByID(remoteStaticDNSEntry.ID, remoteStaticDNSEntry)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot UPDATE StaticDNSEntries using url: %v - %v", err, alert)
	}

	// Retrieve the StaticDNSEntries to check StaticDNSEntries name got updated
	resp, _, err = TOSession.GetStaticDNSEntryByID(remoteStaticDNSEntry.ID, nil)
	if err != nil {
		t.Errorf("cannot GET StaticDNSEntries by name: '$%s', %v", firstStaticDNSEntry.Host, err)
	}
	respStaticDNSEntry := resp[0]
	if respStaticDNSEntry.Address != expectedAddress {
		t.Errorf("results do not match actual: %s, expected: %s", respStaticDNSEntry.Address, expectedAddress)
	}

}

func UpdateTestStaticDNSEntriesInvalidAddress(t *testing.T) {

	expectedAlerts := []tc.Alerts{tc.Alerts{[]tc.Alert{tc.Alert{"'address' must be a valid IPv4 address", "error"}}}, tc.Alerts{[]tc.Alert{tc.Alert{"'address' must be a valid DNS name", "error"}}}, tc.Alerts{[]tc.Alert{tc.Alert{"'address' must be a valid IPv6 address", "error"}}}}

	// A_RECORD
	firstStaticDNSEntry := testData.StaticDNSEntries[0]
	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	resp, _, err := TOSession.GetStaticDNSEntriesByHost(firstStaticDNSEntry.Host, nil)
	if err != nil {
		t.Errorf("cannot GET StaticDNSEntries by name: '%s', %v", firstStaticDNSEntry.Host, err)
	}
	remoteStaticDNSEntry := resp[0]
	expectedAddress := "test.testdomain.net."
	remoteStaticDNSEntry.Address = expectedAddress
	var alert tc.Alerts
	var status int
	alert, _, status, err = TOSession.UpdateStaticDNSEntryByID(remoteStaticDNSEntry.ID, remoteStaticDNSEntry)
	t.Log("Status Code [expect 400]: ", status)
	if err != nil {
		t.Logf("cannot UPDATE StaticDNSEntries using url: %v - %v\n", err, alert)
	}
	if !reflect.DeepEqual(alert, expectedAlerts[0]) {
		t.Errorf("got alerts: %v but expected alerts: %v", alert, expectedAlerts[0])
	}

	// CNAME_RECORD
	secondStaticDNSEntry := testData.StaticDNSEntries[1]
	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	resp, _, err = TOSession.GetStaticDNSEntriesByHost(secondStaticDNSEntry.Host, nil)
	if err != nil {
		t.Errorf("cannot GET StaticDNSEntries by name: '%s', %v", secondStaticDNSEntry.Host, err)
	}
	remoteStaticDNSEntry = resp[0]
	expectedAddress = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	remoteStaticDNSEntry.Address = expectedAddress
	alert, _, status, err = TOSession.UpdateStaticDNSEntryByID(remoteStaticDNSEntry.ID, remoteStaticDNSEntry)
	t.Log("Status Code [expect 400]: ", status)
	if err != nil {
		t.Logf("cannot UPDATE StaticDNSEntries using url: %v - %v\n", err, alert)
	}
	if !reflect.DeepEqual(alert, expectedAlerts[1]) {
		t.Errorf("got alerts: %v but expected alerts: %v", alert, expectedAlerts[1])
	}

	// AAAA_RECORD
	thirdStaticDNSEntry := testData.StaticDNSEntries[2]
	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	resp, _, err = TOSession.GetStaticDNSEntriesByHost(thirdStaticDNSEntry.Host, nil)
	if err != nil {
		t.Errorf("cannot GET StaticDNSEntries by name: '%s', %v", thirdStaticDNSEntry.Host, err)
	}
	remoteStaticDNSEntry = resp[0]
	expectedAddress = "192.168.0.1"
	remoteStaticDNSEntry.Address = expectedAddress
	alert, _, status, err = TOSession.UpdateStaticDNSEntryByID(remoteStaticDNSEntry.ID, remoteStaticDNSEntry)
	t.Log("Status Code [expect 400]: ", status)
	if err != nil {
		t.Logf("cannot UPDATE StaticDNSEntries using url: %v - %v\n", err, alert)
	}
	if !reflect.DeepEqual(alert, expectedAlerts[2]) {
		t.Errorf("got alerts: %v but expected alerts: %v", alert, expectedAlerts[2])
	}
}

func GetTestStaticDNSEntries(t *testing.T) {

	for _, staticDNSEntry := range testData.StaticDNSEntries {
		resp, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host, nil)
		if err != nil {
			t.Errorf("cannot GET StaticDNSEntries by name: %v - %v", err, resp)
		}
	}
}

func DeleteTestStaticDNSEntries(t *testing.T) {

	for _, staticDNSEntry := range testData.StaticDNSEntries {
		// Retrieve the StaticDNSEntries by name so we can get the id for the Update
		resp, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host, nil)
		if err != nil {
			t.Errorf("cannot GET StaticDNSEntries by name: %v - %v", staticDNSEntry.Host, err)
		}
		if len(resp) > 0 {
			respStaticDNSEntry := resp[0]

			_, _, err := TOSession.DeleteStaticDNSEntryByID(respStaticDNSEntry.ID)
			if err != nil {
				t.Errorf("cannot DELETE StaticDNSEntry by name: '%s' %v", respStaticDNSEntry.Host, err)
			}

			// Retrieve the StaticDNSEntry to see if it got deleted
			staticDNSEntries, _, err := TOSession.GetStaticDNSEntriesByHost(staticDNSEntry.Host, nil)
			if err != nil {
				t.Errorf("error deleting StaticDNSEntrie name: %s", err.Error())
			}
			if len(staticDNSEntries) > 0 {
				t.Errorf("expected StaticDNSEntry name: %s to be deleted", staticDNSEntry.Host)
			}
		}
	}
}
