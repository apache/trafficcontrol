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
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
)

func TestCacheGroups(t *testing.T) {
	CreateTestTypes(t)
	CreateTestCacheGroups(t)
	GetTestCacheGroups(t)
	CheckCacheGroupsAuthentication(t)
	UpdateTestCacheGroups(t)
	DeleteTestCacheGroups(t)
	DeleteTestTypes(t)
}

func CreateTestCacheGroups(t *testing.T) {
	failed := false

	for _, cg := range testData.CacheGroups {
		// get the typeID
		typeResp, _, err := TOSession.GetTypeByName(*cg.Type)
		if err != nil {
			t.Error("could not lookup a typeID for this cachegroup")
			failed = true
		}
		cg.TypeID = &typeResp[0].ID

		if cg.ParentName != nil && *cg.ParentName != "" {
			// get parent cachegroup ID (must already be created)
			resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.ParentName)
			if err != nil {
				t.Errorf("cannot GET CacheGroup by name: %v - %v\n", err, resp)
				failed = true
			}
			cg.ParentCachegroupID = resp[0].ID
			cg.ParentName = nil // to guarantee that parent IDs aren't looked up by name
		}

		if cg.SecondaryParentName != nil && *cg.SecondaryParentName != "" {
			// get secondary parent cachegroup ID (must already be created)
			resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.SecondaryParentName)
			if err != nil {
				t.Errorf("cannot GET CacheGroup by name: %v - %v\n", err, resp)
				failed = true
			}
			cg.SecondaryParentCachegroupID = resp[0].ID
			cg.SecondaryParentName = nil // to guarantee that parent IDs aren't looked up by name
		}

		_, _, err = TOSession.CreateCacheGroupNullable(cg)
		if err != nil {
			t.Errorf("could not CREATE cachegroups: %v, request: %v\n", err, cg)
			failed = true
		}
	}
	if !failed {
		log.Debugln("CreateTestCacheGroups() PASSED: ")
	}
}

func GetTestCacheGroups(t *testing.T) {
	failed := false
	for _, cg := range testData.CacheGroups {
		resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v\n", err, resp)
			failed = true
		}
	}
	if !failed {
		log.Debugln("GetTestCacheGroups() PASSED: ")
	}
}

func UpdateTestCacheGroups(t *testing.T) {
	failed := false
	firstCG := testData.CacheGroups[0]
	resp, _, err := TOSession.GetCacheGroupNullableByName(*firstCG.Name)
	if err != nil {
		t.Errorf("cannot GET CACHEGROUP by name: %v - %v\n", *firstCG.Name, err)
		failed = true
	}
	cg := resp[0]
	expectedShortName := "blah"
	cg.ShortName = &expectedShortName

	// fix the type id for test
	typeResp, _, err := TOSession.GetTypeByID(*cg.TypeID)
	if err != nil {
		t.Error("could not lookup a typeID for this cachegroup")
		failed = true
	}
	cg.TypeID = &typeResp[0].ID

	updResp, _, err := TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v\n", err, updResp)
		failed = true
	}

	// Retrieve the CacheGroup to check CacheGroup name got updated
	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: '$%s', %v\n", *firstCG.Name, err)
		failed = true
	}
	cg = resp[0]
	if *cg.ShortName != expectedShortName {
		t.Errorf("results do not match actual: %s, expected: %s\n", *cg.ShortName, expectedShortName)
		failed = true
	}

	// test coordinate updates
	expectedLat := 7.0
	expectedLong := 8.0
	cg.Latitude = &expectedLat
	cg.Longitude = &expectedLong
	updResp, _, err = TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v\n", err, updResp)
		failed = true
	}

	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v\n", *cg.ID, err)
		failed = true
	}
	cg = resp[0]
	if *cg.Latitude != expectedLat {
		t.Errorf("failed to update latitude (expected = %f, actual = %f)\n", expectedLat, *cg.Latitude)
		failed = true
	}
	if *cg.Longitude != expectedLong {
		t.Errorf("failed to update longitude (expected = %f, actual = %f)\n", expectedLong, *cg.Longitude)
		failed = true
	}

	// test localizationMethods
	expectedMethods := []tc.LocalizationMethod{tc.LocalizationMethodGeo}
	cg.LocalizationMethods = &expectedMethods
	updResp, _, err = TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v\n", err, updResp)
		failed = true
	}

	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v\n", *cg.ID, err)
		failed = true
	}
	cg = resp[0]
	if !reflect.DeepEqual(expectedMethods, *cg.LocalizationMethods) {
		t.Errorf("failed to update localizationMethods (expected = %v, actual = %v\n", expectedMethods, *cg.LocalizationMethods)
		failed = true
	}

	if !failed {
		log.Debugln("UpdateTestCacheGroups() PASSED: ")
	}
}

