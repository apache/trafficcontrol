package cdnfederation

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
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault/backends/disabled"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAddTenancyStmt(t *testing.T) {
	where := ""
	output := addTenancyStmt(where)
	if expected := "WHERE ds.tenant_id = ANY(:tenantIDs)"; output != expected {
		t.Errorf("Incorrect statement from blank WHERE; want: '%s', got: '%s", expected, output)
	}
	where = "WHERE cname=:cname"
	output = addTenancyStmt(where)
	if expected := "WHERE cname=:cname AND ds.tenant_id = ANY(:tenantIDs)"; output != expected {
		t.Errorf("Incorrect statement from blank WHERE; want: '%s', got: '%s", expected, output)
	}
}

func TestParamColumnInfo(t *testing.T) {
	params := paramColumnInfo(api.Version{Major: 4})
	if l := len(params); l != 3 {
		t.Errorf("Expected versions prior to 5 to support 3 query params, found support for: %d", l)
	}
	for _, param := range [3]string{"cname", "id", "name"} {
		if _, ok := params[param]; !ok {
			t.Errorf("Expected versions prior to 5 to support the '%s' query param, but support for such wasn't found", param)
		}
	}

	params = paramColumnInfo(api.Version{Major: 5})
	if l := len(params); l != 5 {
		t.Errorf("Expected versions 5 and later to support 5 query params, found support for: %d", l)
	}
	for _, param := range [5]string{"cname", "dsID", "id", "name", "xmlID"} {
		if _, ok := params[param]; !ok {
			t.Errorf("Expected versions 5 and later to support the '%s' query param, but support for such wasn't found", param)
		}
	}
}

func getMockTx(t *testing.T) (sqlmock.Sqlmock, *sqlx.Tx, *sqlx.DB) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	mock.ExpectBegin()

	return mock, db.MustBegin(), db
}

func cleanup(t *testing.T, mock sqlmock.Sqlmock, db *sqlx.DB) {
	t.Helper()
	mock.ExpectClose()
	db.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %v", err)
	}
}

func gettingUserTenantListFails(t *testing.T) {
	mock, tx, db := getMockTx(t)
	defer cleanup(t, mock, db)
	defer func() {
		mock.ExpectRollback()
		tx.Rollback()
	}()

	err := errors.New("unknown failure")
	mock.ExpectQuery("WITH RECURSIVE").WillReturnError(err)
	_, _, code, userErr, sysErr := getCDNFederations(&api.Info{Tx: tx, User: &auth.CurrentUser{TenantID: 1}})

	if code != http.StatusInternalServerError {
		t.Errorf("Incorrect response code when getting user tenants fails; want: %d, got: %d", http.StatusInternalServerError, code)
	}
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
	if sysErr == nil {
		t.Fatal("Expected a system error but didn't get one")
	}

	// You can't use `errors.Is` here because sqlmock doesn't wrap the error you
	// give it, so we have to resort to comparing text and praying there's no
	// weird coincidence going on behind the scenes.
	if !strings.Contains(sysErr.Error(), err.Error()) {
		t.Errorf("Incorrect system error returned; want: %v, got: %v", err, sysErr)
	}
}

func buildingQueryPartsFails(t *testing.T) {
	mock, tx, db := getMockTx(t)
	defer cleanup(t, mock, db)
	defer func() {
		mock.ExpectRollback()
		tx.Rollback()
	}()

	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)

	mock.ExpectQuery("WITH RECURSIVE").WillReturnRows(rows)

	inf := api.Info{
		Params: map[string]string{
			"dsID": "not an integer",
		},
		Tx:      tx,
		User:    &auth.CurrentUser{TenantID: 1},
		Version: &api.Version{Major: 5},
	}
	_, _, code, userErr, sysErr := getCDNFederations(&inf)
	if code != http.StatusBadRequest {
		t.Errorf("Incorrect response code when getting user tenants fails; want: %d, got: %d", http.StatusBadRequest, code)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}
	if userErr == nil {
		t.Fatal("Expected a user-facing error, but didn't get one")
	}
	if !strings.Contains(userErr.Error(), "dsID") {
		t.Errorf("Incorrect user error; expected it to mention 'dsID', but it didn't: %v", userErr)
	}
}

