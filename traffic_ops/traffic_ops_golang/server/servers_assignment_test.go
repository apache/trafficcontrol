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
	"reflect"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAssignDsesToServer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
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
