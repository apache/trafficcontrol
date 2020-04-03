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

	"github.com/apache/trafficcontrol/lib/go-util"

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

func TestAddSearchableWhereClause(t *testing.T) {
	var testCases = []struct {
		description                string
		expectErrors               bool
		expectedWhereClause        string
		startingWhereClause        string
		parameters                 map[string]string
		searchQueryParamsToSQLCols map[string]WhereColumnInfo
	}{
		{
			description:         "Failure: Missing search parameter",
			expectErrors:        true,
			startingWhereClause: "",
			expectedWhereClause: "",
			parameters: map[string]string{
				SearchValueParam: "test",
			},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{},
		},
		{
			description:         "Failure: Missing search parameter value",
			expectErrors:        true,
			startingWhereClause: "",
			expectedWhereClause: "",
			parameters: map[string]string{
				SearchColumnsParam: "key1",
			},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{},
		},
		{
			description:         "Success: Existing Where Clause",
			expectErrors:        false,
			startingWhereClause: "\nWHERE foo=:bar",
			expectedWhereClause: "\nWHERE foo=:bar AND (UPPER(t.col) LIKE UPPER(:key1_search) OR UPPER(t1.col) LIKE UPPER(:key2_search)) ",
			parameters: map[string]string{
				SearchValueParam:   "test",
				SearchColumnsParam: "key1,key2",
			},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{
				"key1": WhereColumnInfo{"t.col", nil},
				"key2": WhereColumnInfo{"t1.col", nil},
			},
		},
		{
			description:         "Success: New Where Clause",
			expectErrors:        false,
			startingWhereClause: "",
			expectedWhereClause: "\nWHERE UPPER(t.col) LIKE UPPER(:key1_search) OR UPPER(t1.col) LIKE UPPER(:key2_search)",
			parameters: map[string]string{
				SearchValueParam:   "test",
				SearchColumnsParam: "key1,key2",
			},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{
				"key1": WhereColumnInfo{"t.col", nil},
				"key2": WhereColumnInfo{"t1.col", nil},
			},
		},
		{
			description:         "Success: No search parameters",
			expectErrors:        false,
			startingWhereClause: "",
			expectedWhereClause: "",
			parameters:          map[string]string{},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{
				"key1": WhereColumnInfo{"t.col", nil},
				"key2": WhereColumnInfo{"t1.col", nil},
			},
		},
		{
			description:         "Success: Defined join operator - AND",
			expectErrors:        false,
			startingWhereClause: "",
			expectedWhereClause: "\nWHERE UPPER(t.col) LIKE UPPER(:key1_search) AND UPPER(t1.col) LIKE UPPER(:key2_search)",
			parameters: map[string]string{
				SearchValueParam:   "test",
				SearchColumnsParam: "key1,key2",
				SearchFilterParam:  "AND",
			},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{
				"key1": WhereColumnInfo{"t.col", nil},
				"key2": WhereColumnInfo{"t1.col", nil},
			},
		},
		{
			description:         "Success: Defined join operator - OR",
			expectErrors:        false,
			startingWhereClause: "",
			expectedWhereClause: "\nWHERE UPPER(t.col) LIKE UPPER(:key1_search) OR UPPER(t1.col) LIKE UPPER(:key2_search)",
			parameters: map[string]string{
				SearchValueParam:   "test",
				SearchColumnsParam: "key1,key2",
				SearchFilterParam:  "OR",
			},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{
				"key1": WhereColumnInfo{"t.col", nil},
				"key2": WhereColumnInfo{"t1.col", nil},
			},
		},
		{
			description:         "Failure: Invalid join operator",
			expectErrors:        true,
			startingWhereClause: "",
			expectedWhereClause: "",
			parameters: map[string]string{
				SearchValueParam:   "test",
				SearchColumnsParam: "key1,key2",
				SearchFilterParam:  "Bogus",
			},
			searchQueryParamsToSQLCols: map[string]WhereColumnInfo{
				"key1": WhereColumnInfo{"t.col", nil},
				"key2": WhereColumnInfo{"t1.col", nil},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			t.Log("Starting test scenario: ", tc.description)
			where, _, errors := AddSearchableWhereClause(tc.startingWhereClause, nil, tc.parameters, tc.searchQueryParamsToSQLCols)
			if len(errors) == 0 && tc.expectErrors {
				t.Error("expected errors on building search where clause and received non")
				return
			} else if len(errors) != 0 && !tc.expectErrors {
				t.Errorf("expected no errors on building search where clause and received: %v", util.JoinErrs(errors))
				return
			}
			if where != tc.expectedWhereClause {
				t.Errorf("expected resulting WHERE Clause %v received %v", tc.expectedWhereClause, where)
			}
		})
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

func TestGetDSIDAndCDNFromName(t *testing.T) {
	var testCases = []struct {
		description  string
		storageError error
		found        bool
	}{
		{
			description:  "Success: DS ID and CDN Name found",
			storageError: nil,
			found:        true,
		},
		{
			description:  "Failure: DS ID or CDN Name not found",
			storageError: nil,
			found:        false,
		},
		{
			description:  "Failure: Storage error getting DS ID or CDN Name",
			storageError: errors.New("error getting the delivery service ID or the CDN name"),
			found:        false,
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
				"id",
				"name",
			})
			mock.ExpectBegin()
			if testCase.storageError != nil {
				mock.ExpectQuery("SELECT").WillReturnError(testCase.storageError)
			} else {
				if testCase.found {
					rows = rows.AddRow(1, "testCdn")
				}
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			}
			mock.ExpectCommit()
			_, _, exists, err := GetDSIDAndCDNFromName(db.MustBegin().Tx, "testDs")
			if testCase.storageError != nil && err == nil {
				t.Errorf("Storage error expected: received no storage error")
			}
			if testCase.storageError == nil && err != nil {
				t.Errorf("Storage error not expected: received storage error")
			}
			if testCase.found != exists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.found, exists)
			}
		})
	}

}

func TestGetCDNIDFromName(t *testing.T) {
	var testCases = []struct {
		description  string
		storageError error
		found        bool
	}{
		{
			description:  "Success: CDN ID found",
			storageError: nil,
			found:        true,
		},
		{
			description:  "Failure: CDN ID not found",
			storageError: nil,
			found:        false,
		},
		{
			description:  "Failure: Storage error getting CDN ID",
			storageError: errors.New("error getting the CDN ID"),
			found:        false,
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
				"id",
			})
			mock.ExpectBegin()
			if testCase.storageError != nil {
				mock.ExpectQuery("SELECT").WillReturnError(testCase.storageError)
			} else {
				if testCase.found {
					rows = rows.AddRow(1)
				}
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			}
			mock.ExpectCommit()
			_, exists, err := GetCDNIDFromName(db.MustBegin().Tx, "testCdn")
			if testCase.storageError != nil && err == nil {
				t.Errorf("Storage error expected: received no storage error")
			}
			if testCase.storageError == nil && err != nil {
				t.Errorf("Storage error not expected: received storage error")
			}
			if testCase.found != exists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.found, exists)
			}
		})
	}

}
