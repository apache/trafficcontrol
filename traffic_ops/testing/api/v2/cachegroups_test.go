package v2

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

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
)

func TestCacheGroups(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups}, func() {
		GetTestCacheGroupsByName(t)
		GetTestCacheGroupsByShortName(t)
		CheckCacheGroupsAuthentication(t)
		UpdateTestCacheGroups(t)
	})
}

func CreateTestCacheGroups(t *testing.T) {

	var err error
	var resp *tc.CacheGroupDetailResponse

	for _, cg := range testData.CacheGroups {

		resp, _, err = TOSession.CreateCacheGroupNullable(cg)
		if err != nil {
			t.Errorf("could not CREATE cachegroups: %v, request: %v", err, cg)
			continue
		}
		if resp == nil {
			t.Error("Traffic Ops returned a null or undefined Cache Group in creation response")
			continue
		}

		// Testing 'join' fields during create
		if cg.ParentName != nil && resp.Response.ParentName == nil {
			t.Error("Parent cachegroup is null in response when it should have a value")
		}
		if cg.SecondaryParentName != nil && resp.Response.SecondaryParentName == nil {
			t.Error("Secondary parent cachegroup is null in response when it should have a value\n")
		}
		if cg.Type != nil && resp.Response.Type == nil {
			t.Error("Type is null in response when it should have a value\n")
		}
		if resp.Response.LocalizationMethods == nil {
			t.Error("Localization methods are null")
		}
		if resp.Response.Fallbacks == nil {
			t.Error("Fallbacks are null")
		}

	}
}

func GetTestCacheGroupsByName(t *testing.T) {
	for _, cg := range testData.CacheGroups {
		resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v", err, resp)
		}
	}
}

func GetTestCacheGroupsByShortName(t *testing.T) {
	for _, cg := range testData.CacheGroups {
		resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.ShortName)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by shortName: %v - %v", err, resp)
		}
	}
}

func UpdateTestCacheGroups(t *testing.T) {
	firstCG := testData.CacheGroups[0]
	resp, _, err := TOSession.GetCacheGroupNullableByName(*firstCG.Name)
	if err != nil {
		t.Errorf("cannot GET CACHEGROUP by name: %v - %v", *firstCG.Name, err)
	}
	cg := resp[0]
	expectedShortName := "blah"
	cg.ShortName = &expectedShortName

	// fix the type id for test
	typeResp, _, err := TOSession.GetTypeByID(*cg.TypeID)
	if err != nil {
		t.Error("could not lookup a typeID for this cachegroup")
	}
	cg.TypeID = &typeResp[0].ID

	updResp, _, err := TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	// Check response to make sure fields aren't null
	if cg.ParentName != nil && updResp.Response.ParentName == nil {
		t.Error("Parent cachegroup is null in response when it should have a value")
	}
	if cg.SecondaryParentName != nil && updResp.Response.SecondaryParentName == nil {
		t.Error("Secondary parent cachegroup is null in response when it should have a value\n")
	}
	if cg.Type != nil && updResp.Response.Type == nil {
		t.Error("Type is null in response when it should have a value\n")
	}
	if updResp.Response.LocalizationMethods == nil {
		t.Error("Localization methods are null")
	}
	if updResp.Response.Fallbacks == nil {
		t.Error("Fallbacks are null")
	}

	// Retrieve the CacheGroup to check CacheGroup name got updated
	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: '%s', %v", *firstCG.Name, err)
	}
	cg = resp[0]
	if *cg.ShortName != expectedShortName {
		t.Errorf("results do not match actual: %s, expected: %s", *cg.ShortName, expectedShortName)
	}

	// test coordinate updates
	expectedLat := 7.0
	expectedLong := 8.0
	cg.Latitude = &expectedLat
	cg.Longitude = &expectedLong
	updResp, _, err = TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	cg = resp[0]
	if *cg.Latitude != expectedLat {
		t.Errorf("failed to update latitude (expected = %f, actual = %f)", expectedLat, *cg.Latitude)
	}
	if *cg.Longitude != expectedLong {
		t.Errorf("failed to update longitude (expected = %f, actual = %f)", expectedLong, *cg.Longitude)
	}

	// test localizationMethods
	expectedMethods := []tc.LocalizationMethod{tc.LocalizationMethodGeo}
	cg.LocalizationMethods = &expectedMethods
	updResp, _, err = TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	cg = resp[0]
	if !reflect.DeepEqual(expectedMethods, *cg.LocalizationMethods) {
		t.Errorf("failed to update localizationMethods (expected = %v, actual = %v)", expectedMethods, *cg.LocalizationMethods)
	}

	// test cachegroup fallbacks

	// Retrieve the CacheGroup to check CacheGroup name got updated
	firstEdgeCGName := "cachegroup1"
	resp, _, err = TOSession.GetCacheGroupNullableByName(firstEdgeCGName)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: '$%s', %v", firstEdgeCGName, err)
	}
	cg = resp[0]
	if *cg.Name != firstEdgeCGName {
		t.Errorf("results do not match actual: %s, expected: %s", *cg.ShortName, firstEdgeCGName)
	}

	// Test adding fallbacks when previously nil
	expectedFallbacks := []string{"fallback1", "fallback2"}
	cg.Fallbacks = &expectedFallbacks
	updResp, _, err = TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	cg = resp[0]
	if !reflect.DeepEqual(expectedFallbacks, *cg.Fallbacks) {
		t.Errorf("failed to update fallbacks (expected = %v, actual = %v)", expectedFallbacks, *cg.Fallbacks)
	}

	// Test adding fallback to existing list
	expectedFallbacks = []string{"fallback1", "fallback2", "fallback3"}
	cg.Fallbacks = &expectedFallbacks
	updResp, _, err = TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v)", err, updResp)
	}

	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	cg = resp[0]
	if !reflect.DeepEqual(expectedFallbacks, *cg.Fallbacks) {
		t.Errorf("failed to update fallbacks (expected = %v, actual = %v)", expectedFallbacks, *cg.Fallbacks)
	}

	// Test removing fallbacks
	expectedFallbacks = []string{}
	cg.Fallbacks = &expectedFallbacks
	updResp, _, err = TOSession.UpdateCacheGroupNullableByID(*cg.ID, cg)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	resp, _, err = TOSession.GetCacheGroupNullableByID(*cg.ID)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	cg = resp[0]
	if !reflect.DeepEqual(expectedFallbacks, *cg.Fallbacks) {
		t.Errorf("failed to update fallbacks (expected = %v, actual = %v)", expectedFallbacks, *cg.Fallbacks)
	}
}

