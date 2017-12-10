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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

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
		FileName:         "TestFile.cfg",
		DBLinesMissing:   []string{"db_line_missing1", "db_line_missing2"},
		DiskLinesMissing: []string{"disk_line_missing1", "disk_line_missing2"},
		ReportTimestamp:  timestamp,
	}

	rows := sqlmock.NewRows([]string{"config_name", "db_lines_missing", "disk_lines_missing", "last_checked"})

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

	var db *sqlx.DB = nil
	hostName := "myedge"

	timestamp := time.Now().UTC().String()
	cfgFileDiffsResponse := CfgFileDiffsResponse{
		Response: []CfgFileDiffs{{
			FileName:         "TestFile.cfg",
			DBLinesMissing:   []string{"db_line_missing1", "db_line_missing2"},
			DiskLinesMissing: []string{"disk_line_missing1", "disk_line_missing2"},
			ReportTimestamp:  timestamp,
		}},
	}

	// Test successful request
	cfgFileDiffsResponseT, err := getCfgDiffsJson(hostName, db,
		func(db *sqlx.DB, hostName string) ([]CfgFileDiffs, error) {
			return cfgFileDiffsResponse.Response, nil
		})

	if err != nil {
		t.Errorf("getCfgDiffs expected: nil error, actual: %v", err)
	}

	if len(cfgFileDiffsResponseT.Response) != 1 {
		t.Errorf("getCfgDiffsJson expected: len(cfgFileDiffsResponseT.Response) == 1, actual: %v", len(cfgFileDiffsResponseT.Response))
	}

	if !reflect.DeepEqual(*cfgFileDiffsResponseT, cfgFileDiffsResponse) {
		t.Errorf("getCfgDiffsJson expected: cfgFileDiffsResponseT == %+v, actual: %+v", cfgFileDiffsResponseT, cfgFileDiffsResponse)
	}

	// Test error case
	cfgFileDiffsResponseT, err = getCfgDiffsJson(hostName, db,
		func(db *sqlx.DB, hostName string) ([]CfgFileDiffs, error) {
			return nil, fmt.Errorf("Intentional Error for testing")
		})

	if err == nil {
		t.Errorf("getCfgDiffsJson expected: non-nil error, actual: nil")
	}

	if cfgFileDiffsResponseT != nil {
		t.Errorf("getCfgFileDiffsJson expected: nil response, actual: %v", cfgFileDiffsResponseT)
	}
}

