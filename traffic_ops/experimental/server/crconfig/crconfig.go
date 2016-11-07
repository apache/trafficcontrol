
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crconfig

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Note: a lot of these structs are generated from the DB. No need to type them all out, there's tools for that.
// a view will generate structs also with get_structs.go

// use this view (pq)
// create view content_routers as select ip_address as ip, ip6_address as ip6, profile.name as profile, cachegroup.name as location,
// status.name as status, server.tcp_port as port, host_name, concat(server.host_name, '.', server.domain_name) as fqdn,
// parameter.value as apiport, cdn.name as cdnname
// from server
// join profile on profile.id = server.profile
// join profile_parameter on profile_parameter.profile = profile.id
// join parameter on parameter.id = profile_parameter.parameter
// join cachegroup on cachegroup.id = server.cachegroup
// join status on status.id = server.status
// join cdn on cdn.id = server.cdn_id
// join type on type.id = server.type
// where type.name = 'CCR' and parameter.name='api.port';
type ContentRouter struct {
	Profile  string `db:"profile" json:"profile"`
	Apiport  int64  `db:"apiport" json:"api.port"`
	Location string `db:"location" json:"location"`
	Ip       string `db:"ip" json:"ip"`
	Status   string `db:"status" json:"status"`
	Ip6      string `db:"ip6" json:"ip6"`
	Port     int64  `db:"port" json:"port"`
	Fqdn     string `db:"fqdn" json:"fqdn"`
	HostName string `db:"host_name" json:"hostName,omitempty"`
	Cdnname  string `db:"cdn" json:"cdn"`
}

// use this view (pq)
// create view monitors as select ip_address as ip, ip6_address as ip6, profile.name as profile, cachegroup.name as location,
// status.name as status, server.tcp_port as port, concat(server.host_name, '.', server.domain_name) as fqdn,
// cdn.name as cdnname, host_name
// from server
// join profile on profile.id = server.profile
// join cachegroup on cachegroup.id = server.cachegroup
// join status on status.id = server.status
// join cdn on cdn.id = server.cdn_id
// join type on type.id = server.type
// where type.name = 'RASCAL';

type Monitor struct {
	Profile  string      `db:"profile" json:"profile"`
	Location string      `db:"location" json:"location"`
	Ip       string      `db:"ip" json:"ip"`
	Status   string      `db:"status" json:"status"`
	Ip6      null.String `db:"ip6" json:"ip6"`
	Port     int64       `db:"port" json:"port"`
	Fqdn     string      `db:"fqdn" json:"fqdn"`
	HostName string      `db:"host_name" json:"hostName,omitempty"`
	Cdnname  string      `db:"cdn" json:"cdn"`
}

type EdgeLocation struct {
	Name      string     `db:"name" json:"name,omitempty"`
	Longitude null.Float `db:"longitude" json:"longitude"`
	Latitude  null.Float `db:"latitude" json:"latitude"`
}

