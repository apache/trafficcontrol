/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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

package client

// GetDeliveryServiceResponse ...
type GetDeliveryServiceResponse struct {
	Version  string            `json:"version"`
	Response []DeliveryService `json:"response"`
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
	ID                   int                    `json:"id"`
	XMLID                string                 `json:"xmlId"`
	Active               bool                   `json:"active"`
	DSCP                 int                    `json:"dscp"`
	Signed               bool                   `json:"signed"`
	QStringIgnore        int                    `json:"qstringIgnore"`
	GeoLimit             int                    `json:"geoLimit"`
	GeoProvider          int                    `json:"geoProvider"`
	HTTPBypassFQDN       string                 `json:"httpBypassFqdn"`
	DNSBypassIP          string                 `json:"dnsBypassIp"`
	DNSBypassIP6         string                 `json:"dnsBypassIp6"`
	DNSBypassCname       string                 `json:"dnsBypassCname"`
	DNSBypassTTL         int                    `json:"dnsBypassTtl"`
	OrgServerFQDN        string                 `json:"orgServerFqdn"`
	Type                 string                 `json:"type"`
	ProfileName          string                 `json:"profileName"`
	ProfileDesc          string                 `json:"profileDescription"`
	CDNName              string                 `json:"cdnName"`
	CCRDNSTTL            int                    `json:"ccrDnsTtl"`
	GlobalMaxMBPS        int                    `json:"globalMaxMbps"`
	GlobalMaxTPS         int                    `json:"globalMaxTps"`
	LongDesc             string                 `json:"longDesc"`
	LongDesc1            string                 `json:"longDesc1"`
	LongDesc2            string                 `json:"longDesc2"`
	MaxDNSAnswers        int                    `json:"maxDnsAnswers"`
	InfoURL              string                 `json:"infoUrl"`
	MissLat              float64                `json:"missLat"`
	MissLong             float64                `json:"missLong"`
	CheckPath            string                 `json:"checkPath"`
	LastUpdated          string                 `json:"lastUpdated"`
	Protocol             int                    `json:"protocol"`
	IPV6RoutingEnabled   bool                   `json:"ipv6RoutingEnabled"`
	RangeRequestHandling int                    `json:"rangeRequestHandling"`
	HeaderRewrite        string                 `json:"headerRewrite"`
	EdgeHeaderRewrite    string                 `json:"edgeHeaderRewrite"`
	MidHeaderRewrite     string                 `json:"midHeaderRewrite"`
	TRResponseHeaders    string                 `json:"trResponseHeaders"`
	RegexRemap           string                 `json:"regexRemap"`
	CacheURL             string                 `json:"cacheurl"`
	RemapText            string                 `json:"remapText"`
	MultiSiteOrigin      int                    `json:"multiSiteOrigin"`
	DisplayName          string                 `json:"displayName"`
	InitialDispersion    int                    `json:"initialDispersion"`
	MatchList            []DeliveryServiceMatch `json:"matchList,omitempty"`
}

// DeliveryServiceMatch ...
type DeliveryServiceMatch struct {
	Type      string `json:"type"`
	SetNumber string `json:"setNumber"`
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
	Page     int                     `json:"page"`
	OrderBy  string                  `json:"orderby"`
	Limit    int                     `json:"limit"`
}

// DeliveryServiceServer ...
type DeliveryServiceServer struct {
	LastUpdated     string `json:"lastUpdated"`
	Server          string `json:"server"`
	DeliveryService string `json:"deliveryService"`
}

// DeliveryServiceSSLKeysResponse ...
type DeliveryServiceSSLKeysResponse struct {
	Response DeliveryServiceSSLKeys `json:"response"`
}

// DeliveryServiceSSLKeys ...
type DeliveryServiceSSLKeys struct {
	CDN             string                            `json:"cdn"`
	DeliveryService string                            `json:"DeliveryService"`
	BusinessUnit    string                            `json:"businessUnit"`
	City            string                            `json:"city"`
	Organization    string                            `json:"organization"`
	Hostname        string                            `json:"hostname"`
	Country         string                            `json:"country"`
	State           string                            `json:"state"`
	Version         string                            `json:"version"`
	Certificate     DeliveryServiceSSLKeysCertificate `json:"certificate"`
}

// DeliveryServiceSSLKeysCertificate ...
type DeliveryServiceSSLKeysCertificate struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
	CSR string `json:"csr"`
}
