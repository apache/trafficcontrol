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
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

const RouterTypeName = "CCR"
const MonitorTypeName = "RASCAL"
const EdgeTypePrefix = "EDGE"
const MidTypePrefix = "MID"

func makeCRConfigServers(cdn string, tx *sql.Tx, cdnDomain string) (
	map[string]tc.CRConfigTrafficOpsServer,
	map[string]tc.CRConfigRouter,
	map[string]tc.CRConfigMonitor,
	error,
) {
	allServers, err := getAllServers(cdn, tx)
	if err != nil {
		return nil, nil, nil, err
	}

	serverDSes, err := getServerDSes(cdn, tx, cdnDomain)
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
				s.CRConfigTrafficOpsServer.DeliveryServices = serverDSes[tc.CacheName(host)]
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
}

const DefaultWeightMultiplier = 1000.0
const DefaultWeight = 0.999

func getAllServers(cdn string, tx *sql.Tx) (map[string]ServerUnion, error) {
	servers := map[string]ServerUnion{}

	serverParams, err := getServerParams(cdn, tx)
	if err != nil {
		return nil, errors.New("Error getting server params: " + err.Error())
	}

	// TODO select deliveryservices as array?
	q := `
select s.host_name, cg.name as cachegroup, concat(s.host_name, '.', s.domain_name) as fqdn, s.xmpp_id as hashid, s.https_port, s.interface_name, s.ip_address, s.ip6_address, s.tcp_port, p.name as profile_name, cast(p.routing_disabled as int), st.name as status, t.name as type
from server as s
inner join cachegroup as cg ON cg.id = s.cachegroup
inner join type as t on t.id = s.type
inner join profile as p ON p.id = s.profile
inner join status as st ON st.id = s.status
where cdn_id = (select id from cdn where name = $1)
and (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, errors.New("Error querying servers: " + err.Error())
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
		if err := rows.Scan(&host, &s.CacheGroup, &s.Fqdn, &hashId, &httpsPort, &s.InterfaceName, &s.Ip, &ip6, &port, &s.Profile, &s.RoutingDisabled, &status, &s.ServerType); err != nil {
			return nil, errors.New("Error scanning server: " + err.Error())
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
		return nil, errors.New("Error iterating router param rows: " + err.Error())
	}

	return servers, nil
}

func getServerDSNames(cdn string, tx *sql.Tx) (map[tc.CacheName][]tc.DeliveryServiceName, error) {
	q := `
select s.host_name, ds.xml_id
from deliveryservice_server as dss
inner join server as s on dss.server = s.id
inner join deliveryservice as ds on ds.id = dss.deliveryservice
inner join type as dt on dt.id = ds.type
inner join profile as p on p.id = s.profile
inner join status as st ON st.id = s.status
where ds.cdn_id = (select id from cdn where name = $1)
and ds.active = true` +
		fmt.Sprintf(" and dt.name != '%s' ", tc.DSTypeAnyMap) + `
and p.routing_disabled = false
and (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, errors.New("Error querying server deliveryservice names: " + err.Error())
	}
	defer rows.Close()

	serverDSes := map[tc.CacheName][]tc.DeliveryServiceName{}
	for rows.Next() {
		ds := ""
		server := ""
		if err := rows.Scan(&server, &ds); err != nil {
			return nil, errors.New("Error scanning server deliveryservice names: " + err.Error())
		}
		serverDSes[tc.CacheName(server)] = append(serverDSes[tc.CacheName(server)], tc.DeliveryServiceName(ds))
	}
	return serverDSes, nil
}

type DSRouteInfo struct {
	IsDNS bool
	IsRaw bool
	Remap string
}

func getServerDSes(cdn string, tx *sql.Tx, domain string) (map[tc.CacheName]map[string][]string, error) {
	serverDSNames, err := getServerDSNames(cdn, tx)
	if err != nil {
		return nil, errors.New("Error getting server deliveryservices: " + err.Error())
	}

	q := `
select ds.xml_id as ds, dt.name as ds_type, ds.routing_name, r.pattern as pattern
from regex as r
inner join type as rt on r.type = rt.id
inner join deliveryservice_regex as dsr on dsr.regex = r.id
inner join deliveryservice as ds on ds.id = dsr.deliveryservice
inner join type as dt on dt.id = ds.type
where ds.cdn_id = (select id from cdn where name = $1)
and ds.active = true` +
		fmt.Sprintf(" and dt.name != '%s' ", tc.DSTypeAnyMap) + `
and rt.name = 'HOST_REGEXP'
order by dsr.set_number asc
`
	rows, err := tx.Query(q, cdn)
	if err != nil {
		return nil, errors.New("Error server deliveryservices: " + err.Error())
	}
	defer rows.Close()

	hostReplacer := strings.NewReplacer(`\`, ``, `.*`, ``)

	dsInfs := map[string][]DSRouteInfo{}
	for rows.Next() {
		ds := ""
		dsType := ""
		dsPattern := ""
		dsRoutingName := ""
		inf := DSRouteInfo{}
		if err := rows.Scan(&ds, &dsType, &dsRoutingName, &dsPattern); err != nil {
			return nil, errors.New("Error scanning server deliveryservices: " + err.Error())
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
	}

	serverDSPatterns := map[tc.CacheName]map[string][]string{}
	for server, dses := range serverDSNames {
		for _, dsName := range dses {
			dsInfList, ok := dsInfs[string(dsName)]
			if !ok {
				log.Warnln("Creating CRConfig: deliveryservice " + string(dsName) + " has no regexes, skipping")
				continue
			}
			for _, dsInf := range dsInfList {
				if !dsInf.IsRaw && !dsInf.IsDNS {
					dsInf.Remap = string(server) + dsInf.Remap
				}
				if _, ok := serverDSPatterns[server]; !ok {
					serverDSPatterns[server] = map[string][]string{}
				}
				serverDSPatterns[server][string(dsName)] = append(serverDSPatterns[server][string(dsName)], dsInf.Remap)
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

func getServerParams(cdn string, tx *sql.Tx) (map[string]ServerParams, error) {
	params := map[string]ServerParams{}

	q := `
select s.host_name, p.name, p.value
from server as s
left join profile_parameter as pp on pp.profile = s.profile
left join parameter as p on p.id = pp.parameter
inner join status as st ON st.id = s.status
where s.cdn_id = (select id from cdn where name = $1)
and ((p.config_file = 'CRConfig.json' and (p.name = 'weight' or p.name = 'weightMultiplier')) or (p.name = 'api.port') or (p.name = 'secure.api.port'))
and (st.name = 'REPORTED' or st.name = 'ONLINE' or st.name = 'ADMIN_DOWN')
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
		if err := rows.Scan(&server, &name, &val); err != nil {
			return nil, errors.New("Error scanning server parameters: " + err.Error())
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

// getCDNInfo returns the CDN domain, and whether DNSSec is enabled
func getCDNInfo(cdn string, tx *sql.Tx) (string, bool, error) {
	domain := ""
	dnssec := false
	if err := tx.QueryRow(`select domain_name, dnssec_enabled from cdn where name = $1`, cdn).Scan(&domain, &dnssec); err != nil {
		return "", false, errors.New("Error querying CDN domain name: " + err.Error())
	}
	return domain, dnssec, nil
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

// getGlobalParam returns the global parameter with the requested name, whether it existed, and any error
func getGlobalParam(tx *sql.Tx, name string) (string, bool, error) {
	val := ""
	if err := tx.QueryRow(`SELECT value FROM parameter WHERE config_file = 'global' and name = $1`, name).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying global parameter '" + name + "': " + err.Error())
	}
	return val, true, nil
}
