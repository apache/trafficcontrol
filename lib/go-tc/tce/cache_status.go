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
