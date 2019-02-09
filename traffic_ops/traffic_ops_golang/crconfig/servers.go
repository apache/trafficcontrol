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

const RouterTypeName = "CCR"
const MonitorTypeName = "RASCAL"
const EdgeTypePrefix = "EDGE"
const MidTypePrefix = "MID"

func makeCRConfigServers(tx *sql.Tx, cdn string, cdnDomain string, live bool) (
	map[string]tc.CRConfigTrafficOpsServer,
	map[string]tc.CRConfigRouter,
	map[string]tc.CRConfigMonitor,
	error,
) {
	allServers, err := getAllServers(tx, cdn, live)
	if err != nil {
		return nil, nil, nil, err
	}

	serverIDNames := allServersToServerIDNames(allServers)

	serverDSes, err := getServerDSes(tx, cdn, cdnDomain, serverIDNames, live)
	if err != nil {
		return nil, nil, nil, errors.New("getting server deliveryservices: " + err.Error())
	}

	servers := map[string]tc.CRConfigTrafficOpsServer{}
	routers := map[string]tc.CRConfigRouter{}
	monitors := map[string]tc.CRConfigMonitor{}
	for host, s := range allServers {
		switch {
		case *s.ServerType == tc.RouterTypeName:
			status := tc.CRConfigRouterStatus(*s.ServerStatus)
			routers[host] = tc.CRConfigRouter{
				APIPort:       s.APIPort,
				FQDN:          s.Fqdn,
				HTTPSPort:     s.HttpsPort,
				IP:            s.Ip,
				IP6:           s.Ip6,
				Location:      s.LocationId,
				Port:          s.Port,
				Profile:       s.Profile,
				SecureAPIPort: s.SecureAPIPort,
				ServerStatus:  &status,
			}
		case *s.ServerType == tc.MonitorTypeName:
			monitors[host] = tc.CRConfigMonitor{
				FQDN:         s.Fqdn,
				HTTPSPort:    s.HttpsPort,
				IP:           s.Ip,
				IP6:          s.Ip6,
				Location:     s.LocationId,
				Port:         s.Port,
				Profile:      s.Profile,
				ServerStatus: s.ServerStatus,
			}
		case strings.HasPrefix(*s.ServerType, tc.EdgeTypePrefix) || strings.HasPrefix(*s.ServerType, tc.MidTypePrefix):
			if s.RoutingDisabled == 0 {
				s.CRConfigTrafficOpsServer.DeliveryServices = serverDSes[host]
			}
			servers[host] = s.CRConfigTrafficOpsServer
		}
	}
	return servers, routers, monitors, nil
}

// ServerUnion has all fields from all servers. This is used to select all server data with a single query, and then convert each to the proper type afterwards.
type ServerUnion struct {
	tc.CRConfigTrafficOpsServer
	APIPort       *string
	SecureAPIPort *string
	ID            int
}

const DefaultWeightMultiplier = 1000.0
const DefaultWeight = 0.999

func getAllServers(tx *sql.Tx, cdn string, live bool) (map[string]ServerUnion, error) {
	servers := map[string]ServerUnion{}

	serverParams, err := getServerParams(tx, cdn, live)
	if err != nil {
		return nil, errors.New("getting server params: " + err.Error())
	}

	// TODO select deliveryservices as array?
	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
),
snapshot_time AS (
  SELECT time as v FROM snapshot sn where sn.cdn = (SELECT v from cdn_name)
)
SELECT
  host_name,
  cachegroup,
  fqdn,
  hashid,
  https_port,
  interface_name,
  ip_address,
  ip6_address,
  tcp_port,
  profile_name,
  routing_disabled,
  status,
  type,
  id
