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
	"testing"
)

func TestGetTMURLHost(t *testing.T) {
	inputOutputs := [][]string{
		{"a", "a"},
		{"example.com", "example.com"},
		{"example.com:42", "example.com:42"},
		{"example.com:80", "example.com:80"},
		{"example.com:443", "example.com:443"},
		{"http://example.com:42", "example.com:42"},
		{"https://example.com:42", "example.com:42"},
		{"http://example.com:80", "example.com:80"},
		{"https://example.com:80", "example.com:80"},
		{"https://example.com:443", "example.com:443"},
		{"http://example.com:443", "example.com:443"},
		{"example.com:42/", "example.com:42"},
		{"http://example.com:42/", "example.com:42"},
		{"https://example.com:42/", "example.com:42"},
		{"http://example.com:80/", "example.com:80"},
		{"https://example.com:80/", "example.com:80"},
		{"https://example.com:443/", "example.com:443"},
		{"http://example.com:443/", "example.com:443"},
		{"example.com:42/a/b", "example.com:42"},
		{"http://example.com:42/a/b", "example.com:42"},
		{"https://example.com:42/a/b", "example.com:42"},
		{"http://example.com:80/a/b", "example.com:80"},
		{"https://example.com:80/a/b", "example.com:80"},
		{"https://example.com:443/a/b", "example.com:443"},
		{"http://example.com:443/a/b", "example.com:443"},
		{"example.com:42/a/b?c=42/d", "example.com:42"},
		{"http://example.com:42/a/b?c=42/d", "example.com:42"},
		{"https://example.com:42/a/b?c=42/d", "example.com:42"},
		{"http://example.com:80/a/b?c=42/d", "example.com:80"},
		{"https://example.com:80/a/b?c=42/d", "example.com:80"},
		{"https://example.com:443/a/b?c=42/d", "example.com:443"},
		{"http://example.com:443/a/b?c=42/d", "example.com:443"},
		{"http://foo.example.com", "foo.example.com"},
		{"https://foo.example.com", "foo.example.com"},
		{"https://foo.example.com/bar", "foo.example.com"},
		{"http://foo.example.com/bar", "foo.example.com"},
		{"foo.example.com/bar", "foo.example.com"},
		{"http:::foo.com$%^&*()__+/abc/def?42", "http:::foo.com$%^&*()__+"},
		{"http:::foo.com$%^&*()__+/", "http:::foo.com$%^&*()__+"},
		{"http:::foo.com$%^&*()__+", "http:::foo.com$%^&*()__+"},
		{"asdf1345978fasf", "asdf1345978fasf"},
	}
	for _, inputOutput := range inputOutputs {
		input := inputOutput[0]
		expected := inputOutput[1]
		if actual := getTMURLHost(input); expected != actual {
			t.Errorf("getTMHostURL expected: '%+v', actual: '%+v'", expected, actual)
		}
	}
}
