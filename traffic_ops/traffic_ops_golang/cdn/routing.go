package cdn

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
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetRouting(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	api.RespWriter(w, r, inf.Tx.Tx)(getRouting(inf.Tx.Tx))
}

const RouterProxyParameter = "tm.traffic_rtr_fwd_proxy"
const RouterRequestTimeout = time.Second * 10
const RouterOnlineStatus = "ONLINE"

func getRouting(tx *sql.Tx) (tc.CDNRouting, error) {
	routers, err := getCDNRouterFQDNs(tx)
	if err != nil {
		return tc.CDNRouting{}, errors.New("getting monitors: " + err.Error())
	}
	return getRoutersRouting(tx, routers)
}

func getRoutersRouting(tx *sql.Tx, routers map[tc.CDNName][]string) (tc.CDNRouting, error) {
	forwardProxy, forwardProxyExists, err := getGlobalParam(tx, RouterProxyParameter)
	if err != nil {
		return tc.CDNRouting{}, errors.New("getting global router proxy parameter: " + err.Error())
	}
	client := &http.Client{Timeout: RouterRequestTimeout}
	if forwardProxyExists {
		proxyURI, err := url.Parse(forwardProxy)
		if err != nil {
			return tc.CDNRouting{}, errors.New("router forward proxy '" + forwardProxy + "' in parameter '" + RouterProxyParameter + "' not a URI: " + err.Error())
		}
		client = &http.Client{Timeout: RouterRequestTimeout, Transport: &http.Transport{Proxy: http.ProxyURL(proxyURI)}}
	}

	dat := RouterData{}
	for cdn, routerFQDNs := range routers {
		for _, routerFQDN := range routerFQDNs {
			crsStats := tc.CRSStats{}
			if err := getCRSStats(routerFQDN, client, &crsStats); err != nil {
				return tc.CDNRouting{}, errors.New("getting crs stats for CDN '" + string(cdn) + "' router '" + routerFQDN + "': " + err.Error())
			}
			dat = addCRSStats(dat, crsStats)
		}
	}
	return sumRouterData(dat), nil
}

type RouterData struct {
	StatTotal tc.CRSStatsStat
	Total     uint64
}

func sumRouterData(d RouterData) tc.CDNRouting {
	if d.Total == 0 {
		return tc.CDNRouting{} // TODO warn?
	}
	return tc.CDNRouting{
		CZ:                float64(d.StatTotal.CZCount) / float64(d.Total) * 100.0,
		Geo:               float64(d.StatTotal.GeoCount) / float64(d.Total) * 100.0,
		DeepCZ:            float64(d.StatTotal.DeepCZCount) / float64(d.Total) * 100.0,
		Miss:              float64(d.StatTotal.MissCount) / float64(d.Total) * 100.0,
		DSR:               float64(d.StatTotal.DSRCount) / float64(d.Total) * 100.0,
		Err:               float64(d.StatTotal.ErrCount) / float64(d.Total) * 100.0,
		StaticRoute:       float64(d.StatTotal.StaticRouteCount) / float64(d.Total) * 100.0,
		Fed:               float64(d.StatTotal.FedCount) / float64(d.Total) * 100.0,
		RegionalDenied:    float64(d.StatTotal.RegionalDeniedCount) / float64(d.Total) * 100.0,
		RegionalAlternate: float64(d.StatTotal.RegionalAlternateCount) / float64(d.Total) * 100.0,
	}
}

func addCRSStats(d RouterData, stats tc.CRSStats) RouterData {
	for _, stat := range stats.Stats.DNSMap {
		d.StatTotal = sumCRSStat(d.StatTotal, stat)
		d.Total += totalCRSStat(stat)
	}
	return d
}

func totalCRSStat(s tc.CRSStatsStat) uint64 {
	return s.CZCount +
		s.GeoCount +
		s.DeepCZCount +
		s.MissCount +
		s.DSRCount +
		s.ErrCount +
		s.StaticRouteCount +
		s.FedCount +
		s.RegionalDeniedCount +
		s.RegionalAlternateCount
}

