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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
)

func GetSSLMultiCertDotConfig(w http.ResponseWriter, r *http.Request) {
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

	dses, err := GetSSLMultiCertDSes(inf.Tx.Tx, cdnName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting dses: "+err.Error()))
		return
	}

	txt := atscfg.MakeSSLMultiCertDotConfig(cdnName, toToolName, toURL, dses)
	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(txt))
}

type SSLMultiCertDSInfo struct {
	Protocol    int
	Type        tc.DSType
	RoutingName string
	CDNDomain   string
}

func GetSSLMultiCertDSesInfo(tx *sql.Tx, cdn tc.CDNName) (map[tc.DeliveryServiceName]SSLMultiCertDSInfo, error) {
	dses := map[tc.DeliveryServiceName]SSLMultiCertDSInfo{}
	qry := `
SELECT
  ds.xml_id,
  COALESCE(ds.protocol, 0),
  tp.name as ds_type,
  ds.routing_name,
  cdn.domain_name as cdn_domain
FROM
  deliveryservice ds
  JOIN type tp on tp.id = ds.type
  JOIN cdn cdn on cdn.id = ds.cdn_id
WHERE
  ds.cdn_id = (select id from cdn where name = $1)
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		dsName := tc.DeliveryServiceName("")
		ds := SSLMultiCertDSInfo{}
		if err := rows.Scan(&dsName, &ds.Protocol, &ds.Type, &ds.RoutingName, &ds.CDNDomain); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		dses[dsName] = ds
	}
	return dses, nil
}

func GetSSLMultiCertDSes(tx *sql.Tx, cdn tc.CDNName) (map[tc.DeliveryServiceName]atscfg.SSLMultiCertDS, error) {
	dsInfos, err := GetSSLMultiCertDSesInfo(tx, cdn)
	if err != nil {
		return nil, errors.New("getting dses info: " + err.Error())
	}

	dsNames := []string{}
	for dsName, _ := range dsInfos {
		dsNames = append(dsNames, string(dsName))
	}

	matchLists, err := deliveryservice.GetDeliveryServicesMatchLists(dsNames, tx)
	if err != nil {
		return nil, errors.New("getting matchlists: " + err.Error())
	}

	dses := map[tc.DeliveryServiceName]atscfg.SSLMultiCertDS{}
	for dsName, dsInfo := range dsInfos {
		ds := atscfg.SSLMultiCertDS{Type: dsInfo.Type, Protocol: dsInfo.Protocol}
		matchList, ok := matchLists[string(dsName)]
		if !ok {
			return nil, errors.New("ds '" + string(dsName) + "' returned no matchlist, cannot create example URLs!")
		}
		ds.ExampleURLs = deliveryservice.MakeExampleURLs(&dsInfo.Protocol, dsInfo.Type, dsInfo.RoutingName, matchList, dsInfo.CDNDomain)
		dses[dsName] = ds
	}
	return dses, nil
}
