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

	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setupServer() Server {
	return Server{
		Cachegroup:     "Cachegroup",
		CachegroupId:   1,
		CdnId:          1,
		CdnName:        "cdnName",
		DomainName:     "domainName",
		Guid:           "guid",
		HostName:       "hostName",
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
}

func TestGetServers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testServer := setupServer()

	rows := sqlmock.NewRows([]string{
		"cachegroup",
		"cachegroupId",
		"cdnId",
		"cdnName",
		"domainName",
		"guid",
		"hostName",
		"httpsPort",
		"id",
		"iloIpAddress",
		"iloIpGateway",
		"iloIpNetMask",
		"iloPassword",
		"iloUsername",
		"interfaceMtu",
		"interfaceName",
		"ip6Address",
		"ip6Gateway",
		"ipAddress",
		"ipGateway",
		"ipNetmask",
		"lastUpdated",
		"mgmtIpAddress",
		"mgmtIpGateway",
		"mgmtIpNetmask",
		"offlineReason",
		"physLocation",
		"physLocationId",
		"profile",
		"profileDesc",
		"profileId",
		"rack",
		"routerHostName",
		"routerPortName",
		"status",
		"statusId",
		"tcpPort",
		"type",
		"typeId",
		"updPending",
		"xmppId",
		"xmppPassword",
	})
	rows = rows.AddRow(
		testServer.Cachegroup,
		testServer.CachegroupId,
		testServer.CdnId,
		testServer.CdnName,
		testServer.DomainName,
		testServer.Guid,
		testServer.HostName,
		testServer.HttpsPort,
		testServer.Id,
		testServer.IloIpAddress,
		testServer.IloIpGateway,
		testServer.IloIpNetmask,
		testServer.IloPassword,
		testServer.IloUsername,
		testServer.InterfaceMtu,
		testServer.InterfaceName,
		testServer.Ip6Address,
		testServer.Ip6Gateway,
		testServer.IpAddress,
		testServer.IpGateway,
		testServer.IpNetmask,
		testServer.LastUpdated,
		testServer.MgmtIpAddress,
		testServer.MgmtIpGateway,
		testServer.MgmtIpNetmask,
		testServer.OfflineReason,
		testServer.PhysLocation,
		testServer.PhysLocationId,
		testServer.Profile,
		testServer.ProfileDesc,
		testServer.ProfileId,
		testServer.Rack,
		testServer.RouterHostName,
		testServer.RouterPortName,
		testServer.Status,
		testServer.StatusId,
		testServer.TcpPort,
		testServer.ServerType,
		testServer.ServerTypeId,
		testServer.UpdPending,
		testServer.XmppId,
		testServer.XmppPasswd,
	)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := url.Values{}
	v.Set("type", "EDGE")
	servers, err := getServers(v, db, PrivLevelAdmin)
	if err != nil {
		t.Errorf("getServers expected: nil error, actual: %v", err)
	}

	if len(servers) != 1 {
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
