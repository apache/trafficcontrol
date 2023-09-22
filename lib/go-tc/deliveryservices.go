package tc

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// DefaultRoutingName is the Routing Name given to Delivery Services upon
// creation through the Traffic Ops API if a Routing Name is not specified in
// the request body.
const DefaultRoutingName = "cdn"

// DefaultMaxRequestHeaderBytes is the Max Header Bytes given to Delivery
// Services upon creation through the Traffic Ops API if Max Request Header
// Bytes is not specified in the request body.
const DefaultMaxRequestHeaderBytes = 0

// MinRangeSliceBlockSize is the minimum allowed value for a Delivery Service's
// Range Slice Block Size, in bytes. This is 256KiB.
const MinRangeSliceBlockSize = 262144

// MaxRangeSliceBlockSize is the maximum allowed value for a Delivery Service's
// Range Slice Block Size, in bytes. This is 32MiB.
const MaxRangeSliceBlockSize = 33554432

// DeliveryServiceActiveState is an "enumerated" type which encodes the valid
// values of a Delivery Service's 'Active' property (v3.0+).
type DeliveryServiceActiveState string

// These names of URL query string parameters are not allowed to be in a
// Delivery Service's "ConsistentHashQueryParams" set, because they collide with
// query string parameters reserved for use by Traffic Router.
const (
	ReservedConsistentHashingQueryParameterFormat              = "format"
	ReservedConsistentHashingQueryParameterTRRED               = "trred"
	ReservedConsistentHashingQueryParameterFakeClientIPAddress = "fakeClientIpAddress"
)

// A DeliveryServiceActiveState describes the availability of Delivery Service
// content from the perspective of caching servers and Traffic Routers.
const (
	// Traffic Router routes clients for this Delivery Service and cache
	// servers are configured to serve its content.
	DSActiveStateActive = DeliveryServiceActiveState("ACTIVE")
	// Traffic Router does not route for this Delivery Service and cache
	// servers cannot serve its content.
	DSActiveStateInactive = DeliveryServiceActiveState("INACTIVE")
	// Traffic Router does not route for this Delivery Service, but cache
	// servers are configured to serve its content.
	DSActiveStatePrimed = DeliveryServiceActiveState("PRIMED")
)

// DeliveryServicesResponseV30 is the type of a response from the
// /api/3.0/deliveryservices Traffic Ops endpoint.
//
// TODO: Move these into the respective clients?
type DeliveryServicesResponseV30 struct {
	Response []DeliveryServiceNullableV30 `json:"response"`
	Alerts
}

// DeliveryServicesResponseV40 is the type of a response from the
// /api/4.0/deliveryservices Traffic Ops endpoint.
type DeliveryServicesResponseV40 struct {
	Response []DeliveryServiceV40 `json:"response"`
	Alerts
}

// DeliveryServicesResponseV41 is the type of a response from the
// /api/4.1/deliveryservices Traffic Ops endpoint.
type DeliveryServicesResponseV41 struct {
	Response []DeliveryServiceV41 `json:"response"`
	Alerts
}

// DeliveryServicesResponseV4 is the type of a response from the
// /api/4.x/deliveryservices Traffic Ops endpoint.
// It always points to the type for the latest minor version of APIv4.
type DeliveryServicesResponseV4 = DeliveryServicesResponseV41

// DeliveryServicesNullableResponse roughly models the structure of responses
// from Traffic Ops to GET requests made to its
// /servers/{{ID}}/deliveryservices and /deliveryservices API endpoints.
//
// "Roughly" because although that's what it's used for, this type cannot
// actually represent those accurately, because its representation is tied to a
// version of the API that no longer exists - DO NOT USE THIS, it WILL drop
// data that the API returns.
//
// Deprecated: Please only use the versioned structures.
type DeliveryServicesNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts
}

// CreateDeliveryServiceNullableResponse roughly models the structure of
// responses from Traffic Ops to POST requests made to its /deliveryservices
// API endpoint.
//
// "Roughly" because although that's what it's used for, this type cannot
// actually represent those accurately, because its representation is tied to a
// version of the API that no longer exists - DO NOT USE THIS, it WILL drop
// data that the API returns.
//
// Deprecated: Please only use the versioned structures.
type CreateDeliveryServiceNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts
}

// UpdateDeliveryServiceNullableResponse roughly models the structure of
// responses from Traffic Ops to PUT requests made to its
// /deliveryservices/{{ID}} API endpoint.
//
// "Roughly" because although that's what it's used for, this type cannot
// actually represent those accurately, because its representation is tied to a
// version of the API that no longer exists - DO NOT USE THIS, it WILL drop
// data that the API returns.
//
// Deprecated: Please only use the versioned structures.
type UpdateDeliveryServiceNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts
}

// DeleteDeliveryServiceResponse is the type of a response from Traffic Ops to
// DELETE requests made to its /deliveryservices/{{ID}} API endpoint.
type DeleteDeliveryServiceResponse struct {
	Alerts
}

// DeliveryService structures represent a Delivery Service as it is exposed
// through the Traffic Ops API at version 1.4 - which no longer exists.
//
// DO NOT USE THIS - it ONLY still exists because it is used in the
// DeliveryServiceRequest type that is still in use by Traffic Ops's Go client
// for API versions 2 and 3 - and even that is incorrect and DOES DROP DATA.
//
// Deprecated: Instead use the appropriate structure for the version of the
// Traffic Ops API being worked with, e.g. DeliveryServiceV4.
type DeliveryService struct {
	DeliveryServiceV13
	MaxOriginConnections      int      `json:"maxOriginConnections" db:"max_origin_connections"`
	ConsistentHashRegex       string   `json:"consistentHashRegex"`
	ConsistentHashQueryParams []string `json:"consistentHashQueryParams"`
}

// DeliveryServiceV13 structures represent a Delivery Service as it is exposed
// through the Traffic Ops API at version 1.3 - which no longer exists.
//
// DO NOT USE THIS - it ONLY still exists because it is nested within the
// structure of the DeliveryService type.
//
// Deprecated: Instead, use the appropriate structure for the version of the
// Traffic Ops API being worked with, e.g. DeliveryServiceV4.
type DeliveryServiceV13 struct {
	DeliveryServiceV11
	DeepCachingType   DeepCachingType `json:"deepCachingType"`
	FQPacingRate      int             `json:"fqPacingRate,omitempty"`
	SigningAlgorithm  string          `json:"signingAlgorithm" db:"signing_algorithm"`
	Tenant            string          `json:"tenant"`
	TRRequestHeaders  string          `json:"trRequestHeaders,omitempty"`
	TRResponseHeaders string          `json:"trResponseHeaders,omitempty"`
}

// DeliveryServiceV11 contains the information relating to a Delivery Service
// that was around in version 1.1 of the Traffic Ops API.
//
// DO NOT USE THIS - it ONLY still exists because it is nested within the
// structure of the DeliveryServiceV13 type.
//
// Deprecated: Instead, use the appropriate structure for the version of the
// Traffic Ops API being worked with, e.g. DeliveryServiceV4.
type DeliveryServiceV11 struct {
	Active                   bool                   `json:"active"`
	AnonymousBlockingEnabled bool                   `json:"anonymousBlockingEnabled"`
	CacheURL                 string                 `json:"cacheurl"`
	CCRDNSTTL                int                    `json:"ccrDnsTtl"`
	CDNID                    int                    `json:"cdnId"`
	CDNName                  string                 `json:"cdnName"`
	CheckPath                string                 `json:"checkPath"`
	DeepCachingType          DeepCachingType        `json:"deepCachingType"`
	DisplayName              string                 `json:"displayName"`
	DNSBypassCname           string                 `json:"dnsBypassCname"`
	DNSBypassIP              string                 `json:"dnsBypassIp"`
	DNSBypassIP6             string                 `json:"dnsBypassIp6"`
	DNSBypassTTL             int                    `json:"dnsBypassTtl"`
	DSCP                     int                    `json:"dscp"`
	EdgeHeaderRewrite        string                 `json:"edgeHeaderRewrite"`
	ExampleURLs              []string               `json:"exampleURLs"`
	GeoLimit                 int                    `json:"geoLimit"`
	GeoProvider              int                    `json:"geoProvider"`
	GlobalMaxMBPS            int                    `json:"globalMaxMbps"`
	GlobalMaxTPS             int                    `json:"globalMaxTps"`
	HTTPBypassFQDN           string                 `json:"httpBypassFqdn"`
	ID                       int                    `json:"id"`
	InfoURL                  string                 `json:"infoUrl"`
	InitialDispersion        float32                `json:"initialDispersion"`
	IPV6RoutingEnabled       bool                   `json:"ipv6RoutingEnabled"`
	LastUpdated              *TimeNoMod             `json:"lastUpdated" db:"last_updated"`
	LogsEnabled              bool                   `json:"logsEnabled"`
	LongDesc                 string                 `json:"longDesc"`
	LongDesc1                string                 `json:"longDesc1"`
	LongDesc2                string                 `json:"longDesc2"`
	MatchList                []DeliveryServiceMatch `json:"matchList,omitempty"`
	MaxDNSAnswers            int                    `json:"maxDnsAnswers"`
	MidHeaderRewrite         string                 `json:"midHeaderRewrite"`
	MissLat                  float64                `json:"missLat"`
	MissLong                 float64                `json:"missLong"`
	MultiSiteOrigin          bool                   `json:"multiSiteOrigin"`
	OrgServerFQDN            string                 `json:"orgServerFqdn"`
	ProfileDesc              string                 `json:"profileDescription"`
	ProfileID                int                    `json:"profileId,omitempty"`
	ProfileName              string                 `json:"profileName"`
	Protocol                 int                    `json:"protocol"`
	QStringIgnore            int                    `json:"qstringIgnore"`
	RangeRequestHandling     int                    `json:"rangeRequestHandling"`
	RegexRemap               string                 `json:"regexRemap"`
	RegionalGeoBlocking      bool                   `json:"regionalGeoBlocking"`
	RemapText                string                 `json:"remapText"`
	RoutingName              string                 `json:"routingName"`
	Signed                   bool                   `json:"signed"`
	TypeID                   int                    `json:"typeId"`
	Type                     DSType                 `json:"type"`
	TRResponseHeaders        string                 `json:"trResponseHeaders"`
	TenantID                 int                    `json:"tenantId"`
	XMLID                    string                 `json:"xmlId"`
}

// DeliveryServiceV31 represents a Delivery Service as they appear in version
// 3.1 of the Traffic Ops API.
//
// Deprecated: API version 3.1 is deprecated - upgrade to DeliveryServiceV4.
type DeliveryServiceV31 struct {
	DeliveryServiceV30
	DeliveryServiceFieldsV31
}

// DeliveryServiceFieldsV31 contains additions to DeliverySservices introduced
// in API v3.1.
//
// Deprecated: API version 3.1 is deprecated.
type DeliveryServiceFieldsV31 struct {
	// MaxRequestHeaderBytes is the maximum size (in bytes) of the request
	// header that is allowed for this Delivery Service.
	MaxRequestHeaderBytes *int `json:"maxRequestHeaderBytes" db:"max_request_header_bytes"`
}

