package cacheassignmentgroup

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)


func getTestCacheAssignmentGroups() []tc.CacheAssignmentGroup {
	cag := []tc.CacheAssignmentGroup{}
	testCase := tc.CacheAssignmentGroup{
		Name:        "CacheAssignmentGroup1",
		Description: "Description1",
		CDNID:       1,
	}
	cag = append(cag, testCase)

	testCase2 := tc.CacheAssignmentGroup{
		Name:        "CacheAssignmentGroup2",
		Description: "Description2",
		CDNID:        2,
	}
	cag = append(cag, testCase2)

	return cag
}

func TestGetCacheAssignmentGroup(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testCase := getTestCacheAssignmentGroups()
	cols := test.ColsFromStructByTag("db", tc.CacheAssignmentGroup{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testCase {
		rows = rows.AddRow(
			ts.ID,
			ts.Name,
			ts.CDNID,
			ts.LastUpdated,
			ts.Description)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	reqInfo := api.APIInfo{Tx: db.MustBegin(), Params: map[string]string{"dsId": "1"}}
	obj := TOCacheAssignmentGroup{
		api.APIInfoImpl{&reqInfo},
		tc.CacheAssignmentGroupNullable{},
	}
	cags, userErr, sysErr, _ := obj.Read()
	if userErr != nil || sysErr != nil {
		t.Errorf("Read expected: no errors, actual: %v %v", userErr, sysErr)
	}

	if len(cags) != 2 {
		t.Errorf("cacheassignmentgroup.Read expected: len(cacheassignmentgroup) == 2, actual: %v", len(cags))
	}

}

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOCacheAssignmentGroup{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("CacheAssignmentGroup must be Creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("CacheAssignmentGroup must be Reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("CacheAssignmentGroup must be Updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("CacheAssignmentGroup must be Deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("CacheAssignmentGroup must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	p := TOCacheAssignmentGroup{}
	errs := util.JoinErrsStr(test.SortErrors(test.SplitErrors(p.Validate())))
	expected := util.JoinErrsStr(test.SortErrors([]error{
		errors.New("'name' cannot be blank"),
		errors.New("'description' cannot be blank"),
		errors.New("'cdnId' cannot be blank"),
		errors.New("'servers' cannot be blank"),
	}))

	if !reflect.DeepEqual(expected, errs) {
		t.Errorf("expected %++v,  got %++v", expected, errs)
	}
}
