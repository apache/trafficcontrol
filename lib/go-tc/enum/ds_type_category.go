package enum

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

// DSTypeCategory is the Delivery Service type category: HTTP or DNS
type DSTypeCategory string

const (
	// DSTypeCategoryHTTP represents an HTTP delivery service
	DSTypeCategoryHTTP = DSTypeCategory("http")
	// DSTypeCategoryDNS represents a DNS delivery service
	DSTypeCategoryDNS = DSTypeCategory("dns")
	// DSTypeCategoryInvalid represents an invalid delivery service type enumeration. Note this is the default construction for a DSTypeCategory.
	DSTypeCategoryInvalid = DSTypeCategory("")
)

// String returns a string representation of this delivery service type.
func (t DSTypeCategory) String() string {
	switch t {
	case DSTypeCategoryHTTP:
		return "HTTP"
	case DSTypeCategoryDNS:
		return "DNS"
	default:
		return "INVALIDDSTYPE"
	}
}

// DSTypeCategoryFromString returns a delivery service type object from its string representation, or DSTypeCategoryInvalid if the string is not a valid type.
func DSTypeCategoryFromString(s string) DSTypeCategory {
	s = strings.ToLower(s)
	switch s {
	case "http":
		return DSTypeCategoryHTTP
	case "dns":
		return DSTypeCategoryDNS
	default:
		return DSTypeCategoryInvalid
	}
}
