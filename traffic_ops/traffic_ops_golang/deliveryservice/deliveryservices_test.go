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
	"net/http"
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"

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
	if oldDetails.OldOrgServerFqdn == nil {
		t.Fatalf("old org server fqdn is nil")
	}
	if *oldDetails.OldOrgServerFqdn != "http://123.34.32.21:9090" {
		t.Errorf("expected old org server fqdn to be http://123.34.32.21:9090, but got %v", *oldDetails.OldOrgServerFqdn)
	}
	if oldDetails.OldRoutingName != "cdn" {
		t.Errorf("expected old routing name to be cdn, but got %v", oldDetails.OldRoutingName)
	}
	if oldDetails.OldCdnName != "foo" {
		t.Errorf("expected old cdn name to be foo, but got %v", oldDetails.OldCdnName)
	}
	if oldDetails.OldCdnId != 1 {
		t.Errorf("expected old cdn id to be 1, but got %v", oldDetails.OldCdnId)
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
