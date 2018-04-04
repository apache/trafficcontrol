package crconfig

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
	"reflect"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ExpectedMakeLocations() (map[string]tc.CRConfigLatitudeLongitude, map[string]tc.CRConfigLatitudeLongitude) {
	return map[string]tc.CRConfigLatitudeLongitude{
			"cache0": tc.CRConfigLatitudeLongitude{Lat: *randFloat64(), Lon: *randFloat64()},
			"cache1": tc.CRConfigLatitudeLongitude{Lat: *randFloat64(), Lon: *randFloat64()},
		},
		map[string]tc.CRConfigLatitudeLongitude{
			"router0": tc.CRConfigLatitudeLongitude{Lat: *randFloat64(), Lon: *randFloat64()},
			"router1": tc.CRConfigLatitudeLongitude{Lat: *randFloat64(), Lon: *randFloat64()},
		}
}

func MockMakeLocations(mock sqlmock.Sqlmock, expectedEdgeLocs map[string]tc.CRConfigLatitudeLongitude, expectedRouterLocs map[string]tc.CRConfigLatitudeLongitude, cdn string) {
	rows := sqlmock.NewRows([]string{"name", "type", "latitude", "longitude"})
	for s, l := range expectedEdgeLocs {
		rows = rows.AddRow(s, EdgeTypePrefix, l.Lat, l.Lon)
	}
	for s, l := range expectedRouterLocs {
		rows = rows.AddRow(s, RouterTypeName, l.Lat, l.Lon)
	}
	mock.ExpectQuery("select").WithArgs(cdn).WillReturnRows(rows)
}

func TestMakeLocations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"

	expectedEdgeLocs, expectedRouterLocs := ExpectedMakeLocations()
	MockMakeLocations(mock, expectedEdgeLocs, expectedRouterLocs, cdn)

	actualEdgeLocs, actualRouterLocs, err := makeLocations(cdn, db)
	if err != nil {
		t.Fatalf("makeLocations expected: nil error, actual: %v", err)
	}

	if !reflect.DeepEqual(expectedEdgeLocs, actualEdgeLocs) {
		t.Errorf("makeLocations expected: %+v, actual: %+v", expectedEdgeLocs, actualEdgeLocs)
	}
	if !reflect.DeepEqual(expectedRouterLocs, actualRouterLocs) {
		t.Errorf("makeLocations expected: %+v, actual: %+v", expectedRouterLocs, actualRouterLocs)
	}
}
