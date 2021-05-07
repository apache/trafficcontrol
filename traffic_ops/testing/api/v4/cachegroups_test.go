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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCacheGroups(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CDNs, Profiles, Statuses, Divisions, Regions, PhysLocations, Servers, Topologies}, func() {
		GetTestCacheGroupsIMS(t)
		GetTestCacheGroupsByNameIMS(t)
		GetTestCacheGroupsByShortNameIMS(t)
		GetTestCacheGroups(t)
		GetTestCacheGroupsByName(t)
		GetTestCacheGroupsByShortName(t)
		GetTestCacheGroupsByTopology(t)
		CheckCacheGroupsAuthentication(t)
		currentTime := time.Now().UTC().Add(-5 * time.Second)
		time := currentTime.Format(time.RFC1123)
		var header http.Header
		header = make(map[string][]string)
		header.Set(rfc.IfModifiedSince, time)
		header.Set(rfc.IfUnmodifiedSince, time)
		UpdateTestCacheGroups(t)
		UpdateTestCacheGroupsWithHeaders(t, header)
		GetTestCacheGroupsAfterChangeIMS(t, header)
		header = make(map[string][]string)
		etag := rfc.ETag(currentTime)
		header.Set(rfc.IfMatch, etag)
		UpdateTestCacheGroupsWithHeaders(t, header)
		GetTestPaginationSupportCg(t)
	})
}

func UpdateTestCacheGroupsWithHeaders(t *testing.T, h http.Header) {
	firstCG := testData.CacheGroups[0]
	resp, _, err := TOSession.GetCacheGroupByName(*firstCG.Name, h)
	if err != nil {
		t.Errorf("cannot GET CACHEGROUP by name: %v - %v", *firstCG.Name, err)
	}
	if len(resp) > 0 {
		cg := resp[0]
		expectedShortName := "blah"
		cg.ShortName = &expectedShortName

		// fix the type id for test
		typeResp, _, err := TOSession.GetTypeByID(*cg.TypeID, h)
		if err != nil {
			t.Fatalf("could not lookup a typeID for this cachegroup: %v", err.Error())
		}
		if len(typeResp) > 0 {
			cg.TypeID = &typeResp[0].ID
			_, reqInf, err := TOSession.UpdateCacheGroup(*cg.ID, cg, h)
			if err == nil {
				t.Errorf("Expected an error showing Precondition Failed, got none")
			}
			if reqInf.StatusCode != http.StatusPreconditionFailed {
				t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
			}
		}
	}
}

