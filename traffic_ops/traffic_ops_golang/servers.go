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
	"strconv"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

	"github.com/jmoiron/sqlx"
)

// ServersPrivLevel - privileges for the /servers endpoint
const ServersPrivLevel = 10

func serversHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		// p PathParams, username string, privLevel int
		ctx := r.Context()
		privLevel, err := getPrivLevel(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}
		pathParams, err := getPathParams(ctx)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		q := r.URL.Query()
		for k, v := range pathParams {
			q.Set(k, v)
		}
		resp, err := getServersResponse(q, db, privLevel)
		if err != nil {
			log.Errorln(err)
			handleErr(err, http.StatusInternalServerError)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			log.Errorln("marshaling response %v", err)
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getServersResponse(v url.Values, db *sqlx.DB, privLevel int) (*tc.ServersResponse, error) {
	servers, err := getServers(v, db, privLevel)
	if err != nil {
		return nil, fmt.Errorf("getting servers response: %v", err)
	}

	resp := tc.ServersResponse{
		Response: servers,
	}
	return &resp, nil
}

func getServers(v url.Values, db *sqlx.DB, privLevel int) ([]tc.Server, error) {

	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]string{
		"cachegroup":   "cg.name",
		"cdn":          "s.cdn_id",
		"id":           "s.id",
		"physLocation": "s.phys_location",
		"profileId":    "s.profile",
		"status":       "st.name",
		"type":         "t.name",
	}

	query, queryValues := BuildQuery(v, selectServersQuery(), queryParamsToSQLCols)

	rows, err = db.NamedQuery(query, queryValues)

	if err != nil {
		return nil, fmt.Errorf("querying: %v", err)
	}
	servers := []tc.Server{}

	const HiddenField = "********"

	defer rows.Close()

	for rows.Next() {
		var s tc.Server
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting servers: %v", err)
		}
		if privLevel < PrivLevelAdmin {
			s.ILOPassword = HiddenField
			s.XMPPPasswd = HiddenField
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func selectServersQuery() string {

	const JumboFrameBPS = 9000

	// COALESCE is needed to default values that are nil in the database
	// because Go does not allow that to marshal into the struct
	selectStmt := `SELECT
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

	return selectStmt
}
