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
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/lib/pq"
)

// RemapDotConfigIncludeInactiveDeliveryServices is whether delivery services with 'active' false are included in the remap.config.
const RemapDotConfigIncludeInactiveDeliveryServices = true

// getProfileData returns the necessary info about the profile, whether it exists, and any error.
func getProfileData(tx *sql.Tx, id int) (ProfileData, bool, error) {
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

// GetProfilesParamData returns a map[profileID][paramName]paramVal
func GetProfilesParamData(tx *sql.Tx, profileIDs []int, configFile string) (map[int]map[string]string, error) {
	qry := `
SELECT
  pr.id,
  p.name,
  p.value
FROM
  profile pr
  JOIN profile_parameter pp on pr.id = pp.profile
  JOIN parameter p on p.id = pp.parameter
WHERE
  pr.id = ANY($1)
  AND p.config_file = $2
`

	rows, err := tx.Query(qry, pq.Array(profileIDs), configFile)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	profilesParams := map[int]map[string]string{}
	for rows.Next() {
		profileID := 0
		name := ""
		val := ""
		if err := rows.Scan(&profileID, &name, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if _, ok := profilesParams[profileID]; !ok {
			profilesParams[profileID] = map[string]string{}
		}
		profilesParams[profileID][name] = val
	}
	return profilesParams, nil
}

func GetProfileParamData(tx *sql.Tx, profileID int, configFile string) (map[string]string, error) {
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

		if _, ok := params[name]; ok {
			log.Warnf("Profile %v has multiple parameters '%v' assigned! ATS config generation ignoring value '%v'!", profileID, name, params[name])
		}

		params[name] = val
	}
	return params, nil
}

func GetServerProfileParamData(tx *sql.Tx, serverName tc.CacheName, configFile string) (map[string]string, error) {
	qry := `
SELECT
  p.name,
  p.value
FROM
  parameter p
  JOIN profile_parameter pp on p.id = pp.parameter
  JOIN server s on s.profile = pp.profile
WHERE
  s.host_name = $1
  AND p.config_file = $2
`
	rows, err := tx.Query(qry, serverName, configFile)
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

		if _, ok := params[name]; ok {
			log.Warnf("Server %v profile has multiple parameters '%v' assigned! ATS config generation ignoring value '%v'!", serverName, name, params[name])
		}

		params[name] = val
	}
	return params, nil
}

// GetATSMajorVersion returns the major version of the given profile's package trafficserver parameter.
// If no parameter exists, this does not return an error, but rather logs a warning and uses DefaultATSVersion.
func GetATSMajorVersionFromServerName(tx *sql.Tx, serverName tc.CacheName) (int, error) {
	atsVersion, _, err := GetServerProfileParamValue(tx, serverName, "package", "trafficserver")
	if err != nil {
		return 0, errors.New("getting profile param value: " + err.Error())
	}
	if len(atsVersion) == 0 {
		atsVersion = atscfg.DefaultATSVersion
		log.Warnln("Parameter package.trafficserver missing for server " + string(serverName) + " profile. Assuming version " + atsVersion)
	}

	atsMajorVer, err := atscfg.GetATSMajorVersionFromATSVersion(atsVersion)
	if err != nil {
		return 0, errors.New("ats version parameter '" + atsVersion + "' on this profile is not a number (config_file 'package', name 'trafficserver')")
	}
	return atsMajorVer, nil
}

type ProfileData struct {
	ID   int
	Name string
}

