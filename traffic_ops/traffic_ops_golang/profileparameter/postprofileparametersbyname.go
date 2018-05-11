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
	"io/ioutil"
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"

	"github.com/lib/pq"
)

func PostProfileParamsByName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		params, _, userErr, sysErr, errCode := api.AllParams(r, nil)
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
		profileName := params["name"]

		profileID, profileExists, err := getProfileIDFromName(profileName, db)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("getting profile '"+profileName+"' ID: "+err.Error()))
			return
		}
		if !profileExists {
			api.HandleErr(w, r, http.StatusBadRequest, errors.New("no profile with that name exists"), nil)
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

// getProfileIDFromName returns the profile's ID, whether a profile with name exists, or any error.
func getProfileIDFromName(name string, db *sql.DB) (int, bool, error) {
	id := 0
	if err := db.QueryRow(`SELECT id from profile where name = $1`, name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, errors.New("querying profile id from name: " + err.Error())
	}
	return id, true, nil
}

// insertParametersForProfile returns the PostResp object, because the ID is needed, and the ID must be associated with the real key (name,value,config_file), so we might as well return the whole object.
func insertParametersForProfile(profileName string, params []tc.ProfileParameterByNamePost, db *sql.DB) ([]tc.ProfileParameterPostRespObj, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("beginning transaction: " + err.Error())
	}
	commitTx := false
	defer FinishTx(tx, &commitTx)
	insertParamsQ := `
INSERT INTO parameter (name, config_file, value, secure)
VALUES (unnest($1::text[]), unnest($2::text[]), unnest($3::text[]), unnest($4::bool[]))
ON CONFLICT(name, config_file, value) DO UPDATE set name=EXCLUDED.name RETURNING id, name, config_file, value, secure;
`
	paramNames := make([]string, len(params))
	paramConfigFiles := make([]string, len(params))
	paramValues := make([]string, len(params))
	paramSecures := make([]bool, len(params))
	for i, param := range params {
		paramNames[i] = param.Name
		paramConfigFiles[i] = param.ConfigFile
		paramValues[i] = param.Value
		if param.Secure != 0 {
			paramSecures[i] = true
		}
	}
	rows, err := tx.Query(insertParamsQ, pq.Array(paramNames), pq.Array(paramConfigFiles), pq.Array(paramValues), pq.Array(paramSecures))
	if err != nil {
		return nil, errors.New("querying post parameters for profile: " + err.Error())
	}
	defer rows.Close()
	ids := make([]int64, 0, len(params))
	insertedObjs := []tc.ProfileParameterPostRespObj{}
	for rows.Next() {
		id := int64(0)
		name := ""
		configFile := ""
		value := ""
		secure := false
		secureNum := 0
		if err := rows.Scan(&id, &name, &configFile, &value, &secure); err != nil {
			return nil, errors.New("scanning new parameter IDs: " + err.Error())
		}
		if secure {
			secureNum = 1
		}
		ids = append(ids, id)
		insertedObjs = append(insertedObjs, tc.ProfileParameterPostRespObj{ID: id, ProfileParameterByNamePost: tc.ProfileParameterByNamePost{Name: name, ConfigFile: configFile, Value: value, Secure: secureNum}})
	}
	insertProfileParamsQ := `
INSERT INTO profile_parameter (profile, parameter)
VALUES ((SELECT id FROM profile WHERE name = $1), unnest($2::int[]))
ON CONFLICT DO NOTHING;
`
	if _, err := tx.Exec(insertProfileParamsQ, profileName, pq.Array(ids)); err != nil {
		return nil, errors.New("inserting profile parameters: " + err.Error())
	}

	commitTx = true
	return insertedObjs, nil
}

// FinishTx commits the transaction if commit is true when it's called, otherwise it rolls back the transaction. This is designed to be called in a defer.
func FinishTx(tx *sql.Tx, commit *bool) {
	if tx == nil {
		return
	}
	if !*commit {
		tx.Rollback()
		return
	}
	tx.Commit()
}
