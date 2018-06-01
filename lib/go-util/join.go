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

func JoinErrsStr(errs []error) string {
	joined := JoinErrs(errs)

	if joined == nil {
		return ""
	}

	return joined.Error()
}

func ErrsToStrs(errs []error) []string {
	errorStrs := []string{}
	for _, errType := range errs {
		et := errType.Error()
		errorStrs = append(errorStrs, et)
	}
	return errorStrs
}

func JoinErrs(errs []error) error {
	return JoinErrsSep(errs, "")
}

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

func CamelToSnakeCase(s string) string {
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(s, "${1}_${2}"))
}
