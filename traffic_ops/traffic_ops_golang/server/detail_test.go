package server

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetDetailServers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testServerDetails := getMockServerDetails()
	mock.ExpectBegin()

	idRows := sqlmock.NewRows([]string{"id"})
	for _, s := range testServerDetails {
		idRows = idRows.AddRow(*s.ID)
	}
	mock.ExpectQuery("SELECT server.id").WillReturnRows(idRows)

	for _, s := range testServerDetails {
		testInterfaces := createServerIntefaces(*s.ID)
		mockServerInterfaces(mock, *s.ID, testInterfaces)
	}

	detailCols := []string{
		"id",
		"cachegroup",
		"cdn_name",
		"deliveryservices",
		"domain_name",
		"guid",
		"host_name",
		"https_port",
		"ilo_ip_address",
		"ilo_ip_gateway",
		"ilo_ip_netmask",
		"ilo_password",
		"ilo_username",
		"service_ip",
		"service_ip6",
		"service_gateway",
		"service_gateway6",
		"service_netmask",
		"interface_name",
		"interface_mtu",
		"mgmt_ip_address",
		"mgmt_ip_gateway",
		"mgmt_ip_netmask",
		"offline_reason",
		"phys_location",
		"profile_name",
		"rack",
		"status",
		"tcp_port",
		"server_type",
		"xmpp_id",
		"xmpp_passwd",
	}
	detailRows := sqlmock.NewRows(detailCols)

	serviceAddress := util.StrPtr("")
	service6Address := util.StrPtr("")
	serviceGateway := util.StrPtr("")
	service6Gateway := util.StrPtr("")
	serviceNetmask := util.StrPtr("")
	serviceInterface := util.StrPtr("")
	serviceMtu := util.StrPtr("")

	for _, sd := range testServerDetails {
		detailRows = detailRows.AddRow(
			sd.ID,
			sd.CacheGroup,
			sd.CDNName,
			[]byte(`{1}`),
			sd.DomainName,
			sd.GUID,
			sd.HostName,
			sd.HTTPSPort,
			sd.ILOIPAddress,
			sd.ILOIPGateway,
			sd.ILOIPNetmask,
			sd.ILOPassword,
			sd.ILOUsername,
			serviceAddress,
			service6Address,
			serviceGateway,
			service6Gateway,
			serviceNetmask,
			serviceInterface,
			serviceMtu,
			sd.MgmtIPAddress,
			sd.MgmtIPGateway,
			sd.MgmtIPNetmask,
			sd.OfflineReason,
			sd.PhysLocation,
			fmt.Sprintf("{%s}", strings.Join(sd.ProfileNames, ",")),
			sd.Rack,
			sd.Status,
			sd.TCPPort,
			sd.Type,
			sd.XMPPID,
			sd.XMPPPasswd,
		)
	}
	mock.ExpectQuery("SELECT server.id ,").WillReturnRows(detailRows)

	hwInfoRows := sqlmock.NewRows([]string{"serverid", "description", "val"})
	hwInfoRows = hwInfoRows.AddRow(1, "desc1", "val1")
	hwInfoRows = hwInfoRows.AddRow(1, "desc2", "val2")
	hwInfoRows = hwInfoRows.AddRow(1, "desc3", "val3")

	mock.ExpectQuery("SELECT serverid").WillReturnRows(hwInfoRows)
	mock.ExpectCommit()

	actualSrvs, err := getDetailServers(db.MustBegin().Tx, &auth.CurrentUser{PrivLevel: 30}, "test", 1, "id", 10, api.Version{Major: 4})
	if err != nil {
		t.Fatalf("an error '%s' occurred during read", err)
	}

	if len(actualSrvs) != 1 {
		t.Fatalf("servers.read expected len(actualSrvs) == 1, actual = %v", len(actualSrvs))
	}

	if len(actualSrvs[0].HardwareInfo) != 3 {
		t.Fatalf("servers.read expected len(actualSrvs[0].HardwareInfo) == 3, actual = %v", len(actualSrvs[0].HardwareInfo))
	}

	srvInts := actualSrvs[0].ServerInterfaces
	if len(srvInts) != 2 {
		t.Fatalf("servers.read expected len(srvInts) == 2, actual = %v", len(srvInts))
	}

	for _, interf := range srvInts {
		if len(interf.IPAddresses) != 4 {
			t.Fatalf("servers.read expected len(interf.IPAddresses) == 4, actual = %v", len(interf.IPAddresses))
		}
	}
}

func getMockServerDetails() []tc.ServerDetailV40 {
	srvData := tc.ServerDetailV40{
		ID:               util.IntPtr(1),
		ServerInterfaces: []tc.ServerInterfaceInfoV40{}, // left empty because it must be written as json above since sqlmock does not support nested arrays
	}
	return []tc.ServerDetailV40{srvData}
}

func createServerIntefaces(cacheID int) []tc.ServerInterfaceInfoV40 {
	return []tc.ServerInterfaceInfoV40{
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
			RouterHostName: "",
			RouterPortName: "",
		},
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
			RouterHostName: "",
			RouterPortName: "",
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
