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
	"fmt"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/jmoiron/sqlx"
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

func CreateChangeLog(level string, action string, i Identifier, user auth.CurrentUser, db *sqlx.DB) error {
	keys, _ := i.GetKeys()
	keysString := "{ "
	for key, value := range keys {
		keysString += key + ":" + fmt.Sprintf("%v", value) + " "
	}
	keysString += "}"
	message := action + " " + i.GetType() + ": " + i.GetAuditName() + " keys: " + keysString
	// if the object has its own log message generation, use it
	if t, ok := i.(ChangeLogger); ok {
		m, err := t.ChangeLogMessage(action)
		if err != nil {
			log.Errorf("error %++v creating log message for %++v", err, t)
			// use the default message in this case
		} else {
			message = m
		}
	}

	query := `INSERT INTO log (level, message, tm_user) VALUES ($1, $2, $3)`
	log.Debugf("about to exec %s with %s", query, message)
	_, err := db.Exec(query, level, message, user.ID)
	if err != nil {
		log.Errorf("received error: %++v from audit log insertion", err)
		return err
	}
	return nil
}
