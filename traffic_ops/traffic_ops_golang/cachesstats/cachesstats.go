package cachesstats

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
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/monitorhlp"
)

func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	api.RespWriter(w, r, inf.Tx.Tx)(getCachesStats(inf.Tx.Tx))
}

const MonitorOnlineStatus = "ONLINE"

func getCachesStats(tx *sql.Tx) ([]CacheData, error) {
	monitors, err := getCDNMonitorFQDNs(tx)
	if err != nil {
		return nil, errors.New("getting monitors: " + err.Error())
	}

	client, err := monitorhlp.GetClient(tx)
	if err != nil {
		return nil, errors.New("getting monitor client: " + err.Error())
	}

	cacheData, err := getCacheData(tx)
	if err != nil {
		return nil, errors.New("getting cache data: " + err.Error())
	}

	for cdn, monitorFQDNs := range monitors {
		if len(monitorFQDNs) == 0 {
			log.Warnln("getCachesStats: cdn '" + string(cdn) + "' has no online monitors, skipping!")
			continue
		}

		success := false
		errs := []error{}
		for _, monitorFQDN := range monitorFQDNs {
			crStates, err := monitorhlp.GetCRStates(monitorFQDN, client)
			if err != nil {
				errs = append(errs, errors.New("getting CRStates for CDN '"+string(cdn)+"' monitor '"+monitorFQDN+"': "+err.Error()))
				continue
			}

			cacheStats, err := getCacheStats(monitorFQDN, client)
			if err != nil {
				errs = append(errs, errors.New("getting CacheStats for CDN '"+string(cdn)+"' monitor '"+monitorFQDN+"': "+err.Error()))
				continue
			}

			cacheData = addHealth(cacheData, crStates)
			cacheData = addStats(cacheData, cacheStats)
			success = true
			break
		}

		if !success {
			return nil, errors.New("getting cache stats from all monitors failed for cdn '" + string(cdn) + "': " + util.JoinErrs(errs).Error())
		}

		// if we succeeded, log the monitor failures but don't return them
		for _, err := range errs {
			log.Errorln(err.Error())
		}
	}
	cacheData = addTotals(cacheData)
	return cacheData, nil
}

// addTotals sums each cachegroup, and adds the sum to an object with the cache-specific fields set to "ALL", and the cachegroup. It then sums all cachegroups, and adds the total to an object with all fields set to "ALL".
// TODO in the next API version, add totals in their own JSON objects, not amidst the cachegroup keys.
func addTotals(data []CacheData) []CacheData {
	all := "ALL"
	cachegroups := map[tc.CacheGroupName]CacheData{}
	total := CacheData{Profile: all, Status: all, Healthy: true, HostName: tc.CacheName(all), CacheGroup: tc.CacheGroupName(all)}
	for _, d := range data {
		cg := cachegroups[tc.CacheGroupName(d.CacheGroup)]
		cg.CacheGroup = d.CacheGroup
		cg.Connections += d.Connections
		cg.KBPS += d.KBPS
		cachegroups[tc.CacheGroupName(d.CacheGroup)] = cg
		total.Connections += d.Connections
		total.KBPS += d.KBPS
	}
	for _, cg := range cachegroups {
		cg.Profile = all
		cg.Status = all
		cg.Healthy = true
		cg.HostName = tc.CacheName(all)
		data = append(data, cg)
	}
	data = append(data, total)
	return data
}

// CRStates contains the Monitor CacheStats needed by Cachedata. It is NOT the full object served by the Monitor, but only the data required by the caches stats endpoint.
type CacheStats struct {
	Caches map[tc.CacheName]CacheStat `json:"caches"`
}

type CacheStat struct {
	BandwidthKBPS []CacheStatData `json:"bandwidth"`
	Connections   []CacheStatData `json:"ats.proxy.process.http.current_client_connections"`
}

type CacheStatData struct {
	Value int64 `json:"value,string"`
}

