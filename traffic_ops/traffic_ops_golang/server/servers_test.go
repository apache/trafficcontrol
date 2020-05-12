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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type ServerAndInterfaces struct {
	Server tc.Server
	Interfaces []tc.ServerInterfaceInfo
}

func getTestServers() []ServerAndInterfaces {
	servers := []ServerAndInterfaces{}
	testServer := tc.Server{
		Cachegroup:     "Cachegroup",
		CachegroupID:   1,
		CDNID:          1,
		CDNName:        "cdnName",
		DomainName:     "domainName",
		GUID:           "guid",
		HostName:       "server1",
		HTTPSPort:      443,
		ID:             1,
		ILOIPAddress:   "iloIpAddress",
		ILOIPGateway:   "iloIpGateway",
		ILOIPNetmask:   "iloIpNetmask",
		ILOPassword:    "iloPassword",
		ILOUsername:    "iloUsername",
		InterfaceMtu:   9500,
		InterfaceName:  "interfaceName",
		IP6Address:     "ip6Address",
		IP6IsService:   false,
		IP6Gateway:     "ip6Gateway",
		IPAddress:      "ipAddress",
		IPIsService:    true,
		IPGateway:      "ipGateway",
		IPNetmask:      "ipNetmask",
		LastUpdated:    tc.TimeNoMod{Time: time.Now()},
		MgmtIPAddress:  "mgmtIpAddress",
		MgmtIPGateway:  "mgmtIpGateway",
		MgmtIPNetmask:  "mgmtIpNetmask",
		OfflineReason:  "offlineReason",
		PhysLocation:   "physLocation",
		PhysLocationID: 1,
		Profile:        "profile",
		ProfileDesc:    "profileDesc",
		ProfileID:      1,
		Rack:           "rack",
		RevalPending:   true,
		RouterHostName: "routerHostName",
		RouterPortName: "routerPortName",
		Status:         "status",
		StatusID:       1,
		TCPPort:        80,
		Type:           "EDGE",
		TypeID:         1,
		UpdPending:     true,
		XMPPID:         "xmppId",
		XMPPPasswd:     "xmppPasswd",
	}

	mtu := uint64(testServer.InterfaceMtu)

	interfaces := []tc.ServerInterfaceInfo{
		tc.ServerInterfaceInfo{
			IPAddresses: []tc.ServerIpAddress{
				tc.ServerIpAddress{
					Address: testServer.IPAddress,
					Gateway: &testServer.IPGateway,
					ServiceAddress: testServer.IPIsService,
				},
			},
			MaxBandwidth: nil,
			Monitor: true,
			MTU: &mtu,
			Name: testServer.InterfaceName,
		},
	}
	servers = append(servers, ServerAndInterfaces{Server: testServer, Interfaces: interfaces})

	testServer2 := testServer
	testServer2.Cachegroup = "cachegroup2"
	testServer2.HostName = "server2"
	testServer2.ID = 2
	servers = append(servers, ServerAndInterfaces{Server: testServer2, Interfaces: interfaces})

	testServer3 := testServer
	testServer3.Cachegroup = "cachegroup3"
	testServer3.HostName = "server3"
	testServer3.ID = 3
	servers = append(servers, ServerAndInterfaces{Server: testServer3, Interfaces: interfaces})

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
	rows.AddRow(testServers[0].TypeID, 5)
	// Make it return a list of atleast one associated ds
	dsrows := sqlmock.NewRows([]string{"array"})
	dsrows.AddRow("{3}")
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectQuery("SELECT ARRAY").WillReturnRows(dsrows)

	snv := tc.ServerNullableV11{
		CDNID:    &testServers[0].CDNID,
		FqdnTime: time.Time{},
		TypeID:   &testServers[0].TypeID,
	}
	sn := tc.ServerNullable{
		ServerNullableV11: snv,
		IPIsService:       nil,
		IP6IsService:      nil,
	}

	s := &TOServer{
		APIInfoImpl: api.APIInfoImpl{
			ReqInfo: &api.APIInfo{
				Params:    nil,
				IntParams: nil,
				User:      nil,
				ReqID:     0,
				Version:   nil,
				Tx:        db.MustBegin(),
				Config:    nil,
			},
		},
		ServerNullable: sn,
	}

	userErr, _, errCode := s.Update()
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

	cols := test.ColsFromStructByTag("db", tc.CommonServerProperties{})
	interfaceCols := []string{"interfaces", "id"}
	rows := sqlmock.NewRows(cols)
	interfaceRows := sqlmock.NewRows(interfaceCols)

	//TODO: drichardson - build helper to add these Rows from the struct values
	//                    or by CSV if types get in the way
	for _, srv := range testServers {
		ts := srv.Server
		rows = rows.AddRow(
			ts.Cachegroup,
			ts.CachegroupID,
			ts.CDNID,
			ts.CDNName,
			ts.DomainName,
			ts.GUID,
			ts.HostName,
			ts.HTTPSPort,
			ts.ID,
			ts.ILOIPAddress,
			ts.ILOIPGateway,
			ts.ILOIPNetmask,
			ts.ILOPassword,
			ts.ILOUsername,
			ts.LastUpdated,
			ts.MgmtIPAddress,
			ts.MgmtIPGateway,
			ts.MgmtIPNetmask,
			ts.OfflineReason,
			ts.PhysLocation,
			ts.PhysLocationID,
			ts.Profile,
			ts.ProfileDesc,
			ts.ProfileID,
			ts.Rack,
			ts.RevalPending,
			ts.RouterHostName,
			ts.RouterPortName,
			ts.Status,
			ts.StatusID,
			ts.TCPPort,
			ts.Type,
			ts.TypeID,
			ts.UpdPending,
			ts.XMPPID,
			ts.XMPPPasswd,
		)
		interfaceRows = interfaceRows.AddRow(
			srv.Interfaces,
			ts.ID,
		)
	}


	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(server.id\\) FROM server").WillReturnRows(unfilteredRows)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectQuery("SELECT").WillReturnRows(interfaceRows)

	v := map[string]string{"cachegroup": "2"}

	user := auth.CurrentUser{}

	servers, _, userErr, sysErr, errCode := getServers(v, db.MustBegin(), &user)
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

	testServers[1].Server.Cachegroup = "parentCacheGroup"
	testServers[1].Server.CachegroupID = 2
	testServers[1].Server.Type = "MID"

	unfilteredCols := []string{"count"}
	unfilteredRows := sqlmock.NewRows(unfilteredCols).AddRow(len(testServers))

	cols := test.ColsFromStructByTag("db", tc.CommonServerProperties{})
	interfaceCols := []string{"interfaces", "id"}
	rows := sqlmock.NewRows(cols)
	interfaceRows := sqlmock.NewRows(interfaceCols)

	for _, srv := range testServers {
		ts := srv.Server
		rows = rows.AddRow(
			ts.Cachegroup,
			ts.CachegroupID,
			ts.CDNID,
			ts.CDNName,
			ts.DomainName,
			ts.GUID,
			ts.HostName,
			ts.HTTPSPort,
			ts.ID,
			ts.ILOIPAddress,
			ts.ILOIPGateway,
			ts.ILOIPNetmask,
			ts.ILOPassword,
			ts.ILOUsername,
			ts.InterfaceMtu,
			ts.InterfaceName,
			ts.IP6Address,
			ts.IP6IsService,
			ts.IP6Gateway,
			ts.IPAddress,
			ts.IPIsService,
			ts.IPNetmask,
			ts.IPGateway,
			ts.LastUpdated,
			ts.MgmtIPAddress,
			ts.MgmtIPGateway,
			ts.MgmtIPNetmask,
			ts.OfflineReason,
			ts.PhysLocation,
			ts.PhysLocationID,
			ts.Profile,
			ts.ProfileDesc,
			ts.ProfileID,
			ts.Rack,
			ts.RevalPending,
			ts.RouterHostName,
			ts.RouterPortName,
			ts.Status,
			ts.StatusID,
			ts.TCPPort,
			ts.Type,
			ts.TypeID,
			ts.UpdPending,
			ts.XMPPID,
			ts.XMPPPasswd,
		)
		interfaceRows = interfaceRows.AddRow(
			srv.Interfaces,
			ts.ID,
		)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(server.id\\) FROM server").WillReturnRows(unfilteredRows)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectQuery("SELECT").WillReturnRows(interfaceRows)
	v := map[string]string{}

	user := auth.CurrentUser{}

	servers, _, userErr, sysErr, errCode := getServers(v, db.MustBegin(), &user)

	if userErr != nil || sysErr != nil {
		t.Errorf("getServers expected: no errors, actual: %v %v with status: %s", userErr, sysErr, http.StatusText(errCode))
	}

	cols2 := test.ColsFromStructByTag("db", tc.Server{})
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

	ts := servers[1]
	rows2 = rows2.AddRow(
		ts.Cachegroup,
		ts.CachegroupID,
		ts.CDNID,
		ts.CDNName,
		ts.DomainName,
		ts.GUID,
		ts.HostName,
		ts.HTTPSPort,
		ts.ID,
		ts.ILOIPAddress,
		ts.ILOIPGateway,
		ts.ILOIPNetmask,
		ts.ILOPassword,
		ts.ILOUsername,
		ts.LastUpdated,
		ts.MgmtIPAddress,
		ts.MgmtIPGateway,
		ts.MgmtIPNetmask,
		ts.OfflineReason,
		ts.PhysLocation,
		ts.PhysLocationID,
		ts.Profile,
		ts.ProfileDesc,
		ts.ProfileID,
		ts.Rack,
		ts.RevalPending,
		ts.RouterHostName,
		ts.RouterPortName,
		ts.Status,
		ts.StatusID,
		ts.TCPPort,
		ts.Type,
		ts.TypeID,
		ts.UpdPending,
		ts.XMPPID,
		ts.XMPPPasswd,
	)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows2)
	mid, userErr, sysErr, errCode := getMidServers(servers, db.MustBegin())

	if userErr != nil || sysErr != nil {
		t.Fatalf("getMidServers expected: no errors, actual: %v %v with status: %s", userErr, sysErr, http.StatusText(errCode))
	}
	if len(mid) != 1 {
		t.Fatalf("getMidServers expected: len(mid) == 1, actual: %v", len(mid))
	}
	if mid[0].Type != "MID" || *(mid[0].CachegroupID) != 2 || *(mid[0].Cachegroup) != "parentCacheGroup" {
		t.Fatalf("getMidServers expected: Type == MID, actual: %v", mid[0].Type)
		t.Fatalf("getMidServers expected: CachegroupID == 2, actual: %v", *(mid[0].CachegroupID))
		t.Fatalf("getMidServers expected: Cachegroup == parentCacheGroup, actual: %v", *(mid[0].Cachegroup))
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