// GetProfileData returns the necessary info about the profile, whether it exists, and any error.
func GetProfileData(tx *sql.Tx, id int) (ProfileData, bool, error) {
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

func GetProfileDS(tx *sql.Tx, profileID int) ([]atscfg.ProfileDS, error) {
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

	dses := []atscfg.ProfileDS{}
	for rows.Next() {
		d := atscfg.ProfileDS{}
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

// GetServerProfileParamValue gets the value of a parameter assigned to a server's Profile, by name and config file.
// Returns the parameter, whether it existed, and any error.
func GetServerProfileParamValue(tx *sql.Tx, serverName tc.CacheName, configFile string, name string) (string, bool, error) {
	qry := `
SELECT
  p.value
FROM
  parameter p
  JOIN profile_parameter pp ON p.id = pp.parameter
  JOIN server s on s.profile = pp.profile
WHERE
  s.host_name = $1
  AND p.config_file = $2
  AND p.name = $3
`
	val := ""
	if err := tx.QueryRow(qry, serverName, configFile, name).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying: " + err.Error())
	}
	return val, true, nil
}

// GetProfileIDFromName returns the profile's ID, whether it exists, and any error.
func GetProfileIDFromName(tx *sql.Tx, profileName string) (int, bool, error) {
	qry := `SELECT id from profile where name = $1`
	id := 0
	if err := tx.QueryRow(qry, profileName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, errors.New("querying: " + err.Error())
	}
	return id, true, nil
}

type Parameter struct {
	Name       string
	ConfigFile string
	Value      string
}

func GetParamsByName(tx *sql.Tx, paramName string) ([]Parameter, error) {
	qry := `
SELECT
  p.value,
  p.config_file
FROM
  parameter p
WHERE
  p.name = $1
`
	rows, err := tx.Query(qry, paramName)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := []Parameter{}
	for rows.Next() {
		pa := Parameter{Name: paramName}
		if err := rows.Scan(&pa.Value, &pa.ConfigFile); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		params = append(params, pa)
	}
	return params, nil
}

func GetServerParamData(tx *sql.Tx, profileID int, configFile string, serverHost string, serverDomain string) (map[string]string, error) {
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

func GetRemapDSData(tx *sql.Tx, serverInfo *atscfg.ServerInfo) ([]atscfg.RemapConfigDSData, error) {
	if tc.CacheTypeFromString(serverInfo.Type) == tc.CacheTypeMid {
		return GetRemapDSDataForMid(tx, serverInfo)
	} else {
		return GetRemapDSDataForEdge(tx, serverInfo)
	}
}

const RemapDSDataQuerySelectFrom = `
SELECT
  ds.xml_id,
  ds.id AS ds_id,
  ds.dscp,
  ds.routing_name,
  ds.signing_algorithm,
  ds.qstring_ignore,
  (SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
    FROM origin o
    WHERE o.deliveryservice = ds.id
    AND o.is_primary) as org_server_fqdn,
  ds.multi_site_origin,
  ds.range_request_handling,
  ds.fq_pacing_rate,
  ds.origin_shield,
  r.pattern,
  retype.name AS re_type,
  dstype.name AS ds_type,
  cdn.domain_name AS domain_name,
  dsr.set_number,
  ds.edge_header_rewrite,
  ds.mid_header_rewrite,
  ds.regex_remap,
  ds.cacheurl,
  ds.remap_text,
  ds.protocol,
  ds.profile,
  ds.anonymous_blocking_enabled,
  ds.active
FROM
  deliveryservice ds
  JOIN deliveryservice_regex dsr ON dsr.deliveryservice = ds.id
  JOIN regex r ON dsr.regex = r.id
  JOIN type retype ON r.type = retype.id
  JOIN type dstype ON ds.type = dstype.id
  JOIN cdn ON cdn.id = ds.cdn_id
`

const RemapDSDataQueryWhereForMid = `
WHERE
  cdn.name = $1
  AND ds.id in (SELECT dss.deliveryservice FROM deliveryservice_server dss)
  AND ds.active = true
`

const RemapDSDataQueryWhereForEdge = `
JOIN deliveryservice_server dss ON dss.deliveryservice = ds.id
WHERE dss.server = $1
`

const RemapDSDataQueryOrderBy = `
ORDER BY
  ds_id,
  re_type,
  set_number
`

func GetRemapDSDataForMid(tx *sql.Tx, serverInfo *atscfg.ServerInfo) ([]atscfg.RemapConfigDSData, error) {
	qry := RemapDSDataQuerySelectFrom + RemapDSDataQueryWhereForMid + RemapDSDataQueryOrderBy
	rows, err := tx.Query(qry, serverInfo.CDN)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := []atscfg.RemapConfigDSData{}
	for rows.Next() {
		d := atscfg.RemapConfigDSData{}
		if err := rows.Scan(&d.Name, &d.ID, &d.DSCP, &d.RoutingName, &d.SigningAlgorithm, &d.QStringIgnore, &d.OriginFQDN, &d.MultiSiteOrigin, &d.RangeRequestHandling, &d.FQPacingRate, &d.OriginShield, &d.Pattern, &d.RegexType, &d.Type, &d.Domain, &d.RegexSetNumber, &d.EdgeHeaderRewrite, &d.MidHeaderRewrite, &d.RegexRemap, &d.CacheURL, &d.RemapText, &d.Protocol, &d.ProfileID, &d.AnonymousBlockingEnabled, &d.Active); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if !RemapDotConfigIncludeInactiveDeliveryServices && !d.Active {
			continue
		}
		d.Type = tc.DSTypeFromString(string(d.Type))
		dses = append(dses, d)
	}
	return dses, nil
}

func GetRemapDSDataForEdge(tx *sql.Tx, server *atscfg.ServerInfo) ([]atscfg.RemapConfigDSData, error) {
	qry := RemapDSDataQuerySelectFrom + RemapDSDataQueryWhereForEdge + RemapDSDataQueryOrderBy
	rows, err := tx.Query(qry, server.ID)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := []atscfg.RemapConfigDSData{}
	for rows.Next() {
		d := atscfg.RemapConfigDSData{}
		if err := rows.Scan(&d.Name, &d.ID, &d.DSCP, &d.RoutingName, &d.SigningAlgorithm, &d.QStringIgnore, &d.OriginFQDN, &d.MultiSiteOrigin, &d.RangeRequestHandling, &d.FQPacingRate, &d.OriginShield, &d.Pattern, &d.RegexType, &d.Type, &d.Domain, &d.RegexSetNumber, &d.EdgeHeaderRewrite, &d.MidHeaderRewrite, &d.RegexRemap, &d.CacheURL, &d.RemapText, &d.Protocol, &d.ProfileID, &d.AnonymousBlockingEnabled, &d.Active); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if !RemapDotConfigIncludeInactiveDeliveryServices && !d.Active {
			continue
		}
		d.Type = tc.DSTypeFromString(string(d.Type))
		dses = append(dses, d)
	}
	return dses, nil
}

func GetServerNameFromID(tx *sql.Tx, id int) (tc.CacheName, bool, error) {
	qry := `SELECT s.host_name FROM server s WHERE s.id = $1`
	name := tc.CacheName("")
	if err := tx.QueryRow(qry, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying: " + err.Error())
	}
	return name, true, nil
}

// GetServerNameFromNameOrID returns the server name from a parameter which may be the name or ID.
// This also checks and verifies the existence of the given server, and returns an appropriate user error if it doesn't exist.
// Returns the name, any user error, any system error, and any error code.
func GetServerNameFromNameOrID(tx *sql.Tx, serverNameOrID string) (tc.CacheName, error, error, int) {
	if serverID, err := strconv.Atoi(serverNameOrID); err == nil {
		serverName, ok, err := dbhelpers.GetServerNameFromID(tx, serverID)
		if err != nil {
			return "", nil, fmt.Errorf("getting server name from id %v: %v", serverID, err), http.StatusInternalServerError
		} else if !ok {
			return "", errors.New("server not found"), nil, http.StatusNotFound
		}
		return tc.CacheName(serverName), nil, nil, http.StatusOK
	}

	serverName := tc.CacheName(serverNameOrID)
	if _, ok, err := dbhelpers.GetServerIDFromName(string(serverName), tx); err != nil {
		return "", nil, fmt.Errorf("checking server name '%v' existence: %v", serverName, err), http.StatusInternalServerError
	} else if !ok {
		return "", errors.New("server not found"), nil, http.StatusNotFound
	}
	return serverName, nil, nil, http.StatusOK
}

// GetServerInfoByID returns the necessary info about the server, whether the server exists, and any error.
func GetServerInfoByID(tx *sql.Tx, id int) (*atscfg.ServerInfo, bool, error) {
	return getServerInfo(tx, ServerInfoQuery()+`WHERE s.id = $1`, []interface{}{id})
}

// GetServerInfoByHost returns the necessary info about the server, whether the server exists, and any error.
func GetServerInfoByHost(tx *sql.Tx, host tc.CacheName) (*atscfg.ServerInfo, bool, error) {
	return getServerInfo(tx, ServerInfoQuery()+` WHERE s.host_name = $1 `, []interface{}{host})
}

// getServerInfo returns the necessary info about the server, whether the server exists, and any error.
func getServerInfo(tx *sql.Tx, qry string, qryParams []interface{}) (*atscfg.ServerInfo, bool, error) {
	s := atscfg.ServerInfo{}
	if err := tx.QueryRow(qry, qryParams...).Scan(&s.CDN, &s.CDNID, &s.ID, &s.HostName, &s.DomainName, &s.IP, &s.ProfileID, &s.ProfileName, &s.Port, &s.HTTPSPort, &s.Type, &s.CacheGroupID, &s.ParentCacheGroupID, &s.SecondaryParentCacheGroupID, &s.ParentCacheGroupType, &s.SecondaryParentCacheGroupType); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, errors.New("querying server info: " + err.Error())
	}
	return &s, true, nil
}

func ServerInfoQuery() string {
	return `
SELECT
  c.name as cdn,
  s.cdn_id,
  s.id,
  s.host_name,
  c.domain_name,
  s.ip_address,
  s.profile AS profile_id,
  p.name AS profile_name,
  s.tcp_port,
  s.https_port,
  t.name as type,
  s.cachegroup,
  COALESCE(cg.parent_cachegroup_id, ` + strconv.Itoa(atscfg.InvalidID) + `) as parent_cachegroup_id,
  COALESCE(cg.secondary_parent_cachegroup_id, ` + strconv.Itoa(atscfg.InvalidID) + `) as secondary_parent_cachegroup_id,
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
`
}

// GetATSMajorVersion returns the major version of the given profile's package trafficserver parameter.
// If no parameter exists, this does not return an error, but rather logs a warning and uses DefaultATSVersion.
func GetATSMajorVersion(tx *sql.Tx, serverProfileID atscfg.ProfileID) (int, error) {
	atsVersion, _, err := GetProfileParamValue(tx, int(serverProfileID), "package", "trafficserver")
	if err != nil {
		return 0, errors.New("getting profile param value: " + err.Error())
	}
	if len(atsVersion) == 0 {
		atsVersion = atscfg.DefaultATSVersion
		log.Warnln("Parameter package.trafficserver missing for profile " + strconv.Itoa(int(serverProfileID)) + ". Assuming version " + atsVersion)
	}
	atsMajorVer, err := atscfg.GetATSMajorVersionFromATSVersion(atsVersion)
	if err != nil {
		return 0, errors.New("ats version parameter '" + atsVersion + "' on this profile is not a number (config_file 'package', name 'trafficserver')")
	}
	return atsMajorVer, nil
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

// GetLocationParams returns a map[configFile]locationParams, and any error. If either param doesn't exist, an empty string is returned without error.
func GetLocationParams(tx *sql.Tx, profileID int) (map[string]atscfg.ConfigProfileParams, error) {
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

	params := map[string]atscfg.ConfigProfileParams{}
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

type TMParams struct {
	URL             string
	ReverseProxyURL string
}

// GetScopeParameters returns a map[cfgFile]scope, from all Parameters with the name 'scope' (irrespective of Profile).
func GetScopeParameters(tx *sql.Tx) (map[string]string, error) {
	rows, err := tx.Query(`SELECT config_file, value FROM parameter p WHERE name = 'scope'`)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	scopes := map[string]string{}
	for rows.Next() {
		cfgFile := ""
		val := ""
		if err := rows.Scan(&cfgFile, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		scopes[cfgFile] = val
	}
	return scopes, nil
}

// GetServerNameAndTypeFromID returns the server's name, type, whether it exists, and any error.
func GetServerNameAndTypeFromID(tx *sql.Tx, id int) (tc.CacheName, tc.CacheType, bool, error) {
	qry := `
SELECT
  s.host_name,
  tp.name
FROM
  server s
  JOIN type tp on s.type = tp.id
WHERE
  s.id = $1
`
	name := tc.CacheName("")
	typ := tc.CacheType("")
	if err := tx.QueryRow(qry, id).Scan(&name, &typ); err != nil {
		if err == sql.ErrNoRows {
			return "", tc.CacheType(""), false, nil
		}
		return "", tc.CacheType(""), false, errors.New("querying: " + err.Error())
	}
	return name, typ, true, nil
}

// GetServerNameAndTypeFromID returns the server's name, type, whether it exists, and any error.
func GetServerTypeFromName(tx *sql.Tx, name tc.CacheName) (tc.CacheType, bool, error) {
	qry := `
SELECT
  tp.name
FROM
  server s
  JOIN type tp on s.type = tp.id
WHERE
  s.host_name = $1
`
	typ := tc.CacheType("")
	if err := tx.QueryRow(qry, name).Scan(&typ); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying: " + err.Error())
	}
	return typ, true, nil
}

// GetServerNameAndDomainFromID returns the server's name, domain, whether it exists, and any error.
func GetServerNameAndDomainFromID(tx *sql.Tx, id int) (tc.CacheName, string, bool, error) {
	qry := `
SELECT
  s.host_name,
  s.domain_name
FROM
  server s
WHERE
  s.id = $1
`
	name := tc.CacheName("")
	domain := ""
	if err := tx.QueryRow(qry, id).Scan(&name, &domain); err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		return "", "", false, errors.New("querying: " + err.Error())
	}
	return name, domain, true, nil
}

// GetServerNameAndDomainFromID returns the server's name, domain, whether it exists, and any error.
func GetServerDomainFromName(tx *sql.Tx, name tc.CacheName) (string, bool, error) {
	qry := `
SELECT
  s.domain_name
FROM
  server s
WHERE
  s.host_name = $1
`
	domain := ""
	if err := tx.QueryRow(qry, name).Scan(&domain); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying: " + err.Error())
	}
	return domain, true, nil
}

func GetToolNameAndURL(tx *sql.Tx) (string, string, error) {
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
		return "", "", errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	toolName := ""
	url := ""
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return "", "", errors.New("scanning: " + err.Error())
		}
		if name == "tm.toolname" {
			toolName = val
		} else if name == "tm.url" {
			url = val
		}
	}
	return toolName, url, nil
}

