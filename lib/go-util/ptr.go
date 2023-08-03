// Package util contains various miscellaneous utilities that are helpful
// throughout components and aren't necessarily related to ATC data structures.
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

import "time"

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// StrPtr returns a pointer to the given string.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func StrPtr(str string) *string {
	return Ptr(str)
}

// IntPtr returns a pointer to the given integer.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func IntPtr(i int) *int {
	return Ptr(i)
}

// UIntPtr returns a pointer to the given unsigned integer.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func UIntPtr(u uint) *uint {
	return Ptr(u)
}

// UInt64Ptr returns a pointer to the given 64-bit unsigned integer.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func UInt64Ptr(u uint64) *uint64 {
	return Ptr(u)
}

// Uint64Ptr returns a pointer to the given 64-bit unsigned integer.
//
// Deprecated. This is just a common mis-casing of UInt64Ptr. These should not
// both exist, and this one - being the less proper casing - is subject to
// removal without warning, as its very existence is likely accidental.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func Uint64Ptr(u uint64) *uint64 {
	return Ptr(u)
}

// Int64Ptr returns a pointer to the given 64-bit integer.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func Int64Ptr(i int64) *int64 {
	return Ptr(i)
}

// BoolPtr returns a pointer to the given boolean.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func BoolPtr(b bool) *bool {
	return Ptr(b)
}

// FloatPtr returns a pointer to the given 64-bit floating-point number.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func FloatPtr(f float64) *float64 {
	return Ptr(f)
}

// InterfacePtr returns a pointer to the given empty interface.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func InterfacePtr(i any) *any {
	return Ptr(i)
}

// TimePtr returns a pointer to the given time.Time value.
//
// Deprecated. This is exactly equivalent to just using Ptr, so duplicated
// functionality like this function will likely be removed before too long.
func TimePtr(t time.Time) *time.Time {
	return Ptr(t)
}

// Coalesce coalesces the given pointer to a concrete value. This is basically
// the inverse operation of Ptr - it safely dereferences its input. If the
// pointer is nil, def is returned as a default value.
func Coalesce[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// CoalesceToDefault coalesces a pointer to the type to which it points. It
// returns the "zero value" of its input's pointed-to type when the input is
// nil. This is equivalent to:
//
//	var x T
//	result := Coalesce(p, x)
//
// ... but can be done on one line without knowing the type of `p`.
func CoalesceToDefault[T any](p *T) T {
	var ret T
	if p != nil {
		ret = *p
	}
	return ret
}

// CopyIfNotNil makes a deep copy of p - unless it's nil, in which case it just
// returns nil.
func CopyIfNotNil[T any](p *T) *T {
	if p == nil {
		return nil
	}
	q := *p
	return &q
}
