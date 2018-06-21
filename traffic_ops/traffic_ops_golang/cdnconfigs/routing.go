package cdnconfigs

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
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

const RouterProfilePrefix = "CCR"
const MonitorType = "RASCAL"
const RouterType = "CCR"
const EdgeTypePrefix = "EDGE"
const MidTypePrefix = "MID"
const DefaultRouterAPIPort = 3333

func Routing(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName := tc.CDNName(inf.Params["name"])

	cfg, err := getCDNConfigRouting(r, inf, cdnName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cdn config routing: "+err.Error()))
		return
	}
	api.WriteResp(w, r, cfg)
}

func getCDNConfigRouting(req *http.Request, inf *api.APIInfo, cdnName tc.CDNName) (tc.TrafficRouterConfig, error) {
	cfg := tc.TrafficRouterConfig{}
	cfg.Stats = getCDNConfigRoutingStats(req, inf, cdnName)
	routerProfileID, ok, err := getRouterProfileID(inf.Tx.Tx, cdnName)
	if err != nil {
		return tc.TrafficRouterConfig{}, errors.New("getting router profile ID: " + err.Error())
	} else if !ok {
		return tc.TrafficRouterConfig{}, errors.New("No CCR router profile found for CDN")
	}
	if cfg.CacheGroups, err = getCDNConfigRoutingCachegroups(inf.Tx.Tx, cdnName); err != nil {
		return tc.TrafficRouterConfig{}, errors.New("getting routing cdn cachegroups: " + err.Error())
	}
	if cfg.Config, err = getCDNConfigRoutingConfig(inf.Tx.Tx, routerProfileID); err != nil {
		return tc.TrafficRouterConfig{}, errors.New("getting routing cdn config: " + err.Error())
	}
	if cfg.TrafficServers, cfg.TrafficMonitors, cfg.TrafficRouters, err = getCDNConfigRoutingServers(inf.Tx.Tx, cdnName); err != nil {
		return tc.TrafficRouterConfig{}, errors.New("getting routing servers: " + err.Error())
	}
	return cfg, nil
}

func getCDNConfigRoutingStats(req *http.Request, inf *api.APIInfo, cdnName tc.CDNName) map[string]interface{} {
	return map[string]interface{}{
		"cdnName":           cdnName,
		"date":              time.Now().Unix(),
		"trafficOpsVersion": inf.Config.Version,
		"trafficOpsPath":    req.URL.Path,
		"trafficOpsHost":    req.Host,
		"trafficOpsUser":    inf.User.UserName,
	}
}

func getCDNConfigRoutingCachegroups(tx *sql.Tx, cdnName tc.CDNName) ([]tc.TMCacheGroup, error) {
	q := `
SELECT distinct(cg.name), cg.latitude, cg.longitude
FROM cachegroup as cg
JOIN server as s ON s.cachegroup = cg.id
WHERE s.cdn_id = (select id from cdn where name = $1)
`
	rows, err := tx.Query(q, cdnName)
	if err != nil {
		return nil, errors.New("querying routing cachegroups: " + err.Error())
	}
	defer rows.Close()
	cgs := []tc.TMCacheGroup{}
	for rows.Next() {
		cg := tc.TMCacheGroup{}
		if err := rows.Scan(&cg.Name, &cg.Coordinates.Latitude, &cg.Coordinates.Longitude); err != nil {
			return nil, errors.New("scanning routing cachegroups: " + err.Error())
		}
		cgs = append(cgs, cg)
	}
	return cgs, nil
}

func getRouterProfileID(tx *sql.Tx, cdnName tc.CDNName) (int, bool, error) {
	q := `
SELECT distinct(p.id)
FROM profile as p
JOIN server as s ON s.profile = p.id
JOIN cdn on cdn.id = s.cdn_id
WHERE cdn.name = $1
AND p.name LIKE '` + RouterProfilePrefix + `%'
`
	routerProfileID := 0
	if err := tx.QueryRow(q, cdnName).Scan(&routerProfileID); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, errors.New("querying router profile ID: " + err.Error())
	}
	return routerProfileID, true, nil
}

func getCDNConfigRoutingConfig(tx *sql.Tx, routerProfileID int) (map[string]interface{}, error) {
	q := `
SELECT name, value
FROM parameter as p
JOIN profile_parameter as pp ON pp.parameter = p.id
WHERE pp.profile = $1
AND p.config_file = 'CRConfig.json'
`
	rows, err := tx.Query(q, routerProfileID)
	if err != nil {
		return nil, errors.New("querying router parameters: " + err.Error())
	}
	defer rows.Close()

	config := map[string]interface{}{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("scanning router parameters: " + err.Error())
		}
		if ival, err := strconv.Atoi(val); err == nil {
			config[name] = ival
		} else {
			config[name] = val
		}
	}
	return config, nil
}

