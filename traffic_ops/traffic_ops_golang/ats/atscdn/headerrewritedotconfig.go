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
	"errors"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/jmoiron/sqlx"
)

func GetEdgeHeaderRewriteDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName, userErr, sysErr, errCode := ats.GetCDNNameFromNameOrID(inf.Tx.Tx, inf.Params["cdn-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	dsName := inf.Params["xml-id"]

	toToolName, toURL, err := ats.GetToolNameAndURL(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting tool name and url: "+err.Error()))
		return
	}

	ds, err := getDeliveryService(inf.Tx, dsName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_xml-id.config text: "+err.Error()))
		return
	}

	assignedEdges, err := getEdges(inf.Tx, dsName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_xml-id.config text: "+err.Error()))
		return
	}

	txt := atscfg.MakeHeaderRewriteDotConfig(tc.CDNName(cdnName), toToolName, toURL, ds, assignedEdges)
	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(txt))
}

func GetMidHeaderRewriteDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName, userErr, sysErr, errCode := ats.GetCDNNameFromNameOrID(inf.Tx.Tx, inf.Params["cdn-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	dsName := inf.Params["xml-id"]

	ds, err := getDeliveryService(inf.Tx, dsName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
		return
	}

	assignedMids, err := getMids(inf.Tx, dsName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
		return
	}

	toToolName, toURL, err := ats.GetToolNameAndURL(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting tool name and url: "+err.Error()))
		return
	}

	txt := atscfg.MakeHeaderRewriteMidDotConfig(tc.CDNName(cdnName), toToolName, toURL, ds, assignedMids)

	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(txt))
}

func getDeliveryService(tx *sqlx.Tx, xmlId string) (atscfg.HeaderRewriteDS, error) {
	qry := `
SELECT
  ds.id,
  tp.name as type,
  ds.max_origin_connections,
  COALESCE(ds.edge_header_rewrite, ''),
  COALESCE(ds.mid_header_rewrite, '')
FROM
  deliveryservice ds
  JOIN type tp on tp.id = ds.type
WHERE
  ds.xml_id = $1
`
	ds := atscfg.HeaderRewriteDS{}
	if err := tx.QueryRow(qry, xmlId).Scan(&ds.ID, &ds.Type, &ds.MaxOriginConnections, &ds.EdgeHeaderRewrite, &ds.MidHeaderRewrite); err != nil {
		return atscfg.HeaderRewriteDS{}, errors.New("scanning: " + err.Error())
	}
	return ds, nil
}

func getEdges(tx *sqlx.Tx, dsName string) ([]atscfg.HeaderRewriteServer, error) {
	qry := `
SELECT
  s.host_name,
  s.domain_name,
  s.tcp_port,
  s.status
FROM
  server s
  JOIN deliveryservice_server dss ON dss.server = s.id
  JOIN deliveryservice ds ON ds.id = dss.deliveryservice
WHERE
  ds.xml_id = $1
`
	rows, err := tx.Query(qry, dsName)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	servers := []atscfg.HeaderRewriteServer{}
	for rows.Next() {
		s := atscfg.HeaderRewriteServer{}
		if err := rows.Scan(&s.Status, &s.HostName, &s.DomainName, &s.Port); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		s.Status = tc.CacheStatusFromString(string(s.Status))

		servers = append(servers, s)
	}
	return servers, nil
}

func getMids(tx *sqlx.Tx, dsName string) ([]atscfg.HeaderRewriteServer, error) {
	qry := `
SELECT
  s.host_name,
  s.domain_name,
  s.tcp_port,
  s.status
FROM
  server s
WHERE s.cachegroup IN (
  SELECT
    cg.parent_cachegroup_id
  FROM
    cachegroup cg
    JOIN server s on s.cachegroup = cg.id
    JOIN deliveryservice_server dss on dss.server = s.id
    JOIN deliveryservice ds on ds.id = dss.deliveryservice
  WHERE
    ds.xml_id = $1
)
`
	rows, err := tx.Query(qry, dsName)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	servers := []atscfg.HeaderRewriteServer{}
	for rows.Next() {
		s := atscfg.HeaderRewriteServer{}
		if err := rows.Scan(&s.Status, &s.HostName, &s.DomainName, &s.Port); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		s.Status = tc.CacheStatusFromString(string(s.Status))

		servers = append(servers, s)
	}
	return servers, nil
}
