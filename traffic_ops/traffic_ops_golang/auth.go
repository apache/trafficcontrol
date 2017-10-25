package main

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
	"database/sql"
	"errors"
	"fmt"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/jmoiron/sqlx"
)

// PrivLevelInvalid - The Default Priv level
const PrivLevelInvalid = -1

// PrivLevelReadOnly - The user cannot do any API updates
const PrivLevelReadOnly = 10

// PrivLevelOperations - The user has minimal privileges
const PrivLevelOperations = 20

// PrivLevelAdmin - The user has full privileges
const PrivLevelAdmin = 30

const PrivLevelKey = "privLevel"
const UserNameKey = "userName"

func preparePrivLevelStmt(db *sqlx.DB) (*sql.Stmt, error) {
	return db.Prepare("SELECT r.priv_level FROM tm_user AS u JOIN role AS r ON u.role = r.id WHERE u.username = $1")
}

// PrivLevel - returns the privilege level of the given user, or PrivLevelInvalid if the user doesn't exist.
func PrivLevel(privLevelStmt *sql.Stmt, user string) int {
	var privLevel int
	err := privLevelStmt.QueryRow(user).Scan(&privLevel)
	switch {
	case err == sql.ErrNoRows:
		log.Errorf("checking user %v priv level: user not in database", user)
		return PrivLevelInvalid
	case err != nil:
		log.Errorf("Error checking user %v priv level: %v", user, err.Error())
		return PrivLevelInvalid
	default:
		return privLevel
	}
}

func getPrivLevel(ctx context.Context) (int, error) {
	val := ctx.Value(PrivLevelKey)
	if val != nil {
		switch v := val.(type) {
		case int:
			return v, nil
		default:
			return PrivLevelInvalid, fmt.Errorf("privLevel found with bad type: %T\n", v)
		}
	}
	return PrivLevelInvalid, errors.New("no privLevel found in Context")
}

func getUserName(ctx context.Context) (string, error) {
	val := ctx.Value(UserNameKey)
	if val != nil {
		switch v := val.(type) {
		case string:
			return v, nil
		default:
			return "-", fmt.Errorf("userName found with bad type: %T\n", v)
		}
	}
	return "-", errors.New("No userName found in Context")
}
