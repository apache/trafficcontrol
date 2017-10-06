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

func getTestHWInfo() []tc.HWInfo {
	hwinfo := []tc.HWInfo{}
	testHWInfo := tc.HWInfo{
		ID:          1,
		ServerID:    1,
		Description: "Description",
		Val:         "Val",
		LastUpdated: "LastUpdated",
	}
	hwinfo = append(hwinfo, testHWInfo)

	testHWInfo2 := testHWInfo
	testHWInfo2.Description = "hwinfo2"
	testHWInfo2.Val = "val2"
	testHWInfo2.ServerID = 2
	hwinfo = append(hwinfo, testHWInfo2)

	return hwinfo
}

func TestGetHWInfo(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testHWInfo := getTestHWInfo()
	cols := test.ColsFromStructByTag("db", tc.HWInfo{})
	rows := sqlmock.NewRows(cols)

	//TODO: drichardson - build helper to add these Rows from the struct values
	//                    or by CSV if types get in the way
	for _, ts := range testHWInfo {
		rows = rows.AddRow(
			ts.Description,
			ts.ID,
			ts.LastUpdated,
			ts.ServerID,
			ts.Val,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := url.Values{}
	v.Set("ServerId", "1")

	hwinfos, err := getHWInfo(v, db, PrivLevelAdmin)
	if err != nil {
		t.Errorf("getHWInfo expected: nil error, actual: %v", err)
	}

	if len(hwinfos) != 2 {
		t.Errorf("getHWInfo expected: len(hwinfos) == 1, actual: %v", len(hwinfos))
	}
}

type SortableHWInfo []tc.HWInfo

func (s SortableHWInfo) Len() int {
	return len(s)
}
func (s SortableHWInfo) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableHWInfo) Less(i, j int) bool {
	return s[i].Description < s[j].Description
}
