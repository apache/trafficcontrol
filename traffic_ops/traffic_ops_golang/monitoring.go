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

	"github.com/lib/pq"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

const MonitoringPrivLevel = PrivLevelReadOnly

const CacheMonitorConfigFile = "rascal.properties"
const MonitorType = "RASCAL"
const RouterType = "CCR"
const MonitorProfilePrefix = "RASCAL"
const MonitorConfigFile = "rascal-config.txt"
const KilobitsPerMegabit = 1000
const DeliveryServiceStatus = "REPORTED"

type MonitoringData struct {
	Servers          *sql.Stmt
	Cachegroups      *sql.Stmt
	Profiles         *sql.Stmt
	DeliveryServices *sql.Stmt
	Config           *sql.Stmt
}

const monitoringServersQuery = `
SELECT
me.host_name as hostName,
CONCAT(me.host_name, '.', me.domain_name) as fqdn,
status.name as status,
cachegroup.name as cachegroup,
COALESCE(me.tcp_port, 0) as port,
me.ip_address as ip,
COALESCE(me.ip6_address, '') as ip6,
profile.name as profile,
me.interface_name as interfaceName,
type.name as type,
COALESCE(me.xmpp_id, '') as hashId
FROM server me
JOIN type type ON type.id = me.type
JOIN status status ON status.id = me.status
JOIN cachegroup cachegroup ON cachegroup.id = me.cachegroup
JOIN profile profile ON profile.id = me.profile
JOIN cdn cdn ON cdn.id = me.cdn_id
WHERE cdn.name = $1
`
const monitoringCachegroupsQuery = `
SELECT name, COALESCE(latitude, 0), COALESCE(longitude, 0)
FROM cachegroup
WHERE id IN
  (SELECT cachegroup FROM server WHERE server.cdn_id =
    (SELECT id FROM cdn WHERE name = $1));
`
const monitoringProfilesQuery = `
SELECT p.name as profile, pr.name, pr.value
FROM parameter pr
JOIN profile p ON p.name = ANY($1)
JOIN profile_parameter pp ON pp.profile = p.id and pp.parameter = pr.id
WHERE pr.config_file = $2;
`
const monitoringDeliveryServicesQuery = `
SELECT ds.xml_id, COALESCE(ds.global_max_tps, 0), COALESCE(ds.global_max_mbps, 0)
FROM deliveryservice ds
JOIN profile profile ON profile.id = ds.profile
WHERE profile.name = ANY($1)
AND ds.active = true
`
const monitoringConfigQuery = `
SELECT pr.name, pr.value
FROM parameter pr
JOIN profile p ON p.name LIKE '` + MonitorProfilePrefix + `%%'
JOIN profile_parameter pp ON pp.profile = p.id and pp.parameter = pr.id
WHERE pr.config_file = '` + MonitorConfigFile + `'
` // TODO remove 'like' in query? Slow?

func monitoringData(db *sql.DB) (*MonitoringData, error) {
	d := MonitoringData{}
	err := error(nil)
	if d.Servers, err = db.Prepare(monitoringServersQuery); err != nil {
		return nil, err
	}
	if d.Cachegroups, err = db.Prepare(monitoringCachegroupsQuery); err != nil {
		return nil, err
	}
	if d.Profiles, err = db.Prepare(monitoringProfilesQuery); err != nil {
		return nil, err
	}
	if d.DeliveryServices, err = db.Prepare(monitoringDeliveryServicesQuery); err != nil {
		return nil, err
	}
	if d.Config, err = db.Prepare(monitoringConfigQuery); err != nil {
		return nil, err
	}
	return &d, nil
}

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

