package dbhelpers

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
	"errors"
	"strings"
	"testing"
	"unicode"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func stripAllWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func TestBuildQuery(t *testing.T) {
	v := map[string]string{"param1": "queryParamv1", "param2": "queryParamv2", "limit": "20", "offset": "10"}

	selectStmt := `SELECT
	t.col1,
	t.col2
FROM table t
`
	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]WhereColumnInfo{
		"param1": WhereColumnInfo{"t.col1", nil},
		"param2": WhereColumnInfo{"t.col2", nil},
	}
	where, orderBy, pagination, queryValues, _ := BuildWhereAndOrderByAndPagination(v, queryParamsToSQLCols)
	query := selectStmt + where + orderBy + pagination
	actualQuery := stripAllWhitespace(query)

	expectedPagination := "\nLIMIT " + v["limit"] + "\nOFFSET " + v["offset"]
	if pagination != expectedPagination {
		t.Errorf("expected: %s for pagination, actual: %s", expectedPagination, pagination)
	}

	if queryValues == nil {
		t.Errorf("expected: nil error, actual: %v", queryValues)
	}
	expectedV1 := v["param1"]
	actualV1 := queryValues["param1"]
	if expectedV1 != actualV1 {
		t.Errorf("expected: %v error, actual: %v", expectedV1, actualV1)
	}

	if strings.Contains(actualQuery, expectedV1) {
		t.Errorf("expected: %v error, actual: %v", actualQuery, expectedV1)
	}

	expectedV2 := v["param2"]
	if strings.Contains(actualQuery, expectedV2) {
		t.Errorf("expected: %v error, actual: %v", actualQuery, expectedV2)
	}

}

func TestGetCacheGroupByName(t *testing.T) {
	var testCases = []struct {
		description  string
		storageError error
		cgExists     bool
	}{
		{
			description:  "Success: Cache Group exists",
			storageError: nil,
			cgExists:     true,
		},
		{
			description:  "Failure: Cache Group does not exist",
			storageError: nil,
			cgExists:     false,
		},
		{
			description:  "Failure: Storage error getting Cache Group",
			storageError: errors.New("error getting the group name"),
			cgExists:     false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			t.Log("Starting test scenario: ", testCase.description)
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()
			rows := sqlmock.NewRows([]string{
				"name",
			})
			mock.ExpectBegin()
			if testCase.storageError != nil {
				mock.ExpectQuery("cachegroup").WillReturnError(testCase.storageError)
			} else {
				if testCase.cgExists {
					rows = rows.AddRow("cachegroup_name")
				}
				mock.ExpectQuery("cachegroup").WillReturnRows(rows)
			}
			mock.ExpectCommit()
			_, exists, err := GetCacheGroupNameFromID(db.MustBegin().Tx, int64(1))
			if testCase.storageError != nil && err == nil {
				t.Errorf("Storage error expected: received no storage error")
			}
			if testCase.storageError == nil && err != nil {
				t.Errorf("Storage error not expected: received storage error")
			}
			if testCase.cgExists != exists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.cgExists, exists)
			}
		})
	}

}
