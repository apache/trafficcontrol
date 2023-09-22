package totestv4

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
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

func CreateTestUsers(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, user := range td.Users {
		resp, _, err := cl.CreateUser(user, toclient.RequestOptions{})
		assert.RequireNoError(t, err, "Could not create user: %v - alerts: %+v", err, resp.Alerts)
	}
}

// ForceDeleteTestUsers forcibly deletes the users from the db.
// NOTE: Special circumstances!  This should *NOT* be done without a really good reason!
// Connects directly to the DB to remove users rather than going through the client.
// This is required here because the DeleteUser action does not really delete users,  but disables them.
func ForceDeleteTestUsers(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	var usernames []string
	for _, user := range td.Users {
		usernames = append(usernames, `'`+user.Username+`'`)
	}

	// there is a constraint that prevents users from being deleted when they have a log
	q := `DELETE FROM log WHERE NOT tm_user = (SELECT id FROM tm_user WHERE username = 'admin')`
	err := execSQL(db, q)
	assert.RequireNoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)

	q = `DELETE FROM tm_user WHERE username IN (` + strings.Join(usernames, ",") + `)`
	err = execSQL(db, q)
	assert.NoError(t, err, "Cannot execute SQL: %v; SQL is %s", err, q)
}

func GetUserID(t *testing.T, cl *toclient.Session, username string) func() int {
	return func() int {
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("username", username)
		users, _, err := cl.GetUsers(opts)
		assert.RequireNoError(t, err, "Get Users Request failed with error:", err)
		assert.RequireEqual(t, 1, len(users.Response), "Expected response object length 1, but got %d", len(users.Response))
		assert.RequireNotNil(t, users.Response[0].ID, "Expected ID to not be nil.")
		return *users.Response[0].ID
	}
}

func execSQL(db *sql.DB, sqlStmt string) error {
	var err error

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin failed %v %v ", err, tx)
	}

	res, err := tx.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("exec failed %v %v", err, res)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed %v %v", err, res)
	}
	return nil
}
