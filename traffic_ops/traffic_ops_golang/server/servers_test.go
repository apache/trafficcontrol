package server

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
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type ServerAndInterfaces struct {
	Server    tc.ServerV5
	Interface tc.ServerInterfaceInfoV40
}

func getTestServers() []ServerAndInterfaces {
	servers := []ServerAndInterfaces{}
	testServerV40 := tc.ServerV5{
		CacheGroup:         "Cache Group",
		CacheGroupID:       1,
		CDNID:              1,
		CDN:                "cdnName",
		DomainName:         "domainName",
		GUID:               util.Ptr("guid"),
		HostName:           "server1",
		HTTPSPort:          util.Ptr(443),
		ID:                 1,
		ILOIPAddress:       util.Ptr("iloIpAddress"),
		ILOIPGateway:       util.Ptr("iloIpGateway"),
		ILOIPNetmask:       util.Ptr("iloIpNetmask"),
		ILOPassword:        util.Ptr("iloPassword"),
		ILOUsername:        util.Ptr("iloUsername"),
		LastUpdated:        time.Now(),
		MgmtIPAddress:      util.Ptr("mgmtIpAddress"),
		MgmtIPGateway:      util.Ptr("mgmtIpGateway"),
		MgmtIPNetmask:      util.Ptr("mgmtIpNetmask"),
		OfflineReason:      util.Ptr("offlineReason"),
		PhysicalLocation:   "physLocation",
		PhysicalLocationID: 1,
		Profiles:           []string{"profile"},
		Rack:               util.Ptr("rack"),
		Status:             "status",
		StatusID:           1,
		TCPPort:            util.Ptr(80),
		Type:               "EDGE",
		TypeID:             1,
		XMPPID:             util.Ptr("xmppId"),
		XMPPPasswd:         util.Ptr("xmppPasswd"),
		StatusLastUpdated:  &(time.Time{}),
		ConfigUpdateTime:   &(time.Time{}),
		ConfigApplyTime:    &(time.Time{}),
		ConfigUpdateFailed: false,
		RevalUpdateTime:    &(time.Time{}),
		RevalApplyTime:     &(time.Time{}),
		RevalUpdateFailed:  false,
	}
	mtu := uint64(9500)

	iface := tc.ServerInterfaceInfoV40{
		ServerInterfaceInfo: tc.ServerInterfaceInfo{
			IPAddresses: []tc.ServerIPAddress{
				{
					Address:        "ip6Address",
					Gateway:        nil,
					ServiceAddress: true,
				},
			},
			MaxBandwidth: nil,
			Monitor:      true,
			MTU:          &mtu,
			Name:         "interfaceName",
		},
		RouterHostName: "routerHostName",
		RouterPortName: "routerPortName",
	}

	servers = append(servers, ServerAndInterfaces{Server: testServerV40, Interface: iface})

	testServer2 := testServerV40
	testServer2.CacheGroup = "cachegroup2"
	testServer2.HostName = "server2"
	testServer2.ID = 2
	servers = append(servers, ServerAndInterfaces{Server: testServer2, Interface: iface})

	testServer3 := testServerV40
	testServer3.CacheGroup = "cachegroup3"
	testServer3.HostName = "server3"
	testServer3.ID = 3
	servers = append(servers, ServerAndInterfaces{Server: testServer3, Interface: iface})

	return servers
}

// Test to make sure that updating the "cdn" of a server already assigned to a DS fails
func TestCheckTypeChangeSafety(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testServers := getTestServers()

	rows := sqlmock.NewRows([]string{"type", "cdn_id"})
	// note here that the cdnid is 5, which is not the same as the initial cdnid of the fist traffic server
	rows.AddRow(testServers[0].Server.TypeID, 5)
	// Make it return a list of atleast one associated ds
	dsrows := sqlmock.NewRows([]string{"array"})
	dsrows.AddRow("{3}")
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectQuery("SELECT ARRAY").WillReturnRows(dsrows)

	s := tc.ServerV5{
		CDNID:  testServers[0].Server.CDNID,
		TypeID: testServers[0].Server.TypeID,
		ID:     testServers[0].Server.ID,
	}

	userErr, _, errCode := checkTypeChangeSafety(s, db.MustBegin())
	if errCode != 409 {
		t.Errorf("Update servers: Expected error code of %v, but got %v", 409, errCode)
	}
	expectedErr := "server cdn can not be updated when it is currently assigned to delivery services"
	if userErr == nil {
		t.Errorf("Update expected error: %v, but got no error with status: %s", expectedErr, http.StatusText(errCode))
	}
}

