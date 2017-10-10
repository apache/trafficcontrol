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

func getTestParameters() []tc.Parameter {
	parameters := []tc.Parameter{}
	testParameter := tc.Parameter{
		ConfigFile:  "global",
		ID:          1,
		LastUpdated: "lastUpdated",
		Name:        "paramname1",
		Secure:      false,
		Value:       "val1",
	}
	parameters = append(parameters, testParameter)

	testParameter2 := testParameter
	testParameter2.Name = "paramname2"
	testParameter2.Value = "val2"
	testParameter2.ConfigFile = "some.config"
	parameters = append(parameters, testParameter2)

	return parameters
}

func TestGetParameters(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testParameters := getTestParameters()
	cols := test.ColsFromStructByTag("db", tc.Parameter{})
	rows := sqlmock.NewRows(cols)

	//TODO: drichardson - build helper to add these Rows from the struct values
	//                    or by CSV if types get in the way
	for _, ts := range testParameters {
		rows = rows.AddRow(
			ts.ConfigFile,
			ts.ID,
			ts.LastUpdated,
			ts.Name,
			ts.Secure,
			ts.Value,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := url.Values{}
	v.Set("dsId", "1")

	parameters, err := getParameters(v, db, PrivLevelAdmin)
	if err != nil {
		t.Errorf("getParameters expected: nil error, actual: %v", err)
	}

	if len(parameters) != 2 {
		t.Errorf("getParameters expected: len(parameters) == 1, actual: %v", len(parameters))
	}

}

type SortableParameters []tc.Parameter

func (s SortableParameters) Len() int {
	return len(s)
}
func (s SortableParameters) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableParameters) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
