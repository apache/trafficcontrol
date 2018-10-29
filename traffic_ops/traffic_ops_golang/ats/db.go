package ats

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
	"database/sql"
	"errors"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

type ProfileData struct {
	ID   int
	Name string
}

// getProfileData returns the necessary info about the profile, whether it exists, and any error.
func getProfileData(tx *sql.Tx, id int) (ProfileData, bool, error) {
	// TODO implement, determine what fields are necessary
	qry := `
SELECT
  p.name
FROM
  profile p
WHERE
  p.id = $1
`
	v := ProfileData{ID: id}
	if err := tx.QueryRow(qry, id).Scan(&v.Name); err != nil {
		if err == sql.ErrNoRows {
			return ProfileData{}, false, nil
		}
		return ProfileData{}, false, errors.New("querying: " + err.Error())
	}
	return v, true, nil
}

func GetNameVersionString(tx *sql.Tx) (string, error) {
	qry := `
SELECT
  p.name,
  p.value
FROM
  parameter p
WHERE
  (p.name = 'tm.toolname' OR p.name = 'tm.url') AND p.config_file = 'global'
`
	rows, err := tx.Query(qry)
	if err != nil {
		return "", errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	toolName := ""
	url := ""
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return "", errors.New("scanning: " + err.Error())
		}
		if name == "tm.toolname" {
			toolName = val
		} else if name == "tm.url" {
			url = val
		}
	}
	return toolName + " (" + url + ")", nil
}

func GetProfileParamData(tx *sql.Tx, profileID int, configFile string) (map[string]string, error) {
	// TODO add another func to return a slice, for things that don't need a map, for performance? Does it make a difference?
	qry := `
SELECT
  p.name,
  p.value
FROM
  parameter p
  join profile_parameter pp on p.id = pp.parameter
  JOIN profile pr on pr.id = pp.profile
WHERE
  pr.id = $1
  AND p.config_file = $2
`
	rows, err := tx.Query(qry, profileID, configFile)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[string]string{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if name == "location" {
			continue
		}
		params[name] = val
	}
	return params, nil
}

type ProfileDS struct {
	Type       tc.DSType
	OriginFQDN *string
}

func GetProfileDS(tx *sql.Tx, profileID int) ([]ProfileDS, error) {
	qry := `
SELECT
  dstype.name AS ds_type,
  (SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
    FROM origin o
    WHERE o.deliveryservice = ds.id
    AND o.is_primary) as org_server_fqdn
FROM
  deliveryservice ds
  JOIN type as dstype ON ds.type = dstype.id
WHERE
  ds.id IN (
    SELECT DISTINCT deliveryservice
    FROM deliveryservice_server
    WHERE server IN (SELECT id FROM server WHERE profile = $1)
  )
`
	rows, err := tx.Query(qry, profileID)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := []ProfileDS{}
	for rows.Next() {
		d := ProfileDS{}
		if err := rows.Scan(&d.Type, &d.OriginFQDN); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		d.Type = tc.DSTypeFromString(string(d.Type))
		dses = append(dses, d)
	}
	return dses, nil
}

// GetProfileParamValue gets the value of a parameter assigned to a profile, by name and config file.
// Returns the parameter, whether it existed, and any error.
func GetProfileParamValue(tx *sql.Tx, profileID int, configFile string, name string) (string, bool, error) {
	qry := `
SELECT
  p.value
FROM
  parameter p
  JOIN profile_parameter pp ON p.id = pp.parameter
WHERE
  pp.profile = $1
  AND p.config_file = $2
  AND p.name = $3
`
	val := ""
	if err := tx.QueryRow(qry, profileID, configFile, name).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying: " + err.Error())
	}
	return val, true, nil
}

func serverParamData(tx *sql.Tx, profileID int, configFile string, serverHost string, serverDomain string) (map[string]string, error) {
	qry := `
SELECT
  p.id,
  p.name,
  p.value
FROM
  parameter p
  join profile_parameter pp on p.id = pp.parameter
  JOIN profile pr on p.id = pp.profile
WHERE
  pr.id = $1
  AND p.config_file = $2
`
	rows, err := tx.Query(qry, profileID, configFile)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[string]string{}
	for rows.Next() {
		id := 0
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if name == "location" {
			continue
		}

		// some files have multiple lines with the same key... handle that with param id.
		key := name
		if _, ok := params[name]; ok {
			key += "__" + strconv.Itoa(id)
		}
		if val == "STRING __HOSTNAME__" {
			val = serverHost + "." + serverDomain
		}
		params[key] = val
	}
	return params, nil
}

