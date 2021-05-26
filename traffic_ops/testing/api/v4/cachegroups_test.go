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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
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
		GetTestCacheGroupsByInvalidId(t)
		GetTestCacheGroupsByInvalidType(t)
		GetTestCacheGroupsByType(t)
		DeleteTestCacheGroupsByInvalidId(t)
	})
}

func UpdateTestCacheGroupsWithHeaders(t *testing.T, h http.Header) {
	if len(testData.CacheGroups) < 1 {
		t.Fatal("Need at least one Cache Group to test updating Cache Groups")
	}
	firstCG := testData.CacheGroups[0]
	if firstCG.Name == nil {
		t.Fatal("Found Cache Group with null or undefined name in testing data")
	}

	opts := client.RequestOptions{
		Header:          h,
		QueryParameters: url.Values{},
	}
	opts.QueryParameters.Set("name", *firstCG.Name)

	resp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group '%s': %v - alerts: %+v", *firstCG.Name, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected exactly one Cache Group to exist with name '%s', but got: %d", *firstCG.Name, len(resp.Response))
	}

	cg := resp.Response[0]
	if cg.TypeID == nil {
		t.Fatalf("Traffic Ops returned Cache Group '%s' with null or undefined typeId", *firstCG.Name)
	}
	if cg.ID == nil {
		t.Fatalf("Traffic Ops returned Cache Group '%s' with null or undefined id", *firstCG.Name)
	}
	expectedShortName := "blah"
	cg.ShortName = &expectedShortName

	// fix the type id for test
	typeOpts := client.NewRequestOptions()
	typeOpts.QueryParameters.Set("id", strconv.Itoa(*cg.TypeID))
	typeResp, _, err := TOSession.GetTypes(typeOpts)
	if err != nil {
		t.Fatalf("Failed to fetch Type #%d: %v - alerts: %+v", *cg.TypeID, err, typeResp.Alerts)
	}
	if len(typeResp.Response) != 1 {
		t.Fatalf("Expected exactly one Type to exist with ID %d, but got: %d", *cg.TypeID, len(typeResp.Response))
	}

	cg.TypeID = &typeResp.Response[0].ID
	_, reqInf, err := TOSession.UpdateCacheGroup(*cg.ID, cg, opts)
	if err == nil {
		t.Errorf("Expected an error showing Precondition Failed, got none")
	}
	if reqInf.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("Expected status code 412, got %v", reqInf.StatusCode)
	}
}

func GetTestCacheGroupsAfterChangeIMS(t *testing.T, header http.Header) {
	opts := client.RequestOptions{
		Header: header,
	}
	resp, reqInf, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 status code, got %v", reqInf.StatusCode)
	}

	currentTime := time.Now().UTC()
	currentTime = currentTime.Add(1 * time.Second)
	timeStr := currentTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, timeStr)
	resp, reqInf, err = TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Errorf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func GetTestCacheGroupsByShortNameIMS(t *testing.T) {
	opts := client.NewRequestOptions()
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, cg := range testData.CacheGroups {
		if cg.ShortName == nil {
			t.Error("found Cache Group with null or undefined 'short name' in test data")
			continue
		}
		opts.QueryParameters.Set("shortName", *cg.ShortName)
		resp, reqInf, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCacheGroupsByNameIMS(t *testing.T) {
	opts := client.NewRequestOptions()
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, time)
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("found Cache Group with null or undefined name in test data")
			continue
		}
		opts.QueryParameters.Set("name", *cg.Name)
		resp, reqInf, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCacheGroupsIMS(t *testing.T) {
	opts := client.NewRequestOptions()
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	opts.Header.Set(rfc.IfModifiedSince, time)
	resp, reqInf, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Errorf("Expected no error, but got %v - alerts: %+v", err, resp.Alerts)
	}
	if reqInf.StatusCode != http.StatusNotModified {
		t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
	}
}