FROM (
SELECT DISTINCT ON (s.host_name)
  s.host_name,
  cg.name as cachegroup,
  concat(s.host_name, '.', s.domain_name) as fqdn,
  s.xmpp_id as hashid,
  s.https_port,
  s.interface_name,
  s.ip_address,
  s.ip6_address,
  s.tcp_port,
  p.name as profile_name,
  cast(p.routing_disabled as int),
  st.name as status,
  t.name as type,
  s.id,
  s.deleted
FROM
  server_snapshot s
  JOIN cachegroup_snapshot cg ON cg.id = s.cachegroup
  JOIN type_snapshot t ON t.id = s.type
  JOIN profile_snapshot p ON p.id = s.profile
  JOIN status_snapshot st ON st.id = s.status
WHERE
  s.cdn_id = (SELECT id FROM cdn_snapshot c where c.name = (select v from cdn_name) and c.last_updated <= (select v from snapshot_time))
  AND (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	if !live {
		qry += `
  AND s.last_updated <= (select v from snapshot_time)
  AND cg.last_updated <= (select v from snapshot_time)
  AND t.last_updated <= (select v from snapshot_time)
  AND p.last_updated <= (select v from snapshot_time)
  AND st.last_updated <= (select v from snapshot_time)
`
	}
	qry += `
ORDER BY
  s.host_name DESC,
  s.last_updated DESC,
  cg.last_updated DESC,
  t.last_updated DESC,
  p.last_updated DESC,
  st.last_updated DESC
) s WHERE s.deleted = false
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying servers: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		port := sql.NullInt64{}
		ip6 := sql.NullString{}
		hashId := sql.NullString{}
		httpsPort := sql.NullInt64{}

		s := ServerUnion{}

		host := ""
		status := ""
		if err := rows.Scan(&host, &s.CacheGroup, &s.Fqdn, &hashId, &httpsPort, &s.InterfaceName, &s.Ip, &ip6, &port, &s.Profile, &s.RoutingDisabled, &status, &s.ServerType, &s.ID); err != nil {
			return nil, errors.New("scanning server: " + err.Error())
		}

		s.LocationId = s.CacheGroup

		serverStatus := tc.CRConfigServerStatus(status)
		s.ServerStatus = &serverStatus
		if port.Valid {
			i := int(port.Int64)
			s.Port = &i
		}

		s.Ip6 = &ip6.String // Don't check valid, assign empty string if null

		if hashId.String != "" {
			s.HashId = &hashId.String
		} else {
			s.HashId = &host
		}

		if httpsPort.Valid {
			i := int(httpsPort.Int64)
			s.HttpsPort = &i
		}

		params, hasParams := serverParams[host]
		if hasParams && params.APIPort != nil {
			s.APIPort = params.APIPort
		}

		if hasParams && params.SecureAPIPort != nil {
			s.SecureAPIPort = params.SecureAPIPort
		}

		weightMultiplier := DefaultWeightMultiplier
		if hasParams && params.WeightMultiplier != nil {
			weightMultiplier = *params.WeightMultiplier
		}
		weight := DefaultWeight
		if hasParams && params.Weight != nil {
			weight = *params.Weight
		}
		hashCount := int(weight * weightMultiplier)
		s.HashCount = &hashCount

		servers[host] = s
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("iterating router param rows: " + err.Error())
	}

	return servers, nil
}

// getServerDSNames returns a map[serverID]dsID
func getServerDSNames(tx *sql.Tx, cdn string, live bool) (map[int][]int, error) {
	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
)
SELECT dss.server, dss.deliveryservice FROM (
SELECT DISTINCT ON (dss.server, dss.deliveryservice) dss.server, dss.deliveryservice, dss.deleted
FROM
  deliveryservice_server_snapshot dss
  JOIN server_snapshot s ON dss.server = s.id
  JOIN deliveryservice_snapshot ds ON ds.id = dss.deliveryservice
  JOIN type_snapshot dt ON dt.id = ds.type
  JOIN profile_snapshot p ON p.id = s.profile
  JOIN status_snapshot st ON st.id = s.status
  JOIN deliveryservice_snapshots dsn ON dsn.deliveryservice = ds.xml_id
  JOIN type_snapshot dt ON dt.id = ds.type
WHERE
  ds.cdn_id = (select id from cdn_snapshot c where c.name = (select v from cdn_name) and c.last_updated <= dsn.time)
  AND ds.active = true
  AND dt.name != '` + string(tc.DSTypeAnyMap) + `'
  AND p.routing_disabled = false
  AND (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
	AND dt.name <> '` + string(tc.DSTypeAnyMap) + `'
`
	if !live {
		qry += `
  AND dss.last_updated <= dsn.time
  AND s.last_updated <= dsn.time
  AND ds.last_updated <= dsn.time
  AND p.last_updated <= dsn.time
  AND st.last_updated <= dsn.time
  AND dt.last_updated <= dsn.time
`
	}
	qry += `
ORDER BY
  dss.server DESC,
  dss.deliveryservice DESC,
  ds.last_updated DESC,
  dt.last_updated DESC,
  dss.last_updated DESC,
  s.last_updated DESC,
  p.last_updated DESC,
  st.last_updated DESC
) dss WHERE dss.deleted = false
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying server deliveryservice names: " + err.Error())
	}
	defer rows.Close()

	serverDSes := map[int][]int{}
	for rows.Next() {
		ds := 0
		server := 0
		if err := rows.Scan(&server, &ds); err != nil {
			return nil, errors.New("scanning server deliveryservice names: " + err.Error())
		}
		serverDSes[server] = append(serverDSes[server], ds)
	}
	return serverDSes, nil
}

type DSRouteInfo struct {
	DSName string
	IsDNS  bool
	IsRaw  bool
	Remap  string
}

// getServerDSes returns a map[serverName][dsName][]regexPattern
func getServerDSes(tx *sql.Tx, cdn string, domain string, serverIDNames map[int]string, live bool) (map[string]map[string][]string, error) {
	serverDSNames, err := getServerDSNames(tx, cdn, live)
	if err != nil {
		return nil, errors.New("getting server deliveryservice names: " + err.Error())
	}

	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
)
SELECT ds_id, ds, ds_type, routing_name, pattern FROM (
SELECT DISTINCT ON (ds.xml_id, dt.name, ds.routing_name, r.pattern)
  ds.id as ds_id,
  ds.xml_id as ds,
  dt.name as ds_type,
  ds.routing_name,
  r.pattern as pattern,
  dsr.set_number,
  ds.deleted
FROM
  regex_snapshot as r
  JOIN type_snapshot rt on r.type = rt.id
  JOIN deliveryservice_regex_snapshot dsr on dsr.regex = r.id
  JOIN deliveryservice_snapshot ds on ds.id = dsr.deliveryservice
  JOIN type_snapshot dt on dt.id = ds.type
  JOIN deliveryservice_snapshots dsn ON dsn.deliveryservice = ds.xml_id
WHERE
  ds.cdn_id = (select id from cdn_snapshot c where c.name = (select v from cdn_name) and c.last_updated <= dsn.time)
  AND ds.active = true
  AND dt.name != '` + string(tc.DSTypeAnyMap) + `'
  AND rt.name = 'HOST_REGEXP'
	AND dt.name <> '` + string(tc.DSTypeAnyMap) + `'
`
	if !live { // TODO use template?
		qry += `
  AND r.last_updated <= dsn.time
  AND rt.last_updated <= dsn.time
  AND dsr.last_updated <= dsn.time
  AND ds.last_updated <= dsn.time
  AND dt.last_updated <= dsn.time
`
	}
	qry += `
ORDER BY
  ds.xml_id DESC,
  dt.name DESC,
  ds.routing_name DESC,
  r.pattern DESC,
  r.last_updated DESC,
  rt.last_updated DESC,
  dsr.last_updated DESC,
  ds.last_updated DESC,
  dt.last_updated DESC
) s where deleted = false
ORDER BY set_number ASC
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	hostReplacer := strings.NewReplacer(`\`, ``, `.*`, ``)

	dsInfs := map[int][]DSRouteInfo{}
	for rows.Next() {
		dsID := 0
		dsName := ""
		dsType := ""
		dsPattern := ""
		dsRoutingName := ""
		inf := DSRouteInfo{}
		if err := rows.Scan(&dsID, &dsName, &dsType, &dsRoutingName, &dsPattern); err != nil {
			return nil, errors.New("scanning server deliveryservices: " + err.Error())
		}
		inf.DSName = dsName
		inf.IsDNS = strings.HasPrefix(dsType, "DNS")
		inf.IsRaw = !strings.Contains(dsPattern, `.*`)
		if !inf.IsRaw {
			host := hostReplacer.Replace(dsPattern)
			if inf.IsDNS {
				inf.Remap = dsRoutingName + host + domain
			} else {
				inf.Remap = host + domain
			}
		} else {
			inf.Remap = dsPattern
		}
		dsInfs[dsID] = append(dsInfs[dsID], inf)
	}

	serverDSPatterns := map[string]map[string][]string{}
	for serverID, dsIDs := range serverDSNames {
		server, ok := serverIDNames[serverID]
		if !ok {
			log.Errorf("Creating CRConfig: getting server dses: server ID %+v not in allServers, skipping!", serverID)
			continue
		}
		for _, dsID := range dsIDs {
			dsInfList, ok := dsInfs[dsID]
			if !ok {
				log.Warnf("Creating CRConfig: deliveryservice %v has no regexes, skipping", dsID)
				continue
			}
			for _, dsInf := range dsInfList {
				if !dsInf.IsRaw && !dsInf.IsDNS {
					dsInf.Remap = string(server) + dsInf.Remap
				}
				if _, ok := serverDSPatterns[server]; !ok {
					serverDSPatterns[server] = map[string][]string{}
				}
				serverDSPatterns[server][dsInf.DSName] = append(serverDSPatterns[server][dsInf.DSName], dsInf.Remap)
			}
		}
	}
	return serverDSPatterns, nil
}

// ServerParams contains parameter data filled in the CRConfig Servers objects. If a given param doesn't exist on the given server, it will be nil.
type ServerParams struct {
	APIPort          *string
	SecureAPIPort    *string
	Weight           *float64
	WeightMultiplier *float64
}

func getServerParams(tx *sql.Tx, cdn string, live bool) (map[string]ServerParams, error) {
	params := map[string]ServerParams{}

	qry := `
WITH cdn_name AS (
  SELECT $1::text as v
),
snapshot_time AS (
  SELECT time as v FROM snapshot sn where sn.cdn = (SELECT v from cdn_name)
)
SELECT server_name, param_name, param_val FROM (
SELECT DISTINCT ON (s.host_name, p.name)
  s.host_name as server_name,
  p.name as param_name,
  p.value as param_val,
  s.deleted as server_deleted
FROM
  server_snapshot s
  LEFT JOIN profile_parameter_snapshot pp ON pp.profile = s.profile
  LEFT JOIN parameter_snapshot p ON p.id = pp.parameter
  JOIN status_snapshot st ON st.id = s.status
WHERE
  s.cdn_id = (SELECT id FROM cdn_snapshot c where c.name = (select v from cdn_name) and c.last_updated <= (select v from snapshot_time))
  AND (
    (p.config_file = 'CRConfig.json' AND (p.name = 'weight' or p.name = 'weightMultiplier'))
    OR (p.name = 'api.port') OR (p.name = 'secure.api.port')
	)
  AND (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	if !live {
		qry += `
  AND s.last_updated <= (select v from snapshot_time)
  AND pp.last_updated <= (select v from snapshot_time)
  AND p.last_updated <= (select v from snapshot_time)
  AND st.last_updated <= (select v from snapshot_time)
`
	}
	qry += `
ORDER BY
  s.host_name DESC,
  p.name DESC,
  s.last_updated DESC,
  pp.last_updated DESC,
  p.last_updated DESC,
  st.last_updated DESC
) s WHERE server_deleted = false
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying server parameters: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		server := ""
		name := ""
		val := ""
		if err := rows.Scan(&server, &name, &val); err != nil {
			return nil, errors.New("scanning server parameters: " + err.Error())
		}

		param := params[server]
		switch name {
		case "api.port":
			param.APIPort = &val
		case "secure.api.port":
			param.SecureAPIPort = &val
		case "weight":
			i, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Warnln("Creating CRConfig: server " + server + " weight param " + val + " not a number, ignoring")
				continue
			}
			param.Weight = &i
		case "weightMultiplier":
			i, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Warnln("Creating CRConfig: server " + server + " weightMultiplier param " + val + " not a number, ignoring")
				continue
			}
			param.WeightMultiplier = &i
		}
		params[server] = param
	}
	return params, nil
}

// // WithDSSnapshotTimes returns DS snapshot times as a "with" query part.
// // Note this exists, so we can fake the times to get a "live" snapshot. Queries could obviously just query the table in other queries, but by using this, WithDSSnapshotTimesLive can be swapped in.
// func WithDSSnapshotTimes() string {
// 	return `
// WITH ds_snapshot_time AS (
//   SELECT deliveryservice, time from deliveryservice_snapshots
// )
// `
// }

// getCDNInfo returns the CDN domain, whether DNSSec is enabled, and whether the CDN has a snapshot, from the _snapshot tables.
// If the CDN has no snapshot, the domain will be blank and DNSSEC enabled will be false.
// The live argument is whether to get the latest data, not the snapshot time. Note this still queries the snapshot tables, so live calls must be preceded by populating the snapshot tables, e.g. with UpdateSnapshotTables.
// If live is true, the returned snapshot time will be now.
func getCDNSnapshotInfo(tx *sql.Tx, cdn string, live bool) (string, bool, bool, error) {
	qryArgs := []interface{}{}
	withCDNSnapshotTimeQueryPart, qryArgs := WithCDNSnapshotTime(cdn, live, qryArgs)
	qry := `
WITH ` + withCDNSnapshotTimeQueryPart + `,
 ` + CDNTable.WithLatest() + `
SELECT
  domain_name,
  dnssec_enabled
FROM ` + CDNTable.SnapshotLatestTable() + ` c
WHERE c.name = $` + strconv.Itoa(len(qryArgs)+1) + `
`
	qryArgs = append(qryArgs, cdn)

	domain := ""
	dnssec := false
	if err := tx.QueryRow(qry, qryArgs...).Scan(&domain, &dnssec); err != nil {
		if err == sql.ErrNoRows {
			return "", false, false, nil
		}
		return "", false, false, errors.New("querying CDN domain name: " + err.Error())
	}
	return domain, dnssec, true, nil
}

// getCDNNameFromID returns the CDN name given the ID, false if the no CDN with the given ID exists, and an error if the database query fails.
func getCDNNameFromID(id int, tx *sql.Tx) (string, bool, error) {
	// TODO change to use snapshot tables
	name := ""
	if err := tx.QueryRow(`SELECT name FROM cdn WHERE id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying CDN name: " + err.Error())
	}
	return name, true, nil
}

// getGlobalParam returns the global parameter with the requested name, whether it existed, and any error
func getGlobalParam(tx *sql.Tx, name string) (string, bool, error) {
	// TODO change to use snapshot tables?
	val := ""
	if err := tx.QueryRow(`SELECT value FROM parameter WHERE config_file = 'global' and name = $1`, name).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying global parameter '" + name + "': " + err.Error())
	}
	return val, true, nil
}

// allServersToServerIDNames returns a map[serverID]serverHostName
func allServersToServerIDNames(servers map[string]ServerUnion) map[int]string {
	serverIDNames := make(map[int]string, len(servers))
	for serverName, server := range servers {
		serverIDNames[server.ID] = serverName
	}
	return serverIDNames
}

// GetCRconfigSnapshotTime returns the latest time of the CRConfig snapshot, whether any snapshots were found, and any error.
// Note this is the max of the CDN's snapshot and the snapshots of all delivery services on that CDN.
func GetCRConfigSnapshotTime(tx *sql.Tx, cdnName tc.CDNName) (time.Time, bool, error) {
	qry := `
WITH cdn_name AS (
  SELECT $1::text AS v
)
SELECT MAX(time) FROM (
SELECT
  MAX(time) as time
FROM
  deliveryservice_snapshots dsn
  JOIN deliveryservice ds ON ds.xml_id = dsn.deliveryservice
  JOIN cdn c ON c.id = ds.cdn_id
WHERE
  c.name = (select v from cdn_name)
UNION ALL
SELECT time FROM snapshot WHERE cdn = (select v from cdn_name)
) t
`
	t := pq.NullTime{}
	if err := tx.QueryRow(qry, cdnName).Scan(&t); err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, false, nil
		}
		return time.Time{}, false, errors.New("Error querying CDN snapshot time: " + err.Error())
	}
	return t.Time, t.Valid, nil
}