func GetProfileParamsByName(tx *sql.Tx, profileName string, configFile string) (map[string][]string, error) {
	qry := `
SELECT
  p.name,
  p.value
FROM
  parameter p
  join profile_parameter pp on p.id = pp.parameter
  JOIN profile pr on pr.id = pp.profile
WHERE
  pr.name = $1
  AND p.config_file = $2
`
	rows, err := tx.Query(qry, profileName, configFile)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[string][]string{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		params[name] = append(params[name], val)
	}
	return params, nil
}

// getCDNNameFromNameOrID returns the CDN name from a parameter which may be the name or ID.
// This also checks and verifies the existence of the given CDN, and returns an appropriate user error if it doesn't exist.
// Returns the name, any user error, any system error, and any error code.
func GetCDNNameFromNameOrID(tx *sql.Tx, cdnNameOrID string) (string, error, error, int) {
	if cdnID, err := strconv.Atoi(cdnNameOrID); err == nil {
		cdnName, ok, err := dbhelpers.GetCDNNameFromID(tx, int64(cdnID))
		if err != nil {
			return "", nil, fmt.Errorf("getting CDN name from id %v: %v", cdnID, err), http.StatusInternalServerError
		} else if !ok {
			return "", errors.New("cdn not found"), nil, http.StatusNotFound
		}
		return string(cdnName), nil, nil, http.StatusOK
	}

	cdnName := cdnNameOrID
	if ok, err := dbhelpers.CDNExists(cdnName, tx); err != nil {
		return "", nil, fmt.Errorf("checking CDN name '%v' existence: %v", cdnName, err), http.StatusInternalServerError
	} else if !ok {
		return "", errors.New("cdn not found"), nil, http.StatusNotFound
	}
	return cdnName, nil, nil, http.StatusOK
}

