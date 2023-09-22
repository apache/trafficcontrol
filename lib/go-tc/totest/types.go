package totest

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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestTypes(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	defer func() {
		err := db.Close()
		assert.NoError(t, err, "unable to close connection to db, error: %v", err)
	}()
	dbQueryTemplate := "INSERT INTO type (name, description, use_in_table) VALUES ('%s', '%s', '%s');"

	for _, typ := range td.Types {
		if typ.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbQueryTemplate, typ.Name, typ.Description, typ.UseInTable))
			assert.RequireNoError(t, err, "could not create Type using database operations: %v", err)
		} else {
			alerts, _, err := cl.CreateType(typ, toclient.RequestOptions{})
			assert.RequireNoError(t, err, "could not create Type: %v - alerts: %+v", err, alerts.Alerts)
		}
	}
}

func DeleteTestTypes(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	dbDeleteTemplate := "DELETE FROM type WHERE name='%s';"

	types, _, err := cl.GetTypes(toclient.RequestOptions{})
	assert.NoError(t, err, "Cannot get Types: %v - alerts: %+v", err, types.Alerts)

	for _, typ := range types.Response {
		if typ.Name == "CHECK_EXTENSION_BOOL" || typ.Name == "CHECK_EXTENSION_NUM" || typ.Name == "CHECK_EXTENSION_OPEN_SLOT" {
			continue
		}

		if typ.UseInTable != "server" {
			err := execSQL(db, fmt.Sprintf(dbDeleteTemplate, typ.Name))
			assert.RequireNoError(t, err, "cannot delete Type using database operations: %v", err)
		} else {
			delResp, _, err := cl.DeleteType(typ.ID, toclient.RequestOptions{})
			assert.RequireNoError(t, err, "cannot delete Type using the API: %v - alerts: %+v", err, delResp.Alerts)
		}

		// Retrieve the Type by name to see if it was deleted.
		opts := toclient.NewRequestOptions()
		opts.QueryParameters.Set("name", typ.Name)
		types, _, err := cl.GetTypes(opts)
		assert.NoError(t, err, "error fetching Types filtered by presumably deleted name: %v - alerts: %+v", err, types.Alerts)
		assert.Equal(t, 0, len(types.Response), "expected Type '%s' to be deleted", typ.Name)
	}
}

func GetTypeId(t *testing.T, cl *toclient.Session, typeName string) int {
	opts := toclient.NewRequestOptions()
	opts.QueryParameters.Set("name", typeName)
	resp, _, err := cl.GetTypes(opts)

	assert.RequireNoError(t, err, "Get Types Request failed with error: %v", err)
	assert.RequireEqual(t, 1, len(resp.Response), "Expected response object length 1, but got %d", len(resp.Response))
	assert.RequireNotNil(t, &resp.Response[0].ID, "Expected id to not be nil")

	return resp.Response[0].ID
}
