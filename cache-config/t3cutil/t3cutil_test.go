package t3cutil

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

func TestSortAndCombineStrs(t *testing.T) {
	type Expected struct {
		InputA   []string
		InputB   []string
		Expected []string
	}
	expecteds := []Expected{
		{
			InputA:   []string{"foo", "bar", "baz"},
			InputB:   []string{"alpha", "bar", "beta"},
			Expected: []string{"alpha", "bar", "baz", "beta", "foo"},
		},
		{
			InputA:   []string{"remap.config", "regex_revalidate.config", "parent.config"},
			InputB:   []string{"regex_revalidate.config"},
			Expected: []string{"parent.config", "regex_revalidate.config", "remap.config"},
		},
		{
			InputA:   []string{"remap.config", "parent.config"},
			InputB:   []string{"regex_revalidate.config"},
			Expected: []string{"parent.config", "regex_revalidate.config", "remap.config"},
		},
		{
			InputA:   []string{},
			InputB:   []string{"remap.config", "regex_revalidate.config", "parent.config"},
			Expected: []string{"parent.config", "regex_revalidate.config", "remap.config"},
		},
		{
			InputA:   []string{"remap.config", "regex_revalidate.config", "parent.config"},
			InputB:   []string{},
			Expected: []string{"parent.config", "regex_revalidate.config", "remap.config"},
		},
		{
			InputA:   []string{},
			InputB:   []string{},
			Expected: []string{},
		},
	}

	for _, ex := range expecteds {
		actual := sortAndCombineStrs(ex.InputA, ex.InputB)
		if !reflect.DeepEqual(actual, ex.Expected) {
			t.Errorf("sortAndCombineStrs(%+v,%+v) expected %+v actual %+v", ex.InputA, ex.InputB, ex.Expected, actual)
		}
	}
}
