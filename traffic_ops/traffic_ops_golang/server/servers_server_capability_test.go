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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"strings"
	"testing"
	"time"
)

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
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
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

	reqInfo := api.APIInfo{Tx: db.MustBegin(), Params: map[string]string{"serverId": "1"}}
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

func TestAssignMultipleServersCapabilities(t *testing.T) {

}

func TestDeleteMultipleServersCapabilities(t *testing.T) {

}
