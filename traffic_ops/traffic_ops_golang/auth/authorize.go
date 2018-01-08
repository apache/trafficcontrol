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
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/jmoiron/sqlx"
)

type CurrentUser struct {
	UserName  string `json:"userName" db:"username"`
	ID        int    `json:"id" db:"id"`
	PrivLevel int    `json:"privLevel" db:"priv_level"`
	TenantID  int    `json:"tenantID" db:"tenant_id"`
}

// PrivLevelInvalid - The Default Priv level
const PrivLevelInvalid = -1

// PrivLevelReadOnly - The user cannot do any API updates
const PrivLevelReadOnly = 10

// PrivLevelOperations - The user has minimal privileges
const PrivLevelOperations = 20

// PrivLevelAdmin - The user has full privileges
const PrivLevelAdmin = 30

// TenantIDInvalid - The default Tenant ID
const TenantIDInvalid = -1

const CurrentUserKey = "currentUser"

// GetCurrentUserFromDB  - returns the id and privilege level of the given user along with the username, or -1 as the id, - as the userName and PrivLevelInvalid if the user doesn't exist.
func GetCurrentUserFromDB(CurrentUserStmt *sqlx.Stmt, user string) CurrentUser {
	var currentUserInfo CurrentUser
	err := CurrentUserStmt.Get(&currentUserInfo, user)
	switch {
	case err == sql.ErrNoRows:
		log.Errorf("checking user %v info: user not in database", user)
		return CurrentUser{"-", -1, PrivLevelInvalid, TenantIDInvalid}
	case err != nil:
		log.Errorf("Error checking user %v info: %v", user, err.Error())
		return CurrentUser{"-", -1, PrivLevelInvalid, TenantIDInvalid}
	default:
		return currentUserInfo
	}
}

func GetCurrentUser(ctx context.Context) (*CurrentUser, error) {
	val := ctx.Value(CurrentUserKey)
	if val != nil {
		switch v := val.(type) {
		case CurrentUser:
			return &v, nil
		default:
			return nil, fmt.Errorf("CurrentUser found with bad type: %T", v)
		}
	}
	return &CurrentUser{"-", -1, PrivLevelInvalid, TenantIDInvalid}, errors.New("No user found in Context")
}
