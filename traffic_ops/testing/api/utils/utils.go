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

package utils

import (
	"encoding/json"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

type ErrorAndMessage struct {
	Error   error
	Message string
}

func FindNeedle(needle string, haystack []string) bool {
	found := false
	for _, s := range haystack {
		if s == needle {
			found = true
			break
		}
	}
	return found
}

func ErrorsToStrings(errs []error) []string {
	errorStrs := []string{}
	for _, errType := range errs {
		et := errType.Error()
		errorStrs = append(errorStrs, et)
	}
	return errorStrs
}

func Compare(t *testing.T, expected []string, alertsStrs []string) {
	sort.Strings(alertsStrs)
	expectedFmt, _ := json.MarshalIndent(expected, "", "  ")
	errorsFmt, _ := json.MarshalIndent(alertsStrs, "", "  ")

	var found bool
	// Compare both directions
	for _, s := range alertsStrs {
		found = FindNeedle(s, expected)
		if !found {
			t.Errorf("\nExpected %s and \n Actual %v must match exactly", string(expectedFmt), string(errorsFmt))
			break
		}
	}

	found = false
	if !found {
		// Compare both directions
		for _, s := range expected {
			found = FindNeedle(s, expected)
			if !found {
				t.Errorf("\nExpected %s and \n Actual %v must match exactly", string(expectedFmt), string(errorsFmt))
				break
			}
		}
	}
}

type CkReqFunc func(*testing.T, toclientlib.ReqInf, interface{}, error)

func CkRequest(c ...CkReqFunc) []CkReqFunc {
	return c
}

func NoError() CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, err error) {
		assert.NoError(t, err, "Expected no error, but got %v", err)
	}
}

func HasError() CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, err error) {
		switch resp := resp.(type) {
		case tc.CacheGroupDetailResponse:
			assert.Error(t, err, "Expected error - but got: %+v %+v", resp.Response, resp.Alerts)
		case tc.CacheGroupsNullableResponse:
			assert.Error(t, err, "Expected error - but got: %+v %+v", resp.Response, resp.Alerts)
		case tc.Alerts:
			assert.Error(t, err, "Expected error - but got: %+v", resp.Alerts)
		default:
			assert.Fail(t, "Failed to find response type")
		}
	}
}

func HasStatus(expectedStatus int) CkReqFunc {
	return func(t *testing.T, reqInf toclientlib.ReqInf, _ interface{}, _ error) {
		assert.Equal(t, expectedStatus, reqInf.StatusCode, "Expected Status Code %d, got %d", expectedStatus, reqInf.StatusCode)
	}
}

func ResponseHasLength(expected int) CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ error) {
		switch resp := resp.(type) {
		case tc.CacheGroupsNullableResponse:
			assert.Equal(t, len(resp.Response), expected, "Expected response object length %d, but got %d", expected, len(resp.Response))
		default:
			assert.Fail(t, "Failed to find response type")
		}
	}
}

func ResponseLengthGreaterOrEqual(expected int) CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ error) {
		switch resp := resp.(type) {
		case tc.CacheGroupsNullableResponse:
			assert.GreaterOrEqual(t, len(resp.Response), expected, "Expected response object length %d, but got %d", expected, len(resp.Response))
		default:
			assert.Fail(t, "Failed to find response type")
		}
	}
}
