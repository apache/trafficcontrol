package tc

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/asaskevich/govalidator"
	validation "github.com/go-ozzo/ozzo-validation"
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

// GetDeliveryServiceResponse is deprecated use DeliveryServicesResponse...
type GetDeliveryServiceResponse struct {
	Response []DeliveryService `json:"response"`
}

// DeliveryServicesResponse ...
type DeliveryServicesResponse struct {
	Response []DeliveryService `json:"response"`
}

// DeliveryServicesNullableResponse ...
type DeliveryServicesNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
}

// CreateDeliveryServiceResponse ...
type CreateDeliveryServiceResponse struct {
	Response []DeliveryService      `json:"response"`
	Alerts   []DeliveryServiceAlert `json:"alerts"`
}

// CreateDeliveryServiceNullableResponse ...
type CreateDeliveryServiceNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts   []DeliveryServiceAlert    `json:"alerts"`
}

// UpdateDeliveryServiceResponse ...
type UpdateDeliveryServiceResponse struct {
	Response []DeliveryService      `json:"response"`
	Alerts   []DeliveryServiceAlert `json:"alerts"`
}

// UpdateDeliveryServiceNullableResponse ...
type UpdateDeliveryServiceNullableResponse struct {
	Response []DeliveryServiceNullable `json:"response"`
	Alerts   []DeliveryServiceAlert    `json:"alerts"`
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

// DeliveryService ...
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

type DeliveryServiceNullableV14 DeliveryServiceNullable // this type alias should always alias the latest minor version of the deliveryservices endpoints

type DeliveryServiceNullable struct {
	DeliveryServiceNullableV13
	ConsistentHashRegex       *string  `json:"consistentHashRegex"`
	ConsistentHashQueryParams []string `json:"consistentHashQueryParams"`
	MaxOriginConnections      *int     `json:"maxOriginConnections" db:"max_origin_connections"`
}

type DeliveryServiceNullableV13 struct {
	DeliveryServiceNullableV12
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

// DeliveryServiceNullable - a version of the deliveryservice that allows for all fields to be null
// TODO move contents to DeliveryServiceNullableV12, fix references, and remove
type DeliveryServiceNullableV11 struct {
	// NOTE: the db: struct tags are used for testing to map to their equivalent database column (if there is one)
	//
	Active                   *bool                   `json:"active" db:"active"`
	AnonymousBlockingEnabled *bool                   `json:"anonymousBlockingEnabled" db:"anonymous_blocking_enabled"`
	CacheURL                 *string                 `json:"cacheurl" db:"cacheurl"`
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
	LongDesc1                *string                 `json:"longDesc1" db:"long_desc_1"`
	LongDesc2                *string                 `json:"longDesc2" db:"long_desc_2"`
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

func requiredIfMatchesTypeName(patterns []string, typeName string) func(interface{}) error {
	return func(value interface{}) error {
		switch v := value.(type) {
		case *int:
			if v != nil {
				return nil
			}
		case *bool:
			if v != nil {
				return nil
			}
		case *string:
			if v != nil {
				return nil
			}
		case *float64:
			if v != nil {
				return nil
			}
		default:
			return fmt.Errorf("validation failure: unknown type %T", value)
		}
		pattern := strings.Join(patterns, "|")
		err := error(nil)
		match := false
		if typeName != "" {
			match, err = regexp.MatchString(pattern, typeName)
			if match {
				return fmt.Errorf("is required if type is '%s'", typeName)
			}
		}
		return err
	}
}

func validateOrgServerFQDN(orgServerFQDN string) bool {
	_, fqdn, port, err := ParseOrgServerFQDN(orgServerFQDN)
	if err != nil || !govalidator.IsHost(*fqdn) || (port != nil && !govalidator.IsPort(*port)) {
		return false
	}
	return true
}

func ParseOrgServerFQDN(orgServerFQDN string) (*string, *string, *string, error) {
	originRegex := regexp.MustCompile(`^(https?)://([^:]+)(:(\d+))?$`)
	matches := originRegex.FindStringSubmatch(orgServerFQDN)
	if len(matches) == 0 {
		return nil, nil, nil, fmt.Errorf("unable to parse invalid orgServerFqdn: '%s'", orgServerFQDN)
	}

	protocol := strings.ToLower(matches[1])
	FQDN := matches[2]

	if len(protocol) == 0 || len(FQDN) == 0 {
		return nil, nil, nil, fmt.Errorf("empty Origin protocol or FQDN parsed from '%s'", orgServerFQDN)
	}

	var port *string
	if len(matches[4]) != 0 {
		port = &matches[4]
	}
	return &protocol, &FQDN, port, nil
}

func (ds *DeliveryServiceNullable) Sanitize() {
	if ds.GeoLimitCountries != nil {
		*ds.GeoLimitCountries = strings.ToUpper(strings.Replace(*ds.GeoLimitCountries, " ", "", -1))
	}
	if ds.ProfileID != nil && *ds.ProfileID == -1 {
		ds.ProfileID = nil
	}
	if ds.EdgeHeaderRewrite != nil && strings.TrimSpace(*ds.EdgeHeaderRewrite) == "" {
		ds.EdgeHeaderRewrite = nil
	}
	if ds.MidHeaderRewrite != nil && strings.TrimSpace(*ds.MidHeaderRewrite) == "" {
		ds.MidHeaderRewrite = nil
	}
	if ds.RoutingName == nil || *ds.RoutingName == "" {
		ds.RoutingName = util.StrPtr(DefaultRoutingName)
	}
	if ds.AnonymousBlockingEnabled == nil {
		ds.AnonymousBlockingEnabled = util.BoolPtr(false)
	}
	signedAlgorithm := SigningAlgorithmURLSig
	if ds.Signed && (ds.SigningAlgorithm == nil || *ds.SigningAlgorithm == "") {
		ds.SigningAlgorithm = &signedAlgorithm
	}
	if !ds.Signed && ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == signedAlgorithm {
		ds.Signed = true
	}
	if ds.MaxOriginConnections == nil || *ds.MaxOriginConnections < 0 {
		ds.MaxOriginConnections = util.IntPtr(0)
	}
	if ds.DeepCachingType == nil {
		s := DeepCachingType("")
		ds.DeepCachingType = &s
	}
	*ds.DeepCachingType = DeepCachingTypeFromString(string(*ds.DeepCachingType))
}

func (ds *DeliveryServiceNullable) validateTypeFields(tx *sql.Tx) error {
	// Validate the TypeName related fields below
	err := error(nil)
	DNSRegexType := "^DNS.*$"
	HTTPRegexType := "^HTTP.*$"
	SteeringRegexType := "^STEERING.*$"
	latitudeErr := "Must be a floating point number within the range +-90"
	longitudeErr := "Must be a floating point number within the range +-180"

	typeName, err := ValidateTypeID(tx, ds.TypeID, "deliveryservice")
	if err != nil {
		return err
	}

	errs := validation.Errors{
		"consistentHashQueryParams": validation.Validate(ds,
			validation.By(func(dsi interface{}) error {
				ds := dsi.(*DeliveryServiceNullable)
				if len(ds.ConsistentHashQueryParams) == 0 || DSType(typeName).IsHTTP() {
					return nil
				}
				return fmt.Errorf("consistentHashQueryParams not allowed for '%s' deliveryservice type", typeName)
			})),
		"initialDispersion": validation.Validate(ds.InitialDispersion,
			validation.By(requiredIfMatchesTypeName([]string{HTTPRegexType}, typeName)),
			validation.By(tovalidate.IsGreaterThanZero)),
		"ipv6RoutingEnabled": validation.Validate(ds.IPV6RoutingEnabled,
			validation.By(requiredIfMatchesTypeName([]string{SteeringRegexType, DNSRegexType, HTTPRegexType}, typeName))),
		"missLat": validation.Validate(ds.MissLat,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName)),
			validation.Min(-90.0).Error(latitudeErr),
			validation.Max(90.0).Error(latitudeErr)),
		"missLong": validation.Validate(ds.MissLong,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName)),
			validation.Min(-180.0).Error(longitudeErr),
			validation.Max(180.0).Error(longitudeErr)),
		"multiSiteOrigin": validation.Validate(ds.MultiSiteOrigin,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"orgServerFqdn": validation.Validate(ds.OrgServerFQDN,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName)),
			validation.NewStringRule(validateOrgServerFQDN, "must start with http:// or https:// and be followed by a valid hostname with an optional port (no trailing slash)")),
		"protocol": validation.Validate(ds.Protocol,
			validation.By(requiredIfMatchesTypeName([]string{SteeringRegexType, DNSRegexType, HTTPRegexType}, typeName))),
		"qstringIgnore": validation.Validate(ds.QStringIgnore,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
		"rangeRequestHandling": validation.Validate(ds.RangeRequestHandling,
			validation.By(requiredIfMatchesTypeName([]string{DNSRegexType, HTTPRegexType}, typeName))),
	}
	toErrs := tovalidate.ToErrors(errs)
	if len(toErrs) > 0 {
		return errors.New(util.JoinErrsStr(toErrs))
	}
	return nil
}

