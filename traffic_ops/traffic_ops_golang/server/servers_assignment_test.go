package server

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
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAssignDsesToServer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	newDses := []int{4, 5, 6}
	pqNewDses := pq.Array(newDses)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE").WithArgs(100).WillReturnResult(sqlmock.NewResult(1, 3))
	mock.ExpectExec("INSERT").WithArgs(pqNewDses, 100).WillReturnResult(sqlmock.NewResult(1, 3))

	//fetch remap config location
	remapConfigLocation := "a/path/to/a/remap.config"
	remapConfigRow := sqlmock.NewRows([]string{"value"})
	remapConfigRow.AddRow(remapConfigLocation + "/")          // verifies we strip off the trailing slash
	mock.ExpectQuery("SELECT").WillReturnRows(remapConfigRow) //remap.config

	//select xmlids and edge_header_rewrite, regex_remap, and cache_url  for each ds
	dsFieldRows := sqlmock.NewRows([]string{"xml_id", "edge_header_rewrite", "regex_remap", "cacheurl"})
	dsFieldRows.AddRow("ds1", nil, "regexRemapPlaceholder", "cacheurlPlaceholder")
	dsFieldRows.AddRow("ds2", "edgeHeaderRewritePlaceholder2", "regexRemapPlaceholder", "cacheurlPlaceholder")
	dsFieldRows.AddRow("ds3", "", nil, "cacheurlPlaceholder")
	mock.ExpectQuery("SELECT").WithArgs(pqNewDses).WillReturnRows(dsFieldRows)

	//prepare the insert and delete parameter slices as they should be constructed in the function
	headerRewritePrefix := "hdr_rw_"
	regexRemapPrefix := "regex_remap_"
	cacheurlPrefix := "cacheurl_"
	configPostfix := ".config"
	insert := []string{regexRemapPrefix + "ds1" + configPostfix, cacheurlPrefix + "ds1" + configPostfix, headerRewritePrefix + "ds2" + configPostfix, regexRemapPrefix + "ds2" + configPostfix, cacheurlPrefix + "ds2" + configPostfix, cacheurlPrefix + "ds3" + configPostfix}
	delete := []string{headerRewritePrefix + "ds1" + configPostfix, headerRewritePrefix + "ds3" + configPostfix, regexRemapPrefix + "ds3" + configPostfix}
	fileNamesPq := pq.Array(insert)
	//insert the parameters
	mock.ExpectExec("INSERT").WithArgs(fileNamesPq, "location", remapConfigLocation).WillReturnResult(sqlmock.NewResult(1, 6))

	//select out the parameterIDs we just inserted
	parameterIDRows := sqlmock.NewRows([]string{"id"})
	parameterIDs := []int64{1, 2, 3, 4, 5, 6}
	for _, i := range parameterIDs {
		parameterIDRows.AddRow(i)
	}
	mock.ExpectQuery("SELECT").WithArgs(fileNamesPq).WillReturnRows(parameterIDRows)

	//insert those ids as profile_parameters
	mock.ExpectExec("INSERT").WithArgs(pqNewDses, pq.Array(parameterIDs)).WillReturnResult(sqlmock.NewResult(6, 6))

	//delete the parameters in the delete list
	mock.ExpectExec("DELETE").WithArgs(pq.Array(delete)).WillReturnResult(sqlmock.NewResult(1, 3))
	mock.ExpectCommit()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	result, err := assignDeliveryServicesToServer(100, newDses, true, tx)
	if err != nil {
		t.Errorf("error assigning deliveryservice: %v", err)
	}
	if !reflect.DeepEqual(result, newDses) {
		t.Errorf("delivery services assigned: Expected %v.   Got  %v", newDses, result)
	}
}

func TestCheckForLastServerInActiveDeliveryServices(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	dsIDs := []int{1, 2, 3}
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "multi_site_origin", "topology"})
	rows.AddRow(1, false, util.Ptr(""))
	mock.ExpectQuery("SELECT").WithArgs(1, tc.DSActiveStateActive, pq.Array(dsIDs), tc.CacheStatusOnline, tc.CacheStatusReported, "EDGE%").WillReturnRows(rows)
	mock.ExpectCommit()

	_, err = checkForLastServerInActiveDeliveryServices(1, "EDGE", dsIDs, db.MustBegin().Tx)
	if err != nil {
		t.Errorf("unable to check server in active DS, got error:%v", err)
	}
}

func TestCheckTenancyAndCDN(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	dsIDs := []int{1}
	serverInfo := tc.ServerInfo{
		Cachegroup:   "testCG",
		CachegroupID: 0,
		CDNID:        1,
		DomainName:   "",
		HostName:     "",
		ID:           1,
		Status:       "",
		Type:         "EDGE",
	}
	user := auth.CurrentUser{
		UserName:     "user",
		ID:           0,
		PrivLevel:    0,
		TenantID:     1,
		Role:         0,
		RoleName:     "admin",
		Capabilities: nil,
		UCDN:         "",
	}
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "cdn_id", "tenant_id", "xml_id", "name"})
	rows.AddRow(1, 1, 1, "test", "ALL")
	mock.ExpectQuery("SELECT deliveryservice.id").WithArgs(pq.Array(dsIDs)).WillReturnRows(rows)

	rows1 := sqlmock.NewRows([]string{"id", "active"})
	rows1.AddRow(1, true)
	mock.ExpectQuery("WITH RECURSIVE").WithArgs(user.TenantID, 1).WillReturnRows(rows1)
	mock.ExpectCommit()

	code, usrErr, sysErr := checkTenancyAndCDN(db.MustBegin().Tx, "ALL", 1, serverInfo, dsIDs, &user)
	if usrErr != nil {
		t.Errorf("unable to check tenancy, either DS doesn't exist or DS-CDN not the same as server-CDN, incorrect user input: %v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to check tenancy, system error: %v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("tenancy and cdn check for a given user failed with status code:%d", code)
	}
}

func TestValidateDSCapabilities(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"server_capability"})
	rows.AddRow([]byte("{eas}"))
	mock.ExpectQuery("SELECT").WithArgs("eas").WillReturnRows(rows)

	dsIDs := []int64{1}
	rows1 := sqlmock.NewRows([]string{"id", "required_capabilities"})
	rows1.AddRow(1, []byte("{eas}"))
	mock.ExpectQuery("SELECT ").WithArgs(pq.Array(dsIDs)).WillReturnRows(rows1)
	mock.ExpectCommit()

	usrErr, sysErr, code := ValidateDSCapabilities([]int{1}, "eas", db.MustBegin().Tx)
	if usrErr != nil {
		t.Errorf("unable to validate DS capability, most likely cache doesn't have the DS capability, error: %v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to validate DS Capability, system error: %v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("DS validation failed with status code:%d", code)
	}
}
