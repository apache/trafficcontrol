package deliveryservicesregexes

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

func TestValidateDSRegexOrder(t *testing.T) {
	expected := `cannot add regex, another regex with the same order exists`
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	cols := []string{"deliveryservice"}
	rows := sqlmock.NewRows(cols)
	rows = rows.AddRow(
		1,
	)
	mock.ExpectBegin()
	mock.ExpectQuery("select").WithArgs(1, 3).WillReturnRows(rows)
	tx := db.MustBegin().Tx
	err = validateDSRegexOrder(tx, 1, 3)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err.Error() != expected {
		t.Fatalf("Expected error was %v, got %v", expected, err.Error())
	}
	mock.ExpectQuery("select").WithArgs(1, 4).WillReturnRows(nil)
	mock.ExpectCommit()
	err = validateDSRegexOrder(tx, 1, 3)
	if err != nil {
		t.Fatalf("Expect no error, got %v", err.Error())
	}
	err = validateDSRegexOrder(tx, 1, -1)
	if err == nil {
		t.Fatal("Expect error saying cannot add regex with order < 0, got nothing")
	}
	if err.Error() != "cannot add regex with order < 0" {
		t.Errorf("Expected error detail to be 'cannot add regex with order <0', got %v", err.Error())
	}
}
