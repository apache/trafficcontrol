package middleware

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
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tocookie"

	"github.com/jmoiron/sqlx"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var debugLogging = flag.Bool("debug", false, "enable debug logging in test")

// TestWrapHeaders checks that appropriate default headers are added to a request
func TestWrapHeaders(t *testing.T) {
	body := "We are here!!"
	f := WrapHeaders(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", ".", nil)
	if err != nil {
		t.Error("Error creating new request")
	}

	// Call to add the headers
	f(w, r)
	if w.Body.String() != body {
		t.Error("Expected body", body, "got", w.Body.String())
	}

	expected := map[string][]string{
		"Access-Control-Allow-Credentials": nil,
		"Access-Control-Allow-Headers":     nil,
		"Access-Control-Allow-Methods":     nil,
		"Access-Control-Allow-Origin":      nil,
		rfc.Vary:                           {rfc.AcceptEncoding},
		"Content-Type":                     nil,
		"Whole-Content-Sha512":             nil,
		"X-Server-Name":                    nil,
		rfc.PermissionsPolicy:              {"interest-cohort=()"},
	}

	if len(expected) != len(w.HeaderMap) {
		t.Error("Expected", len(expected), "header, got", len(w.HeaderMap))
	}
	m := w.Header()
	for k := range expected {
		if _, ok := m[k]; !ok {
			t.Error("Expected header", k, "not found")
		} else if len(expected[k]) > 0 && !reflect.DeepEqual(expected[k], m[k]) {
			t.Errorf("expected: %v, actual: %v", expected[k], m[k])
		}
	}
}

// TestWrapPanicRecover checks that a recovered panic returns a 500
func TestWrapPanicRecover(t *testing.T) {
	f := WrapPanicRecover(func(w http.ResponseWriter, r *http.Request) {
		var foo *string
		bar := *foo // will throw nil dereference panic
		w.Write([]byte(bar))
	})
	f = WrapHeaders(f)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "/", nil)
	if err != nil {
		t.Error("Error creating new request")
	}

	// Call to wrap the panic recovery
	f(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Error("expected panic recovery to return a 500, got", w.Code)
	}
}

// TestGzip checks that if Accept-Encoding contains "gzip" that the body is indeed gzip'd
func TestGzip(t *testing.T) {
	body := "am I gzip'd?"
	gz := bytes.Buffer{}
	zw := gzip.NewWriter(&gz)

	if _, err := zw.Write([]byte(body)); err != nil {
		t.Error("Error gzipping", err)
	}

	if err := zw.Close(); err != nil {
		t.Error("Error closing gzipper", err)
	}

	f := WrapHeaders(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	})

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "/", nil)
	if err != nil {
		t.Error("Error creating new request")
	}

	f(w, r)

	// body should not be gzip'd
	if !bytes.Equal(w.Body.Bytes(), []byte(body)) {
		t.Error("Expected body to be NOT gzip'd!")
	}

	// Call with gzip
	w = httptest.NewRecorder()
	r.Header.Add("Accept-Encoding", "gzip")
	f(w, r)
	if !bytes.Equal(w.Body.Bytes(), gz.Bytes()) {
		t.Error("Expected body to be gzip'd!")
	}
}

func newRWPair(t *testing.T, cookie *http.Cookie) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "/api/4.0/blah", nil)
	if err != nil {
		t.Fatalf("Failed to create new request: %v", err)
	}

	if cookie != nil {
		r.Header.Add("Cookie", tocookie.Name+"="+cookie.Value)
	}

	return w, r
}

