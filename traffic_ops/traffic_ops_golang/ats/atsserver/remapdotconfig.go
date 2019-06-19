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
	"io"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

func GetServerConfigRemap(w http.ResponseWriter, r *http.Request) {
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

	toToolName, toURL, err := ats.GetToolNameAndURL(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting tool name and url: "+err.Error()))
		return
	}

	atsMajorVersion, err := ats.GetATSMajorVersionFromServerName(inf.Tx.Tx, serverName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ATS major version: "+err.Error()))
		return
	}

	cacheURLConfigParams, err := ats.GetServerProfileParamData(inf.Tx.Tx, serverName, "cacheurl.config")
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting cacheurl.config params: "+err.Error()))
		return
	}

	serverInfo, ok, err := ats.GetServerInfoByHost(inf.Tx.Tx, serverName)
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("server not found"), nil)
		return
	}

	remapDSData, err := ats.GetRemapDSData(inf.Tx.Tx, serverInfo)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("Getting remap ds data: "+err.Error()))
		return
	}

	dsProfilesCacheKeyConfigParams, err := ats.GetProfilesParamData(inf.Tx.Tx, atscfg.DSProfileIDs(remapDSData), "cachekey.config")
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("Getting profiles param data for cachekey: "+err.Error()))
		return
	}

	serverPackageParamData, err := ats.GetServerParamData(inf.Tx.Tx, int(serverInfo.ProfileID), "package", serverInfo.HostName, serverInfo.DomainName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("Getting server package param data: "+err.Error()))
		return
	}

	txt := atscfg.MakeRemapDotConfig(serverName, toToolName, toURL, atsMajorVersion, cacheURLConfigParams, dsProfilesCacheKeyConfigParams, serverPackageParamData, serverInfo, remapDSData)

	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	io.WriteString(w, txt)
}
