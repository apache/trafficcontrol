package utils

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

var testSet = map[string][]string{
	"test0-0.1.0-1234.01ab23cd": []string{"test0", "0.1.0", "1234", "01ab23cd"},
	"test1-0.2.0":               []string{"test1", "0.2.0", "", ""},
	"test2":                     []string{"test2", "", "", ""},
}

func TestSplitRPMVersion(t *testing.T) {
	for s, e := range testSet {
		t.Logf("Testing %s %v", s, e)
		n, v, c, h := splitRPMVersion(s)

		if n != e[0] {
			t.Errorf("expected name '%s', got '%s'", n, e[0])
		}
		if v != e[1] {
			t.Errorf("expected version '%s', got '%s'", v, e[1])
		}
		if c != e[2] {
			t.Errorf("expected commits '%s', got '%s'", c, e[2])
		}
		if h != e[3] {
			t.Errorf("expected commitHash '%s', got '%s'", h, e[3])
		}
	}
}
