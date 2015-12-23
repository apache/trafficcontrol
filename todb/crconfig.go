package todb

import (
	"fmt"
	"gopkg.in/guregu/null.v3"
	"reflect"
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
	Profile  string      `db:"profile" json:"profile"`
	Location string      `db:"location" json:"location"`
	Ip       string      `db:"ip" json:"ip"`
	Status   string      `db:"status" json:"status"`
	Ip6      null.String `db:"ip6" json:"ip6"`
	Port     int64       `db:"port" json:"port"`
	Fqdn     string      `db:"fqdn" json:"fqdn"`
	Cdnname  string      `db:"cdnname" json:"cdnname"`
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
	ParamMap map[string]string `json:"misc"`
	Soa      Soa               `json:"soa"`
	Ttls     Ttls              `json:"ttls"`
}

type MactchListEntry struct {
	MatchType string `json:"match-type"`
	Regex     string `json:"regex"`
}

type MatchSetEntry struct {
	Matchlist []MactchListEntry `json:"matchlist"`
	Protocol  string            `json:"protocol"`
}

type MissLocation struct {
	Longitude null.Float `db:"longitude" json:"long"`
	Latitude  null.Float `db:"latitude" json:"lat"`
}

type Dispersion struct {
	Limit    int      `json:"limit"`
	Shuffled null.Int `json:"shuffled"`
}

type CrDeliveryService struct {
	CoverageZoneOnly     string            `json:"coverageZoneOnly"`
	Domains              []string          `json:"domains"`
	IP6RoutingEnabled    string            `json:"ip6RoutingEnabled"`
	MatchSets            []MatchSetEntry   `json:"matchsets"`
	MaxDNSIpsForLocation null.Int          `json:"maxDnsIpsForLocation"`
	MissLocation         MissLocation      `json:"missLocation"`
	Soa                  Soa               `json:"soa"`
	TTL                  null.Int          `json:"ttl"`
	Ttls                 Ttls              `json:"ttls"`
	ResponseHeaders      map[string]string `json:"responseHeaders"`
	Dispersion           Dispersion        `json:"dispersion"`
	GeoEnabled           map[string]string `json:"geoEnabled"`
}

// create view crconfig_ds_data as select xml_id, profile, ccr_dns_ttl, global_max_mbps, global_max_tps,
// max_dns_answers, miss_lat, miss_long, protocoltype.name as protocol, ipv6_routing_enabled,
// tr_request_headers, tr_response_headers, initial_dispersion, dns_bypass_cname,
// dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, cdn.name as cdn_name,
// regex.pattern as match_pattern, regextype.name as match_type
// from deliveryservice
// join cdn on cdn.id = deliveryservice.cdn_id
// join deliveryservice_regex on deliveryservice_regex.deliveryservice = deliveryservice.id
// join regex on regex.id = deliveryservice_regex.regex
// join type as protocoltype on protocoltype.id = deliveryservice.type
// join type as regextype on regextype.id = regex.type;

type CrconfigDsData struct {
	XmlId              string      `db:"xml_id" json:"xmlId"`
	Profile            int64       `db:"profile" json:"profile"`
	CcrDnsTtl          null.Int    `db:"ccr_dns_ttl" json:"ccrDnsTtl"`
	GlobalMaxMbps      null.Int    `db:"global_max_mbps" json:"globalMaxMbps"`
	GlobalMaxTps       null.Int    `db:"global_max_tps" json:"globalMaxTps"`
	MaxDnsAnswers      null.Int    `db:"max_dns_answers" json:"maxDnsAnswers"`
	MissLat            null.Float  `db:"miss_lat" json:"missLat"`
	MissLong           null.Float  `db:"miss_long" json:"missLong"`
	Protocol           string      `db:"protocol" json:"protocol"`
	Ipv6RoutingEnabled null.Int    `db:"ipv6_routing_enabled" json:"ipv6RoutingEnabled"`
	TrRequestHeaders   null.String `db:"tr_request_headers" json:"trRequestHeaders"`
	TrResponseHeaders  null.String `db:"tr_response_headers" json:"trResponseHeaders"`
	InitialDispersion  null.Int    `db:"initial_dispersion" json:"initialDispersion"`
	DnsBypassCname     null.String `db:"dns_bypass_cname" json:"dnsBypassCname"`
	DnsBypassIp        null.String `db:"dns_bypass_ip" json:"dnsBypassIp"`
	DnsBypassIp6       null.String `db:"dns_bypass_ip6" json:"dnsBypassIp6"`
	DnsBypassTtl       null.Int    `db:"dns_bypass_ttl" json:"dnsBypassTtl"`
	GeoLimit           int64       `db:"geo_limit" json:"geoLimit"`
	CdnName            null.String `db:"cdn_name" json:"cdnName"`
	MatchPattern       string      `db:"match_pattern" json:"matchPattern"`
	MatchType          string      `db:"match_type" json:"matchType"`
}