// GetServerURISignedDSes returns a list of delivery service names which have the given server assigned and have URI signing enabled, and any error.
func GetServerURISignedDSes(tx *sql.Tx, serverHostName string, serverPort int) ([]tc.DeliveryServiceName, error) {
	qry := `
SELECT
  ds.xml_id
FROM
  deliveryservice ds
  JOIN deliveryservice_server dss ON ds.id = dss.deliveryservice
  JOIN server s ON s.id = dss.server
WHERE
  s.host_name = $1
  AND s.tcp_port = $2
  AND ds.signing_algorithm = 'uri_signing'
`
	rows, err := tx.Query(qry, serverHostName, serverPort)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := []tc.DeliveryServiceName{}
	for rows.Next() {
		ds := tc.DeliveryServiceName("")
		if err := rows.Scan(&ds); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		dses = append(dses, ds)
	}
	return dses, nil
}

type ConfigProfileParams struct {
	FileNameOnDisk string
	Location       string
	URL            string
	APIURI         string
}

// GetLocationParams returns a map[configFile]locationParams, and any error. If either param doesn't exist, an empty string is returned without error.
func GetLocationParams(tx *sql.Tx, profileID int) (map[string]ConfigProfileParams, error) {
	qry := `
SELECT
  p.name,
  p.config_file,
  p.value
FROM
  parameter p
  JOIN profile_parameter pp ON pp.parameter = p.id
WHERE
  pp.profile = $1
`
	rows, err := tx.Query(qry, profileID)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[string]ConfigProfileParams{}
	for rows.Next() {
		name := ""
		file := ""
		val := ""
		if err := rows.Scan(&name, &file, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if name == "location" {
			p := params[file]
			p.FileNameOnDisk = file
			p.Location = val
			params[file] = p
		} else if name == "URL" {
			p := params[file]
			p.URL = val
			params[file] = p
		}
	}
	return params, nil
}

type TMParams struct {
	URL             string
	ReverseProxyURL string
}

// GetTMParams returns the global "tm.url" and "tm.rev_proxy.url" parameters, and any error. If either param doesn't exist, an empty string is returned without error.
func GetTMParams(tx *sql.Tx) (TMParams, error) {
	rows, err := tx.Query(`SELECT name, value from parameter where config_file = 'global' AND (name = 'tm.url' OR name = 'tm.rev_proxy.url')`)
	if err != nil {
		return TMParams{}, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	p := TMParams{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return TMParams{}, errors.New("scanning: " + err.Error())
		}
		if name == "tm.url" {
			p.URL = val
		} else if name == "tm.rev_proxy.url" {
			p.ReverseProxyURL = val
		} else {
			return TMParams{}, errors.New("querying got unexpected parameter: " + name + " (value: '" + val + "')") // should never happen
		}
	}
	return p, nil
}

type ServerInfo struct {
	CacheGroupID                  int
	CDN                           tc.CDNName
	CDNID                         int
	DomainName                    string
	HostName                      string
	ID                            int
	IP                            string
	ParentCacheGroupID            int
	ParentCacheGroupType          string
	ProfileID                     int
	ProfileName                   string
	Port                          int
	SecondaryParentCacheGroupID   int
	SecondaryParentCacheGroupType string
	TypeName                      string
}

// getServerInfo returns the necessary info about the server, whether the server exists, and any error.
func getServerInfo(tx *sql.Tx, id int) (ServerInfo, bool, error) {
	// TODO separate this into only what's requried for each config file, and create interfaces for funcs?
	qry := `
SELECT
  c.name as cdn,
  s.cdn_id,
  s.host_name,
  c.domain_name,
  s.ip_address,
  s.profile AS profile_id,
  p.name AS profile_name,
  s.tcp_port,
  t.name as type,
  s.cachegroup,
  COALESCE(cg.parent_cachegroup_id, -1),
  COALESCE(cg.secondary_parent_cachegroup_id, -1),
  COALESCE(parentt.name, '') as parent_cachegroup_type,
  COALESCE(sparentt.name, '') as secondary_parent_cachegroup_type
FROM
  server s
  JOIN cdn c ON s.cdn_id = c.id
  JOIN type t ON s.type = t.id
  JOIN profile p ON p.id = s.profile
  JOIN cachegroup cg on s.cachegroup = cg.id
  LEFT JOIN type parentt on parentt.id = (select type from cachegroup where id = cg.parent_cachegroup_id)
  LEFT JOIN type sparentt on sparentt.id = (select type from cachegroup where id = cg.secondary_parent_cachegroup_id)
WHERE
  s.id = $1
`
	s := ServerInfo{ID: id}
	if err := tx.QueryRow(qry, id).Scan(&s.CDN, &s.CDNID, &s.HostName, &s.DomainName, &s.IP, &s.ProfileID, &s.ProfileName, &s.Port, &s.TypeName, &s.CacheGroupID, &s.ParentCacheGroupID, &s.SecondaryParentCacheGroupID, &s.ParentCacheGroupType, &s.SecondaryParentCacheGroupType); err != nil {
		if err == sql.ErrNoRows {
			return ServerInfo{}, false, nil
		}
		return ServerInfo{}, false, errors.New("querying server info: " + err.Error())
	}
	return s, true, nil
}

// GetFirstScopeParameter returns the value of the arbitrarily-first parameter with the name 'scope' and the given config file, whether a parameter was found, and any error.
func GetFirstScopeParameter(tx *sql.Tx, cfgFile string) (string, bool, error) {
	v := ""
	if err := tx.QueryRow(`SELECT p.value FROM parameter p WHERE p.config_file = $1 AND p.name = 'scope'`, cfgFile).Scan(&v); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying first scope parameter: " + err.Error())
	}
	return v, true, nil
}

