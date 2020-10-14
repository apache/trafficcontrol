package deliveryservicesregexes

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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestValidateDSRegexOrderExisting(t *testing.T) {
	expected := `'setNumber' cannot add regex, another regex with the same order exists`
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	cols := []string{"name", "use_in_table"}
	rows := sqlmock.NewRows(cols)
	rows = rows.AddRow(
		"HTTP",
		"regex",
	)
	cols2 := []string{"deliveryservice"}
	rows2 := sqlmock.NewRows(cols2)
	rows2 = rows2.AddRow(
		1,
	)

	regex := tc.DeliveryServiceRegexPost{Type: 33, SetNumber: 3, Pattern: ".*"}
	mock.ExpectBegin()
	mock.ExpectQuery("select").WithArgs(1, regex.SetNumber).WillReturnRows(rows2)
	mock.ExpectQuery("SELECT").WithArgs(regex.Type).WillReturnRows(rows)
	mock.ExpectCommit()
	tx := db.MustBegin().Tx
	err = validateDSRegex(tx, regex, 1)
	if err == nil {
		t.Fatalf("Expected error '%v' but got none", expected)
	}
}

func TestValidateDSRegex(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	cols := []string{"name", "use_in_table"}
	rows := sqlmock.NewRows(cols)
	rows = rows.AddRow(
		"HTTP",
		"regex",
	)
	cols2 := []string{"deliveryservice"}
	rows2 := sqlmock.NewRows(cols2)

	regex := tc.DeliveryServiceRegexPost{Type: 33, SetNumber: 3, Pattern: ".*"}
	mock.ExpectBegin()
	mock.ExpectQuery("select").WithArgs(1, regex.SetNumber).WillReturnRows(rows2)
	mock.ExpectQuery("SELECT").WithArgs(regex.Type).WillReturnRows(rows)
	mock.ExpectCommit()
	tx := db.MustBegin().Tx
	err = validateDSRegex(tx, regex, 1)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err.Error())
	}
}
