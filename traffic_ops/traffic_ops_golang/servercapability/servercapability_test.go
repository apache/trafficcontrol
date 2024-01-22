package servercapability

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
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestSCs() []tc.ServerCapability {
	scs := []tc.ServerCapability{}
	testSC := tc.ServerCapability{
		Name:        "test",
		LastUpdated: &tc.TimeNoMod{Time: time.Now()},
	}
	scs = append(scs, testSC)

	testSC1 := testSC
	testSC1.Name = "blah"
	scs = append(scs, testSC1)

	return scs
}

func TestReadSCs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testSCs := getTestSCs()
	rows := sqlmock.NewRows([]string{"name", "last_updated"})
	for _, ts := range testSCs {
		rows = rows.AddRow(
			ts.Name,
			ts.LastUpdated,
		)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	reqInfo := api.Info{Tx: db.MustBegin(), Params: map[string]string{"name": "test"}}
	obj := TOServerCapability{
		APIInfoImpl:      api.APIInfoImpl{ReqInfo: &reqInfo},
		ServerCapability: tc.ServerCapability{},
	}

	scs, userErr, sysErr, _, _ := obj.Read(nil, false)
	if userErr != nil || sysErr != nil {
		t.Errorf("Read expected: no errors, actual: %v %v", userErr, sysErr)
	}

	if len(scs) != 2 {
		t.Errorf("Server Capability.Read expected: len(sc) == 1, actual: %v", len(scs))
	}
}

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOServerCapability{}

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
	testTOSC := TOServerCapability{}
	if strings.Index(testTOSC.SelectQuery(), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(testTOSC.InsertQuery(), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(testTOSC.updateQuery(), "UPDATE") != 0 {
		t.Errorf("expected updateQuery to start with UPDATE")
	}
	if strings.Index(testTOSC.DeleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}
}

func TestValidate(t *testing.T) {
	var errs []error
	// Negative test case
	testSC := tc.ServerCapability{
		Name: "",
	}
	testTOSC := TOServerCapability{}
	testTOSC.ServerCapability = testSC
	err, _ := testTOSC.Validate()
	errs = test.SortErrors(test.SplitErrors(err))
	if len(errs) < 0 {
		t.Errorf(`expected errors: %v,  got no errors`, errs)
	}

	// Positive test case
	testSCs := getTestSCs()
	for _, val := range testSCs {
		testTOSC.ServerCapability = val
		err, _ := testTOSC.Validate()
		errs = test.SortErrors(test.SplitErrors(err))
	}
	if len(errs) > 0 {
		t.Errorf(`expected no errors,  got errors: %v`, errs)
	}
}