type DSData struct {
	Type       tc.DSType
	OriginFQDN *string
}

func GetDSData(tx *sql.Tx, serverID int) ([]DSData, error) {
	qry := `
SELECT
  dstype.name AS ds_type,
  (SELECT o.protocol::text || \'://\' || o.fqdn || rtrim(concat(\':\', o.port::text), \':\')
    FROM origin o
    WHERE o.deliveryservice = ds.id
    AND o.is_primary) as org_server_fqdn
FROM
  deliveryservice ds
  JOIN type as dstype ON ds.type = dstype.id
WHERE
  ds.id IN (
    SELECT DISTINCT deliveryservice
    FROM deliveryservice_server
    WHERE server = $1
  )
`
	rows, err := tx.Query(qry, serverID)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := []DSData{}
	for rows.Next() {
		d := DSData{}
		if err := rows.Scan(&d.Type, &d.OriginFQDN); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		d.Type = tc.DSTypeFromString(string(d.Type))
		dses = append(dses, d)
	}
	return dses, nil
}

// type IPAllowAccess struct {
// 	SourceIP string
// 	Action   string
// 	Method   string
// }

// func getIPAllowData(tx *sql.Tx, serverInfo *ServerInfo, fileName string) ([]IPAllowAccess, error) {
// 	allow := []IPAllowAccess{}

// 	// localhost is trusted
// 	allow = append(allow, IPAllowAccess{
// 		SourceIP: "127.0.0.1",
// 		Action:   "ip_allow",
// 		Method:   "ALL",
// 	})
// 	allow = append(allow, IPAllowAccess{
// 		SourceIP: "::1",
// 		Action:   "ip_allow",
// 		Method:   "ALL",
// 	})

// 	coalesceMasklenV4 := 24
// 	coalesceNumberV4 := 5
// 	coalesceMasklenV6 := 48
// 	coalesceNumberV6 := 5

// 	params, err := GetProfileParamData(tx, serverInfo.ProfileID, "ip_allow.config")
// 	if err != nil {
// 		return nil, errors.New("getting profile param data: " + err.Error())
// 	}

// 	for name, val := range params {
// 		switch name {
// 		case "purge_allow_ip":
// 			allow = append(allow, IPAllowAccess{
// 				SourceIP: val,
// 				Action:   "ip_allow",
// 				Method:   "ALL",
// 			})
// 		case "coalesce_masklen_v4":
// 			coalesceMasklenV4 := val
// 		case "coalesce_number_v4":
// 			coalseceNumberV4 := val
// 		case "coalesce_masklen_v6":
// 			coalesceMasklenV6 := val
// 		case "coalesce_number_v6":
// 			coalseceNumberV6 := val
// 		}
// 	}

// 	if tc.CacheTypeFromString(serverInfo.TypeName) == tc.CacheTypeMid {
// 		allowedServers, err := getAllowedServers(tx, serverInfo.CDNID, serverInfo.CacheGroupID)
// 		if err != nil {
// 			return nil, errors.New("getting allowed servers: " + err.Error())
// 		}
// 		ipv4 = 42 //				my $ipv4 = NetAddr::IP->new( $allow_row->ip_address, $allow_row->ip_netmask );
// 	} else {

