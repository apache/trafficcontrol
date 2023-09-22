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
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"

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
			_, exists, err := GetCacheGroupNameFromID(db.MustBegin().Tx, 1)
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

// createServerInterfaces takes in a cache id and creates the interfaces/ipaddresses for it
func createServerIntefaces(cacheID int) []tc.ServerInterfaceInfoV40 {
	return []tc.ServerInterfaceInfoV40{
		{
			ServerInterfaceInfo: tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "5.6.7.8",
						Gateway:        util.StrPtr("5.6.7.0/24"),
						ServiceAddress: true,
					},
					{
						Address:        "2020::4",
						Gateway:        util.StrPtr("fd53::9"),
						ServiceAddress: true,
					},
					{
						Address:        "5.6.7.9",
						Gateway:        util.StrPtr("5.6.7.0/24"),
						ServiceAddress: false,
					},
					{
						Address:        "2021::4",
						Gateway:        util.StrPtr("fd53::9"),
						ServiceAddress: false,
					},
				},
				MaxBandwidth: util.Uint64Ptr(2500),
				Monitor:      true,
				MTU:          util.Uint64Ptr(1500),
				Name:         "interfaceName" + strconv.Itoa(cacheID),
			},
			RouterHostName: "",
			RouterPortName: "",
		},
		{
			ServerInterfaceInfo: tc.ServerInterfaceInfo{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "6.7.8.9",
						Gateway:        util.StrPtr("6.7.8.0/24"),
						ServiceAddress: true,
					},
					{
						Address:        "2021::4",
						Gateway:        util.StrPtr("fd54::9"),
						ServiceAddress: true,
					},
					{
						Address:        "6.6.7.9",
						Gateway:        util.StrPtr("6.6.7.0/24"),
						ServiceAddress: false,
					},
					{
						Address:        "2022::4",
						Gateway:        util.StrPtr("fd53::9"),
						ServiceAddress: false,
					},
				},
				MaxBandwidth: util.Uint64Ptr(1500),
				Monitor:      false,
				MTU:          util.Uint64Ptr(1500),
				Name:         "interfaceName2" + strconv.Itoa(cacheID),
			},
			RouterHostName: "",
			RouterPortName: "",
		},
	}
}

func mockServerInterfaces(mock sqlmock.Sqlmock, cacheID int, serverInterfaces []tc.ServerInterfaceInfoV40) {
	interfaceRows := sqlmock.NewRows([]string{"max_bandwidth", "monitor", "mtu", "name", "server", "router_host_name", "router_port_name"})
	ipAddressRows := sqlmock.NewRows([]string{"address", "gateway", "service_address", "interface", "server"})
	for _, interf := range serverInterfaces {
		interfaceRows = interfaceRows.AddRow(*interf.MaxBandwidth, interf.Monitor, *interf.MTU, interf.Name, cacheID, interf.RouterHostName, interf.RouterPortName)
		for _, ip := range interf.IPAddresses {
			ipAddressRows = ipAddressRows.AddRow(ip.Address, *ip.Gateway, ip.ServiceAddress, interf.Name, cacheID)
		}
	}

	mock.ExpectQuery("SELECT (.+) FROM interface").WillReturnRows(interfaceRows)
	mock.ExpectQuery("SELECT (.+) FROM ip_address").WillReturnRows(ipAddressRows)
}

func TestGetServerInterfaces(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unable to open mock db: %v", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cacheID := 1
	serverInterfaces := createServerIntefaces(cacheID)
	mock.ExpectBegin()
	mockServerInterfaces(mock, cacheID, serverInterfaces)

	dbCtx, cancelTx := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	serversMap, err := GetServersInterfaces([]int{cacheID}, tx)
	if err != nil {
		t.Fatal(err)
	}
	if len(serversMap) != 1 {
		t.Fatalf("expected to get a single server, got %v", len(serversMap))
	}
	if interfacesMap, ok := serversMap[cacheID]; ok {
		if len(interfacesMap) != len(serverInterfaces) {
			t.Fatalf("expected cache %v to have %v interfaces, got %v", cacheID, len(serverInterfaces), len(interfacesMap))
		}

		for _, interf := range serverInterfaces {
			if calculatedInterface, ok := interfacesMap[interf.Name]; ok {
				if !reflect.DeepEqual(calculatedInterface, interf) {
					t.Fatalf("expected %v to match %v", calculatedInterface, interf)
				}

			} else {
				t.Fatalf("expected map to contain interface %v, but did not", interf.Name)
			}
		}
	} else {
		t.Fatalf("Cache %v not found in servers map", cacheID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
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
			expectedIDValue := 3
			if testCase.storageError != nil {
				mock.ExpectQuery("SELECT").WillReturnError(testCase.storageError)
			} else {
				if testCase.found {
					rows = rows.AddRow(expectedIDValue)
				}
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			}
			mock.ExpectCommit()
			id, exists, err := GetCDNIDFromName(db.MustBegin().Tx, "testCdn")
			if testCase.storageError != nil && err == nil && id == expectedIDValue {
				t.Errorf("Storage error expected: received no storage error")
			}
			if testCase.storageError == nil && err != nil && id == expectedIDValue {
				t.Errorf("Storage error not expected: received storage error")
			}
			if exists && testCase.storageError == nil && err == nil && id != expectedIDValue {
				t.Errorf("Expected ID %d, but got %d", expectedIDValue, id)
			}
			if testCase.found != exists && id == 0 {
				t.Errorf("Expected return exists: %t, actual %t", testCase.found, exists)
			}
		})
	}

}

