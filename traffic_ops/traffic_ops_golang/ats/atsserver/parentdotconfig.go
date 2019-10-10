package atsserver

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
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"

	"github.com/lib/pq"
)

func GetParentDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"server-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	idOrHost := strings.TrimSuffix(inf.Params["server-name-or-id"], ".json")
	hostName := ""
	isHost := false
	id, err := strconv.Atoi(idOrHost)
	if err != nil {
		isHost = true
		hostName = idOrHost
	}

	serverInfo, ok, err := &atscfg.ServerInfo{}, false, error(nil)
	if isHost {
		serverInfo, ok, err = getServerInfoByHost(inf.Tx.Tx, hostName)
	} else {
		serverInfo, ok, err = getServerInfoByID(inf.Tx.Tx, id)
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server info: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server not found"), nil)
		return
	}

	atsMajorVer, err := GetATSMajorVersion(inf.Tx.Tx, int(serverInfo.ProfileID))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ATS major version: "+err.Error()))
		return
	}

	toolName, toURL, err := ats.GetToolNameAndURL(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting toolname and TO url parameters: "+err.Error()))
		return
	}

	parentConfigDSes := []atscfg.ParentConfigDSTopLevel{}
	if serverInfo.IsTopLevelCache() {
		parentConfigDSes, err = getParentConfigDSTopLevel(inf.Tx.Tx, serverInfo.CDN)
	} else {
		parentConfigDSes, err = getParentConfigDS(inf.Tx.Tx, serverInfo.ID)
	}
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server params: "+err.Error()))
		return
	}

	serverParams, err := getParentConfigServerProfileParams(inf.Tx.Tx, serverInfo.ID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server params: "+err.Error()))
		return
	}

	parentInfos, err := getParentInfo(inf.Tx.Tx, serverInfo)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server parent info: "+err.Error()))
		return
	}

	text := atscfg.MakeParentDotConfig(serverInfo, atsMajorVer, toolName, toURL, parentConfigDSes, serverParams, parentInfos)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}

type ParentConfigDSSortByName []atscfg.ParentConfigDS

func (s ParentConfigDSSortByName) Len() int      { return len(([]atscfg.ParentConfigDS)(s)) }
func (s ParentConfigDSSortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ParentConfigDSSortByName) Less(i, j int) bool {
	// TODO make this match the Perl sort "foreach my $ds ( sort @{ $data->{dslist} } )" ?
	return strings.Compare(string(s[i].Name), string(s[j].Name)) < 0
}

// getServerInfo returns the necessary info about the server, whether the server exists, and any error.
func getServerInfoByID(tx *sql.Tx, id int) (*atscfg.ServerInfo, bool, error) {
	return getServerInfo(tx, ServerInfoQuery()+`WHERE s.id = $1`, []interface{}{id})
}

// getServerInfo returns the necessary info about the server, whether the server exists, and any error.
func getServerInfoByHost(tx *sql.Tx, host string) (*atscfg.ServerInfo, bool, error) {
	return getServerInfo(tx, ServerInfoQuery()+` WHERE s.host_name = $1 `, []interface{}{host})
}

