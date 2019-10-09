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
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

func GetHostingDotConfig(w http.ResponseWriter, r *http.Request) {
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

	multiParams, err := GetServerParams(inf.Tx.Tx, serverName, atscfg.HostingConfigParamConfigFile)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server '"+string(serverName)+"' + hosting parameters: "+err.Error()))
		return
	}

	params := map[string]string{}
	for name, vals := range multiParams {
		if len(vals) == 0 {
			log.Warnln("hosting config got no parameters for '" + name + "'")
			continue
		}
		if len(vals) > 1 {
			log.Errorln("hosting config parameter name '"+name+"' got multiple values: %+v - using first!", name, vals)
		}
		params[name] = vals[0]
	}

	origins, err := GetServerHostingOrigins(inf.Tx.Tx, serverName, serverType)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting server '"+string(serverName)+"' hosting origins: "+err.Error()))
		return
	}

	log.Errorf("hosting config DEBUG params %+v\n", params)
	log.Errorf("hosting config DEBUG multiParams %+v\n", multiParams)

	txt := atscfg.MakeHostingDotConfig(serverName, toToolName, toURL, params, origins)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(txt))
}

// GetServerHostingOrigins returns the list of origins on delivery services assigned to the given server, to be used in the ATS config file.
// It returns only LIVE_NATNL delivery services, for mids; and only LIVE and LIVE_NATNL services for edges.
func GetServerHostingOrigins(tx *sql.Tx, serverName tc.CacheName, serverType tc.CacheType) ([]string, error) {
	qry := `
SELECT
  DISTINCT(SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')) as org_server_fqdn
FROM
  deliveryservice ds
  JOIN deliveryservice_server dss on dss.deliveryservice = ds.id
  JOIN server s on s.id = dss.server
  LEFT JOIN origin o on (o.deliveryservice = ds.id AND o.is_primary)
`
	if strings.HasPrefix(string(serverType), tc.MidTypePrefix) {
		// Note mids only include active DSes, but edges include inactive DSes as well.
		qry += `
WHERE
  s.cdn_id = (select cdn_id from server where host_name = $1)
  AND ds.type IN (SELECT id FROM type WHERE name like '%` + tc.DSTypeLiveNationalSuffix + `')
  AND ds.active = true
  AND ds.cdn_id = s.cdn_id
`
	} else {
		qry += `
WHERE
  s.host_name = $1
  AND ds.cdn_id = s.cdn_id
  AND ds.type IN (SELECT id FROM type WHERE (name LIKE '%` + tc.DSTypeLiveSuffix + `' OR name LIKE '%` + tc.DSTypeLiveNationalSuffix + `'))
`
	}
	// Note the 'ds.cdn_id = s.cdn_id' in the query shouldn't be necessary, but it is, because there's no DB constraint.

	log.Errorln("hosting config DEBUG qry QQ" + qry + "QQ")

	rows, err := tx.Query(qry, serverName)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	origins := []string{}
	for rows.Next() {
		origin := ""
		if err := rows.Scan(&origin); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		origins = append(origins, origin)
	}
	return origins, nil
}