// DeliveryServiceV40 is a Delivery Service as it appears in version 4.0 of the
// Traffic Ops API.
type DeliveryServiceV40 struct {
	DeliveryServiceFieldsV31
	DeliveryServiceFieldsV30
	DeliveryServiceFieldsV15
	DeliveryServiceFieldsV14
	DeliveryServiceFieldsV13
	DeliveryServiceNullableFieldsV11

	// TLSVersions is the list of explicitly supported TLS versions for cache
	// servers serving the Delivery Service's content.
	TLSVersions       []string              `json:"tlsVersions" db:"tls_versions"`
	GeoLimitCountries GeoLimitCountriesType `json:"geoLimitCountries"`
}

// DeliveryServiceV41 is a Delivery Service as it appears in version 4.1 of the
// Traffic Ops API.
type DeliveryServiceV41 struct {
	DeliveryServiceV40

	// Regional indicates whether the Delivery Service's MaxOriginConnections is
	// only per Cache Group, rather than divided over all Cache Servers in child
	// Cache Groups of the Origin.
	Regional             bool     `json:"regional" db:"regional"`
	RequiredCapabilities []string `json:"requiredCapabilities" db:"required_capabilities"`
}

// DeliveryServiceV4 is a Delivery Service as it appears in version 4 of the
// Traffic Ops API - it always points to the highest minor version in APIv4.
type DeliveryServiceV4 = DeliveryServiceV41

// These are the TLS Versions known by Apache Traffic Control to exist.
const (
	// Deprecated: TLS version 1.0 is known to be insecure.
	TLSVersion10 = "1.0"
	// Deprecated: TLS version 1.1 is known to be insecure.
	TLSVersion11 = "1.1"
	TLSVersion12 = "1.2"
	TLSVersion13 = "1.3"
)

func newerTLSVersionsDisallowedMessage(old string, newer []string) string {
	l := len(newer)
	if l < 1 {
		return ""
	}

	var msg strings.Builder
	msg.WriteString("old TLS version ")
	msg.WriteString(old)
	msg.WriteString(" is allowed, but newer version")
	if l > 1 {
		msg.WriteRune('s')
	}
	msg.WriteRune(' ')
	msg.WriteString(newer[0])
	if l > 1 {
		msg.WriteString(", ")
		if l > 2 {
			msg.WriteString(newer[1])
			msg.WriteString(", and ")
			msg.WriteString(newer[2])
		} else {
			msg.WriteString("and ")
			msg.WriteString(newer[1])
		}
		msg.WriteString(" are ")
	} else {
		msg.WriteString(" is ")
	}
	msg.WriteString("disallowed; this configuration may be insecure")

	return msg.String()
}

func tlsVersionsAlerts(versions []string, protocol int) Alerts {
	messages := []string{}

	if len(versions) > 0 {
		messages = append(messages, "setting TLS Versions that are explicitly supported may break older clients that can't use the specified versions")
	} else {
		return Alerts{Alerts: []Alert{}}
	}

	found := map[string]bool{
		TLSVersion10: false,
		TLSVersion11: false,
		TLSVersion12: false,
		TLSVersion13: false,
	}

	for _, v := range versions {
		switch v {
		case TLSVersion10:
			found[TLSVersion10] = true
		case TLSVersion11:
			found[TLSVersion11] = true
		case TLSVersion12:
			found[TLSVersion12] = true
		case TLSVersion13:
			found[TLSVersion13] = true
		default:
			messages = append(messages, "unknown TLS version '"+v+"' - possible typo")
		}
	}

	if found[TLSVersion10] {
		var newerDisallowed []string
		if !found[TLSVersion11] {
			newerDisallowed = append(newerDisallowed, TLSVersion11)
		}
		if !found[TLSVersion12] {
			newerDisallowed = append(newerDisallowed, TLSVersion12)
		}
		if !found[TLSVersion13] {
			newerDisallowed = append(newerDisallowed, TLSVersion13)
		}
		msg := newerTLSVersionsDisallowedMessage(TLSVersion10, newerDisallowed)
		if msg != "" {
			messages = append(messages, msg)
		}
	} else if found[TLSVersion11] {
		var newerDisallowed []string
		if !found[TLSVersion12] {
			newerDisallowed = append(newerDisallowed, TLSVersion12)
		}
		if !found[TLSVersion13] {
			newerDisallowed = append(newerDisallowed, TLSVersion13)
		}
		msg := newerTLSVersionsDisallowedMessage(TLSVersion11, newerDisallowed)
		if msg != "" {
			messages = append(messages, msg)
		}
	} else if found[TLSVersion12] {
		var newerDisallowed []string
		if !found[TLSVersion13] {
			newerDisallowed = append(newerDisallowed, TLSVersion13)
		}
		msg := newerTLSVersionsDisallowedMessage(TLSVersion12, newerDisallowed)
		if msg != "" {
			messages = append(messages, msg)
		}
	}

	if protocol == DSProtocolHTTP {
		messages = append(messages, "tlsVersions has no effect on Delivery Services with Protocol '0' (HTTP_ONLY)")
	}

	return CreateAlerts(WarnLevel, messages...)
}

// TLSVersionsAlerts generates warning-level alerts for the Delivery Service's
// TLS versions array. It will warn if newer versions are disallowed while
// older, less secure versions are allowed, if there are unrecognized versions
// present, if the Delivery Service's Protocol does not make use of TLS
// Versions, and whenever TLSVersions are explicitly set at all.
//
// This does NOT verify that the Delivery Service's TLS versions are _valid_,
// it ONLY creates warnings based on conditions that are possibly detrimental
// to CDN operation, but can, in fact, work.
func (ds DeliveryServiceV40) TLSVersionsAlerts() Alerts {
	vers := ds.TLSVersions
	return tlsVersionsAlerts(vers, util.Coalesce(ds.Protocol, 3))
}

// TLSVersionsAlerts generates warning-level alerts for the Delivery Service's
// TLS versions array. It will warn if newer versions are disallowed while
// older, less secure versions are allowed, if there are unrecognized versions
// present, if the Delivery Service's Protocol does not make use of TLS
// Versions, and whenever TLSVersions are explicitly set at all.
//
// This does NOT verify that the Delivery Service's TLS versions are _valid_,
// it ONLY creates warnings based on conditions that are possibly detrimental
// to CDN operation, but can, in fact, work.
func (ds DeliveryServiceV41) TLSVersionsAlerts() Alerts {
	return ds.DeliveryServiceV40.TLSVersionsAlerts()
}

// DeliveryServiceV30 represents a Delivery Service as they appear in version
// 3.0 of the Traffic Ops API.
//
// Deprecated: API version 3.0 is deprecated - upgrade to DeliveryServiceV4.
type DeliveryServiceV30 struct {
	DeliveryServiceNullableV15
	DeliveryServiceFieldsV30
}

// DeliveryServiceFieldsV30 contains additions to Delivery Services introduced
// in API v3.0.
//
// Deprecated: API version 3.0 is deprecated - upgrade to DeliveryServiceV4.
type DeliveryServiceFieldsV30 struct {
	// FirstHeaderRewrite is a "header rewrite rule" used by ATS at the first
	// caching layer encountered in the Delivery Service's Topology, or nil if
	// there is no such rule. This has no effect on Delivery Services that don't
	// employ Topologies.
	FirstHeaderRewrite *string `json:"firstHeaderRewrite" db:"first_header_rewrite"`
	// InnerHeaderRewrite is a "header rewrite rule" used by ATS at all caching
	// layers encountered in the Delivery Service's Topology except the first
	// and last, or nil if there is no such rule. This has no effect on Delivery
	// Services that don't employ Topologies.
	InnerHeaderRewrite *string `json:"innerHeaderRewrite" db:"inner_header_rewrite"`
	// LastHeaderRewrite is a "header rewrite rule" used by ATS at the first
	// caching layer encountered in the Delivery Service's Topology, or nil if
	// there is no such rule. This has no effect on Delivery Services that don't
	// employ Topologies.
	LastHeaderRewrite *string `json:"lastHeaderRewrite" db:"last_header_rewrite"`
	// ServiceCategory defines a category to which a Delivery Service may
	// belong, which will cause HTTP Responses containing content for the
	// Delivery Service to have the "X-CDN-SVC" header with a value that is the
	// XMLID of the Delivery Service.
	ServiceCategory *string `json:"serviceCategory" db:"service_category"`
	// Topology is the name of the Topology used by the Delivery Service, or nil
	// if no Topology is used.
	Topology *string `json:"topology" db:"topology"`
}

// DeliveryServiceNullableV30 is the aliased structure that we should be using
// for all API 3.x Delivery Service operations.
//
// Again, this type is an alias that refers to the LATEST MINOR VERSION of API
// version 3 - NOT API version 3.0 as the name might imply.
//
// This type should always alias the latest 3.x minor version struct. For
// example, if you wanted to create a DeliveryServiceV32 struct, you would do
// the following:
//
//	type DeliveryServiceNullableV30 DeliveryServiceV32
//	DeliveryServiceV32 = DeliveryServiceV31 + the new fields.
//
// Deprecated: API version 3 is deprecated - upgrade to DeliveryServiceV4.
type DeliveryServiceNullableV30 DeliveryServiceV31

// DeliveryServiceNullable  represents a Delivery Service as they appeared in
// version 1.5 - and coincidentally also version 2.0 - of the Traffic Ops API.
//
// Deprecated: All API versions for which this could be used to represent
// structures are deprecated - upgrade to DeliveryServiceV4.
type DeliveryServiceNullable DeliveryServiceNullableV15

// DeliveryServiceNullableV15 represents a Delivery Service as they appeared in
// version 1.5 of the Traffic Ops API - which no longer exists.
//
// Because the structure of Delivery Services did not change between Traffic
// Ops API versions 1.5 and 2.0, this is also used in many places to represent
// an APIv2 Delivery Service.
//
// Deprecated: All API versions for which this could be used to represent
// structures are deprecated - upgrade to DeliveryServiceV4.
type DeliveryServiceNullableV15 struct {
	DeliveryServiceNullableV14
	DeliveryServiceFieldsV15
}

// DeliveryServiceFieldsV15 contains additions to Delivery Services introduced
// in Traffic Ops API v1.5.
//
// Deprecated: API version 1.5 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceFieldsV15 struct {
	// EcsEnabled describes whether or not the Traffic Router's EDNS0 Client
	// Subnet extensions should be enabled when serving DNS responses for this
	// Delivery Service. Even if this is true, the Traffic Router may still
	// have the extensions unilaterally disabled in its own configuration.
	EcsEnabled bool `json:"ecsEnabled" db:"ecs_enabled"`
	// RangeSliceBlockSize defines the size of range request blocks - or
	// "slices" - used by the "slice" plugin. This has no effect if
	// RangeRequestHandling does not point to exactly 3. This may never legally
	// point to a value less than zero.
	RangeSliceBlockSize *int `json:"rangeSliceBlockSize" db:"range_slice_block_size"`
}