func TestGetServersByCachegroup(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testServers := getTestServers()

	unfilteredCols := []string{"count"}
	unfilteredRows := sqlmock.NewRows(unfilteredCols).AddRow(len(testServers))

	cols := []string{"cachegroup", "cachegroup_id", "cdn_id", "cdn_name", "domain_name", "guid", "host_name",
		"https_port", "id", "ilo_ip_address", "ilo_ip_gateway", "ilo_ip_netmask", "ilo_password", "ilo_username",
		"last_updated", "mgmt_ip_address", "mgmt_ip_gateway", "mgmt_ip_netmask", "offline_reason", "phys_location",
		"phys_location_id", "profile_name", "rack", "revalidate_update_time", "revalidate_apply_time",
		"revalidate_update_failed", "status", "status_id", "tcp_port", "server_type", "server_type_id",
		"config_update_time", "config_apply_time", "config_update_failed", "xmpp_id", "xmpp_passwd",
		"status_last_updated"}
	interfaceCols := []string{"max_bandwidth", "monitor", "mtu", "name", "server", "router_host_name", "router_port_name"}
	rows := sqlmock.NewRows(cols)
	interfaceRows := sqlmock.NewRows(interfaceCols)

	ipCols := []string{"address", "gateway", "service_address", "server", "interface"}
	ipRows := sqlmock.NewRows(ipCols)

	//TODO: drichardson - build helper to add these Rows from the struct values
	//                    or by CSV if types get in the way
	for _, srv := range testServers {
		ts := srv.Server
		rows = rows.AddRow(
			ts.CacheGroup,
			ts.CacheGroupID,
			ts.CDNID,
			ts.CDN,
			ts.DomainName,
			*ts.GUID,
			ts.HostName,
			*ts.HTTPSPort,
			ts.ID,
			*ts.ILOIPAddress,
			*ts.ILOIPGateway,
			*ts.ILOIPNetmask,
			*ts.ILOPassword,
			*ts.ILOUsername,
			ts.LastUpdated,
			*ts.MgmtIPAddress,
			*ts.MgmtIPGateway,
			*ts.MgmtIPNetmask,
			*ts.OfflineReason,
			ts.PhysicalLocation,
			ts.PhysicalLocationID,
			fmt.Sprintf("{%s}", strings.Join(ts.Profiles, ",")),
			*ts.Rack,
			*ts.RevalUpdateTime,
			*ts.RevalApplyTime,
			ts.RevalUpdateFailed,
			ts.Status,
			ts.StatusID,
			*ts.TCPPort,
			ts.Type,
			ts.TypeID,
			*ts.ConfigUpdateTime,
			*ts.ConfigApplyTime,
			ts.ConfigUpdateFailed,
			*ts.XMPPID,
			*ts.XMPPPasswd,
			*ts.StatusLastUpdated,
		)
		interfaceRows = interfaceRows.AddRow(
			srv.Interface.MaxBandwidth,
			srv.Interface.Monitor,
			srv.Interface.MTU,
			srv.Interface.Name,
			ts.ID,
			srv.Interface.RouterHostName,
			srv.Interface.RouterPortName,
		)

		for _, ip := range srv.Interface.IPAddresses {
			ipRows = ipRows.AddRow(
				ip.Address,
				ip.Gateway,
				ip.ServiceAddress,
				ts.ID,
				srv.Interface.Name,
			)
		}
	}

	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT COUNT\\(s.id\\) FROM s")
	mock.ExpectPrepare("SELECT COUNT\\(s.id\\) FROM s")
	mock.ExpectQuery("SELECT COUNT\\(s.id\\) FROM s").WillReturnRows(unfilteredRows)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectQuery("SELECT").WillReturnRows(interfaceRows)
	mock.ExpectQuery("SELECT").WillReturnRows(ipRows)

	v := map[string]string{"cachegroup": "2"}

	user := auth.CurrentUser{}

	version := api.Version{Major: 4, Minor: 0}

	servers, _, userErr, sysErr, errCode, _ := getServers(nil, v, db.MustBegin(), &user, false, version, false)
	if userErr != nil || sysErr != nil {
		t.Errorf("getServers expected: no errors, actual: %v %v with status: %s", userErr, sysErr, http.StatusText(errCode))
	}

	if len(servers) != 3 {
		t.Errorf("getServers expected: len(servers) == 3, actual: %v", len(servers))
	}
}

