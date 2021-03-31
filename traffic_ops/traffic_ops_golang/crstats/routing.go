package crstats

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
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
)

type (
	RouterResp struct {
		Error error
		Stats tc.CRSStats
	}
	RouterData struct {
		StatTotal tc.CRSStatsStat
		Total     uint64
	}
)

const (
	RouterProxyParameter = "tm.traffic_rtr_fwd_proxy"
	RouterRequestTimeout = time.Second * 10
	RouterOnlineStatus   = "ONLINE"
	HTTP                 = "HTTP"
	DNS                  = "DNS"
)

func getRoutersRouting(tx *sql.Tx, routers map[tc.CDNName][]string, statType *string, hostRegexs []string) (tc.Routing, error) {
	forwardProxy, forwardProxyExists, err := dbhelpers.GetGlobalParam(tx, RouterProxyParameter)
	fmt.Println(forwardProxy, forwardProxyExists, err)
	if err != nil {
		return tc.Routing{}, errors.New("getting global router proxy parameter: " + err.Error())
	}
	client := &http.Client{Timeout: RouterRequestTimeout}
	if forwardProxyExists {
		fmt.Println("Doubt you are here?")
		proxyURI, err := url.Parse(forwardProxy)
		if err != nil {
			return tc.Routing{}, errors.New("router forward proxy '" + forwardProxy + "' in parameter '" + RouterProxyParameter + "' not a URI: " + err.Error())
		}
		clientTransport := &http.Transport{Proxy: http.ProxyURL(proxyURI)}
		// Disable HTTP/2. Go Transport Proxy does not support H2 Servers, and if the server does support it, the client will fail.
		// See https://github.com/golang/go/issues/26479 "We only support http1 proxies currently."
		clientTransport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
		client = &http.Client{Timeout: RouterRequestTimeout, Transport: clientTransport}
	}

	var hostRegex *regexp.Regexp
	fmt.Println("Host", hostRegex, "Client", client)
	if len(hostRegexs) > 0 {
		hostRegex, err = regexp.Compile(strings.Join(hostRegexs, "|"))
		if err != nil {
			return tc.Routing{}, errors.New("getting regex from host patterns: " + err.Error())
		}
	}

	count := 0
	for _, routerFQDNs := range routers {
		fmt.Println(routerFQDNs)
		count += len(routerFQDNs)
	}

	resp := make(chan RouterResp, count)

	wg := sync.WaitGroup{}
	fmt.Println(count)
	wg.Add(count)

	for cdn, routerFQDNs := range routers {
		fmt.Println(cdn, routerFQDNs)
		for _, routerFQDN := range routerFQDNs {
			go getCRSStats(resp, &wg, routerFQDN, string(cdn), client)
		}
	}

	wg.Wait()
	close(resp)

	dat := RouterData{}
	//testmap := map[string]tc.CRSStatsStat{
	//	"video.demo1.mycdn.ciab.test": tc.CRSStatsStat {10, 20, 0, 0, 0, 0, 0, 0, 0, 0},
	//}
	//test := tc.CRSStats{
	//	App:   tc.CRSStatsApp{
	//		BuildTimestamp: "2021-03-24",
	//		Name:           "traffic_router",
	//		DeployDir:      "/opt/traffic_router",
	//		GitRevision:    "36a47e9a3",
	//		Version:        "5.1.0",
	//	},
	//	Stats: tc.CRSStatsStats{
	//		DNSMap:           nil,
	//		HTTPMap:          testmap,
	//		TotalDNSCount:    0,
	//		TotalHTTPCount:   2,
	//		TotalDSMissCount: 0,
	//		AppStartTime:     1616606046447,
	//		AverageDnsTime:   0,
	//		AverageHttpTime:  1616606088088,
	//		UpdateTracker:    tc.CRSStatsUpdateTracker{
	//			LastHttpsCertificatesCheck:           0,
	//			LastGeolocationDatabaseUpdaterUpdate: 2,
	//			LastCacheStateCheck:                  0,
	//			LastCacheStateChange:                 1616606408121,
	//			LastNetworkUpdaterUpdate:             0,
	//			LastHTTPSCertificatesUpdate:          1616606543653,
	//			LastConfigCheck:                      1616606069501,
	//			LastConfigChange:                     0,
	//			LastHTTPSCertificatesFetchFail:       1616606410547,
	//			LastNetworkUpdaterCheck:              1616606544371,
	//			NewDNSSECKeysFound:                   1616606109280,
	//			LastGeolocationDatabaseUpdaterCheck:  0,
	//			LastHTTPSCertificatesFetchSuccess:    0,
	//			LastSteeringWatcherCheck:             0,
	//			LastDNSSECKeysCheck:                  0,
	//			LastFederationsWatcherCheck:          1616606408507,
	//			LastHTTPSCertificatesFetchAttempt:    1616606487260,
	//		},
	//	},
	//}

	fmt.Println("Resp", &resp)
	for r := range resp {
		fmt.Println("Data", dat, "r", r.Stats)
		//if r.Error != nil {
		//	return tc.Routing{}, err
		//}
		fmt.Println("You here?")
		//r.Stats = test
		dat = addCRSStats(dat, r.Stats, statType, hostRegex)
		fmt.Println("Data1", dat)
	}
	return sumRouterData(dat), nil
}

