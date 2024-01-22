package server

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
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestSSCsV5() []tc.ServerServerCapabilityV5 {
	sscs := []tc.ServerServerCapabilityV5{}
	testSSC := tc.ServerServerCapabilityV5{
		LastUpdated:      util.Ptr(time.Now()),
		Server:           util.StrPtr("test"),
		ServerID:         util.IntPtr(1),
		ServerCapability: util.StrPtr("test"),
	}
	sscs = append(sscs, testSSC)

	testSSC1 := testSSC
	testSSC1.ServerCapability = util.Ptr("blah")
	sscs = append(sscs, testSSC1)

	return sscs
}

func TestReadSCsV5(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testSCs := getTestSSCsV5()
	rows := sqlmock.NewRows([]string{"server_capability", "server", "last_updated"})

	for _, ts := range testSCs {
		rows = rows.AddRow(
			ts.ServerCapability,
			ts.ServerID,
			ts.LastUpdated)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	reqInfo := api.Info{Tx: db.MustBegin(), Params: map[string]string{"serverId": "1"}}
	obj := TOServerServerCapabilityV5{
		api.APIInfoImpl{ReqInfo: &reqInfo},
		tc.ServerServerCapabilityV5{},
	}
	sscs, userErr, sysErr, _, _ := obj.Read(nil, false)
	if userErr != nil || sysErr != nil {
		t.Errorf("Read expected: no errors, actual: %v %v", userErr, sysErr)
	}

	if len(sscs) != 2 {
		t.Errorf("ServerServerCapabilityV5.Read expected: len(scs) == 1, actual: %v", len(sscs))
	}
}

func TestInterfacesV5(t *testing.T) {
	var i interface{}
	i = &TOServerServerCapabilityV5{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("ServerServerCapabilityV5 must be Creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("ServerServerCapabilityV5 must be Reader")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("ServerServerCapabilityV5 must be Deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("ServerServerCapabilityV5 must be Identifier")
	}
}

func TestFuncsV5(t *testing.T) {
	if strings.Index(scSelectQuery(), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(scInsertQuery(), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(scDeleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}
}

func TestValidateV5(t *testing.T) {
	testSSC := tc.ServerServerCapabilityV5{
		LastUpdated:      util.Ptr(time.Now()),
		Server:           util.StrPtr("test1"),
		ServerID:         util.IntPtr(1),
		ServerCapability: util.StrPtr("abc"),
	}
	testTOSSC := TOServerServerCapabilityV5{
		ServerServerCapabilityV5: testSSC,
	}

	err, _ := testTOSSC.Validate()
	errs := test.SortErrors(test.SplitErrors(err))

	if len(errs) > 0 {
		t.Errorf(`expected no errors,  got %v`, errs)
	}
}

func TestCheckExistingServerV5(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"host_name"})
	rows.AddRow("test")
	mock.ExpectQuery("SELECT host_name").WithArgs(1).WillReturnRows(rows)

	rows1 := sqlmock.NewRows([]string{"name"})
	rows1.AddRow("ALL")
	mock.ExpectQuery("SELECT name").WithArgs(1).WillReturnRows(rows1)

	rows2 := sqlmock.NewRows([]string{"username", "soft", "shared_usernames"})
	rows2.AddRow("user1", false, []byte("{}"))
	mock.ExpectQuery("SELECT c.username, c.soft").WithArgs("ALL").WillReturnRows(rows2)
	mock.ExpectCommit()

	testSCCs := getTestSSCsV5()
	var sids []int64
	sids = append(sids, int64(*testSCCs[0].ServerID))
	code, usrErr, sysErr := checkExistingServer(db.MustBegin().Tx, sids, "user1")
	if usrErr != nil {
		t.Errorf("server not found, error:%v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to check if server exists, error:%v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("existing server check failed, expected:%d, got:%d", http.StatusOK, code)
	}
}

func TestCheckServerTypeV5(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testSCCs := getTestSSCsV5()
	testSCCs[1].ServerID = util.Ptr(2)
	testSCCs[1].Server = util.Ptr("foo")

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"array_agg"})
	var sids []int64
	for i, _ := range testSCCs {
		sids = append(sids, int64(*testSCCs[i].ServerID))
	}
	rows.AddRow([]byte("{1,2}"))
	mock.ExpectQuery("SELECT array_agg").WithArgs(pq.Array(sids)).WillReturnRows(rows)
	mock.ExpectCommit()

	code, usrErr, sysErr := checkServerType(db.MustBegin().Tx, sids)
	if usrErr != nil {
		t.Errorf("mismatch in server type, error:%v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to check if server type exists, error:%v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("server type check failed, expected:%d, got:%d", http.StatusOK, code)
	}
}

func getTestSSCs() []tc.ServerServerCapability {
	sscs := []tc.ServerServerCapability{}
	testSSC := tc.ServerServerCapability{
		LastUpdated:      &tc.TimeNoMod{Time: time.Now()},
		Server:           util.StrPtr("test"),
		ServerID:         util.IntPtr(1),
		ServerCapability: util.StrPtr("test"),
	}
	sscs = append(sscs, testSSC)

	testSSC1 := testSSC
	testSSC1.ServerCapability = util.Ptr("blah")
	sscs = append(sscs, testSSC1)

	return sscs
}

func TestReadSCs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testSCs := getTestSSCs()
	rows := sqlmock.NewRows([]string{"server_capability", "server", "last_updated"})

	for _, ts := range testSCs {
		rows = rows.AddRow(
			ts.ServerCapability,
			ts.ServerID,
			ts.LastUpdated)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	reqInfo := api.Info{Tx: db.MustBegin(), Params: map[string]string{"serverId": "1"}}
	obj := TOServerServerCapability{
		api.APIInfoImpl{ReqInfo: &reqInfo},
		tc.ServerServerCapability{},
	}
	sscs, userErr, sysErr, _, _ := obj.Read(nil, false)
	if userErr != nil || sysErr != nil {
		t.Errorf("Read expected: no errors, actual: %v %v", userErr, sysErr)
	}

	if len(sscs) != 2 {
		t.Errorf("ServerServerCapability.Read expected: len(scs) == 1, actual: %v", len(sscs))
	}
}

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOServerServerCapability{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("ServerServerCapability must be Creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("ServerServerCapability must be Reader")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("ServerServerCapability must be Deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("ServerServerCapability must be Identifier")
	}
}

func TestFuncs(t *testing.T) {
	if strings.Index(scSelectQuery(), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(scInsertQuery(), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(scDeleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}
}

func TestValidate(t *testing.T) {
	testSSC := tc.ServerServerCapability{
		LastUpdated:      &tc.TimeNoMod{Time: time.Now()},
		Server:           util.StrPtr("test1"),
		ServerID:         util.IntPtr(1),
		ServerCapability: util.StrPtr("abc"),
	}
	testTOSSC := TOServerServerCapability{
		ServerServerCapability: testSSC,
	}

	err, _ := testTOSSC.Validate()
	errs := test.SortErrors(test.SplitErrors(err))

	if len(errs) > 0 {
		t.Errorf(`expected no errors,  got %v`, errs)
	}
}

func TestCheckExistingServer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"host_name"})
	rows.AddRow("test")
	mock.ExpectQuery("SELECT host_name").WithArgs(1).WillReturnRows(rows)

	rows1 := sqlmock.NewRows([]string{"name"})
	rows1.AddRow("ALL")
	mock.ExpectQuery("SELECT name").WithArgs(1).WillReturnRows(rows1)

	rows2 := sqlmock.NewRows([]string{"username", "soft", "shared_usernames"})
	rows2.AddRow("user1", false, []byte("{}"))
	mock.ExpectQuery("SELECT c.username, c.soft").WithArgs("ALL").WillReturnRows(rows2)
	mock.ExpectCommit()

	testSCCs := getTestSSCs()
	var sids []int64
	sids = append(sids, int64(*testSCCs[0].ServerID))
	code, usrErr, sysErr := checkExistingServer(db.MustBegin().Tx, sids, "user1")
	if usrErr != nil {
		t.Errorf("server not found, error:%v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to check if server exists, error:%v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("existing server check failed, expected:%d, got:%d", http.StatusOK, code)
	}
}

func TestCheckServerType(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testSCCs := getTestSSCs()
	testSCCs[1].ServerID = util.Ptr(2)
	testSCCs[1].Server = util.Ptr("foo")

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"array_agg"})
	var sids []int64
	for i, _ := range testSCCs {
		sids = append(sids, int64(*testSCCs[i].ServerID))
	}
	rows.AddRow([]byte("{1,2}"))
	mock.ExpectQuery("SELECT array_agg").WithArgs(pq.Array(sids)).WillReturnRows(rows)
	mock.ExpectCommit()

	code, usrErr, sysErr := checkServerType(db.MustBegin().Tx, sids)
	if usrErr != nil {
		t.Errorf("mismatch in server type, error:%v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to check if server type exists, error:%v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("server type check failed, expected:%d, got:%d", http.StatusOK, code)
	}
}