// GetServerCapabilities returns the list of capabilities assigned to the given servers.
func GetServerCapabilitiesByID(tx *sql.Tx, serverIDs []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	qry := `
SELECT
  sc.server,
  sc.server_capability
FROM
  server_server_capability sc
WHERE
  sc.server = ANY($1)
`
	rows, err := tx.Query(qry, pq.Array(serverIDs))
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	serverCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	for rows.Next() {
		id := 0
		cap := atscfg.ServerCapability("")
		if err := rows.Scan(&id, &cap); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if _, ok := serverCaps[id]; !ok {
			serverCaps[id] = map[atscfg.ServerCapability]struct{}{}
		}
		serverCaps[id][cap] = struct{}{}
	}
	return serverCaps, nil
}

// GetDeliveryServiceRequiredCapabilities returns the list of required capabilities assigned to the given delivery services.
func GetDeliveryServiceRequiredCapabilities(tx *sql.Tx, dses []int) (map[int]map[atscfg.ServerCapability]struct{}, error) {
	qry := `
SELECT
  dsc.deliveryservice_id,
  dsc.required_capability
FROM
  deliveryservices_required_capability dsc
WHERE
  dsc.deliveryservice_id = ANY($1)
`
	rows, err := tx.Query(qry, pq.Array(dses))
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dsCaps := map[int]map[atscfg.ServerCapability]struct{}{}
	for rows.Next() {
		id := 0
		cap := atscfg.ServerCapability("")
		if err := rows.Scan(&id, &cap); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if _, ok := dsCaps[id]; !ok {
			dsCaps[id] = map[atscfg.ServerCapability]struct{}{}
		}
		dsCaps[id][cap] = struct{}{}
	}
	return dsCaps, nil
}
