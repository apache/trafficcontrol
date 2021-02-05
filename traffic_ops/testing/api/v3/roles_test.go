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
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
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
	if len(testData.Roles) > 0 {
		t.Logf("testData.Roles contains: %+v\n", testData.Roles)
		firstRole := testData.Roles[0]
		// Retrieve the Role by role so we can get the id for the Update
		resp, _, status, err := TOSession.GetRoleByNameWithHdr(*firstRole.Name, header)
		t.Log("Status Code: ", status)
		if err != nil {
			t.Errorf("cannot GET Role by role: %v - %v", firstRole.Name, err)
		}
		t.Logf("got response: %+v\n", resp)
		if len(resp) > 0 {
			remoteRole := resp[0]
			expectedRole := "new admin2"
			remoteRole.Name = &expectedRole
			_, reqInf, status, _ := TOSession.UpdateRoleByIDWithHdr(*remoteRole.ID, remoteRole, header)
			if status != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		}
	}
}

func GetTestRolesIMSAfterChange(t *testing.T, header http.Header) {
	role := testData.Roles[roleGood]
	_, reqInf, _, err := TOSession.GetRoleByNameWithHdr(*role.Name, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}
	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, timeStr)
	_, reqInf, _, err = TOSession.GetRoleByNameWithHdr(*role.Name, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestRolesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	role := testData.Roles[roleGood]
	_, reqInf, _, err := TOSession.GetRoleByNameWithHdr(*role.Name, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}

}

func CreateTestRoles(t *testing.T) {
	expectedAlerts := []string{
		"",
		"can not add non-existent capabilities: [invalid-capability]",
		"",
	}
	for i, role := range testData.Roles {
		var alerts tc.Alerts
		alerts, _, status, err := TOSession.CreateRole(role)
		t.Log("Status Code: ", status)
		t.Log("Response: ", alerts)
		if err != nil {
			t.Logf("error: %v", err)
		}
		if expectedAlerts[i] == "" && err != nil {
			t.Errorf("expected: no error, actual: %v", err)
		} else if len(expectedAlerts[i]) > 0 && err == nil {
			t.Errorf("expected: error containing '%s', actual: nil", expectedAlerts[i])
		} else if err != nil && !strings.Contains(err.Error(), expectedAlerts[i]) {
			t.Errorf("expected: error containing '%s', actual: %v", expectedAlerts[i], err)
		}
	}
}

func SortTestRoles(t *testing.T) {
	var header http.Header
	var sortedList []string
	resp, _, _, err := TOSession.GetRolesWithHdr(header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	for i, _ := range resp {
		sortedList = append(sortedList, *resp[i].Name)
	}

	res := sort.SliceIsSorted(sortedList, func(p, q int) bool {
		return sortedList[p] < sortedList[q]
	})
	if res != true {
		t.Errorf("list is not sorted by their names: %v", sortedList)
	}
}

func UpdateTestRoles(t *testing.T) {
	t.Logf("testData.Roles contains: %+v\n", testData.Roles)
	firstRole := testData.Roles[0]
	// Retrieve the Role by role so we can get the id for the Update
	resp, _, status, err := TOSession.GetRoleByName(*firstRole.Name)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v", firstRole.Name, err)
	}
	t.Logf("got response: %+v", resp)
	if len(resp) < 1 {
		t.Fatal("got empty response if GET role by name")
	}
	remoteRole := resp[0]
	expectedRole := "new admin2"
	remoteRole.Name = &expectedRole
	var alert tc.Alerts
	alert, _, status, err = TOSession.UpdateRoleByID(*remoteRole.ID, remoteRole)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot UPDATE Role by id: %v - %v", err, alert)
	}

	// Retrieve the Role to check role got updated
	resp, _, status, err = TOSession.GetRoleByID(*remoteRole.ID)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v", firstRole.Name, err)
	}
	respRole := resp[0]
	if *respRole.Name != expectedRole {
		t.Errorf("results do not match actual: %s, expected: %s", *respRole.Name, expectedRole)
	}

	// Set the name back to the fixture value so we can delete it after
	remoteRole.Name = firstRole.Name
	alert, _, status, err = TOSession.UpdateRoleByID(*remoteRole.ID, remoteRole)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot UPDATE Role by id: %v - %v", err, alert)
	}

}

func GetTestRoles(t *testing.T) {
	role := testData.Roles[roleGood]
	resp, _, status, err := TOSession.GetRoleByName(*role.Name)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v", err, resp)
	}

}

func VerifyGetRolesOrder(t *testing.T) {
	params := map[string]string{
		"orderby":   "name",
		"sortOrder": "desc",
	}
	descResp, _, status, err := TOSession.GetRoleByQueryParams(params)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v", err, descResp)
	}
	params["sortOrder"] = "asc"
	ascResp, _, status, err := TOSession.GetRoleByQueryParams(params)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by role: %v - %v", err, ascResp)
	}

	if reflect.DeepEqual(descResp, ascResp) {
		t.Errorf("Role responses for descending and ascending are the same: %v - %v", descResp, ascResp)
	}

	// reverse the descending-sorted response and compare it to the ascending-sorted one
	for start, end := 0, len(descResp)-1; start < end; start, end = start+1, end-1 {
		descResp[start], descResp[end] = descResp[end], descResp[start]
	}
	if !reflect.DeepEqual(descResp, ascResp) {
		t.Errorf("Role responses are not equal after reversal: %v - %v", descResp, ascResp)
	}
}

func DeleteTestRoles(t *testing.T) {

	role := testData.Roles[roleGood]
	// Retrieve the Role by name so we can get the id
	resp, _, status, err := TOSession.GetRoleByName(*role.Name)
	t.Log("Status Code: ", status)
	if err != nil {
		t.Errorf("cannot GET Role by name: %v - %v", role.Name, err)
	}
	respRole := resp[0]

	delResp, _, status, err := TOSession.DeleteRoleByID(*respRole.ID)
	t.Log("Status Code: ", status)

	if err != nil {
		t.Errorf("cannot DELETE Role by role: %v - %v", err, delResp)
	}

	// Retrieve the Role to see if it got deleted
	roleResp, _, status, err := TOSession.GetRoleByName(*role.Name)
	t.Log("Status Code: ", status)

	if err != nil {
		t.Errorf("error deleting Role role: %s", err.Error())
	}
	if len(roleResp) > 0 {
		t.Errorf("expected Role : %s to be deleted", *role.Name)
	}

}
