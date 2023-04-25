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

func TestDeleteCDNByName(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM cdn").WithArgs("cdn1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = deleteCDNByName(db.MustBegin().Tx, "cdn1")
	if err != nil {
		t.Fatalf("no error expected while deleting CDN by name, but got: %v", err)
	}
}

func TestCDNUsed(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{"?column?"}
	rows := sqlmock.NewRows(cols)
	rows.AddRow(5)

	mock.ExpectBegin()
	mock.ExpectQuery("WITH cdn_id as").WithArgs("cdn1").WillReturnRows(rows)
	mock.ExpectCommit()

	unused, err := cdnUnused(db.MustBegin().Tx, "cdn1")
	if err != nil {
		t.Fatalf("no error expected in call to cdnUnused, but got: %v", err)
	}
	if unused {
		t.Errorf("expected CDN to be used, but is unused")
	}
}
