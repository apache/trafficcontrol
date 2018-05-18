package ip

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

func getTestIPs() []v13.IP {
	ips := []v13.IP{}
	testIP := v13.IP{
		ID:               1,
		Server:           "atlanta-org-1",
		ServerID:         2,
		Type:             "IP_SECONDARY",
		TypeID:           12,
		Interface:        "eth0",
		InterfaceID:      5,
		IP6Address:       "2001::2/64",
		IP6Gateway:       "2001::1",
		IPAddress:        "192.168.0.2",
		IPAddressNetmask: "192.168.0.2/24",
		IPNetmask:        "255.255.255.0",
		IPGateway:        "192.168.0.1",
		LastUpdated:      tc.TimeNoMod{Time: time.Now()},
	}
	ips = append(ips, testIP)

	testIP2 := testIP
	testIP2.ID = 2
	testIP2.IP6Address = "2002::2/64"
	testIP2.IP6Gateway = "2002::1"
	testIP2.IPAddress = "192.168.1.2"
	testIP2.IPAddressNetmask = "192.168.1.2/24"
	testIP2.IPNetmask = "255.255.255.0"
	testIP2.IPGateway = "192.168.1.1"
	ips = append(ips, testIP2)

	return ips
}

func TestReadIPs(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	refType := GetRefType()

	testIPs := getTestIPs()
	cols := test.ColsFromStructByTag("db", v13.IP{})
	var index int
	for i, col := range cols {
		if col == "ipv4" {
			index = i
			break
		}
	}
	cols = append(cols[:index], cols[index+1:]...)
	rows := sqlmock.NewRows(cols)

	for _, ts := range testIPs {
		rows = rows.AddRow(
			ts.ID,
			ts.Server,
			ts.ServerID,
			ts.Type,
			ts.TypeID,
			ts.Interface,
			ts.InterfaceID,
			ts.IP6Address,
			ts.IP6Gateway,
			ts.IPAddress,
			ts.IPNetmask,
			ts.IPGateway,
			ts.LastUpdated,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := map[string]string{"id": "1"}

	ips, errs, _ := refType.Read(db, v, auth.CurrentUser{})
	if len(errs) > 0 {
		t.Errorf("IP.Read expected: no errors, actual: %v", errs)
	}

	if len(ips) != 2 {
		t.Errorf("IP.Read expected: len(ips) == 2, actual: %v", len(ips))
	}
}

func TestFuncs(t *testing.T) {
	if strings.Index(SelectQuery(), "SELECT") != 0 {
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
	i = &TOIP{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("TOIP must be creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("TOIP must be reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("TOIP must be updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("TOIP must be deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("TOIP must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	// empty InterfaceName
	c := TOIP{}
	errs := test.SortErrors(c.Validate(nil))

	expectedErrs := []error{
		errors.New(`'serverId' is required`),
		errors.New(`'typeId' is required`),
	}

	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}

	// valid
	serverId := 1
	typeId := 6
	c = TOIP{ServerID: &serverId, TypeID: &typeId}
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

	cols = []string{"use_in_table"}
	rows = sqlmock.NewRows(cols)
	rows = rows.AddRow("ip")
	mock.ExpectQuery("select").WillReturnRows(rows)

	errs = c.Validate(db)
	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}
}
