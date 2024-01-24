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
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestDeliveryServiceServer() TODeliveryServiceServer {
	return TODeliveryServiceServer{
		APIInfoImpl: api.APIInfoImpl{},
		DeliveryServiceServer: tc.DeliveryServiceServer{
			Server:          nil,
			DeliveryService: nil,
			LastUpdated:     nil,
		},
		TenantIDs:          nil,
		DeliveryServiceIDs: nil,
		ServerIDs:          nil,
		CDN:                "",
	}
}

func TestValidateDSSAssignments(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	mock.ExpectBegin()

	serverInfo := make([]tc.ServerInfo, 0)
	s := tc.ServerInfo{
		Cachegroup:   "cg1",
		CachegroupID: 20,
		CDNID:        10,
		DomainName:   "test",
		HostName:     "blah",
		ID:           100,
		Status:       "ONLINE",
		Type:         "EDGE",
	}
	serverInfo = append(serverInfo, s)
	s2 := s
	s2.ID = 200
	s2.HostName = "blah2"
	serverInfo = append(serverInfo, s2)

	dsInfo := DSInfo{Active: true,
		ID:                   1,
		Name:                 "ds1",
		Type:                 tc.DSTypeDNS,
		EdgeHeaderRewrite:    nil,
		MidHeaderRewrite:     nil,
		RegexRemap:           nil,
		SigningAlgorithm:     nil,
		CacheURL:             nil,
		MaxOriginConnections: nil,
		Topology:             util.Ptr("topology1"),
		CDNID:                util.Ptr(10),
		UseMultiSiteOrigin:   false}

	// Try to assign non-ORG servers to a topology based DS (with required capabilities)
	userErr, sysErr, sc := validateDSSAssignments(db.MustBegin().Tx, dsInfo, serverInfo, false)

	if sysErr != nil {
		t.Errorf("expected no system error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Errorf("expected error while trying to assign EDGE server to a topology, but got nothing")
	}
	if sc != http.StatusBadRequest {
		t.Errorf("expected status code to be 400, but got %d instead", sc)
	}

	// Try to assign ORG servers without required capabilities to a topology based DS (with required capabilities)
	for i, _ := range serverInfo {
		serverInfo[i].Type = "ORG"
	}
	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"array_agg", "array_agg"})
	rows.AddRow([]byte("{20,21}"), []byte("{cg1,cg2}"))
	mock.ExpectQuery("SELECT ARRAY(c.id)*").WithArgs("topology1").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"required_capabilities"})
	rows.AddRow("{reqCap1}")
	mock.ExpectQuery("SELECT required_capabilities*").WithArgs(1).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"host_name", "capabilities"})
	rows.AddRow("blah", "{reqCap2, reqCap3}")
	mock.ExpectQuery("SELECT s.host_name*").WithArgs(pq.StringArray{}).WillReturnRows(rows)
	userErr, sysErr, sc = validateDSSAssignments(db.MustBegin().Tx, dsInfo, serverInfo, false)

	if sysErr != nil {
		t.Errorf("expected no system error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Errorf("expected error while trying to assign server without a required capability, but got nothing")
	}
	if sc != http.StatusBadRequest {
		t.Errorf("expected status code to be 400, but got %d instead", sc)
	}

	// Try to assign ORG servers with required capabilities to a topology based DS (with required capabilities)
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"array_agg", "array_agg"})
	rows.AddRow([]byte("{20,21}"), []byte("{cg1,cg2}"))
	mock.ExpectQuery("SELECT ARRAY(c.id)*").WithArgs("topology1").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"required_capabilities"})
	rows.AddRow("{reqCap1}")
	mock.ExpectQuery("SELECT required_capabilities*").WithArgs(1).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"host_name", "capabilities"})
	rows.AddRow("blah", "{reqCap1, reqCap2, reqCap3}")
	mock.ExpectQuery("SELECT s.host_name*").WithArgs(pq.StringArray{}).WillReturnRows(rows)
	userErr, sysErr, sc = validateDSSAssignments(db.MustBegin().Tx, dsInfo, serverInfo, false)

	if userErr != nil || sysErr != nil {
		t.Errorf("expected no errors, but got userErr: %v, sysErr: %v", userErr, sysErr)
	}
	if sc != http.StatusOK {
		t.Errorf("expected status code to be 200, but got %d instead", sc)
	}

	// Try to assign EDGE servers without required capabilities to a DS (with required capabilities)
	dsInfo.Topology = nil
	for i, _ := range serverInfo {
		serverInfo[i].Type = "EDGE"
	}

	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"required_capabilities"})
	rows.AddRow("{reqCap1}")
	mock.ExpectQuery("SELECT required_capabilities*").WithArgs(1).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"host_name", "capabilities"})
	rows.AddRow("blah", "{reqCap2, reqCap3}")
	mock.ExpectQuery("SELECT s.host_name*").WithArgs(pq.StringArray{"blah", "blah2"}).WillReturnRows(rows)
	userErr, sysErr, sc = validateDSSAssignments(db.MustBegin().Tx, dsInfo, serverInfo, false)

	if sysErr != nil {
		t.Errorf("expected no system error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Errorf("expected error while trying to assign server without a required capability, but got nothing")
	}
	if sc != http.StatusBadRequest {
		t.Errorf("expected status code to be 400, but got %d instead", sc)
	}

	// Try to assign EDGE servers with required capabilities to a DS (with required capabilities)
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"required_capabilities"})
	rows.AddRow("{reqCap1}")
	mock.ExpectQuery("SELECT required_capabilities*").WithArgs(1).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"host_name", "capabilities"})
	rows.AddRow("blah", "{reqCap1, reqCap2, reqCap3}")
	mock.ExpectQuery("SELECT s.host_name*").WithArgs(pq.StringArray{"blah", "blah2"}).WillReturnRows(rows)
	userErr, sysErr, sc = validateDSSAssignments(db.MustBegin().Tx, dsInfo, serverInfo, false)

	if userErr != nil || sysErr != nil {
		t.Errorf("expected no errors, but got userErr: %v, sysErr: %v", userErr, sysErr)
	}
	if sc != http.StatusOK {
		t.Errorf("expected status code to be 200, but got %d instead", sc)
	}
}

