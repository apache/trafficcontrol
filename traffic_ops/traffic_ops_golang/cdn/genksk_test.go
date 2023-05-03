package cdn

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

	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetKSKParams(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name", "value"})
	rows.AddRow("test", "2")
	mock.ExpectBegin()
	mock.ExpectQuery("WITH cdn_profile_id").WithArgs("test").WillReturnRows(rows)
	mock.ExpectCommit()

	ttl, mult, err := getKSKParams(db.MustBegin().Tx, "test")
	if ttl != nil {
		t.Errorf("expected: nil, got: %v", ttl)
	}
	if *mult != 2 {
		t.Errorf("expected: 2, got: %v", *mult)
	}
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestGetKSKParamsDNSKey(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	rows1 := sqlmock.NewRows([]string{"name", "value"})
	rows1.AddRow("tld.ttls.DNSKEY", "5")
	mock.ExpectBegin()
	mock.ExpectQuery("WITH cdn_profile_id").WithArgs("test").WillReturnRows(rows1)
	mock.ExpectCommit()

	ttl1, mult1, err1 := getKSKParams(db.MustBegin().Tx, "test")
	if *ttl1 != 5 {
		t.Errorf("expected: 5, got: %v", *ttl1)
	}
	if mult1 != nil {
		t.Errorf("expected: nil, got: %v", mult1)
	}
	if err1 != nil {
		t.Errorf("%s", err1)
	}
}
