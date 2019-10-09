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
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetCapacity(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	api.RespWriter(w, r, inf.Tx.Tx)(getCapacity(inf.Tx.Tx))
}

const MonitorProxyParameter = "tm.traffic_mon_fwd_proxy"
const MonitorRequestTimeout = time.Second * 10
const MonitorOnlineStatus = "ONLINE"

// CRStates contains the Monitor CRStates members needed for health. It is NOT the full object served by the Monitor, but only the data required by this endpoint.
type CRStates struct {
	Caches map[tc.CacheName]Available `json:"caches"`
}

type Available struct {
	IsAvailable bool `json:"isAvailable"`
}

// CRConfig contains the Monitor CRConfig members needed for health. It is NOT the full object served by the Monitor, but only the data required by this endpoint.
type CRConfig struct {
	ContentServers map[tc.CacheName]CRConfigServer `json:"contentServers"`
}

type CRConfigServer struct {
	CacheGroup tc.CacheGroupName `json:"locationId"`
	Status     tc.CacheStatus    `json:"status"`
	Type       tc.CacheType      `json:"type"`
	Profile    string            `json:"profile"`
}

func getCapacity(tx *sql.Tx) (CapacityResp, error) {
	monitors, err := getCDNMonitorFQDNs(tx)
	if err != nil {
		return CapacityResp{}, errors.New("getting monitors: " + err.Error())
	}

	return getMonitorsCapacity(tx, monitors)
}

type CapacityResp struct {
	AvailablePercent   float64 `json:"availablePercent"`
	UnavailablePercent float64 `json:"unavailablePercent"`
	UtilizedPercent    float64 `json:utilizedPercent"`
	MaintenancePercent float64 `json:maintenancePercent"`
}

type CapData struct {
	Available   float64
	Unavailable float64
	Utilized    float64
	Maintenance float64
	Capacity    float64
}

func getMonitorsCapacity(tx *sql.Tx, monitors map[tc.CDNName][]string) (CapacityResp, error) {
	monitorForwardProxy, monitorForwardProxyExists, err := getGlobalParam(tx, MonitorProxyParameter)
	if err != nil {
		return CapacityResp{}, errors.New("getting global monitor proxy parameter: " + err.Error())
	}
	client := &http.Client{Timeout: MonitorRequestTimeout}
	if monitorForwardProxyExists {
		proxyURI, err := url.Parse(monitorForwardProxy)
		if err != nil {
			return CapacityResp{}, errors.New("monitor forward proxy '" + monitorForwardProxy + "' in parameter '" + MonitorProxyParameter + "' not a URI: " + err.Error())
		}
		clientTransport := &http.Transport{Proxy: http.ProxyURL(proxyURI)}
		if proxyURI.Scheme == "https" {
			// TM does not support HTTP/2 and golang when connecting to https will use HTTP/2 by default causing a conflict
			// The result will be an unsupported scheme error
			// Setting TLSNextProto to any empty map will disable using HTTP/2 per https://golang.org/src/net/http/doc.go
			clientTransport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
		}
		client = &http.Client{Timeout: MonitorRequestTimeout, Transport: clientTransport}
	}

	thresholds, err := getEdgeProfileHealthThresholdBandwidth(tx)
	if err != nil {
		return CapacityResp{}, errors.New("getting profile thresholds: " + err.Error())
	}

	cap, err := getCapacityData(monitors, thresholds, client)
	if err != nil {
		return CapacityResp{}, errors.New("getting capacity from monitors: " + err.Error())
	} else if cap.Capacity == 0 {
		return CapacityResp{}, errors.New("capacity was zero!") // avoid divide-by-zero below.
	}

	return CapacityResp{
		UtilizedPercent:    (cap.Available * 100) / cap.Capacity,
		UnavailablePercent: (cap.Unavailable * 100) / cap.Capacity,
		MaintenancePercent: (cap.Maintenance * 100) / cap.Capacity,
		AvailablePercent:   ((cap.Capacity - cap.Unavailable - cap.Maintenance - cap.Available) * 100) / cap.Capacity,
	}, nil
}

// getCapacityData attempts to get the CDN capacity from each monitor. If one fails, it tries the next.
// The first monitor for which all data requests succeed is used.
// Only if all monitors for a CDN fail is an error returned, from the last monitor tried.
func getCapacityData(monitors map[tc.CDNName][]string, thresholds map[string]float64, client *http.Client) (CapData, error) {
	cap := CapData{}
	for cdn, monitorFQDNs := range monitors {
		err := error(nil)
		for _, monitorFQDN := range monitorFQDNs {
			crStates := CRStates{}
			crConfig := CRConfig{}
			cacheStats := CacheStats{}
			if crStates, err = getCRStates(monitorFQDN, client); err != nil {
				err = errors.New("getting CRStates for CDN '" + string(cdn) + "' monitor '" + monitorFQDN + "': " + err.Error())
				log.Warnln("getCapacity failed to get CRStates from cdn '" + string(cdn) + " monitor '" + monitorFQDN + "', trying next monitor: " + err.Error())
				continue
			}
			if crConfig, err = getCRConfig(monitorFQDN, client); err != nil {
				err = errors.New("getting CRConfig for CDN '" + string(cdn) + "' monitor '" + monitorFQDN + "': " + err.Error())
				log.Warnln("getCapacity failed to get CRConfig from cdn '" + string(cdn) + " monitor '" + monitorFQDN + "', trying next monitor: " + err.Error())
				continue
			}
			if err := getCacheStats(monitorFQDN, client, []string{"kbps", "maxKbps"}, &cacheStats); err != nil {
				err = errors.New("getting cache stats for CDN '" + string(cdn) + "' monitor '" + monitorFQDN + "': " + err.Error())
				log.Warnln("getCapacity failed to get CacheStats from cdn '" + string(cdn) + " monitor '" + monitorFQDN + "', trying next monitor: " + err.Error())
				continue
			}
			cap = addCapacity(cap, cacheStats, crStates, crConfig, thresholds)
			break
		}
		if err != nil {
			return CapData{}, err
		}
	}
	return cap, nil
}

