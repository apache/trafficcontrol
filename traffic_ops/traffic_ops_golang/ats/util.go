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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	// "github.com/dspinhirne/netaddr-go"
)

func getServerScope(tx *sql.Tx, cfgFile string, serverType string) (tc.ATSConfigMetaDataConfigFileScope, error) {
	switch {
	case cfgFile == "cache.config" && tc.CacheTypeFromString(serverType) == tc.CacheTypeMid:
		return tc.ATSConfigMetaDataConfigFileScopeServers, nil
	default:
		return getScope(tx, cfgFile)
	}
}

// getScope returns the ATSConfigMetaDataConfigFileScope for the given config file, and potentially the given server. If the config is not a Server scope, i.e. was part of an endpoint which does not include a server name or id, the server may be nil.
func getScope(tx *sql.Tx, cfgFile string) (tc.ATSConfigMetaDataConfigFileScope, error) {
	switch {
	case cfgFile == "ip_allow.config":
		fallthrough
	case cfgFile == "parent.config":
		fallthrough
	case cfgFile == "hosting.config":
		fallthrough
	case cfgFile == "packages":
		fallthrough
	case cfgFile == "chkconfig":
		fallthrough
	case cfgFile == "remap.config":
		fallthrough
	case strings.HasPrefix(cfgFile, "to_ext_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeServers, nil
	case cfgFile == "12M_facts":
		fallthrough
	case cfgFile == "50-ats.rules":
		fallthrough
	case cfgFile == "astats.config":
		fallthrough
	case cfgFile == "cache.config":
		fallthrough
	case cfgFile == "drop_qstring.config":
		fallthrough
	case cfgFile == "logs_xml.config":
		fallthrough
	case cfgFile == "logging.config":
		fallthrough
	case cfgFile == "plugin.config":
		fallthrough
	case cfgFile == "records.config":
		fallthrough
	case cfgFile == "storage.config":
		fallthrough
	case cfgFile == "volume.config":
		fallthrough
	case cfgFile == "sysctl.conf":
		fallthrough
	case strings.HasPrefix(cfgFile, "url_sig_") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "uri_signing_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeProfiles, nil

	case cfgFile == "bg_fetch.config":
		fallthrough
	case cfgFile == "regex_revalidate.config":
		fallthrough
	case cfgFile == "ssl_multicert.config":
		fallthrough
	case strings.HasPrefix(cfgFile, "cacheurl") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "hdr_rw_") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "regex_remap_") && strings.HasSuffix(cfgFile, ".config"):
		fallthrough
	case strings.HasPrefix(cfgFile, "set_dscp_") && strings.HasSuffix(cfgFile, ".config"):
		return tc.ATSConfigMetaDataConfigFileScopeCDNs, nil
	}

	scope, ok, err := GetFirstScopeParameter(tx, cfgFile)
	if err != nil {
		return tc.ATSConfigMetaDataConfigFileScopeInvalid, errors.New("getting scope parameter: " + err.Error())
	}
	if !ok {
		scope = string(tc.ATSConfigMetaDataConfigFileScopeServers)
	}
	return tc.ATSConfigMetaDataConfigFileScope(scope), nil
}

// GetIPv4CIDRs takes a list of IP addresses, and returns a list of strings, of reasonably compact CIDRs which encompass the given IPs.
// func GetIPv4CIDRs(ips []string) ([]string, error) {
// 	// for i, ip := range ips {
// 	// 	ips[i] = ip + "/32"
// 	// }

// 	ipList, err := netaddr.NewIPv4NetList(ips)
// 	if err != nil {
// 		return nil, errors.New("parsing IPs: " + err.Error())
// 	}
// 	ipList = ipList.Summ()

// 	cidrs := []string{}
// 	for _, ip := range ipList {
// 		cidrs = append(cidrs, ip.String())
// 	}
// 	return cidrs, nil
// }

// func GetIPv6CIDRs(ips []string) ([]string, error) {
// 	ipList, err := netaddr.NewIPv6NetList(ips)
// 	if err != nil {
// 		return nil, errors.New("parsing IPs: " + err.Error())
// 	}
// 	ipList.Summ()