func sumRouterData(d RouterData) tc.Routing {
	if d.Total == 0 {
		return tc.Routing{}
	}
	return tc.Routing{
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

func addCRSStats(d RouterData, stats tc.CRSStats, statType *string, hostRegex *regexp.Regexp) RouterData {
	matchingHost := func(host string) bool {
		if hostRegex == nil {
			return true
		}
		return hostRegex.MatchString(host)
	}
	// DNSMap
	if statType == nil || *statType == "DNS" {
		for host, stat := range stats.Stats.DNSMap {
			fmt.Println("MH-dns", host, matchingHost(host))
			if matchingHost(host) {
				fmt.Println(d.StatTotal, stat)
				d.StatTotal = sumCRSStat(d.StatTotal, stat)
				d.Total += totalCRSStat(stat)
			}
		}
	}

	// HTTPMap
	if statType == nil || *statType == "HTTP" {
		for host, stat := range stats.Stats.HTTPMap {
			fmt.Println("MH-http", host, matchingHost(host))
			if matchingHost(host) {
				fmt.Println(d.StatTotal, stat)
				d.StatTotal = sumCRSStat(d.StatTotal, stat)
				d.Total += totalCRSStat(stat)
			}
		}
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

func getCRSStats(respond chan<- RouterResp, wg *sync.WaitGroup, routerFQDN, cdn string, client *http.Client) {
	defer wg.Done()
	r := RouterResp{}
	resp, err := client.Get("http://" + routerFQDN + "/crs/stats")
	if err != nil {
		r.Error = fmt.Errorf("getting crs stats for CDN %s router %s: %v", cdn, routerFQDN, err)
		respond <- r
		return
	}
	stats := tc.CRSStats{}
	fmt.Println("Response", resp)
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		r.Error = fmt.Errorf("decoding stats from CDN %s router %s: %v", cdn, routerFQDN, err)
		respond <- r
		return
	}
	r.Stats = stats
	fmt.Println(r.Stats)
	respond <- r
}

// getCDNRouterFQDNs returns an FQDN, including port, of an online router for each CDN, for each router. If a CDN has no online routers, that CDN will not have an entry in the map. The port returned is the API port.
func getCDNRouterFQDNs(tx *sql.Tx, requiredCDN *string) (map[tc.CDNName][]string, error) {
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
		fmt.Println(fqdn)
		//fqdn = "cdn-tr-atl-01.cdn.comcast.net:3333"
		fqdn = "cdn-tr-oshob-02.nightly.cdnlab.comcast.net:3333"
		fmt.Println(fqdn)
		if requiredCDN != nil && *requiredCDN != cdn {
			continue
		}
		routers[tc.CDNName(cdn)] = append(routers[tc.CDNName(cdn)], fqdn)
	}
	fmt.Println(routers)
	return routers, nil
}
