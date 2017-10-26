// Copyright 2015-present Basho Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build timeseries

package riak

import (
	"fmt"
	"testing"
	"time"
)

var runTimeseriesTests = true

// NB: the following is 1443806900 seconds, 103ms after the epoch
var tsMillis = 103 * time.Millisecond
var tsTimestamp = time.Unix(1443806900, tsMillis.Nanoseconds())

const tsQuery = `select * from %v where region = 'South Atlantic' and state = 'South Carolina' and (time > %v and time < %v)`
const tsTable = `WeatherByRegion`
const tsTableDefinition = `
	CREATE TABLE %s (
		region varchar not null,
		state varchar not null,
		time timestamp not null,
		weather varchar not null,
		temperature double,
		uv_index sint64,
		observed boolean not null,
		binary blob not null,
		PRIMARY KEY((region, state, quantum(time, 15, 'm')), region, state, time)
	)`

func init() {
	cmd, err := NewFetchBucketTypePropsCommandBuilder().
		WithBucketType(tsTable).
		Build()
	if err != nil {
		runTimeseriesTests = false
	}

	cluster := integrationTestsBuildCluster()
	defer func() {
		cluster.Stop()
	}()

	if err = cluster.Execute(cmd); err != nil {
		runTimeseriesTests = false
	}
}

// TsFetchRow
func TestTsFetchRowNotFound(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsFetchRowCommandBuilder()
	key := make([]TsCell, 3)

	key[0] = NewStringTsCell("South Atlantic")
	key[1] = NewStringTsCell("South Carolina")
	key[2] = NewTimestampTsCell(time.Now())

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithTable(tsTable).WithKey(key).WithTimeout(time.Second * 30).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*TsFetchRowCommand); ok {
		rsp := scmd.Response
		if rsp == nil {
			t.Errorf("expected non-nil Response")
		}

		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := true, scmd.Response.IsNotFound; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// TsQuery
func TestTsDescribeTable(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsQueryCommandBuilder()

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithQuery("DESCRIBE " + tsTable).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}

	if scmd, ok := cmd.(*TsQueryCommand); ok {
		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if got, want := len(scmd.Response.Columns), 5; !(got >= want) {
			t.Errorf("expected %v to be greater than or equal to %v", got, want)
		}
		if got, want := len(scmd.Response.Rows), 5; !(got >= want) {
			t.Errorf("expected %v to be greater than or equal to %v", got, want)
		}
	} else {
		t.FailNow()
	}
}

