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
	"fmt"
	"testing"
	"time"
)

func ExamplePtr() {
	ptr := Ptr("testquest")
	fmt.Println(*ptr)
	// Output: testquest
}

func ExampleStrPtr() {
	ptr := StrPtr("testquest")
	fmt.Println(*ptr)
	// Output: testquest
}

func ExampleIntPtr() {
	ptr := IntPtr(5)
	fmt.Println(*ptr)
	// Output: 5
}

func ExampleUIntPtr() {
	ptr := UIntPtr(5)
	fmt.Println(*ptr)
	// Output: 5
}

func ExampleUInt64Ptr() {
	ptr := UInt64Ptr(5)
	fmt.Println(*ptr)
	// Output: 5
}

func ExampleUint64Ptr() {
	ptr := Uint64Ptr(5)
	fmt.Println(*ptr)
	// Output: 5
}

func ExampleInt64Ptr() {
	ptr := Int64Ptr(5)
	fmt.Println(*ptr)
	// Output: 5
}

func ExampleBoolPtr() {
	ptr := BoolPtr(true)
	fmt.Println(*ptr)
	// Output: true
}

func ExampleFloatPtr() {
	ptr := FloatPtr(5.0)
	fmt.Println(*ptr)
	// Output: 5
}

func ExampleInterfacePtr() {
	ptr := InterfacePtr(1 + 2i)
	fmt.Println(*ptr)
	// Output: (1+2i)
}

func ExampleTimePtr() {
	ptr := TimePtr(time.Time{})
	fmt.Println(*ptr)
	// Output: 0001-01-01 00:00:00 +0000 UTC
}

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
	s := Ptr("9001")
	copiedS := Coalesce(s, "9000")
	if copiedS != "9001" {
		t.Errorf("Coalescing non-nil should've given %s, got: %s", *s, copiedS)
	}
}

func TestCoalesceToDefault(t *testing.T) {
	var i *int
	copiedI := CoalesceToDefault(i)
	var intDefault int
	if copiedI != intDefault {
		t.Errorf("Coalescing nil should've given the default value (%d), got: %d", intDefault, copiedI)
	}
	s := Ptr("9001")
	copiedS := CoalesceToDefault(s)
	if copiedS != "9001" {
		t.Errorf("Coalescing non-nil should've given %s, got: %s", *s, copiedS)
	}
}