func DeleteTestCacheGroups(t *testing.T) {
	failed := false
	var mids []tc.CacheGroupNullable

	// delete the edge caches.
	for _, cg := range testData.CacheGroups {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v\n", *cg.Name, err)
			failed = true
		}
		// Mids are parents and need to be deleted only after the children
		// cachegroups are deleted.
		if *cg.Type == "MID_LOC" {
			mids = append(mids, cg)
			continue
		}
		if len(resp) > 0 {
			respCG := resp[0]
			_, _, err := TOSession.DeleteCacheGroupByID(*respCG.ID)
			if err != nil {
				t.Errorf("cannot DELETE CacheGroup by name: '%s' %v\n", *respCG.Name, err)
				failed = true
			}
			// Retrieve the CacheGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
			if err != nil {
				t.Errorf("error deleting CacheGroup by name: %s\n", err.Error())
				failed = true
			}
			if len(cgs) > 0 {
				t.Errorf("expected CacheGroup name: %s to be deleted\n", *cg.Name)
				failed = true
			}
		}
	}
	// now delete the mid tier caches
	for _, cg := range mids {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v\n", *cg.Name, err)
			failed = true
		}
		if len(resp) > 0 {
			respCG := resp[0]
			_, _, err := TOSession.DeleteCacheGroupByID(*respCG.ID)
			if err != nil {
				t.Errorf("cannot DELETE CacheGroup by name: '%s' %v\n", *respCG.Name, err)
				failed = true
			}

			// Retrieve the CacheGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
			if err != nil {
				t.Errorf("error deleting CacheGroup name: %s\n", err.Error())
				failed = true
			}
			if len(cgs) > 0 {
				t.Errorf("expected CacheGroup name: %s to be deleted\n", *cg.Name)
				failed = true
			}
		}
	}

	if !failed {
		log.Debugln("DeleteTestCacheGroups() PASSED: ")
	}
}

func CheckCacheGroupsAuthentication(t *testing.T) {
	failed := false
	errFormat := "expected error from %s when unauthenticated"

	cg := testData.CacheGroups[0]

	resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: %v - %v\n", *cg.Name, err)
		failed = true
	}
	cg = resp[0]

	errors := make([]utils.ErrorAndMessage, 0)

	_, _, err = NoAuthTOSession.CreateCacheGroupNullable(cg)
	errors = append(errors, utils.ErrorAndMessage{err, fmt.Sprintf(errFormat, "CreateCacheGroup")})

	_, _, err = NoAuthTOSession.GetCacheGroupsNullable()
	errors = append(errors, utils.ErrorAndMessage{err, fmt.Sprintf(errFormat, "GetCacheGroups")})

	_, _, err = NoAuthTOSession.GetCacheGroupNullableByName(*cg.Name)
	errors = append(errors, utils.ErrorAndMessage{err, fmt.Sprintf(errFormat, "GetCacheGroupByName")})

	_, _, err = NoAuthTOSession.GetCacheGroupNullableByID(*cg.ID)
	errors = append(errors, utils.ErrorAndMessage{err, fmt.Sprintf(errFormat, "GetCacheGroupByID")})

	_, _, err = NoAuthTOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	errors = append(errors, utils.ErrorAndMessage{err, fmt.Sprintf(errFormat, "UpdateCacheGroupByID")})

	_, _, err = NoAuthTOSession.DeleteCacheGroupByID(*cg.ID)
	errors = append(errors, utils.ErrorAndMessage{err, fmt.Sprintf(errFormat, "DeleteCacheGroupByID")})

	for _, err := range errors {
		if err.Error == nil {
			t.Error(err.Message)
			failed = true
		}
	}

	if !failed {
		log.Debugln("TestCacheGroupsAuthentication() PASSED: ")
	}
}
