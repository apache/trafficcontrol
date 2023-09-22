package auth

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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("creating new sqlmock: %v", err)
	}
	defer db.Close()
	expectedRoles := []role{
		{
			Capabilities: []string{"foo", "bar"},
			ID:           1,
			Name:         "foo_role",
			PrivLevel:    42,
		},
	}
	expectedUsers := map[string]user{
		"user1": {
			CurrentUser: CurrentUser{
				UserName:     "user1",
				ID:           1,
				PrivLevel:    42,
				TenantID:     1,
				Role:         1,
				RoleName:     "foo_role",
				Capabilities: []string{"foo", "bar"},
				UCDN:         "ucdn1",
				perms: map[string]struct{}{
					"foo": {},
					"bar": {},
				},
			},
			LocalPasswd: util.StrPtr("foo"),
			Token:       util.StrPtr("bar"),
		},
	}
	roleRows := sqlmock.NewRows([]string{"capabilities", "role", "role_name", "priv_level"})
	userRows := sqlmock.NewRows([]string{"id", "local_passwd", "role", "tenant_id", "token", "ucdn", "username"})

	for _, r := range expectedRoles {
		roleRows.AddRow("{"+strings.Join(r.Capabilities, ",")+"}", r.ID, r.Name, r.PrivLevel)
	}
	for _, u := range expectedUsers {
		userRows.AddRow(u.ID, u.LocalPasswd, u.Role, u.TenantID, u.Token, u.UCDN, u.UserName)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT.+").WillReturnRows(roleRows)
	mock.ExpectQuery("SELECT.+").WillReturnRows(userRows)
	mock.ExpectCommit()
	actualUsers, err := getUsers(db, 10*time.Second)
	if err != nil {
		t.Fatalf("getUsers expected: nil error, actual: %v", err)
	}
	if !reflect.DeepEqual(expectedUsers, actualUsers) {
		t.Errorf("getUsers expected: %v, actual: %v", expectedUsers, actualUsers)
	}
}
