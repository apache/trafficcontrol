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
	"reflect"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestServers() []tc.Server {
	servers := []tc.Server{}
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
		Type:           "EDGE",
		TypeID:         1,
		UpdPending:     true,
		XMPPID:         "xmppId",
		XMPPPasswd:     "xmppPasswd",
	}
	servers = append(servers, testServer)

	testServer2 := testServer
	testServer2.Cachegroup = "cachegroup2"
	testServer2.HostName = "server2"
	servers = append(servers, testServer2)

	testServer3 := testServer
	testServer3.Cachegroup = "cachegroup3"
	testServer3.HostName = "server3"
	servers = append(servers, testServer2)

	return servers
}

func TestGetServersByCachegroup(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testServers := getTestServers()
	cols := test.ColsFromStructByTag("db", tc.Server{})
	rows := sqlmock.NewRows(cols)

	//TODO: drichardson - build helper to add these Rows from the struct values
	//                    or by CSV if types get in the way
	for _, ts := range testServers {
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
			ts.Type,
			ts.TypeID,
			ts.UpdPending,
			ts.XMPPID,
			ts.XMPPPasswd,
		)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	v := url.Values{}
	v.Set("cachegroup", "cachegroup2")

	servers, err := getServers(v, db, PrivLevelAdmin)
	log.Debugln("%v-->", servers)
	if err != nil {
		t.Errorf("getServers expected: nil error, actual: %v", err)
	}

	if len(servers) != 3 {
		t.Errorf("getServers expected: len(servers) == 3, actual: %v", len(servers))
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

func TestAssignDsesToServer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	newDses := []int{4, 5, 6}
	pqNewDses := pq.Array(newDses)

	mock.ExpectBegin()
	mock.ExpectPrepare("DELETE").ExpectExec().WithArgs(100).WillReturnResult(sqlmock.NewResult(1, 3))
	mock.ExpectPrepare("INSERT").ExpectExec().WithArgs(pqNewDses, 100).WillReturnResult(sqlmock.NewResult(1, 3))

	//dses query:
	dsesRows := sqlmock.NewRows([]string{"deliveryservice"})
	for _, i := range newDses {
		dsesRows.AddRow(i)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(dsesRows)

	//fetch remap config location
	remapConfigLocation := "a/path/to/a/remap.config"
	remapConfigRow := sqlmock.NewRows([]string{"value"})
	remapConfigRow.AddRow(remapConfigLocation + "/")          // verifies we strip off the trailing slash
	mock.ExpectQuery("SELECT").WillReturnRows(remapConfigRow) //remap.config

	//select xmlids and edge_header_rewrite, regex_remap, and cache_url  for each ds
	dsFieldRows := sqlmock.NewRows([]string{"xml_id", "edge_header_rewrite", "regex_remap", "cacheurl"})
	dsFieldRows.AddRow("ds1", nil, "regexRemapPlaceholder", "cacheurlPlaceholder")
	dsFieldRows.AddRow("ds2", "edgeHeaderRewritePlaceholder2", "regexRemapPlaceholder", "cacheurlPlaceholder")
	dsFieldRows.AddRow("ds3", "", nil, "cacheurlPlaceholder")
	mock.ExpectPrepare("SELECT").ExpectQuery().WithArgs(pqNewDses).WillReturnRows(dsFieldRows)

	//prepare the insert and delete parameter slices as they should be constructed in the function
	headerRewritePrefix := "hdr_rw_"
	regexRemapPrefix := "regex_remap_"
	cacheurlPrefix := "cacheurl_"
	configPostfix := ".config"
	insert := []string{regexRemapPrefix + "ds1" + configPostfix, cacheurlPrefix + "ds1" + configPostfix, headerRewritePrefix + "ds2" + configPostfix, regexRemapPrefix + "ds2" + configPostfix, cacheurlPrefix + "ds2" + configPostfix, cacheurlPrefix + "ds3" + configPostfix}
	delete := []string{headerRewritePrefix + "ds1" + configPostfix, headerRewritePrefix + "ds3" + configPostfix, regexRemapPrefix + "ds3" + configPostfix}
	fileNamesPq := pq.Array(insert)
	//insert the parameters
	mock.ExpectPrepare("INSERT").ExpectExec().WithArgs(fileNamesPq, "location", remapConfigLocation).WillReturnResult(sqlmock.NewResult(1, 6))

	//select out the parameterIds we just inserted
	parameterIdRows := sqlmock.NewRows([]string{"id"})
	parameterIds := []int64{1, 2, 3, 4, 5, 6}
	for _, i := range parameterIds {
		parameterIdRows.AddRow(i)
	}
	mock.ExpectPrepare("SELECT").ExpectQuery().WithArgs(fileNamesPq).WillReturnRows(parameterIdRows)

	//insert those ids as profile_parameters
	mock.ExpectPrepare("INSERT").ExpectExec().WithArgs(pqNewDses, pq.Array(parameterIds)).WillReturnResult(sqlmock.NewResult(6, 6))

	//delete the parameters in the delete list
	mock.ExpectPrepare("DELETE").ExpectExec().WithArgs(pq.Array(delete)).WillReturnResult(sqlmock.NewResult(1, 3))
	mock.ExpectCommit()

	result, err := assignDeliveryServicesToServer(100, newDses, true, db)
	if err != nil {
		t.Errorf("error assigning deliveryservice: %v", err)
	}
	if !reflect.DeepEqual(result, newDses) {
		t.Errorf("delivery services assigned: Expected %v.   Got  %v", newDses, result)
	}
}
