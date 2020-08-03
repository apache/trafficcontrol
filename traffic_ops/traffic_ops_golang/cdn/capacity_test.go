package cdn

import (
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

func TestGetServiceInterfaces(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	cols := []string{"host_name", "interface"}
	rows := sqlmock.NewRows(cols)
	rows = rows.AddRow(
		"host1",
		"eth1",
	)
	rows = rows.AddRow(
		"host2",
		"eth2",
	)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	m, err := getServiceInterfaces(db.MustBegin().Tx)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err.Error())
	}
	if len(m) != 2 {
		t.Errorf("Expected a result of length %v, got %v instead", 2, len(m))
	}
	if m["host1"] != "eth1" {
		t.Errorf("Expected host1 to have service interface eth1, got %v instead", m["host1"])
	}
	if m["host2"] != "eth2" {
		t.Errorf("Expected host2 to have service interface eth2, got %v instead", m["host2"])
	}
}

func TestGetStatsFromServiceInterface(t *testing.T) {
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
	kbps, maxKbps, err := getStatsFromServiceInterface(c)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err.Error())
	}
	if kbps != 24.5 || maxKbps != 66.8 {
		t.Errorf("Expected kbps to be 24.5, got %v; Expected maxKbps to be 66.8, got %v", kbps, maxKbps)
	}
}
