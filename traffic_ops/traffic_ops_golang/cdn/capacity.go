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
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/monitorhlp"
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

func getCapacity(tx *sql.Tx) (CapacityResp, error) {
	monitors, err := monitorhlp.GetURLs(tx)
	if err != nil {
		return CapacityResp{}, errors.New("getting monitors: " + err.Error())
	}
	if len(monitors) == 0 {
		return CapacityResp{}, errors.New("no monitors found")
	}
	return getMonitorsCapacity(tx, monitors)
}

type CapacityResp struct {
	AvailablePercent   float64 `json:"availablePercent"`
	UnavailablePercent float64 `json:"unavailablePercent"`
	UtilizedPercent    float64 `json:"utilizedPercent"`
	MaintenancePercent float64 `json:"maintenancePercent"`
}

type CapData struct {
	Available   float64
	Unavailable float64
	Utilized    float64
	Maintenance float64
	Capacity    float64
}

func getMonitorsCapacity(tx *sql.Tx, monitors map[tc.CDNName][]string) (CapacityResp, error) {
	client, err := monitorhlp.GetClient(tx)
	if err != nil {
		return CapacityResp{}, errors.New("getting TM client: " + err.Error())
	}

	thresholds, err := getEdgeProfileHealthThresholdBandwidth(tx)
	if err != nil {
		return CapacityResp{}, errors.New("getting profile thresholds: " + err.Error())
	}

	cap, err := getCapacityData(monitors, thresholds, client, tx)
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
func getCapacityData(monitors map[tc.CDNName][]string, thresholds map[string]float64, client *http.Client, tx *sql.Tx) (CapData, error) {
	cap := CapData{}
	for cdn, monitorFQDNs := range monitors {
		err := error(nil)
		for _, monitorFQDN := range monitorFQDNs {
			crStates := tc.CRStates{}
			crConfig := tc.CRConfig{}
			cacheStats := tc.Stats{}
			if crStates, err = monitorhlp.GetCRStates(monitorFQDN, client); err != nil {
				err = errors.New("getting CRStates for CDN '" + string(cdn) + "' monitor '" + monitorFQDN + "': " + err.Error())
				log.Warnln("getCapacity failed to get CRStates from cdn '" + string(cdn) + " monitor '" + monitorFQDN + "', trying next monitor: " + err.Error())
				continue
			}
			if crConfig, err = monitorhlp.GetCRConfig(monitorFQDN, client); err != nil {
				err = errors.New("getting CRConfig for CDN '" + string(cdn) + "' monitor '" + monitorFQDN + "': " + err.Error())
				log.Warnln("getCapacity failed to get CRConfig from cdn '" + string(cdn) + " monitor '" + monitorFQDN + "', trying next monitor: " + err.Error())
				continue
			}
			statsToFetch := []string{tc.StatNameKBPS, tc.StatNameMaxKBPS}
			var monitorEndpoint string
			if cacheStats, monitorEndpoint, err = monitorhlp.GetCacheStats(monitorFQDN, client, statsToFetch); err != nil {
				log.Warnln("getCapacity failed to get '" + monitorEndpoint + "' from cdn '" + string(cdn) + "', Error: " + err.Error() + ", trying CacheStats")
				legacyCacheStats, monitorEndpoint, err := monitorhlp.GetLegacyCacheStats(monitorFQDN, client, statsToFetch)
				if err != nil {
					log.Warnln("getCapacity failed to get '" + monitorEndpoint + "' from cdn '" + string(cdn) + "', Error: " + err.Error())
					continue
				}
				cacheStats = monitorhlp.UpgradeLegacyStats(legacyCacheStats)
			}

			cap = addCapacity(cap, cacheStats, crStates, crConfig, thresholds, tx)
			break
		}
		if err != nil {
			return CapData{}, err
		}
	}
	return cap, nil
}

func addCapacity(cap CapData, cacheStats tc.Stats, crStates tc.CRStates, crConfig tc.CRConfig, thresholds map[string]float64, tx *sql.Tx) CapData {
	for cacheName, stats := range cacheStats.Caches {
		cache, ok := crConfig.ContentServers[(cacheName)]
		if !ok {
			continue
		}
		if cache.ServerType == nil || cache.ServerStatus == nil || cache.Profile == nil {
			log.Warnln("addCapacity got cache with nil values! Skipping!")
			continue
		}
		if !strings.HasPrefix(*cache.ServerType, string(tc.CacheTypeEdge)) {
			continue
		}
		kbps, maxKbps, err := getStats(stats)
		if err != nil {
			log.Errorf("couldn't get stats for %v. err: %v", cacheName, err.Error())
			continue
		}
		if string(*cache.ServerStatus) == string(tc.CacheStatusReported) || string(*cache.ServerStatus) == string(tc.CacheStatusOnline) {
			if crStates.Caches[tc.CacheName(cacheName)].IsAvailable {
				cap.Available += kbps
			} else {
				cap.Unavailable += kbps
			}
		} else if string(*cache.ServerStatus) == string(tc.CacheStatusAdminDown) {
			cap.Maintenance += kbps
		} else {
			continue // don't add capacity for OFFLINE or other statuses
		}
		cap.Capacity += maxKbps - thresholds[*cache.Profile]
	}
	return cap
}

func getStats(stats tc.ServerStats) (float64, float64, error) {
	kbpsRaw, ok := stats.Stats[tc.StatNameKBPS]
	if !ok {
		return 0, 0, errors.New("no kbps stats")
	}
	maxKbpsRaw, ok := stats.Stats[tc.StatNameMaxKBPS]
	if !ok {
		return 0, 0, errors.New("no maxKbpsR stats")
	}
	if len(kbpsRaw) < 1 ||
		len(maxKbpsRaw) < 1 {
		return 0, 0, errors.New("no kbps/maxKbps stats to return")
	}
	kbps, ok := util.ToNumeric(kbpsRaw[0].Val)
	if !ok {
		return 0, 0, errors.New("unable to convert kbps to a float")
	}
	maxKbps, ok := util.ToNumeric(maxKbpsRaw[0].Val)
	if !ok {
		return 0, 0, errors.New("unable to convert maxKbps to a float")
	}
	return kbps, maxKbps, nil
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
