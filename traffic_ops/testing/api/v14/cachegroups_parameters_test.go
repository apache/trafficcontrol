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

package v14

import (
	"fmt"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestCacheGroupParameters(t *testing.T) {
	WithObjs(t, []TCObj{Types, Parameters, CacheGroups, CacheGroupParameters}, func() {
		GetTestCacheGroupParameters(t)
	})
}

func CreateTestCacheGroupParameters(t *testing.T) {
	firstCacheGroup := testData.CacheGroups[0]
	cacheGroupResp, _, err := TOSession.GetCacheGroupNullableByName(*firstCacheGroup.Name)
	if err != nil {
		t.Errorf("cannot GET Cache Group by name: %v - %v\n", firstCacheGroup.Name, err)
	}

	if cacheGroupResp == nil {
		t.Fatalf("Cache Groups response should not be nil")
	}

	firstParameter := testData.Parameters[0]
	paramResp, _, err := TOSession.GetParameterByName(firstParameter.Name)
	if err != nil {
		t.Errorf("cannot GET Parameter by name: %v - %v\n", firstParameter.Name, err)
	}
	if paramResp == nil {
		t.Fatalf("Parameter response should not be nil")
	}

	cacheGroupID := cacheGroupResp[0].ID
	parameterID := paramResp[0].ID

	resp, _, err := TOSession.CreateCacheGroupParameter(*cacheGroupID, parameterID)
	if err != nil {
		t.Errorf("could not CREATE cache group parameter: %v\n", err)
	}
	if resp == nil {
		t.Fatalf("Cache Group Parameter response should not be nil")
	}
	testData.CacheGroupParameters = append(testData.CacheGroupParameters, resp.Response...)
}

func GetTestCacheGroupParameters(t *testing.T) {
	for _, cgp := range testData.CacheGroupParameters {
		resp, _, err := TOSession.GetCacheGroupParameters(cgp.CacheGroupID)
		if err != nil {
			t.Errorf("cannot GET Parameter by cache group: %v - %v\n", err, resp)
		}
	}
}

func DeleteTestCacheGroupParameters(t *testing.T) {
	for _, cgp := range testData.CacheGroupParameters {
		DeleteTestCacheGroupParameter(t, cgp)
	}
}

func DeleteTestCacheGroupParameter(t *testing.T, cgp tc.CacheGroupParameter) {

	delResp, _, err := TOSession.DeleteCacheGroupParameter(cgp.CacheGroupID, cgp.ParameterID)
	if err != nil {
		t.Errorf("cannot DELETE Parameter by cache group: %v - %v\n", err, delResp)
	}

	// Retrieve the Cache Group Parameter to see if it got deleted
	queryParams := fmt.Sprintf("parameterId=%d", cgp.ParameterID)

	parameters, _, err := TOSession.GetCacheGroupParametersByQueryParams(cgp.CacheGroupID, queryParams)
	if err != nil {
		t.Errorf("error deleting Parameter name: %s\n", err.Error())
	}
	if parameters == nil {
		t.Fatalf("Cache Group Parameters response should not be nil")
	}
	if len(parameters) > 0 {
		t.Errorf("expected Parameter: %d to be to be disassociated from Cache Group: %d\n", cgp.ParameterID, cgp.CacheGroupID)
	}

}
