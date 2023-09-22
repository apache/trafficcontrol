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
	idRow = []string{
		"id",
	}
)

type mockStorageReturn struct {
	storageErr error
	empty      bool
	query      string
	id         int
}

func TestGetImportProfile(t *testing.T) {
	var testCases = []struct {
		description       string
		mockStorageReturn mockStorageReturn
		profile           tc.ProfileExportImportNullable
		returnedErr       bool
		returnedID        int
	}{
		{
			description: "Success: Import profile successful",
			mockStorageReturn: mockStorageReturn{
				id:    1,
				query: "profile",
			},
			profile:     generateExportImportProfile("profile", "test profile", "cdn", "type"),
			returnedErr: false,
			returnedID:  1,
		},
		{
			description: "Failure: Import profile didn't insert row",
			mockStorageReturn: mockStorageReturn{
				empty: true,
				query: "profile",
			},
			profile:     generateExportImportProfile("profile", "test profile", "cdn", "type"),
			returnedErr: true,
			returnedID:  0,
		},
		{
			description: "Failure: Import profile storage error",
			mockStorageReturn: mockStorageReturn{
				query:      "profile",
				storageErr: errors.New("storage error"),
			},
			profile:     generateExportImportProfile("profile", "test profile", "cdn", "type"),
			returnedErr: true,
			returnedID:  0,
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
			msr := testCase.mockStorageReturn
			if msr.storageErr != nil {
				mock.ExpectQuery(msr.query).WillReturnError(msr.storageErr)
			} else {
				rows := sqlmock.NewRows(idRow)
				if !msr.empty {
					rows.AddRow(msr.id)
				}
				mock.ExpectQuery(msr.query).WillReturnRows(rows)
			}
			profileID, err := importProfile(&testCase.profile, db.MustBegin().Tx)
			mock.ExpectCommit()
			if testCase.returnedID != profileID {
				t.Errorf("Expected profile id %v on return: received %v", testCase.returnedID, profileID)
			}
			if testCase.returnedErr && err == nil {
				t.Errorf("Expected profile import to return error: received nil error")
			}
			if !testCase.returnedErr && err != nil {
				t.Errorf("Expected profile import to not return error: received error %v", err)
			}
		})
	}
}

func TestGetImportProfileParameters(t *testing.T) {
	param1 := generateExportImportParameter("cf", "param1", "v")
	param2 := generateExportImportParameter("cf", "param2", "v")

	var testCases = []struct {
		description                string
		mockStorageReturns         map[string][]mockStorageReturn
		parameters                 []tc.ProfileExportImportParameterNullable
		returnedErr                bool
		returnedNewParameters      int
		returnedExistingParameters int
	}{
		{
			description: "Success: All import parameters new",
			mockStorageReturns: map[string][]mockStorageReturn{
				*param1.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						empty: true,
						query: "SELECT",
					},
					// Insert Returns
					mockStorageReturn{
						id:    1,
						query: "INSERT INTO parameter",
					},
				},
				*param2.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						empty: true,
						query: "SELECT",
					},
					// Insert Returns
					mockStorageReturn{
						id:    2,
						query: "INSERT INTO parameter",
					},
				},
			},
			parameters: []tc.ProfileExportImportParameterNullable{
				param1,
				param2,
			},
			returnedNewParameters: 2,
		},
		{
			description: "Success: All parameters exisiting",
			mockStorageReturns: map[string][]mockStorageReturn{
				*param1.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						query: "SELECT",
						id:    1,
					},
				},
				*param2.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						query: "SELECT",
						id:    2,
					},
				},
			},
			parameters: []tc.ProfileExportImportParameterNullable{
				param1,
				param2,
			},
			returnedExistingParameters: 2,
		},
		{
			description: "Success: Mix of existing/new parameters",
			mockStorageReturns: map[string][]mockStorageReturn{
				*param1.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						empty: true,
						query: "SELECT",
					},
					// Insert Returns
					mockStorageReturn{
						id:    1,
						query: "INSERT INTO parameter",
					},
				},
				*param2.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						query: "SELECT",
						id:    2,
					},
				},
			},
			parameters: []tc.ProfileExportImportParameterNullable{
				param1,
				param2,
			},
			returnedNewParameters:      1,
			returnedExistingParameters: 1,
		},
		{
			description: "Success: Dup of existing",
			mockStorageReturns: map[string][]mockStorageReturn{
				*param1.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						empty: true,
						query: "SELECT",
					},
					// Insert Returns
					mockStorageReturn{
						id:    1,
						query: "INSERT INTO parameter",
					},
				},
				*param2.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						query: "SELECT",
						id:    2,
					},
				},
			},
			parameters: []tc.ProfileExportImportParameterNullable{
				param1,
				param2,
				param2,
			},
			returnedNewParameters:      1,
			returnedExistingParameters: 1,
		},
		{
			description: "Fail: Storage error selecting param",
			mockStorageReturns: map[string][]mockStorageReturn{
				*param1.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						empty:      true,
						query:      "SELECT",
						storageErr: errors.New("storage error"),
					},
				},
			},
			parameters: []tc.ProfileExportImportParameterNullable{
				param1,
			},
			returnedErr: true,
		},
		{
			description: "Fail: Storage error inserting param",
			mockStorageReturns: map[string][]mockStorageReturn{
				*param1.Name: []mockStorageReturn{
					// Select Returns
					mockStorageReturn{
						empty: true,
						query: "SELECT",
					},
					// Insert Returns
					mockStorageReturn{
						query:      "INSERT INTO parameter",
						storageErr: errors.New("storage error"),
					},
				},
			},
			parameters: []tc.ProfileExportImportParameterNullable{
				param1,
			},
			returnedErr: true,
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
			for _, param := range testCase.parameters {
				mockStorageReturns := testCase.mockStorageReturns[*param.Name]
				for _, msr := range mockStorageReturns {
					if msr.storageErr != nil {
						mock.ExpectQuery(msr.query).WillReturnError(msr.storageErr)
					} else {
						rows := sqlmock.NewRows(idRow)
						if !msr.empty {
							rows.AddRow(msr.id)
						}
						mock.ExpectQuery(msr.query).
							WithArgs(param.Name, param.ConfigFile, param.Value).
							WillReturnRows(rows)
					}
				}
			}

			mock.ExpectExec("profile_parameter").WillReturnResult(sqlmock.NewResult(1, int64(len(testCase.parameters))))

			newParams, existingParams, err := importProfileParameters(1, testCase.parameters, db.MustBegin().Tx)

			mock.ExpectCommit()
			if testCase.returnedNewParameters != newParams {
				t.Errorf("Expected %v new parameters on return: received %v", testCase.returnedNewParameters, newParams)
			}
			if testCase.returnedExistingParameters != existingParams {
				t.Errorf("Expected %v existing parameters on return: received %v", testCase.returnedExistingParameters, existingParams)
			}
			if testCase.returnedErr && err == nil {
				t.Errorf("Expected profile import to return error: received nil error")
			}
			if !testCase.returnedErr && err != nil {
				t.Errorf("Expected profile import to not return error: received error %v", err)
			}
		})
	}
}
