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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
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
	err = validateDSRegex(tx, regex, 1, false)
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
	err = validateDSRegex(tx, regex, 1, true)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err.Error())
	}
}

func TestUpdateImmutableRegex(t *testing.T) {
	expected := `'setNumber' cannot update regex with set number 0 and type HOST_REGEXP`
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
		"HOST_REGEXP",
		"regex",
	)
	cols2 := []string{"deliveryservice"}
	rows2 := sqlmock.NewRows(cols2)

	cols3 := []string{"name"}
	rows3 := sqlmock.NewRows(cols3)
	rows3 = rows3.AddRow("HOST_REGEXP")
	regex := tc.DeliveryServiceRegexPost{Type: 33, SetNumber: 0, Pattern: ".*"}
	mock.ExpectBegin()
	mock.ExpectQuery("select").WithArgs(1, regex.SetNumber).WillReturnRows(rows2)
	mock.ExpectQuery("select").WithArgs(regex.Type).WillReturnRows(rows3)
	mock.ExpectQuery("SELECT").WithArgs(regex.Type).WillReturnRows(rows)
	mock.ExpectCommit()
	tx := db.MustBegin().Tx
	err = validateDSRegex(tx, regex, 1, false)
	if err == nil {
		t.Fatalf("Expected error forbidding updates to regex with set number 0 and type HOST_REGEXP, but got none")
	}
	if err.Error() != expected {
		t.Fatalf("expected error detail to be %v, but got %v instead", expected, err.Error())
	}
}

func TestGetCurrentDetails(t *testing.T) {
	expected := `cannot change/ delete a regex with an order of 0 and type name of HOST_REGEXP`
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	cols := []string{"set_number", "name"}
	rows := sqlmock.NewRows(cols)
	rows = rows.AddRow(
		0,
		"HOST_REGEXP",
	)
	dsID := 1
	regexID := 3
	mock.ExpectBegin()
	mock.ExpectQuery("select").WithArgs(dsID, regexID).WillReturnRows(rows)
	mock.ExpectCommit()
	tx := db.MustBegin().Tx
	err = getCurrentDetails(tx, dsID, regexID)
	if err == nil {
		t.Fatalf("Expected error forbidding updates to regex with set number 0 and type HOST_REGEXP, but got none")
	}
	if err.Error() != expected {
		t.Fatalf("expected error detail to be %v, but got %v instead", expected, err.Error())
	}
}