// DeliveryServiceNullableV14 represents a Delivery Service as they appeared in
// version 1.4 of the Traffic Ops API - which no longer exists.
//
// Deprecated: API version 1.4 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceNullableV14 struct {
	DeliveryServiceNullableV13
	DeliveryServiceFieldsV14
}

// DeliveryServiceFieldsV14 contains additions to Delivery Services introduced
// in Traffic Ops API v1.4.
//
// Deprecated: API version 1.4 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceFieldsV14 struct {
	// ConsistentHashRegex is used by Traffic Router to extract meaningful parts
	// of a client's request URI for HTTP-routed Delivery Services before
	// hashing the request to find a cache server to which to direct the client.
	ConsistentHashRegex *string `json:"consistentHashRegex"`
	// ConsistentHashQueryParams is a list of al of the query string parameters
	// which ought to be considered by Traffic Router in client request URIs for
	// HTTP-routed Delivery Services in the hashing process.
	ConsistentHashQueryParams []string `json:"consistentHashQueryParams"`
	// MaxOriginConnections defines the total maximum  number of connections
	// that the highest caching layer ("Mid-tier" in a non-Topology context) is
	// allowed to have concurrently open to the Delivery Service's Origin. This
	// may never legally point to a value less than 0.
	MaxOriginConnections *int `json:"maxOriginConnections" db:"max_origin_connections"`
}

// DeliveryServiceNullableV13 represents a Delivery Service as they appeared in
// version 1.3 of the Traffic Ops API - which no longer exists.
//
// Deprecated: API version 1.3 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceNullableV13 struct {
	DeliveryServiceNullableV12
	DeliveryServiceFieldsV13
}

// DeliveryServiceFieldsV13 contains additions to Delivery Services introduced
// in Traffic Ops API v1.3.
//
// Deprecated: API version 1.3 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceFieldsV13 struct {
	// DeepCachingType may only legally point to 'ALWAYS' or 'NEVER', which
	// define whether "deep caching" may or may not be used for this Delivery
	// Service, respectively.
	DeepCachingType *DeepCachingType `json:"deepCachingType" db:"deep_caching_type"`
	// FQPacingRate sets the maximum bytes per second a cache server will deliver
	// on any single TCP connection for this Delivery Service. This may never
	// legally point to a value less than zero.
	FQPacingRate *int `json:"fqPacingRate" db:"fq_pacing_rate"`
	// SigningAlgorithm is the name of the algorithm used to sign CDN URIs for
	// this Delivery Service's content, or nil if no URI signing is done for the
	// Delivery Service. This may only point to the values "url_sig" or
	// "uri_signing".
	SigningAlgorithm *string `json:"signingAlgorithm" db:"signing_algorithm"`
	// Tenant is the Tenant to which the Delivery Service belongs.
	Tenant *string `json:"tenant"`
	// TRResponseHeaders is a set of headers (separated by CRLF pairs as per the
	// HTTP spec) and their values (separated by a colon as per the HTTP spec)
	// which will be sent by Traffic Router in HTTP responses to client requests
	// for this Delivery Service's content. This has no effect on DNS-routed or
	// un-routed Delivery Service Types.
	TRResponseHeaders *string `json:"trResponseHeaders"`
	// TRRequestHeaders is an "array" of HTTP headers which should be logged
	// from client HTTP requests for this Delivery Service's content by Traffic
	// Router, separated by newlines. This has no effect on DNS-routed or
	// un-routed Delivery Service Types.
	TRRequestHeaders *string `json:"trRequestHeaders"`
}

// DeliveryServiceNullableV12 represents a Delivery Service as they appeared in
// version 1.2 of the Traffic Ops API - which no longer exists.
//
// Deprecated: API version 1.2 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceNullableV12 struct {
	DeliveryServiceNullableV11
}

// DeliveryServiceNullableV11 represents a Delivery Service as they appeared in
// version 1.1 of the Traffic Ops API - which no longer exists.
//
// Deprecated: API version 1.1 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceNullableV11 struct {
	DeliveryServiceNullableFieldsV11
	DeliveryServiceRemovedFieldsV11
}

// GeoLimitCountriesType is the type alias that is used to represent the GeoLimitCountries attribute of the DeliveryService struct.
type GeoLimitCountriesType []string

// UnmarshalJSON will unmarshal a byte slice into type GeoLimitCountriesType.
func (g *GeoLimitCountriesType) UnmarshalJSON(data []byte) error {
	var err error
	var initial = make([]string, 0)
	var initialStr string
	if err = json.Unmarshal(data, &initial); err != nil {
		if err = json.Unmarshal(data, &initialStr); err != nil {
			return err
		}
		if strings.Contains(initialStr, ",") {
			initial = strings.Split(initialStr, ",")
		} else {
			initial = append(initial, initialStr)
		}
	}

	if initial == nil || len(initial) == 0 {
		g = nil
		return nil
	}
	*g = initial
	return nil

}

// MarshalJSON will marshal a GeoLimitCountriesType into a byte slice.
func (g GeoLimitCountriesType) MarshalJSON() ([]byte, error) {
	arr := ([]string)(g)
	return json.Marshal(arr)
}