func getCDNConfigRoutingServers(tx *sql.Tx, cdnName tc.CDNName) ([]tc.TrafficServer, []tc.TrafficMonitor, []tc.TrafficRouter, error) {
	apiPorts, err := getRouterAPIPorts(tx, cdnName)
	if err != nil {
		return nil, nil, nil, errors.New("getting router API ports: " + err.Error())
	}

	q := `
SELECT
s.host_name,
s.domain_name,
cg.name as cachegroup,
COALESCE(s.xmpp_id, '') as hashid,
s.interface_name,
s.ip_address,
COALESCE(s.ip6_address, ''),
COALESCE(s.tcp_port, 0),
p.name as profile,
st.name as status,
t.name as type
FROM server as s
JOIN cachegroup as cg on cg.id = s.cachegroup
JOIN profile as p on p.id = s.profile
JOIN status as st on st.id = s.status
JOIN type as t on t.id = s.type
WHERE s.cdn_id = (select id from cdn where name = $1)
`
	rows, err := tx.Query(q, cdnName)
	if err != nil {
		return nil, nil, nil, errors.New("querying config routing servers: " + err.Error())
	}
	defer rows.Close()

	servers := []tc.TrafficServer{}
	routers := []tc.TrafficRouter{}
	monitors := []tc.TrafficMonitor{}
	for rows.Next() {
		s := tc.TrafficServer{}
		domain := ""
		if err := rows.Scan(&s.HostName, &domain, &s.CacheGroup, &s.HashID, &s.InterfaceName, &s.IP, &s.IP6, &s.Port, &s.Profile, &s.ServerStatus, &s.Type); err != nil {
			return nil, nil, nil, errors.New("scanning router servers: " + err.Error())
		}
		s.FQDN = s.HostName + "." + domain

		if s.Type == MonitorType {
			monitors = append(monitors, serverToMonitor(s))
		} else if s.Type == RouterType {
			router := serverToRouter(s)
			if apiPort, ok := apiPorts[router.HostName]; ok {
				router.APIPort = apiPort
			}
			routers = append(routers, router)
		} else if strings.HasPrefix(s.Type, EdgeTypePrefix) || strings.HasPrefix(s.Type, MidTypePrefix) {
			s.Deliveryservices = []tc.TSDeliveryService{}
			servers = append(servers, s)
		}
	}
	return servers, monitors, routers, nil
}

func getRouterAPIPorts(tx *sql.Tx, cdnName tc.CDNName) (map[string]int, error) {
	q := `
SELECT host_name, pa.value as api_port
FROM server as s
JOIN profile as p on s.profile = p.id
JOIN profile_parameter as pp on pp.profile = p.id
JOIN parameter as pa on pa.id = pp.parameter
WHERE s.type = (select id from type where name = '` + RouterType + `')
AND s.cdn_id = (select id from cdn where name = $1)
AND pa.name = 'api.port'
`
	rows, err := tx.Query(q, cdnName)
	if err != nil {
		return nil, errors.New("querying config routing server api ports: " + err.Error())
	}
	defer rows.Close()

	ports := map[string]int{}
	for rows.Next() {
		host := ""
		portStr := ""
		if err := rows.Scan(&host, &portStr); err != nil {
			return nil, errors.New("scanning router server api ports: " + err.Error())
		}
		// TODO warn if api.port is not an int?
		if port, err := strconv.Atoi(portStr); err != nil {
			ports[host] = port
		}
	}
	return ports, nil
}

func serverToRouter(s tc.TrafficServer) tc.TrafficRouter {
	return tc.TrafficRouter{
		Port:         s.Port,
		IP6:          s.IP6,
		IP:           s.IP,
		FQDN:         s.FQDN,
		Profile:      s.Profile,
		Location:     s.CacheGroup,
		ServerStatus: s.ServerStatus,
		HostName:     s.HostName,
		APIPort:      DefaultRouterAPIPort,
	}
}

func serverToMonitor(s tc.TrafficServer) tc.TrafficMonitor {
	return tc.TrafficMonitor{
		Port:         s.Port,
		IP6:          s.IP6,
		IP:           s.IP,
		HostName:     s.HostName,
		FQDN:         s.FQDN,
		Profile:      s.Profile,
		Location:     s.CacheGroup,
		ServerStatus: s.ServerStatus,
	}
}
