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

func TestCopyIntIfNotNil(t *testing.T) {
	var i *int
	copiedI := copyIntIfNotNil(i)
	if copiedI != nil {
		t.Errorf("Copying a nil int should've given nil, got: %d", *copiedI)
	}
	i = new(int)
	*i = 9000
	copiedI = copyIntIfNotNil(i)
	if copiedI == nil {
		t.Errorf("Copied pointer to %d was nil", *i)
	} else {
		if *copiedI != *i {
			t.Errorf("Incorrectly copied int pointer; expected: %d, got: %d", *i, *copiedI)
		}
		*i = 9001
		if *copiedI == *i {
			t.Error("Expected copy to be 'deep' but modifying the original int changed the copy")
		}
	}
}

func TestCopyFloatIfNotNil(t *testing.T) {
	var f *float64
	copiedF := copyFloatIfNotNil(f)
	if copiedF != nil {
		t.Errorf("Copying a nil int should've given nil, got: %f", *copiedF)
	}
	f = new(float64)
	*f = 9000
	copiedF = copyFloatIfNotNil(f)
	if copiedF == nil {
		t.Errorf("Copied pointer to %f was nil", *f)
	} else {
		if *copiedF != *f {
			t.Errorf("Incorrectly copied float64 pointer; expected: %f, got: %f", *f, *copiedF)
		}
		*f = 9001
		if *copiedF == *f {
			t.Error("Expected copy to be 'deep' but modifying the original float changed the copy")
		}
	}
}

func TestCopyBoolIfNotNil(t *testing.T) {
	var b *bool
	copiedB := copyBoolIfNotNil(b)
	if copiedB != nil {
		t.Errorf("Copying a nil int should've given nil, got: %t", *copiedB)
	}
	b = new(bool)
	*b = true
	copiedB = copyBoolIfNotNil(b)
	if copiedB == nil {
		t.Errorf("Copied pointer to %t was nil", *b)
	} else {
		if *copiedB != *b {
			t.Errorf("Incorrectly copied bool pointer; expected: %t, got: %t", *b, *copiedB)
		}
		*b = false
		if *copiedB == *b {
			t.Error("Expected copy to be 'deep' but modifying the original bool changed the copy")
		}
	}
}

func TestCopyStringIfNotNil(t *testing.T) {
	var s *string
	copiedS := copyStringIfNotNil(s)
	if copiedS != nil {
		t.Errorf("Copying a nil string should've given nil, got: %s", *copiedS)
	}
	s = new(string)
	*s = "test string"
	copiedS = copyStringIfNotNil(s)
	if copiedS == nil {
		t.Errorf("Copied pointer to '%s' was nil", *s)
	} else {
		if *copiedS != *s {
			t.Errorf("Incorrectly copied string pointer; expected: '%s', got: '%s'", *s, *copiedS)
		}
		*s = "a different test string"
		if *copiedS == *s {
			t.Error("Expected copy to be 'deep' but modifying the original string changed the copy")
		}
	}
}