// DeliveryServiceNullableFieldsV11 contains properties that Delivery Services
// as they appeared in Traffic Ops API v1.1 had, AND were not removed by ANY
// later API version.
//
// Deprecated: API version 1.1 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceNullableFieldsV11 struct {
	// Active dictates whether the Delivery Service is routed by Traffic Router.
	Active *bool `json:"active" db:"active"`
	// AnonymousBlockingEnabled sets whether or not anonymized IP addresses
	// (e.g. Tor exit nodes) should be restricted from accessing the Delivery
	// Service's content.
	AnonymousBlockingEnabled *bool `json:"anonymousBlockingEnabled" db:"anonymous_blocking_enabled"`
	// CCRDNSTTL sets the Time-to-Live - in seconds - for DNS responses for this
	// Delivery Service from Traffic Router.
	CCRDNSTTL *int `json:"ccrDnsTtl" db:"ccr_dns_ttl"`
	// CDNID is the integral, unique identifier for the CDN to which the
	// Delivery Service belongs.
	CDNID *int `json:"cdnId" db:"cdn_id"`
	// CDNName is the name of the CDN to which the Delivery Service belongs.
	CDNName *string `json:"cdnName"`
	// CheckPath is a path which may be requested of the Delivery Service's
	// origin to ensure it's working properly.
	CheckPath *string `json:"checkPath" db:"check_path"`
	// DisplayName is a human-friendly name that might be used in some UIs
	// somewhere.
	DisplayName *string `json:"displayName" db:"display_name"`
	// DNSBypassCNAME is a fully qualified domain name to be used in a CNAME
	// record presented to clients in bypass scenarios.
	DNSBypassCNAME *string `json:"dnsBypassCname" db:"dns_bypass_cname"`
	// DNSBypassIP is an IPv4 address to be used in an A record presented to
	// clients in bypass scenarios.
	DNSBypassIP *string `json:"dnsBypassIp" db:"dns_bypass_ip"`
	// DNSBypassIP6 is an IPv6 address to be used in an AAAA record presented to
	// clients in bypass scenarios.
	DNSBypassIP6 *string `json:"dnsBypassIp6" db:"dns_bypass_ip6"`
	// DNSBypassTTL sets the Time-to-Live - in seconds - of DNS responses from
	// the Traffic Router that contain records for bypass destinations.
	DNSBypassTTL *int `json:"dnsBypassTtl" db:"dns_bypass_ttl"`
	// DSCP sets the Differentiated Services Code Point for IP packets
	// transferred between clients, origins, and cache servers when obtaining
	// and serving content for this Delivery Service.
	// See Also: https://en.wikipedia.org/wiki/Differentiated_services
	DSCP *int `json:"dscp" db:"dscp"`
	// EdgeHeaderRewrite is a "header rewrite rule" used by ATS at the Edge-tier
	// of caching. This has no effect on Delivery Services that don't use a
	// Topology.
	EdgeHeaderRewrite *string `json:"edgeHeaderRewrite" db:"edge_header_rewrite"`
	// ExampleURLs is a list of all of the URLs from which content may be
	// requested from the Delivery Service.
	ExampleURLs []string `json:"exampleURLs"`
	// GeoLimit defines whether or not access to a Delivery Service's content
	// should be limited based on the requesting client's geographic location.
	// Despite that this is a pointer to an arbitrary integer, the only valid
	// values are 0 (which indicates that content should not be limited
	// geographically), 1 (which indicates that content should only be served to
	// clients whose IP addresses can be found within a Coverage Zone File), and
	// 2 (which indicates that content should be served to clients whose IP
	// addresses can be found within a Coverage Zone File OR are allowed access
	// according to the "array" in GeoLimitCountries).
	GeoLimit *int `json:"geoLimit" db:"geo_limit"`
	// GeoLimitCountries is an "array" of "country codes" that itemizes the
	// countries within which the Delivery Service's content ought to be made
	// available. This has no effect if GeoLimit is not a pointer to exactly the
	// value 2.
	GeoLimitCountries *string `json:"geoLimitCountries" db:"geo_limit_countries"`
	// GeoLimitRedirectURL is a URL to which clients will be redirected if their
	// access to the Delivery Service's content is blocked by GeoLimit rules.
	GeoLimitRedirectURL *string `json:"geoLimitRedirectURL" db:"geolimit_redirect_url"`
	// GeoProvider names the type of database to be used for providing IP
	// address-to-geographic-location mapping for this Delivery Service. The
	// only valid values to which it may point are 0 (which indicates the use of
	// a MaxMind GeoIP2 database) and 1 (which indicates the use of a Neustar
	// GeoPoint IP address database).
	GeoProvider *int `json:"geoProvider" db:"geo_provider"`
	// GlobalMaxMBPS defines a maximum number of MegaBytes Per Second which may
	// be served for the Delivery Service before redirecting clients to bypass
	// locations.
	GlobalMaxMBPS *int `json:"globalMaxMbps" db:"global_max_mbps"`
	// GlobalMaxTPS defines a maximum number of Transactions Per Second which
	// may be served for the Delivery Service before redirecting clients to
	// bypass locations.
	GlobalMaxTPS *int `json:"globalMaxTps" db:"global_max_tps"`
	// HTTPBypassFQDN is a network location to which clients will be redirected
	// in bypass scenarios using HTTP "Location" headers and appropriate
	// redirection response codes.
	HTTPBypassFQDN *string `json:"httpBypassFqdn" db:"http_bypass_fqdn"`
	// ID is an integral, unique identifier for the Delivery Service.
	ID *int `json:"id" db:"id"`
	// InfoURL is a URL to which operators or clients may be directed to obtain
	// further information about a Delivery Service.
	InfoURL *string `json:"infoUrl" db:"info_url"`
	// InitialDispersion sets the number of cache servers within the first
	// caching layer ("Edge-tier" in a non-Topology context) across which
	// content will be dispersed per Cache Group.
	InitialDispersion *int `json:"initialDispersion" db:"initial_dispersion"`
	// IPV6RoutingEnabled controls whether or not routing over IPv6 should be
	// done for this Delivery Service.
	IPV6RoutingEnabled *bool `json:"ipv6RoutingEnabled" db:"ipv6_routing_enabled"`
	// LastUpdated is the time and date at which the Delivery Service was last
	// updated.
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	// LogsEnabled controls nothing. It is kept only for legacy compatibility.
	LogsEnabled *bool `json:"logsEnabled" db:"logs_enabled"`
	// LongDesc is a description of the Delivery Service, having arbitrary
	// length.
	LongDesc *string `json:"longDesc" db:"long_desc"`
	// LongDesc1 is a description of the Delivery Service, having arbitrary
	// length.
	LongDesc1 *string `json:"longDesc1,omitempty" db:"long_desc_1"`
	// LongDesc2 is a description of the Delivery Service, having arbitrary
	// length.
	LongDesc2 *string `json:"longDesc2,omitempty" db:"long_desc_2"`
	// MatchList is a list of Regular Expressions used for routing the Delivery
	// Service. Order matters, and the array is not allowed to be sparse.
	MatchList *[]DeliveryServiceMatch `json:"matchList"`
	// MaxDNSAnswers sets the maximum number of records which should be returned
	// by Traffic Router in DNS responses to requests for resolving names for
	// this Delivery Service.
	MaxDNSAnswers *int `json:"maxDnsAnswers" db:"max_dns_answers"`
	// MidHeaderRewrite is a "header rewrite rule" used by ATS at the Mid-tier
	// of caching. This has no effect on Delivery Services that don't use a
	// Topology.
	MidHeaderRewrite *string `json:"midHeaderRewrite" db:"mid_header_rewrite"`
	// MissLat is a latitude to default to for clients of this Delivery Service
	// when geolocation attempts fail.
	MissLat *float64 `json:"missLat" db:"miss_lat"`
	// MissLong is a longitude to default to for clients of this Delivery
	// Service when geolocation attempts fail.
	MissLong *float64 `json:"missLong" db:"miss_long"`
	// MultiSiteOrigin determines whether or not the Delivery Service makes use
	// of "Multi-Site Origin".
	MultiSiteOrigin *bool `json:"multiSiteOrigin" db:"multi_site_origin"`
	// OriginShield is a field that does nothing. It is kept only for legacy
	// compatibility reasons.
	OriginShield *string `json:"originShield" db:"origin_shield"`
	// OrgServerFQDN is the URL - NOT Fully Qualified Domain Name - of the
	// origin of the Delivery Service's content.
	OrgServerFQDN *string `json:"orgServerFqdn" db:"org_server_fqdn"`
	// ProfileDesc is the Description of the Profile used by the Delivery
	// Service, if any.
	ProfileDesc *string `json:"profileDescription"`
	// ProfileID is the integral, unique identifier of the Profile used by the
	// Delivery Service, if any.
	ProfileID *int `json:"profileId" db:"profile"`
	// ProfileName is the Name of the Profile used by the Delivery Service, if
	// any.
	ProfileName *string `json:"profileName"`
	// Protocol defines the protocols by which caching servers may communicate
	// with clients. The valid values to which it may point are 0 (which implies
	// that only HTTP may be used), 1 (which implies that only HTTPS may be
	// used), 2 (which implies that either HTTP or HTTPS may be used), and 3
	// (which implies that clients using HTTP must be redirected to use HTTPS,
	// while communications over HTTPS may proceed as normal).
	Protocol *int `json:"protocol" db:"protocol"`
	// QStringIgnore sets how query strings in HTTP requests to cache servers
	// from clients are treated. The only valid values to which it may point are
	// 0 (which implies that all caching layers will pass query strings in
	// upstream requests and use them in the cache key), 1 (which implies that
	// all caching layers will pass query strings in upstream requests, but not
	// use them in cache keys), and 2 (which implies that the first encountered
	// caching layer - "Edge-tier" in a non-Topology context - will strip query
	// strings, effectively preventing them from being passed in upstream
	// requests, and not use them in the cache key).
	QStringIgnore *int `json:"qstringIgnore" db:"qstring_ignore"`
	// RangeRequestHandling defines how HTTP GET requests with a Range header
	// will be handled by cache servers serving the Delivery Service's content.
	// The only valid values to which it may point are 0 (which implies that
	// Range requests will not be cached at all), 1 (which implies that the
	// background_fetch plugin will be used to service the range request while
	// caching the whole object), 2 (which implies that the cache_range_requests
	// plugin will be used to cache ranges as unique objects), and 3 (which
	// implies that the slice plugin will be used to slice range based requests
	// into deterministic chunks.)
	RangeRequestHandling *int `json:"rangeRequestHandling" db:"range_request_handling"`
	// Regex Remap is a raw line to be inserted into "regex_remap.config" on the
	// cache server. Care is necessitated in its use, because the input is in no
	// way restricted, validated, or limited in scope to the Delivery Service.
	RegexRemap *string `json:"regexRemap" db:"regex_remap"`
	// RegionalGeoBlocking defines whether or not whatever Regional Geo Blocking
	// rules are configured on the Traffic Router serving content for this
	// Delivery Service will have an effect on the traffic of this Delivery
	// Service.
	RegionalGeoBlocking *bool `json:"regionalGeoBlocking" db:"regional_geo_blocking"`
	// RemapText is raw text to insert in "remap.config" on the cache servers
	// serving content for this Delivery Service. Care is necessitated in its
	// use, because the input is in no way restricted, validated, or limited in
	// scope to the Delivery Service.
	RemapText *string `json:"remapText" db:"remap_text"`
	// RoutingName defines the lowest-level DNS label used by the Delivery
	// Service, e.g. if trafficcontrol.apache.org were a Delivery Service, it
	// would have a RoutingName of "trafficcontrol".
	RoutingName *string `json:"routingName" db:"routing_name"`
	// Signed is a legacy field. It is allowed to be `true` if and only if
	// SigningAlgorithm is not nil.
	Signed bool `json:"signed"`
	// SSLKeyVersion incremented whenever Traffic Portal generates new SSL keys
	// for the Delivery Service, effectively making it a "generational" marker.
	SSLKeyVersion *int `json:"sslKeyVersion" db:"ssl_key_version"`
	// TenantID is the integral, unique identifier for the Tenant to which the
	// Delivery Service belongs.
	TenantID *int `json:"tenantId" db:"tenant_id"`
	// Type describes how content is routed and cached for this Delivery Service
	// as well as what other properties have any meaning.
	Type *DSType `json:"type"`
	// TypeID is an integral, unique identifier for the Tenant to which the
	// Delivery Service belongs.
	TypeID *int `json:"typeId" db:"type"`
	// XMLID is a unique identifier that is also the second lowest-level DNS
	// label used by the Delivery Service. For example, if a Delivery Service's
	// content may be requested from video.demo1.mycdn.ciab.test, it may be
	// inferred that the Delivery Service's XMLID is demo1.
	XMLID *string `json:"xmlId" db:"xml_id"`
}

// DeliveryServiceRemovedFieldsV11 contains properties of Delivery Services as
// they appeared in version 1.1 of the Traffic Ops API that were later removed
// in some other API version.
//
// Deprecated: API version 1.1 no longer exists, this type ONLY still exists
// because newer structures nest it, so removing it would be a breaking change
// - please upgrade to DeliveryServiceV4.
type DeliveryServiceRemovedFieldsV11 struct {
	CacheURL *string `json:"cacheurl" db:"cacheurl"`
}

// RemoveLD1AndLD2 removes the Long Description 1 and Long Description 2 fields
// from a V4.0 Delivery Service, and returns the resulting struct.
func (ds *DeliveryServiceV40) RemoveLD1AndLD2() DeliveryServiceV40 {
	ds.LongDesc1 = nil
	ds.LongDesc2 = nil
	return *ds
}

// RemoveLD1AndLD2 removes the Long Description 1 and Long Description 2 fields
// from a V 4.x DS, and returns the resulting struct.
func (ds *DeliveryServiceV4) RemoveLD1AndLD2() DeliveryServiceV4 {
	ds.LongDesc1 = nil
	ds.LongDesc2 = nil
	return *ds
}

// DowngradeToV31 converts a 4.x DS to a 3.1 DS.
func (ds DeliveryServiceV4) DowngradeToV31() DeliveryServiceNullableV30 {
	nullableFields := ds.DeliveryServiceNullableFieldsV11
	geoLimitCountries := ([]string)(ds.GeoLimitCountries)
	geo := strings.Join(geoLimitCountries, ",")
	nullableFields.GeoLimitCountries = &geo
	return DeliveryServiceNullableV30{
		DeliveryServiceV30: DeliveryServiceV30{
			DeliveryServiceNullableV15: DeliveryServiceNullableV15{
				DeliveryServiceNullableV14: DeliveryServiceNullableV14{
					DeliveryServiceNullableV13: DeliveryServiceNullableV13{
						DeliveryServiceNullableV12: DeliveryServiceNullableV12{
							DeliveryServiceNullableV11: DeliveryServiceNullableV11{
								DeliveryServiceNullableFieldsV11: nullableFields,
							},
						},
						DeliveryServiceFieldsV13: ds.DeliveryServiceFieldsV13,
					},
					DeliveryServiceFieldsV14: ds.DeliveryServiceFieldsV14,
				},
				DeliveryServiceFieldsV15: ds.DeliveryServiceFieldsV15,
			},
			DeliveryServiceFieldsV30: ds.DeliveryServiceFieldsV30,
		},
		DeliveryServiceFieldsV31: ds.DeliveryServiceFieldsV31,
	}
}

// UpgradeToV4 converts the 3.x DS to a 4.x DS.
func (ds DeliveryServiceNullableV30) UpgradeToV4() DeliveryServiceV4 {
	var geo GeoLimitCountriesType
	if ds.GeoLimitCountries != nil {
		str := *ds.GeoLimitCountries
		geo = strings.Split(str, ",")
	}
	return DeliveryServiceV4{
		DeliveryServiceV40: DeliveryServiceV40{
			DeliveryServiceFieldsV31:         ds.DeliveryServiceFieldsV31,
			DeliveryServiceFieldsV30:         ds.DeliveryServiceFieldsV30,
			DeliveryServiceFieldsV15:         ds.DeliveryServiceFieldsV15,
			DeliveryServiceFieldsV14:         ds.DeliveryServiceFieldsV14,
			DeliveryServiceFieldsV13:         ds.DeliveryServiceFieldsV13,
			DeliveryServiceNullableFieldsV11: ds.DeliveryServiceNullableFieldsV11,
			TLSVersions:                      nil,
			GeoLimitCountries:                geo,
		},
	}
}

func jsonScan(src interface{}, dest interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected deliveryservice in byte array form; got %T", src)
	}
	return json.Unmarshal(b, dest)
}

