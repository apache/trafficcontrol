package role

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
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault/backends/disabled"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func stringAddr(s string) *string {
	return &s
}
func intAddr(i int) *int {
	return &i
}

//removed sqlmock based ReadRoles test due to sqlmock / pq.Array() type incompatibility issue.

func TestFuncs(t *testing.T) {
	if strings.Index(selectQuery(), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(insertQuery(), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(updateQuery(), "UPDATE") != 0 {
		t.Errorf("expected updateQuery to start with UPDATE")
	}
	if strings.Index(deleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}

}
func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TORole{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("role must be creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("role must be reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("role must be updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("role must be deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("role must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	// invalid name, empty domainname
	n := "not_a_valid_role"
	reqInfo := api.Info{}
	role := tc.Role{}
	role.Name = &n
	r := TORole{
		APIInfoImpl: api.APIInfoImpl{ReqInfo: &reqInfo},
		Role:        role,
	}
	userErr, _ := r.Validate()
	errs := util.JoinErrsStr(test.SortErrors(test.SplitErrors(userErr)))

	expectedErrs := util.JoinErrsStr([]error{
		errors.New(`'description' cannot be blank`),
		errors.New(`'privLevel' is required`),
	})

	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}

	//  name,  domainname both valid
	role = tc.Role{}
	role.Name = stringAddr("this is a valid name")
	role.Description = stringAddr("this is a description")
	role.PrivLevel = intAddr(30)
	r = TORole{
		APIInfoImpl: api.APIInfoImpl{ReqInfo: &reqInfo},
		Role:        role,
	}
	userErr, sysErr := r.Validate()
	if userErr != nil {
		t.Errorf("expected nil user error, got: %s", userErr)
	}
	if sysErr != nil {
		t.Errorf("expected nil system error, got: %v", sysErr)
	}

}

func TestCreateWithEmptyPermissions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "", strings.NewReader(`{"name":"role", "description":"description"}`))
	if err != nil {
		t.Error("Error creating new request")
	}

	addRequestContext(r, db)

	columns := []string{"id", "last_updated"}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WithArgs("role", "description", auth.PrivLevelAdmin).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, time.Now()))
	mock.ExpectCommit()

	Create(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("Got status code %d but want 201", w.Code)
	}

	resp := struct {
		tc.Alerts
		Response tc.RoleV4 `json:"response"`
	}{}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Error decoding response")
	}

	if resp.Response.Permissions == nil {
		t.Error("Permissions should be empty not nil")
	}
}

func TestUpdateWithEmptyPermissions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "", strings.NewReader(`{"name":"new_role", "description":"new_description"}`))
	if err != nil {
		t.Error("Error creating new request")
	}

	addRequestContext(r, db)

	r.URL.RawQuery = "name=role"

	mock.ExpectBegin()

	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("(?i)SELECT").WillReturnRows(sqlmock.NewRows([]string{"last_updated"}).AddRow(time.Now()))
	mock.ExpectQuery("SET").WithArgs("new_role", "new_description", "role").WillReturnRows(sqlmock.NewRows([]string{"last_updated"}).AddRow(time.Now()))

	mock.ExpectExec("DELETE").WithArgs("new_role").WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectExec("INSERT").WithArgs(1, nil).WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	Update(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Got status code %d but want 200", w.Code)
	}

	resp := struct {
		tc.Alerts
		Response tc.RoleV4 `json:"response"`
	}{}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Error decoding response")
	}

	if resp.Response.Permissions == nil {
		t.Fatal("Permissions should be empty not nil")
	}
}

func addRequestContext(r *http.Request, db *sqlx.DB) {
	cfg := config.Config{ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{DBQueryTimeoutSeconds: 20}, UseIMS: true}
	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.CurrentUserKey,
		auth.CurrentUser{UserName: "username", ID: 1, PrivLevel: auth.PrivLevelAdmin})
	ctx = context.WithValue(ctx, api.PathParamsKey, map[string]string{"id": "1"})
	ctx = context.WithValue(ctx, api.DBContextKey, db)
	ctx = context.WithValue(ctx, api.ConfigContextKey, &cfg)
	ctx = context.WithValue(ctx, api.ReqIDContextKey, uint64(0))
	var tv trafficvault.TrafficVault = &disabled.Disabled{}
	ctx = context.WithValue(ctx, api.TrafficVaultContextKey, tv)

	*r = *r.WithContext(ctx)
}
