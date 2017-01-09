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
		return "INVALIDCACHETYPE"
	}
}

// CacheTypeFromString returns a cache type object from its string representation, or CacheTypeInvalid if the string is not a valid type.
func CacheTypeFromString(s string) CacheType {
	s = strings.ToLower(s)
	if strings.HasPrefix(s, "edge") {
		return CacheTypeEdge
	}
	if strings.HasPrefix(s, "mid") {
		return CacheTypeMid
	}
	return CacheTypeInvalid
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
		return "INVALIDDSTYPE"
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

// CacheStatus represents the Traffic Server status set in Traffic Ops (online, offline, admin_down, reported). The string values of this type should match the Traffic Ops values.
type CacheStatus string

const (
	// CacheStatusAdminDown represents a cache which has been administratively marked as down, but which should still appear in the CDN (Traffic Server, Traffic Monitor, Traffic Router).
	CacheStatusAdminDown = CacheStatus("ADMIN_DOWN")
	// CacheStatusOnline represents a cache which has been marked as Online in Traffic Ops, irrespective of monitoring. Traffic Monitor will always flag these caches as available.
	CacheStatusOnline = CacheStatus("ONLINE")
	// CacheStatusOffline represents a cache which has been marked as Offline in Traffic Ops. These caches will not be returned in any endpoint, and Traffic Monitor acts like they don't exist.
	CacheStatusOffline = CacheStatus("OFFLINE")
	// CacheStatusReported represents a cache which has been marked as Reported in Traffic Ops. These caches are polled for health and returned in endpoints as available or unavailable based on bandwidth, response time, and other factors. The vast majority of caches should be Reported.
	CacheStatusReported = CacheStatus("REPORTED")
	// CacheStatusInvalid represents an invalid status enumeration.
	CacheStatusInvalid = CacheStatus("")
)

// String returns a string representation of this cache status
func (t CacheStatus) String() string {
	switch t {
	case CacheStatusAdminDown:
		fallthrough
	case CacheStatusOnline:
		fallthrough
	case CacheStatusOffline:
		fallthrough
	case CacheStatusReported:
		return string(t)
	default:
		return "INVALIDCACHESTATUS"
	}
}

// CacheStatusFromString returns a CacheStatus from its string representation, or CacheStatusInvalid if the string is not a valid type.
func CacheStatusFromString(s string) CacheStatus {
	s = strings.ToLower(s)
	switch s {
	case "admin_down":
		fallthrough
	case "admindown":
		return CacheStatusAdminDown
	case "offline":
		return CacheStatusOffline
	case "online":
		return CacheStatusOnline
	case "reported":
		return CacheStatusReported
	default:
		return CacheStatusInvalid
	}
}
