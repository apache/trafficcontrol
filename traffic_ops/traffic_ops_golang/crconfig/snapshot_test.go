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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/monitoring"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ExpectedGetSnapshot(crc *tc.CRConfig) ([]byte, error) {
	return json.Marshal(crc)
}

func ExpectedGetMontioringSnapshot(crc *tc.CRConfig, tx *sql.Tx) ([]byte, error) {
	tm, _ := monitoring.GetMonitoringJSON(tx, *crc.Stats.CDNName, true)
	return json.Marshal(tm)
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

	for i := 0; i < len(SnapshotTables)*2; i++ {
		mock.ExpectExec("INSERT").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// actual snapshot insert, after the x_snapshot tables
	mock.ExpectExec("INSERT").WithArgs().WillReturnResult(sqlmock.NewResult(1, 1))

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	defer tx.Commit()

	if err := Snapshot(tx, tc.CDNName(cdn)); err != nil {
		t.Fatalf("GetSnapshot err expected: nil, actual: %v", err)
	}
}
