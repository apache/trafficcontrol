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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
)

func WithProfileData(w http.ResponseWriter, r *http.Request, contentType string, makeCfg func(tx *sql.Tx, cfg *config.Config, profile ats.ProfileData, fileName string) (string, error)) {
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
	}
	if !ok {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	text, err := makeCfg(inf.Tx.Tx, inf.Config, profileData, fileName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("making config: "+err.Error()))
		return
	}

	if text == "" {
		// TODO replicates old Perl; verify required.
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("not found"), nil)
		return
	}

	if contentType != "" {
		w.Header().Set(tc.ContentType, contentType)
	}
	w.Write([]byte(text))
}
