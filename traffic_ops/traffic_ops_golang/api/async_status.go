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
	"errors"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	AsyncSucceeded = "SUCCEEDED"
	AsyncFailed    = "FAILED"
	AsyncPending   = "PENDING"
)

const CurrentAsyncEndpoint = "/api/4.0/async_status/"

type AsyncStatus struct {
	Id        int        `json:"id, omitempty" db:"id"`
	Status    string     `json:"status, omitempty" db:"status"`
	StartTime time.Time  `json:"start_time, omitempty" db:"start_time"`
	EndTime   *time.Time `json:"end_time, omitempty" db:"end_time"`
	Message   *string    `json:"message, omitempty" db:"message"`
}

const SelectAsyncStatusQuery = `SELECT id, status, message, start_time, end_time from async_status WHERE id = $1`
const InsertAsyncStatusQuery = `INSERT INTO async_status (status, message) VALUES ($1, $2) RETURNING id`
const UpdateAsyncStatusEndTimeQuery = `UPDATE async_status SET status = $1, message = $2, end_time = now() WHERE id = $3`
const UpdateAsyncStatusQuery = `UPDATE async_status SET status = $1, message = $2 WHERE id = $3`

func GetAsyncStatus(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := NewInfo(r, []string{"id"}, nil)
	if userErr != nil || sysErr != nil {
		HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	asyncStatusId := inf.Params["id"]

	rows, err := inf.Tx.Tx.Query(SelectAsyncStatusQuery, asyncStatusId)
	if err != nil {
		HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	defer rows.Close()

	asyncStatus := AsyncStatus{}
	rowCount := 0
	for rows.Next() {
		rowCount++
		err := rows.Scan(&asyncStatus.Id, &asyncStatus.Status, &asyncStatus.Message, &asyncStatus.StartTime, &asyncStatus.EndTime)
		if err != nil {
			HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
			return
		}
	}

	if rowCount == 0 {
		HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("async status not found"))
		return
	}

	WriteResp(w, r, asyncStatus)
}

func UpdateAsyncStatus(db *sqlx.DB, newStatus string, newMessage string, asyncStatusId int, finished bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	q := UpdateAsyncStatusQuery
	if finished {
		q = UpdateAsyncStatusEndTimeQuery
	}
	_, err = tx.Exec(q, newStatus, newMessage, asyncStatusId)
	if err != nil {
		return err
	}

	return nil
}
