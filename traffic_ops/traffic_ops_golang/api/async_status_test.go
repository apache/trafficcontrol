package api

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

func TestInsertAsyncStatus(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	expectedMessage := "test async message"
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)
	mock.ExpectQuery("INSERT").WithArgs(AsyncPending, expectedMessage).WillReturnRows(rows)

	asyncId, errCode, userErr, sysErr := InsertAsyncStatus(db.MustBegin().Tx, expectedMessage)

	if userErr != nil {
		t.Fatalf("userError was expected to be nil but got %v", userErr)
		return
	}
	if sysErr != nil {
		t.Fatalf("sysErr was expected to be nil but got %v", sysErr)
		return
	}
	if errCode != http.StatusOK {
		t.Fatalf("errCode was expected to be %v but got %v", http.StatusOK, errCode)
		return
	}
	if asyncId != 1 {
		t.Fatalf("asyncId was expected to be 1 but got %v", asyncId)
		return
	}
}

func TestUpdateAsyncStatus(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	expectedMessage := "test updated async message"
	expectedStatus := AsyncPending
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(expectedStatus, expectedMessage, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	updateErr := UpdateAsyncStatus(db, expectedStatus, expectedMessage, 1, false)

	if updateErr != nil {
		t.Fatalf("updateErr was expected to be nil but got %v", updateErr)
		return
	}
}

func TestUpdateAsyncStatusFinished(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	expectedMessage := "test job complete"
	expectedStatus := AsyncSucceeded
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(expectedStatus, expectedMessage, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	updateErr := UpdateAsyncStatus(db, expectedStatus, expectedMessage, 1, true)

	if updateErr != nil {
		t.Fatalf("updateErr was expected to be nil but got %v", updateErr)
		return
	}
}
