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

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

type ChangeLogger interface {
	ChangeLogMessage(action string) (string, error)
}

const (
	ApiChange = "APICHANGE"
	Updated   = "Updated"
	Created   = "Created"
	Deleted   = "Deleted"
)

func CreateChangeLog(action string, i Identifier, user *auth.CurrentUser, tx *sql.Tx) error {
	t, ok := i.(ChangeLogger)
	if !ok {
		keys, _ := i.GetKeys()
		return CreateChangeLogBuildMsg(action, user, tx, i.GetType(), i.GetAuditName(), keys)
	}
	msg, err := t.ChangeLogMessage(action)
	if err != nil {
		log.Errorf("%++v creating log message for %++v", err, t)
		keys, _ := i.GetKeys()
		return CreateChangeLogBuildMsg(action, user, tx, i.GetType(), i.GetAuditName(), keys)
	}
	return CreateChangeLogRawErr(msg, user, tx)
}

func CreateChangeLogBuildMsg(action string, user *auth.CurrentUser, tx *sql.Tx, objType string, auditName string, keys map[string]interface{}) error {
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
	return CreateChangeLogRawErr(msg, user, tx)
}

func CreateChangeLogRawErr(msg string, user *auth.CurrentUser, tx *sql.Tx) error {
	if _, err := tx.Exec(`INSERT INTO log (message, "user") VALUES ($1, $2)`, msg, user.ID); err != nil {
		return errors.New("Inserting change log message '" + msg + "' user '" + user.UserName + "': " + err.Error())
	}
	return nil
}

func CreateChangeLogRawTx(msg string, user *auth.CurrentUser, tx *sql.Tx) {
	if _, err := tx.Exec(`INSERT INTO log (message, "user") VALUES ($1, $2)`, msg, user.ID); err != nil {
		log.Errorln("Inserting change log message '" + msg + "' user '" + user.UserName + "': " + err.Error())
	}
}
