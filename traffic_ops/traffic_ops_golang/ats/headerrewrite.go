package ats

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
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/jmoiron/sqlx"
)

func GetEdgeHeaderRewriteDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName, userErr, sysErr, errCode := getCDNNameFromNameOrID(inf.Tx.Tx, inf.Params["cdn-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	text, err := headerComment(inf.Tx.Tx, "CDN "+cdnName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_xml-id.config text: "+err.Error()))
		return
	}

	ds, err := getDeliveryService(inf.Tx, inf.Params["xml-id"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
		return
	}

	maxOriginConnections := *ds.MaxOriginConnections

	dsType, err := deliveryservice.GetDeliveryServiceType(*ds.ID, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_xml-id.config text: "+err.Error()))
		return
	}
	usesMids := dsType.UsesMidCache()

	// write a header rewrite rule if maxOriginConnections > 0 and the ds does NOT use mids
	if maxOriginConnections > 0 && !usesMids {
		dsOnlineEdgeCount, err := getOnlineDSEdgeCount(inf.Tx, *ds.ID)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ds server count: "+err.Error()))
			return
		}
		if dsOnlineEdgeCount > 0 {
			maxOriginConnectionsPerEdge := int(math.Round(float64(maxOriginConnections) / float64(dsOnlineEdgeCount)))
			text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerEdge)
			if ds.EdgeHeaderRewrite == nil {
				text += " [L]"
			} else {
				text += "\n"
			}
		}
	}

	// write the contents of ds.EdgeHeaderRewrite to hdr_rw_xml-id.config replacing any instances of __RETURN__ (surrounded by spaces or not) with \n
	if ds.EdgeHeaderRewrite != nil {
		var re = regexp.MustCompile(`\s*__RETURN__\s*`)
		text += re.ReplaceAllString(*ds.EdgeHeaderRewrite, "\n")
	}

	text += "\n"

	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(text))
}

func GetMidHeaderRewriteDotConfig(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"cdn-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	cdnName, userErr, sysErr, errCode := getCDNNameFromNameOrID(inf.Tx.Tx, inf.Params["cdn-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	text, err := headerComment(inf.Tx.Tx, "CDN "+cdnName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
		return
	}

	ds, err := getDeliveryService(inf.Tx, inf.Params["xml-id"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
		return
	}

	maxOriginConnections := *ds.MaxOriginConnections

	dsType, err := deliveryservice.GetDeliveryServiceType(*ds.ID, inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
		return
	}
	usesMids := dsType.UsesMidCache()

	// write a header rewrite rule if maxOriginConnections > 0 and the ds DOES use mids
	if maxOriginConnections > 0 && usesMids {
		dsOnlineMidCount, err := getOnlineDSMidCount(inf.Tx, *ds.ID)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ds server count: "+err.Error()))
			return
		}
		if dsOnlineMidCount > 0 {
			maxOriginConnectionsPerMid := int(math.Round(float64(maxOriginConnections) / float64(dsOnlineMidCount)))
			text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerMid)
			if ds.MidHeaderRewrite == nil {
				text += " [L]"
			} else {
				text += "\n"
			}
		}
	}

	// write the contents of ds.MidHeaderRewrite to hdr_rw_mid_xml-id.config replacing any instances of __RETURN__ (surrounded by spaces or not) with \n
	if ds.MidHeaderRewrite != nil {
		var re = regexp.MustCompile(`\s*__RETURN__\s*`)
		text += re.ReplaceAllString(*ds.MidHeaderRewrite, "\n")
	}

	text += "\n"

	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(text))
}

func getDeliveryService(tx *sqlx.Tx, xmlId string) (tc.DeliveryServiceNullable, error) {
	qry := `SELECT id, cdn_id, max_origin_connections, edge_header_rewrite, mid_header_rewrite FROM deliveryservice WHERE xml_id = $1`
	ds := tc.DeliveryServiceNullable{}
	if err := tx.QueryRow(qry, xmlId).Scan(&ds.ID, &ds.CDNID, &ds.MaxOriginConnections, &ds.EdgeHeaderRewrite, &ds.MidHeaderRewrite); err != nil {
		return tc.DeliveryServiceNullable{}, err
	}
	return ds, nil
}

// getOnlineDSEdgeCount gets the count of online or reported edges assigned to a delivery service
func getOnlineDSEdgeCount(tx *sqlx.Tx, dsID int) (int, error) {
	count := 0
	qry := `SELECT count(1)
	  	FROM deliveryservice_server 
		JOIN server ON deliveryservice_server.server = server.id 
		JOIN status ON server.status = status.id
		WHERE deliveryservice_server.deliveryservice = $1 AND status.name IN ('REPORTED', 'ONLINE')`
	if err := tx.QueryRow(qry, dsID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// getOnlineDSMidCount gets the count of online or reported mids employed by the delivery service
// 1. get the cache groups of the edges assigned to the ds
// 2. get the parent cachegroups for those cachegroups (found in 1)
// 3. get the servers that belong to those cachegroups that are a) mids and b) online/reported
func getOnlineDSMidCount(tx *sqlx.Tx, dsID int) (int, error) {
	count := 0
	qry := `SELECT COUNT(1)
FROM server AS s 
JOIN type AS t ON s.type = t.id
JOIN status AS st ON s.status = st.id
WHERE t.name = 'MID' AND st.name IN ('ONLINE', 'REPORTED') AND s.cachegroup IN (
    SELECT cg.parent_cachegroup_id FROM cachegroup AS cg 
    WHERE cg.id IN (
        SELECT s.cachegroup FROM server AS s 
        WHERE s.id IN (
            SELECT server FROM deliveryservice_server WHERE deliveryservice = $1)))`
	if err := tx.QueryRow(qry, dsID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
