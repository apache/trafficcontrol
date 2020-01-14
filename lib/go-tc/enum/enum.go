// Package enum contains enumerations and strongly typed names.
//
// These enums should be treated as enumerables, and MUST NOT be cast as anything else (integer, strings, etc). Enums MUST NOT be compared to strings or integers via casting. Enumerable data SHOULD be stored as the enumeration, not as a string or number. The *only* reason they are internally represented as strings, is to make them implicitly serialize to human-readable JSON. They should not be treated as strings. Casting or storing strings or numbers defeats the purpose of enum safety and conveniences.
//
// When storing enumumerable data in memory, it SHOULD be converted to and stored as an enum via the corresponding `FromString` function, checked whether the conversion failed and Invalid values handled, and valid data stored as the enum. This guarantees stored data is valid, and catches invalid input as soon as possible.
//
// When adding new enum types, enums should be internally stored as strings, so they implicitly serialize as human-readable JSON, unless the performance or memory of integers is necessary (it almost certainly isn't). Enums should always have the "invalid" value as the empty string (or 0), so default-initialized enums are invalid.
// Enums should always have a FromString() conversion function, to convert input data to enums. Conversion functions should usually be case-insensitive, and may ignore underscores or hyphens, depending on the use case.
//
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
	"strconv"
	"strings"
)

// CDNName is the name of a CDN in Traffic Control.
type CDNName string

// TrafficMonitorName is the hostname of a Traffic Monitor peer.
type TrafficMonitorName string

// CacheName is the hostname of a CDN cache.
type CacheName string

// CacheGroupName is the name of a CDN cachegroup.
type CacheGroupName string

// DeliveryServiceName is the name of a CDN delivery service.
type DeliveryServiceName string

const AlgorithmConsistentHash = "consistent_hash"

const MonitorTypeName = "RASCAL"
const MonitorProfilePrefix = "RASCAL"
const RouterTypeName = "CCR"
const EdgeTypePrefix = "EDGE"
const MidTypePrefix = "MID"

const OriginTypeName = "ORG"

const CacheGroupOriginTypeName = "ORG_LOC"

const GlobalProfileName = "GLOBAL"

func (c CacheName) String() string {
	return string(c)
}

func (t TrafficMonitorName) String() string {
	return string(t)
}

func (d DeliveryServiceName) String() string {
	return string(d)
}

// These are prefixed "QueryStringIgnore" even though the values don't always indicate ignoring, because the database column is named "qstring_ignore"

const QueryStringIgnoreUseInCacheKeyAndPassUp = 0
const QueryStringIgnoreIgnoreInCacheKeyAndPassUp = 1
const QueryStringIgnoreDropAtEdge = 2

const RangeRequestHandlingDontCache = 0
const RangeRequestHandlingBackgroundFetch = 1
const RangeRequestHandlingCacheRangeRequest = 2

const SigningAlgorithmURLSig = "url_sig"
const SigningAlgorithmURISigning = "uri_signing"

type FederationResolverType string

const (
	FederationResolverType4       = FederationResolverType("RESOLVE4")
	FederationResolverType6       = FederationResolverType("RESOLVE6")
	FederationResolverTypeInvalid = FederationResolverType("")
)

func (t FederationResolverType) String() string {
	switch t {
	case FederationResolverType4:
		fallthrough
	case FederationResolverType6:
		return string(t)
	default:
		return "INVALID"
	}
}

func FederationResolverTypeFromString(s string) FederationResolverType {
	switch strings.ToLower(s) {
	case "resolve4":
		return FederationResolverType4
	case "resolve6":
		return FederationResolverType6
	default:
		return FederationResolverTypeInvalid
	}
}

// QStringIgnore is an entry in the delivery_service table qstring_ignore column, and represents how to treat the URL query string for requests to that delivery service.
// This enum's String function returns the numeric representation, because it is a legacy database value, and the number should be kept for both database and API JSON uses. For the same reason, this enum has no FromString function.
type QStringIgnore int

const (
	QStringIgnoreUseInCacheKeyAndPassUp    QStringIgnore = 0
	QStringIgnoreIgnoreInCacheKeyAndPassUp QStringIgnore = 1
	QStringIgnoreDrop                      QStringIgnore = 2
)

// String returns the string number of the QStringIgnore value.
// Note this returns the number, not a human-readable value, because QStringIgnore is a legacy database sigil, and both database and API JSON uses should use the number. This also returns 'INVALID' for unknown values, to fail fast in the event of bad data.
func (e QStringIgnore) String() string {
	switch e {
	case QStringIgnoreUseInCacheKeyAndPassUp:
		fallthrough
	case QStringIgnoreIgnoreInCacheKeyAndPassUp:
		fallthrough
	case QStringIgnoreDrop:
		return strconv.Itoa(int(e))
	default:
		return "INVALID"
	}
}
