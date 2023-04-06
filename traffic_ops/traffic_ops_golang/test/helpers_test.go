package test

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
)

type thingy struct {
	Cat string `json:"cat" db:"cat"`
	Mat bool   `json:"mat"`
	Hat int    `json:"hat" db:"hat"`
}

// Extract the tag annotations from a struct into a string array
func TestColsFromStructByTag(t *testing.T) {
	res := ColsFromStructByTag("db", thingy{})
	cols := []string{"cat", "hat"}

	if len(res) != len(cols) {
		t.Errorf("Expected %d columns, got %d", len(cols), len(res))
	}
	for i, c := range cols {
		if c != res[i] {
			t.Errorf("Expected %s, got %s", c, res[i])
		}
	}
}

func TestColsFromStructByTagExclude(t *testing.T) {
	res := ColsFromStructByTagExclude("db", thingy{}, []string{"cat"})
	cols := []string{"hat"}

	if len(res) != len(cols) {
		t.Errorf("Expected %d columns, got %d", len(cols), len(res))
	}
	for i, c := range cols {
		if c != res[i] {
			t.Errorf("Expected %s, got %s", c, res[i])
		}
	}
}

func TestInsertAtStr(t *testing.T) {
	cols := []string{"zero", "two", "four"}
	expectedCols := []string{"zero", "one", "two", "three", "four", "five", "six"}
	colMap := map[string][]string{
		"zero": {"one"},
		"two":  {"three"},
		"four": {"five", "six"},
	}

	newCols := InsertAtStr(cols, colMap)
	if newCols == nil {
		t.Fatal("expected new columns to be not nil")
	}

	if len(newCols) != len(expectedCols) {
		t.Fatalf("expected same amount of columns got %d and %d", len(newCols), len(expectedCols))
	}

	for i, _ := range expectedCols {
		if (newCols)[i] != expectedCols[i] {
			t.Fatalf("expected col %v to be the same, got %s and %s", i, newCols[i], expectedCols[i])
		}
	}

	newCols = InsertAtStr(nil, colMap)
	if newCols != nil {
		t.Fatal("expected no new columns when an argument is nil")
	}
	newCols = InsertAtStr(cols, nil)
	if !reflect.DeepEqual(newCols, cols) {
		t.Fatal("expected the same columns when insertMap is nil")
	}
	newCols = InsertAtStr(nil, nil)
	if newCols != nil {
		t.Fatal("expected no new columns when both arguments are nil")
	}
}
