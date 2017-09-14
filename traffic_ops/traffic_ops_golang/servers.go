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
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tcstructs"
)

const ServersPrivLevel = 10

func serversHandler(db *sql.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		q := r.URL.Query()
		resp, err := getServersResponse(q, db, privLevel)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getServersResponse(q url.Values, db *sql.DB, privLevel int) (*tcstructs.ServersResponse, error) {
	servers, err := getServers(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting servers response: %v", err)
	}

	resp := tcstructs.ServersResponse{
		Response: servers,
	}
	return &resp, nil
}

func getServers(q url.Values, db *sql.DB, privLevel int) ([]tcstructs.Server, error) {

	rows, err := db.Query(selectServersQuery())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	servers := []tcstructs.Server{}

	const HiddenField = "********"

	for rows.Next() {
		var cachegroup sql.NullString
		var cachegroupId sql.NullInt64
		var cdnId sql.NullInt64
		var cdnName sql.NullString
		var domainName sql.NullString
		var guid sql.NullString
		var hostName sql.NullString
		var httpsPort sql.NullInt64
		var id sql.NullInt64
		var iloIpAddress sql.NullString
		var iloIpGateway sql.NullString
		var iloIpNetmask sql.NullString
		var iloPassword sql.NullString
		var iloUsername sql.NullString
		var interfaceMtu sql.NullInt64
		var interfaceName sql.NullString
		var ip6Address sql.NullString
		var ip6Gateway sql.NullString
		var ipAddress sql.NullString
		var ipGateway sql.NullString
		var ipNetmask sql.NullString
		var lastUpdated sql.NullString
		var mgmtIpAddress sql.NullString
		var mgmtIpGateway sql.NullString
		var mgmtIpNetmask sql.NullString
		var offlineReason sql.NullString
		var physLocation sql.NullString
		var physLocationId sql.NullInt64
		var profile sql.NullString
		var profileDesc sql.NullString
		var profileId sql.NullInt64
		var rack sql.NullString
		var revalPending sql.NullBool
		var routerHostName sql.NullString
		var routerPortName sql.NullString
		var status sql.NullString
		var statusId sql.NullInt64
		var tcpPort sql.NullInt64
		var serverType sql.NullString
		var serverTypeId sql.NullInt64
		var updPending sql.NullBool
		var xmppId sql.NullString
		var xmppPasswd sql.NullString
		if err := rows.Scan(&cachegroup, &cachegroupId, &cdnId, &cdnName, &domainName, &guid, &hostName, &httpsPort, &id, &iloIpAddress, &iloIpGateway, &iloIpNetmask, &iloPassword, &iloUsername, &interfaceMtu, &interfaceName, &ip6Address, &ip6Gateway, &ipAddress, &ipGateway, &ipNetmask, &lastUpdated, &mgmtIpAddress, &mgmtIpGateway, &mgmtIpNetmask, &offlineReason, &physLocation, &physLocationId, &profile, &profileDesc, &profileId, &rack, &revalPending, &routerHostName, &routerPortName, &status, &statusId, &tcpPort, &serverType, &serverTypeId, &updPending, &xmppId, &xmppPasswd); err != nil {
			return nil, err
		}

		var iloPasswordField string
		var xmppPasswdField string
		if privLevel < PrivLevelAdmin {
			iloPasswordField = HiddenField
			xmppPasswdField = HiddenField
		} else {
			iloPasswordField = iloPassword.String
			xmppPasswdField = iloPassword.String
		}

		servers = append(servers, tcstructs.Server{
			Cachegroup:     cachegroup.String,
			CachegroupId:   int(cachegroupId.Int64),
			CdnId:          int(cdnId.Int64),
			CdnName:        cdnName.String,
			DomainName:     domainName.String,
			Guid:           guid.String,
			HostName:       hostName.String,
			HttpsPort:      int(httpsPort.Int64),
			Id:             int(id.Int64),
			IloIpAddress:   iloIpAddress.String,
			IloIpGateway:   iloIpGateway.String,
			IloIpNetmask:   iloIpNetmask.String,
			IloPassword:    iloPasswordField,
			IloUsername:    iloUsername.String,
			InterfaceMtu:   int(interfaceMtu.Int64),
			InterfaceName:  interfaceName.String,
			Ip6Address:     ip6Address.String,
			Ip6Gateway:     ip6Gateway.String,
			IpAddress:      ipAddress.String,
			IpGateway:      ipGateway.String,
			IpNetmask:      ipNetmask.String,
			LastUpdated:    lastUpdated.String,
			MgmtIpAddress:  mgmtIpAddress.String,
			MgmtIpGateway:  mgmtIpGateway.String,
			MgmtIpNetmask:  mgmtIpNetmask.String,
			OfflineReason:  offlineReason.String,
			PhysLocation:   physLocation.String,
			PhysLocationId: int(physLocationId.Int64),
			Profile:        profile.String,
			ProfileDesc:    profileDesc.String,
			ProfileId:      int(profileId.Int64),
			Rack:           rack.String,
			RevalPending:   revalPending.Bool,
			RouterHostName: routerHostName.String,
			RouterPortName: routerPortName.String,
			Status:         status.String,
			StatusId:       int(statusId.Int64),
			TcpPort:        int(tcpPort.Int64),
			ServerType:     serverType.String,
			ServerTypeId:   int(serverTypeId.Int64),
			UpdPending:     bool(updPending.Bool),
			XmppId:         xmppId.String,
			XmppPasswd:     xmppPasswdField,
		})
	}
	return servers, nil
}

func selectServersQuery() string {

	const JumboFrameBPS = 9000
	//COALESCE is needed to default values that are nil in the database
	// because Go does not allow that to marshal into the struct
	query := `SELECT
cg.name as cachegroup,
s.cachegroup as cachegroup_id,
s.cdn_id,
cdn.name as cdn_name,
s.domain_name,
COALESCE(s.guid, '') as guid,
s.host_name,
COALESCE(s.https_port, 0) as https_port,
s.id,
COALESCE(s.ilo_ip_address, '') as ilo_ip_address,
COALESCE(s.ilo_ip_gateway, '') as ilo_ip_gateway,
COALESCE(s.ilo_ip_netmask, '') as ilo_ip_netmask,
COALESCE(s.ilo_password, '') as ilo_password,
COALESCE(s.ilo_username, '') as ilo_username,
COALESCE(s.interface_mtu, ` + strconv.Itoa(JumboFrameBPS) + `) as interface_mtu,
COALESCE(s.interface_name, '') as interface_name,
COALESCE(s.ip6_address, '') as ip6_address,
COALESCE(s.ip6_gateway, '') as ip6_gateway,
s.ip_address,
s.ip_gateway,
s.ip_netmask,
s.last_updated,
COALESCE(s.mgmt_ip_address, '') as mgmt_ip_address,
COALESCE(s.mgmt_ip_gateway, '') as mgmt_ip_gateway,
COALESCE(s.mgmt_ip_netmask, '') as mgmt_ip_netmask,
COALESCE(s.offline_reason, '') as offline_reason,
pl.name as phys_location,
s.phys_location as phys_location_id,
p.name as profile,
p.description as profile_desc,
s.profile as profile_id,
COALESCE(s.rack, '') as rack,
s.reval_pending,
COALESCE(s.router_host_name, '') as router_host_name,
COALESCE(s.router_port_name, '') as router_port_name,
st.name as status,
s.status as status_id,
COALESCE(s.tcp_port, 0) as tcp_port,
t.name as server_type,
s.type as server_type_id,
s.upd_pending as upd_pending,
COALESCE(s.xmpp_id, '') as xmpp_id,
COALESCE(s.xmpp_passwd, '') as xmpp_passwd

FROM server s

JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id`
	return query
}