func TestGetMidServers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testServers := getTestServers()
	testServers = testServers[0:2]

	testServers[1].Server.CacheGroup = "parentCacheGroup"
	testServers[1].Server.CacheGroupID = 2
	testServers[1].Server.Type = "MID"

	unfilteredCols := []string{"count"}
	unfilteredRows := sqlmock.NewRows(unfilteredCols).AddRow(len(testServers))

	cols := []string{"cachegroup", "cachegroup_id", "cdn_id", "cdn_name", "domain_name", "guid", "host_name",
		"https_port", "id", "ilo_ip_address", "ilo_ip_gateway", "ilo_ip_netmask", "ilo_password", "ilo_username",
		"last_updated", "mgmt_ip_address", "mgmt_ip_gateway", "mgmt_ip_netmask", "offline_reason", "phys_location",
		"phys_location_id", "profile_name", "rack", "revalidate_update_time", "revalidate_apply_time",
		"revalidate_update_failed", "status", "status_id", "tcp_port", "server_type", "server_type_id",
		"config_update_time", "config_apply_time", "config_update_failed", "xmpp_id", "xmpp_passwd",
		"status_last_updated"}
	interfaceCols := []string{"max_bandwidth", "monitor", "mtu", "name", "server", "router_host_name", "router_port_name"}
	rows := sqlmock.NewRows(cols)
	interfaceRows := sqlmock.NewRows(interfaceCols)

	ipCols := []string{"address", "gateway", "service_address", "server", "interface"}
	ipRows := sqlmock.NewRows(ipCols)

	for _, srv := range testServers {
		ts := srv.Server
		rows = rows.AddRow(
			ts.CacheGroup,
			ts.CacheGroupID,
			ts.CDNID,
			ts.CDN,
			ts.DomainName,
			*ts.GUID,
			ts.HostName,
			*ts.HTTPSPort,
			ts.ID,
			*ts.ILOIPAddress,
			*ts.ILOIPGateway,
			*ts.ILOIPNetmask,
			*ts.ILOPassword,
			*ts.ILOUsername,
			ts.LastUpdated,
			*ts.MgmtIPAddress,
			*ts.MgmtIPGateway,
			*ts.MgmtIPNetmask,
			*ts.OfflineReason,
			ts.PhysicalLocation,
			ts.PhysicalLocationID,
			fmt.Sprintf("{%s}", strings.Join(ts.Profiles, ",")),
			*ts.Rack,
			*ts.RevalUpdateTime,
			*ts.RevalApplyTime,
			ts.RevalUpdateFailed,
			ts.Status,
			ts.StatusID,
			*ts.TCPPort,
			ts.Type,
			ts.TypeID,
			*ts.ConfigUpdateTime,
			*ts.ConfigApplyTime,
			ts.ConfigUpdateFailed,
			*ts.XMPPID,
			*ts.XMPPPasswd,
			*ts.StatusLastUpdated,
		)
		interfaceRows = interfaceRows.AddRow(
			srv.Interface.MaxBandwidth,
			srv.Interface.Monitor,
			srv.Interface.MTU,
			srv.Interface.Name,
			ts.ID,
			srv.Interface.RouterHostName,
			srv.Interface.RouterPortName,
		)

		for _, ip := range srv.Interface.IPAddresses {
			ipRows = ipRows.AddRow(
				ip.Address,
				ip.Gateway,
				ip.ServiceAddress,
				ts.ID,
				srv.Interface.Name,
			)
		}
	}
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT COUNT\\(s.id\\) FROM s")
	mock.ExpectPrepare("SELECT COUNT\\(s.id\\) FROM s")
	mock.ExpectQuery("SELECT COUNT\\(s.id\\) FROM s").WillReturnRows(unfilteredRows)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectQuery("SELECT").WillReturnRows(interfaceRows)
	mock.ExpectQuery("SELECT").WillReturnRows(ipRows)
	v := map[string]string{}

	user := auth.CurrentUser{}
	version := api.Version{Major: 4, Minor: 0}
	servers, _, userErr, sysErr, errCode, _ := getServers(nil, v, db.MustBegin(), &user, false, version, false)

	if userErr != nil || sysErr != nil {
		t.Errorf("getServers expected: no errors, actual: %v %v with status: %s", userErr, sysErr, http.StatusText(errCode))
	}

	cols2 := []string{"cachegroup", "cachegroup_id", "cdn_id", "cdn_name", "domain_name", "guid", "host_name",
		"https_port", "id", "ilo_ip_address", "ilo_ip_gateway", "ilo_ip_netmask", "ilo_password", "ilo_username",
		"last_updated", "mgmt_ip_address", "mgmt_ip_gateway", "mgmt_ip_netmask", "offline_reason", "phys_location",
		"phys_location_id", "profile_name", "rack", "revalidate_update_time", "revalidate_apply_time",
		"revalidate_update_failed", "status", "status_id", "tcp_port", "server_type", "server_type_id",
		"config_update_time", "config_apply_time", "config_update_failed", "xmpp_id", "xmpp_passwd",
		"status_last_updated"}
	rows2 := sqlmock.NewRows(cols2)

	cgs := []tc.CacheGroup{}
	testCG1 := tc.CacheGroup{
		ID:                          1,
		Name:                        "Cachegroup",
		ShortName:                   "cg1",
		Latitude:                    38.7,
		Longitude:                   90.7,
		ParentCachegroupID:          2,
		SecondaryParentCachegroupID: 2,
		LocalizationMethods: []tc.LocalizationMethod{
			tc.LocalizationMethodDeepCZ,
			tc.LocalizationMethodCZ,
			tc.LocalizationMethodGeo,
		},
		Type:        "EDGE_LOC",
		TypeID:      6,
		LastUpdated: tc.TimeNoMod{Time: time.Now()},
		Fallbacks: []string{
			"cachegroup2",
			"cachegroup3",
		},
		FallbackToClosest: true,
	}
	cgs = append(cgs, testCG1)
	testCG2 := tc.CacheGroup{
		ID:                          2,
		Name:                        "parentCacheGroup",
		ShortName:                   "pg1",
		Latitude:                    38.7,
		Longitude:                   90.7,
		ParentCachegroupID:          1,
		SecondaryParentCachegroupID: 1,
		LocalizationMethods: []tc.LocalizationMethod{
			tc.LocalizationMethodDeepCZ,
			tc.LocalizationMethodCZ,
			tc.LocalizationMethodGeo,
		},
		Type:        "MID_LOC",
		TypeID:      7,
		LastUpdated: tc.TimeNoMod{Time: time.Now()},
	}
	cgs = append(cgs, testCG2)

	serverMap := make(map[int]tc.ServerV5, len(servers))
	serverIDs := make([]int, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, server.ID)
		serverMap[server.ID] = server
	}

	var ts tc.ServerV5
	for _, s := range servers {
		if s.HostName == "server2" {
			ts = s
			break
		}
	}
	ts.ID++
	rows2 = rows2.AddRow(
		ts.CacheGroup,
		ts.CacheGroupID,
		ts.CDNID,
		ts.CDN,
		ts.DomainName,
		*ts.GUID,
		ts.HostName,
		*ts.HTTPSPort,
		ts.ID,
		*ts.ILOIPAddress,
		*ts.ILOIPGateway,
		*ts.ILOIPNetmask,
		*ts.ILOPassword,
		*ts.ILOUsername,
		ts.LastUpdated,
		*ts.MgmtIPAddress,
		*ts.MgmtIPGateway,
		*ts.MgmtIPNetmask,
		*ts.OfflineReason,
		ts.PhysicalLocation,
		ts.PhysicalLocationID,
		fmt.Sprintf("{%s}", strings.Join(ts.Profiles, ",")),
		*ts.Rack,
		*ts.RevalUpdateTime,
		*ts.RevalApplyTime,
		ts.RevalUpdateFailed,
		ts.Status,
		ts.StatusID,
		*ts.TCPPort,
		ts.Type,
		ts.TypeID,
		*ts.ConfigUpdateTime,
		*ts.ConfigApplyTime,
		ts.ConfigUpdateFailed,
		*ts.XMPPID,
		*ts.XMPPPasswd,
		*ts.StatusLastUpdated,
	)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows2)
	mid, userErr, sysErr, errCode := getMidServers(serverIDs, serverMap, 0, 0, db.MustBegin(), false)

	if userErr != nil || sysErr != nil {
		t.Fatalf("getMidServers expected: no errors, actual: %v %v with status: %s", userErr, sysErr, http.StatusText(errCode))
	}
	if len(mid) != 1 {
		t.Fatalf("getMidServers expected: len(mid) == 1, actual: %v", len(mid))
	}
	if serverMap[mid[0]].Type != "MID" {
		t.Errorf("getMidServers expected: Type == MID, actual: %v", serverMap[mid[0]].Type)
	}

	if actual := serverMap[mid[0]].CacheGroup; actual != "parentCacheGroup" {
		t.Errorf("getMidServers expected: Cachegroup == parentCacheGroup, actual: %s", actual)
	}

	if actual := serverMap[mid[0]].CacheGroupID; actual != 2 {
		t.Errorf("getMidServers expected: CachegroupID == 2, actual: %d", actual)
	}
}

