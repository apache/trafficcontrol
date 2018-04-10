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
	"fmt"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
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

func CreateChangeLog(level string, action string, i Identifier, user auth.CurrentUser, db *sql.DB) error {
	t, ok := i.(ChangeLogger)
	if !ok {
		keys, _ := i.GetKeys()
		return CreateChangeLogBuildMsg(level, action, user, db, i.GetType(), i.GetAuditName(), keys)
	}
	msg, err := t.ChangeLogMessage(action)
	if err != nil {
		log.Errorf("%++v creating log message for %++v", err, t)
		keys, _ := i.GetKeys()
		return CreateChangeLogBuildMsg(level, action, user, db, i.GetType(), i.GetAuditName(), keys)
	}
	return CreateChangeLogMsg(level, user, db, msg)
}

func CreateChangeLogBuildMsg(level string, action string, user auth.CurrentUser, db *sql.DB, objType string, auditName string, keys map[string]interface{}) error {
	keyStr := "{ "
	for key, value := range keys {
		keyStr += key + ":" + fmt.Sprintf("%v", value) + " "
	}
	keyStr += "}"
	msg := action + " " + objType + ": " + auditName + " keys: " + keyStr
	return CreateChangeLogMsg(level, user, db, msg)
}

func CreateChangeLogMsg(level string, user auth.CurrentUser, db *sql.DB, msg string) error {
	query := `INSERT INTO log (level, message, tm_user) VALUES ($1, $2, $3)`
	log.Debugf("about to exec %s with %s", query, msg)
	if _, err := db.Exec(query, level, msg, user.ID); err != nil {
		log.Errorf("received error: %++v from audit log insertion", err)
		return err
	}
	return nil
}

func CreateChangeLogRaw(level string, message string, user auth.CurrentUser, db *sql.DB) error {
	if _, err := db.Exec(`INSERT INTO log (level, message, tm_user) VALUES ($1, $2, $3)`, level, message, user.ID); err != nil {
		return fmt.Errorf("inserting change log level '%v' message '%v' user '%v': %v", level, message, user.ID, err)
	}
	return nil
}
