package todb

import (
	"fmt"
	"gopkg.in/guregu/null.v3"
	"strings"
)

// use this view
// create view content_routers as select ip_address as ip, ip6_address as ip6, profile.name as profile, cachegroup.name as location,
// status.name as status, server.tcp_port as port, concat(server.host_name, ".", server.domain_name) as fqdn,
// parameter.value as apiport, cdn.name as cdnname
// from server
// join profile on profile.id = server.profile
// join profile_parameter on profile_parameter.profile = profile.id
// join parameter on parameter.id = profile_parameter.parameter
// join cachegroup on cachegroup.id = server.cachegroup
// join status on status.id = server.status
// join cdn on cdn.id = server.cdn_id
// join type on type.id = server.type
// where type.name = "CCR" and parameter.name="api.port";

type ContentRouter struct {
	Profile  string `db:"profile" json:"profile"`
	Apiport  int64  `db:"apiport" json:"api.port"`
	Location string `db:"location" json:"location"`
	Ip       string `db:"ip" json:"ip"`
	Status   string `db:"status" json:"status"`
	Ip6      string `db:"ip6" json:"ip6"`
	Port     int64  `db:"port" json:"port"`
	Fqdn     string `db:"fqdn" json:"fqdn"`
	Cdnname  string `db:"cdnname" json:"cdnname"`
}

// use this view
// create view monitors as select ip_address as ip, ip6_address as ip6, profile.name as profile, cachegroup.name as location,
// status.name as status, server.tcp_port as port, concat(server.host_name, ".", server.domain_name) as fqdn,
// cdn.name as cdnname
// from server
// join profile on profile.id = server.profile
// join cachegroup on cachegroup.id = server.cachegroup
// join status on status.id = server.status
// join cdn on cdn.id = server.cdn_id
// join type on type.id = server.type
// where type.name = "RASCAL";

type Monitor struct {
	Profile  string `db:"profile" json:"profile"`
	Location string `db:"location" json:"location"`
	Ip       string `db:"ip" json:"ip"`
	Status   string `db:"status" json:"status"`
	Ip6      string `db:"ip6" json:"ip6"`
	Port     int64  `db:"port" json:"port"`
	Fqdn     string `db:"fqdn" json:"fqdn"`
	Cdnname  string `db:"cdnname" json:"cdnname"`
}

type EdgeLocation struct {
	Name      string     `db:"name" json:"name"`
	Longitude null.Float `db:"longitude" json:"longitude"`
	Latitude  null.Float `db:"latitude" json:"latitude"`
}

// use this view
// create view crconfig_params as select distinct cdn.name as cdn_name, cdn.id as cdn_id,
// server.profile as profile_id,
// server.type as stype, parameter.name as pname,
// parameter.config_file as cfile, parameter.value as pvalue
// from server
// join cdn on cdn.id = server.cdn_id
// join profile on profile.id = server.profile
// join profile_parameter on profile_parameter.profile = server.profile
// join parameter on parameter.id = profile_parameter.parameter
// where server.type in (select id from type where name in ("EDGE", "MID", "CCR"))
// and parameter.config_file = 'CRConfig.json';
type CRConfigParam struct {
	CdnName        string `db:"cdn_name"`
	CdnId          int64  `db:"cdn_id"`
	ProfileId      int64  `db:"profile_id"`
	ServerType     int64  `db:"stype"`
	ParameterName  string `db:"pname"`
	ConfigFile     string `db:"cfile"`
	ParameterValue string `db:"pvalue"`
}

type Soa struct {
	Admin   string `json:"admin"`
	Expire  string `json:"expire"`
	Minimum string `json:"minimum"`
	Refresh string `json:"refresh"`
	Retry   string `json:"retry"`
}

type Ttls struct {
	A      string `json:"A"`
	AAAA   string `json:"AAAA"`
	DNSKEY string `json:"DNSKEY"`
	DS     string `json:"DS"`
	NS     string `json:"NS"`
	SOA    string `json:"SOA"`
}

