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
	"errors"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

func GetConfigMetaData(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"server-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	serverName, userErr, sysErr, errCode := ats.GetServerNameFromNameOrID(inf.Tx.Tx, inf.Params["server-name-or-id"])
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	server, ok, err := ats.GetServerInfoByHost(inf.Tx.Tx, serverName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting server info: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("server not found"))
		return
	}

	tmParams, err := ats.GetTMParams(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting tm.url parameter: "+err.Error()))
		return
	}

	scopeParams, err := ats.GetScopeParameters(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting scope: "+err.Error()))
		return
	}

	locationParams, err := ats.GetLocationParams(inf.Tx.Tx, int(server.ProfileID))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting location parameters: "+err.Error()))
		return
	}

	uriSignedDSes, err := ats.GetServerURISignedDSes(inf.Tx.Tx, server.HostName, server.Port)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting server uri-signed dses: "+err.Error()))
		return
	}

	txt := atscfg.MakeMetaConfig(serverName, server, tmParams.URL, tmParams.ReverseProxyURL, locationParams, uriSignedDSes, scopeParams)
	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.Write([]byte(txt))
}