// getServerInfo returns the necessary info about the server, whether the server exists, and any error.
func getServerInfo(tx *sql.Tx, qry string, qryParams []interface{}) (*atscfg.ServerInfo, bool, error) {
	s := atscfg.ServerInfo{}
	if err := tx.QueryRow(qry, qryParams...).Scan(&s.CDN, &s.CDNID, &s.ID, &s.HostName, &s.DomainName, &s.IP, &s.ProfileID, &s.ProfileName, &s.Port, &s.Type, &s.CacheGroupID, &s.ParentCacheGroupID, &s.SecondaryParentCacheGroupID, &s.ParentCacheGroupType, &s.SecondaryParentCacheGroupType); err != nil {
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
func GetATSMajorVersion(tx *sql.Tx, serverProfileID int) (int, error) {
	atsVersion, _, err := ats.GetProfileParamValue(tx, serverProfileID, "package", "trafficserver")
	if err != nil {
		return 0, errors.New("getting profile param value: " + err.Error())
	}
	if len(atsVersion) == 0 {
		atsVersion = atscfg.DefaultATSVersion
		log.Warnln("Parameter package.trafficserver missing for profile " + strconv.Itoa(int(serverProfileID)) + ". Assuming version " + atsVersion)
	}
	return atscfg.GetATSMajorVersionFromATSVersion(atsVersion)
}

type ParentConfigDS struct {
	Name            tc.DeliveryServiceName
	QStringIgnore   tc.QStringIgnore
	OriginFQDN      string
	MultiSiteOrigin bool
	OriginShield    string
	Type            tc.DSType

	QStringHandling string
}

type ParentConfigDSTopLevel struct {
	ParentConfigDS
	MSOAlgorithm                       string
	MSOParentRetry                     string
	MSOUnavailableServerRetryResponses string
	MSOMaxSimpleRetries                string
	MSOMaxUnavailableServerRetries     string
}

func ParentConfigDSQuerySelect() string {
	return `
SELECT
  ds.xml_id,
  COALESCE(ds.qstring_ignore, ` + tc.QStringIgnoreUseInCacheKeyAndPassUp.String() + `),
  COALESCE((SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
    FROM origin o
    WHERE o.deliveryservice = ds.id
    AND o.is_primary), '') as org_server_fqdn,
  COALESCE(ds.multi_site_origin, false),
  COALESCE(ds.origin_shield, ''),
  ARRAY(SELECT required_capability FROM deliveryservices_required_capability dsrc WHERE dsrc.deliveryservice_id = ds.id),
  dt.name AS ds_type
`
}

const ParentConfigDSQueryFromTopLevel = `
FROM
  deliveryservice ds
  JOIN type as dt ON ds.type = dt.id
  JOIN cdn ON cdn.id = ds.cdn_id
` // TODO Perl does 'JOIN deliveryservice_regex dsr ON dsr.deliveryservice = ds.id   JOIN regex r ON dsr.regex = r.id   JOIN type as rt ON r.type = rt.id' and orders by, but doesn't use; ensure it isn't necessary

const ParentConfigDSQueryFrom = ParentConfigDSQueryFromTopLevel + `
`

const ParentConfigDSQueryOrder = `
ORDER BY ds.id
` // TODO: perl does 'ORDER BY ds.id, rt.name, dsr.set_number' - but doesn't actually use regexes - ensure it isn't necessary

const ParentConfigDSQueryWhere = `
WHERE ds.id in (SELECT DISTINCT(dss.deliveryservice) FROM deliveryservice_server dss where dss.server = $1)
`

const ParentConfigDSQueryWhereTopLevel = `
WHERE
  cdn.name = $1
  AND ds.id in (SELECT deliveryservice_server.deliveryservice FROM deliveryservice_server)
  AND ds.active = true
`

func ParentConfigDSQuery() string {
	return ParentConfigDSQuerySelect() +
		ParentConfigDSQueryFrom +
		ParentConfigDSQueryWhere +
		ParentConfigDSQueryOrder
}

func ParentConfigDSQueryTopLevel() string {
	return ParentConfigDSQuerySelect() +
		ParentConfigDSQueryFromTopLevel +
		ParentConfigDSQueryWhereTopLevel +
		ParentConfigDSQueryOrder
}

func getParentConfigDSTopLevel(tx *sql.Tx, cdnName tc.CDNName) ([]atscfg.ParentConfigDSTopLevel, error) {
	dses, err := getParentConfigDSRaw(tx, ParentConfigDSQueryTopLevel(), []interface{}{cdnName})
	if err != nil {
		return nil, errors.New("getting top level raw parent config ds: " + err.Error())
	}
	topDSes := []atscfg.ParentConfigDSTopLevel{}
	for _, ds := range dses {
		topDSes = append(topDSes, ds)
	}

	dsesWithParams, err := getParentConfigDSParamsTopLevel(tx, topDSes)
	if err != nil {
		return nil, errors.New("getting top level ds params: " + err.Error())
	}

	return dsesWithParams, nil
}

func getParentConfigDS(tx *sql.Tx, serverID int) ([]atscfg.ParentConfigDSTopLevel, error) {
	dses, err := getParentConfigDSRaw(tx, ParentConfigDSQuery(), []interface{}{serverID})
	if err != nil {
		return nil, errors.New("getting raw parent config ds: " + err.Error())
	}

	dsesWithParams, err := getParentConfigDSParams(tx, dses)
	if err != nil {
		return nil, errors.New("getting ds params: " + err.Error())
	}
	return dsesWithParams, nil
}

func getParentConfigServerProfileParams(tx *sql.Tx, serverID int) (map[string]string, error) {
	qry := `
SELECT
  pa.name,
  pa.value
FROM
  parameter pa
  JOIN profile_parameter pp ON pp.parameter = pa.id
  JOIN profile pr ON pr.id = pp.profile
  JOIN server s on s.profile = pr.id
WHERE
  s.id = $1
  AND pa.config_file = 'parent.config'
  AND pa.name IN (
    '` + atscfg.ParentConfigParamQStringHandling + `',
    '` + atscfg.ParentConfigParamAlgorithm + `',
    '` + atscfg.ParentConfigParamQString + `'
  )
`
	rows, err := tx.Query(qry, serverID)
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
		params[name] = val
	}
	return params, nil
}

// getParentConfigDSRaw returns a ParentConfigDSTopLevel, but all fields in addition to ParentConfigDS will be defaulted. This is because a ParentConfigDSTopLevel is returned to share the same interface, but it doesn't actually have top level data.
func getParentConfigDSRaw(tx *sql.Tx, qry string, qryParams []interface{}) ([]atscfg.ParentConfigDSTopLevel, error) {
	rows, err := tx.Query(qry, qryParams...)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()
	dses := []atscfg.ParentConfigDSTopLevel{}
	for rows.Next() {
		d := atscfg.ParentConfigDS{RequiredCapabilities: map[atscfg.ServerCapability]struct{}{}}
		requiredCaps := []string{}
		if err := rows.Scan(&d.Name, &d.QStringIgnore, &d.OriginFQDN, &d.MultiSiteOrigin, &d.OriginShield, pq.Array(&requiredCaps), &d.Type); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if d.OriginFQDN == "" {
			// TODO skip ANY_MAP DSes? Why? Did Perl, I didn't see it?
			log.Errorf("parent.config generation: getting parent config ds: server %+v has no origin, skipping!\n", d.Name)
			continue
		}
		d.Type = tc.DSTypeFromString(string(d.Type))
		for _, cap := range requiredCaps {

			d.RequiredCapabilities[atscfg.ServerCapability(cap)] = struct{}{}
		}
		dses = append(dses, atscfg.ParentConfigDSTopLevel{ParentConfigDS: d})
	}

	return dses, nil
}

func parentConfigDSesToNames(dses []atscfg.ParentConfigDS) []string {
	names := []string{}
	for _, ds := range dses {
		names = append(names, string(ds.Name))
	}
	return names
}

func parentConfigDSesToNamesTopLevel(dses []atscfg.ParentConfigDSTopLevel) []string {
	names := []string{}
	for _, ds := range dses {
		names = append(names, string(ds.Name))
	}
	return names
}

const ParentConfigDSParamsQuerySelect = `
SELECT
  ds.xml_id,
  pa.name,
  pa.value
`
const ParentConfigDSParamsQueryFrom = `
FROM
  parameter pa
  JOIN profile_parameter pp ON pp.parameter = pa.id
  JOIN profile pr ON pr.id = pp.profile
  JOIN deliveryservice ds on ds.profile = pr.id
`
const ParentConfigDSParamsQueryWhere = `
WHERE
  pa.config_file = 'parent.config'
  AND ds.xml_id = ANY($1)
  AND pa.name IN (
    '` + atscfg.ParentConfigParamQStringHandling + `'
  )
`

var ParentConfigDSParamsQueryWhereTopLevel = `
WHERE
  pa.config_file = 'parent.config'
  AND ds.xml_id = ANY($1)
  AND pa.name IN (
    '` + atscfg.ParentConfigParamQStringHandling + `',
    '` + atscfg.ParentConfigParamMSOAlgorithm + `',
    '` + atscfg.ParentConfigParamMSOParentRetry + `',
    '` + atscfg.ParentConfigParamUnavailableServerRetryResponses + `',
    '` + atscfg.ParentConfigParamMaxSimpleRetries + `',
    '` + atscfg.ParentConfigParamMaxUnavailableServerRetries + `'
  )
`

const ParentConfigDSParamsQuery = ParentConfigDSParamsQuerySelect + ParentConfigDSParamsQueryFrom + ParentConfigDSParamsQueryWhere

var ParentConfigDSParamsQueryTopLevel = ParentConfigDSParamsQuerySelect + ParentConfigDSParamsQueryFrom + ParentConfigDSParamsQueryWhereTopLevel

func getParentConfigDSParams(tx *sql.Tx, dses []atscfg.ParentConfigDSTopLevel) ([]atscfg.ParentConfigDSTopLevel, error) {
	params, err := getParentConfigDSParamsRaw(tx, ParentConfigDSParamsQuery, parentConfigDSesToNamesTopLevel(dses))
	if err != nil {
		return nil, err
	}
	for i, ds := range dses {
		dsParams, ok := params[ds.Name]
		if !ok {
			continue
		}
		if v, ok := dsParams[atscfg.ParentConfigParamQStringHandling]; ok {
			ds.QStringHandling = v
			dses[i] = ds
		}
	}
	return dses, nil
}

func getParentConfigDSParamsTopLevel(tx *sql.Tx, dses []atscfg.ParentConfigDSTopLevel) ([]atscfg.ParentConfigDSTopLevel, error) {
	params, err := getParentConfigDSParamsRaw(tx, ParentConfigDSParamsQueryTopLevel, parentConfigDSesToNamesTopLevel(dses))
	if err != nil {
		return nil, err
	}
	for i, ds := range dses {
		dsParams := params[ds.Name] // it's acceptable for this to not exist, if there are no params for the DS. If so, we still need to continue below, to set all the defaults.
		if v, ok := dsParams[atscfg.ParentConfigParamQStringHandling]; ok {
			ds.QStringHandling = v
		}
		if v, ok := dsParams[atscfg.ParentConfigParamMSOAlgorithm]; ok && strings.TrimSpace(v) != "" {
			ds.MSOAlgorithm = v
		} else {
			ds.MSOAlgorithm = atscfg.ParentConfigDSParamDefaultMSOAlgorithm
		}
		if v, ok := dsParams[atscfg.ParentConfigParamMSOParentRetry]; ok {
			ds.MSOParentRetry = v
		} else {
			ds.MSOParentRetry = atscfg.ParentConfigDSParamDefaultMSOParentRetry
		}
		if v, ok := dsParams[atscfg.ParentConfigParamUnavailableServerRetryResponses]; ok {
			ds.MSOUnavailableServerRetryResponses = v
		} else {
			ds.MSOUnavailableServerRetryResponses = atscfg.ParentConfigDSParamDefaultMSOUnavailableServerRetryResponses
		}
		if v, ok := dsParams[atscfg.ParentConfigParamMaxSimpleRetries]; ok {
			ds.MSOMaxSimpleRetries = v
		} else {
			ds.MSOMaxSimpleRetries = atscfg.ParentConfigDSParamDefaultMaxSimpleRetries
		}
		if v, ok := dsParams[atscfg.ParentConfigParamMaxUnavailableServerRetries]; ok {
			ds.MSOMaxUnavailableServerRetries = v
		} else {
			ds.MSOMaxUnavailableServerRetries = atscfg.ParentConfigDSParamDefaultMaxUnavailableServerRetries
		}
		dses[i] = ds
	}
	return dses, nil
}

func getParentConfigDSParamsRaw(tx *sql.Tx, qry string, dsNames []string) (map[tc.DeliveryServiceName]map[string]string, error) {
	rows, err := tx.Query(qry, pq.Array(dsNames))
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[tc.DeliveryServiceName]map[string]string{}
	for rows.Next() {
		dsName := tc.DeliveryServiceName("")
		pName := ""
		pVal := ""
		if err := rows.Scan(&dsName, &pName, &pVal); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if _, ok := params[dsName]; !ok {
			params[dsName] = map[string]string{}
		}
		params[dsName][pName] = pVal
	}
	return params, nil
}

func getParentInfo(tx *sql.Tx, server *atscfg.ServerInfo) (map[atscfg.OriginHost][]atscfg.ParentInfo, error) {
	parentInfos := map[atscfg.OriginHost][]atscfg.ParentInfo{}

	serverDomain, ok, err := getCDNDomainByProfileID(tx, server.ProfileID)
	if err != nil {
		return nil, errors.New("getting CDN domain from profile ID: " + err.Error())
	} else if !ok || serverDomain == "" {
		return parentInfos, nil // TODO warn? Perl doesn't.
	}

	profileCaches, originServers, err := getServerParentCacheGroupProfiles(tx, server)
	if err != nil {
		return nil, errors.New("getting server parent cachegroup profiles: " + err.Error())
	}

	return atscfg.MakeParentInfo(server, serverDomain, profileCaches, originServers), nil
}

// getServerParentCacheGroupProfiles gets the profile information for servers belonging to the parent cachegroup, and secondary parent cachegroup, of the cachegroup of each server.
func getServerParentCacheGroupProfiles(tx *sql.Tx, server *atscfg.ServerInfo) (map[atscfg.ProfileID]atscfg.ProfileCache, map[atscfg.OriginHost][]atscfg.CGServer, error) {
	// TODO make this more efficient - should be a single query - this was transliterated from Perl - it's extremely inefficient.

	profileCaches := map[atscfg.ProfileID]atscfg.ProfileCache{}
	originServers := map[atscfg.OriginHost][]atscfg.CGServer{} // "deliveryServices" in Perl

	qry := ""
	if server.IsTopLevelCache() {
		// multisite origins take all the org groups in to account
		qry = `
WITH parent_cachegroup_ids AS (
  SELECT cg.id as v
  FROM cachegroup cg
  JOIN type on type.id = cg.type
  WHERE type.name = '` + tc.CacheGroupOriginTypeName + `'
)
`
	} else {
		qry = `
WITH server_cachegroup_ids AS (
  SELECT cachegroup as v FROM server WHERE id = $2
),
parent_cachegroup_ids AS (
  SELECT parent_cachegroup_id as v
  FROM cachegroup WHERE id IN (SELECT v from server_cachegroup_ids)
  UNION ALL
  SELECT secondary_parent_cachegroup_id as v
  FROM cachegroup WHERE id IN (SELECT v from server_cachegroup_ids)
)
`
	}

	qry += `
SELECT
  s.id,
  s.host_name,
  s.ip_address,
  s.tcp_port,
  s.cachegroup,
  s.status,
  s.type,
  s.profile,
  s.cdn_id,
  stype.name as type_name,
  ARRAY(SELECT server_capability FROM server_server_capability ssc WHERE ssc.server = s.id),
  s.domain_name
FROM
  server s
  JOIN type stype ON s.type = stype.id
  JOIN cachegroup cg ON cg.id = s.cachegroup
  JOIN cdn on s.cdn_id = cdn.id
  JOIN status st ON st.id = s.status
WHERE
  cg.id IN (SELECT v FROM parent_cachegroup_ids)
  AND (stype.name = '` + tc.OriginTypeName + `' OR stype.name LIKE '` + tc.EdgeTypePrefix + `%' OR stype.name LIKE '` + tc.MidTypePrefix + `%')
  AND (st.name = '` + string(tc.CacheStatusReported) + `' OR st.name = '` + string(tc.CacheStatusOnline) + `')
  AND cdn.name = $1
`

	// TODO move qry, qryParams to separate funcs/consts
	qryParams := []interface{}{}
	if server.IsTopLevelCache() {
		qryParams = []interface{}{server.CDN}
	} else {
		qryParams = []interface{}{server.CDN, server.ID}
	}

	rows, err := tx.Query(qry, qryParams...)
	if err != nil {
		return nil, nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	cgServerIDs := []int{}
	cgServers := []atscfg.CGServer{}
	for rows.Next() {
		s := atscfg.CGServer{Capabilities: map[atscfg.ServerCapability]struct{}{}}
		caps := []string{}
		if err := rows.Scan(&s.ServerID, &s.ServerHost, &s.ServerIP, &s.ServerPort, &s.CacheGroupID, &s.Status, &s.Type, &s.ProfileID, &s.CDN, &s.TypeName, pq.Array(&caps), &s.Domain); err != nil {
			return nil, nil, errors.New("scanning: " + err.Error())
		}
		for _, cap := range caps {
			s.Capabilities[atscfg.ServerCapability(cap)] = struct{}{}
		}
		cgServers = append(cgServers, s)
		cgServerIDs = append(cgServerIDs, int(s.ServerID))
	}

	serverCapabilities, err := ats.GetServerCapabilitiesByID(tx, cgServerIDs)
	if err != nil {
		return nil, nil, errors.New("getting server capabilities: " + err.Error())
	}

	cgServerDSes, err := getServerDSes(tx, cgServerIDs)
	if err != nil {
		return nil, nil, errors.New("getting cachegroup server deliveryservices: " + err.Error())
	}

	profileParams, err := getParentConfigServerCacheProfileParams(tx, cgServerIDs) // TODO change to take cg IDs directly?
	if err != nil {
		return nil, nil, errors.New("getting cachegroup server profile params: " + err.Error())
	}

	allDSMap := map[atscfg.DeliveryServiceID]struct{}{}
	for _, dses := range cgServerDSes {
		for _, ds := range dses {
			allDSMap[ds] = struct{}{}
		}
	}
	allDSes := []int{}
	for ds, _ := range allDSMap {
		allDSes = append(allDSes, int(ds))
	}

	dsRequiredCapabilities, err := ats.GetDeliveryServiceRequiredCapabilities(tx, allDSes)
	if err != nil {
		return nil, nil, errors.New("getting DS required capabilities: " + err.Error())
	}

	dsOrigins, err := getDSOrigins(tx, allDSes)
	if err != nil {
		return nil, nil, errors.New("getting deliveryservice origins: " + err.Error())
	}

	for _, cgServer := range cgServers {
		if cgServer.TypeName == tc.OriginTypeName {
			dses := cgServerDSes[cgServer.ServerID]
			for _, ds := range dses {
				orgURI := dsOrigins[ds]
				if atscfg.HasRequiredCapabilities(serverCapabilities[int(cgServer.ServerID)], dsRequiredCapabilities[int(ds)]) {
					originServers[atscfg.OriginHost(orgURI.Host)] = append(originServers[atscfg.OriginHost(orgURI.Host)], cgServer)
				}
			}
		} else {
			originServers[atscfg.DeliveryServicesAllParentsKey] = append(originServers[atscfg.DeliveryServicesAllParentsKey], cgServer)
		}

		if _, profileCachesHasProfile := profileCaches[cgServer.ProfileID]; !profileCachesHasProfile {
			defaultProfileCache := atscfg.DefaultProfileCache()
			if profileCache, profileParamsHasProfile := profileParams[cgServer.ProfileID]; !profileParamsHasProfile {
				log.Warnf("cachegroup has server with profile %+v but that profile has no parameters", cgServer.ProfileID)
				profileCaches[cgServer.ProfileID] = defaultProfileCache
			} else {
				profileCaches[cgServer.ProfileID] = profileCache
			}
		}
	}
	return profileCaches, originServers, nil
}

func getServerDSes(tx *sql.Tx, serverIDs []int) (map[atscfg.ServerID][]atscfg.DeliveryServiceID, error) {
	sds := map[atscfg.ServerID][]atscfg.DeliveryServiceID{}
	if len(serverIDs) == 0 {
		return sds, nil
	}
	qry := `
SELECT
  dss.server,
  dss.deliveryservice
FROM
  deliveryservice_server dss
WHERE
  dss.server = ANY($1)
`
	rows, err := tx.Query(qry, pq.Array(serverIDs))
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		sID := atscfg.ServerID(0)
		dsID := atscfg.DeliveryServiceID(0)
		if err := rows.Scan(&sID, &dsID); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		sds[sID] = append(sds[sID], dsID)
	}
	return sds, nil
}

func getDSOrigins(tx *sql.Tx, dsIDs []int) (map[atscfg.DeliveryServiceID]*atscfg.OriginURI, error) {
	origins := map[atscfg.DeliveryServiceID]*atscfg.OriginURI{}
	if len(dsIDs) == 0 {
		return origins, nil
	}
	qry := `
SELECT
  ds.id,
  o.protocol::text,
  o.fqdn,
  COALESCE(o.port::text, '')
FROM
  deliveryservice ds
  JOIN origin o ON o.deliveryservice = ds.id
WHERE
  ds.id = ANY($1)
  AND o.is_primary
`
	rows, err := tx.Query(qry, pq.Array(dsIDs))
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		id := atscfg.DeliveryServiceID(0)
		uri := &atscfg.OriginURI{}
		if err := rows.Scan(&id, &uri.Scheme, &uri.Host, &uri.Port); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if uri.Port == "" {
			if uri.Scheme == "http" {
				uri.Port = "80"
			} else if uri.Scheme == "https" {
				uri.Port = "443"
			} else {
				log.Warnf("parent.config generation: origin had unknown scheme '" + uri.Scheme + "' and no port; leaving port empty")
			}
		}
		origins[id] = uri
	}
	return origins, nil
}