// TODO change to use the ParamMap, instead of parsing the URL
func monitoringHandler(d *MonitoringData) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p ParamMap) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		cdnName := p["cdn"]

		resp, err := getMonitoringJson(cdnName, d)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getMonitoringServers(stmt *sql.Stmt, cdn string) ([]Monitor, []Cache, []Router, error) {
	rows, err := stmt.Query(cdn)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	monitors := []Monitor{}
	caches := []Cache{}
	routers := []Router{}

	for rows.Next() {
		ttype := ""
		hashId := ""
		interfaceName := ""
		s := BasicServer{}
		if err := rows.Scan(&s.HostName, &s.FQDN, &s.Status, &s.Cachegroup, &s.Port, &s.IP, &s.IP6, &s.Profile, &interfaceName, &ttype, &hashId); err != nil {
			return nil, nil, nil, err
		}
		if ttype == MonitorType {
			monitors = append(monitors, Monitor{BasicServer: s})
		} else if strings.HasPrefix(ttype, "EDGE") || strings.HasPrefix(ttype, "MID") {
			caches = append(caches, Cache{
				BasicServer:   s,
				InterfaceName: interfaceName,
				Type:          ttype,
				HashID:        hashId,
			})
		} else if ttype == RouterType {
			routers = append(routers, Router{
				Type:    ttype,
				Profile: s.Profile,
			})
		}
	}
	return monitors, caches, routers, nil
}

func getCachegroups(stmt *sql.Stmt, cdn string) ([]Cachegroup, error) {
	rows, err := stmt.Query(cdn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cachegroups := []Cachegroup{}
	for rows.Next() {
		cg := Cachegroup{}
		if err := rows.Scan(&cg.Name, &cg.Coordinates.Latitude, &cg.Coordinates.Longitude); err != nil {
			return nil, err
		}
		cachegroups = append(cachegroups, cg)
	}
	return cachegroups, nil
}

func getProfiles(stmt *sql.Stmt, caches []Cache, routers []Router) ([]Profile, error) {
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

	rows, err := stmt.Query(pq.Array(profileNames), CacheMonitorConfigFile)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		profileName := ""
		name := ""
		value := ""
		if err := rows.Scan(&profileName, &name, &value); err != nil {
			return nil, err
		}
		if name == "" {
			return nil, fmt.Errorf("null name") // TODO continue and warn?
		}
		profile := profiles[profileName]
		if profile.Parameters == nil {
			profile.Parameters = map[string]interface{}{}
		}

		if valNum, err := strconv.Atoi(value); err == nil {
			profile.Parameters[name] = valNum
		} else {
			profile.Parameters[name] = value
		}
		profiles[profileName] = profile

	}

	profilesArr := []Profile{} // TODO make for efficiency?
	for _, profile := range profiles {
		profilesArr = append(profilesArr, profile)
	}
	return profilesArr, nil
}

func getDeliveryServices(stmt *sql.Stmt, routers []Router) ([]DeliveryService, error) {
	profileNames := []string{}
	for _, router := range routers {
		profileNames = append(profileNames, router.Profile)
	}

	rows, err := stmt.Query(pq.Array(profileNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dses := []DeliveryService{}
	for rows.Next() {
		mbps := 0.0
		ds := DeliveryService{Status: DeliveryServiceStatus}
		if err := rows.Scan(&ds.XMLID, &ds.TotalTPSThreshold, &mbps); err != nil {
			return nil, err
		}
		ds.TotalKBPSThreshold = mbps * KilobitsPerMegabit
		dses = append(dses, ds)
	}
	return dses, nil
}

func getConfig(stmt *sql.Stmt) (map[string]interface{}, error) {
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cfg := map[string]interface{}{}

	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, err
		}
		if valNum, err := strconv.Atoi(val); err == nil {
			cfg[name] = valNum
		} else {
			cfg[name] = val
		}
	}
	return cfg, nil
}

func getMonitoringJson(cdnName string, d *MonitoringData) (*MonitoringResponse, error) {
	monitors, caches, routers, err := getMonitoringServers(d.Servers, cdnName)
	if err != nil {
		return nil, fmt.Errorf("error getting servers: %v", err)
	}

	cachegroups, err := getCachegroups(d.Cachegroups, cdnName)
	if err != nil {
		return nil, fmt.Errorf("error getting cachegroups: %v", err)
	}

	profiles, err := getProfiles(d.Profiles, caches, routers)
	if err != nil {
		return nil, fmt.Errorf("error getting profiles: %v", err)
	}

	deliveryServices, err := getDeliveryServices(d.DeliveryServices, routers)
	if err != nil {
		return nil, fmt.Errorf("error getting deliveryservices: %v", err)
	}

	config, err := getConfig(d.Config)
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
