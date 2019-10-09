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
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func randBool() *bool {
	b := rand.Int()%2 == 0
	return &b
}
func randStr() *string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	num := 100
	s := ""
	for i := 0; i < num; i++ {
		s += string(chars[rand.Intn(len(chars))])
	}
	return &s
}
func randInt() *int {
	i := rand.Int()
	return &i
}
func randInt64() *int64 {
	i := int64(rand.Int63())
	return &i
}
func randFloat64() *float64 {
	f := rand.Float64()
	return &f
}

func randServer() tc.CRConfigTrafficOpsServer {
	status := tc.CRConfigServerStatus(*randStr())
	cachegroup := randStr()
	return tc.CRConfigTrafficOpsServer{
		CacheGroup:      cachegroup,
		Fqdn:            randStr(),
		HashCount:       randInt(),
		HashId:          randStr(),
		HttpsPort:       randInt(),
		InterfaceName:   randStr(),
		Ip:              randStr(),
		Ip6:             randStr(),
		LocationId:      cachegroup,
		Port:            randInt(),
		Profile:         randStr(),
		ServerStatus:    &status,
		ServerType:      randStr(),
		RoutingDisabled: *randInt64(),
	}
}

func ExpectedGetServerParams() map[string]ServerParams {
	return map[string]ServerParams{
		"cache0": ServerParams{
			APIPort:          randStr(),
			SecureAPIPort:    randStr(),
			Weight:           randFloat64(),
			WeightMultiplier: randFloat64(),
		},
		"cache1": ServerParams{
			APIPort:          randStr(),
			Weight:           randFloat64(),
			WeightMultiplier: randFloat64(),
		},
	}
}

