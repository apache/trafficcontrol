package tc

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

// enum.go contains enumerations and strongly typed names.

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CDNName is the name of a CDN in Traffic Control.
type CDNName string

// TrafficMonitorName is the (short) hostname of a Traffic Monitor peer.
type TrafficMonitorName string

// String implements the fmt.Stringer interface.
func (t TrafficMonitorName) String() string {
	return string(t)
}

// CacheName is the (short) hostname of a cache server.
type CacheName string

// String implements the fmt.Stringer interface.
func (c CacheName) String() string {
	return string(c)
}

// CacheGroupName is the name of a Cache Group.
type CacheGroupName string

// DeliveryServiceName is the name of a Delivery Service.
//
// This has no attached semantics, and so it may be encountered in situations
// where it is the actual Display Name, but most often (as far as this author
// knows), it actually refers to a Delivery Service's XMLID.
type DeliveryServiceName string

// String implements the fmt.Stringer interface.
func (d DeliveryServiceName) String() string {
	return string(d)
}

// TopologyName is the name of a Topology.
type TopologyName string

// CacheType is the type (or tier) of a cache server.
type CacheType string

// The allowable values for a CacheType.
const (
	// CacheTypeEdge represents an edge-tier cache server.
	CacheTypeEdge = CacheType("EDGE")
	// CacheTypeMid represents a mid-tier cache server.
	CacheTypeMid = CacheType("MID")
	// CacheTypeInvalid represents an CacheType. Note this is the default
	// construction for a CacheType.
	CacheTypeInvalid = CacheType("")
)

// String returns a string representation of this CacheType, implementing the
// fmt.Stringer interface.
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

// CacheTypeFromString returns a CacheType structure from its string
// representation, or CacheTypeInvalid if the string is not a valid CacheType.
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

// IsValidCacheType returns true if the given string represents a valid cache type.
func IsValidCacheType(s string) bool {
	return CacheTypeFromString(s) != CacheTypeInvalid
}

// InterfaceName is the name of a server interface.
type InterfaceName string

// OriginLocationType is the Name of a Cache Group which represents an
// "Origin Location".
//
// There is no enforcement in Traffic Control anywhere that ensures than an
// ORG_LOC-Type Cache Group will only actually contain Origin servers.
//
// Deprecated: Prefer CacheGroupOriginTypeName for consistency.
const OriginLocationType = "ORG_LOC"

// AlgorithmConsistentHash is the name of the Multi-Site Origin hashing
// algorithm that performs consistent hashing on a set of parents.
const AlgorithmConsistentHash = "consistent_hash"

// MonitorTypeName is the Name of the Type which must be assigned to a server
// for it to be treated as a Traffic Monitor instance by ATC.
//
// "Rascal" is a legacy name for Traffic Monitor.
//
// Note that there is, in general, no guarantee that a Type with this name
// exists in Traffic Ops at any given time.
const MonitorTypeName = "RASCAL"

// MonitorProfilePrefix is a prefix which MUST appear on the Names of Profiles
// used by Traffic Monitor instances as servers in Traffic Ops, or they may not
// be treated properly.
//
// "Rascal" is a legacy name for Traffic Monitor.
//
// Deprecated: This should not be a requirement for TM instances to be treated
// properly, and new code should not check this.
const MonitorProfilePrefix = "RASCAL"

// RouterTypeName is the Name of the Type which must be assigned to a server
// for it to be treated as a Traffic Router instance by ATC.
//
// "CCR" is an acronym for a legacy name for Traffic Router.
//
// Note that there is, in general, no guarantee that a Type with this name
// exists in Traffic Ops at any given time.
const RouterTypeName = "CCR"

// EdgeTypePrefix is a prefix which MUST appear on the Names of Types used by
// edge-tier cache servers in order for them to be recognized as edge-tier
// cache servers by ATC.
const EdgeTypePrefix = "EDGE"

