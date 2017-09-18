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
	"net/url"
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/tostructs"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestServers() []tostructs.Server {
	servers := []tostructs.Server{}
	testServer := tostructs.Server{
		Cachegroup:     "Cachegroup",
		CachegroupID:   1,
		CdnID:          1,
		CdnName:        "cdnName",
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
		IP6Gateway:     "ip6Gateway",
		IPAddress:      "ipAddress",
		IPGateway:      "ipGateway",
		IPNetmask:      "ipNetmask",
		LastUpdated:    "lastUpdated",
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
		ServerType:     "EDGE",
		ServerTypeID:   1,
		UpdPending:     true,
		XMPPID:         "xmppId",
		XMPPPasswd:     "xmppPasswd",
	}
	servers = append(servers, testServer)

	testServer2 := testServer
	testServer2.Cachegroup = "cachegroup2"
	testServer2.HostName = "server2"
	servers = append(servers, testServer2)

	return servers
}

func TestGetServersByDsId(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testServers := getTestServers()
	cols := ColsFromStructByTag("db", tostructs.Server{})
	rows := sqlmock.NewRows(cols)

	//TODO: drichardson - build helper to add these Rows from the struct values
	//                    or by CSV if types get in the way
	for _, ts := range testServers {
		rows = rows.AddRow(
			ts.Cachegroup,
			ts.CachegroupID,
			ts.CdnID,
			ts.CdnName,
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
			ts.IP6Gateway,
			ts.IPAddress,
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
			ts.ServerType,
			ts.ServerTypeID,
			ts.UpdPending,
			ts.XMPPID,
			ts.XMPPPasswd,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := url.Values{}
	v.Set("dsId", "1")

	servers, err := getServers(v, db, PrivLevelAdmin)
	if err != nil {
		t.Errorf("getServers expected: nil error, actual: %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("getServers expected: len(servers) == 1, actual: %v", len(servers))
	}

}

type SortableServers []tostructs.Server

func (s SortableServers) Len() int {
	return len(s)
}
func (s SortableServers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableServers) Less(i, j int) bool {
	return s[i].HostName < s[j].HostName
}
