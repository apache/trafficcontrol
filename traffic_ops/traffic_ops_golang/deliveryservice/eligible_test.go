package deliveryservice

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
	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetEligibleServers(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testServers := getMockDSServers()

	mock.ExpectBegin()

	idRows := sqlmock.NewRows([]string{"id"})
	for _, s := range testServers {
		idRows = idRows.AddRow(*s.ID)
	}
	mock.ExpectQuery("SELECT s.id FROM server s (.+)").WithArgs(1).WillReturnRows(idRows)

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
		"server_capabilities",
		"deliveryservice_capabilities"}
	eligbleRows := sqlmock.NewRows(cols)

	for _, s := range testServers {
		eligbleRows = eligbleRows.AddRow(
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
			[]byte(`{""}`),
			[]byte(`{""}`),
		)
	}
	mock.ExpectQuery("SELECT s.id ,").WillReturnRows(eligbleRows)

	mock.ExpectCommit()

	actualSrvs, err := getEligibleServers(db.MustBegin().Tx, 1)
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
		if len(interf.IPAddresses) != 4 {
			t.Fatalf("servers.read expected len(interf.IPAddresses) == 4, actual = %v", len(interf.IPAddresses))
		}
	}

}

func getMockDSServers() []tc.DSServerV4 {
	srv := tc.DSServerV4{}
	srv.ID = util.IntPtr(1)
	srv.Cachegroup = util.StrPtr("cgTest")
	srv.CachegroupID = util.IntPtr(1)
	srv.CDNID = util.IntPtr(1)
	srv.CDNName = util.StrPtr("cdnTest")
	srv.DomainName = util.StrPtr("domain")
	srv.ServerInterfaces = &[]tc.ServerInterfaceInfoV40{}
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
			RouterHostName: "",
			RouterPortName: "",
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
