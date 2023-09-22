package deliveryservice

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
	"database/sql/driver"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetDetails(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	rows := sqlmock.NewRows([]string{"routing_name", "ssl_key_version", "name", "id", "origin_server_fqdn"})
	rows.AddRow("cdn", 1, "foo", 1, "http://123.34.32.21:9090")

	rows2 := sqlmock.NewRows([]string{"ds_name", "type", "pattern", "coalesce"})
	rows2.AddRow("testDS", "HOST_REGEXP", ".*\\.testDS\\..*", 0)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT ds.routing_name, ds.ssl_key_version, cdn.name, cdn.id").WillReturnRows(rows)
	mock.ExpectQuery("SELECT ds.xml_id as ds_name, t.name as type, r.pattern").WillReturnRows(rows2)

	oldDetails, userErr, sysErr, code := getOldDetails(1, db.MustBegin().Tx)
	if userErr != nil || sysErr != nil {
		t.Fatalf("didn't expect an error but got user err %v, sys err %v", userErr, sysErr)
	}
	if code != http.StatusOK {
		t.Fatalf("expected status OK 200, but got %d", code)
	}
	if oldDetails.OldOrgServerFQDN == nil {
		t.Fatalf("old org server fqdn is nil")
	}
	if *oldDetails.OldOrgServerFQDN != "http://123.34.32.21:9090" {
		t.Errorf("expected old org server fqdn to be http://123.34.32.21:9090, but got %v", *oldDetails.OldOrgServerFQDN)
	}
	if oldDetails.OldRoutingName != "cdn" {
		t.Errorf("expected old routing name to be cdn, but got %v", oldDetails.OldRoutingName)
	}
	if oldDetails.OldCDNName != "foo" {
		t.Errorf("expected old cdn name to be foo, but got %v", oldDetails.OldCDNName)
	}
	if oldDetails.OldCDNID != 1 {
		t.Errorf("expected old cdn id to be 1, but got %v", oldDetails.OldCDNID)
	}
	if *oldDetails.OldSSLKeyVersion != 1 {
		t.Errorf("expected old ssl_key_version to be 1, but got %v", oldDetails.OldSSLKeyVersion)
	}
}

func TestGetOldDetailsError(t *testing.T) {
	expected := `querying delivery service 1 host name: no such delivery service exists`
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	rows := sqlmock.NewRows([]string{""})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT ds.routing_name, ds.ssl_key_version, cdn.name, cdn.id").WillReturnRows(rows)
	_, userErr, _, code := getOldDetails(1, db.MustBegin().Tx)
	if userErr == nil {
		t.Fatalf("expected error %v, but got none", expected)
	}
	if userErr.Error() != expected {
		t.Errorf("expected error %v, but got %v", expected, userErr.Error())
	}
	if code != http.StatusNotFound {
		t.Errorf("expected error code : %d, but got : %d", http.StatusNotFound, code)
	}
}

func TestGetDeliveryServicesMatchLists(t *testing.T) {
	// test to make sure that the DS matchlists query orders by set_number
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .+ ORDER BY dsr.set_number")

	GetDeliveryServicesMatchLists([]string{"foo"}, db.MustBegin().Tx)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: %s", err)
	}
}

func TestGetDSTLSVersions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error opening a stub database connection: %v", err)
	}
	defer func() {
		if err := mockDB.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close sqlx DB handle: %v", err)
		}
	}()

	rows := sqlmock.NewRows([]string{"tls_version"})
	expected := []string{"1.0", "1.1", "1.2", "1.3"}
	rows.AddRow(fmt.Sprintf("{%s}", strings.Join(expected, ",")))

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	vers, err := GetDSTLSVersions(0, db.MustBegin().Tx)
	if err != nil {
		t.Errorf("Unexpected error getting DS TLS Versions: %v", err)
	} else if len(vers) != 4 {
		t.Errorf("Expected to get 4 TLS versions, got: %d", len(vers))
	} else if !reflect.DeepEqual(expected, vers) {
		t.Errorf("Incorrect TLS versions returned, expected: %+v - actual: %+v", expected, vers)
	}
}

