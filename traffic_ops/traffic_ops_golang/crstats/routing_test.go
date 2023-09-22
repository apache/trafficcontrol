package crstats

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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetCDNRouterFQDNs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{
		"host_name",
		"domain_name",
		"port",
		"cdn"})

	rows.AddRow("host1", "test", 2500, "ott")
	mock.ExpectQuery("SELECT").WithArgs("ott").WillReturnRows(rows)
	mock.ExpectCommit()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	result, err := getCDNRouterFQDNs(tx, util.StrPtr("ott"))
	if err != nil {
		t.Fatalf("error Getting CDN router FQDN: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected to receive one item in the response, but got %d", len(result))
	}
	if _, ok := result["ott"]; !ok {
		t.Fatal("expected to get an item with 'ott' key, but got nothing")
	}
}

func TestGetCDNRouterFQDNsWithoutCDN(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{
		"host_name",
		"domain_name",
		"port",
		"cdn"})

	rows.AddRow("host1", "test", 2500, "ott")
	rows.AddRow("host2", "test2", 3500, "newCDN")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatal("creating transaction: ", err)
	}

	result, err := getCDNRouterFQDNs(tx, nil)
	if err != nil {
		t.Fatalf("error Getting CDN router FQDN: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected to receive two items in the response, but got %d", len(result))
	}
	if _, ok := result["ott"]; !ok {
		t.Fatal("expected to get an item with 'ott' key, but got nothing")
	}
	if _, ok := result["newCDN"]; !ok {
		t.Fatal("expected to get an item with 'newCDN' key, but got nothing")
	}
}
