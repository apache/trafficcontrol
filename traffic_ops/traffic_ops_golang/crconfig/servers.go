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

const RouterTypeName = "CCR"
const MonitorTypeName = "RASCAL"
const EdgeTypePrefix = "EDGE"
const MidTypePrefix = "MID"

// makeCRConfigServers returns:
// - the map of cache servers
// - the map of routers
// - the map of monitors
// - the last time anything was modified _not including_ the Delivery Service Server mappings.
//   - if you need the last time including DSS mappings, it can be computed from the returned
//     value and the passed serverDSNames.
// - any error
func makeCRConfigServers(cdn string, tx *sql.Tx, cdnDomain string, serverDSNames map[tc.CacheName][]ServerDS) (
	map[string]tc.CRConfigServer,
	map[string]tc.CRConfigRouter,
	map[string]tc.CRConfigMonitor,
	time.Time,
	error,
) {
	allServers, lastUpdated, err := getAllServers(cdn, tx)
	if err != nil {
		return nil, nil, nil, time.Time{}, err
	}

	serverDSes, _, err := getServerDSes(cdn, tx, cdnDomain, serverDSNames)
	if err != nil {
		return nil, nil, nil, time.Time{}, errors.New("getting server deliveryservices: " + err.Error())
	}
	// if serverDSLastUpdated.After(lastUpdated) {
	// 	lastUpdated = serverDSLastUpdated
	// }

	servers := map[string]tc.CRConfigServer{}
	routers := map[string]tc.CRConfigRouter{}
	monitors := map[string]tc.CRConfigMonitor{}

	for host, s := range allServers {
		switch {
		case *s.ServerType == RouterTypeName:
			status := tc.CRConfigRouterStatus(*s.ServerStatus)
			routers[host] = tc.CRConfigRouter{
				APIPort:      s.APIPort,
				FQDN:         s.Fqdn,
				HTTPSPort:    s.HttpsPort,
				IP:           s.Ip,
				IP6:          s.Ip6,
				Location:     s.LocationId,
				Port:         s.Port,
				Profile:      s.Profile,
				ServerStatus: &status,
			}
		case *s.ServerType == MonitorTypeName:
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
		case strings.HasPrefix(*s.ServerType, EdgeTypePrefix) || strings.HasPrefix(*s.ServerType, MidTypePrefix):
			if s.RoutingDisabled == 0 {
				s.CRConfigServer.DeliveryServices = serverDSes[tc.CacheName(host)]

				latestDSModified := time.Time{}
				for _, ds := range serverDSNames[tc.CacheName(host)] {
					if ds.Modified.After(latestDSModified) {
						latestDSModified = ds.Modified
					}
				}
				s.CRConfigServer.DeliveryServicesModified = latestDSModified
			}
			servers[host] = s.CRConfigServer
		}
	}
	return servers, routers, monitors, lastUpdated, nil
}

// ServerUnion has all fields from all servers. This is used to select all server data with a single query, and then convert each to the proper type afterwards.
type ServerUnion struct {
	tc.CRConfigServer
	APIPort *string
}

const DefaultWeightMultiplier = 1000.0
const DefaultWeight = 0.999

