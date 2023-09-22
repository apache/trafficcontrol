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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

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
	serverStatusRow := sqlmock.NewRows([]string{"id", "host_name", "type", "server_reval_pending", "use_reval_pending",
		"server_upd_pending", "status", "parent_upd_pending", "parent_reval_pending",
		"config_update_time", "config_apply_time", "config_update_failed", "revalidate_update_time",
		"revalidate_apply_time", "revalidate_update_failed"})
	serverStatusRow.AddRow(1, "host_name_1", "EDGE", true, true,
		true, "ONLINE", true, false,
		time.Now(), time.Now(), false, time.Now(), time.Now(), false)

	mock.ExpectQuery("SELECT").WillReturnRows(serverStatusRow)
	mock.ExpectCommit()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	result, err, _ := getServerUpdateStatus(tx, "host_name_1")
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

func TestGetServerUpdateStatuses(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectBegin()
	revalPendingRows := sqlmock.NewRows([]string{"value"})
	revalPendingRows.AddRow(true)
	mock.ExpectQuery("SELECT").WillReturnRows(revalPendingRows)

	serverInfoRows := sqlmock.NewRows([]string{"id", "host_name", "type", "cdn_id", "status",
		"cachegroup", "config_update_time", "config_apply_time", "config_update_failed", "revalidate_update_time",
		"revalidate_apply_time", "revalidate_update_failed"})
	tenSecAfter := time.UnixMilli(10000)
	epoch := time.UnixMilli(0)
	serverInfoRows.AddRow(1, "edge1", tc.CacheTypeEdge.String(), 1, tc.CacheStatusReported.String(), 1, tenSecAfter, tenSecAfter, false, tenSecAfter, tenSecAfter, false)
	serverInfoRows.AddRow(2, "mid1", tc.CacheTypeMid.String(), 1, tc.CacheStatusReported.String(), 2, tenSecAfter, epoch, false, tenSecAfter, tenSecAfter, false)
	serverInfoRows.AddRow(3, "edge2", tc.CacheTypeEdge.String(), 2, tc.CacheStatusReported.String(), 1, tenSecAfter, tenSecAfter, false, tenSecAfter, tenSecAfter, false)
	serverInfoRows.AddRow(4, "mid2", tc.CacheTypeMid.String(), 2, tc.CacheStatusReported.String(), 2, tenSecAfter, tenSecAfter, false, tenSecAfter, tenSecAfter, false)
	serverInfoRows.AddRow(5, "mid3", tc.CacheTypeMid.String(), 2, tc.CacheStatusReported.String(), 3, tenSecAfter, tenSecAfter, false, tenSecAfter, epoch, false)
	mock.ExpectQuery("SELECT").WillReturnRows(serverInfoRows)

	cachegroupRows := sqlmock.NewRows([]string{"id", "parent_cachegroup_id", "secondary_parent_cachegroup_id"})
	cachegroupRows.AddRow(1, 2, nil)
	cachegroupRows.AddRow(2, nil, nil)
	cachegroupRows.AddRow(3, nil, nil)
	mock.ExpectQuery("SELECT").WillReturnRows(cachegroupRows)

	topologyCachegroupRows := sqlmock.NewRows([]string{"id", "array_agg"})
	topologyCachegroupRows.AddRow(1, "{3}")
	mock.ExpectQuery("SELECT").WillReturnRows(topologyCachegroupRows)

	mock.ExpectCommit()

	expected := map[string][]tc.ServerUpdateStatusV5{
		"edge1": {
			{
				HostName:               "edge1",
				UpdatePending:          false,
				RevalPending:           false,
				UseRevalPending:        true,
				HostId:                 1,
				Status:                 tc.CacheStatusReported.String(),
				ParentPending:          true,
				ParentRevalPending:     false,
				ConfigUpdateTime:       &tenSecAfter,
				ConfigApplyTime:        &tenSecAfter,
				ConfigUpdateFailed:     util.Ptr(false),
				RevalidateUpdateTime:   &tenSecAfter,
				RevalidateApplyTime:    &tenSecAfter,
				RevalidateUpdateFailed: util.Ptr(false),
			},
		},
		"mid1": {
			{
				HostName:               "mid1",
				UpdatePending:          true,
				RevalPending:           false,
				UseRevalPending:        true,
				HostId:                 2,
				Status:                 tc.CacheStatusReported.String(),
				ParentPending:          false,
				ParentRevalPending:     false,
				ConfigUpdateTime:       &tenSecAfter,
				ConfigApplyTime:        &epoch,
				ConfigUpdateFailed:     util.Ptr(false),
				RevalidateUpdateTime:   &tenSecAfter,
				RevalidateApplyTime:    &tenSecAfter,
				RevalidateUpdateFailed: util.Ptr(false),
			},
		},
		"edge2": {
			{
				HostName:               "edge2",
				UpdatePending:          false,
				RevalPending:           false,
				UseRevalPending:        true,
				HostId:                 3,
				Status:                 tc.CacheStatusReported.String(),
				ParentPending:          false,
				ParentRevalPending:     true,
				ConfigUpdateTime:       &tenSecAfter,
				ConfigApplyTime:        &tenSecAfter,
				ConfigUpdateFailed:     util.Ptr(false),
				RevalidateUpdateTime:   &tenSecAfter,
				RevalidateApplyTime:    &tenSecAfter,
				RevalidateUpdateFailed: util.Ptr(false),
			},
		},
		"mid2": {
			{
				HostName:               "mid2",
				UpdatePending:          false,
				RevalPending:           false,
				UseRevalPending:        true,
				HostId:                 4,
				Status:                 tc.CacheStatusReported.String(),
				ParentPending:          false,
				ParentRevalPending:     false,
				ConfigUpdateTime:       &tenSecAfter,
				ConfigApplyTime:        &tenSecAfter,
				ConfigUpdateFailed:     util.Ptr(false),
				RevalidateUpdateTime:   &tenSecAfter,
				RevalidateApplyTime:    &tenSecAfter,
				RevalidateUpdateFailed: util.Ptr(false),
			},
		},
		"mid3": {
			{
				HostName:               "mid3",
				UpdatePending:          false,
				RevalPending:           true,
				UseRevalPending:        true,
				HostId:                 5,
				Status:                 tc.CacheStatusReported.String(),
				ParentPending:          false,
				ParentRevalPending:     false,
				ConfigUpdateTime:       &tenSecAfter,
				ConfigApplyTime:        &tenSecAfter,
				ConfigUpdateFailed:     util.Ptr(false),
				RevalidateUpdateTime:   &tenSecAfter,
				RevalidateApplyTime:    &epoch,
				RevalidateUpdateFailed: util.Ptr(false),
			},
		},
	}
	actual, err := getServerUpdateStatuses(mockDB, 20*time.Second)
	if err != nil {
		t.Fatalf("unexpected error getting server update statuses: %s", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getting server update statuses - expected: %+v, actual: %+v", expected, actual)
	}
}
