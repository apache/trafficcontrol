package monitorhlp

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
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

const MonitorProxyParameter = "tm.traffic_mon_fwd_proxy"
const MonitorRequestTimeout = time.Second * 10
const MonitorOnlineStatus = "ONLINE"

// GetClient returns the http.Client for making requests to the Traffic Monitor. This should always be used, rather than creating a default http.Client, to ensure any monitor forward proxy parameter is used correctly.
func GetClient(tx *sql.Tx) (*http.Client, error) {
	monitorForwardProxy, monitorForwardProxyExists, err := dbhelpers.GetGlobalParam(tx, MonitorProxyParameter)
	if err != nil {
		return nil, errors.New("getting global monitor proxy parameter: " + err.Error())
	}
	client := &http.Client{Timeout: MonitorRequestTimeout}
	if monitorForwardProxyExists {
		proxyURI, err := url.Parse(monitorForwardProxy)
		if err != nil {
			return nil, errors.New("monitor forward proxy '" + monitorForwardProxy + "' in parameter '" + MonitorProxyParameter + "' not a URI: " + err.Error())
		}

		clientTransport := &http.Transport{Proxy: http.ProxyURL(proxyURI)}
		// Disable HTTP/2. Go Transport Proxy does not support H2 Servers, and if the server does support it, the client will fail.
		// See https://github.com/golang/go/issues/26479 "We only support http1 proxies currently."
		clientTransport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
		client = &http.Client{Timeout: MonitorRequestTimeout, Transport: clientTransport}
	}
	return client, nil
}

// GetURLs returns an FQDN, including port, of an online monitor for each CDN. If a CDN has no online monitors, that CDN will not have an entry in the map. If a CDN has multiple online monitors, an arbitrary one will be returned.
func GetURLs(tx *sql.Tx) (map[tc.CDNName]string, error) {
	rows, err := tx.Query(`
SELECT s.host_name, s.domain_name, s.tcp_port, c.name as cdn
FROM server as s
JOIN type as t ON s.type = t.id
JOIN status as st ON st.id = s.status
JOIN cdn as c ON c.id = s.cdn_id
WHERE t.name = '` + tc.MonitorTypeName + `'
AND st.name = '` + MonitorOnlineStatus + `'
`)
	if err != nil {
		return nil, errors.New("querying monitors: " + err.Error())
	}
	defer rows.Close()
	monitors := map[tc.CDNName]string{}
	for rows.Next() {
		host := ""
		domain := ""
		port := sql.NullInt64{}
		cdn := ""
		if err := rows.Scan(&host, &domain, &port, &cdn); err != nil {
			return nil, errors.New("scanning monitors: " + err.Error())
		}
		fqdn := host + "." + domain
		if port.Valid {
			fqdn += ":" + strconv.FormatInt(port.Int64, 10)
		}
		monitors[tc.CDNName(cdn)] = fqdn
	}
	return monitors, nil
}

func GetCRStates(monitorFQDN string, client *http.Client) (tc.CRStates, error) {
	path := `/publish/CrStates`
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return tc.CRStates{}, errors.New("getting CRStates from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()

	crs := tc.CRStates{}
	if err := json.NewDecoder(resp.Body).Decode(&crs); err != nil {
		return tc.CRStates{}, errors.New("decoding CRStates from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return crs, nil
}

func GetCRConfig(monitorFQDN string, client *http.Client) (tc.CRConfig, error) {
	path := `/publish/CrConfig`
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return tc.CRConfig{}, errors.New("getting CRConfig from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()
	crs := tc.CRConfig{}
	if err := json.NewDecoder(resp.Body).Decode(&crs); err != nil {
		return tc.CRConfig{}, errors.New("decoding CRConfig from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return crs, nil
}

// GetCacheStats gets the cache stats from the given monitor. The stats parameters is which stats to get;
// if stats is empty or nil, all stats are fetched.
func GetCacheStats(monitorFQDN string, client *http.Client, stats []string) (tc.Stats, error) {
	path := `/publish/CacheStatsNew?hc=1`
	if len(stats) > 0 {
		path += `&stats=` + strings.Join(stats, `,`)
	}
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return tc.Stats{}, errors.New("getting CacheStatsNew from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()
	cacheStats := tc.Stats{}
	if err := json.NewDecoder(resp.Body).Decode(&cacheStats); err != nil {
		return tc.Stats{}, errors.New("decoding CacheStatsNew from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return cacheStats, nil
}

// GetLegacyCacheStats gets the pre ATCv5.0 cache stats from the given monitor. The stats parameters is which stats to
// get; if stats is empty or nil, all stats are fetched.
func GetLegacyCacheStats(monitorFQDN string, client *http.Client, stats []string) (tc.LegacyStats, error) {
	path := `/publish/CacheStats?hc=1`
	if len(stats) > 0 {
		path += `&stats=` + strings.Join(stats, `,`)
	}
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return tc.LegacyStats{}, errors.New("getting CacheStats from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()
	cacheStats := tc.LegacyStats{}
	if err := json.NewDecoder(resp.Body).Decode(&cacheStats); err != nil {
		return tc.LegacyStats{}, errors.New("decoding CacheStats from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return cacheStats, nil
}

// UpgradeLegacyStats will take LegacyStats and transform them to Stats. It assumes all stats that go in
// Stats.Caches[cacheName] exist in Stats and not Interfaces
func UpgradeLegacyStats(legacyStats tc.LegacyStats) tc.Stats {
	stats := tc.Stats{
		CommonAPIData: legacyStats.CommonAPIData,
		Caches:        make(map[string]tc.ServerStats, len(legacyStats.Caches)),
	}

	for cacheName, cache := range legacyStats.Caches {
		stats.Caches[string(cacheName)] = tc.ServerStats{
			Interfaces: nil,
			Stats:      make(map[string][]tc.ResultStatVal, len(cache)),
		}
		for statName, stat := range cache {
			stats.Caches[string(cacheName)].Stats[statName] = stat
		}
	}

	return stats
}
