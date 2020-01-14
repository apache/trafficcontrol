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