// getAllServers returns:
//  - the map of server names to servers
//  - the last updated time of any server (including both the server table and any params used)
//  - any error
func getAllServers(cdn string, tx *sql.Tx) (map[string]ServerUnion, time.Time, error) {
	servers := map[string]ServerUnion{}

	serverParams, err := getServerParams(cdn, tx)
	if err != nil {
		return nil, time.Time{}, errors.New("Error getting server params: " + err.Error())
	}

	// TODO select deliveryservices as array?
	qry := `
SELECT
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
  GREATEST(s.last_updated, cg.last_updated, t.last_updated, p.last_updated, st.last_updated, cdn.last_updated) as last_updated
FROM
  server s
  JOIN cachegroup cg ON cg.id = s.cachegroup
  JOIN type t ON t.id = s.type
  JOIN profile p ON p.id = s.profile
  JOIN status st ON st.id = s.status
  JOIN cdn on cdn.id = s.cdn_id
WHERE
  cdn.name = $1
  AND (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, time.Time{}, errors.New("Error querying servers: " + err.Error())
	}
	defer rows.Close()

	lastUpdated := time.Time{}

	for rows.Next() {
		port := sql.NullInt64{}
		ip6 := sql.NullString{}
		hashId := sql.NullString{}
		httpsPort := sql.NullInt64{}

		s := ServerUnion{}

		host := ""
		status := ""
		serverLastUpdated := time.Time{}
		if err := rows.Scan(&host, &s.CacheGroup, &s.Fqdn, &hashId, &httpsPort, &s.InterfaceName, &s.Ip, &ip6, &port, &s.Profile, &s.RoutingDisabled, &status, &s.ServerType, &serverLastUpdated); err != nil {
			return nil, time.Time{}, errors.New("Error scanning server: " + err.Error())
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
			s.APIPort = &params.APIPort.S
			if params.APIPort.M.After(lastUpdated) {
				lastUpdated = params.APIPort.M
			}
		}

		weightMultiplier := DefaultWeightMultiplier
		if hasParams && params.WeightMultiplier != nil {
			weightMultiplier = params.WeightMultiplier.F
			if params.WeightMultiplier.M.After(lastUpdated) {
				lastUpdated = params.WeightMultiplier.M
			}
		}
		weight := DefaultWeight
		if hasParams && params.Weight != nil {
			weight = params.Weight.F
			if params.Weight.M.After(lastUpdated) {
				lastUpdated = params.Weight.M
			}
		}
		hashCount := int(weight * weightMultiplier)
		s.HashCount = &hashCount

		servers[host] = s
		if serverLastUpdated.After(lastUpdated) {
			lastUpdated = serverLastUpdated
		}
	}
	if err := rows.Err(); err != nil {
		return nil, time.Time{}, errors.New("Error iterating router param rows: " + err.Error())
	}

	return servers, lastUpdated, nil
}

type ServerDS struct {
	DS       tc.DeliveryServiceName
	Modified time.Time
}

func getServerDSNames(cdn string, tx *sql.Tx) (map[tc.CacheName][]ServerDS, error) {
	q := `
select s.host_name, ds.xml_id, dss.last_updated
from deliveryservice_server as dss
inner join server as s on dss.server = s.id
inner join deliveryservice as ds on ds.id = dss.deliveryservice
inner join profile as p on p.id = s.profile
inner join status as st ON st.id = s.status
where ds.cdn_id = (select id from cdn where name = $1)
and ds.active = true
and p.routing_disabled = false
and (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, errors.New("Error querying server deliveryservice names: " + err.Error())
	}
	defer rows.Close()

	serverDSes := map[tc.CacheName][]ServerDS{}
	for rows.Next() {
		ds := ""
		server := ""
		updated := time.Time{}
		if err := rows.Scan(&server, &ds, &updated); err != nil {
			return nil, errors.New("Error scanning server deliveryservice names: " + err.Error())
		}
		serverDSes[tc.CacheName(server)] = append(serverDSes[tc.CacheName(server)], ServerDS{DS: tc.DeliveryServiceName(ds), Modified: updated})
	}
	return serverDSes, nil
}

type DSRouteInfo struct {
	IsDNS bool
	IsRaw bool
	Remap string
}

