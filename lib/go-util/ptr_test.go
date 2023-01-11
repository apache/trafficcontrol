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

func ExampleCopyIfNotNil() {
	var i *int
	fmt.Println(CopyIfNotNil(i))

	s := Ptr("9000")
	fmt.Println(*CopyIfNotNil(s))

	// Output: <nil>
	// 9000
}

func ExampleCoalesce() {
	var i *int
	fmt.Println(Coalesce(i, 9000))

	s := Ptr("testquest")
	fmt.Println(Coalesce(s, "9001"))
	// Output: 9000
	// testquest
}

func ExampleCoalesceToDefault() {
	var i *int
	fmt.Println(CoalesceToDefault(i))

	s := Ptr("testquest")
	fmt.Println(CoalesceToDefault(s))
	// Output: 0
	// testquest
}
