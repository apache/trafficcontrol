package request

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

func TestGetAssignee(t *testing.T) {
	req := assignmentRequest{
		AssigneeID: nil,
		Assignee:   nil,
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("opening mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	// check simple case, no Assignee or ID means no change
	mock.ExpectBegin()
	_, _, userErr, sysErr := getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}
	if req.AssigneeID != nil {
		t.Errorf("assignee ID was somehow set to: %d", *req.AssigneeID)
	}
	if req.Assignee != nil {
		t.Errorf("assignee was somehow set to: %s", *req.Assignee)
	}

	expectID := 12
	expectName := "test assignee"

	req.Assignee = &expectName

	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(expectID)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id").WillReturnRows(rows)

	// check case where getting Assignee ID from username
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}

	if req.Assignee == nil {
		t.Error("Expected assignee to not be nil after getting assignee")
	} else if *req.Assignee != expectName {
		t.Errorf("Incorrect assignee; expected: '%s', got: '%s'", expectName, *req.Assignee)
	}

	if req.AssigneeID == nil {
		t.Error("Expected assignee ID to not be nil after getting assignee")
	} else if *req.AssigneeID != expectID {
		t.Errorf("Incorrect assignee ID; expected: %d, got: %d", expectID, *req.AssigneeID)
	}

	req.Assignee = nil
	req.AssigneeID = &expectID

	rows = sqlmock.NewRows([]string{"username"})
	rows.AddRow(expectName)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT username").WillReturnRows(rows)

	// check case where getting username from Assignee ID
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}

	if req.Assignee == nil {
		t.Error("Expected assignee to not be nil after getting assignee")
	} else if *req.Assignee != expectName {
		t.Errorf("Incorrect assignee; expected: '%s', got: '%s'", expectName, *req.Assignee)
	}

	if req.AssigneeID == nil {
		t.Error("Expected assignee ID to not be nil after getting assignee")
	} else if *req.AssigneeID != expectID {
		t.Errorf("Incorrect assignee ID; expected: %d, got: %d", expectID, *req.AssigneeID)
	}

	req.Assignee = new(string)
	*req.Assignee = expectName + " - but not actually"
	req.AssigneeID = &expectID
	rows = sqlmock.NewRows([]string{"username"})
	rows.AddRow(expectName)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT username").WillReturnRows(rows)

	// check that Assignee ID has precedence over Assignee
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr != nil {
		t.Errorf("unexpected user error: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("unexpected system error: %v", sysErr)
	}

	if req.Assignee == nil {
		t.Error("Expected assignee to not be nil after getting assignee")
	} else if *req.Assignee != expectName {
		t.Errorf("Incorrect assignee; expected: '%s', got: '%s'", expectName, *req.Assignee)
	}

	if req.AssigneeID == nil {
		t.Error("Expected assignee ID to not be nil after getting assignee")
	} else if *req.AssigneeID != expectID {
		t.Errorf("Incorrect assignee ID; expected: %d, got: %d", expectID, *req.AssigneeID)
	}

	req.Assignee = nil
	req.AssigneeID = &expectID
	rows = sqlmock.NewRows([]string{"username"})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT username").WillReturnRows(rows)

	// check that looking for ID of non-existent username is an error
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr == nil {
		t.Error("Expected a user error, but didn't get one")
	}

	req.Assignee = &expectName
	req.AssigneeID = nil
	rows = sqlmock.NewRows([]string{"id"})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id").WillReturnRows(rows)

	// check that looking for username of non-existent Assignee is an error
	_, _, userErr, sysErr = getAssignee(&req, "test", db.MustBegin().Tx)
	if userErr == nil {
		t.Error("Expected a user error, but didn't get one")
	}
}
