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
	"github.com/apache/trafficcontrol/lib/go-tc/tce"
	"reflect"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

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
	matches := []tc.DeliveryServiceMatch{tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`}}
	actual := MakeExampleURLs(util.IntPtr(0), tce.DSTypeHTTP, "routing-name", matches, "domain-name.invalid")
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
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.ds-name.invalid`},
	}
	actual = MakeExampleURLs(util.IntPtr(0), tce.DSTypeDNS, "routing-name", matches, "domain-name.invalid")
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
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.ds-name.invalid`},
	}
	actual = MakeExampleURLs(util.IntPtr(2), tce.DSTypeDNS, "routing-name", matches, "domain-name.invalid")
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
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.ds-name.invalid`},
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 1, Pattern: `fqdn.two.ds-name.invalid`},
	}
	actual = MakeExampleURLs(util.IntPtr(2), tce.DSTypeDNS, "different-routing-name", matches, "different-domain-name.invalid")
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
		tc.DeliveryServiceMatch{Type: tce.DSMatchTypeHostRegex, SetNumber: 0, Pattern: `\.*ds-name\.*`},
	}
	actual = MakeExampleURLs(util.IntPtr(1), tce.DSTypeDNS, "routing-name", matches, "domain-name.invalid")
	if len(expected) != len(actual) {
		t.Fatalf("MakeExampleURLs urls expected %v actual %v", len(expected), len(actual))
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("MakeExampleURLs expected %v actual %v", expected, actual)
	}
}
