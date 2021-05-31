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
	"strconv"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestCacheGroupParameters(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CacheGroupParameters}, func() {
		GetTestCacheGroupParameters(t)
		GetTestCacheGroupParametersIMS(t)
		CreateTestCacheGroupParametersMulAssignments(t)
		GetTestPaginationSupportCgParameters(t)
	})
}

func CreateTestCacheGroupParameters(t *testing.T) {
	if len(testData.CacheGroups) < 1 || len(testData.Parameters) < 1 {
		t.Fatal("Need at least one Cache Group and one Parameter to test associating Parameters to Cache Groups")
	}
	firstCacheGroup := testData.CacheGroups[0]
	if firstCacheGroup.Name == nil {
		t.Fatal("Found Cache Group with null or undefined name in test data")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstCacheGroup.Name)
	cacheGroupResp, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group '%s': %v - alerts: %+v", *firstCacheGroup.Name, err, cacheGroupResp.Alerts)
	}
	if len(cacheGroupResp.Response) != 1 {
		t.Fatalf("Expected exactly one Cache Group named '%s' to exist, but found %d", *firstCacheGroup.Name, len(cacheGroupResp.Response))
	}

	// Get Parameter to assign to Cache Group
	firstParameter := testData.Parameters[0]
	opts.QueryParameters.Set("name", firstParameter.Name)
	paramResp, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Errorf("cannot get Parameter '%s': %v - alerts: %+v", firstParameter.Name, err, paramResp.Alerts)
	}
	if len(paramResp.Response) < 1 {
		t.Fatalf("Expected at least one Parameter to exist with Name '%s'", firstParameter.Name)
	}

	// Assign Parameter to Cache Group
	cacheGroupID := cacheGroupResp.Response[0].ID
	if cacheGroupID == nil {
		t.Fatalf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *firstCacheGroup.Name)
	}
	parameterID := paramResp.Response[0].ID
	resp, _, err := TOSession.CreateCacheGroupParameter(*cacheGroupID, parameterID, client.RequestOptions{})
	if err != nil {
		t.Errorf("could not create cache group parameter: %v - alerts: %+v", err, resp.Alerts)
	}
	if resp.Response == nil {
		t.Fatal("Cache Group Parameter response should not be nil")
	}
	testData.CacheGroupParameterRequests = append(testData.CacheGroupParameterRequests, resp.Response...)
}

