package api

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
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type tester struct {
	ID        int
	error     error           //only for testing
	errorType tc.ApiErrorType //only for testing
}

//Identifier interface functions
func (i *tester) GetID() int {
	return i.ID
}

func (i *tester) GetType() string {
	return "tester"
}

func (i *tester) GetAuditName() string {
	return "testerInstance:" + strconv.Itoa(i.ID)
}

//Validator interface function
func (v *tester) Validate(db *sqlx.DB) []error {
	if v.ID < 1 {
		return []error{errors.New("ID is too low")}
	}
	return []error{}
}

//Inserter interface functions
func (i *tester) Insert(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	return i.error, i.errorType
}

func (i *tester) SetID(newID int) {
	i.ID = newID
}

//Updater interface functions
func (i *tester) Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	return i.error, i.errorType
}

//Deleter interface functions
func (i *tester) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	return i.error, i.errorType
}

//used for testing purposes only
func (t *tester) SetError(newError error, newErrorType tc.ApiErrorType) {
	t.error = newError
	t.errorType = newErrorType
}

func TestCreateHandler(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "", strings.NewReader(`{"ID":1}`))
	if err != nil {
		t.Error("Error creating new request")
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.CurrentUserKey,
		auth.CurrentUser{UserName: "username", ID: 1, PrivLevel: auth.PrivLevelAdmin})

	// Add our context to the request
	r = r.WithContext(ctx)

	typeRef := tester{}
	createFunc := CreateHandler(&typeRef, db)

	//verifies we get the right changelog insertion
	expectedMessage := Created + " " + typeRef.GetType() + ": " + typeRef.GetAuditName() + " id: 1"
	mock.ExpectExec("INSERT").WithArgs(ApiChange, expectedMessage, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	createFunc(w, r)

	//verifies the body is in the expected format
	body := `{"response":{"ID":1},"alerts":[{"text":"tester was created.","level":"success"}]}`
	if w.Body.String() != body {
		t.Error("Expected body", body, "got", w.Body.String())
	}
}

func TestUpdateHandler(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "", strings.NewReader(`{"ID":1}`))
	if err != nil {
		t.Error("Error creating new request")
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.CurrentUserKey,
		auth.CurrentUser{UserName: "username", ID: 1, PrivLevel: auth.PrivLevelAdmin})
	ctx = context.WithValue(ctx, PathParamsKey, PathParams{"id": "1"})
	// Add our context to the request
	r = r.WithContext(ctx)

	typeRef := tester{}
	updateFunc := UpdateHandler(&typeRef, db)

	//verifies we get the right changelog insertion
	expectedMessage := Updated + " " + typeRef.GetType() + ": " + typeRef.GetAuditName() + " id: 1"
	mock.ExpectExec("INSERT").WithArgs(ApiChange, expectedMessage, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	updateFunc(w, r)

	//verifies the body is in the expected format
	body := `{"response":{"ID":1},"alerts":[{"text":"tester was updated.","level":"success"}]}`
	if w.Body.String() != body {
		t.Error("Expected body", body, "got", w.Body.String())
	}
}

func TestDeleteHandler(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "", strings.NewReader(`{"ID":1}`))
	if err != nil {
		t.Error("Error creating new request")
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, auth.CurrentUserKey,
		auth.CurrentUser{UserName: "username", ID: 1, PrivLevel: auth.PrivLevelAdmin})
	ctx = context.WithValue(ctx, PathParamsKey, PathParams{"id": "1"})
	// Add our context to the request
	r = r.WithContext(ctx)

	typeRef := tester{}
	deleteFunc := DeleteHandler(&typeRef, db)

	//verifies we get the right changelog insertion
	expectedMessage := Deleted + " " + typeRef.GetType() + ": " + typeRef.GetAuditName() + " id: 1"
	mock.ExpectExec("INSERT").WithArgs(ApiChange, expectedMessage, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	deleteFunc(w, r)

	//verifies the body is in the expected format
	body := `{"alerts":[{"text":"tester was deleted.","level":"success"}]}`
	if w.Body.String() != body {
		t.Error("Expected body", body, "got", w.Body.String())
	}
}
