package sar

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

func TestSplitLast(t *testing.T) {

	type InputExpected struct {
		InputStr   string
		InputDelim string
		Expected   []string
	}

	inputExpecteds := []InputExpected{
		{"", ":", []string{""}},
		{":", ":", []string{"", ""}},
		{"192.168.2.1:42", ":", []string{"192.168.2.1", "42"}},
		{"2001:DB8::1:42", ":", []string{"2001:DB8::1", "42"}},
		{"[2001:DB8::1]:42", ":", []string{"[2001:DB8::1]", "42"}},
	}

	for _, ie := range inputExpecteds {
		actual := SplitLast(ie.InputStr, ie.InputDelim)
		if !reflect.DeepEqual(ie.Expected, actual) {
			for _, actual := range actual {
				t.Errorf("actual '" + actual + "'")
			}
			t.Errorf("input '%v' delim '%v' expected '%v' actual '%v'", ie.InputStr, ie.InputDelim, ie.Expected, actual)
		}
	}
}