type CRConfig struct {
	ContentRouters   []ContentRouter              `json:"contentRouters"`
	Monitors         []Monitor                    `json:"monitors"`
	EdgeLocations    []EdgeLocation               `json:"edgeLocations"`
	Config           Config                       `json:"config"`
	DeliveryServices map[string]CrDeliveryService `json:"deliveryServices"`
}

func boolString(in interface{}) string {
	if reflect.TypeOf(in) == nil {
		return "false"
	}

	return "false"
}

func GetCRConfig(cdnName string) (interface{}, error) {

	// contentRouters section
	crQuery := "select * from content_routers where cdnname=\"" + cdnName + "\""
	fmt.Println(crQuery)
	crs := []ContentRouter{}
	err := globalDB.Select(&crs, crQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// monitors section
	mQuery := "select * from monitors where cdnname=\"" + cdnName + "\""
	fmt.Println(mQuery)
	ms := []Monitor{}
	err = globalDB.Select(&ms, mQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// edgeLocations section
	eQuery := "select name,longitude,latitude from cachegroup where type in (select id from type where name=\"EDGE_LOC\")"
	fmt.Println(eQuery)
	edges := []EdgeLocation{}
	err = globalDB.Select(&edges, eQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// stats section
	// TODO JvD

	// config section
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
		ParamMap: miscMap,
		Soa: Soa{
			Admin:   pmap["tld.soa.admin"],
			Expire:  pmap["tld.soa.expire"],
			Minimum: pmap["tld.soa.minimum"],
			Refresh: pmap["tld.soa.refresh"],
			Retry:   pmap["tld.soa.retry"],
		},
		Ttls: Ttls{
			A:      pmap["tld.ttls.A"],
			AAAA:   pmap["tld.ttls.AAAA"],
			DNSKEY: pmap["tld.ttls.DNSKEY"],
			DS:     pmap["tld.ttls.DS"],
			NS:     pmap["tld.ttls.NS"],
			SOA:    pmap["tld.ttls.SOA"],
		},
	}

	// deliveryServices Section
	dQuery := "select * from crconfig_ds_data where cdn_name=\"" + cdnName + "\""
	fmt.Println(">>>> ", dQuery)
	ds := []CrconfigDsData{}
	err = globalDB.Select(&ds, dQuery)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dsMap := make(map[string]CrDeliveryService)
	for _, deliveryService := range ds {
		GeoMap := make(map[string]string) // TODO: should this 1->USA, and 2->CA be hardcoded here?
		CzfOnly := "false"
		if deliveryService.GeoLimit == 1 {
			CzfOnly = "true"
		} else if deliveryService.GeoLimit == 2 {
			GeoMap["countryCode"] = "USA"
		} else if deliveryService.GeoLimit == 3 {
			GeoMap["countryCode"] = "CA"
		}
		// respHdrs := make(map[string]string)
		// domains := make([]string, 0, 0)
		// msets := make([]MatchSetEntry, 0, 0)
		dsMap[deliveryService.XmlId] = CrDeliveryService{
			CoverageZoneOnly: CzfOnly,
			// Domains:              domains,
			IP6RoutingEnabled: boolString(deliveryService.Ipv6RoutingEnabled),
			// MatchSets:            msets,
			MaxDNSIpsForLocation: deliveryService.MaxDnsAnswers,
			MissLocation: MissLocation{
				Latitude:  deliveryService.MissLat,
				Longitude: deliveryService.MissLong,
			},
			Soa: Soa{
				Admin:   pmap["tld.soa.admin"],
				Expire:  pmap["tld.soa.expire"],
				Minimum: pmap["tld.soa.minimum"],
				Refresh: pmap["tld.soa.refresh"],
				Retry:   pmap["tld.soa.retry"],
			},
			TTL: deliveryService.CcrDnsTtl,
			Ttls: Ttls{
				A:      pmap["tld.ttls.A"],
				AAAA:   pmap["tld.ttls.AAAA"],
				DNSKEY: pmap["tld.ttls.DNSKEY"],
				DS:     pmap["tld.ttls.DS"],
				NS:     pmap["tld.ttls.NS"],
				SOA:    pmap["tld.ttls.SOA"],
			},
			// ResponseHeaders: respHdrs,
			Dispersion: Dispersion{
				Shuffled: deliveryService.InitialDispersion,
				Limit:    1,
			},
			GeoEnabled: GeoMap,
		}
		//https://code.google.com/p/go/issues/detail?id=3117
		dService := dsMap[deliveryService.XmlId]
		domains := dService.Domains
		if domains == nil {
			domains = make([]string, 0, 0)
		}
		dService.Domains = append(domains, "hey") // put real domain here
		// put matchlists in
		// put headers in
		dsMap[deliveryService.XmlId] = dService
	}

	// contentServers Section

	return CRConfig{
		ContentRouters:   crs,
		Monitors:         ms,
		EdgeLocations:    edges,
		Config:           cfg,
		DeliveryServices: dsMap,
	}, nil
}
