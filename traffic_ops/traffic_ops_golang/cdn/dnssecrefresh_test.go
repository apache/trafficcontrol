package cdn

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

	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetDNSSECKeyRefreshParams_test(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{"cdn_name", "cdn_domain", "cdn_dnssec_enabled", "parameter_name", "parameter_value"}
	rows := sqlmock.NewRows(cols)
	rows.AddRow("test", "test.com", false, "", "")

	mock.ExpectBegin()
	mock.ExpectQuery("WITH cdn_profile_ids").WillReturnRows(rows)
	mock.ExpectCommit()

	params, err := getDNSSECKeyRefreshParams(db.MustBegin().Tx)
	for _, v := range params {
		if v.CDNName != "test" {
			t.Errorf("Expected cdn name: test, got: %s", v.CDNName)
		}
		if v.CDNDomain != "test.com" {
			t.Errorf("Expected cdn domain: test.com, got: %s", v.CDNDomain)
		}
		if v.DNSSECEnabled != false {
			t.Errorf("Expected DNSSEC to not be enabled")
		}
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
	}

}