func sumCRSStat(a, b tc.CRSStatsStat) tc.CRSStatsStat {
	return tc.CRSStatsStat{
		CZCount:                a.CZCount + b.CZCount,
		GeoCount:               a.GeoCount + b.GeoCount,
		DeepCZCount:            a.DeepCZCount + b.DeepCZCount,
		MissCount:              a.MissCount + b.MissCount,
		DSRCount:               a.DSRCount + b.DSRCount,
		ErrCount:               a.ErrCount + b.ErrCount,
		StaticRouteCount:       a.StaticRouteCount + b.StaticRouteCount,
		FedCount:               a.FedCount + b.FedCount,
		RegionalDeniedCount:    a.RegionalDeniedCount + b.RegionalDeniedCount,
		RegionalAlternateCount: a.RegionalAlternateCount + b.RegionalAlternateCount,
	}
}

func getCRSStats(routerFQDN string, client *http.Client, stats interface{}) error {
	path := `/crs/stats`
	resp, err := client.Get("http://" + routerFQDN + path)
	if err != nil {
		return errors.New("getting crs stats from router '" + routerFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(stats); err != nil {
		return errors.New("decoding stats from router '" + routerFQDN + "': " + err.Error())
	}
	return nil
}

func getRouterForwardProxy(tx *sql.Tx) (string, error) {
	forwardProxy, forwardProxyExists, err := getGlobalParam(tx, RouterProxyParameter)
	if err != nil {
		return "", errors.New("getting global router proxy parameter: " + err.Error())
	} else if !forwardProxyExists {
		forwardProxy = ""
	}
	return forwardProxy, nil
}

// getCDNRouterFQDNs returns an FQDN, including port, of an online router for each CDN, for each router. If a CDN has no online routers, that CDN will not have an entry in the map. The port returned is the API port.
func getCDNRouterFQDNs(tx *sql.Tx) (map[tc.CDNName][]string, error) {
	rows, err := tx.Query(`
SELECT s.host_name, s.domain_name, max(pa.value) as port, c.name as cdn
FROM server as s
JOIN type as t ON s.type = t.id
JOIN status as st ON st.id = s.status
JOIN cdn as c ON c.id = s.cdn_id
JOIN profile as pr ON s.profile = pr.id
JOIN profile_parameter as pp ON pp.profile = pr.id
LEFT JOIN parameter as pa ON (pp.parameter = pa.id AND pa.name = 'api.port' AND pa.config_file = 'server.xml')
WHERE t.name = '` + tc.RouterTypeName + `'
AND st.name = '` + RouterOnlineStatus + `'
GROUP BY s.host_name, s.domain_name, c.name
`)
	if err != nil {
		return nil, errors.New("querying routers: " + err.Error())
	}
	defer rows.Close()
	routers := map[tc.CDNName][]string{}
	for rows.Next() {
		host := ""
		domain := ""
		port := sql.NullInt64{}
		cdn := ""
		if err := rows.Scan(&host, &domain, &port, &cdn); err != nil {
			return nil, errors.New("scanning routers: " + err.Error())
		}
		fqdn := host + "." + domain
		if port.Valid {
			fqdn += ":" + strconv.FormatInt(port.Int64, 10)
		}
		routers[tc.CDNName(cdn)] = append(routers[tc.CDNName(cdn)], fqdn)
	}
	return routers, nil
}

// getGlobalParams returns the value of the global param, whether it existed, or any error
func getGlobalParam(tx *sql.Tx, name string) (string, bool, error) {
	return getParam(tx, name, "global")
}

// getGlobalParams returns the value of the param, whether it existed, or any error.
func getParam(tx *sql.Tx, name string, configFile string) (string, bool, error) {
	val := ""
	if err := tx.QueryRow(`select value from parameter where name = $1 and config_file = $2`, name, configFile).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("Error querying global paramter '" + name + "': " + err.Error())
	}
	return val, true, nil
}