func CreateTestCacheGroups(t *testing.T) {

	for _, cg := range testData.CacheGroups {

		resp, _, err := TOSession.CreateCacheGroup(cg, client.RequestOptions{})
		if err != nil {
			t.Errorf("could not create Cache Group: %v - alerts: %+v", err, resp.Alerts)
			continue
		}

		// Testing 'join' fields during create
		if cg.ParentName != nil && resp.Response.ParentName == nil {
			t.Error("Parent cachegroup is null in response when it should have a value")
		}
		if cg.SecondaryParentName != nil && resp.Response.SecondaryParentName == nil {
			t.Error("Secondary parent cachegroup is null in response when it should have a value")
		}
		if cg.Type != nil && resp.Response.Type == nil {
			t.Error("Type is null in response when it should have a value")
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
	resp, _, err := TOSession.GetCacheGroups(client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot get Cache Groups: %v - alerts: %+v", err, resp.Alerts)
	}
	expectedCachegroups := make(map[string]struct{}, len(testData.CacheGroups))
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("Found Cache Group in testing data with null or undefined name")
			continue
		}
		expectedCachegroups[*cg.Name] = struct{}{}
	}
	foundCachegroups := make(map[string]struct{}, len(expectedCachegroups))
	for _, cg := range resp.Response {
		if cg.Name == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined name")
			continue
		}
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
	opts := client.NewRequestOptions()
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("found Cache Group with null or undefined name in test data")
			continue
		}
		opts.QueryParameters.Set("name", *cg.Name)
		resp, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("cannot get CacheGroup by name '%s': %v - alerts: %+v", *cg.Name, err, resp)
		}
		if len(resp.Response) < 1 {
			t.Errorf("Expected exactly one Cache Group with name '%s', but got none", *cg.Name)
			continue
		}
		if len(resp.Response) > 1 {
			t.Errorf("Expected exactly one Cache Group with name '%s', but got %d", *cg.Name, len(resp.Response))
			t.Log("Testing will proceed using the first Cache Group found in the response")
		}
		respCG := resp.Response[0]
		if respCG.Name == nil {
			t.Errorf("Cache Group as returned by Traffic Ops had null or undefined name (should be '%s')", *cg.Name)
		}
		if *respCG.Name != *cg.Name {
			t.Errorf("name expected: %s, actual: %s", *cg.Name, *respCG.Name)
		}
	}
}

func GetTestCacheGroupsByShortName(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, cg := range testData.CacheGroups {
		if cg.ShortName == nil {
			t.Error("found Cache Group with null or undefined 'short name' in test data")
			continue
		}
		opts.QueryParameters.Set("shortName", *cg.ShortName)
		resp, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("cannot get Cache Group by 'short name': %v - alerts: %+v", err, resp.Alerts)
		}
		if len(resp.Response) > 1 {
			t.Errorf("Expected exactly one Cache Group with name '%s', but got %d", *cg.ShortName, len(resp.Response))
			t.Log("Testing will proceed using the first Cache Group found in the response")
		}
		respCG := resp.Response[0]
		if respCG.ShortName == nil {
			t.Errorf("Cache Group as returned by Traffic Ops had null or undefined name (should be '%s')", *cg.ShortName)
		}
		if *respCG.ShortName != *cg.ShortName {
			t.Errorf("short name expected: %s, actual: %s", *cg.ShortName, *respCG.ShortName)
		}
	}
}

func GetTestCacheGroupsByTopology(t *testing.T) {
	opts := client.NewRequestOptions()
	for _, top := range testData.Topologies {
		opts.QueryParameters.Set("topology", top.Name)
		resp, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("cannot get Cache Groups by Topology '%s': %v - alerts: %v", top.Name, err, resp.Alerts)
		}
		expectedCGs := topologyCachegroups(top)
		for _, cg := range resp.Response {
			if cg.Name == nil {
				t.Errorf("Traffic Ops returned a Cache Group in Topology '%s' with null or undefined name", top.Name)
				continue
			}
			if _, exists := expectedCGs[*cg.Name]; !exists {
				t.Errorf("GET cachegroups by topology - expected one of: %v, actual: %s", expectedCGs, *cg.Name)
			}
		}
	}
}