// NOTE: the driver.Valuer and sql.Scanner interface implementations are
// necessary for Delivery Service Requests which store and read raw JSON
// from the database.

// Value implements the database/sql/driver.Valuer interface by marshaling the
// struct to JSON to pass back as an encoding/json.RawMessage.
func (ds *DeliveryServiceNullable) Value() (driver.Value, error) {
	return json.Marshal(ds)
}

// Scan implements the database/sql.Scanner interface.
//
// This expects src to be an encoding/json.RawMessage and unmarshals that into
// the DeliveryServiceNullable.
func (ds *DeliveryServiceNullable) Scan(src interface{}) error {
	return jsonScan(src, ds)
}

// Value implements the database/sql/driver.Valuer interface by marshaling the
// struct to JSON to pass back as an encoding/json.RawMessage.
func (ds *DeliveryServiceV4) Value() (driver.Value, error) {
	return json.Marshal(ds)
}

// Scan implements the database/sql.Scanner interface.
//
// This expects src to be an encoding/json.RawMessage and unmarshals that into
// the DeliveryServiceV4.
func (ds *DeliveryServiceV4) Scan(src interface{}) error {
	return jsonScan(src, ds)
}

// DeliveryServiceMatch structures are the type of each entry in a Delivery
// Service's Match List.
type DeliveryServiceMatch struct {
	Type      DSMatchType `json:"type"`
	SetNumber int         `json:"setNumber"`
	Pattern   string      `json:"pattern"`
}

// DeliveryServiceHealthResponse is the type of a response from Traffic Ops to
// a request for a Delivery Service's "health".
type DeliveryServiceHealthResponse struct {
	Response DeliveryServiceHealth `json:"response"`
	Alerts
}

// DeliveryServiceHealth represents the "health" of a Delivery Service by the
// number of cache servers responsible for serving its content that are
// determined to be "online"/"healthy" and "offline"/"unhealthy".
type DeliveryServiceHealth struct {
	TotalOnline  int                         `json:"totalOnline"`
	TotalOffline int                         `json:"totalOffline"`
	CacheGroups  []DeliveryServiceCacheGroup `json:"cacheGroups"`
}

// DeliveryServiceCacheGroup breaks down the "health" of a Delivery Service by
// the number of cache servers responsible for serving its content within a
// specific Cache Group that are determined to be "online"/"healthy" and
// "offline"/"unhealthy".
type DeliveryServiceCacheGroup struct {
	Online  int `json:"online"`
	Offline int `json:"offline"`
	// The name of the Cache Group represented by this data.
	Name string `json:"name"`
}

// DeliveryServiceCapacityResponse is the type of a response from Traffic Ops to
// a request for a Delivery Service's "capacity".
type DeliveryServiceCapacityResponse struct {
	Response DeliveryServiceCapacity `json:"response"`
	Alerts
}

// DeliveryServiceCapacity represents the "capacity" of a Delivery Service as
// the portions of the pool of cache servers responsible for serving its
// content that are available for servicing client requests.
type DeliveryServiceCapacity struct {
	// The portion of cache servers that are ready, willing, and able to
	// service client requests.
	AvailablePercent float64 `json:"availablePercent"`
	// The portion of cache servers that are read and willing, but not able to
	// service client requests, generally because Traffic Monitor deems them
	// "unhealthy".
	UnavailablePercent float64 `json:"unavailablePercent"`
	// The portion of cache servers that are actively involved in the flow of
	// Delivery Service content.
	UtilizedPercent float64 `json:"utilizedPercent"`
	// The portion of cache servers that are not yet ready to service client
	// requests because they are undergoing maintenance.
	MaintenancePercent float64 `json:"maintenancePercent"`
}

// A FederationDeliveryServiceNullable is an association between a Federation
// and a Delivery Service.
type FederationDeliveryServiceNullable struct {
	ID    *int    `json:"id" db:"id"`
	CDN   *string `json:"cdn" db:"cdn"`
	Type  *string `json:"type" db:"type"`
	XMLID *string `json:"xmlId" db:"xml_id"`
}

// FederationDeliveryServicesResponse is the type of a response from Traffic
// Ops to a request made to its /federations/{{ID}}/deliveryservices endpoint.
type FederationDeliveryServicesResponse struct {
	Response []FederationDeliveryServiceNullable `json:"response"`
	Alerts
}

// DeliveryServiceUserPost represents a legacy concept that no longer exists in
// Apache Traffic Control.
//
// DO NOT USE THIS - it ONLY still exists because it is still in use by Traffic
// Ops's Go client for API versions 2 and 3, despite that those API versions do
// not include the concepts and functionality for which this structure was
// created.
//
// Deprecated: All Go clients for API versions that still erroneously link to
// this symbol are deprecated, and this structure serves no known purpose.
type DeliveryServiceUserPost struct {
	UserID           *int   `json:"userId"`
	DeliveryServices *[]int `json:"deliveryServices"`
	Replace          *bool  `json:"replace"`
}

// UserDeliveryServicePostResponse represents a legacy concept that no longer
// exists in Apache Traffic Control.
//
// DO NOT USE THIS - it ONLY still exists because it is still in use by Traffic
// Ops's Go client for API versions 2 and 3, despite that those API versions do
// not include the concepts and functionality for which this structure was
// created.
//
// Deprecated: All Go clients for API versions that still erroneously link to
// this symbol are deprecated, and this structure serves no known purpose.
type UserDeliveryServicePostResponse struct {
	Alerts   []Alert                 `json:"alerts"`
	Response DeliveryServiceUserPost `json:"response"`
}

// A DSServerIDs is a description of relationships between a Delivery Service
// and zero or more servers, as well as how that relationship may have been
// recently modified.
type DSServerIDs struct {
	DeliveryServiceID *int  `json:"dsId" db:"deliveryservice"`
	ServerIDs         []int `json:"servers"`
	Replace           *bool `json:"replace"`
}

// DeliveryserviceserverResponse - not to be confused with DSServerResponseV40
// or DSServerResponse- is the type of a response from Traffic Ops to a request
// to the /deliveryserviceserver endpoint to associate servers with a Delivery
// Service.
type DeliveryserviceserverResponse struct {
	Response DSServerIDs `json:"response"`
	Alerts
}

// A CachegroupPostDSReq is a request to associate some Cache Group with a set
// of zero or more Delivery Services.
type CachegroupPostDSReq struct {
	DeliveryServices []int `json:"deliveryServices"`
}

// CacheGroupPostDSResp is the type of the `response` property of a response
// from Traffic Ops to a POST request made to its
// /cachegroups/{{ID}}/deliveryservices API endpoint.
type CacheGroupPostDSResp struct {
	ID               util.JSONIntStr `json:"id"`
	ServerNames      []CacheName     `json:"serverNames"`
	DeliveryServices []int           `json:"deliveryServices"`
}

// CacheGroupPostDSRespResponse is the type of a response from Traffic Ops to a
// POST request made to its /cachegroups/{{ID}}/deliveryservices API endpoint.
type CacheGroupPostDSRespResponse struct {
	Alerts
	Response CacheGroupPostDSResp `json:"response"`
}

// AssignedDsResponse is the type of the `response` property of a response from
// Traffic Ops to a POST request made to its /servers/{{ID}}/deliveryservices
// API endpoint.
type AssignedDsResponse struct {
	ServerID int   `json:"serverId"`
	DSIds    []int `json:"dsIds"`
	Replace  bool  `json:"replace"`
}

// DeliveryServiceSafeUpdateRequest represents a request to update the "safe" fields of a
// Delivery Service.
type DeliveryServiceSafeUpdateRequest struct {
	DisplayName *string `json:"displayName"`
	InfoURL     *string `json:"infoUrl"`
	LongDesc    *string `json:"longDesc"`
	LongDesc1   *string `json:"longDesc1,omitempty"`
}

// Validate implements the github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidator
// interface.
func (r *DeliveryServiceSafeUpdateRequest) Validate(*sql.Tx) error {
	if r.DisplayName == nil {
		return errors.New("displayName: cannot be null/missing")
	}
	return nil
}

// DeliveryServiceSafeUpdateResponse represents Traffic Ops's response to a PUT
// request to its /deliveryservices/{{ID}}/safe endpoint.
// Deprecated: Please only use versioned structures.
type DeliveryServiceSafeUpdateResponse struct {
	Alerts
	// Response contains the representation of the Delivery Service after it has been updated.
	Response []DeliveryServiceNullable `json:"response"`
}

// DeliveryServiceSafeUpdateResponseV30 represents Traffic Ops's response to a PUT
// request to its /api/3.0/deliveryservices/{{ID}}/safe endpoint.
type DeliveryServiceSafeUpdateResponseV30 struct {
	Alerts
	// Response contains the representation of the Delivery Service after it has
	// been updated.
	Response []DeliveryServiceNullableV30 `json:"response"`
}

// DeliveryServiceSafeUpdateResponseV40 represents Traffic Ops's response to a PUT
// request to its /api/4.0/deliveryservices/{{ID}}/safe endpoint.
type DeliveryServiceSafeUpdateResponseV40 struct {
	Alerts
	// Response contains the representation of the Delivery Service after it has
	// been updated.
	Response []DeliveryServiceV40 `json:"response"`
}

// DeliveryServiceSafeUpdateResponseV4 represents TrafficOps's response to a
// PUT request to its /api/4.x/deliveryservices/{{ID}}/safe endpoint.
// This is always a type alias for the structure of a response in the latest
// minor APIv4 version.
type DeliveryServiceSafeUpdateResponseV4 = DeliveryServiceSafeUpdateResponseV40

