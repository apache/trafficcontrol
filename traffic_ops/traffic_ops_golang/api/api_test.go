package api

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func ExampleVersion_String() {
	// Because api.Info objects use pointers to Versions, this handles nil
	// without needing the caller to do it - because that's annoying.
	var v *Version
	fmt.Println(v)
	v = &Version{Major: 4, Minor: 20}
	fmt.Println(v.String())
	// Output: {{null}}
	// 4.20
}

func TestCamelCase(t *testing.T) {
	testStrings := []string{"hello_world", "trailing_underscore_", "w_h_a_t____"}
	expected := []string{"helloWorld", "trailingUnderscore", "wHAT"}
	for i, str := range testStrings {
		if toCamelCase(str) != expected[i] {
			t.Errorf("expected: %v error, actual: %v", expected[i], toCamelCase(str))
		}
	}
}

// TestRespWrittenAfterErrFails tests that a WriteResp called after HandleErr will not be written.
func TestRespWrittenAfterErrFails(t *testing.T) {
	w := &MockHTTPResponseWriter{}
	r := &http.Request{URL: &url.URL{}}
	tx := (*sql.Tx)(nil)

	expectedCode := http.StatusUnauthorized
	expectedUserErr := errors.New("user unauthorized")

	HandleErr(w, r, tx, expectedCode, expectedUserErr, nil)
	WriteResp(w, r, "should not be written")

	actualCode := w.Code
	statusVal := r.Context().Value(tc.StatusKey)
	statusInt, ok := statusVal.(int)
	if ok {
		actualCode = statusInt
	}

	if actualCode != expectedCode {
		t.Errorf("code expected: %+v, actual %+v", expectedCode, actualCode)
	}

	alerts := tc.Alerts{}
	if err := json.Unmarshal(w.Body, &alerts); err != nil {
		t.Fatalf("unmarshalling actual body: %v", err)
	}
	for _, alert := range alerts.Alerts {
		if string(alert.Level) != tc.ErrorLevel.String() {
			t.Errorf("alert level expected: '%s', actual: '%s'", tc.ErrorLevel.String(), alert.Level)
		}
	}
}

func TestWriteAlertsObjEmpty(t *testing.T) {
	w := &MockHTTPResponseWriter{}
	r := &http.Request{URL: &url.URL{}}
	a := tc.Alerts{}
	code := http.StatusAlreadyReported

	WriteAlertsObj(w, r, code, a, code)

	resp := struct {
		Response interface{} `json:"response"`
	}{code}
	serialized, _ := json.Marshal(resp)
	if !bytes.Equal(append(serialized[:], '\n'), w.Body[:]) {
		t.Error("expected response to only include object")
	}
	writeAlertsCodeTest(t, *w, code)
}

func TestWriteAlertsObj(t *testing.T) {
	w := &MockHTTPResponseWriter{}
	r := &http.Request{URL: &url.URL{}}
	a := tc.CreateAlerts(tc.WarnLevel, "test")
	code := http.StatusAlreadyReported

	WriteAlertsObj(w, r, code, a, code)

	resp := struct {
		tc.Alerts
		Response interface{} `json:"response"`
	}{a, code}
	serialized, _ := json.Marshal(resp)
	if !bytes.Equal(append(serialized[:], '\n'), w.Body[:]) {
		t.Error("expected response to include alert")
	}
	writeAlertsCodeTest(t, *w, code)
}

func writeAlertsCodeTest(t *testing.T, w MockHTTPResponseWriter, code int) {
	if w.Body == nil || len(w.Body) == 0 {
		t.Error("expected response body to be written to")
	}
	if w.Code != code {
		t.Errorf("expected response code %v, got %v", code, w.Code)
	}
}

func TestWriteResp(t *testing.T) {
	apiWriteTest(t, func(w http.ResponseWriter, r *http.Request) {
		WriteResp(w, r, "foo")
	})
}

func TestWriteRespRaw(t *testing.T) {
	apiWriteTest(t, func(w http.ResponseWriter, r *http.Request) {
		WriteRespRaw(w, r, "foo")
	})
}

func TestWriteRespVals(t *testing.T) {
	apiWriteTest(t, func(w http.ResponseWriter, r *http.Request) {
		WriteRespVals(w, r, "foo", map[string]interface{}{"a": "b"})
	})
}

func TestRespWriter(t *testing.T) {
	apiWriteTest(t, func(w http.ResponseWriter, r *http.Request) {
		RespWriter(w, r, nil)("foo", nil)
	})
}

func TestRespWriterVals(t *testing.T) {
	apiWriteTest(t, func(w http.ResponseWriter, r *http.Request) {
		RespWriterVals(w, r, nil, map[string]interface{}{"a": "b"})("foo", nil)
	})
}

