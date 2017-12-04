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

// CRConfig is JSON-serializable as the CRConfig used by Traffic Control.
type CRConfig struct {
	Config           CRConfigConfig                       `json:"config,omitempty"`
	ContentServers   map[string]CRConfigTrafficOpsServer  `json:"contentServers,omitempty"`
	ContentRouters   map[string]CRConfigRouter            `json:"contentRouters,omitempty"`
	DeliveryServices map[string]CRConfigDeliveryService   `json:"deliveryServices,omitempty"`
	EdgeLocations    map[string]CRConfigLatitudeLongitude `json:"edgeLocations,omitempty"`
	Monitors         map[string]CRConfigMonitor           `json:"monitors,omitempty"`
	Stats            CRConfigStats                        `json:"stats,omitempty"`
}

type CRConfigConfig struct {
	APICacheControlMaxAge                      *string      `json:"api.cache-control.max_age,omitempty"`
	ConsistentDNSRouting                       *string      `json:"consistent.dns.routing,omitempty"`
	CoverageZonePollingIntervalSeconds         *string      `json:"coveragezone.polling.interval,omitempty"`
	CoverageZonePollingURL                     *string      `json:"coveragezone.polling.url,omitempty"`
	DNSSecDynamicResponseExpiration            *string      `json:"dnssec.dynamic.response.expiration,omitempty"`
	DNSSecEnabled                              *string      `json:"dnssec.enabled,omitempty"`
	DomainName                                 *string      `json:"domain_name,omitempty"`
	FederationMappingPollingIntervalSeconds    *string      `json:"federationmapping.polling.interval,omitempty"`
	FederationMappingPollingURL                *string      `json:"federationmapping.polling.url"`
	GeoLocationPollingInterval                 *string      `json:"geolocation.polling.interval,omitempty"`
	GeoLocationPollingURL                      *string      `json:"geolocation.polling.url,omitempty"`
	KeyStoreMaintenanceIntervalSeconds         *string      `json:"keystore.maintenance.interval,omitempty"`
	NeustarPollingIntervalSeconds              *string      `json:"neustar.polling.interval,omitempty"`
	NeustarPollingURL                          *string      `json:"neustar.polling.url,omitempty"`
	SOA                                        *SOA         `json:"soa,omitempty"`
	DNSSecInceptionSeconds                     *string      `json:"dnssec.inception,omitempty"`
	Ttls                                       *CRConfigTTL `json:"ttls,omitempty"`
	Weight                                     *string      `json:"weight,omitempty"`
	ZoneManagerCacheMaintenanceIntervalSeconds *string      `json:"zonemanager.cache.maintenance.interval,omitempty"`
	ZoneManagerThreadpoolScale                 *string      `json:"zonemanager.threadpool.scale,omitempty"`
}

type CRConfigTTL struct {
	ASeconds      *string `json:"A,omitempty"`
	AAAASeconds   *string `json:"AAAA,omitempty"`
	DNSkeySeconds *string `json:"DNSKEY,omitempty"`
	DSSeconds     *string `json:"DS,omitempty"`
	NSSeconds     *string `json:"NS,omitempty"`
	SOASeconds    *string `json:"SOA,omitempty"`
}

type CRConfigRouterStatus string

type CRConfigRouter struct {
	APIPort      *string               `json:"apiPort,omitempty"`
	FQDN         *string               `json:"fqdn,omitempty"`
	HTTPSPort    *int                  `json:"httpsPort,omitempty"`
	IP           *string               `json:"ip,omitempty"`
	IP6          *string               `json:"ip6,omitempty"`
	Location     *string               `json:"location,omitempty"`
	Port         *int                  `json:"port,omitempty"`
	Profile      *string               `json:"profile,omitempty"`
	ServerStatus *CRConfigRouterStatus `json:"status,omitempty"`
}

type CRConfigServerStatus string

type CRConfigTrafficOpsServer struct {
	CacheGroup       *string               `json:"cacheGroup,omitempty"`
	Fqdn             *string               `json:"fqdn,omitempty"`
	HashCount        *int                  `json:"hashCount,omitempty"`
	HashId           *string               `json:"hashId,omitempty"`
	HttpsPort        *int                  `json:"httpsPort,omitempty"`
	InterfaceName    *string               `json:"interfaceName,omitempty"`
	Ip               *string               `json:"ip,omitempty"`
	Ip6              *string               `json:"ip6,omitempty"`
	LocationId       *string               `json:"locationId,omitempty"`
	Port             *int                  `json:"port,omitempty"`
	Profile          *string               `json:"profile,omitempty"`
	ServerStatus     *CRConfigServerStatus `json:"status,omitempty"`
	ServerType       *string               `json:"type,omitempty"`
	DeliveryServices map[string][]string   `json:"deliveryServices,omitempty"`
}

//TODO: drichardson - reconcile this with the DeliveryService struct in deliveryservices.go
type CRConfigDeliveryService struct {
	CoverageZoneOnly    *string                          `json:"coverageZoneOnly,omitempty"`
	Dispersion          *CRConfigDispersion              `json:"dispersion,omitempty"`
	Domains             []string                         `json:"domains,omitempty"`
	GeoLocationProvider *string                          `json:"geoLocationProvider,omitempty"`
	MatchSets           []MatchSet                       `json:"matchSets,omitempty"`
	MissLocation        *CRConfigLatitudeLongitude       `json:"missLocation,omitempty"`
	Protocol            *CRConfigDeliveryServiceProtocol `json:"protocol,omitempty"`
	RegionalGeoBlocking *string                          `json:"regionalGeoBlocking,omitempty"`
	ResponseHeaders     map[string]string                `json:"responseHeaders,omitempty"`
	Soa                 *SOA                             `json:"soa,omitempty"`
	SSLEnabled          *string                          `json:"sslEnabled,omitempty"`
	TTL                 *int                             `json:"ttl,omitempty"`
	TTLs                *CRConfigTTL                     `json:"ttls,omitempty"`
}

type CRConfigDispersion struct {
	Limit    int     `json:"limit,omitempty"`
	Shuffled *string `json:"shuffled,omitempty"`
}

type CRConfigLatitudeLongitude struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type CRConfigDeliveryServiceProtocol struct {
	AcceptHTTP      bool `json:"acceptHttp,string,omitempty"`
	AcceptHTTPS     bool `json:"acceptHttps,string,omitempty"`
	RedirectOnHTTPS bool `json:"redirectOnHttps,string,omitempty"`
}

type CRConfigMonitor struct {
	FQDN         *string               `json:"fqdn,omitempty"`
	HTTPSPort    *int                  `json:"httpsPort,omitempty"`
	IP           *string               `json:"ip,omitempty"`
	IP6          *string               `json:"ip6,omitempty"`
	Location     *string               `json:"location,omitempty"`
	Port         *int                  `json:"port,omitempty"`
	Profile      *string               `json:"profile,omitempty"`
	ServerStatus *CRConfigServerStatus `json:"status,omitempty"`
}

type CRConfigStats struct {
	CDNName         *string `json:"CDN_name,omitempty"`
	DateUnixSeconds *int64  `json:"date,omitempty"`
	TMHost          *string `json:"tm_host,omitempty"`
	TMPath          *string `json:"tm_path,omitempty"`
	TMUser          *string `json:"tm_user,omitempty"`
	TMVersion       *string `json:"tm_version,omitempty"`
}