// DeliveryServiceV50 is a Delivery Service as it appears in version 5.0 of the
// Traffic Ops API.
type DeliveryServiceV50 struct {
	// Active dictates whether the Delivery Service is routed by Traffic Router,
	// and whether cache servers have its configuration downloaded.
	Active DeliveryServiceActiveState `json:"active" db:"active"`
	// AnonymousBlockingEnabled sets whether or not anonymized IP addresses
	// (e.g. Tor exit nodes) should be restricted from accessing the Delivery
	// Service's content.
	AnonymousBlockingEnabled bool `json:"anonymousBlockingEnabled" db:"anonymous_blocking_enabled"`
	// CCRDNSTTL sets the Time-to-Live - in seconds - for DNS responses for this
	// Delivery Service from Traffic Router.
	CCRDNSTTL *int `json:"ccrDnsTtl" db:"ccr_dns_ttl"`
	// CDNID is the integral, unique identifier for the CDN to which the
	// Delivery Service belongs.
	CDNID int `json:"cdnId" db:"cdn_id"`
	// CDNName is the name of the CDN to which the Delivery Service belongs.
	CDNName *string `json:"cdnName"`
	// CheckPath is a path which may be requested of the Delivery Service's
	// origin to ensure it's working properly.
	CheckPath *string `json:"checkPath" db:"check_path"`
	// ConsistentHashQueryParams is a list of al of the query string parameters
	// which ought to be considered by Traffic Router in client request URIs for
	// HTTP-routed Delivery Services in the hashing process.
	ConsistentHashQueryParams []string `json:"consistentHashQueryParams"`
	// ConsistentHashRegex is used by Traffic Router to extract meaningful parts
	// of a client's request URI for HTTP-routed Delivery Services before
	// hashing the request to find a cache server to which to direct the client.
	ConsistentHashRegex *string `json:"consistentHashRegex"`
	// DeepCachingType may only be 'ALWAYS' or 'NEVER', which
	// define whether "deep caching" may or may not be used for this Delivery
	// Service, respectively.
	DeepCachingType DeepCachingType `json:"deepCachingType" db:"deep_caching_type"`
	// DisplayName is a human-friendly name that might be used in some UIs
	// somewhere.
	DisplayName string `json:"displayName" db:"display_name"`
	// DNSBypassCNAME is a fully qualified domain name to be used in a CNAME
	// record presented to clients in bypass scenarios.
	DNSBypassCNAME *string `json:"dnsBypassCname" db:"dns_bypass_cname"`
	// DNSBypassIP is an IPv4 address to be used in an A record presented to
	// clients in bypass scenarios.
	DNSBypassIP *string `json:"dnsBypassIp" db:"dns_bypass_ip"`
	// DNSBypassIP6 is an IPv6 address to be used in an AAAA record presented to
	// clients in bypass scenarios.
	DNSBypassIP6 *string `json:"dnsBypassIp6" db:"dns_bypass_ip6"`
	// DNSBypassTTL sets the Time-to-Live - in seconds - of DNS responses from
	// the Traffic Router that contain records for bypass destinations.
	DNSBypassTTL *int `json:"dnsBypassTtl" db:"dns_bypass_ttl"`
	// DSCP sets the Differentiated Services Code Point for IP packets
	// transferred between clients, origins, and cache servers when obtaining
	// and serving content for this Delivery Service.
	// See Also: https://en.wikipedia.org/wiki/Differentiated_services
	DSCP int `json:"dscp" db:"dscp"`
	// EcsEnabled describes whether or not the Traffic Router's EDNS0 Client
	// Subnet extensions should be enabled when serving DNS responses for this
	// Delivery Service. Even if this is true, the Traffic Router may still
	// have the extensions unilaterally disabled in its own configuration.
	EcsEnabled bool `json:"ecsEnabled" db:"ecs_enabled"`
	// EdgeHeaderRewrite is a "header rewrite rule" used by ATS at the Edge-tier
	// of caching. This has no effect on Delivery Services that don't use a
	// Topology.
	EdgeHeaderRewrite *string `json:"edgeHeaderRewrite" db:"edge_header_rewrite"`
	// ExampleURLs is a list of all of the URLs from which content may be
	// requested from the Delivery Service.
	ExampleURLs []string `json:"exampleURLs"`
	// FirstHeaderRewrite is a "header rewrite rule" used by ATS at the first
	// caching layer encountered in the Delivery Service's Topology, or nil if
	// there is no such rule. This has no effect on Delivery Services that don't
	// employ Topologies.
	FirstHeaderRewrite *string `json:"firstHeaderRewrite" db:"first_header_rewrite"`
	// FQPacingRate sets the maximum bytes per second a cache server will deliver
	// on any single TCP connection for this Delivery Service. This may never
	// legally point to a value less than zero.
	FQPacingRate *int `json:"fqPacingRate" db:"fq_pacing_rate"`
	// GeoLimit defines whether or not access to a Delivery Service's content
	// should be limited based on the requesting client's geographic location.
	// The only valid values are 0 (which indicates that content should not be
	// limited geographically), 1 (which indicates that content should only be
	// served to clients whose IP addresses can be found within a Coverage Zone
	// File), and 2 (which indicates that content should be served to clients
	// whose IP addresses can be found within a Coverage Zone File OR are
	// allowed access according to the array in GeoLimitCountries).
	GeoLimit int `json:"geoLimit" db:"geo_limit"`
	// GeoLimitCountries is an "array" of "country codes" that itemizes the
	// countries within which the Delivery Service's content ought to be made
	// available. This has no effect if GeoLimit is not a pointer to exactly the
	// value 2.
	GeoLimitCountries []string `json:"geoLimitCountries"`
	// GeoLimitRedirectURL is a URL to which clients will be redirected if their
	// access to the Delivery Service's content is blocked by GeoLimit rules.
	GeoLimitRedirectURL *string `json:"geoLimitRedirectURL" db:"geolimit_redirect_url"`
	// GeoProvider names the type of database to be used for providing IP
	// address-to-geographic-location mapping for this Delivery Service. The
	// only valid values are 0 (which indicates the use of a MaxMind GeoIP2
	// database) and 1 (which indicates the use of a Neustar GeoPoint IP address
	// database).
	GeoProvider int `json:"geoProvider" db:"geo_provider"`
	// GlobalMaxMBPS defines a maximum number of MegaBytes Per Second which may
	// be served for the Delivery Service before redirecting clients to bypass
	// locations.
	GlobalMaxMBPS *int `json:"globalMaxMbps" db:"global_max_mbps"`
	// GlobalMaxTPS defines a maximum number of Transactions Per Second which
	// may be served for the Delivery Service before redirecting clients to
	// bypass locations.
	GlobalMaxTPS *int `json:"globalMaxTps" db:"global_max_tps"`
	// HTTPBypassFQDN is a network location to which clients will be redirected
	// in bypass scenarios using HTTP "Location" headers and appropriate
	// redirection response codes.
	HTTPBypassFQDN *string `json:"httpBypassFqdn" db:"http_bypass_fqdn"`
	// ID is an integral, unique identifier for the Delivery Service.
	ID *int `json:"id" db:"id"`
	// InfoURL is a URL to which operators or clients may be directed to obtain
	// further information about a Delivery Service.
	InfoURL *string `json:"infoUrl" db:"info_url"`
	// InitialDispersion sets the number of cache servers within the first
	// caching layer ("Edge-tier" in a non-Topology context) across which
	// content will be dispersed per Cache Group.
	InitialDispersion *int `json:"initialDispersion" db:"initial_dispersion"`
	// InnerHeaderRewrite is a "header rewrite rule" used by ATS at all caching
	// layers encountered in the Delivery Service's Topology except the first
	// and last, or nil if there is no such rule. This has no effect on Delivery
	// Services that don't employ Topologies.
	InnerHeaderRewrite *string `json:"innerHeaderRewrite" db:"inner_header_rewrite"`
	// IPV6RoutingEnabled controls whether or not routing over IPv6 should be
	// done for this Delivery Service.
	IPV6RoutingEnabled *bool `json:"ipv6RoutingEnabled" db:"ipv6_routing_enabled"`
	// LastHeaderRewrite is a "header rewrite rule" used by ATS at the first
	// caching layer encountered in the Delivery Service's Topology, or nil if
	// there is no such rule. This has no effect on Delivery Services that don't
	// employ Topologies.
	LastHeaderRewrite *string `json:"lastHeaderRewrite" db:"last_header_rewrite"`
	// LastUpdated is the time and date at which the Delivery Service was last
	// updated.
	LastUpdated time.Time `json:"lastUpdated" db:"last_updated"`
	// LogsEnabled controls nothing. It is kept only for legacy compatibility.
	LogsEnabled bool `json:"logsEnabled" db:"logs_enabled"`
	// LongDesc is a description of the Delivery Service, having arbitrary
	// length.
	LongDesc string `json:"longDesc" db:"long_desc"`
	// MatchList is a list of Regular Expressions used for routing the Delivery
	// Service. Order matters, and the array is not allowed to be sparse.
	MatchList []DeliveryServiceMatch `json:"matchList"`
	// MaxDNSAnswers sets the maximum number of records which should be returned
	// by Traffic Router in DNS responses to requests for resolving names for
	// this Delivery Service.
	MaxDNSAnswers *int `json:"maxDnsAnswers" db:"max_dns_answers"`
	// MaxOriginConnections defines the total maximum  number of connections
	// that the highest caching layer ("Mid-tier" in a non-Topology context) is
	// allowed to have concurrently open to the Delivery Service's Origin. This
	// may never legally point to a value less than 0.
	MaxOriginConnections *int `json:"maxOriginConnections" db:"max_origin_connections"`
	// MaxRequestHeaderBytes is the maximum size (in bytes) of the request
	// header that is allowed for this Delivery Service.
	MaxRequestHeaderBytes *int `json:"maxRequestHeaderBytes" db:"max_request_header_bytes"`
	// MidHeaderRewrite is a "header rewrite rule" used by ATS at the Mid-tier
	// of caching. This has no effect on Delivery Services that don't use a
	// Topology.
	MidHeaderRewrite *string `json:"midHeaderRewrite" db:"mid_header_rewrite"`
	// MissLat is a latitude to default to for clients of this Delivery Service
	// when geolocation attempts fail.
	MissLat *float64 `json:"missLat" db:"miss_lat"`
	// MissLong is a longitude to default to for clients of this Delivery
	// Service when geolocation attempts fail.
	MissLong *float64 `json:"missLong" db:"miss_long"`
	// MultiSiteOrigin determines whether or not the Delivery Service makes use
	// of "Multi-Site Origin".
	MultiSiteOrigin bool `json:"multiSiteOrigin" db:"multi_site_origin"`
	// OriginShield is a field that does nothing. It is kept only for legacy
	// compatibility reasons.
	OriginShield *string `json:"originShield" db:"origin_shield"`
	// OrgServerFQDN is the URL - NOT Fully Qualified Domain Name - of the
	// origin of the Delivery Service's content.
	OrgServerFQDN *string `json:"orgServerFqdn" db:"org_server_fqdn"`
	// ProfileDesc is the Description of the Profile used by the Delivery
	// Service, if any.
	ProfileDesc *string `json:"profileDescription"`
	// ProfileID is the integral, unique identifier of the Profile used by the
	// Delivery Service, if any.
	ProfileID *int `json:"profileId" db:"profile"`
	// ProfileName is the Name of the Profile used by the Delivery Service, if
	// any.
	ProfileName *string `json:"profileName"`
	// Protocol defines the protocols by which caching servers may communicate
	// with clients. The valid values are 0 (which implies that only HTTP may be
	// used), 1 (which implies that only HTTPS may be used), 2 (which implies
	// that either HTTP or HTTPS may be used), and 3 (which implies that clients
	// using HTTP must be redirected to use HTTPS, while communications over
	// HTTPS may proceed as normal).
	Protocol *int `json:"protocol" db:"protocol"`
	// QStringIgnore sets how query strings in HTTP requests to cache servers
	// from clients are treated. The only valid values are 0 (which implies that
	// all caching layers will pass query strings in upstream requests and use
	// them in the cache key), 1 (which implies that all caching layers will
	// pass query strings in upstream requests, but not use them in cache keys),
	// and 2 (which implies that the first encountered caching layer -
	// "Edge-tier" in a non-Topology context - will strip query strings,
	// effectively preventing them from being passed in upstream requests, and
	// not use them in the cache key).
	QStringIgnore *int `json:"qstringIgnore" db:"qstring_ignore"`
	// RangeRequestHandling defines how HTTP GET requests with a Range header
	// will be handled by cache servers serving the Delivery Service's content.
	// The only valid values are 0 (which implies that Range requests will not
	// be cached at all), 1 (which implies that the background_fetch plugin will
	// be used to service the range request while caching the whole object), 2
	// (which implies that the cache_range_requests plugin will be used to cache
	// ranges as unique objects), and 3 (which implies that the slice plugin
	// will be used to slice range based requests into deterministic chunks.)
	RangeRequestHandling *int `json:"rangeRequestHandling" db:"range_request_handling"`
	// RangeSliceBlockSize defines the size of range request blocks - or
	// "slices" - used by the "slice" plugin. This has no effect if
	// RangeRequestHandling does not point to exactly 3. This may never legally
	// point to a value less than zero.
	RangeSliceBlockSize *int `json:"rangeSliceBlockSize" db:"range_slice_block_size"`
	// Regex Remap is a raw line to be inserted into "regex_remap.config" on the
	// cache server. Care is necessitated in its use, because the input is in no
	// way restricted, validated, or limited in scope to the Delivery Service.
	RegexRemap *string `json:"regexRemap" db:"regex_remap"`
	// Regional indicates whether the Delivery Service's MaxOriginConnections is
	// only per Cache Group, rather than divided over all Cache Servers in child
	// Cache Groups of the Origin.
	Regional bool `json:"regional" db:"regional"`
	// RegionalGeoBlocking defines whether or not whatever Regional Geo Blocking
	// rules are configured on the Traffic Router serving content for this
	// Delivery Service will have an effect on the traffic of this Delivery
	// Service.
	RegionalGeoBlocking bool `json:"regionalGeoBlocking" db:"regional_geo_blocking"`
	// RemapText is raw text to insert in "remap.config" on the cache servers
	// serving content for this Delivery Service. Care is necessitated in its
	// use, because the input is in no way restricted, validated, or limited in
	// scope to the Delivery Service.
	RemapText *string `json:"remapText" db:"remap_text"`
	// RequiredCapabilities is an array of capabilities required for this delivery service.
	RequiredCapabilities []string `json:"requiredCapabilities" db:"required_capabilities"`
	// RoutingName defines the lowest-level DNS label used by the Delivery
	// Service, e.g. if trafficcontrol.apache.org were a Delivery Service, it
	// would have a RoutingName of "trafficcontrol".
	RoutingName string `json:"routingName" db:"routing_name"`
	// ServiceCategory defines a category to which a Delivery Service may
	// belong, which will cause HTTP Responses containing content for the
	// Delivery Service to have the "X-CDN-SVC" header with a value that is the
	// XMLID of the Delivery Service.
	ServiceCategory *string `json:"serviceCategory" db:"service_category"`
	// Signed is a legacy field. It is allowed to be `true` if and only if
	// SigningAlgorithm is not nil.
	Signed bool `json:"signed"`
	// SigningAlgorithm is the name of the algorithm used to sign CDN URIs for
	// this Delivery Service's content, or nil if no URI signing is done for the
	// Delivery Service. This may only point to the values "url_sig" or
	// "uri_signing" when it is not `nil`.
	SigningAlgorithm *string `json:"signingAlgorithm" db:"signing_algorithm"`
	// SSLKeyVersion incremented whenever Traffic Portal generates new SSL keys
	// for the Delivery Service, effectively making it a "generational" marker.
	SSLKeyVersion *int `json:"sslKeyVersion" db:"ssl_key_version"`
	// Tenant is the Tenant to which the Delivery Service belongs.
	Tenant *string `json:"tenant"`
	// TenantID is the integral, unique identifier for the Tenant to which the
	// Delivery Service belongs.
	TenantID int `json:"tenantId" db:"tenant_id"`
	// TLSVersions is the list of explicitly supported TLS versions for cache
	// servers serving the Delivery Service's content.
	TLSVersions []string `json:"tlsVersions" db:"tls_versions"`
	// Topology is the name of the Topology used by the Delivery Service, or nil
	// if no Topology is used.
	Topology *string `json:"topology" db:"topology"`
	// TRResponseHeaders is a set of headers (separated by CRLF pairs as per the
	// HTTP spec) and their values (separated by a colon as per the HTTP spec)
	// which will be sent by Traffic Router in HTTP responses to client requests
	// for this Delivery Service's content. This has no effect on DNS-routed or
	// un-routed Delivery Service Types.
	TRResponseHeaders *string `json:"trResponseHeaders"`
	// TRRequestHeaders is an "array" of HTTP headers which should be logged
	// from client HTTP requests for this Delivery Service's content by Traffic
	// Router, separated by newlines. This has no effect on DNS-routed or
	// un-routed Delivery Service Types.
	TRRequestHeaders *string `json:"trRequestHeaders"`
	// Type describes how content is routed and cached for this Delivery Service
	// as well as what other properties have any meaning.
	Type *string `json:"type"`
	// TypeID is an integral, unique identifier for the Tenant to which the
	// Delivery Service belongs.
	TypeID int `json:"typeId" db:"type"`
	// XMLID is a unique identifier that is also the second lowest-level DNS
	// label used by the Delivery Service. For example, if a Delivery Service's
	// content may be requested from video.demo1.mycdn.ciab.test, it may be
	// inferred that the Delivery Service's XMLID is demo1.
	XMLID string `json:"xmlId" db:"xml_id"`
}

