package main

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
	"reflect"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetCfgDiffs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hostName := "myedge"

	timestamp := time.Now().UTC().String()
	cfgFileDiffs1 := CfgFileDiffs{
		FileName: "TestFile.cfg",
		DBLinesMissing: []string{ "db_line_missing1", "db_line_missing2", },
		DiskLinesMissing: []string{ "disk_line_missing1", "disk_line_missing2", },
		ReportTimestamp: timestamp,
	}

	rows := sqlmock.NewRows([]string{"config_name", "db_lines_missing", "disk_lines_missing", "last_checked", })

	dbLinesMissingJson, err := json.Marshal(cfgFileDiffs1.DBLinesMissing)
	diskLinesMissingJson, err := json.Marshal(cfgFileDiffs1.DiskLinesMissing)
	rows = rows.AddRow(cfgFileDiffs1.FileName, dbLinesMissingJson, diskLinesMissingJson, cfgFileDiffs1.ReportTimestamp)
	
	
	mock.ExpectQuery("SELECT").WithArgs(hostName).WillReturnRows(rows)

	cfgFileDiffs, err := getCfgDiffs(db, hostName)
	if err != nil {
		t.Errorf("getCfgDiffs expected: nil error, actual: %v", err)
	}

	if len(cfgFileDiffs) != 1 {
		t.Errorf("getCfgDiffs expected: len(cfgFileDiffs) == 1, actual: %v", len(cfgFileDiffs))
	}
	sqlCfgFileDiffs := cfgFileDiffs[0]
	if !reflect.DeepEqual(sqlCfgFileDiffs, cfgFileDiffs1) {
		t.Errorf("getCfgDiffs expected: cfgFileDiffs == %+v, actual: %+v", cfgFileDiffs1, sqlCfgFileDiffs)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetCfgDiffsJson(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hostName := "myedge"

	timestamp := time.Now().UTC().String()
	cfgFileDiffsResponse := CfgFileDiffsResponse{
		Response: []CfgFileDiffs {{
			FileName: "TestFile.cfg",
			DBLinesMissing: []string{ "db_line_missing1", "db_line_missing2", },
			DiskLinesMissing: []string{ "disk_line_missing1", "disk_line_missing2", },
			ReportTimestamp: timestamp,
		},},
	}

	rows := sqlmock.NewRows([]string{"config_name", "db_lines_missing", "disk_lines_missing", "last_checked", })

	dbLinesMissingJson, err := json.Marshal(cfgFileDiffsResponse.Response[0].DBLinesMissing)
	diskLinesMissingJson, err := json.Marshal(cfgFileDiffsResponse.Response[0].DiskLinesMissing)
	rows = rows.AddRow(cfgFileDiffsResponse.Response[0].FileName, dbLinesMissingJson, diskLinesMissingJson, cfgFileDiffsResponse.Response[0].ReportTimestamp)
	
	
	mock.ExpectQuery("SELECT").WithArgs(hostName).WillReturnRows(rows)

	cfgFileDiffsResponseT, err := getCfgDiffsJson(hostName, db)
	if err != nil {
		t.Errorf("getCfgDiffs expected: nil error, actual: %v", err)
	}

	if len(cfgFileDiffsResponseT.Response) != 1 {
		t.Errorf("getCfgDiffsJson expected: len(cfgFileDiffsResponseT.Response) == 1, actual: %v", len(cfgFileDiffsResponseT.Response))
	}
	
	if !reflect.DeepEqual(*cfgFileDiffsResponseT, cfgFileDiffsResponse) {
		t.Errorf("getCfgDiffsJson expected: cfgFileDiffsResponseT == %+v, actual: %+v", cfgFileDiffsResponseT, cfgFileDiffsResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}