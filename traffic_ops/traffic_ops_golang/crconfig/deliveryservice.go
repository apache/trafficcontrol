package crconfig

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
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/lib/pq"
)

const CDNSOAMinimum = 30 * time.Second
const CDNSOAExpire = 604800 * time.Second
const CDNSOARetry = 7200 * time.Second
const CDNSOARefresh = 28800 * time.Second
const CDNSOAAdmin = "traffic_ops"
const DefaultTLDTTLSOA = 86400 * time.Second
const DefaultTLDTTLNS = 3600 * time.Second

const GeoProviderMaxmindStr = "maxmindGeolocationService"
const GeoProviderNeustarStr = "neustarGeolocationService"

func makeDSes(tx *sql.Tx, cdn string, domain string, live bool) (map[string]tc.CRConfigDeliveryService, error) {
	dses := map[string]tc.CRConfigDeliveryService{}

	admin := CDNSOAAdmin
	expireSecondsStr := strconv.Itoa(int(CDNSOAExpire / time.Second))
	minimumSecondsStr := strconv.Itoa(int(CDNSOAMinimum / time.Second))
	refreshSecondsStr := strconv.Itoa(int(CDNSOARefresh / time.Second))
	retrySecondsStr := strconv.Itoa(int(CDNSOARetry / time.Second))
	cdnSOA := &tc.SOA{
		Admin:          &admin,
		ExpireSeconds:  &expireSecondsStr,
		MinimumSeconds: &minimumSecondsStr,
		RefreshSeconds: &refreshSecondsStr,
		RetrySeconds:   &retrySecondsStr,
	}

	// Note the CRConfig omits acceptHTTP if it's true
	falsePtr := false
	protocol0 := &tc.CRConfigDeliveryServiceProtocol{AcceptHTTPS: false, RedirectOnHTTPS: false}
	protocol1 := &tc.CRConfigDeliveryServiceProtocol{AcceptHTTP: &falsePtr, AcceptHTTPS: true, RedirectOnHTTPS: false}
	protocol2 := &tc.CRConfigDeliveryServiceProtocol{AcceptHTTPS: true, RedirectOnHTTPS: false}
	protocol3 := &tc.CRConfigDeliveryServiceProtocol{AcceptHTTPS: true, RedirectOnHTTPS: true}
	protocolDefault := protocol0

	geoProvider0 := GeoProviderMaxmindStr
	geoProvider1 := GeoProviderNeustarStr
	geoProviderDefault := geoProvider0

	serverParams, err := getServerProfileParams(tx, cdn, live)
	if err != nil {
		return nil, errors.New("getting deliveryservice parameters: " + err.Error())
	}

	dsParams, err := getDSParams(serverParams)
	if err != nil {
		return nil, errors.New("getting deliveryservice server parameters: " + err.Error())
	}

	dsmatchsets, dsdomains, err := getDSRegexesDomains(tx, cdn, domain, live)
	if err != nil {
		return nil, errors.New("getting regex matchsets: " + err.Error())
	}

	staticDNSEntries, err := getStaticDNSEntries(tx, cdn, live)
	if err != nil {
		return nil, errors.New("getting static DNS entries: " + err.Error())
	}

	// TODO fix cdn_id subquery to get distinct latest
	qry := `
SELECT
  xml_id,
  miss_lat,
  miss_long,
  protocol,
  ttl,
  routing_name,
  geo_provider,
  type,
  geo_limit,
  geo_limit_countries,
  geolimit_redirect_url,
  initial_dispersion,
  regional_geo_blocking,
  tr_response_headers,
  max_dns_answers,
  profile,
  dns_bypass_ip,
  dns_bypass_ip6,
  dns_bypass_ttl,
  dns_bypass_cname,
  http_bypass_fqdn,
  ipv6_routing_enabled,
  deep_caching_type,
  ds_tr_request_headers,
  ds_tr_response_headers,
  anonymous_blocking_enabled,
	consistent_hash_regex,
	query_keys
FROM (
SELECT DISTINCT ON (d.xml_id)
  d.xml_id,
  d.miss_lat,
  d.miss_long,
  d.protocol,
  d.ccr_dns_ttl as ttl,
  d.routing_name,
  d.geo_provider,
  t.name as type,
  d.geo_limit,
  d.geo_limit_countries,
  d.geolimit_redirect_url,
  d.initial_dispersion,
  d.regional_geo_blocking,
  d.tr_response_headers,
  d.max_dns_answers,
  p.name as profile,
  d.dns_bypass_ip,
  d.dns_bypass_ip6,
  d.dns_bypass_ttl,
  d.dns_bypass_cname,
  d.http_bypass_fqdn,
  d.ipv6_routing_enabled,
  d.deep_caching_type,
  d.tr_request_headers as ds_tr_request_headers,
  d.tr_response_headers as ds_tr_response_headers,
  d.anonymous_blocking_enabled,
	d.consistent_hash_regex,
	(SELECT ARRAY_AGG(name ORDER BY name)
		FROM deliveryservice_consistent_hash_query_param
		WHERE deliveryservice_id = d.id) AS query_keys, -- TODO change to a join, for readability?
  d.deleted
FROM
  deliveryservice_snapshot d
  JOIN type_snapshot t ON t.id = d.type
  JOIN deliveryservice_snapshots dsn ON dsn.deliveryservice = d.xml_id
  LEFT OUTER JOIN profile_snapshot p ON p.id = d.profile `
	if !live {
		qry += ` AND p.last_updated <= dsn.time `
	}
	qry += `
WHERE
  d.cdn_id = (select id from cdn_snapshot c where c.name = $1 and c.last_updated <= dsn.time)
  AND d.active = true
  AND t.name != '` + string(tc.DSTypeAnyMap) + `'
`
	if !live {
		qry += `
  AND d.last_updated <= dsn.time
  AND t.last_updated <= dsn.time
`
	}
	qry += `
ORDER BY
  d.xml_id DESC,
  d.last_updated DESC,
  t.last_updated DESC,
  p.last_updated DESC
) v where deleted = false
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying deliveryservices: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		ds := tc.CRConfigDeliveryService{
			ConsistentHashQueryParams: []string{},
			Protocol:                  &tc.CRConfigDeliveryServiceProtocol{},
			ResponseHeaders:           map[string]string{},
			Soa:                       cdnSOA,
			TTLs:                      &tc.CRConfigTTL{},
		}

		missLat := sql.NullFloat64{}
		missLon := sql.NullFloat64{}
		protocol := sql.NullInt64{}
		ttl := sql.NullInt64{}
		geoProvider := sql.NullInt64{}
		ttype := ""
		geoLimit := sql.NullInt64{}
		geoLimitCountries := sql.NullString{}
		geoLimitRedirectURL := sql.NullString{}
		dispersion := sql.NullInt64{}
		geoBlocking := false
		trRespHdrsStr := sql.NullString{}
		xmlID := ""
		maxDNSAnswers := sql.NullInt64{}
		profile := sql.NullString{}
		dnsBypassIP := sql.NullString{}
		dnsBypassIP6 := sql.NullString{}
		dnsBypassTTL := sql.NullInt64{}
		dnsBypassCName := sql.NullString{}
		httpBypassFQDN := sql.NullString{}
		ip6RoutingEnabled := sql.NullBool{}
		deepCachingType := sql.NullString{}
		trRequestHeaders := sql.NullString{}
		trResponseHeaders := sql.NullString{}
		anonymousBlocking := false
		consistentHashRegex := sql.NullString{}
		err := rows.Scan(
			&xmlID,
			&missLat,
			&missLon,
			&protocol,
			&ds.TTL,
			&ds.RoutingName,
			&geoProvider,
			&ttype,
			&geoLimit,
			&geoLimitCountries,
			&geoLimitRedirectURL,
			&dispersion,
			&geoBlocking,
			&trRespHdrsStr,
			&maxDNSAnswers,
			&profile,
			&dnsBypassIP,
			&dnsBypassIP6,
			&dnsBypassTTL,
			&dnsBypassCName,
			&httpBypassFQDN,
			&ip6RoutingEnabled,
			&deepCachingType,
			&trRequestHeaders,
			&trResponseHeaders,
			&anonymousBlocking,
			&consistentHashRegex,
			pq.Array(&ds.ConsistentHashQueryParams),
		)
		if err != nil {
			return nil, errors.New("scanning deliveryservice: " + err.Error())
		}

		// TODO prevent (lat XOR lon) in the Tx and UI
		if missLat.Valid && missLon.Valid {
			ds.MissLocation = &tc.CRConfigLatitudeLongitudeShort{Lat: missLat.Float64, Lon: missLon.Float64}
		} else if missLat.Valid {
			log.Warnln("delivery service " + xmlID + " has miss latitude but not longitude: omitting miss lat-lon from CRConfig")
		} else if missLon.Valid {
			log.Warnln("delivery service " + xmlID + " has miss longitude but not latitude: omitting miss lat-lon from CRConfig")
		}
		if ttl.Valid {
			ttl := int(ttl.Int64)
			ds.TTL = &ttl
		}

		protocolStr := getProtocolStr(ttype)

		ds.Protocol = protocolDefault
		if protocol.Valid {
			switch protocol.Int64 {
			case 0:
				ds.Protocol = protocol0
			case 1:
				ds.Protocol = protocol1
			case 2:
				ds.Protocol = protocol2
			case 3:
				ds.Protocol = protocol3
			}
		}

		ds.GeoLocationProvider = &geoProviderDefault
		if geoProvider.Valid {
			switch geoProvider.Int64 {
			case 0:
				ds.GeoLocationProvider = &geoProvider0
			case 1:
				ds.GeoLocationProvider = &geoProvider1
			}
		}

		if ds.Protocol.AcceptHTTPS {
			ds.SSLEnabled = true
		}

		if deepCachingType.Valid {
			// TODO change to omit Valid check, default to the default DeepCachingType (NEVER). I'm pretty sure that's what should happen, but the Valid check emulates the old Perl CRConfig generation
			t := tc.DeepCachingTypeFromString(deepCachingType.String)
			ds.DeepCachingType = &t
		}

		ds.GeoLocationProvider = &geoProviderDefault

		if matchsets, ok := dsmatchsets[xmlID]; ok {
			ds.MatchSets = matchsets
		} else {
			log.Warnln("no regex matchsets for delivery service: " + xmlID)
		}
		if domains, ok := dsdomains[xmlID]; ok {
			ds.Domains = domains
		} else {
			log.Warnln("no host regex for delivery service: " + xmlID)
		}

		switch geoLimit.Int64 { // No Valid check - default false and set countries, if null
		case 0:
			ds.CoverageZoneOnly = false
		case 1:
			ds.CoverageZoneOnly = true
			if protocolStr == "HTTP" {
				ds.GeoLimitRedirectURL = &geoLimitRedirectURL.String // No Valid check - empty string, if null
			}
		default:
			ds.CoverageZoneOnly = false
			if protocolStr == "HTTP" {
				ds.GeoLimitRedirectURL = &geoLimitRedirectURL.String // No Valid check - empty string, if null
			}
			if geoLimitCountries.Valid {
				for _, code := range strings.Split(geoLimitCountries.String, ",") {
					ds.GeoEnabled = append(ds.GeoEnabled, tc.CRConfigGeoEnabled{CountryCode: strings.TrimSpace(code)})
				}
			}
		}

		nsSeconds := DefaultTLDTTLNS
		soaSeconds := DefaultTLDTTLSOA
		if profile.Valid {
			if sval, ok := dsParams["tld.ttls.SOA"]; ok {
				if val, err := strconv.Atoi(sval); err == nil {
					soaSeconds = time.Duration(val) * time.Second
				} else {
					log.Errorln("delivery service " + xmlID + " profile " + profile.String + " param tld.ttls.SOA '" + sval + "' not a number - skipping")
				}
			}
			if sval, ok := dsParams["tld.ttls.NS"]; ok {
				if val, err := strconv.Atoi(sval); err == nil {
					nsSeconds = time.Duration(val) * time.Second
				} else {
					log.Errorln("delivery service " + xmlID + " profile " + profile.String + " param tld.ttls.NS '" + sval + "' not a number - skipping")
				}
			}
		}
		nsSecondsStr := strconv.Itoa(int(nsSeconds / time.Second))
		soaSecondsStr := strconv.Itoa(int(soaSeconds / time.Second))
		ttlStr := ""
		if ds.TTL != nil {
			ttlStr = strconv.Itoa(*ds.TTL)
		}
		ds.TTLs = &tc.CRConfigTTL{
			ASeconds:    &ttlStr,
			AAAASeconds: &ttlStr,
			NSSeconds:   &nsSecondsStr,
			SOASeconds:  &soaSecondsStr,
		}

		if protocolStr == "DNS" {
			bypassDest := &tc.CRConfigBypassDestination{}
			if dnsBypassIP.String != "" {
				bypassDest.IP = &dnsBypassIP.String
			}
			if dnsBypassIP6.String != "" {
				bypassDest.IP6 = &dnsBypassIP6.String
			}
			if dnsBypassTTL.Valid {
				i := int(dnsBypassTTL.Int64)
				bypassDest.TTL = &i
			}
			if dnsBypassCName.Valid && dnsBypassCName.String != "" {
				bypassDest.CName = &dnsBypassCName.String
			}
			if *bypassDest != (tc.CRConfigBypassDestination{}) {
				if ds.BypassDestination == nil {
					ds.BypassDestination = map[string]*tc.CRConfigBypassDestination{}
				}
				ds.BypassDestination["DNS"] = bypassDest
			}
			if maxDNSAnswers.Valid {
				i := int(maxDNSAnswers.Int64)
				ds.MaxDNSIPsForLocation = &i
			}
		} else if protocolStr == "HTTP" {
			if httpBypassFQDN.String != "" {
				if ds.BypassDestination == nil {
					ds.BypassDestination = map[string]*tc.CRConfigBypassDestination{}
				}
				hostPort := strings.Split(httpBypassFQDN.String, ":")
				bypass := &tc.CRConfigBypassDestination{FQDN: &hostPort[0]}
				if len(hostPort) > 1 {
					bypass.Port = &hostPort[1]
				}
				ds.BypassDestination["HTTP"] = bypass
			}
			geoBlockingStr := "false"
			if geoBlocking {
				geoBlockingStr = "true"
			}
			ds.RegionalGeoBlocking = &geoBlockingStr

			anonymousBlockingStr := "false"
			if anonymousBlocking {
				anonymousBlockingStr = "true"
			}
			ds.AnonymousBlockingEnabled = &anonymousBlockingStr
			if dispersion.Valid {
				ds.Dispersion = &tc.CRConfigDispersion{Limit: int(dispersion.Int64), Shuffled: true}
			}
		}

		if consistentHashRegex.Valid && consistentHashRegex.String != "" {
			ds.ConsistentHashRegex = &consistentHashRegex.String
		}

		ds.IP6RoutingEnabled = &ip6RoutingEnabled.Bool // No Valid check, false if null

		if trResponseHeaders.Valid && trResponseHeaders.String != "" {
			trResponseHeaders.String = strings.Replace(trResponseHeaders.String, "__RETURN__", "\n", -1)
			hdrs := strings.Split(trResponseHeaders.String, "\n")
			for _, hdr := range hdrs {
				nameVal := strings.Split(hdr, `:`)
				name := strings.TrimSpace(nameVal[0])
				val := ""
				if len(nameVal) > 1 {
					val = strings.Trim(nameVal[1], " \n\"")
				}
				ds.ResponseHeaders[name] = val
			}
		}

		if trRequestHeaders.Valid && trRequestHeaders.String != "" {
			trRequestHeaders.String = strings.Replace(trRequestHeaders.String, "__RETURN__", "\n", -1)
			hdrs := strings.Split(trRequestHeaders.String, "\n")
			for _, hdr := range hdrs {
				nameVal := strings.Split(hdr, `:`)
				name := strings.TrimSpace(nameVal[0])
				ds.RequestHeaders = append(ds.RequestHeaders, name)
			}
		}

		ds.StaticDNSEntries = staticDNSEntries[tc.DeliveryServiceName(xmlID)]

		dses[xmlID] = ds
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("iterating deliveryservice rows: " + err.Error())
	}

	return dses, nil
}

func getStaticDNSEntries(tx *sql.Tx, cdn string, live bool) (map[tc.DeliveryServiceName][]tc.CRConfigStaticDNSEntry, error) {
	entries := map[tc.DeliveryServiceName][]tc.CRConfigStaticDNSEntry{}

	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
)
SELECT ds, name, ttl, value, type FROM (
SELECT DISTINCT ON (e.host, e.address, e.deliveryservice, e.cachegroup)
  d.xml_id as ds,
  e.host as name,
  e.ttl,
  e.address as value,
  t.name as type,
  e.deleted
FROM
  staticdnsentry_snapshot e
  JOIN deliveryservice_snapshot d on d.id = e.deliveryservice
  JOIN type_snapshot t on t.id = e.type
  JOIN deliveryservice_snapshots dsn ON dsn.deliveryservice = d.xml_id
WHERE
  d.cdn_id = (select id from cdn_snapshot c where c.name = (select v from cdn_name) and c.last_updated <= dsn.time)
  AND d.active = true
`
	if !live {
		qry += `
  AND e.last_updated <= dsn.time
  AND d.last_updated <= dsn.time
  AND t.last_updated <= dsn.time
`
	}
	qry += `
ORDER BY
  e.host DESC,
  e.address DESC,
  e.deliveryservice DESC,
  e.cachegroup DESC,
  e.last_updated DESC,
  d.last_updated DESC,
  t.last_updated DESC
) v where deleted = false
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying static DNS entries: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		ds := ""
		name := ""
		ttl := 0
		value := ""
		ttype := ""
		if err := rows.Scan(&ds, &name, &ttl, &value, &ttype); err != nil {
			return nil, errors.New("scanning static DNS entries: " + err.Error())
		}
		ttype = strings.Replace(ttype, "_RECORD", "", -1)
		entries[tc.DeliveryServiceName(ds)] = append(entries[tc.DeliveryServiceName(ds)], tc.CRConfigStaticDNSEntry{Name: name, TTL: ttl, Value: value, Type: ttype})
	}
	return entries, nil
}

func getProtocolStr(dsType string) string {
	if strings.HasPrefix(dsType, "DNS") {
		return "DNS"
	}
	return "HTTP"
}

func getDSRegexesDomains(tx *sql.Tx, cdn string, domain string, live bool) (map[string][]*tc.MatchSet, map[string][]string, error) {
	dsmatchsets := map[string][]*tc.MatchSet{}
	domains := map[string][]string{}
	patternToHostReplacer := strings.NewReplacer(`\`, ``, `.*`, ``, `.`, ``)
	qry := `
SELECT pattern, type, dstype, set_number, dsname FROM (
SELECT DISTINCT ON (dsname, pattern, type, set_number)
  r.pattern,
  t.name as type,
  dt.name as dstype,
  COALESCE(dr.set_number, 0) as set_number,
  d.xml_id as dsname,
  r.deleted
FROM
  regex_snapshot as r
  JOIN deliveryservice_regex_snapshot dr on r.id = dr.regex
  JOIN deliveryservice_snapshot d on d.id = dr.deliveryservice
  JOIN type_snapshot t on t.id = r.type
  JOIN type_snapshot dt on dt.id = d.type
  JOIN deliveryservice_snapshots dsn ON dsn.deliveryservice = d.xml_id
WHERE
  d.cdn_id = (select id from cdn_snapshot c where c.name = $1 and c.last_updated <= dsn.time)
  AND d.active = true
`
	if !live {
		qry += `
  AND r.last_updated <= dsn.time
  AND dr.last_updated <= dsn.time
  AND d.last_updated <= dsn.time
  AND t.last_updated <= dsn.time
  AND dt.last_updated <= dsn.time
`
	}
	qry += `
ORDER BY
  dsname DESC,
  pattern DESC,
  type DESC,
  set_number DESC,
  dr.last_updated DESC,
  dr.last_updated DESC,
  d.last_updated DESC,
  t.last_updated DESC,
  dt.last_updated DESC
) v where deleted = false
ORDER BY set_number ASC
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, nil, errors.New("querying deliveryservices: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		pattern := ""
		ttype := ""
		dstype := ""
		setnum := 0
		dsname := ""
		if err := rows.Scan(&pattern, &ttype, &dstype, &setnum, &dsname); err != nil {
			return nil, nil, errors.New("scanning deliveryservice regexes: " + err.Error())
		}

		protocolStr := getProtocolStr(dstype)

		for len(dsmatchsets[dsname]) <= setnum {
			dsmatchsets[dsname] = append(dsmatchsets[dsname], nil) // TODO change to not insert empties? Current behavior emulates old Perl CRConfig
		}

		matchType := ""
		switch ttype {
		case "HOST_REGEXP":
			matchType = "HOST"
		case "PATH_REGEXP":
			matchType = "PATH"
		case "HEADER_REGEXP":
			matchType = "HEADER"
		default:
			log.Infoln("unknown delivery service '" + dsname + "' regex type: " + ttype + " - skipping") // info, not warn or err, because this is normal for STEERING_REGEXP (and maybe others in the future)
			continue
		}

		if dsmatchsets[dsname][setnum] == nil {
			dsmatchsets[dsname][setnum] = &tc.MatchSet{}
		}
		matchset := dsmatchsets[dsname][setnum]
		matchset.Protocol = protocolStr
		matchset.MatchList = append(matchset.MatchList, tc.MatchList{MatchType: matchType, Regex: pattern})

		if ttype == "HOST_REGEXP" && setnum == 0 {
			domains[dsname] = append(domains[dsname], patternToHostReplacer.Replace(pattern)+"."+domain)
		}
	}
	return dsmatchsets, domains, nil
}