// MidTypePrefix is a prefix which MUST appear on the Names of Types used by
// mid-tier cache servers in order for them to be recognized as mid-tier cache
// servers by ATC.
const MidTypePrefix = "MID"

// OriginTypeName is the Name of the Type which must be assigned to a server
// for it to be treated as an Origin server by ATC.
//
// Note that there is, in general, no guarantee that a Type with this name
// exists in Traffic Ops at any given time.
const OriginTypeName = "ORG"

// These are the Names of the Types that must be used by various kinds of Cache
// Groups to ensure proper behavior.
//
// Note that there is, in general, no guarantee that a Type with any of these
// names exist in Traffic Ops at any given time.
//
// Note also that there is no enforcement in Traffic Control that a particular
// Type of Cache Group contains only or any of a particular Type or Types of
// server(s).
const (
	CacheGroupEdgeTypeName   = EdgeTypePrefix + "_LOC"
	CacheGroupMidTypeName    = MidTypePrefix + "_LOC"
	CacheGroupOriginTypeName = OriginTypeName + "_LOC"
)

// GlobalProfileName is the Name of a Profile that is treated specially in some
// respects by some components of ATC.
//
// Note that there is, in general, no guarantee that a Profile with this Name
// exists in Traffic Ops at any given time, nor that it will have or not have
// any particular set of assigned Parameters, nor that any set of Parameters
// that do happen to be assigned to it should it exist will have or not have
// particular ConfigFile or Secure or Value values, nor that any such
// Parameters should they exist would have Values that match any given pattern,
// nor that it has any particular Profile Type, nor that it is or is not
// assigned to any Delivery Service or server or server(s), nor that those
// servers be or not be of any particular Type should they exist. Traffic Ops
// promises much, but guarantees little.
const GlobalProfileName = "GLOBAL"

// ParameterName represents the name of a Traffic Ops Parameter.
//
// This has no additional attached semantics.
type ParameterName string

// UseRevalPendingParameterName is the name of a Parameter which tells whether
// or not Traffic Ops should use pending content invalidation jobs separately
// from traditionally "queued updates".
//
// Note that there is no guarantee that a Parameter with this name exists in
// Traffic Ops at any given time, nor that it will have or not have any
// particular ConfigFile or Secure or Value value, nor that should such a
// Parameter exist that its value will match or not match any given pattern,
// nor that it will or will not be assigned to any particular Profile or
// Profiles or Cache Groups.
//
// Deprecated: UseRevalPending was a feature flag introduced for ATC version 2,
// and newer code should just assume that pending revalidations are going to be
// fetched by t3c.
const UseRevalPendingParameterName = ParameterName("use_reval_pending")

// RefetchEnabled is the name of a Parameter used to determine if the
// Refetch feature is enabled. If enabled, this allows a consumer of the TO API
// to submit Refetch InvalidationJob types. These will subsequently be treated
// as a MISS by cache servers. Previously, the only capability was Refresh
// which was, in turn, treated as a STALE by cache servers. This value should
// be used with caution, since coupled with regex, could cause significant
// performance impacts by implementing cache servers if used incorrectly.
//
// Note that there is no guarantee that a Parameter with this name exists in
// Traffic Ops at any given time, and while it's implementation relies on
// a boolean Value, this it not guaranteed either.
const RefetchEnabled = ParameterName("refetch_enabled")

// ConfigFileName is a Parameter ConfigFile value.
//
// This has no additional attached semantics, and so while it is known to most
// frequently refer to a Parameter's ConfigFile, it may also be used to refer
// to the name of a literal configuration file within or without the context of
// Traffic Control, with unknown specificity (relative or absolute path?
// file:// URL?) and/or restrictions.
type ConfigFileName string