// DeliveryServiceV5 is the type of a Delivery Service as it appears in the
// latest minor version of Traffic Ops API version 5.
type DeliveryServiceV5 = DeliveryServiceV50

// TLSVersionsAlerts generates warning-level alerts for the Delivery Service's
// TLS versions array. It will warn if newer versions are disallowed while
// older, less secure versions are allowed, if there are unrecognized versions
// present, if the Delivery Service's Protocol does not make use of TLS
// Versions, and whenever TLSVersions are explicitly set at all.
//
// This does NOT verify that the Delivery Service's TLS versions are _valid_,
// it ONLY creates warnings based on conditions that are possibly detrimental
// to CDN operation, but can, in fact, work.
func (ds DeliveryServiceV5) TLSVersionsAlerts() Alerts {
	return tlsVersionsAlerts(ds.TLSVersions, util.Coalesce(ds.Protocol, 3))
}

// Value implements the database/sql/driver.Valuer interface by marshaling the
// struct to JSON to pass back as an encoding/json.RawMessage.
func (ds *DeliveryServiceV5) Value() (driver.Value, error) {
	return json.Marshal(ds)
}

// Scan implements the database/sql.Scanner interface.
//
// This expects src to be an encoding/json.RawMessage and unmarshals that into
// the DeliveryServiceV5.
func (ds *DeliveryServiceV5) Scan(src interface{}) error {
	return jsonScan(src, ds)
}

// Downgrade downgrades an APIv5 Delivery Service into an APIv4 Delivery Service
// of the latest minor version.
func (ds DeliveryServiceV5) Downgrade() DeliveryServiceV4 {
	downgraded := DeliveryServiceV4{
		DeliveryServiceV40: DeliveryServiceV40{
			DeliveryServiceFieldsV31: DeliveryServiceFieldsV31{
				MaxRequestHeaderBytes: util.CopyIfNotNil(ds.MaxRequestHeaderBytes),
			},
			DeliveryServiceFieldsV30: DeliveryServiceFieldsV30{
				FirstHeaderRewrite: util.CopyIfNotNil(ds.FirstHeaderRewrite),
				InnerHeaderRewrite: util.CopyIfNotNil(ds.InnerHeaderRewrite),
				LastHeaderRewrite:  util.CopyIfNotNil(ds.LastHeaderRewrite),
				ServiceCategory:    util.CopyIfNotNil(ds.ServiceCategory),
				Topology:           util.CopyIfNotNil(ds.Topology),
			},
			DeliveryServiceFieldsV15: DeliveryServiceFieldsV15{
				EcsEnabled:          ds.EcsEnabled,
				RangeSliceBlockSize: util.CopyIfNotNil(ds.RangeSliceBlockSize),
			},
			DeliveryServiceFieldsV14: DeliveryServiceFieldsV14{
				ConsistentHashQueryParams: make([]string, len(ds.ConsistentHashQueryParams)),
				ConsistentHashRegex:       util.CopyIfNotNil(ds.ConsistentHashRegex),
				MaxOriginConnections:      util.CopyIfNotNil(ds.MaxOriginConnections),
			},
			DeliveryServiceFieldsV13: DeliveryServiceFieldsV13{
				DeepCachingType:   new(DeepCachingType),
				FQPacingRate:      util.CopyIfNotNil(ds.FQPacingRate),
				SigningAlgorithm:  util.CopyIfNotNil(ds.SigningAlgorithm),
				Tenant:            util.CopyIfNotNil(ds.Tenant),
				TRResponseHeaders: util.CopyIfNotNil(ds.TRResponseHeaders),
				TRRequestHeaders:  util.CopyIfNotNil(ds.TRRequestHeaders),
			},
			DeliveryServiceNullableFieldsV11: DeliveryServiceNullableFieldsV11{
				Active:                   new(bool),
				AnonymousBlockingEnabled: util.BoolPtr(ds.AnonymousBlockingEnabled),
				CCRDNSTTL:                util.CopyIfNotNil(ds.CCRDNSTTL),
				CDNID:                    util.IntPtr(ds.CDNID),
				CDNName:                  util.CopyIfNotNil(ds.CDNName),
				CheckPath:                util.CopyIfNotNil(ds.CheckPath),
				DisplayName:              util.StrPtr(ds.DisplayName),
				DNSBypassCNAME:           util.CopyIfNotNil(ds.DNSBypassCNAME),
				DNSBypassIP:              util.CopyIfNotNil(ds.DNSBypassIP),
				DNSBypassIP6:             util.CopyIfNotNil(ds.DNSBypassIP6),
				DNSBypassTTL:             util.CopyIfNotNil(ds.DNSBypassTTL),
				DSCP:                     util.IntPtr(ds.DSCP),
				EdgeHeaderRewrite:        util.CopyIfNotNil(ds.EdgeHeaderRewrite),
				GeoLimit:                 util.IntPtr(ds.GeoLimit),
				GeoLimitRedirectURL:      util.CopyIfNotNil(ds.GeoLimitRedirectURL),
				GeoProvider:              util.IntPtr(ds.GeoProvider),
				GlobalMaxMBPS:            util.CopyIfNotNil(ds.GlobalMaxMBPS),
				GlobalMaxTPS:             util.CopyIfNotNil(ds.GlobalMaxTPS),
				HTTPBypassFQDN:           util.CopyIfNotNil(ds.HTTPBypassFQDN),
				ID:                       util.CopyIfNotNil(ds.ID),
				InfoURL:                  util.CopyIfNotNil(ds.InfoURL),
				InitialDispersion:        util.CopyIfNotNil(ds.InitialDispersion),
				IPV6RoutingEnabled:       util.CopyIfNotNil(ds.IPV6RoutingEnabled),
				LastUpdated:              TimeNoModFromTime(ds.LastUpdated),
				LogsEnabled:              util.BoolPtr(ds.LogsEnabled),
				LongDesc:                 util.StrPtr(ds.LongDesc),
				MaxDNSAnswers:            util.CopyIfNotNil(ds.MaxDNSAnswers),
				MidHeaderRewrite:         util.CopyIfNotNil(ds.MidHeaderRewrite),
				MissLat:                  util.CopyIfNotNil(ds.MissLat),
				MissLong:                 util.CopyIfNotNil(ds.MissLong),
				MultiSiteOrigin:          util.BoolPtr(ds.MultiSiteOrigin),
				OriginShield:             util.CopyIfNotNil(ds.OriginShield),
				OrgServerFQDN:            util.CopyIfNotNil(ds.OrgServerFQDN),
				ProfileDesc:              util.CopyIfNotNil(ds.ProfileDesc),
				ProfileID:                util.CopyIfNotNil(ds.ProfileID),
				ProfileName:              util.CopyIfNotNil(ds.ProfileName),
				Protocol:                 util.CopyIfNotNil(ds.Protocol),
				QStringIgnore:            util.CopyIfNotNil(ds.QStringIgnore),
				RangeRequestHandling:     util.CopyIfNotNil(ds.RangeRequestHandling),
				RegexRemap:               util.CopyIfNotNil(ds.RegexRemap),
				RegionalGeoBlocking:      util.BoolPtr(ds.RegionalGeoBlocking),
				RemapText:                util.CopyIfNotNil(ds.RemapText),
				RoutingName:              util.StrPtr(ds.RoutingName),
				Signed:                   ds.Signed,
				SSLKeyVersion:            util.CopyIfNotNil(ds.SSLKeyVersion),
				TenantID:                 util.IntPtr(ds.TenantID),
				Type:                     (*DSType)(util.CopyIfNotNil(ds.Type)),
				TypeID:                   util.IntPtr(ds.TypeID),
				XMLID:                    util.StrPtr(ds.XMLID),
			},
			TLSVersions: make([]string, len(ds.TLSVersions)),
		},
		Regional:             ds.Regional,
		RequiredCapabilities: make([]string, len(ds.RequiredCapabilities)),
	}

	*downgraded.Active = ds.Active == DSActiveStateActive
	copy(downgraded.ConsistentHashQueryParams, ds.ConsistentHashQueryParams)
	if ds.ExampleURLs != nil {
		downgraded.ExampleURLs = make([]string, len(ds.ExampleURLs))
		copy(downgraded.ExampleURLs, ds.ExampleURLs)
	}
	if len(ds.GeoLimitCountries) > 0 {
		countries := make([]string, len(ds.GeoLimitCountries))
		copy(countries, ds.GeoLimitCountries)
		downgraded.GeoLimitCountries = GeoLimitCountriesType(countries)
	}
	if ds.MatchList != nil {
		downgraded.MatchList = new([]DeliveryServiceMatch)
		*downgraded.MatchList = make([]DeliveryServiceMatch, len(ds.MatchList))
		copy(*downgraded.MatchList, ds.MatchList)
	}
	copy(downgraded.TLSVersions, ds.TLSVersions)
	if len(ds.RequiredCapabilities) > 0 {
		copy(downgraded.RequiredCapabilities, ds.RequiredCapabilities)
	}
	return downgraded
}