// 	}

// }

type ServerIPs struct {
	HostName    string
	IPv4        string
	IPv4Netmask string
	IPv6        string
}

// getAllowedServers returns all servers which are allowed to talk to this server
func getAllowedServers(tx *sql.Tx, serverCDNID int, serverCacheGroupID int) ([]ServerIPs, error) {
	// TODO rename. GetAllowedEdgesForMid?
	// TODO use constants for query RASCAL and EDGE
	qry := `
WITH edge_locs AS (
  SELECT id FROM cachegroup WHERE parent_cachegroup_id = $1 OR secondary_parent_cachegroup_id = $1
)
SELECT
  s.host_name,
  s.ip_address,
  s.ip_netmask,
  s.ip6_address
FROM
  server s
  JOIN type tp ON s.type = tp.id
  JOIN cdn on cdn.id = $2
WHERE
  tp.name = 'RASCAL'
  OR (tp.name like 'EDGE%' AND tp.use_in_table = 'server' AND s.cachegroup IN (SELECT id FROM edge_locs))
`
	rows, err := tx.Query(qry, serverCacheGroupID, serverCDNID)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	ips := []ServerIPs{}
	for rows.Next() {
		i := ServerIPs{}
		if err := rows.Scan(&i.HostName, &i.IPv4, &i.IPv4Netmask, &i.IPv6); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		ips = append(ips, i)
	}
	return ips, nil
}

// type ParentDS struct {
// 	RemapLine            map[string]string // map[from]to
// 	HeaderRewriteFile    *string
// 	MidHeaderRewriteFile *string

// 	CacheURL             *string
// 	DomainName           *string
// 	DSCP                 int
// 	EdgeHeaderRewrite    *string
// 	FQPacingRate         *int
// 	MidHeaderRewrite     *string
// 	MultiSiteOrigin      *bool
// 	OriginServerFQDN     *string
// 	QStringIgnore        *int
// 	RangeRequestHandling *string
// 	RemapText            *string
// 	RegexRemap           *string
// 	SigningAlgorithm     *string
// 	Signed               *bool
// 	Type                 tc.DSType
// 	XMLID                string
// }

// func GetParentDSData(tx *sql.Tx, profileID int) ([]ProfileDSInfo, error) {
// 	qry := `
// SELECT
//   deliveryservice.xml_id,
//   deliveryservice.id AS ds_id,
//   deliveryservice.dscp,
//   deliveryservice.routing_name,
//   deliveryservice.signing_algorithm,
//   deliveryservice.qstring_ignore,
//   (SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
//     FROM origin o
//     WHERE o.deliveryservice = deliveryservice.id
//     AND o.is_primary) as org_server_fqdn,
//   deliveryservice.multi_site_origin,
//   deliveryservice.range_request_handling,
//   deliveryservice.fq_pacing_rate,
//   deliveryservice.origin_shield,
//   regex.pattern,
//   retype.name AS re_type,
//   dstype.name AS ds_type,
//   cdn.domain_name AS domain_name,
//   deliveryservice_regex.set_number,
//   deliveryservice.edge_header_rewrite,
//   deliveryservice.mid_header_rewrite,
//   deliveryservice.regex_remap,
//   deliveryservice.cacheurl,
//   deliveryservice.remap_text,
//   deliveryservice.protocol,
//   deliveryservice.profile,
//   deliveryservice.anonymous_blocking_enabled
// FROM
//   deliveryservice
//   JOIN deliveryservice_regex ON deliveryservice_regex.deliveryservice = deliveryservice.id
//   JOIN regex ON deliveryservice_regex.regex = regex.id
//   JOIN type retype ON regex.type = retype.id
//   JOIN type dstype ON deliveryservice.type = dstype.id
//   JOIN cdn ON cdn.id = deliveryservice.cdn_id
// WHERE
//   cdn.name = $1
//   AND deliveryservice.id in (SELECT deliveryservice_server.deliveryservice FROM deliveryservice_server)
//   AND deliveryservice.active = true
// ORDER BY
//   ds_id,
//   re_type,
//   set_number
// `

// 	dsDataByProfile, err := GetDSDataByProfile(tx, profileID)
// 	if err != nil {
// 		return nil, fmt.Errorf("getting ds data by profile for %+v: %+v", profileID, err)
// 	}

