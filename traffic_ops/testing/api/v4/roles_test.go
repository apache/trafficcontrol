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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

const (
	roleGood         = 0
	roleInvalidCap   = 1
	roleNeedCap      = 2
	roleBadPrivLevel = 3
)

func TestRoles(t *testing.T) {
	WithObjs(t, []TCObj{Roles}, func() {
		GetTestRolesIMS(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		SortTestRoles(t)
		UpdateTestRoles(t)
		GetTestRoles(t)
		UpdateTestRolesWithHeaders(t, header)
		GetTestRolesIMSAfterChange(t, header)
		VerifyGetRolesOrder(t)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestRolesWithHeaders(t, header)
	})
}

func UpdateTestRolesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Roles) < 1 {
		t.Fatal("Need at least one Role to test updating a Role with HTTP headers")
	}
	firstRole := testData.Roles[0]
	if firstRole.Name == nil {
		t.Fatal("Found a Role in the testing data with null or undefined name")
	}

	// Retrieve the Role by role so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.Header = header
	opts.QueryParameters.Set("name", *firstRole.Name)
	resp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role '%s' by name: %v - alerts: %+v", *firstRole.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Role to exist with name '%s', found: %d", *firstRole.Name, len(resp.Response))
	}
	remoteRole := resp.Response[0]
	if remoteRole.ID == nil {
		t.Fatal("Traffic Ops returned a representation for a Role with null or undefined ID")
	}

	expectedRole := "new admin2"
	remoteRole.Name = &expectedRole
	opts.QueryParameters.Del("name")
	_, reqInf, _ := TOSession.UpdateRole(*remoteRole.ID, remoteRole, opts)
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestRolesIMSAfterChange(t *testing.T, header http.Header) {
	if len(testData.Roles) < roleGood+1 {
		t.Fatalf("Need at least %d Roles to test getting Roles with IMS change", roleGood+1)
	}
	role := testData.Roles[roleGood]
	if role.Name == nil {
		t.Fatal("Found a Role in the testing data with null or undefined name")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *role.Name)
	opts.Header = header
	resp, reqInf, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)

	resp, reqInf, err = TOSession.GetRoles(opts)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestRolesIMS(t *testing.T) {
	if len(testData.Roles) < roleGood+1 {
		t.Fatalf("Need at least %d Roles to test getting Roles with IMS change", roleGood+1)
	}
	role := testData.Roles[roleGood]
	if role.Name == nil {
		t.Fatal("Found a Role in the testing data with null or undefined name")
	}

	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)

	resp, reqInf, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}

}

// This will break if anyone adds a Role or rearranges Roles in the testing data.
func CreateTestRoles(t *testing.T) {
	if len(testData.Roles) > 3 {
		t.Fatal("Too many Roles in the test data. Tests can only handle 3")
	}

	expectedAlerts := []string{
		"",
		"can not add non-existent capabilities: [invalid-capability]",
		"",
	}
	for i, role := range testData.Roles {
		alerts, _, err := TOSession.CreateRole(role, client.RequestOptions{})
		if err != nil {
			if expectedAlerts[i] == "" {
				t.Errorf("Unexpected error creating a Role: %v - alerts: %+v", err, alerts.Alerts)
			} else if !alertsHaveError(alerts.Alerts, expectedAlerts[i]) {
				t.Errorf("expected: error containing '%s', actual: %v - alerts: %+v", expectedAlerts[i], err, alerts.Alerts)
			}
		} else if expectedAlerts[i] != "" {
			t.Errorf("expected: error containing '%s', actual: nil", expectedAlerts[i])
		}
	}
}

