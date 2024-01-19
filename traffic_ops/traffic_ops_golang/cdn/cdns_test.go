package cdn

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
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestCDNs() []tc.CDN {
	cdns := []tc.CDN{}
	testCDN := tc.CDN{
		DNSSECEnabled: false,
		DomainName:    "domainName",
		ID:            1,
		Name:          "cdn1",
		LastUpdated:   tc.TimeNoMod{Time: time.Now()},
		TTLOverride:   50,
	}
	cdns = append(cdns, testCDN)

	testCDN2 := testCDN
	testCDN2.Name = "cdn2"
	testCDN2.DomainName = "domain.net"
	cdns = append(cdns, testCDN2)

	return cdns
}

func TestReadCDNs(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testCDNs := getTestCDNs()
	cols := test.ColsFromStructByTag("db", tc.CDN{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testCDNs {
		rows = rows.AddRow(
			ts.DNSSECEnabled,
			ts.DomainName,
			ts.ID,
			ts.LastUpdated,
			ts.Name,
			ts.TTLOverride,
		)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	reqInfo := api.Info{Tx: db.MustBegin(), Params: map[string]string{"dsId": "1"}, Version: &api.Version{Major: 5, Minor: 0}}
	obj := TOCDN{
		api.APIInfoImpl{ReqInfo: &reqInfo},
		tc.CDNNullable{},
	}
	cdns, userErr, sysErr, _, _ := obj.Read(nil, false)
	if userErr != nil || sysErr != nil {
		t.Errorf("Read expected: no errors, actual: %v %v", userErr, sysErr)
	}

	if len(cdns) != 2 {
		t.Errorf("cdn.Read expected: len(cdns) == 2, actual: %v", len(cdns))
	}
}

func TestFuncs(t *testing.T) {
	apiVersion := &api.Version{Major: 4, Minor: 1}
	if strings.Index(selectQuery(apiVersion), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(insertQuery(apiVersion), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(updateQuery(apiVersion), "UPDATE") != 0 {
		t.Errorf("expected updateQuery to start with UPDATE")
	}
	if strings.Index(deleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}
}

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOCDN{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("cdn must be creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("cdn must be reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("cdn must be updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("cdn must be deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("cdn must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	// invalid name, empty domainname
	n := "not_a_valid_cdn"
	reqInfo := api.Info{Tx: nil, Params: map[string]string{"dsId": "1"}, Version: &api.Version{Major: 5, Minor: 0}}
	c := TOCDN{CDNNullable: tc.CDNNullable{Name: &n}, APIInfoImpl: api.APIInfoImpl{ReqInfo: &reqInfo}}
	err, _ := c.Validate()
	errs := util.JoinErrsStr(test.SortErrors(test.SplitErrors(err)))

	expectedErrs := util.JoinErrsStr([]error{
		errors.New(`'domainName' cannot be blank`),
		errors.New(`'name' invalid characters found - Use alphanumeric . or - .`),
	})

	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}

	//  name,  domainname both valid
	n = "This.is.2.a-Valid---CDNNAME."
	d := `awesome-cdn.example.net`
	c = TOCDN{CDNNullable: tc.CDNNullable{Name: &n, DomainName: &d}, APIInfoImpl: api.APIInfoImpl{ReqInfo: &reqInfo}}
	err, _ = c.Validate()
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
}

func TestTOCDNUpdate(t *testing.T) {
	// Create a new mock database and retrieve a mock connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Create a new sqlx database object using the mock connection
	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	// Define the columns for the database rows
	cols := []string{"name"}

	// Create a new sqlmock.Rows object with the defined columns
	rows := sqlmock.NewRows(cols)
	rows.AddRow("testcdn")

	// Expect a transaction begin
	mock.ExpectBegin()

	// Expect a query to select a name from the database with an argument of 1
	mock.ExpectQuery("SELECT name").WithArgs(1).WillReturnRows(rows)

	// Redefine the columns for the database rows
	cols = []string{"username", "soft", "shared_usernames"}

	// Create a new sqlmock.Rows object with the new columns
	rows = sqlmock.NewRows(cols)

	// Expect a query to select username from the database with an argument of "testcdn"
	// and return an error of "sql.ErrNoRows"
	mock.ExpectQuery("SELECT c.username").WithArgs("testcdn").WillReturnError(sql.ErrNoRows)

	// Redefine the columns for the database rows
	cols = []string{"last_updated"}

	// Create a new sqlmock.Rows object with the new columns
	rows = sqlmock.NewRows(cols)
	rows.AddRow(time.Now())

	// Expect a query to select last_updated from the database with an argument of 1
	mock.ExpectQuery("select last_updated").WithArgs(1).WillReturnRows(rows)

	// Create a new api.Info object with required information
	reqInfo := api.Info{
		Tx:     db.MustBegin(),
		Params: map[string]string{"dsId": "1"},
		User: &auth.CurrentUser{
			UserName: "admin",
		},
		Version: &api.Version{Major: 4, Minor: 1},
	}

	// Define variables for the CDN update
	domainName := "example.com"
	id := 1
	lastUpdated := tc.TimeNoMod{Time: time.Now()}
	dnsSecEnabled := false
	name := "testcdn"
	ttlOverride := 0

	// Create a new TOCDN object with the defined variables
	cdn := &TOCDN{
		api.APIInfoImpl{ReqInfo: &reqInfo},
		tc.CDNNullable{
			DNSSECEnabled: &dnsSecEnabled,
			DomainName:    &domainName,
			Name:          &name,
			ID:            &id,
			TTLOverride:   &ttlOverride,
			LastUpdated:   &lastUpdated,
		},
	}

	// Redefine the columns for the database rows
	cols = []string{"last_updated"}

	// Create a new sqlmock.Rows object with the new columns
	rows = sqlmock.NewRows(cols)
	rows.AddRow(time.Now())

	// Expect a query to update the CDN in the database with the CDN object values
	mock.ExpectQuery("UPDATE cdn SET").WithArgs(*cdn.DNSSECEnabled, *cdn.DomainName, *cdn.Name, *cdn.TTLOverride, *cdn.ID).WillReturnRows(rows)

	// Expect a transaction commit
	mock.ExpectCommit()

	// Call the update method on the CDN object and retrieve the error values
	userErr, sysErr, errCode := cdn.Update(nil)

	// Check if there are any unexpected errors
	if userErr != nil || sysErr != nil {
		t.Errorf("Unexpected error: userErr=%v, sysErr=%v, errCode=%d", userErr, sysErr, errCode)
	}

	// Check if the domain name is updated correctly
	if *cdn.DomainName != "example.com" {
		t.Errorf("Unexpected domain name: %s", *cdn.DomainName)
	}

	// Check if the TTL override value is set correctly
	if cdn.TTLOverride == nil {
		t.Errorf("Unexpected TTL override value: %d", *cdn.TTLOverride)
	}
}
