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
	"strconv"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
)

type testIdentifier struct {
	ID int
}

func (i testIdentifier) GetKeyFieldsInfo() []KeyFieldInfo {
	return []KeyFieldInfo{{"id", GetIntKey}}
}

func (i *testIdentifier) SetKeys(keys map[string]interface{}) {
	id, _ := keys["id"].(int)
	i.ID = id
}

func (i *testIdentifier) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{"id": i.ID}, true
}

func (i *testIdentifier) GetType() string {
	return "tester"
}

func (i *testIdentifier) GetAuditName() string {
	return "testerInstance"
}

func TestCreateChangeLog(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	i := testIdentifier{}

	keys, _ := i.GetKeys()
	expectedMessage := strings.ToUpper(i.GetType()) + ": " + i.GetAuditName() + ", ID: " + strconv.Itoa(keys["id"].(int)) + ", ACTION: " + Created + " " + i.GetType() + ", keys: { id:" + strconv.Itoa(keys["id"].(int)) + " }"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT").WithArgs(ApiChange, expectedMessage, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	user := auth.CurrentUser{ID: 1}
	err = CreateChangeLog(ApiChange, Created, &i, &user, db.MustBegin().Tx)
	if err != nil {
		t.Fatal(err)
	}
}
