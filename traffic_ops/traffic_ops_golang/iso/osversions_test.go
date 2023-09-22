package iso

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
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetOSVersions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	// Temporary directory prefix. Use the top-level `t` to ensure
	// there's no `/` symbols in the name.
	tmpPrefix := t.Name()

	t.Run("valid-file", func(t *testing.T) {
		expected := tc.OSVersionsResponse{
			"TempleOS": "temple503",
		}

		dir, err := ioutil.TempDir("", tmpPrefix)
		if err != nil {
			log.Fatalf("error creating tempdir: %v", err)
		}
		// Clean up temp dir + file
		defer os.RemoveAll(dir)

		// Create config file within temp dir
		fd, err := os.Create(filepath.Join(dir, "osversions.json"))
		if err != nil {
			t.Fatalf("error creating tempfile: %v", err)
		}
		defer fd.Close()
		if err = json.NewEncoder(fd).Encode(expected); err != nil {
			t.Fatal(err)
		}

		dbCtx, cancel := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
		defer cancel()

		// Setup mock DB to return row for SELECT query on parameter table.
		mock.ExpectBegin()
		cols := []string{"value"}
		rows := sqlmock.NewRows(cols)
		rows = rows.AddRow(dir) // return temp dir
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		mock.ExpectCommit()

		tx, err := db.BeginTxx(dbCtx, nil)
		if err != nil {
			t.Fatalf("BeginTxx() err: %v", err)
		}
		defer tx.Commit()

		got, err := getOSVersions(tx)
		if err != nil {
			t.Fatalf("getOSVersions() err = %v; expected nil", err)
		}
		t.Logf("getOSVersions(): %#v", got)

		if lenGot, lenExp := len(got), len(expected); lenGot != lenExp {
			t.Fatalf("incorrect map length: got %d map entries, expected %d", lenGot, lenExp)
		}
		for k, expectedVal := range expected {
			if gotVal := got[k]; gotVal != expectedVal {
				t.Fatalf("incorrect map entry for key %q: got %q, expected %q", k, gotVal, expectedVal)
			}
		}
	})

	t.Run("invalid-file", func(t *testing.T) {
		dir, err := ioutil.TempDir("", tmpPrefix)
		if err != nil {
			log.Fatalf("error creating tempdir: %v", err)
		}
		// Clean up temp dir + file
		defer os.RemoveAll(dir)

		// No config file is created within temp dir

		dbCtx, cancel := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
		defer cancel()

		// Setup mock DB to return row for SELECT query on parameter table.
		mock.ExpectBegin()
		cols := []string{"value"}
		rows := sqlmock.NewRows(cols)
		rows = rows.AddRow(dir) // return temp dir, which is empty
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		mock.ExpectCommit()

		tx, err := db.BeginTxx(dbCtx, nil)
		if err != nil {
			t.Fatalf("BeginTxx() err: %v", err)
		}
		defer tx.Commit()

		_, err = getOSVersions(tx)
		if err == nil {
			t.Fatalf("getOSVersions() err = %v; expected non-nil", err)
		}
		t.Logf("getOSVersions() err (expected) = %v", err)
	})
}

func TestOsversionsCfgPath(t *testing.T) {
	cases := []struct {
		name           string
		parameterValue string
		expectedPath   string
	}{
		{
			"default",
			"", // No db parameter entry
			filepath.Join(cfgDefaultDir, cfgFilename),
		},
		{
			"override",
			"/this/is/not/the/default/",
			filepath.Join("/this/is/not/the/default", cfgFilename),
		},
		{
			"override-cwd",
			".", // CWD
			cfgFilename,
		},
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() err: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dbCtx, cancel := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
			defer cancel()

			// Setup mock DB to return rows for SELECT query on parameter table.
			// If parameterValue is empty, no rows will be returned.
			mock.ExpectBegin()
			cols := []string{"value"}
			rows := sqlmock.NewRows(cols)
			if tc.parameterValue != "" {
				rows = rows.AddRow(tc.parameterValue)
			}
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			mock.ExpectCommit()

			tx, err := db.BeginTxx(dbCtx, nil)
			if err != nil {
				t.Fatalf("BeginTxx() err: %v", err)
			}
			defer tx.Commit()

			got, err := osversionCfgPath(tx)

			if err != nil {
				t.Fatalf("osversionCfgPath() err: %v, expected: nil", err)
			}

			if got != tc.expectedPath {
				t.Fatalf("osversionCfgPath(): %q, expected: %q", got, tc.expectedPath)
			}
			t.Logf("osversionCfgPath(): %q", got)
		})
	}
}
