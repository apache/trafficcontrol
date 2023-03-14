package servers

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
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCheckLastAvailableEdgeOrOrigin(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	// check DS with no topology or mso
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"available", "available"})
	rows.AddRow(true, true)
	mock.ExpectQuery("SELECT").WithArgs(1, 2).WillReturnRows(rows)
	sc, userErr, sysErr := checkLastAvailableEdgeOrOrigin(1, 2, false, false, db.MustBegin().Tx)
	if sysErr != nil {
		t.Errorf("expected no system error, but got %v", sysErr)
	}
	if userErr == nil {
		t.Errorf("expected error because removing the given server would result in active DS with no REPORTED/ ONLINE EDGE servers, but got nothing")
	}
	if sc != http.StatusConflict {
		t.Errorf("expected 409 status code, but got %d", sc)
	}

	// check DS with topology, but no MSO
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"available", "available"})
	rows.AddRow(true, true)
	mock.ExpectQuery("SELECT").WithArgs(1, 2).WillReturnRows(rows)
	sc, userErr, sysErr = checkLastAvailableEdgeOrOrigin(1, 2, false, true, db.MustBegin().Tx)
	if userErr != nil || sysErr != nil {
		t.Errorf("expected no error, but got userErr: %v, sysErr: %v", userErr, sysErr)
	}
	if sc != http.StatusOK {
		t.Errorf("ecpected status code 200, but got %d", sc)
	}

	// check DS with MSO, but no topology
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"available", "available"})
	rows.AddRow(false, true)
	mock.ExpectQuery("SELECT").WithArgs(1, 2).WillReturnRows(rows)
	sc, userErr, sysErr = checkLastAvailableEdgeOrOrigin(1, 2, true, false, db.MustBegin().Tx)
	if sysErr != nil {
		t.Errorf("expected no system error, but got %v", sysErr)
	}
	if userErr == nil {
		t.Errorf("expected error because removing the given server would result in active DS with no REPORTED/ ONLINE EDGE servers, but got nothing")
	}
	if sc != http.StatusConflict {
		t.Errorf("expected 409 status code, but got %d", sc)
	}
}

func TestDeleteDSServer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"server"})
	mock.ExpectQuery("DELETE").WithArgs(1, 2).WillReturnRows(rows)
	exists, err := deleteDSServer(db.MustBegin().Tx, 1, 2)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if exists {
		t.Errorf("expected exists to be false, but got true")
	}

	rows = sqlmock.NewRows([]string{"server"})
	rows.AddRow(2)
	mock.ExpectBegin()
	mock.ExpectQuery("DELETE").WithArgs(1, 2).WillReturnRows(rows)
	exists, err = deleteDSServer(db.MustBegin().Tx, 1, 2)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if !exists {
		t.Errorf("expected exists to be true, but got false")
	}
}
