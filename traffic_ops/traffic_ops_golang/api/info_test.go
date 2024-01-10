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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault/backends/disabled"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func buildContextWithout(t *testing.T, key any, beginnable bool) context.Context {
	t.Helper()
	ctx := context.Background()

	if key != DBContextKey {
		d, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open a stub database connection: %v", err)
		}
		db := sqlx.NewDb(d, "sqlmock")
		if beginnable {
			mock.ExpectBegin()
		}
		ctx = context.WithValue(ctx, DBContextKey, db)
	}
	if key != ConfigContextKey {
		ctx = context.WithValue(ctx, ConfigContextKey, &config.Config{ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{DBQueryTimeoutSeconds: 1}})
	}
	if key != TrafficVaultContextKey {
		ctx = context.WithValue(ctx, TrafficVaultContextKey, &disabled.Disabled{})
	}
	if key != ReqIDContextKey {
		ctx = context.WithValue(ctx, ReqIDContextKey, uint64(1))
	}
	if key != auth.CurrentUserKey {
		ctx = context.WithValue(ctx, auth.CurrentUserKey, auth.CurrentUser{})
	}
	if key != PathParamsKey {
		ctx = context.WithValue(ctx, PathParamsKey, map[string]string{})
	}
	return ctx
}

func buildRequest(t *testing.T, withoutContext any, beginnable bool, params map[string]string) *http.Request {
	t.Helper()
	var values url.Values = make(map[string][]string, len(params))
	for p, v := range params {
		values[p] = []string{v}
	}

	r := httptest.NewRequest(http.MethodConnect, "/api/5.0/ping?"+values.Encode(), nil)

	return r.WithContext(buildContextWithout(t, withoutContext, beginnable))
}

func testNewInfo_MissingContextKey(key any) func(*testing.T) {
	return func(t *testing.T) {
		r := buildRequest(t, key, false, nil)
		_, sysErr, _, code := NewInfo(r, nil, nil)
		if sysErr == nil {
			t.Errorf("Expected non-nil system error, got: nil")
		} else {
			t.Log("Received expected system error:", sysErr)
		}
		if code != http.StatusInternalServerError {
			t.Errorf("Incorrect status code for context missing '%v'; want: %d, got: %d", key, http.StatusInternalServerError, code)
		}
	}
}

