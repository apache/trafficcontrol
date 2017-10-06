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

import (
	"time"
)

// CRConfig is JSON-serializable as the CRConfig used by Traffic Control. However, it also contains diff timestamps, for the last update time of each field. These can be used to return only fields which have changed since a given time.
type CRConfig struct {
	Config               Config                       `json:"config,omitempty"`
	ConfigTime           time.Time                    `json:"-"`
	ContentServers       map[string]TrafficOpsServer  `json:"contentServers,omitempty"`
	ContentServersTime   map[string]time.Time         `json:"-"`
	ContentRouters       map[string]Router            `json:"contentRouters,omitempty"`
	ContentRoutersTime   map[string]time.Time         `json:"-"`
	DeliveryServices     map[string]DeliveryService   `json:"deliveryServices,omitempty"`
	DeliveryServicesTime map[string]time.Time         `json:"-"`
	EdgeLocations        map[string]LatitudeLongitude `json:"edgeLocations,omitempty"`
	EdgeLocationsTime    map[string]time.Time         `json:"-"`
	Monitors             map[string]Monitor           `json:"monitors,omitempty"`
	MonitorsTime         map[string]time.Time         `json:"-"`
	Stats                Stats                        `json:"stats,omitempty"`
	StatsTime            time.Time                    `json:"-"`
}

type Config struct {
	APICacheControlMaxAge                          *string   `json:"api.cache-control.max_age,omitempty"`
	APICacheControlMaxAgeLastTime                  time.Time `json:"-"`
	ConsistentDNSRouting                           *string   `json:"consistent.dns.routing,omitempty"`
	ConsistentDNSRoutingTime                       time.Time `json:"-"`
	CoverageZonePollingIntervalSeconds             *string   `json:"coveragezone.polling.interval,omitempty"`
	CoverageZonePollingIntervalSecondsTime         time.Time `json:"-"`
	CoverageZonePollingURL                         *string   `json:"coveragezone.polling.url,omitempty"`
	CoverageZonePollingURLTime                     time.Time `json:"-"`
	DNSSecDynamicResponseExpiration                *string   `json:"dnssec.dynamic.response.expiration,omitempty"`
	DNSSecDynamicResponseExpirationTime            time.Time `json:"-"`
	DNSSecEnabled                                  *string   `json:"dnssec.enabled,omitempty"`
	DNSSecEnabledTime                              time.Time `json:"-"`
	DomainName                                     *string   `json:"domain_name,omitempty"`
	DomainNameTime                                 time.Time `json:"-"`
	FederationMappingPollingIntervalSeconds        *string   `json:"federationmapping.polling.interval,omitempty"`
	FederationMappingPollingIntervalSecondsTime    time.Time `json:"-"`
	FederationMappingPollingURL                    *string   `json:"federationmapping.polling.url"`
	FederationMappingPollingURLTime                time.Time `json:"-"`
	GeoLocationPollingInterval                     *string   `json:"geolocation.polling.interval,omitempty"`
	GeoLocationPollingIntervalTime                 time.Time `json:"-"`
	GeoLocationPollingURL                          *string   `json:"geolocation.polling.url,omitempty"`
	GeoLocationPollingURLTime                      time.Time `json:"-"`
	KeyStoreMaintenanceIntervalSeconds             *string   `json:"keystore.maintenance.interval,omitempty"`
	KeyStoreMaintenanceIntervalSecondsTime         time.Time `json:"-"`
	NeustarPollingIntervalSeconds                  *string   `json:"neustar.polling.interval,omitempty"`
	NeustarPollingIntervalSecondsTime              time.Time `json:"-"`
	NeustarPollingURL                              *string   `json:"neustar.polling.url,omitempty"`
	NeustarPollingURLTime                          time.Time `json:"-"`
	SOA                                            *SOA      `json:"soa,omitempty"`
	SOATime                                        time.Time `json:"-"`
	DNSSecInceptionSeconds                         *string   `json:"dnssec.inception,omitempty"`
	DNSSecInceptionSecondsTime                     time.Time `json:"-"`
	Ttls                                           *TTL      `json:"ttls,omitempty"`
	TtlsTime                                       time.Time `json:"-"`
	Weight                                         *string   `json:"weight,omitempty"`
	WeightTime                                     time.Time `json:"-"`
	ZoneManagerCacheMaintenanceIntervalSeconds     *string   `json:"zonemanager.cache.maintenance.interval,omitempty"`
	ZoneManagerCacheMaintenanceIntervalSecondsTime time.Time `json:"-"`
	ZoneManagerThreadpoolScale                     *string   `json:"zonemanager.threadpool.scale,omitempty"`
	ZoneManagerThreadpoolScaleTime                 time.Time `json:"-"`
}

