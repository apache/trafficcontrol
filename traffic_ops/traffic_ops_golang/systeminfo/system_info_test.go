package systeminfo

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
	"testing"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"encoding/json"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"
)

func TestGetSystemInfo(t *testing.T) {

	lastUpdated := tc.TimeNoMod{}
	lastUpdated.Scan(time.Now())
	configFile := "global"
	secureFlag := false
	firstID := 1
	firstParam := "paramname1"
	firstVal := "val1"

	secondID := 1
	secondParam := "paramname2"
	secondVal := "val2"

	var sysInfoParameters = []tc.ParameterNullable{

		tc.ParameterNullable{
			ConfigFile:  &configFile,
			ID:          &firstID,
			LastUpdated: &lastUpdated,
			Name:        &firstParam,
			Profiles:    json.RawMessage(`["foo","bar"]`),
			Secure:      &secureFlag,
			Value:       &firstVal,
		},

		tc.ParameterNullable{
			ConfigFile:  &configFile,
			ID:          &secondID,
			LastUpdated: &lastUpdated,
			Name:        &secondParam,
			Profiles:    json.RawMessage(`["foo","bar"]`),
			Secure:      &secureFlag,
			Value:       &secondVal,
		},
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := test.ColsFromStructByTag("db", tc.ParameterNullable{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range sysInfoParameters {
		rows = rows.AddRow(
			ts.ConfigFile,
			ts.ID,
			ts.LastUpdated,
			ts.Name,
			ts.Profiles,
			ts.Secure,
			ts.Value,
		)
	}

	mock.ExpectQuery("SELECT.*WHERE p.config_file='global'").WillReturnRows(rows)
	refType := GetRefType()
	v := map[string]string{"id": "1"}
	parameters, errs, _ := refType.Read(db, v, auth.CurrentUser{})
	if len(errs) > 0 {
		t.Errorf("parameter.Read expected: no errors, actual: %v", errs)
	}

	if len(parameters) != 2 {
		t.Errorf("parameter.Read expected: len(parameters) == 2, actual: %v", len(parameters))
	}
}

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOParameter{}

	if _, ok := i.(api.Reader); !ok {
		t.Errorf("cdn must be reader")
	}
}
