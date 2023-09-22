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

package v3

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestCacheGroupParameters(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CacheGroupParameters}, func() {
		GetTestCacheGroupParameters(t)
		GetTestCacheGroupParametersIMS(t)
	})
}

func CreateTestCacheGroupParameters(t *testing.T) {
	// Get Cache Group to assign parameter to
	firstCacheGroup := testData.CacheGroups[0]
	cacheGroupResp, _, err := TOSession.GetCacheGroupNullableByName(*firstCacheGroup.Name)
	if err != nil {
		t.Errorf("cannot GET Cache Group by name: %v - %v", firstCacheGroup.Name, err)
	}
	if cacheGroupResp == nil {
		t.Fatal("Cache Groups response should not be nil")
	}

	// Get Parameter to assign to Cache Group
	firstParameter := testData.Parameters[0]
	paramResp, _, err := TOSession.GetParameterByName(firstParameter.Name)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v", firstParameter.Name, err)
	}
	if paramResp == nil {
		t.Fatal("Parameter response should not be nil")
	}

	// Assign Parameter to Cache Group
	cacheGroupID := cacheGroupResp[0].ID
	parameterID := paramResp[0].ID
	resp, _, err := TOSession.CreateCacheGroupParameter(*cacheGroupID, parameterID)
	if err != nil {
		t.Errorf("could not CREATE cache group parameter: %v", err)
	}
	if resp == nil {
		t.Fatal("Cache Group Parameter response should not be nil")
	}
	testData.CacheGroupParameterRequests = append(testData.CacheGroupParameterRequests, resp.Response...)
}

func GetTestCacheGroupParameters(t *testing.T) {
	for _, cgp := range testData.CacheGroupParameterRequests {
		resp, _, err := TOSession.GetCacheGroupParameters(cgp.CacheGroupID)
		if err != nil {
			t.Errorf("cannot GET Parameter by cache group: %v - %v", err, resp)
		}
		if resp == nil {
			t.Fatal("Cache Group Parameters response should not be nil")
		}
	}
}

func GetTestCacheGroupParametersIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	for _, cgp := range testData.CacheGroupParameterRequests {
		_, reqInf, err := TOSession.GetCacheGroupParametersWithHdr(cgp.CacheGroupID, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func DeleteTestCacheGroupParameters(t *testing.T) {
	for _, cgp := range testData.CacheGroupParameterRequests {
		DeleteTestCacheGroupParameter(t, cgp)
	}
}

func DeleteTestCacheGroupParameter(t *testing.T, cgp tc.CacheGroupParameterRequest) {

	delResp, _, err := TOSession.DeleteCacheGroupParameter(cgp.CacheGroupID, cgp.ParameterID)
	if err != nil {
		t.Fatalf("cannot DELETE Parameter by cache group: %v - %v", err, delResp)
	}

	// Retrieve the Cache Group Parameter to see if it got deleted
	queryParams := fmt.Sprintf("?parameterId=%d", cgp.ParameterID)

	parameters, _, err := TOSession.GetCacheGroupParametersByQueryParams(cgp.CacheGroupID, queryParams)
	if err != nil {
		t.Errorf("error deleting Parameter name: %s", err.Error())
	}
	if parameters == nil {
		t.Fatal("Cache Group Parameters response should not be nil")
	}
	if len(parameters) > 0 {
		t.Errorf("expected Parameter: %d to be to be disassociated from Cache Group: %d", cgp.ParameterID, cgp.CacheGroupID)
	}

	// Attempt to delete it again and it should return an error now
	_, _, err = TOSession.DeleteCacheGroupParameter(cgp.CacheGroupID, cgp.ParameterID)
	if err == nil {
		t.Error("expected error when deleting unassociated cache group parameter")
	}

	// Attempt to delete using a non existing cache group
	_, _, err = TOSession.DeleteCacheGroupParameter(-1, cgp.ParameterID)
	if err == nil {
		t.Error("expected error when deleting cache group parameter with non existing cache group")
	}

	// Attempt to delete using a non existing parameter
	_, _, err = TOSession.DeleteCacheGroupParameter(cgp.CacheGroupID, -1)
	if err == nil {
		t.Error("expected error when deleting cache group parameter with non existing parameter")
	}
}
