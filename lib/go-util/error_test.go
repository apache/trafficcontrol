package util

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
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func ExampleWrapError() {
	err := WrapError("querying for cdns", sql.ErrNoRows)
	fmt.Println(err)
	fmt.Println(err == sql.ErrNoRows, errors.Is(err, sql.ErrNoRows))
	// Output: querying for cdns: sql: no rows in result set
	// false true
}

func BenchmarkErrorStringConcatenation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errors.New("querying for cdns: " + sql.ErrNoRows.Error())
	}
}

func BenchmarkErrorWrapping(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WrapError("querying for cdns", sql.ErrNoRows)
	}
}

func BenchmarkErrorf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Errorf("querying for cdns: %w", sql.ErrNoRows)
	}
}

func BenchmarkNestedWraps(b *testing.B) {
	var err error = sql.ErrNoRows
	for i := 0; i < b.N; i++ {
		err = WrapError("querying for cdns", err)
	}
}

func BenchmarkNestedStringConcatenations(b *testing.B) {
	var err error = sql.ErrNoRows
	for i := 0; i < b.N; i++ {
		err = errors.New("querying for cdns: " + err.Error())
	}
}

func BenchmarkNestedErrorf(b *testing.B) {
	var err error = sql.ErrNoRows
	for i := 0; i < b.N; i++ {
		err = fmt.Errorf("querying for cdns: %w", err)
	}
}