func TestV3Validations(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	goodInterface := tc.ServerInterfaceInfo{
		IPAddresses: []tc.ServerIPAddress{
			{
				Address:        "127.0.0.1/32",
				Gateway:        nil,
				ServiceAddress: true,
			},
		},
		MaxBandwidth: nil,
		Monitor:      true,
		MTU:          nil,
		Name:         "eth0",
	}

	testServer := tc.ServerV30{
		CommonServerProperties: tc.CommonServerProperties{
			CDNID:          util.IntPtr(1),
			HostName:       util.StrPtr("test"),
			DomainName:     util.StrPtr("quest"),
			PhysLocationID: new(int),
			ProfileID:      new(int),
			StatusID:       new(int),
			TypeID:         new(int),
			UpdPending:     new(bool),
			CachegroupID:   new(int),
		},
		Interfaces: []tc.ServerInterfaceInfo{goodInterface},
	}

	typeCols := []string{"name", "use_in_table"}
	cdnCols := []string{"cdn"}
	ipCols := []string{"id", "address"}
	typeRows := sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows := sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	ipRows := sqlmock.NewRows(ipCols)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT cdn").WillReturnRows(cdnRows)
	mock.ExpectQuery("SELECT s.ID").WillReturnRows(ipRows)

	tx := db.MustBegin().Tx

	_, err, _ = validateV3(&testServer, tx)
	if err != nil {
		t.Errorf("Unexpected error validating test server: %v", err)
	}

	testServer.Interfaces = []tc.ServerInterfaceInfo{}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server with no interfaces to be invalid")
	} else {
		t.Logf("Got expected error validating server with no interfaces: %v", err)
	}

	testServer.Interfaces = nil

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server with nil interfaces to be invalid")
	} else {
		t.Logf("Got expected error validating server with nil interfaces: %v", err)
	}

	badIface := goodInterface
	var badMTU uint64 = 1279
	badIface.MTU = &badMTU
	testServer.Interfaces = []tc.ServerInterfaceInfo{badIface}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server an MTU < 1280 to be invalid")
	} else {
		t.Logf("Got expected error validating server with an MTU < 1280: %v", err)
	}

	badIface.MTU = nil
	badIface.IPAddresses = []tc.ServerIPAddress{}
	testServer.Interfaces = []tc.ServerInterfaceInfo{badIface}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server with no IP addresses to be invalid")
	} else {
		t.Logf("Got expected error validating server with no IP addresses: %v", err)
	}

	badIface.IPAddresses = nil
	testServer.Interfaces = []tc.ServerInterfaceInfo{badIface}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server with nil IP addresses to be invalid")
	} else {
		t.Logf("Got expected error validating server with nil IP addresses: %v", err)
	}

	badIface = goodInterface
	badIP := tc.ServerIPAddress{
		Address:        "127.0.0.1/32",
		Gateway:        nil,
		ServiceAddress: false,
	}
	badIface.IPAddresses = []tc.ServerIPAddress{badIP}
	testServer.Interfaces = []tc.ServerInterfaceInfo{badIface}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server with no service addresses to be invalid")
	} else {
		t.Logf("Got expected error validating server with no service addresses: %v", err)
	}

	testServer.Interfaces = []tc.ServerInterfaceInfo{goodInterface, goodInterface}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server with too many interfaces with service addresses to be invalid")
	} else {
		t.Logf("Got expected error validating server with too many interfaces with service addresses: %v", err)
	}

	badIface = goodInterface
	badIface.IPAddresses = append(badIface.IPAddresses, tc.ServerIPAddress{
		Address:        "1.2.3.4/1",
		Gateway:        nil,
		ServiceAddress: true,
	})
	testServer.Interfaces = []tc.ServerInterfaceInfo{badIface}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err, _ = validateV3(&testServer, tx)
	if err == nil {
		t.Errorf("Expected a server with no service addresses to be invalid")
	} else {
		t.Logf("Got expected error validating server with no service addresses: %v", err)
	}
}