func everythingWorks(t *testing.T) {
	mock, tx, db := getMockTx(t)
	defer cleanup(t, mock, db)
	defer func() {
		mock.ExpectCommit()
		tx.Commit()
	}()

	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)

	mock.ExpectQuery("WITH RECURSIVE").WillReturnRows(rows)

	fedRows := sqlmock.NewRows([]string{"tenant_id", "id", "cname", "ttl", "description", "last_updated", "ds_id", "xml_id"})
	fed := tc.CDNFederationV5{
		CName:       "test.quest.",
		Description: util.Ptr("a non-blank description"),
		DeliveryService: &tc.CDNFederationDeliveryService{
			ID:    1,
			XMLID: "some-xmlid",
		},
		ID:          1,
		LastUpdated: time.Time{}.Add(time.Hour),
		TTL:         5,
	}

	fedRows.AddRow(1, fed.ID, fed.CName, fed.TTL, fed.Description, fed.LastUpdated, fed.DeliveryService.ID, fed.DeliveryService.XMLID)
	mock.ExpectQuery("SELECT").WillReturnRows(fedRows)

	feds, _, _, userErr, sysErr := getCDNFederations(&api.Info{Tx: tx, User: &auth.CurrentUser{TenantID: 1}, Version: &api.Version{Major: 5}})
	if userErr != nil {
		t.Errorf("Unexpected user-facing error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error: %v", sysErr)
	}
	if l := len(feds); l != 1 {
		t.Fatalf("Expected one federation to be returned; got: %d", l)
	}
	if !reflect.DeepEqual(feds[0], fed) {
		t.Errorf("expected a federation like '%#v', but instead found: %#v", fed, feds[0])
	}
}

func TestGetCDNFederations(t *testing.T) {
	t.Run("getting user Tenant list fails", gettingUserTenantListFails)
	t.Run("building where/orderby/pagination fails", buildingQueryPartsFails)
	t.Run("everything works", everythingWorks)
}

func wrapContext(r *http.Request, key any, value any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

func testingInf(t *testing.T, body []byte) (*http.Request, sqlmock.Sqlmock, *sqlx.DB) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	r := httptest.NewRequest(http.MethodPost, "/api/5.0/cdns/ALL/federations", bytes.NewReader(body))
	r = wrapContext(r, api.ConfigContextKey, &config.Config{ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{DBQueryTimeoutSeconds: 1000}})
	r = wrapContext(r, api.DBContextKey, db)
	r = wrapContext(r, api.TrafficVaultContextKey, &disabled.Disabled{})
	r = wrapContext(r, api.ReqIDContextKey, uint64(0))
	r = wrapContext(r, auth.CurrentUserKey, auth.CurrentUser{})
	r = wrapContext(r, api.PathParamsKey, make(map[string]string))

	mock.ExpectBegin()

	return r, mock, db
}

func TestCreate(t *testing.T) {
	newFed := tc.CDNFederationV5{
		CName:       "test.quest.",
		TTL:         420,
		Description: nil,
	}
	bts, err := json.Marshal(newFed)
	if err != nil {
		t.Fatalf("marshaling testing request body: %v", err)
	}

	newFed.ID = 1
	newFed.LastUpdated = time.Time{}.Add(time.Hour)

	r, mock, db := testingInf(t, bts)
	defer cleanup(t, mock, db)

	rows := sqlmock.NewRows([]string{"id", "last_updated"})
	rows.AddRow(newFed.ID, newFed.LastUpdated)
	mock.ExpectQuery("INSERT").WillReturnRows(rows)

	f := api.Wrap(Create, nil, nil)
	w := httptest.NewRecorder()
	f(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("Incorrect response code; want: %d, got: %d", http.StatusCreated, w.Code)
	}

	var created tc.CDNFederationV5Response
	err = json.Unmarshal(w.Body.Bytes(), &created)
	if err != nil {
		t.Fatalf("Unmarshaling response: %v", err)
	}

	if created.Response != newFed {
		t.Errorf("Didn't create the expected Federation; want: %#v, got: %#v", newFed, created.Response)
	}
}

func TestValidate(t *testing.T) {
	fed := tc.CDNFederationV5{
		CName: "test.quest.",
		TTL:   60,
	}
	err := validate(fed)
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}

	fed.TTL--
	err = validate(fed)
	if err == nil {
		t.Fatal("Expected an error for TTL below minimum, but didn't get one")
	}
	if !strings.Contains(err.Error(), "ttl") {
		t.Errorf("Expected error message to mention 'ttl': %v", err)
	}

	fed.TTL = 60
	fed.CName = "test.quest"
	err = validate(fed)
	if err == nil {
		t.Fatal("Expected an error for a CNAME without a terminating '.', but didn't get one")
	}
	if !strings.Contains(err.Error(), "cname") {
		t.Errorf("Expected error message to mention 'cname': %v", err)
	}

}
