package crconfig

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
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ExpectedGetConfigParams(domain string) []CRConfigConfigParameter {
	return []CRConfigConfigParameter{
		{"tld.ttls.foo" + test.RandStr(), test.RandStr()},
		{"tld.soa.bar" + test.RandStr(), test.RandStr()},
		{"domain_name", domain},
	}
}

func MockGetConfigParams(mock sqlmock.Sqlmock, expected []CRConfigConfigParameter, cdn string) {
	rows := sqlmock.NewRows([]string{"name", "value"})
	for _, param := range expected {
		n := param.Name
		v := param.Value
		rows = rows.AddRow(n, v)
	}
	mock.ExpectQuery("select").WithArgs(cdn, tc.RouterTypeName).WillReturnRows(rows)
}

func TestGetConfigParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"
	domain := "mycdn.invalid"
	mock.ExpectBegin()
	expected := ExpectedGetConfigParams(domain)
	MockGetConfigParams(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getConfigParams(cdn, tx)
	if err != nil {
		t.Fatalf("getConfigParams err expected: nil, actual: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getConfigParams expected: %+v, actual: %+v", expected, actual)
	}
}

const soaPrefix = "tld.soa."
const ttlPrefix = "tld.ttls."

func ExpectedMakeCRConfigConfig(expectedGetConfigParams []CRConfigConfigParameter, expectedDNSSECEnabled bool) map[string]interface{} {
	m := map[string]interface{}{}
	soa := map[string]string{}
	ttl := map[string]string{}
	for _, param := range expectedGetConfigParams {
		n := param.Name
		v := param.Value
		if strings.HasPrefix(n, soaPrefix) {
			soa[n[len(soaPrefix):]] = v
		} else if strings.HasPrefix(n, ttlPrefix) {
			ttl[n[len(ttlPrefix):]] = v
		} else {
			m[n] = v
		}
	}
	m["soa"] = soa
	m["ttls"] = ttl
	if expectedDNSSECEnabled {
		m["dnssec.enabled"] = "true"
	} else {
		m["dnssec.enabled"] = "false"
	}
	return m
}

func TestMakeCRConfigConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"
	domain := "mycdn.invalid"
	dnssecEnabled := true

	mock.ExpectBegin()
	expectedGetConfigParams := ExpectedGetConfigParams(domain)
	MockGetConfigParams(mock, expectedGetConfigParams, cdn)

	expected := ExpectedMakeCRConfigConfig(expectedGetConfigParams, dnssecEnabled)
	mock.ExpectCommit()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := makeCRConfigConfig(cdn, tx, dnssecEnabled, domain)

	if err != nil {
		t.Fatalf("makeCRConfigConfig err expected: nil, actual: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("makeCRConfigConfig expected: %+v, actual: %+v", expected, actual)
	}
}

func TestCreateMaxmindDefaultOverrideObj(t *testing.T) {
	errs := ""
	separator := ", "
	testCases := []string{
		"US;12.345,-12.345",
		"US",
		"US;12.345",
		"US;abc,-12.345",
		"US;1,abc",
	}

	for _, v := range testCases {
		_, err := createMaxmindDefaultOverrideObj(v)
		if err != nil {
			errs += util.JoinErrsStr([]error{err}) + separator
		}
	}
	errs = errs[:len(errs)-len(separator)]

	expectedErrs := util.JoinErrsStr([]error{
		errors.New(`malformed maxmind.default.override parameter: 'US'`),
		errors.New(`malformed maxmind.default.override parameter coordinates 'US;12.345'`),
		errors.New(`malformed maxmind.default.override parameter coordinates, latitude not a number: 'US;abc,-12.345'`),
		errors.New(`malformed maxmind.default.override parameter coordinates, longitude not an number: 'US;1,abc'`),
	})

	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}
}
