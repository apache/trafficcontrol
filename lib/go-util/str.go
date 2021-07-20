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
	"strings"
	"unicode"
)

// RemoveStrDuplicates removes duplicates from strings, considering a map of already-seen duplicates.
// Returns the strings which are unique, and not already present in seen; and a map of the unique strings in inputStrings and seenStrings.
//
// This can be used, for example, to remove duplicates from multiple lists of strings, in order, using a shared map of seen strings.
func RemoveStrDuplicates(inputStrings []string, seenStrings map[string]struct{}) ([]string, map[string]struct{}) {
	if seenStrings == nil {
		seenStrings = make(map[string]struct{})
	}
	uniqueStrings := []string{}
	for _, str := range inputStrings {
		if _, ok := seenStrings[str]; !ok {
			uniqueStrings = append(uniqueStrings, str)
			seenStrings[str] = struct{}{}
		}
	}
	return uniqueStrings, seenStrings
}

// StrInArray returns whether s is one of the strings in strs.
//
// Deprecated: This function is totally identical to ContainsStr, but this one
// is not used in any known ATC code, while ContainsStr is. New code should use
// ContainsStr so that this duplicate can be removed.
func StrInArray(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

// RemoveStrFromArray removes a specific string from a string slice.
func RemoveStrFromArray(strs []string, s string) []string {
	newStrArray := []string{}
	for _, str := range strs {
		if str != s {
			newStrArray = append(newStrArray, str)
		}
	}
	return newStrArray
}

// ContainsStr returns whether x is one of the elements of a.
func ContainsStr(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// StripAllWhitespace returns s with all whitespace removed, as defined by unicode.IsSpace.
func StripAllWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}