type SOA struct {
	Admin              *string `json:"admin,omitempty"`
	AdminTime          time.Time
	ExpireSeconds      *string `json:"expire,omitempty"`
	ExpireSecondsTime  time.Time
	MinimumSeconds     *string `json:"minimum,omitempty"`
	MinimumSecondsTime time.Time
	RefreshSeconds     *string `json:"refresh,omitempty"`
	RefreshSecondsTime time.Time
	RetrySeconds       *string `json:"retry,omitempty"`
	RetrySecondsTime   time.Time
}

type TTL struct {
	ASeconds          *string `json:"A,omitempty"`
	ASecondsTime      time.Time
	AAAASeconds       *string `json:"AAAA,omitempty"`
	AAAASecondsTime   time.Time
	DNSkeySeconds     *string `json:"DNSKEY,omitempty"`
	DNSkeySecondsTime time.Time
	DSSeconds         *string `json:"DS,omitempty"`
	DSSecondsTime     time.Time
	NSSeconds         *string `json:"NS,omitempty"`
	NSSecondsTime     time.Time
	SOASeconds        *string `json:"SOA,omitempty"`
	SOASecondsTime    time.Time
}

type Router struct {
	APIPort       *string `json:"apiPort,omitempty"`
	APIPortTime   time.Time
	FQDN          *string `json:"fqdn,omitempty"`
	FQDNTime      time.Time
	HTTPSPort     *int `json:"httpsPort,omitempty"`
	HTTPSPortTime time.Time
	IP            *string `json:"ip,omitempty"`
	IPTime        time.Time
	IP6           *string `json:"ip6,omitempty"`
	IP6Time       time.Time
	Location      *string `json:"location,omitempty"`
	LocationTime  time.Time
	Port          *int `json:"port,omitempty"`
	PortTime      time.Time
	Profile       *string `json:"profile,omitempty"`
	ProfileTime   time.Time
	Status        *Status `json:"status,omitempty"`
	StatusTime    time.Time
}

type Status string

type TrafficOpsServer struct {
	CacheGroup        *string   `json:"cacheGroup,omitempty"`
	CacheGroupTime    time.Time `json:"-"`
	Fqdn              *string   `json:"fqdn,omitempty"`
	FqdnTime          time.Time `json:"-"`
	HashCount         *int      `json:"hashCount,omitempty"`
	HashCountTime     time.Time `json:"-"`
	HashId            *string   `json:"hashId,omitempty"`
	HashIdTime        time.Time `json:"-"`
	HttpsPort         *int      `json:"httpsPort,omitempty"`
	HttpsPortTime     time.Time `json:"-"`
	InterfaceName     *string   `json:"interfaceName,omitempty"`
	InterfaceNameTime time.Time `json:"-"`
	Ip                *string   `json:"ip,omitempty"`
	IpTime            time.Time `json:"-"`
	Ip6               *string   `json:"ip6,omitempty"`
	Ip6Time           time.Time `json:"-"`
	LocationId        *string   `json:"locationId,omitempty"`
	LocationIdTime    time.Time `json:"-"`
	Port              *int      `json:"port,omitempty"`
	PortTime          time.Time `json:"-"`
	Profile           *string   `json:"profile,omitempty"`
	ProfileTime       time.Time `json:"-"`
	Status            *Status   `json:"status,omitempty"`
	StatusTime        time.Time `json:"-"`
	ServerType        *string   `json:"type,omitempty"`
	ServerTypeTime    time.Time `json:"-"`
}

