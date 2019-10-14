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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

func GetServerNameAndTypeFromNameOrID(tx *sql.Tx, nameOrID string) (tc.CacheName, tc.CacheType, bool, error) {
	if id, err := strconv.Atoi(nameOrID); err == nil {
		return ats.GetServerNameAndTypeFromID(tx, id)
	}
	name := tc.CacheName(nameOrID)
	typ, ok, err := ats.GetServerTypeFromName(tx, name)
	return name, typ, ok, err
}

func GetServerNameAndDomainFromNameOrID(tx *sql.Tx, nameOrID string) (tc.CacheName, string, bool, error) {
	if id, err := strconv.Atoi(nameOrID); err == nil {
		return ats.GetServerNameAndDomainFromID(tx, id)
	}
	name := tc.CacheName(nameOrID)
	typ, ok, err := ats.GetServerDomainFromName(tx, name)
	return name, typ, ok, err
}

func GetServerCacheConfigData(tx *sql.Tx, serverName tc.CacheName, serverType tc.CacheType) (map[tc.DeliveryServiceName]atscfg.ServerCacheConfigDS, error) {
	qry := `
SELECT
  ds.xml_id,
  (o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')) as org_server_fqdn,
  dt.name AS ds_type
FROM
  deliveryservice ds
  JOIN type dt ON ds.type = dt.id
  JOIN cdn ON cdn.id = ds.cdn_id
  JOIN deliveryservice_server dss on dss.deliveryservice = ds.id
  LEFT JOIN origin o on (o.deliveryservice = ds.id AND o.is_primary)
`
	if strings.HasPrefix(string(serverType), tc.MidTypePrefix) {
		// Note inactive DSes are omitted from Mids, but not Edges
		// See https://github.com/apache/trafficcontrol/issues/3746
		qry += `
WHERE
  cdn.id = (select cdn_id from server where host_name = $1)
  AND ds.active = true
`
	} else {
		qry += `
WHERE
  dss.server = (select id from server where host_name = $1)
`
	}

	qry += `
  AND dt.name = '` + string(tc.DSTypeHTTPNoCache) + `'
`

	rows, err := tx.Query(qry, serverName)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := map[tc.DeliveryServiceName]atscfg.ServerCacheConfigDS{}
	for rows.Next() {
		dsName := tc.DeliveryServiceName("")
		ds := atscfg.ServerCacheConfigDS{}
		if err := rows.Scan(&dsName, &ds.OrgServerFQDN, &ds.Type); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		ds.Type = tc.DSTypeFromString(string(ds.Type))
		dses[dsName] = ds
	}
	return dses, nil
}
