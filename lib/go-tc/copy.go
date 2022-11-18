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

// copyStringIfNotNil makes a deep copy of s - unless it's nil, in which case it
// just returns nil.
func copyStringIfNotNil(s *string) *string {
	if s == nil {
		return nil
	}
	ret := new(string)
	*ret = *s
	return ret
}

// coalesceString coalesces a possibly nil pointer to a string to a concrete
// string, using the provided default value in case `nil` is encountered.
//
// This can be thought of as roughly the inverse of
// github.com/apache/trafficcontrol/lib/go-util.StrPtr.
func coalesceString(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// copyIntIfNotNil makes a deep copy of i - unless it's nil, in which case it
// just returns nil.
func copyIntIfNotNil(i *int) *int {
	if i == nil {
		return nil
	}
	ret := new(int)
	*ret = *i
	return ret
}

// coalesceInt coalesces a possibly nil pointer to an integer to a concrete
// integer, using the provided default value in case `nil` is encountered.
//
// This can be thought of as roughly the inverse of
// github.com/apache/trafficcontrol/lib/go-util.IntPtr.
func coalesceInt(i *int, def int) int {
	if i == nil {
		return def
	}
	return *i
}

// copyBoolIfNotNil makes a deep copy of b - unless it's nil, in which case it
// just returns nil.
func copyBoolIfNotNil(b *bool) *bool {
	if b == nil {
		return nil
	}
	ret := new(bool)
	*ret = *b
	return ret
}

// coalesceBool coalesces a possibly nil pointer to a boolean to a concrete
// boolean, using the provided default value in case `nil` is encountered.
//
// This can be thought of as roughly the inverse of
// github.com/apache/trafficcontrol/lib/go-util.BoolPtr.
func coalesceBool(b *bool, def bool) bool {
	if b == nil {
		return def
	}
	return *b
}

// copyFloatIfNotNil makes a deep copy of f - unless it's nil, in which case it
// just returns nil.
func copyFloatIfNotNil(f *float64) *float64 {
	if f == nil {
		return nil
	}
	ret := new(float64)
	*ret = *f
	return ret
}

// coalesceFloat coalesces a possibly nil pointer to a float64 to a concrete
// float64, using the provided default value in case `nil` is encountered.
//
// This can be thought of as roughly the inverse of
// github.com/apache/trafficcontrol/lib/go-util.FloatPtr.
func coalesceFloat(f *float64, def float64) float64 {
	if f == nil {
		return def
	}
	return *f
}
