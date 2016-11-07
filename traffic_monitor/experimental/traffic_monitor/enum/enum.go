// Package enum contains enumerations and strongly typed names.
// The names are an experiment with strong typing of string types. The primary goal is to make code more self-documenting, especially map keys. If peole don't like it, we can get rid of it.
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


import (
	"strings"
)

// TrafficMonitorName is the hostname of a Traffic Monitor peer.
type TrafficMonitorName string

// CacheName is the hostname of a CDN cache.
type CacheName string

// CacheGroupName is the name of a CDN cachegroup.
type CacheGroupName string

// DeliveryServiceName is the name of a CDN delivery service.
type DeliveryServiceName string

// CacheType is the type (or tier) of a CDN cache.
type CacheType string

const (
	// CacheTypeEdge represents an edge cache.
	CacheTypeEdge = CacheType("EDGE")
	// CacheTypeMid represents a mid cache.
	CacheTypeMid = CacheType("MID")
	// CacheTypeInvalid represents an cache type enumeration. Note this is the default construction for a CacheType.
	CacheTypeInvalid = CacheType("")
)

// String returns a string representation of this cache type.
func (t CacheType) String() string {
	switch t {
	case CacheTypeEdge:
		return "EDGE"
	case CacheTypeMid:
		return "MID"
	default:
		return "INVALID"
	}
}

// CacheTypeFromString returns a cache type object from its string representation, or CacheTypeInvalid if the string is not a valid type.
func CacheTypeFromString(s string) CacheType {
	s = strings.ToLower(s)
	switch s {
	case "edge":
		return CacheTypeEdge
	case "mid":
		return CacheTypeMid
	default:
		return CacheTypeInvalid
	}
}

// DSType is the Delivery Service type. HTTP, DNS, etc.
type DSType string

const (
	// DSTypeHTTP represents an HTTP delivery service
	DSTypeHTTP = DSType("http")
	// DSTypeDNS represents a DNS delivery service
	DSTypeDNS = DSType("dns")
	// DSTypeInvalid represents an invalid delivery service type enumeration. Note this is the default construction for a DSType.
	DSTypeInvalid = DSType("")
)

// String returns a string representation of this delivery service type.
func (t DSType) String() string {
	switch t {
	case DSTypeHTTP:
		return "HTTP"
	case DSTypeDNS:
		return "DNS"
	default:
		return "INVALID"
	}
}

// DSTypeFromString returns a delivery service type object from its string representation, or DSTypeInvalid if the string is not a valid type.
func DSTypeFromString(s string) DSType {
	s = strings.ToLower(s)
	switch s {
	case "http":
		return DSTypeHTTP
	case "dns":
		return DSTypeDNS
	default:
		return DSTypeInvalid
	}
}
