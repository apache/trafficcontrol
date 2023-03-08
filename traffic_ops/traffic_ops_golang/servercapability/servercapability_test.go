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
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
	"time"
)

func getTestSCs() []tc.ServerCapabilityV4 {
	scs := []tc.ServerCapabilityV4{}
	testSC := tc.ServerCapabilityV4{
		ServerCapability: tc.ServerCapability{
			Name:        "test",
			LastUpdated: &tc.TimeNoMod{Time: time.Now()},
		},
		Description: "test servers",
	}
	scs = append(scs, testSC)

	testSC1 := testSC
	testSC1.Name = "blah"
	testSC1.Description = "blah servers"
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
	fmt.Println(testSCs)
	rows := sqlmock.NewRows([]string{"name", "last_updated", "description"})

	for _, ts := range testSCs {
		rows = rows.AddRow(
			ts.ServerCapability.Name,
			ts.ServerCapability.LastUpdated,
			ts.Description)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	reqInfo := api.APIInfo{Tx: db.MustBegin(), Params: map[string]string{"name": "test"}}
	obj := TOServerCapability{
		APIInfoImpl: api.APIInfoImpl{ReqInfo: &reqInfo},
		//ServerCapability: tc.ServerCapability{},
	}

	scs, userErr, sysErr, _, _ := obj.Read(nil, false)
	if userErr != nil || sysErr != nil {
		t.Errorf("Read expected: no errors, actual: %v %v", userErr, sysErr)
	}

	if len(scs) != 1 {
		t.Errorf("Server Capability.Read expected: len(sc) == 1, actual: %v", len(scs))
	}
}
