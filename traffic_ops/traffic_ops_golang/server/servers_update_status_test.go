package server

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
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/v6/lib/go-tc"
	"github.com/apache/trafficcontrol/v6/traffic_ops/traffic_ops_golang/config"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetServerUpdateStatus(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	serverStatusRow := sqlmock.NewRows([]string{"id", "host_name", "type", "server_reval_pending", "use_reval_pending", "upd_pending", "status", "parent_upd_pending", "parent_reval_pending"})
	serverStatusRow.AddRow(1, "host_name_1", "EDGE", true, true, true, "ONLINE", true, false)

	mock.ExpectQuery("SELECT").WillReturnRows(serverStatusRow)
	mock.ExpectCommit()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	result, err := getServerUpdateStatus(tx, &config.Config{ConfigTrafficOpsGolang: config.ConfigTrafficOpsGolang{DBQueryTimeoutSeconds: 20}}, "host_name_1")
	if err != nil {
		t.Errorf("getServerUpdateStatus: %v", err)
	}

	expected := []tc.ServerUpdateStatus{{
		HostName:           "host_name_1",
		UpdatePending:      true,
		RevalPending:       true,
		HostId:             1,
		Status:             "ONLINE",
		ParentPending:      true,
		ParentRevalPending: false,
	}}

	reflect.DeepEqual(expected, result)
}
