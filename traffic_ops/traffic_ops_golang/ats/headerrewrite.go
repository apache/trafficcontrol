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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/jmoiron/sqlx"
	"math"
	"net/http"
	"regexp"
	"strconv"
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

	where := "WHERE ds.xml_id = '" + inf.Params["xml-id"] + "'"
	query := deliveryservice.GetDSSelectQuery() + where
	dses, errs, _ := deliveryservice.GetDeliveryServices(query, nil, inf.Tx)

	if len(errs) > 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_xml-id.config text: "+err.Error()))
		return
	}

	if len(dses) > 0 {
		ds := dses[0]
		maxOriginConnections := *ds.MaxOriginConnections

		dsType, err := deliveryservice.GetDeliveryServiceType(*ds.ID, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_xml-id.config text: "+err.Error()))
			return
		}
		usesMids := dsType.UsesMidCache()

		// write a header rewrite rule if maxOriginConnections > 0 and the ds does NOT use mids
		if maxOriginConnections > 0 && !usesMids {
			dsEdgeCount, err := getDSEdgeCount(inf.Tx, *ds.ID)
			maxOriginConnectionsPerEdge := int(math.Round(float64(maxOriginConnections) / float64(dsEdgeCount)))
			if err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ds server count: "+err.Error()))
				return
			}
			if ds.EdgeHeaderRewrite == nil {
				text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerEdge) + " [L]"
			} else {
				text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerEdge) + "\n"

			}
		}

		// write the contents of ds.EdgeHeaderRewrite to hdr_rw_xml-id.config replacing any instances of __RETURN__ (surrounded by spaces or not) with \n
		if ds.EdgeHeaderRewrite != nil {
			var re = regexp.MustCompile(`\s*__RETURN__\s*`)
			text += re.ReplaceAllString(*ds.EdgeHeaderRewrite, "\n")
		}
	}

	text += "\n"

	w.Header().Set("Content-Type", "text/plain")
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

	where := "WHERE ds.xml_id = '" + inf.Params["xml-id"] + "'"
	query := deliveryservice.GetDSSelectQuery() + where
	dses, errs, _ := deliveryservice.GetDeliveryServices(query, nil, inf.Tx)

	if len(errs) > 0 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
		return
	}

	if len(dses) > 0 {
		ds := dses[0]
		maxOriginConnections := *ds.MaxOriginConnections

		dsType, err := deliveryservice.GetDeliveryServiceType(*ds.ID, inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting hdr_rw_mid_xml-id.config text: "+err.Error()))
			return
		}
		usesMids := dsType.UsesMidCache()

		// write a header rewrite rule if maxOriginConnections > 0 and the ds DOES use mids
		if maxOriginConnections > 0 && usesMids {
			dsMidCount, err := getDSMidCount(inf.Tx, *ds.CDNID)
			maxOriginConnectionsPerMid := int(math.Round(float64(maxOriginConnections) / float64(dsMidCount)))
			if err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ds server count: "+err.Error()))
				return
			}
			if ds.MidHeaderRewrite == nil {
				text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerMid) + " [L]"
			} else {
				text += "cond %{REMAP_PSEUDO_HOOK}\nset-config proxy.config.http.origin_max_connections " + strconv.Itoa(maxOriginConnectionsPerMid) + "\n"

			}
		}

		// write the contents of ds.MidHeaderRewrite to hdr_rw_mid_xml-id.config replacing any instances of __RETURN__ (surrounded by spaces or not) with \n
		if ds.MidHeaderRewrite != nil {
			var re = regexp.MustCompile(`\s*__RETURN__\s*`)
			text += re.ReplaceAllString(*ds.MidHeaderRewrite, "\n")
		}
	}

	text += "\n"

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}

func getDSEdgeCount(tx *sqlx.Tx, dsID int) (int, error) {
	qry := `SELECT count(1) FROM deliveryservice_server WHERE deliveryservice = $1`
	count := 0
	if err := tx.QueryRow(qry, dsID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func getDSMidCount(tx *sqlx.Tx, cdnID int) (int, error) {
	qry := `SELECT count(1) FROM server WHERE type = (SELECT id FROM type WHERE name = 'MID') AND cdn_id = $1`
	count := 0
	if err := tx.QueryRow(qry, cdnID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
