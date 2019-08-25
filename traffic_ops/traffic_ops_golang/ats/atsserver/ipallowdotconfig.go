package atsserver

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
	"errors"
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

func GetIPAllowDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id-or-host"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	serverName, serverType, ok, err := GetServerNameAndTypeFromNameOrID(inf.Tx.Tx, inf.Params["id-or-host"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server name from ID: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	toToolName, toURL, err := ats.GetToolNameAndURL(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting tool name and url: "+err.Error()))
		return
	}

	params, err := GetServerParams(inf.Tx.Tx, serverName, atscfg.IPAllowConfigFileName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server '"+string(serverName)+"' + ip_allow parameters: "+err.Error()))
		return
	}

	childServers := map[tc.CacheName]atscfg.IPAllowServer{}
	if strings.HasPrefix(string(serverType), tc.MidTypePrefix) {
		if childServers, err = GetChildServers(inf.Tx.Tx, serverName); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting child servers from mid '"+string(serverName)+"': "+err.Error()))
			return
		}
	}

	txt := atscfg.MakeIPAllowDotConfig(serverName, serverType, toToolName, toURL, params, childServers)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(txt))
}

// GetChildServers returns the child servers of the given Mid serverName. This should not be called with an Edge server.
func GetChildServers(tx *sql.Tx, serverName tc.CacheName) (map[tc.CacheName]atscfg.IPAllowServer, error) {
	qry := `
SELECT
  s.host_name,
  s.ip_address,
  COALESCE(s.ip6_address, '')
FROM
  server s
  JOIN type tp on tp.id = s.type
  JOIN cachegroup cg on cg.id = s.cachegroup
WHERE
  (tp.name = '` + tc.MonitorTypeName + `' OR tp.name LIKE '` + tc.EdgeTypePrefix + `%')
  AND cg.id IN (
    SELECT
      cg2.id
    FROM
     server s2
     JOIN cachegroup cg2 ON (cg2.parent_cachegroup_id = s2.cachegroup OR cg2.secondary_parent_cachegroup_id = s2.cachegroup)
    WHERE
      s2.host_name = $1
  )
`
	rows, err := tx.Query(qry, serverName)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	servers := map[tc.CacheName]atscfg.IPAllowServer{}
	for rows.Next() {
		svName := tc.CacheName("")
		sv := atscfg.IPAllowServer{}
		if err := rows.Scan(&svName, &sv.IPAddress, &sv.IP6Address); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		servers[svName] = sv
	}
	return servers, nil
}

func GetServerParams(tx *sql.Tx, serverName tc.CacheName, configFile string) (map[string][]string, error) {
	qry := `
SELECT
  pa.name,
  pa.value
FROM
  parameter pa
  JOIN profile_parameter pp ON pp.parameter = pa.id
  JOIN profile pr ON pr.id = pp.profile
  JOIN server s ON s.profile = pr.id
WHERE
  s.host_name = $1
  AND pa.config_file = $2
`
	rows, err := tx.Query(qry, serverName, configFile)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[string][]string{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		params[name] = append(params[name], val)
	}
	return params, nil
}
