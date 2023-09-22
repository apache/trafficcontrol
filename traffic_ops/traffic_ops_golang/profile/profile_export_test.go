package profile

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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	exProRows = []string{
		"profile_description",
		"profile_name",
		"profile_type",
		"cdn",
		"parm_name",
		"parm_config_file",
		"parm_value",
	}
)

func TestGetExportProfileResponse(t *testing.T) {
	var testCases = []struct {
		description  string
		storageError error
		profile      tc.ProfileExportImportNullable
		parameters   []tc.ProfileExportImportParameterNullable
	}{
		{
			description:  "Success: Read export profile successful",
			storageError: nil,
			profile:      generateExportImportProfile("profile", "test profile", "cdn", "type"),
			parameters: []tc.ProfileExportImportParameterNullable{
				generateExportImportParameter("config", "param1", "val1"),
				generateExportImportParameter("config", "param2", "val2"),
			},
		},
		{
			description:  "Success: Read export profile with no parameters successful",
			storageError: nil,
			profile:      generateExportImportProfile("profile", "test profile", "cdn", "type"),
			parameters:   []tc.ProfileExportImportParameterNullable{},
		},
		{
			description:  "Failure: Storage error reading profile",
			storageError: errors.New("Storage error"),
			profile:      tc.ProfileExportImportNullable{},
			parameters:   []tc.ProfileExportImportParameterNullable{},
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
			mock.ExpectBegin()
			if testCase.storageError != nil {
				mock.ExpectQuery("profile").WillReturnError(testCase.storageError)
			} else {
				rows := sqlmock.NewRows(exProRows)
				if len(testCase.parameters) == 0 {
					rows = rows.AddRow(
						testCase.profile.Description,
						testCase.profile.Name,
						testCase.profile.Type,
						testCase.profile.CDNName,
						nil,
						nil,
						nil,
					)
				} else {
					for _, param := range testCase.parameters {
						rows = rows.AddRow(
							testCase.profile.Description,
							testCase.profile.Name,
							testCase.profile.Type,
							testCase.profile.CDNName,
							param.Name,
							param.ConfigFile,
							param.Value,
						)
					}
				}
				mock.ExpectQuery("profile").WillReturnRows(rows)
			}
			mock.ExpectCommit()
			exportProfileResponse, err := getExportProfileResponse(1, db.MustBegin())
			if testCase.storageError != nil {
				if err == nil {
					t.Errorf("Read error expected: received no error")
				}
				if exportProfileResponse != nil {
					t.Errorf("Export Profile response expected to be nil: received non nil")
				}
			} else {
				if exportProfileResponse.Profile != testCase.profile {
					t.Errorf("Returned profile does not match expected")
				}
				for _, param := range testCase.parameters {
					found := false
					for _, rParam := range exportProfileResponse.Parameters {
						if rParam == param {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected to find parameter %v in return but did not", param.Name)
					}
				}
			}
		})
	}
}

func generateExportImportParameter(configFile, param, val string) tc.ProfileExportImportParameterNullable {
	return tc.ProfileExportImportParameterNullable{
		ConfigFile: &configFile,
		Name:       &param,
		Value:      &val,
	}
}

func generateExportImportProfile(name, description, cdnName, profileType string) tc.ProfileExportImportNullable {
	return tc.ProfileExportImportNullable{
		Name:        &name,
		CDNName:     &cdnName,
		Description: &description,
		Type:        &profileType,
	}
}