// getDSParams takes a map[serverProfile][paramName]paramVal and returns a map[paramName]paramVal.
// The returned map of parameter values is used for DS settings for the current CDN.
// If any profiles have conflicting parameters, an error is returned.
func getDSParams(serverParams map[string]map[string]string) (map[string]string, error) {
	dsParamNames := map[string]struct{}{
		"tld.soa.admin":     struct{}{},
		"tld.soa.expire":    struct{}{},
		"tld.soa.minimum":   struct{}{},
		"tld.soa.refresh":   struct{}{},
		"tld.soa.retry":     struct{}{},
		"tld.ttls.SOA":      struct{}{},
		"tld.ttls.NS":       struct{}{},
		"LogRequestHeaders": struct{}{},
	}
	dsParams := map[string]string{}
	dsParamsOriginalProfile := map[string]string{} // map[paramName]profile - used exclusively for the error message
	for profile, profileParams := range serverParams {
		for paramName, _ := range dsParamNames {
			paramVal, profileHasParam := profileParams[paramName]
			if !profileHasParam {
				continue
			}
			if dsParamVal, ok := dsParams[paramName]; ok && dsParamVal != paramVal {
				return nil, errors.New("profiles " + profile + " and " + dsParamsOriginalProfile[paramName] + " have conflicting values '" + paramVal + "' and '" + dsParamVal + "'")
			}
			dsParams[paramName] = paramVal
			dsParamsOriginalProfile[paramName] = profile
		}
	}
	return dsParams, nil
}