func (ds *DeliveryServiceNullable) Validate(tx *sql.Tx) error {
	ds.Sanitize()
	neverOrAlways := validation.NewStringRule(tovalidate.IsOneOfStringICase("NEVER", "ALWAYS"),
		"must be one of 'NEVER' or 'ALWAYS'")
	isDNSName := validation.NewStringRule(govalidator.IsDNSName, "must be a valid hostname")
	noPeriods := validation.NewStringRule(tovalidate.NoPeriods, "cannot contain periods")
	noSpaces := validation.NewStringRule(tovalidate.NoSpaces, "cannot contain spaces")
	errs := tovalidate.ToErrors(validation.Errors{
		"active":              validation.Validate(ds.Active, validation.NotNil),
		"cdnId":               validation.Validate(ds.CDNID, validation.Required),
		"deepCachingType":     validation.Validate(ds.DeepCachingType, neverOrAlways),
		"displayName":         validation.Validate(ds.DisplayName, validation.Required, validation.Length(1, 48)),
		"dscp":                validation.Validate(ds.DSCP, validation.NotNil, validation.Min(0)),
		"geoLimit":            validation.Validate(ds.GeoLimit, validation.NotNil),
		"geoProvider":         validation.Validate(ds.GeoProvider, validation.NotNil),
		"logsEnabled":         validation.Validate(ds.LogsEnabled, validation.NotNil),
		"regionalGeoBlocking": validation.Validate(ds.RegionalGeoBlocking, validation.NotNil),
		"routingName":         validation.Validate(ds.RoutingName, isDNSName, noPeriods, validation.Length(1, 48)),
		"typeId":              validation.Validate(ds.TypeID, validation.Required, validation.Min(1)),
		"xmlId":               validation.Validate(ds.XMLID, noSpaces, noPeriods, validation.Length(1, 48)),
	})
	if err := ds.validateTypeFields(tx); err != nil {
		errs = append(errs, errors.New("type fields: "+err.Error()))
	}
	if len(errs) == 0 {
		return nil
	}
	return util.JoinErrs(errs)
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
	Type      DSMatchType `json:"type"`
	SetNumber int         `json:"setNumber"`
	Pattern   string      `json:"pattern"`
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

type CachegroupPostDSReq struct {
	DeliveryServices []int64 `json:"deliveryServices"`
}

type CacheGroupPostDSResp struct {
	ID               util.JSONIntStr `json:"id"`
	ServerNames      []CacheName     `json:"serverNames"`
	DeliveryServices []int64         `json:"deliveryServices"`
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