func TestServerExists(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hostName := "myedge"

	// Test Expecting True Response
	rows := sqlmock.NewRows([]string{"host_name"}).AddRow("true")

	mock.ExpectQuery("SELECT EXISTS").WithArgs(hostName).WillReturnRows(rows)

	result, err := serverExists(db, hostName)
	if err != nil {
		t.Errorf("serverExists expected: nil error, actual: %v", err)
	}

	if result != true {
		t.Errorf("serverExists expected: result == true, actual: %v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	// Test Expecting False Response
	rows = sqlmock.NewRows([]string{"host_name"}).AddRow("false")

	mock.ExpectQuery("SELECT EXISTS").WithArgs(hostName).WillReturnRows(rows)

	result, err = serverExists(db, hostName)
	if err != nil {
		t.Errorf("serverExists expected: nil error, actual: %v", err)
	}

	if result != false {
		t.Errorf("serverExists expected: result == false, actual: %v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestInsertCfgDiffs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hostName := "myedge"
	timestamp := time.Now().UTC().String()

	cfgFileDiffs := CfgFileDiffs{
		FileName:         "TestFile.cfg",
		DBLinesMissing:   []string{"db_line_missing1", "db_line_missing2"},
		DiskLinesMissing: []string{"disk_line_missing1", "disk_line_missing2"},
		ReportTimestamp:  timestamp,
	}

	// Since "insertCfgDiffs" Marshals the json, we must store the unmarshalled json here.
	//		This will need to be updated if the above text gets changed
	dbLinesMissingJson := []uint8{91, 34, 100, 98, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 49, 34, 44, 34, 100, 98, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 50, 34, 93}
	diskLinesMissingJson := []uint8{91, 34, 100, 105, 115, 107, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 49, 34, 44, 34, 100, 105, 115, 107, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 50, 34, 93}

	mock.ExpectExec("INSERT INTO").WithArgs(
		hostName,
		cfgFileDiffs.FileName,
		dbLinesMissingJson,
		diskLinesMissingJson,
		cfgFileDiffs.ReportTimestamp).WillReturnResult(sqlmock.NewResult(1, 1))

	err = insertCfgDiffs(db, hostName, cfgFileDiffs)
	if err != nil {
		t.Errorf("insertCfgDiffs expected: nil error, actual: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestUpdateCfgDiiffs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hostName := "myedge"
	timestamp := time.Now().UTC().String()

	cfgFileDiffs := CfgFileDiffs{
		FileName:         "TestFile.cfg",
		DBLinesMissing:   []string{"db_line_missing1", "db_line_missing2"},
		DiskLinesMissing: []string{"disk_line_missing1", "disk_line_missing2"},
		ReportTimestamp:  timestamp,
	}

	// Since "updateCfgDiffs" Marshals the json, we must store the unmarshalled json here.
	//		This will need to be updated if the above text gets changed
	dbLinesMissingJson := []uint8{91, 34, 100, 98, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 49, 34, 44, 34, 100, 98, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 50, 34, 93}
	diskLinesMissingJson := []uint8{91, 34, 100, 105, 115, 107, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 49, 34, 44, 34, 100, 105, 115, 107, 95, 108, 105, 110, 101, 95, 109, 105, 115, 115, 105, 110, 103, 50, 34, 93}

	// Test Update Successful
	mock.ExpectExec("UPDATE").WithArgs(
		dbLinesMissingJson,
		diskLinesMissingJson,
		cfgFileDiffs.ReportTimestamp,
		hostName,
		cfgFileDiffs.FileName).WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := updateCfgDiffs(db, hostName, cfgFileDiffs)
	if err != nil {
		t.Errorf("updateCfgDiffs expected: nil error, actual: %v", err)
	}

	if result != true {
		t.Errorf("updateCfgDiffs expected: result == true, actual: %v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	// Test Update Unsuccessful
	mock.ExpectExec("UPDATE").WithArgs(
		dbLinesMissingJson,
		diskLinesMissingJson,
		cfgFileDiffs.ReportTimestamp,
		hostName,
		cfgFileDiffs.FileName).WillReturnResult(sqlmock.NewResult(0, 0))

	result, err = updateCfgDiffs(db, hostName, cfgFileDiffs)
	if err != nil {
		t.Errorf("updateCfgDiffs expected: nil error, actual: %v", err)
	}

	if result != false {
		t.Errorf("updateCfgDiffs expected: result == false, actual: %v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func serverExistsError(db *sqlx.DB, hostName string) (bool, error) {
	return false, fmt.Errorf("Intentional Error")
}
func serverExistsFalse(db *sqlx.DB, hostName string) (bool, error) {
	return false, nil
}
func serverExistsTrue(db *sqlx.DB, hostName string) (bool, error) {
	return true, nil
}
func updateCfgDiffsError(db *sqlx.DB, hostname string, diffs CfgFileDiffs) (bool, error) {
	return false, fmt.Errorf("Intentional Error")
}
func updateCfgDiffsTrue(db *sqlx.DB, hostname string, diffs CfgFileDiffs) (bool, error) {
	return true, nil
}
func updateCfgDiffsFalse(db *sqlx.DB, hostname string, diffs CfgFileDiffs) (bool, error) {
	return false, nil
}
func insertCfgDiffsError(db *sqlx.DB, hostname string, diffs CfgFileDiffs) error {
	return fmt.Errorf("Intentional Error")
}
func insertCfgDiffsSuccess(db *sqlx.DB, hostname string, diffs CfgFileDiffs) error {
	return nil
}

func TestPutCfgDiffs(t *testing.T) {
	var db *sqlx.DB = nil
	hostName := "myedge"
	timestamp := time.Now().UTC().String()

	cfgFileDiffs := CfgFileDiffs{
		FileName:         "TestFile.cfg",
		DBLinesMissing:   []string{"db_line_missing1", "db_line_missing2"},
		DiskLinesMissing: []string{"disk_line_missing1", "disk_line_missing2"},
		ReportTimestamp:  timestamp,
	}

	// Test when server request has error
	code, err := putCfgDiffs(db, hostName, cfgFileDiffs, serverExistsError, updateCfgDiffsError, insertCfgDiffsError)

	if code != -1 {
		t.Errorf("putCfgDiffs expected: -1 code, actual: %v", code)
	}
	if err == nil {
		t.Errorf("putCfgDiffs expected: non-nil error, actual: nil")
	}

	// Test when the server doesn't exist
	code, err = putCfgDiffs(db, hostName, cfgFileDiffs, serverExistsFalse, updateCfgDiffsError, insertCfgDiffsError)

	if code != 0 {
		t.Errorf("putCfgDiffs expected: 0 code, actual: %v", code)
	}
	if err != nil {
		t.Errorf("putCfgDiffs expected: nil error, actual: %v", err)
	}

	// Test when the server exists and the update query fails
	code, err = putCfgDiffs(db, hostName, cfgFileDiffs, serverExistsTrue, updateCfgDiffsError, insertCfgDiffsError)

	if code != -1 {
		t.Errorf("putCfgDiffs expected: -1 code, actual: %v", code)
	}
	if err == nil {
		t.Errorf("putCfgDiffs expected: non-nil error, actual: nil")
	}

	// Test when the server exists and the update is successful
	code, err = putCfgDiffs(db, hostName, cfgFileDiffs, serverExistsTrue, updateCfgDiffsTrue, insertCfgDiffsError)

	if code != 2 {
		t.Errorf("putCfgDiffs expected: 2 code, actual: %v", code)
	}
	if err != nil {
		t.Errorf("putCfgDiffs expected: non-nil error, actual: %v", err)
	}

	// Test when the server exists and the update was unsuccessful and the insert had an error
	code, err = putCfgDiffs(db, hostName, cfgFileDiffs, serverExistsTrue, updateCfgDiffsFalse, insertCfgDiffsError)

	if code != 1 {
		t.Errorf("putCfgDiffs expected: 1 code, actual: %v", code)
	}
	if err == nil {
		t.Errorf("putCfgDiffs expected: non-nil error, actual: nil")
	}

	// Test when the server exists and the update was unsuccessful and the insert was successful
	code, err = putCfgDiffs(db, hostName, cfgFileDiffs, serverExistsTrue, updateCfgDiffsFalse, insertCfgDiffsSuccess)

	if code != 1 {
		t.Errorf("putCfgDiffs expected: 1 code, actual: %v", code)
	}
	if err != nil {
		t.Errorf("putCfgDiffs expected: nil error, actual: %v", err)
	}

}

func TestGetCfgDiffsHandler(t *testing.T) {

}

func TestPutCfgDiffsHandler(t *testing.T) {

}
