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

package util

import (
	"time"
)

func Ptr[T any](v T) *T {
	return &v
}

func StrPtr(str string) *string {
	return &str
}

func IntPtr(i int) *int {
	return &i
}

func UIntPtr(u uint) *uint {
	return &u
}

func UInt64Ptr(u uint64) *uint64 {
	return &u
}

func Uint64Ptr(u uint64) *uint64 {
	return &u
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func FloatPtr(f float64) *float64 {
	return &f
}

func InterfacePtr(i interface{}) *interface{} {
	return &i
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
