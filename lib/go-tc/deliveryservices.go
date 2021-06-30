package tc

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-util"
)

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

const DefaultRoutingName = "cdn"
const DefaultMaxRequestHeaderBytes = 0
const MinRangeSliceBlockSize = 262144   // 265Kib
const MaxRangeSliceBlockSize = 33554432 // 32Mib

// GetDeliveryServiceResponse is deprecated use DeliveryServicesResponse...
type GetDeliveryServiceResponse struct {
	Response []DeliveryService `json:"response"`
}

// DeliveryServicesResponse ...
// Deprecated: use DeliveryServicesNullableResponse instead
type DeliveryServicesResponse struct {
	Response []DeliveryService `json:"response"`
	Alerts
}

// DeliveryServicesResponseV30 is the type of a response from the
// /api/3.0/deliveryservices Traffic Ops endpoint.
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

// DeliveryServicesResponseV4 is the type of a response from the
// /api/4.x/deliveryservices Traffic Ops endpoint.
// It always points to the type for the latest minor version of APIv4.
type DeliveryServicesResponseV4 = DeliveryServicesResponseV40

// DeliveryServicesNullableResponse ...
// Deprecated: Please only use the versioned structures.
type DeliveryServicesNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts
}

// CreateDeliveryServiceResponse ...
// Deprecated: use CreateDeliveryServiceNullableResponse instead
type CreateDeliveryServiceResponse struct {
	Response []DeliveryService `json:"response"`
	Alerts
}

// CreateDeliveryServiceNullableResponse ...
// Deprecated: Please only use the versioned structures.
type CreateDeliveryServiceNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts
}

// UpdateDeliveryServiceResponse ...
// Deprecated: use UpdateDeliveryServiceNullableResponse instead
type UpdateDeliveryServiceResponse struct {
	Response []DeliveryService `json:"response"`
	Alerts
}

// UpdateDeliveryServiceNullableResponse ...
// Deprecated: Please only use the versioned structures.
type UpdateDeliveryServiceNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts
}

// DeleteDeliveryServiceResponse ...
type DeleteDeliveryServiceResponse struct {
	Alerts
}

// Deprecated: use DeliveryServiceNullable instead
type DeliveryService struct {
	DeliveryServiceV13
	MaxOriginConnections      int      `json:"maxOriginConnections" db:"max_origin_connections"`
	ConsistentHashRegex       string   `json:"consistentHashRegex"`
	ConsistentHashQueryParams []string `json:"consistentHashQueryParams"`
}

type DeliveryServiceV13 struct {
	DeliveryServiceV11
	DeepCachingType   DeepCachingType `json:"deepCachingType"`
	FQPacingRate      int             `json:"fqPacingRate,omitempty"`
	SigningAlgorithm  string          `json:"signingAlgorithm" db:"signing_algorithm"`
	Tenant            string          `json:"tenant"`
	TRRequestHeaders  string          `json:"trRequestHeaders,omitempty"`
	TRResponseHeaders string          `json:"trResponseHeaders,omitempty"`
}

// DeliveryServiceV11 contains the information relating to a delivery service
// that was around in version 1.1 of the API.
// TODO move contents to DeliveryServiceV12, fix references, and remove
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

type DeliveryServiceV31 struct {
	DeliveryServiceV30
	DeliveryServiceFieldsV31
}

// DeliveryServiceFieldsV31 contains additions to delivery services in api v3.1
type DeliveryServiceFieldsV31 struct {
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
}

// DeliveryServiceV4 is a Delivery Service as it appears in version 4 of the
// Traffic Ops API - it always points to the highest minor version in APIv4.
type DeliveryServiceV4 = DeliveryServiceV40

type DeliveryServiceV30 struct {
	DeliveryServiceNullableV15
	DeliveryServiceFieldsV30
}