func GetTestCacheGroupsAfterChangeIMS(t *testing.T, header http.Header) {
	_, reqInf, err := TOSession.GetCacheGroups(nil, header)
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
	_, reqInf, err = TOSession.GetCacheGroups(nil, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestCacheGroupsByShortNameIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, cg := range testData.CacheGroups {
		if cg.ShortName == nil {
			t.Error("found Cache Group with null or undefined 'short name' in test data")
			continue
		}
		_, reqInf, err := TOSession.GetCacheGroupByShortName(*cg.ShortName, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCacheGroupsByNameIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("found Cache Group with null or undefined name in test data")
			continue
		}
		_, reqInf, err := TOSession.GetCacheGroupByName(*cg.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCacheGroupsIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	_, reqInf, err := TOSession.GetCacheGroups(nil, header)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err.Error())
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestCacheGroups(t *testing.T) {

	var err error
	var resp *tc.CacheGroupDetailResponse

	for _, cg := range testData.CacheGroups {

		resp, _, err = TOSession.CreateCacheGroup(cg)
		if err != nil {
			t.Errorf("could not CREATE cachegroups: %v, request: %v", err, cg)
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

func GetTestCacheGroups(t *testing.T) {
	resp, _, err := TOSession.GetCacheGroups(nil, nil)
	if err != nil {
		t.Errorf("cannot GET CacheGroups %v - %v", err, resp)
	}
	expectedCachegroups := make(map[string]struct{})
	for _, cg := range testData.CacheGroups {
		expectedCachegroups[*cg.Name] = struct{}{}
	}
	foundCachegroups := make(map[string]struct{})
	for _, cg := range resp {
		if _, expected := expectedCachegroups[*cg.Name]; !expected {
			t.Errorf("got unexpected cachegroup: %s", *cg.Name)
		}
		if _, found := foundCachegroups[*cg.Name]; !found {
			foundCachegroups[*cg.Name] = struct{}{}
		} else {
			t.Errorf("GET returned duplicate cachegroup: %s", *cg.Name)
		}
	}
}

func GetTestCacheGroupsByName(t *testing.T) {
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("found Cache Group with null or undefined name in test data")
			continue
		}
		resp, _, err := TOSession.GetCacheGroupByName(*cg.Name, nil)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v", err, resp)
		}
		if *resp[0].Name != *cg.Name {
			t.Errorf("name expected: %s, actual: %s", *cg.Name, *resp[0].Name)
		}
	}
}

func GetTestCacheGroupsByShortName(t *testing.T) {
	for _, cg := range testData.CacheGroups {
		if cg.ShortName == nil {
			t.Error("found Cache Group with null or undefined 'short name' in test data")
			continue
		}
		resp, _, err := TOSession.GetCacheGroupByShortName(*cg.ShortName, nil)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by shortName: %v - %v", err, resp)
		}
		if *resp[0].ShortName != *cg.ShortName {
			t.Errorf("short name expected: %s, actual: %s", *cg.ShortName, *resp[0].ShortName)
		}
	}
}

func GetTestCacheGroupsByTopology(t *testing.T) {
	for _, top := range testData.Topologies {
		qparams := url.Values{}
		qparams.Set("topology", top.Name)
		resp, _, err := TOSession.GetCacheGroups(qparams, nil)
		if err != nil {
			t.Errorf("cannot GET CacheGroups by topology: %v - %v", err, resp)
		}
		expectedCGs := topologyCachegroups(top)
		for _, cg := range resp {
			if _, exists := expectedCGs[*cg.Name]; !exists {
				t.Errorf("GET cachegroups by topology - expected one of: %v, actual: %s", expectedCGs, *cg.Name)
			}
		}
	}
}

func topologyCachegroups(top tc.Topology) map[string]struct{} {
	res := make(map[string]struct{})
	for _, node := range top.Nodes {
		res[node.Cachegroup] = struct{}{}
	}
	return res
}

func UpdateTestCacheGroups(t *testing.T) {
	if len(testData.CacheGroups) < 1 {
		t.Fatal("Need at least one Cache Group to test updating Cache Groups")
	}
	firstCG := testData.CacheGroups[0]
	if firstCG.Name == nil {
		t.Fatal("Cache Group selected for testing had a null or undefined name")
	}
	resp, _, err := TOSession.GetCacheGroupByName(*firstCG.Name, nil)
	if err != nil {
		t.Fatalf("cannot GET CACHEGROUP by name: %v - %v", *firstCG.Name, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg := resp[0]
	if cg.TypeID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined typeId")
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined id")
	}
	expectedShortName := "blah"
	cg.ShortName = &expectedShortName

	// fix the type id for test
	typeResp, _, err := TOSession.GetTypeByID(*cg.TypeID, nil)
	if err != nil {
		t.Error("could not lookup a typeID for this cachegroup")
	}
	if len(typeResp) == 0 {
		t.Fatal("got an empty response for types")
	}
	cg.TypeID = &typeResp[0].ID
	updResp, _, err := TOSession.UpdateCacheGroup(*cg.ID, cg, nil)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	if updResp == nil {
		t.Fatal("could not update cachegroup by ID, got nil response")
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
	resp, _, err = TOSession.GetCacheGroupByID(*cg.ID, nil)
	if err != nil {
		t.Fatalf("cannot GET CacheGroup by name: '%s', %v", *firstCG.Name, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if *cg.ShortName != expectedShortName {
		t.Errorf("results do not match actual: %s, expected: %s", *cg.ShortName, expectedShortName)
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined id")
	}

	// test coordinate updates
	expectedLat := 7.0
	expectedLong := 8.0
	cg.Latitude = &expectedLat
	cg.Longitude = &expectedLong
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, nil)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	if updResp == nil {
		t.Fatal("could not update cachegroup by ID, got nil response")
	}
	resp, _, err = TOSession.GetCacheGroupByID(*cg.ID, nil)

	if err != nil {
		t.Fatalf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if cg.Latitude == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined latitude")
	}
	if cg.Longitude == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined longitude")
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined id")
	}
	if *cg.Latitude != expectedLat {
		t.Errorf("failed to update latitude (expected = %f, actual = %f)", expectedLat, *cg.Latitude)
	}
	if *cg.Longitude != expectedLong {
		t.Errorf("failed to update longitude (expected = %f, actual = %f)", expectedLong, *cg.Longitude)
	}

	// test localizationMethods
	expectedMethods := []tc.LocalizationMethod{tc.LocalizationMethodGeo}
	cg.LocalizationMethods = &expectedMethods
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, nil)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	if updResp == nil {
		t.Fatal("could not update cachegroup by ID, got nil response")
	}
	resp, _, err = TOSession.GetCacheGroupByID(*cg.ID, nil)

	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if !reflect.DeepEqual(expectedMethods, *cg.LocalizationMethods) {
		t.Errorf("failed to update localizationMethods (expected = %v, actual = %v)", expectedMethods, *cg.LocalizationMethods)
	}

	// test cachegroup fallbacks

	// Retrieve the CacheGroup to check CacheGroup name got updated
	firstEdgeCGName := "cachegroup1"
	resp, _, err = TOSession.GetCacheGroupByName(firstEdgeCGName, nil)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: '$%s', %v", firstEdgeCGName, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if cg.Name == nil {
		t.Error("Cache Group returned from Traffic Ops had null or undefined name")
	} else if *cg.Name != firstEdgeCGName {
		t.Errorf("results do not match actual: %s, expected: %s", *cg.ShortName, firstEdgeCGName)
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined id")
	}

	// Test adding fallbacks when previously nil
	expectedFallbacks := []string{"fallback1", "fallback2"}
	cg.Fallbacks = &expectedFallbacks
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, nil)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	if updResp == nil {
		t.Fatal("could not update cachegroup by ID, got nil response")
	}
	resp, _, err = TOSession.GetCacheGroupByID(*cg.ID, nil)

	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if cg.Fallbacks == nil {
		t.Error("Cache Group returned by Traffic Ops had null or undefined fallbacks")
	} else if !reflect.DeepEqual(expectedFallbacks, *cg.Fallbacks) {
		t.Errorf("failed to update fallbacks (expected = %v, actual = %v)", expectedFallbacks, *cg.Fallbacks)
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined ID")
	}

	// Test adding fallback to existing list
	expectedFallbacks = []string{"fallback1", "fallback2", "fallback3"}
	cg.Fallbacks = &expectedFallbacks
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, nil)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v)", err, updResp)
	}

	if updResp == nil {
		t.Fatal("could not update cachegroup by ID, got nil response")
	}
	resp, _, err = TOSession.GetCacheGroupByID(*cg.ID, nil)

	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if cg.Fallbacks == nil {
		t.Error("Cache Group returned by Traffic Ops had null or undefined fallbacks")
	} else if !reflect.DeepEqual(expectedFallbacks, *cg.Fallbacks) {
		t.Errorf("failed to update fallbacks (expected = %v, actual = %v)", expectedFallbacks, *cg.Fallbacks)
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined ID")
	}

	// Test removing fallbacks
	expectedFallbacks = []string{}
	cg.Fallbacks = &expectedFallbacks
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, nil)
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - %v", err, updResp)
	}

	if updResp == nil {
		t.Fatal("could not update cachegroup by ID, got nil response")
	}
	resp, _, err = TOSession.GetCacheGroupByID(*cg.ID, nil)

	if err != nil {
		t.Errorf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if cg.Fallbacks == nil {
		t.Error("Cache Group returned by Traffic Ops had null or undefined fallbacks")
	} else if !reflect.DeepEqual(expectedFallbacks, *cg.Fallbacks) {
		t.Errorf("failed to update fallbacks (expected = %v, actual = %v)", expectedFallbacks, *cg.Fallbacks)
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined ID")
	}

	const topologyEdgeCGName = "topology-edge-cg-01"
	resp, _, err = TOSession.GetCacheGroupByName(topologyEdgeCGName, nil)
	if err != nil {
		t.Fatalf("cannot GET CacheGroup by name: '$%s', %v", topologyEdgeCGName, err)
	}
	if len(resp) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp[0]
	if cg.TypeID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined typeId")
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined ID")
	}

	var cacheGroupEdgeType, cacheGroupMidType tc.Type
	types, _, err := TOSession.GetTypes(nil)
	if err != nil {
		t.Fatalf("unable to get types: %s", err.Error())
	}
	for _, typeObject := range types {
		switch typeObject.Name {
		case tc.CacheGroupEdgeTypeName:
			cacheGroupEdgeType = typeObject
		case tc.CacheGroupMidTypeName:
			cacheGroupMidType = typeObject
		}
	}
	if *cg.TypeID != cacheGroupEdgeType.ID {
		t.Fatalf("expected cachegroup %s to have type %s, actual type was %s", topologyEdgeCGName, tc.CacheGroupEdgeTypeName, *cg.Type)
	}
	*cg.TypeID = cacheGroupMidType.ID
	_, reqInfo, err := TOSession.UpdateCacheGroup(*cg.ID, cg, nil)
	if err == nil {
		t.Fatalf("expected an error when updating the type of cache group %s because it is assigned to a topology, actual error was nil", *cg.Name)
	}
	if reqInfo.StatusCode < http.StatusBadRequest || reqInfo.StatusCode >= http.StatusInternalServerError {
		t.Fatalf("expected to receive status code %d but received status code %d: %s", http.StatusBadRequest, reqInfo.StatusCode, err.Error())
	}
}