type Config struct {
	// API_cache_control_max_age              string `json:"api.cache-control.max-age"`
	// Consistent_dns_routing                 string `json:"consistent.dns.routing"`
	// Coveragezone_polling_interval          string `json:"coveragezone.polling.interval"`
	// Coveragezone_polling_url               string `json:"coveragezone.polling.url"`
	// Dnssec_dynamic_response_expiration     string `json:"dnssec.dynamic.response.expiration"`
	// Dnssec_enabled                         string `json:"dnssec.enabled"`
	// DomainName                             string `json:"domain_name"`
	// Federationmapping_polling_interval     string `json:"federationmapping.polling.interval"`
	// Federationmapping_polling_url          string `json:"federationmapping.polling.url"`
	// Geolocation_polling_interval           string `json:"geolocation.polling.interval"`
	// Geolocation_polling_url                string `json:"geolocation.polling.url"`
	// Keystore_maintenance_interval          string `json:"keystore.maintenance.interval"`
	ParamMap map[string]string `json:"misc"`
	Soa      Soa               `json:"soa"`
	// Dnssec_inception                       string `json:"title-vi.dnssec.inception"`
	Ttls Ttls `json:"ttls"`
	// Weight                                 string `json:"weight"`
	// Zonemanager_cache_maintenance_interval string `json:"zonemanager.cache.maintenance.interval"`
	// Zonemanager_threadpool_scale           string `json:"zonemanager.threadpool.scale"`
}

type CRConfig struct {
	ContentRouters []ContentRouter `json:"contentRouters"`
	Monitors       []Monitor       `json:"monitors"`
	EdgeLocations  []EdgeLocation  `json:"edgeLocations"`
	Config         Config          `json:"config"`
}

func GetCRConfig(cdnName string) (interface{}, error) {

	crQuery := "select * from content_routers where cdnname=\"" + cdnName + "\""
	fmt.Println(crQuery)
	crs := []ContentRouter{}
	err := globalDB.Select(&crs, crQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	mQuery := "select * from monitors where cdnname=\"" + cdnName + "\""
	fmt.Println(mQuery)
	ms := []Monitor{}
	err = globalDB.Select(&ms, mQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	eQuery := "select name,longitude,latitude from cachegroup where type in (select id from type where name=\"EDGE_LOC\")"
	fmt.Println(eQuery)
	edges := []EdgeLocation{}
	err = globalDB.Select(&edges, eQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	pQuery := "select * from crconfig_params where cdn_name=\"" + cdnName + "\""
	fmt.Println(pQuery)
	params := []CRConfigParam{}
	err = globalDB.Select(&params, pQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pmap := make(map[string]string)
	miscMap := make(map[string]string)
	for _, row := range params {
		pmap[row.ParameterName] = row.ParameterValue
		if !strings.HasPrefix(row.ParameterName, "tld.") {
			miscMap[row.ParameterName] = row.ParameterValue
		}
	}
	cfg := Config{
		// API_cache_control_max_age:          pmap["api.cache-control.max-age"],
		// Consistent_dns_routing:             pmap[""],
		// Coveragezone_polling_interval:      pmap[""],
		// Coveragezone_polling_url:           pmap["coveragezone.polling.url"],
		// Dnssec_dynamic_response_expiration: pmap[""],
		// Dnssec_enabled:                     pmap[""],
		// DomainName:                         pmap[""],
		// Federationmapping_polling_interval: pmap[""],
		// Federationmapping_polling_url:      pmap[""],
		// Geolocation_polling_interval:       pmap[""],
		// Geolocation_polling_url:            pmap[""],
		// Keystore_maintenance_interval:      pmap[""],
		ParamMap: miscMap,
		Soa: Soa{
			Admin:   pmap["tld.soa.admin"],
			Expire:  pmap["tld.soa.expire"],
			Minimum: pmap["tld.soa.minimum"],
			Refresh: pmap["tld.soa.refresh"],
			Retry:   pmap["tld.soa.retry"],
		},
		// Dnssec_inception: "",
		Ttls: Ttls{
			A:      pmap["tld.ttl.A"],
			AAAA:   pmap["tld.ttl.AAA"],
			DNSKEY: pmap["tld.ttl.DNSKEY"],
			DS:     pmap["tld.ttl.DS"],
			NS:     pmap["tld.ttl.NS"],
			SOA:    pmap["tld.ttl.SOA"],
		},
		// Weight: "",
		// Zonemanager_cache_maintenance_interval: "",
		// Zonemanager_threadpool_scale:           "",
	}

	return CRConfig{
		ContentRouters: crs,
		Monitors:       ms,
		EdgeLocations:  edges,
		Config:         cfg,
	}, nil
}
