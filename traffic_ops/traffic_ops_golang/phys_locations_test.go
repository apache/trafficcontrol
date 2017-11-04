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

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestPhysLocations() []tc.PhysLocation {
	physLocations := []tc.PhysLocation{}
	testCase := tc.PhysLocation{
		ID:           1,
		Name:         "physLocation1",
		ShortName:    "pl1",
		Address:      "1118 S. Grant St.",
		City:         "Denver",
		State:        "CO",
		Zip:          "80210",
		RegionId:     1,
        POC:          "Dennis Thompson",
        Phone:        "303-210-0000",
        Email:        "d.t@gmail.com",
        Comments:     "",
        RegionName:   "Central",
		LastUpdated:  "2015-12-10 15:43:45-07",
	}
	physLocations = append(physLocations, testCase)

	testCase2 := testCase
	testCase2.Name = "physLocation2"
	physLocations = append(physLocations, testCase2)

	return physLocations
}

func TestGetPhysLocations(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testCase := getTestPhysLocations()
	cols := test.ColsFromStructByTag("db", tc.PhysLocation{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testCase {
		rows = rows.AddRow(
			ts.ID,
			ts.Name,
			ts.ShortName,
			ts.Address,
			ts.City,
			ts.State,
			ts.Zip,
			ts.RegionId,
			ts.POC,
			ts.Phone,
			ts.Email,
			ts.Comments,
			ts.RegionName,
			ts.LastUpdated,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := url.Values{}
	v.Set("dsId", "1")

	servers, err := getPhysLocations(v, db)
	if err != nil {
		t.Errorf("getPhysLocations expected: nil error, actual: %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("getPhysLocations expected: len(servers) == 2, actual: %v", len(servers))
	}

}

type SortablePhysLocations []tc.PhysLocation

func (s SortablePhysLocations) Len() int {
	return len(s)
}
func (s SortablePhysLocations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortablePhysLocations) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
