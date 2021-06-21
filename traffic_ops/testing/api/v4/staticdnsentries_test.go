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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
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
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestStaticDNSEntries(t)
		UpdateTestStaticDNSEntries(t)
		UpdateTestStaticDNSEntriesWithHeaders(t, header)
		GetTestStaticDNSEntriesIMSAfterChange(t, header)
		UpdateTestStaticDNSEntriesInvalidAddress(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestStaticDNSEntriesWithHeaders(t, header)
	})
}

func UpdateTestStaticDNSEntriesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.StaticDNSEntries) < 1 {
		t.Error("Need at least one Static DNS Entry to test updating a Static DNS Entry with an HTTP Header")
		return
	}
	firstStaticDNSEntry := testData.StaticDNSEntries[0]

	opts := client.NewRequestOptions()
	opts.Header = header
	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	opts.QueryParameters.Set("host", firstStaticDNSEntry.Host)
	resp, _, err := TOSession.GetStaticDNSEntries(opts)
	if err != nil {
		t.Errorf("cannot get Static DNS Entries filtered by host name '%s': %v - alerts: %+v", firstStaticDNSEntry.Host, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Errorf("Expected at least one Static DNS Entry to exist with host name '%s'", firstStaticDNSEntry.Host)
		return
	}
	remoteStaticDNSEntry := resp.Response[0]
	expectedAddress := "192.168.0.2"
	remoteStaticDNSEntry.Address = expectedAddress

	opts.QueryParameters.Del("host")
	_, reqInf, _ := TOSession.UpdateStaticDNSEntry(remoteStaticDNSEntry.ID, remoteStaticDNSEntry, opts)
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %d", reqInf.StatusCode)
	}
}

