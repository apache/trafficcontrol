package intf

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
	"strings"
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestInterfaces() []v13.Interface {
	interfaces := []v13.Interface{}
	testInterface := v13.Interface{
		ID:            1,
		Server:        "atlanta-org-1",
		ServerID:      2,
		InterfaceName: "eth0",
		InterfaceMtu:  9000,
		LastUpdated:   tc.TimeNoMod{Time: time.Now()},
	}
	interfaces = append(interfaces, testInterface)

	testInterface2 := testInterface
	testInterface2.ID = 2
	testInterface2.InterfaceName = "eth1"
	interfaces = append(interfaces, testInterface2)

	return interfaces
}

func TestReadInterfaces(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	refType := GetRefType()

	testInterfaces := getTestInterfaces()
	cols := test.ColsFromStructByTag("db", v13.Interface{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testInterfaces {
		rows = rows.AddRow(
			ts.ID,
			ts.Server,
			ts.ServerID,
			ts.InterfaceName,
			ts.InterfaceMtu,
			ts.LastUpdated,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := map[string]string{"id": "1"}

	interfaces, errs, _ := refType.Read(db, v, auth.CurrentUser{})
	if len(errs) > 0 {
		t.Errorf("Interface.Read expected: no errors, actual: %v", errs)
	}

	if len(interfaces) != 2 {
		t.Errorf("Interface.Read expected: len(interfaces) == 2, actual: %v", len(interfaces))
	}
}

func TestFuncs(t *testing.T) {
	if strings.Index(selectQuery(), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(insertQuery(), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(updateQuery(), "UPDATE") != 0 {
		t.Errorf("expected updateQuery to start with UPDATE")
	}
	if strings.Index(deleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}

}
func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOInterface{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("TOInterface must be creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("TOInterface must be reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("TOInterface must be updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("TOInterface must be deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("TOInterface must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	// empty InterfaceName
	c := TOInterface{}
	errs := test.SortErrors(c.Validate(nil))

	expectedErrs := []error{
		errors.New(`'interfaceName' is required`),
		errors.New(`'serverId' is required`),
	}

	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}

	// valid
	serverId := 1
	name := "eth0"
	mtu := 9000
	c = TOInterface{ServerID: &serverId, InterfaceName: &name, InterfaceMtu: &mtu}
	expectedErrs = []error{}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cols := []string{"id"}
	rows := sqlmock.NewRows(cols)

	rows = rows.AddRow(1)
	mock.ExpectQuery("select").WillReturnRows(rows)

	errs = c.Validate(db)
	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}
}