// getDSProfileParams returns a map[dsname]map[paramname]paramvalue
func getServerProfileParams(tx *sql.Tx, cdn string, live bool) (map[string]map[string]string, error) {
	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
),
snapshot_time AS (
  SELECT time as v FROM snapshot sn where sn.cdn = (SELECT v from cdn_name)
)
SELECT name, value, profile FROM (
SELECT DISTINCT ON (pa.name, pa.value, pr.name)
  pa.name,
  pa.value,
  pr.name as profile,
  pa.deleted
FROM
  profile_snapshot pr
  JOIN profile_parameter_snapshot as pp on pp.profile = pr.id
  JOIN parameter_snapshot pa on pa.id = pp.parameter
WHERE
  pr.id in (
    SELECT DISTINCT(profile) FROM (
    SELECT DISTINCT ON (s.ip_address, profile)
      s.profile,
      s.deleted
    FROM
      server_snapshot s
    WHERE
      s.cdn_id = (SELECT id FROM cdn_snapshot c WHERE c.name = (SELECT v from cdn_name) AND c.last_updated <= (SELECT v from snapshot_time))
`
	if !live {
		qry += `AND s.last_updated <= (SELECT v from snapshot_time)`
	}
	qry += `
    ORDER BY
      s.ip_address DESC,
      s.profile DESC,
      s.last_updated DESC
    ) v where deleted = false
  )
`
	if !live {
		qry += `
  AND pr.last_updated <= (SELECT v from snapshot_time)
  AND pp.last_updated <= (SELECT v from snapshot_time)
  AND pa.last_updated <= (SELECT v from snapshot_time)
`
	}
	qry += `
ORDER BY
  pa.name DESC,
  pa.value DESC,
  pr.name DESC,
  pr.last_updated DESC,
  pp.last_updated DESC,
  pa.last_updated DESC
) v where deleted = false
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying deliveryservices: " + err.Error())
	}
	defer rows.Close()

	params := map[string]map[string]string{}
	debugCount := 0
	for rows.Next() {
		debugCount++
		name := ""
		val := ""
		profile := ""
		if err := rows.Scan(&name, &val, &profile); err != nil {
			return nil, errors.New("scanning deliveryservice parameters: " + err.Error())
		}
		if _, ok := params[profile]; !ok {
			params[profile] = map[string]string{}
		}
		params[profile][name] = val
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("iterating deliveryservice parameter rows: " + err.Error())
	}
	return params, nil
}
