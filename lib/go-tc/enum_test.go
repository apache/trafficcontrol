package tc

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
	"encoding/json"
	"testing"
)

func TestDeepCachingType(t *testing.T) {
	var d DeepCachingType
	c := d.String()
	if c != "NEVER" {
		t.Errorf(`Default "%s" expected to be "NEVER"`, c)
	}

	tests := map[string]string{
		"":            "NEVER",
		"NEVER":       "NEVER",
		"ALWAYS":      "ALWAYS",
		"never":       "NEVER",
		"always":      "ALWAYS",
		"Never":       "NEVER",
		"AlwayS":      "ALWAYS",
		" ":           "INVALID",
		"NEVERALWAYS": "INVALID",
		"ALWAYSNEVER": "INVALID",
	}

	for in, exp := range tests {
		dc := DeepCachingTypeFromString(in)
		got, err := json.Marshal(dc)
		if err != nil {
			t.Errorf("%v", err)
		}

		if string(got) != `"`+exp+`"` {
			t.Errorf("for %s,  expected %s,  got %s", in, exp, got)
		}
		var new DeepCachingType
		json.Unmarshal(got, &new)
		if new != dc {
			t.Errorf("Expected %v,  got %v", dc, new)
		}
	}
}