// DeliveryServiceFieldsV30 contains additions to delivery services in api v3.0
type DeliveryServiceFieldsV30 struct {
	Topology           *string `json:"topology" db:"topology"`
	FirstHeaderRewrite *string `json:"firstHeaderRewrite" db:"first_header_rewrite"`
	InnerHeaderRewrite *string `json:"innerHeaderRewrite" db:"inner_header_rewrite"`
	LastHeaderRewrite  *string `json:"lastHeaderRewrite" db:"last_header_rewrite"`
	ServiceCategory    *string `json:"serviceCategory" db:"service_category"`
}

// DeliveryServiceNullableV30 is the aliased structure that we should be using for all api 3.x delivery structure operations
// This type should always alias the latest 3.x minor version struct. For ex, if you wanted to create a DeliveryServiceV32 struct, you would do the following:
// type DeliveryServiceNullableV30 DeliveryServiceV32
// DeliveryServiceV32 = DeliveryServiceV31 + the new fields
type DeliveryServiceNullableV30 DeliveryServiceV31

// Deprecated: Use versioned structures only from now on.
type DeliveryServiceNullable DeliveryServiceNullableV15
type DeliveryServiceNullableV15 struct {
	DeliveryServiceNullableV14
	DeliveryServiceFieldsV15
}

// DeliveryServiceFieldsV15 contains additions to delivery services in api v1.5
type DeliveryServiceFieldsV15 struct {
	EcsEnabled          bool `json:"ecsEnabled" db:"ecs_enabled"`
	RangeSliceBlockSize *int `json:"rangeSliceBlockSize" db:"range_slice_block_size"`
}

type DeliveryServiceNullableV14 struct {
	DeliveryServiceNullableV13
	DeliveryServiceFieldsV14
}

// DeliveryServiceFieldsV14 contains additions to delivery services in api v1.4
type DeliveryServiceFieldsV14 struct {
	ConsistentHashRegex       *string  `json:"consistentHashRegex"`
	ConsistentHashQueryParams []string `json:"consistentHashQueryParams"`
	MaxOriginConnections      *int     `json:"maxOriginConnections" db:"max_origin_connections"`
}

type DeliveryServiceNullableV13 struct {
	DeliveryServiceNullableV12
	DeliveryServiceFieldsV13
}

// DeliveryServiceFieldsV13 contains additions to delivery services in api v1.3
type DeliveryServiceFieldsV13 struct {
	DeepCachingType   *DeepCachingType `json:"deepCachingType" db:"deep_caching_type"`
	FQPacingRate      *int             `json:"fqPacingRate" db:"fq_pacing_rate"`
	SigningAlgorithm  *string          `json:"signingAlgorithm" db:"signing_algorithm"`
	Tenant            *string          `json:"tenant"`
	TRResponseHeaders *string          `json:"trResponseHeaders"`
	TRRequestHeaders  *string          `json:"trRequestHeaders"`
}

type DeliveryServiceNullableV12 struct {
	DeliveryServiceNullableV11
}

// DeliveryServiceNullableV11 is a version of the deliveryservice that allows
// for all fields to be null.
// TODO move contents to DeliveryServiceNullableV12, fix references, and remove
type DeliveryServiceNullableV11 struct {
	DeliveryServiceNullableFieldsV11
	DeliveryServiceRemovedFieldsV11
}

