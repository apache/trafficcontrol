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

	"github.com/apache/incubator-trafficcontrol.dew/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"

	"github.com/jmoiron/sqlx"

	"database/sql"
	"strings"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/lib/pq"
)

// ServersPrivLevel - privileges for the /servers endpoint
const ServersPrivLevel = 10

func serversHandler(db *sqlx.DB) AuthRegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams, username string, privLevel int) {

		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		q := r.URL.Query()
		for k, v := range p {
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

		w.Header().Set(api.ContentType, api.ApplicationJson)
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
		"profileId":    "s.profileId",
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

func assignDeliveryServicesToServerHandler(db *sqlx.DB) AuthRegexHandlerFunc {

	return func(w http.ResponseWriter, r *http.Request, params PathParams, username string, privLevel int) {
		handleErr := tc.GetHandleErrorFunc(w, r)

		var dsList []int

		err := json.NewDecoder(r.Body).Decode(&dsList)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		q := r.URL.Query()

		replaceQueryParameter := q["replace"][0]
		replace, err := strconv.ParseBool(replaceQueryParameter) //accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. for replace url parameter documentation
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}

		serverPathParameter := params["server"]
		server, err := strconv.Atoi(serverPathParameter)
		if err != nil {
			handleErr(err, http.StatusBadRequest)
			return
		}

		assignedDSes, err := assignDeliveryServicesToServer(server, dsList, replace, db)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		resp := struct {
			Response []int `json:"response"`
			tc.Alerts
		}{assignedDSes, tc.CreateAlerts(tc.SuccessLevel, "successfully assigned dses to server")}
		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		fmt.Fprintf(w, "%s", respBts)
	}
}

type Parameter struct {
	ConfigFile string
	Name       string
	Value      string
}

func assignDeliveryServicesToServer(server int, dses []int, replace bool, db *sqlx.DB) ([]int, error) {
	//transaction rollback in this functions requires err to be set to the proper error or nil before returning
	tx, err := db.Beginx()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if err != nil {
		log.Error.Printf("could not begin transaction: %v\n", err)
		return nil, tc.DBError
	}

	if replace {
		//delete currently assigned dses from server
		deleteCurrent, err := tx.Prepare("DELETE FROM deliveryservice_server WHERE server = $1")
		if err != nil {
			log.Error.Printf("could not prepare deliveryservice_server delete statement: %s\n", err)
			return nil, tc.DBError
		}
		_, err = deleteCurrent.Exec(server)
		if err != nil {
			log.Error.Printf("could not delete old deliveryservice_server associations for server: %s\n", err)
			return nil, tc.DBError
		}
	}

	//assign new dses
	dsPqArray := pq.Array(dses)
	// The common table expressions (CTEs) used below allow for inserting every deliveryService in dses with the server without a loop
	// the result of the select (for server = 100, dses = [1,2,3]) used by the insert essentially looks like:
	//          | server
	//    ---------------
	//      1   |   100
	//      2   |   100
	//      3   |   100
	// UNNEST is used to turn the array of ds values into, essentially, * FROM ( VALUES (1),(2),(3) )
	// this allows for a single insert query instead of a loop over the dses.
	// This pattern is used for the other bulk inserts as well.
	bulkInsert, err := tx.Prepare(`INSERT INTO deliveryservice_server (deliveryservice, server)
	WITH
	q1 AS (SELECT UNNEST($1::bigint[])),
	q2 AS ( SELECT * FROM (VALUES ($2::bigint)) AS server )
	SELECT * FROM q1,q2 ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Error.Printf("could not prepare deliveryservice_server bulk insert: %s\n", err)
		return nil, tc.DBError
	}
	_, err = bulkInsert.Exec(dsPqArray, server)
	if err != nil {
		log.Error.Printf("could not execute deliveryservice_server bulk insert: %s\n", err)
		return nil, tc.DBError
	}
	//select dses assigned now.
	var newDses []int
	tx.Select(&newDses, "SELECT deliveryservice FROM deliveryservice_server where server = $1", server)

	//need remap config location
	row := tx.QueryRow("SELECT value FROM parameter WHERE name = 'location' AND config_file = '" + ats.RemapFile + "'")
	var atsConfigLocation string
	row.Scan(&atsConfigLocation)
	if strings.HasSuffix(atsConfigLocation, "/") {
		atsConfigLocation = atsConfigLocation[:len(atsConfigLocation)-1]
	}

	//we need dses: xmlids and edge_header_rewrite, regex_remap, and cache_url
	selectDsFieldsQuery, err := tx.Prepare("SELECT xml_id, edge_header_rewrite, regex_remap, cacheurl FROM deliveryservice WHERE id = ANY($1::bigint[])")
	if err != nil {
		log.Error.Printf("could not prepare ds fields query: %s\n", err)
		return nil, tc.DBError
	}
	rows, err := selectDsFieldsQuery.Query(dsPqArray)
	if err != nil {
		log.Error.Printf("could not execute ds fields select query: %s\n", err)
		return nil, tc.DBError
	}
	defer rows.Close()

	//create new parameters here as necessary:
	//loop over ds results and build file parameters we need to insert / select
	//for all of: header rewrite, regex_remap, cache_url
	// if ds has it add parameter to insert list
	// other wise add to delete list
	//TODO: DylanVolz this may need to be extended refactored, there are potentially other parameters that need this like urlSigKeys...
	insert := []string{}
	delete := []string{}
	for rows.Next() {
		var XmlId sql.NullString
		var EdgeHeaderRewrite sql.NullString
		var RegexRemap sql.NullString
		var CacheUrl sql.NullString

		if err := rows.Scan(&XmlId, &EdgeHeaderRewrite, &RegexRemap, &CacheUrl); err != nil {
			log.Error.Printf("could not scan ds fields row: %s\n", err)
			return nil, tc.DBError
		}
		if XmlId.Valid && len(XmlId.String) > 0 {
			//param := "hdr_rw_" + XmlId.String + ".config"
			param := ats.GetConfigFile(ats.HeaderRewritePrefix, XmlId.String)
			if EdgeHeaderRewrite.Valid && len(EdgeHeaderRewrite.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
			param = ats.GetConfigFile(ats.RegexRemapPrefix, XmlId.String)
			if RegexRemap.Valid && len(RegexRemap.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
			param = ats.GetConfigFile(ats.CacheUrlPrefix, XmlId.String)
			if CacheUrl.Valid && len(CacheUrl.String) > 0 {
				insert = append(insert, param)
			} else {
				delete = append(delete, param)
			}
		}

	}

	//insert the parameters we selected above:
	insertParams, err := tx.Prepare(`INSERT INTO parameter (config_file, name, value)
	WITH
	q1 AS (SELECT UNNEST($1::text[]) AS config_file),
	q2 AS (SELECT * FROM (VALUES ($2) ) AS name),
	q3 AS (SELECT * FROM (VALUES ($3) ) AS value)
	 SELECT * FROM q1,q2,q3 ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Error.Printf("could not prepare parameter bulk insert query: %s\n", err)
		return nil, tc.DBError
	}
	fileNamePqArray := pq.Array(insert)
	if _, err = insertParams.Exec(fileNamePqArray, "location", atsConfigLocation); err != nil {
		log.Error.Printf("could not execute parameter bulk insert: %s\n", err)
		return nil, tc.DBError
	}

	//select the ids associated with the parameters we created above (may be able to get them from insert above to optimize)
	selectParameterIds, err := tx.Prepare("SELECT id FROM parameter WHERE name = 'location' AND config_file IN ($1)")
	if err != nil {
		log.Error.Printf("could not prepare parameter id select query: %s\n", err)
		return nil, tc.DBError
	}
	rows, err = selectParameterIds.Query(fileNamePqArray)
	if err != nil {
		log.Error.Printf("could not execute parameter id select query: %s\n", err)
		return nil, tc.DBError
	}
	parameterIds := []int64{}
	for rows.Next() {
		var Id int64
		if err := rows.Scan(&Id); err != nil {
			log.Error.Printf("could not scan parameter id: %s\n", err)
			return nil, tc.DBError
		}
		parameterIds = append(parameterIds, Id)
	}

	//associate all parameter ids with the profiles associated with all servers associated with assigned dses.
	insertProfileParams, err := tx.Prepare(`INSERT INTO profile_parameter (profile, parameter)
	WITH
	q1 AS ( SELECT DISTINCT profile FROM server LEFT JOIN deliveryservice_server ON server.id = deliveryservice_server.server WHERE deliveryservice_server.deliveryservice = ANY($1::bigint[]) ),
	q2 AS (SELECT UNNEST($2::bigint[]) AS parameter)
	SELECT * FROM q1,q2
	ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Error.Printf("could not prepare profile_parameter bulk insert: %s\n", err)
		return nil, tc.DBError
	}
	if _, err = insertProfileParams.Exec(dsPqArray, pq.Array(parameterIds)); err != nil {
		log.Error.Printf("could not execute profile_parameter bulk insert: %s\n", err)
		return nil, tc.DBError
	}

	//process delete list
	deleteTx, err := tx.Prepare(`DELETE FROM parameter WHERE name = 'location' AND config_file = ANY($1)`)
	if err != nil {
		log.Error.Printf("could not prepare parameter delete query: %s\n", err)
		return nil, tc.DBError
	}
	if _, err = deleteTx.Exec(pq.Array(delete)); err != nil {
		log.Error.Printf("could not execute parameter delete query: %s\n", err)
		return nil, tc.DBError
	}

	return newDses, nil
}