// 	dsInfos := []ProfileDSInfo{}

// 	for _, ds := range dsDataByProfile {
// 		dsInfo := ProfileDSInfo{}

// 		if ds.RegexType == tc.DSMatchTypeHostRegex {
// 			hostRegex := ds.Pattern
// 			mapTo := (*string)(nil)
// 			if ds.OriginServerFQDN != nil {
// 				mapTo = util.StrPtr(strings.TrimSuffix(*ds.OriginServerFQDN, "/") + "/")
// 			}

// 			httpMapFrom := "http://" + hostRegex + "/"
// 			httpsMapFrom := "https://" + hostRegex + "/"
// 			if strings.HasSuffix(hostRegex, `.*`) {
// 				re := strings.Replace(hostRegex, `/`, ``)
// 				re := strings.Replace(hostRegex, `.*`, ``)
// 				hName := `__http__`
// 				if ds.Type.IsDNS() {
// 					hName = ds.RoutingName
// 				}
// 				portStr := ":" + "__SERVER_TCP_PORT__"
// 				httpMapFrom = "http://" + hName + re + ds.DomainName + portStr + "/"
// 				httpsMapFrom = "https://" + hName + re + ds.DomainName + "/" // TODO verify https shouldn't get the port string?
// 			}

// 			if ds.Protocol != nil || *ds.Protocol == tc.DSProtocolHTTP { // "or", not "and": default to HTTP if protocol is nil.
// 				dsInfo.RemapLine[httpMapFrom] = mapTo
// 			} else if *ds.Protocol == tc.DSProtocolHTTPS || *ds.Protocol == tc.DSProtocolHTTPToHTTPS {
// 				dsInfo.RemapLine[httpsMapFrom] = mapTo
// 			} else if *ds.Protocol == tc.DSProtocolHTTPAndHTTPS {
// 				dsInfo.RemapLine[httpsMapFrom] = mapTo
// 				dsInfo.RemapLine[httpsMapFrom] = mapTo
// 			}
// 		}

// 		dsInfo.DSCP = ds.DSCP
// 		dsInfo.OriginServerFQDN = ds.OriginServerFQDN
// 		dsInfo.Type = ds.Type
// 		dsInfo.DomainName = ds.DomainName
// 		dsInfo.Signed = ds.SigningAlgorithm != nil && *ds.SigningAlgorithm == "url_sig" // TODO change to enum, once PR is merged
// 		dsInfo.SigningAlgorithm = ds.SigningAlgorithm
// 		dsInfo.QStringIgnore = ds.QStringIgnore
// 		dsInfo.XMLID = ds.XMLID
// 		dsInfo.EdgeHeaderRewrite = ds.EdgeHeaderRewrite
// 		dsInfo.MidHeaderRewrite = ds.MidHeaderRewrite
// 		dsInfo.RegexRemap = ds.RegexRemap
// 		dsInfo.RangeRequestHandling = ds.RangeRequestHandling
// 		dsInfo.FQPacingRate = ds.FQPacingRate
// 		dsInfo.OriginShield = ds.OriginShield
// 		dsInfo.Cacheurl = ds.Cacheurl
// 		dsInfo.RemapText = ds.RemapText
// 		dsInfo.MultiSiteOrigin = ds.MultiSiteOrigin

// 		if ds.EdgeHeaderRewrite != nil && *ds.EdgeHeaderRewrite != "" {
// 			dsInfo.HeaderRewriteFile = util.StrPtr("hdr_rw_" + ds.XMLID + ".config")
// 		}
// 		if ds.MidHeaderRewrite != nil && *ds.MidHeaderRewrite != "" {
// 			dsInfo.MidHeaderRewriteFile = util.StrPtr("hdr_rw_mid_" + ds.XMLID + ".config")
// 		}
// 		if ds.CacheURL != nil && *ds.CacheURL != "" {
// 			dsInfo.CacheURLFile = util.StrPtr("cacheurl_" + ds.XMLID + ".config")
// 		}
// 		if ds.Profile != nil {
// 			// my $dsparamrs = $self->db->resultset('ProfileParameter')->search( { profile => $row->{'profile'} }, { prefetch => [ 'profile', 'parameter' ] } );
// 			// while ( my $prow = $dsparamrs->next ) {
// 			// 	$dsinfo->{dslist}->[$j]->{'param'}->{ $prow->parameter->config_file }->{ $prow->parameter->name } = $prow->parameter->value;
// 			// }
// 		}

// 	}

// }
