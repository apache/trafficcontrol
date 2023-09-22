package topology

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
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/test"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func randTopology() tc.CRConfigTopology {
	return tc.CRConfigTopology{
		Nodes: test.RandStrArray(),
	}
}

func ExpectedMakeTops() map[string]tc.CRConfigTopology {
	return map[string]tc.CRConfigTopology{
		"top1": randTopology(),
		"top2": randTopology(),
	}
}

func MockMakeTops(mock sqlmock.Sqlmock, expected map[string]tc.CRConfigTopology) {
	rows := sqlmock.NewRows([]string{
		"name",
		"nodes"})

	for topName, top := range expected {
		nodes := "{" + strings.Join(top.Nodes, ",") + "}"
		rows = rows.AddRow(
			topName,
			nodes)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
}

func TestMakeTops(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	expected := ExpectedMakeTops()
	mock.ExpectBegin()
	MockMakeTops(mock, expected)
	mock.ExpectCommit()

	dbCtx, cancelTx := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatal("creating transaction: ", err)
	}

	actual, err := MakeTopologies(tx)
	if err != nil {
		t.Fatal("makeTopologies expected: nil error, actual: ", err)
	}

	if err = db.Close(); err != nil {
		t.Fatal("closing db: ", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf("makeTopologies len expected: %v, actual: %v", len(expected), len(actual))
	}

	for topName, top := range expected {
		actualTop, ok := actual[topName]
		if !ok {
			t.Errorf("makeTopologies expected: %v, actual: missing", topName)
			continue
		}
		expectedBts, _ := json.MarshalIndent(top, " ", " ")
		actualBts, _ := json.MarshalIndent(actualTop, " ", " ")
		if !reflect.DeepEqual(expectedBts, actualBts) {
			t.Errorf("makeDSes ds %+v expected: %+v\n\nactual: %+v\n\n\n", topName, string(expectedBts), string(actualBts))
		}
	}
}