// GlobalConfigFileName is ConfigFile value which can cause a Parameter to be
// handled specially by Traffic Control components under certain circumstances.
//
// Note that there is no guarantee that a Parameter with this ConfigFile value
// exists in Traffic Ops at any given time, nor that any particular number of
// such Parameters is allowed or forbidden, nor that any such existing
// Parameters will have or not have any particular Name or Secure or Value
// value, nor that should such a Parameter exist that its value will match or
// not match any given pattern, nor that it will or will not be assigned to any
// particular Profile or Profiles or Cache Groups.
const GlobalConfigFileName = ConfigFileName("global")

// The allowed values for a Delivery Service's Query String Handling.
//
// These are prefixed "QueryStringIgnore" even though the values don't always
// indicate ignoring, because the database column is named "qstring_ignore".
const (
	QueryStringIgnoreUseInCacheKeyAndPassUp    = 0
	QueryStringIgnoreIgnoreInCacheKeyAndPassUp = 1
	QueryStringIgnoreDropAtEdge                = 2
)

// The allowed values for a Delivery Service's Range Request Handling.
const (
	RangeRequestHandlingDontCache         = 0
	RangeRequestHandlingBackgroundFetch   = 1
	RangeRequestHandlingCacheRangeRequest = 2
	RangeRequestHandlingSlice             = 3
)

// A DSTypeCategory defines the routing method for a Delivery Service, i.e.
// HTTP or DNS.
type DSTypeCategory string

const (
	// DSTypeCategoryHTTP represents an HTTP-routed Delivery Service.
	DSTypeCategoryHTTP = DSTypeCategory("http")
	// DSTypeCategoryDNS represents a DNS-routed Delivery Service.
	DSTypeCategoryDNS = DSTypeCategory("dns")
	// DSTypeCategoryInvalid represents an invalid Delivery Service routing
	// type. Note this is the default construction for a DSTypeCategory.
	DSTypeCategoryInvalid = DSTypeCategory("")
)

// String returns a string representation of this DSTypeCategory, implementing
// the fmt.Stringer interface.
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

// DSTypeCategoryFromString returns a DSTypeCategory from its string
// representation, or DSTypeCategoryInvalid if the string is not a valid
// DSTypeCategory.
//
// This is cAsE-iNsEnSiTiVe.
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

// GetDSTypeCategory returns the delivery service type category (either HTTP or DNS) of the given delivery service type.
func GetDSTypeCategory(dsType string) string {
	if strings.HasPrefix(dsType, "DNS") {
		return "DNS"
	}
	return "HTTP"
}

// These are the allowable values for the Signing Algorithm property of a
// Delivery Service.
const (
	SigningAlgorithmURLSig     = "url_sig"
	SigningAlgorithmURISigning = "uri_signing"
)

// These are the allowable values for the Protocol property of a Delivery
// Service.
const (
	// Indicates content will only be served using the unsecured HTTP Protocol.
	DSProtocolHTTP = 0
	// Indicates content will only be served using the secured HTTPS Protocol.
	DSProtocolHTTPS = 1
	// Indicates content will only be served over both HTTP and HTTPS.
	DSProtocolHTTPAndHTTPS = 2
	// Indicates content will only be served using the secured HTTPS Protocol,
	// and that unsecured HTTP requests will be directed to use HTTPS instead.
	DSProtocolHTTPToHTTPS = 3
)

// CacheStatus is a Name of some Status.
//
// More specifically, it is used here to enumerate the Statuses that are
// understood and acted upon in specific ways by Traffic Monitor.
//
// Note that the Statuses captured in this package as CacheStatus values are in
// no way restricted to use by cache servers, despite the name.
type CacheStatus string

