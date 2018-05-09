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

// getServerDSesModifieds returns a map[ds]time of the lastest modified time for any server on each DS.
func getServerDSesModifieds(serverDSNames map[tc.CacheName][]ServerDS) map[tc.DeliveryServiceName]time.Time {
	dsModifieds := map[tc.DeliveryServiceName]time.Time{}
	for _, dses := range serverDSNames {
		for _, ds := range dses {
			if ds.Modified.After(dsModifieds[ds.DS]) {
				dsModifieds[ds.DS] = ds.Modified
			}
		}
	}
	return dsModifieds
}

func makeDSes(cdn string, domain string, serverDSNames map[tc.CacheName][]ServerDS, tx *sql.Tx) (map[string]tc.CRConfigDeliveryService, error) {
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

	q := `
SELECT
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
d.tr_request_headers,
d.tr_response_headers,
d.anonymous_blocking_enabled,
d.last_updated
FROM deliveryservice as d
JOIN type as t ON t.id = d.type
LEFT OUTER JOIN profile as p ON p.id = d.profile
WHERE d.cdn_id = (SELECT id FROM cdn WHERE name = $1)
AND d.active = true
AND t.name != '` + string(tc.DSTypeAnyMap) + `'
`

	serverParams, err := getServerProfileParams(cdn, tx)
	if err != nil {
		return nil, errors.New("getting deliveryservice parameters: " + err.Error())
	}
	dsParams, err := getDSParams(serverParams)
	if err != nil {
		return nil, errors.New("getting deliveryservice server parameters: " + err.Error())
	}
	dsmatchsets, dsdomains, dsMatchsetsNewestModifieds, err := getDSRegexesDomains(cdn, domain, tx)
	if err != nil {
		return nil, errors.New("getting regex matchsets: " + err.Error())
	}
	staticDNSEntries, staticDNSEntriesNewestModifieds, err := getStaticDNSEntries(cdn, tx)
	if err != nil {
		return nil, errors.New("getting static DNS entries: " + err.Error())
	}

	dsServerModifieds := getServerDSesModifieds(serverDSNames)

	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, errors.New("querying deliveryservices: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		ds := tc.CRConfigDeliveryService{
			CRConfigDeliveryServiceV11: tc.CRConfigDeliveryServiceV11{
				Protocol:        &tc.CRConfigDeliveryServiceProtocol{},
				ResponseHeaders: map[string]string{},
				Soa:             cdnSOA,
				TTLs:            &tc.CRConfigTTL{},
			},
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
		if err := rows.Scan(&xmlID, &missLat, &missLon, &protocol, &ds.TTL, &ds.RoutingName, &geoProvider, &ttype, &geoLimit, &geoLimitCountries, &geoLimitRedirectURL, &dispersion, &geoBlocking, &trRespHdrsStr, &maxDNSAnswers, &profile, &dnsBypassIP, &dnsBypassIP6, &dnsBypassTTL, &dnsBypassCName, &httpBypassFQDN, &ip6RoutingEnabled, &deepCachingType, &trRequestHeaders, &trResponseHeaders, &anonymousBlocking, &ds.Modified); err != nil {
			return nil, errors.New("scanning deliveryservice: " + err.Error())
		}

		ds.AnyModified = ds.Modified

		ds.ServersModified = dsServerModifieds[tc.DeliveryServiceName(xmlID)]

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
			if matchsetModified, ok := dsMatchsetsNewestModifieds[xmlID]; ok && matchsetModified.After(ds.AnyModified) {
				ds.AnyModified = matchsetModified
			}
		} else {
			log.Warnln("no regex matchsets for delivery service: " + xmlID)
		}
		if domains, ok := dsdomains[xmlID]; ok {
			ds.Domains = domains
			// don't need to check domain modified here, since it's the same SQL data as matchsets - there may be a matchset without a domain, but there will never be a domain without a matchset
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
				if val, err := strconv.Atoi(sval.Val); err == nil {
					soaSeconds = time.Duration(val) * time.Second
					if sval.Modified.After(ds.AnyModified) {
						ds.AnyModified = sval.Modified
					}
				} else {
					log.Errorln("delivery service " + xmlID + " profile " + profile.String + " param tld.ttls.SOA '" + sval.Val + "' not a number - skipping")
				}
			}
			if sval, ok := dsParams["tld.ttls.NS"]; ok {
				if val, err := strconv.Atoi(sval.Val); err == nil {
					nsSeconds = time.Duration(val) * time.Second
					if sval.Modified.After(ds.AnyModified) {
						ds.AnyModified = sval.Modified
					}
				} else {
					log.Errorln("delivery service " + xmlID + " profile " + profile.String + " param tld.ttls.NS '" + sval.Val + "' not a number - skipping")
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

		ds.IP6RoutingEnabled = &ip6RoutingEnabled.Bool // No Valid check, false if null

		if trResponseHeaders.Valid && trResponseHeaders.String != "" {
			hdrs := strings.Split(trResponseHeaders.String, `__RETURN__`)
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
			hdrs := strings.Split(trRequestHeaders.String, `__RETURN__`)
			for _, hdr := range hdrs {
				nameVal := strings.Split(hdr, `:`)
				name := strings.TrimSpace(nameVal[0])
				ds.RequestHeaders = append(ds.RequestHeaders, name)
			}
		}

		ds.StaticDNSEntries = staticDNSEntries[tc.DeliveryServiceName(xmlID)]
		if sdeModified, ok := staticDNSEntriesNewestModifieds[tc.DeliveryServiceName(xmlID)]; ok && sdeModified.After(ds.AnyModified) {
			ds.AnyModified = sdeModified
		}

		dses[xmlID] = ds
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("iterating deliveryservice rows: " + err.Error())
	}

	return dses, nil
}

func getStaticDNSEntries(cdn string, tx *sql.Tx) (map[tc.DeliveryServiceName][]tc.CRConfigStaticDNSEntry, map[tc.DeliveryServiceName]time.Time, error) {
	entries := map[tc.DeliveryServiceName][]tc.CRConfigStaticDNSEntry{}
	newestModified := map[tc.DeliveryServiceName]time.Time{}
	q := `
 SELECT d.xml_id as ds, e.host as name, e.ttl, e.address as value, t.name as type, greatest(e.last_updated, d.last_updated, t.last_updated, cdn.last_updated)
FROM staticdnsentry as e
JOIN deliveryservice as d on d.id = e.deliveryservice
JOIN type as t on t.id = e.type
JOIN cdn on cdn.id = d.cdn_id
WHERE cdn.name = $1
AND d.active = true
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, nil, errors.New("querying static DNS entries: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		ds := ""
		name := ""
		ttl := 0
		value := ""
		ttype := ""
		modified := time.Time{}
		if err := rows.Scan(&ds, &name, &ttl, &value, &ttype, &modified); err != nil {
			return nil, nil, errors.New("scanning static DNS entries: " + err.Error())
		}
		ttype = strings.Replace(ttype, "_RECORD", "", -1)
		entries[tc.DeliveryServiceName(ds)] = append(entries[tc.DeliveryServiceName(ds)], tc.CRConfigStaticDNSEntry{
			Name:  name,
			TTL:   ttl,
			Value: value,
			Type:  ttype,
		})
		if modified.After(newestModified[tc.DeliveryServiceName(ds)]) {
			newestModified[tc.DeliveryServiceName(ds)] = modified
		}
	}
	return entries, newestModified, nil
}

func getProtocolStr(dsType string) string {
	if strings.HasPrefix(dsType, "DNS") {
		return "DNS"
	}
	return "HTTP"
}

// getDSRegexesDomains returns a map[ds][]matchests, a map[ds][]domain, and a map[ds]lastestModifiedTime
func getDSRegexesDomains(cdn string, domain string, tx *sql.Tx) (map[string][]*tc.MatchSet, map[string][]string, map[string]time.Time, error) {
	dsmatchsets := map[string][]*tc.MatchSet{}
	domains := map[string][]string{}
	modifieds := map[string]time.Time{}
	patternToHostReplacer := strings.NewReplacer(`\`, ``, `.*`, ``, `.`, ``)
	q := `
SELECT r.pattern, t.name as type, dt.name as dstype, COALESCE(dr.set_number, 0), d.xml_id as dsname, greatest(r.last_updated, d.last_updated, dr.last_updated, t.last_updated, dt.last_updated, cdn.last_updated) as modified
FROM regex as r
JOIN deliveryservice_regex as dr ON r.id = dr.regex
JOIN deliveryservice as d ON d.id = dr.deliveryservice
JOIN type as t ON t.id = r.type
JOIN type as dt ON dt.id = d.type
JOIN cdn ON d.cdn_id = cdn.id
WHERE cdn.name = $1
AND d.active = true
ORDER BY dr.set_number asc
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, nil, nil, errors.New("querying deliveryservices: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		pattern := ""
		ttype := ""
		dstype := ""
		setnum := 0
		dsname := ""
		modified := time.Time{}
		if err := rows.Scan(&pattern, &ttype, &dstype, &setnum, &dsname, &modified); err != nil {
			return nil, nil, nil, errors.New("scanning deliveryservice regexes: " + err.Error())
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
		modifieds[dsname] = modified
	}
	return dsmatchsets, domains, modifieds, nil
}

// getDSParams takes a map[serverProfile][paramName]paramVal and returns a map[paramName]paramVal.
// The returned map of parameter values is used for DS settings for the current CDN.
// If any profiles have conflicting parameters, an error is returned.
func getDSParams(serverParams map[string]map[string]ParamVal) (map[string]ParamVal, error) {
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
	dsParams := map[string]ParamVal{}
	dsParamsOriginalProfile := map[string]string{} // map[paramName]profile - used exclusively for the error message
	for profile, profileParams := range serverParams {
		for paramName, _ := range dsParamNames {
			paramVal, profileHasParam := profileParams[paramName]
			if !profileHasParam {
				continue
			}
			if dsParamVal, ok := dsParams[paramName]; ok && dsParamVal.Val != paramVal.Val {
				return nil, errors.New("profiles " + profile + " and " + dsParamsOriginalProfile[paramName] + " have conflicting values '" + paramVal.Val + "' and '" + dsParamVal.Val + "'")
			}
			dsParams[paramName] = paramVal
			dsParamsOriginalProfile[paramName] = profile
		}
	}
	return dsParams, nil
}

type ParamVal struct {
	Val      string
	Modified time.Time
}

// getDSProfileParams returns a map[dsname]map[paramname]val
func getServerProfileParams(cdn string, tx *sql.Tx) (map[string]map[string]ParamVal, error) {
	q := `
SELECT parameter.name, parameter.value, profile.name as profile, greatest(parameter.last_updated, profile.last_updated, pp.last_updated, server.last_updated, cdn.last_updated) as modified
FROM profile
JOIN profile_parameter as pp ON pp.profile = profile.id
JOIN parameter ON parameter.id = pp.parameter
JOIN server ON server.profile = profile.id
JOIN cdn ON cdn.id = server.id
WHERE cdn.name = $1
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, errors.New("querying deliveryservices: " + err.Error())
	}
	defer rows.Close()

	params := map[string]map[string]ParamVal{}
	debugCount := 0
	for rows.Next() {
		debugCount++
		name := ""
		val := ""
		profile := ""
		modified := time.Time{}
		if err := rows.Scan(&name, &val, &profile, &modified); err != nil {
			return nil, errors.New("scanning deliveryservice parameters: " + err.Error())
		}
		if _, ok := params[profile]; !ok {
			params[profile] = map[string]ParamVal{}
		}
		params[profile][name] = ParamVal{Val: val, Modified: modified}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("iterating deliveryservice parameter rows: " + err.Error())
	}
	return params, nil
}