type SortableServers []tc.Server

func (s SortableServers) Len() int {
	return len(s)
}
func (s SortableServers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableServers) Less(i, j int) bool {
	return s[i].HostName < s[j].HostName
}

func TestUpdateStatusLastUpdatedTime(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	lastUpdated := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WithArgs(lastUpdated, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	sysErr, _, code := updateStatusLastUpdatedTime(1, &lastUpdated, db.MustBegin().Tx)
	if sysErr != nil {
		t.Errorf("unable to update time, system error: %v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("updated time failed with status code:%d", code)
	}
}

func TestCreateInterfaces(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testInterface := getTestServers()[0].Interface
	var iface []tc.ServerInterfaceInfoV40
	iface = append(iface, testInterface)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO interface").
		WithArgs(iface[0].MaxBandwidth, iface[0].Monitor, iface[0].MTU, iface[0].Name, 1, iface[0].RouterHostName, iface[0].RouterPortName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO ip_address").
		WithArgs(iface[0].IPAddresses[0].Address, iface[0].IPAddresses[0].Gateway, iface[0].Name, 1, iface[0].IPAddresses[0].ServiceAddress).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	usrErr, sysErr, code := createInterfaces(1, iface, db.MustBegin().Tx)
	if usrErr != nil {
		t.Errorf("unable to create interface, user error: %v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to create interface, system error: %v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("unable to create interface, failed with status code:%d", code)
	}
}

func TestDeleteInterfaces(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM ip_address").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM interface").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	usrErr, sysErr, code := deleteInterfaces(1, db.MustBegin().Tx)
	if usrErr != nil {
		t.Errorf("unable to delete interface, user error: %v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to delete interface, system error: %v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("unable to delete interface, failed with status code:%d", code)
	}
}

func TestInsertServerProfile(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	profileName := []string{"traffic_ops", "global"}
	priority := []int{0, 1}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").WithArgs(1, pq.Array(profileName), pq.Array(priority)).WillReturnResult(sqlmock.NewResult(2, 2))
	mock.ExpectCommit()

	usrErr, sysErr, code := insertServerProfile(1, profileName, db.MustBegin().Tx)
	if usrErr != nil {
		t.Errorf("unable to insert profile, user error: %v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to insert profile, system error: %v", sysErr)
	}
	if code != http.StatusOK {
		t.Errorf("unable to insert profile, failed with status code:%d", code)
	}
}

func TestCreateServerV4(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	s4 := getTestServers()[0].Server.Downgrade()

	mock.ExpectBegin()
	rows0 := sqlmock.NewRows([]string{"id"})
	rows0.AddRow(1)
	mock.ExpectQuery("SELECT").WithArgs(s4.ProfileNames[0]).WillReturnRows(rows0)

	rows := sqlmock.NewRows([]string{
		"cachegroup",
		"cachegroup_id",
		"cdn_id",
		"cdn_name",
		"domain_name",
		"guid",
		"host_name",
		"https_port",
		"id",
		"ilo_ip_address",
		"ilo_ip_gateway",
		"ilo_ip_netmask",
		"ilo_password",
		"ilo_username",
		"last_updated",
		"mgmt_ip_address",
		"mgmt_ip_gateway",
		"mgmt_ip_netmask",
		"offline_reason",
		"phys_location",
		"phys_location_id",
		"profile_name",
		"rack",
		"status",
		"status_id",
		"tcp_port",
		"server_type",
		"server_type_id",
	})
	rows.AddRow(
		*s4.Cachegroup,
		*s4.CachegroupID,
		*s4.CDNID,
		*s4.CDNName,
		*s4.DomainName,
		*s4.GUID,
		*s4.HostName,
		*s4.HTTPSPort,
		*s4.ID,
		*s4.ILOIPAddress,
		*s4.ILOIPGateway,
		*s4.ILOIPNetmask,
		*s4.ILOPassword,
		*s4.ILOUsername,
		*s4.LastUpdated,
		*s4.MgmtIPAddress,
		*s4.MgmtIPGateway,
		*s4.MgmtIPNetmask,
		*s4.OfflineReason,
		*s4.PhysLocation,
		*s4.PhysLocationID,
		fmt.Sprintf("{%s}", strings.Join(s4.ProfileNames, ",")),
		*s4.Rack,
		*s4.Status,
		*s4.StatusID,
		*s4.TCPPort,
		s4.Type,
		*s4.TypeID,
	)
	mock.ExpectQuery("INSERT INTO server").
		WithArgs(*s4.CachegroupID, *s4.CDNID, *s4.DomainName, *s4.HostName, *s4.HTTPSPort, *s4.ILOIPAddress,
			*s4.ILOIPNetmask, *s4.ILOIPGateway, *s4.ILOUsername, *s4.ILOPassword, *s4.MgmtIPAddress,
			*s4.MgmtIPNetmask, *s4.MgmtIPGateway, *s4.OfflineReason, *s4.PhysLocationID, 1, *s4.Rack,
			*s4.StatusID, *s4.TCPPort, *s4.TypeID, *s4.XMPPID, *s4.XMPPPasswd).
		WillReturnRows(rows)
	mock.ExpectCommit()

	sid, err := createServerV4(db.MustBegin(), s4)
	if err != nil {
		t.Errorf("unable to create v4 server, error: %v", err)
	}
	if sid != int64(*s4.ID) {
		t.Errorf("mismatched server ID, expected: %d, got: %d", *s4.ID, sid)
	}
}

func TestCreateServerV3(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	s3 := tc.ServerV30{
		CommonServerProperties: tc.CommonServerProperties{
			Cachegroup:       util.Ptr("mb"),
			CachegroupID:     util.Ptr(1),
			CDNID:            util.Ptr(1),
			CDNName:          util.Ptr("ALL"),
			DeliveryServices: nil,
			DomainName:       util.Ptr(""),
			FQDN:             nil,
			FqdnTime:         time.Time{},
			GUID:             util.Ptr(""),
			HostName:         util.Ptr("test"),
			HTTPSPort:        util.Ptr(8443),
			ID:               util.Ptr(1),
			ILOIPAddress:     util.Ptr(""),
			ILOIPGateway:     util.Ptr(""),
			ILOIPNetmask:     util.Ptr(""),
			ILOPassword:      util.Ptr(""),
			ILOUsername:      util.Ptr(""),
			LastUpdated:      util.Ptr(tc.TimeNoMod{}),
			MgmtIPAddress:    util.Ptr(""),
			MgmtIPGateway:    util.Ptr(""),
			MgmtIPNetmask:    util.Ptr(""),
			OfflineReason:    util.Ptr(""),
			PhysLocation:     util.Ptr("boulder"),
			PhysLocationID:   util.Ptr(1),
			Profile:          util.Ptr("GLOBAL"),
			ProfileDesc:      nil,
			ProfileID:        util.Ptr(1),
			Rack:             util.Ptr(""),
			RevalPending:     util.Ptr(false),
			Status:           util.Ptr("ACTIVE"),
			StatusID:         util.Ptr(1),
			TCPPort:          util.Ptr(8080),
			Type:             "EDGE",
			TypeID:           util.Ptr(2),
			UpdPending:       util.Ptr(false),
			XMPPID:           util.Ptr(""),
			XMPPPasswd:       util.Ptr(""),
		},
		RouterHostName:    util.Ptr(""),
		RouterPortName:    util.Ptr(""),
		Interfaces:        []tc.ServerInterfaceInfo{},
		StatusLastUpdated: util.Ptr(time.Time{}),
	}

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(1)
	mock.ExpectQuery("INSERT INTO server").
		WithArgs(*s3.CachegroupID, *s3.CDNID, *s3.DomainName, *s3.HostName, *s3.HTTPSPort, *s3.ILOIPAddress,
			*s3.ILOIPNetmask, *s3.ILOIPGateway, *s3.ILOUsername, *s3.ILOPassword, *s3.MgmtIPAddress,
			*s3.MgmtIPNetmask, *s3.MgmtIPGateway, *s3.OfflineReason, *s3.PhysLocationID, *s3.ProfileID,
			*s3.Rack, *s3.StatusID, *s3.TCPPort, *s3.TypeID, *s3.XMPPID, *s3.XMPPPasswd, *s3.StatusLastUpdated).
		WillReturnRows(rows)
	mock.ExpectCommit()

	sid, err := createServerV3(db.MustBegin(), s3)
	if err != nil {
		t.Errorf("unable to create v3 server, error: %v", err)
	}
	if sid != int64(*s3.ID) {
		t.Errorf("mismatched server ID, expected: %d, got: %d", *s3.ID, sid)
	}
}

func TestUpdateServer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	server := getTestServers()[0].Server

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{
		"cachegroup",
		"cachegroup_id",
		"cdn_id",
		"cdn_name",
		"domain_name",
		"guid",
		"host_name",
		"https_port",
		"id",
		"ilo_ip_address",
		"ilo_ip_gateway",
		"ilo_ip_netmask",
		"ilo_password",
		"ilo_username",
		"last_updated",
		"mgmt_ip_address",
		"mgmt_ip_gateway",
		"mgmt_ip_netmask",
		"offline_reason",
		"phys_location",
		"phys_location_id",
		"profile_name",
		"rack",
		"status",
		"status_id",
		"tcp_port",
		"server_type",
		"server_type_id",
		"status_last_updated",
	})
	rows.AddRow(
		server.CacheGroup,
		server.CacheGroupID,
		server.CDNID,
		server.CDN,
		server.DomainName,
		*server.GUID,
		server.HostName,
		*server.HTTPSPort,
		server.ID,
		*server.ILOIPAddress,
		*server.ILOIPGateway,
		*server.ILOIPNetmask,
		*server.ILOPassword,
		*server.ILOUsername,
		server.LastUpdated,
		*server.MgmtIPAddress,
		*server.MgmtIPGateway,
		*server.MgmtIPNetmask,
		*server.OfflineReason,
		server.PhysicalLocation,
		server.PhysicalLocationID,
		fmt.Sprintf("{%s}", strings.Join(server.Profiles, ",")),
		*server.Rack,
		server.Status,
		server.StatusID,
		*server.TCPPort,
		server.Type,
		server.TypeID,
		*server.StatusLastUpdated,
	)
	mock.ExpectQuery("UPDATE server SET").
		WithArgs(server.CacheGroupID, server.CDNID, server.DomainName, server.HostName, *server.HTTPSPort, *server.ILOIPAddress,
			*server.ILOIPNetmask, *server.ILOIPGateway, *server.ILOUsername, *server.ILOPassword, *server.MgmtIPAddress,
			*server.MgmtIPNetmask, *server.MgmtIPGateway, *server.OfflineReason, server.PhysicalLocationID, 1, *server.Rack,
			server.StatusID, *server.TCPPort, server.TypeID, *server.XMPPPasswd, *server.StatusLastUpdated, server.ID).
		WillReturnRows(rows)
	mock.ExpectCommit()

	sid, code, usrErr, sysErr := updateServer(db.MustBegin(), server)
	if usrErr != nil {
		t.Errorf("unable to update v4 server, user error: %v", usrErr)
	}
	if sysErr != nil {
		t.Errorf("unable to update v4 server, system error: %v", sysErr)
	}
	if sid != int64(server.ID) {
		t.Errorf("updated incorrect server, expected: %d, got: %d", server.ID, sid)
	}
	if code != http.StatusOK {
		t.Errorf("failed to update server with id: %d, expected: %d, got: %d", server.ID, http.StatusOK, code)
	}
}