func TestNewInfo(t *testing.T) {
	for _, key := range []any{
		DBContextKey,
		ConfigContextKey,
		TrafficVaultContextKey,
		ReqIDContextKey,
		auth.CurrentUserKey,
	} {
		t.Run(fmt.Sprintf("missing '%v' context key", key), testNewInfo_MissingContextKey(key))
	}

	// TODO: This really should return a system-internal error. But I'm testing
	// the current behavior, soo... ¯\_(ツ)_/¯
	r := buildRequest(t, nil, false, nil)
	_, sysErr, userErr, code := NewInfo(r, nil, nil)
	if userErr == nil {
		t.Errorf("Expected non-nil user error, got: nil")
	} else {
		t.Log("Received expected user error:", userErr)
	}
	if code != http.StatusInternalServerError {
		t.Errorf("Incorrect status code for unable to start a database transaction; want: %d, got: %d", http.StatusInternalServerError, code)
	}

	// TODO: shouldn't this be a user-facing error? It returns
	// http.StatusBadRequest but a nil user error... why?
	r = buildRequest(t, nil, false, nil)
	_, sysErr, userErr, code = NewInfo(r, []string{"testquest"}, nil)
	if sysErr == nil {
		t.Error("Expected a system error; got: nil")
	} else {
		t.Log("Recieved expected system error:", sysErr)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
	if code != http.StatusBadRequest {
		t.Errorf("Incorrect status code for missing required parameter; want: %d, got: %d", http.StatusBadRequest, code)
	}

	r = buildRequest(t, nil, true, nil)
	_, sysErr, userErr, _ = NewInfo(r, nil, nil)
	if userErr != nil {
		t.Error("Unexpected user error:", userErr)
	}
	if sysErr != nil {
		t.Error("Unexpected system error:", sysErr)
	}
}

func TestInfo_CheckPrecondition(t *testing.T) {
	var inf Info
	code, userErr, sysErr := inf.CheckPrecondition("anything", nil)
	if code != http.StatusInternalServerError {
		t.Errorf("incorrect status code for unitialized info structure; want: %d, got: %d", http.StatusInternalServerError, code)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
	if sysErr == nil {
		t.Errorf("Expected a system error; got: nil")
	} else if !errors.Is(sysErr, NilRequestError) {
		t.Errorf("Incorrect system error; want: %v, got: %v", NilRequestError, sysErr)
	} else {
		t.Log("Received expected system error:", sysErr)
	}

	inf.request = httptest.NewRequest(http.MethodConnect, "/", nil)
	code, userErr, sysErr = inf.CheckPrecondition("anything", nil)
	if code != http.StatusOK {
		t.Errorf("incorrect status code for a request with no precondition headers; want: %d, got: %d", http.StatusOK, code)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}

	// We have to use `time.Now` because for whatever reason the ETag parser
	// refuses to aknowledge timestamps more than 20 years in either direction
	// from whatever `time.Now` returns at the time the parser runs.
	// TODO: Fix that?
	etag := rfc.ETag(time.Now())
	inf.request.Header.Add(rfc.IfMatch, etag)
	code, userErr, sysErr = inf.CheckPrecondition("anything", nil)
	if code != http.StatusInternalServerError {
		t.Errorf("incorrect status code for nil db transaction; want: %d, got: %d", http.StatusInternalServerError, code)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
	if sysErr == nil {
		t.Errorf("Expected a system error; got: nil")
	} else if !errors.Is(sysErr, NilTransactionError) {
		t.Errorf("Incorrect system error; want: %v, got: %v", NilTransactionError, sysErr)
	} else {
		t.Log("Received expected system error:", sysErr)
	}

	d, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}
	db := sqlx.NewDb(d, "sqlmock")
	mock.ExpectBegin()
	inf.Tx = db.MustBegin()

	code, userErr, sysErr = inf.CheckPrecondition("anything", nil)
	if code != http.StatusInternalServerError {
		t.Errorf("incorrect status code for database query error; want: %d, got: %d", http.StatusInternalServerError, code)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
	if sysErr == nil {
		t.Errorf("Expected a system error; got: nil")
	} else {
		t.Log("Received expected system error:", sysErr)
	}

	rows := sqlmock.NewRows([]string{"last_updated"})
	rows.AddRow(time.Now().Add(time.Hour))
	mock.ExpectQuery("^anything$").WillReturnRows(rows)
	code, userErr, sysErr = inf.CheckPrecondition("anything", nil)
	if code != http.StatusPreconditionFailed {
		t.Errorf("incorrect status code for ETag match failure; want: %d, got: %d", http.StatusPreconditionFailed, code)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}
	if userErr == nil {
		t.Errorf("Expected a user-facing error; got: nil")
	} else if !errors.Is(userErr, ResourceModifiedError) {
		t.Errorf("Incorrect user-facing error; want: %v, got: %v", ResourceModifiedError, userErr)
	} else {
		t.Log("Received expected user-facing error:", userErr)
	}

	rows = sqlmock.NewRows([]string{"last_updated"})
	rows.AddRow(time.Now().Add(-time.Hour))
	mock.ExpectQuery("^anything$").WillReturnRows(rows)
	code, userErr, sysErr = inf.CheckPrecondition("anything", nil)
	if code != http.StatusOK {
		t.Errorf("incorrect status code for ETag match success (with no unmodified-since time); want: %d, got: %d", http.StatusOK, code)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}

	ius := rfc.FormatHTTPDate(time.Now())
	rows = sqlmock.NewRows([]string{"last_updated"})
	rows.AddRow(time.Now().Add(time.Hour))
	mock.ExpectQuery("^anything$").WillReturnRows(rows)
	inf.request.Header.Add(rfc.IfUnmodifiedSince, ius)
	inf.request.Header.Del(rfc.IfMatch)
	code, userErr, sysErr = inf.CheckPrecondition("anything", nil)
	if code != http.StatusPreconditionFailed {
		t.Errorf("incorrect status code for ETag match failure; want: %d, got: %d", http.StatusPreconditionFailed, code)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}
	if userErr == nil {
		t.Errorf("Expected a user-facing error; got: nil")
	} else if !errors.Is(userErr, ResourceModifiedError) {
		t.Errorf("Incorrect user-facing error; want: %v, got: %v", ResourceModifiedError, userErr)
	} else {
		t.Log("Received expected user-facing error:", userErr)
	}

	rows = sqlmock.NewRows([]string{"last_updated"})
	rows.AddRow(time.Now().Add(-time.Hour))
	mock.ExpectQuery("^anything$").WillReturnRows(rows)
	code, userErr, sysErr = inf.CheckPrecondition("anything", nil)
	if code != http.StatusOK {
		t.Errorf("incorrect status code for ETag match success (with no unmodified-since time); want: %d, got: %d", http.StatusOK, code)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
}

func TestInfo_Close(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open stub database connection: %v", err)
	}
	mock.ExpectBegin()

	called := false
	inf := Info{
		CancelTx: func() {
			called = true
		},
		Tx: sqlx.NewDb(db, "sqlmock").MustBegin(),
	}

	mock.ExpectCommit()
	inf.Close()
	if !called {
		t.Errorf("Expected cancel function to be called when the Info is closed; but it wasn't")
	}

	called = false
	mock.ExpectBegin()
	inf.Tx = sqlx.NewDb(db, "sqlmock").MustBegin()
	mock.ExpectCommit().WillReturnError(errors.New("testquest"))
	inf.Close()
	if !called {
		t.Errorf("Expected cancel function to be called when the Info is closed (even if an error occurs commiting the transaction); but it wasn't")
	}
}

func TestInfo_WriteOKResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := Info{
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

func TestInfo_WriteOKResponseWithSummary(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := Info{
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

func TestInfo_WriteNotModifiedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := Info{
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

func TestInfo_WriteSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	inf := Info{
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

func TestInfo_WriteCreatedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	inf := Info{
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

func TestInfo_RequestHeaders(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("test", "quest")

	inf := Info{
		request: r,
	}
	testHdr := inf.RequestHeaders().Get("test")
	if testHdr != "quest" {
		t.Errorf("should have retrieved the 'test' header (expected value: 'quest'), but found that header to have value: '%s'", testHdr)
	}
}

func TestInfo_SetLastModified(t *testing.T) {
	w := httptest.NewRecorder()
	inf := Info{w: w}
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

func TestInfo_DecodeBody(t *testing.T) {
	inf := Info{
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

func ExampleInfo_SendMail() {
	inf := Info{
		Config: &config.Config{
			SMTP: &config.ConfigSMTP{
				// Note that we're not actually sending an email in this
				// example. In fact it's explicitly disabled!
				Enabled: false,
			},
		},
	}

	code, _, err := inf.SendMail(rfc.EmailAddress{}, []byte("anything"))
	fmt.Println(code, "-", err.Error())
	// Output: 500 - SMTP is not enabled; mail cannot be sent
}

func TestInfo_IsResourceAuthorizedToUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open stub database connection: %v", err)
	}

	mock.ExpectBegin()
	inf := Info{
		Tx:   sqlx.NewDb(db, "sqlmock").MustBegin(),
		User: &auth.CurrentUser{},
	}

	mock.ExpectQuery(".").WillReturnError(errors.New("testquest"))
	ok, err := inf.IsResourceAuthorizedToCurrentUser(-1)
	if ok {
		t.Errorf("Expected a query failure to report the user is not authorized")
	}
	if err == nil {
		t.Errorf("Expected an error when a query error occurs, but didn't get one")
	} else {
		t.Logf("Received expected error: %v", err)
	}
}

func TestInfo_CreateChangeLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open stub database connection: %v", err)
	}

	mock.ExpectBegin()
	inf := Info{
		Tx: sqlx.NewDb(db, "sqlmock").MustBegin(),
		User: &auth.CurrentUser{
			ID: 1,
		},
	}

	msg := "anything"
	mock.ExpectExec("INSERT INTO log").WithArgs(ApiChange, msg, inf.User.ID)
	inf.CreateChangeLog(msg)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}

	mock.ExpectExec("INSERT INTO log").WithArgs(ApiChange, msg, inf.User.ID).WillReturnError(errors.New("testquest"))
	inf.CreateChangeLog(msg)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestInfo_UseIMS(t *testing.T) {
	r := httptest.NewRequest(http.MethodConnect, "", nil)
	r.Header.Add(rfc.IfModifiedSince, "doesn't matter")
	inf := Info{
		request: r,
	}
	use := inf.UseIMS()
	if use {
		t.Error("expected IMS to not be used if the Info structure has no Config")
	}

	inf = Info{
		Config: &config.Config{
			UseIMS: true,
		},
	}
	use = inf.UseIMS()
	if use {
		t.Error("expected IMS to not be used if the Info structure has no request")
	}

	inf.request = r
	inf.Config.UseIMS = false
	use = inf.UseIMS()
	if use {
		t.Error("expected IMS to not be used if it's configured not to")
	}

	inf.Config.UseIMS = true
	inf.request.Header.Del(rfc.IfModifiedSince)
	use = inf.UseIMS()
	if use {
		t.Error("expected IMS to not be used if the request has no If-Modified-Since header")
	}

	inf.request.Header.Add(rfc.IfModifiedSince, "literally anything")
	use = inf.UseIMS()
	if !use {
		t.Error("expected IMS to be used when it's configured to and the request includes an If-Modified-Since header")
	}
}

func ExampleInfo_CreateInfluxClient() {
	inf := Info{
		Config: &config.Config{
			InfluxEnabled: false,
			ConfigInflux:  &config.ConfigInflux{
				// ...
			},
		},
	}

	// These are BOTH `nil` when Influx is disabled!
	client, err := inf.CreateInfluxClient()
	fmt.Println(client, err)
	// Output: <nil> <nil>
}

// TODO: the situation where the actual influx library would return a nil client
// but also a nil error is untestable - because the library doesn't ever do that
// under any circumstances. So... figure out how to test that, or change the
// code to not handle unnecessary circumstances.
func TestInfo_CreateInfluxClient(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open stub database connection: %v", err)
	}

	mock.ExpectBegin()
	inf := Info{
		Config: &config.Config{
			ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{
				ReadTimeout: 2,
			},
			InfluxEnabled: true,
			Version:       "9.9.9",
		},
		Tx:   sqlx.NewDb(db, "sqlmock").MustBegin(),
		User: &auth.CurrentUser{},
	}
	client, err := inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error when Influx is not configured: %v", err)
	}
	if client != nil {
		t.Error("Client should be nil when Influx is not configured")
	}

	inf.Config.ConfigInflux = &config.ConfigInflux{
		User:     "user",
		Password: "password",
		Secure:   new(bool),
	}

	mock.ExpectQuery("^SELECT").WillReturnError(errors.New("testquest"))
	client, err = inf.CreateInfluxClient()
	if err == nil {
		t.Error("Expected to get an error when the influx server query fails, but didn't")
	} else {
		t.Log("Received expected error:", err)
	}
	if client != nil {
		t.Error("Client should be nil when querying for influx servers fails")
	}

	rows := sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("not a valid hostname, to make the Influx library return an error", 1, 2)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	_, err = inf.CreateInfluxClient()
	if err == nil {
		t.Error("Expected an error trying to create a client for an invalid URL, but didn't get one")
	} else {
		t.Logf("Received expected error: %v", err)
	}

	rows = sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("test.quest", 1, 2)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	client, err = inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil when everything's supposed to succeed")
	}

	rows = sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("test.quest", 0, 2)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	client, err = inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil - even when the influx server's TCP port is incorrectly configured")
	}

	inf.Config.ConfigInflux.Secure = nil
	rows = sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("test.quest", 1, 2)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	client, err = inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil - even influx configuration leaves 'secure' null/undefined")
	}

	inf.Config.ConfigInflux.Secure = util.Ptr(true)
	rows = sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("test.quest", 1, nil)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	client, err = inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil - even when the server's HTTPS port is not configured")
	}

	rows = sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("test.quest", 1, -1)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	client, err = inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil - even when the server's HTTPS port is negative")
	}

	rows = sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("test.quest", 1, 65536)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	client, err = inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil - even when the server's HTTPS port is out of the valid port number range")
	}

	rows = sqlmock.NewRows([]string{"fqdn", "tcp_port", "https_port"})
	rows.AddRow("test.quest", 1, 2)
	mock.ExpectQuery("^SELECT").WillReturnRows(rows)
	client, err = inf.CreateInfluxClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("client should not be nil when everything's supposed to be working (but with HTTPS this time)")
	}
}

func ExampleInfo_DefaultSort() {
	inf := Info{
		Params: map[string]string{},
	}

	inf.DefaultSort("testquest")
	fmt.Println(inf.Params["orderby"])

	inf.DefaultSort("id")
	fmt.Println(inf.Params["orderby"])

	// Output: testquest
	// testquest
}
