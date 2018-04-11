package tc

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

//
// GetDeliveryServiceResponse is deprecated use DeliveryServicesResponse...
type GetDeliveryServiceResponse struct {
	Response []DeliveryService `json:"response"`
}

// DeliveryServicesResponse ...
type DeliveryServicesResponse struct {
	Response []DeliveryService `json:"response"`
}

// CreateDeliveryServiceResponse ...
type CreateDeliveryServiceResponse struct {
	Response []DeliveryService      `json:"response"`
	Alerts   []DeliveryServiceAlert `json:"alerts"`
}

// UpdateDeliveryServiceResponse ...
type UpdateDeliveryServiceResponse struct {
	Response []DeliveryService      `json:"response"`
	Alerts   []DeliveryServiceAlert `json:"alerts"`
}

// DeliveryServiceResponse ...
type DeliveryServiceResponse struct {
	Response DeliveryService        `json:"response"`
	Alerts   []DeliveryServiceAlert `json:"alerts"`
}

// DeleteDeliveryServiceResponse ...
type DeleteDeliveryServiceResponse struct {
	Alerts []DeliveryServiceAlert `json:"alerts"`
}

// DeliveryService ...
type DeliveryService struct {
	Active               bool                   `json:"active"`
	CacheURL             string                 `json:"cacheurl"`
	CCRDNSTTL            int                    `json:"ccrDnsTtl"`
	CDNID                int                    `json:"cdnId"`
	CDNName              string                 `json:"cdnName"`
	CheckPath            string                 `json:"checkPath"`
	DeepCachingType      DeepCachingType        `json:"deepCachingType"`
	DisplayName          string                 `json:"displayName"`
	DNSBypassCname       string                 `json:"dnsBypassCname"`
	DNSBypassIP          string                 `json:"dnsBypassIp"`
	DNSBypassIP6         string                 `json:"dnsBypassIp6"`
	DNSBypassTTL         int                    `json:"dnsBypassTtl"`
	DSCP                 int                    `json:"dscp"`
	EdgeHeaderRewrite    string                 `json:"edgeHeaderRewrite"`
	ExampleURLs          []string               `json:"exampleURLs"`
	GeoLimit             int                    `json:"geoLimit"`
	FQPacingRate         int                    `json:"fqPacingRate"`
	GeoProvider          int                    `json:"geoProvider"`
	GlobalMaxMBPS        int                    `json:"globalMaxMbps"`
	GlobalMaxTPS         int                    `json:"globalMaxTps"`
	HTTPBypassFQDN       string                 `json:"httpBypassFqdn"`
	ID                   int                    `json:"id"`
	InfoURL              string                 `json:"infoUrl"`
	InitialDispersion    float32                `json:"initialDispersion"`
	IPV6RoutingEnabled   bool                   `json:"ipv6RoutingEnabled"`
	LastUpdated          *TimeNoMod             `json:"lastUpdated" db:"last_updated"`
	LogsEnabled          bool                   `json:"logsEnabled"`
	LongDesc             string                 `json:"longDesc"`
	LongDesc1            string                 `json:"longDesc1"`
	LongDesc2            string                 `json:"longDesc2"`
	MatchList            []DeliveryServiceMatch `json:"matchList,omitempty"`
	MaxDNSAnswers        int                    `json:"maxDnsAnswers"`
	MidHeaderRewrite     string                 `json:"midHeaderRewrite"`
	MissLat              float64                `json:"missLat"`
	MissLong             float64                `json:"missLong"`
	MultiSiteOrigin      bool                   `json:"multiSiteOrigin"`
	OrgServerFQDN        string                 `json:"orgServerFqdn"`
	ProfileDesc          string                 `json:"profileDescription"`
	ProfileID            int                    `json:"profileId,omitempty"`
	ProfileName          string                 `json:"profileName"`
	Protocol             int                    `json:"protocol"`
	QStringIgnore        int                    `json:"qstringIgnore"`
	RangeRequestHandling int                    `json:"rangeRequestHandling"`
	RegexRemap           string                 `json:"regexRemap"`
	RegionalGeoBlocking  bool                   `json:"regionalGeoBlocking"`
	RemapText            string                 `json:"remapText"`
	RoutingName          string                 `json:"routingName"`
	SigningAlgorithm     string                 `json:"signingAlgorithm" db:"signing_algorithm"`
	TypeID               int                    `json:"typeId"`
	Type                 string                 `json:"type"`
	TRResponseHeaders    string                 `json:"trResponseHeaders"`
	TenantID             int                    `json:"tenantId,omitempty"`
	XMLID                string                 `json:"xmlId"`
}

