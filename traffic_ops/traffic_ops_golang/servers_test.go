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
	"sort"
	"testing"

	. "github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tcstructs"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestServers() []Server {
	servers := []Server{}
	testServer := Server{
		Cachegroup:     "Cachegroup",
		CachegroupId:   1,
		CdnId:          1,
		CdnName:        "cdnName",
		DomainName:     "domainName",
		Guid:           "guid",
		HostName:       "server1",
		HttpsPort:      443,
		Id:             1,
		IloIpAddress:   "iloIpAddress",
		IloIpGateway:   "iloIpGateway",
		IloIpNetmask:   "iloIpNetmask",
		IloPassword:    "iloPassword",
		IloUsername:    "iloUsername",
		InterfaceMtu:   "interfaceMtu",
		InterfaceName:  "interfaceName",
		Ip6Address:     "ip6Address",
		Ip6Gateway:     "ip6Gateway",
		IpAddress:      "ipAddress",
		IpGateway:      "ipGateway",
		IpNetmask:      "ipNetmask",
		LastUpdated:    "lastUpdated",
		MgmtIpAddress:  "mgmtIpAddress",
		MgmtIpGateway:  "mgmtIpGateway",
		MgmtIpNetmask:  "mgmtIpNetmask",
		OfflineReason:  "offlineReason",
		PhysLocation:   "physLocation",
		PhysLocationId: 1,
		Profile:        "profile",
		ProfileDesc:    "profileDesc",
		ProfileId:      1,
		Rack:           "rack",
		RouterHostName: "routerHostName",
		RouterPortName: "routerPortName",
		Status:         "status",
		StatusId:       1,
		TcpPort:        80,
		ServerType:     "EDGE",
		ServerTypeId:   1,
		UpdPending:     true,
		XmppId:         "xmppId",
		XmppPasswd:     "xmppPasswd",
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
	cols := ColsFromStructByTag("db", Server{})
	rows := sqlmock.NewRows(cols)

	//TODO: drichardson - build helper to add these Rows from the struct values
	//                    or by CSV if types get in the way
	for _, ts := range testServers {
		rows = rows.AddRow(
			ts.Cachegroup,
			ts.CachegroupId,
			ts.CdnId,
			ts.CdnName,
			ts.DomainName,
			ts.Guid,
			ts.HostName,
			ts.HttpsPort,
			ts.Id,
			ts.IloIpAddress,
			ts.IloIpGateway,
			ts.IloIpNetmask,
			ts.IloPassword,
			ts.IloUsername,
			ts.InterfaceMtu,
			ts.InterfaceName,
			ts.Ip6Address,
			ts.Ip6Gateway,
			ts.IpAddress,
			ts.IpNetmask,
			ts.IpGateway,
			ts.LastUpdated,
			ts.MgmtIpAddress,
			ts.MgmtIpGateway,
			ts.MgmtIpNetmask,
			ts.OfflineReason,
			ts.PhysLocation,
			ts.PhysLocationId,
			ts.Profile,
			ts.ProfileDesc,
			ts.ProfileId,
			ts.Rack,
			ts.RouterHostName,
			ts.RouterPortName,
			ts.Status,
			ts.StatusId,
			ts.TcpPort,
			ts.ServerType,
			ts.ServerTypeId,
			ts.UpdPending,
			ts.XmppId,
			ts.XmppPasswd,
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

type SortableServers []Server

func (s SortableServers) Len() int {
	return len(s)
}
func (s SortableServers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortableServers) Less(i, j int) bool {
	return s[i].HostName < s[j].HostName
}

func sortServers(p []Server) []Server {
	sort.Sort(SortableServers(p))
	return p
}
