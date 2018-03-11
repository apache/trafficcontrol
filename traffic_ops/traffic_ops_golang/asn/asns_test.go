package asn

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
	"reflect"
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestASNs() []TOASN {
	ASNs := []TOASN{}
	i := 1
	c := "Yukon"
	testCase := TOASN{
		ASN:          &i,
		Cachegroup:   &c,
		CachegroupID: &i,
		ID:           &i,
		LastUpdated:  &tc.TimeNoMod{Time: time.Now()},
	}
	ASNs = append(ASNs, testCase)

	testCase2 := testCase
	*testCase2.ASN = 2
	ASNs = append(ASNs, testCase2)

	return ASNs
}

func TestGetASNs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testCase := getTestASNs()
	cols := test.ColsFromStructByTag("db", TOASN{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testCase {
		rows = rows.AddRow(
			*ts.ASN,
			*ts.CachegroupID,
			*ts.ID,
			*ts.LastUpdated,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := map[string]string{"dsId": "1"}

	asns, errs, _ := refType.Read(db, v, auth.CurrentUser{})

	if len(errs) > 0 {
		t.Errorf("asn.Read expected: no errors, actual: %v", errs)
	}

	if len(asns) != 2 {
		t.Errorf("asn.Read expected: len(asns) == 2, actual: %v", len(asns))
	}

}

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOASN{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("asn must be creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("asn must be reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("asn must be updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("asn must be deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("asn must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	i := -99
	asn := TOASN{ASN: &i, CachegroupID: &i}

	errs := test.SortErrors(asn.Validate(nil))
	expected := []error{
		errors.New(`'asn' must be no less than 0`),
		errors.New(`'cachegroupId' must be no less than 0`),
	}
	if !reflect.DeepEqual(expected, errs) {
		t.Errorf(`expected %v,  got %v`, expected, errs)
	}
}