func TestWrapAuth(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	userName := "user1"
	id := 1
	secret := "secret"

	rows := sqlmock.NewRows([]string{"priv_level", "username", "id", "tenant_id"})
	rows.AddRow(30, "user1", 1, 1)
	mock.ExpectQuery("SELECT").WithArgs(userName).WillReturnRows(rows)

	authBase := AuthBase{secret, nil}

	cookie := tocookie.GetCookie(userName, time.Minute, secret)

	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := auth.GetCurrentUser(ctx)
		if err != nil {
			t.Fatalf("unable to get privLevel: %v", err)
			return
		}

		respBts, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("unable to marshal: %v", err)
			return
		}

		w.Header().Set(rfc.ContentType, rfc.ApplicationJSON)
		fmt.Fprintf(w, "%s", respBts)
	}

	authWrapper := authBase.GetWrapper(15)

	f := authWrapper(handler)

	w, r := newRWPair(t, cookie)

	expected := auth.CurrentUser{UserName: userName, ID: id, PrivLevel: 30, TenantID: 1}

	expectedBody, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("unable to marshal: %v", err)
	}

	r = r.WithContext(context.WithValue(context.Background(), api.DBContextKey, db))
	r = r.WithContext(context.WithValue(r.Context(), api.ConfigContextKey, &config.Config{ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{DBQueryTimeoutSeconds: 20}}))

	f(w, r)

	if !bytes.Equal(w.Body.Bytes(), expectedBody) {
		t.Errorf("received: %s\n expected: %s\n", w.Body.Bytes(), expectedBody)
	}

	w, r = newRWPair(t, nil)

	f(w, r)

	expectedError := `{"alerts":[{"text":"unauthorized, please log in.","level":"error"}]}` + "\n"

	if *debugLogging {
		fmt.Printf("received: %s\n expected: %s\n", w.Body.Bytes(), expectedError)
	}

	if !bytes.Equal(w.Body.Bytes(), []byte(expectedError)) {
		t.Errorf("received: %s\n expected: %s\n", w.Body.Bytes(), expectedError)
	}
}

func TestRequiredPermissionsMiddleware(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	userName := "user1"
	secret := "secret"

	rows := sqlmock.NewRows([]string{"priv_level", "username", "id", "tenant_id", "capabilities"})
	rows.AddRow(30, userName, 1, 1, "{foo}")
	mock.ExpectQuery("SELECT").WithArgs(userName).WillReturnRows(rows)

	authBase := AuthBase{secret, nil}

	cookie := tocookie.GetCookie(userName, time.Minute, secret)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success\n"))
	}

	authWrapper := authBase.GetWrapper(0)

	f := authWrapper(WrapHeaders(RequiredPermissionsMiddleware([]string{"foo"})(handler)))

	w, r := newRWPair(t, cookie)

	dbctx := context.WithValue(context.Background(), api.DBContextKey, db)
	r = r.WithContext(dbctx)
	conf := config.Config{
		ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{
			DBQueryTimeoutSeconds: 20,
		},
		RoleBasedPermissions: true,
	}
	r = r.WithContext(context.WithValue(r.Context(), api.ConfigContextKey, &conf))

	f(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected a 200 OK response when the user had all the required Permissions, got: %d", w.Code)
	}

	w, r = newRWPair(t, cookie)
	r = r.WithContext(dbctx)
	r = r.WithContext(context.WithValue(r.Context(), api.ConfigContextKey, &conf))
	rows = sqlmock.NewRows([]string{"priv_level", "username", "id", "tenant_id", "capabilities"})
	rows.AddRow(30, "user1", 1, 1, "{}")
	mock.ExpectQuery("SELECT").WithArgs(userName).WillReturnRows(rows)

	f(w, r)

	result := w.Result()
	if result.StatusCode != http.StatusForbidden {
		t.Errorf("Expected a 403 Forbidden response when the user was missing the required Permissions, got: %d", result.StatusCode)
	}
	rawResp, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	var alerts tc.Alerts
	if err := json.Unmarshal(rawResp, &alerts); err != nil {
		t.Errorf("Failed to read response recorder body: %v", err)
	}
	if !strings.Contains(alerts.ErrorString(), "foo") {
		t.Errorf("Expected an error-level alert mentioning the missing Permission, got: %s", alerts.ErrorString())
	}
}