type DeliveryServiceNullableFieldsV11 struct {
	Active                   *bool                   `json:"active" db:"active"`
	AnonymousBlockingEnabled *bool                   `json:"anonymousBlockingEnabled" db:"anonymous_blocking_enabled"`
	CCRDNSTTL                *int                    `json:"ccrDnsTtl" db:"ccr_dns_ttl"`
	CDNID                    *int                    `json:"cdnId" db:"cdn_id"`
	CDNName                  *string                 `json:"cdnName"`
	CheckPath                *string                 `json:"checkPath" db:"check_path"`
	DisplayName              *string                 `json:"displayName" db:"display_name"`
	DNSBypassCNAME           *string                 `json:"dnsBypassCname" db:"dns_bypass_cname"`
	DNSBypassIP              *string                 `json:"dnsBypassIp" db:"dns_bypass_ip"`
	DNSBypassIP6             *string                 `json:"dnsBypassIp6" db:"dns_bypass_ip6"`
	DNSBypassTTL             *int                    `json:"dnsBypassTtl" db:"dns_bypass_ttl"`
	DSCP                     *int                    `json:"dscp" db:"dscp"`
	EdgeHeaderRewrite        *string                 `json:"edgeHeaderRewrite" db:"edge_header_rewrite"`
	GeoLimit                 *int                    `json:"geoLimit" db:"geo_limit"`
	GeoLimitCountries        *string                 `json:"geoLimitCountries" db:"geo_limit_countries"`
	GeoLimitRedirectURL      *string                 `json:"geoLimitRedirectURL" db:"geolimit_redirect_url"`
	GeoProvider              *int                    `json:"geoProvider" db:"geo_provider"`
	GlobalMaxMBPS            *int                    `json:"globalMaxMbps" db:"global_max_mbps"`
	GlobalMaxTPS             *int                    `json:"globalMaxTps" db:"global_max_tps"`
	HTTPBypassFQDN           *string                 `json:"httpBypassFqdn" db:"http_bypass_fqdn"`
	ID                       *int                    `json:"id" db:"id"`
	InfoURL                  *string                 `json:"infoUrl" db:"info_url"`
	InitialDispersion        *int                    `json:"initialDispersion" db:"initial_dispersion"`
	IPV6RoutingEnabled       *bool                   `json:"ipv6RoutingEnabled" db:"ipv6_routing_enabled"`
	LastUpdated              *TimeNoMod              `json:"lastUpdated" db:"last_updated"`
	LogsEnabled              *bool                   `json:"logsEnabled" db:"logs_enabled"`
	LongDesc                 *string                 `json:"longDesc" db:"long_desc"`
	LongDesc1                *string                 `json:"longDesc1,omitempty" db:"long_desc_1"`
	LongDesc2                *string                 `json:"longDesc2,omitempty" db:"long_desc_2"`
	MatchList                *[]DeliveryServiceMatch `json:"matchList"`
	MaxDNSAnswers            *int                    `json:"maxDnsAnswers" db:"max_dns_answers"`
	MidHeaderRewrite         *string                 `json:"midHeaderRewrite" db:"mid_header_rewrite"`
	MissLat                  *float64                `json:"missLat" db:"miss_lat"`
	MissLong                 *float64                `json:"missLong" db:"miss_long"`
	MultiSiteOrigin          *bool                   `json:"multiSiteOrigin" db:"multi_site_origin"`
	OriginShield             *string                 `json:"originShield" db:"origin_shield"`
	OrgServerFQDN            *string                 `json:"orgServerFqdn" db:"org_server_fqdn"`
	ProfileDesc              *string                 `json:"profileDescription"`
	ProfileID                *int                    `json:"profileId" db:"profile"`
	ProfileName              *string                 `json:"profileName"`
	Protocol                 *int                    `json:"protocol" db:"protocol"`
	QStringIgnore            *int                    `json:"qstringIgnore" db:"qstring_ignore"`
	RangeRequestHandling     *int                    `json:"rangeRequestHandling" db:"range_request_handling"`
	RegexRemap               *string                 `json:"regexRemap" db:"regex_remap"`
	RegionalGeoBlocking      *bool                   `json:"regionalGeoBlocking" db:"regional_geo_blocking"`
	RemapText                *string                 `json:"remapText" db:"remap_text"`
	RoutingName              *string                 `json:"routingName" db:"routing_name"`
	Signed                   bool                    `json:"signed"`
	SSLKeyVersion            *int                    `json:"sslKeyVersion" db:"ssl_key_version"`
	TenantID                 *int                    `json:"tenantId" db:"tenant_id"`
	Type                     *DSType                 `json:"type"`
	TypeID                   *int                    `json:"typeId" db:"type"`
	XMLID                    *string                 `json:"xmlId" db:"xml_id"`
	ExampleURLs              []string                `json:"exampleURLs"`
}