// DeliveryServiceNullable - a version of the deliveryservice that allows for all fields to be null
type DeliveryServiceNullable struct {
	// NOTE: the db: struct tags are used for testing to map to their equivalent database column (if there is one)
	//
	Active                   *bool                   `json:"active" db:"active"`
	CacheURL                 *string                 `json:"cacheurl" db:"cacheurl"`
	CCRDNSTTL                *int                    `json:"ccrDnsTtl" db:"ccr_dns_ttl"`
	CDNID                    *int                    `json:"cdnId" db:"cdn_id"`
	CDNName                  *string                 `json:"cdnName"`
	CheckPath                *string                 `json:"checkPath" db:"check_path"`
	DeepCachingType          *DeepCachingType        `json:"deepCachingType" db:"deep_caching_type"`
	DisplayName              *string                 `json:"displayName" db:"display_name"`
	DNSBypassCNAME           *string                 `json:"dnsBypassCname" db:"dns_bypass_cname"`
	DNSBypassIP              *string                 `json:"dnsBypassIp" db:"dns_bypass_ip"`
	DNSBypassIP6             *string                 `json:"dnsBypassIp6" db:"dns_bypass_ip6"`
	DNSBypassTTL             *int                    `json:"dnsBypassTtl" db:"dns_bypass_ttl"`
	DSCP                     *int                    `json:"dscp" db:"dscp"`
	EdgeHeaderRewrite        *string                 `json:"edgeHeaderRewrite" db:"edge_header_rewrite"`
	FQPacingRate             *int                    `json:"fqPacingRate" db:"fq_pacing_rate"`
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
	LongDesc1                *string                 `json:"longDesc1" db:"long_desc_1"`
	LongDesc2                *string                 `json:"longDesc2" db:"long_desc_2"`
	MatchList                *[]DeliveryServiceMatch `json:"matchList"`
	MaxDNSAnswers            *int                    `json:"maxDnsAnswers" db:"max_dns_answers"`
	MidHeaderRewrite         *string                 `json:"midHeaderRewrite" db:"mid_header_rewrite"`
	MissLat                  *float64                `json:"missLat" db:"miss_lat"`
	MissLong                 *float64                `json:"missLong" db:"miss_long"`
	MultiSiteOrigin          *bool                   `json:"multiSiteOrigin" db:"multi_site_origin"`
	MultiSiteOriginAlgorithm *int                    `json:"multiSiteOriginAlgorithm" db:"multi_site_origin_algorithm"`
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
	SigningAlgorithm         *string                 `json:"signingAlgorithm" db:"signing_algorithm"`
	SSLKeyVersion            *int                    `json:"sslKeyVersion" db:"ssl_key_version"`
	TRRequestHeaders         *string                 `json:"trRequestHeaders" db:"tr_request_headers"`
	TRResponseHeaders        *string                 `json:"trResponseHeaders" db:"tr_response_headers"`
	TenantID                 *int                    `json:"tenantId" db:"tenant_id"`
	TypeName                 *string                 `json:"typeName"`
	TypeID                   *int                    `json:"typeId" db:"type"`
	XMLID                    *string                 `json:"xmlId" db:"xml_id"`
}

// Value implements the driver.Valuer interface
// marshals struct to json to pass back as a json.RawMessage
func (d *DeliveryServiceNullable) Value() (driver.Value, error) {
	b, err := json.Marshal(d)
	return b, err
}

// Scan implements the sql.Scanner interface
// expects json.RawMessage and unmarshals to a deliveryservice struct
func (d *DeliveryServiceNullable) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected deliveryservice in byte array form; got %T", src)
	}
	return json.Unmarshal(b, d)
}

// DeliveryServiceMatch ...
type DeliveryServiceMatch struct {
	Type      string `json:"type"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}

// DeliveryServiceAlert ...
type DeliveryServiceAlert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
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

// DeliveryServiceHealthResponse ...
type DeliveryServiceHealthResponse struct {
	Response DeliveryServiceHealth `json:"response"`
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

// DeliveryServiceCapacityResponse ...
type DeliveryServiceCapacityResponse struct {
	Response DeliveryServiceCapacity `json:"response"`
}

// DeliveryServiceCapacity ...
type DeliveryServiceCapacity struct {
	AvailablePercent   float64 `json:"availablePercent"`
	UnavailablePercent float64 `json:"unavailablePercent"`
	UtilizedPercent    float64 `json:"utilizedPercent"`
	MaintenancePercent float64 `json:"maintenancePercent"`
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

// DeliveryServiceServerResponse ...
type DeliveryServiceServerResponse struct {
	Response []DeliveryServiceServer `json:"response"`
	Size     int                     `json:"size"`
	OrderBy  string                  `json:"orderby"`
	Limit    int                     `json:"limit"`
}

// DeliveryServiceServer ...
type DeliveryServiceServer struct {
	LastUpdated     string `json:"lastUpdated"`
	Server          int    `json:"server"`
	DeliveryService int    `json:"deliveryService"`
}