type DeliveryService struct {
	CoverageZoneOnly        *string                  `json:"coverageZoneOnly,omitempty"`
	CoverageZoneOnlyTime    time.Time                `json:"-"`
	Dispersion              *Dispersion              `json:"dispersion,omitempty"`
	DispersionTime          time.Time                `json:"-"`
	Domains                 []string                 `json:"domains,omitempty"`
	DomainsTime             time.Time                `json:"-"`
	GeoLocationProvider     *string                  `json:"geoLocationProvider,omitempty"`
	GeoLocationProviderTime time.Time                `json:"-"`
	MatchSets               []MatchSet               `json:"matchSets,omitempty"`
	MatchSetsTime           time.Time                `json:"-"`
	MissLocation            *LatitudeLongitude       `json:"missLocation,omitempty"`
	MissLocationTime        time.Time                `json:"-"`
	Protocol                *DeliveryServiceProtocol `json:"protocol,omitempty"`
	ProtocolTime            time.Time                `json:"-"`
	RegionalGeoBlocking     *string                  `json:"regionalGeoBlocking,omitempty"`
	RegionalGeoBlockingTime time.Time                `json:"-"`
	ResponseHeaders         map[string]string        `json:"responseHeaders,omitempty"`
	ResponseHeadersTime     time.Time                `json:"-"`
	Soa                     *SOA                     `json:"soa,omitempty"`
	SoaTime                 time.Time                `json:"-"`
	SSLEnabled              *string                  `json:"sslEnabled,omitempty"`
	SSLEnabledTime          time.Time                `json:"-"`
	TTL                     *int                     `json:"ttl,omitempty"`
	TTLTime                 time.Time                `json:"-"`
	TTLs                    *TTL                     `json:"ttls,omitempty"`
	TTLsTime                time.Time                `json:"-"`
}

// MissLocation ...
type MissLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitudef"`
}

// MatchSet ...
type MatchSet struct {
	Protocol  string      `json:"protocol"`
	MatchList []MatchList `json:"matchList"`
}

// MatchList ...
type MatchList struct {
	Regex     string `json:"regex"`
	MatchType string `json:"matchType"`
}

// BypassDestination ...
type BypassDestination struct {
	FQDN string `json:"fqdn"`
	Type string `json:"type"`
	Port int    `json:"Port"`
}

// TTLS ...
type TTLS struct {
	Arecord    int `json:"A"`
	SoaRecord  int `json:"SOA"`
	NsRecord   int `json:"NS"`
	AaaaRecord int `json:"AAAA"`
}

type Dispersion struct {
	Limit    int     `json:"limit,omitempty"`
	Shuffled *string `json:"shuffled,omitempty"`
}

type LatitudeLongitude struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type DeliveryServiceProtocol struct {
	AcceptHTTP      bool `json:"acceptHttp,string,omitempty"`
	AcceptHTTPS     bool `json:"acceptHttps,string,omitempty"`
	RedirectOnHTTPS bool `json:"redirectOnHttps,string,omitempty"`
}

type Monitor struct {
	FQDN          *string   `json:"fqdn,omitempty"`
	FQDNTime      time.Time `json:"-"`
	HTTPSPort     *int      `json:"httpsPort,omitempty"`
	HTTPSPortTime time.Time `json:"-"`
	IP            *string   `json:"ip,omitempty"`
	IPTime        time.Time `json:"-"`
	IP6           *string   `json:"ip6,omitempty"`
	IP6Time       time.Time `json:"-"`
	Location      *string   `json:"location,omitempty"`
	LocationTime  time.Time `json:"-"`
	Port          *int      `json:"port,omitempty"`
	PortTime      time.Time `json:"-"`
	Profile       *string   `json:"profile,omitempty"`
	ProfileTime   time.Time `json:"-"`
	Status        *Status   `json:"status,omitempty"`
	StatusTime    time.Time `json:"-"`
}