func DeleteTestCacheGroups(t *testing.T) {
	var parentlessCacheGroups []tc.CacheGroupNullable

	// delete the edge caches.
	for _, cg := range testData.CacheGroups {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupByName(*cg.Name, nil)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name: %v - %v", *cg.Name, err)
		}
		cg = resp[0]

		// Cachegroups that are parents (usually mids but sometimes edges)
		// need to be deleted only after the children cachegroups are deleted.
		if cg.ParentCachegroupID == nil && cg.SecondaryParentCachegroupID == nil {
			parentlessCacheGroups = append(parentlessCacheGroups, cg)
			continue
		}
		if len(resp) > 0 {
			respCG := resp[0]
			_, _, err := TOSession.DeleteCacheGroupByID(*respCG.ID)
			if err != nil {
				t.Errorf("cannot DELETE CacheGroup by name: '%s' %v", *respCG.Name, err)
			}
			// Retrieve the CacheGroup to see if it got deleted
			cgs, _, err := TOSession.GetCacheGroupByName(*cg.Name, nil)
			if err != nil {
				t.Errorf("error deleting CacheGroup by name: %s", err.Error())
			}
			if len(cgs) > 0 {
				t.Errorf("expected CacheGroup name: %s to be deleted", *cg.Name)
			}
		}
	}

	// now delete the parentless cachegroups
	for _, cg := range parentlessCacheGroups {
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroupByName(*cg.Name, nil)
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
			cgs, _, err := TOSession.GetCacheGroupByName(*cg.Name, nil)
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

	resp, _, err := TOSession.GetCacheGroupByName(*cg.Name, nil)
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: %v - %v", *cg.Name, err)
	}
	cg = resp[0]

	if _, _, err = NoAuthTOSession.CreateCacheGroup(cg); err == nil {
		t.Error(fmt.Errorf(errFormat, "CreateCacheGroup"))
	}
	if _, _, err = NoAuthTOSession.GetCacheGroups(nil, nil); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheGroups"))
	}
	if _, _, err = NoAuthTOSession.GetCacheGroupByName(*cg.Name, nil); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheGroupByName"))
	}
	if _, _, err = NoAuthTOSession.GetCacheGroupByID(*cg.ID, nil); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheGroupByID"))
	}
	if _, _, err = NoAuthTOSession.UpdateCacheGroup(*cg.ID, cg, nil); err == nil {
		t.Error(fmt.Errorf(errFormat, "UpdateCacheGroup"))
	}
	if _, _, err = NoAuthTOSession.DeleteCacheGroupByID(*cg.ID); err == nil {
		t.Error(fmt.Errorf(errFormat, "DeleteCacheGroupByID"))
	}
}