// These are the allowable values of a CacheStatus.
//
// Note that there is no guarantee that a Status by any of these Names exists
// in Traffic Ops at any given time, nor that such a Status - should it exist
// - have any given Description or ID.
const (
	// CacheStatusAdminDown represents a cache server which has been
	// administratively marked as "down", but which should still appear in the
	// CDN.
	CacheStatusAdminDown = CacheStatus("ADMIN_DOWN")
	// CacheStatusOnline represents a cache server which should always be
	// considered online/available/healthy, irrespective of monitoring.
	// Non-cache servers also typically use this Status instead of REPORTED.
	CacheStatusOnline = CacheStatus("ONLINE")
	// CacheStatusOffline represents a cache server which should always be
	// considered offline/unavailable/unhealthy, irrespective of monitoring.
	CacheStatusOffline = CacheStatus("OFFLINE")
	// CacheStatusReported represents a cache server which is polled for health
	// by Traffic Monitor. The vast majority of cache servers should have this
	// Status.
	CacheStatusReported = CacheStatus("REPORTED")
	// CacheStatusPreProd represents a cache server that is not deployed to "production",
	// but is ready for it.
	CacheStatusPreProd = CacheStatus("PRE_PROD")
	// CacheStatusInvalid represents an unrecognized Status value. Note that
	// this is not actually "invalid", because Statuses may have any unique
	// name, not just those captured as CacheStatus values in this package.
	CacheStatusInvalid = CacheStatus("")
)

// String returns a string representation of this CacheStatus, implementing the
// fmt.Stringer interface.
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

// CacheStatusFromString returns a CacheStatus from its string representation,
// or CacheStatusInvalid if the string is not a valid CacheStatus.
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

// Protocol represents an ATC-supported content delivery protocol.
//
// Deprecated: This does not accurately model the Protocol of a Delivery
// Service, use DSProtocolHTTP, DSProtocolHTTPS, DSProtocolHTTPToHTTPS, and
// DSProtocolHTTPAndHTTPS instead.
type Protocol string

// The allowable values of a Protocol.
//
// Deprecated: These do not accurately model the Protocol of a Delivery
// Service, use DSProtocolHTTP, DSProtocolHTTPS, DSProtocolHTTPToHTTPS, and
// DSProtocolHTTPAndHTTPS instead.
const (
	// ProtocolHTTP represents the HTTP/1.1 protocol as specified in RFC2616.
	ProtocolHTTP = Protocol("http")
	// ProtocolHTTPS represents the HTTP/1.1 protocol over a TCP connection secured by TLS.
	ProtocolHTTPS = Protocol("https")
	// ProtocolHTTPtoHTTPS represents a redirection of unsecured HTTP requests to HTTPS.
	ProtocolHTTPtoHTTPS = Protocol("http to https")
	// ProtocolHTTPandHTTPS represents the use of both HTTP and HTTPS.
	ProtocolHTTPandHTTPS = Protocol("http and https")
	// ProtocolInvalid represents an invalid Protocol.
	ProtocolInvalid = Protocol("")
)

// String implements the fmt.Stringer interface.
func (p Protocol) String() string {
	switch p {
	case ProtocolHTTP:
		fallthrough
	case ProtocolHTTPS:
		fallthrough
	case ProtocolHTTPtoHTTPS:
		fallthrough
	case ProtocolHTTPandHTTPS:
		return string(p)
	default:
		return "INVALIDPROTOCOL"
	}
}