func TestConfigRoleBasedPermissionsHandling(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	userName := "user1"
	secret := "secret"
	var rows *sqlmock.Rows
	resetRows := func(privLevel int, caps ...string) {
		rows = sqlmock.NewRows([]string{"priv_level", "username", "id", "tenant_id", "capabilities"})
		rows.AddRow(privLevel, userName, 1, 1, fmt.Sprintf("{%s}", strings.Join(caps, ",")))
		mock.ExpectQuery("SELECT").WithArgs(userName).WillReturnRows(rows)
	}
	resetRows(3, "foo")
	cookie := tocookie.GetCookie(userName, time.Minute, secret)

	conf := config.Config{
		ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{
			DBQueryTimeoutSeconds: 20,
		},
		RoleBasedPermissions: false,
	}
	ctx := context.WithValue(context.Background(), api.DBContextKey, db)
	ctx = context.WithValue(ctx, api.ConfigContextKey, &conf)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("successs\n"))
	}
	f := WrapHeaders(AuthBase{secret, nil}.GetWrapper(5)(RequiredPermissionsMiddleware([]string{"foo"})(handler)))

	w, r := newRWPair(t, cookie)
	r = r.WithContext(ctx)
	resetRows(3, "foo")

	f(w, r)
	result := w.Result()
	if result.StatusCode != http.StatusForbidden {
		t.Errorf("Expected a 403 Forbidden response when the user has insufficient PrivLevel and RoleBasedPermissions is configured to false; got: %d", result.StatusCode)
	}

	conf.RoleBasedPermissions = true
	w, r = newRWPair(t, cookie)
	r = r.WithContext(ctx)

	f(w, r)
	result = w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected a user with the right Permissions for an endpoint to get a 200 OK response regardless of PrivLevel when RoleBasedPermissions is configured to true; got: %d", result.StatusCode)
	}

	resetRows(30)
	w, r = newRWPair(t, cookie)
	r = r.WithContext(ctx)

	f(w, r)
	result = w.Result()
	if result.StatusCode != http.StatusForbidden {
		t.Errorf("Expected a user with the wrong Permissions for an endpoint to get a 403 Forbidden response regardless of PrivLevel, got: %d", result.StatusCode)
	}
	var alerts tc.Alerts
	if err := json.NewDecoder(result.Body).Decode(&alerts); err != nil {
		t.Errorf("Failed to read and decode response body: %v", err)
	} else {
		errStr := alerts.ErrorString()
		if !strings.Contains(errStr, "foo") {
			t.Errorf("Expected the reason the user was denied access to be a missing Permission, actual: %s", errStr)
		}
	}

	conf.RoleBasedPermissions = false
	w, r = newRWPair(t, cookie)
	r = r.WithContext(ctx)
	resetRows(30)
	f(w, r)
	result = w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected a user with the right PrivLevel for an endpoint to get a 200 OK response regardless of Permissions when RoleBasedPermissions is configured to false; got: %d", result.StatusCode)
	}
}

func TestNoOpWhenNoPermissionsRequired(t *testing.T) {
	respBts := "success"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(respBts))
	}
	f := RequiredPermissionsMiddleware([]string{})(handler)

	w, r := newRWPair(t, nil)
	f(w, r)
	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected checking for required Permissions to be a no-op when no Permissions are required, but response had status code %d", result.StatusCode)
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	} else if string(body) != respBts {
		t.Errorf("Expected normal response '%s' from endpoint, but got: %s", respBts, string(body))
	}

	f = RequiredPermissionsMiddleware(nil)(handler)

	w, r = newRWPair(t, nil)
	f(w, r)
	result = w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected checking for required Permissions to be a no-op when nil Permissions are required, but response had status code %d", result.StatusCode)
	}
	body, err = ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	} else if string(body) != respBts {
		t.Errorf("Expected normal response '%s' from endpoint, but got: %s", respBts, string(body))
	}
}

func TestGetCookieToken(t *testing.T) {
	var cookies []http.Cookie
	var jwtToken jwt.Token
	var jwtSigned []byte

	authUser := "foobar"
	httpCookie := tocookie.GetCookie(authUser, 0, "fOObAR.")

	jwtToken, _ = jwt.NewBuilder().Claim(api.MojoCookie, httpCookie.Value).Build()
	jwtSigned, _ = jwt.Sign(jwtToken, jwa.HS256, []byte("fOObAR."))

	mojoCookie := http.Cookie{Name: httpCookie.Name, Value: httpCookie.Value}
	accessToken := http.Cookie{Name: "access_token", Value: string(jwtSigned)}
	bearerToken := "Bearer " + string(jwtSigned)
	cookies = append(cookies, mojoCookie, accessToken, http.Cookie{})

	getUserFromCookie := func(cookieToken string) {
		secret := "fOObAR."
		user := ""
		cookie, userErr, sysErr := tocookie.Parse(secret, cookieToken)
		if userErr == nil && sysErr == nil {
			user = cookie.AuthData
		}
		if user != "foobar" {
			t.Errorf("Error: Unable to parse user from cookie. Expected: %v Got: %v", authUser, user)
		}
	}

	r, err := http.NewRequest("GET", "https://localhost:8888", nil)
	if err == nil && r != nil {
		for i := range cookies {
			if cookies[i].Name != "" {
				r.AddCookie(&cookies[i])
				cookieToken := getCookieToken(r)
				getUserFromCookie(cookieToken)
			} else {
				r.Header.Add("Authorization", bearerToken)
				cookieToken := getCookieToken(r)
				getUserFromCookie(cookieToken)
			}
		}
	}
}
