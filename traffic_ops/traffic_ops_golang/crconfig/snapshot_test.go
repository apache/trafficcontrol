package crconfig

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
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/traffic_ops/traffic_ops_golang/monitoring"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ExpectedGetSnapshot(crc *tc.CRConfig) ([]byte, error) {
	return json.Marshal(crc)
}

func ExpectedGetMontioringSnapshot(crc *tc.CRConfig, tx *sql.Tx) ([]byte, error) {
	tm, _ := monitoring.GetMonitoringJSON(tx, *crc.Stats.CDNName)
	return json.Marshal(tm)
}

func MockGetSnapshot(mock sqlmock.Sqlmock, expected []byte, cdn string) {
	rows := sqlmock.NewRows([]string{"snapshot"})
	rows = rows.AddRow(expected)
	rows = rows.AddRow(expected)
	mock.ExpectQuery("SELECT").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetSnapshot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	crc := &tc.CRConfig{}
	crc.Stats.CDNName = &cdn

	mock.ExpectBegin()
	expected, err := ExpectedGetSnapshot(crc)
	if err != nil {
		t.Fatalf("GetSnapshot creating expected err expected: nil, actual: %v", err)
	}
	MockGetSnapshot(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, exists, err := GetSnapshot(tx, cdn)
	if err != nil {
		t.Fatalf("GetSnapshot err expected: nil, actual: %v", err)
	}
	if !exists {
		t.Fatalf("GetSnapshot exists expected: true, actual: false")
	}

	if !reflect.DeepEqual(string(expected), actual) {
		t.Errorf("GetSnapshot expected: %+v, actual: %+v", string(expected), actual)
	}
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Any struct{}

// Match satisfies sqlmock.Argument interface
func (a Any) Match(v driver.Value) bool {
	return true
}

func MockSnapshot(mock sqlmock.Sqlmock, expected []byte, expectedtm []byte, cdn string) {
	mock.ExpectExec("insert").WithArgs(cdn, expected, AnyTime{}, expectedtm).WillReturnResult(sqlmock.NewResult(1, 1))
}

func TestSnapshot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	crc := &tc.CRConfig{}
	crc.Stats.CDNName = &cdn
	mock.ExpectBegin()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	expected, err := ExpectedGetSnapshot(crc)
	if err != nil {
		t.Fatalf("GetSnapshot creating expected err expected: nil, actual: %v", err)
	}

	expectedtm, err := ExpectedGetMontioringSnapshot(crc, tx)
	if err != nil {
		t.Fatalf("GetSnapshotMonitor creating expected err expected: nil, actual: %v", err)
	}

	tm, _ := monitoring.GetMonitoringJSON(tx, *crc.Stats.CDNName)
	MockSnapshot(mock, expected, expectedtm, cdn)
	mock.ExpectCommit()

	defer tx.Commit()

	if err := Snapshot(tx, crc, tm); err != nil {
		t.Fatalf("GetSnapshot err expected: nil, actual: %v", err)
	}
}
