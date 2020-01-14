package tce

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

import "strings"

type DSMatchType string

const (
	DSMatchTypeHostRegex     DSMatchType = "HOST_REGEXP"
	DSMatchTypePathRegex     DSMatchType = "PATH_REGEXP"
	DSMatchTypeSteeringRegex DSMatchType = "STEERING_REGEXP"
	DSMatchTypeHeaderRegex   DSMatchType = "HEADER_REGEXP"
	DSMatchTypeInvalid       DSMatchType = ""
)

// String returns a string representation of this delivery service match type.
func (t DSMatchType) String() string {
	switch t {
	case DSMatchTypeHostRegex:
		fallthrough
	case DSMatchTypePathRegex:
		fallthrough
	case DSMatchTypeSteeringRegex:
		fallthrough
	case DSMatchTypeHeaderRegex:
		return string(t)
	default:
		return "INVALID_MATCH_TYPE"
	}
}

// DSMatchTypeFromString returns a delivery service match type object from its string representation, or DSMatchTypeInvalid if the string is not a valid type.
func DSMatchTypeFromString(s string) DSMatchType {
	s = strings.ToLower(strings.Replace(s, "_", "", -1))
	switch s {
	case "hostregexp":
		return DSMatchTypeHostRegex
	case "pathregexp":
		return DSMatchTypePathRegex
	case "steeringregexp":
		return DSMatchTypeSteeringRegex
	case "headerregexp":
		return DSMatchTypeHeaderRegex
	default:
		return DSMatchTypeInvalid
	}
}
