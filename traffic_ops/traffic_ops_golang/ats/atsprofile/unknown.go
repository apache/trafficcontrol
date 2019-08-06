package atsprofile

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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
)

func GetUnknown(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"profile-name-or-id"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	profileNameOrID := inf.Params["profile-name-or-id"]
	profileID, err := strconv.Atoi(profileNameOrID)
	if err != nil {
		profileName := profileNameOrID
		ok := false
		if profileID, ok, err = ats.GetProfileIDFromName(inf.Tx.Tx, profileName); err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting profile id from name: "+err.Error()))
			return
		} else if !ok {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("Resource not found."), nil)
			return
		}
	}

	fileName := strings.TrimSuffix(inf.Params["file"], ".json")

	profileData, ok, err := ats.GetProfileData(inf.Tx.Tx, profileID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting profile data: "+err.Error()))
		return
	} else if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	txt, userErr, sysErr, errCode := makeUnknown(inf.Tx.Tx, profileData, fileName)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}

	w.Header().Set(tc.ContentType, tc.ContentTypeTextPlain)
	w.Write([]byte(txt))
}

// makeUnknown returns the text of the unknown config, any user error, any system error, and the HTTP code to return if there was an error.
func makeUnknown(tx *sql.Tx, profile ats.ProfileData, fileName string) (string, error, error, int) {
	scopeParams, err := ats.GetParamsByName(tx, "scope")
	if err != nil {
		return "", nil, errors.New("getting scope parameters: " + err.Error()), http.StatusInternalServerError
	}

	inScope := false
	for _, scopeParam := range scopeParams {
		if scopeParam.ConfigFile != fileName {
			continue
		}
		if scopeParam.Value != "profiles" {
			continue
		}
		inScope = true
		break
	}

	if !inScope {
		return "", errors.New("Error - incorrect file scope for route used.  Please use the servers route."), nil, http.StatusBadRequest
	}

	toolName, toURL, err := ats.GetToolNameAndURL(tx)
	if err != nil {
		return "", nil, errors.New("getting tool name and URL: " + err.Error()), http.StatusInternalServerError
	}

	paramData, err := ats.GetProfileParamData(tx, profile.ID, fileName)
	if err != nil {
		return "", nil, errors.New("getting profile param data: " + err.Error()), http.StatusInternalServerError
	}

	return atscfg.MakeUnknownConfig(profile.Name, paramData, toolName, toURL), nil, nil, http.StatusOK
}
