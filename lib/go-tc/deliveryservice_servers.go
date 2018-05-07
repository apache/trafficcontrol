package tc

import "time"

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

// DeliveryServiceServerResponse ...
type DeliveryServiceServerResponse struct {
	Response []DeliveryServiceServer `json:"response"`
	Size     int                     `json:"size"`
	OrderBy  string                  `json:"orderby"`
	Limit    int                     `json:"limit"`
}

// DeliveryServiceServer ...
type DeliveryServiceServer struct {
	Server          *int             `json:"server"`
	DeliveryService *int             `json:"deliveryService"`
	LastUpdated     *TimeNoMod       `json:"lastUpdated" db:"last_updated"`
}


type DssServer struct {
	Cachegroup       *string              `json:"cachegroup" db:"cachegroup"`
	CachegroupID     *int                 `json:"cachegroupId" db:"cachegroup_id"`
	CDNID            *int                 `json:"cdnId" db:"cdn_id"`
	CDNName          *string              `json:"cdnName" db:"cdn_name"`
	DeliveryServices *map[string][]string `json:"deliveryServices,omitempty"`
	DomainName       *string              `json:"domainName" db:"domain_name"`
	FQDN             *string              `json:"fqdn,omitempty"`
	FqdnTime         time.Time            `json:"-"`
	GUID             *string              `json:"guid" db:"guid"`
	HostName         *string              `json:"hostName" db:"host_name"`
	HTTPSPort        *int                 `json:"httpsPort" db:"https_port"`
	ID               *int                 `json:"id" db:"id"`
	ILOIPAddress     *string              `json:"iloIpAddress" db:"ilo_ip_address"`
	ILOIPGateway     *string              `json:"iloIpGateway" db:"ilo_ip_gateway"`
	ILOIPNetmask     *string              `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	ILOPassword      *string              `json:"iloPassword" db:"ilo_password"`
	ILOUsername      *string              `json:"iloUsername" db:"ilo_username"`
	InterfaceMtu     *int                 `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName    *string              `json:"interfaceName" db:"interface_name"`
	IP6Address       *string              `json:"ip6Address" db:"ip6_address"`
	IP6Gateway       *string              `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress        *string              `json:"ipAddress" db:"ip_address"`
	IPGateway        *string              `json:"ipGateway" db:"ip_gateway"`
	IPNetmask        *string              `json:"ipNetmask" db:"ip_netmask"`
	LastUpdated      *TimeNoMod           `json:"lastUpdated" db:"last_updated"`
	MgmtIPAddress    *string              `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIPGateway    *string              `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIPNetmask    *string              `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason    *string              `json:"offlineReason" db:"offline_reason"`
	PhysLocation     *string              `json:"physLocation" db:"phys_location"`
	PhysLocationID   *int                 `json:"physLocationId" db:"phys_location_id"`
	Profile          *string              `json:"profile" db:"profile"`
	ProfileDesc      *string              `json:"profileDesc" db:"profile_desc"`
	ProfileID        *int                 `json:"profileId" db:"profile_id"`
	Rack             *string              `json:"rack" db:"rack"`
	RouterHostName   *string              `json:"routerHostName" db:"router_host_name"`
	RouterPortName   *string              `json:"routerPortName" db:"router_port_name"`
	Status           *string              `json:"status" db:"status"`
	StatusID         *int                 `json:"statusId" db:"status_id"`
	TCPPort          *int                 `json:"tcpPort" db:"tcp_port"`
	Type             string               `json:"type" db:"server_type"`
	TypeID           *int                 `json:"typeId" db:"server_type_id"`
	UpdPending       *bool                `json:"updPending" db:"upd_pending"`
}

// DeliveryServiceNullable - a version of the deliveryservice that allows for all fields to be null
type DssDeliveryService struct {
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
	SigningAlgorithm         *string                 `json:"signingAlgorithm" db:"signing_algorithm"`
	SSLKeyVersion            *int                    `json:"sslKeyVersion" db:"ssl_key_version"`
	TRRequestHeaders         *string                 `json:"trRequestHeaders" db:"tr_request_headers"`
	TRResponseHeaders        *string                 `json:"trResponseHeaders" db:"tr_response_headers"`
	TenantID                 *int                    `json:"tenantId" db:"tenant_id"`
	TypeName                 *string                 `json:"typeName"`
	TypeID                   *int                    `json:"typeId" db:"type"`
	XMLID                    *string                 `json:"xmlId" db:"xml_id"`
}
