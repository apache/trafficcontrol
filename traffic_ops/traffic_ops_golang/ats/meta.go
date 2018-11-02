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
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetConfigMetaData(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	server, ok, err := getServerInfo(inf.Tx.Tx, inf.IntParams["id"])
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting server info: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, nil, errors.New("server not found"))
		return
	}

	tmParams, err := GetTMParams(inf.Tx.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting tm.url parameter: "+err.Error()))
		return
	}
	if tmParams.URL == "" {
		log.Warnln("ats.GetConfigMetadata: global tm.url parameter missing or empty!")
	}

	atsData := tc.ATSConfigMetaData{
		Info: tc.ATSConfigMetaDataInfo{
			ProfileID:         server.ProfileID,
			TOReverseProxyURL: tmParams.ReverseProxyURL,
			TOURL:             tmParams.URL,
			ServerIPv4:        server.IP,
			ServerPort:        server.Port,
			ServerName:        server.HostName,
			CDNID:             server.CDNID,
			CDNName:           string(server.CDN),
			ServerID:          server.ID,
			ProfileName:       server.ProfileName,
		},
	}

	locationParams, err := GetLocationParams(inf.Tx.Tx, server.ProfileID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting location parameters: "+err.Error()))
		return
	}

	if locationParams["remap.config"].Location != "" {
		configLocation := locationParams["remap.config"].Location
		uriSignedDSes, err := GetServerURISignedDSes(inf.Tx.Tx, server.HostName, server.Port)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting server uri-signed dses: "+err.Error()))
			return
		}
		for _, ds := range uriSignedDSes {
			cfgName := "uri_signing_" + string(ds) + ".config"
			// If there's already a parameter for it, don't clobber it. The user may wish to override the location.
			if _, ok := locationParams[cfgName]; !ok {
				p := locationParams[cfgName]
				p.FileNameOnDisk = cfgName
				p.Location = configLocation
			}
		}
	}

	for cfgFile, cfgParams := range locationParams {
		atsCfg := tc.ATSConfigMetaDataConfigFile{
			FileNameOnDisk: cfgParams.FileNameOnDisk,
			Location:       cfgParams.Location,
		}

		scope, err := getServerScope(inf.Tx.Tx, cfgFile, server.Type)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("GetConfigMetaData getting scope: "+err.Error()))
			return
		}
		if cfgParams.URL != "" {
			scope = tc.ATSConfigMetaDataConfigFileScopeCDNs
		}
		atsCfg.Scope = string(scope)

		if cfgParams.URL != "" {
			atsCfg.URL = cfgParams.URL
		} else {
			scopeID := ""
			if scope == tc.ATSConfigMetaDataConfigFileScopeCDNs {
				scopeID = string(server.CDN)
			} else if scope == tc.ATSConfigMetaDataConfigFileScopeProfiles {
				scopeID = server.ProfileName
			} else { // ATSConfigMetaDataConfigFileScopeServers
				scopeID = server.HostName
			}
			atsCfg.APIURI = "/api/1.2/" + string(scope) + "/" + scopeID + "/configfiles/ats/" + cfgFile
		}

		atsData.ConfigFiles = append(atsData.ConfigFiles, atsCfg)
	}

	api.WriteRespRaw(w, r, atsData)
}
