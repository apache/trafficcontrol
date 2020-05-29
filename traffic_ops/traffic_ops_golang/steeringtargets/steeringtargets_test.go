package steeringtargets

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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"testing"
)

func TestInvalidSteeringTargetType(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	rows := sqlmock.NewRows([]string{
		"name",
		"use_in_table",
	})
	rows = rows.AddRow("HTTP", "server")
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)
	tx := db.MustBegin()

	expected := `type is not a valid steering_target type`
	var dsID, targetID uint64
	dsID = 0
	targetID = 1
	typeID := 1
	st := tc.SteeringTargetNullable{
		DeliveryService:   nil,
		DeliveryServiceID: &dsID,
		Target:            nil,
		TargetID:          &targetID,
		Type:              nil,
		TypeID:            &typeID,
		Value:             nil,
	}
	stObj := &TOSteeringTargetV11{
		APIInfoImpl: api.APIInfoImpl{
			ReqInfo: &api.APIInfo{
				Params:    nil,
				IntParams: nil,
				User:      nil,
				ReqID:     0,
				Version:   nil,
				Tx:        tx,
				Config:    nil,
			},
		},
		SteeringTargetNullable: st,
		DSTenantID:             nil,
		LastUpdated:            nil,
	}

	userErr, _, statusCode := stObj.Create()
	if statusCode != http.StatusBadRequest {
		t.Errorf("expected status code %v,  got %v", http.StatusBadRequest, statusCode)
	}
	if userErr == nil {
		t.Fatal("expected user error to say that type is invalid, got no error instead")
	}
	if userErr.Error() != expected {
		t.Errorf("Expected error details %v, got %v", expected, userErr.Error())
	}
}
