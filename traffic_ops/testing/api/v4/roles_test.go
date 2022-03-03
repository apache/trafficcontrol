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
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/routing"
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
		CheckAllRoutePermissionsExist(t)
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

func CheckAllRoutePermissionsExist(t *testing.T) {
	fake := routing.ServerData{Config: config.NewFakeConfig()}
	routes, _, err := routing.Routes(fake)
	if err != nil {
		t.Fatalf("error getting routes: %v\n", err)
	}

	db, err := OpenConnection()
	if err != nil {
		t.Fatalf("error opening db connection: %v\n", err)
	}
	defer log.Close(db, "closing connection")
	query := "SELECT DISTINCT cap_name from role_capability ORDER BY cap_name"
	perms := make(map[string]struct{})

	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("error querying db: %v\n", err)
	}
	defer log.Close(rows, "closing query")

	for rows.Next() {
		capName := ""
		if err := rows.Scan(&capName); err != nil {
			t.Fatalf("unable to scan row: %v\n", err)
		}
		if _, ok := perms[capName]; ok {
			continue
		}
		perms[capName] = struct{}{}
	}

	var missing []string
	for _, route := range routes {
		if route.Version.Major != 4 {
			continue
		}
		for _, perm := range route.RequiredPermissions {
			if _, ok := perms[perm]; !ok {
				missing = append(missing, fmt.Sprintf("%v (%v)", perm, route.Path))
			}
		}
	}

	if len(missing) > 0 {
		t.Fatalf("found several permissions in routes: %v\n", strings.Join(missing, ","))
	}
}

func UpdateTestRolesWithHeaders(t *testing.T, header http.Header) {
	if len(testData.Roles) < 1 {
		t.Fatal("Need at least one Role to test updating a Role with HTTP headers")
	}
	firstRole := testData.Roles[0]

	// Retrieve the Role by role so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.Header = header
	opts.QueryParameters.Set("name", firstRole.Name)
	resp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role '%s' by name: %v - alerts: %+v", firstRole.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Role to exist with name '%s', found: %d", firstRole.Name, len(resp.Response))
	}
	remoteRole := resp.Response[0]
	expectedDescription := "new description"
	remoteRole.Description = expectedDescription
	opts.QueryParameters.Del("name")
	_, reqInf, err := TOSession.UpdateRole(remoteRole.Name, remoteRole, opts)
	if err == nil {
		t.Errorf("updating role with name: %s, expected an error stating resource was modified, but got nothing", remoteRole.Name)
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestRolesIMSAfterChange(t *testing.T, header http.Header) {
	if len(testData.Roles) < roleGood+1 {
		t.Fatalf("Need at least %d Roles to test getting Roles with IMS change", roleGood+1)
	}
	role := testData.Roles[roleGood]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", role.Name)
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
	for _, role := range testData.Roles {
		_, _, err := TOSession.CreateRole(role, client.RequestOptions{})
		if err != nil {
			t.Errorf("no error expected, but got %v", err)
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
		sortedList = append(sortedList, role.Name)
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
	// Retrieve the Role by role so we can get the id for the Update
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", firstRole.Name)
	resp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role '%s' by name: %v - alerts: %+v", firstRole.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Role to exist with name '%s', found: %d", firstRole.Name, len(resp.Response))
	}
	remoteRole := resp.Response[0]
	expectedDescription := "new description"
	remoteRole.Description = expectedDescription
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateRole(remoteRole.Name, remoteRole, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot update Role: %v - alerts: %+v", err, alert.Alerts)
	}

	// Retrieve the Role to check role got updated
	opts.QueryParameters.Del("name")
	opts.QueryParameters.Set("name", remoteRole.Name)
	resp, _, err = TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role by ID after update: %v - alerts: %+v", err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Role to exist with name %s, found: %d", remoteRole.Name, len(resp.Response))
	}
	respRole := resp.Response[0]

	if respRole.Description != expectedDescription {
		t.Errorf("results do not match actual: %s, expected: %s", respRole.Description, expectedDescription)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRole.Name = firstRole.Name
	alert, _, err = TOSession.UpdateRole(remoteRole.Name, remoteRole, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Role: %v - alerts: %+v", err, alert.Alerts)
	}

}

func GetTestRoles(t *testing.T) {
	if len(testData.Roles) < roleGood+1 {
		t.Fatalf("Need at least %d Roles to test getting Roles with IMS change", roleGood+1)
	}
	role := testData.Roles[roleGood]

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", role.Name)
	resp, _, err := TOSession.GetRoles(opts)
	if err != nil {
		t.Errorf("cannot get Role '%s' by name: %v - alerts: %+v", role.Name, err, resp)
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
	for _, r := range testData.Roles {
		_, _, err := TOSession.DeleteRole(r.Name, client.NewRequestOptions())
		if err != nil {
			t.Errorf("expected no error while deleting role %s, but got %v", r.Name, err)
		}
	}

	resp, reqInf, err := TOSession.DeleteRole(tc.AdminRoleName, client.NewRequestOptions())
	if err == nil {
		t.Errorf("Expected an error trying to delete the '%s' Role, but didn't get one", tc.AdminRoleName)
	}
	if reqInf.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected a %d response, got: %d %s", http.StatusBadRequest, reqInf.StatusCode, http.StatusText(reqInf.StatusCode))
	}
	if !strings.Contains(resp.ErrorString(), tc.AdminRoleName) {
		t.Errorf("Expected an error-level alert that mentions the special '%s' Role, got: %s", tc.AdminRoleName, resp.ErrorString())
	}
}