// ProtocolFromString parses a string and returns the corresponding Protocol.
func ProtocolFromString(s string) Protocol {
	switch strings.Replace(strings.ToLower(s), "_", " ", -1) {
	case "http":
		return ProtocolHTTP
	case "https":
		return ProtocolHTTPS
	case "http to https":
		return ProtocolHTTPtoHTTPS
	case "http and https":
		return ProtocolHTTPandHTTPS
	default:
		return ProtocolInvalid
	}
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (p *Protocol) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return errors.New("Protocol cannot be null")
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("JSON %s not quoted: %v", data, err)
	}
	*p = ProtocolFromString(s)
	if *p == ProtocolInvalid {
		return fmt.Errorf("%s is not a (supported) Protocol", s)
	}
	return nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (p Protocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// LocalizationMethod represents an enabled localization method for a Cache
// Group. The string values of this type should match the Traffic Ops values.
type LocalizationMethod string

// These are the allowable values of a LocalizationMethod.
const (
	LocalizationMethodCZ      = LocalizationMethod("CZ")
	LocalizationMethodDeepCZ  = LocalizationMethod("DEEP_CZ")
	LocalizationMethodGeo     = LocalizationMethod("GEO")
	LocalizationMethodInvalid = LocalizationMethod("INVALID")
)

// String returns a string representation of this LocalizationMethod,
// implementing the fmt.Stringer interface.
func (m LocalizationMethod) String() string {
	switch m {
	case LocalizationMethodCZ:
		return string(m)
	case LocalizationMethodDeepCZ:
		return string(m)
	case LocalizationMethodGeo:
		return string(m)
	default:
		return "INVALID"
	}
}

// LocalizationMethodFromString parses and returns a LocalizationMethod from
// its string representation.
//
// This is cAsE-iNsEnSiTiVe.
func LocalizationMethodFromString(s string) LocalizationMethod {
	switch strings.ToLower(s) {
	case "cz":
		return LocalizationMethodCZ
	case "deep_cz":
		return LocalizationMethodDeepCZ
	case "geo":
		return LocalizationMethodGeo
	default:
		return LocalizationMethodInvalid
	}
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (m *LocalizationMethod) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return errors.New("LocalizationMethod cannot be null")
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return errors.New(string(data) + " JSON not quoted")
	}
	*m = LocalizationMethodFromString(s)
	if *m == LocalizationMethodInvalid {
		return errors.New(s + " is not a LocalizationMethod")
	}
	return nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (m LocalizationMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// Scan implements the database/sql.Scanner interface.
func (m *LocalizationMethod) Scan(value interface{}) error {
	if value == nil {
		return errors.New("LocalizationMethod cannot be null")
	}
	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return errors.New("failed to scan LocalizationMethod: " + err.Error())
	}

	switch v := sv.(type) {
	case []byte:
		*m = LocalizationMethodFromString(string(v))
		if *m == LocalizationMethodInvalid {
			return errors.New(string(v) + " is not a valid LocalizationMethod")
		}
		return nil
	default:
		return fmt.Errorf("failed to scan LocalizationMethod, unsupported input type: %T", value)
	}
}

// DeepCachingType represents a Delivery Service's Deep Caching Type. The
// string values of this type should match the Traffic Ops values.
type DeepCachingType string

// These are the allowable values of a DeepCachingType. Note that unlike most
// "enumerated" constants exported by this package, the default construction
// yields a valid value.
const (
	DeepCachingTypeNever   = DeepCachingType("") // default value
	DeepCachingTypeAlways  = DeepCachingType("ALWAYS")
	DeepCachingTypeInvalid = DeepCachingType("INVALID")
)

// String returns a string representation of this DeepCachingType, implementing
// the fmt.Stringer interface.
func (t DeepCachingType) String() string {
	switch t {
	case DeepCachingTypeAlways:
		return string(t)
	case DeepCachingTypeNever:
		return "NEVER"
	default:
		return "INVALID"
	}
}

// DeepCachingTypeFromString returns a DeepCachingType from its string
// representation, or DeepCachingTypeInvalid if the string is not a valid
// DeepCachingTypeFromString.
func DeepCachingTypeFromString(s string) DeepCachingType {
	switch strings.ToLower(s) {
	case "always":
		return DeepCachingTypeAlways
	case "never":
		return DeepCachingTypeNever
	case "":
		// default when omitted
		return DeepCachingTypeNever
	default:
		return DeepCachingTypeInvalid
	}
}

// UnmarshalJSON unmarshals a JSON representation of a DeepCachingType (i.e. a
// string) or returns an error if the DeepCachingType is invalid.
//
// This implements the encoding/json.Unmarshaler interface.
func (t *DeepCachingType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = DeepCachingTypeNever
		return nil
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return errors.New(string(data) + " JSON not quoted")
	}
	*t = DeepCachingTypeFromString(s)
	if *t == DeepCachingTypeInvalid {
		return errors.New(string(data) + " is not a DeepCachingType")
	}
	return nil
}