func SortTestRoles(t *testing.T) {
	resp, _, err := TOSession.GetRoles(client.RequestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, but got: %v - alerts: %+v", err, resp.Alerts)
	}

	sortedList := make([]string, 0, len(resp.Response))
	for _, role := range resp.Response {
		if role.Name == nil {
			t.Error("Traffic Ops returned a representation for a Role with null or undefined name")
			continue
		}
		sortedList = append(sortedList, *role.Name)
	}

	if !sort.StringsAreSorted(sortedList) {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestRoles(t *testing.T) {
	if len(testData.Roles) < 1 {
		t.Fatalf("Need at least on Role to test updating Roles")
	}
	firstRole := testData.Roles[0]
	if firstRole.Name == nil {
		t.Fatal("Found Role with null or undefined name in test data")
	}

	// Retrieve the Role by role so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstRole.Name)
	resp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role '%s' by name: %v - alerts: %+v", *firstRole.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Role to exist with name '%s', found: %d", *firstRole.Name, len(resp.Response))
	}
	remoteRole := resp.Response[0]
	if remoteRole.ID == nil {
		t.Fatal("Role returned from Traffic Ops had null or undefined ID")
	}

	expectedRole := "new admin2"
	remoteRole.Name = &expectedRole
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateRole(*remoteRole.ID, remoteRole, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot update Role: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Role to check role got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("id", strconv.Itoa(*remoteRole.ID))
	resp, _, err = TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role by ID after update: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Role to exist with ID %d, found: %d", *remoteRole.ID, len(resp.Response))
	}
	respRole := resp.Response[0]
	if respRole.Name == nil {
		t.Fatal("Traffic Ops returned a representation for a Role that had null or undefined name")
	}

	if *respRole.Name != expectedRole {
		t.Errorf("results do not match actual: %s, expected: %s", *respRole.Name, expectedRole)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRole.Name = firstRole.Name
	alert, _, err = TOSession.UpdateRole(*remoteRole.ID, remoteRole, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Role: %v - alerts: %+v", err, alert.Alerts)
	}

}

func GetTestRoles(t *testing.T) {
	if len(testData.Roles) < roleGood+1 {
		t.Fatalf("Need at least %d Roles to test getting Roles with IMS change", roleGood+1)
	}
	role := testData.Roles[roleGood]
	if role.Name == nil {
		t.Fatal("Found a Role in the testing data with null or undefined name")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *role.Name)
	resp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role '%s' by name: %v - alerts: %+v", *role.Name, err, resp)
	}

}

func VerifyGetRolesOrder(t *testing.T) {
	opts := client.RequestOptions{
		QueryParameters: url.Values{
			"orderby":   {"name"},
			"sortOrder": {"desc"},
		},
	}
	descResp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Roles: %v - alerts: %+v", err, descResp.Alerts)
	}

	opts.QueryParameters.Set("sortOrder", "asc")
	ascResp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Roles: %v - alerts: %+v", err, ascResp.Alerts)
	}

	if reflect.DeepEqual(descResp.Response, ascResp.Response) {
		t.Errorf("Role responses for descending and ascending are the same: %v - %v", descResp.Response, ascResp.Response)
	}

	// reverse the descending-sorted response and compare it to the ascending-sorted one
	for start, end := 0, len(descResp.Response)-1; start < end; start, end = start+1, end-1 {
		descResp.Response[start], descResp.Response[end] = descResp.Response[end], descResp.Response[start]
	}
	if !reflect.DeepEqual(descResp.Response, ascResp.Response) {
		t.Errorf("Role responses are not equal after reversal: %v - %v", descResp, ascResp)
	}
}

func DeleteTestRoles(t *testing.T) {
	if len(testData.Roles) < roleGood+1 {
		t.Fatalf("Need at least %d Roles to test getting Roles with IMS change", roleGood+1)
	}
	role := testData.Roles[roleGood]
	if role.Name == nil {
		t.Fatal("Found a Role in the testing data with null or undefined name")
	}

	// Retrieve the Role by name so we can get the id
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *role.Name)
	resp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role '%s' by name: %v - alerts: %+v", *role.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Role to exist with name '%s', found: %d", *role.Name, len(resp.Response))
	}
	respRole := resp.Response[0]
	if respRole.ID == nil {
		t.Fatal("Traffic Ops returned a representation for a Role that had null or undefined ID")
	}

	delResp, _, err := TOSession.DeleteRole(*respRole.ID, client.RequestOptions{})

	if err != nil {
		t.Errorf("cannot delete Role: %v - alerts: %+v", err, delResp.Alerts)
	}

	// Retrieve the Role to see if it got deleted
	roleResp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("error fetching Role after deletion: %v - alerts: %+v", err, roleResp.Alerts)
	}
	if len(roleResp.Response) > 0 {
		t.Errorf("expected Role '%s' to be deleted, but it was found in Traffic Ops", *role.Name)
	}

}
