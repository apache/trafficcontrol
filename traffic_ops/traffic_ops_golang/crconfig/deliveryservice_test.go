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
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func randDS() tc.CRConfigDeliveryService {
	// truePtr := true
	falseStrPtr := "false"
	// numStr := "42"
	ttlAdmin := "traffic_ops"
	ttlExpire := "604800"
	ttlMinimum := "30"
	ttlRefresh := "28800"
	ttlRetry := "7200"
	ttl := randInt()
	ttlStr := strconv.Itoa(*ttl)
	ttlNS := "3600"
	ttlSOA := "86400"
	geoProviderStr := GeoProviderMaxmindStr
	ecsEnabled := false
	return tc.CRConfigDeliveryService{
		AnonymousBlockingEnabled:  &falseStrPtr,
		CoverageZoneOnly:          false,
		ConsistentHashQueryParams: []string{},
		Dispersion: &tc.CRConfigDispersion{
			Limit:    42,
			Shuffled: true,
		},
		// Domains: []string{"foo"},
		GeoLocationProvider: &geoProviderStr,
		// MatchSets:            randMatchsetArr(),
		MissLocation: &tc.CRConfigLatitudeLongitudeShort{
			Lat: *randFloat64(),
			Lon: *randFloat64(),
		},
		Protocol: &tc.CRConfigDeliveryServiceProtocol{
			// AcceptHTTP: &truePtr,
			AcceptHTTPS:     false,
			RedirectOnHTTPS: false,
		},
		RegionalGeoBlocking:  &falseStrPtr,
		ResponseHeaders:      nil,
		RequestHeaders:       nil,
		RequiredCapabilities: randStrArray(),
		Soa: &tc.SOA{
			Admin:          &ttlAdmin,
			ExpireSeconds:  &ttlExpire,
			MinimumSeconds: &ttlMinimum,
			RefreshSeconds: &ttlRefresh,
			RetrySeconds:   &ttlRetry,
		},
		SSLEnabled: false,
		EcsEnabled: &ecsEnabled,
		Topology:   randStr(),
		TTL:        ttl,
		TTLs: &tc.CRConfigTTL{
			ASeconds:    &ttlStr,
			AAAASeconds: &ttlStr,
			NSSeconds:   &ttlNS,
			SOASeconds:  &ttlSOA,
		},
		// MaxDNSIPsForLocation: randInt(),
		IP6RoutingEnabled: randBool(),
		RoutingName:       randStr(),
		BypassDestination: map[string]*tc.CRConfigBypassDestination{
			"HTTP": &tc.CRConfigBypassDestination{
				// IP: randStr(),
				// IP6: randStr(),
				// CName: randStr(),
				// TTL: randInt(),
				FQDN: randStr(),
				// Port: randStr(),
			},
		},
		DeepCachingType: nil,
		GeoEnabled:      nil,
		// GeoLimitRedirectURL: randStr(),
		StaticDNSEntries: []tc.CRConfigStaticDNSEntry{
			tc.CRConfigStaticDNSEntry{
				Name:  *randStr(),
				TTL:   *randInt(),
				Type:  *randStr(),
				Value: *randStr(),
			},
		},
	}
}

func ExpectedMakeDSes() map[string]tc.CRConfigDeliveryService {
	return map[string]tc.CRConfigDeliveryService{
		"ds1": randDS(),
		"ds2": randDS(),
	}
}

