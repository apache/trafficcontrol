package asn

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
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestASNs() []tc.ASN {
	ASNs := []tc.ASN{}
	testCase := tc.ASN{
		ASN:         1,
		Cachegroup:  "Yukon",
		ID:          1,
		LastUpdated: tc.Time{Time: time.Now()},
	}
	ASNs = append(ASNs, testCase)

	testCase2 := testCase
	testCase2.ASN = 2
	ASNs = append(ASNs, testCase2)

	return ASNs
}

func TestGetASNs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testCase := getTestASNs()
	cols := test.ColsFromStructByTag("db", tc.ASN{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testCase {
		rows = rows.AddRow(
			ts.ASN,
			ts.Cachegroup,
			ts.CachegroupID,
			ts.ID,
			ts.LastUpdated,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := map[string]string{"dsId": "1"}

	servers, errs, errType := getASNs(v, db)
	if len(errs) > 0 {
		t.Errorf("getASNs expected: no errors, actual: %v with type %s", errs, errType.String())
	}

	if len(servers) != 2 {
		t.Errorf("getASNs expected: len(servers) == 1, actual: %v", len(servers))
	}

}

type SortableASNs []tc.ASN

func (s SortableASNs) Len() int {
	return len(s)
}
func (s SortableASNs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableASNs) Less(i, j int) bool {
	return s[i].ASN < s[j].ASN
}