// use this view (pq)
// create view crconfig_params as select distinct cdn.name as cdn_name, cdn.id as cdn_id,
// server.profile as profile_id,
// server.type as stype, parameter.name as pname,
// parameter.config_file as cfile, parameter.value as pvalue
// from server
// join cdn on cdn.id = server.cdn_id
// join profile on profile.id = server.profile
// join profile_parameter on profile_parameter.profile = server.profile
// join parameter on parameter.id = profile_parameter.parameter
// where server.type in (select id from type where name in ('EDGE', 'MID', 'CCR'))
// and parameter.config_file = 'CRConfig.json';
type CRConfigParam struct {
	CdnName        string `db:"cdn"`
	Profile        string `db:"profile"`
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
	MatchList []MactchListEntry `json:"matchlist"`
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

type CrStaticDnsEntry struct {
	Name  string `json:"name"`
	Ttl   string `json:"ttl"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type CrDeliveryService struct {
	CoverageZoneOnly     string             `json:"coverageZoneOnly"`
	Domains              []string           `json:"domains"`
	IP6RoutingEnabled    string             `json:"ip6RoutingEnabled"`
	MatchSets            []MatchSetEntry    `json:"matchsets"`
	MaxDNSIpsForLocation null.Int           `json:"maxDnsIpsForLocation"`
	MissLocation         MissLocation       `json:"missLocation"`
	Soa                  Soa                `json:"soa"`
	TTL                  null.Int           `json:"ttl"`
	StaticDnsEntries     []CrStaticDnsEntry `json:"staticDnsEntries,omitempty"`
	Ttls                 Ttls               `json:"ttls"`
	ResponseHeaders      map[string]string  `json:"responseHeaders,omitempty"`
	RequestHeaders       []string           `json:"requestHeaders,omitempty"`
	Dispersion           Dispersion         `json:"dispersion"`
	GeoEnabled           map[string]string  `json:"geoEnabled"`
}

// use this view
// create view crconfig_ds_data as select xml_id, profile, ccr_dns_ttl, global_max_mbps, global_max_tps,
// max_dns_answers, miss_lat, miss_long, protocoltype.name as protocol, ipv6_routing_enabled,
// tr_request_headers, tr_response_headers, initial_dispersion, dns_bypass_cname,
// dns_bypass_ip, dns_bypass_ip6, dns_bypass_ttl, geo_limit, cdn.name as cdn_name,
// regex.pattern as match_pattern, regextype.name as match_type, deliveryservice_regex.set_number,
// staticdnsentry.host as sdns_host, staticdnsentry.address as sdns_address,
// staticdnsentry.ttl as sdns_ttl, sdnstype.name as sdns_type
// from deliveryservice
// join cdn on cdn.id = deliveryservice.cdn_id
// left outer join staticdnsentry on deliveryservice.id = staticdnsentry.deliveryservice
// join deliveryservice_regex on deliveryservice_regex.deliveryservice = deliveryservice.id
// join regex on regex.id = deliveryservice_regex.regex
// join type as protocoltype on protocoltype.id = deliveryservice.type
// join type as regextype on regextype.id = regex.type
// left outer join type as sdnstype on sdnstype.id = staticdnsentry.type;

type CrconfigDsData struct {
	Name               string      `db:"name" json:"name"`
	Profile            int64       `db:"profile" json:"profile"`
	DnsTtl             null.Int    `db:"dns_ttl" json:"dnsTtl"`
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
	CdnName            null.String `db:"cdn" json:"cdn"`
	MatchPattern       string      `db:"match_pattern" json:"matchPattern"`
	MatchType          string      `db:"match_type" json:"matchType"`
	SetNumber          int64       `db:"set_number" json:"setNumber"`
	SdnsHost           null.String `db:"sdns_host" json:"SdnsHost"`
	SdnsAddress        null.String `db:"sdns_address" json:"SdnsAddress"`
	SdnsTtl            null.String `db:"sdns_ttl" json:"SdnsTtl"`
	SdnsType           null.String `db:"sdns_type" json:"SdnsType"`
}

// use this view
// create view content_servers as select distinct  host_name as host_name, profile.name as profile,
// type.name as type, cachegroup.name as location_id, ip_address as ip, cdn.name as cdnname,
// status.name as status, cachegroup.name as cache_group, ip6_address as ip6, tcp_port as port,
// concat(host_name, '.', domain_name) as fqdn, interface_name, parameter.value as hash_count
// from server
// join profile on profile.id = server.profile
// join profile_parameter on profile_parameter.profile = profile.id
// join parameter on parameter.id = profile_parameter.parameter
// join cachegroup on cachegroup.id = server.cachegroup
// join type on type.id = server.type
// join status on status.id = server.status
// join cdn on cdn.id = server.cdn_id
// and parameter.name = 'weight'
// and server.status in (select id from status where name='REPORTED' or name='ONLINE')
// and server.type=(select id from type where name='EDGE');
// this struct is for the db query
type CrContentServer struct {
	HostName      string      `db:"host_name" json:"hostName"`
	Profile       string      `db:"profile" json:"profile"`
	Type          string      `db:"type" json:"type"`
	Location      string      `db:"location" json:"location"`
	Ip            string      `db:"ip" json:"ip"`
	Status        string      `db:"status" json:"status"`
	CacheGroup    string      `db:"cache_group" json:"cacheGroup"`
	Ip6           null.String `db:"ip6" json:"ip6"`
	Port          null.Int    `db:"port" json:"port"`
	Fqdn          string      `db:"fqdn" json:"fqdn"`
	InterfaceName string      `db:"interface_name" json:"interfaceName"`
	HashCount     string      `db:"hash_count" json:"hashCount"`
	CdnName       null.String `db:"cdn" json:"cdn"`
}

type ContentServerDomainList []string
type ContentServerDsMap map[string]ContentServerDomainList

type ContentServer struct {
	Fqdn             string             `json:"fqdn"`
	HashCount        int                `json:"hashCount"`
	HashID           string             `json:"hashId"`
	InterfaceName    string             `json:"interfaceName"`
	IP               string             `json:"ip"`
	IP6              null.String        `json:"ip6"`
	LocationID       string             `json:"locationId"`
	Port             null.Int           `json:"port"`
	Profile          string             `json:"profile"`
	Status           string             `json:"status"`
	Type             string             `json:"type"`
	DeliveryServices ContentServerDsMap `json:"deliveryServices"`
}

// use this view
// create or replace view cr_deliveryservice_server as select distinct regex.pattern as
// pattern, xml_id, deliveryservice.id as ds_id, server.id as srv_id,
// cdn.name as cdnname, server.host_name as server_name
// from deliveryservice
// join deliveryservice_regex on deliveryservice_regex.deliveryservice = deliveryservice.id
// join regex on regex.id = deliveryservice_regex.regex
// join deliveryservice_server on deliveryservice.id = deliveryservice_server.deliveryservice
// join server on server.id = deliveryservice_server.server
// join cdn on cdn.id = server.cdn_id
// where deliveryservice.type != (select id from type where name='ANY_MAP');
type CrDeliveryserviceServer struct {
	Pattern    string      `db:"pattern" json:"pattern"`
	Name       string      `db:"name" json:"name"`
	DsId       int64       `db:"ds_id" json:"id"`
	SrvId      int64       `db:"srv_id" json:"srvId"`
	ServerName string      `db:"server_name" json:"servername"`
	Cdnname    null.String `db:"cdn" json:"cdn"`
	DsType     string      `db:"ds_type" json:"dsType"`
}

type Stats struct {
	CDNName   string `json:"cdn"`
	Date      int    `json:"date"`
	TmHost    string `json:"tm_host"`
	TmPath    string `json:"tm_path"`
	TmUser    string `json:"tm_user"`
	TmVersion string `json:"tm_version"`
}

type CRConfig struct {
	ContentRouters   map[string]ContentRouter     `json:"contentRouters"`
	Monitors         map[string]Monitor           `json:"monitors"`
	EdgeLocations    map[string]EdgeLocation      `json:"edgeLocations"`
	Stats            Stats                        `json:"stats"`
	Config           Config                       `json:"config"`
	DeliveryServices map[string]CrDeliveryService `json:"deliveryServices"`
	ContentServers   map[string]ContentServer     `json:"contentServers"`
}

func boolString(in interface{}) string {
	if reflect.TypeOf(in) == nil {
		return "false"
	}
	if in == true {
		return "true"
	}

	return "false"
}

func contentRoutersSection(cdnName string, db *sqlx.DB) (map[string]ContentRouter, error) {
	crQuery := "select * from content_routers where cdn='" + cdnName + "'"
	crs := []ContentRouter{}
	err := db.Select(&crs, crQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	retMap := make(map[string]ContentRouter)
	for _, row := range crs {
		out := row
		out.HostName = "" // omitempty will make it dissapear
		retMap[row.HostName] = out
	}

	return retMap, nil
}

func monitorSection(cdnName string, db *sqlx.DB) (map[string]Monitor, error) {
	mQuery := "select * from monitors where cdn='" + cdnName + "'"
	ms := []Monitor{}
	err := db.Select(&ms, mQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	retMap := make(map[string]Monitor)
	for _, row := range ms {
		out := row
		out.HostName = "" // omitempty will make it dissapear
		retMap[row.HostName] = out
	}

	return retMap, nil
}

func edgeLocationSection(cdnName string, db *sqlx.DB) (map[string]EdgeLocation, error) {
	eQuery := "select name,longitude,latitude from cachegroup where type in (select id from type where name='EDGE_LOC')"
	edges := []EdgeLocation{}
	err := db.Select(&edges, eQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	retMap := make(map[string]EdgeLocation)
	for _, row := range edges {
		out := row
		out.Name = "" // omitempty will make it dissapear
		retMap[row.Name] = out
	}

	return retMap, nil
}

func configSection(cdnName string, db *sqlx.DB) (Config, map[string]string, error) {
	pQuery := "select * from crconfig_params where cdn='" + cdnName + "'"
	params := []CRConfigParam{}
	err := db.Select(&params, pQuery)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("configSection error selecting params: %v", err)
		return Config{}, nil, err
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
	return cfg, pmap, nil
}

func genReqHeaderList(inString string) []string {
	if inString == "" {
		return nil
	}
	retArray := make([]string, 0, 0)
	for _, header := range strings.Split(inString, "__RETURN__") {
		retArray = append(retArray, header)
	}
	return retArray
}

func genRespHeaderList(inString string) map[string]string {
	if inString == "" {
		return nil
	}
	retMap := make(map[string]string)
	for _, line := range strings.Split(inString, "__RETURN__") {
		fields := strings.Split(line, ":")
		retMap[fields[0]] = fields[1]
	}
	return retMap
}

func deliveryServicesSection(cdnName string, pmap map[string]string, db *sqlx.DB) (map[string]CrDeliveryService, error) {
	dQuery := "select * from crconfig_ds_data where cdn='" + cdnName + "'"
	ds := []CrconfigDsData{}
	err := db.Select(&ds, dQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	reqHeaderList := genReqHeaderList(pmap["LogRequestHeaders"])
	ccrDomain := pmap["domain_name"]
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
		if _, ok := dsMap[deliveryService.Name]; !ok { // there are multiple rows for each DS, only create the struct once
			dsMap[deliveryService.Name] = CrDeliveryService{
				CoverageZoneOnly:     CzfOnly,
				IP6RoutingEnabled:    boolString(deliveryService.Ipv6RoutingEnabled),
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
				TTL: deliveryService.DnsTtl,
				Ttls: Ttls{
					A:      pmap["tld.ttls.A"],
					AAAA:   pmap["tld.ttls.AAAA"],
					DNSKEY: pmap["tld.ttls.DNSKEY"],
					DS:     pmap["tld.ttls.DS"],
					NS:     pmap["tld.ttls.NS"],
					SOA:    pmap["tld.ttls.SOA"],
				},
				ResponseHeaders: genRespHeaderList(deliveryService.TrResponseHeaders.String), // TODO JvD test this
				RequestHeaders:  reqHeaderList,
				Dispersion: Dispersion{
					Shuffled: deliveryService.InitialDispersion,
					Limit:    1,
				},
				GeoEnabled: GeoMap,
			}
		}
		dService := dsMap[deliveryService.Name]
		if deliveryService.MatchType == "HOST_REGEXP" && deliveryService.SetNumber == 0 { // TODO JvD: why / how is this an array?
			// if dService.Domains == nil {
			// 	dService.Domains = make([]string, 0, 0)
			// }
			dsDomain := deliveryService.MatchPattern + "." + ccrDomain
			dsDomain = strings.Replace(dsDomain, ".*\\.", "", 1) // XXX check to see if this should be smarter??
			dsDomain = strings.Replace(dsDomain, "\\..*", "", 1) // XXX check to see if this should be smarter??
			dService.Domains = append(dService.Domains, dsDomain)
		}
		// TODO JvD: add support of set entry 1, 2, 3
		// if dService.MatchSets == nil {
		// 	dService.MatchSets = make([]MatchSetEntry, 0, 10)
		// }
		mType := deliveryService.MatchType
		mType = strings.Replace(mType, "_REGEXP", "", 1)
		mle := MactchListEntry{
			Regex:     deliveryService.MatchPattern,
			MatchType: mType,
		}
		ml := make([]MactchListEntry, 0, 10)
		ml = append(ml, mle)
		mse := MatchSetEntry{
			Protocol:  deliveryService.Protocol,
			MatchList: ml,
		}
		dService.MatchSets = append(dService.MatchSets, mse)

		if deliveryService.TrRequestHeaders.String != "" { // TODO JvD: test this.
			dService.RequestHeaders = append(dService.RequestHeaders, genReqHeaderList(deliveryService.TrRequestHeaders.String)...)
		}

		// if dService.StaticDnsEntries == nil {
		// 	dService.StaticDnsEntries = make([]CrStaticDnsEntry, 0, 10)
		// }
		if deliveryService.SdnsHost.String != "" {
			SdnsEntry := CrStaticDnsEntry{
				Value: deliveryService.SdnsAddress.String,
				Name:  deliveryService.SdnsHost.String,
				Ttl:   deliveryService.SdnsTtl.String,
				Type:  deliveryService.SdnsType.String,
			}
			dService.StaticDnsEntries = append(dService.StaticDnsEntries, SdnsEntry)
		}
		dsMap[deliveryService.Name] = dService
	}
	return dsMap, nil
}

func contentServersSection(cdnName string, ccrDomain string, db *sqlx.DB) (map[string]ContentServer, error) {
	csQuery := "select * from content_servers where cdn='" + cdnName + "'"
	cServers := []CrContentServer{}
	err := db.Select(&cServers, csQuery)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("contentServersSection error selecting content_servers: %v", err)
		return nil, err
	}
	dsServerQuery := "select * from cr_deliveryservice_server"
	dsServers := []CrDeliveryserviceServer{}
	err = db.Select(&dsServers, dsServerQuery)
	if err != nil {
		log.Println("ERROR: >> ", err)
		err = fmt.Errorf("contentServersSection error selecting cr_deliveryservice_server: %v", err)
		return nil, err
	}
	dsMap := make(map[string]ContentServerDsMap)
	for _, row := range dsServers {
		if dsMap[row.ServerName] == nil {
			dsMap[row.ServerName] = make(ContentServerDsMap)
		}
		// if dsMap[row.ServerName][row.Name] == nil {
		// 	dsMap[row.ServerName][row.Name] = make(ContentServerDomainList, 0, 10)
		// }
		pattern := row.Pattern
		if strings.HasSuffix(pattern, "\\..*") {
			pattern = strings.Replace(pattern, ".*\\.", "", 1)
			pattern = strings.Replace(pattern, "\\..*", "", 1)
			if strings.HasPrefix(row.DsType, "HTTP") {
				pattern = row.ServerName + "." + pattern + "." + ccrDomain
			} else {
				pattern = "edge." + pattern + "." + ccrDomain
			}
		}
		dsMap[row.ServerName][row.Name] = append(dsMap[row.ServerName][row.Name], pattern)
	}

	retMap := make(map[string]ContentServer)
	for _, row := range cServers {
		hCount, _ := strconv.Atoi(row.HashCount)
		hCount = hCount * 1000 // TODO JvD
		retMap[row.HostName] = ContentServer{
			Fqdn:             row.Fqdn,
			HashCount:        hCount,
			HashID:           row.HostName,
			InterfaceName:    row.InterfaceName,
			IP:               row.Ip,
			IP6:              row.Ip6,
			LocationID:       row.CacheGroup,
			Port:             row.Port,
			Profile:          row.Profile,
			Status:           row.Status,
			Type:             row.Status,
			DeliveryServices: dsMap[row.HostName],
		}
	}

	return retMap, nil
}

func statsSection(cdnName string) (Stats, error) {
	stats := Stats{
		CDNName:   cdnName,
		Date:      int(time.Now().Unix()),
		TmHost:    "dummy",
		TmPath:    "dummy",
		TmUser:    "jvd",
		TmVersion: "2.0!",
	}
	return stats, nil
}

func GetCRConfig(cdnName string, db *sqlx.DB) (interface{}, error) {
	crs, err := contentRoutersSection(cdnName, db)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("GetCRConfig error getting contentRoutersSection: %v", err)
		return nil, err
	}
	ms, err := monitorSection(cdnName, db)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("GetCRConfig error getting monitorSection: %v", err)
		return nil, err
	}
	edges, err := edgeLocationSection(cdnName, db)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("GetCRConfig error getting edgeLocationSection: %v", err)
		return nil, err
	}
	cfg, pmap, err := configSection(cdnName, db)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("GetCRConfig error getting configSection: %v", err)
		return nil, err
	}
	dsMap, err := deliveryServicesSection(cdnName, pmap, db)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("GetCRConfig error getting deliveryServicesSection: %v", err)
		return nil, err
	}
	cServermap, err := contentServersSection(cdnName, pmap["domain_name"], db)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("GetCRConfig error getting contentServersSection: %v", err)
		return nil, err
	}
	stats, err := statsSection(cdnName)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("GetCRConfig error getting statsSection: %v", err)
		return nil, err
	}

	return CRConfig{
		ContentRouters:   crs,
		Monitors:         ms,
		EdgeLocations:    edges,
		Config:           cfg,
		Stats:            stats,
		DeliveryServices: dsMap,
		ContentServers:   cServermap,
	}, nil
}