func getParentConfigServerCacheProfileParams(tx *sql.Tx, serverIDs []int) (map[atscfg.ProfileID]atscfg.ProfileCache, error) {
	qry := `
SELECT
  pr.id,
  pa.name,
  pa.value
FROM
  parameter pa
  JOIN profile_parameter pp ON pp.parameter = pa.id
  JOIN profile pr ON pr.id = pp.profile
  JOIN server s on s.profile = pr.id
WHERE
  s.id = ANY($1)
  AND pa.config_file = 'parent.config'
  AND pa.name IN (
    '` + atscfg.ParentConfigCacheParamWeight + `',
    '` + atscfg.ParentConfigCacheParamPort + `',
    '` + atscfg.ParentConfigCacheParamUseIP + `',
    '` + atscfg.ParentConfigCacheParamRank + `',
    '` + atscfg.ParentConfigCacheParamNotAParent + `'
  )
`
	rows, err := tx.Query(qry, pq.Array(serverIDs))
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	type Param struct {
		ProfileID atscfg.ProfileID
		Name      string
		Val       string
	}

	params := []Param{}
	for rows.Next() {
		p := Param{}
		if err := rows.Scan(&p.ProfileID, &p.Name, &p.Val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		params = append(params, p)
	}

	sParams := map[atscfg.ProfileID]atscfg.ProfileCache{} // TODO change to map of pointers? Does efficiency matter?
	for _, param := range params {
		profileCache, ok := sParams[param.ProfileID]
		if !ok {
			profileCache = atscfg.DefaultProfileCache()
		}
		switch param.Name {
		case atscfg.ParentConfigCacheParamWeight:
			// f, err := strconv.ParseFloat(param.Val, 64)
			// if err != nil {
			// 	log.Errorln("parent.config generation: weight param is not a float, skipping! : " + err.Error())
			// } else {
			// 	profileCache.Weight = f
			// }
			// TODO validate float?
			profileCache.Weight = param.Val
		case atscfg.ParentConfigCacheParamPort:
			i, err := strconv.ParseInt(param.Val, 10, 64)
			if err != nil {
				log.Errorln("parent.config generation: port param is not an integer, skipping! : " + err.Error())
			} else {
				profileCache.Port = int(i)
			}
		case atscfg.ParentConfigCacheParamUseIP:
			profileCache.UseIP = param.Val == "1"
		case atscfg.ParentConfigCacheParamRank:
			i, err := strconv.ParseInt(param.Val, 10, 64)
			if err != nil {
				log.Errorln("parent.config generation: rank param is not an integer, skipping! : " + err.Error())
			} else {
				profileCache.Rank = int(i)
			}

		case atscfg.ParentConfigCacheParamNotAParent:
			profileCache.NotAParent = param.Val != "false"
		default:
			return nil, errors.New("query returned unexpected param: " + param.Name)
		}
		sParams[param.ProfileID] = profileCache
	}
	return sParams, nil
}

func getServerParams(tx *sql.Tx, serverID int) (map[string]string, error) {
	qry := `
SELECT
  pa.name
  pa.value
FROM
  parameter pa
  JOIN profile_parameter pp ON pp.parameter = pa.id
  JOIN profile pr ON pr.id = pp.profile
  JOIN server s on s.profile = pr.id
WHERE
  s.id = $1
  AND pa.config_file = 'parent.config'
  AND pa.name IN (
    '` + atscfg.ParentConfigParamQStringHandling + `',
    '` + atscfg.ParentConfigParamAlgorithm + `',
    '` + atscfg.ParentConfigParamQString + `'
  )
`
	rows, err := tx.Query(qry, serverID)
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
		params[name] = val
	}
	return params, nil
}

type ParentConfigServerParams struct {
	QString         string
	Algorithm       string
	QStringHandling bool
}

func getCDNDomainByProfileID(tx *sql.Tx, profileID atscfg.ProfileID) (string, bool, error) {
	qry := `SELECT domain_name from cdn where id = (select cdn from profile where id = $1)`
	val := ""
	if err := tx.QueryRow(qry, profileID).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying: " + err.Error())
	}
	return val, true, nil
}
