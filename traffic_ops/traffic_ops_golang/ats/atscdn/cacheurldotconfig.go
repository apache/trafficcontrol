package atscdn

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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

func GetCacheURLDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName, cdnExists, err := GetCDNNameFromNameOrID(inf.Tx.Tx, inf.Params["cdn-name-or-id"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cdn name from id: "+err.Error()))
	} else if !cdnExists {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("cdn not found."), nil)
		return
	}

	toToolName, toURL, err := ats.GetToolNameAndURL(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting tool name and url: "+err.Error()))
		return
	}

	fileName := inf.Params["filename"] // note this is the cacheurl{name}.config, not the full filename
	fullFileName := "cacheurl" + fileName + ".config"

	dses, err := GetCacheURLDSes(inf.Tx.Tx, cdnName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cdn dses: "+err.Error()))
		return
	}

	txt := atscfg.MakeCacheURLDotConfig(cdnName, toToolName, toURL, fullFileName, dses)
	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(txt))
}

// TODO test for nil origin, nil qstring ignore
// TODO test performance - we could break up the cacheurl configs, as this is only needed

func GetCacheURLDSes(tx *sql.Tx, cdn tc.CDNName) (map[tc.DeliveryServiceName]atscfg.CacheURLDS, error) {
	dses := map[tc.DeliveryServiceName]atscfg.CacheURLDS{}
	qry := `
SELECT
  ds.xml_id,
  COALESCE(ds.qstring_ignore, 0),
  COALESCE((SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
    FROM origin o
    WHERE o.deliveryservice = ds.id
    AND o.is_primary), '') as org_server_fqdn,
  COALESCE(ds.cacheurl, '')
FROM
  deliveryservice ds
  JOIN deliveryservice_server dss on ds.id = dss.deliveryservice
WHERE
  ds.cdn_id = (select id from cdn where name = $1)
  AND ds.active = true
`
	// note the dss inner join is intentional, to remove dses with no servers
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		dsName := tc.DeliveryServiceName("")
		ds := atscfg.CacheURLDS{}
		if err := rows.Scan(&dsName, &ds.QStringIgnore, &ds.OrgServerFQDN, &ds.CacheURL); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		dses[dsName] = ds
	}
	return dses, nil
}
