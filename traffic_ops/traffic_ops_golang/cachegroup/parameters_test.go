package cachegroup

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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	cgpRows = []string{
		"config_file",
		"id",
		"last_updated",
		"name",
		"value",
		"secure",
	}
)

func TestReadCacheGroupParameters(t *testing.T) {

	var testCases = []struct {
		description       string
		storageError      error
		expectedUserError bool
		params            map[string]string
		cgParams          []tc.ParameterNullable
	}{
		{
			description:       "Success: Read Cache Group Parameters",
			storageError:      nil,
			expectedUserError: false,
			params: map[string]string{
				"id": "1",
			},
			cgParams: []tc.ParameterNullable{
				generateParameter("global", "param1", "val1", false, 1),
				generateParameter("global", "param2", "val2", false, 2),
			},
		},
		{
			description:       "Success: Read Cache Group Parameters with parameter id",
			storageError:      nil,
			expectedUserError: false,
			params: map[string]string{
				"id":          "1",
				"parameterId": "1",
			},
			cgParams: []tc.ParameterNullable{
				generateParameter("global", "param1", "val1", false, 1),
			},
		},
		{
			description:       "Success: Read Cache Group Parameters no data",
			storageError:      nil,
			expectedUserError: false,
			params: map[string]string{
				"id": "1",
			},
			cgParams: []tc.ParameterNullable{},
		},
		{
			description:       "Failure: Storage Error reading Cache Group Parameters",
			storageError:      errors.New("failure getting cache group parameters"),
			expectedUserError: false,
			params: map[string]string{
				"id": "1",
			},
			cgParams: []tc.ParameterNullable{
				generateParameter("global", "param1", "val1", false, 1),
				generateParameter("global", "param2", "val2", false, 2),
			},
		},
		{
			description:       "Failure: User Error invalid params",
			storageError:      nil,
			expectedUserError: true,
			params: map[string]string{
				"id": "not_an_id",
			},
			cgParams: []tc.ParameterNullable{
				generateParameter("global", "param1", "val1", false, 1),
				generateParameter("global", "param2", "val2", false, 2),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			t.Log("Starting test scenario: ", testCase.description)
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer mockDB.Close()
			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()
			rows := sqlmock.NewRows(cgpRows)
			for _, cgParam := range testCase.cgParams {
				rows = rows.AddRow(
					cgParam.ConfigFile,
					cgParam.ID,
					cgParam.LastUpdated,
					cgParam.Name,
					cgParam.Value,
					cgParam.Secure,
				)
			}
			mock.ExpectBegin()
			if testCase.storageError != nil {
				mock.ExpectQuery("SELECT").WillReturnError(testCase.storageError)
			} else {
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			}
			mock.ExpectCommit()

			reqInfo := api.APIInfo{Tx: db.MustBegin(), Params: testCase.params}
			obj := TOCacheGroupParameter{
				api.APIInfoImpl{&reqInfo},
				tc.ParameterNullable{},
			}
			parameters, userErr, sysErr, _ := obj.Read()

			if testCase.storageError != nil {
				if sysErr == nil {
					t.Errorf("Read error expected: received no sysErr")
				}
			} else if testCase.expectedUserError {
				if userErr == nil {
					t.Errorf("User error expected: received no userErr")
				}
			} else {
				if userErr != nil || sysErr != nil {
					t.Errorf("Read expected: no errors, actual: %v %v", userErr, sysErr)
				}
				if len(parameters) != len(testCase.cgParams) {
					t.Errorf("cdn.Read expected: len(parameters) == %v, actual: %v", len(testCase.cgParams), len(parameters))
				}
			}
		})
	}
}

func generateParameter(configFile, param, val string, secureFlag bool, id int) tc.ParameterNullable {
	lastUpdated := tc.TimeNoMod{}
	lastUpdated.Scan(time.Now())
	testParameter := tc.ParameterNullable{
		ConfigFile:  &configFile,
		ID:          &id,
		LastUpdated: &lastUpdated,
		Name:        &param,
		Secure:      &secureFlag,
		Value:       &val,
	}
	return testParameter
}
