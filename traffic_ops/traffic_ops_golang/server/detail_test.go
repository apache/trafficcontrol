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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
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
	cols := []string{"cachegroup",
		"cdn_name",
		"deliveryservices",
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
		"interfaces",
		"mgmt_ip_address",
		"mgmt_ip_gateway",
		"mgmt_ip_netmask",
		"offline_reason",
		"phys_location",
		"profile",
		"profile_desc",
		"rack",
		"router_host_name",
		"router_port_name",
		"status",
		"tcp_port",
		"server_type",
		"xmpp_id",
		"xmpp_passwd",
	}
	rows := sqlmock.NewRows(cols)

	for _, sd := range testServerDetails {
		rows = rows.AddRow(
			sd.CacheGroup,
			sd.CDNName,
			[]byte(`{1}`),
			sd.DomainName,
			sd.GUID,
			sd.HostName,
			sd.HTTPSPort,
			sd.ID,
			sd.ILOIPAddress,
			sd.ILOIPGateway,
			sd.ILOIPNetmask,
			sd.ILOPassword,
			sd.ILOUsername,
			[]byte(`{"{\"ipAddresses\" : [{\"address\" : \"127.0.0.0\", \"gateway\" : null, \"service_address\" : true}], \"max_bandwidth\" : null, \"monitor\" : true, \"mtu\" : 1500, \"name\" : \"eth0\"}"}`),
			sd.MgmtIPAddress,
			sd.MgmtIPGateway,
			sd.MgmtIPNetmask,
			sd.OfflineReason,
			sd.PhysLocation,
			sd.Profile,
			sd.ProfileDesc,
			sd.Rack,
			sd.RouterHostName,
			sd.RouterPortName,
			sd.Status,
			sd.TCPPort,
			sd.Type,
			sd.XMPPID,
			sd.XMPPPasswd,
		)
	}

	hwCols := []string{"serverid", "description", "val"}
	hwRows := sqlmock.NewRows(hwCols)
	hwRows = hwRows.AddRow(1, "desc1", "val1")
	hwRows = hwRows.AddRow(1, "desc2", "val2")
	hwRows = hwRows.AddRow(1, "desc3", "val3")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT cg.name").WillReturnRows(rows)
	mock.ExpectQuery("SELECT serverid").WillReturnRows(hwRows)
	mock.ExpectCommit()

	actualSrvs, err := getDetailServers(db.MustBegin().Tx, &auth.CurrentUser{PrivLevel: 30}, "test", 1, "id", 10, api.Version{3, 0})
	if err != nil {
		t.Fatalf("an error '%s' occurred during read", err)
	}

	if len(actualSrvs) != 1 {
		t.Fatalf("servers.read expected len(actualSrvs) == 1, actual = %v", len(actualSrvs))
	}

	srvInts := *(actualSrvs[0]).ServerInterfaces
	if len(srvInts) != 1 {
		t.Fatalf("servers.read expected len(srvInts) == 1, actual = %v", len(srvInts))
	}

	if len(srvInts[0].IpAddresses) != 1 {
		t.Fatalf("servers.read expected len(srvInts[0].IpAddresses) == 1, actual = %v", len(srvInts[0].IpAddresses))
	}

	if len(actualSrvs[0].HardwareInfo) != 3 {
		t.Fatalf("servers.read expected len(actualSrvs[0].HardwareInfo) == 3, actual = %v", len(actualSrvs[0].HardwareInfo))
	}

	if !srvInts[0].IpAddresses[0].ServiceAddress {
		t.Fatalf("srvInts[0].IpAddresses[0].ServiceAddress expected to be true, actual = %v", srvInts[0].IpAddresses[0].ServiceAddress)
	}
}

func getMockServerDetails() []tc.ServerDetailV30 {
	srvData := tc.ServerDetailV30{
		tc.ServerDetail{
			ID: util.IntPtr(1),
		},
		&[]tc.ServerInterfaceInfo{}, // left empty because it must be written as json above since sqlmock does not support nested arrays
	}
	return []tc.ServerDetailV30{srvData}
}