func CreateTestCacheGroupParametersMulAssignments(t *testing.T) {
	if len(testData.CacheGroups) < 3 || len(testData.Parameters) < 3 {
		t.Fatal("Need at least three Cache Group and three Parameter to test associating Parameters to Cache Groups")
	}
	firstCacheGroup := testData.CacheGroups[1]
	secondCacheGroup := testData.CacheGroups[2]
	if firstCacheGroup.Name == nil {
		t.Fatal("Found Cache Group1 with null or undefined name in test data")
	}
	if secondCacheGroup.Name == nil {
		t.Fatal("Found Cache Group2 with null or undefined name in test data")
	}

	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("name", *firstCacheGroup.Name)
	cacheGroupResp1, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group '%s': %v - alerts: %+v", *firstCacheGroup.Name, err, cacheGroupResp1.Alerts)
	}
	if len(cacheGroupResp1.Response) != 1 {
		t.Fatalf("Expected exactly one Cache Group named '%s' to exist, but found %d", *firstCacheGroup.Name, len(cacheGroupResp1.Response))
	}
	opts.QueryParameters.Set("name", *secondCacheGroup.Name)
	cacheGroupResp2, _, err := TOSession.GetCacheGroups(opts)
	if err != nil {
		t.Fatalf("cannot get Cache Group '%s': %v - alerts: %+v", *secondCacheGroup.Name, err, cacheGroupResp2.Alerts)
	}
	if len(cacheGroupResp2.Response) != 1 {
		t.Fatalf("Expected exactly one Cache Group named '%s' to exist, but found %d", *secondCacheGroup.Name, len(cacheGroupResp2.Response))
	}

	// Get Parameter to assign to Cache Group
	firstParameter := testData.Parameters[1]
	secondParameter := testData.Parameters[2]
	opts.QueryParameters.Set("name", firstParameter.Name)
	paramResp1, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Errorf("cannot get Parameter '%s': %v - alerts: %+v", firstParameter.Name, err, paramResp1.Alerts)
	}
	if len(paramResp1.Response) != 1 {
		t.Fatalf("Expected exactly one Parameter named '%s' to exist, but found %d", firstParameter.Name, len(paramResp1.Response))
	}
	opts.QueryParameters.Set("name", secondParameter.Name)
	paramResp2, _, err := TOSession.GetParameters(opts)
	if err != nil {
		t.Errorf("cannot get Parameter '%s': %v - alerts: %+v", secondParameter.Name, err, paramResp2.Alerts)
	}
	if len(paramResp2.Response) < 1 {
		t.Fatalf("Expected exactly one Parameter named '%s' to exist, but found %d", secondParameter.Name, len(paramResp2.Response))
	}

	// Assign Parameter to Cache Group
	cacheGroupID1 := cacheGroupResp1.Response[0].ID
	cacheGroupID2 := cacheGroupResp2.Response[0].ID
	if cacheGroupID1 == nil {
		t.Fatalf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *firstCacheGroup.Name)
	}
	if cacheGroupID2 == nil {
		t.Fatalf("Traffic Ops returned Cache Group '%s' with null or undefined ID", *secondCacheGroup.Name)
	}

	parameterID1 := paramResp1.Response[0].ID
	parameterID2 := paramResp2.Response[0].ID
	pp := tc.CacheGroupParameterCreationRequest{
		CacheGroupID: *cacheGroupID1,
		ParameterID:  parameterID1,
	}
	pp2 := tc.CacheGroupParameterCreationRequest{
		CacheGroupID: *cacheGroupID2,
		ParameterID:  parameterID2,
	}

	ppSlice := []tc.CacheGroupParameterCreationRequest{
		pp,
		pp2,
	}
	resp, _, err := TOSession.CreateMultipleCacheGroupParameter(ppSlice, client.RequestOptions{})
	if err != nil {
		t.Errorf("could not create cache group parameter: %v - alerts: %+v", err, resp.Alerts)
	}
	if resp.Response == nil {
		t.Fatal("Cache Group Parameter response should not be nil")
	}
	testData.CacheGroupParameterRequests = append(testData.CacheGroupParameterRequests, resp.Response...)
}

func GetTestCacheGroupParameters(t *testing.T) {
	for _, cgp := range testData.CacheGroupParameterRequests {
		resp, _, err := TOSession.GetCacheGroupParameters(cgp.CacheGroupID, client.RequestOptions{})
		if err != nil {
			t.Errorf("cannot get Parameter by Cache Group #%d: %v - alerts: %+v", cgp.CacheGroupID, err, resp.Alerts)
		}
		if len(resp.Response) < 1 {
			t.Errorf("Expected Cache Group #%d to have at least one associated Parameter, but found none", cgp.CacheGroupID)
		}
	}
}

