package trafficstats

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

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestQueryStatsSummary(t *testing.T) {
	type testStruct struct {
		version uint64
	}

	var testData = []testStruct{
		{4},
		{5},
	}

	query := "SELECT cdn_name, deliveryservice_name, stat_name, stat_value, summary_time, stat_date FROM stats_summary"
	queryValues := map[string]interface{}{
		"lastSummaryDate": "true",
	}

	for i, _ := range testData {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()

		db := sqlx.NewDb(mockDB, "sqlmock")
		defer db.Close()

		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{
			"cdn_name",
			"deliveryservice_name",
			"stat_name",
			"stat_value",
			"summary_time",
			"stat_date",
		})

		rows.AddRow("cdn1", "all", "daily_maxgbps", 5, time.Now().AddDate(0, 0, -5), time.Now().AddDate(0, 0, -5).Truncate(24*time.Hour))
		rows.AddRow("cdn2", "all", "daily_byteserved", 1000, time.Now().AddDate(0, 0, -10), time.Now().AddDate(0, 0, -10).Truncate(24*time.Hour))

		mock.ExpectQuery("SELECT cdn_name, deliveryservice_name, stat_name, stat_value, summary_time, stat_date FROM stats_summary").WithArgs().WillReturnRows(rows)

		statsSummaries1, err1 := queryStatsSummary(db.MustBegin(), testData[i].version, query, queryValues)

		assert.NoError(t, err1)
		assert.Equal(t, len(statsSummaries1), 2)
	}

}
