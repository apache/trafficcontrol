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
	"testing"
)

func TestHostNameToShort(t *testing.T) {
	type InputExpected struct {
		Input    string
		Expected string
	}

	inputs := []InputExpected{
		{"", ""},
		{".", ""},
		{".foo", ""},
		{"foo", "foo"},
		{"foo.", "foo"},
		{"foo.bar", "foo"},
		{"foo.bar.", "foo"},
		{"foo.bar.baz", "foo"},
		{"foo.bar.baz.", "foo"},
		{`!@#$%^&*()_+{}|:"<>?.bar.baz`, `!@#$%^&*()_+{}|:"<>?`},
		{"`.bar.baz", "`"},
	}

	for _, ie := range inputs {
		actual := HostNameToShort(ie.Input)
		if actual != ie.Expected {
			t.Errorf("HostNameToShort(%v) expected '%v' actual '%v'", ie.Input, ie.Expected, actual)
		}
	}
}