func GetTestStaticDNSEntriesIMSAfterChange(t *testing.T, header http.Header) {
	opts := client.NewRequestOptions()
	opts.Header = header
	for _, staticDNSEntry := range testData.StaticDNSEntries {
		opts.QueryParameters.Set("host", staticDNSEntry.Host)
		resp, reqInf, err := TOSession.GetStaticDNSEntries(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
		}
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)

	for _, staticDNSEntry := range testData.StaticDNSEntries {
		opts.QueryParameters.Set("host", staticDNSEntry.Host)
		resp, reqInf, err := TOSession.GetStaticDNSEntries(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestStaticDNSEntriesIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)

	for _, staticDNSEntry := range testData.StaticDNSEntries {
		opts.QueryParameters.Set("host", staticDNSEntry.Host)
		resp, reqInf, err := TOSession.GetStaticDNSEntries(opts)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func CreateTestStaticDNSEntries(t *testing.T) {
	for _, staticDNSEntry := range testData.StaticDNSEntries {
		resp, _, err := TOSession.CreateStaticDNSEntry(staticDNSEntry, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Static DNS Entry: %v - alerts: %+v", err, resp.Alerts)
		}
	}

}

func SortTestStaticDNSEntries(t *testing.T) {
	resp, _, err := TOSession.GetStaticDNSEntries(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	sortedList := make([]string, 0, len(resp.Response))
	for _, sde := range resp.Response {
		sortedList = append(sortedList, sde.Host)
	}

	if !sort.StringsAreSorted(sortedList) {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestStaticDNSEntries(t *testing.T) {
	if len(testData.StaticDNSEntries) < 1 {
		t.Fatal("Need at least one Static DNS Entry to test updating a Static DNS Entry")
	}
	firstStaticDNSEntry := testData.StaticDNSEntries[0]

	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("host", firstStaticDNSEntry.Host)
	resp, _, err := TOSession.GetStaticDNSEntries(opts)
	if err != nil {
		t.Errorf("cannot get Static DNS Entries by host name '%s': %v - alerts: %+v", firstStaticDNSEntry.Host, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one Static DNS Entry to exist with host name '%s'", firstStaticDNSEntry.Host)
	}

	remoteStaticDNSEntry := resp.Response[0]
	expectedAddress := "192.168.0.2"
	remoteStaticDNSEntry.Address = expectedAddress

	alert, _, err := TOSession.UpdateStaticDNSEntry(remoteStaticDNSEntry.ID, remoteStaticDNSEntry, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot updated Static DNS Entry: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the StaticDNSEntries to check StaticDNSEntries name got updated
	opts.QueryParameters.Del("host")
	opts.QueryParameters.Set("id", strconv.Itoa(remoteStaticDNSEntry.ID))
	resp, _, err = TOSession.GetStaticDNSEntries(opts)
	if err != nil {
		t.Errorf("cannot get Static DNS Entries filtered by ID %d: %v - alerts: %+v", remoteStaticDNSEntry.ID, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Static DNS Entry to exist with ID %d, found: %d", remoteStaticDNSEntry.ID, len(resp.Response))
	}
	respStaticDNSEntry := resp.Response[0]
	if respStaticDNSEntry.Address != expectedAddress {
		t.Errorf("results do not match actual: %s, expected: %s", respStaticDNSEntry.Address, expectedAddress)
	}

}

func UpdateTestStaticDNSEntriesInvalidAddress(t *testing.T) {
	if len(testData.StaticDNSEntries) < 3 {
		t.Fatal("Need at least three Static DNS Entries to test updating a Static DNS Entry with an invalid address, DNS name, and CNAME record")
	}

	expectedAlerts := []string{
		"'address' must be a valid IPv4 address",
		"'address' must be a valid DNS name",
		"'address' for type: CNAME_RECORD must have a trailing period",
		"'address' must be a valid IPv6 address",
	}

	// A_RECORD
	firstStaticDNSEntry := testData.StaticDNSEntries[0]

	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("host", firstStaticDNSEntry.Host)
	resp, _, err := TOSession.GetStaticDNSEntries(opts)
	if err != nil {
		t.Errorf("cannot get Static DNS Entries filtered by host name '%s': %v - alerts: %+v", firstStaticDNSEntry.Host, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one Static DNS Entry to exist with host name '%s'", firstStaticDNSEntry.Host)
	}
	remoteStaticDNSEntry := resp.Response[0]
	expectedAddress := "test.testdomain.net."
	remoteStaticDNSEntry.Address = expectedAddress
	alerts, _, err := TOSession.UpdateStaticDNSEntry(remoteStaticDNSEntry.ID, remoteStaticDNSEntry, client.RequestOptions{})
	if err == nil {
		t.Errorf("making invalid update to static DNS entry - expected: error, actual: nil")
	} else if !alertsHaveError(alerts.Alerts, expectedAlerts[0]) {
		t.Errorf("Expected an error-level alert containing '%s', but didn't find it - error: %v - alerts: %+v", expectedAlerts[0], err, alerts.Alerts)
	}

	// CNAME_RECORD
	secondStaticDNSEntry := testData.StaticDNSEntries[1]

	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	opts.QueryParameters.Set("host", secondStaticDNSEntry.Host)
	resp, _, err = TOSession.GetStaticDNSEntries(opts)
	if err != nil {
		t.Errorf("cannot get Static DNS Entries by host name '%s': %v - alerts: %+v", secondStaticDNSEntry.Host, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one Static DNS Entry to exist with host name '%s'", secondStaticDNSEntry.Host)
	}

	remoteStaticDNSEntry = resp.Response[0]
	expectedAddress = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	remoteStaticDNSEntry.Address = expectedAddress

	alerts, _, err = TOSession.UpdateStaticDNSEntry(remoteStaticDNSEntry.ID, remoteStaticDNSEntry, client.RequestOptions{})
	if err == nil {
		t.Errorf("making invalid update to static DNS entry - expected: error, actual: nil")
	} else if !alertsHaveError(alerts.Alerts, expectedAlerts[1]) {
		t.Errorf("Expected an error-level alert containing '%s', but didn't find it - error: %v - alerts: %+v", expectedAlerts[1], err, alerts.Alerts)
	}

	//CNAME_RECORD: missing a trailing period
	expectedAddressMissingPeriod := "cdn.test.com"
	remoteStaticDNSEntry.Address = expectedAddressMissingPeriod
	alerts, _, err = TOSession.UpdateStaticDNSEntry(remoteStaticDNSEntry.ID, remoteStaticDNSEntry, client.RequestOptions{})
	if err == nil {
		t.Errorf("making invalid update to static DNS entry - expected: error, actual: nil")
	} else if !alertsHaveError(alerts.Alerts, expectedAlerts[2]) {
		t.Errorf("Expected an error-level alert containing '%s', but didn't find it - error: %v - alerts: %+v", expectedAlerts[2], err, alerts.Alerts)
	}

	// AAAA_RECORD
	thirdStaticDNSEntry := testData.StaticDNSEntries[2]

	// Retrieve the StaticDNSEntries by name so we can get the id for the Update
	opts.QueryParameters.Set("host", thirdStaticDNSEntry.Host)
	resp, _, err = TOSession.GetStaticDNSEntries(opts)
	if err != nil {
		t.Errorf("cannot get Static DNS Entries filtered by host name '%s': %v - alerts: %+v", thirdStaticDNSEntry.Host, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected at least one Static DNS Entry to exist with host name '%s'", thirdStaticDNSEntry.Host)
	}

	remoteStaticDNSEntry = resp.Response[0]
	expectedAddress = "192.168.0.1"
	remoteStaticDNSEntry.Address = expectedAddress
	alerts, _, err = TOSession.UpdateStaticDNSEntry(remoteStaticDNSEntry.ID, remoteStaticDNSEntry, client.RequestOptions{})
	if err == nil {
		t.Errorf("making invalid update to static DNS entry - expected: error, actual: nil")
	} else if !alertsHaveError(alerts.Alerts, expectedAlerts[3]) {
		t.Errorf("Expected an error-level alert containing '%s', but didn't find it - error: %v - alerts: %+v", expectedAlerts[3], err, alerts.Alerts)
	}
}

func GetTestStaticDNSEntries(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, staticDNSEntry := range testData.StaticDNSEntries {
		opts.QueryParameters.Set("host", staticDNSEntry.Host)
		resp, _, err := TOSession.GetStaticDNSEntries(opts)
		if err != nil {
			t.Errorf("cannot get Static DNS Entries filtered by host name: %v - alerts: %+v", err, resp.Alerts)
		}
	}
}

// This test will break if any two Static DNS Entries share a host name (not sure if that's legal)
func DeleteTestStaticDNSEntries(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, staticDNSEntry := range testData.StaticDNSEntries {
		// Retrieve the StaticDNSEntries by name so we can get the id for the Update
		opts.QueryParameters.Set("host", staticDNSEntry.Host)
		resp, _, err := TOSession.GetStaticDNSEntries(opts)
		if err != nil {
			t.Errorf("cannot get Static DNS Entries filtered by host name '%s': %v - alerts: %+v", staticDNSEntry.Host, err, resp.Alerts)
		}
		if len(resp.Response) > 0 {
			respStaticDNSEntry := resp.Response[0]

			alerts, _, err := TOSession.DeleteStaticDNSEntry(respStaticDNSEntry.ID, client.RequestOptions{})
			if err != nil {
				t.Errorf("cannot delete Static DNS Entry for host name '%s': %v - alerts: %+v", respStaticDNSEntry.Host, err, alerts.Alerts)
			}

			// Retrieve the StaticDNSEntry to see if it got deleted
			staticDNSEntries, _, err := TOSession.GetStaticDNSEntries(opts)
			if err != nil {
				t.Errorf("error fetching Static DNS Entry after supposed deletion: %v - alerts: %+v", err, staticDNSEntries.Alerts)
			}
			if len(staticDNSEntries.Response) > 0 {
				t.Errorf("expected Static DNS Entry with host name '%s' to be deleted, but it was found in Traffic Ops", staticDNSEntry.Host)
			}
		}
	}
}