func topologyCachegroups(top tc.Topology) map[string]struct{} {
	res := make(map[string]struct{}, len(top.Nodes))
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
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstCG.Name)
	resp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group by name '%s': %v - alerts: %+v", *firstCG.Name, err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg := resp.Response[0]
	if cg.TypeID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined typeId")
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined id")
	}
	expectedShortName := "blah"
	cg.ShortName = &expectedShortName

	// fix the type id for test
	typeOpts := client.NewRequestOptions()
	typeOpts.QueryParameters.Set("id", strconv.Itoa(*cg.TypeID))
	typeResp, _, err := TOSession.GetTypes(typeOpts)
	if err != nil {
		t.Errorf("could not lookup an ID for the Type of this Cache Group: %v - alerts: %+v", err, typeResp.Alerts)
	}
	if len(typeResp.Response) == 0 {
		t.Fatal("got an empty response for types")
	}
	cg.TypeID = &typeResp.Response[0].ID
	opts.QueryParameters = url.Values{}
	updResp, _, err := TOSession.UpdateCacheGroup(*cg.ID, cg, opts)
	if err != nil {
		t.Errorf("cannot update Cache Group: %v - alerts: %+v", err, updResp.Alerts)
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
	opts.QueryParameters.Set("id", strconv.Itoa(*cg.ID))
	resp, _, err = TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot GET CacheGroup by name: '%s', %v", *firstCG.Name, err)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
	if cg.ShortName == nil {
		t.Fatal("Traffic Ops returned a Cache Group with nullor undefined short name")
	}
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
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by id: %v - alerts: %+v", err, updResp.Alerts)
	}

	resp, _, err = TOSession.GetCacheGroups(opts)

	if err != nil {
		t.Fatalf("cannot GET CacheGroup by id: '%d', %v", *cg.ID, err)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
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
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot UPDATE CacheGroup by ID %d: %v - alerts: %+v", *cg.ID, err, updResp.Alerts)
	}

	opts.QueryParameters.Set("id", strconv.Itoa(*cg.ID))
	resp, _, err = TOSession.GetCacheGroups(opts)

	if err != nil {
		t.Errorf("cannot GET CacheGroup by ID %d: %v - alerts: %+v", *cg.ID, err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
	if cg.LocalizationMethods == nil {
		t.Fatal("Traffic Ops returned Cache Group with null or undefined localizationMethods")
	}
	if !reflect.DeepEqual(expectedMethods, *cg.LocalizationMethods) {
		t.Errorf("failed to update localizationMethods (expected = %v, actual = %v)", expectedMethods, *cg.LocalizationMethods)
	}

	// test cachegroup fallbacks

	// Retrieve the CacheGroup to check CacheGroup name got updated
	firstEdgeCGName := "cachegroup1"
	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("name", firstEdgeCGName)
	resp, _, err = TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Errorf("cannot get Cache Group by name '%s': %v - alerts: %+v", firstEdgeCGName, err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
	if cg.Name == nil {
		t.Error("Cache Group returned from Traffic Ops had null or undefined name")
	} else if *cg.Name != firstEdgeCGName {
		t.Errorf("results do not match actual: %s, expected: %s", *cg.Name, firstEdgeCGName)
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned from Traffic Ops had null or undefined id")
	}

	// Test adding fallbacks when previously nil
	expectedFallbacks := []string{"fallback1", "fallback2"}
	cg.Fallbacks = &expectedFallbacks
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Cache Group by ID: %v - alerts: %+v", err, updResp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("id", strconv.Itoa(*cg.ID))
	resp, _, err = TOSession.GetCacheGroups(opts)

	if err != nil {
		t.Errorf("cannot GET CacheGroup by ID #%d: %v - alerts: %+v", *cg.ID, err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
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
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Cache Group by ID: %v - %+v)", err, updResp.Alerts)
	}

	opts.QueryParameters.Set("id", strconv.Itoa(*cg.ID))
	resp, _, err = TOSession.GetCacheGroups(opts)

	if err != nil {
		t.Errorf("cannot get Cache Group by id #%d: %v - alerts: %+v", *cg.ID, err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
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
	updResp, _, err = TOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	if err != nil {
		t.Errorf("cannot update Cache Group by ID: %v - alerts: %+v", err, updResp.Alerts)
	}

	opts.QueryParameters.Set("id", strconv.Itoa(*cg.ID))
	resp, _, err = TOSession.GetCacheGroups(opts)

	if err != nil {
		t.Errorf("cannot get Cache Group by ID: '%d', %v", *cg.ID, err)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
	if cg.Fallbacks == nil {
		t.Error("Cache Group returned by Traffic Ops had null or undefined fallbacks")
	} else if !reflect.DeepEqual(expectedFallbacks, *cg.Fallbacks) {
		t.Errorf("failed to update fallbacks (expected = %v, actual = %v)", expectedFallbacks, *cg.Fallbacks)
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined ID")
	}

	const topologyEdgeCGName = "topology-edge-cg-01"
	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("name", topologyEdgeCGName)
	resp, _, err = TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group by name '%s': %v - alerts: %+v", topologyEdgeCGName, err, resp.Alerts)
	}
	if len(resp.Response) == 0 {
		t.Fatal("got an empty response for cachegroups")
	}
	cg = resp.Response[0]
	if cg.TypeID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined typeId")
	}
	if cg.ID == nil {
		t.Fatal("Cache Group returned by Traffic Ops had null or undefined ID")
	}

	var cacheGroupEdgeType, cacheGroupMidType tc.Type
	types, _, err := TOSession.GetTypes(client.RequestOptions{})
	if err != nil {
		t.Fatalf("unable to get Types: %v - alerts: %+v", err, types.Alerts)
	}
	for _, typeObject := range types.Response {
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
	updResp, reqInfo, err := TOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{})
	if err == nil {
		t.Fatalf("expected an error when updating the type of cache group %s because it is assigned to a topology, actual error was nil", *cg.Name)
	}
	if reqInfo.StatusCode != http.StatusBadRequest {
		msg := "expected to receive status code %d but received status code %d - error was: %v - alerts: %+v"
		t.Fatalf(msg, http.StatusBadRequest, reqInfo.StatusCode, err, updResp.Alerts)
	}
}

func DeleteTestCacheGroups(t *testing.T) {
	var parentlessCacheGroups []tc.CacheGroupNullable
	opts := client.NewRequestOptions()

	// delete the edge caches.
	for _, cg := range testData.CacheGroups {
		if cg.Name == nil {
			t.Error("Found a Cache Group with null or undefined name")
			continue
		}
		// Retrieve the CacheGroup by name so we can get the id for the Update
		opts.QueryParameters.Set("name", *cg.Name)
		resp, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("cannot GET CacheGroup by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)
		}
		if len(resp.Response) < 1 {
			t.Errorf("Could not find test data Cache Group '%s' in Traffic Ops", *cg.Name)
			continue
		}
		cg = resp.Response[0]

		// Cachegroups that are parents (usually mids but sometimes edges)
		// need to be deleted only after the children cachegroups are deleted.
		if cg.ParentCachegroupID == nil && cg.SecondaryParentCachegroupID == nil {
			parentlessCacheGroups = append(parentlessCacheGroups, cg)
			continue
		}

		// TODO: Typo here? cg is already reassigned to resp.Response[0] - is respCG supposed to be different?
		respCG := resp.Response[0]
		if respCG.ID == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined ID")
			continue
		}
		if respCG.Name == nil {
			t.Error("Traffic Ops returned a Cache Group with null or undefined name")
			continue
		}
		alerts, _, err := TOSession.DeleteCacheGroup(*respCG.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Cache Group: %v - alerts: %+v", err, alerts)
		}
		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *respCG.Name)
		cgs, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("error deleting Cache Group by name: %v - alerts: %+v", err, cgs.Alerts)
		}
		if len(cgs.Response) > 0 {
			t.Errorf("expected CacheGroup name: %s to be deleted", *cg.Name)
		}
	}

	opts = client.NewRequestOptions()
	// now delete the parentless cachegroups
	for _, cg := range parentlessCacheGroups {
		// nil check for cg.Name occurs prior to insertion into parentlessCacheGroups
		opts.QueryParameters.Set("name", *cg.Name)
		// Retrieve the CacheGroup by name so we can get the id for the Update
		resp, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("cannot get Cache Group by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)
		}
		if len(resp.Response) < 1 {
			t.Errorf("Cache Group '%s' somehow stopped existing since the last time we ask Traffic Ops about it", *cg.Name)
			continue
		}

		respCG := resp.Response[0]
		if respCG.ID == nil {
			t.Errorf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *cg.Name)
			continue
		}
		delResp, _, err := TOSession.DeleteCacheGroup(*respCG.ID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot delete Cache Group '%s': %v - alerts: %+v", *respCG.Name, err, delResp.Alerts)
		}

		// Retrieve the CacheGroup to see if it got deleted
		opts.QueryParameters.Set("name", *cg.Name)
		cgs, _, err := TOSession.GetCacheGroups(opts)
		if err != nil {
			t.Errorf("error attempting to fetch Cache Group '%s' after deletion: %v - alerts: %+v", *cg.Name, err, cgs.Alerts)
		}
		if len(cgs.Response) > 0 {
			t.Errorf("expected Cache Group '%s' to be deleted", *cg.Name)
		}
	}
}