type Stats struct {
	CDNName             *string   `json:"CDN_name,omitempty"`
	CDNNameTime         time.Time `json:"-"`
	DateUnixSeconds     *int64    `json:"date,omitempty"`
	DateUnixSecondsTime time.Time `json:"-"`
	TMHost              *string   `json:"tm_host,omitempty"`
	TMHostTime          time.Time `json:"-"`
	TMPath              *string   `json:"tm_path,omitempty"`
	TMPathTime          time.Time `json:"-"`
	TMUser              *string   `json:"tm_user,omitempty"`
	TMUserTime          time.Time `json:"-"`
	TMVersion           *string   `json:"tm_version,omitempty"`
	TMVersionTime       time.Time `json:"-"`
}

// TrafficRouter ...
type TrafficRouter struct {
	Port     int    `json:"port"`
	IP6      string `json:"ip6"`
	IP       string `json:"ip"`
	FQDN     string `json:"fqdn"`
	Profile  string `json:"profile"`
	Location string `json:"location"`
	Status   string `json:"status"`
	APIPort  int    `json:"apiPort"`
}

// TrafficRouterConfig is the json unmarshalled without any changes
// note all structs are local to this file _except_ the TrafficRouterConfig struct.
type TrafficRouterConfig struct {
	TrafficServers   []TrafficServer        `json:"trafficServers,omitempty"`
	TrafficMonitors  []TrafficMonitor       `json:"trafficMonitors,omitempty"`
	TrafficRouters   []TrafficRouter        `json:"trafficRouters,omitempty"`
	CacheGroups      []TMCacheGroup         `json:"cacheGroups,omitempty"`
	DeliveryServices []TRDeliveryService    `json:"deliveryServices,omitempty"`
	Stats            map[string]interface{} `json:"stats,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
}

// TrafficRouterConfigMap ...
type TrafficRouterConfigMap struct {
	TrafficServer   map[string]TrafficServer
	TrafficMonitor  map[string]TrafficMonitor
	TrafficRouter   map[string]TrafficRouter
	CacheGroup      map[string]TMCacheGroup
	DeliveryService map[string]TRDeliveryService
	Config          map[string]interface{}
	Stat            map[string]interface{}
}

// TRDeliveryService ...
// TODO JvD: move to deliveryservice.go ??
type TRDeliveryService struct {
	XMLID             string            `json:"xmlId"`
	Domains           []string          `json:"domains"`
	RoutingName       string            `json:"routingName"`
	MissLocation      MissLocation      `json:"missCoordinates"`
	CoverageZoneOnly  bool              `json:"coverageZoneOnly"`
	MatchSets         []MatchSet        `json:"matchSets"`
	TTL               int               `json:"ttl"`
	TTLs              TTLS              `json:"ttls"`
	BypassDestination BypassDestination `json:"bypassDestination"`
	StatcDNSEntries   []StaticDNS       `json:"statitDnsEntries"`
	Soa               SOA               `json:"soa"`
}

// StaticDNS ...
type StaticDNS struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

// Coordinates ...
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TrafficServer ...
type TrafficServer struct {
	Profile          string              `json:"profile"`
	IP               string              `json:"ip"`
	Status           string              `json:"status"`
	CacheGroup       string              `json:"cacheGroup"`
	IP6              string              `json:"ip6"`
	Port             int                 `json:"port"`
	HostName         string              `json:"hostName"`
	FQDN             string              `json:"fqdn"`
	InterfaceName    string              `json:"interfaceName"`
	Type             string              `json:"type"`
	HashID           string              `json:"hashId"`
	DeliveryServices []tsdeliveryService `json:"deliveryServices,omitempty"` // the deliveryServices key does not exist on mids
}

type tsdeliveryService struct {
	Xmlid  string   `json:"xmlId"`
	Remaps []string `json:"remaps"`
}