func GetTestCacheGroupParametersIMS(t *testing.T) {
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)

	opts := client.NewRequestOptions()
	opts.Header.Set(rfc.IfModifiedSince, time)

	for _, cgp := range testData.CacheGroupParameterRequests {
		resp, reqInf, err := TOSession.GetCacheGroupParameters(cgp.CacheGroupID, opts)
		if err != nil {
			t.Errorf("Expected no error fetching Parameters for a Cache Group, but got %v - alerts: %+v", err, resp.Alerts)
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Errorf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func DeleteTestCacheGroupParameters(t *testing.T) {
	for _, cgp := range testData.CacheGroupParameterRequests {
		DeleteTestCacheGroupParameter(t, cgp)
	}
}

func DeleteTestCacheGroupParameter(t *testing.T, cgp tc.CacheGroupParameterRequest) {

	delResp, _, err := TOSession.DeleteCacheGroupParameter(cgp.CacheGroupID, cgp.ParameterID, client.RequestOptions{})
	if err != nil {
		t.Fatalf("cannot delete Parameter by Cache Group ID: %v - alerts: %+v", err, delResp)
	}

	// Retrieve the Cache Group Parameter to see if it got deleted
	opts := client.NewRequestOptions()
	opts.QueryParameters.Add("parameterId", strconv.Itoa(cgp.ParameterID))

	parameters, _, err := TOSession.GetCacheGroupParameters(cgp.CacheGroupID, opts)
	if err != nil {
		t.Errorf("error getting Parameters by Cache Group ID after dissociation: %s - alerts: %+v", err, parameters.Alerts)
	}
	if parameters.Response == nil {
		t.Fatal("Cache Group Parameters response should not be nil")
	}
	if len(parameters.Response) > 0 {
		t.Errorf("expected Parameter: %d to be to be disassociated from Cache Group: %d", cgp.ParameterID, cgp.CacheGroupID)
	}

	// Attempt to delete it again and it should return an error now
	_, _, err = TOSession.DeleteCacheGroupParameter(cgp.CacheGroupID, cgp.ParameterID, client.RequestOptions{})
	if err == nil {
		t.Error("expected error when deleting unassociated cache group parameter")
	}

	// Attempt to delete using a non existing cache group
	_, _, err = TOSession.DeleteCacheGroupParameter(-1, cgp.ParameterID, client.RequestOptions{})
	if err == nil {
		t.Error("expected error when deleting cache group parameter with non existing cache group")
	}

	// Attempt to delete using a non existing parameter
	_, _, err = TOSession.DeleteCacheGroupParameter(cgp.CacheGroupID, -1, client.RequestOptions{})
	if err == nil {
		t.Error("expected error when deleting cache group parameter with non existing parameter")
	}
}

func GetTestPaginationSupportCgParameters(t *testing.T) {
	opts := client.NewRequestOptions()
	opts.QueryParameters.Set("orderby", "id")
	resp, _, err := TOSession.GetAllCacheGroupParameters(opts)
	if err != nil {
		t.Fatalf("cannot get cachegroup parameters: %v - alerts: %+v", err, resp.Alerts)
	}
	cachegroupParameters := resp.Response
	if len(cachegroupParameters.CacheGroupParameters) < 3 {
		t.Fatalf("Need at least 3 cachegroup parameters in Traffic Ops to test pagination support, found: %d", len(cachegroupParameters.CacheGroupParameters))
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	respWithLimit, _, err := TOSession.GetAllCacheGroupParameters(opts)
	if err != nil {
		t.Fatalf("cannot get cachegroup parameters with limits: %v - alerts: %+v", err, respWithLimit.Alerts)
	}
	cachegroupParametersWithLimit := respWithLimit.Response
	if !reflect.DeepEqual(cachegroupParameters.CacheGroupParameters[:1], cachegroupParametersWithLimit.CacheGroupParameters) {
		t.Error("expected GET cachegroup parameters with limit = 1 to return first result")
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "1")
	respWithOffset, _, err := TOSession.GetAllCacheGroupParameters(opts)
	if err != nil {
		t.Fatalf("cannot get cachegroup parameters with offset: %v - alerts: %+v", err, respWithOffset.Alerts)
	}
	cachegroupParametersWithOffset := respWithOffset.Response
	if !reflect.DeepEqual(cachegroupParameters.CacheGroupParameters[1:2], cachegroupParametersWithOffset.CacheGroupParameters) {
		t.Error("expected GET cachegroup parameters with limit = 1, offset = 1 to return second result")
	}

	opts.QueryParameters.Set("orderby", "id")
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "2")
	respWithPage, _, err := TOSession.GetAllCacheGroupParameters(opts)
	if err != nil {
		t.Fatalf("cannot get cachegroup parameters with page: %v - alerts: %+v", err, respWithPage.Alerts)
	}
	cachegroupParametersWithPage := respWithPage.Response
	if !reflect.DeepEqual(cachegroupParameters.CacheGroupParameters[1:2], cachegroupParametersWithPage.CacheGroupParameters) {
		t.Error("expected GET cachegroup parameters with limit = 1, page = 2 to return second result")
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "-2")
	resp, _, err = TOSession.GetAllCacheGroupParameters(opts)
	if !alertsHaveError(resp.Alerts.Alerts, "limit parameter must be bigger than -1") {
		t.Errorf("expected GET cachegroup parameters to return an error for limit is not bigger than -1, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("offset", "0")
	resp, _, err = TOSession.GetAllCacheGroupParameters(opts)
	if !alertsHaveError(resp.Alerts.Alerts, "offset parameter must be a positive integer") {
		t.Errorf("expected GET cachegroup parameters to return an error for offset is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}

	opts.QueryParameters = url.Values{}
	opts.QueryParameters.Set("limit", "1")
	opts.QueryParameters.Set("page", "0")
	resp, _, err = TOSession.GetAllCacheGroupParameters(opts)
	if !alertsHaveError(resp.Alerts.Alerts, "page parameter must be a positive integer") {
		t.Errorf("expected GET cachegroup parameters to return an error for page is not a positive integer, actual error: %v - alerts: %+v", err, resp.Alerts)
	}
}
