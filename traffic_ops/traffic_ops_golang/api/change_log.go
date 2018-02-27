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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/jmoiron/sqlx"
)

type ChangeLog struct {
	ID          int     `json:"id" db:"id"`
	Level       string  `json:"level" db:"level"`
	Message     string  `json:"message" db:"message"`
	TMUser      int     `json:"tmUser" db:"tm_user"`
	TicketNum   string  `json:"ticketNum" db:"ticketnum"`
	LastUpdated tc.Time `json:"lastUpdated" db:"last_updated"`
}

const (
	ApiChange = "APICHANGE"
	Updated   = "Updated"
	Created   = "Created"
	Deleted   = "Deleted"
)

func CreateChangeLog(level string, action string, i Identifier, user auth.CurrentUser, db *sqlx.DB) error {
	query := `INSERT INTO log (level, message, tm_user) VALUES ($1, $2, $3)`
	id, _ := i.GetID()
	message := action + " " + i.GetType() + ": " + i.GetAuditName() + " id: " + strconv.Itoa(id)
	log.Debugf("about to exec ", query, " with ", message)
	_, err := db.Exec(query, level, message, user.ID)
	if err != nil {
		log.Errorf("received error: %++v from audit log insertion", err)
		return err
	}
	return nil
}