func CheckCacheGroupsAuthentication(t *testing.T) {
	if len(testData.CacheGroups) < 1 {
		t.Fatalf("Need at least one Cache Group to test Cache Group API authentication")
	}
	cg := testData.CacheGroups[0]
	if cg.Name == nil {
		t.Fatal("Cache Group selected from testing data had null or undefined name")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *cg.Name)
	resp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Errorf("cannot get Cache Group by name '%s': %v - alerts: %+v", *cg.Name, err, resp.Alerts)
	}
	if len(resp.Response) != 1 {
		t.Fatalf("Expected exactly one Cache Group with name '%s', but found: %d", *cg.Name, len(resp.Response))
	}
	cg = resp.Response[0]

	const errFormat = "expected error from %s when unauthenticated"
	if _, _, err = NoAuthTOSession.CreateCacheGroup(cg, client.RequestOptions{}); err == nil {
		t.Error(fmt.Errorf(errFormat, "CreateCacheGroup"))
	}
	if _, _, err = NoAuthTOSession.GetCacheGroups(client.RequestOptions{}); err == nil {
		t.Error(fmt.Errorf(errFormat, "GetCacheGroups"))
	}
	if cg.Name == nil {
		t.Error("Traffic Ops returned a Cache Group with a null or undefined name")
	} else {
		opts.QueryParameters.Set("name", *cg.Name)
		_, _, err = NoAuthTOSession.GetCacheGroups(opts)
		if err == nil {
			t.Error(fmt.Errorf(errFormat, "GetCacheGroups filtered by Name"))
		}
	}
	opts.QueryParameters = url.Values{}
	if cg.ID == nil {
		t.Error("Traffic Ops returned a Cache Group with a null or undefined name")
	} else {
		opts.QueryParameters.Set("id", strconv.Itoa(*cg.ID))
		_, _, err = NoAuthTOSession.GetCacheGroups(opts)
		if err == nil {
			t.Error(fmt.Errorf(errFormat, "GetCacheGroups filtered by ID"))
		}
		if _, _, err = NoAuthTOSession.UpdateCacheGroup(*cg.ID, cg, client.RequestOptions{}); err == nil {
			t.Error(fmt.Errorf(errFormat, "UpdateCacheGroup"))
		}
		if _, _, err = NoAuthTOSession.DeleteCacheGroup(*cg.ID, client.RequestOptions{}); err == nil {
			t.Error(fmt.Errorf(errFormat, "DeleteCacheGroupByID"))
		}
	}
}

