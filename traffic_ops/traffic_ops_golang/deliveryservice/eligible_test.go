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
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
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
	cols := []string{"cachegroup",
		"cachegroup_id",
		"cdn_id",
		"cdn_name",
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
		"last_updated",
		"mgmt_ip_address",
		"mgmt_ip_gateway",
		"mgmt_ip_netmask",
		"offline_reason",
		"phys_location",
		"phys_location_id",
		"profile",
		"profile_desc",
		"profile_id",
		"rack",
		"router_host_name",
		"router_port_name",
		"status",
		"status_id",
		"tcp_port",
		"server_type",
		"server_type_id",
		"upd_pending",
		"server_capabilities",
		"deliveryservice_capabilities"}

	rows := sqlmock.NewRows(cols)

	for _, s := range testServers {
		rows = rows.AddRow(
			s.Cachegroup,
			s.CachegroupID,
			s.CDNID,
			s.CDNName,
			s.DomainName,
			s.GUID,
			s.HostName,
			s.HTTPSPort,
			s.ID,
			s.ILOIPAddress,
			s.ILOIPGateway,
			s.ILOIPNetmask,
			s.ILOPassword,
			s.ILOUsername,
			[]byte(`{"{\"ipAddresses\" : [{\"address\" : \"127.0.0.0\", \"gateway\" : null, \"service_address\" : true}], \"max_bandwidth\" : null, \"monitor\" : true, \"mtu\" : 1500, \"name\" : \"eth0\"}"}`),
			s.LastUpdated,
			s.MgmtIPAddress,
			s.MgmtIPGateway,
			s.MgmtIPNetmask,
			s.OfflineReason,
			s.PhysLocation,
			s.PhysLocationID,
			s.Profile,
			s.ProfileDesc,
			s.ProfileID,
			s.Rack,
			s.RouterHostName,
			s.RouterPortName,
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

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	actualSrvs, err := getEligibleServers(db.MustBegin().Tx, 1)
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
}

func getMockDSServers() []tc.DSServer {
	base := tc.DSServerBase{
		Cachegroup:   util.StrPtr("cgTest"),
		CachegroupID: util.IntPtr(1),
		CDNID:        util.IntPtr(1),
		CDNName:      util.StrPtr("cdnTest"),
		DomainName:   util.StrPtr("domain"),
	}
	srv := tc.DSServer{
		base,
		&[]tc.ServerInterfaceInfo{}, // left empty because it must be written as json above since sqlmock does not support nested arrays
	}
	srvsExpected := []tc.DSServer{srv}
	return srvsExpected
}
