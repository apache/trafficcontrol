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
	"regexp"
	"strings"
)

// JoinErrsStr returns a string representation of all of the passed errors.
//
// This is equivalent to calling JoinErrs(errs).Error(), but in the case that
// JoinErrs returns nil that would panic. This checks for that case and instead
// returns an empty string.
func JoinErrsStr(errs []error) string {
	joined := JoinErrs(errs)

	if joined == nil {
		return ""
	}

	return joined.Error()
}

// ErrsToStrs converts a slice of errors to a slice of their string
// representations.
func ErrsToStrs(errs []error) []string {
	errorStrs := []string{}
	for _, errType := range errs {
		et := errType.Error()
		errorStrs = append(errorStrs, et)
	}
	return errorStrs
}

// JoinErrs joins the passed errors into a single error with a message that is
// a combination of their error messages.
//
// This is equivalent to calling JoinErrsSep(errs, "").
func JoinErrs(errs []error) error {
	return JoinErrsSep(errs, "")
}

// JoinErrsSep joins the passed errors into a single error with a message that
// is a combination of their error messages, joined by the given separator.
//
// If the given separator is an empty string, the default (", ") is used.
//
// Note that this DOES NOT preserve error identity. For example:
//
//	err := JoinErrsSep([]error{sql.ErrNoRows, errors.New("foo")})
//	fmt.Println(errors.Is(err, sql.ErrNoRows))
//	// Output: false
func JoinErrsSep(errs []error, separator string) error {
	if separator == "" {
		separator = ", "
	}

	joinedErrors := ""

	for _, err := range errs {
		if err != nil {
			joinedErrors += err.Error() + separator
		}
	}

	if len(joinedErrors) == 0 {
		return nil
	}

	joinedErrors = joinedErrors[:len(joinedErrors)-len(separator)] // strip trailing separator

	return fmt.Errorf("%s", joinedErrors)
}

// CamelToSnakeCase returns a case transformation of the input string from
// assumed "camelCase" (or "PascalCase") to "snake_case".
//
// Note that the transformation applied to strings that contain non-"word"
// characters is undefined. Also, this doesn't handle names that contain
// abbreviations, initialisms, and/or acronyms very well in general, like e.g.
// "IPAddress".
func CamelToSnakeCase(s string) string {
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(s, "${1}_${2}"))
}