func TestMakeExampleURLs(t *testing.T) {
	expected := []string{
		`http://routing-name.ds-name.domain-name.invalid`,
	}
	matches := []tc.DeliveryServiceMatch{tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`}}
	actual := MakeExampleURLs(util.IntPtr(0), tc.DSTypeHTTP, "routing-name", matches, "domain-name.invalid")
	if len(expected) != len(actual) {
		t.Fatalf("MakeExampleURLs urls expected %v, actual %v", len(expected), len(actual))
	} else if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("MakeExampleURLs expected %v actual %v", expected, actual)
	}

	expected = []string{
		`http://routing-name.ds-name.domain-name.invalid`,
		`http://fqdn.ds-name.invalid`,
	}
	matches = []tc.DeliveryServiceMatch{
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.ds-name.invalid`},
	}
	actual = MakeExampleURLs(util.IntPtr(0), tc.DSTypeDNS, "routing-name", matches, "domain-name.invalid")
	if len(expected) != len(actual) {
		t.Fatalf("MakeExampleURLs urls expected %v actual %v", len(expected), len(actual))
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("MakeExampleURLs expected %v actual %v", expected, actual)
	}

	expected = []string{
		`http://routing-name.ds-name.domain-name.invalid`,
		`https://routing-name.ds-name.domain-name.invalid`,
		`http://fqdn.ds-name.invalid`,
		`https://fqdn.ds-name.invalid`,
	}
	matches = []tc.DeliveryServiceMatch{
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.ds-name.invalid`},
	}
	actual = MakeExampleURLs(util.IntPtr(2), tc.DSTypeDNS, "routing-name", matches, "domain-name.invalid")
	if len(expected) != len(actual) {
		t.Fatalf("MakeExampleURLs urls expected %v actual %v", len(expected), len(actual))
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("MakeExampleURLs expected %v actual %v", expected, actual)
	}

	expected = []string{
		`http://different-routing-name.ds-name.different-domain-name.invalid`,
		`https://different-routing-name.ds-name.different-domain-name.invalid`,
		`http://fqdn.ds-name.invalid`,
		`https://fqdn.ds-name.invalid`,
		`http://fqdn.two.ds-name.invalid`,
		`https://fqdn.two.ds-name.invalid`,
	}
	matches = []tc.DeliveryServiceMatch{
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.ds-name.invalid`},
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.two.ds-name.invalid`},
	}
	actual = MakeExampleURLs(util.IntPtr(2), tc.DSTypeDNS, "different-routing-name", matches, "different-domain-name.invalid")
	if len(expected) != len(actual) {
		t.Fatalf("MakeExampleURLs urls expected %v actual %v", len(expected), len(actual))
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("MakeExampleURLs expected %v actual %v", expected, actual)
	}

	expected = []string{
		`https://routing-name.ds-name.domain-name.invalid`,
	}
	matches = []tc.DeliveryServiceMatch{
		tc.DeliveryServiceMatch{Type: tc.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
	}
	actual = MakeExampleURLs(util.IntPtr(1), tc.DSTypeDNS, "routing-name", matches, "domain-name.invalid")
	if len(expected) != len(actual) {
		t.Fatalf("MakeExampleURLs urls expected %v actual %v", len(expected), len(actual))
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("MakeExampleURLs expected %v actual %v", expected, actual)
	}
}