// Upgrade upgrades an APIv4 Delivery Service into an APIv5 Delivery Service of
// the latest minor version.
func (ds DeliveryServiceV4) Upgrade() DeliveryServiceV5 {
	upgraded := DeliveryServiceV5{
		AnonymousBlockingEnabled:  util.CoalesceToDefault(ds.AnonymousBlockingEnabled),
		CCRDNSTTL:                 util.CopyIfNotNil(ds.CCRDNSTTL),
		CDNID:                     util.Coalesce(ds.CDNID, -1),
		CDNName:                   util.CopyIfNotNil(ds.CDNName),
		CheckPath:                 util.CopyIfNotNil(ds.CheckPath),
		ConsistentHashQueryParams: make([]string, len(ds.ConsistentHashQueryParams)),
		ConsistentHashRegex:       util.CopyIfNotNil(ds.ConsistentHashRegex),
		DeepCachingType:           util.CoalesceToDefault(ds.DeepCachingType),
		DisplayName:               util.CoalesceToDefault(ds.DisplayName),
		DNSBypassCNAME:            util.CopyIfNotNil(ds.DNSBypassCNAME),
		DNSBypassIP:               util.CopyIfNotNil(ds.DNSBypassIP),
		DNSBypassIP6:              util.CopyIfNotNil(ds.DNSBypassIP6),
		DNSBypassTTL:              util.CopyIfNotNil(ds.DNSBypassTTL),
		DSCP:                      util.Coalesce(ds.DSCP, -1),
		EcsEnabled:                ds.EcsEnabled,
		EdgeHeaderRewrite:         util.CopyIfNotNil(ds.EdgeHeaderRewrite),
		ExampleURLs:               nil,
		FirstHeaderRewrite:        util.CopyIfNotNil(ds.FirstHeaderRewrite),
		FQPacingRate:              util.CopyIfNotNil(ds.FQPacingRate),
		GeoLimit:                  util.Coalesce(ds.GeoLimit, -1),
		GeoLimitCountries:         make([]string, len(ds.GeoLimitCountries)),
		GeoLimitRedirectURL:       util.CopyIfNotNil(ds.GeoLimitRedirectURL),
		GeoProvider:               util.Coalesce(ds.GeoProvider, -1),
		GlobalMaxMBPS:             util.CopyIfNotNil(ds.GlobalMaxMBPS),
		GlobalMaxTPS:              util.CopyIfNotNil(ds.GlobalMaxTPS),
		HTTPBypassFQDN:            util.CopyIfNotNil(ds.HTTPBypassFQDN),
		ID:                        util.CopyIfNotNil(ds.ID),
		InfoURL:                   util.CopyIfNotNil(ds.InfoURL),
		InitialDispersion:         util.CopyIfNotNil(ds.InitialDispersion),
		InnerHeaderRewrite:        util.CopyIfNotNil(ds.InnerHeaderRewrite),
		IPV6RoutingEnabled:        util.CopyIfNotNil(ds.IPV6RoutingEnabled),
		LastHeaderRewrite:         util.CopyIfNotNil(ds.LastHeaderRewrite),
		LogsEnabled:               util.Coalesce(ds.LogsEnabled, false),
		LongDesc:                  util.CoalesceToDefault(ds.LongDesc),
		MatchList:                 nil,
		MaxDNSAnswers:             util.CopyIfNotNil(ds.MaxDNSAnswers),
		MaxOriginConnections:      util.CopyIfNotNil(ds.MaxOriginConnections),
		MaxRequestHeaderBytes:     util.CopyIfNotNil(ds.MaxRequestHeaderBytes),
		MidHeaderRewrite:          util.CopyIfNotNil(ds.MidHeaderRewrite),
		MissLat:                   util.CopyIfNotNil(ds.MissLat),
		MissLong:                  util.CopyIfNotNil(ds.MissLong),
		MultiSiteOrigin:           util.CoalesceToDefault(ds.MultiSiteOrigin),
		OriginShield:              util.CopyIfNotNil(ds.OriginShield),
		OrgServerFQDN:             util.CopyIfNotNil(ds.OrgServerFQDN),
		ProfileDesc:               util.CopyIfNotNil(ds.ProfileDesc),
		ProfileID:                 util.CopyIfNotNil(ds.ProfileID),
		ProfileName:               util.CopyIfNotNil(ds.ProfileName),
		Protocol:                  util.CopyIfNotNil(ds.Protocol),
		QStringIgnore:             util.CopyIfNotNil(ds.QStringIgnore),
		RangeRequestHandling:      util.CopyIfNotNil(ds.RangeRequestHandling),
		RangeSliceBlockSize:       util.CopyIfNotNil(ds.RangeSliceBlockSize),
		RegexRemap:                util.CopyIfNotNil(ds.RegexRemap),
		Regional:                  ds.Regional,
		RegionalGeoBlocking:       util.CoalesceToDefault(ds.RegionalGeoBlocking),
		RemapText:                 util.CopyIfNotNil(ds.RemapText),
		RequiredCapabilities:      make([]string, len(ds.RequiredCapabilities)),
		RoutingName:               util.CoalesceToDefault(ds.RoutingName),
		ServiceCategory:           util.CopyIfNotNil(ds.ServiceCategory),
		Signed:                    ds.Signed,
		SigningAlgorithm:          util.CopyIfNotNil(ds.SigningAlgorithm),
		SSLKeyVersion:             util.CopyIfNotNil(ds.SSLKeyVersion),
		Tenant:                    util.CopyIfNotNil(ds.Tenant),
		TenantID:                  util.Coalesce(ds.TenantID, -1),
		TLSVersions:               make([]string, len(ds.TLSVersions)),
		Topology:                  util.CopyIfNotNil(ds.Topology),
		TRResponseHeaders:         util.CopyIfNotNil(ds.TRResponseHeaders),
		TRRequestHeaders:          util.CopyIfNotNil(ds.TRRequestHeaders),
		Type:                      (*string)(util.CopyIfNotNil(ds.Type)),
		TypeID:                    util.Coalesce(ds.TypeID, -1),
		XMLID:                     util.CoalesceToDefault(ds.XMLID),
	}

	if ds.Active == nil || !*ds.Active {
		upgraded.Active = DSActiveStatePrimed
	} else {
		upgraded.Active = DSActiveStateActive
	}
	copy(upgraded.ConsistentHashQueryParams, ds.ConsistentHashQueryParams)
	if ds.ExampleURLs != nil {
		upgraded.ExampleURLs = make([]string, len(ds.ExampleURLs))
		copy(upgraded.ExampleURLs, ds.ExampleURLs)
	}
	copy(upgraded.GeoLimitCountries, ds.GeoLimitCountries)
	if ds.LastUpdated != nil {
		upgraded.LastUpdated = ds.LastUpdated.Time
	}
	if ds.MatchList != nil && len(*ds.MatchList) > 0 {
		upgraded.MatchList = make([]DeliveryServiceMatch, len(*ds.MatchList))
		copy(upgraded.MatchList, *ds.MatchList)
	}
	copy(upgraded.TLSVersions, ds.TLSVersions)
	if len(ds.RequiredCapabilities) > 0 {
		copy(upgraded.RequiredCapabilities, ds.RequiredCapabilities)
	}

	return upgraded
}

// DeliveryServicesResponseV5 is the type of a response from the
// /deliveryservices Traffic Ops endpoint in version 5 of its API.
type DeliveryServicesResponseV5 struct {
	Alerts
	Response []DeliveryServiceV5 `json:"response"`
}

// DeliveryServiceResponseV5 is the type of a response for API endponts
// returning a single Delivery Service in Traffic Ops API version 5.
type DeliveryServiceResponseV5 struct {
	Alerts
	Response DeliveryServiceV5 `json:"response"`
}
