package login

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

import "context"
import "net/http"
import "net/http/httptest"
import "testing"
import "time"

import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tocookie"

import "github.com/jmoiron/sqlx"
import sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

var testUser = auth.CurrentUser{
	UserName:     "admin",
	ID:           1,
	PrivLevel:    30,
	TenantID:     1,
	Role:         1,
	Capabilities: nil,
}

func TestLogout(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{
		"r.priv_level",
		"r.id as role",
		"u.id",
		"u.username",
		"COALESCE(u.tenant_id, -1) AS tenant_id",
		"ARRAY(SELECT rc.cap_name FROM role_capability AS rc WHERE rc.role_id=r.id) AS capabilities",
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows(cols)
	rows.AddRow(
		testUser.PrivLevel,
		testUser.Role,
		testUser.ID,
		testUser.UserName,
		testUser.TenantID,
		testUser.Capabilities,
	)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	cookie := tocookie.GetCookie(testUser.UserName, 24*time.Hour, "secret")
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/1.4/logout", nil)
	if err != nil {
		t.Fatalf("Failed to create a request: %v", err)
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, api.DBContextKey, db)
	conf := config.Config{}
	conf.ConfigTrafficOpsGolang.DBQueryTimeoutSeconds = 100
	ctx = context.WithValue(ctx, api.ConfigContextKey, &conf)
	ctx = context.WithValue(ctx, api.ReqIDContextKey, uint64(1))
	ctx = context.WithValue(ctx, api.APIRespWrittenKey, false)
	ctx = context.WithValue(ctx, auth.CurrentUserKey, testUser)
	ctx = context.WithValue(ctx, api.PathParamsKey, map[string]string{})
	ctx, _ = context.WithDeadline(ctx, time.Now().Add(24*time.Hour))
	req = req.WithContext(ctx)

	req.AddCookie(cookie)
	LogoutHandler("test")(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected response code %d, got %d", http.StatusOK, rr.Code)
	}

	expected := `{"alerts":[{"text":"You are logged out.","level":"success"}]}
`
	if rr.Body.String() != expected {
		t.Errorf("Expected response body:\n\t%sbut got:\n\t%s", expected, rr.Body.String())
	}

	cookieFound := false
	for _, c := range rr.Result().Cookies() {
		if c.Name != tocookie.Name {
			continue
		}
		cookieFound = true

		if c.Path != "/" {
			t.Errorf("Expected cookie path to be '/', but got: %s", c.Path)
		}

		if !c.HttpOnly {
			t.Errorf("Expected cookie to be HTTP-only, but it wasn't")
		}

		if time.Second < time.Since(c.Expires) || -time.Second > time.Since(c.Expires) {
			t.Errorf("Expected cookie expiration to be within one second of now, but was %s", time.Since(c.Expires))
		}
		break

		parsedCookie, err := tocookie.Parse("secret", c.Value)
		if err != nil {
			t.Errorf("Failed to parse cookie value: %v", err)
			break
		}

		if parsedCookie.ExpiresUnix != c.Expires.Unix() {
			t.Errorf("Expected encoded expiration to be %d, but it was %d", c.Expires.Unix(), parsedCookie.ExpiresUnix)
		}

		if parsedCookie.AuthData != testUser.UserName {
			t.Errorf("Incorrect user parsed from cookie; expected '%s' but got: %s", testUser.UserName, parsedCookie.AuthData)
		}
	}

	if !cookieFound {
		t.Errorf("Expected handler to set the '%s' cookie, but it didn't", tocookie.Name)
	}
}
