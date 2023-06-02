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
	"runtime"
	"time"
)

// Stacktrace is a helper function which returns the current stacktrace.
// It wraps runtime.Stack, which requires a sufficiently long buffer.
func Stacktrace() []byte {
	initialBufSize := 1024
	buf := make([]byte, initialBufSize)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, len(buf)*2)
	}
}

// SliceToSet converts a slice to a map whose keys are the slice members, that is, a set.
// Note duplicates will be lost, as is the nature of a set.
func SliceToSet[T comparable](sl []T) map[T]struct{} {
	st := map[T]struct{}{}
	for _, val := range sl {
		st[val] = struct{}{}
	}
	return st
}

// ConvertTimeFormat converts the input time to the supplied format.
func ConvertTimeFormat(t time.Time, format string) (*time.Time, error) {
	formattedTime, err := time.Parse(format, t.Format(format))
	return &formattedTime, err
}
