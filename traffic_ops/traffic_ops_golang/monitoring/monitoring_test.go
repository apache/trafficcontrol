package monitoring

import (
	"context"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/parameter"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

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

func TestGetMonitoringServers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cdn := "mycdn"

	monitor := createMockMonitor()
	router := createMockRouter()
	cache := createMockCache("test")
	// Different caches with the 'same' interfaces (in value only)
	otherCache := createMockCache("test")
	otherCache.Type = "MID"
	cacheID := uint64(1)
	otherCacheID := uint64(2)
	cache3 := createMockCache("test")
	cache3.Type = "MID"
	cache3.Status = string(tc.CacheStatusOffline) // should be ignored
	cache3ID := uint64(3)

	mock.ExpectBegin()
	setupMockGetMonitoringServers(mock, monitor, router, []Cache{cache, otherCache, cache3}, []uint64{cacheID, otherCacheID, cache3ID}, cdn)

	dbCtx, f := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer f()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	monitors, caches, routers, err := getMonitoringServers(tx, cdn)
	if err != nil {
		t.Fatalf("getMonitoringServers expected: nil error, actual: %v", err)
	}

	if len(caches) != 2 {
		t.Fatalf("got %v caches, expecting 2", len(caches))
	}

	for _, cacheServer := range caches {
		var testCache Cache
		if cacheServer.Type == cache.Type {
			testCache = cache
		} else {
			testCache = otherCache
		}

		for _, interf := range testCache.Interfaces {
			if len(interf.IPAddresses) != 4 {
				t.Errorf("cache: %v, interface: %v, expected 4 ip addresses, got %v", testCache.HostName, interf.Name, len(interf.IPAddresses))
			}
		}
	}

	if len(monitors) != 1 {
		t.Fatalf("getMonitoringServers expected: len(monitors) == 1, actual: %v", len(monitors))
	}
	sqlMonitor := monitors[0]
	if sqlMonitor != monitor {
		t.Errorf("getMonitoringServers expected: monitor == %+v, actual: %+v", monitor, sqlMonitor)
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

	// Now test by setting the monitor without an IPv4 address
	monitor.IP = ""

	mock.ExpectBegin()
	setupMockGetMonitoringServersWithoutIPv4(mock, monitor, router, []Cache{cache, otherCache, cache3}, []uint64{cacheID, otherCacheID, cache3ID}, cdn)

	dbCtx, f = context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer f()
	tx, err = db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	monitors, _, _, err = getMonitoringServers(tx, cdn)
	if err != nil {
		t.Fatalf("getMonitoringServers expected: nil error, actual: %v", err)
	}

	if len(monitors) != 1 {
		t.Fatalf("getMonitoringServers expected: len(monitors) == 1, actual: %v", len(monitors))
	}
	sqlMonitor = monitors[0]
	if sqlMonitor != monitor {
		t.Errorf("getMonitoringServers expected: monitor == %+v, actual: %+v", monitor, sqlMonitor)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetProfileWithParams(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cdn := "mycdn"
	mock.ExpectBegin()
	expected := ExpectedGetParams()
	mockGetParams(mock, expected, cdn)
	mock.ExpectCommit()

	dbCtx, f := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	defer f()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}
	defer tx.Commit()

	actual, err := getConfig(tx, cdn)
	if err != nil {
		t.Fatalf("getConfig err expected: nil, actual: %v", err)
	}

	// Should be just one
	for k, v := range actual {
		if *expected[0].Name != k {
			t.Fatalf("Expected param name %s doesn't match actual %s", *expected[0].Name, k)
		}
		if *expected[0].Value != strconv.Itoa(v.(int)) {
			t.Fatalf("Expected param value %s doesn't match actual %s", *expected[0].Value, v)
		}
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

	dbCtx, f := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	defer f()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlCachegroups, err := getCachegroups(tx, cdn)
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

func TestGetProfiles(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	cache := createMockCache("test")
	router := createMockRouter()

	profiles := []Profile{
		{
			Name: router.Profile,
			Type: RouterType,
			Parameters: map[string]interface{}{
				"param0": "param0Val",
				"param1": "param1Val",
			},
		},
		{
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
	rows1 := sqlmock.NewRows([]string{"profile", "name", "value"})
	for _, profile := range profiles {
		for paramName, paramVal := range profile.Parameters {
			if profile.Name == "routerProfile" {
				rows = rows.AddRow(profile.Name, paramName, paramVal)
			} else if profile.Name == "cacheProfile" {
				rows1 = rows1.AddRow(profile.Name, paramName, paramVal)
			}
		}
	}

	caches := []Cache{cache}
	routers := []Router{router}

	profileNames := []string{"routerProfile"}
	mock.ExpectQuery("SELECT").WithArgs(pq.Array(profileNames), CacheMonitorConfigFile).WillReturnRows(rows)

	profileNames = []string{"cacheProfile"}
	mock.ExpectQuery("SELECT").WithArgs(pq.Array(profileNames), CacheMonitorConfigFile).WillReturnRows(rows1)

	dbCtx, f := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	defer f()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlProfiles, err := getProfiles(tx, caches, routers)
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

	deliveryservice := DeliveryService{
		XMLID:              "myDsid",
		TotalTPSThreshold:  42.42,
		Status:             DeliveryServiceStatus,
		TotalKBPSThreshold: 24.24,
		Type:               "HTTP",
		Topology:           "foo",
		HostRegexes:        []string{`.*\.example\..*`},
	}

	deliveryservices := []DeliveryService{deliveryservice}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"xml_id", "global_max_tps", "global_max_mbps", "ds_type", "topology", "host_regexes"})
	for _, deliveryservice := range deliveryservices {
		rows = rows.AddRow(deliveryservice.XMLID, deliveryservice.TotalTPSThreshold, deliveryservice.TotalKBPSThreshold/KilobitsPerMegabit,
			deliveryservice.Type, deliveryservice.Topology, "{"+strings.Join(deliveryservice.HostRegexes, ",")+"}")
	}

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	dbCtx, f := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	defer f()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlDeliveryservices, err := getDeliveryServices(tx, "cdn")
	if err != nil {
		t.Errorf("getDeliveryServices expected: nil error, actual: %v", err)
	}

	if len(deliveryservices) != len(sqlDeliveryservices) {
		t.Fatalf("getDeliveryServices expected: %+v actual: %+v", deliveryservices, sqlDeliveryservices)
	}

	for i, sqlDeliveryservice := range sqlDeliveryservices {
		deliveryservice := deliveryservices[i]
		if !reflect.DeepEqual(deliveryservice, sqlDeliveryservice) {
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

	mock.ExpectQuery("SELECT").WithArgs(tc.MonitorProfilePrefix+"%%", MonitorConfigFile, "mycdn").WillReturnRows(rows)

	dbCtx, f := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	defer f()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlConfig, err := getConfig(tx, "mycdn")
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
		monitor := createMockMonitor()
		cache := createMockCache("test")
		router := createMockRouter()

		setupMockGetMonitoringServers(mock, monitor, router, []Cache{cache}, []uint64{1}, cdn)
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
		cache := createMockCache("test")
		router := createMockRouter()

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
		rows1 := sqlmock.NewRows([]string{"profile", "name", "value"})
		for _, profile := range profiles {
			for paramName, paramVal := range profile.Parameters {
				if profile.Name == "routerProfile" {
					rows = rows.AddRow(profile.Name, paramName, paramVal)
				} else if profile.Name == "cacheProfile" {
					rows1 = rows1.AddRow(profile.Name, paramName, paramVal)
				}
			}
		}

		// caches := []Cache{cache}
		// routers := []Router{router}

		profileNames := []string{"routerProfile"}
		mock.ExpectQuery("SELECT").WithArgs(pq.Array(profileNames), CacheMonitorConfigFile).WillReturnRows(rows)
		resp.Response.Profiles = profiles

		profileNames = []string{"cacheProfile"}
		mock.ExpectQuery("SELECT").WithArgs(pq.Array(profileNames), CacheMonitorConfigFile).WillReturnRows(rows1)
	}
	{
		//
		// getDeliveryServices
		//

		deliveryservice := DeliveryService{
			XMLID:              "myDsid",
			TotalTPSThreshold:  42.42,
			Status:             DeliveryServiceStatus,
			TotalKBPSThreshold: 24.24,
			Type:               "HTTP",
			Topology:           "foo",
			HostRegexes:        []string{`.*\.example\..*`},
		}

		deliveryservices := []DeliveryService{deliveryservice}
		// routers := []Router{router}

		rows := sqlmock.NewRows([]string{"xml_id", "global_max_tps", "global_max_mbps", "ds_type", "topology", "host_regexes"})
		for _, deliveryservice := range deliveryservices {
			rows = rows.AddRow(deliveryservice.XMLID, deliveryservice.TotalTPSThreshold, deliveryservice.TotalKBPSThreshold/KilobitsPerMegabit,
				deliveryservice.Type, deliveryservice.Topology, "{"+strings.Join(deliveryservice.HostRegexes, ",")+"}")
		}

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
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
	{
		//
		// topologies
		//
		topologies := map[string]tc.CRConfigTopology{
			"foo": {Nodes: []string{"cg1"}},
		}

		rows := sqlmock.NewRows([]string{"name", "nodes"})
		for name, nodes := range topologies {
			rows = rows.AddRow(name, "{"+strings.Join(nodes.Nodes, ",")+"}")
		}
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		resp.Response.Topologies = topologies
	}

	dbCtx, f := context.WithTimeout(context.TODO(), time.Duration(10)*time.Second)
	defer f()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		t.Fatalf("creating transaction: %v", err)
	}

	sqlResp, err := GetMonitoringJSON(tx, cdn)
	if err != nil {
		t.Fatalf("GetMonitoringJSON expected: nil error, actual: %v", err)
	}
	for _, cache := range resp.Response.TrafficServers {
		cache.Interfaces = sortInterfaces(cache.Interfaces)
	}
	for _, cache := range sqlResp.TrafficServers {
		cache.Interfaces = sortInterfaces(cache.Interfaces)
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
	if !reflect.DeepEqual(sqlResp.Topologies, resp.Response.Topologies) {
		t.Errorf("GetMonitoringJSON expected Topologies: %+v actual: %+v", resp.Response.Topologies, sqlResp.Topologies)
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

type SortableInterfaces []tc.ServerInterfaceInfo

func (s SortableInterfaces) Len() int {
	return len(s)
}
func (s SortableInterfaces) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableInterfaces) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func sortInterfaces(i []tc.ServerInterfaceInfo) []tc.ServerInterfaceInfo {
	sort.Sort(SortableInterfaces(i))
	return i
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

func createMockMonitor() Monitor {
	return Monitor{
		BasicServer: BasicServer{
			CommonServerProperties: CommonServerProperties{
				Profile:    "monitorProfile",
				Status:     "monitorStatus",
				Port:       8081,
				Cachegroup: "monitorCachegroup",
				HostName:   "monitorHost",
				FQDN:       "monitorFqdn.me",
			},
			IP:  "5.6.7.10",
			IP6: "2020::10",
		},
	}
}

func setupMockGetMonitoringServersWithoutIPv4(mock sqlmock.Sqlmock, monitor Monitor, router Router, caches []Cache, cacheIDs []uint64, cdn string) {
	serverRows := sqlmock.NewRows([]string{"hostName", "fqdn", "status", "cachegroup", "port", "profile", "type", "hashId", "serverID"})
	interfaceRows := sqlmock.NewRows([]string{"name", "max_bandwidth", "mtu", "monitor", "server"})
	ipAddressRows := sqlmock.NewRows([]string{"address", "gateway", "service_address", "server", "interface"})
	dssRows := sqlmock.NewRows([]string{"host_name", "xml_id"})
	serverRows = serverRows.AddRow(monitor.HostName, monitor.FQDN, monitor.Status, monitor.Cachegroup, monitor.Port, monitor.Profile, MonitorType, "noHash", 5)
	for index, cache := range caches {
		serverRows = serverRows.AddRow(cache.HostName, cache.FQDN, cache.Status, cache.Cachegroup, cache.Port, cache.Profile, cache.Type, cache.HashID, cacheIDs[index])
		dssRows = dssRows.AddRow(cache.HostName, "xml_id_foo")

		interfaceRows = interfaceRows.AddRow("none", nil, 1500, false, 0)
		for _, interf := range cache.Interfaces {
			interfaceRows = interfaceRows.AddRow(interf.Name, interf.MaxBandwidth, interf.MTU, interf.Monitor, cacheIDs[index])

			for _, ip := range interf.IPAddresses {
				ipAddressRows = ipAddressRows.AddRow(ip.Address, ip.Gateway, ip.ServiceAddress, cacheIDs[index], interf.Name)
				//Create two orphaned records
				ipAddressRows = ipAddressRows.AddRow("0.0.0.0", "0.0.0.0", false, 0, interf.Name)
				ipAddressRows = ipAddressRows.AddRow("0.0.0.0", "0.0.0.0", false, cacheIDs[index], "none")
			}
		}
	}
	// Add an interface and only ipv6 ip address for the monitor cache
	interfaceRows = interfaceRows.AddRow("monitorCacheInterface", nil, 1500, false, 5)
	ipAddressRows = ipAddressRows.AddRow("2020::10", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false, 5, "monitorCacheInterface")
	serverRows = serverRows.AddRow("noHostname", "noFqdn", "noStatus", "noGroup", 0, router.Profile, RouterType, "noHashid", 3)
	mock.ExpectQuery("SELECT (.+) FROM interface i (.+)").WithArgs(cdn).WillReturnRows(interfaceRows)
	mock.ExpectQuery("SELECT (.+) FROM ip_address ip (.+)").WillReturnRows(ipAddressRows)
	mock.ExpectQuery("SELECT (.+) FROM deliveryservice_server AS dss (.+)").WillReturnRows(dssRows)
	mock.ExpectQuery("SELECT (.+) FROM server me (.+)").WithArgs(cdn).WillReturnRows(serverRows)
}

func setupMockGetMonitoringServers(mock sqlmock.Sqlmock, monitor Monitor, router Router, caches []Cache, cacheIDs []uint64, cdn string) {
	serverRows := sqlmock.NewRows([]string{"hostName", "fqdn", "status", "cachegroup", "port", "profile", "type", "hashId", "serverID"})
	interfaceRows := sqlmock.NewRows([]string{"name", "max_bandwidth", "mtu", "monitor", "server"})
	ipAddressRows := sqlmock.NewRows([]string{"address", "gateway", "service_address", "server", "interface"})
	dssRows := sqlmock.NewRows([]string{"host_name", "xml_id"})
	serverRows = serverRows.AddRow(monitor.HostName, monitor.FQDN, monitor.Status, monitor.Cachegroup, monitor.Port, monitor.Profile, MonitorType, "noHash", 5)
	for index, cache := range caches {
		serverRows = serverRows.AddRow(cache.HostName, cache.FQDN, cache.Status, cache.Cachegroup, cache.Port, cache.Profile, cache.Type, cache.HashID, cacheIDs[index])
		dssRows = dssRows.AddRow(cache.HostName, "xml_id_foo")

		interfaceRows = interfaceRows.AddRow("none", nil, 1500, false, 0)
		for _, interf := range cache.Interfaces {
			interfaceRows = interfaceRows.AddRow(interf.Name, interf.MaxBandwidth, interf.MTU, interf.Monitor, cacheIDs[index])

			for _, ip := range interf.IPAddresses {
				ipAddressRows = ipAddressRows.AddRow(ip.Address, ip.Gateway, ip.ServiceAddress, cacheIDs[index], interf.Name)
				//Create two orphaned records
				ipAddressRows = ipAddressRows.AddRow("0.0.0.0", "0.0.0.0", false, 0, interf.Name)
				ipAddressRows = ipAddressRows.AddRow("0.0.0.0", "0.0.0.0", false, cacheIDs[index], "none")
			}
		}
	}
	// Add an interface and ip addresses for the monitor cache
	interfaceRows = interfaceRows.AddRow("monitorCacheInterface", nil, 1500, false, 5)
	ipAddressRows = ipAddressRows.AddRow("5.6.7.10", "10.0.0.0", true, 5, "monitorCacheInterface")
	ipAddressRows = ipAddressRows.AddRow("2020::10", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false, 5, "monitorCacheInterface")
	serverRows = serverRows.AddRow("noHostname", "noFqdn", "noStatus", "noGroup", 0, router.Profile, RouterType, "noHashid", 3)
	mock.ExpectQuery("SELECT (.+) FROM interface i (.+)").WithArgs(cdn).WillReturnRows(interfaceRows)
	mock.ExpectQuery("SELECT (.+) FROM ip_address ip (.+)").WillReturnRows(ipAddressRows)
	mock.ExpectQuery("SELECT (.+) FROM deliveryservice_server AS dss (.+)").WillReturnRows(dssRows)
	mock.ExpectQuery("SELECT (.+) FROM server me (.+)").WithArgs(cdn).WillReturnRows(serverRows)
}

func mockGetParams(mock sqlmock.Sqlmock, expected []parameter.TOParameter, cdn string) {
	rows := sqlmock.NewRows([]string{"name", "value"})
	for _, param := range expected {
		n := param.Name
		v := param.Value
		rows = rows.AddRow(*n, *v)
	}
	mock.ExpectQuery("SELECT").WithArgs(tc.MonitorProfilePrefix+"%%", MonitorConfigFile, cdn).WillReturnRows(rows)
}

func createMockCache(interfaceName string) Cache {
	return Cache{
		CommonServerProperties: CommonServerProperties{
			Profile:    "cacheProfile",
			Status:     "REPORTED",
			Port:       8081,
			Cachegroup: "cacheCachegroup",
			HostName:   "cacheHost",
			FQDN:       "cacheFqdn.me",
		},
		Interfaces: []tc.ServerInterfaceInfo{
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "5.6.7.8",
						Gateway:        util.StrPtr("5.6.7.0/24"),
						ServiceAddress: true,
					},
					{
						Address:        "2020::4",
						Gateway:        util.StrPtr("fd53::9"),
						ServiceAddress: true,
					},
					{
						Address:        "5.6.7.9",
						Gateway:        util.StrPtr("5.6.7.0/24"),
						ServiceAddress: false,
					},
					{
						Address:        "2021::4",
						Gateway:        util.StrPtr("fd53::9"),
						ServiceAddress: false,
					},
				},
				MaxBandwidth: util.UInt64Ptr(2500),
				Monitor:      true,
				MTU:          util.UInt64Ptr(1500),
				Name:         interfaceName + "1",
			},
			{
				IPAddresses: []tc.ServerIPAddress{
					{
						Address:        "6.7.8.9",
						Gateway:        util.StrPtr("6.7.8.0/24"),
						ServiceAddress: true,
					},
					{
						Address:        "2021::4",
						Gateway:        util.StrPtr("fd54::9"),
						ServiceAddress: true,
					},
					{
						Address:        "6.6.7.9",
						Gateway:        util.StrPtr("6.6.7.0/24"),
						ServiceAddress: false,
					},
					{
						Address:        "2022::4",
						Gateway:        util.StrPtr("fd53::9"),
						ServiceAddress: false,
					},
				},
				MaxBandwidth: util.UInt64Ptr(1500),
				Monitor:      false,
				MTU:          util.UInt64Ptr(1500),
				Name:         interfaceName + "2",
			},
		},
		Type:             "EDGE",
		HashID:           "cacheHash",
		DeliveryServices: []tc.TSDeliveryService{{XmlId: "xml_id_foo"}},
	}
}

func createMockRouter() Router {
	return Router{
		Type:    RouterType,
		Profile: "routerProfile",
	}
}

func ExpectedGetParams() []parameter.TOParameter {
	name := "peers.polling.interval"
	value := "3000"
	return []parameter.TOParameter{
		{
			APIInfoImpl: api.APIInfoImpl{ReqInfo: nil},
			ParameterNullable: tc.ParameterNullable{
				ConfigFile:  nil,
				ID:          nil,
				LastUpdated: nil,
				Name:        &name,
				Profiles:    nil,
				Secure:      nil,
				Value:       &value,
			},
		},
	}
}