func TestGetSCInfo(t *testing.T) {
	var testCases = []struct {
		description   string
		name          string
		expectedError error
		exists        bool
	}{
		{
			description:   "Success: Get valid SC",
			name:          "hdd",
			expectedError: nil,
			exists:        true,
		},
		{
			description:   "Failure: SC not in DB",
			name:          "disk",
			expectedError: sql.ErrNoRows,
			exists:        false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"count"})
			if testCase.exists {
				rows = rows.AddRow(1)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			scExists, err := GetSCInfo(db.MustBegin().Tx, testCase.name)
			if testCase.exists != scExists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.exists, scExists)
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("getSCInfo expected: %s, actual: %s", testCase.expectedError, err)
			}
		})
	}
}

func TestServiceCategoryExists(t *testing.T) {
	var testCases = []struct {
		description   string
		name          string
		expectedError error
		exists        bool
	}{
		{
			description:   "Success: Get valid Service Category",
			name:          "testServiceCategory1",
			expectedError: nil,
			exists:        true,
		},
		{
			description:   "Failure: Service Category not in DB",
			name:          "testServiceCategory2",
			expectedError: sql.ErrNoRows,
			exists:        false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"count"})
			if testCase.exists {
				rows = rows.AddRow(1)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			scExists, err := ServiceCategoryExists(db.MustBegin().Tx, testCase.name)
			if testCase.exists != scExists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.exists, scExists)
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("ServiceCategoryExists expected: %s, actual: %s", testCase.expectedError, err)
			}
		})
	}
}

func TestASNExists(t *testing.T) {
	var testCases = []struct {
		description   string
		id            string
		expectedError error
		exists        bool
	}{
		{
			description:   "Success: Get valid ASN",
			id:            "1",
			expectedError: nil,
			exists:        true,
		},
		{
			description:   "Failure: ASN not in DB",
			id:            "10",
			expectedError: sql.ErrNoRows,
			exists:        false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"count"})
			if testCase.exists {
				rows = rows.AddRow(1)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			asnExists, err := ASNExists(db.MustBegin().Tx, testCase.id)
			if testCase.exists != asnExists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.exists, asnExists)
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("getSCInfo expected: %s, actual: %s", testCase.expectedError, err)
			}
		})
	}
}

func TestCacheGroupExistsExists(t *testing.T) {
	var testCases = []struct {
		description   string
		name          string
		expectedError error
		exists        bool
	}{
		{
			description:   "Success: Get valid Cache Group",
			name:          "testCacheGroup1",
			expectedError: nil,
			exists:        true,
		},
		{
			description:   "Failure: Cache Group not in DB",
			name:          "testCacheGroup2",
			expectedError: sql.ErrNoRows,
			exists:        false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"count"})
			if testCase.exists {
				rows = rows.AddRow(1)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			cgExists, err := CacheGroupExists(db.MustBegin().Tx, testCase.name)
			if testCase.exists != cgExists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.exists, cgExists)
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("CacheGroupExists expected: %s, actual: %s", testCase.expectedError, err)
			}
		})
	}
}

func TestDivisionExists(t *testing.T) {
	var testCases = []struct {
		description   string
		id            string
		expectedError error
		exists        bool
	}{
		{
			description:   "Success: Get valid Division",
			id:            "1",
			expectedError: nil,
			exists:        true,
		},
		{
			description:   "Failure: Division not in DB",
			id:            "10",
			expectedError: sql.ErrNoRows,
			exists:        false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"count"})
			if testCase.exists {
				rows = rows.AddRow(1)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			divisionExists, err := DivisionExists(db.MustBegin().Tx, testCase.id)
			if testCase.exists != divisionExists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.exists, divisionExists)
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("DivisionExists Error. expected: %s, actual: %s", testCase.expectedError, err)
			}
		})
	}
}

func TestProfileExists(t *testing.T) {
	var testCases = []struct {
		description   string
		id            string
		expectedError error
		exists        bool
	}{
		{
			description:   "Success: Get valid Profile",
			id:            "1",
			expectedError: nil,
			exists:        true,
		},
		{
			description:   "Failure: Profile not in DB",
			id:            "5",
			expectedError: sql.ErrNoRows,
			exists:        false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"SELECT count(name)"})
			if testCase.exists {
				rows = rows.AddRow(1)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			profileExists, err := ProfileExists(db.MustBegin().Tx, testCase.id)
			if testCase.exists != profileExists {
				t.Errorf("Expected return exists: %t, actual %t", testCase.exists, profileExists)
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("getSCInfo expected: %s, actual: %s", testCase.expectedError, err)
			}
		})
	}
}
