package profileparameter

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
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func PostProfileParamsByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_, intParams, userErr, sysErr, errCode := api.AllParams(r, []string{"id"}, []string{"id"})
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		bts, err := ioutil.ReadAll(r.Body)
		if err != nil {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("body read failed"), nil)
			return
		}
		bts = bytes.TrimLeft(bts, " \n\t\r")
		if len(bts) == 0 {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("no body"), nil)
			return
		}
		profParams := []tc.ProfileParameterByNamePost{}
		if bts[0] == '[' {
			if err := json.Unmarshal(bts, &profParams); err != nil {
				api.HandleErr(w, r, http.StatusBadRequest, errors.New("malformed JSON"), nil)
				return
			}
		} else {
			param := tc.ProfileParameterByNamePost{}
			if err := json.Unmarshal(bts, &param); err != nil {
				api.HandleErr(w, r, http.StatusInternalServerError, errors.New("posting profile parameters by name: "+err.Error()), nil)
				return
			}
			profParams = append(profParams, param)
		}

		profileID := intParams["id"]
		profileName, profileExists, err := getProfileNameFromID(profileID, db)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, fmt.Errorf("getting profile ID %d: "+err.Error(), profileID))
			return
		}
		if !profileExists {
			api.HandleErr(w, r, http.StatusBadRequest, fmt.Errorf("no profile with ID %d exists", profileID), nil)
			return
		}

		insertedObjs, err := insertParametersForProfile(profileName, profParams, db)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("posting profile parameters by name: "+err.Error()))
			return
		}

		// TODO create helper func
		resp := struct {
			Response tc.ProfileParameterPostResp `json:"response"`
			tc.Alerts
		}{tc.ProfileParameterPostResp{Parameters: insertedObjs, ProfileName: profileName, ProfileID: profileID}, tc.CreateAlerts(tc.SuccessLevel, "Assign parameters successfully to profile "+profileName)}
		respBts, err := json.Marshal(resp)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("posting profile parameters by name: "+err.Error()))
			return
		}
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		w.Write(respBts)
	}
}

// getProfileIDFromName returns the profile's name, whether a profile with ID exists, or any error.
func getProfileNameFromID(id int, db *sql.DB) (string, bool, error) {
	name := ""
	if err := db.QueryRow(`SELECT name from profile where id = $1`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying profile name from id: " + err.Error())
	}
	return name, true, nil
}
