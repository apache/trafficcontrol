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
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDeliveryServiceServerCapabilityInterfaces(t *testing.T) {
	var i interface{}
	i = &ServerCapability{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("DeliveryServiceServerCapability must be Creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("DeliveryServiceServerCapability must be Reader")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("DeliveryServiceServerCapability must be Deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("DeliveryServiceServerCapability must be Identifier")
	}
}

func TestCreateDeliveryServiceServerCapability(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	sc := ServerCapability{
		api.APIInfoImpl{
			ReqInfo: &api.APIInfo{Tx: db.MustBegin()},
		},
		tc.DeliveryServiceServerCapability{},
	}

	rows := sqlmock.NewRows([]string{"server_capability", "deliveryservice_id", "last_updated"}).AddRow(
		util.StrPtr("mem"),
		util.IntPtr(1),
		time.Now(),
	)
	mock.ExpectQuery("INSERT*").WillReturnRows(rows)

	userErr, sysErr, errCode := sc.Create()
	if userErr != nil {
		t.Fatalf(userErr.Error())
	}
	if sysErr != nil {
		t.Fatalf(sysErr.Error())
	}
	if got, want := errCode, http.StatusOK; got != want {
		t.Fatalf(fmt.Sprintf("got %d; expected http status code %d", got, want))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestReadDeliveryServiceServerCapability(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	capability := tc.DeliveryServiceServerCapability{
		ServerCapability:  util.StrPtr("mem"),
		DeliveryServiceID: util.IntPtr(1),
		XMLID:             util.StrPtr("ds1"),
	}

	mock.ExpectBegin()
	sc := ServerCapability{
		api.APIInfoImpl{
			ReqInfo: &api.APIInfo{Tx: db.MustBegin()},
		},
		capability,
	}

	rows := sqlmock.NewRows([]string{"server_capability", "deliveryservice_id", "xml_id", "last_updated"}).AddRow(
		capability.ServerCapability,
		capability.DeliveryServiceID,
		capability.XMLID,
		time.Now(),
	)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, userErr, sysErr, errCode := sc.Read()
	if userErr != nil {
		t.Fatalf(userErr.Error())
	}
	if sysErr != nil {
		t.Fatalf(sysErr.Error())
	}
	if got, want := errCode, http.StatusOK; got != want {
		t.Fatalf(fmt.Sprintf("got %d; expected http status code %d", got, want))
	}
	if got, want := len(results), 1; got != want {
		t.Errorf("got %d; expected %d server capabilities assigned to deliveryservices", got, want)
	}

	for _, result := range results {
		cap, ok := result.(tc.DeliveryServiceServerCapability)
		if ok {
			if got, want := *cap.DeliveryServiceID, 1; got != want {
				t.Errorf("got %d; expected %d ", got, want)
			}
			if got, want := *cap.XMLID, "ds1"; got != want {
				t.Errorf("got %s; expected %s ", got, want)
			}
			if got, want := *cap.ServerCapability, "mem"; got != want {
				t.Errorf("got %s; expected %s ", got, want)
			}
		}
	}
}

func TestDeleteDeliveryServiceServerCapability(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	sc := ServerCapability{
		api.APIInfoImpl{
			ReqInfo: &api.APIInfo{Tx: db.MustBegin()},
		},
		tc.DeliveryServiceServerCapability{},
	}

	userErr, sysErr, errCode := sc.Delete()
	if userErr != nil {
		t.Fatalf(userErr.Error())
	}
	if sysErr != nil {
		t.Fatalf(sysErr.Error())
	}
	if got, want := errCode, http.StatusOK; got != want {
		t.Fatalf(fmt.Sprintf("got %d; expected http status code %d", got, want))
	}
}
