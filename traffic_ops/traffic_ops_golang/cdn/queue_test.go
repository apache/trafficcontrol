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

func TestQueueUpdates(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE public.server").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	_, err = queueUpdates(db.MustBegin(), "UPDATE public.server SET config_update_time = now()", nil)
	if err != nil {
		t.Fatalf("expected no error while queueing updates, but got: %v", err)
	}
}
