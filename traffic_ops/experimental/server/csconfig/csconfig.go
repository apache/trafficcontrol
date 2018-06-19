
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csconfig

import (
	"fmt"
	"github.com/apache/trafficcontrol/traffic_ops/experimental/server/api"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
	"log"
	// "reflect"
	// "strconv"
	// "strings"
	"time"
)

// Note: a lot of these structs are generated from the DB. No need to type them all out, there's tools for that.
// a view will generate structs also with get_structs.go

// create or replace view csconfig_params as select  name, value, config_file,
// profile_parameter.profile as profile
// from profile_parameter
// join parameter on parameter.id = profile_parameter.parameter
type CsconfigParam struct {
	Name       string `db:"name" json:"name"`
	Value      string `db:"value" json:"value"`
	ConfigFile string `db:"config_file" json:"configFile"`
	Profile    string `db:"profile" json:"profile"`
}

// create or replace view csconfig_remap as
// select deliveryservice.*, regex.pattern as r_pattern, rtype.name as r_type,
// server.id as server_id
// from server
// join deliveryservice_server on deliveryservice_server.server = server.id
// join deliveryservice on deliveryservice.id = deliveryservice_server.deliveryservice
// join deliveryservice_regex on deliveryservice_regex.deliveryservice = deliveryservice.id
// join regex on regex.id = deliveryservice_regex.regex
// join type as rtype on regex.type = rtype.id
type CsconfigRemap struct {
	Id                   int64       `db:"id" json:"id"`
	XmlId                string      `db:"xml_id" json:"xmlId"`
	Active               int64       `db:"active" json:"active"`
	Dscp                 int64       `db:"dscp" json:"dscp"`
	Signed               null.Int    `db:"signed" json:"signed"`
	QstringIgnore        null.Int    `db:"qstring_ignore" json:"qstringIgnore"`
	GeoLimit             null.Int    `db:"geo_limit" json:"geoLimit"`
	HttpBypassFqdn       null.String `db:"http_bypass_fqdn" json:"httpBypassFqdn"`
	DnsBypassIp          null.String `db:"dns_bypass_ip" json:"dnsBypassIp"`
	DnsBypassIp6         null.String `db:"dns_bypass_ip6" json:"dnsBypassIp6"`
	DnsBypassTtl         null.Int    `db:"dns_bypass_ttl" json:"dnsBypassTtl"`
	OrgServerFqdn        null.String `db:"org_server_fqdn" json:"orgServerFqdn"`
	Type                 int64       `db:"type" json:"type"`
	Profile              int64       `db:"profile" json:"profile"`
	CdnId                int64       `db:"cdn_id" json:"cdnId"`
	CcrDnsTtl            null.Int    `db:"ccr_dns_ttl" json:"ccrDnsTtl"`
	GlobalMaxMbps        null.Int    `db:"global_max_mbps" json:"globalMaxMbps"`
	GlobalMaxTps         null.Int    `db:"global_max_tps" json:"globalMaxTps"`
	LongDesc             null.String `db:"long_desc" json:"longDesc"`
	LongDesc1            null.String `db:"long_desc_1" json:"longDesc1"`
	LongDesc2            null.String `db:"long_desc_2" json:"longDesc2"`
	MaxDnsAnswers        null.Int    `db:"max_dns_answers" json:"maxDnsAnswers"`
	InfoUrl              null.String `db:"info_url" json:"infoUrl"`
	MissLat              null.Float  `db:"miss_lat" json:"missLat"`
	MissLong             null.Float  `db:"miss_long" json:"missLong"`
	CheckPath            null.String `db:"check_path" json:"checkPath"`
	LastUpdated          time.Time   `db:"last_updated" json:"lastUpdated"`
	Protocol             null.Int    `db:"protocol" json:"protocol"`
	SslKeyVersion        null.Int    `db:"ssl_key_version" json:"sslKeyVersion"`
	Ipv6RoutingEnabled   null.Int    `db:"ipv6_routing_enabled" json:"ipv6RoutingEnabled"`
	RangeRequestHandling null.Int    `db:"range_request_handling" json:"rangeRequestHandling"`
	EdgeHeaderRewrite    null.String `db:"edge_header_rewrite" json:"edgeHeaderRewrite"`
	OriginShield         null.String `db:"origin_shield" json:"originShield"`
	MidHeaderRewrite     null.String `db:"mid_header_rewrite" json:"midHeaderRewrite"`
	RegexRemap           null.String `db:"regex_remap" json:"regexRemap"`
	Cacheurl             null.String `db:"cacheurl" json:"cacheurl"`
	RemapText            null.String `db:"remap_text" json:"remapText"`
	MultiSiteOrigin      null.Int    `db:"multi_site_origin" json:"multiSiteOrigin"`
	DisplayName          string      `db:"display_name" json:"displayName"`
	TrResponseHeaders    null.String `db:"tr_response_headers" json:"trResponseHeaders"`
	InitialDispersion    null.Int    `db:"initial_dispersion" json:"initialDispersion"`
	DnsBypassCname       null.String `db:"dns_bypass_cname" json:"dnsBypassCname"`
	TrRequestHeaders     null.String `db:"tr_request_headers" json:"trRequestHeaders"`
	RPattern             null.String `db:"r_pattern" json:"rPattern"`
	RType                string      `db:"r_type" json:"rType"`
	ServerId             int64       `db:"server_id" json:"serverId"`
}

type CsConfig struct {
	Remaps []CsconfigRemap `json:"remaps"`
	Params []CsconfigParam `json:"allParams"`
}

func getCSConfigParams(profile string, db *sqlx.DB) ([]CsconfigParam, error) {
	ret := []CsconfigParam{}
	arg := CsconfigParam{Profile: profile}
	nstmt, err := db.PrepareNamed(`select * from csconfig_params where profile=:profile`)
	err = nstmt.Select(&ret, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nstmt.Close()
	return ret, nil
}

// \todo add port (Servers PK is a compound key, host_name and port
func getCSConfigRemap(serverName string, db *sqlx.DB) ([]CsconfigRemap, error) {
	ret := []CsconfigRemap{}
	arg := api.Servers{HostName: serverName}
	nstmt, err := db.PrepareNamed(`select * from csconfig_remap where host_name=:host_name`)
	err = nstmt.Select(&ret, arg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nstmt.Close()
	return ret, nil
}

func GetCSConfig(hostName string, port int64, db *sqlx.DB) (interface{}, error) {

	// stats, err := statsSection(cdnName)

	serverInterface, err := api.GetServer(hostName, port, db)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	server, ok := serverInterface.(api.Servers)
	if !ok {
		err = fmt.Errorf("GetServer returned a non-server")
		log.Println(err)
		return nil, err
	}

	params, err := getCSConfigParams(server.Links.ProfilesLink.ID, db)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	remaps, err := getCSConfigRemap(server.HostName, db) // TODO(take port, part of PK)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return CsConfig{
		Remaps: remaps,
		Params: params,
	}, nil
}
