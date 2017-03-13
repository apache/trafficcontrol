package crconfig

import (
	"time"
)

// CRConfig is JSON-serializable as the CRConfig used by Traffic Control. However, it also contains diff timestamps, for the last update time of each field. These can be used to return only fields which have changed since a given time.
type CRConfig struct {
	Config               Config                       `json:"config,omitempty"`
	ConfigTime           time.Time                    `json:"-"`
	ContentServers       map[string]Server            `json:"contentServers,omitempty"`
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

type MatchSet struct {
	Protocol  string      `json:"protocol"`
	MatchList []MatchType `json:"matchlist"`
}

type MatchType struct {
	MatchType string `json:"match-type"`
	Regex     string `json:"regex"`
}

type Config struct {
	APICacheControlMaxAge                          *int      `json:"api.cache-control.max_age,string,omitempty"`
	APICacheControlMaxAgeLastTime                  time.Time `json:"-"`
	ConsistentDNSRouting                           *bool     `json:"consistent.dns.routing,string,omitempty"`
	ConsistentDNSRoutingTime                       time.Time `json:"-"`
	CoverageZonePollingIntervalSeconds             *int      `json:"coveragezone.polling.interval,string,omitempty"`
	CoverageZonePollingIntervalSecondsTime         time.Time `json:"-"`
	CoverageZonePollingURL                         *string   `json:"coveragezone.polling.url,omitempty"`
	CoverageZonePollingURLTime                     time.Time `json:"-"`
	DNSSecDynamicResponseExpiration                *string   `json:"dnssec.dynamic.response.expiration,omitempty"`
	DNSSecDynamicResponseExpirationTime            time.Time `json:"-"`
	DNSSecEnabled                                  *bool     `json:"dnssec.enabled,string,omitempty"`
	DNSSecEnabledTime                              time.Time `json:"-"`
	DomainName                                     *string   `json:"domain_name,omitempty"`
	DomainNameTime                                 time.Time `json:"-"`
	FederationMappingPollingIntervalSeconds        *int      `json:"federationmapping.polling.interval,string,omitempty"`
	FederationMappingPollingIntervalSecondsTime    time.Time `json:"-"`
	FederationMappingPollingURL                    *string   `json:"federationmapping.polling.url"`
	FederationMappingPollingURLTime                time.Time `json:"-"`
	GeoLocationPollingInterval                     *int      `json:"geolocation.polling.interval,string,omitempty"`
	GeoLocationPollingIntervalTime                 time.Time `json:"-"`
	GeoLocationPollingURL                          *string   `json:"geolocation.polling.url,omitempty"`
	GeoLocationPollingURLTime                      time.Time `json:"-"`
	KeyStoreMaintenanceIntervalSeconds             *int      `json:"keystore.maintenance.interval,string,omitempty"`
	KeyStoreMaintenanceIntervalSecondsTime         time.Time `json:"-"`
	NeustarPollingIntervalSeconds                  *int      `json:"neustar.polling.interval,string,omitempty"`
	NeustarPollingIntervalSecondsTime              time.Time `json:"-"`
	NeustarPollingURL                              *string   `json:"neustar.polling.url,omitempty"`
	NeustarPollingURLTime                          time.Time `json:"-"`
	SOA                                            *SOA      `json:"soa,omitempty"`
	SOATime                                        time.Time `json:"-"`
	DNSSecInceptionSeconds                         *int      `json:"dnssec.inception,string,omitempty"`
	DNSSecInceptionSecondsTime                     time.Time `json:"-"`
	Ttls                                           *TTL      `json:"ttls,omitempty"`
	TtlsTime                                       time.Time `json:"-"`
	Weight                                         *float64  `json:"weight,string,omitempty"`
	WeightTime                                     time.Time `json:"-"`
	ZoneManagerCacheMaintenanceIntervalSeconds     *int      `json:"zonemanager.cache.maintenance.interval,string,omitempty"`
	ZoneManagerCacheMaintenanceIntervalSecondsTime time.Time `json:"-"`
	ZoneManagerThreadpoolScale                     *float64  `json:"zonemanager.threadpool.scale,string,omitempty"`
	ZoneManagerThreadpoolScaleTime                 time.Time `json:"-"`
}

type SOA struct {
	Admin              *string `json:"admin,omitempty"`
	AdminTime          time.Time
	ExpireSeconds      *int `json:"expire,string,omitempty"`
	ExpireSecondsTime  time.Time
	MinimumSeconds     *int `json:"minimum,string,omitempty"`
	MinimumSecondsTime time.Time
	RefreshSeconds     *int `json:"refresh,string,omitempty"`
	RefreshSecondsTime time.Time
	RetrySeconds       *int `json:"retry,string,omitempty"`
	RetrySecondsTime   time.Time
}

type TTL struct {
	ASeconds          *int `json:"A,string,omitempty"`
	ASecondsTime      time.Time
	AAAASeconds       *int `json:"AAAA,string,omitempty"`
	AAAASecondsTime   time.Time
	DNSkeySeconds     *int `json:"DNSKEY,string,omitempty"`
	DNSkeySecondsTime time.Time
	DSSeconds         *int `json:"DS,string,omitempty"`
	DSSecondsTime     time.Time
	NSSeconds         *int `json:"NS,string,omitempty"`
	NSSecondsTime     time.Time
	SOASeconds        *int `json:"SOA,string,omitempty"`
	SOASecondsTime    time.Time
}

type Router struct {
	APIPort       *int `json:"apiPort,string,omitempty"`
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

type Server struct {
	CacheGroup           *string             `json:"cacheGroup,omitempty"`
	CacheGroupTime       time.Time           `json:"-"`
	DeliveryServices     map[string][]string `json:"deliveryServices,omitempty"`
	DeliveryServicesTime time.Time           `json:"-"`
	Fqdn                 *string             `json:"fqdn,omitempty"`
	FqdnTime             time.Time           `json:"-"`
	HashCount            *int                `json:"hashCount,omitempty"`
	HashCountTime        time.Time           `json:"-"`
	HashId               *string             `json:"hashId,omitempty"`
	HashIdTime           time.Time           `json:"-"`
	HttpsPort            *int                `json:"httpsPort,omitempty"`
	HttpsPortTime        time.Time           `json:"-"`
	InterfaceName        *string             `json:"interfaceName,omitempty"`
	InterfaceNameTime    time.Time           `json:"-"`
	Ip                   *string             `json:"ip,omitempty"`
	IpTime               time.Time           `json:"-"`
	Ip6                  *string             `json:"ip6,omitempty"`
	Ip6Time              time.Time           `json:"-"`
	LocationId           *string             `json:"locationId,omitempty"`
	LocationIdTime       time.Time           `json:"-"`
	Port                 *int                `json:"port,omitempty"`
	PortTime             time.Time           `json:"-"`
	Profile              *string             `json:"profile,omitempty"`
	ProfileTime          time.Time           `json:"-"`
	Status               *Status             `json:"status,omitempty"`
	StatusTime           time.Time           `json:"-"`
	ServerType           *string             `json:"type,omitempty"`
	ServerTypeTime       time.Time           `json:"-"`
}

type DeliveryService struct {
	CoverageZoneOnly        *bool                    `json:"coverageZoneOnly,string,omitempty"`
	CoverageZoneOnlyTime    time.Time                `json:"-"`
	Dispersion              *Dispersion              `json:"dispersion,omitempty"`
	DispersionTime          time.Time                `json:"-"`
	Domains                 []string                 `json:"domains,omitempty"`
	DomainsTime             time.Time                `json:"-"`
	GeoLocationProvider     *string                  `json:"geoLocationProvider,omitempty"`
	GeoLocationProviderTime time.Time                `json:"-"`
	MatchSets               []MatchSet               `json:"matchSets,omitempty"`
	MatchSetsTime           time.Time                `json:"-"`
	MissLocation            *LatLon                  `json:"missLocation,omitempty"`
	MissLocationTime        time.Time                `json:"-"`
	Protocol                *DeliveryServiceProtocol `json:"protocol,omitempty"`
	ProtocolTime            time.Time                `json:"-"`
	RegionalGeoBlocking     *bool                    `json:"regionalGeoBlocking,string,omitempty"`
	RegionalGeoBlockingTime time.Time                `json:"-"`
	ResponseHeaders         map[string]string        `json:"responseHeaders,omitempty"`
	ResponseHeadersTime     time.Time                `json:"-"`
	Soa                     *SOA                     `json:"soa,omitempty"`
	SoaTime                 time.Time                `json:"-"`
	SSLEnabled              *bool                    `json:"sslEnabled,string,omitempty"`
	SSLEnabledTime          time.Time                `json:"-"`
	TTL                     *int                     `json:"ttl,string,omitempty"`
	TTLTime                 time.Time                `json:"-"`
	TTLs                    *TTL                     `json:"ttls,omitempty"`
	TTLsTime                time.Time                `json:"-"`
}
type Dispersion struct {
	Limit    int  `json:"limit,omitempty"`
	Shuffled bool `json:"shuffled,string,omitempty"`
}
type LatLon struct {
	Lat float64 `json:"lat,string"`
	Lon float64 `json:"lon,string"`
}

type LatitudeLongitude struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type DeliveryServiceProtocol struct {
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
