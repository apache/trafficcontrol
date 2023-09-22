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
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/lib/pq"
)

const (
	getUsersQuery = `
		SELECT
			u.id,
			u.local_passwd,
			u.role,
			u.tenant_id,
			u.token,
			u.ucdn,
			u.username
		FROM
			tm_user AS u
	`
	getRolesQuery = `
		SELECT
			ARRAY(SELECT rc.cap_name FROM role_capability AS rc WHERE rc.role_id=r.id) AS capabilities,
			r.id as role,
			r.name as role_name,
			r.priv_level
		FROM role r
	`
)

type user struct {
	CurrentUser
	LocalPasswd *string
	Token       *string
}

type role struct {
	Capabilities pq.StringArray
	ID           int
	Name         string
	PrivLevel    int
}

type users struct {
	userMap          map[string]user
	usernamesByToken map[string]string
	*sync.RWMutex
	initialized bool
	enabled     bool // note: enabled is only written to once at startup, before serving requests, so it doesn't need synchronized access
}

var usersCache = users{RWMutex: &sync.RWMutex{}}

func usersCacheIsEnabled() bool {
	if usersCache.enabled {
		usersCache.RLock()
		defer usersCache.RUnlock()
		return usersCache.initialized
	}
	return false
}

// getUserFromCache returns the user with the given username and a boolean indicating whether the user exists.
func getUserFromCache(username string) (user, bool) {
	usersCache.RLock()
	defer usersCache.RUnlock()
	u, exists := usersCache.userMap[username]
	return u, exists
}

// getUserNameFromCacheByToken returns the username with the given token and a boolean indicating whether a matching token was found.
func getUserNameFromCacheByToken(token string) (string, bool) {
	usersCache.RLock()
	defer usersCache.RUnlock()
	t, exists := usersCache.usernamesByToken[token]
	return t, exists
}

var once = sync.Once{}

// InitUsersCache attempts to initialize the in-memory users data (if enabled) then
// starts a goroutine to periodically refresh the in-memory data from the database.
func InitUsersCache(interval time.Duration, db *sql.DB, timeout time.Duration) {
	once.Do(func() {
		if interval <= 0 {
			return
		}
		usersCache.enabled = true
		refreshUsersCache(db, timeout)
		startUsersCacheRefresher(interval, db, timeout)
	})
}

func startUsersCacheRefresher(interval time.Duration, db *sql.DB, timeout time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			refreshUsersCache(db, timeout)
		}
	}()
}

func refreshUsersCache(db *sql.DB, timeout time.Duration) {
	newUsers, err := getUsers(db, timeout)
	if err != nil {
		log.Errorf("refreshing users cache: %s", err.Error())
		return
	}
	usersCache.Lock()
	defer usersCache.Unlock()
	usersCache.userMap = newUsers
	usersCache.usernamesByToken = createTokenToUsernameMap(newUsers)
	usersCache.initialized = true
	log.Infof("refreshed users cache (len = %d)", len(usersCache.userMap))
}

func createTokenToUsernameMap(users map[string]user) map[string]string {
	tokenToUserName := make(map[string]string)
	for username, u := range users {
		if u.Token == nil || u.RoleName == disallowed {
			continue
		}
		tokenToUserName[*u.Token] = username
	}
	return tokenToUserName
}

func getUsers(db *sql.DB, timeout time.Duration) (map[string]user, error) {
	dbCtx, dbClose := context.WithTimeout(context.Background(), timeout)
	defer dbClose()
	roles := make(map[int]role)
	newUsers := make(map[string]user)

	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		return nil, errors.New("beginning users transaction: " + err.Error())
	}
	defer func() {
		if err := tx.Commit(); err != nil && err != sql.ErrTxDone {
			log.Errorln("committing users transaction: " + err.Error())
		}
	}()
	rolesRows, err := tx.QueryContext(dbCtx, getRolesQuery)
	if err != nil {
		return nil, errors.New("querying roles: " + err.Error())
	}
	defer log.Close(rolesRows, "closing role rows")
	for rolesRows.Next() {
		r := role{}
		if err := rolesRows.Scan(&r.Capabilities, &r.ID, &r.Name, &r.PrivLevel); err != nil {
			return nil, errors.New("scanning roles: " + err.Error())
		}
		roles[r.ID] = r
	}
	if err = rolesRows.Err(); err != nil {
		return nil, errors.New("iterating over role rows: " + err.Error())
	}

	rows, err := tx.QueryContext(dbCtx, getUsersQuery)
	if err != nil {
		return nil, errors.New("querying users: " + err.Error())
	}
	defer log.Close(rows, "closing users rows")
	for rows.Next() {
		u := user{}
		if err := rows.Scan(&u.ID, &u.LocalPasswd, &u.Role, &u.TenantID, &u.Token, &u.UCDN, &u.UserName); err != nil {
			return nil, errors.New("scanning users: " + err.Error())
		}
		r := roles[u.Role]
		u.RoleName = r.Name
		u.PrivLevel = r.PrivLevel
		u.Capabilities = r.Capabilities
		u.perms = make(map[string]struct{}, len(u.Capabilities))
		for _, perm := range u.Capabilities {
			u.perms[perm] = struct{}{}
		}
		newUsers[u.UserName] = u
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New("iterating over user rows: " + err.Error())
	}
	return newUsers, nil
}
