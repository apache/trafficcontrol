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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type ServerAndInterfaces struct {
	Server    tc.ServerV40
	Interface tc.ServerInterfaceInfoV40
}

func getTestServers() []ServerAndInterfaces {
	servers := []ServerAndInterfaces{}
	testServer := tc.ServerV40{
		CommonServerPropertiesV40: tc.CommonServerPropertiesV40{
			Cachegroup:     util.StrPtr("Cachegroup"),
			CachegroupID:   util.IntPtr(1),
			CDNID:          util.IntPtr(1),
			CDNName:        util.StrPtr("cdnName"),
			DomainName:     util.StrPtr("domainName"),
			GUID:           util.StrPtr("guid"),
			HostName:       util.StrPtr("server1"),
			HTTPSPort:      util.IntPtr(443),
			ID:             util.IntPtr(1),
			ILOIPAddress:   util.StrPtr("iloIpAddress"),
			ILOIPGateway:   util.StrPtr("iloIpGateway"),
			ILOIPNetmask:   util.StrPtr("iloIpNetmask"),
			ILOPassword:    util.StrPtr("iloPassword"),
			ILOUsername:    util.StrPtr("iloUsername"),
			LastUpdated:    &tc.TimeNoMod{Time: time.Now()},
			MgmtIPAddress:  util.StrPtr("mgmtIpAddress"),
			MgmtIPGateway:  util.StrPtr("mgmtIpGateway"),
			MgmtIPNetmask:  util.StrPtr("mgmtIpNetmask"),
			OfflineReason:  util.StrPtr("offlineReason"),
			PhysLocation:   util.StrPtr("physLocation"),
			PhysLocationID: util.IntPtr(1),
			ProfileNames:   []string{"profile"},
			Rack:           util.StrPtr("rack"),
			RevalPending:   util.BoolPtr(true),
			Status:         util.StrPtr("status"),
			StatusID:       util.IntPtr(1),
			TCPPort:        util.IntPtr(80),
			Type:           "EDGE",
			TypeID:         util.IntPtr(1),
			UpdPending:     util.BoolPtr(true),
			XMPPID:         util.StrPtr("xmppId"),
			XMPPPasswd:     util.StrPtr("xmppPasswd"),
		},
		StatusLastUpdated: &(time.Time{}),
		ConfigUpdateTime:  &(time.Time{}),
		ConfigApplyTime:   &(time.Time{}),
		RevalUpdateTime:   &(time.Time{}),
		RevalApplyTime:    &(time.Time{}),
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

	servers = append(servers, ServerAndInterfaces{Server: testServer, Interface: iface})

	testServer2 := testServer
	testServer2.Cachegroup = util.StrPtr("cachegroup2")
	testServer2.HostName = util.StrPtr("server2")
	testServer2.ID = util.IntPtr(2)
	servers = append(servers, ServerAndInterfaces{Server: testServer2, Interface: iface})

	testServer3 := testServer
	testServer3.Cachegroup = util.StrPtr("cachegroup3")
	testServer3.HostName = util.StrPtr("server3")
	testServer3.ID = util.IntPtr(3)
	servers = append(servers, ServerAndInterfaces{Server: testServer3, Interface: iface})

	return servers
}

// Test to make sure that updating the "cdn" of a server already assigned to a DS fails
func TestUpdateServer(t *testing.T) {
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
	rows.AddRow(*testServers[0].Server.TypeID, 5)
	// Make it return a list of atleast one associated ds
	dsrows := sqlmock.NewRows([]string{"array"})
	dsrows.AddRow("{3}")
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectQuery("SELECT ARRAY").WillReturnRows(dsrows)

	s := tc.CommonServerPropertiesV40{
		CDNID:    testServers[0].Server.CDNID,
		FqdnTime: time.Time{},
		TypeID:   testServers[0].Server.TypeID,
		ID:       testServers[0].Server.ID,
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

	cols := test.ColsFromStructByTag("db", tc.CommonServerPropertiesV40{})
	cols = append(cols, "status_last_updated", "config_update_time", "config_apply_time", "revalidate_update_time", "revalidate_apply_time")
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
			*ts.Cachegroup,
			*ts.CachegroupID,
			*ts.CDNID,
			*ts.CDNName,
			*ts.DomainName,
			*ts.GUID,
			*ts.HostName,
			*ts.HTTPSPort,
			*ts.ID,
			*ts.ILOIPAddress,
			*ts.ILOIPGateway,
			*ts.ILOIPNetmask,
			*ts.ILOPassword,
			*ts.ILOUsername,
			*ts.LastUpdated,
			*ts.MgmtIPAddress,
			*ts.MgmtIPGateway,
			*ts.MgmtIPNetmask,
			*ts.OfflineReason,
			*ts.PhysLocation,
			*ts.PhysLocationID,
			fmt.Sprintf("{%s}", strings.Join(ts.ProfileNames, ",")),
			*ts.Rack,
			*ts.RevalPending,
			*ts.RevalUpdateTime,
			*ts.RevalApplyTime,
			*ts.Status,
			*ts.StatusID,
			*ts.TCPPort,
			ts.Type,
			*ts.TypeID,
			*ts.UpdPending,
			*ts.ConfigUpdateTime,
			*ts.ConfigApplyTime,
			*ts.XMPPID,
			*ts.XMPPPasswd,
			*ts.StatusLastUpdated,
		)
		interfaceRows = interfaceRows.AddRow(
			srv.Interface.MaxBandwidth,
			srv.Interface.Monitor,
			srv.Interface.MTU,
			srv.Interface.Name,
			*ts.ID,
			srv.Interface.RouterHostName,
			srv.Interface.RouterPortName,
		)

		for _, ip := range srv.Interface.IPAddresses {
			ipRows = ipRows.AddRow(
				ip.Address,
				ip.Gateway,
				ip.ServiceAddress,
				*ts.ID,
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

	servers, _, userErr, sysErr, errCode, _ := getServers(nil, v, db.MustBegin(), &user, false, version)
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

	testServers[1].Server.Cachegroup = util.StrPtr("parentCacheGroup")
	testServers[1].Server.CachegroupID = util.IntPtr(2)
	testServers[1].Server.Type = "MID"

	unfilteredCols := []string{"count"}
	unfilteredRows := sqlmock.NewRows(unfilteredCols).AddRow(len(testServers))

	cols := test.ColsFromStructByTag("db", tc.CommonServerPropertiesV40{})
	cols = append(cols, "status_last_updated", "config_update_time", "config_apply_time", "revalidate_update_time", "revalidate_apply_time")
	interfaceCols := []string{"max_bandwidth", "monitor", "mtu", "name", "server", "router_host_name", "router_port_name"}
	rows := sqlmock.NewRows(cols)
	interfaceRows := sqlmock.NewRows(interfaceCols)

	ipCols := []string{"address", "gateway", "service_address", "server", "interface"}
	ipRows := sqlmock.NewRows(ipCols)

	for _, srv := range testServers {
		ts := srv.Server
		rows = rows.AddRow(
			*ts.Cachegroup,
			*ts.CachegroupID,
			*ts.CDNID,
			*ts.CDNName,
			*ts.DomainName,
			*ts.GUID,
			*ts.HostName,
			*ts.HTTPSPort,
			*ts.ID,
			*ts.ILOIPAddress,
			*ts.ILOIPGateway,
			*ts.ILOIPNetmask,
			*ts.ILOPassword,
			*ts.ILOUsername,
			*ts.LastUpdated,
			*ts.MgmtIPAddress,
			*ts.MgmtIPGateway,
			*ts.MgmtIPNetmask,
			*ts.OfflineReason,
			*ts.PhysLocation,
			*ts.PhysLocationID,
			fmt.Sprintf("{%s}", strings.Join(ts.ProfileNames, ",")),
			*ts.Rack,
			*ts.RevalPending,
			*ts.RevalUpdateTime,
			*ts.RevalApplyTime,
			*ts.Status,
			*ts.StatusID,
			*ts.TCPPort,
			ts.Type,
			*ts.TypeID,
			*ts.UpdPending,
			*ts.ConfigUpdateTime,
			*ts.ConfigApplyTime,
			*ts.XMPPID,
			*ts.XMPPPasswd,
			*ts.StatusLastUpdated,
		)
		interfaceRows = interfaceRows.AddRow(
			srv.Interface.MaxBandwidth,
			srv.Interface.Monitor,
			srv.Interface.MTU,
			srv.Interface.Name,
			*ts.ID,
			srv.Interface.RouterHostName,
			srv.Interface.RouterPortName,
		)

		for _, ip := range srv.Interface.IPAddresses {
			ipRows = ipRows.AddRow(
				ip.Address,
				ip.Gateway,
				ip.ServiceAddress,
				*ts.ID,
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
	servers, _, userErr, sysErr, errCode, _ := getServers(nil, v, db.MustBegin(), &user, false, version)

	if userErr != nil || sysErr != nil {
		t.Errorf("getServers expected: no errors, actual: %v %v with status: %s", userErr, sysErr, http.StatusText(errCode))
	}

	cols2 := test.ColsFromStructByTag("db", tc.CommonServerPropertiesV40{})
	cols2 = append(cols2, "status_last_updated", "config_update_time", "config_apply_time", "revalidate_update_time", "revalidate_apply_time")
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

	serverMap := make(map[int]tc.ServerV40, len(servers))
	serverIDs := make([]int, 0, len(servers))
	for _, server := range servers {
		if server.ID == nil {
			t.Fatal("Found server with nil ID")
		}
		serverIDs = append(serverIDs, *server.ID)
		serverMap[*server.ID] = server
	}

	var ts tc.ServerV40
	for _, s := range servers {
		if s.HostName != nil && *s.HostName == "server2" {
			ts = s
			break
		}
	}
	*ts.ID = *ts.ID + 1
	rows2 = rows2.AddRow(
		*ts.Cachegroup,
		*ts.CachegroupID,
		*ts.CDNID,
		*ts.CDNName,
		*ts.DomainName,
		*ts.GUID,
		*ts.HostName,
		*ts.HTTPSPort,
		*ts.ID,
		*ts.ILOIPAddress,
		*ts.ILOIPGateway,
		*ts.ILOIPNetmask,
		*ts.ILOPassword,
		*ts.ILOUsername,
		*ts.LastUpdated,
		*ts.MgmtIPAddress,
		*ts.MgmtIPGateway,
		*ts.MgmtIPNetmask,
		*ts.OfflineReason,
		*ts.PhysLocation,
		*ts.PhysLocationID,
		fmt.Sprintf("{%s}", strings.Join(ts.ProfileNames, ",")),
		*ts.Rack,
		*ts.RevalPending,
		*ts.RevalUpdateTime,
		*ts.RevalApplyTime,
		*ts.Status,
		*ts.StatusID,
		*ts.TCPPort,
		ts.Type,
		*ts.TypeID,
		*ts.UpdPending,
		*ts.ConfigUpdateTime,
		*ts.ConfigApplyTime,
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

	if serverMap[mid[0]].Cachegroup == nil {
		t.Error("getMidServers expected: Cachegroup == parentCacheGroup, actual: nil")
	} else if *(serverMap[mid[0]].Cachegroup) != "parentCacheGroup" {
		t.Errorf("getMidServers expected: Cachegroup == parentCacheGroup, actual: %v", *(serverMap[mid[0]].Cachegroup))
	}

	if serverMap[mid[0]].CachegroupID == nil {
		t.Error("getMidServers expected: CachegroupID == 2, actual: nil")
	} else if *(serverMap[mid[0]].CachegroupID) != 2 {
		t.Errorf("getMidServers expected: CachegroupID == 2, actual: %v", *(serverMap[mid[0]].CachegroupID))
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

	_, err = validateV3(&testServer, tx)
	if err != nil {
		t.Errorf("Unexpected error validating test server: %v", err)
	}

	testServer.Interfaces = []tc.ServerInterfaceInfo{}

	typeRows = sqlmock.NewRows(typeCols).AddRow("EDGE", "server")
	cdnRows = sqlmock.NewRows(cdnCols).AddRow(*testServer.CDNID)
	mock.ExpectQuery("SELECT name, use_in_table").WillReturnRows(typeRows)
	mock.ExpectQuery("SELECT").WillReturnRows(cdnRows)

	_, err = validateV3(&testServer, tx)
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

	_, err = validateV3(&testServer, tx)
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

	_, err = validateV3(&testServer, tx)
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

	_, err = validateV3(&testServer, tx)
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

	_, err = validateV3(&testServer, tx)
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

	_, err = validateV3(&testServer, tx)
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

	_, err = validateV3(&testServer, tx)
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

	_, err = validateV3(&testServer, tx)
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