func TestWriteRespAlert(t *testing.T) {
	apiWriteTest(t, func(w http.ResponseWriter, r *http.Request) {
		WriteRespAlert(w, r, tc.ErrorLevel, "foo error")
	})
}

func TestWriteRespAlertObj(t *testing.T) {
	apiWriteTest(t, func(w http.ResponseWriter, r *http.Request) {
		WriteRespAlertObj(w, r, tc.ErrorLevel, "foo error", "bar")
	})
}

// apiWriteTest tests that an API write func succeeds and writes a body and a 200.
func apiWriteTest(t *testing.T, write func(w http.ResponseWriter, r *http.Request)) {
	w := &MockHTTPResponseWriter{}
	r := &http.Request{URL: &url.URL{}}

	write(w, r)

	if w.Code == 0 {
		w.Code = http.StatusOK // emulate behavior of w.Write
	}

	actualCode := w.Code
	statusVal := r.Context().Value(tc.StatusKey)
	statusInt, ok := statusVal.(int)
	if ok {
		actualCode = statusInt
	}

	expectedCode := http.StatusOK

	if actualCode != expectedCode {
		t.Errorf("code expected: %+v, actual %+v", expectedCode, actualCode)
	}

	if len(w.Body) == 0 {
		t.Errorf("body len expected: >0, actual 0")
	}
}

type MockHTTPResponseWriter struct {
	Code int
	Body []byte
}

func (i *MockHTTPResponseWriter) WriteHeader(rc int) {
	i.Code = rc
}

func (i *MockHTTPResponseWriter) Write(b []byte) (int, error) {
	i.Body = append(i.Body, b...)
	return len(b), nil
}

func (i *MockHTTPResponseWriter) Header() http.Header {
	return http.Header{}
}

func TestParseRestrictFKConstraint(t *testing.T) {
	var testCases = []struct {
		description        string
		storageError       pq.Error
		expectedReturnCode int
	}{
		{
			description: "FK Constraint Error",
			storageError: pq.Error{
				Message: "update or delete on table \"foo\" violates foreign key constraint \"fk_foo_bar\" on table \"bar\"",
			},
			expectedReturnCode: http.StatusBadRequest,
		},
		{
			description: "FK Constraint Error with underscores in table name",
			storageError: pq.Error{
				Message: "update or delete on table \"foo_ser\" violates foreign key constraint \"fk_foo_bar\" on table \"bar_cap\"",
			},
			expectedReturnCode: http.StatusBadRequest,
		},
		{
			description: "Non FK Constraint Error",
			storageError: pq.Error{
				Message: "connection error",
			},
			expectedReturnCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			t.Log("Starting test scenario: ", tc.description)
			_, _, sc := parseRestrictFKConstraint(&tc.storageError)
			if sc != tc.expectedReturnCode {
				t.Errorf("code expected: %v, actual %v", tc.expectedReturnCode, sc)
			}
		})
	}
}

