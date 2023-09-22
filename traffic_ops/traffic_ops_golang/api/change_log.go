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
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
)

type ChangeLog struct {
	ID          int          `json:"id" db:"id"`
	Level       string       `json:"level" db:"level"`
	Message     string       `json:"message" db:"message"`
	TMUser      int          `json:"tmUser" db:"tm_user"`
	TicketNum   string       `json:"ticketNum" db:"ticketnum"`
	LastUpdated tc.TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type ChangeLogger interface {
	ChangeLogMessage(action string) (string, error)
}

const (
	ApiChange = "APICHANGE"
	Updated   = "Updated"
	Created   = "Created"
	Deleted   = "Deleted"
)

func CreateChangeLog(level string, action string, i Identifier, user *auth.CurrentUser, tx *sql.Tx) error {
	t, ok := i.(ChangeLogger)
	if !ok {
		keys, _ := i.GetKeys()
		return CreateChangeLogBuildMsg(level, action, user, tx, i.GetType(), i.GetAuditName(), keys)
	}
	msg, err := t.ChangeLogMessage(action)
	if err != nil {
		log.Errorf("%++v creating log message for %++v", err, t)
		keys, _ := i.GetKeys()
		return CreateChangeLogBuildMsg(level, action, user, tx, i.GetType(), i.GetAuditName(), keys)
	}
	return CreateChangeLogRawErr(level, msg, user, tx)
}

func CreateChangeLogBuildMsg(level string, action string, user *auth.CurrentUser, tx *sql.Tx, objType string, auditName string, keys map[string]interface{}) error {
	keyStr := "{ "
	for key, value := range keys {
		keyStr += key + ":" + fmt.Sprintf("%v", value) + " "
	}
	keyStr += "}"
	id, ok := keys["id"]
	if !ok {
		id = "N/A"
	}
	msg := fmt.Sprintf("%v: %v, ID: %v, ACTION: %v %v, keys: %v", strings.ToTitle(objType), auditName, id, strings.Title(action), objType, keyStr)
	return CreateChangeLogRawErr(level, msg, user, tx)
}

func CreateChangeLogRawErr(level string, msg string, user *auth.CurrentUser, tx *sql.Tx) error {
	if _, err := tx.Exec(`INSERT INTO log (level, message, tm_user) VALUES ($1, $2, $3)`, level, msg, user.ID); err != nil {
		return errors.New("Inserting change log level '" + level + "' message '" + msg + "' user '" + user.UserName + "': " + err.Error())
	}
	return nil
}

func CreateChangeLogRawTx(level string, msg string, user *auth.CurrentUser, tx *sql.Tx) {
	if _, err := tx.Exec(`INSERT INTO log (level, message, tm_user) VALUES ($1, $2, $3)`, level, msg, user.ID); err != nil {
		log.Errorln("Inserting change log level '" + level + "' message '" + msg + "' user '" + user.UserName + "': " + err.Error())
	}
}