func DeleteTestCacheGroups(t *testing.T) {
	var mids []tc.CacheGroupNullable

	// delete the edge caches.
	for _, cg := range testData.CacheGroups {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v", *cg.Name, err)
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
				t.Errorf("cannot DELETE CacheGroup by name: '%s' %v", *respCG.Name, err)
			}
			// Retrieve the CacheGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
			if err != nil {
				t.Errorf("error deleting CacheGroup by name: %s", err.Error())
			}
			if len(cgs) > 0 {
				t.Errorf("expected CacheGroup name: %s to be deleted", *cg.Name)
			}
		}
	}
	// now delete the mid tier caches
	for _, cg := range mids {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v", *cg.Name, err)
		}
		if len(resp) > 0 {
			respCG := resp[0]
			_, _, err := TOSession.DeleteCacheGroupByID(*respCG.ID)
			if err != nil {
				t.Errorf("cannot DELETE CacheGroup by name: '%s' %v", *respCG.Name, err)
			}

			// Retrieve the CacheGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
			if err != nil {
				t.Errorf("error deleting CacheGroup name: %s", err.Error())
			}
			if len(cgs) > 0 {
				t.Errorf("expected CacheGroup name: %s to be deleted", *cg.Name)
			}
		}
	}
}

func CheckCacheGroupsAuthentication(t *testing.T) {
	errFormat := "expected error from %s when unauthenticated"

	cg := testData.CacheGroups[0]

	resp, _, err := TOSession.GetCacheGroupNullableByName(*cg.Name)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: %v - %v", *cg.Name, err)
	}
	cg = resp[0]

	if _, _, err = NoAuthTOSession.CreateCacheGroupNullable(cg); err == nil {
		t.Error(fmt.Errorf(errFormat, "CreateCacheGroup"))
	}
	if _, _, err = NoAuthTOSession.GetCacheGroupsNullable(); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheGroups"))
	}
	if _, _, err = NoAuthTOSession.GetCacheGroupNullableByName(*cg.Name); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheGroupByName"))
	}
	if _, _, err = NoAuthTOSession.GetCacheGroupNullableByID(*cg.ID); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheGroupByID"))
	}
	if _, _, err = NoAuthTOSession.UpdateCacheGroupNullableByID(*cg.ID, cg); err == nil {
		t.Error(fmt.Errorf(errFormat, "UpdateCacheGroupByID"))
	}
	if _, _, err = NoAuthTOSession.DeleteCacheGroupByID(*cg.ID); err == nil {
		t.Error(fmt.Errorf(errFormat, "DeleteCacheGroupByID"))
	}
}