func TestReadGetDeliveryServices(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	u := auth.CurrentUser{
		TenantID: 1,
	}

	dsRow := []struct {
		key   string
		value driver.Value
	}{
		{"active", tc.DSActiveStateActive},
		{"anonymous_blocking_enabled", false},
		{"ccr_dns_ttl", nil},
		{"cdn_id", 1},
		{"cdnName", "test"},
		{"check_path", ""},
		{"consistent_hash_regex", ""},
		{"deep_caching_type", "NEVER"},
		{"display_name", "Demo 1"},
		{"dns_bypass_cname", nil},
		{"dns_bypass_ip", nil},
		{"dns_bypass_ip6", nil},
		{"dns_bypass_ttl", nil},
		{"dscp", 1},
		{"ecs_enabled", false},
		{"edge_header_rewrite", nil},
		{"first_header_rewrite", nil},
		{"geolimit_redirect_url", nil},
		{"geo_limit", 0},
		{"geo_limit_countries", nil},
		{"geo_provider", 0},
		{"global_max_mbps", nil},
		{"global_max_tps", nil},
		{"fq_pacing_rate", nil},
		{"http_bypass_fqdn", nil},
		{"id", 1},
		{"info_url", nil},
		{"initial_dispersion", 1},
		{"inner_header_rewrite", nil},
		{"ipv6_routing_enabled", true},
		{"last_header_rewrite", nil},
		{"last_updated", time.Now()},
		{"logs_enabled", false},
		{"long_desc", ""},
		{"long_desc_1", nil},
		{"long_desc_2", nil},
		{"max_dns_answers", nil},
		{"max_origin_connections", nil},
		{"max_request_header_bytes", nil},
		{"mid_header_rewrite", nil},
		{"miss_lat", 0.0},
		{"miss_long", 0.0},
		{"multi_site_origin", false},
		{"org_server_fqdn", "origin.infra.ciab.test"},
		{"origin_shield", nil},
		{"profileID", nil},
		{"profile_name", nil},
		{"profile_description", nil},
		{"protocol", 0},
		{"qstring_ignore", 0},
		{"query_keys", "{}"},
		{"range_request_handling", 0},
		{"regex_remap", nil},
		{"regional", false},
		{"regional_geo_blocking", false},
		{"remap_text", nil},
		{"required_capabilities", nil},
		{"routing_name", "video"},
		{"service_category", nil},
		{"signing_algorithm", nil},
		{"range_slice_block_size", nil},
		{"ssl_key_version", nil},
		{"tenant_id", 1},
		{"tenant.name", "test"},
		{"tls_versions", "{}"},
		{"topology", nil},
		{"tr_request_headers", nil},
		{"tr_response_headers", nil},
		{"name", "test"},
		{"type_id", 1},
		{"xml_id", "demo1"},
		{"cdn_domain", "mycdn.ciab.test"},
	}
	keys, values := []string{}, []driver.Value{}
	for _, cell := range dsRow {
		keys = append(keys, cell.key)
		values = append(values, cell.value)
	}

	mock.ExpectBegin()
	tenantRows := sqlmock.NewRows([]string{"id"})
	tenantRows.AddRow(u.TenantID)
	mock.ExpectQuery("WITH RECURSIVE").WillReturnRows(tenantRows)
	dsRows := sqlmock.NewRows(keys)
	dsRows.AddRow(values...)
	mock.ExpectQuery("^SELECT.*ORDER BY ds.xml_id$").WillReturnRows(dsRows)
	regexRows := sqlmock.NewRows([]string{"ds_name", "type", "pattern", "set_number"})
	regexRows.AddRow("demo1", "hostregexp", "", 0)
	mock.ExpectQuery("SELECT ds\\.xml_id as ds_name, t\\.name as type, r\\.pattern, COALESCE\\(dsr\\.set_number, 0\\) FROM regex").WillReturnRows(regexRows)
	mock.ExpectCommit()

	_, userErr, sysErr, _, _ := readGetDeliveryServices(nil, nil, db.MustBegin(), &u, false, api.Version{Major: 5, Minor: 0})
	if userErr != nil {
		t.Errorf("Unexpected user error reading Delivery Services: %v", userErr)
	}
	if sysErr != nil {
		t.Errorf("Unexpected system error reading Delivery Services: %v", sysErr)
	}
}

func TestRequiredIfTypeMatchesName(t *testing.T) {
	ds := &tc.DeliveryServiceV50{
		OrgServerFQDN:         new(string),
		Protocol:              new(int),
		InitialDispersion:     new(int),
		MissLat:               new(float64),
		MissLong:              new(float64),
		RangeRequestHandling:  new(int),
		QStringIgnore:         new(int),
		MaxRequestHeaderBytes: new(int),
		IPV6RoutingEnabled:    new(bool),
	}
	*ds.InitialDispersion = 1
	fn := requiredIfMatchesTypeName([]string{httpTypeRegexp, dnsTypeRegexp}, "HTTP")
	err := fn(ds.OrgServerFQDN)
	if err == nil {
		t.Error("Failed to raise an error when the orgserver fqdn is empty")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}