// MarshalJSON marshals into a JSON representation, implementing the
// encoding/json.Marshaler interface.
func (t DeepCachingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// A SteeringType is the Name of the Type of a Steering Target.
type SteeringType string

// These are the allowable values of a SteeringType.
//
// Note that, in general, there is no guarantee that a Type by any of these
// Names exists in Traffic Ops at any given time, nor that any such Types
// - should they exist - will have any particular UseInTable value, nor that
// the Types assigned to Steering Target relationships will be representable
// by these values.
const (
	SteeringTypeOrder     SteeringType = "STEERING_ORDER"
	SteeringTypeWeight    SteeringType = "STEERING_WEIGHT"
	SteeringTypeGeoOrder  SteeringType = "STEERING_GEO_ORDER"
	SteeringTypeGeoWeight SteeringType = "STEERING_GEO_WEIGHT"
	SteeringTypeInvalid   SteeringType = ""
)

// SteeringTypeFromString parses a string to return the corresponding
// SteeringType.
//
// Warning: This is cAsE-iNsEnSiTiVe, but some components of Traffic Ops may
// compare the Names of Steering Target Types in case-sensitive ways, so this
// may obscure values that will be later detected as invalid, depending on how
// it's used.
func SteeringTypeFromString(s string) SteeringType {
	s = strings.ToLower(strings.Replace(s, "_", "", -1))
	switch s {
	case "steeringorder":
		return SteeringTypeOrder
	case "steeringweight":
		return SteeringTypeWeight
	case "steeringgeoorder":
		return SteeringTypeGeoOrder
	case "steeringgeoweight":
		return SteeringTypeGeoWeight
	default:
		return SteeringTypeInvalid
	}
}

// String returns a string representation of this SteeringType, implementing
// the fmt.Stringer interface.
func (t SteeringType) String() string {
	switch t {
	case SteeringTypeOrder:
		fallthrough
	case SteeringTypeWeight:
		fallthrough
	case SteeringTypeGeoOrder:
		fallthrough
	case SteeringTypeGeoWeight:
		return string(t)
	default:
		return "INVALID"
	}
}

// A FederationResolverType is the Name of a Type of a Federation Resolver
// Mapping.
type FederationResolverType string

// These are the allowable values of a FederationResolverType.
//
// Note that, in general, there is no guarantee that a Type by any of these
// Names exists in Traffic Ops at any given time, nor that any such Types
// - should they exist - will have any particular UseInTable value, nor that
// the Types assigned to Federation Resolver Mappings will be representable
// by these values.
const (
	FederationResolverType4       = FederationResolverType("RESOLVE4")
	FederationResolverType6       = FederationResolverType("RESOLVE6")
	FederationResolverTypeInvalid = FederationResolverType("")
)

// String imlpements the fmt.Stringer interface.
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

// FederationResolverTypeFromString parses a string and returns the
// corresponding FederationResolverType.
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

// DSType is the Name of a Type used by a Delivery Service.
type DSType string

// These are the allowable values for a DSType.
//
// Note that, in general, there is no guarantee that a Type by any of these
// Names exists in Traffic Ops at any given time, nor that any such Types
// - should they exist - will have any particular UseInTable value, nor that
// the Types assigned to Delivery Services will be representable by these
// values.
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

// String returns a string representation of this DSType, implementing the
// fmt.Stringer interface.
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

// DSTypeFromString returns a DSType from its string representation, or
// DSTypeInvalid if the string is not a valid DSType.
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

// UsesDNSSECKeys returns whether a Delivery Service of a Type that has a Name
// that is this DSType uses or needs DNSSEC keys.
func (t DSType) UsesDNSSECKeys() bool {
	return t.IsDNS() || t.IsHTTP() || t.IsSteering()
}

// IsHTTP returns whether a Delivery Service of a Type that has a Name
// that is this DSType is HTTP-routed.
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

// IsDNS returns whether a Delivery Service of a Type that has a Name
// that is this DSType is DNS-routed.
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

// IsSteering returns whether a Delivery Service of a Type that has a Name
// that is this DSType is Steering-routed.
func (t DSType) IsSteering() bool {
	switch t {
	case DSTypeSteering:
		fallthrough
	case DSTypeClientSteering:
		return true
	}
	return false
}

// HasSSLKeys returns whether Dlivery Services of this Type can have SSL keys.
func (t DSType) HasSSLKeys() bool {
	return t.IsHTTP() || t.IsDNS() || t.IsSteering()
}

// IsLive returns whether Delivery Services of this Type are "live".
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

// IsNational returns whether Delivery Services of this Type are "national".
func (t DSType) IsNational() bool {
	switch t {
	case DSTypeDNSLiveNational:
		fallthrough
	case DSTypeHTTPLiveNational:
		return true
	}
	return false
}

// UsesMidCache returns whether Delivery Services of this Type use mid-tier
// cache servers.
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

// DSTypeLiveNationalSuffix is the suffix that Delivery Services which are both
// "live" and "national" MUST have at the end of the Name(s) of the Type(s)
// they use in order to be treated properly by ATC.
//
// Deprecated: Use DSType.IsLive and DSType.IsNational instead.
const DSTypeLiveNationalSuffix = "_LIVE_NATNL"

// DSTypeLiveSuffix is the suffix that Delivery Services which are "live" - but
// not "national" (maybe?) - MUST have at the end of the Name(s) of the Type(s)
// they use in order to be treated properly by ATC.
//
// Deprecated: Use DSType.IsLive and DSType.IsNational instead.
const DSTypeLiveSuffix = "_LIVE"

// A QStringIgnore defines how to treat the URL query string for requests
// for a given Delivery Service's content.
//
// This enum's String function returns the numeric representation, because it
// is a legacy database value, and the number should be kept for both database
// and API JSON uses. For the same reason, this enum has no FromString
// function.
//
// One should normally use the QueryStringIgnoreUseInCacheKeyAndPassUp,
// QueryStringIgnoreIgnoreInCacheKeyAndPassUp, and
// QueryStringIgnoreDropAtEdge constants instead, as they have the same type as
// the Delivery Service field they represent.
type QStringIgnore int

// These are the allowable values for a QStringIgnore.
const (
	QStringIgnoreUseInCacheKeyAndPassUp    QStringIgnore = 0
	QStringIgnoreIgnoreInCacheKeyAndPassUp QStringIgnore = 1
	QStringIgnoreDrop                      QStringIgnore = 2
)

// String returns the string number of the QStringIgnore value, implementing
// the fmt.Stringer interface.
//
// Note this returns the number, not a human-readable value, because
// QStringIgnore is a legacy database sigil, and both database and API JSON
// uses should use the number. This also returns 'INVALID' for unknown values,
// to fail fast in the event of bad data.
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

// A DSMatchType is the Name of a Type of a Regular Expression ("Regex") used
// by a Delivery Service.
type DSMatchType string

// These are the allowed values for a DSMatchType.
//
// Note that, in general, there is no guarantee that a Type by any of these
// Names exists in Traffic Ops at any given time, nor that any such Types
// - should they exist - will have any particular UseInTable value, nor that
// the Types assigned to Delivery Service Regexes will be representable
// by these values.
const (
	DSMatchTypeHostRegex     DSMatchType = "HOST_REGEXP"
	DSMatchTypePathRegex     DSMatchType = "PATH_REGEXP"
	DSMatchTypeSteeringRegex DSMatchType = "STEERING_REGEXP"
	DSMatchTypeHeaderRegex   DSMatchType = "HEADER_REGEXP"
	DSMatchTypeInvalid       DSMatchType = ""
)

// String returns a string representation of this DSMatchType, implementing the
// fmt.Stringer interface.
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

// DSMatchTypeFromString returns a DSMatchType from its string representation,
// or DSMatchTypeInvalid if the string is not a valid type.
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