func addCapacity(cap CapData, cacheStats CacheStats, crStates CRStates, crConfig CRConfig, thresholds map[string]float64) CapData {
	for cacheName, stats := range cacheStats.Caches {
		cache, ok := crConfig.ContentServers[cacheName]
		if !ok {
			continue
		}
		if !strings.HasPrefix(string(cache.Type), string(tc.CacheTypeEdge)) {
			continue
		}
		if len(stats.KBPS) < 1 || len(stats.MaxKBPS) < 1 {
			continue
		}
		if cache.Status == "REPORTED" || cache.Status == "ONLINE" {
			if crStates.Caches[cacheName].IsAvailable {
				cap.Available += float64(stats.KBPS[0].Value)
			} else {
				cap.Unavailable += float64(stats.KBPS[0].Value)
			}
		} else if cache.Status == "ADMIN_DOWN" {
			cap.Maintenance += float64(stats.KBPS[0].Value)
		} else {
			continue // don't add capacity for OFFLINE or other statuses
		}
		cap.Capacity += float64(stats.MaxKBPS[0].Value) - thresholds[cache.Profile]
	}
	return cap
}

func getEdgeProfileHealthThresholdBandwidth(tx *sql.Tx) (map[string]float64, error) {
	rows, err := tx.Query(`
SELECT pr.name as profile, pa.name, pa.config_file, pa.value
FROM parameter as pa
JOIN profile_parameter as pp ON pp.parameter = pa.id
JOIN profile as pr ON pp.profile = pr.id
JOIN server as s ON s.profile = pr.id
JOIN cdn as c ON c.id = s.cdn_id
JOIN type as t ON s.type = t.id
WHERE t.name LIKE 'EDGE%'
AND pa.config_file = 'rascal-config.txt'
AND pa.name = 'health.threshold.availableBandwidthInKbps'
`)
	if err != nil {
		return nil, errors.New("querying thresholds: " + err.Error())
	}
	defer rows.Close()
	profileThresholds := map[string]float64{}
	for rows.Next() {
		profile := ""
		threshStr := ""
		if err := rows.Scan(&profile, &threshStr); err != nil {
			return nil, errors.New("scanning thresholds: " + err.Error())
		}
		threshStr = strings.TrimPrefix(threshStr, ">")
		thresh, err := strconv.ParseFloat(threshStr, 64)
		if err != nil {
			return nil, errors.New("profile '" + profile + "' health.threshold.availableBandwidthInKbps is not a number")
		}
		profileThresholds[profile] = thresh
	}
	return profileThresholds, nil
}

func getCRStates(monitorFQDN string, client *http.Client) (CRStates, error) {
	path := `/publish/CrStates`
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return CRStates{}, errors.New("getting CRStates from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()

	crs := CRStates{}
	if err := json.NewDecoder(resp.Body).Decode(&crs); err != nil {
		return CRStates{}, errors.New("decoding CRStates from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return crs, nil
}

func getCRConfig(monitorFQDN string, client *http.Client) (CRConfig, error) {
	path := `/publish/CrConfig`
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return CRConfig{}, errors.New("getting CRConfig from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()
	crs := CRConfig{}
	if err := json.NewDecoder(resp.Body).Decode(&crs); err != nil {
		return CRConfig{}, errors.New("decoding CRConfig from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return crs, nil
}

// CacheStats contains the Monitor CacheStats needed by Cachedata. It is NOT the full object served by the Monitor, but only the data required by the caches stats endpoint.
type CacheStats struct {
	Caches map[tc.CacheName]CacheStat `json:"caches"`
}

type CacheStat struct {
	KBPS    []CacheStatData `json:"kbps"`
	MaxKBPS []CacheStatData `json:"maxKbps"`
}

type CacheStatData struct {
	Value float64 `json:"value,string"`
}

// getCacheStats gets the cache stats from the given monitor. It takes stats, a slice of stat names; and cacheStats, an object to deserialize stats into. The cacheStats type must be of the form struct {caches map[tc.CacheName]struct{statName []struct{value float64}}} with the desired stats, with appropriate member names or tags.
func getCacheStats(monitorFQDN string, client *http.Client, stats []string, cacheStats interface{}) error {
	path := `/publish/CacheStats`
	if len(stats) > 0 {
		path += `?stats=` + strings.Join(stats, `,`)
	}
	resp, err := client.Get("http://" + monitorFQDN + path)
	if err != nil {
		return errors.New("getting CacheStats from Monitor '" + monitorFQDN + "': " + err.Error())
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(cacheStats); err != nil {
		return errors.New("decoding CacheStats from monitor '" + monitorFQDN + "': " + err.Error())
	}
	return nil
}

// getCDNMonitors returns an FQDN, including port, of an online monitor for each CDN. If a CDN has no online monitors, that CDN will not have an entry in the map. If a CDN has multiple online monitors, an arbitrary one will be returned.
func getCDNMonitorFQDNs(tx *sql.Tx) (map[tc.CDNName][]string, error) {
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