// TsQuery
func TestTsCreateTable(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsQueryCommandBuilder()

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	query := fmt.Sprintf(tsTableDefinition, fmt.Sprintf("%v%v", tsTable, time.Now().Unix()))
	cmd, err = sbuilder.WithQuery(query).Build()
	if err != nil {
		t.Log(query)
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Log(query)
		t.Fatal(err.Error())
	}

	if scmd, ok := cmd.(*TsQueryCommand); ok {
		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 0, len(scmd.Response.Columns); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 0, len(scmd.Response.Rows); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// TsStoreRows
func TestTsStoreRow(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsStoreRowsCommandBuilder()
	row := make([]TsCell, 7)

	row[0] = NewStringTsCell("South Atlantic")
	row[1] = NewStringTsCell("South Carolina")
	row[2] = NewTimestampTsCell(tsTimestamp)
	row[3] = NewStringTsCell("hot")
	row[4] = NewDoubleTsCell(23.5)
	row[5] = NewSint64TsCell(10)
	row[6] = NewBooleanTsCell(true)

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithTable(tsTable).WithRows([][]TsCell{row}).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*TsStoreRowsCommand); ok {
		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := true, scmd.Response; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// TsStoreRows
func TestTsStoreRows(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsStoreRowsCommandBuilder()
	row := make([]TsCell, 7)

	row[0] = NewStringTsCell("South Atlantic")
	row[1] = NewStringTsCell("South Carolina")
	row[2] = NewTimestampTsCell(tsTimestamp.Add(-1 * time.Hour))
	row[3] = NewStringTsCell("windy")
	row[4] = NewDoubleTsCell(19.8)
	row[5] = NewSint64TsCell(10)
	row[6] = NewBooleanTsCell(true)

	row2 := row
	row[2] = NewTimestampTsCell(tsTimestamp.Add(-2 * time.Hour))
	row[3] = NewStringTsCell("cloudy")
	row[4] = NewDoubleTsCell(19.1)
	row[5] = NewSint64TsCell(15)
	row[6] = NewBooleanTsCell(false)

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithTable(tsTable).WithRows([][]TsCell{row, row2}).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*TsStoreRowsCommand); ok {
		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := true, scmd.Response; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// TsFetchRow
func TestTsFetchRow(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsFetchRowCommandBuilder()
	key := make([]TsCell, 3)

	key[0] = NewStringTsCell("South Atlantic")
	key[1] = NewStringTsCell("South Carolina")
	key[2] = NewTimestampTsCell(tsTimestamp)

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithTable(tsTable).WithKey(key).WithTimeout(time.Second * 30).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*TsFetchRowCommand); ok {
		rsp := scmd.Response
		if rsp == nil {
			t.Errorf("expected non-nil Response")
		}

		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := false, rsp.IsNotFound; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := 7, len(rsp.Columns); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := 7, len(rsp.Row); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		} else {
			if expected, actual := "TIMESTAMP", rsp.Row[2].GetDataType(); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			} else {
				if expected, actual := tsTimestamp, rsp.Row[2].GetTimeValue(); expected != actual {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			}

			if expected, actual := "VARCHAR", rsp.Row[3].GetDataType(); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}

			if expected, actual := "DOUBLE", rsp.Row[4].GetDataType(); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}

			if expected, actual := "SINT64", rsp.Row[5].GetDataType(); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}

			if expected, actual := "BOOLEAN", rsp.Row[6].GetDataType(); expected != actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		}
	} else {
		t.FailNow()
	}
}

// TsQuery
func TestTsQuery(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsQueryCommandBuilder()
	upperBoundMs := ToUnixMillis(tsTimestamp.Add(time.Second))
	lowerBoundMs := ToUnixMillis(tsTimestamp.Add(-3601 * time.Second))
	query := fmt.Sprintf(tsQuery, tsTable, lowerBoundMs, upperBoundMs)

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithQuery(query).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*TsQueryCommand); ok {
		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 7, len(scmd.Response.Columns); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 1, len(scmd.Response.Rows); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// TsListKeys
func TestTsListKeys(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsListKeysCommandBuilder()
	sbuilder.WithAllowListing()

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithTable(tsTable).WithTimeout(time.Second * 30).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}

	if scmd, ok := cmd.(*TsListKeysCommand); ok {
		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := true, len(scmd.Response.Keys) > 0; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := 3, len(scmd.Response.Keys[0]); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := "VARCHAR", scmd.Response.Keys[0][0].GetDataType(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := "VARCHAR", scmd.Response.Keys[0][1].GetDataType(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := "TIMESTAMP", scmd.Response.Keys[0][2].GetDataType(); expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}

// TsDeleteRow
func TestTsDeleteRow(t *testing.T) {
	if !runTimeseriesTests {
		t.SkipNow()
	}

	var err error
	var cmd Command
	sbuilder := NewTsDeleteRowCommandBuilder()
	key := make([]TsCell, 3)

	key[0] = NewStringTsCell("South Atlantic")
	key[1] = NewStringTsCell("South Carolina")
	key[2] = NewTimestampTsCell(tsTimestamp)

	cluster := integrationTestsBuildCluster()
	defer func() {
		if err := cluster.Stop(); err != nil {
			t.Error(err.Error())
		}
	}()

	cmd, err = sbuilder.WithTable(tsTable).WithKey(key).WithTimeout(time.Second * 30).Build()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		t.Fatal(err.Error())
	}
	if scmd, ok := cmd.(*TsDeleteRowCommand); ok {
		if expected, actual := true, scmd.success; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		if expected, actual := true, scmd.Response; expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	} else {
		t.FailNow()
	}
}
