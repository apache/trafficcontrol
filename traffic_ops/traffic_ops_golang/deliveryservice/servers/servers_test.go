package servers

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
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestValidateDSSAssignments(t *testing.T) {
	expected := `server and delivery service CDNs do not match`
	cdnID := 1
	ds := DSInfo{
		ID:    0,
		CDNID: &cdnID,
	}
	var servers []tc.ServerInfo
	server := tc.ServerInfo{
		HostName: "serverHost",
		CDNID:    0,
		Type:     "",
	}
	servers = append(servers, server)
	userErr, _, _ := validateDSS(nil, ds, servers)
	if userErr == nil {
		t.Fatalf("Expected user error with mismatching ds and server CDN IDs, got no error instead")
	}
	if userErr.Error() != expected {
		t.Errorf("Expected error details %v, got %v", expected, userErr.Error())
	}
	servers[0].CDNID = 1
	userErr, _, _ = validateDSS(nil, ds, servers)
	if userErr != nil {
		t.Fatalf("Expected no user error, got %v", userErr.Error())
	}
}

func TestReadServers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testServers := getMockDSServers()
	dsID := 1
	mock.ExpectBegin()

	idRows := sqlmock.NewRows([]string{"id"})
	for _, s := range testServers {
		idRows = idRows.AddRow(*s.ID)
	}
	mock.ExpectQuery("SELECT s.id FROM server s (.+)").WithArgs(dsID).WillReturnRows(idRows)

	for _, s := range testServers {
		testInterfaces := createServerIntefaces(*s.ID)
		mockServerInterfaces(mock, *s.ID, testInterfaces)
	}

	cols := []string{
		"id",
		"cachegroup",
		"cachegroup_id",
		"cdn_id",
		"cdn_name",
		"domain_name",
		"guid",
		"host_name",
		"https_port",
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
		"upd_pending",
		"asns",
	}

	rows := sqlmock.NewRows(cols)

	for _, s := range testServers {
		rows = rows.AddRow(
			s.ID,
			s.Cachegroup,
			s.CachegroupID,
			s.CDNID,
			s.CDNName,
			s.DomainName,
			s.GUID,
			s.HostName,
			s.HTTPSPort,
			s.ILOIPAddress,
			s.ILOIPGateway,
			s.ILOIPNetmask,
			s.ILOPassword,
			s.ILOUsername,
			s.LastUpdated,
			s.MgmtIPAddress,
			s.MgmtIPGateway,
			s.MgmtIPNetmask,
			s.OfflineReason,
			s.PhysLocation,
			s.PhysLocationID,
			fmt.Sprintf("{%s}", strings.Join(s.ProfileNames, ",")),
			s.Rack,
			s.Status,
			s.StatusID,
			s.TCPPort,
			s.Type,
			s.TypeID,
			s.UpdPending,
			[]byte(`{1,2}`),
		)
	}

	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	actualSrvs, err := read(db.MustBegin(), dsID, &auth.CurrentUser{PrivLevel: 30})
	if err != nil {
		t.Fatalf("an error '%s' occurred during read", err)
	}

	if len(actualSrvs) != 1 {
		t.Fatalf("servers.read expected len(actualSrvs) == 1, actual = %v", len(actualSrvs))
	}

	srvInts := *(actualSrvs[0]).ServerInterfaces
	if len(srvInts) != 2 {
		t.Fatalf("servers.read expected len(srvInts) == 2, actual = %v", len(srvInts))
	}

	for _, interf := range srvInts {
		if interf.RouterHostName != "router1" && interf.RouterHostName != "router2" {
			t.Errorf("RouterHostName %s does't match router1 or router2", interf.RouterHostName)
		}
		if interf.RouterPortName != "9090" && interf.RouterPortName != "9091" {
			t.Errorf("RouterPortName %s does't match 9090 or 9091", interf.RouterPortName)
		}
		if len(interf.IPAddresses) != 4 {
			t.Fatalf("servers.read expected len(interf.IPAddresses) == 4, actual = %v", len(interf.IPAddresses))
		}
	}
}

func getMockDSServers() []tc.DSServerV4 {
	base := tc.DSServerBaseV4{
		ID:           util.IntPtr(1),
		Cachegroup:   util.StrPtr("cgTest"),
		CachegroupID: util.IntPtr(1),
		CDNID:        util.IntPtr(1),
		CDNName:      util.StrPtr("cdnTest"),
		DomainName:   util.StrPtr("domain"),
	}
	srvV40 := tc.DSServerV40{
		DSServerBaseV4:   base,
		ServerInterfaces: &[]tc.ServerInterfaceInfoV40{}, // left empty because it must be written as json above since sqlmock does not support nested arrays
	}
	srv := tc.DSServerV4{
		DSServerV40: srvV40,
	}
	srvsExpected := []tc.DSServerV4{srv}
	return srvsExpected
}

func createServerIntefaces(cacheID int) []tc.ServerInterfaceInfoV40 {
	return []tc.ServerInterfaceInfoV40{
		{
			ServerInterfaceInfo: tc.ServerInterfaceInfo{
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
				MaxBandwidth: util.Uint64Ptr(2500),
				Monitor:      true,
				MTU:          util.Uint64Ptr(1500),
				Name:         "interfaceName" + strconv.Itoa(cacheID),
			},
			RouterHostName: "router1",
			RouterPortName: "9090",
		},
		{
			ServerInterfaceInfo: tc.ServerInterfaceInfo{
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
				MaxBandwidth: util.Uint64Ptr(1500),
				Monitor:      false,
				MTU:          util.Uint64Ptr(1500),
				Name:         "interfaceName2" + strconv.Itoa(cacheID),
			},
			RouterHostName: "router2",
			RouterPortName: "9091",
		},
	}
}

func mockServerInterfaces(mock sqlmock.Sqlmock, cacheID int, serverInterfaces []tc.ServerInterfaceInfoV40) {
	interfaceRows := sqlmock.NewRows([]string{"max_bandwidth", "monitor", "mtu", "name", "server", "router_host_name", "router_port_name"})
	ipAddressRows := sqlmock.NewRows([]string{"address", "gateway", "service_address", "interface", "server"})
	for _, interf := range serverInterfaces {
		interfaceRows = interfaceRows.AddRow(*interf.MaxBandwidth, interf.Monitor, *interf.MTU, interf.Name, cacheID, interf.RouterHostName, interf.RouterPortName)
		for _, ip := range interf.IPAddresses {
			ipAddressRows = ipAddressRows.AddRow(ip.Address, *ip.Gateway, ip.ServiceAddress, interf.Name, cacheID)
		}
	}

	mock.ExpectQuery("SELECT (.+) FROM interface").WillReturnRows(interfaceRows)
	mock.ExpectQuery("SELECT (.+) FROM ip_address").WillReturnRows(ipAddressRows)
}
