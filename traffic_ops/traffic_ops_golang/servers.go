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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
	"github.com/jmoiron/sqlx"
)

const ServersPrivLevel = 10

func serversHandler(db *sqlx.DB) AuthRegexHandlerFunc {
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

func ServersQuery() string {
	query := `SELECT
cg.name as cachegroup,
s.cachegroup as cachegroupId,
s.cdn_id as cdnId,
cdn.name as cdnName,
s.domain_name as domainName,
COALESCE(s.guid, '') as guid,
s.host_name as hostName,
COALESCE(s.https_port, 0) as httpsPort,
s.id as id,
COALESCE(s.ilo_ip_address, '') as iloIpAddress,
COALESCE(s.ilo_ip_gateway, '') as iloIpGateway,
COALESCE(s.ilo_ip_netmask, '') as iloIpNetmask,
COALESCE(s.ilo_password, '') as iloPassword,
COALESCE(s.ilo_username, '') as iloUsername,
s.interface_mtu as interfaceMtu,
s.interface_name as interfaceName,
s.ip6_address as ip6Address,
s.ip6_gateway as ip6Gateway,
s.ip_address as ipAddress,
s.ip_gateway as ipGateway,
s.ip_netmask as ipNetmask,
s.last_updated as lastUpdated,
s.mgmt_ip_address as mgmtIpAddress,
s.mgmt_ip_gateway as mgmtIpGateway,
s.mgmt_ip_netmask as mgmtIpNetmask,
s.offline_reason as offlineReason,
pl.name as physLocation,
s.phys_location as physLocationId,
p.name as profile,
p.description as profileDesc,
s.profile as profileId,
s.rack as rack,
s.router_host_name as routerHostName,
s.router_port_name as routerPortName,
st.name as status,
s.status as statusId,
s.tcp_port as tcpPort,
t.name as serverType,
s.type as serverTypeId,
s.upd_pending as updPending,
s.xmpp_id as xmppId,
s.xmpp_passwd as xmppPasswd
FROM server s
JOIN cachegroup cg ON s.cachegroup = cg.id
JOIN cdn cdn ON s.cdn_id = cdn.id
JOIN phys_location pl ON s.phys_location = pl.id
JOIN profile p ON s.profile = p.id
JOIN status st ON s.status = st.id
JOIN type t ON s.type = t.id`

	return query
}

func getServers(q url.Values, db *sqlx.DB, privLevel int) ([]Server, error) {

	rows, err := db.Queryx(ServersQuery())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	servers := []Server{}

	const HiddenField = "********"
	for rows.Next() {
		var s Server
		err = rows.StructScan(&s)
		if privLevel < PrivLevelAdmin {
			s.IloPassword = HiddenField
			s.XmppPasswd = HiddenField
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func getServersResponse(q url.Values, db *sqlx.DB, privLevel int) (*ServersResponse, error) {
	servers, err := getServers(q, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("error getting servers: %v", err)
	}

	resp := ServersResponse{
		Response: servers,
	}
	return &resp, nil
}