func TestAPIInfo_WriteOKResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := APIInfo{
		request: r,
		w:       w,
	}
	code, userErr, sysErr := inf.WriteOKResponse("test")
	if code != http.StatusOK {
		t.Errorf("WriteOKResponse should return a %d %s code, got: %d %s", http.StatusOK, http.StatusText(http.StatusOK), code, http.StatusText(code))
	}
	if userErr != nil {
		t.Errorf("Unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}

	if w.Code != http.StatusOK {
		t.Errorf("incorrect response status code; want: %d, got: %d", http.StatusOK, w.Code)
	}
}
func TestAPIInfo_WriteOKResponseWithSummary(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := APIInfo{
		request: r,
		w:       w,
	}
	code, userErr, sysErr := inf.WriteOKResponseWithSummary("test", 42)
	if code != http.StatusOK {
		t.Errorf("WriteOKResponseWithSummary should return a %d %s code, got: %d %s", http.StatusOK, http.StatusText(http.StatusOK), code, http.StatusText(code))
	}
	if userErr != nil {
		t.Errorf("Unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}

	if w.Code != http.StatusOK {
		t.Errorf("incorrect response status code; want: %d, got: %d", http.StatusOK, w.Code)
	}
}
func TestAPIInfo_WriteNotModifiedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := APIInfo{
		request: r,
		w:       w,
	}
	code, userErr, sysErr := inf.WriteNotModifiedResponse(time.Time{})
	if code != http.StatusNotModified {
		t.Errorf("WriteNotModifiedResponse should return a %d %s code, got: %d %s", http.StatusNotModified, http.StatusText(http.StatusNotModified), code, http.StatusText(code))
	}
	if userErr != nil {
		t.Errorf("Unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}

	if w.Code != http.StatusNotModified {
		t.Errorf("incorrect response status code; want: %d, got: %d", http.StatusNotModified, w.Code)
	}
}

func TestAPIInfo_WriteSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := APIInfo{
		request: r,
		w:       w,
	}
	code, userErr, sysErr := inf.WriteSuccessResponse("test", "quest")
	if code != http.StatusOK {
		t.Errorf("WriteSuccessResponse should return a %d %s code, got: %d %s", http.StatusOK, http.StatusText(http.StatusOK), code, http.StatusText(code))
	}
	if userErr != nil {
		t.Errorf("Unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}

	if w.Code != http.StatusOK {
		t.Errorf("incorrect response status code; want: %d, got: %d", http.StatusOK, w.Code)
	}
	var alerts tc.Alerts
	if err := json.NewDecoder(w.Body).Decode(&alerts); err != nil {
		t.Fatalf("couldn't decode response body: %v", err)
	}

	if len(alerts.Alerts) != 1 {
		t.Fatalf("expected exactly one alert; got: %d", len(alerts.Alerts))
	}
	alert := alerts.Alerts[0]
	if alert.Level != tc.SuccessLevel.String() {
		t.Errorf("Incorrect alert level; want: %s, got: %s", tc.SuccessLevel, alert.Level)
	}
	if alert.Text != "quest" {
		t.Errorf("Incorrect alert text; want: 'quest', got: '%s'", alert.Text)
	}
}

func TestAPIInfo_WriteCreatedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	inf := APIInfo{
		request: r,
		Version: &Version{Major: 420, Minor: 9001},
		w:       w,
	}
	code, userErr, sysErr := inf.WriteCreatedResponse("test", "quest", "mypath")
	if code != http.StatusCreated {
		t.Errorf("WriteCreatedResponse should return a %d %s code, got: %d %s", http.StatusCreated, http.StatusText(http.StatusCreated), code, http.StatusText(code))
	}
	if userErr != nil {
		t.Errorf("Unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("incorrect response status code; want: %d, got: %d", http.StatusCreated, w.Code)
	}
	if locHdr := w.Header().Get(rfc.Location); locHdr != "/api/420.9001/mypath" {
		t.Errorf("incorrect '%s' header value; want: '/api/420.9001/mypath', got: '%s'", rfc.Location, locHdr)
	}
	var alerts tc.Alerts
	if err := json.NewDecoder(w.Body).Decode(&alerts); err != nil {
		t.Fatalf("couldn't decode response body: %v", err)
	}

	if len(alerts.Alerts) != 1 {
		t.Fatalf("expected exactly one alert; got: %d", len(alerts.Alerts))
	}
	alert := alerts.Alerts[0]
	if alert.Level != tc.SuccessLevel.String() {
		t.Errorf("Incorrect alert level; want: %s, got: %s", tc.SuccessLevel, alert.Level)
	}
	if alert.Text != "quest" {
		t.Errorf("Incorrect alert text; want: 'quest', got: '%s'", alert.Text)
	}
}

func TestAPIInfo_RequestHeaders(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("test", "quest")

	inf := APIInfo{
		request: r,
	}
	testHdr := inf.RequestHeaders().Get("test")
	if testHdr != "quest" {
		t.Errorf("should have retrieved the 'test' header (expected value: 'quest'), but found that header to have value: '%s'", testHdr)
	}

}

func TestAPIInfo_SetLastModified(t *testing.T) {
	w := httptest.NewRecorder()
	inf := APIInfo{w: w}
	tm := time.Now().Truncate(time.Second).UTC()
	inf.SetLastModified(tm)

	wLMHdr := w.Header().Get(rfc.LastModified)
	lm, err := time.Parse(time.RFC1123, wLMHdr)
	if err != nil {
		t.Fatalf("Failed to parse the response writer's '%s' header as an RFC1123 timestamp: %v", rfc.LastModified, err)
	}

	// For unknown reasons, our API always adds a second to the truncated time
	// value for LastModified headers. I suspect it's a poor attempt at rounding
	// - for which the `Round` method ought to be used instead.
	if expected := tm.Add(time.Second); lm != expected {
		t.Errorf("Incorrect time set as '%s' header; want: %s, got: %s", rfc.LastModified, expected.Format(time.RFC3339Nano), lm.Format(time.RFC3339Nano))
	}
}

func TestAPIInfo_DecodeBody(t *testing.T) {
	inf := APIInfo{
		request: httptest.NewRequest(http.MethodConnect, "/", strings.NewReader(`{"test": "quest"}`)),
	}

	var out struct {
		Test string `json:"test"`
	}
	if err := inf.DecodeBody(&out); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if out.Test != "quest" {
		t.Errorf(`incorrect request body parsed; want: {"test": "quest"}, got: {"test": "%s"}`, out.Test)
	}
}
