package monitoring

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
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetMonitoringServers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cdn := "mycdn"

	monitor := Monitor{
		BasicServer: BasicServer{Profile: "monitorProfile",
			Status:     "monitorStatus",
			IP:         "1.2.3.4",
			IP6:        "2001::4",
			Port:       8081,
			Cachegroup: "monitorCachegroup",
			HostName:   "monitorHost",
			FQDN:       "monitorFqdn.me",
		},
	}

	cacheType := "EDGE"
	cache := Cache{
		BasicServer: BasicServer{
			Profile:    "cacheProfile",
			Status:     "cacheStatus",
			IP:         "1.2.3.4",
			IP6:        "2001::4",
			Port:       8081,
			Cachegroup: "cacheCachegroup",
			HostName:   "cacheHost",
			FQDN:       "cacheFqdn.me",
		},
		InterfaceName: "cacheInterface",
		Type:          cacheType,
		HashID:        "cacheHash",
	}

	router := Router{
		Type:    RouterType,
		Profile: "routerProfile",
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"hostName", "fqdn", "status", "cachegroup", "port", "ip", "ip6", "profile", "interfaceName", "type", "hashId"})
	rows = rows.AddRow(monitor.HostName, monitor.FQDN, monitor.Status, monitor.Cachegroup, monitor.Port, monitor.IP, monitor.IP6, monitor.Profile, "noInterface", MonitorType, "noHash")
	rows = rows.AddRow(cache.HostName, cache.FQDN, cache.Status, cache.Cachegroup, cache.Port, cache.IP, cache.IP6, cache.Profile, cache.InterfaceName, cache.Type, cache.HashID)
	rows = rows.AddRow("noHostname", "noFqdn", "noStatus", "noGroup", 0, "noIp", "noIp6", router.Profile, "noInterface", RouterType, "noHashid")

	mock.ExpectQuery("SELECT").WithArgs(cdn).WillReturnRows(rows)

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	monitors, caches, routers, err := getMonitoringServers(tx, cdn, true)
	if err != nil {
		t.Errorf("getMonitoringServers expected: nil error, actual: %v", err)
	}

	if len(monitors) != 1 {
		t.Errorf("getMonitoringServers expected: len(monitors) == 1, actual: %v", len(monitors))
	}
	sqlMonitor := monitors[0]
	if sqlMonitor != monitor {
		t.Errorf("getMonitoringServers expected: monitor == %+v, actual: %+v", monitor, sqlMonitor)
	}

	if len(caches) != 1 {
		t.Errorf("getMonitoringServers expected: len(caches) == 1, actual: %v", len(caches))
	}
	sqlCache := caches[0]
	if sqlCache != cache {
		t.Errorf("getMonitoringServers expected: cache == %+v, actual: %+v", cache, sqlCache)
	}

	if len(routers) != 1 {
		t.Errorf("getMonitoringServers expected: len(routers) == 1, actual: %v", len(routers))
	}
	sqlRouter := routers[0]
	if sqlRouter != router {
		t.Errorf("getMonitoringServers expected: router == %+v, actual: %+v", router, sqlRouter)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetCachegroups(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cdn := "mycdn"

	cachegroup := Cachegroup{
		Name: "myCachegroup",
		Coordinates: Coordinates{
			Latitude:  42.42,
			Longitude: 24.24,
		},
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"name", "latitude", "longitude"})
	rows = rows.AddRow(cachegroup.Name, cachegroup.Coordinates.Latitude, cachegroup.Coordinates.Longitude)

	mock.ExpectQuery("SELECT").WithArgs(cdn).WillReturnRows(rows)

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlCachegroups, err := getCachegroups(tx, cdn, true)
	if err != nil {
		t.Errorf("getCachegroups expected: nil error, actual: %v", err)
	}

	if len(sqlCachegroups) != 1 {
		t.Errorf("getCachegroups expected: len(monitors) == 1, actual: %v", len(sqlCachegroups))
	}
	sqlCachegroup := sqlCachegroups[0]
	if sqlCachegroup != cachegroup {
		t.Errorf("getMonitoringServers expected: cachegroup == %+v, actual: %+v", cachegroup, sqlCachegroup)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

type SortableProfiles []Profile

func (s SortableProfiles) Len() int {
	return len(s)
}
func (s SortableProfiles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableProfiles) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func sortProfiles(p []Profile) []Profile {
	sort.Sort(SortableProfiles(p))
	return p
}

type SortableMonitors []Monitor

func (s SortableMonitors) Len() int {
	return len(s)
}
func (s SortableMonitors) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableMonitors) Less(i, j int) bool {
	return s[i].HostName < s[j].HostName
}

func sortMonitors(p []Monitor) []Monitor {
	sort.Sort(SortableMonitors(p))
	return p
}

type SortableCaches []Cache

func (s SortableCaches) Len() int {
	return len(s)
}
func (s SortableCaches) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableCaches) Less(i, j int) bool {
	return s[i].HostName < s[j].HostName
}

