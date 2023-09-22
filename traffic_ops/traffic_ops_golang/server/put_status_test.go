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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestInvalidStatusForDeliveryServicesAlertText(t *testing.T) {
	type testStruct struct {
		dsId     []int
		expected string
	}

	var testData = []testStruct{
		{[]int{0}, " #0 with no 'ONLINE' or 'REPORTED' EDGE servers"},
		{[]int{0, 1}, "s #0 and #1 with no 'ONLINE' or 'REPORTED' EDGE servers"},
		{[]int{0, 1, 2}, "s #0, #1, and #2 with no 'ONLINE' or 'REPORTED' EDGE servers"},
	}

	for i, _ := range testData {
		desc := InvalidStatusForDeliveryServicesAlertText("", "EDGE", testData[i].dsId)
		if testData[i].expected != desc {
			t.Errorf("strings don't't match, got:%s; expected:%s", desc, testData[i].expected)
		}
	}
}

func TestCheckExistingStatusInfo(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	lastUpdated := time.Now()
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"status", "status_last_updated"})
	rows.AddRow(1, lastUpdated)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	status, statusLastUpdated := checkExistingStatusInfo(1, db.MustBegin().Tx)
	if status != 1 {
		t.Errorf("Expected server status to be 1, got %v", status)
	}

	if statusLastUpdated != lastUpdated {
		t.Errorf("Expected status time: %s, got: %s", lastUpdated, statusLastUpdated)
	}
}

func TestUpdateServerStatusAndOfflineReason(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	lastUpdated := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(2, "no longer needed", lastUpdated, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	reason := util.Ptr("no longer needed")
	err = updateServerStatusAndOfflineReason(2, 2, 1, lastUpdated, reason, db.MustBegin().Tx)
	if err != nil {
		t.Errorf("unable to change the status of the server, error: %s", err)
	}
}