func GetTestPaginationSupportCg(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Groups: %v - alerts: %+v", err, resp.Alerts)
	}
	cachegroup := resp.Response
	if len(cachegroup) < 3 {
		t.Fatalf("Need at least 3 Cache Groups in Traffic Ops to test pagination support, found: %d", len(cachegroup))
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	cachegroupWithLimit, _, err := TOSession.GetCacheGroups(opts)

	if !reflect.DeepEqual(cachegroup[:1], cachegroupWithLimit.Response) {
		t.Error("expected GET Cachegroups with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "1")
	cachegroupsWithOffset, _, err := TOSession.GetCacheGroups(opts)
	if !reflect.DeepEqual(cachegroup[1:2], cachegroupsWithOffset.Response) {
		t.Error("expected GET cachegroup with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "2")
	cachegroupWithPage, _, err := TOSession.GetCacheGroups(opts)
	if !reflect.DeepEqual(cachegroup[1:2], cachegroupWithPage.Response) {
		t.Error("expected GET cachegroup with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetCacheGroups(opts)
	if err == nil {
		t.Error("expected GET cachegroup to return an error when limit is not bigger than -1")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be bigger than -1") {
		t.Errorf("expected GET cachegroup to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetCacheGroups(opts)
	if err == nil {
		t.Error("expected GET cachegroup to return an error when offset is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET cachegroup to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetCacheGroups(opts)
	if err == nil {
		t.Error("expected GET cachegroup to return an error when page is not a positive integer")
	} else if !alertsHaveError(resp.Alerts.Alerts, "must be a positive integer") {
		t.Errorf("expected GET cachegroup to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}

func GetTestCacheGroupsByInvalidId(t *testing.T) {
	opts := client.NewRequestOptions()
	// Retrieve the CacheGroup to check CacheGroup name got updated
	opts.QueryParameters.Set("id", "10000")
	resp, _, _ := TOSession.GetCacheGroups(opts)
	if len(resp.Response) > 0 {
		t.Errorf("Expected 0 response, but got many %v", resp)
	}
}

func GetTestCacheGroupsByInvalidType(t *testing.T) {
	opts := client.NewRequestOptions()
	// Retrieve the CacheGroup to check CacheGroup name got updated
	opts.QueryParameters.Set("type", "10000")
	resp, _, _ := TOSession.GetCacheGroups(opts)
	if len(resp.Response) > 0 {
		t.Errorf("Expected 0 response, but got many %v", resp)
	}
}

func GetTestCacheGroupsByType(t *testing.T) {
	if len(testData.CacheGroups) < 1 {
		t.Fatal("Need at least one Cache Group to test updating Cache Groups")
	}
	firstCG := testData.CacheGroups[0]

	if firstCG.Name == nil {
		t.Fatal("Found Cache Group with null or undefined name in testing data")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstCG.Name)

	resp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group '%s': %v - alerts: %+v", *firstCG.Name, err, resp.Alerts)
	}
	if len(resp.Response) < 1 {
		t.Fatalf("Expected exactly one Cache Group to exist with name '%s', but got: %d", *firstCG.Name, len(resp.Response))
	}

	cg := resp.Response[0]
	if cg.TypeID == nil {
		t.Fatalf("Traffic Ops returned Cache Group '%s' with null or undefined typeId", *firstCG.Name)
	}

	opts = client.NewRequestOptions()
	opts.QueryParameters.Set("type", strconv.Itoa(*cg.TypeID))
	resp, _, _ = TOSession.GetCacheGroups(opts)
	if len(resp.Response) < 1 {
		t.Fatalf("Expected atleast one Cache Group by type ID '%d', but got: %d", *cg.TypeID, len(resp.Response))
	}
}

func DeleteTestCacheGroupsByInvalidId(t *testing.T) {

	alerts, reqInf, err := TOSession.DeleteCacheGroup(111111, client.RequestOptions{})
	if err == nil {
		t.Errorf("Expected no cachegroup with that id found - but got alerts: %+v", alerts)
	}
	if reqInf.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %v", reqInf.StatusCode)
	}
}