func getCacheStats(monitorFQDN string, client *http.Client) (CacheStats, error) {
	path := `/publish/CacheStats?stats=ats.proxy.process.http.current_client_connections,bandwidth`
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return CacheStats{}, errors.New("getting CacheStats from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()

	cs := CacheStats{}
	if err := json.NewDecoder(resp.Body).Decode(&cs); err != nil {
		return CacheStats{}, errors.New("decoding CacheStats from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return cs, nil
}

func addStats(cacheData []CacheData, stats CacheStats) []CacheData {
	if stats.Caches == nil {
		return cacheData // TODO warn?
	}
	for i, cache := range cacheData {
		stat, ok := stats.Caches[cache.HostName]
		if !ok {
			continue
		}
		if len(stat.BandwidthKBPS) > 0 {
			cache.KBPS = uint64(stat.BandwidthKBPS[0].Value)
		}
		if len(stat.Connections) > 0 {
			cache.Connections = uint64(stat.Connections[0].Value)
		}
		cacheData[i] = cache
	}
	return cacheData
}

func addHealth(cacheData []CacheData, crStates tc.CRStates) []CacheData {
	if crStates.Caches == nil {
		return cacheData // TODO warn?
	}
	for i, cache := range cacheData {
		crsCache, ok := crStates.Caches[cache.HostName]
		if !ok {
			continue
		}
		cache.Healthy = crsCache.IsAvailable
		cacheData[i] = cache
	}
	return cacheData
}

type CacheData struct {
	HostName    tc.CacheName      `json:"hostname"`
	CacheGroup  tc.CacheGroupName `json:"cachegroup"`
	Status      string            `json:"status"`
	Profile     string            `json:"profile"`
	IP          *string           `json:"ip"`
	Healthy     bool              `json:"healthy"`
	KBPS        uint64            `json:"kbps"`
	Connections uint64            `json:"connections"`
}

const CacheProfilePrefixEdge = "EDGE"
const CacheProfilePrefixMID = "MID"

// getCacheData gets the cache data from the servers table. Note this only gets from the database, and thus does not set the Healthy member.
func getCacheData(tx *sql.Tx) ([]CacheData, error) {
	qry := `
SELECT
  s.host_name,
  cg.name as cachegroup,
  st.name as status,
  p.name as profile,
  (select address from ip_address where s.id = ip_address.server and service_address = true AND family(address) = 4) as ip
FROM
  server s
  JOIN cachegroup cg ON s.cachegroup = cg.id
  JOIN status st ON s.status = st.id
  JOIN profile p ON s.profile = p.id
WHERE
  p.name LIKE '` + CacheProfilePrefixEdge + `%' OR p.name LIKE '` + CacheProfilePrefixMID + `%'
`
	rows, err := tx.Query(qry)
	if err != nil {
		return nil, errors.New("querying cache data: " + err.Error())
	}
	defer rows.Close()
	data := []CacheData{}
	for rows.Next() {
		d := CacheData{}
		if err := rows.Scan(&d.HostName, &d.CacheGroup, &d.Status, &d.Profile, &d.IP); err != nil {
			return nil, errors.New("scanning cache data: " + err.Error())
		}
		data = append(data, d)
	}
	return data, nil
}

// getCDNMonitors returns an FQDN, including port, of an online monitor for each CDN. If a CDN has no online monitors, that CDN will not have an entry in the map. If a CDN has multiple online monitors, an arbitrary one will be returned.
func getCDNMonitorFQDNs(tx *sql.Tx) (map[tc.CDNName][]string, error) {
	qry := `
SELECT
  s.host_name,
  s.domain_name,
  s.tcp_port,
  c.name as cdn
FROM
  server s
  JOIN type t ON s.type = t.id
  JOIN status st ON st.id = s.status
  JOIN cdn c ON c.id = s.cdn_id
WHERE
  t.name = '` + tc.MonitorTypeName + `'
  AND st.name = '` + MonitorOnlineStatus + `'
`
	rows, err := tx.Query(qry)
	if err != nil {
		return nil, errors.New("querying monitors: " + err.Error())
	}
	defer rows.Close()
	monitors := map[tc.CDNName][]string{}
	for rows.Next() {
		host := ""
		domain := ""
		port := sql.NullInt64{}
		cdn := tc.CDNName("")
		if err := rows.Scan(&host, &domain, &port, &cdn); err != nil {
			return nil, errors.New("scanning monitors: " + err.Error())
		}
		fqdn := host + "." + domain
		if port.Valid {
			fqdn += ":" + strconv.FormatInt(port.Int64, 10)
		}
		monitors[cdn] = append(monitors[cdn], fqdn)
	}
	return monitors, nil
}
