package servicecategory

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
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetServiceCategory(t *testing.T) {

	type testStruct struct {
		useIms   bool
		expected int
	}

	var testData = []testStruct{
		// When useIMS is set to false in system Config
		{false, 200},
		// When useIMS is set to true in system Config
		{true, 200},
	}

	var header http.Header
	lastUpdated := time.Now()
	params := map[string]string{}

	for i, _ := range testData {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		db := sqlx.NewDb(mockDB, "sqlmock")
		defer db.Close()

		header = make(http.Header)
		ifModifiedSince := time.Now().AddDate(0, 00, 01)
		timeValue := ifModifiedSince.Format("Mon, 02 Jan 2006 15:04:05 MST")
		header.Set(rfc.IfModifiedSince, timeValue)

		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"name", "last_updated"})
		rows.AddRow("testObj1", lastUpdated.AddDate(0, 0, -5))
		mock.ExpectQuery("SELECT name, last_updated FROM service_category").WithArgs().WillReturnRows(rows)

		_, _, code, _, _ := GetServiceCategory(db.MustBegin(), params, testData[i].useIms, header)

		if testData[i].expected != code {
			t.Errorf("Expected status code result doesn't match, got: %v; expected: %v", code, testData[i].expected)
		}

	}

}