func sortCaches(p []Cache) []Cache {
	sort.Sort(SortableCaches(p))
	return p
}

type SortableCachegroups []Cachegroup

func (s SortableCachegroups) Len() int {
	return len(s)
}
func (s SortableCachegroups) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableCachegroups) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func sortCachegroups(p []Cachegroup) []Cachegroup {
	sort.Sort(SortableCachegroups(p))
	return p
}

type SortableDeliveryServices []DeliveryService

func (s SortableDeliveryServices) Len() int {
	return len(s)
}
func (s SortableDeliveryServices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableDeliveryServices) Less(i, j int) bool {
	return s[i].XMLID < s[j].XMLID
}

func sortDeliveryServices(p []DeliveryService) []DeliveryService {
	sort.Sort(SortableDeliveryServices(p))
	return p
}

func TestGetProfiles(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cache := Cache{
		BasicServer: BasicServer{
			Profile:    "cacheProfile",
			Status:     "cacheStatus",
			IP:         "1.2.3.4",
			IP6:        "2001::4",
			Port:       8081,
			Cachegroup: "cacheCachegroup",
			HostName:   "cacheHost",
			FQDN:       "cacheFqdn.me",
		},
		InterfaceName: "cacheInterface",
		Type:          "EDGE",
		HashID:        "cacheHash",
	}

	router := Router{
		Type:    RouterType,
		Profile: "routerProfile",
	}

	profiles := []Profile{
		Profile{
			Name: router.Profile,
			Type: RouterType,
			Parameters: map[string]interface{}{
				"param0": "param0Val",
				"param1": "param1Val",
			},
		},
		Profile{
			Name: cache.Profile,
			Type: "myType2",
			Parameters: map[string]interface{}{
				"2param0": "2param0Val",
				"2param1": "2param1Val",
			},
		},
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"profile", "name", "value"})
	for _, profile := range profiles {
		for paramName, paramVal := range profile.Parameters {
			rows = rows.AddRow(profile.Name, paramName, paramVal)
		}
	}

	caches := []Cache{cache}
	routers := []Router{router}
	profileNames := []string{"cacheProfile"}
	cdn := "mycdn"

	mock.ExpectQuery("SELECT").WithArgs(cdn, pq.Array(profileNames), CacheMonitorConfigFile).WillReturnRows(rows)

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlProfiles, err := getProfiles(tx, cdn, caches, routers, true)
	if err != nil {
		t.Errorf("getProfiles expected: nil error, actual: %v", err)
	}

	if len(sqlProfiles) != len(profiles) {
		t.Errorf("getProfiles expected: %+v actual: %+v", profiles, sqlProfiles)
	}

	profiles = sortProfiles(profiles)
	sqlProfiles = sortProfiles(sqlProfiles)

	for i, profile := range profiles {
		if profile.Name != sqlProfiles[i].Name {
			t.Errorf("getProfiles expected: profiles[%v].Name %v, actual: %v", i, profile.Name, sqlProfiles[i].Name)
		}
		for paramName, paramVal := range profile.Parameters {
			val, ok := sqlProfiles[i].Parameters[paramName]
			if !ok {
				t.Errorf("getProfiles expected: profiles[%v].Parameters[%v] = %v, actual: %v", i, paramName, paramVal, val)
			}
			if val != paramVal {
				t.Errorf("getProfiles expected: profiles[%v].Parameters[%v] = %v, actual: %v", i, paramName, paramVal, val)
			}
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetDeliveryServices(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	router := Router{
		Type:    RouterType,
		Profile: "routerProfile",
	}

	deliveryservice := DeliveryService{
		XMLID:              "myDsid",
		TotalTPSThreshold:  42.42,
		Status:             DeliveryServiceStatus,
		TotalKBPSThreshold: 24.24,
	}

	deliveryservices := []DeliveryService{deliveryservice}
	routers := []Router{router}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"xml_id", "global_max_tps", "global_max_mbps"})
	for _, deliveryservice := range deliveryservices {
		rows = rows.AddRow(deliveryservice.XMLID, deliveryservice.TotalTPSThreshold, deliveryservice.TotalKBPSThreshold/KilobitsPerMegabit)
	}

	profileNames := []string{router.Profile}

	mock.ExpectQuery("SELECT").WithArgs(pq.Array(profileNames)).WillReturnRows(rows)

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlDeliveryservices, err := getDeliveryServices(tx, routers, true)
	if err != nil {
		t.Errorf("getProfiles expected: nil error, actual: %v", err)
	}

	if len(deliveryservices) != len(sqlDeliveryservices) {
		t.Errorf("getProfiles expected: %+v actual: %+v", deliveryservices, sqlDeliveryservices)
	}

	for i, sqlDeliveryservice := range sqlDeliveryservices {
		deliveryservice := deliveryservices[i]
		if deliveryservice != sqlDeliveryservice {
			t.Errorf("getDeliveryServices expected: %v, actual: %v", deliveryservice, sqlDeliveryservice)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetConfig(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	config := map[string]interface{}{
		"name0": "val0",
		"name1": "val1",
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"name", "value"})
	for name, val := range config {
		rows = rows.AddRow(name, val)
	}

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	cdn := "mycdn"

	sqlConfig, err := getConfig(tx, cdn, true)
	if err != nil {
		t.Errorf("getProfiles expected: nil error, actual: %v", err)
	}

	if !reflect.DeepEqual(config, sqlConfig) {
		t.Errorf("getConfig expected: %+v actual: %+v", config, sqlConfig)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetMonitoringJSON(t *testing.T) {
	resp := MonitoringResponse{}
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cdn := "mycdn"

	mock.ExpectBegin()
	{
		//
		// getMonitoringServers
		//
		monitor := Monitor{
			BasicServer: BasicServer{Profile: "monitorProfile",
				Status:     "monitorStatus",
				IP:         "1.2.3.4",
				IP6:        "2001::4",
				Port:       8081,
				Cachegroup: "monitorCachegroup",
				HostName:   "monitorHost",
				FQDN:       "monitorFqdn.me",
			},
		}

		cacheType := "EDGE"
		cache := Cache{
			BasicServer: BasicServer{
				Profile:    "cacheProfile",
				Status:     "cacheStatus",
				IP:         "1.2.3.4",
				IP6:        "2001::4",
				Port:       8081,
				Cachegroup: "cacheCachegroup",
				HostName:   "cacheHost",
				FQDN:       "cacheFqdn.me",
			},
			InterfaceName: "cacheInterface",
			Type:          cacheType,
			HashID:        "cacheHash",
		}

		router := Router{
			Type:    RouterType,
			Profile: "routerProfile",
		}

		rows := sqlmock.NewRows([]string{"hostName", "fqdn", "status", "cachegroup", "port", "ip", "ip6", "profile", "interfaceName", "type", "hashId"})
		rows = rows.AddRow(monitor.HostName, monitor.FQDN, monitor.Status, monitor.Cachegroup, monitor.Port, monitor.IP, monitor.IP6, monitor.Profile, "noInterface", MonitorType, "noHash")
		rows = rows.AddRow(cache.HostName, cache.FQDN, cache.Status, cache.Cachegroup, cache.Port, cache.IP, cache.IP6, cache.Profile, cache.InterfaceName, cache.Type, cache.HashID)
		rows = rows.AddRow("noHostname", "noFqdn", "noStatus", "noGroup", 0, "noIp", "noIp6", router.Profile, "noInterface", RouterType, "noHashid")

		mock.ExpectQuery("SELECT").WithArgs(cdn).WillReturnRows(rows)
		resp.Response.TrafficServers = []Cache{cache}
		resp.Response.TrafficMonitors = []Monitor{monitor}
	}
	{
		//
		// getCachegroups
		//
		cachegroup := Cachegroup{
			Name: "myCachegroup",
			Coordinates: Coordinates{
				Latitude:  42.42,
				Longitude: 24.24,
			},
		}

		rows := sqlmock.NewRows([]string{"name", "latitude", "longitude"})
		rows = rows.AddRow(cachegroup.Name, cachegroup.Coordinates.Latitude, cachegroup.Coordinates.Longitude)

		mock.ExpectQuery("SELECT").WithArgs(cdn).WillReturnRows(rows)
		resp.Response.Cachegroups = []Cachegroup{cachegroup}
	}
	{
		//
		// getProfiles
		//
		cache := Cache{
			BasicServer: BasicServer{
				Profile:    "cacheProfile",
				Status:     "cacheStatus",
				IP:         "1.2.3.4",
				IP6:        "2001::4",
				Port:       8081,
				Cachegroup: "cacheCachegroup",
				HostName:   "cacheHost",
				FQDN:       "cacheFqdn.me",
			},
			InterfaceName: "cacheInterface",
			Type:          "EDGE",
			HashID:        "cacheHash",
		}

		router := Router{
			Type:    RouterType,
			Profile: "routerProfile",
		}

		profiles := []Profile{
			Profile{
				Name: router.Profile,
				Type: RouterType,
				Parameters: map[string]interface{}{
					"param0": "param0Val",
					"param1": "param1Val",
				},
			},
			Profile{
				Name: cache.Profile,
				Type: "EDGE",
				Parameters: map[string]interface{}{
					"2param0": "2param0Val",
					"2param1": "2param1Val",
				},
			},
		}

		rows := sqlmock.NewRows([]string{"profile", "name", "value"})
		for _, profile := range profiles {
			for paramName, paramVal := range profile.Parameters {
				rows = rows.AddRow(profile.Name, paramName, paramVal)
			}
		}

		// caches := []Cache{cache}
		// routers := []Router{router}

		profileNames := []string{"cacheProfile"}

		mock.ExpectQuery("SELECT").WithArgs(cdn, pq.Array(profileNames), CacheMonitorConfigFile).WillReturnRows(rows)
		resp.Response.Profiles = profiles
	}
	{
		//
		// getDeliveryServices
		//
		router := Router{
			Type:    RouterType,
			Profile: "routerProfile",
		}

		deliveryservice := DeliveryService{
			XMLID:              "myDsid",
			TotalTPSThreshold:  42.42,
			Status:             DeliveryServiceStatus,
			TotalKBPSThreshold: 24.24,
		}

		deliveryservices := []DeliveryService{deliveryservice}
		// routers := []Router{router}

		rows := sqlmock.NewRows([]string{"xml_id", "global_max_tps", "global_max_mbps"})
		for _, deliveryservice := range deliveryservices {
			rows = rows.AddRow(deliveryservice.XMLID, deliveryservice.TotalTPSThreshold, deliveryservice.TotalKBPSThreshold/KilobitsPerMegabit)
		}

		profileNames := []string{router.Profile}

		mock.ExpectQuery("SELECT").WithArgs(pq.Array(profileNames)).WillReturnRows(rows)
		resp.Response.DeliveryServices = deliveryservices
	}
	{
		//
		// getConfig
		//
		config := map[string]interface{}{
			"name0": "val0",
			"name1": "val1",
		}

		rows := sqlmock.NewRows([]string{"name", "value"})
		for name, val := range config {
			rows = rows.AddRow(name, val)
		}

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		resp.Response.Config = config
	}

	dbCtx, _ := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlResp, err := GetMonitoringJSON(tx, cdn, true)
	if err != nil {
		t.Errorf("GetMonitoringJSON expected: nil error, actual: %v", err)
	}

	resp.Response.TrafficServers = sortCaches(resp.Response.TrafficServers)
	sqlResp.TrafficServers = sortCaches(sqlResp.TrafficServers)
	resp.Response.TrafficMonitors = sortMonitors(resp.Response.TrafficMonitors)
	sqlResp.TrafficMonitors = sortMonitors(sqlResp.TrafficMonitors)
	resp.Response.Cachegroups = sortCachegroups(resp.Response.Cachegroups)
	sqlResp.Cachegroups = sortCachegroups(sqlResp.Cachegroups)
	resp.Response.Profiles = sortProfiles(resp.Response.Profiles)
	sqlResp.Profiles = sortProfiles(sqlResp.Profiles)
	resp.Response.DeliveryServices = sortDeliveryServices(resp.Response.DeliveryServices)
	sqlResp.DeliveryServices = sortDeliveryServices(sqlResp.DeliveryServices)

	if !reflect.DeepEqual(sqlResp.TrafficServers, resp.Response.TrafficServers) {
		t.Errorf("GetMonitoringJSON expected TrafficServers: %+v actual: %+v", resp.Response.TrafficServers, sqlResp.TrafficServers)
	}
	if !reflect.DeepEqual(sqlResp.TrafficMonitors, resp.Response.TrafficMonitors) {
		t.Errorf("GetMonitoringJSON expected TrafficMonitors: %+v actual: %+v", resp.Response.TrafficMonitors, sqlResp.TrafficMonitors)
	}
	if !reflect.DeepEqual(sqlResp.Cachegroups, resp.Response.Cachegroups) {
		t.Errorf("GetMonitoringJSON expected Cachegroups: %+v actual: %+v", resp.Response.Cachegroups, sqlResp.Cachegroups)
	}
	if !reflect.DeepEqual(sqlResp.Profiles, resp.Response.Profiles) {
		t.Errorf("GetMonitoringJSON expected Profiles: %+v actual: %+v", resp.Response.Profiles, sqlResp.Profiles)
	}
	if !reflect.DeepEqual(sqlResp.DeliveryServices, resp.Response.DeliveryServices) {
		t.Errorf("GetMonitoringJSON expected DeliveryServices: %+v actual: %+v", resp.Response.DeliveryServices, sqlResp.DeliveryServices)
	}
	if !reflect.DeepEqual(sqlResp.Config, resp.Response.Config) {
		t.Errorf("GetMonitoringJSON expected Config: %+v actual: %+v", resp.Response.Config, sqlResp.Config)
	}

}