func GetTestPaginationSupportCg(t *testing.T) {
	qparams := url.Values{}
	qparams.Set("orderby", "id")
	cachegroup, _, err := TOSession.GetCacheGroups(qparams, nil)
	if err != nil {
		t.Fatalf("cannot GET Cachegroup: %v", err)
	}

	qparams = url.Values{}
	qparams.Set("orderby", "id")
	qparams.Set("limit", "1")
	cachegroupWithLimit, _, err := TOSession.GetCacheGroups(qparams, nil)
	if !reflect.DeepEqual(cachegroup[:1], cachegroupWithLimit) {
		t.Error("expected GET Cachegroups with limit = 1 to return first result")
	}

	qparams = url.Values{}
	qparams.Set("orderby", "id")
	qparams.Set("limit", "1")
	qparams.Set("offset", "1")
	cachegroupsWithOffset, _, err := TOSession.GetCacheGroups(qparams, nil)
	if !reflect.DeepEqual(cachegroup[1:2], cachegroupsWithOffset) {
		t.Error("expected GET cachegroup with limit = 1, offset = 1 to return second result")
	}

	qparams = url.Values{}
	qparams.Set("orderby", "id")
	qparams.Set("limit", "1")
	qparams.Set("page", "2")
	cachegroupWithPage, _, err := TOSession.GetCacheGroups(qparams, nil)
	if !reflect.DeepEqual(cachegroup[1:2], cachegroupWithPage) {
		t.Error("expected GET cachegroup with limit = 1, page = 2 to return second result")
	}

	qparams = url.Values{}
	qparams.Set("limit", "-2")
	_, _, err = TOSession.GetCacheGroups(qparams, nil)
	if err == nil {
		t.Error("expected GET cachegroup to return an error when limit is not bigger than -1")
	} else if !strings.Contains(err.Error(), "must be bigger than -1") {
		t.Errorf("expected GET cachegroup to return an error for limit is not bigger than -1, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("offset", "0")
	_, _, err = TOSession.GetCacheGroups(qparams, nil)
	if err == nil {
		t.Error("expected GET cachegroup to return an error when offset is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET cachegroup to return an error for offset is not a positive integer, actual error: " + err.Error())
	}

	qparams = url.Values{}
	qparams.Set("limit", "1")
	qparams.Set("page", "0")
	_, _, err = TOSession.GetCacheGroups(qparams, nil)
	if err == nil {
		t.Error("expected GET cachegroup to return an error when page is not a positive integer")
	} else if !strings.Contains(err.Error(), "must be a positive integer") {
		t.Errorf("expected GET cachegroup to return an error for page is not a positive integer, actual error: " + err.Error())
	}
}
