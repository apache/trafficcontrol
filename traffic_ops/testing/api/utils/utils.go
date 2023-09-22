// Package utils provides helpful utilities for easing the process of testing
// Traffic Ops API clients. The functions here don't depend on any particular
// API version, so they need not be copy/pasted along with the entire testing
// suites.
package utils

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
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
	v3client "github.com/apache/trafficcontrol/v8/traffic_ops/v3-client"
	v4client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
	v5client "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

// FindNeedle searches a "haystack" slice of values for the the "needle" value,
// returning true if the value is in the haystack, false otherwise.
func FindNeedle[T comparable](needle T, haystack []T) bool {
	found := false
	for _, s := range haystack {
		if s == needle {
			found = true
			break
		}
	}
	return found
}

// ErrorsToStrings converts a slice of errors to a slice of their error
// messages.
func ErrorsToStrings(errs []error) []string {
	errorStrs := make([]string, 0, len(errs))
	for _, err := range errs {
		errorStrs = append(errorStrs, err.Error())
	}
	return errorStrs
}

// Compare compares a set of expected alert messages to those actually received.
// It checks for the existence of each expected alert, not that they appear in
// any particular order. It also checks that no unexpected alert strings exist.
// Note that this isn't particularly efficient, which is fine because it's only
// meant to be used in testing.
func Compare(t *testing.T, expected []string, alertsStrs []string) {
	t.Helper()
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

// CreateV3Session creates a session for client v4 using the passed in username and password.
func CreateV3Session(t *testing.T, trafficOpsURL string, username string, password string, toReqTimeout int) *v3client.Session {
	t.Helper()
	userSession, _, err := v3client.LoginWithAgent(trafficOpsURL, username, password, true, "to-api-v3-client-tests", false, time.Second*time.Duration(toReqTimeout))
	assert.RequireNoError(t, err, "Could not login with user %v: %v", username, err)
	return userSession
}

// CreateV4Session creates a session for client v4 using the passed in username and password.
func CreateV4Session(t *testing.T, trafficOpsURL string, username string, password string, toReqTimeout int) *v4client.Session {
	t.Helper()
	userSession, _, err := v4client.LoginWithAgent(trafficOpsURL, username, password, true, "to-api-v4-client-tests", false, time.Second*time.Duration(toReqTimeout))
	assert.RequireNoError(t, err, "Could not login with user %v: %v", username, err)
	return userSession
}

// CreateV5Session creates a session for client v5 using the passed in username and password.
func CreateV5Session(t *testing.T, trafficOpsURL, username, password string, toReqTimeout int) *v5client.Session {
	t.Helper()
	userSession, _, err := v5client.LoginWithAgent(trafficOpsURL, username, password, true, "to-api-v5-client-tests", false, time.Second*time.Duration(toReqTimeout))
	assert.RequireNoError(t, err, "Could not login with user %v: %v", username, err)
	return userSession
}

// V3TestData represents the data needed for testing the v3 api endpoints.
type V3TestData struct {
	EndpointID     func() int
	ClientSession  *v3client.Session
	RequestParams  url.Values
	RequestHeaders http.Header
	RequestBody    map[string]interface{}
	PreReqFuncs    []func()
	Expectations   []CkReqFunc
}

// V3TestDataT represents the data needed for testing the v3 api endpoints.
type V3TestDataT[B any] struct {
	EndpointID     func() int
	ClientSession  *v3client.Session
	RequestParams  url.Values
	RequestHeaders http.Header
	RequestBody    B
	Expectations   []CkReqFunc
}

// V4TestData represents the data needed for testing the v4 api endpoints.
type V4TestData struct {
	EndpointID    func() int
	ClientSession *v4client.Session
	RequestOpts   v4client.RequestOptions
	RequestBody   map[string]interface{}
	PreReqFuncs   []func()
	Expectations  []CkReqFunc
}

// V5TestData represents the data needed for testing the v5 api endpoints.
type V5TestData struct {
	EndpointID    func() int
	ClientSession *v5client.Session
	RequestOpts   v5client.RequestOptions
	RequestBody   map[string]interface{}
	PreReqFuncs   []func()
	Expectations  []CkReqFunc
}

type clientSession interface {
	v4client.Session | v5client.Session
}

type requestOpts interface {
	v4client.RequestOptions | v5client.RequestOptions
}

// TestData represents the data needed for testing the api endpoints.
type TestData[C clientSession, R requestOpts, B any] struct {
	EndpointID    func() int
	ClientSession *C
	RequestOpts   R
	RequestBody   B
	Expectations  []CkReqFunc
}

// V3TestCase is the type of the V3TestData struct.
// Uses nested map to represent the method being tested and the test's description.
type V3TestCase map[string]map[string]V3TestData

// V3TestCaseT is the type of the V3TestDataT struct.
// Uses nested map to represent the method being tested and the test's description.
type V3TestCaseT[B any] map[string]map[string]V3TestDataT[B]

// V4TestCase is the type of the V4TestData struct.
// Uses nested map to represent the method being tested and the test's description.
type V4TestCase map[string]map[string]V4TestData

// V5TestCase is a map of test names to maps of HTTP request method descriptions
// to V5TestData structures.
// Uses nested map to represent the method being tested and the test's description.
type V5TestCase map[string]map[string]V5TestData

type TestCase[C clientSession, R requestOpts, B any] map[string]map[string]TestData[C, R, B]

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
		t.Helper()
		assert.NoError(t, err, "Expected no error. Got: %v", err)
	}
}

// HasError checks that an error was returned (i.e. not `nil`).
func HasError() CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, alerts tc.Alerts, err error) {
		t.Helper()
		assert.Error(t, err, "Expected error. Got: %v", alerts)
	}
}

// HasStatus checks that the status code from the request is as expected.
func HasStatus(expectedStatus int) CkReqFunc {
	return func(t *testing.T, reqInf toclientlib.ReqInf, _ interface{}, _ tc.Alerts, _ error) {
		t.Helper()
		assert.Equal(t, expectedStatus, reqInf.StatusCode, "Expected Status Code: %d Got: %d", expectedStatus, reqInf.StatusCode)
	}
}

// HasAlertLevel checks that the alert from the request matches the expected level.
func HasAlertLevel(expectedAlertLevel string) CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, _ interface{}, alerts tc.Alerts, _ error) {
		t.Helper()
		assert.RequireNotNil(t, alerts, "Expected alerts to not be nil.")
		found := false
		for _, alert := range alerts.Alerts {
			if expectedAlertLevel == alert.Level {
				found = true
			}
		}
		assert.Equal(t, true, found, "Expected to find Alert Level: %s", expectedAlertLevel)
	}
}

// ResponseHasLength checks that the length of the response is as expected.
// Determines that response is a slice before checking the length of the reflected value.
func ResponseHasLength(expected int) CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		t.Helper()
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
		t.Helper()
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
