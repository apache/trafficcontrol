package v14

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
	"testing"
)

func TestCacheAssignmentGroups(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, CacheAssignmentGroups}, func() {
		GetTestCacheAssignmentGroups(t)
		CheckCacheAssignmentGroupsAuthentication(t)
		UpdateTestCacheAssignmentGroups(t)
	})
}
func GetTestCacheAssignmentGroups(t *testing.T) {
	for _, cag := range testData.CacheAssignmentGroups {
		resp, _, err := TOSession.GetCacheAssignmentGroupByName(cag.Name)
		if err != nil {
			t.Errorf("cannot GET CacheAssignmentGroup by name: %v - %v\n", err, resp)
		}
	}
}

func CheckCacheAssignmentGroupsAuthentication(t *testing.T) {
	errFormat := "expected error from %s when unauthenticated"

	cag := testData.CacheAssignmentGroups[0]

	resp, _, err := TOSession.GetCacheAssignmentGroupByName(cag.Name)
	if err != nil || len(resp) == 0 {
		t.Errorf("cannot GET CacheAssignmentGroup by name: %v - %v\n", cag.Name, err)
	}
	cag = resp[0]

	if _, _, err = NoAuthTOSession.CreateCacheAssignmentGroup(cag); err == nil {
		t.Error(fmt.Errorf(errFormat, "CreateCacheAssignmentGroup"))
	}
	if _, _, err = NoAuthTOSession.GetCacheAssignmentGroups(); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheAssignmentGroups"))
	}
	if _, _, err = NoAuthTOSession.GetCacheAssignmentGroupByName(cag.Name); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheAssignmentGroupByName"))
	}
	if _, _, err = NoAuthTOSession.GetCacheAssignmentGroupByID(cag.ID); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheAssignmentGroupByID"))
	}
	if _, _, err = NoAuthTOSession.UpdateCacheAssignmentGroupByID(cag.ID, cag); err == nil {
		t.Error(fmt.Errorf(errFormat, "UpdateCacheAssignmentGroupByID"))
	}
	if _, _, err = NoAuthTOSession.DeleteCacheAssignmentGroupByID(cag.ID); err == nil {
		t.Error(fmt.Errorf(errFormat, "DeleteCacheAssignmentGroupByID"))
	}
}

func UpdateTestCacheAssignmentGroups(t *testing.T) {
	firstCAG := testData.CacheAssignmentGroups[0]
	resp, _, err := TOSession.GetCacheAssignmentGroupByName(firstCAG.Name)
	if err != nil {
		t.Errorf("cannot GET CacheAssignmentGroup by name: %v - %v\n", firstCAG.Name, err)
	}
	cag := resp[0]
	expectedServers := []int{5}
	expectedName := "blah"
	cag.Name = expectedName
	cag.Servers = expectedServers

	updResp, _, err := TOSession.UpdateCacheAssignmentGroupByID(cag.ID, cag)
	if err != nil {
		t.Errorf("cannot UPDATE CacheAssignmentGroup by id: %v - %v\n", err, updResp)
	}

	// Retrieve the CacheAssignmentGroup to check name got updated
	resp, _, err = TOSession.GetCacheAssignmentGroupByID(cag.ID)
	if err != nil {
		t.Errorf("cannot GET CacheAssignmentGroup by name: '%s', %v\n", firstCAG.Name, err)
	}
	cag = resp[0]
	if cag.Name != expectedName {
		t.Errorf("results do not match actual: %s, expected: %s\n", cag.Name, expectedName)
	}

	for idx, val := range expectedServers {
		if cag.Servers[idx] != val {
			t.Errorf("results do not match actual: idx: %d %d, expected: %d\n", idx, cag.Servers[idx], expectedServers)
		}
	}
}

func CreateTestCacheAssignmentGroups(t *testing.T) {
	for _, cag := range testData.CacheAssignmentGroups {
		_, _, err := TOSession.CreateCacheAssignmentGroup(cag)
		if err != nil {
			t.Errorf("could not CREATE cache assignment groups: %v, request: %v\n", err, cag)
		}
	}
}

func DeleteTestCacheAssignmentGroups(t *testing.T) {
	for _, cag := range testData.CacheAssignmentGroups {
		resp, _, err := TOSession.GetCacheAssignmentGroupByName(cag.Name)
		if err != nil {
			t.Errorf("cannot GET CacheAssignmentGroup by name: %v - %v\n", cag.Name, err)
		}

		if len(resp) > 0 {
			respCAG := resp[0]
			_, _, err := TOSession.DeleteCacheAssignmentGroupByID(respCAG.ID)
			if err != nil {
				t.Errorf("cannot DELETE CacheAssignmentGroup by name: '%s' %v\n", respCAG.Name, err)
			}

			// Retrieve the CacheAssignmentGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheAssignmentGroupByName(cag.Name)
			if err != nil {
				t.Errorf("error deleting CacheAssignmentGroup by name: %s\n", err.Error())
			}
			if len(cgs) > 0 {
				t.Errorf("expected CacheAssignmentGroup name: %s to be deleted\n", cag.Name)
			}
		}
	}
}