// 	cidrs := []string{}
// 	for _, ip := range ipList {
// 		cidrs = append(cidrs, ip.String())
// 	}
// 	return cidrs, nil
// }

// type ProfileDS struct {
// 	CacheURL                 *string
// 	DeepCachingType          DeepCachingType
// 	DomainName               *string
// 	DSCP                     int
// 	EdgeHeaderRewrite        *string
// 	FQPacingRate             *int
// 	ID                       int
// 	MidHeaderRewrite         *string
// 	MultiSiteOrigin          *bool
// 	MultiSiteOriginAlgorithm *string
// 	OriginServerFQDN         *string
// 	OriginShield             *string
// 	Pattern                  *string
// 	Profile                  *int
// 	Protocol                 *int
// 	QStringIgnore            *int
// 	RangeRequestHandling     *string
// 	RemapText                *string
// 	RegexRemap               *string
// 	RegexType                *string
// 	RoutingName              string
// 	SigningAlgorithm         *string
// 	SSLKeyVersion            *int
// 	Type                     tc.DeliveryServiceType
// 	XMLID                    string
// }

// func GetDSDataByProfile(tx *sql.Tx, profileID int) ([]ProfileDS, error) {
// 	qry := `
// SELECT
//   ds.id,
//   ds.xml_id,
//   ds.dscp,
//   ds.routing_name,
//   ds.signing_algorithm,
//   ds.qstring_ignore,
//   (SELECT o.protocol::text || \'://\' || o.fqdn || rtrim(concat(\':\', o.port::text), \':\')
//     FROM origin o
//     WHERE o.deliveryservice = ds.id
//     AND o.is_primary) as org_server_fqdn,
//   ds.origin_shield,
//   regex.pattern AS pattern,
//   retype.name AS re_type,
//   dstype.name AS ds_type,
//   cdn.domain_name AS domain_name,
//   ds.profile,
//   ds.protocol,
//   ds.ssl_key_version,
//   ds.range_request_handling,
//   ds.fq_pacing_rate,
//   ds.edge_header_rewrite,
//   ds.mid_header_rewrite,
//   ds.regex_remap,
//   ds.cacheurl,
//   ds.remap_text,
//   ds.multi_site_origin,
//   ds.multi_site_origin_algorithm
// FROM
//   deliveryservice ds
//   JOIN deliveryservice_regex ON deliveryservice_regex.deliveryservice = ds.id
//   JOIN regex ON deliveryservice_regex.regex = regex.id
//   JOIN type as retype ON regex.type = retype.id
//   JOIN type as dstype ON ds.type = dstype.id
//   JOIN cdn ON cdn.id = ds.cdn_id
// WHERE
//   ds.id IN (
//     SELECT DISTINCT deliveryservice
//     FROM deliveryservice_server
//     WHERE server IN (SELECT id FROM server WHERE profile = $1)
//   )
// `
// 	rows, err := tx.Query(qry, profileID)
// 	if err != nil {
// 		return nil, errors.New("querying: " + err.Error())
// 	}
// 	defer rows.Close()

// 	dses := []ProfileDS{}
// 	for rows.Next() {
// 		d := ProfileDS{}
// 		if err := rows.Scan(&d.ID, &d.XMLID, &d.DSCP, &d.RoutingName, &d.SigningAlgorithm, &d.QstringIgnore, &d.OriginServerFQDN, &d.OriginShield, &d.Pattern, &d.RegexType, &d.Type, &d.Domain, &d.Profile, &d.Protocol, &d.SslKeyVersion, &d.RangeRequestHandling, &d.FqPacingRate, &d.EdgeHeaderRewrite, &d.MidHeaderRewrite, &d.RegexRemap, &d.Cacheurl, &d.RemapText, &d.MultiSiteOrigin, &d.MultiSiteOriginAlgorithm); err != nil {
// 			return nil, errors.New("scanning: " + err.Error())
// 		}
// 		d.Type = tc.DeliveryServiceTypeFromString(d.Type)
// 		dses = append(dses, d)
// 	}
// 	return dses, nil
// }