func TestHasAvailableEdgesCurrentlyAssigned(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()
	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("edge1")

	mock.ExpectQuery("SELECT t.name AS name*").WithArgs(1).WillReturnRows(rows)
	assigned, err := hasAvailableEdgesCurrentlyAssigned(db.MustBegin().Tx, 1)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
	if !assigned {
		t.Errorf("expected 'hasAvailableEdgesCurrentlyAssigned' to return true, but got false")
	}
}

func TestReadDSS(t *testing.T) {
	//func (dss *TODeliveryServiceServer) readDSS(h http.Header, tx *sqlx.Tx, user *auth.CurrentUser, params map[string]string, intParams map[string]int, dsIDs []int64, serverIDs []int64, useIMS bool) (*tc.DeliveryServiceServerResponse, error, *time.Time)
	dss := getTestDeliveryServiceServer()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"})
	rows.AddRow(10)
	rows.AddRow(20)

	mock.ExpectBegin()
	mock.ExpectQuery("WITH RECURSIVE*").WithArgs(10).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"server", "deliveryservice", "last_updated"})
	rows.AddRow(1, 2, time.Now())

	mock.ExpectQuery("SELECT*").WithArgs(pq.Int64Array{10, 20}).WillReturnRows(rows)
	response, err, _ := dss.readDSS(nil, db.MustBegin(), &auth.CurrentUser{PrivLevel: 30, TenantID: 10}, nil, nil, nil, nil, false)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
	if response == nil {
		t.Fatalf("expected a valid response, but got nothing")
	}
	if len(response.Response) != 1 {
		t.Fatalf("expected response to have 1 deliveryserviceServer, but got %d", len(response.Response))
	}
	if response.Response[0].Server == nil || response.Response[0].DeliveryService == nil {
		t.Fatalf("expected valid values for server and deliveryservice, but got nil instead. server: %v, deliveryservice: %v", response.Response[0].Server, response.Response[0].DeliveryService)
	}
	if *response.Response[0].Server != 1 || *response.Response[0].DeliveryService != 2 {
		t.Errorf("expected server to be 1 and deliveryservice to be 2, but got server: %d, deliveryservice: %d instead", *response.Response[0].Server, *response.Response[0].DeliveryService)
	}
}