// DeliveryServiceRemovedFieldsV11 contains additions to delivery services in api v1.1 that were later removed
// Deprecated: used for backwards compatibility  with ATC <v5.1
type DeliveryServiceRemovedFieldsV11 struct {
	CacheURL *string `json:"cacheurl" db:"cacheurl"`
}

// RemoveLD1AndLD2 removes the Long Description 1 and Long Description 2 fields from a V 4.x DS, and returns the resulting struct.
func (ds *DeliveryServiceV4) RemoveLD1AndLD2() DeliveryServiceV4 {
	ds.LongDesc1 = nil
	ds.LongDesc2 = nil
	return *ds
}

// DowngradeToV3 converts the 4.x DS to a 3.x DS
func (ds *DeliveryServiceV4) DowngradeToV3() DeliveryServiceNullableV30 {
	return DeliveryServiceNullableV30{
		DeliveryServiceV30: DeliveryServiceV30{
			DeliveryServiceNullableV15: DeliveryServiceNullableV15{
				DeliveryServiceNullableV14: DeliveryServiceNullableV14{
					DeliveryServiceNullableV13: DeliveryServiceNullableV13{
						DeliveryServiceNullableV12: DeliveryServiceNullableV12{
							DeliveryServiceNullableV11: DeliveryServiceNullableV11{
								DeliveryServiceNullableFieldsV11: ds.DeliveryServiceNullableFieldsV11,
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

// UpgradeToV4 converts the 3.x DS to a 4.x DS
func (ds *DeliveryServiceNullableV30) UpgradeToV4() DeliveryServiceV4 {
	return DeliveryServiceV4{
		DeliveryServiceFieldsV31:         ds.DeliveryServiceFieldsV31,
		DeliveryServiceFieldsV30:         ds.DeliveryServiceFieldsV30,
		DeliveryServiceFieldsV15:         ds.DeliveryServiceFieldsV15,
		DeliveryServiceFieldsV14:         ds.DeliveryServiceFieldsV14,
		DeliveryServiceFieldsV13:         ds.DeliveryServiceFieldsV13,
		DeliveryServiceNullableFieldsV11: ds.DeliveryServiceNullableFieldsV11,
	}
}

func jsonValue(v interface{}) (driver.Value, error) {
	b, err := json.Marshal(v)
	return b, err
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

// Value implements the driver.Valuer interface --
// marshals struct to json to pass back as a json.RawMessage.
func (ds *DeliveryServiceNullable) Value() (driver.Value, error) {
	return jsonValue(ds)
}

// Scan implements the sql.Scanner interface --
// expects json.RawMessage and unmarshals to a DeliveryServiceNullable struct.
func (ds *DeliveryServiceNullable) Scan(src interface{}) error {
	return jsonScan(src, ds)
}

// Value implements the driver.Valuer interface --
// marshals struct to json to pass back as a json.RawMessage.
func (ds *DeliveryServiceV4) Value() (driver.Value, error) {
	return jsonValue(ds)
}

// Scan implements the sql.Scanner interface --
// expects json.RawMessage and unmarshals to a DeliveryServiceV4 struct.
func (ds *DeliveryServiceV4) Scan(src interface{}) error {
	return jsonScan(src, ds)
}

// DeliveryServiceMatch ...
type DeliveryServiceMatch struct {
	Type      DSMatchType `json:"type"`
	SetNumber int         `json:"setNumber"`
	Pattern   string      `json:"pattern"`
}

// DeliveryServiceStateResponse ...
type DeliveryServiceStateResponse struct {
	Response DeliveryServiceState `json:"response"`
}

// DeliveryServiceState ...
type DeliveryServiceState struct {
	Enabled  bool                    `json:"enabled"`
	Failover DeliveryServiceFailover `json:"failover"`
}

// DeliveryServiceFailover ...
type DeliveryServiceFailover struct {
	Locations   []string                   `json:"locations"`
	Destination DeliveryServiceDestination `json:"destination"`
	Configured  bool                       `json:"configured"`
	Enabled     bool                       `json:"enabled"`
}

// DeliveryServiceDestination ...
type DeliveryServiceDestination struct {
	Location string `json:"location"`
	Type     string `json:"type"`
}

// DeliveryServiceHealthResponse is the type of a response from Traffic Ops to
// a request for a Delivery Service's "health".
type DeliveryServiceHealthResponse struct {
	Response DeliveryServiceHealth `json:"response"`
	Alerts
}

// DeliveryServiceHealth ...
type DeliveryServiceHealth struct {
	TotalOnline  int                         `json:"totalOnline"`
	TotalOffline int                         `json:"totalOffline"`
	CacheGroups  []DeliveryServiceCacheGroup `json:"cacheGroups"`
}

// DeliveryServiceCacheGroup ...
type DeliveryServiceCacheGroup struct {
	Online  int    `json:"online"`
	Offline int    `json:"offline"`
	Name    string `json:"name"`
}

// DeliveryServiceCapacityResponse is the type of a response from Traffic Ops to
// a request for a Delivery Service's "capacity".
type DeliveryServiceCapacityResponse struct {
	Response DeliveryServiceCapacity `json:"response"`
	Alerts
}

// DeliveryServiceCapacity ...
type DeliveryServiceCapacity struct {
	AvailablePercent   float64 `json:"availablePercent"`
	UnavailablePercent float64 `json:"unavailablePercent"`
	UtilizedPercent    float64 `json:"utilizedPercent"`
	MaintenancePercent float64 `json:"maintenancePercent"`
}

type DeliveryServiceMatchesResp []DeliveryServicePatterns

type DeliveryServicePatterns struct {
	Patterns []string            `json:"patterns"`
	DSName   DeliveryServiceName `json:"dsName"`
}

type DeliveryServiceMatchesResponse struct {
	Response []DeliveryServicePatterns `json:"response"`
}

// DeliveryServiceRoutingResponse ...
type DeliveryServiceRoutingResponse struct {
	Response DeliveryServiceRouting `json:"response"`
}

// DeliveryServiceRouting ...
type DeliveryServiceRouting struct {
	StaticRoute       int     `json:"staticRoute"`
	Miss              int     `json:"miss"`
	Geo               float64 `json:"geo"`
	Err               int     `json:"err"`
	CZ                float64 `json:"cz"`
	DSR               float64 `json:"dsr"`
	Fed               int     `json:"fed"`
	RegionalAlternate int     `json:"regionalAlternate"`
	RegionalDenied    int     `json:"regionalDenied"`
}

type UserAvailableDS struct {
	ID          *int    `json:"id" db:"id"`
	DisplayName *string `json:"displayName" db:"display_name"`
	XMLID       *string `json:"xmlId" db:"xml_id"`
	TenantID    *int    `json:"-"` // tenant is necessary to check authorization, but not serialized
}

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

type DeliveryServiceUserPost struct {
	UserID           *int   `json:"userId"`
	DeliveryServices *[]int `json:"deliveryServices"`
	Replace          *bool  `json:"replace"`
}

type UserDeliveryServicePostResponse struct {
	Alerts   []Alert                 `json:"alerts"`
	Response DeliveryServiceUserPost `json:"response"`
}

type UserDeliveryServicesNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
}

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

type CachegroupPostDSReq struct {
	DeliveryServices []int `json:"deliveryServices"`
}

type CacheGroupPostDSResp struct {
	ID               util.JSONIntStr `json:"id"`
	ServerNames      []CacheName     `json:"serverNames"`
	DeliveryServices []int           `json:"deliveryServices"`
}

type CacheGroupPostDSRespResponse struct {
	Alerts
	Response CacheGroupPostDSResp `json:"response"`
}

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

// Validate implements the github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api.ParseValidator
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
