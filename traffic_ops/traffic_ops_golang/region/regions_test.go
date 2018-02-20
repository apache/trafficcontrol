package region

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
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestRegions() []tc.Region {
	regions := []tc.Region{}
	testCase := tc.Region{
		DivisionName: "west",
		ID:           1,
		Name:         "region1",
		LastUpdated:  tc.Time{Time: time.Now()},
	}
	regions = append(regions, testCase)

	testCase2 := testCase
	testCase2.Name = "region2"
	regions = append(regions, testCase2)

	return regions
}

func TestReadRegions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	refType := GetRefType()

	testRegions := getTestRegions()
	cols := test.ColsFromStructByTag("db", tc.Region{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testRegions {
		rows = rows.AddRow(
			ts.Division,
			ts.ID,
			ts.LastUpdated,
			ts.Name,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := map[string]string{"dsId": "1"}

	servers, errs, _ := refType.Read(db, v, auth.CurrentUser{})
	if len(errs) > 0 {
		t.Errorf("cdn.Read expected: no errors, actual: %v", errs)
	}

	if len(servers) != 2 {
		t.Errorf("cdn.Read expected: len(servers) == 2, actual: %v", len(servers))
	}
}