func TestValidate(t *testing.T) {
	dss := getTestDeliveryServiceServer()
	err := dss.Validate(nil)
	if err == nil {
		t.Errorf("expected error about deliveryservice and server not being present, but got nothing")
	}
	dss.Server = util.Ptr(1)
	dss.DeliveryService = util.Ptr(2)
	err = dss.Validate(nil)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestSetKeys(t *testing.T) {
	dss := getTestDeliveryServiceServer()
	keys := make(map[string]interface{})
	keys["server"] = 1
	keys["deliveryservice"] = 2
	dss.SetKeys(keys)
	if dss.DeliveryService == nil || dss.Server == nil {
		t.Fatalf("expected both server and deliveryservice to be not nil")
	}
	if *dss.DeliveryService != 2 {
		t.Errorf("expected deliveryservice key to be 2, but got %d", *dss.DeliveryService)
	}
	if *dss.Server != 1 {
		t.Errorf("expected server key to be 1, but got %d", *dss.Server)
	}
}

func TestGetAuditName(t *testing.T) {
	dss := getTestDeliveryServiceServer()

	auditName := dss.GetAuditName()
	if auditName != "unknown" {
		t.Errorf("expected audit name to be 'unknown', but got %s", auditName)
	}

	dss.DeliveryServiceServer.Server = util.Ptr(1)
	dss.DeliveryServiceServer.DeliveryService = util.Ptr(2)
	auditName = dss.GetAuditName()
	if auditName != "2-1" {
		t.Errorf("expected audit name to be '2-1', but got %s", auditName)
	}
}

func TestGetKeys(t *testing.T) {
	dss := getTestDeliveryServiceServer()
	dss.Server = util.Ptr(1)
	dss.DeliveryService = util.Ptr(2)
	keys, exists := dss.GetKeys()
	if keys == nil {
		t.Fatalf("expected function to return a valid map of keys, but got nothing")
	}
	if !exists {
		t.Fatalf("expected function to return a true boolean for exists, got false")
	}
	if serverID, ok := keys["server"]; !ok {
		t.Fatalf("expected returned keys to have 'server' key, but key wasn't present")
	} else if serverID.(int) != 1 {
		t.Errorf("expected serverID to be 1, but got %d", serverID.(int))
	}

	if dsID, ok := keys["deliveryservice"]; !ok {
		t.Fatalf("expected returned keys to have 'deliveryservice' key, but key wasn't present")
	} else if dsID.(int) != 2 {
		t.Errorf("expected dsID to be 2, but got %d", dsID.(int))
	}

	// check with nil values for server and deliveryservice
	dss.DeliveryServiceServer.Server = nil
	dss.DeliveryServiceServer.DeliveryService = nil
	keys, exists = dss.GetKeys()
	if keys == nil {
		t.Fatalf("expected function to return a valid map of keys, but got nothing")
	}
	if exists {
		t.Fatalf("expected function to return a false boolean for exists, got true")
	}
	if dsID, ok := keys["deliveryservice"]; !ok {
		t.Fatalf("expected returned keys to have 'deliveryservice' key, but key wasn't present")
	} else if dsID.(int) != 0 {
		t.Errorf("expected dsID to be 0, but got %d", dsID.(int))
	}
	if _, ok := keys["server"]; ok {
		t.Errorf("'server' key was not expected to be present")
	}
	dss.DeliveryServiceServer.DeliveryService = util.Ptr(2)
	keys, exists = dss.GetKeys()
	if keys == nil {
		t.Fatalf("expected function to return a valid map of keys, but got nothing")
	}
	if exists {
		t.Fatalf("expected function to return a false boolean for exists, got true")
	}
	if serverID, ok := keys["server"]; !ok {
		t.Fatalf("expected returned keys to have 'server' key, but key wasn't present")
	} else if serverID.(int) != 0 {
		t.Errorf("expected serverID to be 0, but got %d", serverID.(int))
	}
}

func TestValidateDSS(t *testing.T) {
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
		"upd_pending"}

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
		)
	}

	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()

	params := make(map[string]int, 0)
	params["id"] = dsID
	inf := api.Info{
		Version: &api.Version{
			Major: 5,
			Minor: 0,
		},
		Tx:        db.MustBegin(),
		IntParams: params,
		User:      &auth.CurrentUser{PrivLevel: 30},
		Config:    &config.Config{RoleBasedPermissions: true},
	}
	actualSrvs, err := read(&inf)
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