func getServerDSes(cdn string, tx *sql.Tx, domain string, serverDSNames map[tc.CacheName][]ServerDS) (map[tc.CacheName]map[string][]string, time.Time, error) {
	qry := `
SELECT
  ds.xml_id as ds,
  dt.name as ds_type,
  ds.routing_name,
  r.pattern as pattern,
  GREATEST(r.last_updated, rt.last_updated, dsr.last_updated, ds.last_updated, dt.last_updated, cdn.last_updated) as last_updated
FROM
  regex r
  JOIN type rt on r.type = rt.id
  JOIN deliveryservice_regex dsr on dsr.regex = r.id
  JOIN deliveryservice ds on ds.id = dsr.deliveryservice
  JOIN type dt on dt.id = ds.type
  JOIN cdn on cdn.id = ds.cdn_id
WHERE
  cdn.name = $1
  AND ds.active = true
  AND rt.name = 'HOST_REGEXP'
ORDER BY
  dsr.set_number asc
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, time.Time{}, errors.New("Error server deliveryservices: " + err.Error())
	}
	defer rows.Close()

	hostReplacer := strings.NewReplacer(`\`, ``, `.*`, ``)

	lastUpdated := time.Time{}
	dsInfs := map[string][]DSRouteInfo{}
	for rows.Next() {
		ds := ""
		dsType := ""
		dsPattern := ""
		dsRoutingName := ""
		dsinfLastUpdated := time.Time{}
		inf := DSRouteInfo{}
		if err := rows.Scan(&ds, &dsType, &dsRoutingName, &dsPattern, &dsinfLastUpdated); err != nil {
			return nil, time.Time{}, errors.New("Error scanning server deliveryservices: " + err.Error())
		}
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
		dsInfs[ds] = append(dsInfs[ds], inf)
		if dsinfLastUpdated.After(lastUpdated) {
			lastUpdated = dsinfLastUpdated
		}
	}

	serverDSPatterns := map[tc.CacheName]map[string][]string{}
	for server, dses := range serverDSNames {
		for _, dsName := range dses {
			dsInfList, ok := dsInfs[string(dsName.DS)]
			if !ok {
				log.Warnln("Creating CRConfig: deliveryservice " + string(dsName.DS) + " has no regexes, skipping")
				continue
			}
			for _, dsInf := range dsInfList {
				if !dsInf.IsRaw && !dsInf.IsDNS {
					dsInf.Remap = string(server) + dsInf.Remap
				}
				if _, ok := serverDSPatterns[server]; !ok {
					serverDSPatterns[server] = map[string][]string{}
				}
				serverDSPatterns[server][string(dsName.DS)] = append(serverDSPatterns[server][string(dsName.DS)], dsInf.Remap)
			}
		}
	}
	return serverDSPatterns, lastUpdated, nil
}

type StrModified struct {
	S string
	M time.Time
}

type FloatModified struct {
	F float64
	M time.Time
}

// ServerParams contains parameter data filled in the CRConfig Servers objects. If a given param doesn't exist on the given server, it will be nil.
type ServerParams struct {
	APIPort          *StrModified
	Weight           *FloatModified
	WeightMultiplier *FloatModified
}

func getServerParams(cdn string, tx *sql.Tx) (map[string]ServerParams, error) {
	params := map[string]ServerParams{}

	q := `
SELECT
  s.host_name,
  p.name,
  p.value,
  GREATEST(s.last_updated, p.last_updated, pp.last_updated) as last_updated
FROM
  server s
  JOIN profile_parameter as pp on pp.profile = s.profile
  JOIN parameter as p on p.id = pp.parameter
  JOIN cdn on cdn.id = s.cdn_id
  INNER JOIN status as st ON st.id = s.status
WHERE
  cdn.name = $1
  AND ((p.config_file = 'CRConfig.json' and (p.name = 'weight' or p.name = 'weightMultiplier')) or (p.name = 'api.port'))
  AND (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, errors.New("Error querying server parameters: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		server := ""
		name := ""
		val := ""
		lastModified := time.Time{}
		if err := rows.Scan(&server, &name, &val, &lastModified); err != nil {
			return nil, errors.New("Error scanning server parameters: " + err.Error())
		}

		param := params[server]
		switch name {
		case "api.port":
			param.APIPort = &StrModified{S: val, M: lastModified}
		case "weight":
			i, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Warnln("Creating CRConfig: server " + server + " weight param " + val + " not a number, ignoring")
				continue
			}
			param.Weight = &FloatModified{F: i, M: lastModified}
		case "weightMultiplier":
			i, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Warnln("Creating CRConfig: server " + server + " weightMultiplier param " + val + " not a number, ignoring")
				continue
			}
			param.WeightMultiplier = &FloatModified{F: i, M: lastModified}
		}
		params[server] = param
	}
	return params, nil
}

// getCDNInfo returns the CDN domain, and whether DNSSec is enabled
func getCDNInfo(cdn string, tx *sql.Tx) (string, bool, time.Time, error) {
	domain := ""
	dnssec := false
	lastUpdated := time.Time{}
	if err := tx.QueryRow(`select domain_name, dnssec_enabled, last_updated from cdn where name = $1`, cdn).Scan(&domain, &dnssec, &lastUpdated); err != nil {
		return "", false, time.Time{}, errors.New("Error querying CDN domain name: " + err.Error())
	}
	return domain, dnssec, lastUpdated, nil
}

// getCDNNameFromID returns the CDN name given the ID, false if the no CDN with the given ID exists, and an error if the database query fails.
func getCDNNameFromID(id int, tx *sql.Tx) (string, bool, error) {
	name := ""
	if err := tx.QueryRow(`select name from cdn where id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("Error querying CDN name: " + err.Error())
	}
	return name, true, nil
}
