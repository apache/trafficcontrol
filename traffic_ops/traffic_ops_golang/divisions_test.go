package main

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
	"net/url"
	"testing"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestDivisions() []tc.Division {
	regions := []tc.Division{}
	testCase := tc.Division{
		ID:          1,
		Name:        "division1",
		LastUpdated: "lastUpdated",
	}
	regions = append(regions, testCase)

	testCase2 := testCase
	testCase2.Name = "region2"
	regions = append(regions, testCase2)

	return regions
}

func TestGetDivisions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testCase := getTestDivisions()
	cols := test.ColsFromStructByTag("db", tc.Division{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testCase {
		rows = rows.AddRow(
			ts.ID,
			ts.LastUpdated,
			ts.Name,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := url.Values{}
	v.Set("dsId", "1")

	servers, err := getDivisions(v, db, PrivLevelAdmin)
	if err != nil {
		t.Errorf("getDivisions expected: nil error, actual: %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("getDivisions expected: len(servers) == 1, actual: %v", len(servers))
	}

}

type SortableDivisions []tc.Division

func (s SortableDivisions) Len() int {
	return len(s)
}
func (s SortableDivisions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableDivisions) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
