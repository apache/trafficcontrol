package deliveryservice

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
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/monitorhlp"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	dsID := inf.IntParams["id"]

	userErr, sysErr, errCode = tenant.CheckID(inf.Tx.Tx, inf.User, dsID)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	ds, cdn, ok, err := dbhelpers.GetDSNameAndCDNFromID(inf.Tx.Tx, dsID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting delivery service name from ID: "+err.Error()))
		return
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, nil)
		return
	}

	health, err := getHealth(inf.Tx.Tx, ds, cdn)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting delivery service health: "+err.Error()))
		return
	}

	api.WriteResp(w, r, health)
}

func getHealth(tx *sql.Tx, ds tc.DeliveryServiceName, cdn tc.CDNName) (tc.HealthData, error) {
	monitorURLs, err := monitorhlp.GetURLs(tx)
	if err != nil {
		return tc.HealthData{}, errors.New("getting monitors: " + err.Error())
	}
	monitors, ok := monitorURLs[cdn]
	if !ok {
		return tc.HealthData{}, nil // TODO emulates old Perl behavior; change to return error?
	}
	return getMonitorHealth(tx, ds, monitors)
}

func getMonitorHealth(tx *sql.Tx, ds tc.DeliveryServiceName, monitorFQDNs []string) (tc.HealthData, error) {
	client, err := monitorhlp.GetClient(tx)
	if err != nil {
		return tc.HealthData{}, errors.New("getting monitor client: " + err.Error())
	}

	totalOnline := uint64(0)
	totalOffline := uint64(0)
	cgData := map[tc.CacheGroupName]tc.HealthDataCacheGroup{}

	errs := []error{}
	for _, monitorFQDN := range monitorFQDNs {
		crStates, err := monitorhlp.GetCRStates(monitorFQDN, client)
		if err != nil {
			errs = append(errs, errors.New("getting CRStates for delivery service '"+string(ds)+"' monitor '"+monitorFQDN+"': "+err.Error()))
			continue
		}
		crConfig, err := monitorhlp.GetCRConfig(monitorFQDN, client)
		if err != nil {
			errs = append(errs, errors.New("getting CRConfig for delivery service '"+string(ds)+"' monitor '"+monitorFQDN+"': "+err.Error()))
			continue
		}
		cgData, totalOnline, totalOffline, err = addHealth(ds, cgData, totalOnline, totalOffline, crStates, crConfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		healthData := tc.HealthData{TotalOffline: totalOffline, TotalOnline: totalOnline, CacheGroups: []tc.HealthDataCacheGroup{}}
		for _, health := range cgData {
			healthData.CacheGroups = append(healthData.CacheGroups, health)
		}
		return healthData, nil
	}
	return tc.HealthData{}, errors.New("getting monitor health: " + util.JoinErrs(errs).Error())
}

// addHealth adds the given cache states to the given data and totals, and returns the new data and totals
func addHealth(ds tc.DeliveryServiceName, data map[tc.CacheGroupName]tc.HealthDataCacheGroup, totalOnline uint64, totalOffline uint64, crStates tc.CRStates, crConfig tc.CRConfig) (map[tc.CacheGroupName]tc.HealthDataCacheGroup, uint64, uint64, error) {

	var deliveryService tc.CRConfigDeliveryService
	var ok bool
	var topology string
	var cacheGroupNameMap = make(map[string]bool)

	if deliveryService, ok = crConfig.DeliveryServices[string(ds)]; !ok {
		return map[tc.CacheGroupName]tc.HealthDataCacheGroup{}, 0, 0, errors.New("delivery service not found in CRConfig")
	}
	if deliveryService.Topology != nil {
		var top tc.CRConfigTopology
		topology = *deliveryService.Topology
		if topology != "" {
			if top, ok = crConfig.Topologies[topology]; !ok {
				return map[tc.CacheGroupName]tc.HealthDataCacheGroup{}, 0, 0, fmt.Errorf("CRConfig topologies does not contain DS topology: %s", topology)
			}
			for _, n := range top.Nodes {
				cacheGroupNameMap[n] = true
			}
		}
	}
	for cacheName, avail := range crStates.Caches {
		var skip bool
		cache, ok := crConfig.ContentServers[string(cacheName)]
		if !ok {
			continue // TODO warn?
		}
		if topology == "" {
			if _, ok := cache.DeliveryServices[string(ds)]; !ok {
				continue
			}
		}
		if cache.ServerStatus == nil || *cache.ServerStatus != tc.CRConfigServerStatus(tc.CacheStatusReported) {
			continue
		}
		if cache.ServerType == nil || !strings.HasPrefix(string(*cache.ServerType), string(tc.CacheTypeEdge)) {
			continue
		}
		if cache.CacheGroup == nil {
			continue // TODO warn?
		}
		if topology != "" {
			if _, ok := cacheGroupNameMap[*cache.CacheGroup]; !ok {
				continue
			}
			cacheCapabilities := make(map[string]struct{}, len(cache.Capabilities))
			for _, cap := range cache.Capabilities {
				cacheCapabilities[cap] = struct{}{}
			}
			for _, rc := range deliveryService.RequiredCapabilities {
				if _, ok = cacheCapabilities[rc]; !ok {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
		}
		cgHealth := data[tc.CacheGroupName(*cache.CacheGroup)]
		cgHealth.Name = tc.CacheGroupName(*cache.CacheGroup)
		if avail.IsAvailable {
			cgHealth.Online++
			totalOnline++
		} else {
			cgHealth.Offline++
			totalOffline++
		}
		data[tc.CacheGroupName(*cache.CacheGroup)] = cgHealth
	}
	return data, totalOnline, totalOffline, nil
}
