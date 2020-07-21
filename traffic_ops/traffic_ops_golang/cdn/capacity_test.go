package cdn

import (
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

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

func TestCheckServiceInterface(t *testing.T) {
	m := make(map[tc.InterfaceName]CacheStat)
	var data1 []CacheStatData
	var data2 []CacheStatData
	kbpsData := CacheStatData{Value: 24.5}
	maxKbpsData := CacheStatData{Value: 66.8}
	data1 = append(data1, kbpsData)
	data2 = append(data2, maxKbpsData)

	c := CacheStat{
		KBPS:    data1,
		MaxKBPS: data2,
	}
	m["notservice"] = c

	kbpsData = CacheStatData{Value: 50.0}
	maxKbpsData = CacheStatData{Value: 100.9}
	data1 = nil
	data2 = nil
	data1 = append(data1, kbpsData)
	data2 = append(data2, maxKbpsData)

	c = CacheStat{
		KBPS:    data1,
		MaxKBPS: data2,
	}

	m["service"] = c

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{"service_address"}
	rows := sqlmock.NewRows(cols)
	rows = rows.AddRow(
		false,
	)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs("notservice", "host").WillReturnRows(rows)
	rows = sqlmock.NewRows(cols)
	rows = rows.AddRow(
		true,
	)

	mock.ExpectQuery("SELECT").WithArgs("service", "host").WillReturnRows(rows)
	mock.ExpectCommit()

	kbps, maxKbps := checkServiceInterface(db.MustBegin().Tx, "host", m)

	if kbps != 50.0 || maxKbps != 100.9 {
		t.Fatalf("Expected kbps = 50.0, got %v, expected maxKbps = 100.9, got %v", kbps, maxKbps)
	}
}