func MockGetServerParams(mock sqlmock.Sqlmock, expected map[string]ServerParams, cdn string) {
	rows := sqlmock.NewRows([]string{"host_name", "name", "value"})
	rows = rows.AddRow("cache0", "api.port", *expected["cache0"].APIPort)
	rows = rows.AddRow("cache0", "secure.api.port", *expected["cache0"].SecureAPIPort)
	rows = rows.AddRow("cache0", "weight", *expected["cache0"].Weight)
	rows = rows.AddRow("cache0", "weightMultiplier", *expected["cache0"].WeightMultiplier)
	rows = rows.AddRow("cache1", "api.port", *expected["cache1"].APIPort)
	rows = rows.AddRow("cache1", "weight", *expected["cache1"].Weight)
	rows = rows.AddRow("cache1", "weightMultiplier", *expected["cache1"].WeightMultiplier)
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetServerParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	mock.ExpectBegin()
	expected := ExpectedGetServerParams()
	MockGetServerParams(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getServerParams(cdn, tx)
	if err != nil {
		t.Fatalf("getServerParams expected: nil error, actual: %v", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("getServerParams len expected: %v, actual: %v", len(expected), len(actual))
	}

	for name, params := range expected {
		actualParams, ok := actual[name]
		if !ok {
			t.Errorf("getServerParams expected: %v, actual: missing", name)
			continue
		}
		if !reflect.DeepEqual(params, actualParams) {
			t.Errorf("getServerParams server %+v expected: %+v, actual: %+v", name, params, actualParams)
		}
	}
}

func ExpectedGetAllServers(params map[string]ServerParams) map[string]ServerUnion {
	expected := map[string]ServerUnion{}
	for name, param := range params {
		s := ServerUnion{
			APIPort:                  param.APIPort,
			SecureAPIPort:            param.SecureAPIPort,
			CRConfigTrafficOpsServer: randServer(),
		}
		i := int(*param.Weight * *param.WeightMultiplier)
		s.HashCount = &i
		expected[name] = s
	}
	return expected
}

func MockGetAllServers(mock sqlmock.Sqlmock, expected map[string]ServerUnion, cdn string) {
	rows := sqlmock.NewRows([]string{"host_name", "cachegroup", "fqdn", "hashid", "https_port", "interface_name", "ip_address", "ip6_address", "tcp_port", "profile_name", "routing_disabled", "status", "type"})
	for name, s := range expected {
		rows = rows.AddRow(name, *s.CacheGroup, *s.Fqdn, *s.HashId, *s.HttpsPort, *s.InterfaceName, *s.Ip, *s.Ip6, *s.Port, *s.Profile, s.RoutingDisabled, *s.ServerStatus, *s.ServerType)
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetAllServers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	mock.ExpectBegin()
	getServerParamsExpected := ExpectedGetServerParams()
	MockGetServerParams(mock, getServerParamsExpected, cdn)

	expected := ExpectedGetAllServers(getServerParamsExpected)
	MockGetAllServers(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getAllServers(cdn, tx)

	if err != nil {
		t.Fatalf("getAllServers expected: nil error, actual: %v", err)
	}

	if len(actual) != len(expected) {
		t.Errorf("getAllServers len expected: %v, actual: %v", len(expected), len(actual))
	}

	for name, server := range expected {
		actualServer, ok := actual[name]
		if !ok {
			t.Errorf("getAllServers expected: %v, actual: missing", name)
			continue
		}
		if !reflect.DeepEqual(server, actualServer) {
			t.Errorf("getAllServers server %v expected: %v, actual: %v", name, server, actualServer)
		}
	}
}

func ExpectedGetServerDSNames() map[tc.CacheName][]tc.DeliveryServiceName {
	return map[tc.CacheName][]tc.DeliveryServiceName{
		"cache0": []tc.DeliveryServiceName{"ds0", "ds1"},
		"cache1": []tc.DeliveryServiceName{"ds0", "ds1"},
	}
}

func MockGetServerDSNames(mock sqlmock.Sqlmock, expected map[tc.CacheName][]tc.DeliveryServiceName, cdn string) {
	rows := sqlmock.NewRows([]string{"host_name", "xml_id"})
	for cache, dses := range expected {
		for _, ds := range dses {
			rows = rows.AddRow(cache, ds)
		}
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetServerDSNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	mock.ExpectBegin()
	expected := ExpectedGetServerDSNames()
	MockGetServerDSNames(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getServerDSNames(cdn, tx)

	if err != nil {
		t.Fatalf("getServerDSNames expected: nil error, actual: %v", err)
	}

	if len(actual) != len(expected) {
		t.Errorf("getServerDSNames len expected: %v, actual: %v", len(expected), len(actual))
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getServerDSNames expected: %v, actual: %v", expected, actual)
	}
}

func ExpectedGetServerDSes(expectedGetServerDSNames map[tc.CacheName][]tc.DeliveryServiceName) map[tc.CacheName]map[string][]string {
	e := map[tc.CacheName]map[string][]string{}
	for cache, dses := range expectedGetServerDSNames {
		e[cache] = map[string][]string{}
		for _, ds := range dses {
			e[cache][string(ds)] = []string{string(ds) + "regex0", string(ds) + "regex1"}
		}
	}
	return e
}

func MockGetServerDSes(mock sqlmock.Sqlmock, expected map[tc.CacheName]map[string][]string, cdn string) {
	rows := sqlmock.NewRows([]string{"ds", "ds_type", "routing_name", "pattern"})
	dsmap := map[string][]string{}
	for _, dses := range expected {
		for ds, patterns := range dses {
			dsmap[ds] = patterns
		}
	}

	for ds, patterns := range dsmap {
		for _, pattern := range patterns {
			rows = rows.AddRow(ds, "DNS", "", pattern)
		}
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetServerDSes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"
	domain := "mydomain"

	mock.ExpectBegin()
	expectedGetServerDSNames := ExpectedGetServerDSNames()
	MockGetServerDSNames(mock, expectedGetServerDSNames, cdn)

	expected := ExpectedGetServerDSes(expectedGetServerDSNames)
	MockGetServerDSes(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getServerDSes(cdn, tx, domain)

	if err != nil {
		t.Fatalf("getServerDSes expected: nil error, actual: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getServerDSes expected: %v, actual: %v", expected, actual)
	}
}

func ExpectedGetCDNInfo() (string, bool) {
	return *randStr(), *randBool()
}

func MockGetCDNInfo(mock sqlmock.Sqlmock, expectedDomain string, expectedDNSSECEnabled bool, cdn string) {
	rows := sqlmock.NewRows([]string{"domain_name", "dnssec_enabled"})
	rows = rows.AddRow(expectedDomain, expectedDNSSECEnabled)
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetCDNInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	mock.ExpectBegin()
	expectedDomain, expectedDNSSECEnabled := ExpectedGetCDNInfo()
	MockGetCDNInfo(mock, expectedDomain, expectedDNSSECEnabled, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actualDomain, actualDNSSECEnabled, err := getCDNInfo(cdn, tx)
	if err != nil {
		t.Fatalf("getCDNInfo expected: nil error, actual: %v", err)
	}

	if expectedDomain != actualDomain {
		t.Errorf("getCDNInfo expected: %v, actual: %v", expectedDomain, actualDomain)
	}
	if expectedDNSSECEnabled != actualDNSSECEnabled {
		t.Errorf("getCDNInfo expected: %v, actual: %v", expectedDNSSECEnabled, actualDNSSECEnabled)
	}
}

func ExpectedGetCDNNameFromID() string {
	return *randStr()
}

func MockGetCDNNameFromID(mock sqlmock.Sqlmock, expected string, cdnID int) {
	rows := sqlmock.NewRows([]string{"name"})
	rows = rows.AddRow(expected)
	mock.ExpectQuery("select").WithArgs(cdnID).WillReturnRows(rows)
}

func TestGetCDNNameFromID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdnID := 42

	mock.ExpectBegin()
	expected := ExpectedGetCDNNameFromID()
	MockGetCDNNameFromID(mock, expected, cdnID)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, exists, err := getCDNNameFromID(cdnID, tx)
	if err != nil {
		t.Fatalf("getCDNNameFromID expected: nil error, actual: %v", err)
	}
	if !exists {
		t.Fatalf("getCDNNameFromID exists expected: true, actual: false")
	}

	if expected != actual {
		t.Errorf("getCDNNameFromID expected: %v, actual: %v", expected, actual)
	}
}