func MockMakeDSes(mock sqlmock.Sqlmock, expected map[string]tc.CRConfigDeliveryService, cdn string) {
	rows := sqlmock.NewRows([]string{
		"anonymous_blocking_enabled",
		"consistent_hash_regex",
		"deep_caching_type",
		"initial_dispersion",
		"dns_bypass_cname",
		"dns_bypass_ip",
		"dns_bypass_ip6",
		"dns_bypass_ttl",
		"query_keys",
		"routing_name",
		"ttl",
		"ecs_enabled",
		"regional_geo_blocking",
		"geo_limit",
		"geo_limit_countries",
		"geeo_limit_redirect_url",
		"geo_provider",
		"http_bypass_fqdn",
		"ipv6_routing_enabled",
		"max_dns_answers",
		"miss_lat",
		"miss_long",
		"profile",
		"protocol",
		"required_capabilities",
		"topology",
		"tr_request_headers",
		"tr_response_headers",
		"tr_response_headers",
		"type",
		"xml_id"})

	for dsName, ds := range expected {
		queryParams := "{" + strings.Join(ds.ConsistentHashQueryParams, ",") + "}"
		rows = rows.AddRow(
			false,
			"",
			nil,
			42,
			"",
			"",
			"",
			0,
			queryParams,
			*ds.RoutingName,
			*ds.TTL,
			*ds.EcsEnabled,
			false,
			0,
			"",
			"",
			0,
			*ds.BypassDestination["HTTP"].FQDN,
			*ds.IP6RoutingEnabled,
			nil,
			ds.MissLocation.Lat,
			ds.MissLocation.Lon,
			"",
			0,
			"{"+strings.Join(ds.RequiredCapabilities, ",")+"}",
			ds.Topology,
			"",
			"",
			"",
			"HTTP",
			dsName)
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestMakeDSes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"
	domain := "mycdn.invalid"

	expected := ExpectedMakeDSes()
	expectedParams := ExpectedGetServerProfileParams(expected)
	expectedDSParams, err := getDSParams(expectedParams)
	if err != nil {
		t.Fatalf("getDSParams error expected: nil, actual: %v", err)
	}
	expectedMatchsets, expectedDomains := ExpectedGetDSRegexesDomains(expectedDSParams)
	expectedStaticDNSEntries := ExpectedGetStaticDNSEntries(expected)

	mock.ExpectBegin()
	MockGetServerProfileParams(mock, expectedParams, cdn)
	MockGetDSRegexesDomains(mock, expectedMatchsets, expectedDomains, cdn)
	MockGetStaticDNSEntries(mock, expectedStaticDNSEntries, cdn)
	MockMakeDSes(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := makeDSes(cdn, domain, tx)
	if err != nil {
		t.Fatalf("makeDSes expected: nil error, actual: %v", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("makeDses len expected: %v, actual: %v", len(expected), len(actual))
	}

	for dsName, ds := range expected {
		actualDS, ok := actual[dsName]
		if !ok {
			t.Errorf("makeDSes expected: %v, actual: missing", dsName)
			continue
		}
		expectedBts, _ := json.MarshalIndent(ds, " ", " ")
		actualBts, _ := json.MarshalIndent(actualDS, " ", " ")
		if !reflect.DeepEqual(expectedBts, actualBts) {
			t.Errorf("makeDSes ds %+v expected: %+v\n\nactual: %+v\n\n\n", dsName, string(expectedBts), string(actualBts))
		}
	}
}

func ExpectedGetServerProfileParams(expectedMakeDSes map[string]tc.CRConfigDeliveryService) map[string]map[string]string {
	expected := map[string]map[string]string{}
	for dsName, _ := range expectedMakeDSes {
		expected[dsName] = map[string]string{
			"param0": "val0",
			"param1": "val1",
		}
	}
	return expected
}

func MockGetServerProfileParams(mock sqlmock.Sqlmock, expected map[string]map[string]string, cdn string) {
	rows := sqlmock.NewRows([]string{"name", "value", "profile"})
	for dsName, params := range expected {
		for param, val := range params {
			rows = rows.AddRow(param, val, dsName)
		}
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetServerProfileParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	expectedMakeDSes := ExpectedMakeDSes()
	expected := ExpectedGetServerProfileParams(expectedMakeDSes)

	mock.ExpectBegin()
	MockGetServerProfileParams(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getServerProfileParams(cdn, tx)
	if err != nil {
		t.Fatalf("getServerProfileParams expected: nil error, actual: %v", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("getServerProfileParams len expected: %v, actual: %v (%+v)", len(expected), len(actual), actual)
	}

	for dsName, expectedParams := range expected {
		actualParams, ok := actual[dsName]
		if !ok {
			t.Errorf("getServerProfileParams expected: %v, actual: missing (actual %+v)", dsName, actual)
			continue
		}
		if !reflect.DeepEqual(expectedParams, actualParams) {
			t.Errorf("getServerProfileParams ds %+v expected: %+v, actual: %+v", dsName, expectedParams, actualParams)
		}
	}
}

func ExpectedGetDSRegexesDomains(expectedDSParams map[string]string) (map[string][]*tc.MatchSet, map[string][]string) {
	matchsets := map[string][]*tc.MatchSet{}
	domains := map[string][]string{}

	setnum := 0
	protocolStr := "HTTP"
	matchType := "HOST_REGEXP"

	domain := "foo"
	if val, ok := expectedDSParams["domain_name"]; ok {
		domain = val
	}

	for dsName, _ := range expectedDSParams {
		pattern := `.*\.` + dsName + `\..*`

		matchsets[dsName][setnum] = &tc.MatchSet{}
		matchset := matchsets[dsName][setnum]
		matchset.Protocol = protocolStr
		matchset.MatchList = append(matchset.MatchList, tc.MatchList{MatchType: matchType, Regex: pattern})

		domains[dsName] = append(domains[dsName], strings.NewReplacer(`\`, ``, `.*`, ``, `.`, ``).Replace(pattern)+"."+domain)
	}
	return matchsets, domains
}

func MockGetDSRegexesDomains(mock sqlmock.Sqlmock, expectedMatchsets map[string][]*tc.MatchSet, expectedDomains map[string][]string, cdn string) {
	rows := sqlmock.NewRows([]string{"pattern", "type", "dstype", "set_number", "xml_id"})
	for dsName, matchsets := range expectedMatchsets {
		for _, matchset := range matchsets {
			for _, matchlist := range matchset.MatchList {
				rows = rows.AddRow(matchlist.Regex, "HOST", "HTTP", 0, dsName)
			}
		}
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetDSRegexesDomains(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"
	domain := "mycdn.invalid"

	expectedMakeDSes := ExpectedMakeDSes()
	expectedServerProfileParams := ExpectedGetServerProfileParams(expectedMakeDSes)
	expectedDSParams, err := getDSParams(expectedServerProfileParams)
	if err != nil {
		t.Fatalf("getDSParams error expected: nil, actual: %v", err)
	}
	expectedMatchsets, expectedDomains := ExpectedGetDSRegexesDomains(expectedDSParams)

	mock.ExpectBegin()
	MockGetDSRegexesDomains(mock, expectedMatchsets, expectedDomains, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actualMatchsets, actualDomains, err := getDSRegexesDomains(cdn, domain, tx)
	if err != nil {
		t.Fatalf("getDSRegexesDomains expected: nil error, actual: %v", err)
	}

	if len(actualMatchsets) != len(expectedMatchsets) {
		t.Fatalf("getDSRegexesDomains len(matchsets) expected: %v, actual: %v", len(expectedMatchsets), len(actualMatchsets))
	}
	if len(actualDomains) != len(expectedDomains) {
		t.Fatalf("getDSRegexesDomains len(matchsets) expected: %v, actual: %v", len(expectedDomains), len(actualDomains))
	}

	if !reflect.DeepEqual(expectedMatchsets, actualMatchsets) {
		t.Errorf("getDSRegexesDomains expected: %+v, actual: %+v", expectedMatchsets, actualMatchsets)
	}
	if !reflect.DeepEqual(expectedDomains, actualDomains) {
		t.Errorf("getDSRegexesDomains expected: %+v, actual: %+v", expectedDomains, actualDomains)
	}
}

func ExpectedGetStaticDNSEntries(expectedMakeDSes map[string]tc.CRConfigDeliveryService) map[tc.DeliveryServiceName][]tc.CRConfigStaticDNSEntry {
	expected := map[tc.DeliveryServiceName][]tc.CRConfigStaticDNSEntry{}
	for dsName, ds := range expectedMakeDSes {
		for _, entry := range ds.StaticDNSEntries {
			expected[tc.DeliveryServiceName(dsName)] = append(expected[tc.DeliveryServiceName(dsName)], tc.CRConfigStaticDNSEntry{Name: entry.Name, TTL: entry.TTL, Value: entry.Value, Type: entry.Type})
		}
	}
	return expected
}

func MockGetStaticDNSEntries(mock sqlmock.Sqlmock, expected map[tc.DeliveryServiceName][]tc.CRConfigStaticDNSEntry, cdn string) {
	rows := sqlmock.NewRows([]string{"ds", "name", "ttl", "value", "type"})
	for dsName, entries := range expected {
		for _, entry := range entries {
			rows = rows.AddRow(dsName, entry.Name, entry.TTL, entry.Value, entry.Type+"_RECORD")
		}
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestGetStaticDNSEntries(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	expectedMakeDSes := ExpectedMakeDSes()
	expected := ExpectedGetStaticDNSEntries(expectedMakeDSes)

	mock.ExpectBegin()
	MockGetStaticDNSEntries(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getStaticDNSEntries(cdn, tx)
	if err != nil {
		t.Fatalf("getStaticDNSEntries expected: nil error, actual: %v", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("getStaticDNSEntries len expected: %v, actual: %v", len(expected), len(actual))
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getDSRegexesDomains expected: %+v, actual: %+v", expected, actual)
	}
}
