package deliveryservice

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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestUpdateDSSafe(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	// test with a DS that exists
	mock.ExpectBegin()
	dsID := 1
	dsr := tc.DeliveryServiceSafeUpdateRequest{
		DisplayName: util.Ptr("displayName"),
		InfoURL:     util.Ptr("http://blah.com"),
		LongDesc:    util.Ptr("longdesc"),
		LongDesc1:   util.Ptr("longdesc1"),
	}
	mock.ExpectExec("UPDATE deliveryservice").WithArgs(*dsr.DisplayName, *dsr.InfoURL, *dsr.LongDesc, *dsr.LongDesc1, dsID).WillReturnResult(sqlmock.NewResult(int64(dsID), 1))
	exists, err := updateDSSafe(db.MustBegin().Tx, dsID, dsr, false)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if !exists {
		t.Errorf("expected DS with id 1 to exist")
	}

	// test with a DS that doesn't exist
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE deliveryservice").WithArgs(*dsr.DisplayName, *dsr.InfoURL, *dsr.LongDesc, *dsr.LongDesc1, 2).WillReturnResult(sqlmock.NewResult(2, 0))
	exists, err = updateDSSafe(db.MustBegin().Tx, 2, dsr, false)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if exists {
		t.Errorf("expected DS with id 2 to not exist")
	}
}
