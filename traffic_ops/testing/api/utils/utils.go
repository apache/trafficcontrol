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
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	v3client "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	v4client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
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

// V3TestCase is the type of the V3TestData struct.
// Uses nested map to represent the method being tested and the test's description.
type V3TestCase map[string]map[string]V3TestData

// V4TestCase is the type of the V4TestData struct.
// Uses nested map to represent the method being tested and the test's description.
type V4TestCase map[string]map[string]V4TestData

// V3TestData represents the data needed for testing the v3 api endpoints.
type V3TestData struct {
	EndpointId     func() int
	ClientSession  *v3client.Session
	RequestParams  url.Values
	RequestHeaders http.Header
	RequestBody    map[string]interface{}
	Expectations   []CkReqFunc
}

// V4TestData represents the data needed for testing the v4 api endpoints.
type V4TestData struct {
	EndpointId    func() int
	ClientSession *v4client.Session
	RequestOpts   v4client.RequestOptions
	RequestBody   map[string]interface{}
	Expectations  []CkReqFunc
}

// CkReqFunc defines the reusable signature for all other functions that perform checks.
// Common parameters that are checked include the request's info, response, alerts, and errors.
type CkReqFunc func(*testing.T, toclientlib.ReqInf, interface{}, tc.Alerts, error)

// CkRequest wraps CkReqFunc functions, to be concise and reduce lexical tokens
// i.e. Instead of `[]CkReqFunc {` we can use `CkRequest(` in test case declarations.
func CkRequest(c ...CkReqFunc) []CkReqFunc {
	return c
}

// NoError checks that no error was returned (i.e. `nil`).
func NoError() CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, _ tc.Alerts, err error) {
		assert.NoError(t, err, "Expected no error. Got: %v", err)
	}
}

// HasError checks that an error was returned (i.e. not `nil`).
func HasError() CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, alerts tc.Alerts, err error) {
		assert.Error(t, err, "Expected error. Got: %v", alerts)
	}
}

// HasStatus checks that the status code from the request is as expected.
func HasStatus(expectedStatus int) CkReqFunc {
	return func(t *testing.T, reqInf toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		assert.Equal(t, expectedStatus, reqInf.StatusCode, "Expected Status Code: %d Got: %d", expectedStatus, reqInf.StatusCode)
	}
}

// ResponseHasLength checks that the length of the response is as expected.
// Determines that response is a slice before checking the length of the reflected value.
func ResponseHasLength(expected int) CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		rt := reflect.TypeOf(resp)
		switch rt.Kind() {
		case reflect.Slice:
			actual := reflect.ValueOf(resp).Len()
			assert.Equal(t, expected, actual, "Expected response object length: %d Got: %d", expected, actual)
		default:
			t.Errorf("Expected response to be an array. Got: %v", rt)
		}
	}
}

// ResponseLengthGreaterOrEqual checks that the response is greater or equal to the expected length.
// Determines that response is a slice before checking the length of the reflected value.
func ResponseLengthGreaterOrEqual(expected int) CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		rt := reflect.TypeOf(resp)
		switch rt.Kind() {
		case reflect.Slice:
			actual := reflect.ValueOf(resp).Len()
			assert.GreaterOrEqual(t, actual, expected, "Expected response object length: %d Got: %d", expected, actual)
		default:
			t.Errorf("Expected response to be an array. Got: %v", rt)
		}
	}
}
