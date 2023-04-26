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
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetDomainsList(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{"id", "name", "description", "domain_name"}
	rows := sqlmock.NewRows(cols)
	rows.AddRow(1, "profile1", "profiledesc1", "profiledomain1")
	rows.AddRow(2, "profile2", "profiledesc2", "profiledomain2")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT p.id").WillReturnRows(rows)
	mock.ExpectCommit()

	domainList, err, sc, _ := getDomainsList(false, nil, db.MustBegin())
	if err != nil {
		t.Fatalf("expected no error while getting domains list, but got: %v", err)
	}
	if sc != http.StatusOK {
		t.Errorf("expected a 200 status, but got %d", sc)
	}
	if len(domainList) != 2 {
		t.Fatalf("expected domains to have a length of 2, but got %d", len(domainList))
	}
	if domainList[0].ProfileID != 1 || domainList[0].ProfileName != "profile1" ||
		domainList[0].ProfileDescription != "profiledesc1" || domainList[0].DomainName != "profiledomain1" {
		t.Errorf("expected: profile ID: 1, profile name: profile1, profile desc: profiledesc1, profile domain: profiledomain1; got: %d, %s, %s, %s",
			domainList[0].ProfileID, domainList[0].ProfileName, domainList[0].ProfileDescription, domainList[0].DomainName)
	}
	if domainList[1].ProfileID != 2 || domainList[1].ProfileName != "profile2" ||
		domainList[1].ProfileDescription != "profiledesc2" || domainList[1].DomainName != "profiledomain2" {
		t.Errorf("expected: profile ID: 2, profile name: profile2, profile desc: profiledesc2, profile domain: profiledomain2; got: %d, %s, %s, %s",
			domainList[1].ProfileID, domainList[1].ProfileName, domainList[1].ProfileDescription, domainList[1].DomainName)
	}
}
