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

import "testing"

func TestCopyIfNotNil(t *testing.T) {
	var i *int
	copiedI := CopyIfNotNil(i)
	if copiedI != nil {
		t.Errorf("Copying nil should've given nil, got: %d", *copiedI)
	}

	s := new(string)
	*s = "9000"
	copiedS := CopyIfNotNil(s)
	if copiedS == nil {
		t.Errorf("Copied pointer to %s was nil", *s)
	} else {
		if *copiedS != *s {
			t.Errorf("Incorrectly copied pointer; expected: %s, got: %s", *s, *copiedS)
		}
		*s = "9001"
		if *copiedS == *s {
			t.Error("Expected copy to be 'deep' but modifying the original changed the copy")
		}
	}
}

func TestCoalesce(t *testing.T) {
	var i *int
	copiedI := Coalesce(i, 9000)
	if copiedI != 9000 {
		t.Errorf("Coalescing nil should've given the default value, got: %d", copiedI)
	}
	s := new(string)
	*s = "9001"
	copiedS := Coalesce(s, "9000")
	if copiedS != "9001" {
		t.Errorf("Coalescing non-nil should've given %s, got: %s", *s, copiedS)
	}
}

func TestCoalesceToDefault(t *testing.T) {

}
