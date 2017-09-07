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

func getServers(v url.Values, db *sqlx.DB, privLevel int) ([]Server, error) {

	var rows *sqlx.Rows
	var err error

	wc := newWhereClause(v)
	query := SelectStatement{
		Select: selectQuery(),
		Where:  wc,
	}
	if wc.Exists() {
		rows, err = db.Queryx(query.String(), wc.Condition.Value)
	} else {
		rows, err = db.Queryx(query.String())
	}

	if err != nil {
		//TODO: drichardson - send back an alert if the Query Count is larger than 1
		//                    Test for bad Query Parameters
		return nil, err
	}
	servers := []Server{}

	const HiddenField = "********"
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer rows.Close()
	for rows.Next() {
		var s Server
		err = rows.StructScan(&s)
		if err != nil {
			return nil, fmt.Errorf("getting servers: %v", err)
		}
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
		return nil, fmt.Errorf("getting servers response: %v", err)
	}

	resp := ServersResponse{
		Response: servers,
	}
	return &resp, nil
}

func selectQuery() string {

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
COALESCE(s.interface_mtu, 9000) as interface_mtu,
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

const (
	EQUAL     = "="
	NOT_EQUAL = "!="
	OR        = "OR"
)

type Condition struct {
	Key     string
	Operand string
	Value   string
}

type SelectStatement struct {
	Select string
	Where  WhereClause
}

func (q *SelectStatement) String() string {
	if q.Where.Exists() {
		return q.Select + q.Where.String()
	} else {
		return q.Select
	}
}

type WhereClause struct {
	Condition Condition
}

func (w *WhereClause) SetCondition(c Condition) Condition {
	w.Condition = c
	return w.Condition
}

func (w *WhereClause) String() string {
	c := w.Condition
	return "\nWHERE " + c.Key + c.Operand + "$1"
}

func (w *WhereClause) Exists() bool {
	if (Condition{}) != w.Condition {
		return true
	} else {
		return false
	}
}

func newWhereClause(v url.Values) WhereClause {

	whereClause := WhereClause{}

	switch {
	case v.Get("cachegroup") != "":
		whereClause.SetCondition(Condition{"s.cachegroup", EQUAL, v.Get("cachegroup")})

	// Support what should have been the cachegroupId as well
	case v.Get("cachegroupId") != "":
		whereClause.SetCondition(Condition{"s.cachegroup", EQUAL, v.Get("cachegroupId")})

	case v.Get("cdn") != "":
		whereClause.SetCondition(Condition{"s.cdn_id", EQUAL, v.Get("cdn")})

	case v.Get("physLocation") != "":
		whereClause.SetCondition(Condition{"s.phys_location", EQUAL, v.Get("physLocation")})

	case v.Get("physLocationId") != "":
		whereClause.SetCondition(Condition{"s.phys_location", EQUAL, v.Get("physLocationId")})

	case v.Get("profileId") != "":
		whereClause.SetCondition(Condition{"s.profile", EQUAL, v.Get("profileId")})

	case v.Get("type") != "":
		whereClause.SetCondition(Condition{"s.type", EQUAL, v.Get("type")})

	case v.Get("typeId") != "":
		whereClause.SetCondition(Condition{"s.type", EQUAL, v.Get("typeId")})
	}
	return whereClause
}
