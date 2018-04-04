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
	"reflect"
	"strings"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ExpectedGetConfigParams() map[string]string {
	return map[string]string{
		"tld.ttls.foo" + *randStr(): *randStr(),
		"tld.soa.bar" + *randStr():  *randStr(),
	}
}

func MockGetConfigParams(mock sqlmock.Sqlmock, expected map[string]string, cdn string) {
	rows := sqlmock.NewRows([]string{"name", "value"})
	for n, v := range expected {
		rows = rows.AddRow(n, v)
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetConfigParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	expected := ExpectedGetConfigParams()
	MockGetConfigParams(mock, expected, cdn)

	actual, err := getConfigParams(cdn, db)
	if err != nil {
		t.Fatalf("getConfigParams err expected: nil, actual: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getConfigParams expected: %+v, actual: %+v", expected, actual)
	}
}

const soaPrefix = "tld.soa."
const ttlPrefix = "tld.ttls."

func ExpectedMakeCRConfigConfig(expectedGetConfigParams map[string]string, expectedDNSSECEnabled bool) map[string]interface{} {
	m := map[string]interface{}{}
	soa := map[string]string{}
	ttl := map[string]string{}
	for n, v := range expectedGetConfigParams {
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
	dnssecEnabled := true

	expectedGetConfigParams := ExpectedGetConfigParams()
	MockGetConfigParams(mock, expectedGetConfigParams, cdn)

	expected := ExpectedMakeCRConfigConfig(expectedGetConfigParams, dnssecEnabled)

	actual, err := makeCRConfigConfig(cdn, db, dnssecEnabled)

	if err != nil {
		t.Fatalf("makeCRConfigConfig err expected: nil, actual: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("makeCRConfigConfig expected: %+v, actual: %+v", expected, actual)
	}
}
