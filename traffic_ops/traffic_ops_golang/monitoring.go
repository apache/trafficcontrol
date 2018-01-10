package main

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
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

const MonitoringPrivLevel = auth.PrivLevelReadOnly

const CacheMonitorConfigFile = "rascal.properties"
const MonitorType = "RASCAL"
const RouterType = "CCR"
const MonitorProfilePrefix = "RASCAL"
const MonitorConfigFile = "rascal-config.txt"
const KilobitsPerMegabit = 1000
const DeliveryServiceStatus = "REPORTED"

type BasicServer struct {
	Profile    string `json:"profile"`
	Status     string `json:"status"`
	IP         string `json:"ip"`
	IP6        string `json:"ip6"`
	Port       int    `json:"port"`
	Cachegroup string `json:"cachegroup"`
	HostName   string `json:"hostname"`
	FQDN       string `json:"fqdn"`
}

type Monitor struct {
	BasicServer
}

type Cache struct {
	BasicServer
	InterfaceName string `json:"interfacename"`
	Type          string `json:"type"`
	HashID        string `json:"hashid"`
}

type Cachegroup struct {
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Profile struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

type Monitoring struct {
	TrafficServers   []Cache                `json:"trafficServers"`
	TrafficMonitors  []Monitor              `json:"trafficMonitors"`
	Cachegroups      []Cachegroup           `json:"cacheGroups"`
	Profiles         []Profile              `json:"profiles"`
	DeliveryServices []DeliveryService      `json:"deliveryServices"`
	Config           map[string]interface{} `json:"config"`
}

type MonitoringResponse struct {
	Response Monitoring `json:"response"`
}

type Router struct {
	Type    string
	Profile string
}

type DeliveryService struct {
	XMLID              string  `json:"xmlId"`
	TotalTPSThreshold  float64 `json:"totalTpsThreshold"`
	Status             string  `json:"status"`
	TotalKBPSThreshold float64 `json:"totalKbpsThreshold"`
}

// TODO change to use the PathParams, instead of parsing the URL
func monitoringHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		cdnName := pathParams["name"]

		resp, err := getMonitoringJSON(cdnName, db)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getMonitoringServers(db *sqlx.DB, cdn string) ([]Monitor, []Cache, []Router, error) {
	query := `SELECT
me.host_name as hostName,
CONCAT(me.host_name, '.', me.domain_name) as fqdn,
status.name as status,
cachegroup.name as cachegroup,
me.tcp_port as port,
me.ip_address as ip,
me.ip6_address as ip6,
profile.name as profile,
me.interface_name as interfaceName,
type.name as type,
me.xmpp_id as hashID
FROM server me
JOIN type type ON type.id = me.type
JOIN status status ON status.id = me.status
JOIN cachegroup cachegroup ON cachegroup.id = me.cachegroup
JOIN profile profile ON profile.id = me.profile
JOIN cdn cdn ON cdn.id = me.cdn_id
WHERE cdn.name = $1`

	rows, err := db.Query(query, cdn)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	monitors := []Monitor{}
	caches := []Cache{}
	routers := []Router{}

	for rows.Next() {
		var hostName sql.NullString
		var fqdn sql.NullString
		var status sql.NullString
		var cachegroup sql.NullString
		var port sql.NullInt64
		var ip sql.NullString
		var ip6 sql.NullString
		var profile sql.NullString
		var interfaceName sql.NullString
		var ttype sql.NullString
		var hashID sql.NullString

		if err := rows.Scan(&hostName, &fqdn, &status, &cachegroup, &port, &ip, &ip6, &profile, &interfaceName, &ttype, &hashID); err != nil {
			return nil, nil, nil, err
		}

		if ttype.String == MonitorType {
			monitors = append(monitors, Monitor{
				BasicServer: BasicServer{
					Profile:    profile.String,
					Status:     status.String,
					IP:         ip.String,
					IP6:        ip6.String,
					Port:       int(port.Int64),
					Cachegroup: cachegroup.String,
					HostName:   hostName.String,
					FQDN:       fqdn.String,
				},
			})
		} else if strings.HasPrefix(ttype.String, "EDGE") || strings.HasPrefix(ttype.String, "MID") {
			caches = append(caches, Cache{
				BasicServer: BasicServer{
					Profile:    profile.String,
					Status:     status.String,
					IP:         ip.String,
					IP6:        ip6.String,
					Port:       int(port.Int64),
					Cachegroup: cachegroup.String,
					HostName:   hostName.String,
					FQDN:       fqdn.String,
				},
				InterfaceName: interfaceName.String,
				Type:          ttype.String,
				HashID:        hashID.String,
			})
		} else if ttype.String == RouterType {
			routers = append(routers, Router{
				Type:    ttype.String,
				Profile: profile.String,
			})
		}
	}
	return monitors, caches, routers, nil
}

func getCachegroups(db *sqlx.DB, cdn string) ([]Cachegroup, error) {
	query := `
SELECT name, latitude, longitude
FROM cachegroup
WHERE id IN
  (SELECT cachegroup FROM server WHERE server.cdn_id =
    (SELECT id FROM cdn WHERE name = $1));`

	rows, err := db.Query(query, cdn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cachegroups := []Cachegroup{}
	for rows.Next() {
		var name sql.NullString
		var lat sql.NullFloat64
		var lon sql.NullFloat64
		if err := rows.Scan(&name, &lat, &lon); err != nil {
			return nil, err
		}
		cachegroups = append(cachegroups, Cachegroup{
			Name: name.String,
			Coordinates: Coordinates{
				Latitude:  lat.Float64,
				Longitude: lon.Float64,
			},
		})
	}
	return cachegroups, nil
}

func getProfiles(db *sqlx.DB, caches []Cache, routers []Router) ([]Profile, error) {
	cacheProfileTypes := map[string]string{}
	profiles := map[string]Profile{}
	profileNames := []string{}
	for _, router := range routers {
		profiles[router.Profile] = Profile{
			Name: router.Profile,
			Type: router.Type,
		}
	}

	for _, cache := range caches {
		if _, ok := cacheProfileTypes[cache.Profile]; !ok {
			cacheProfileTypes[cache.Profile] = cache.Type
			profiles[cache.Profile] = Profile{
				Name: cache.Profile,
				Type: cache.Type,
			}
			profileNames = append(profileNames, cache.Profile)
		}
	}

	query := `
SELECT p.name as profile, pr.name, pr.value
FROM parameter pr
JOIN profile p ON p.name = ANY($1)
JOIN profile_parameter pp ON pp.profile = p.id and pp.parameter = pr.id
WHERE pr.config_file = $2;
`
	rows, err := db.Query(query, pq.Array(profileNames), CacheMonitorConfigFile)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var profileName sql.NullString
		var name sql.NullString
		var value sql.NullString
		if err := rows.Scan(&profileName, &name, &value); err != nil {
			return nil, err
		}
		if name.String == "" {
			return nil, fmt.Errorf("null name") // TODO continue and warn?
		}
		profile := profiles[profileName.String]
		if profile.Parameters == nil {
			profile.Parameters = map[string]interface{}{}
		}

		if valNum, err := strconv.Atoi(value.String); err == nil {
			profile.Parameters[name.String] = valNum
		} else {
			profile.Parameters[name.String] = value.String
		}
		profiles[profileName.String] = profile

	}

	profilesArr := []Profile{} // TODO make for efficiency?
	for _, profile := range profiles {
		profilesArr = append(profilesArr, profile)
	}
	return profilesArr, nil
}

func getDeliveryServices(db *sqlx.DB, routers []Router) ([]DeliveryService, error) {
	profileNames := []string{}
	for _, router := range routers {
		profileNames = append(profileNames, router.Profile)
	}

	query := `
SELECT ds.xml_id, ds.global_max_tps, ds.global_max_mbps
FROM deliveryservice ds
JOIN profile profile ON profile.id = ds.profile
WHERE profile.name = ANY($1)
AND ds.active = true
`
	rows, err := db.Query(query, pq.Array(profileNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dses := []DeliveryService{}

	for rows.Next() {
		var xmlid sql.NullString
		var tps sql.NullFloat64
		var mbps sql.NullFloat64
		if err := rows.Scan(&xmlid, &tps, &mbps); err != nil {
			return nil, err
		}
		dses = append(dses, DeliveryService{
			XMLID:              xmlid.String,
			TotalTPSThreshold:  tps.Float64,
			Status:             DeliveryServiceStatus,
			TotalKBPSThreshold: mbps.Float64 * KilobitsPerMegabit,
		})
	}
	return dses, nil
}

func getConfig(db *sqlx.DB) (map[string]interface{}, error) {
	// TODO remove 'like' in query? Slow?
	query := fmt.Sprintf(`
SELECT pr.name, pr.value
FROM parameter pr
JOIN profile p ON p.name LIKE '%s%%'
JOIN profile_parameter pp ON pp.profile = p.id and pp.parameter = pr.id
WHERE pr.config_file = '%s'
`, MonitorProfilePrefix, MonitorConfigFile)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cfg := map[string]interface{}{}

	for rows.Next() {
		var name sql.NullString
		var val sql.NullString
		if err := rows.Scan(&name, &val); err != nil {
			return nil, err
		}
		if valNum, err := strconv.Atoi(val.String); err == nil {
			cfg[name.String] = valNum
		} else {
			cfg[name.String] = val.String
		}
	}
	return cfg, nil
}

func getMonitoringJSON(cdnName string, db *sqlx.DB) (*MonitoringResponse, error) {
	monitors, caches, routers, err := getMonitoringServers(db, cdnName)
	if err != nil {
		return nil, fmt.Errorf("error getting servers: %v", err)
	}

	cachegroups, err := getCachegroups(db, cdnName)
	if err != nil {
		return nil, fmt.Errorf("error getting cachegroups: %v", err)
	}

	profiles, err := getProfiles(db, caches, routers)
	if err != nil {
		return nil, fmt.Errorf("error getting profiles: %v", err)
	}

	deliveryServices, err := getDeliveryServices(db, routers)
	if err != nil {
		return nil, fmt.Errorf("error getting deliveryservices: %v", err)
	}

	config, err := getConfig(db)
	if err != nil {
		return nil, fmt.Errorf("error getting config: %v", err)
	}

	resp := MonitoringResponse{
		Response: Monitoring{
			TrafficServers:   caches,
			TrafficMonitors:  monitors,
			Cachegroups:      cachegroups,
			Profiles:         profiles,
			DeliveryServices: deliveryServices,
			Config:           config,
		},
	}
	return &resp, nil
}
