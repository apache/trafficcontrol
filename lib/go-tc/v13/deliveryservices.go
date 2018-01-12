package v13

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
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// This struct supports /api/1.3

// DeliveryServicesResponse ...
type DeliveryServicesResponse struct {
	Response []DeliveryService `json:"response"`
}

// DeliveryService ...
type DeliveryService struct {
	// NOTE: Fields that are pointers (with an asterisk '*') are required
	// for existence checking
	//
	// NOTE: the db: struct tags are used for testing to map to their equivalent database column (if there is one)
	//
	Active                   *bool                  `json:"active" db:"active"`
	CacheURL                 string                 `json:"cacheurl" db:"cacheurl"`
	CCRDNSTTL                int                    `json:"ccrDnsTtl" db:"ccr_dns_ttl"`
	CDNID                    *int                   `json:"cdnId" db:"cdn_id"`
	CheckPath                string                 `json:"checkPath" db:"check_path"`
	DisplayName              *string                `json:"displayName" db:"display_name"`
	CDNName                  string                 `json:"cdnName"`
	DNSBypassCNAME           string                 `json:"dnsBypassCname" db:"dns_bypass_cname"`
	DNSBypassIP              string                 `json:"dnsBypassIp" db:"dns_bypass_ip"`
	DNSBypassIP6             string                 `json:"dnsBypassIp6" db:"dns_bypass_ip6"`
	DNSBypassTTL             int                    `json:"dnsBypassTtl" db:"dns_bypass_ttl"`
	DSCP                     *int                   `json:"dscp" db:"dscp"`
	EdgeHeaderRewrite        string                 `json:"edgeHeaderRewrite" db:"edge_header_rewrite"`
	GeoLimit                 *int                   `json:"geoLimit" db:"geo_limit"`
	GeoLimitCountries        string                 `json:"geoLimitCountries" db:"geo_limit_countries"`
	GeoLimitRedirectURL      string                 `json:"geoLimitRedirectUrl" db:"geolimit_redirect_url"`
	GeoProvider              *int                   `json:"geoProvider" db:"geo_provider"`
	GlobalMaxMBPS            *int                   `json:"globalMaxMbps" db:"global_max_mbps"`
	GlobalMaxTPS             *int                   `json:"globalMaxTps" db:"global_max_tps"`
	HTTPBypassFQDN           string                 `json:"httpBypassFqdn" db:"http_bypass_fqdn"`
	ID                       int                    `json:"id" db:"id"`
	InfoURL                  string                 `json:"infoUrl" db:"info_url"`
	InitialDispersion        *int                   `json:"initialDispersion" db:"initial_dispersion"`
	IPV6RoutingEnabled       bool                   `json:"ipv6RoutingEnabled" db:"ipv6_routing_enabled"`
	LastUpdated              tc.Time                `json:"lastUpdated" db:"last_updated""`
	LogsEnabled              *bool                  `json:"logsEnabled" db:"logs_enabled"`
	LongDesc                 string                 `json:"longDesc" db:"long_desc"`
	LongDesc1                string                 `json:"longDesc1" db:"long_desc_1"`
	LongDesc2                string                 `json:"longDesc2" db:"long_desc_2"`
	MatchList                []DeliveryServiceMatch `json:"matchList,omitempty"`
	MaxDNSAnswers            int                    `json:"maxDnsAnswers" db:"max_dns_answers"`
	MidHeaderRewrite         string                 `json:"midHeaderRewrite" db:"mid_header_rewrite"`
	MissLat                  float64                `json:"missLat" db:"miss_lat"`
	MissLong                 float64                `json:"missLong" db:"miss_long"`
	MultiSiteOrigin          bool                   `json:"multiSiteOrigin" db:"multi_site_origin"`
	MultiSiteOriginAlgorithm int                    `json:"multiSiteOriginAlgorithm" db:"multi_site_origin_algorithm"`
	OriginShield             string                 `json:"originShield" db:"origin_shield"`
	OrgServerFQDN            string                 `json:"orgServerFqdn" db:"org_server_fqdn"`
	ProfileDesc              string                 `json:"profileDescription"`
	ProfileID                int                    `json:"profileId,omitempty" db:"profile"`
	ProfileName              string                 `json:"profileName"`
	Protocol                 int                    `json:"protocol" db:"protocol"`
	QStringIgnore            int                    `json:"qstringIgnore" db:"qstring_ignore"`
	RangeRequestHandling     int                    `json:"rangeRequestHandling" db:"range_request_handling"`
	RegexRemap               string                 `json:"regexRemap" db:"regex_remap"`
	RegionalGeoBlocking      *bool                  `json:"regionalGeoBlocking" db:"regional_geo_blocking"`
	RemapText                string                 `json:"remapText" db:"remap_text"`
	RoutingName              string                 `json:"routingName" db:"routing_name"`
	SigningAlgorithm         string                 `json:"signingAlgorithm" db:"signing_algorithm"`
	SSLKeyVersion            int                    `json:"sslKeyVersion" db:"ssl_key_version"`
	TRRequestHeaders         string                 `json:"trRequestHeaders" db:"tr_request_headers"`
	TRResponseHeaders        string                 `json:"trResponseHeaders" db:"tr_response_headers"`
	TenantID                 int                    `json:"tenantId" db:"tenant_id"`
	TypeName                 string                 `json:"typeName"`
	TypeID                   *int                   `json:"typeId" db:"type"`
	XMLID                    *string                `json:"xmlId" db:"xml_id"`
}

// DeliveryServiceMatch ...
type DeliveryServiceMatch struct {
	Type      string `json:"type"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}
