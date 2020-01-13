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

// DSType is the Delivery Service type.
type DSType string

const (
	DSTypeClientSteering   DSType = "CLIENT_STEERING"
	DSTypeDNS              DSType = "DNS"
	DSTypeDNSLive          DSType = "DNS_LIVE"
	DSTypeDNSLiveNational  DSType = "DNS_LIVE_NATNL"
	DSTypeHTTP             DSType = "HTTP"
	DSTypeHTTPLive         DSType = "HTTP_LIVE"
	DSTypeHTTPLiveNational DSType = "HTTP_LIVE_NATNL"
	DSTypeHTTPNoCache      DSType = "HTTP_NO_CACHE"
	DSTypeSteering         DSType = "STEERING"
	DSTypeAnyMap           DSType = "ANY_MAP"
	DSTypeInvalid          DSType = ""
)

// String returns a string representation of this delivery service type.
func (t DSType) String() string {
	switch t {
	case DSTypeHTTPNoCache:
		fallthrough
	case DSTypeDNS:
		fallthrough
	case DSTypeDNSLive:
		fallthrough
	case DSTypeHTTP:
		fallthrough
	case DSTypeDNSLiveNational:
		fallthrough
	case DSTypeAnyMap:
		fallthrough
	case DSTypeHTTPLive:
		fallthrough
	case DSTypeSteering:
		fallthrough
	case DSTypeHTTPLiveNational:
		fallthrough
	case DSTypeClientSteering:
		return string(t)
	default:
		return "INVALID"
	}
}

// DSTypeFromString returns a delivery service type object from its string representation, or DSTypeInvalid if the string is not a valid type.
func DSTypeFromString(s string) DSType {
	s = strings.ToLower(strings.Replace(s, "_", "", -1))
	switch s {
	case "httpnocache":
		return DSTypeHTTPNoCache
	case "dns":
		return DSTypeDNS
	case "dnslive":
		return DSTypeDNSLive
	case "http":
		return DSTypeHTTP
	case "dnslivenatnl":
		return DSTypeDNSLiveNational
	case "anymap":
		return DSTypeAnyMap
	case "httplive":
		return DSTypeHTTPLive
	case "steering":
		return DSTypeSteering
	case "httplivenatnl":
		return DSTypeHTTPLiveNational
	case "clientsteering":
		return DSTypeClientSteering
	default:
		return DSTypeInvalid
	}
}

// IsHTTP returns whether the DSType is an HTTP category.
func (t DSType) IsHTTP() bool {
	switch t {
	case DSTypeHTTP:
		fallthrough
	case DSTypeHTTPLive:
		fallthrough
	case DSTypeHTTPLiveNational:
		fallthrough
	case DSTypeHTTPNoCache:
		return true
	}
	return false
}

// IsDNS returns whether the DSType is a DNS category.
func (t DSType) IsDNS() bool {
	switch t {
	case DSTypeDNS:
		fallthrough
	case DSTypeDNSLive:
		fallthrough
	case DSTypeDNSLiveNational:
		return true
	}
	return false
}

// IsSteering returns whether the DSType is a Steering category
func (t DSType) IsSteering() bool {
	switch t {
	case DSTypeSteering:
		fallthrough
	case DSTypeClientSteering:
		return true
	}
	return false
}

// HasSSLKeys returns whether delivery services of this type have SSL keys.
func (t DSType) HasSSLKeys() bool {
	return t.IsHTTP() || t.IsDNS() || t.IsSteering()
}

// IsLive returns whether delivery services of this type are "live".
func (t DSType) IsLive() bool {
	switch t {
	case DSTypeDNSLive:
		fallthrough
	case DSTypeDNSLiveNational:
		fallthrough
	case DSTypeHTTPLive:
		fallthrough
	case DSTypeHTTPLiveNational:
		return true
	}
	return false
}

// IsLive returns whether delivery services of this type are "national".
func (t DSType) IsNational() bool {
	switch t {
	case DSTypeDNSLiveNational:
		fallthrough
	case DSTypeHTTPLiveNational:
		return true
	}
	return false
}

// UsesMidCache returns whether delivery services of this type use mid-tier caches
func (t DSType) UsesMidCache() bool {
	switch t {
	case DSTypeDNSLive:
		fallthrough
	case DSTypeHTTPLive:
		fallthrough
	case DSTypeHTTPNoCache:
		return false
	}
	return true
}

const DSTypeLiveNationalSuffix = "_LIVE_NATNL"
const DSTypeLiveSuffix = "_LIVE"
